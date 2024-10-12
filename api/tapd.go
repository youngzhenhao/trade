package api

import (
	"encoding/hex"
	"errors"
	"github.com/lightninglabs/taproot-assets/proof"
	"github.com/lightninglabs/taproot-assets/taprpc"
	"github.com/lightninglabs/taproot-assets/taprpc/mintrpc"
	"github.com/lightninglabs/taproot-assets/taprpc/universerpc"
	"strconv"
	"trade/config"
	"trade/models"
	"trade/utils"
)

func GetAssetInfo(id string) (*models.AssetIssuanceLeaf, error) {
	return assetLeafIssuanceInfo(id)
}

func MintAsset(name string, assetTypeIsCollectible bool, assetMetaData *Meta, amount int, newGroupedAsset bool) string {
	Metastr := assetMetaData.ToJsonStr()
	response, err := mintAsset(false, assetTypeIsCollectible, name, Metastr, false, amount, newGroupedAsset, false, "", "", false)
	if err != nil {
		return utils.MakeJsonResult(false, "mintAsset error. "+err.Error(), "")
	}
	return utils.MakeJsonResult(true, "", response)
}

func FinalizeBatch(feeRate int) string {
	response, err := finalizeBatch(false, feeRate)
	if err != nil {
		return utils.MakeJsonResult(false, err.Error(), nil)
	}
	return utils.MakeJsonResult(true, "", response)
}

func AddGroupAsset(name string, assetTypeIsCollectible bool, assetMetaData *Meta, amount int, groupKey string) string {
	Metastr := assetMetaData.ToJsonStr()
	response, err := mintAsset(false, assetTypeIsCollectible, name, Metastr, false, amount, false, true, groupKey, "", false)
	if err != nil {
		return utils.MakeJsonResult(false, "mintAsset error. "+err.Error(), "")
	}
	return utils.MakeJsonResult(true, "", response)
}

func NewAddr(assetId string, amt int) string {
	proofCourierAddr := config.GetLoadConfig().ApiConfig.Tapd.UniverseHost
	response, err := newAddr(assetId, amt, proofCourierAddr)
	if err != nil {
		return utils.MakeJsonResult(false, err.Error(), "")
	}
	return utils.MakeJsonResult(true, "", response)
}

func NewAddrAndGetResponse(assetId string, amt int) (*taprpc.Addr, error) {
	proofCourierAddr := config.GetLoadConfig().ApiConfig.Tapd.UniverseHost
	return newAddr(assetId, amt, proofCourierAddr)
}

func NewAddrAndGetStringResponse(assetId string, amt int) (string, error) {
	proofCourierAddr := config.GetLoadConfig().ApiConfig.Tapd.UniverseHost
	response, err := newAddr(assetId, amt, proofCourierAddr)
	if err != nil {
		return "", err
	}
	return response.Encoded, nil
}

func SendAsset(addr string, feeRate int) string {
	response, err := sendAsset(addr, feeRate)
	if err != nil {
		return utils.MakeJsonResult(false, err.Error(), "")
	}
	return utils.MakeJsonResult(true, "", response)
}

func SendAssetBool(addr string, feeRate int) (bool, error) {
	_, err := sendAsset(addr, feeRate)
	if err != nil {
		return false, utils.AppendErrorInfo(err, "sendAsset")
	}
	return true, nil
}

func SendAssetAndGetResponse(addr string, feeRate int) (*taprpc.SendAssetResponse, error) {
	return sendAsset(addr, feeRate)
}

func SendAssetAddrSliceAndGetResponse(addrSlice []string, feeRate int) (*taprpc.SendAssetResponse, error) {
	return sendAssetAddrSlice(addrSlice, feeRate)
}

func DecodeAddr(addr string) string {
	response, err := decodeAddr(addr)
	if err != nil {
		return utils.MakeJsonResult(false, err.Error(), "")
	}
	return utils.MakeJsonResult(true, "", response)
}

func GetDecodedAddrInfo(addr string) (*taprpc.Addr, error) {
	return decodeAddr(addr)
}

