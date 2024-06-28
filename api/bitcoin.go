package api

import (
	"encoding/hex"
	"fmt"
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/wire"
	"strconv"
	"trade/models"
)

type Verbosity int

const (
	VerbosityHexEncoded Verbosity = iota
	VerbosityJson
	VerbosityJsonWithFeeAndPrevout
)

func EstimateSmartFeeAndGetResult(network models.Network, blocks int) (feeResult *btcjson.EstimateSmartFeeResult, err error) {
	return estimateSmartFee(network, int64(blocks), &btcjson.EstimateModeUnset)
}

func GetTransaction(network models.Network, txid string) (*btcjson.GetTransactionResult, error) {
	response, err := getTransaction(network, txid)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func GetRawTransactionMsgTx(network models.Network, txid string) (*GetRawTransactionResponseMsgTx, error) {
	response, err := getRawTransaction(network, txid)
	if err != nil {
		return nil, err
	}
	rawTransactionMsgTx := ProcessRawTransactionMsgTx(response.MsgTx())
	return rawTransactionMsgTx, nil
}

func ProcessRawTransactionMsgTx(transaction *wire.MsgTx) *GetRawTransactionResponseMsgTx {
	var txin []TransactionMsgTxIn
	var txout []TransactionMsgTxOut
	for _, in := range transaction.TxIn {
		var witness []string
		for _, wit := range in.Witness {
			witness = append(witness, hex.EncodeToString(wit))
		}
		txin = append(txin, TransactionMsgTxIn{
			PreviousOutPoint: TransactionMsgTxInPreviousOutPoint{
				Hash:  in.PreviousOutPoint.Hash.String(),
				Index: int(in.PreviousOutPoint.Index),
			},
			SignatureScript: hex.EncodeToString(in.SignatureScript),
			Witness:         witness,
			Sequence:        int(in.Sequence),
		})
	}
	for _, out := range transaction.TxOut {
		txout = append(txout, TransactionMsgTxOut{
			Value:    int(out.Value),
			PkScript: hex.EncodeToString(out.PkScript),
		})
	}
	result := GetRawTransactionResponseMsgTx{
		Version:  int(transaction.Version),
		TxIn:     txin,
		TxOut:    txout,
		LockTime: int(transaction.LockTime),
	}
	return &result
}

func DecodeScript(network models.Network, encodedPubKeyScript string) (transaction *btcjson.DecodeScriptResult, err error) {
	return decodeScript(network, encodedPubKeyScript)
}

func PostGetRawTransaction(network models.Network, txid string, verbosity int) (result *PostGetRawTransactionResult, err error) {
	if verbosity == int(VerbosityJson) {
		return PostGetRawTransactionWithoutFeeAndPrevout(network, txid)
	} else if verbosity == int(VerbosityJsonWithFeeAndPrevout) {
		return PostGetRawTransactionWithFeeAndPrevout(network, txid)
	} else {
		return nil, fmt.Errorf("invalid verbosity: %d", verbosity)
	}
}

func PostGetRawTransactionWithoutFeeAndPrevout(network models.Network, txid string) (result *PostGetRawTransactionResult, err error) {
	return postGetRawTransaction(network, txid, VerbosityJson)
}

func PostGetRawTransactionWithFeeAndPrevout(network models.Network, txid string) (result *PostGetRawTransactionResult, err error) {
	return postGetRawTransaction(network, txid, VerbosityJsonWithFeeAndPrevout)
}

func GetAddressByTxidAndIndex(network models.Network, txid string, index int) (address string, err error) {
	response, err := PostGetRawTransactionWithFeeAndPrevout(network, txid)
	if err != nil {
		return "", err
	}
	vout := (*response).Vout
	if !(len(vout) > index) {
		return "", fmt.Errorf("invalid index: %d", index)
	}
	return vout[index].ScriptPubKey.Address, nil
}

func GetTransactionByTxid(network models.Network, txid string) (transaction *PostGetRawTransactionResult, err error) {
	response, err := PostGetRawTransactionWithFeeAndPrevout(network, txid)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func GetAddressByOutpoint(network models.Network, outpoint string) (address string, err error) {
	txid, indexStr := OutpointToTransactionAndIndex(outpoint)
	if txid == "" || indexStr == "" {
		return "", fmt.Errorf("invalid outpoint: %s", outpoint)
	}
	index, err := strconv.Atoi(indexStr)
	if err != nil {
		return "", err
	}
	return GetAddressByTxidAndIndex(network, txid, index)
}

func GetTransactionByOutpoint(network models.Network, outpoint string) (transaction *PostGetRawTransactionResult, err error) {
	txid, indexStr := OutpointToTransactionAndIndex(outpoint)
	if txid == "" || indexStr == "" {
		return nil, fmt.Errorf("invalid outpoint: %s", outpoint)
	}
	_, err = strconv.Atoi(indexStr)
	if err != nil {
		return nil, err
	}
	return GetTransactionByTxid(network, txid)
}

func GetAddressesByOutpointSlice(network models.Network, outpoints []string) (addresses map[string]string, err error) {
	return GetAddressesBatchProcess(network, outpoints)
}

func OutpointsToTxidsAndIndexes(outpoints []string) (txidIndex map[string]int) {
	txidIndex = make(map[string]int)
	for _, outpoint := range outpoints {
		txid, indexStr := OutpointToTransactionAndIndex(outpoint)
		if txid == "" || indexStr == "" {
			continue
		}
		index, err := strconv.Atoi(indexStr)
		if err != nil {
			continue
		}
		txidIndex[txid] = index
	}
	return txidIndex
}

func OutpointsToTxids(outpoints []string) (txids []string) {
	for _, outpoint := range outpoints {
		txid, indexStr := OutpointToTransactionAndIndex(outpoint)
		if txid == "" || indexStr == "" {
			continue
		}
		_, err := strconv.Atoi(indexStr)
		if err != nil {
			continue
		}
		txids = append(txids, txid)
	}
	return txids
}

func PostGetRawTransactionsWithFeeAndPrevoutByOutpoints(network models.Network, outpoints []string) (*[]PostGetRawTransactionResponse, error) {
	response, err := postGetRawTransactions(network, outpoints, VerbosityJsonWithFeeAndPrevout)
	if err != nil {
		return nil, err
	}
	result := ProcessGetRawTransactions(response)
	return result, nil
}

func ProcessGetRawTransactions(response *[]PostGetRawTransactionResponse) *[]PostGetRawTransactionResponse {
	var result []PostGetRawTransactionResponse
	for _, transaction := range *response {
		if transaction.Result == nil || transaction.Error != nil {
			continue
		}
		result = append(result, transaction)
	}
	return &result
}

func GetAddressesBatchProcess(network models.Network, outpoints []string) (outpointAddress map[string]string, err error) {
	outpointAddress = make(map[string]string)
	response, err := PostGetRawTransactionsWithFeeAndPrevoutByOutpoints(network, outpoints)
	if err != nil {
		return nil, err
	}
	for _, transaction := range *response {
		var index int
		if transaction.Result == nil || transaction.Error != nil {
			continue
		}
		outpoint := transaction.ID
		txid, indexStr := OutpointToTransactionAndIndex(outpoint)
		if txid == "" || indexStr == "" {
			continue
		}
		index, err = strconv.Atoi(indexStr)
		if err != nil {
			continue
		}
		out := transaction.Result.Vout
		if !(index < len(out)) {
			continue
		}
		outpointAddress[outpoint] = out[index].ScriptPubKey.Address
	}
	return outpointAddress, nil
}

func GetTransactionsBatchProcess(network models.Network, outpoints []string) (*[]PostGetRawTransactionResponse, error) {
	response, err := PostGetRawTransactionsWithFeeAndPrevoutByOutpoints(network, outpoints)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func GetTransactionsByOutpointSlice(network models.Network, outpoints []string) (*[]PostGetRawTransactionResponse, error) {
	return GetTransactionsBatchProcess(network, outpoints)
}
