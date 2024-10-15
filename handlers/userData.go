package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"trade/models"
	"trade/services"
)

func GetUserData(c *gin.Context) {
	username := c.Query("username")
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
	c.JSON(http.StatusOK, models.JsonResult{
		Success: true,
		Error:   "",
		Code:    models.SUCCESS,
		Data:    userData,
	})
}
