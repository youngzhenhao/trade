package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"net/http"
	"strconv"
	"trade/models"
	"trade/services"
	"trade/utils"
)

func GetAllFairLaunchInfo(c *gin.Context) {
	allFairLaunch, err := services.GetAllFairLaunchInfos()
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   "Can not get all fair launch infos. " + err.Error(),
			Code:    models.GetAllFairLaunchInfosErr,
			Data:    nil,
		})
		return
	}
	c.JSON(http.StatusOK, models.JsonResult{
		Success: true,
		Error:   models.SuccessErr,
		Code:    models.SUCCESS,
		Data:    allFairLaunch,
	})
}

func GetFairLaunchInfo(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   "id is not valid int. " + err.Error(),
			Code:    models.FairLaunchInfoIdInvalidErr,
			Data:    nil,
		})
		return
	}
	fairLaunch, err := services.GetFairLaunchInfo(id)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   "Can not get fair launch info. " + err.Error(),
			Code:    models.GetFairLaunchInfoErr,
			Data:    nil,
		})
		return
	}
	c.JSON(http.StatusOK, models.JsonResult{
		Success: true,
		Error:   models.SuccessErr,
		Code:    models.SUCCESS,
		Data:    fairLaunch,
	})
}

func GetMintedInfo(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   "id is not valid int. " + err.Error(),
			Code:    models.FairLaunchMintedInfoIdInvalidErr,
			Data:    nil,
		})
		return
	}
	minted, err := services.GetFairLaunchMintedInfosByFairLaunchId(id)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   "Can not get fair launch minted info. " + err.Error(),
			Code:    models.GetFairLaunchMintedInfosByFairLaunchIdErr,
			Data:    nil,
		})
		return
	}
	c.JSON(http.StatusOK, models.JsonResult{
		Success: true,
		Error:   models.SuccessErr,
		Code:    models.SUCCESS,
		Data:    minted,
	})
}

func SetFairLaunchInfo(c *gin.Context) {
	var fairLaunchInfo *models.FairLaunchInfo
	// @dev: Use MustGet. alice ONLY FOR TEST
	username := c.MustGet("username").(string)
	userId, err := services.NameToId(username)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   "Query user id by name." + err.Error(),
			Code:    models.NameToIdErr,
			Data:    nil,
		})
		return
	}
	// @dev: Use SetFairLaunchInfoRequest c.ShouldBind
	var setFairLaunchInfoRequest models.SetFairLaunchInfoRequest
	err = c.ShouldBindJSON(&setFairLaunchInfoRequest)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   "Should Bind JSON setFairLaunchInfoRequest. " + err.Error(),
			Code:    models.ShouldBindJsonErr,
			Data:    nil,
		})
		return
	}
	imageData := setFairLaunchInfoRequest.ImageData
	name := setFairLaunchInfoRequest.Name
	assetType := setFairLaunchInfoRequest.AssetType
	amount := setFairLaunchInfoRequest.Amount
	reserved := setFairLaunchInfoRequest.Reserved
	mintQuantity := setFairLaunchInfoRequest.MintQuantity
	startTime := setFairLaunchInfoRequest.StartTime
	endTime := setFairLaunchInfoRequest.EndTime
	description := setFairLaunchInfoRequest.Description
	feeRate := setFairLaunchInfoRequest.FeeRate
	// @dev: Process struct, update later
	// @notice: State is 0 now
	fairLaunchInfo, err = services.ProcessFairLaunchInfo(imageData, name, assetType, amount, reserved, mintQuantity, startTime, endTime, description, feeRate, userId, username)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   "Process fair launch info. " + err.Error(),
			Code:    models.ProcessFairLaunchInfoErr,
			Data:    nil,
		})
		return
	}
	// @dev: Update db, State models.FairLaunchStateNoPay
	err = services.SetFairLaunchInfo(fairLaunchInfo)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   "Set fair launch error. " + err.Error(),
			Code:    models.SetFairLaunchInfoErr,
			Data:    nil,
		})
		return
	}
	c.JSON(http.StatusOK, models.JsonResult{
		Success: true,
		Error:   models.SuccessErr,
		Code:    models.SUCCESS,
		Data:    nil,
	})
}

