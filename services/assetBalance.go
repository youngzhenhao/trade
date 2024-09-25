package services

import (
	"errors"
	"math"
	"sort"
	"time"
	"trade/models"
	"trade/services/btldb"
)

func GetAssetBalancesByUserId(userId int) (*[]models.AssetBalance, error) {
	return btldb.ReadAssetBalancesByUserId(userId)
}

func GetAssetBalancesByUserIdNonZero(userId int) (*[]models.AssetBalance, error) {
	return btldb.ReadAssetBalancesByUserIdNonZero(userId)
}

func ProcessAssetBalanceSetRequest(userId int, username string, assetBalanceSetRequest *models.AssetBalanceSetRequest) *models.AssetBalance {
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
		Username:     username,
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
	if assetBalanceByInvoice.Username != old.Username {
		return true
	}
	return false
}

func CheckAssetBalanceIfUpdate(assetBalance *models.AssetBalance, userId int) (*models.AssetBalance, error) {
	if assetBalance == nil {
		return nil, errors.New("nil asset balance")
	}
	assetBalanceByAssetId, err := btldb.ReadAssetBalanceByAssetIdAndUserId(assetBalance.AssetID, userId)
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
	assetBalanceByAssetId.Username = assetBalance.Username
	return assetBalanceByAssetId, nil
}

func CreateOrUpdateAssetBalance(balance *models.AssetBalance, userId int) (err error) {
	var assetBalance *models.AssetBalance
	assetBalance, err = CheckAssetBalanceIfUpdate(balance, userId)
	return btldb.UpdateAssetBalance(assetBalance)
}

func ProcessAssetBalanceSetRequestSlice(userId int, username string, assetBalanceSetRequestSlice *[]models.AssetBalanceSetRequest) *[]models.AssetBalance {
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
			Username:     username,
		})
	}
	return &assetBalances
}

func CreateOrUpdateAssetBalances(balances *[]models.AssetBalance, userId int) (err error) {
	var assetBalances []models.AssetBalance
	var assetBalance *models.AssetBalance
	for _, balance := range *balances {
		assetBalance, err = CheckAssetBalanceIfUpdate(&balance, userId)
		if err != nil {
			return err
		}
		assetBalances = append(assetBalances, *assetBalance)
	}
	return btldb.UpdateAssetBalances(&assetBalances)
}

type UserAssetBalance struct {
	UserId        int                    `json:"user_id"`
	AssetBalances *[]models.AssetBalance `json:"asset_balances"`
}

type UsernameAssetBalance struct {
	Username      string                 `json:"username"`
	AssetBalances *[]models.AssetBalance `json:"asset_balances"`
}

func GetAllAssetBalances() (*[]models.AssetBalance, error) {
	return btldb.ReadAllAssetBalances()
}

func GetAllAssetBalancesNonZeroUpdatedAtDesc() (*[]models.AssetBalance, error) {
	return btldb.ReadAllAssetBalancesNonZeroUpdatedAtDesc()
}

func GetAllAssetBalancesNonZero() (*[]models.AssetBalance, error) {
	return btldb.ReadAllAssetBalancesNonZero()
}

func GetAllAssetBalancesNonZeroByAssetId(assetId string) (*[]models.AssetBalance, error) {
	return btldb.ReadAllAssetBalancesNonZeroByAssetId(assetId)
}

// Deprecated: Use GetAssetIdAndBalancesByAssetIdLimitAndOffset instead
func GetAllAssetBalancesNonZeroLimit(limit int) (*[]models.AssetBalance, error) {
	return btldb.ReadAllAssetBalancesNonZeroLimit(limit)
}

