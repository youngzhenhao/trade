package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
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
		return nil, err
	}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			return
		}
	}(response.Body)
	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	var mempoolGetRecommendedFeesResponse MempoolGetRecommendedFeesResponse
	if err = json.Unmarshal(bodyBytes, &mempoolGetRecommendedFeesResponse); err != nil {
		return nil, err
	}
	return &mempoolGetRecommendedFeesResponse, nil
}
