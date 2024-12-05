package satBackQueue

import (
	"bytes"
	"encoding/json"
	"gorm.io/gorm"
	"io"
	"math/rand"
	"net/http"
	"time"
	"trade/btlLog"
	"trade/middleware"
	"trade/models"
	"trade/utils"
)

const host = "172.27.16.10:7040"

type queueTopic string

const (
	claimAsset         queueTopic = "claimAsset"
	purchasePresaleNFT queueTopic = "purchasePresaleNFT"
)

func (q queueTopic) String() string {
	return string(q)
}

type Request struct {
	Data string `json:"data"`
}

type Response struct {
	Qid string `json:"qid"`
	Rid string `json:"rid"`
}

type FeeInfo struct {
	ID       uint   `json:"id"`
	NpubKey  string `json:"npub_key"`
	AssetsID string `json:"assets_id"`
	HandFee  int    `json:"hand_fee"`
}

func Push(topic queueTopic, qid string, request Request) (Response, error) {
	body, err := Post(topic, qid, request)
	if err != nil {
		return Response{}, err
	}
	var response Response
	err = json.Unmarshal(body, &response)
	if err != nil {
		return Response{}, err
	}
	return response, nil
}

func Post(topic queueTopic, qid string, data any) ([]byte, error) {
	url := "http://" + host + "/q/" + topic.String() + "?qid=" + qid
	requestJsonBytes, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	payload := bytes.NewBuffer(requestJsonBytes)
	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		return nil, err
	}

	//req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {

		func(url string, method string, req *http.Request, res *http.Response, err error) {
			var errInfo string
			if err != nil {
				errInfo = err.Error()
			}
			var requestHeader, responseHeader, responseBody []byte
			requestHeader, _ = json.Marshal(req.Header)
			if res != nil {
				responseHeader, _ = json.Marshal(res.Header)
				responseBody, _ = io.ReadAll(res.Body)
			}
			var restRecord = models.RestRecord{
				Method:         "POST",
				Url:            url,
				RequestHeader:  string(requestHeader),
				Data:           string(requestJsonBytes),
				ResponseHeader: string(responseHeader),
				ResponseBody:   string(responseBody),
				Error:          errInfo,
			}
			_ = middleware.DB.Model(&models.RestRecord{}).Create(&restRecord).Error
		}(url, "POST", req, res, err)

		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			return
		}
	}(res.Body)

	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			return
		}
	}(req.Body)

	body, err := io.ReadAll(res.Body)

	func(url string, method string, req *http.Request, res *http.Response, body []byte, err error) {
		var errInfo string
		if err != nil {
			errInfo = err.Error()
		}
		var requestHeader, responseHeader, responseBody []byte
		requestHeader, _ = json.Marshal(req.Header)
		responseHeader, _ = json.Marshal(res.Header)
		responseBody = body
		var restRecord = models.RestRecord{
			Method:         "POST",
			Url:            url,
			RequestHeader:  string(requestHeader),
			Data:           string(requestJsonBytes),
			ResponseHeader: string(responseHeader),
			ResponseBody:   string(responseBody),
			Error:          errInfo,
		}
		_ = middleware.DB.Model(&models.RestRecord{}).Create(&restRecord).Error
	}(url, "POST", req, res, body, err)

	return body, nil
}

func PushClaimAsset(info FeeInfo) (string, string, Response, error) {
	topic := claimAsset
	seed := rand.New(rand.NewSource(time.Now().UnixNano()))
	qid := utils.RandString(seed, 64)
	data, err := json.Marshal(info)
	if err != nil {
		return "", "", Response{}, err
	}
	request := Request{
		Data: string(data),
	}
	var response Response
	response, err = Push(topic, qid, request)
	return qid, string(data), response, err
}

func PushPurchasePresaleNFT(info FeeInfo) (string, string, Response, error) {
	topic := purchasePresaleNFT
	seed := rand.New(rand.NewSource(time.Now().UnixNano()))
	qid := utils.RandString(seed, 64)
	data, err := json.Marshal(info)
	if err != nil {
		return "", "", Response{}, err
	}
	request := Request{
		Data: string(data),
	}
	var response Response
	response, err = Push(topic, qid, request)
	return qid, string(data), response, err
}

