package handlers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"trade/models"
	"trade/services"
)

func GetFairLaunchFollowByUserId(c *gin.Context) {
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
	fairLaunchFollows, err := services.GetFairLaunchFollowsByUserId(userId)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   err.Error(),
			Code:    models.GetFairLaunchFollowsByUserIdErr,
			Data:    nil,
		})
		return
	}
	c.JSON(http.StatusOK, models.JsonResult{
		Success: true,
		Error:   models.SuccessErr,
		Code:    models.SUCCESS,
		Data:    fairLaunchFollows,
	})
}

func SetFollowFairLaunchInfo(c *gin.Context) {
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
	var fairLaunchFollowSetRequest models.FairLaunchFollowSetRequest
	err = c.BindJSON(&fairLaunchFollowSetRequest)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   "Should Bind JSON setFairLaunchInfoRequest. " + err.Error(),
			Code:    models.ShouldBindJsonErr,
			Data:    nil,
		})
		return
	}
	isValid, err := services.IsFairLaunchInfoIdAndAssetIdValid(fairLaunchFollowSetRequest.FairLaunchInfoId, fairLaunchFollowSetRequest.AssetId)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   err.Error(),
			Code:    models.IsFairLaunchInfoIdAndAssetIdValidErr,
			Data:    nil,
		})
		return
	}
	if !isValid {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   errors.New("fair launch info asset id is invalid").Error(),
			Code:    models.FairLaunchInfoAssetIdInvalidErr,
			Data:    nil,
		})
		return
	}
	fairLaunchFollow := services.ProcessFairLaunchFollowSetRequest(userId, username, fairLaunchFollowSetRequest)
	err = services.SetFollowFairLaunchInfo(&fairLaunchFollow)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   err.Error(),
			Code:    models.SetFollowFairLaunchInfoErr,
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

func SetUnfollowFairLaunchInfo(c *gin.Context) {
	assetId := c.Param("asset_id")
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
	err = services.SetUnfollowFairLaunchInfo(userId, assetId)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   err.Error(),
			Code:    models.SetUnfollowFairLaunchInfoErr,
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

func GetAllFairLaunchFollowSimplified(c *gin.Context) {
	fairLaunchFollows, err := services.GetAllFairLaunchFollowSimplified()
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   err.Error(),
			Code:    models.GetAllFairLaunchFollowSimplifiedErr,
			Data:    nil,
		})
		return
	}
	c.JSON(http.StatusOK, models.JsonResult{
		Success: true,
		Error:   models.SuccessErr,
		Code:    models.SUCCESS,
		Data:    fairLaunchFollows,
	})
}

func GetFairLaunchInfoIsFollowed(c *gin.Context) {
	assetId := c.Param("asset_id")
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
	isFairLaunchFollowed := services.IsFairLaunchFollowed(userId, assetId)
	c.JSON(http.StatusOK, models.JsonResult{
		Success: true,
		Error:   models.SuccessErr,
		Code:    models.SUCCESS,
		Data:    isFairLaunchFollowed,
	})
}
