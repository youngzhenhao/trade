package services

import "gorm.io/gorm"

type NftPresaleOfflinePurchaseData struct {
	gorm.Model
	NftNo          string `json:"nft_no"`
	Name           string `json:"name"`
	NpubKey        string `json:"npub_key"`
	InvitationCode string `json:"invitation_code"`
	AssetId        string `json:"asset_id"`
	AssetName      string `json:"asset_name"`
}