type fairLaunchMintedInfoRecord struct {
	Id           uint   `json:"id"`
	AssetID      string `json:"asset_id"`
	Username     string `json:"username"`
	MintedGasFee int    `json:"minted_gas_fee"`
}

func GetNotPushedClaimAsset() ([]FeeInfo, error) {
	var err error
	var feeInfos []FeeInfo
	var fairLaunchMintedInfoRecords []fairLaunchMintedInfoRecord
	err = middleware.DB.
		Model(&models.FairLaunchMintedInfo{}).
		Select("id, asset_id, username, minted_gas_fee").
		Where("state > ? and is_pushed_queue = ?", models.FairLaunchMintedStatePaidPending, false).
		Order("id desc").
		Scan(&fairLaunchMintedInfoRecords).Error
	if err != nil {
		return feeInfos, utils.AppendErrorInfo(err, "Scan fairLaunchMintedInfoRecords")
	}
	for _, _fairLaunchMintedInfoRecord := range fairLaunchMintedInfoRecords {
		feeInfo := FeeInfo{
			ID:       _fairLaunchMintedInfoRecord.Id,
			NpubKey:  _fairLaunchMintedInfoRecord.Username,
			AssetsID: _fairLaunchMintedInfoRecord.AssetID,
			HandFee:  _fairLaunchMintedInfoRecord.MintedGasFee,
		}
		feeInfos = append(feeInfos, feeInfo)
	}
	return feeInfos, nil
}

type nftPresaleRecord struct {
	Id            uint   `json:"id"`
	AssetId       string `json:"asset_id"`
	BuyerUsername string `json:"buyer_username"`
	Price         int    `json:"price"`
}

func GetNotPushedPurchasePresaleNFT() ([]FeeInfo, error) {
	var err error
	var feeInfos []FeeInfo
	var nftPresaleRecords []nftPresaleRecord
	err = middleware.DB.
		Model(&models.NftPresale{}).
		Select("id, asset_id, buyer_username, price").
		Where("state > ? and is_pushed_queue = ?", models.NftPresaleStatePaidPending, false).
		Order("id desc").
		Scan(&nftPresaleRecords).Error
	if err != nil {
		return feeInfos, utils.AppendErrorInfo(err, "Scan nftPresaleRecords")
	}
	for _, _nftPresaleRecord := range nftPresaleRecords {
		feeInfo := FeeInfo{
			ID:       _nftPresaleRecord.Id,
			NpubKey:  _nftPresaleRecord.BuyerUsername,
			AssetsID: _nftPresaleRecord.AssetId,
			HandFee:  _nftPresaleRecord.Price,
		}
		feeInfos = append(feeInfos, feeInfo)
	}
	return feeInfos, nil
}

type PushQueueRecord struct {
	gorm.Model
	InfoID    uint       `json:"info_id" gorm:"index"`
	NpubKey   string     `json:"npub_key" gorm:"type:varchar(255);index"`
	AssetsID  string     `json:"assets_id" gorm:"type:varchar(255);index"`
	HandFee   int        `json:"hand_fee" gorm:"index"`
	Topic     queueTopic `json:"topic" gorm:"type:varchar(255);index"`
	Qid       string     `json:"qid" gorm:"type:varchar(255);index"`
	Data      string     `json:"data" gorm:"type:varchar(255);index"`
	IsSuccess bool       `json:"is_success" gorm:"index"`
	Rid       string     `json:"rid" gorm:"type:varchar(255);index"`
	Error     string     `json:"error" gorm:"type:varchar(255);index"`
}

