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
	"trade/models"
	"trade/utils"
)

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

func getBitcoinConnConfig(network models.Network) (*rpcclient.ConnConfig, error) {
	var ip string
	var port int
	var wallet string
	var host string
	var user string
	var pass string
	var httpPostMode bool
	var disableTLS bool
	switch network {
	case models.Mainnet:
		ip = config.GetLoadConfig().ApiConfig.Bitcoind.Mainnet.Ip
		port = config.GetLoadConfig().ApiConfig.Bitcoind.Mainnet.Port
		wallet = config.GetLoadConfig().ApiConfig.Bitcoind.Mainnet.Wallet
		user = config.GetLoadConfig().ApiConfig.Bitcoind.Mainnet.RpcUser
		pass = config.GetLoadConfig().ApiConfig.Bitcoind.Mainnet.RpcPasswd
		httpPostMode = config.GetLoadConfig().ApiConfig.Bitcoind.Mainnet.HttpPostMode
		disableTLS = config.GetLoadConfig().ApiConfig.Bitcoind.Mainnet.DisableTLS
	case models.Testnet:
		ip = config.GetLoadConfig().ApiConfig.Bitcoind.Testnet.Ip
		port = config.GetLoadConfig().ApiConfig.Bitcoind.Testnet.Port
		wallet = config.GetLoadConfig().ApiConfig.Bitcoind.Testnet.Wallet
		user = config.GetLoadConfig().ApiConfig.Bitcoind.Testnet.RpcUser
		pass = config.GetLoadConfig().ApiConfig.Bitcoind.Testnet.RpcPasswd
		httpPostMode = config.GetLoadConfig().ApiConfig.Bitcoind.Testnet.HttpPostMode
		disableTLS = config.GetLoadConfig().ApiConfig.Bitcoind.Testnet.DisableTLS
	case models.Regtest:
		ip = config.GetLoadConfig().ApiConfig.Bitcoind.Regtest.Ip
		port = config.GetLoadConfig().ApiConfig.Bitcoind.Regtest.Port
		wallet = config.GetLoadConfig().ApiConfig.Bitcoind.Regtest.Wallet
		user = config.GetLoadConfig().ApiConfig.Bitcoind.Regtest.RpcUser
		pass = config.GetLoadConfig().ApiConfig.Bitcoind.Regtest.RpcPasswd
		httpPostMode = config.GetLoadConfig().ApiConfig.Bitcoind.Regtest.HttpPostMode
		disableTLS = config.GetLoadConfig().ApiConfig.Bitcoind.Regtest.DisableTLS
	default:
		return nil, errors.New("invalid api network")
	}
	if wallet == "" {
		host = fmt.Sprintf("%s:%d", ip, port)
	} else {
		host = fmt.Sprintf("%s:%d/wallet/%s", ip, port, wallet)
	}
	connConfig := rpcclient.ConnConfig{
		Host:         host,
		User:         user,
		Pass:         pass,
		HTTPPostMode: httpPostMode,
		DisableTLS:   disableTLS,
	}
	return &connConfig, nil
}

