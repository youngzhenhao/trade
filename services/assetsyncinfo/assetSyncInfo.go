package assetsyncinfo

import (
	"bytes"
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/btcsuite/btcd/wire"
	"github.com/lightninglabs/taproot-assets/asset"
	"github.com/lightninglabs/taproot-assets/commitment"
	"github.com/lightninglabs/taproot-assets/proof"
	"github.com/lightninglabs/taproot-assets/taprpc"
	"go.etcd.io/etcd/api/v3/v3rpc/rpctypes"
	"gorm.io/gorm"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"trade/api"
	"trade/config"
	"trade/models"
	"trade/services/servicesrpc"
)

const mainnetUniverse = "universe.lightning.finance:10029"
const regtestUniverse = "132.232.109.84:8443"

type SyncInfoRequest struct {
	Id       string `json:"id"`
	Universe string `json:"universe"`
	GroupKey string `json:"group_key"`
}

var (
	AssetNotFoundErr = errors.New("asset not found")
	AssetRequestErr  = errors.New("asset request error")
	SeverError       = errors.New("server error")
)

// GetAssetSyncInfo returns the asset sync information.
func GetAssetSyncInfo(req *SyncInfoRequest) (*models.AssetSyncInfo, error) {
	//  根据资产id查询数据库
	id := req.Id
	AssetSyncInfo, err := ReadAssetSyncInfoByAssetID(id)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	//TODO: 区分资产和组key

	if AssetSyncInfo.AssetId != "" {
		return AssetSyncInfo, nil
	}
	// 如果资产同步信息不存在，则调用资产查询接口获取资产同步信息并更新数据库
	assetSyncInfo, err := getAssetInfoFromLeaves(req.Id)
	if err != nil && !errors.Is(err, AssetNotFoundErr) {
		return nil, err
	}
	// 更新数据库
	if assetSyncInfo != nil {
		assetSyncInfo.Universe = config.GetConfig().ApiConfig.Tapd.UniverseHost
		err = CreateAssetSyncInfo(assetSyncInfo)
		if err != nil {
			fmt.Println(err)
		}
		return assetSyncInfo, nil
	}
	// 尝试从其他节点同步资产信息
	Universes := []string{}
	switch config.GetConfig().NetWork {
	case "mainnet":
		Universes = append(Universes, mainnetUniverse)
	case "regtest":
		Universes = append(Universes, regtestUniverse)
	default:
	}
	if req.Universe != "" {
		Universes = append(Universes, req.Universe)
	}

	for _, universe := range Universes {
		if isSocketValid(universe) {
			_, err := servicesrpc.SyncAsset(universe, id, false, "issuance")
			if err != nil {
				continue
			}
			assetSyncInfo, err = getAssetInfoFromLeaves(req.Id)
			if assetSyncInfo != nil {
				assetSyncInfo.Universe = config.GetConfig().ApiConfig.Tapd.UniverseHost
				err = CreateAssetSyncInfo(assetSyncInfo)
				if err != nil {
					fmt.Println(err)
				}
				return assetSyncInfo, nil
			}
		}
	}
	return nil, AssetNotFoundErr
}

func getAssetInfoFromLeaves(assetId string) (*models.AssetSyncInfo, error) {
	response, err := servicesrpc.GetAssetLeaves(assetId, false, "issuance")
	if err != nil {
		var etcdError *rpctypes.EtcdError
		if errors.As(err, &etcdError) {
			fmt.Println(err.Error())
			return nil, SeverError
		}
		fmt.Println(err.Error())
		return nil, AssetRequestErr
	}
	if response.Leaves == nil {
		return nil, AssetNotFoundErr
	}

	var assetSyncInfo models.AssetSyncInfo
	assetSyncInfo.AssetId = hex.EncodeToString(response.Leaves[0].Asset.AssetGenesis.AssetId)
	assetSyncInfo.Name = response.Leaves[0].Asset.AssetGenesis.Name
	assetSyncInfo.AssetType = models.AssetType(response.Leaves[0].Asset.AssetGenesis.AssetType)
	assetSyncInfo.Amount = response.Leaves[0].Asset.Amount
	assetSyncInfo.Point = response.Leaves[0].Asset.AssetGenesis.GenesisPoint
	assetSyncInfo.Universe = config.GetConfig().ApiConfig.Tapd.UniverseHost

	//获取元数据
	meta, err := servicesrpc.GetAssetMeta(assetSyncInfo.AssetId, false)
	if err != nil {
		return nil, err
	}
	m := api.Meta{}
	m.GetMetaFromStr(string(meta.GetData()))
	assetSyncInfo.Meta = &m.Description
	//获取组信息
	if response.Leaves[0].Asset.AssetGroup != nil {
		groupKey := hex.EncodeToString(response.Leaves[0].Asset.AssetGroup.RawGroupKey)
		assetSyncInfo.GroupKey = &groupKey
		if m.GroupName != "" {
			assetSyncInfo.GroupName = &m.GroupName
		}
	}
	//获取创建资产
	decode := newDecodeProofOffline()
	decodeProof, err := decode.decodeProof(context.Background(), &taprpc.DecodeProofRequest{
		RawProof: response.Leaves[0].Proof,
	})
	if err != nil {
		return nil, err
	}

	assetSyncInfo.CreateHeight = int64(decodeProof.DecodedProof.Asset.ChainAnchor.BlockHeight)
	info, err := servicesrpc.GetBlockInfo(decodeProof.DecodedProof.Asset.ChainAnchor.AnchorBlockHash)
	if err != nil {
		return nil, err
	}
	msgBlock := &wire.MsgBlock{}
	blockReader := bytes.NewReader(info.RawBlock)
	err = msgBlock.Deserialize(blockReader)
	if err != nil {
		return nil, err
	}
	assetSyncInfo.CreateTime = &msgBlock.Header.Timestamp
	return &assetSyncInfo, nil
}

