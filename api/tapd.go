package api

import (
	"encoding/hex"
	"errors"
	"github.com/lightninglabs/taproot-assets/taprpc"
	"github.com/lightninglabs/taproot-assets/taprpc/mintrpc"
	"github.com/lightninglabs/taproot-assets/taprpc/universerpc"
	"strconv"
	"strings"
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
	response, err := newAddr(assetId, amt)
	if err != nil {
		return utils.MakeJsonResult(false, err.Error(), "")
	}
	return utils.MakeJsonResult(true, "", response)
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
	return transaction + strconv.Itoa(index)
}

func OutpointToTransactionAndIndex(outpoint string) (transaction string, index string) {
	result := strings.Split(outpoint, ":")
	return result[0], result[1]
}

func BatchTxidAnchorToAssetId(batchTxidAnchor string) (string, error) {
	assets, _ := listAssets(true, true, false)
	for _, asset := range assets.Assets {
		txid, _ := OutpointToTransactionAndIndex(asset.GetChainAnchor().GetAnchorOutpoint())
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
	universeHost := config.GetLoadConfig().ApiConfig.Tapd.UniverseHost
	_, err := SyncAssetIssuanceAndGetResponse(universeHost, assetId)
	return err
}
