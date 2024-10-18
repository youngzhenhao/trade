package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"trade/models"
	"trade/services"
)

func AddNftPresaleWhitelists(c *gin.Context) {
	var nftPresaleWhitelistSetRequests []models.NftPresaleWhitelistSetRequest
	err := c.ShouldBindJSON(&nftPresaleWhitelistSetRequests)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   err.Error(),
			Code:    models.ShouldBindJsonErr,
			Data:    nil,
		})
		return
	}
	err = services.AddWhitelistsByRequests(&nftPresaleWhitelistSetRequests)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   err.Error(),
			Code:    models.AddWhitelistsByRequestsErr,
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
