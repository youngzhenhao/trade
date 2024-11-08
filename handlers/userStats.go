package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
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
		userStats, err := services.GetUserStatsYaml(_newBool, _dayBool, _monthBool)
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
