package services

import (
	"trade/middleware"
	"trade/models"
)

type NftPresaleInfo struct {
	ID              uint                   `json:"id"`
	AssetId         string                 `json:"asset_id"`
	Meta            string                 `json:"meta"`
	GroupKey        string                 `json:"group_key"`
	Price           int                    `json:"price"`
	Info            string                 `json:"info"`
	BuyerUsername   string                 `json:"buyer_username"`
	BoughtTime      int                    `json:"bought_time"`
	PaidId          int                    `json:"paid_id"`
	PaidSuccessTime int                    `json:"paid_success_time"`
	State           models.NftPresaleState `json:"state"`
}

func GetPurchasedNftPresaleInfo() ([]NftPresaleInfo, error) {
	db := middleware.DB
	var nftPresaleInfos []NftPresaleInfo
	err := db.Table("nft_presales").
		Select("id, asset_id, meta, group_key, price, info, buyer_username, bought_time, paid_id, paid_success_time, state").
		Where("state > ?", models.NftPresaleStatePaidPending).
		Scan(&nftPresaleInfos).
		Error
	if err != nil {
		return nil, err
	}
	return nftPresaleInfos, nil
}