func SetFairLaunchMintedInfo(c *gin.Context) {
	var fairLaunchMintedInfo *models.FairLaunchMintedInfo
	var mintFairLaunchRequest models.MintFairLaunchRequest
	// @notice: only receive id and number
	err := c.ShouldBindJSON(&mintFairLaunchRequest)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   "Should Bind JSON mintFairLaunchRequest. " + err.Error(),
			Code:    models.ShouldBindJsonErr,
			Data:    nil,
		})
		return
	}
	// @dev: Ensure time is valid
	isTimeRight, err := services.IsFairLaunchMintTimeRight(mintFairLaunchRequest.FairLaunchInfoID)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   "Is FairLaunch Mint Time Right. " + err.Error(),
			Code:    models.IsFairLaunchMintTimeRightErr,
			Data:    nil,
		})
		return
	}
	if !isTimeRight {
		err = errors.New("It is not Right FairLaunch Mint Time now")
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   err.Error(),
			Code:    models.IsTimeRightErr,
			Data:    nil,
		})
		return
	}
	fairLaunchInfoID := mintFairLaunchRequest.FairLaunchInfoID
	isFairLaunchIssued := services.IsFairLaunchIssued(fairLaunchInfoID)
	if !isFairLaunchIssued {
		err = errors.New("FairLaunch is not Issued.")
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   err.Error(),
			Code:    models.IsFairLaunchIssued,
			Data:    nil,
		})
		return
	}
	// @dev: Use MustGet. bob ONLY FOR TEST
	username := c.MustGet("username").(string)
	// @dev: userId
	userId, err := services.NameToId(username)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   "Name to id error. " + err.Error(),
			Code:    models.NameToIdErr,
			Data:    nil,
		})
		return
	}
	mintedNumber := mintFairLaunchRequest.MintedNumber
	addr := mintFairLaunchRequest.EncodedAddr
	mintedFeeRateSatPerKw := mintFairLaunchRequest.MintedFeeRateSatPerKw
	fairLaunchMintedInfo, err = services.ProcessFairLaunchMintedInfo(fairLaunchInfoID, mintedNumber, mintedFeeRateSatPerKw, addr, userId, username)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   "Process FairLaunchMintedInfo " + err.Error(),
			Code:    models.ProcessFairLaunchMintedInfoErr,
			Data:    nil,
		})
		return
	}
	// @dev: Update db, State models.FairLaunchMintedStateNoPay
	err = services.SetFairLaunchMintedInfo(fairLaunchMintedInfo)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   "Set fair launch minted info. " + err.Error(),
			Code:    models.SetFairLaunchMintedInfoErr,
			Data:    nil,
		})
		return
	}
	c.JSON(http.StatusOK, models.JsonResult{
		Success: true,
		Error:   models.SuccessErr,
		Code:    models.SUCCESS,
		Data:    nil,
	})
}

func QueryInventory(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   "strconv string to int." + err.Error(),
			Code:    models.NameToIdErr,
			Data:    nil,
		})
		return
	}
	inventory, err := services.GetInventoryCouldBeMintedByFairLaunchInfoId(id)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   "Get inventory could be minted by fair launch info id." + err.Error(),
			Code:    models.GetInventoryCouldBeMintedByFairLaunchInfoIdErr,
			Data:    nil,
		})
		return
	}
	c.JSON(http.StatusOK, models.JsonResult{
		Success: true,
		Error:   models.SuccessErr,
		Code:    models.SUCCESS,
		Data:    inventory,
	})
}

func QueryMintIsAvailable(c *gin.Context) {
	var mintFairLaunchRequest models.MintFairLaunchRequest
	err := c.ShouldBindJSON(&mintFairLaunchRequest)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   "Should Bind JSON mintFairLaunchRequest. " + err.Error(),
			Code:    models.ShouldBindJsonErr,
			Data:    nil,
		})
		return
	}
	fairLaunchInfoID := mintFairLaunchRequest.FairLaunchInfoID
	mintedNumber := mintFairLaunchRequest.MintedNumber
	// @dev: calculated FeeRate
	feeRate, err := services.UpdateAndCalculateGasFeeRateByMempool(mintedNumber)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   "Calculate Gas FeeRate. " + err.Error(),
			Code:    models.UpdateAndCalculateGasFeeRateByMempoolErr,
			Data:    nil,
		})
		return
	}
	calculatedFeeRateSatPerKw := feeRate.SatPerKw.FastestFee + services.FeeRateSatPerBToSatPerKw(2)
	calculatedFeeRateSatPerB := feeRate.SatPerB.FastestFee + 2
	calculatedFee := services.GetMintedTransactionGasFee(calculatedFeeRateSatPerKw)
	mintedAmount, err := services.GetAmountCouldBeMintByMintedNumber(fairLaunchInfoID, mintedNumber)
	if err != nil {
		// @dev: Do not return
	}
	isMintAvailable := mintedAmount > 0
	numberAndAmountCouldBeMint, err := services.GetNumberAndAmountCouldBeMint(fairLaunchInfoID)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   "Get Number And Amount Of Inventory Could Be Minted. " + err.Error(),
			Code:    models.GetNumberAndAmountOfInventoryCouldBeMintedErr,
			Data:    nil,
		})
		return
	}
	c.JSON(http.StatusOK, models.JsonResult{
		Success: true,
		Error:   models.SUCCESS.Error(),
		Code:    models.SUCCESS,
		Data: gin.H{
			"is_mint_available":              isMintAvailable,
			"inventory_amount":               mintedAmount,
			"calculated_fee_rate_sat_per_kw": calculatedFeeRateSatPerKw,
			"calculated_fee_rate_sat_per_b":  calculatedFeeRateSatPerB,
			"calculated_fee":                 calculatedFee,
			"available_number":               numberAndAmountCouldBeMint.Number,
		},
	})
}

