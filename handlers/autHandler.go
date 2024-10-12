package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"trade/middleware"
	"trade/models"
	"trade/services"
)

func LoginHandler(c *gin.Context) {
	var creds models.User
	fmt.Println("login start")
	if err := c.ShouldBindJSON(&creds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	token, err := services.Login(creds)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	// @dev: Update user ip by client ip
	ip := c.ClientIP()
	fmt.Println("login ip:", ip)
	path := c.Request.URL.Path
	go middleware.InsertLoginInfo(creds.Username, ip, path)

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func RefreshTokenHandler(c *gin.Context) {
	var creds models.User
	if err := c.ShouldBindJSON(&creds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	flag := services.SplitStringAndVerifyChecksum(creds.Password)
	if !flag {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Checksum error"})
		return
	}
	token, err := services.ValidateUserAndGenerateToken(creds)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	ip := c.ClientIP()
	path := c.Request.URL.Path
	go middleware.InsertLoginInfo(creds.Username, ip, path)

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func UserInfoHandler(c *gin.Context) {
	username := c.MustGet("username").(string)
	c.JSON(http.StatusOK, gin.H{"username": username})
}
