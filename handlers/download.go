package handlers

import (
	"github.com/gin-gonic/gin"

	"net/http"
	"path"
	"trade/models"
	"trade/services"
)

func CsvDownloadCaptcha(c *gin.Context) {
	name := c.Query("name")
	captcha := c.Query("captcha")
	err := services.RedisGetVerify("download", captcha)
	if err != nil {
		c.JSON(http.StatusOK, Result2{
			Errno:  models.RedisGetVerifyErr.Code(),
			ErrMsg: err.Error(),
			Data:   nil,
		})
		return
	}
	_path := path.Join("./", "csv", name)
	Download(c, _path)
	return
}
