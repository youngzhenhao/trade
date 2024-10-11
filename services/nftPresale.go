package services

import (
	"encoding/hex"
	"errors"
	"github.com/lightninglabs/taproot-assets/taprpc"
	"strconv"
	"strings"
	"trade/api"
	"trade/btlLog"
	"trade/models"
	"trade/services/btldb"
	"trade/services/custodyAccount"
	"trade/utils"
)

func CreateNftPresale(nftPresale *models.NftPresale) error {
	return btldb.CreateNftPresale(nftPresale)
}

func CreateNftPresales(nftPresales *[]models.NftPresale) error {
	return btldb.CreateNftPresales(nftPresales)
}

func ReadNftPresale(id uint) (*models.NftPresale, error) {
	return btldb.ReadNftPresale(id)
}

func ReadNftPresaleByAssetId(assetId string) (*models.NftPresale, error) {
	return btldb.ReadNftPresaleByAssetId(assetId)
}

func ReadAllNftPresales() (*[]models.NftPresale, error) {
	return btldb.ReadAllNftPresales()
}

func ReadNftPresalesByNftPresaleState(nftPresaleState models.NftPresaleState) (*[]models.NftPresale, error) {
	return btldb.ReadNftPresalesByNftPresaleState(nftPresaleState)
}

func ReadNftPresalesBetweenNftPresaleState(stateStart models.NftPresaleState, stateEnd models.NftPresaleState) (*[]models.NftPresale, error) {
	return btldb.ReadNftPresalesBetweenNftPresaleState(stateStart, stateEnd)
}

func ReadNftPresalesByBuyerUserId(userId int) (*[]models.NftPresale, error) {
	return btldb.ReadNftPresalesByBuyerUserId(userId)
}

func UpdateNftPresale(nftPresale *models.NftPresale) error {
	return btldb.UpdateNftPresale(nftPresale)
}

func UpdateNftPresales(nftPresales *[]models.NftPresale) error {
	return btldb.UpdateNftPresales(nftPresales)
}

func DeleteNftPresale(id uint) error {
	return btldb.DeleteNftPresale(id)
}

func ProcessNftPresale(nftPresaleSetRequest *models.NftPresaleSetRequest) *models.NftPresale {
	var assetId string
	assetId = nftPresaleSetRequest.AssetId
	var name string
	var assetType string
	var groupKey string
	var amount int
	var meta string
	assetInfo, err := api.GetAssetInfoApi(assetId)
	if err != nil {
		// @dev: Do not return
		btlLog.PreSale.Error("api GetAssetInfoApi err")
	} else {
		name = assetInfo.Name
		assetType = assetInfo.AssetType
		if assetInfo.GroupKey != nil {
			groupKey = *assetInfo.GroupKey
		}
		amount = int(assetInfo.Amount)
		if assetInfo.Meta != nil {
			meta = *assetInfo.Meta
		}
	}
	groupKeyByAssetId, err := api.GetGroupKeyByAssetId(assetId)
	if err != nil {
		btlLog.PreSale.Error("api GetGroupKeyByAssetId err")
	} else {
		groupKey = groupKeyByAssetId
	}
	return &models.NftPresale{
		AssetId:    assetId,
		Name:       name,
		AssetType:  assetType,
		Meta:       meta,
		GroupKey:   groupKey,
		Amount:     amount,
		Price:      nftPresaleSetRequest.Price,
		Info:       nftPresaleSetRequest.Info,
		LaunchTime: utils.GetTimestamp(),
		State:      models.NftPresaleStateLaunched,
	}
}

func ProcessNftPresales(nftPresaleSetRequests *[]models.NftPresaleSetRequest) *[]models.NftPresale {
	if nftPresaleSetRequests == nil {
		return nil
	}
	var nftPresales []models.NftPresale
	for _, nftPresaleSetRequest := range *nftPresaleSetRequests {
		nftPresales = append(nftPresales, *(ProcessNftPresale(&nftPresaleSetRequest)))
	}
	return &nftPresales
}