// isSocketValid checks if the given socket is valid or not.
func isSocketValid(socket string) bool {
	host, port, err := net.SplitHostPort(socket)

	if err != nil {
		return false
	}
	_, err = net.LookupHost(host)
	if err != nil {
		return false
	}
	p, err := strconv.Atoi(port)
	if err != nil {
		return false
	}
	return p > 0 && p < 65536
}

type decodeProofOffline struct {
	//withPrevWitnesses and withMetaReveal need an online node
	withPrevWitnesses bool
	withMetaReveal    bool
}

func newDecodeProofOffline() *decodeProofOffline {
	return &decodeProofOffline{
		withPrevWitnesses: false,
		withMetaReveal:    false,
	}
}

func (d *decodeProofOffline) decodeProof(ctx context.Context,
	req *taprpc.DecodeProofRequest) (*taprpc.DecodeProofResponse, error) {

	if req.WithPrevWitnesses || req.WithMetaReveal {
		return nil, fmt.Errorf("unable to marshal proof: cannot set WithPrevWitnesses" +
			"WithMetaReveal when decoding offline")
	}

	var (
		proofReader = bytes.NewReader(req.RawProof)
		rpcProof    *taprpc.DecodedProof
	)
	switch {
	case proof.IsSingleProof(req.RawProof):
		var p proof.Proof
		err := p.Decode(proofReader)
		if err != nil {
			return nil, fmt.Errorf("unable to decode proof: %w",
				err)
		}

		rpcProof, err = d.marshalProof(
			ctx, &p, d.withMetaReveal, d.withPrevWitnesses,
		)
		if err != nil {
			return nil, fmt.Errorf("unable to marshal proof: %w",
				err)
		}

		rpcProof.NumberOfProofs = 1

	case proof.IsProofFile(req.RawProof):
		if err := proof.CheckMaxFileSize(req.RawProof); err != nil {
			return nil, fmt.Errorf("invalid proof file: %w", err)
		}

		var proofFile proof.File
		if err := proofFile.Decode(proofReader); err != nil {
			return nil, fmt.Errorf("unable to decode proof file: "+
				"%w", err)
		}

		latestProofIndex := uint32(proofFile.NumProofs() - 1)
		if req.ProofAtDepth > latestProofIndex {
			return nil, fmt.Errorf("invalid depth %d is greater "+
				"than latest proof index of %d",
				req.ProofAtDepth, latestProofIndex)
		}

		// Default to latest proof.
		index := latestProofIndex - req.ProofAtDepth
		p, err := proofFile.ProofAt(index)
		if err != nil {
			return nil, err
		}

		rpcProof, err = d.marshalProof(
			ctx, p, req.WithPrevWitnesses,
			req.WithMetaReveal,
		)
		if err != nil {
			return nil, fmt.Errorf("unable to marshal proof: %w",
				err)
		}

		rpcProof.ProofAtDepth = req.ProofAtDepth
		rpcProof.NumberOfProofs = uint32(proofFile.NumProofs())

	default:
		return nil, fmt.Errorf("invalid raw proof, could not " +
			"identify decoding format")
	}

	return &taprpc.DecodeProofResponse{
		DecodedProof: rpcProof,
	}, nil
}

