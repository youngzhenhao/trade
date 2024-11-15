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
	count, err := services.GetActiveUserCountBetween(start, end)
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
	records, err := services.GetUserActiveRecord(start, end, _limit, _offset)
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

type Result2 struct {
	Errno  int         `json:"errno"`
	ErrMsg string      `json:"errmsg"`
	Data   interface{} `json:"data"`
}

func GetDateLoginCount(c *gin.Context) {
	start := c.Query("start")
	end := c.Query("end")
	count, err := services.GetDateLoginCount(start, end)
	if err != nil {
		c.JSON(http.StatusOK, Result2{
			Errno:  models.GetDateLoginCountErr.Code(),
			ErrMsg: err.Error(),
			Data:   0,
		})
		return
	}
	c.JSON(http.StatusOK, Result2{
		Errno:  0,
		ErrMsg: models.SUCCESS.Error(),
		Data:   count,
	})
}

func GetDateIpLoginRecord(c *gin.Context) {
	start := c.Query("start")
	end := c.Query("end")
	page := c.Query("page")
	size := c.Query("size")
	_new := c.Query("new")
	records := new([]services.DateIpLoginRecord)
	_page, err := strconv.Atoi(page)
	var pageNumber int
	if err != nil {
		c.JSON(http.StatusOK, Result2{
			Errno:  models.InvalidQueryParamErr.Code(),
			ErrMsg: err.Error() + "(" + page + ")",
			Data: gin.H{
				"total_page": pageNumber,
				"records":    records,
			},
		})
		return
	}
	_size, err := strconv.Atoi(size)
	if err != nil {
		c.JSON(http.StatusOK, Result2{
			Errno:  models.InvalidQueryParamErr.Code(),
			ErrMsg: err.Error() + "(" + size + ")",
			Data: gin.H{
				"total_page": pageNumber,
				"records":    records,
			},
		})
		return
	}

	if _new == "1" {
		pageNumber, err = services.GetNewUserPageNumber(start, end, _size)
	} else {
		pageNumber, err = services.GetDateIpLoginPageNumber(start, end, _size)
	}

	if _page > pageNumber {
		err = errors.New("page is out of range(" + strconv.Itoa(pageNumber) + ")")
		c.JSON(http.StatusOK, Result2{
			Errno:  models.PageNumberOutOfRangeErr.Code(),
			ErrMsg: err.Error(),
			Data: gin.H{
				"total_page": pageNumber,
				"records":    records,
			},
		})
		return
	}
	if _page < 1 || _size < 1 {
		err = errors.New("page or size is invalid(" + strconv.Itoa(_page) + "," + strconv.Itoa(_size) + ")")
		c.JSON(http.StatusOK, Result2{
			Errno:  models.NegativeValueErr.Code(),
			ErrMsg: err.Error() + "(" + size + ")",
			Data: gin.H{
				"total_page": pageNumber,
				"records":    records,
			},
		})
		return
	}
	limit, offset := services.PageAndSizeToLimitAndOffset(uint(_page), uint(_size))

	if _new == "1" {
		records, err = services.GetNewUserRecord(start, end, int(limit), int(offset))
	} else {
		records, err = services.GetDateIpLoginRecord(start, end, int(limit), int(offset))
	}

	if err != nil {
		c.JSON(http.StatusOK, Result2{
			Errno:  models.GetDateIpLoginRecordErr.Code(),
			ErrMsg: err.Error(),
			Data: gin.H{
				"total_page": pageNumber,
				"records":    records,
			},
		})
		return
	}
	c.JSON(http.StatusOK, Result2{
		Errno:  0,
		ErrMsg: models.SUCCESS.Error(),
		Data: gin.H{
			"total_page": pageNumber,
			"records":    records,
		},
	})
}

func GetDateIpLoginRecordCount(c *gin.Context) {
	start := c.Query("start")
	end := c.Query("end")
	count, err := services.GetDateIpLoginCount(start, end)
	if err != nil {
		c.JSON(http.StatusOK, Result2{
			Errno:  models.GetDateIpLoginCountErr.Code(),
			ErrMsg: err.Error(),
			Data:   0,
		})
		return
	}
	c.JSON(http.StatusOK, Result2{
		Errno:  0,
		ErrMsg: models.SUCCESS.Error(),
		Data:   count,
	})
}

func GetNewUserCount(c *gin.Context) {
	start := c.Query("start")
	end := c.Query("end")
	count, err := services.GetNewUserCount(start, end)
	if err != nil {
		c.JSON(http.StatusOK, Result2{
			Errno:  models.GetNewUserCountErr.Code(),
			ErrMsg: err.Error(),
			Data:   0,
		})
		return
	}
	c.JSON(http.StatusOK, Result2{
		Errno:  0,
		ErrMsg: models.SUCCESS.Error(),
		Data:   count,
	})
}
