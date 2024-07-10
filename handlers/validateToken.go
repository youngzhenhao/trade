package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"trade/models"
)

func ValidateTokenPing(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"error": models.SuccessErr,
	})
}
