package handlers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"trade/api"
	"trade/models"
	"trade/services"
)

func GetAllIdoPublishInfo(c *gin.Context) {
	allIdoPublishInfos, err := services.GetAllIdoPublishInfos()
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   err.Error(),
			Code:    models.GetAllIdoPublishInfosErr,
			Data:    nil,
		})
		return
	}
	c.JSON(http.StatusOK, models.JsonResult{
		Success: true,
		Error:   "",
		Code:    models.SUCCESS,
		Data:    allIdoPublishInfos,
	})
}

func GetIdoPublishedInfo(c *gin.Context) {
	var idoPublishInfos *[]models.IdoPublishInfo
	var err error
	idoPublishInfos, err = services.GetIdoPublishedInfos()
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   err.Error(),
			Code:    models.GetIdoPublishedInfosErr,
			Data:    nil,
		})
		return
	}
	idoPublishInfos = services.ProcessIdoPublishedInfos(idoPublishInfos)
	c.JSON(http.StatusOK, models.JsonResult{
		Success: true,
		Error:   "",
		Code:    models.SUCCESS,
		Data:    idoPublishInfos,
	})
}

func GetOwnIdoPublishInfo(c *gin.Context) {
	var idoPublishInfos *[]models.IdoPublishInfo
	var err error
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
	idoPublishInfos, err = services.GetOwnIdoPublishInfosByUserId(userId)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   err.Error(),
			Code:    models.GetOwnIdoPublishInfoErr,
			Data:    nil,
		})
		return
	}
	c.JSON(http.StatusOK, models.JsonResult{
		Success: true,
		Error:   "",
		Code:    models.SUCCESS,
		Data:    idoPublishInfos,
	})
}

func GetOwnIdoParticipateInfo(c *gin.Context) {
	var idoParticipateInfos *[]models.IdoParticipateInfo
	var err error
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
	idoParticipateInfos, err = services.GetOwnIdoParticipateInfosByUserId(userId)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   err.Error(),
			Code:    models.GetOwnIdoParticipateInfoErr,
			Data:    nil,
		})
		return
	}
	c.JSON(http.StatusOK, models.JsonResult{
		Success: true,
		Error:   "",
		Code:    models.SUCCESS,
		Data:    idoParticipateInfos,
	})
}

func GetIdoPublishInfo(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   err.Error(),
			Code:    models.IdAtoiErr,
			Data:    nil,
		})
		return
	}
	idoPublishInfo, err := services.GetIdoPublishInfo(id)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   err.Error(),
			Code:    models.GetIdoParticipateInfoErr,
			Data:    nil,
		})
		return
	}
	c.JSON(http.StatusOK, models.JsonResult{
		Success: true,
		Error:   "",
		Code:    models.SUCCESS,
		Data:    idoPublishInfo,
	})
}

func GetIdoPublishInfoByAssetId(c *gin.Context) {
	assetId := c.Param("asset_id")
	idoPublishInfos, err := services.GetIdoPublishInfosByAssetId(assetId)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   err.Error(),
			Code:    models.GetIdoParticipateInfosByAssetIdErr,
			Data:    nil,
		})
		return
	}
	c.JSON(http.StatusOK, models.JsonResult{
		Success: true,
		Error:   "",
		Code:    models.SUCCESS,
		Data:    idoPublishInfos,
	})
}

func GetIdoParticipateInfo(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   err.Error(),
			Code:    models.IdAtoiErr,
			Data:    nil,
		})
		return
	}
	idoInfo, err := services.GetIdoParticipateInfo(id)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   err.Error(),
			Code:    models.GetIdoParticipateInfoErr,
			Data:    nil,
		})
		return
	}
	c.JSON(http.StatusOK, models.JsonResult{
		Success: true,
		Error:   "",
		Code:    models.SUCCESS,
		Data:    idoInfo,
	})
}

func QueryIdoParticipateIsAvailable(c *gin.Context) {
	//TODO
}

func SetIdoPublishInfo(c *gin.Context) {
	var idoPublishInfo *models.IdoPublishInfo
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
	var publishIdoRequest models.PublishIdoRequest
	err = c.ShouldBindJSON(&publishIdoRequest)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   err.Error(),
			Code:    models.ShouldBindJsonErr,
			Data:    nil,
		})
		return
	}
	assetID := publishIdoRequest.AssetID
	err = api.SyncAssetIssuance(assetID)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   err.Error(),
			Code:    models.SyncAssetIssuanceErr,
			Data:    nil,
		})
		return
	}
	totalAmount := publishIdoRequest.TotalAmount
	minimumQuantity := publishIdoRequest.MinimumQuantity
	unitPrice := publishIdoRequest.UnitPrice
	startTime := publishIdoRequest.StartTime
	endTime := publishIdoRequest.EndTime
	// @dev: SatPerKw
	feeRate := publishIdoRequest.FeeRate
	idoPublishInfo, err = services.ProcessIdoPublishInfo(userId, assetID, totalAmount, minimumQuantity, unitPrice, startTime, endTime, feeRate)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   err.Error(),
			Code:    models.ProcessIdoPublishInfoErr,
			Data:    nil,
		})
		return
	}
	err = services.SetIdoPublishInfo(idoPublishInfo)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   err.Error(),
			Code:    models.SetIdoPublishInfoErr,
			Data:    nil,
		})
		return
	}
	c.JSON(http.StatusOK, models.JsonResult{
		Success: true,
		Error:   "",
		Code:    models.SUCCESS,
		Data:    nil,
	})
}

func SetIdoParticipateInfo(c *gin.Context) {
	var idoParticipateInfo *models.IdoParticipateInfo
	var participateIdoRequest models.ParticipateIdoRequest
	err := c.ShouldBindJSON(&participateIdoRequest)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   err.Error(),
			Code:    models.ShouldBindJsonErr,
			Data:    nil,
		})
		return
	}
	idoPublishInfoID := participateIdoRequest.IdoPublishInfoID
	boughtAmount := participateIdoRequest.BoughtAmount
	// @dev: SatPerKw
	feeRate := participateIdoRequest.FeeRate
	encodedAddr := participateIdoRequest.EncodedAddr
	isTimeRight, err := services.IsIdoParticipateTimeRight(idoPublishInfoID)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   err.Error(),
			Code:    models.IsIdoParticipateTimeRightErr,
			Data:    nil,
		})
		return
	}
	if !isTimeRight {
		err = errors.New("it is not right participate time now")
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   err.Error(),
			Code:    models.IsNotRightTime,
			Data:    nil,
		})
		return
	}
	isIdoPublished := services.IsIdoPublished(idoPublishInfoID)
	if !isIdoPublished {
		err = errors.New("ido is not published")
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   err.Error(),
			Code:    models.IdoIsNotPublished,
			Data:    nil,
		})
		return
	}
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
	idoParticipateInfo, err = services.ProcessIdoParticipateInfo(userId, idoPublishInfoID, boughtAmount, feeRate, encodedAddr)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   err.Error(),
			Code:    models.ProcessIdoParticipateInfoErr,
			Data:    nil,
		})
		return
	}
	err = services.SetIdoParticipateInfo(idoParticipateInfo)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   err.Error(),
			Code:    models.SetIdoParticipateInfoErr,
			Data:    nil,
		})
		return
	}
	c.JSON(http.StatusOK, models.JsonResult{
		Success: true,
		Error:   "",
		Code:    models.SUCCESS,
		Data:    nil,
	})
}
