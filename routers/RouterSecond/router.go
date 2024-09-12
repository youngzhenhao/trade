package RouterSecond

import (
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"trade/config"
	"trade/handlers"
	"trade/handlers/SecondHander"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	award := r.Group("/award")
	award.POST("/PutInSatoshiAward", SecondHander.PutInSatoshiAward)
	award.POST("/PutAssetAward", SecondHander.PutAssetAward)
	username := base64.StdEncoding.EncodeToString([]byte(config.GetLoadConfig().AdminUser.Username))
	password := base64.StdEncoding.EncodeToString([]byte(config.GetLoadConfig().AdminUser.Password))
	authorized := r.Group("/fair_launch/auth_op", gin.BasicAuth(gin.Accounts{
		username: password,
	}))
	authorized.POST("/refund", handlers.RefundUserFirstMintByUsernameAndAssetId)
	return r
}
