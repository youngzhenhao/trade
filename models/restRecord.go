package models

import "gorm.io/gorm"

type RestRecord struct {
	gorm.Model
	Method         string `json:"method"`
	Url            string `json:"url"`
	RequestHeader  string `json:"header"`
	Data           string `json:"data"`
	Payload        string `json:"payload"`
	RequestBody    string `json:"request_body"`
	ResponseHeader string `json:"response_header"`
	ResponseBody   string `json:"response_body"`
	Error          string `json:"error"`
}
