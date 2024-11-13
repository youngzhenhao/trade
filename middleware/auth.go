package middleware

import (
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"net/http"
	"strings"
	"sync"
	"time"
	"trade/btlLog"
	"trade/models"
	"trade/utils"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if the request is for the mini
		ClientType := c.GetHeader("ClientType")
		if ClientType == "mini" {
			authHeader := c.GetHeader("Authorization")
			if authHeader == "" || !strings.HasPrefix(authHeader, "Basic ") {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
				return
			}
			// 获取 Base64 编码的部分
			encodedCredentials := strings.TrimPrefix(authHeader, "Basic ")
			decodedCredentials, err := base64.StdEncoding.DecodeString(encodedCredentials)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header"})
				return
			}
			// 分割用户名和密码
			credentials := strings.SplitN(string(decodedCredentials), ":", 2)
			if len(credentials) != 2 {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
				return
			}
			username := credentials[0]
			password := credentials[1]
			if username == "admin" {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Illegal login"})
				return
			}
			err = ValidateMiniUser(username, password)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
				return
			}
			// Update user's recent IP address
			ip := c.ClientIP()
			path := c.Request.URL.Path
			go InsertLoginInfo(username, ip, path)
			{
				go RecodeDateIpLogin(username, ip, time.Now().Format(time.DateOnly))
			}
			// Store the username in the context of the request
			c.Set("username", username)
			c.Next()
			return
		}

		// Check if the request is authorized
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":         "value of Authorization in header is null",
				"Authorization": authHeader,
			})
			return
		}
		parts := strings.Fields(authHeader)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":         "Authorization header format must be 'Bearer {token}'",
				"Authorization": authHeader,
			})
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
		{
			go RecodeDateIpLogin(claims.Username, ip, time.Now().Format(time.DateOnly))
		}
		// Store the username in the context of the request
		c.Set("username", claims.Username)
		c.Next()
	}
}
func ValidateMiniUser(username, password string) error {
	// Check if the username and password are correct
	var user models.User
	result := DB.Where("user_name = ?", username).First(&user)
	if result.Error != nil {
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// If there are other database errors, an error is returned
			return result.Error
		} else {
			user.Username = username
			password, err := hashPassword(password)
			if err != nil {
				return err
			}
			user.Password = password
			err = DB.Create(&user).Error
			if err != nil {
				return err
			}
		}
	}
	if !CheckPassword(user.Password, password) {
		return errors.New("invalid credentials")
	}
	return nil
}
func CheckPassword(hashedPassword, password string) bool {
	// bcrypt.CompareHashAndPassword Compare the hashed password with the password entered by the user. If there is a match, nil is returned.
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
func hashPassword(password string) (string, error) {
	// Passwords are encrypted using the bcrypt algorithm
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

var LoginInfoMutex = sync.Mutex{}

// InsertLoginInfo 记录登录信息
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

// RecodeDateIpLogin
func RecodeDateIpLogin(username string, date string, ip string) {
	dateIpLogin := models.DateIpLogin{
		Username: username,
		Date:     date,
		Ip:       ip,
	}
	if err := func(d *models.DateIpLogin) error {
		if d.Username == "" || d.Date == "" || d.Ip == "" {
			return errors.New("username, date or ip is null")
		}
		return DB.Create(d).Error
	}(&dateIpLogin); err != nil {
		btlLog.DateIpLogin.Error("%v", err)
	}
}
