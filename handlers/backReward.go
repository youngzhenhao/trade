package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"trade/models"
	"trade/services"
)

func GetBackRewards(c *gin.Context) {
	username := c.Query("username")
	backRewards, err := services.GetBackRewards(username)
	if err != nil {
		c.JSON(http.StatusOK, Result2{
			Errno:  models.GetBackRewardsErr.Code(),
			ErrMsg: err.Error(),
			Data:   backRewards,
		})
		return
	}
	c.JSON(http.StatusOK, Result2{
		Errno:  0,
		ErrMsg: models.SUCCESS.Error(),
		Data:   backRewards,
	})
}
