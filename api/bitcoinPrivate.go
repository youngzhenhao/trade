package api

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/rpcclient"
	"io"
	"net/http"
	"strconv"
	"strings"
	"trade/config"
)

func getBitcoinConnConfig() *rpcclient.ConnConfig {
	ip := config.GetLoadConfig().ApiConfig.Bitcoind.Ip
	port := config.GetLoadConfig().ApiConfig.Bitcoind.Port
	wallet := config.GetLoadConfig().ApiConfig.Bitcoind.Wallet
	host := fmt.Sprintf("%s:%d/wallet/%s", ip, port, wallet)
	return &rpcclient.ConnConfig{
		Host:         host,
		User:         config.GetLoadConfig().ApiConfig.Bitcoind.RpcUser,
		Pass:         config.GetLoadConfig().ApiConfig.Bitcoind.RpcPasswd,
		HTTPPostMode: config.GetLoadConfig().ApiConfig.Bitcoind.HTTPPostMode,
		DisableTLS:   config.GetLoadConfig().ApiConfig.Bitcoind.DisableTLS,
	}
}

func estimateSmartFee(confTarget int64, mode *btcjson.EstimateSmartFeeMode) (feeResult *btcjson.EstimateSmartFeeResult, err error) {
	connCfg := getBitcoinConnConfig()
	client, err := rpcclient.New(connCfg, nil)
	if err != nil {
		return
	}
	defer client.Shutdown()
	feeResult, err = client.EstimateSmartFee(confTarget, mode)
	if err != nil {
		return
	}
	return feeResult, nil
}

