package servicesrpc

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/lightninglabs/taproot-assets/fn"
	"github.com/lightninglabs/taproot-assets/proof"
	"github.com/lightninglabs/taproot-assets/taprpc"
	"github.com/lightninglabs/taproot-assets/taprpc/universerpc"
	"strconv"
	"strings"
	"trade/config"
	"trade/utils"
)

func GetAssetLeaves(ID string, isGroup bool, proofType string) (*universerpc.AssetLeafResponse, error) {
	requset := universerpc.ID{}
	var p universerpc.ProofType
	switch proofType {
	case "issuance":
		p = universerpc.ProofType_PROOF_TYPE_ISSUANCE
	case "transfer":
		p = universerpc.ProofType_PROOF_TYPE_TRANSFER
	default:
		return nil, fmt.Errorf("unknown proof type: %s", proofType)
	}
	requset.ProofType = p

	if isGroup {
		groupId := universerpc.ID_GroupKeyStr{
			GroupKeyStr: ID,
		}
		requset.Id = &groupId
	} else {
		assetId := universerpc.ID_AssetIdStr{
			AssetIdStr: ID,
		}
		requset.Id = &assetId
	}

	leaves, err := getAssetLeaves(&requset)
	if err != nil {
		return nil, err
	}
	return leaves, nil

}

func GetAssetMeta(ID string, isHash bool) (*taprpc.AssetMeta, error) {
	var request taprpc.FetchAssetMetaRequest
	if isHash {
		assetHast := taprpc.FetchAssetMetaRequest_MetaHashStr{
			MetaHashStr: ID,
		}
		request.Asset = &assetHast
	} else {
		assetId := taprpc.FetchAssetMetaRequest_AssetIdStr{
			AssetIdStr: ID,
		}
		request.Asset = &assetId
	}
	assetMeta, err := getAssetMeta(&request)
	if err != nil {
		return nil, err
	}
	return assetMeta, nil
}

