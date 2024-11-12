package handlers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"trade/btlLog"
	"trade/models"
	"trade/services"
)

func GetUserStats(c *gin.Context) {
	yaml := c.Query("yaml")
	_new := c.Query("new")
	_day := c.Query("day")
	_month := c.Query("month")
	isYaml, err := strconv.ParseBool(yaml)
	if err != nil {
		btlLog.UserStats.Error("ParseBool err:%v", err)
	}
	_newInt, err := strconv.Atoi(_new)
	if err != nil {
		btlLog.UserStats.Error("Atoi("+_new+") err:%v", err)
	}
	_newBool := _newInt == 1
	_dayInt, err := strconv.Atoi(_day)
	if err != nil {
		btlLog.UserStats.Error("Atoi("+_day+") err:%v", err)
	}
	_dayBool := _dayInt == 1
	_monthInt, err := strconv.Atoi(_month)
	if err != nil {
		btlLog.UserStats.Error("Atoi("+_month+") err:%v", err)
	}
	_monthBool := _monthInt == 1
	if isYaml {
		userStats, err := services.GetUserStatsYaml(_newBool, _dayBool, _monthBool)
		if err != nil {
			c.JSON(http.StatusOK, models.JsonResult{
				Success: false,
				Error:   err.Error(),
				Code:    models.GetUserStatsYamlErr,
				Data:    nil,
			})
			return
		}
		c.String(http.StatusOK, "%s", userStats)
	} else {
		userStats, err := services.GetUserStats(_newBool, _dayBool, _monthBool)
		if err != nil {
			c.JSON(http.StatusOK, models.JsonResult{
				Success: false,
				Error:   err.Error(),
				Code:    models.GetUserStatsErr,
				Data:    nil,
			})
			return
		}
		c.JSON(http.StatusOK, userStats)
	}
}

func GetSpecifiedDateUserStats(c *gin.Context) {
	day := c.Query("day")
	if !(len(day) == 0 || len(day) == len("20060102")) {
		err := errors.New("date format err, day length wrong")
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   err.Error(),
			Code:    models.DateFormatErr,
			Data:    nil,
		})
		return
	}
	userStats, err := services.GetSpecifiedDateUserStats(day)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   err.Error(),
			Code:    models.GetSpecifiedDateUserStatsErr,
			Data:    nil,
		})
		return
	}
	c.JSON(http.StatusOK, userStats)
}

func DownloadCsv(c *gin.Context) {
	_new := c.Query("new")
	_day := c.Query("day")
	_month := c.Query("month")

	_newInt, err := strconv.Atoi(_new)
	if err != nil {
		btlLog.UserStats.Error("Atoi("+_new+") err:%v", err)
	}
	_newBool := _newInt == 1
	_dayInt, err := strconv.Atoi(_day)
	if err != nil {
		btlLog.UserStats.Error("Atoi("+_day+") err:%v", err)
	}
	_dayBool := _dayInt == 1
	_monthInt, err := strconv.Atoi(_month)
	if err != nil {
		btlLog.UserStats.Error("Atoi("+_month+") err:%v", err)
	}
	_monthBool := _monthInt == 1
	allUserStats, err := services.GetUserStats(true, true, true)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   err.Error(),
			Code:    models.GetUserStatsErr,
			Data:    nil,
		})
		return
	}
	if _newBool {
		newUserToday := allUserStats.NewUserToday
		if newUserToday == nil || len(*newUserToday) == 0 {
			c.String(http.StatusOK, "%s", "no data here")
			return
		}
		if _dayBool || _monthBool {
			err = errors.New("too many query params")
			c.JSON(http.StatusOK, models.JsonResult{
				Success: false,
				Error:   err.Error(),
				Code:    models.TooManyQueryParamsErr,
				Data:    nil,
			})
			return
		}
		newCsv, err := services.StatsUserInfoToCsv("new", allUserStats.NewUserToday)
		if err != nil {
			c.JSON(http.StatusOK, models.JsonResult{
				Success: false,
				Error:   err.Error(),
				Code:    models.StatsUserInfoToCsvErr,
				Data:    nil,
			})
			return
		}
		Download(c, newCsv)
		return
	}
	if _dayBool {
		dailyActiveUser := allUserStats.DailyActiveUser
		if dailyActiveUser == nil || len(*dailyActiveUser) == 0 {
			c.String(http.StatusOK, "%s", "no data here")
			return
		}
		if _monthBool {
			err = errors.New("too many query params")
			c.JSON(http.StatusOK, models.JsonResult{
				Success: false,
				Error:   err.Error(),
				Code:    models.TooManyQueryParamsErr,
				Data:    nil,
			})
			return
		}
		newCsv, err := services.StatsUserInfoToCsv("day", allUserStats.DailyActiveUser)
		if err != nil {
			c.JSON(http.StatusOK, models.JsonResult{
				Success: false,
				Error:   err.Error(),
				Code:    models.StatsUserInfoToCsvErr,
				Data:    nil,
			})
			return
		}
		Download(c, newCsv)
		return
	}
	if _monthBool {
		monthlyActiveUser := allUserStats.MonthlyActiveUser
		if monthlyActiveUser == nil || len(*monthlyActiveUser) == 0 {
			c.String(http.StatusOK, "%s", "no data here")
			return
		}
		newCsv, err := services.StatsUserInfoToCsv("month", monthlyActiveUser)
		if err != nil {
			c.JSON(http.StatusOK, models.JsonResult{
				Success: false,
				Error:   err.Error(),
				Code:    models.StatsUserInfoToCsvErr,
				Data:    nil,
			})
			return
		}
		Download(c, newCsv)
		return
	}
	return
}

