package localQuery

import (
	"errors"
	"gorm.io/gorm"
	"trade/middleware"
	"trade/models/custodyModels"
)

type ChannelQueryQuest struct {
	AssetId string `json:"assetId"`
}
type ChannelQueryResp struct {
	TotalAmount float64 `json:"totalAmount"`
}

var DbError = errors.New("db error")

func QueryChannelAssetInfo(quest *ChannelQueryQuest) (*ChannelQueryResp, error) {
	db := middleware.DB

	var total float64
	if quest.AssetId == "00" {
		err := db.Model(&custodyModels.AccountBtcBalance{}).Select("SUM(amount) as total").Scan(&total).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, DbError
		}
	} else {
		q := db.Where("asset_id =?", quest.AssetId)
		// 查询总金额
		err := q.Model(&custodyModels.AccountBalance{}).Select("SUM(amount) as total").Scan(&total).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, DbError
		}
	}

	var resp ChannelQueryResp
	resp.TotalAmount = total
	return &resp, nil
}
