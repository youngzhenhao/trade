package cashbackQueue

import (
	"bytes"
	"encoding/json"
	"io"
	"math/rand"
	"net/http"
	"time"
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

type request struct {
	Data string `json:"data"`
}

type Response struct {
	Qid string `json:"qid"`
	Rid string `json:"rid"`
}

type FeeInfo struct {
	ID       int    `json:"id"`
	NpubKey  string `json:"npub_key"`
	AssetsID string `json:"assets_id"`
	HandFee  int    `json:"hand_fee"`
}

func Push(topic queueTopic, qid string, request request) (Response, error) {
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
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			return
		}
	}(res.Body)
	body, err := io.ReadAll(res.Body)
	return body, nil
}

func PushClaimAsset(info FeeInfo) (Response, error) {
	topic := claimAsset
	seed := rand.New(rand.NewSource(time.Now().UnixNano()))
	qid := utils.RandString(seed, 64)
	data, err := json.Marshal(info)
	if err != nil {
		return Response{}, err
	}
	request := request{
		Data: string(data),
	}
	return Push(topic, qid, request)
}

func PushPurchasePresaleNFT(info FeeInfo) (Response, error) {
	topic := purchasePresaleNFT
	seed := rand.New(rand.NewSource(time.Now().UnixNano()))
	qid := utils.RandString(seed, 64)
	data, err := json.Marshal(info)
	if err != nil {
		return Response{}, err
	}
	request := request{
		Data: string(data),
	}
	return Push(topic, qid, request)
}
