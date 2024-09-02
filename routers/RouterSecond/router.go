package RouterSecond

import (
	"github.com/gin-gonic/gin"
	"trade/handlers/SecondHander"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	addrReceive := r.Group("/award")
	addrReceive.POST("/PutInSatoshiAward", SecondHander.PutInSatoshiAward)
	return r
}
