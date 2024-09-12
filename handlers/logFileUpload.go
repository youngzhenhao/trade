package handlers

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"os"
	"path"
	"strconv"
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
			Code:    models.FileSizeTooLargeErr,
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

func GetAllLogFilesByDeviceId(c *gin.Context) {
	allLogFiles, err := services.GetAllLogFiles()
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   err.Error(),
			Code:    models.GetAllLogFilesErr,
			Data:    nil,
		})
		return
	}
	allDeviceIdMapLogFileUploads := services.LogFileUploadsToDeviceIdMapLogFileUploads(allLogFiles)
	c.JSON(http.StatusOK, models.JsonResult{
		Success: true,
		Error:   models.SUCCESS.Error(),
		Code:    models.SUCCESS,
		Data:    allDeviceIdMapLogFileUploads,
	})
}

func GetAllLogFiles(c *gin.Context) {
	allLogFiles, err := services.GetAllLogFiles()
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   err.Error(),
			Code:    models.GetAllLogFilesErr,
			Data:    nil,
		})
		return
	}
	c.JSON(http.StatusOK, models.JsonResult{
		Success: true,
		Error:   models.SUCCESS.Error(),
		Code:    models.SUCCESS,
		Data:    allLogFiles,
	})
}

func DownloadLogFileById(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   err.Error(),
			Code:    models.IdAtoiErr,
			Data:    nil,
		})
		return
	}
	logFile, err := services.GetFileUpload(uint(id))
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   err.Error(),
			Code:    models.GetFileUploadErr,
			Data:    nil,
		})
		return
	}
	filename := logFile.OriginFileName
	filePath := logFile.FileSavePath
	var file *os.File
	file, err = os.Open(filePath)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   err.Error(),
			Code:    models.OsOpenFileErr,
			Data:    nil,
		})
		return
	}
	defer func(file *os.File) {
		err = file.Close()
		if err != nil {
			fmt.Println(err)
			return
		}
	}(file)
	c.Writer.Header().Set("Content-Disposition", "attachment; filename="+filename)
	c.Writer.Header().Set("Content-Type", "application/octet-stream")
	c.Writer.WriteHeader(http.StatusOK)
	_, err = io.Copy(c.Writer, file)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   err.Error(),
			Code:    models.IoCopyFIleErr,
			Data:    nil,
		})
		return
	}
}
