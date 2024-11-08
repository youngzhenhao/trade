package services

import (
	"errors"
	"trade/models"
	"trade/services/btldb"
	"trade/utils"
)

func GetAssetListsByUserId(userId int) (*[]models.AssetList, error) {
	return btldb.ReadAssetListsByUserId(userId)
}

func GetAssetListsByUserIdNonZero(userId int) (*[]models.AssetList, error) {
	return btldb.ReadAssetListsByUserIdNonZero(userId)
}

func ProcessAssetListSetRequest(userId int, username string, assetListSetRequest *models.AssetListSetRequest) *models.AssetList {
	var assetList models.AssetList
	assetList = models.AssetList{
		Version:          assetListSetRequest.Version,
		GenesisPoint:     assetListSetRequest.GenesisPoint,
		Name:             assetListSetRequest.Name,
		MetaHash:         assetListSetRequest.MetaHash,
		AssetID:          assetListSetRequest.AssetID,
		AssetType:        assetListSetRequest.AssetType,
		OutputIndex:      assetListSetRequest.OutputIndex,
		Amount:           assetListSetRequest.Amount,
		LockTime:         assetListSetRequest.LockTime,
		RelativeLockTime: assetListSetRequest.RelativeLockTime,
		ScriptKey:        assetListSetRequest.ScriptKey,
		AnchorOutpoint:   assetListSetRequest.AnchorOutpoint,
		TweakedGroupKey:  assetListSetRequest.TweakedGroupKey,
		DeviceId:         assetListSetRequest.DeviceId,
		UserId:           userId,
		Username:         username,
	}
	return &assetList
}

func IsAssetListChanged(assetListByAssetId *models.AssetList, old *models.AssetList) bool {
	if assetListByAssetId.Version != old.Version {
		return true
	}
	if assetListByAssetId.GenesisPoint != old.GenesisPoint {
		return true
	}
	if assetListByAssetId.Name != old.Name {
		return true
	}
	if assetListByAssetId.MetaHash != old.MetaHash {
		return true
	}
	if assetListByAssetId.AssetID != old.AssetID {
		return true
	}
	if assetListByAssetId.AssetType != old.AssetType {
		return true
	}
	if assetListByAssetId.OutputIndex != old.OutputIndex {
		return true
	}
	if assetListByAssetId.Amount != old.Amount {
		return true
	}
	if assetListByAssetId.LockTime != old.LockTime {
		return true
	}
	if assetListByAssetId.RelativeLockTime != old.RelativeLockTime {
		return true
	}
	if assetListByAssetId.ScriptKey != old.ScriptKey {
		return true
	}
	if assetListByAssetId.AnchorOutpoint != old.AnchorOutpoint {
		return true
	}
	if assetListByAssetId.TweakedGroupKey != old.TweakedGroupKey {
		return true
	}
	if assetListByAssetId.DeviceId != old.DeviceId {
		return true
	}
	if assetListByAssetId.UserId != old.UserId {
		return true
	}
	if assetListByAssetId.Username != old.Username {
		return true
	}
	return false
}

func CheckAssetListIfUpdate(assetList *models.AssetList, userId int) (*models.AssetList, error) {
	if assetList == nil {
		return nil, errors.New("nil asset list")
	}
	assetListByAssetId, err := btldb.ReadAssetListByAssetIdAndUserId(assetList.AssetID, userId)
	if err != nil {
		return assetList, nil
	}
	if !IsAssetListChanged(assetListByAssetId, assetList) {
		return assetListByAssetId, nil
	}
	assetListByAssetId.Version = assetList.Version
	assetListByAssetId.GenesisPoint = assetList.GenesisPoint
	assetListByAssetId.Name = assetList.Name
	assetListByAssetId.MetaHash = assetList.MetaHash
	assetListByAssetId.AssetID = assetList.AssetID
	assetListByAssetId.AssetType = assetList.AssetType
	assetListByAssetId.OutputIndex = assetList.OutputIndex
	assetListByAssetId.Amount = assetList.Amount
	assetListByAssetId.LockTime = assetList.LockTime
	assetListByAssetId.RelativeLockTime = assetList.RelativeLockTime
	assetListByAssetId.ScriptKey = assetList.ScriptKey
	assetListByAssetId.AnchorOutpoint = assetList.AnchorOutpoint
	assetListByAssetId.TweakedGroupKey = assetList.TweakedGroupKey
	assetListByAssetId.DeviceId = assetList.DeviceId
	assetListByAssetId.UserId = assetList.UserId
	assetListByAssetId.Username = assetList.Username
	return assetListByAssetId, nil
}

func CreateOrUpdateAssetList(list *models.AssetList, userId int) (err error) {
	var assetList *models.AssetList
	assetList, err = CheckAssetListIfUpdate(list, userId)
	return btldb.UpdateAssetList(assetList)
}

func ProcessAssetListSetRequestSlice(userId int, username string, assetListSetRequestSlice *[]models.AssetListSetRequest) *[]models.AssetList {
	var assetLists []models.AssetList
	for _, assetListRequest := range *assetListSetRequestSlice {
		assetLists = append(assetLists, *ProcessAssetListSetRequest(userId, username, &assetListRequest))
	}
	return &assetLists
}

