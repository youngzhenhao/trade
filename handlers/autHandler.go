package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
	"trade/middleware"
	"trade/models"
	"trade/services"
)

func GetNonceHandler(c *gin.Context) {
	var getNonce models.NonceRequest
	if err := c.ShouldBindJSON(&getNonce); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	decryptUserName, err := services.ValidAndDecrypt(getNonce.Username)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	nonce, err := services.GenerateNonce()
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return

	}
	nonceStored, err := services.StoreNonceInRedis(decryptUserName, nonce)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"nonce": nonceStored})
}
func GetDeviceIdHandler(c *gin.Context) {
	var getNonce models.NonceRequest
	if err := c.ShouldBindJSON(&getNonce); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	decryptUserName, err := services.ValidAndDecrypt(getNonce.Username)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	encryptDeviceID, encodedSalt, err := services.ProcessDeviceRequest(getNonce.Nonce, decryptUserName)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"encryptDeviceID": encryptDeviceID,
		"encodedSalt":     encodedSalt,
	})
}

func LoginHandler(c *gin.Context) {
	var creds models.User
	if err := c.ShouldBindJSON(&creds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	token, err := services.Login(&creds)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	// @dev: Update user ip by client ip
	ip := c.ClientIP()
	path := c.Request.URL.Path
	go middleware.InsertLoginInfo(creds.Username, ip, path)
	{
		go middleware.RecodeDateIpLogin(creds.Username, time.Now().Format(time.DateOnly), ip)
		go middleware.RecodeDateLogin(creds.Username, time.Now().Format(time.DateOnly))
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

	ip := c.ClientIP()
	path := c.Request.URL.Path
	go middleware.InsertLoginInfo(creds.Username, ip, path)
	{
		go middleware.RecodeDateIpLogin(creds.Username, time.Now().Format(time.DateOnly), ip)
		go middleware.RecodeDateLogin(creds.Username, time.Now().Format(time.DateOnly))
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}

func UserInfoHandler(c *gin.Context) {
	username := c.MustGet("username").(string)
	c.JSON(http.StatusOK, gin.H{"username": username})
}
