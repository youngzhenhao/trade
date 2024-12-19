package services

import (
	"errors"
	"gorm.io/gorm"
	"trade/middleware"
	"trade/utils"
)

type NftPresaleOfflinePurchaseData struct {
	gorm.Model
	NftNo          string `json:"nft_no"`
	Name           string `json:"name"`
	NpubKey        string `json:"npub_key"`
	InvitationCode string `json:"invitation_code"`
	AssetId        string `json:"asset_id"`
	AssetName      string `json:"asset_name"`
}

func GetNftPresaleOfflinePurchaseData(nftNo string, npubKey string, invitationCode string, assetId string) (*[]NftPresaleOfflinePurchaseData, error) {
	var err error
	var where string
	if nftNo != "" {
		where += "nft_no = '" + nftNo + "'"
		if npubKey != "" {
			where += " AND npub_key = '" + npubKey + "'"
		}
		if invitationCode != "" {
			where += " AND invitation_code = '" + invitationCode + "'"
		}
		if assetId != "" {
			where += " AND asset_id = '" + assetId + "'"
		}
	} else {
		if npubKey != "" {
			where += "npub_key = '" + npubKey + "'"
			if invitationCode != "" {
				where += " AND invitation_code = '" + invitationCode + "'"
			}
			if assetId != "" {
				where += " AND asset_id = '" + assetId + "'"
			}
		} else {
			if invitationCode != "" {
				where += "invitation_code = '" + invitationCode + "'"
				if assetId != "" {
					where += " AND asset_id = '" + assetId + "'"
				}
			} else {
				if assetId != "" {
					where += "asset_id = '" + assetId + "'"
				} else {
					where = ""
				}
			}
		}
	}
	if where == "" {
		err = errors.New("no query condition")
		return nil, err
	}
	var nftPresaleOfflinePurchaseDatas []NftPresaleOfflinePurchaseData
	err = middleware.DB.
		Model(&NftPresaleOfflinePurchaseData{}).
		Where(where).
		Find(&nftPresaleOfflinePurchaseDatas).
		Error

	if err != nil {
		return nil, utils.AppendErrorInfo(err, "select NftPresaleOfflinePurchaseData")
	}

	return &nftPresaleOfflinePurchaseDatas, nil
}