func (d *decodeProofOffline) marshalProof(ctx context.Context, p *proof.Proof,
	withPrevWitnesses, withMetaReveal bool) (*taprpc.DecodedProof, error) {

	var (
		rpcMeta        *taprpc.AssetMeta
		rpcGenesis     = p.GenesisReveal
		rpcGroupKey    = p.GroupKeyReveal
		anchorOutpoint = wire.OutPoint{
			Hash:  p.AnchorTx.TxHash(),
			Index: p.InclusionProof.OutputIndex,
		}
		txMerkleProof  = p.TxMerkleProof
		inclusionProof = p.InclusionProof
		splitRootProof = p.SplitRootProof
	)

	var txMerkleProofBuf bytes.Buffer
	if err := txMerkleProof.Encode(&txMerkleProofBuf); err != nil {
		return nil, fmt.Errorf("unable to encode serialized Bitcoin "+
			"merkle proof: %w", err)
	}

	var inclusionProofBuf bytes.Buffer
	if err := inclusionProof.Encode(&inclusionProofBuf); err != nil {
		return nil, fmt.Errorf("unable to encode inclusion proof: %w",
			err)
	}

	if inclusionProof.CommitmentProof == nil {
		return nil, fmt.Errorf("inclusion proof is missing " +
			"commitment proof")
	}
	tsSibling, tsHash, err := commitment.MaybeEncodeTapscriptPreimage(
		inclusionProof.CommitmentProof.TapSiblingPreimage,
	)
	if err != nil {
		return nil, fmt.Errorf("error encoding tapscript sibling: %w",
			err)
	}

	tapProof, err := inclusionProof.CommitmentProof.DeriveByAssetInclusion(
		&p.Asset,
	)
	if err != nil {
		return nil, fmt.Errorf("error deriving inclusion proof: %w",
			err)
	}
	merkleRoot := tapProof.TapscriptRoot(tsHash)

	var exclusionProofs [][]byte
	for _, exclusionProof := range p.ExclusionProofs {
		var exclusionProofBuf bytes.Buffer
		err := exclusionProof.Encode(&exclusionProofBuf)
		if err != nil {
			return nil, fmt.Errorf("unable to encode exclusion "+
				"proofs: %w", err)
		}
		exclusionProofs = append(
			exclusionProofs, exclusionProofBuf.Bytes(),
		)
	}

	var splitRootProofBuf bytes.Buffer
	if splitRootProof != nil {
		err := splitRootProof.Encode(&splitRootProofBuf)
		if err != nil {
			return nil, fmt.Errorf("unable to encode split root "+
				"proof: %w", err)
		}
	}

	rpcAsset, err := d.marshalChainAsset(ctx, &asset.ChainAsset{
		Asset:                  &p.Asset,
		AnchorTx:               &p.AnchorTx,
		AnchorBlockHash:        p.BlockHeader.BlockHash(),
		AnchorBlockHeight:      p.BlockHeight,
		AnchorOutpoint:         anchorOutpoint,
		AnchorInternalKey:      p.InclusionProof.InternalKey,
		AnchorMerkleRoot:       merkleRoot[:],
		AnchorTapscriptSibling: tsSibling,
	}, withPrevWitnesses)
	if err != nil {
		return nil, err
	}

	if withMetaReveal {
		//metaHash := rpcAsset.AssetGenesis.MetaHash
		//if len(metaHash) == 0 {
		//	return nil, fmt.Errorf("asset does not contain meta " +
		//		"data")
		//}
		//
		//rpcMeta, err = r.FetchAssetMeta(
		//	ctx, &taprpc.FetchAssetMetaRequest{
		//		Asset: &taprpc.FetchAssetMetaRequest_MetaHash{
		//			MetaHash: metaHash,
		//		},
		//	},
		//)
		//if err != nil {
		//	return nil, err
		//}
	}

	decodedAssetID := p.Asset.ID()
	var genesisReveal *taprpc.GenesisReveal
	if rpcGenesis != nil {
		genesisReveal = &taprpc.GenesisReveal{
			GenesisBaseReveal: &taprpc.GenesisInfo{
				GenesisPoint: rpcGenesis.FirstPrevOut.String(),
				Name:         rpcGenesis.Tag,
				MetaHash:     rpcGenesis.MetaHash[:],
				AssetId:      decodedAssetID[:],
				OutputIndex:  rpcGenesis.OutputIndex,
				AssetType:    taprpc.AssetType(p.Asset.Type),
			},
		}
	}

	var GroupKeyReveal taprpc.GroupKeyReveal
	if rpcGroupKey != nil {
		GroupKeyReveal = taprpc.GroupKeyReveal{
			RawGroupKey:   rpcGroupKey.RawKey[:],
			TapscriptRoot: rpcGroupKey.TapscriptRoot,
		}
	}

	return &taprpc.DecodedProof{
		Asset:               rpcAsset,
		MetaReveal:          rpcMeta,
		TxMerkleProof:       txMerkleProofBuf.Bytes(),
		InclusionProof:      inclusionProofBuf.Bytes(),
		ExclusionProofs:     exclusionProofs,
		SplitRootProof:      splitRootProofBuf.Bytes(),
		NumAdditionalInputs: uint32(len(p.AdditionalInputs)),
		ChallengeWitness:    p.ChallengeWitness,
		IsBurn:              p.Asset.IsBurn(),
		GenesisReveal:       genesisReveal,
		GroupKeyReveal:      &GroupKeyReveal,
	}, nil
}

