package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"trade/btlLog"
	"trade/models"
	"trade/services"
)

func GetUserData(c *gin.Context) {
	username := c.Query("username")
	yaml := c.Query("yaml")
	isYaml, err := strconv.ParseBool(yaml)
	if err != nil {
		btlLog.UserData.Error("ParseBool err:%v", err)
	}
	if isYaml {
		userData, err := services.GetUserDataYaml(username)
		if err != nil {
			c.JSON(http.StatusOK, models.JsonResult{
				Success: false,
				Error:   err.Error(),
				Code:    models.GetUserDataYamlErr,
				Data:    nil,
			})
			return
		}
		c.String(http.StatusOK, "%s", userData)
	} else {
		userData, err := services.GetUserData(username)
		if err != nil {
			c.JSON(http.StatusOK, models.JsonResult{
				Success: false,
				Error:   err.Error(),
				Code:    models.GetUserDataErr,
				Data:    nil,
			})
			return
		}
		c.JSON(http.StatusOK, userData)
	}
}
