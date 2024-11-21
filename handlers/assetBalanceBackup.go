package handlers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"trade/models"
	"trade/services"
)

func GetAssetBalanceBackup(c *gin.Context) {
	username := c.MustGet("username").(string)
	assetBalanceBackup := services.GetAssetBalanceBackup(username)
	c.JSON(http.StatusOK, models.JsonResult{
		Success: true,
		Error:   models.SuccessErr,
		Code:    models.SUCCESS,
		Data:    assetBalanceBackup.Hash,
	})
}

func UpdateAssetBalanceBackup(c *gin.Context) {
	username := c.MustGet("username").(string)
	hash := c.Query("hash")
	if len(hash) != 64 {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   errors.New("invalid hash length").Error(),
			Code:    models.InvalidHashLengthErr,
			Data:    hash,
		})
	}
	err := services.UpdateAssetBalanceBackup(username, hash)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   err.Error(),
			Code:    models.UpdateAssetBalanceBackupErr,
			Data:    hash,
		})
	}
	c.JSON(http.StatusOK, models.JsonResult{
		Success: true,
		Error:   models.SuccessErr,
		Code:    models.SUCCESS,
		Data:    hash,
	})
}