func estimateSmartFee(network models.Network, confTarget int64, mode *btcjson.EstimateSmartFeeMode) (feeResult *btcjson.EstimateSmartFeeResult, err error) {
	connCfg, err := getBitcoinConnConfig(network)
	if err != nil {
		return nil, err
	}
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

func getRawTransaction(network models.Network, txid string) (transaction *btcutil.Tx, err error) {
	connCfg, err := getBitcoinConnConfig(network)
	if err != nil {
		return nil, err
	}
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

func getTransaction(network models.Network, txid string) (transaction *btcjson.GetTransactionResult, err error) {
	connCfg, err := getBitcoinConnConfig(network)
	if err != nil {
		return nil, err
	}
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

func decodeScript(network models.Network, encodedPubKeyScript string) (transaction *btcjson.DecodeScriptResult, err error) {
	connCfg, err := getBitcoinConnConfig(network)
	if err != nil {
		return nil, err
	}
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

type PostGetRawTransactionResponse struct {
	Result *PostGetRawTransactionResult `json:"result"`
	Error  *BitcoindRpcResponseError    `json:"error"`
	ID     string                       `json:"id"`
}

type BitcoindRpcResponseError struct {
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

func getUri(network models.Network) (string, error) {
	var user string
	var password string
	var ip string
	var port int
	var wallet string
	var uri string
	switch network {
	case models.Mainnet:
		user = config.GetLoadConfig().ApiConfig.Bitcoind.Mainnet.RpcUser
		password = config.GetLoadConfig().ApiConfig.Bitcoind.Mainnet.RpcPasswd
		ip = config.GetLoadConfig().ApiConfig.Bitcoind.Mainnet.Ip
		port = config.GetLoadConfig().ApiConfig.Bitcoind.Mainnet.Port
		wallet = config.GetLoadConfig().ApiConfig.Bitcoind.Mainnet.Wallet
	case models.Testnet:
		user = config.GetLoadConfig().ApiConfig.Bitcoind.Testnet.RpcUser
		password = config.GetLoadConfig().ApiConfig.Bitcoind.Testnet.RpcPasswd
		ip = config.GetLoadConfig().ApiConfig.Bitcoind.Testnet.Ip
		port = config.GetLoadConfig().ApiConfig.Bitcoind.Testnet.Port
		wallet = config.GetLoadConfig().ApiConfig.Bitcoind.Testnet.Wallet
	case models.Regtest:
		user = config.GetLoadConfig().ApiConfig.Bitcoind.Regtest.RpcUser
		password = config.GetLoadConfig().ApiConfig.Bitcoind.Regtest.RpcPasswd
		ip = config.GetLoadConfig().ApiConfig.Bitcoind.Regtest.Ip
		port = config.GetLoadConfig().ApiConfig.Bitcoind.Regtest.Port
		wallet = config.GetLoadConfig().ApiConfig.Bitcoind.Regtest.Wallet
	default:
		return "", errors.New("invalid network")
	}
	if wallet == "" {
		uri = fmt.Sprintf("http://%s:%s@%s:%d", user, password, ip, port)
	} else {
		uri = fmt.Sprintf("http://%s:%s@%s:%d/wallet/%s", user, password, ip, port, wallet)
	}
	return uri, nil
}

func postGetRawTransaction(network models.Network, txid string, verbosity Verbosity) (result *PostGetRawTransactionResult, err error) {
	uri, err := getUri(network)
	if err != nil {
		return nil, err
	}
	requestBodyRaw := fmt.Sprintf(`{"jsonrpc":"1.0","id":"%s","method":"getrawtransaction","params":["%s",%d]}`, txid, txid, verbosity)
	payload := strings.NewReader(requestBodyRaw)
	req, err := http.NewRequest("POST", uri, payload)
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

func postCallBitcoindToGetRawTransaction(network models.Network, requestBodyRaw string) (*[]PostGetRawTransactionResponse, error) {
	uri, err := getUri(network)
	if err != nil {
		return nil, err
	}
	payload := strings.NewReader(requestBodyRaw)
	req, err := http.NewRequest("POST", uri, payload)
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
	var response []PostGetRawTransactionResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		var responseSingle PostGetRawTransactionResponse
		err = json.Unmarshal(body, &responseSingle)
		if err != nil {
			return nil, err
		}
		response = append(response, responseSingle)
	}
	return &response, nil
}

func postCallBitcoindToDecodeRawTransaction(network models.Network, requestBodyRaw string) (*[]PostDecodeRawTransactionResponse, error) {
	uri, err := getUri(network)
	if err != nil {
		return nil, err
	}
	payload := strings.NewReader(requestBodyRaw)
	req, err := http.NewRequest("POST", uri, payload)
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
	var response []PostDecodeRawTransactionResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		var responseSingle PostDecodeRawTransactionResponse
		err = json.Unmarshal(body, &responseSingle)
		if err != nil {
			return nil, err
		}
		response = append(response, responseSingle)
	}
	return &response, nil
}

func outpointSliceToRequestBodyRawString(outpoints []string, verbosity Verbosity) (request string) {
	request = "["
	for i, outpoint := range outpoints {
		txid, indexStr := utils.OutpointToTransactionAndIndex(outpoint)
		if txid == "" || indexStr == "" {
			continue
		}
		_, err := strconv.Atoi(indexStr)
		if err != nil {
			continue
		}
		element := fmt.Sprintf("{\"jsonrpc\":\"1.0\",\"id\":\"%s\",\"method\":\"getrawtransaction\",\"params\":[\"%s\",%d]}", outpoint, txid, verbosity)
		request += element
		if i != len(outpoints)-1 {
			request += ","
		}
	}
	request += "]"
	return request
}

func postGetRawTransactions(network models.Network, outpoints []string, verbosity Verbosity) (*[]PostGetRawTransactionResponse, error) {
	requestBodyRaw := outpointSliceToRequestBodyRawString(outpoints, verbosity)
	response, err := postCallBitcoindToGetRawTransaction(network, requestBodyRaw)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func rawTransactionHexSliceToRequestBodyRawString(rawTransactions []string) (request string) {
	request = "["
	for i, transaction := range rawTransactions {
		element := fmt.Sprintf("{\"jsonrpc\":\"1.0\",\"id\":\"%s\",\"method\":\"decoderawtransaction\",\"params\":[\"%s\"]}", transaction, transaction)
		request += element
		if i != len(rawTransactions)-1 {
			request += ","
		}
	}
	request += "]"
	return request
}

func postDecodeRawTransactions(network models.Network, rawTransactions []string) (*[]PostDecodeRawTransactionResponse, error) {
	requestBodyRaw := rawTransactionHexSliceToRequestBodyRawString(rawTransactions)
	response, err := postCallBitcoindToDecodeRawTransaction(network, requestBodyRaw)
	if err != nil {
		return nil, err
	}
	return response, nil
}

type PostDecodeRawTransactionResponse struct {
	Result *PostDecodeRawTransactionResult `json:"result"`
	Error  *BitcoindRpcResponseError       `json:"error"`
	ID     string                          `json:"id"`
}

type PostDecodeRawTransactionResult struct {
	Txid     string                           `json:"txid"`
	Hash     string                           `json:"hash"`
	Version  int                              `json:"version"`
	Size     int                              `json:"size"`
	Vsize    int                              `json:"vsize"`
	Weight   int                              `json:"weight"`
	Locktime int                              `json:"locktime"`
	Vin      []DecodeRawTransactionResultVin  `json:"vin"`
	Vout     []DecodeRawTransactionResultVout `json:"vout"`
}

type DecodeRawTransactionResultVin struct {
	Txid        string                                 `json:"txid"`
	Vout        int                                    `json:"vout"`
	ScriptSig   DecodeRawTransactionResultVinScriptSig `json:"scriptSig"`
	Txinwitness []string                               `json:"txinwitness"`
	Sequence    int64                                  `json:"sequence"`
}

type DecodeRawTransactionResultVinScriptSig struct {
	Asm string `json:"asm"`
	Hex string `json:"hex"`
}

type DecodeRawTransactionResultVout struct {
	Value        float64                                    `json:"value"`
	N            int                                        `json:"n"`
	ScriptPubKey DecodeRawTransactionResultVoutScriptPubKey `json:"scriptPubKey"`
}

type DecodeRawTransactionResultVoutScriptPubKey struct {
	Asm     string `json:"asm"`
	Desc    string `json:"desc"`
	Hex     string `json:"hex"`
	Address string `json:"address"`
	Type    string `json:"type"`
}

func postDecodeRawTransaction(network models.Network, transaction string) (*PostDecodeRawTransactionResponse, error) {
	uri, err := getUri(network)
	if err != nil {
		return nil, err
	}
	request := fmt.Sprintf("{\"jsonrpc\":\"1.0\",\"id\":\"%s\",\"method\":\"decoderawtransaction\",\"params\":[\"%s\"]}", transaction, transaction)
	payload := strings.NewReader(request)
	req, err := http.NewRequest("POST", uri, payload)
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
	var response PostDecodeRawTransactionResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

func txidsToRequestBodyRawString(txids []string, verbosity Verbosity) (request string) {
	request = "["
	for i, txid := range txids {
		if txid == "" {
			continue
		}
		element := fmt.Sprintf("{\"jsonrpc\":\"1.0\",\"id\":\"%s\",\"method\":\"getrawtransaction\",\"params\":[\"%s\",%d]}", txid, txid, verbosity)
		request += element
		if i != len(txids)-1 {
			request += ","
		}
	}
	request += "]"
	return request
}
