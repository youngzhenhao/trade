package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"trade/models"
	"trade/services"
)

func SetBtcUtxo(c *gin.Context) {
	username := c.MustGet("username").(string)
	var err error

	var requests []models.UnspentUtxo
	err = c.ShouldBindJSON(&requests)
	if err != nil {
		c.JSON(http.StatusOK, Result2{
			Errno:  models.ShouldBindJsonErr.Code(),
			ErrMsg: err.Error(),
			Data:   nil,
		})
		return
	}

	err = services.SetBtcUtxo(username, &requests)
	if err != nil {
		c.JSON(http.StatusOK, Result2{
			Errno:  models.SetBtcUtxoErr.Code(),
			ErrMsg: err.Error(),
			Data:   nil,
		})
		return
	}

	c.JSON(http.StatusOK, Result2{
		Errno:  0,
		ErrMsg: models.SUCCESS.Error(),
		Data:   nil,
	})
}
