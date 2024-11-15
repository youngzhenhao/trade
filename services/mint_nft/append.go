package mint_nft

import (
	"fmt"
	"sort"
	"strconv"
	"trade/api"
	"trade/btlLog"
	"trade/services"
	"trade/utils"
)

func MintAppend(groupKey string, groupName string, description string, imgPathPrefix string, imgSuffix string, feeRate uint, start int, end int) error {
	if groupKey == "" {
		return fmt.Errorf("groupKey(%d) is required!\n", groupKey)
	}
	if feeRate > 50 {
		return fmt.Errorf("feeRate(%d) is too high!\n", feeRate)
	}
	if start < 1 || end < 1 || end < start {
		return fmt.Errorf("in valid start or end (%d,%d)\n", start, end)
	}
	for ; start <= end; start++ {
		name := fmt.Sprintf("%s#%03d", groupName, start)
		meta := api.NewMetaWithGroupName(description, groupName)
		// TODO: imgPath need to modify
		imgPath := imgPathPrefix + strconv.Itoa(start) + imgSuffix
		_, err := meta.LoadImage(imgPath)
		if err != nil {
			return fmt.Errorf("\nMint %s LoadImage\n%v", name, err)
		}
		mintResponse, err := api.MintNftAssetAppend(name, meta, groupKey)
		if err != nil {
			return fmt.Errorf("\nMint %s MintNftAssetFirst\n%v", name, err)
		}
		btlLog.PreSale.Info("\nMint %s MintNftAssetFirst\n%v", name, utils.ValueJsonString(mintResponse))
	}
	// Auto fee rate
	feeRateSatPerKw := services.FeeRateSatPerBToSatPerKw(int(feeRate))
	finalizeResponse, err := api.FinalizeBatchAndGetResponse(feeRateSatPerKw)
	if err != nil {
		return fmt.Errorf("\nMint %s FinalizeBatchAndGetResponse\n%v", fmt.Sprintf("%d-%d", start, end), err)
	}
	btlLog.PreSale.Info("\nMint %s FinalizeBatchAndGetResponse\n%v", fmt.Sprintf("%d-%d", start, end), utils.ValueJsonString(finalizeResponse))
	batchTxidAnchor := finalizeResponse.GetBatch().GetBatchTxid()
	assetIdAndNames, err := api.BatchTxidAnchorToAssetIdAndNames(batchTxidAnchor)
	if err != nil {
		return fmt.Errorf("\nMint %s BatchTxidAnchorToAssetIdAndNames\n%v", fmt.Sprintf("%d-%d", start, end), err)
	}
	btlLog.PreSale.Info("\nMint %s BatchTxidAnchorToAssetIdAndNames\n%v", fmt.Sprintf("%d-%d", start, end), utils.ValueJsonString(assetIdAndNames))
	sort.Slice(*assetIdAndNames, func(i, j int) bool {
		return (*assetIdAndNames)[i].Name < (*assetIdAndNames)[j].Name
	})
	btlLog.PreSale.Info("asset id:\n%s\n", utils.ValueJsonString(assetIdAndNames))
	return nil
}
