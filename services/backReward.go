package services

import (
	"trade/middleware"
	"trade/models"
	"trade/utils"
)

type RecordType int

const (
	_ RecordType = iota
	RecordTypeFairLaunch
	RecordTypeNftPresale
)

type backReward struct {
	User       string `json:"user"`
	Fee        uint64 `json:"fee"`
	RecordType string `json:"record_type"`
	RecordId   uint   `json:"record_id"`
}

func (r RecordType) String() string {
	backRewardMap := map[RecordType]string{
		RecordTypeFairLaunch: "FairLaunch",
		RecordTypeNftPresale: "NftPresale",
	}
	return backRewardMap[r]
}

type fairLaunchMintedInfoRecord struct {
	Id           uint   `json:"id"`
	Username     string `json:"username"`
	MintedGasFee int    `json:"minted_gas_fee"`
}

func fairLaunchMintedInfoRecordsToBackRewards(fairLaunchMintedInfoRecords *[]fairLaunchMintedInfoRecord) *[]backReward {
	if fairLaunchMintedInfoRecords == nil {
		return new([]backReward)
	}
	var backRewards []backReward
	for _, record := range *fairLaunchMintedInfoRecords {
		backRewards = append(backRewards, backReward{
			User:       record.Username,
			Fee:        uint64(record.MintedGasFee),
			RecordType: RecordTypeFairLaunch.String(),
			RecordId:   record.Id,
		})
	}
	return &backRewards
}

type nftPresaleRecord struct {
	Id            uint   `json:"id"`
	BuyerUsername string `json:"buyer_username"`
	Price         int    `json:"price"`
}

func nftPresaleRecordsToBackRewards(nftPresaleRecords *[]nftPresaleRecord) *[]backReward {
	if nftPresaleRecords == nil {
		return new([]backReward)
	}
	var backRewards []backReward
	for _, record := range *nftPresaleRecords {
		backRewards = append(backRewards, backReward{
			User:       record.BuyerUsername,
			Fee:        uint64(record.Price),
			RecordType: RecordTypeNftPresale.String(),
			RecordId:   record.Id,
		})
	}
	return &backRewards
}

func GetBackRewards(username string) (*[]backReward, error) {
	backRewards := new([]backReward)
	var fairLaunchMintedInfoRecords []fairLaunchMintedInfoRecord
	var err error
	err = middleware.DB.
		Model(&models.FairLaunchMintedInfo{}).
		Select("id, username, minted_gas_fee").
		Where("username = ? and state > ?", username, models.FairLaunchMintedStatePaidPending).
		Order("id desc").
		Scan(&fairLaunchMintedInfoRecords).Error
	if err != nil {
		return backRewards, utils.AppendErrorInfo(err, "Scan FairLaunchMintedInfo")
	}
	var nftPresaleRecords []nftPresaleRecord
	err = middleware.DB.
		Model(&models.NftPresale{}).
		Select("id, buyer_username, price").
		Where("buyer_username = ? and state > ?", username, models.NftPresaleStatePaidPending).
		Order("id desc").
		Scan(&nftPresaleRecords).Error
	if err != nil {
		return backRewards, utils.AppendErrorInfo(err, "Scan NftPresale")
	}
	fRecords := fairLaunchMintedInfoRecordsToBackRewards(&fairLaunchMintedInfoRecords)
	*backRewards = append(*backRewards, *fRecords...)
	nRecords := nftPresaleRecordsToBackRewards(&nftPresaleRecords)
	*backRewards = append(*backRewards, *nRecords...)
	return backRewards, nil
}