func GetNftPresaleByAssetId(assetId string) (*models.NftPresale, error) {
	return ReadNftPresaleByAssetId(assetId)
}

func GetLaunchedNftPresales() (*[]models.NftPresale, error) {
	return ReadNftPresalesByNftPresaleState(models.NftPresaleStateLaunched)
}

func GetNftPresalesByBuyerUserId(userId int) (*[]models.NftPresale, error) {
	return btldb.ReadNftPresalesByBuyerUserId(userId)
}

func IsNftPresalePurchasable(nftPresale *models.NftPresale) bool {
	if nftPresale == nil {
		return false
	}
	return nftPresale.State == models.NftPresaleStateLaunched
}

func IsNftPresaleAddrValid(nftPresale *models.NftPresale, addr *taprpc.Addr) (bool, error) {
	var err error
	if nftPresale == nil || addr == nil {
		err = errors.New("nftPresale or addr is nil")
		return false, err
	}
	addrAssetId := hex.EncodeToString(addr.AssetId)
	if addrAssetId != nftPresale.AssetId {
		err = errors.New("addrAssetId(" + addrAssetId + ") is not equal nftPresale.AssetId(" + nftPresale.AssetId + ")")
		return false, err
	}
	addrAssetType := addr.AssetType.String()
	if strings.ToLower(addrAssetType) != strings.ToLower(nftPresale.AssetType) {
		err = errors.New("addrAssetType(" + strings.ToLower(addrAssetType) + "[ToLower]) is not equal nftPresale.AssetType(" + strings.ToLower(nftPresale.AssetType) + "[ToLower])")
		return false, err
	}
	addrAmount := addr.Amount
	if addrAmount != uint64(nftPresale.Amount) {
		err = errors.New("addrAmount(" + strconv.FormatUint(addrAmount, 10) + ") is not equal nftPresale.Amount(" + strconv.Itoa(nftPresale.Amount) + ")")
		return false, err
	}
	addrGroupKey := hex.EncodeToString(addr.GroupKey)
	var isGroupKeyEqual bool
	// Without prefix
	if len(nftPresale.GroupKey) == 64 {
		isGroupKeyEqual = strings.Contains(addrGroupKey, nftPresale.GroupKey)
		// @dev: With prefix 0x02 or 0x03
	} else if len(nftPresale.GroupKey) == 66 {
		isGroupKeyEqual = addrGroupKey == nftPresale.GroupKey
	} else {
		isGroupKeyEqual = addrGroupKey == nftPresale.GroupKey
	}
	if !isGroupKeyEqual {
		err = errors.New("addrGroupKey(" + addrGroupKey + ") is not equal or contains nftPresale.GroupKey(" + nftPresale.GroupKey + ")")
		return false, err
	}
	return true, nil
}

func UpdateNftPresaleByPurchaseInfo(userId int, username string, deviceId string, addr string, nftPresale *models.NftPresale) error {
	var err error
	if nftPresale == nil {
		err = errors.New("nftPresale is nil")
		return err
	}
	nftPresale.BuyerUserId = userId
	nftPresale.BuyerUsername = username
	nftPresale.BuyerDeviceId = deviceId
	nftPresale.ReceiveAddr = addr
	nftPresale.BoughtTime = utils.GetTimestamp()
	nftPresale.State = models.NftPresaleStateBoughtNotPay
	err = UpdateNftPresale(nftPresale)
	if err != nil {
		return utils.AppendErrorInfo(err, "UpdateNftPresale")
	}
	return nil
}