func MintFairLaunchReserved(c *gin.Context) {
	var mintFairLaunchReservedRequest models.MintFairLaunchReservedRequest
	err := c.ShouldBindJSON(&mintFairLaunchReservedRequest)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   "Should Bind JSON mintFairLaunchReservedRequest. " + err.Error(),
			Code:    models.ShouldBindJsonErr,
			Data:    nil,
		})
		return
	}
	username := c.MustGet("username").(string)
	userId, err := services.NameToId(username)
	assetId := mintFairLaunchReservedRequest.AssetID
	addr := mintFairLaunchReservedRequest.EncodedAddr
	fairLaunchInfo, err := services.GetFairLaunchInfoByAssetId(assetId)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   "Get FairLaunchInfo By AssetId. " + err.Error(),
			Code:    models.GetFairLaunchInfoByAssetIdErr,
			Data:    nil,
		})
		return
	}
	id := int(fairLaunchInfo.ID)
	isTimeRight, err := services.IsFairLaunchMintTimeRight(id)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   "Is FairLaunch Mint Time Right. " + err.Error(),
			Code:    models.IsFairLaunchMintTimeRightErr,
			Data:    nil,
		})
		return
	}
	if !isTimeRight {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   "It is not Valid Time Now. ",
			Code:    models.IsTimeRightErr,
			Data:    nil,
		})
		return
	}
	fairLaunch, err := services.GetFairLaunchInfo(id)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   "Can not get fair launch info. " + err.Error(),
			Code:    models.GetFairLaunchInfoErr,
			Data:    nil,
		})
		return
	}
	if userId != fairLaunch.UserID {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   "Invalid user id. ",
			Code:    models.InvalidUserIdErr,
			Data:    nil,
		})
		return
	}
	response, err := services.SendFairLaunchReserved(fairLaunch, addr)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   "Send FairLaunch Reserved. " + err.Error(),
			Code:    models.SendFairLaunchReservedErr,
			Data:    nil,
		})
		return
	}
	outpoint := services.ProcessSendFairLaunchReservedResponse(response)
	err = services.UpdateFairLaunchInfoIsReservedSent(fairLaunch, outpoint)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   "Update FairLaunchInfo IsReservedSent. " + err.Error(),
			Code:    models.UpdateFairLaunchInfoIsReservedSentErr,
			Data:    nil,
		})
		return
	}
	// @dev: Record paid fee
	txid, _ := utils.OutpointToTransactionAndIndex(outpoint)
	err = services.CreateFairLaunchIncomeOfServerPaySendReservedFee(fairLaunchInfo.AssetID, int(fairLaunchInfo.ID), txid)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   "Create FairLaunch Income Of Server Pay Send Reserved Fee. " + err.Error(),
			Code:    models.CreateFairLaunchIncomeOfServerPaySendReservedFeeErr,
			Data:    nil,
		})
		return
	}
	c.JSON(http.StatusOK, models.JsonResult{
		Success: true,
		Error:   models.SuccessErr,
		Code:    models.SUCCESS,
		Data: gin.H{
			"anchor_outpoint": outpoint,
		},
	})
}

func GetIssuedFairLaunchInfo(c *gin.Context) {
	var fairLaunchInfos *[]models.FairLaunchInfo
	var err error
	fairLaunchInfos, err = services.GetIssuedFairLaunchInfos()
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   "Get Issued FairLaunchInfos. " + err.Error(),
			Code:    models.GetIssuedFairLaunchInfosErr,
			Data:    nil,
		})
		return
	}
	fairLaunchInfos = services.ProcessIssuedFairLaunchInfos(fairLaunchInfos)
	fairLaunchInfos = services.GetSortedFairLaunchInfosByMintedRate(fairLaunchInfos)
	c.JSON(http.StatusOK, models.JsonResult{
		Success: true,
		Error:   models.SuccessErr,
		Code:    models.SUCCESS,
		Data:    fairLaunchInfos,
	})
}

