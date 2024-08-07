package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"trade/models"
	"trade/services"
	"trade/services/btldb"
)

func LoginHandler(c *gin.Context) {
	var creds models.User
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
	username := creds.Username
	ip := c.ClientIP()
	user, err := services.UpdateUserIpByUsername(username, ip)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	record := models.LoginRecord{
		UserId:            user.ID,
		RecentIpAddresses: ip,
	}
	err = btldb.CreateLoginRecord(&record)
	if err != nil {
		fmt.Println(err)
	}

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
	username := creds.Username
	ip := c.ClientIP()
	user, err := services.UpdateUserIpByUsername(username, ip)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	record := models.LoginRecord{
		UserId:            user.ID,
		RecentIpAddresses: ip,
	}
	err = btldb.CreateLoginRecord(&record)
	if err != nil {
		fmt.Println(err)
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func UserInfoHandler(c *gin.Context) {
	username := c.MustGet("username").(string)
	c.JSON(http.StatusOK, gin.H{"username": username})
}