func getRawTransaction(txid string) (transaction *btcutil.Tx, err error) {
	connCfg := getBitcoinConnConfig()
	client, err := rpcclient.New(connCfg, nil)
	if err != nil {
		return
	}
	defer client.Shutdown()
	var blockHash chainhash.Hash
	err = chainhash.Decode(&blockHash, txid)
	if err != nil {
		return nil, err
	}
	var response *btcutil.Tx
	response, err = client.GetRawTransaction(&blockHash)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func getTransaction(txid string) (transaction *btcjson.GetTransactionResult, err error) {
	connCfg := getBitcoinConnConfig()
	client, err := rpcclient.New(connCfg, nil)
	if err != nil {
		return
	}
	defer client.Shutdown()
	var blockHash chainhash.Hash
	err = chainhash.Decode(&blockHash, txid)
	if err != nil {
		return nil, err
	}
	var response *btcjson.GetTransactionResult
	response, err = client.GetTransaction(&blockHash)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func decodeScript(encodedPubKeyScript string) (transaction *btcjson.DecodeScriptResult, err error) {
	connCfg := getBitcoinConnConfig()
	client, err := rpcclient.New(connCfg, nil)
	if err != nil {
		return
	}
	defer client.Shutdown()
	var response *btcjson.DecodeScriptResult
	decodeString, err := hex.DecodeString(encodedPubKeyScript)
	if err != nil {
		return nil, err
	}
	response, err = client.DecodeScript(decodeString)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func getRawTransactionAndDecodeOutputScript(txid string) (err error) {
	//connCfg := getBitcoinConnConfig()
	//client, err := rpcclient.New(connCfg, nil)
	//if err != nil {
	//	return
	//}
	//defer client.Shutdown()
	//var blockHash chainhash.Hash
	//err = chainhash.Decode(&blockHash, txid)
	//if err != nil {
	//	return nil, err
	//}
	//var response *btcjson.GetTransactionResult
	//response, err = client.GetTransaction(&blockHash)
	//if err != nil {
	//	return nil, err
	//}
	//return response, nil
	return nil
}

type GetRawTransactionResponse struct {
	Txid          string            `json:"txid"`
	Hash          string            `json:"hash"`
	Version       int               `json:"version"`
	Size          int               `json:"size"`
	Vsize         int               `json:"vsize"`
	Weight        int               `json:"weight"`
	Locktime      int               `json:"locktime"`
	Vin           []TransactionVin  `json:"vin"`
	Vout          []TransactionVout `json:"vout"`
	Fee           float64           `json:"fee"`
	Hex           string            `json:"hex"`
	Blockhash     string            `json:"blockhash"`
	Confirmations int               `json:"confirmations"`
	Time          int               `json:"time"`
	Blocktime     int               `json:"blocktime"`
}

type TransactionVinScriptSig struct {
	Asm string `json:"asm"`
	Hex string `json:"hex"`
}

type TransactionVinPrevoutScriptPubKey struct {
	Asm     string `json:"asm"`
	Desc    string `json:"desc"`
	Hex     string `json:"hex"`
	Address string `json:"address"`
	Type    string `json:"type"`
}

type TransactionVinPrevout struct {
	Generated    bool                              `json:"generated"`
	Height       int                               `json:"height"`
	Value        float64                           `json:"value"`
	ScriptPubKey TransactionVinPrevoutScriptPubKey `json:"scriptPubKey"`
}

type TransactionVin struct {
	Txid        string                  `json:"txid"`
	Vout        int                     `json:"vout"`
	ScriptSig   TransactionVinScriptSig `json:"scriptSig"`
	Txinwitness []string                `json:"txinwitness"`
	Prevout     TransactionVinPrevout   `json:"prevout"`
	Sequence    int                     `json:"sequence"`
}

type TransactionVoutScriptPubKey struct {
	Asm     string `json:"asm"`
	Desc    string `json:"desc"`
	Hex     string `json:"hex"`
	Address string `json:"address"`
	Type    string `json:"type"`
}

type TransactionVout struct {
	Value        float64                     `json:"value"`
	N            int                         `json:"n"`
	ScriptPubKey TransactionVoutScriptPubKey `json:"scriptPubKey"`
}

type GetRawTransactionResponseMsgTx struct {
	Version  int                   `json:"Version"`
	TxIn     []TransactionMsgTxIn  `json:"TxIn"`
	TxOut    []TransactionMsgTxOut `json:"TxOut"`
	LockTime int                   `json:"LockTime"`
}

type TransactionMsgTxIn struct {
	PreviousOutPoint TransactionMsgTxInPreviousOutPoint `json:"PreviousOutPoint"`
	SignatureScript  string                             `json:"SignatureScript"`
	Witness          []string                           `json:"Witness"`
	Sequence         int                                `json:"Sequence"`
}

type TransactionMsgTxInPreviousOutPoint struct {
	Hash  string `json:"Hash"`
	Index int    `json:"Index"`
}

type TransactionMsgTxOut struct {
	Value    int    `json:"Value"`
	PkScript string `json:"PkScript"`
}

func GetUri() string {
	user := config.GetLoadConfig().ApiConfig.Bitcoind.RpcUser
	password := config.GetLoadConfig().ApiConfig.Bitcoind.RpcPasswd
	ip := config.GetLoadConfig().ApiConfig.Bitcoind.Ip
	port := config.GetLoadConfig().ApiConfig.Bitcoind.Port
	wallet := config.GetLoadConfig().ApiConfig.Bitcoind.Wallet
	url := fmt.Sprintf("http://%s:%s@%s:%d/wallet/%s", user, password, ip, port, wallet)
	return url
}

func postGetRawTransaction(txid string, verbosity Verbosity) (result *PostGetRawTransactionResult, err error) {
	url := GetUri()
	requestBodyRaw := fmt.Sprintf(`{"jsonrpc":"1.0","id":1,"method":"getrawtransaction","params":["%s",%d]}`, txid, verbosity)
	payload := strings.NewReader(requestBodyRaw)
	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		return nil, err
	}
	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(res.Body)
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var response PostGetRawTransactionResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}
	if response.Error != nil {
		return nil, errors.New(strconv.Itoa(response.Error.Code) + response.Error.Message)
	}
	return response.Result, nil
}

type PostGetRawTransactionResponse struct {
	Result *PostGetRawTransactionResult `json:"result"`
	Error  *PostGetRawTransactionError  `json:"error"`
	ID     int                          `json:"id"`
}

type PostGetRawTransactionError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type PostGetRawTransactionResult struct {
	Txid          string                     `json:"txid"`
	Hash          string                     `json:"hash"`
	Version       int                        `json:"version"`
	Size          int                        `json:"size"`
	Vsize         int                        `json:"vsize"`
	Weight        int                        `json:"weight"`
	Locktime      int                        `json:"locktime"`
	Vin           []RawTransactionResultVin  `json:"vin"`
	Vout          []RawTransactionResultVout `json:"vout"`
	Fee           float64                    `json:"fee"`
	Hex           string                     `json:"hex"`
	Blockhash     string                     `json:"blockhash"`
	Confirmations int                        `json:"confirmations"`
	Time          int                        `json:"time"`
	Blocktime     int                        `json:"blocktime"`
}

type RawTransactionResultVin struct {
	Txid        string                           `json:"txid"`
	Vout        int                              `json:"vout"`
	ScriptSig   RawTransactionResultVinScriptSig `json:"scriptSig"`
	Txinwitness []string                         `json:"txinwitness"`
	Prevout     RawTransactionResultVinPrevout   `json:"prevout"`
	Sequence    int64                            `json:"sequence"`
}

type RawTransactionResultVinPrevout struct {
	Generated    bool                                       `json:"generated"`
	Height       int                                        `json:"height"`
	Value        float64                                    `json:"value"`
	ScriptPubKey RawTransactionResultVinPrevoutScriptPubKey `json:"scriptPubKey"`
}

type RawTransactionResultVinPrevoutScriptPubKey struct {
	Asm     string `json:"asm"`
	Desc    string `json:"desc"`
	Hex     string `json:"hex"`
	Address string `json:"address"`
	Type    string `json:"type"`
}

type RawTransactionResultVinScriptSig struct {
	Asm string `json:"asm"`
	Hex string `json:"hex"`
}

type RawTransactionResultVout struct {
	Value        float64                              `json:"value"`
	N            int                                  `json:"n"`
	ScriptPubKey RawTransactionResultVoutScriptPubKey `json:"scriptPubKey"`
}

type RawTransactionResultVoutScriptPubKey struct {
	Asm     string `json:"asm"`
	Desc    string `json:"desc"`
	Hex     string `json:"hex"`
	Address string `json:"address"`
	Type    string `json:"type"`
}