func MintAssetAndGetResponse(name string, assetTypeIsCollectible bool, assetMetaData *Meta, amount int, newGroupedAsset bool) (*mintrpc.MintAssetResponse, error) {
	return mintAsset(false, assetTypeIsCollectible, name, assetMetaData.ToJsonStr(), false, amount, newGroupedAsset, false, "", "", false)
}

func FinalizeBatchAndGetResponse(feeRate int) (*mintrpc.FinalizeBatchResponse, error) {
	return finalizeBatch(false, feeRate)
}

func GetListAssetsResponse(withWitness bool, includeSpent bool, includeLeased bool) (*taprpc.ListAssetResponse, error) {
	return listAssets(withWitness, includeSpent, includeLeased)
}

func TransactionAndIndexToOutpoint(transaction string, index int) (outpoint string) {
	return transaction + ":" + strconv.Itoa(index)
}

func BatchTxidAnchorToAssetId(batchTxidAnchor string) (string, error) {
	assets, _ := listAssets(true, true, false)
	for _, asset := range assets.Assets {
		txid, _ := utils.OutpointToTransactionAndIndex(asset.GetChainAnchor().GetAnchorOutpoint())
		if batchTxidAnchor == txid {
			return hex.EncodeToString(asset.GetAssetGenesis().AssetId), nil
		}
	}
	err := errors.New("no asset found for batch txid")
	return "", err
}

func QueryAssetType(assetType int) (string, error) {
	if assetType == 0 {
		return taprpc.AssetType_NORMAL.String(), nil
	} else if assetType == 1 {
		return taprpc.AssetType_COLLECTIBLE.String(), nil
	}
	return "", errors.New("not a valid asset type code")
}

func ListBalancesAndGetResponse(isGroupByAssetIdOrGroupKey bool) (*taprpc.ListBalancesResponse, error) {
	return listBalances(isGroupByAssetIdOrGroupKey)
}

func ListTransfersAndGetResponse() (*taprpc.ListTransfersResponse, error) {
	return listTransfers()
}

func GetAllOutPointsOfListTransfersResponse(listTransfersResponse *taprpc.ListTransfersResponse) []string {
	var allOutPoints []string
	for _, listTransfer := range listTransfersResponse.Transfers {
		for _, input := range listTransfer.Inputs {
			allOutPoints = append(allOutPoints, input.AnchorPoint)
		}
		for _, output := range listTransfer.Outputs {
			allOutPoints = append(allOutPoints, output.Anchor.Outpoint)
		}
	}
	return allOutPoints
}

