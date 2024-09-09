package handlers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"path"
	"trade/config"
	"trade/models"
	"trade/services"
	"trade/utils"
)

func UploadLogFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   err.Error(),
			Code:    models.FormFileErr,
			Data:    nil,
		})
		return
	}
	if file.Size > 15*1024*1024 {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   errors.New("file too large, its size is more than 15MB").Error(),
			Code:    models.OsGetPwdErr,
			Data:    nil,
		})
		return
	}
	deviceId := c.PostForm("device_id")
	if deviceId == "" {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   errors.New("device_id is null").Error(),
			Code:    models.DeviceIdIsNullErr,
			Data:    nil,
		})
		return
	}
	info := c.PostForm("info")
	var pwd string
	pwd, err = os.Getwd()
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   err.Error(),
			Code:    models.OsGetPwdErr,
			Data:    nil,
		})
		return
	}
	network := config.GetLoadConfig().NetWork
	timeStr := utils.GetNowTimeStringWithHyphens()
	dst := path.Join(pwd, "log_files", network, deviceId, timeStr+"-"+file.Filename)
	err = c.SaveUploadedFile(file, dst)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   err.Error(),
			Code:    models.SaveUploadedFileErr,
			Data:    nil,
		})
		return
	}
	err = services.CreateLogFileUpload(&models.LogFileUpload{
		DeviceId:       deviceId,
		OriginFileName: file.Filename,
		FileSavePath:   dst,
		Info:           info,
	})
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   err.Error(),
			Code:    models.CreateLogFileUploadErr,
			Data:    nil,
		})
		return
	}
	c.JSON(http.StatusOK, models.JsonResult{
		Success: true,
		Error:   models.SUCCESS.Error(),
		Code:    models.SUCCESS,
		Data:    nil,
	})
}
