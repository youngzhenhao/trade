package services

import (
	"encoding/json"
	"fmt"
	"github.com/lightninglabs/taproot-assets/taprpc"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"trade/config"
	"trade/utils"
)

// @notice: Deprecated
func PostPhoneToNewAddr(remotePort string, assetId string, amount int) (*taprpc.Addr, error) {
	frpsForwardSocket := fmt.Sprintf("%s:%s", config.GetLoadConfig().FrpsServer, remotePort)
	targetUrl := "http://" + frpsForwardSocket + "/newAddr"
	payload := url.Values{"asset_id": {assetId}, "amount": {strconv.Itoa(amount)}}
	response, err := http.PostForm(targetUrl, payload)
	if err != nil {
		utils.LogError("http.PostForm.", err)
		return nil, err
	}
	bodyBytes, _ := io.ReadAll(response.Body)
	var addrResponse struct {
		Success bool   `json:"success"`
		Error   string `json:"error"`
		Data    *taprpc.Addr
	}
	if err := json.Unmarshal(bodyBytes, &addrResponse); err != nil {
		utils.LogError("PPTNA json.Unmarshal.", err)
		return nil, err
	}
	return addrResponse.Data, nil
}