func (d *decodeProofOffline) marshalChainAsset(ctx context.Context, a *asset.ChainAsset,
	withWitness bool) (*taprpc.Asset, error) {

	rpcAsset, err := taprpc.MarshalAsset(
		ctx, a.Asset, a.IsSpent, withWitness, nil,
	)
	if err != nil {
		return nil, err
	}

	var anchorTxBytes []byte
	if a.AnchorTx != nil {
		var anchorTxBuf bytes.Buffer
		err := a.AnchorTx.Serialize(&anchorTxBuf)
		if err != nil {
			return nil, fmt.Errorf("unable to serialize anchor "+
				"tx: %w", err)
		}
		anchorTxBytes = anchorTxBuf.Bytes()
	}

	rpcAsset.ChainAnchor = &taprpc.AnchorInfo{
		AnchorTx:         anchorTxBytes,
		AnchorBlockHash:  a.AnchorBlockHash.String(),
		AnchorOutpoint:   a.AnchorOutpoint.String(),
		InternalKey:      a.AnchorInternalKey.SerializeCompressed(),
		MerkleRoot:       a.AnchorMerkleRoot,
		TapscriptSibling: a.AnchorTapscriptSibling,
		BlockHeight:      a.AnchorBlockHeight,
	}

	if a.AnchorLeaseOwner != [32]byte{} {
		rpcAsset.LeaseOwner = a.AnchorLeaseOwner[:]
		rpcAsset.LeaseExpiry = a.AnchorLeaseExpiry.UTC().Unix()
	}

	return rpcAsset, nil
}

// InsertIssuanceProof 向本地数据库插入资产所有证明
func InsertIssuanceProof(id string) error {
	Id := asset.ID{}
	copy(Id[:], id)
	proofs, err := FetchProofs(Id)
	if err != nil {
		return fmt.Errorf("failed to fetch proofs: %w", err)
	}
	for _, p := range proofs {
		err := servicesrpc.InsertProof(p)
		if err != nil {
			return fmt.Errorf("failed to insert proof: %w", err)
		}
	}
	return nil
}

func FetchProofs(id asset.ID) ([]*proof.AnnotatedProof, error) {
	assetID := hex.EncodeToString(id[:])
	assetPath := filepath.Join(config.GetConfig().ApiConfig.Tapd.Dir, "data", config.GetConfig().NetWork, "proofs", assetID)
	entries, err := os.ReadDir(assetPath)
	if err != nil {
		return nil, fmt.Errorf("unable to read dir %s: %w", assetPath,
			err)
	}
	proofs := make([]*proof.AnnotatedProof, len(entries))
	for idx := range entries {
		// We'll skip any files that don't end with our suffix, this
		// will include directories as well, so we don't need to check
		// for those.
		fileName := entries[idx].Name()
		if !strings.HasSuffix(fileName, proof.TaprootAssetsFileSuffix) {
			continue
		}
		parts := strings.Split(strings.ReplaceAll(
			fileName, proof.TaprootAssetsFileSuffix, "",
		), "-")
		if len(parts) != 3 {
			return nil, fmt.Errorf("malformed proof file name "+
				"'%s', expected two parts, got %d", fileName,
				len(parts))
		}
		fullPath := filepath.Join(assetPath, fileName)
		proofFile, err := os.ReadFile(fullPath)
		if err != nil {
			return nil, fmt.Errorf("unable to read proof: %w", err)
		}
		proofs[idx] = &proof.AnnotatedProof{
			Blob: proofFile,
		}
	}
	return proofs, nil
}
