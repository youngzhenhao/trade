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

func MintAppend(groupKey string, groupName string, description string, imgPathPrefix string, imgSuffix string, attributesPathPrefix string, attributesSuffix string, feeRate uint, start int, end int) error {
	if groupKey == "" {
		return fmt.Errorf("groupKey(%s) is required!\n", groupKey)
	}
	if feeRate > 50 {
		return fmt.Errorf("feeRate(%d) is too high!\n", feeRate)
	}
	if start < 1 || end < 1 || end < start {
		return fmt.Errorf("in valid start or end (%d,%d)\n", start, end)
	}
	for i := start; i <= end; i++ {
		// including start and end
		name := fmt.Sprintf("%s#%d", groupName, i)

		// TODO: Path need to modify
		attributesPath := attributesPathPrefix + groupName + "#" + strconv.Itoa(i) + attributesSuffix
		imgPath := imgPathPrefix + groupName + "#" + strconv.Itoa(i%10000) + imgSuffix

		attributes, err := GetAttributesFromFile(attributesPath)
		if err != nil {
			return fmt.Errorf("GetAttributesFromFile %s \n%v", attributesPath, err)
		}
		meta := api.NewMetaWithAttributes(description, groupName, attributes)

		_, err = meta.LoadImage(imgPath)
		if err != nil {
			return fmt.Errorf("\nMint %s LoadImage\n%v", name, err)
		}
		mintResponse, err := api.MintNftAssetAppend(name, meta, groupKey)
		if err != nil {
			return fmt.Errorf("\nMint %s MintNftAssetAppend\n%v", name, err)
		}
		btlLog.MintNft.Info("\nMint %s MintNftAssetAppend\n%v", name, utils.ValueJsonString(mintResponse))
	}
	// Auto fee rate
	feeRateSatPerKw := services.FeeRateSatPerBToSatPerKw(int(feeRate))
	finalizeResponse, err := api.FinalizeBatchAndGetResponse(feeRateSatPerKw)
	if err != nil {
		_, _err := api.CancelBatchAndGetResponse()
		if _err != nil {
			fmt.Printf("%v\n", fmt.Errorf("\nMint %s CancelBatchAndGetResponse\n%v", fmt.Sprintf("%d-%d", start, end), err))
		}

		return fmt.Errorf("\nMint %s FinalizeBatchAndGetResponse\n%v", fmt.Sprintf("%d-%d", start, end), err)
	}
	btlLog.MintNft.Info("\nMint %s FinalizeBatchAndGetResponse\n%v", fmt.Sprintf("%d-%d", start, end), utils.ValueJsonString(finalizeResponse))
	batchTxidAnchor := finalizeResponse.GetBatch().GetBatchTxid()
	assetIdAndNames, err := api.BatchTxidAnchorToAssetIdAndNames(batchTxidAnchor)
	if err != nil {
		return fmt.Errorf("\nMint %s BatchTxidAnchorToAssetIdAndNames\n%v", fmt.Sprintf("%d-%d", start, end), err)
	}
	btlLog.MintNft.Info("\nMint %s BatchTxidAnchorToAssetIdAndNames\n%v", fmt.Sprintf("%d-%d", start, end), utils.ValueJsonString(assetIdAndNames))
	sort.Slice(*assetIdAndNames, func(i, j int) bool {
		return (*assetIdAndNames)[i].Name < (*assetIdAndNames)[j].Name
	})
	btlLog.MintNft.Info("asset id:\n%s\n", utils.ValueJsonString(assetIdAndNames))
	btlLog.MintNft.Info("batch txid: %s\n", batchTxidAnchor)
	return nil
}
