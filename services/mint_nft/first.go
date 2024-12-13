package mint_nft

import (
	"fmt"
	"trade/api"
	"trade/btlLog"
	"trade/services"
	"trade/utils"
)

func MintFirst(groupName string, description string, attributesPath string, imgPath string, feeRate uint) error {
	if feeRate > 50 {
		return fmt.Errorf("feeRate(%d) is too high!\n", feeRate)
	}
	id := 0
	attributes, err := GetAttributesFromFile(attributesPath)
	if err != nil {
		return fmt.Errorf("GetAttributesFromFile %s \n%v", attributesPath, err)
	}

	meta := api.NewMetaWithAttributes(description, groupName, attributes)
	name := fmt.Sprintf("%s#%d", groupName, id)

	_, err = meta.LoadImage(imgPath)
	if err != nil {
		return fmt.Errorf("\nMint %s LoadImage\n%v", name, err)
	}
	mintResponse, err := api.MintNftAssetFirst(name, meta)
	if err != nil {
		return fmt.Errorf("\nMint %s MintNftAssetFirst\n%v", name, err)
	}
	btlLog.MintNft.Info("\nMint %s MintNftAssetFirst\n%v", name, utils.ValueJsonString(mintResponse))
	// Auto fee rate
	feeRateSatPerKw := services.FeeRateSatPerBToSatPerKw(int(feeRate))
	finalizeResponse, err := api.FinalizeBatchAndGetResponse(feeRateSatPerKw)
	if err != nil {
		_, _err := api.CancelBatchAndGetResponse()
		if _err != nil {
			fmt.Printf("%v\n", fmt.Errorf("\nMint %s CancelBatchAndGetResponse\n%v", name, err))
		}

		return fmt.Errorf("\nMint %s FinalizeBatchAndGetResponse\n%v", name, err)

	}
	btlLog.MintNft.Info("\nMint %s FinalizeBatchAndGetResponse\n%v", name, utils.ValueJsonString(finalizeResponse))
	batchTxidAnchor := finalizeResponse.GetBatch().GetBatchTxid()
	assetId, err := api.BatchTxidAnchorToAssetId(batchTxidAnchor)
	if err != nil {
		return fmt.Errorf("\nMint %s BatchTxidAnchorToAssetId\n%v", name, err)
	}
	btlLog.MintNft.Info("\nMint %s BatchTxidAnchorToAssetId\n%v", name, assetId)
	groupKey, err := api.BatchTxidAnchorToGroupKey(batchTxidAnchor)
	if err != nil {
		return fmt.Errorf("\nMint %s BatchTxidAnchorToGroupKey\n%v", name, err)
	}
	btlLog.MintNft.Info("\nMint %s BatchTxidAnchorToGroupKey\n%v", name, groupKey)
	btlLog.MintNft.Info("asset id: %s\n", assetId)
	btlLog.MintNft.Info("group key: %s\n", groupKey)
	btlLog.MintNft.Info("batch txid: %s\n", batchTxidAnchor)
	return nil
}
