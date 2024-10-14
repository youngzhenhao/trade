package services

import (
	"trade/api"
	"trade/btlLog"
	"trade/models"
	"trade/services/btldb"
	"trade/utils"
)

func CreateAssetMeta(assetMeta *models.AssetMeta) error {
	return btldb.CreateAssetMeta(assetMeta)
}

func CreateAssetMetas(assetMetas *[]models.AssetMeta) error {
	return btldb.CreateAssetMetas(assetMetas)
}

func ReadAssetMeta(id uint) (*models.AssetMeta, error) {
	return btldb.ReadAssetMeta(id)
}

func ReadAssetMetaByAssetId(assetId string) (*models.AssetMeta, error) {
	return btldb.ReadAssetMetaByAssetId(assetId)
}

func ReadAllAssetMetas() (*[]models.AssetMeta, error) {
	return btldb.ReadAllAssetMetas()
}

func UpdateAssetMeta(assetMeta *models.AssetMeta) error {
	return btldb.UpdateAssetMeta(assetMeta)
}

func UpdateAssetMetas(assetMetas *[]models.AssetMeta) error {
	return btldb.UpdateAssetMetas(assetMetas)
}

func DeleteAssetMeta(id uint) error {
	return btldb.DeleteAssetMeta(id)
}

// @dev:

func GetAssetMetaByAssetId(assetId string) (*models.AssetMeta, error) {
	return ReadAssetMetaByAssetId(assetId)
}

func SetAssetMeta(assetMeta *models.AssetMeta) error {
	return CreateAssetMeta(assetMeta)
}

func SetAssetMetas(assetMetas *[]models.AssetMeta) error {
	return CreateAssetMetas(assetMetas)
}

func GetImageDataFromAssetMeta(assetMeta models.AssetMeta) string {
	var meta api.Meta
	meta.GetMetaFromStr(assetMeta.AssetMeta)
	return meta.ImageData
}

// GetAssetMetaImageDataByAssetId
// @Description: Get assetMeta image data by asset id
func GetAssetMetaImageDataByAssetId(assetId string) (string, error) {
	assetMeta, err := GetAssetMetaByAssetId(assetId)
	if err != nil {
		return "", utils.AppendErrorInfo(err, "GetAssetMetaByAssetId")
	}
	imageData := GetImageDataFromAssetMeta(*assetMeta)
	return imageData, nil
}

func StoreAssetMetaIfNotExist(assetId string) error {
	// @dev: Store asset meta
	_, err := GetAssetMetaByAssetId(assetId)
	if err != nil {
		// @dev: not found asset meta from db, find by api
		var assetMetaStr string
		assetMeta, err := api.FetchAssetMetaByAssetId(assetId)
		if err != nil {
			btlLog.PreSale.Error("api FetchAssetMetaByAssetId err:%v", err)
		} else {
			assetMetaStr = assetMeta.Data
		}
		return CreateAssetMeta(&models.AssetMeta{
			AssetID:   assetId,
			AssetMeta: assetMetaStr,
		})
	}
	return nil
}

func StoreAssetMetasIfNotExist(assetIds []string) error {
	// @dev: Store asset meta
	var assetMetas []models.AssetMeta
	for _, assetId := range assetIds {
		_, err := GetAssetMetaByAssetId(assetId)
		if err != nil {
			// @dev: not found asset meta from db, find by api
			var assetMetaStr string
			fetchAssetMeta, err := api.FetchAssetMetaByAssetId(assetId)
			if err != nil {
				btlLog.PreSale.Error("api FetchAssetMetaByAssetId err:%v", err)
			} else {
				assetMetaStr = fetchAssetMeta.Data
			}
			assetMeta := models.AssetMeta{
				AssetID:   assetId,
				AssetMeta: assetMetaStr,
			}
			assetMetas = append(assetMetas, assetMeta)
		}
	}
	if len(assetMetas) == 0 {
		return nil
	}
	return CreateAssetMetas(&assetMetas)
}