func ProcessListTransfersResponse(network models.Network, listTransfersResponse *taprpc.ListTransfersResponse, deviceId string) *[]models.AssetTransferProcessedSetRequest {
	var assetTransferProcessed []models.AssetTransferProcessedSetRequest
	allOutpoints := GetAllOutPointsOfListTransfersResponse(listTransfersResponse)
	response, err := GetAddressesByOutpointSlice(network, allOutpoints)
	if err != nil {
		return nil
	}
	addressMap := response
	for _, listTransfer := range listTransfersResponse.Transfers {
		var txid string
		txid, err = utils.GetTxidFromOutpoint(listTransfer.Outputs[0].Anchor.Outpoint)
		if err != nil {
			return nil
		}
		var assetTransferProcessedInput []models.AssetTransferProcessedInput
		for _, input := range listTransfer.Inputs {
			inOp := input.AnchorPoint
			assetTransferProcessedInput = append(assetTransferProcessedInput, models.AssetTransferProcessedInput{
				Address:     addressMap[inOp],
				Amount:      int(input.Amount),
				AnchorPoint: inOp,
				ScriptKey:   hex.EncodeToString(input.ScriptKey),
			})
		}
		var assetTransferProcessedOutput []models.AssetTransferProcessedOutput
		for _, output := range listTransfer.Outputs {
			outOp := output.Anchor.Outpoint
			assetTransferProcessedOutput = append(assetTransferProcessedOutput, models.AssetTransferProcessedOutput{
				Address:                addressMap[outOp],
				Amount:                 int(output.Amount),
				AnchorOutpoint:         outOp,
				AnchorValue:            int(output.Anchor.Value),
				AnchorInternalKey:      hex.EncodeToString(output.Anchor.InternalKey),
				AnchorTaprootAssetRoot: hex.EncodeToString(output.Anchor.TaprootAssetRoot),
				AnchorMerkleRoot:       hex.EncodeToString(output.Anchor.MerkleRoot),
				AnchorTapscriptSibling: hex.EncodeToString(output.Anchor.TapscriptSibling),
				AnchorNumPassiveAssets: int(output.Anchor.NumPassiveAssets),
				ScriptKey:              hex.EncodeToString(output.ScriptKey),
				ScriptKeyIsLocal:       output.ScriptKeyIsLocal,
				NewProofBlob:           hex.EncodeToString(output.NewProofBlob),
				SplitCommitRootHash:    hex.EncodeToString(output.SplitCommitRootHash),
				OutputType:             output.OutputType.String(),
				AssetVersion:           output.AssetVersion.String(),
			})
		}
		assetTransferProcessed = append(assetTransferProcessed, models.AssetTransferProcessedSetRequest{
			Txid:               txid,
			AssetID:            hex.EncodeToString(listTransfer.Inputs[0].AssetId),
			TransferTimestamp:  int(listTransfer.TransferTimestamp),
			AnchorTxHash:       hex.EncodeToString(listTransfer.AnchorTxHash),
			AnchorTxHeightHint: int(listTransfer.AnchorTxHeightHint),
			AnchorTxChainFees:  int(listTransfer.AnchorTxChainFees),
			Inputs:             assetTransferProcessedInput,
			Outputs:            assetTransferProcessedOutput,
			DeviceID:           deviceId,
		})
	}
	return &assetTransferProcessed
}

func ListTransfersAndGetProcessedResponse(network models.Network, deviceId string) (*[]models.AssetTransferProcessedSetRequest, error) {
	transfers, err := ListTransfersAndGetResponse()
	if err != nil {
		return nil, err
	}
	processedListTransfers := ProcessListTransfersResponse(network, transfers, deviceId)
	return processedListTransfers, nil
}

type ListBalancesShortResponse struct {
	Name    string `json:"name"`
	AssetId string `json:"assetId"`
	Balance int    `json:"balance"`
}

func ListBalancesAndGetShortResponse() (*[]ListBalancesShortResponse, error) {
	var listBalancesShortResponses []ListBalancesShortResponse
	response, err := listBalances(true)
	if err != nil {
		return nil, err
	}
	for _, balance := range (*response).AssetBalances {
		listBalancesShortResponses = append(listBalancesShortResponses, ListBalancesShortResponse{
			Name:    balance.AssetGenesis.Name,
			AssetId: hex.EncodeToString(balance.AssetGenesis.AssetId),
			Balance: int(balance.Balance),
		})
	}
	return &listBalancesShortResponses, nil
}

