package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"net/http"
	"strconv"
	"trade/models"
	"trade/services"
)

func GetAllFairLaunchInfo(c *gin.Context) {
	allFairLaunch, err := services.GetAllFairLaunchInfos()
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   "Can not get all fair launch infos. " + err.Error(),
			Data:    nil,
		})
		return
	}
	c.JSON(http.StatusOK, models.JsonResult{
		Success: true,
		Error:   "",
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
			Data:    nil,
		})
		return
	}
	fairLaunch, err := services.GetFairLaunchInfo(id)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   "Can not get fair launch info. " + err.Error(),
			Data:    nil,
		})
		return
	}
	c.JSON(http.StatusOK, models.JsonResult{
		Success: true,
		Error:   "",
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
			Data:    nil,
		})
		return
	}
	minted, err := services.GetFairLaunchMintedInfosByFairLaunchId(id)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   "Can not get fair launch minted info. " + err.Error(),
			Data:    nil,
		})
		return
	}
	c.JSON(http.StatusOK, models.JsonResult{
		Success: true,
		Error:   "",
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
	fairLaunchInfo, err = services.ProcessFairLaunchInfo(imageData, name, assetType, amount, reserved, mintQuantity, startTime, endTime, description, feeRate, userId)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   "Process fair launch info. " + err.Error(),
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
			Data:    nil,
		})
		return
	}
	c.JSON(http.StatusOK, models.JsonResult{
		Success: true,
		Error:   "",
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
			Data:    nil,
		})
		return
	}
	if !isTimeRight {
		err = errors.New("It is not Right FairLaunch Mint Time now")
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   err.Error(),
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
			Data:    nil,
		})
		return
	}
	mintedNumber := mintFairLaunchRequest.MintedNumber
	addr := mintFairLaunchRequest.EncodedAddr
	mintedFeeRateSatPerKw := mintFairLaunchRequest.MintedFeeRateSatPerKw
	fairLaunchMintedInfo, err = services.ProcessFairLaunchMintedInfo(fairLaunchInfoID, mintedNumber, mintedFeeRateSatPerKw, addr, userId)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   "Process FairLaunchMintedInfo " + err.Error(),
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
			Data:    nil,
		})
		return
	}
	c.JSON(http.StatusOK, models.JsonResult{
		Success: true,
		Error:   "",
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
			Data:    nil,
		})
		return
	}
	inventory, err := services.GetInventoryCouldBeMintedByFairLaunchInfoId(id)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   "Get inventory could be minted by fair launch info id." + err.Error(),
			Data:    nil,
		})
		return
	}
	c.JSON(http.StatusOK, models.JsonResult{
		Success: true,
		Error:   "",
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
			Data:    nil,
		})
		return
	}
	calculatedFeeRateSatPerKw := feeRate.SatPerKw.FastestFee
	calculatedFeeRateSatPerB := feeRate.SatPerB.FastestFee
	inventoryAmount, err := services.GetAmountOfInventoryCouldBeMintedByMintedNumber(fairLaunchInfoID, mintedNumber)
	isMintAvailable := inventoryAmount > 0
	inventoryNumberAndAmount, err := services.GetNumberAndAmountOfInventoryCouldBeMinted(fairLaunchInfoID)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   "Get Number And Amount Of Inventory Could Be Minted. " + err.Error(),
			Data:    nil,
		})
		return
	}
	c.JSON(http.StatusOK, models.JsonResult{
		Success: true,
		Error:   "",
		Data: gin.H{
			"is_mint_available":              isMintAvailable,
			"inventory_amount":               inventoryAmount,
			"calculated_fee_rate_sat_per_kw": calculatedFeeRateSatPerKw,
			"calculated_fee_rate_sat_per_b":  calculatedFeeRateSatPerB,
			"available_number":               inventoryNumberAndAmount.Number,
		},
	})
}

func MintFairLaunchReserved(c *gin.Context) {
	var mintFairLaunchReservedRequest models.MintFairLaunchReservedRequest
	err := c.ShouldBindJSON(&mintFairLaunchReservedRequest)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   "Should Bind JSON mintFairLaunchRequest. " + err.Error(),
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
			Data:    nil,
		})
		return
	}
	if !isTimeRight {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   "It is not Valid Time Now. ",
			Data:    nil,
		})
		return
	}
	fairLaunch, err := services.GetFairLaunchInfo(id)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   "Can not get fair launch info. " + err.Error(),
			Data:    nil,
		})
		return
	}
	if userId != fairLaunch.UserID {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   "Invalid user id. ",
			Data:    nil,
		})
		return
	}
	response, err := services.SendFairLaunchReserved(fairLaunch, addr)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   "Send FairLaunch Reserved. " + err.Error(),
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
			Data:    nil,
		})
		return
	}
	c.JSON(http.StatusOK, models.JsonResult{
		Success: true,
		Error:   "",
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
			Data:    nil,
		})
		return
	}
	fairLaunchInfos = services.ProcessIssuedFairLaunchInfos(fairLaunchInfos)
	c.JSON(http.StatusOK, models.JsonResult{
		Success: true,
		Error:   "",
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
			Data:    nil,
		})
		return
	}
	fairLaunchInfos, err = services.GetOwnFairLaunchInfosByUserId(userId)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   "Get Own Set FairLaunchInfos By UserId. " + err.Error(),
			Data:    nil,
		})
		return
	}
	c.JSON(http.StatusOK, models.JsonResult{
		Success: true,
		Error:   "",
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
			Data:    nil,
		})
		return
	}
	fairLaunchMintedInfos, err = services.GetOwnFairLaunchMintedInfosByUserId(userId)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   "Get Own FairLaunchMintedInfos By UserId. " + err.Error(),
			Data:    nil,
		})
		return
	}
	c.JSON(http.StatusOK, models.JsonResult{
		Success: true,
		Error:   "",
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
			Data:    nil,
		})
		return
	}
	c.JSON(http.StatusOK, models.JsonResult{
		Success: true,
		Error:   "",
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
			Data:    nil,
		})
		return
	}
	inventoryNumberAndAmount, err := services.GetNumberAndAmountOfInventoryCouldBeMinted(int(fairLaunch.ID))
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   "Get Number And Amount Of Inventory Could Be Minted. " + err.Error(),
			Data:    nil,
		})
		return
	}
	c.JSON(http.StatusOK, models.JsonResult{
		Success: true,
		Error:   "",
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
			Data:    nil,
		})
		return
	}
	fairLaunchInfos, err = services.GetFairLaunchInfoSimplifiedByUserIdIssued(userId)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   "Get Own Set FairLaunchInfos By UserId. " + err.Error(),
			Data:    nil,
		})
		return
	}
	c.JSON(http.StatusOK, models.JsonResult{
		Success: true,
		Error:   "",
		Data:    fairLaunchInfos,
	})
}

func GetClosedFairLaunchInfo(c *gin.Context) {
	// TODO:
}

func GetNotStartedFairLaunchInfo(c *gin.Context) {
	// TODO:

}
