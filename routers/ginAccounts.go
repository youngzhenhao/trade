package routers

import "github.com/gin-gonic/gin"

var (
	basicAuthAccounts = gin.Accounts{
		"foo":   "bar",
		"admin": "123456",
	}
)
