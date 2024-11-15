package mint_nft

import (
	"fmt"
	"trade/api"
	"trade/btlLog"
	"trade/services"
	"trade/utils"
)

func MintFirst(groupName string, description string, imgPath string, feeRate uint) error {
	if feeRate > 50 {
		return fmt.Errorf("feeRate(%d) is too high!\n", feeRate)
	}
	id := 1
	meta := api.NewMetaWithGroupName(description, groupName)
	name := fmt.Sprintf("%s#%03d", groupName, id)
	_, err := meta.LoadImage(imgPath)
	if err != nil {
		return fmt.Errorf("\nMint %s LoadImage\n%v", name, err)
	}
	mintResponse, err := api.MintNftAssetFirst(name, meta)
	if err != nil {
		return fmt.Errorf("\nMint %s MintNftAssetFirst\n%v", name, err)
	}
	btlLog.PreSale.Info("\nMint %s MintNftAssetFirst\n%v", name, utils.ValueJsonString(mintResponse))
	// Auto fee rate
	feeRateSatPerKw := services.FeeRateSatPerBToSatPerKw(int(feeRate))
	finalizeResponse, err := api.FinalizeBatchAndGetResponse(feeRateSatPerKw)
	if err != nil {
		return fmt.Errorf("\nMint %s FinalizeBatchAndGetResponse\n%v", name, err)
	}
	btlLog.PreSale.Info("\nMint %s FinalizeBatchAndGetResponse\n%v", name, utils.ValueJsonString(finalizeResponse))
	batchTxidAnchor := finalizeResponse.GetBatch().GetBatchTxid()
	assetId, err := api.BatchTxidAnchorToAssetId(batchTxidAnchor)
	if err != nil {
		return fmt.Errorf("\nMint %s BatchTxidAnchorToAssetId\n%v", name, err)
	}
	btlLog.PreSale.Info("\nMint %s BatchTxidAnchorToAssetId\n%v", name, assetId)
	groupKey, err := api.BatchTxidAnchorToGroupKey(batchTxidAnchor)
	if err != nil {
		return fmt.Errorf("\nMint %s BatchTxidAnchorToGroupKey\n%v", name, err)
	}
	btlLog.PreSale.Info("\nMint %s BatchTxidAnchorToGroupKey\n%v", name, groupKey)
	btlLog.PreSale.Info("asset id: %s\n", assetId)
	btlLog.PreSale.Info("group key: %s\n", groupKey)
	btlLog.PreSale.Info("batch txid: %s\n", batchTxidAnchor)
	return nil
}