func CreateOrUpdateAssetLists(lists *[]models.AssetList, userId int) (err error) {
	var assetLists []models.AssetList
	var assetList *models.AssetList
	for _, balance := range *lists {
		assetList, err = CheckAssetListIfUpdate(&balance, userId)
		if err != nil {
			return err
		}
		assetLists = append(assetLists, *assetList)
	}
	return btldb.UpdateAssetLists(&assetLists)
}

func GetAssetListByAssetIdAndUsername(assetId string, username string) (*models.AssetList, error) {
	return btldb.ReadAssetListByAssetIdAndUsername(assetId, username)
}

func IsAssetListRecordExist(assetId string, username string) (bool, error) {
	assetList, err := GetAssetListByAssetIdAndUsername(assetId, username)
	if err != nil {
		return false, utils.AppendErrorInfo(err, "GetAssetListByAssetIdAndUsername")
	} else {
		return assetList.Amount != 0, nil
	}
}

type UserAssetList struct {
	UserId     int                 `json:"user_id"`
	AssetLists *[]models.AssetList `json:"asset_balances"`
}

type UsernameAssetList struct {
	Username   string              `json:"username"`
	AssetLists *[]models.AssetList `json:"asset_balances"`
}

func GetAllAssetLists() (*[]models.AssetList, error) {
	return btldb.ReadAllAssetLists()
}

func GetAllAssetListsNonZeroUpdatedAtDesc() (*[]models.AssetList, error) {
	return btldb.ReadAllAssetListsNonZeroUpdatedAtDesc()
}

func GetAllAssetListsNonZero() (*[]models.AssetList, error) {
	return btldb.ReadAllAssetListsNonZero()
}

func GetAllAssetListsNonZeroByAssetId(assetId string) (*[]models.AssetList, error) {
	return btldb.ReadAllAssetListsNonZeroByAssetId(assetId)
}

func GetAllAssetListsNonZeroLimitAndOffset(limit int, offset int) (*[]models.AssetList, error) {
	return btldb.ReadAllAssetListsNonZeroLimitAndOffset(limit, offset)
}

func AssetListsToUserMapAssetLists(assetLists *[]models.AssetList) *map[int]*[]models.AssetList {
	userMapAssetLists := make(map[int]*[]models.AssetList)
	for _, assetList := range *assetLists {
		balances, ok := userMapAssetLists[assetList.UserId]
		if !ok {
			userMapAssetLists[assetList.UserId] = &[]models.AssetList{assetList}
		} else {
			*balances = append(*balances, assetList)
		}
	}
	return &userMapAssetLists
}

func AssetListsToUsernameMapAssetLists(assetLists *[]models.AssetList) *map[string]*[]models.AssetList {
	usernameMapBalances := make(map[string]*[]models.AssetList)
	for _, assetList := range *assetLists {
		balances, ok := usernameMapBalances[assetList.Username]
		if !ok {
			usernameMapBalances[assetList.Username] = &[]models.AssetList{assetList}
		} else {
			*balances = append(*balances, assetList)
		}
	}
	return &usernameMapBalances
}

func UserMapAssetListsToUserAssetLists(userMapAssetLists *map[int]*[]models.AssetList) *[]UserAssetList {
	var userAssetLists []UserAssetList
	for userId, assetLists := range *userMapAssetLists {
		userAssetLists = append(userAssetLists, UserAssetList{
			UserId:     userId,
			AssetLists: assetLists,
		})
	}
	return &userAssetLists
}

func UsernameMapAssetListsToUsernameAssetLists(userMapAssetLists *map[string]*[]models.AssetList) *[]UsernameAssetList {
	var usernameAssetLists []UsernameAssetList
	for username, assetLists := range *userMapAssetLists {
		usernameAssetLists = append(usernameAssetLists, UsernameAssetList{
			Username:   username,
			AssetLists: assetLists,
		})
	}
	return &usernameAssetLists
}

func AssetListsToUserAssetLists(assetLists *[]models.AssetList) *[]UserAssetList {
	userMapAssetLists := AssetListsToUserMapAssetLists(assetLists)
	userAssetLists := UserMapAssetListsToUserAssetLists(userMapAssetLists)
	return userAssetLists
}

func AssetListsToUsernameAssetLists(assetLists *[]models.AssetList) *[]UsernameAssetList {
	usernameMapAssetLists := AssetListsToUsernameMapAssetLists(assetLists)
	usernameAssetLists := UsernameMapAssetListsToUsernameAssetLists(usernameMapAssetLists)
	return usernameAssetLists
}

func GetAllUserAssetLists() (*[]UserAssetList, error) {
	allAssetLists, err := GetAllAssetLists()
	if err != nil {
		return nil, err
	}
	userAssetLists := AssetListsToUserAssetLists(allAssetLists)
	return userAssetLists, nil
}

func GetAllUsernameAssetLists() (*[]UsernameAssetList, error) {
	allAssetLists, err := GetAllAssetLists()
	if err != nil {
		return nil, err
	}
	usernameAssetLists := AssetListsToUsernameAssetLists(allAssetLists)
	return usernameAssetLists, nil
}
