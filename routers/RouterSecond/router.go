package RouterSecond

import (
	"github.com/gin-gonic/gin"
	"trade/handlers/SecondHander"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	award := r.Group("/award")
	award.POST("/PutInSatoshiAward", SecondHander.PutInSatoshiAward)
	award.POST("/PutAssetAward", SecondHander.PutAssetAward)
	return r
}
