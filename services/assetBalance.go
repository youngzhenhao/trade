package services

import (
	"errors"
	"trade/models"
)

func GetAssetBalancesByUserId(userId int) (*[]models.AssetBalance, error) {
	return ReadAssetBalancesByUserId(userId)
}

func ProcessAssetBalanceSetRequest(userId int, assetBalanceSetRequest *models.AssetBalanceSetRequest) *models.AssetBalance {
	var assetBalance models.AssetBalance
	assetBalance = models.AssetBalance{
		GenesisPoint: assetBalanceSetRequest.GenesisPoint,
		Name:         assetBalanceSetRequest.Name,
		MetaHash:     assetBalanceSetRequest.MetaHash,
		AssetID:      assetBalanceSetRequest.AssetID,
		AssetType:    assetBalanceSetRequest.AssetType,
		OutputIndex:  assetBalanceSetRequest.OutputIndex,
		Version:      assetBalanceSetRequest.Version,
		Balance:      assetBalanceSetRequest.Balance,
		DeviceId:     assetBalanceSetRequest.DeviceId,
		UserId:       userId,
	}
	return &assetBalance
}

func IsAssetBalanceChanged(assetBalanceByInvoice *models.AssetBalance, old *models.AssetBalance) bool {
	if assetBalanceByInvoice == nil || old == nil {
		return true
	}
	if assetBalanceByInvoice.GenesisPoint != old.GenesisPoint {
		return true
	}
	if assetBalanceByInvoice.Name != old.Name {
		return true
	}
	if assetBalanceByInvoice.MetaHash != old.MetaHash {
		return true
	}
	if assetBalanceByInvoice.AssetID != old.AssetID {
		return true
	}
	if assetBalanceByInvoice.AssetType != old.AssetType {
		return true
	}
	if assetBalanceByInvoice.OutputIndex != old.OutputIndex {
		return true
	}
	if assetBalanceByInvoice.Version != old.Version {
		return true
	}
	if assetBalanceByInvoice.Balance != old.Balance {
		return true
	}
	if assetBalanceByInvoice.DeviceId != old.DeviceId {
		return true
	}
	if assetBalanceByInvoice.UserId != old.UserId {
		return true
	}
	return false
}

func CheckAssetBalanceIfUpdate(assetBalance *models.AssetBalance) (*models.AssetBalance, error) {
	if assetBalance == nil {
		return nil, errors.New("nil asset balance")
	}
	assetBalanceByAssetId, err := ReadAssetBalanceByAssetID(assetBalance.AssetID)
	if err != nil {
		return assetBalance, nil
	}
	if !IsAssetBalanceChanged(assetBalanceByAssetId, assetBalance) {
		return assetBalanceByAssetId, nil
	}
	assetBalanceByAssetId.GenesisPoint = assetBalance.GenesisPoint
	assetBalanceByAssetId.Name = assetBalance.Name
	assetBalanceByAssetId.MetaHash = assetBalance.MetaHash
	assetBalanceByAssetId.AssetID = assetBalance.AssetID
	assetBalanceByAssetId.AssetType = assetBalance.AssetType
	assetBalanceByAssetId.OutputIndex = assetBalance.OutputIndex
	assetBalanceByAssetId.Version = assetBalance.Version
	assetBalanceByAssetId.Balance = assetBalance.Balance
	assetBalanceByAssetId.DeviceId = assetBalance.DeviceId
	assetBalanceByAssetId.UserId = assetBalance.UserId
	return assetBalanceByAssetId, nil
}

func CreateOrUpdateAssetBalance(lock *models.AssetBalance) (err error) {
	var assetBalance *models.AssetBalance
	assetBalance, err = CheckAssetBalanceIfUpdate(lock)
	return UpdateAssetBalance(assetBalance)
}

func ProcessAssetBalanceSetRequestSlice(userId int, assetBalanceSetRequestSlice *[]models.AssetBalanceSetRequest) *[]models.AssetBalance {
	var assetBalances []models.AssetBalance
	for _, assetBalanceRequest := range *assetBalanceSetRequestSlice {
		assetBalances = append(assetBalances, models.AssetBalance{
			GenesisPoint: assetBalanceRequest.GenesisPoint,
			Name:         assetBalanceRequest.Name,
			MetaHash:     assetBalanceRequest.MetaHash,
			AssetID:      assetBalanceRequest.AssetID,
			AssetType:    assetBalanceRequest.AssetType,
			OutputIndex:  assetBalanceRequest.OutputIndex,
			Version:      assetBalanceRequest.Version,
			Balance:      assetBalanceRequest.Balance,
			DeviceId:     assetBalanceRequest.DeviceId,
			UserId:       userId,
		})
	}
	return &assetBalances
}

func CreateOrUpdateAssetBalances(balances *[]models.AssetBalance) (err error) {
	var assetBalances []models.AssetBalance
	var assetBalance *models.AssetBalance
	for _, balance := range *balances {
		assetBalance, err = CheckAssetBalanceIfUpdate(&balance)
		if err != nil {
			return err
		}
		assetBalances = append(assetBalances, *assetBalance)
	}
	return UpdateAssetBalances(&assetBalances)
}