func SyncAsset(universe string, id string, isGroupKey bool, proofType string) (*universerpc.SyncResponse, error) {
	request := universerpc.SyncRequest{}
	var p universerpc.ProofType
	switch proofType {
	case "issuance":
		p = universerpc.ProofType_PROOF_TYPE_ISSUANCE
	case "transfer":
		p = universerpc.ProofType_PROOF_TYPE_TRANSFER
	default:
		return nil, fmt.Errorf("unknown proof type: %s", proofType)
	}

	if isGroupKey {
		groupKey := universerpc.ID_GroupKeyStr{
			GroupKeyStr: id,
		}
		request.SyncTargets = append(request.SyncTargets, &universerpc.SyncTarget{
			Id: &universerpc.ID{Id: &groupKey,
				ProofType: p},
		})
	} else {
		assetId := universerpc.ID_AssetIdStr{
			AssetIdStr: id,
		}
		request.SyncTargets = append(request.SyncTargets, &universerpc.SyncTarget{
			Id: &universerpc.ID{Id: &assetId,
				ProofType: p},
		})
	}
	request.UniverseHost = universe
	request.SyncMode = universerpc.UniverseSyncMode_SYNC_ISSUANCE_ONLY
	response, err := syncAsset(&request)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func InsertProof(annotatedProof *proof.AnnotatedProof) error {
	// Decode annotated proof into proof file.
	proofFile := &proof.File{}
	err := proofFile.Decode(bytes.NewReader(annotatedProof.Blob))
	if err != nil {
		return err
	}
	// Iterate over each proof in the proof file and submit to the courier
	// service.
	for i := 0; i < proofFile.NumProofs(); i++ {
		transitionProof, err := proofFile.ProofAt(uint32(i))
		if err != nil {
			return err
		}
		proofAsset := transitionProof.Asset

		// Construct asset leaf.
		rpcAsset, err := taprpc.MarshalAsset(
			context.Background(), &proofAsset, true, true, nil, fn.None[uint32](),
		)
		if err != nil {
			return err
		}

		var proofBuf bytes.Buffer
		if err := transitionProof.Encode(&proofBuf); err != nil {
			return fmt.Errorf("error encoding proof file: %w", err)
		}

		assetLeaf := universerpc.AssetLeaf{
			Asset: rpcAsset,
			Proof: proofBuf.Bytes(),
		}

		// Construct universe key.
		outPoint := transitionProof.OutPoint()
		assetKey := universerpc.MarshalAssetKey(
			outPoint, proofAsset.ScriptKey.PubKey,
		)
		assetID := proofAsset.ID()

		var (
			groupPubKey      *btcec.PublicKey
			groupPubKeyBytes []byte
		)
		if proofAsset.GroupKey != nil {
			groupPubKey = &proofAsset.GroupKey.GroupPubKey
			groupPubKeyBytes = groupPubKey.SerializeCompressed()
		}

		universeID := universerpc.MarshalUniverseID(
			assetID[:], groupPubKeyBytes,
		)
		universeKey := universerpc.UniverseKey{
			Id:      universeID,
			LeafKey: assetKey,
		}
		// Submit proof to courier.
		err = insertProof(&universerpc.AssetProof{
			Key:       &universeKey,
			AssetLeaf: &assetLeaf,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func getAssetLeaves(request *universerpc.ID) (*universerpc.AssetLeafResponse, error) {
	tapdconf := config.GetConfig().ApiConfig.Tapd
	grpcHost := tapdconf.Host + ":" + strconv.Itoa(tapdconf.Port)
	tlsCertPath := tapdconf.TlsCertPath
	macaroonPath := tapdconf.MacaroonPath
	conn, connClose := utils.GetConn(grpcHost, tlsCertPath, macaroonPath)
	defer connClose()

	client := universerpc.NewUniverseClient(conn)
	response, err := client.AssetLeaves(context.Background(), request)
	return response, err
}

func getAssetMeta(request *taprpc.FetchAssetMetaRequest) (*taprpc.AssetMeta, error) {
	tapdconf := config.GetConfig().ApiConfig.Tapd

	grpcHost := tapdconf.Host + ":" + strconv.Itoa(tapdconf.Port)
	tlsCertPath := tapdconf.TlsCertPath
	macaroonPath := tapdconf.MacaroonPath

	conn, connClose := utils.GetConn(grpcHost, tlsCertPath, macaroonPath)
	defer connClose()

	client := taprpc.NewTaprootAssetsClient(conn)
	response, err := client.FetchAssetMeta(context.Background(), request)
	return response, err
}

func syncAsset(request *universerpc.SyncRequest) (*universerpc.SyncResponse, error) {
	tapdconf := config.GetConfig().ApiConfig.Tapd

	grpcHost := tapdconf.Host + ":" + strconv.Itoa(tapdconf.Port)
	tlsCertPath := tapdconf.TlsCertPath
	macaroonPath := tapdconf.MacaroonPath

	conn, connClose := utils.GetConn(grpcHost, tlsCertPath, macaroonPath)
	defer connClose()

	client := universerpc.NewUniverseClient(conn)
	response, err := client.SyncUniverse(context.Background(), request)
	return response, err
}

func insertProof(request *universerpc.AssetProof) error {
	tapdconf := config.GetConfig().ApiConfig.Tapd

	grpcHost := tapdconf.Host + ":" + strconv.Itoa(tapdconf.Port)
	tlsCertPath := tapdconf.TlsCertPath
	macaroonPath := tapdconf.MacaroonPath

	conn, connClose := utils.GetConn(grpcHost, tlsCertPath, macaroonPath)
	defer connClose()

	client := universerpc.NewUniverseClient(conn)
	_, err := client.InsertProof(context.Background(), request)
	if err != nil {
		return err
	}
	return nil
}

func NewAddr(assetId string, amt int, proofCourierAddr string) (*taprpc.Addr, error) {
	tapdconf := config.GetConfig().ApiConfig.Tapd
	grpcHost := tapdconf.Host + ":" + strconv.Itoa(tapdconf.Port)
	tlsCertPath := tapdconf.TlsCertPath
	macaroonPath := tapdconf.MacaroonPath
	conn, connClose := utils.GetConn(grpcHost, tlsCertPath, macaroonPath)
	defer connClose()

	client := taprpc.NewTaprootAssetsClient(conn)
	_assetIdByteSlice, _ := hex.DecodeString(assetId)
	if !strings.HasPrefix(proofCourierAddr, "universerpc://") {
		proofCourierAddr = "universerpc://" + proofCourierAddr
	}
	request := &taprpc.NewAddrRequest{
		AssetId:          _assetIdByteSlice,
		Amt:              uint64(amt),
		ProofCourierAddr: proofCourierAddr,
	}
	response, err := client.NewAddr(context.Background(), request)
	if err != nil {
		return nil, err
	}
	return response, nil
}
