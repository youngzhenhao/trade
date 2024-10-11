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

// GetLaunchedNftPresales
// @Description: Get launched nftPresales
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

func UpdateNftPresaleByPurchaseInfo(userId int, username string, deviceId string, addr string, scriptKey string, internalKey string, nftPresale *models.NftPresale) error {
	var err error
	if nftPresale == nil {
		err = errors.New("nftPresale is nil")
		return err
	}
	nftPresale.BuyerUserId = userId
	nftPresale.BuyerUsername = username
	nftPresale.BuyerDeviceId = deviceId
	nftPresale.ReceiveAddr = addr
	nftPresale.AddrScriptKey = scriptKey
	nftPresale.AddrInternalKey = internalKey
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
	scriptKey := hex.EncodeToString(decodedAddrInfo.ScriptKey)
	internalKey := hex.EncodeToString(decodedAddrInfo.InternalKey)
	err = UpdateNftPresaleByPurchaseInfo(userId, username, deviceId, addr, scriptKey, internalKey, nftPresale)
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
		AddrScriptKey:   nftPresale.AddrScriptKey,
		AddrInternalKey: nftPresale.AddrInternalKey,
		PayMethod:       nftPresale.PayMethod,
		LaunchTime:      nftPresale.LaunchTime,
		BoughtTime:      nftPresale.BoughtTime,
		PaidId:          nftPresale.PaidId,
		PaidSuccessTime: nftPresale.PaidSuccessTime,
		SentTime:        nftPresale.SentTime,
		SentTxid:        nftPresale.SentTxid,
		SentOutpoint:    nftPresale.SentOutpoint,
		SentAddress:     nftPresale.SentAddress,
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

// @dev: Get all nftPresales

func GetAllNftPresaleStateBoughtNotPay() (*[]models.NftPresale, error) {
	return btldb.ReadNftPresalesByNftPresaleState(models.NftPresaleStateBoughtNotPay)
}

func GetAllNftPresaleStatePaidPending() (*[]models.NftPresale, error) {
	return btldb.ReadNftPresalesByNftPresaleState(models.NftPresaleStatePaidPending)
}

func GetAllNftPresaleStatePaidNotSend() (*[]models.NftPresale, error) {
	return btldb.ReadNftPresalesByNftPresaleState(models.NftPresaleStatePaidNotSend)
}

func GetAllNftPresaleStateSentPending() (*[]models.NftPresale, error) {
	return btldb.ReadNftPresalesByNftPresaleState(models.NftPresaleStateSentPending)
}

// @dev: Operations

// IncreaseNftPresaleProcessNumber
// @Description: Increase NftPresale process number
func IncreaseNftPresaleProcessNumber(nftPresale *models.NftPresale) (err error) {
	nftPresale.ProcessNumber += 1
	return UpdateNftPresale(nftPresale)
}

func StorePaidIdThenChangeStateAndClearProcessNumber(paidId int, nftPresale *models.NftPresale) error {
	nftPresale.PayMethod = models.FeePaymentMethodCustodyAccount
	nftPresale.PaidId = paidId
	nftPresale.State = models.NftPresaleStatePaidPending
	nftPresale.ProcessNumber = 0
	err := UpdateNftPresale(nftPresale)
	if err != nil {
		return utils.AppendErrorInfo(err, "UpdateNftPresale")
	}
	return nil
}

func SetNftPresaleFail(nftPresale *models.NftPresale) error {
	nftPresale.State = models.NftPresaleStateFailOrCanceled
	err := UpdateNftPresale(nftPresale)
	if err != nil {
		return utils.AppendErrorInfo(err, "UpdateNftPresale")
	}
	return nil
}

func ChangeNftPresaleStateAndUpdatePaidSuccessTimeThenClearProcessNumber(nftPresale *models.NftPresale) error {
	nftPresale.State = models.NftPresaleStatePaidNotSend
	nftPresale.PaidSuccessTime = utils.GetTimestamp()
	nftPresale.ProcessNumber = 0
	err := UpdateNftPresale(nftPresale)
	if err != nil {
		return utils.AppendErrorInfo(err, "UpdateNftPresale")
	}
	return nil
}

func ChangeNftPresaleStateAndClearProcessNumber(state models.NftPresaleState, nftPresale *models.NftPresale) error {
	nftPresale.State = state
	nftPresale.ProcessNumber = 0
	err := UpdateNftPresale(nftPresale)
	if err != nil {
		return utils.AppendErrorInfo(err, "UpdateNftPresale")
	}
	return nil
}

func UpdateNftPresaleAfterSent(outpoint string, txid string, address string, nftPresale *models.NftPresale) error {
	nftPresale.SentOutpoint = outpoint
	nftPresale.SentTxid = txid
	nftPresale.SentAddress = address
	nftPresale.State = models.NftPresaleStateSentPending
	nftPresale.ProcessNumber = 0
	nftPresale.SentTime = utils.GetTimestamp()
	err := UpdateNftPresale(nftPresale)
	if err != nil {
		return utils.AppendErrorInfo(err, "UpdateNftPresale")
	}
	return nil
}

func UpdateNftPresaleBySendAssetResponse(nftPresale *models.NftPresale, sendAssetResponse *taprpc.SendAssetResponse) error {
	var err error
	scriptKey := nftPresale.AddrScriptKey
	internalKey := nftPresale.AddrInternalKey
	var outpoint string
	var txid string
	var address string
	outpoint, err = SendAssetResponseScriptKeyAndInternalKeyToOutpoint(sendAssetResponse, scriptKey, internalKey)
	if err != nil {
		btlLog.PreSale.Error("SendAssetResponseScriptKeyAndInternalKeyToOutpoint:%v", err)
	} else {
		txid, _ = utils.GetTransactionAndIndexByOutpoint(outpoint)
		address, err = api.GetListChainTransactionsOutpointAddress(outpoint)
		if err != nil {
			btlLog.PreSale.Error("GetListChainTransactionsOutpointAddress:%v", err)
		}
	}
	err = UpdateNftPresaleAfterSent(outpoint, txid, address, nftPresale)
	if err != nil {
		btlLog.PreSale.Error("UpdateNftPresaleAfterSent:%v", err)
	}
	return nil
}

func UpdateNftPresaleAfterConfirmed(nftPresale *models.NftPresale) error {
	nftPresale.State = models.NftPresaleStateSent
	nftPresale.ProcessNumber = 0
	err := UpdateNftPresale(nftPresale)
	if err != nil {
		return utils.AppendErrorInfo(err, "UpdateNftPresale")
	}
	return nil
}

// @dev: Process

func ProcessNftPresaleStateBoughtNotPayService(nftPresale *models.NftPresale) error {
	// @dev: 1. Pay fee
	paidId, err := PayGasFee(nftPresale.BuyerUserId, nftPresale.Price)
	if err != nil {
		return utils.AppendErrorInfo(err, "PayGasFee for nftPresale")
	}
	// @dev: 2. Store paidId; Change state; Clear ProcessNumber
	err = StorePaidIdThenChangeStateAndClearProcessNumber(paidId, nftPresale)
	if err != nil {
		return utils.AppendErrorInfo(err, "StorePaidIdThenChangeStateAndClearProcessNumber")
	}
	return nil
}

func ProcessNftPresaleStatePaidPendingService(nftPresale *models.NftPresale) error {
	// @dev: 1. Is fee paid
	var isFeePaid bool
	isFeePaid, err := IsFeePaid(nftPresale.PaidId)
	if err != nil {
		if errors.Is(err, models.CustodyAccountPayInsideMissionFaild) {
			err = SetNftPresaleFail(nftPresale)
			if err != nil {
				return utils.AppendErrorInfo(err, "SetNftPresaleFail")
			}
		}
	}
	// @dev: Fee has not been paid
	if isFeePaid {
		// @dev: Change state; clear Process Number
		err = ChangeNftPresaleStateAndUpdatePaidSuccessTimeThenClearProcessNumber(nftPresale)
		if err != nil {
			return utils.AppendErrorInfo(err, "ChangeNftPresaleStateAndUpdatePaidSuccessTimeThenClearProcessNumber")
		}
		return nil
	}
	return nil
}

func ProcessNftPresaleStatePaidNotSendService(nftPresale *models.NftPresale) error {
	var err error
	// @dev: Check if confirmed balance enough
	minBalance := 2000
	if !IsWalletBalanceEnough(minBalance) {
		err = errors.New("lnd wallet balance is not enough(less than " + strconv.Itoa(minBalance) + ")")
		return err
	}
	// @dev: Check if asset balance enough
	if !IsAssetBalanceEnough(nftPresale.AssetId, nftPresale.Amount) {
		err = errors.New("nft presale asset(" + strconv.Itoa(int(nftPresale.ID)) + ") balance is not enough")
		return utils.AppendErrorInfo(err, "IsAssetBalanceEnough")
	}
	// @dev: Check if asset utxo is enough
	if !IsAssetUtxoEnough(nftPresale.AssetId, nftPresale.Amount) {
		err = errors.New("nft presale asset(" + strconv.Itoa(int(nftPresale.ID)) + ") utxo is not enough")
		return utils.AppendErrorInfo(err, "IsAssetUtxoEnough")
	}
	if nftPresale.ReceiveAddr == "" {
		err = errors.New("nft presale asset(" + strconv.Itoa(int(nftPresale.ID)) + ")'s receive addr(" + nftPresale.ReceiveAddr + ") is null")
		return err
	}
	// @dev: Send Asset
	addrs := []string{nftPresale.ReceiveAddr}
	// @dev: Get fee rate
	feeRate, err := UpdateAndGetFeeRateResponseTransformed()
	if err != nil {
		return utils.AppendErrorInfo(err, "UpdateAndGetFeeRateResponseTransformed("+strconv.Itoa(int(nftPresale.ID))+")")
	}
	feeRateSatPerKw := feeRate.SatPerKw.FastestFee
	// @dev: Send and get response
	response, err := api.SendAssetAddrSliceAndGetResponse(addrs, feeRateSatPerKw)
	if err != nil {
		return utils.AppendErrorInfo(err, "SendAssetAddrSliceAndGetResponse("+strconv.Itoa(int(nftPresale.ID))+")")
	}
	// @dev: Update info
	err = UpdateNftPresaleBySendAssetResponse(nftPresale, response)
	if err != nil {
		return utils.AppendErrorInfo(err, "UpdateNftPresaleBySendAssetResponse("+strconv.Itoa(int(nftPresale.ID))+")")
	}
	return nil
}

func ProcessNftPresaleStateSentPendingService(nftPresale *models.NftPresale) error {
	var err error
	if nftPresale.SentOutpoint == "" {
		err = errors.New("no outpoint generated, asset of presale(" + strconv.Itoa(int(nftPresale.ID)) + ") may has not been sent")
		return err
	}
	// @dev: 1.Is Transaction Confirmed
	if IsTransactionConfirmed(nftPresale.SentTxid) {
		// @dev: Change state and Clear ProcessNumber
		err = UpdateNftPresaleAfterConfirmed(nftPresale)
		if err != nil {
			return utils.AppendErrorInfo(err, "UpdateNftPresaleAfterConfirmed")
		}
		return nil
	} else {
		err = errors.New("nftPresale.SentTxid(" + nftPresale.SentTxid + ") is not confirmed")
		return err
	}
}

// @dev: Process all

func ProcessAllNftPresaleStateBoughtNotPayService() (*[]ProcessionResult, error) {
	var processionResults []ProcessionResult
	nftPresales, err := GetAllNftPresaleStateBoughtNotPay()
	if err != nil {
		return nil, utils.AppendErrorInfo(err, "GetAllNftPresaleStateBoughtNotPay")
	}
	for _, nftPresale := range *nftPresales {
		{
			err = IncreaseNftPresaleProcessNumber(&nftPresale)
			if err != nil {
				// @dev: Do nothing
			}
			err = ProcessNftPresaleStateBoughtNotPayService(&nftPresale)
			if err != nil {
				processionResults = append(processionResults, ProcessionResult{
					id: int(nftPresale.ID),
					JsonResult: models.JsonResult{
						Success: false,
						Error:   err.Error(),
						Data:    nil,
					},
				})
				continue
			} else {
				processionResults = append(processionResults, ProcessionResult{
					id: int(nftPresale.ID),
					JsonResult: models.JsonResult{
						Success: true,
						Error:   "",
						Data:    nil,
					},
				})
			}
		}
	}
	if processionResults == nil || len(processionResults) == 0 {
		err = errors.New("procession results null")
		return nil, err
	}
	return &processionResults, nil
}

func ProcessAllNftPresaleStatePaidPendingService() (*[]ProcessionResult, error) {
	var processionResults []ProcessionResult
	nftPresales, err := GetAllNftPresaleStatePaidPending()
	if err != nil {
		return nil, utils.AppendErrorInfo(err, "GetAllNftPresaleStatePaidPending")
	}
	for _, nftPresale := range *nftPresales {
		{
			err = IncreaseNftPresaleProcessNumber(&nftPresale)
			if err != nil {
				// @dev: Do nothing
			}
			err = ProcessNftPresaleStatePaidPendingService(&nftPresale)
			if err != nil {
				processionResults = append(processionResults, ProcessionResult{
					id: int(nftPresale.ID),
					JsonResult: models.JsonResult{
						Success: false,
						Error:   err.Error(),
						Data:    nil,
					},
				})
				continue
			} else {
				processionResults = append(processionResults, ProcessionResult{
					id: int(nftPresale.ID),
					JsonResult: models.JsonResult{
						Success: true,
						Error:   "",
						Data:    nil,
					},
				})
			}
		}
	}
	if processionResults == nil || len(processionResults) == 0 {
		err = errors.New("procession results null")
		return nil, err
	}
	return &processionResults, nil
}

func ProcessAllNftPresaleStatePaidNotSendService() (*[]ProcessionResult, error) {
	var processionResults []ProcessionResult
	nftPresales, err := GetAllNftPresaleStatePaidNotSend()
	if err != nil {
		return nil, utils.AppendErrorInfo(err, "GetAllNftPresaleStatePaidNotSend")
	}
	for _, nftPresale := range *nftPresales {
		{
			err = IncreaseNftPresaleProcessNumber(&nftPresale)
			if err != nil {
				// @dev: Do nothing
			}
			err = ProcessNftPresaleStatePaidNotSendService(&nftPresale)
			if err != nil {
				processionResults = append(processionResults, ProcessionResult{
					id: int(nftPresale.ID),
					JsonResult: models.JsonResult{
						Success: false,
						Error:   err.Error(),
						Data:    nil,
					},
				})
				continue
			} else {
				processionResults = append(processionResults, ProcessionResult{
					id: int(nftPresale.ID),
					JsonResult: models.JsonResult{
						Success: true,
						Error:   "",
						Data:    nil,
					},
				})
			}
		}

	}
	if processionResults == nil || len(processionResults) == 0 {
		err = errors.New("procession results null")
		return nil, err
	}
	return &processionResults, nil
}

func ProcessAllNftPresaleStateSentPendingService() (*[]ProcessionResult, error) {
	var processionResults []ProcessionResult
	nftPresales, err := GetAllNftPresaleStateSentPending()
	if err != nil {
		return nil, utils.AppendErrorInfo(err, "GetAllNftPresaleStateSentPending")
	}
	for _, nftPresale := range *nftPresales {
		{
			err = IncreaseNftPresaleProcessNumber(&nftPresale)
			if err != nil {
				// @dev: Do nothing
			}
			err = ProcessNftPresaleStateSentPendingService(&nftPresale)
			if err != nil {
				processionResults = append(processionResults, ProcessionResult{
					id: int(nftPresale.ID),
					JsonResult: models.JsonResult{
						Success: false,
						Error:   err.Error(),
						Data:    nil,
					},
				})
				continue
			} else {
				processionResults = append(processionResults, ProcessionResult{
					id: int(nftPresale.ID),
					JsonResult: models.JsonResult{
						Success: true,
						Error:   "",
						Data:    nil,
					},
				})
			}
		}

	}
	if processionResults == nil || len(processionResults) == 0 {
		err = errors.New("procession results null")
		return nil, err
	}
	return &processionResults, nil
}

// @dev: Scheduled task

func ProcessNftPresaleBoughtNotPay() {
	processionResult, err := ProcessAllNftPresaleStateBoughtNotPayService()
	if err != nil {
		return
	}
	if processionResult == nil || len(*processionResult) == 0 {
		return
	}
	// @dev: Do not use PrintProcessionResult
	err = utils.WriteToLogFile("./trade.presale.log", "[PRESALE.BNP]", "\n"+utils.ValueJsonString(processionResult))
	if err != nil {
		utils.LogError("WriteToLogFile ./trade.presale.log", err)
	}
}

func ProcessNftPresalePaidPending() {
	processionResult, err := ProcessAllNftPresaleStatePaidPendingService()
	if err != nil {
		return
	}
	if processionResult == nil || len(*processionResult) == 0 {
		return
	}
	// @dev: Do not use PrintProcessionResult
	err = utils.WriteToLogFile("./trade.presale.log", "[PRESALE.PPD]", "\n"+utils.ValueJsonString(processionResult))
	if err != nil {
		utils.LogError("WriteToLogFile ./trade.presale.log", err)
	}
}

func ProcessNftPresalePaidNotSend() {
	processionResult, err := ProcessAllNftPresaleStatePaidNotSendService()
	if err != nil {
		return
	}
	if processionResult == nil || len(*processionResult) == 0 {
		return
	}
	// @dev: Do not use PrintProcessionResult
	err = utils.WriteToLogFile("./trade.presale.log", "[PRESALE.PNS]", "\n"+utils.ValueJsonString(processionResult))
	if err != nil {
		utils.LogError("WriteToLogFile ./trade.presale.log", err)
	}
}

func ProcessNftPresaleSentPending() {
	processionResult, err := ProcessAllNftPresaleStateSentPendingService()
	if err != nil {
		return
	}
	if processionResult == nil || len(*processionResult) == 0 {
		return
	}
	// @dev: Do not use PrintProcessionResult
	err = utils.WriteToLogFile("./trade.presale.log", "[PRESALE.SPD]", "\n"+utils.ValueJsonString(processionResult))
	if err != nil {
		utils.LogError("WriteToLogFile ./trade.presale.log", err)
	}
}