// BuyNftPresale
// @Description: Buy presale nft
func BuyNftPresale(userId int, username string, buyNftPresaleRequest models.BuyNftPresaleRequest) error {
	assetId := buyNftPresaleRequest.AssetId
	nftPresale, err := GetNftPresaleByAssetId(assetId)
	if err != nil {
		return utils.AppendErrorInfo(err, "GetNftPresaleByAssetId")
	}
	// @dev: 1. Check if state is launched so that nft is purchasable
	if !IsNftPresalePurchasable(nftPresale) {
		err = errors.New("nft(" + nftPresale.AssetId + ") is not purchasable, its state is " + nftPresale.State.String() + ", ")
		return utils.AppendErrorInfo(err, "IsNftPresalePurchasable")
	}
	// @dev: 2. Check if account balance is enough
	price := nftPresale.Price
	{
		accountBalance, err := custodyAccount.GetAccountBalance(uint(userId))
		if err != nil {
			return utils.AppendErrorInfo(err, "GetAccountBalance")
		}
		isEnough := accountBalance >= int64(price)
		if !isEnough {
			err = errors.New("user(" + strconv.Itoa(userId) + ")'s account balance(" + strconv.FormatInt(accountBalance, 10) + ") not enough to pay nft presale price" + "(" + strconv.Itoa(nftPresale.Price) + ")")
			return utils.AppendErrorInfo(err, "IsAccountBalanceEnough")
		}
	}
	addr := buyNftPresaleRequest.ReceiveAddr
	// @dev: Decode addr and check if encoded addr is valid
	decodedAddrInfo, err := api.GetDecodedAddrInfo(addr)
	if err != nil {
		return utils.AppendErrorInfo(err, "GetDecodedAddrInfo")
	}
	// @dev: Return error if addr is invalid
	isNftPresaleAddrValid, err := IsNftPresaleAddrValid(nftPresale, decodedAddrInfo)
	if err != nil || !isNftPresaleAddrValid {
		return utils.AppendErrorInfo(err, "IsNftPresaleAddrValid")
	}
	// @dev: 4. Update info
	deviceId := buyNftPresaleRequest.DeviceId
	err = UpdateNftPresaleByPurchaseInfo(userId, username, deviceId, addr, nftPresale)
	if err != nil {
		return utils.AppendErrorInfo(err, "UpdateNftPresaleByPurchaseInfo")
	}
	return nil
}

func NftPresaleToNftPresaleSimplified(nftPresale *models.NftPresale) *models.NftPresaleSimplified {
	if nftPresale == nil {
		return nil
	}
	return &models.NftPresaleSimplified{
		ID:              nftPresale.ID,
		UpdatedAt:       nftPresale.UpdatedAt,
		AssetId:         nftPresale.AssetId,
		Name:            nftPresale.Name,
		AssetType:       nftPresale.AssetType,
		Meta:            nftPresale.Meta,
		GroupKey:        nftPresale.GroupKey,
		Amount:          nftPresale.Amount,
		Price:           nftPresale.Price,
		Info:            nftPresale.Info,
		BuyerUserId:     nftPresale.BuyerUserId,
		BuyerUsername:   nftPresale.BuyerUsername,
		BuyerDeviceId:   nftPresale.BuyerDeviceId,
		ReceiveAddr:     nftPresale.ReceiveAddr,
		PayMethod:       nftPresale.PayMethod,
		LaunchTime:      nftPresale.LaunchTime,
		BoughtTime:      nftPresale.BoughtTime,
		PaidId:          nftPresale.PaidId,
		PaidSuccessTime: nftPresale.PaidSuccessTime,
		SentTime:        nftPresale.SentTime,
		State:           nftPresale.State,
		ProcessNumber:   nftPresale.ProcessNumber,
	}
}

func NftPresaleSliceToNftPresaleSimplifiedSlice(nftPresales *[]models.NftPresale) *[]models.NftPresaleSimplified {
	if nftPresales == nil {
		return nil
	}
	var nftPresaleSimplifiedSlice []models.NftPresaleSimplified
	nftPresaleSimplifiedSlice = make([]models.NftPresaleSimplified, 0)
	for _, nftPresale := range *nftPresales {
		nftPresaleSimplifiedSlice = append(nftPresaleSimplifiedSlice, *(NftPresaleToNftPresaleSimplified(&nftPresale)))
	}
	return &nftPresaleSimplifiedSlice
}

// TODO: scheduled task Process NftPresale
// TODO: Refer fair launch