func SyncAssetIssuanceAndGetResponse(universeHost string, assetId string) (*universerpc.SyncResponse, error) {
	//universeHost := "mainnet.universe.lightning.finance:10029"
	if universeHost == "" {
		return nil, errors.New("universe host is empty")
	}
	_proofType := universerpc.ProofType_PROOF_TYPE_ISSUANCE
	var targets []*universerpc.SyncTarget
	universeID := &universerpc.ID{
		Id: &universerpc.ID_AssetIdStr{
			AssetIdStr: assetId,
		},
		ProofType: _proofType,
	}
	targets = append(targets, &universerpc.SyncTarget{
		Id: universeID,
	})
	response, err := syncUniverse(universeHost, targets, universerpc.UniverseSyncMode_SYNC_FULL)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func SyncAssetIssuance(assetId string) error {
	var universeHost string
	switch config.GetLoadConfig().NetWork {
	case "mainnet":
		universeHost = "mainnet.universe.lightning.finance:10029"
	case "testnet":
	case "regtest":
	default:
	}
	_, err := SyncAssetIssuanceAndGetResponse(universeHost, assetId)
	return err
}

func AddrReceivesAndGetResponse() (*taprpc.AddrReceivesResponse, error) {
	return addrReceives()
}

func AddrReceivesResponseToAddrReceiveEventSetRequests(addrReceivesResponse *taprpc.AddrReceivesResponse, deviceId string) *[]models.AddrReceiveEventSetRequest {
	var addrReceiveEvents []models.AddrReceiveEventSetRequest
	for _, event := range addrReceivesResponse.Events {
		addrReceiveEvents = append(addrReceiveEvents, models.AddrReceiveEventSetRequest{
			CreationTimeUnixSeconds: int(event.CreationTimeUnixSeconds),
			Addr: models.AddrReceiveEventSetRequestAddr{
				Encoded:          event.Addr.Encoded,
				AssetID:          hex.EncodeToString(event.Addr.AssetId),
				Amount:           int(event.Addr.Amount),
				ScriptKey:        hex.EncodeToString(event.Addr.ScriptKey),
				InternalKey:      hex.EncodeToString(event.Addr.InternalKey),
				TaprootOutputKey: hex.EncodeToString(event.Addr.TaprootOutputKey),
				ProofCourierAddr: event.Addr.ProofCourierAddr,
			},
			Status:             event.Status.String(),
			Outpoint:           event.Outpoint,
			UtxoAmtSat:         int(event.UtxoAmtSat),
			ConfirmationHeight: int(event.ConfirmationHeight),
			HasProof:           event.HasProof,
			DeviceID:           deviceId,
		})
	}
	return &addrReceiveEvents
}

func AddrReceivesAndGetEventSetRequests(deviceId string) (*[]models.AddrReceiveEventSetRequest, error) {
	response, err := AddrReceivesAndGetResponse()
	if err != nil {
		return nil, err
	}
	return AddrReceivesResponseToAddrReceiveEventSetRequests(response, deviceId), nil
}

func GetAssetNameByAssetId(assetId string) string {
	assetInfo, err := GetAssetInfo(assetId)
	if err != nil || assetInfo == nil {
		return ""
	}
	return assetInfo.Name
}

func ListUtxosAndGetResponse() (*taprpc.ListUtxosResponse, error) {
	return listUtxos()
}

func FetchAssetMetaAndGetResponse(assetId string) (*taprpc.AssetMeta, error) {
	return fetchAssetMetaByAssetId(assetId)
}

type AssetInfoApi struct {
	AssetId      string  `json:"asset_Id"`
	Name         string  `json:"name"`
	Point        string  `json:"point"`
	AssetType    string  `json:"assetType"`
	GroupName    *string `json:"group_name"`
	GroupKey     *string `json:"group_key"`
	Amount       uint64  `json:"amount"`
	Meta         *string `json:"meta"`
	CreateHeight int64   `json:"create_height"`
	CreateTime   int64   `json:"create_time"`
	Universe     string  `json:"universe"`
}

func QueryAssetRoots(assetId string) *universerpc.QueryRootResponse {
	return queryAssetRoots(assetId)
}

func GetAssetInfoApi(id string) (*AssetInfoApi, error) {
	root := QueryAssetRoots(id)
	if root == nil || root.IssuanceRoot.Id == nil {
		return nil, errors.New("query asset roots err")
	}
	queryId := id
	isGroup := false
	if groupKey, ok := root.IssuanceRoot.Id.Id.(*universerpc.ID_GroupKey); ok {
		isGroup = true
		queryId = hex.EncodeToString(groupKey.GroupKey)
	}
	response, err := assetLeaves(isGroup, queryId, universerpc.ProofType_PROOF_TYPE_ISSUANCE)
	if err != nil {
		return nil, err
	}
	if response.Leaves == nil {
		return nil, errors.New("response leaves null err")
	}
	var blob proof.Blob
	for index, leaf := range response.Leaves {
		if hex.EncodeToString(leaf.Asset.AssetGenesis.GetAssetId()) == id {
			blob = response.Leaves[index].Proof
			break
		}
	}
	if len(blob) == 0 {
		return nil, errors.New("blob length zero err")
	}
	p, _ := blob.AsSingleProof()
	assetId := p.Asset.ID().String()
	assetName := p.Asset.Tag
	assetPoint := p.Asset.FirstPrevOut.String()
	assetType := p.Asset.Type.String()
	amount := p.Asset.Amount
	createHeight := p.BlockHeight
	createTime := p.BlockHeader.Timestamp
	var (
		newMeta Meta
		m       = ""
	)
	if p.MetaReveal != nil {
		m = string(p.MetaReveal.Data)
	}
	newMeta.GetMetaFromStr(m)
	var assetInfo = AssetInfoApi{
		AssetId:      assetId,
		Name:         assetName,
		Point:        assetPoint,
		AssetType:    assetType,
		GroupName:    &newMeta.GroupName,
		Amount:       amount,
		Meta:         &newMeta.Description,
		CreateHeight: int64(createHeight),
		CreateTime:   createTime.Unix(),
		Universe:     "localhost",
	}
	if isGroup {
		assetInfo.GroupKey = &queryId
	}
	return &assetInfo, nil
}

func GetGroupKeyByAssetId(assetId string) (string, error) {
	response, err := GetListAssetsResponse(false, true, false)
	if err != nil {
		return "", err
	}
	for _, asset := range response.Assets {
		if assetId == hex.EncodeToString(asset.AssetGenesis.AssetId) {
			return hex.EncodeToString(asset.AssetGroup.TweakedGroupKey), nil
		}
	}
	err = errors.New("asset group key not found")
	return "", err
}

type AssetKeys struct {
	OpStr          string `json:"op_str"`
	ScriptKeyBytes string `json:"script_key_bytes"`
}

func AssetLeafKeysAndGetResponse(isGroup bool, id string, proofType universerpc.ProofType) (*universerpc.AssetLeafKeyResponse, error) {
	return assetLeafKeys(isGroup, id, proofType)
}

func AssetLeafKeyResponseToAssetKeys(response *universerpc.AssetLeafKeyResponse) *[]AssetKeys {
	if response == nil {
		return nil
	}
	var assetKeys []AssetKeys
	for _, key := range response.AssetKeys {
		assetKeys = append(assetKeys, AssetKeys{
			OpStr:          key.Outpoint.(*universerpc.AssetKey_OpStr).OpStr,
			ScriptKeyBytes: hex.EncodeToString(key.GetScriptKeyBytes()),
		})
	}
	return &assetKeys
}

func AssetLeafKeys(isGroup bool, id string, proofType universerpc.ProofType) (*[]AssetKeys, error) {
	response, err := assetLeafKeys(isGroup, id, proofType)
	if err != nil {
		return nil, err
	}
	var assetKeys *[]AssetKeys
	assetKeys = AssetLeafKeyResponseToAssetKeys(response)
	return assetKeys, nil
}

type AssetMeta struct {
	Data     string `json:"data"`
	Type     string `json:"type"`
	MetaHash string `json:"meta_hash"`
}

// FetchAssetMetaByAssetId
// @Description: Fetch asset meta by asset id
func FetchAssetMetaByAssetId(assetId string) (*AssetMeta, error) {
	response, err := FetchAssetMetaAndGetResponse(assetId)
	if err != nil {
		return nil, err
	}
	assetMeta := AssetMeta{
		Data:     string(response.Data),
		Type:     response.Type.String(),
		MetaHash: hex.EncodeToString(response.MetaHash),
	}
	return &assetMeta, nil
}

func QueryProofAndGetResponse(isGroup bool, id string, outpoint string, scriptKey string, proofType universerpc.ProofType) (*universerpc.AssetProofResponse, error) {
	return queryProof(isGroup, id, outpoint, scriptKey, proofType)
}

// @dev: Has not been used now
type QueryProofAndResponse struct {
	Req struct {
		ID struct {
			GroupKey  string `json:"group_key"`
			ProofType string `json:"proof_type"`
		} `json:"id"`
		LeafKey struct {
			Op struct {
				HashStr string `json:"hash_str"`
				Index   int    `json:"index"`
			} `json:"op"`
			ScriptKeyStr string `json:"script_key_str"`
		} `json:"leaf_key"`
	} `json:"req"`
	UniverseRoot struct {
		ID struct {
			GroupKey  string `json:"group_key"`
			ProofType string `json:"proof_type"`
		} `json:"id"`
		MssmtRoot struct {
			RootHash string `json:"root_hash"`
			RootSum  string `json:"root_sum"`
		} `json:"mssmt_root"`
		AssetName        string `json:"asset_name"`
		AmountsByAssetID struct {
		} `json:"amounts_by_asset_id"`
	} `json:"universe_root"`
	UniverseInclusionProof string `json:"universe_inclusion_proof"`
	AssetLeaf              struct {
		Asset struct {
			Version      string `json:"version"`
			AssetGenesis struct {
				GenesisPoint string `json:"genesis_point"`
				Name         string `json:"name"`
				MetaHash     string `json:"meta_hash"`
				AssetID      string `json:"asset_id"`
				AssetType    string `json:"asset_type"`
				OutputIndex  int    `json:"output_index"`
			} `json:"asset_genesis"`
			Amount           string `json:"amount"`
			LockTime         int    `json:"lock_time"`
			RelativeLockTime int    `json:"relative_lock_time"`
			ScriptVersion    int    `json:"script_version"`
			ScriptKey        string `json:"script_key"`
			ScriptKeyIsLocal bool   `json:"script_key_is_local"`
			AssetGroup       struct {
				RawGroupKey     string `json:"raw_group_key"`
				TweakedGroupKey string `json:"tweaked_group_key"`
				AssetWitness    string `json:"asset_witness"`
				TapscriptRoot   string `json:"tapscript_root"`
			} `json:"asset_group"`
			ChainAnchor   interface{} `json:"chain_anchor"`
			PrevWitnesses []struct {
				PrevID struct {
					AnchorPoint string `json:"anchor_point"`
					AssetID     string `json:"asset_id"`
					ScriptKey   string `json:"script_key"`
					Amount      string `json:"amount"`
				} `json:"prev_id"`
				TxWitness       []string    `json:"tx_witness"`
				SplitCommitment interface{} `json:"split_commitment"`
			} `json:"prev_witnesses"`
			IsSpent                bool        `json:"is_spent"`
			LeaseOwner             string      `json:"lease_owner"`
			LeaseExpiry            string      `json:"lease_expiry"`
			IsBurn                 bool        `json:"is_burn"`
			ScriptKeyDeclaredKnown bool        `json:"script_key_declared_known"`
			ScriptKeyHasScriptPath bool        `json:"script_key_has_script_path"`
			DecimalDisplay         interface{} `json:"decimal_display"`
		} `json:"asset"`
		Proof string `json:"proof"`
	} `json:"asset_leaf"`
	MultiverseRoot struct {
		RootHash string `json:"root_hash"`
		RootSum  string `json:"root_sum"`
	} `json:"multiverse_root"`
	MultiverseInclusionProof string `json:"multiverse_inclusion_proof"`
}

// QueryProofToGetAssetId
// @Description: Query proof to get asset id
func QueryProofToGetAssetId(groupKey string, outpoint string, scriptKey string) (string, error) {
	response, err := QueryProofAndGetResponse(true, groupKey, outpoint, scriptKey, universerpc.ProofType_PROOF_TYPE_ISSUANCE)
	if err != nil {
		return "", err
	}
	assetId := hex.EncodeToString(response.AssetLeaf.Asset.AssetGenesis.AssetId)
	return assetId, nil
}