func GetOwnFairLaunchInfo(c *gin.Context) {
	var fairLaunchInfos *[]models.FairLaunchInfo
	var err error
	username := c.MustGet("username").(string)
	userId, err := services.NameToId(username)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   "Query user id by name." + err.Error(),
			Code:    models.NameToIdErr,
			Data:    nil,
		})
		return
	}
	fairLaunchInfos, err = services.GetOwnFairLaunchInfosByUserId(userId)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   "Get Own Set FairLaunchInfos By UserId. " + err.Error(),
			Code:    models.GetOwnFairLaunchInfosByUserIdErr,
			Data:    nil,
		})
		return
	}
	c.JSON(http.StatusOK, models.JsonResult{
		Success: true,
		Error:   models.SuccessErr,
		Code:    models.SUCCESS,
		Data:    fairLaunchInfos,
	})
}

func GetOwnFairLaunchMintedInfo(c *gin.Context) {
	var fairLaunchMintedInfos *[]models.FairLaunchMintedInfo
	var err error
	username := c.MustGet("username").(string)
	userId, err := services.NameToId(username)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   "Query user id by name." + err.Error(),
			Code:    models.NameToIdErr,
			Data:    nil,
		})
		return
	}
	fairLaunchMintedInfos, err = services.GetOwnFairLaunchMintedInfosByUserId(userId)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   "Get Own FairLaunchMintedInfos By UserId. " + err.Error(),
			Code:    models.GetOwnFairLaunchMintedInfosByUserIdErr,
			Data:    nil,
		})
		return
	}
	c.JSON(http.StatusOK, models.JsonResult{
		Success: true,
		Error:   models.SuccessErr,
		Code:    models.SUCCESS,
		Data:    fairLaunchMintedInfos,
	})
}

func GetFairLaunchInfoByAssetId(c *gin.Context) {
	assetId := c.Param("asset_id")
	fairLaunch, err := services.GetFairLaunchInfoByAssetId(assetId)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   "Can not get fair launch info By AssetId. " + err.Error(),
			Code:    models.GetFairLaunchInfoByAssetIdErr,
			Data:    nil,
		})
		return
	}
	c.JSON(http.StatusOK, models.JsonResult{
		Success: true,
		Error:   models.SuccessErr,
		Code:    models.SUCCESS,
		Data:    fairLaunch,
	})
}

func GetFairLaunchInventoryMintNumberAssetId(c *gin.Context) {
	assetId := c.Param("asset_id")
	fairLaunch, err := services.GetFairLaunchInfoByAssetId(assetId)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   "Can not get fair launch info By AssetId. " + err.Error(),
			Code:    models.GetFairLaunchInfoByAssetIdErr,
			Data:    nil,
		})
		return
	}
	inventoryNumberAndAmount, err := services.GetNumberAndAmountOfInventoryCouldBeMinted(int(fairLaunch.ID))
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   "Get Number And Amount Of Inventory Could Be Minted. " + err.Error(),
			Code:    models.GetNumberAndAmountOfInventoryCouldBeMintedErr,
			Data:    nil,
		})
		return
	}
	c.JSON(http.StatusOK, models.JsonResult{
		Success: true,
		Error:   models.SuccessErr,
		Code:    models.SUCCESS,
		Data:    inventoryNumberAndAmount,
	})
}

func GetOwnFairLaunchInfoIssuedSimplified(c *gin.Context) {
	var fairLaunchInfos *[]services.FairLaunchInfoSimplified
	var err error
	username := c.MustGet("username").(string)
	userId, err := services.NameToId(username)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   "Query user id by name." + err.Error(),
			Code:    models.NameToIdErr,
			Data:    nil,
		})
		return
	}
	fairLaunchInfos, err = services.GetFairLaunchInfoSimplifiedByUserIdIssued(userId)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   "Get Own Set FairLaunchInfos By UserId. " + err.Error(),
			Code:    models.GetFairLaunchInfoSimplifiedByUserIdIssuedErr,
			Data:    nil,
		})
		return
	}
	c.JSON(http.StatusOK, models.JsonResult{
		Success: true,
		Error:   models.SuccessErr,
		Code:    models.SUCCESS,
		Data:    fairLaunchInfos,
	})
}