type UserAssetBalance struct {
	UserId        int                    `json:"user_id"`
	AssetBalances *[]models.AssetBalance `json:"asset_balances"`
}

func GetAllAssetBalances() (*[]models.AssetBalance, error) {
	return ReadAllAssetBalances()
}

func AssetBalancesToUserMapAssetBalances(assetBalances *[]models.AssetBalance) *map[int]*[]models.AssetBalance {
	userMapAssetBalances := make(map[int]*[]models.AssetBalance)
	for _, assetBalance := range *assetBalances {
		balances, ok := userMapAssetBalances[assetBalance.UserId]
		if !ok {
			userMapAssetBalances[assetBalance.UserId] = &[]models.AssetBalance{assetBalance}
		} else {
			*balances = append(*balances, assetBalance)
		}
	}
	return &userMapAssetBalances
}

func UserMapAssetBalancesToUserAssetBalances(userMapAssetBalances *map[int]*[]models.AssetBalance) *[]UserAssetBalance {
	var userAssetBalances []UserAssetBalance
	for userId, assetBalances := range *userMapAssetBalances {
		userAssetBalances = append(userAssetBalances, UserAssetBalance{
			UserId:        userId,
			AssetBalances: assetBalances,
		})
	}
	return &userAssetBalances
}

func AssetBalancesToUserAssetBalances(assetBalances *[]models.AssetBalance) *[]UserAssetBalance {
	userMapAssetBalances := AssetBalancesToUserMapAssetBalances(assetBalances)
	userAssetBalances := UserMapAssetBalancesToUserAssetBalances(userMapAssetBalances)
	return userAssetBalances
}

// GetAllUserAssetBalances
// @Description: Get all asset balances by userId
func GetAllUserAssetBalances() (*[]UserAssetBalance, error) {
	allAssetBalances, err := GetAllAssetBalances()
	if err != nil {
		return nil, err
	}
	userAssetBalances := AssetBalancesToUserAssetBalances(allAssetBalances)
	return userAssetBalances, nil
}

type AssetIdAndBalance struct {
	AssetId       string                 `json:"asset_id"`
	AssetBalances *[]models.AssetBalance `json:"asset_balances"`
}

func AssetBalancesToAssetIdMapAssetBalances(assetBalances *[]models.AssetBalance) *map[string]*[]models.AssetBalance {
	AssetIdMapAssetBalances := make(map[string]*[]models.AssetBalance)
	for _, assetBalance := range *assetBalances {
		balances, ok := AssetIdMapAssetBalances[assetBalance.AssetID]
		if !ok {
			AssetIdMapAssetBalances[assetBalance.AssetID] = &[]models.AssetBalance{assetBalance}
		} else {
			*balances = append(*balances, assetBalance)
		}
	}
	return &AssetIdMapAssetBalances
}

func AssetIdMapAssetBalancesToAssetIdAndBalances(AssetIdMapAssetBalances *map[string]*[]models.AssetBalance) *[]AssetIdAndBalance {
	var assetIdAndBalances []AssetIdAndBalance
	for assetId, assetBalances := range *AssetIdMapAssetBalances {
		assetIdAndBalances = append(assetIdAndBalances, AssetIdAndBalance{
			AssetId:       assetId,
			AssetBalances: assetBalances,
		})
	}
	return &assetIdAndBalances
}

func AssetBalancesToAssetIdAndBalances(assetBalances *[]models.AssetBalance) *[]AssetIdAndBalance {
	assetIdMapAssetBalances := AssetBalancesToAssetIdMapAssetBalances(assetBalances)
	assetIdAndBalances := AssetIdMapAssetBalancesToAssetIdAndBalances(assetIdMapAssetBalances)
	return assetIdAndBalances
}

// GetAllAssetIdAndBalances
// @Description: Get all asset balances by assetId
func GetAllAssetIdAndBalances() (*[]AssetIdAndBalance, error) {
	allAssetBalances, err := GetAllAssetBalances()
	if err != nil {
		return nil, err
	}
	assetIdAndBalances := AssetBalancesToAssetIdAndBalances(allAssetBalances)
	return assetIdAndBalances, nil
}

type AssetIdAndUserAssetBalance struct {
	AssetId          string              `json:"asset_id"`
	UserAssetBalance *[]UserAssetBalance `json:"user_asset_balance"`
}

// GetAllAssetIdAndUserAssetBalances
// @Description: Get all asset balances by assetId and userId
func GetAllAssetIdAndUserAssetBalances() (*[]AssetIdAndUserAssetBalance, error) {
	var assetIdAndUserAssetBalances []AssetIdAndUserAssetBalance
	allAssetBalances, err := GetAllAssetBalances()
	if err != nil {
		return nil, err
	}
	assetIdAndBalances := AssetBalancesToAssetIdAndBalances(allAssetBalances)
	for _, assetIdAndBalance := range *assetIdAndBalances {
		userAssetBalances := AssetBalancesToUserAssetBalances(assetIdAndBalance.AssetBalances)

		assetIdAndUserAssetBalances = append(assetIdAndUserAssetBalances, AssetIdAndUserAssetBalance{
			AssetId:          assetIdAndBalance.AssetId,
			UserAssetBalance: userAssetBalances,
		})
	}
	return &assetIdAndUserAssetBalances, nil
}