func GetAndPushClaimAsset() {
	topic := claimAsset
	feeInfos, err := GetNotPushedClaimAsset()
	if err != nil {
		btlLog.PushQueue.Error("%v", utils.AppendErrorInfo(err, "GetNotPushedClaimAsset"))
	}
	for _, feeInfo := range feeInfos {
		var response Response
		var qid, data string
		qid, data, response, err = PushClaimAsset(feeInfo)

		if err != nil {
			var pushQueueRecord = PushQueueRecord{
				InfoID:    feeInfo.ID,
				NpubKey:   feeInfo.NpubKey,
				AssetsID:  feeInfo.AssetsID,
				HandFee:   feeInfo.HandFee,
				Topic:     topic,
				Qid:       qid,
				Data:      data,
				IsSuccess: false,
				Rid:       "",
				Error:     err.Error(),
			}
			_err := middleware.DB.Model(&PushQueueRecord{}).Create(&pushQueueRecord)
			if _err != nil {
				btlLog.PushQueue.Error("Create _err:\n%v\nPQR:\n%v", _err.Error, utils.ValueJsonString(pushQueueRecord))
			}
		} else {
			var pushQueueRecord = PushQueueRecord{
				InfoID:    feeInfo.ID,
				NpubKey:   feeInfo.NpubKey,
				AssetsID:  feeInfo.AssetsID,
				HandFee:   feeInfo.HandFee,
				Topic:     topic,
				Qid:       qid,
				Data:      data,
				IsSuccess: true,
				Rid:       response.Rid,
				Error:     "",
			}
			_err := middleware.DB.Model(&PushQueueRecord{}).Create(&pushQueueRecord).Error
			if _err != nil {
				btlLog.PushQueue.Error("Create _err:\n%v\nPQR:\n%v", _err, utils.ValueJsonString(pushQueueRecord))
			}
			_err = middleware.DB.Model(&models.FairLaunchMintedInfo{}).
				Where("id = ?", feeInfo.ID).
				Update("is_pushed_queue", true).Error
			if _err != nil {
				btlLog.PushQueue.Error("Update FairLaunchMintedInfo _err:\n%v\nid:\n%v", _err, feeInfo.ID)
			}
		}
	}
	return
}

func GetAndPushPurchasePresaleNFT() {
	topic := purchasePresaleNFT
	feeInfos, err := GetNotPushedPurchasePresaleNFT()
	if err != nil {
		btlLog.PushQueue.Error("%v", utils.AppendErrorInfo(err, "GetNotPushedPurchasePresaleNFT"))
	}
	for _, feeInfo := range feeInfos {
		var response Response
		var qid, data string
		qid, data, response, err = PushPurchasePresaleNFT(feeInfo)
		if err != nil {
			var pushQueueRecord = PushQueueRecord{
				InfoID:    feeInfo.ID,
				NpubKey:   feeInfo.NpubKey,
				AssetsID:  feeInfo.AssetsID,
				HandFee:   feeInfo.HandFee,
				Topic:     topic,
				Qid:       qid,
				Data:      data,
				IsSuccess: false,
				Rid:       "",
				Error:     err.Error(),
			}
			_err := middleware.DB.Model(&PushQueueRecord{}).Create(&pushQueueRecord)
			if _err != nil {
				btlLog.PushQueue.Error("Create _err:\n%v\nPQR:\n%v", _err.Error, utils.ValueJsonString(pushQueueRecord))
			}
		} else {
			var pushQueueRecord = PushQueueRecord{
				InfoID:    feeInfo.ID,
				NpubKey:   feeInfo.NpubKey,
				AssetsID:  feeInfo.AssetsID,
				HandFee:   feeInfo.HandFee,
				Topic:     topic,
				Qid:       qid,
				Data:      data,
				IsSuccess: true,
				Rid:       response.Rid,
				Error:     "",
			}
			_err := middleware.DB.Model(&PushQueueRecord{}).Create(&pushQueueRecord).Error
			if _err != nil {
				btlLog.PushQueue.Error("Create _err:\n%v\nPQR:\n%v", _err, utils.ValueJsonString(pushQueueRecord))
			}
			_err = middleware.DB.Model(&models.NftPresale{}).
				Where("id = ?", feeInfo.ID).
				Update("is_pushed_queue", true).Error
			if _err != nil {
				btlLog.PushQueue.Error("Update NftPresale _err:\n%v\nid:\n%v", _err, feeInfo.ID)
			}
		}
	}
	return
}