func GetClosedFairLaunchInfo(c *gin.Context) {
	var fairLaunchInfos *[]models.FairLaunchInfo
	var err error
	fairLaunchInfos, err = services.GetClosedFairLaunchInfo()
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   "Get Closed FairLaunchInfos. " + err.Error(),
			Code:    models.GetClosedFairLaunchInfoErr,
			Data:    nil,
		})
		return
	}
	c.JSON(http.StatusOK, models.JsonResult{
		Success: true,
		Error:   models.SuccessErr,
		Code:    models.SUCCESS,
		Data:    fairLaunchInfos,
	})
}

func GetNotStartedFairLaunchInfo(c *gin.Context) {
	var fairLaunchInfos *[]models.FairLaunchInfo
	var err error
	fairLaunchInfos, err = services.GetNotStartedFairLaunchInfo()
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   "Get Not Started FairLaunchInfos. " + err.Error(),
			Code:    models.GetNotStartedFairLaunchInfoErr,
			Data:    nil,
		})
		return
	}
	c.JSON(http.StatusOK, models.JsonResult{
		Success: true,
		Error:   models.SuccessErr,
		Code:    models.SUCCESS,
		Data:    fairLaunchInfos,
	})
}

func GetHotFairLaunchInfo(c *gin.Context) {
	var fairLaunchInfos *[]models.FairLaunchInfo
	var err error
	fairLaunchInfos, err = services.GetIssuedFairLaunchInfos()
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   "Get Issued FairLaunchInfos. " + err.Error(),
			Code:    models.GetAllFairLaunchInfosErr,
			Data:    nil,
		})
		return
	}
	fairLaunchInfos = services.GetSortedFairLaunchInfosByMintedRate(fairLaunchInfos)
	c.JSON(http.StatusOK, models.JsonResult{
		Success: true,
		Error:   models.SuccessErr,
		Code:    models.SUCCESS,
		Data:    fairLaunchInfos,
	})
}

func GetFollowedFairLaunchInfo(c *gin.Context) {
	username := c.MustGet("username").(string)
	userId, err := services.NameToId(username)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   err.Error(),
			Code:    models.NameToIdErr,
			Data:    nil,
		})
		return
	}
	fairLaunchInfos, err := services.GetFollowedFairLaunchInfo(userId)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   err.Error(),
			Code:    models.GetFollowedFairLaunchInfoErr,
			Data:    nil,
		})
		return
	}
	c.JSON(http.StatusOK, models.JsonResult{
		Success: true,
		Error:   models.SuccessErr,
		Code:    models.SUCCESS,
		Data:    fairLaunchInfos,
	})
}

func GetFairLaunchInfoPlusByAssetId(c *gin.Context) {
	assetId := c.Param("asset_id")
	fairLaunch, err := services.GetFairLaunchInfoByAssetId(assetId)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   "Can not get fair launch info By AssetId. " + err.Error(),
			Code:    models.GetFairLaunchInfoByAssetIdErr,
			Data:    nil,
		})
		return
	}
	// @dev: Get holder number
	holderNumber, err := services.GetAssetHolderNumberAssetBalance(assetId)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   err.Error(),
			Code:    models.GetAssetHolderNumberAssetBalanceErr,
			Data:    nil,
		})
		return
	}
	fairLaunchPlusInfo := services.ProcessToFairLaunchPlusInfo(fairLaunch, holderNumber)
	c.JSON(http.StatusOK, models.JsonResult{
		Success: true,
		Error:   models.SuccessErr,
		Code:    models.SUCCESS,
		Data:    fairLaunchPlusInfo,
	})
}

func RefundUserFirstMintByUsernameAndAssetId(c *gin.Context) {
	var refundUserFirstMintRequest services.RefundUserFirstMintRequest
	err := c.BindJSON(&refundUserFirstMintRequest)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   err.Error(),
			Code:    models.ShouldBindJsonErr,
			Data:    nil,
		})
		return
	}
	refundResult, err := services.RefundUserFirstMintByUsernameAndAssetId(refundUserFirstMintRequest.Usernames, refundUserFirstMintRequest.AssetId)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   err.Error(),
			Code:    models.RefundUserFirstMintByUsernameAndAssetIdErr,
			Data:    nil,
		})
		return
	}
	c.JSON(http.StatusOK, models.JsonResult{
		Success: true,
		Error:   models.SuccessErr,
		Code:    models.SUCCESS,
		Data:    refundResult,
	})
}
