package models

import "encoding/json"

type JsonResult struct {
	Success bool    `json:"success"`
	Error   string  `json:"error"`
	Code    ErrCode `json:"code"`
	Data    any     `json:"data"`
}

type ErrCode int

const (
	DefaultErr ErrCode = -1
	SUCCESS    ErrCode = 200
)

func MakeJsonErrorResult(code ErrCode, errorString string, data any) string {
	jsonResult := JsonResult{
		Error: errorString,
		Code:  code,
		Data:  data,
	}
	if code == SUCCESS {
		jsonResult.Success = true
	} else {
		jsonResult.Success = false
	}
	jsonStr, err := json.Marshal(jsonResult)
	if err != nil {
		return MakeJsonErrorResult(DefaultErr, err.Error(), nil)
	}
	return string(jsonStr)
}
