package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"trade/models"
	"trade/services"
)

func GetAssetMeta(c *gin.Context) {
	var assetIds []string
	err := c.ShouldBindJSON(&assetIds)
	if err != nil {
		c.JSON(http.StatusOK, Result2{
			Errno:  models.ShouldBindJsonErr.Code(),
			ErrMsg: err.Error(),
			Data:   0,
		})
		return
	}
	datas, errs := services.GetAssetMeta(assetIds)
	c.JSON(http.StatusOK, Result2{
		Errno:  0,
		ErrMsg: models.SUCCESS.Error(),
		Data: gin.H{
			"datas": datas,
			"errs":  errs,
		},
	})
}