func Download(c *gin.Context, path string) {
	file, err := os.Open(path)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "File not found"})
		return
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			return
		}
	}(file)
	filename := filepath.Base(path)
	c.Writer.Header().Set("Content-Disposition", "attachment; filename="+filename)
	c.Writer.Header().Set("Content-Type", "application/octet-stream")
	c.Writer.WriteHeader(http.StatusOK)
	_, err = io.Copy(c.Writer, file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}

func GetActiveUserCount(c *gin.Context) {
	start := c.Query("start")
	end := c.Query("end")
	_start, err := strconv.Atoi(start)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   err.Error() + "(" + start + ")",
			Code:    models.InvalidQueryParamErr,
			Data:    nil,
		})
		return
	}
	_end, err := strconv.Atoi(end)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   err.Error() + "(" + end + ")",
			Code:    models.InvalidQueryParamErr,
			Data:    nil,
		})
		return
	}
	count, err := services.GetActiveUserCountBetween(_start, _end)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   err.Error(),
			Code:    models.GetActiveUserCountBetweenErr,
			Data:    nil,
		})
		return
	}
	c.JSON(http.StatusOK, models.JsonResult{
		Success: false,
		Error:   models.SUCCESS.Error(),
		Code:    models.SUCCESS,
		Data:    count,
	})
}

func GetActiveUserRecord(c *gin.Context) {
	start := c.Query("start")
	end := c.Query("end")
	limit := c.Query("limit")
	offset := c.Query("offset")
	_start, err := strconv.Atoi(start)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   err.Error() + "(" + start + ")",
			Code:    models.InvalidQueryParamErr,
			Data:    nil,
		})
		return
	}
	_end, err := strconv.Atoi(end)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   err.Error() + "(" + end + ")",
			Code:    models.InvalidQueryParamErr,
			Data:    nil,
		})
		return
	}
	_limit, err := strconv.Atoi(limit)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   err.Error() + "(" + limit + ")",
			Code:    models.InvalidQueryParamErr,
			Data:    nil,
		})
		return
	}
	_offset, err := strconv.Atoi(offset)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   err.Error() + "(" + offset + ")",
			Code:    models.InvalidQueryParamErr,
			Data:    nil,
		})
		return
	}
	records, err := services.GetUserActiveRecord(_start, _end, _limit, _offset)
	if err != nil {
		c.JSON(http.StatusOK, models.JsonResult{
			Success: false,
			Error:   err.Error(),
			Code:    models.GetUserActiveRecordErr,
			Data:    nil,
		})
		return
	}
	c.JSON(http.StatusOK, models.JsonResult{
		Success: false,
		Error:   models.SUCCESS.Error(),
		Code:    models.SUCCESS,
		Data:    records,
	})
}
