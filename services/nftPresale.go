package services

import (
	"encoding/hex"
	"errors"
	"github.com/lightninglabs/taproot-assets/taprpc"
	"strconv"
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
		err = errors.New("addrAssetId is not equal nftPresale.AssetId")
		return false, err
	}
	addrAssetType := addr.AssetType.String()
	if addrAssetType != nftPresale.AssetType {
		err = errors.New("addrAssetType is not equal nftPresale.AssetType")
		return false, err
	}
	addrAmount := int(addr.Amount)
	if addrAmount != nftPresale.Amount {
		err = errors.New("addrAmount is not equal nftPresale.Amount")
		return false, err
	}
	addrGroupKey := hex.EncodeToString(addr.GroupKey)
	if addrGroupKey != nftPresale.GroupKey {
		err = errors.New("addrGroupKey is not equal nftPresale.GroupKey")
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

// TODO: Process NftPresale
