package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"sync"
	"trade/models"
	"trade/utils"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			return
		}
		parts := strings.Fields(authHeader)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header format must be 'Bearer {token}'"})
			return
		}
		tokenString := parts[1]
		claims, err := ValidateToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		if claims.Username == "admin" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Illegal login"})
			return
		}
		// Update user's recent IP address
		ip := c.ClientIP()
		path := c.Request.URL.Path
		go InsertLoginInfo(claims.Username, ip, path)
		// Store the username in the context of the request
		c.Set("username", claims.Username)
		c.Next()
	}
}

// InsertLoginInfo 记录登录信息
var LoginInfoMutex = sync.Mutex{}

func InsertLoginInfo(userName, ip, path string) {
	LoginInfoMutex.Lock()
	defer LoginInfoMutex.Unlock()

	var user models.User
	err := DB.Where("user_name = ?", userName).First(&user).Error
	if err != nil {
		fmt.Println("InsertLoginInfo Query Error:", err)
		return
	}
	user.RecentIpAddresses = ip
	user.RecentLoginTime = utils.GetTimestamp()
	err = DB.Save(&user).Error
	if err != nil {
		fmt.Println("InsertLoginInfo Save Error:", err)
		return
	}

	record := models.LoginRecord{
		UserId:            user.ID,
		RecentIpAddresses: user.RecentIpAddresses,
		Path:              path,
		LoginTime:         user.RecentLoginTime,
	}
	err = DB.Create(&record).Error
	if err != nil {
		fmt.Println(err, "insertLoginInfo", user.ID, ",", ip, ",", path)
	}
}
