package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"trade/utils"
)

type MempoolGetRecommendedFeesResponse struct {
	FastestFee  int `json:"fastestFee"`
	HalfHourFee int `json:"halfHourFee"`
	HourFee     int `json:"hourFee"`
	EconomyFee  int `json:"economyFee"`
	MinimumFee  int `json:"minimumFee"`
}

func MempoolGetRecommendedFees() (*MempoolGetRecommendedFeesResponse, error) {
	url := "https://mempool.space/api/v1/fees/recommended"
	client := &http.Client{}
	var jsonData []byte
	request, err := http.NewRequest("GET", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, utils.AppendErrorInfo(err, "NewRequest")
	}
	response, err := client.Do(request)
	if err != nil {
		return nil, utils.AppendErrorInfo(err, "client.Do")
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			return
		}
	}(response.Body)
	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, utils.AppendErrorInfo(err, "ReadAll")
	}
	var mempoolGetRecommendedFeesResponse MempoolGetRecommendedFeesResponse
	if err = json.Unmarshal(bodyBytes, &mempoolGetRecommendedFeesResponse); err != nil {
		return nil, utils.AppendErrorInfo(err, "Unmarshal")
	}
	return &mempoolGetRecommendedFeesResponse, nil
}