func GetAllAssetBalancesNonZeroLimitAndOffset(limit int, offset int) (*[]models.AssetBalance, error) {
	return btldb.ReadAllAssetBalancesNonZeroLimitAndOffset(limit, offset)
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

func AssetBalancesToUsernameMapAssetBalances(assetBalances *[]models.AssetBalance) *map[string]*[]models.AssetBalance {
	usernameMapBalances := make(map[string]*[]models.AssetBalance)
	for _, assetBalance := range *assetBalances {
		balances, ok := usernameMapBalances[assetBalance.Username]
		if !ok {
			usernameMapBalances[assetBalance.Username] = &[]models.AssetBalance{assetBalance}
		} else {
			*balances = append(*balances, assetBalance)
		}
	}
	return &usernameMapBalances
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

func UsernameMapAssetBalancesToUsernameAssetBalances(userMapAssetBalances *map[string]*[]models.AssetBalance) *[]UsernameAssetBalance {
	var usernameAssetBalances []UsernameAssetBalance
	for username, assetBalances := range *userMapAssetBalances {
		usernameAssetBalances = append(usernameAssetBalances, UsernameAssetBalance{
			Username:      username,
			AssetBalances: assetBalances,
		})
	}
	return &usernameAssetBalances
}

func AssetBalancesToUserAssetBalances(assetBalances *[]models.AssetBalance) *[]UserAssetBalance {
	userMapAssetBalances := AssetBalancesToUserMapAssetBalances(assetBalances)
	userAssetBalances := UserMapAssetBalancesToUserAssetBalances(userMapAssetBalances)
	return userAssetBalances
}

func AssetBalancesToUsernameAssetBalances(assetBalances *[]models.AssetBalance) *[]UsernameAssetBalance {
	usernameMapAssetBalances := AssetBalancesToUsernameMapAssetBalances(assetBalances)
	usernameAssetBalances := UsernameMapAssetBalancesToUsernameAssetBalances(usernameMapAssetBalances)
	return usernameAssetBalances
}

// GetAllUserAssetBalances
// @dev: UserId
// @Description: Get all asset balances by userId
func GetAllUserAssetBalances() (*[]UserAssetBalance, error) {
	allAssetBalances, err := GetAllAssetBalances()
	if err != nil {
		return nil, err
	}
	userAssetBalances := AssetBalancesToUserAssetBalances(allAssetBalances)
	return userAssetBalances, nil
}

// GetAllUsernameAssetBalances
// @dev: Username
// @Description: Get all username asset balances
func GetAllUsernameAssetBalances() (*[]UsernameAssetBalance, error) {
	allAssetBalances, err := GetAllAssetBalances()
	if err != nil {
		return nil, err
	}
	usernameAssetBalances := AssetBalancesToUsernameAssetBalances(allAssetBalances)
	return usernameAssetBalances, nil
}

type AssetBalanceSimplifiedWithoutUsername struct {
	Name     string `json:"name" gorm:"type:varchar(255)"`
	AssetID  string `json:"asset_id" gorm:"type:varchar(255)"`
	Balance  int    `json:"balance"`
	DeviceId string `json:"device_id" gorm:"type:varchar(255)"`
}

type UsernameAssetBalanceSimplified struct {
	Username      string                                   `json:"username"`
	AssetBalances *[]AssetBalanceSimplifiedWithoutUsername `json:"asset_balances"`
}

func AssetBalanceToAssetBalanceSimplifiedWithoutUsername(assetBalance models.AssetBalance) AssetBalanceSimplifiedWithoutUsername {
	return AssetBalanceSimplifiedWithoutUsername{
		Name:     assetBalance.Name,
		AssetID:  assetBalance.AssetID,
		Balance:  assetBalance.Balance,
		DeviceId: assetBalance.DeviceId,
	}
}

func AssetBalanceSliceToAssetBalanceSimplifiedSliceWithoutUsername(assetBalances *[]models.AssetBalance) *[]AssetBalanceSimplifiedWithoutUsername {
	if assetBalances == nil {
		return nil
	}
	var assetBalanceSimplified []AssetBalanceSimplifiedWithoutUsername
	for _, assetBalance := range *assetBalances {
		assetBalanceSimplified = append(assetBalanceSimplified, AssetBalanceToAssetBalanceSimplifiedWithoutUsername(assetBalance))
	}
	return &assetBalanceSimplified
}

func UsernameAssetBalanceToUsernameAssetBalanceSimplified(usernameAssetBalance UsernameAssetBalance) *UsernameAssetBalanceSimplified {
	var usernameAssetBalanceSimplified UsernameAssetBalanceSimplified
	usernameAssetBalanceSimplified.Username = usernameAssetBalance.Username
	usernameAssetBalanceSimplified.AssetBalances = AssetBalanceSliceToAssetBalanceSimplifiedSliceWithoutUsername(usernameAssetBalance.AssetBalances)
	return &usernameAssetBalanceSimplified
}

func GetAllUsernameAssetBalanceSimplified() (*[]UsernameAssetBalanceSimplified, error) {
	allUsernameAssetBalances, err := GetAllUsernameAssetBalances()
	var usernameAssetBalanceSimplifiedSlice []UsernameAssetBalanceSimplified
	if err != nil {
		return nil, err
	}
	for _, usernameAssetBalance := range *allUsernameAssetBalances {
		usernameAssetBalanceSimplified := UsernameAssetBalanceToUsernameAssetBalanceSimplified(usernameAssetBalance)
		usernameAssetBalanceSimplifiedSlice = append(usernameAssetBalanceSimplifiedSlice, *usernameAssetBalanceSimplified)
	}
	return &usernameAssetBalanceSimplifiedSlice, nil
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
// @dev
func GetAllAssetIdAndBalances() (*[]AssetIdAndBalance, error) {
	allAssetBalances, err := GetAllAssetBalances()
	if err != nil {
		return nil, err
	}
	assetIdAndBalances := AssetBalancesToAssetIdAndBalances(allAssetBalances)
	return assetIdAndBalances, nil
}

type AssetBalanceSimplified struct {
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name" gorm:"type:varchar(255)"`
	AssetID   string    `json:"asset_id" gorm:"type:varchar(255)"`
	Balance   int       `json:"balance"`
	DeviceId  string    `json:"device_id" gorm:"type:varchar(255)"`
	Username  string    `json:"username" gorm:"type:varchar(255)"`
}

type AssetIdAndBalanceSimplified struct {
	AssetId       string                    `json:"asset_id"`
	AssetBalances *[]AssetBalanceSimplified `json:"asset_balances"`
}

func AssetBalanceToAssetBalanceSimplified(assetBalance models.AssetBalance) AssetBalanceSimplified {
	return AssetBalanceSimplified{
		UpdatedAt: assetBalance.UpdatedAt,
		Name:      assetBalance.Name,
		AssetID:   assetBalance.AssetID,
		Balance:   assetBalance.Balance,
		DeviceId:  assetBalance.DeviceId,
		Username:  assetBalance.Username,
	}
}

func AssetBalanceSliceToAssetBalanceSimplifiedSlice(assetBalances *[]models.AssetBalance) *[]AssetBalanceSimplified {
	if assetBalances == nil {
		return nil
	}
	var assetBalanceSimplified []AssetBalanceSimplified
	for _, assetBalance := range *assetBalances {
		assetBalanceSimplified = append(assetBalanceSimplified, AssetBalanceToAssetBalanceSimplified(assetBalance))
	}
	return &assetBalanceSimplified
}

func AssetIdAndBalanceToAssetIdAndBalanceSimplified(assetIdAndBalance AssetIdAndBalance) AssetIdAndBalanceSimplified {
	return AssetIdAndBalanceSimplified{
		AssetId:       assetIdAndBalance.AssetId,
		AssetBalances: AssetBalanceSliceToAssetBalanceSimplifiedSlice(assetIdAndBalance.AssetBalances),
	}
}

func AssetIdAndBalanceSliceToAssetIdAndBalanceSimplifiedSlice(assetIdAndBalances *[]AssetIdAndBalance) *[]AssetIdAndBalanceSimplified {
	if assetIdAndBalances == nil {
		return nil
	}
	var assetIdAndBalanceSimplified []AssetIdAndBalanceSimplified
	for _, assetIdAndBalance := range *assetIdAndBalances {
		assetIdAndBalanceSimplified = append(assetIdAndBalanceSimplified, AssetIdAndBalanceToAssetIdAndBalanceSimplified(assetIdAndBalance))
	}
	return &assetIdAndBalanceSimplified
}

func GetAssetIdSliceFromAssetIdAndBalanceSimplifiedSliceSort(assetIdAndBalanceSimplifiedSlice *[]AssetIdAndBalanceSimplified) []string {
	var assetIdSlice []string
	for _, assetIdAndBalance := range *assetIdAndBalanceSimplifiedSlice {
		assetIdSlice = append(assetIdSlice, assetIdAndBalance.AssetId)
	}
	// @dev: Sort string slice
	sort.Strings(assetIdSlice)
	return assetIdSlice
}

func AssetIdMapBalanceSimplifiedToAssetIdSlice(assetIdMapBalanceSimplified *map[string]*[]AssetBalanceSimplified) []string {
	var assetIdSlice []string
	for assetId, _ := range *assetIdMapBalanceSimplified {
		assetIdSlice = append(assetIdSlice, assetId)
	}
	// @dev: Sort string slice
	sort.Strings(assetIdSlice)
	return assetIdSlice
}

func AssetIdAndBalanceSimplifiedSliceToAssetIdMapBalanceSimplified(assetIdAndBalanceSimplified *[]AssetIdAndBalanceSimplified) *map[string]*[]AssetBalanceSimplified {
	if assetIdAndBalanceSimplified == nil {
		return nil
	}
	assetIdMapBalanceSimplified := make(map[string]*[]AssetBalanceSimplified)
	for _, assetIdAndBalance := range *assetIdAndBalanceSimplified {
		assetIdMapBalanceSimplified[assetIdAndBalance.AssetId] = assetIdAndBalance.AssetBalances
	}
	return &assetIdMapBalanceSimplified
}

func SortAssetIdAndBalanceSimplifiedSlice(assetIdAndBalanceSimplified *[]AssetIdAndBalanceSimplified) *[]AssetIdAndBalanceSimplified {
	if assetIdAndBalanceSimplified == nil {
		return nil
	}
	assetIdMapBalanceSimplified := *AssetIdAndBalanceSimplifiedSliceToAssetIdMapBalanceSimplified(assetIdAndBalanceSimplified)
	assetIdSlice := GetAssetIdSliceFromAssetIdAndBalanceSimplifiedSliceSort(assetIdAndBalanceSimplified)
	var assetIdAndBalanceSimplifiedSort []AssetIdAndBalanceSimplified
	for _, assetId := range assetIdSlice {
		assetIdAndBalanceSimplifiedSort = append(assetIdAndBalanceSimplifiedSort, AssetIdAndBalanceSimplified{
			AssetId:       assetId,
			AssetBalances: assetIdMapBalanceSimplified[assetId],
		})
	}
	return &assetIdAndBalanceSimplifiedSort
}

// GetAllAssetIdAndBalanceSimplifiedSort
// @Description: Get all asset id and balance simplified
func GetAllAssetIdAndBalanceSimplifiedSort() (*[]AssetIdAndBalanceSimplified, error) {
	assetIdAndBalances, err := GetAllAssetIdAndBalances()
	if err != nil {
		return nil, err
	}
	assetIdAndBalanceSimplified := AssetIdAndBalanceSliceToAssetIdAndBalanceSimplifiedSlice(assetIdAndBalances)
	// @dev: Sort by asset id
	assetIdAndBalanceSimplified = SortAssetIdAndBalanceSimplifiedSlice(assetIdAndBalanceSimplified)
	return assetIdAndBalanceSimplified, nil
}

type AssetIdAndUserAssetBalance struct {
	AssetId          string              `json:"asset_id"`
	UserAssetBalance *[]UserAssetBalance `json:"user_asset_balance"`
}

// GetAllAssetIdAndUserAssetBalances
// @Description: Get all asset balances by assetId and userId
func GetAllAssetIdAndUserAssetBalances() (*[]AssetIdAndUserAssetBalance, error) {
	var assetIdAndUserAssetBalances []AssetIdAndUserAssetBalance
	allAssetBalances, err := GetAllAssetBalancesNonZero()
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

type AssetHolderNumber struct {
	AssetId   string `json:"asset_id"`
	HolderNum int    `json:"holder_num"`
}

func AllAssetIdAndUserAssetBalancesToAssetHolderInfos(assetIdAndUserAssetBalances *[]AssetIdAndUserAssetBalance) *[]AssetHolderNumber {
	var assetHolderInfos []AssetHolderNumber
	for _, asset := range *assetIdAndUserAssetBalances {
		assetHolderInfos = append(assetHolderInfos, AssetHolderNumber{
			AssetId:   asset.AssetId,
			HolderNum: len(*(asset.UserAssetBalance)),
		})
	}
	return &assetHolderInfos
}

func GetAssetHolderInfosByAssetBalances() (*[]AssetHolderNumber, error) {
	assetIdAndUserAssetBalance, err := GetAllAssetIdAndUserAssetBalances()
	if err != nil {
		return nil, err
	}
	assetHolderInfos := AllAssetIdAndUserAssetBalancesToAssetHolderInfos(assetIdAndUserAssetBalance)
	return assetHolderInfos, nil
}

func GetAssetHolderNumberByAssetIdWithAssetBalances(assetId string) (int, error) {
	assetHolderInfos, err := GetAssetHolderInfosByAssetBalances()
	if err != nil {
		return 0, err
	}
	for _, asset := range *assetHolderInfos {
		if asset.AssetId == assetId {
			return asset.HolderNum, nil
		}
	}
	// @dev: Asset holder info not found
	return 0, nil
}

// GetAssetHolderNumberAssetBalance
// @Description: Use asset balances
func GetAssetHolderNumberAssetBalance(assetId string) (int, error) {
	return GetAssetHolderNumberByAssetIdWithAssetBalances(assetId)
}

func GetAssetIdAndBalanceSimplifiedByAssetIdUpdatedAtDesc(assetId string) (*AssetIdAndBalance, error) {
	allAssetBalances, err := GetAllAssetBalancesNonZeroUpdatedAtDesc()
	if err != nil {
		return nil, err
	}
	assetIdMapAssetBalances := AssetBalancesToAssetIdMapAssetBalances(allAssetBalances)
	assetBalances, ok := (*assetIdMapAssetBalances)[assetId]
	if !ok {
		return &AssetIdAndBalance{
			AssetId:       assetId,
			AssetBalances: &[]models.AssetBalance{},
		}, nil
	}
	return &AssetIdAndBalance{
		AssetId:       assetId,
		AssetBalances: assetBalances,
	}, nil
}

// GetAssetIdAndBalancesByAssetId
// @Description: Get assetId and balances by assetId
// @dev
func GetAssetIdAndBalancesByAssetId(assetId string) (*AssetIdAndBalance, error) {
	allAssetBalances, err := GetAllAssetBalancesNonZero()
	if err != nil {
		return nil, err
	}
	assetIdMapAssetBalances := AssetBalancesToAssetIdMapAssetBalances(allAssetBalances)
	assetBalances, ok := (*assetIdMapAssetBalances)[assetId]
	if !ok {
		return &AssetIdAndBalance{
			AssetId:       assetId,
			AssetBalances: &[]models.AssetBalance{},
		}, nil
	}
	return &AssetIdAndBalance{
		AssetId:       assetId,
		AssetBalances: assetBalances,
	}, nil
}

func GetAssetIdAndAssetBalancesByAssetId(assetId string) (*AssetIdAndBalance, error) {
	assetBalances, err := GetAllAssetBalancesNonZeroByAssetId(assetId)
	if err != nil {
		return nil, err
	}
	if assetBalances == nil || len(*(assetBalances)) == 0 {
		return &AssetIdAndBalance{
			AssetId:       assetId,
			AssetBalances: &[]models.AssetBalance{},
		}, nil
	}
	return &AssetIdAndBalance{
		AssetId:       assetId,
		AssetBalances: assetBalances,
	}, nil
}

func GetAssetIdAndBalancesByAssetIdLimitAndOffset(assetId string, limit int, offset int) (*AssetIdAndBalance, error) {
	// @dev: Do not use GetAllAssetBalancesNonZeroLimitAndOffset(limit, offset)
	// @dev: Do not use AssetBalancesToAssetIdMapAssetBalances(allAssetBalances)
	assetBalances, err := btldb.ReadAssetBalanceByAssetIdNonZeroLimitAndOffset(assetId, limit, offset)
	if err != nil {
		return nil, err
	}
	if assetBalances == nil || len(*(assetBalances)) == 0 {
		return &AssetIdAndBalance{
			AssetId:       assetId,
			AssetBalances: &[]models.AssetBalance{},
		}, nil
	}
	return &AssetIdAndBalance{
		AssetId:       assetId,
		AssetBalances: assetBalances,
	}, nil
}

func GetAssetBalanceByAssetIdNonZero(assetId string) (*[]models.AssetBalance, error) {
	return btldb.ReadAssetBalanceByAssetIdNonZero(assetId)
}

// GetAssetBalanceByAssetIdNonZeroLength
// @Description: Get asset balance by asset id non-zero length
func GetAssetBalanceByAssetIdNonZeroLength(assetId string) (int, error) {
	response, err := GetAssetBalanceByAssetIdNonZero(assetId)
	if err != nil {
		return 0, err
	}
	if response == nil || len(*(response)) == 0 {
		return 0, nil
	}
	return len(*response), nil
}

// IsLimitAndOffsetValid
// @dev: Check limit and offset is valid by total amount
func IsLimitAndOffsetValid(assetId string, limit int, offset int) (bool, error) {
	if !(limit > 0 && offset >= 0) {
		return false, errors.New("invalid limit or offset")
	}
	recordsNum, err := GetAssetBalanceByAssetIdNonZeroLength(assetId)
	if err != nil {
		return false, err
	}
	if recordsNum == 0 && offset == 0 {
		return true, nil
	}
	return recordsNum > offset, nil
}

type GetAssetHolderBalancePageNumberRequest struct {
	AssetId  string `json:"asset_id"`
	PageSize int    `json:"page_size"`
}

func GetAssetHolderBalancePageNumberByPageSize(assetId string, pageSize int) (pageNumber int, err error) {
	recordsNum, err := GetAssetBalanceByAssetIdNonZeroLength(assetId)
	if err != nil {
		return 0, err
	}
	return int(math.Ceil(float64(recordsNum) / float64(pageSize))), nil
}

// @dev: Use receives and transfers
// @dev: Rat stands for Receices and transfers
type UserAssetBalanceByRat struct {
	UserId             int `json:"user_id"`
	AssetBalanceAmount int `json:"asset_balance_amount"`
}

type AssetIdAndUserAssetBalanceByRat struct {
	AssetId          string                   `json:"asset_id"`
	UserAssetBalance *[]UserAssetBalanceByRat `json:"user_asset_balance"`
}

// @dev: Use to maps to compute
// GetAllAddressAmountMapByRatPositiveAmount
func GetAssetIdAndUserAssetBalanceByRat() *[]AssetIdAndUserAssetBalanceByRat {
	// TODO: Compute asset Balance by receives and transfers' maps
	return nil
}

// GetAllAddressAmountMapByRat
// @Description: Get all address amount map by receives and transfers
func GetAllAddressAmountMapByRat(network models.Network) (*map[string]*AssetIdAndAmount, error) {
	addressAmountMap := make(map[string]*AssetIdAndAmount)
	receivesAddressAmountMap, err := AllAssetReceivesToAddressAmountMap(network)
	if err != nil {
		return nil, err
	}
	if receivesAddressAmountMap != nil {
		for address, assetIdAndAmount := range *receivesAddressAmountMap {
			_, ok := addressAmountMap[address]
			if !ok {
				addressAmountMap[address] = &AssetIdAndAmount{
					AssetId: assetIdAndAmount.AssetId,
				}
			}
			if (*(addressAmountMap[address])).AssetId == assetIdAndAmount.AssetId {
				(*(addressAmountMap[address])).Amount += assetIdAndAmount.Amount
			}
		}
	}
	transfersAddressAmountMap, err := AllAssetTransferCombinedToAddressAmountMap()
	if err != nil {
		return nil, err
	}
	if transfersAddressAmountMap != nil {
		for address, assetIdAndAmount := range *transfersAddressAmountMap {
			_, ok := addressAmountMap[address]
			if !ok {
				addressAmountMap[address] = &AssetIdAndAmount{
					AssetId: assetIdAndAmount.AssetId,
				}
			}
			if (*(addressAmountMap[address])).AssetId == assetIdAndAmount.AssetId {
				(*(addressAmountMap[address])).Amount += assetIdAndAmount.Amount
			}
		}
	}
	return &addressAmountMap, nil
}

// TODO: Get all address amount map by receives and transfers, then store data in db

// GetAllAddressAmountMapByRatPositiveAmount
// @Description: Filter zero and negative amount of asset address
// @dev: UTXO
func GetAllAddressAmountMapByRatPositiveAmount(network models.Network) (*map[string]*AssetIdAndAmount, error) {
	addressAmountMap := make(map[string]*AssetIdAndAmount)
	allAddressAmountMapByRat, err := GetAllAddressAmountMapByRat(network)
	if err != nil {
		return nil, err
	}
	for address, assetIdAndAmount := range *allAddressAmountMapByRat {
		if assetIdAndAmount.Amount > 0 {
			addressAmountMap[address] = &AssetIdAndAmount{
				AssetId: assetIdAndAmount.AssetId,
				Amount:  assetIdAndAmount.Amount,
			}
		}
	}
	return &addressAmountMap, nil
}

func GetAssetBalanceByAssetIdAndUserId(assetId string, userId int) (*models.AssetBalance, error) {
	return btldb.ReadAssetBalanceByAssetIdAndUserId(assetId, userId)
}
