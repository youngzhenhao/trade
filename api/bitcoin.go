package api

import (
	"encoding/hex"
	"fmt"
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/wire"
	"strconv"
)

type Verbosity int

const (
	VerbosityHexEncoded Verbosity = iota
	VerbosityJson
	VerbosityJsonWithFeeAndPrevout
)

func EstimateSmartFeeAndGetResult(blocks int) (feeResult *btcjson.EstimateSmartFeeResult, err error) {
	return estimateSmartFee(int64(blocks), &btcjson.EstimateModeUnset)
}

func GetTransaction(txid string) (*btcjson.GetTransactionResult, error) {
	response, err := getTransaction(txid)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func GetRawTransactionMsgTx(txid string) (*GetRawTransactionResponseMsgTx, error) {
	response, err := getRawTransaction(txid)
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

func DecodeScript(encodedPubKeyScript string) (transaction *btcjson.DecodeScriptResult, err error) {
	return decodeScript(encodedPubKeyScript)
}

func PostGetRawTransaction(txid string, verbosity int) (result *PostGetRawTransactionResult, err error) {
	if verbosity == int(VerbosityJson) {
		return PostGetRawTransactionWithoutFeeAndPrevout(txid)
	} else if verbosity == int(VerbosityJsonWithFeeAndPrevout) {
		return PostGetRawTransactionWithFeeAndPrevout(txid)
	} else {
		return nil, fmt.Errorf("invalid verbosity: %d", verbosity)
	}
}

func PostGetRawTransactionWithoutFeeAndPrevout(txid string) (result *PostGetRawTransactionResult, err error) {
	return postGetRawTransaction(txid, VerbosityJson)
}

func PostGetRawTransactionWithFeeAndPrevout(txid string) (result *PostGetRawTransactionResult, err error) {
	return postGetRawTransaction(txid, VerbosityJsonWithFeeAndPrevout)
}

func GetAddressByTxidAndIndex(txid string, index int) (address string, err error) {
	response, err := PostGetRawTransactionWithFeeAndPrevout(txid)
	if err != nil {
		return "", err
	}
	vout := (*response).Vout
	if !(len(vout) > index) {
		return "", fmt.Errorf("invalid index: %d", index)
	}
	return vout[index].ScriptPubKey.Address, nil
}

func GetAddressByOutpoint(outpoint string) (address string, err error) {
	txid, indexStr := OutpointToTransactionAndIndex(outpoint)
	if txid == "" || indexStr == "" {
		return "", fmt.Errorf("invalid outpoint: %s", outpoint)
	}
	index, err := strconv.Atoi(indexStr)
	if err != nil {
		return "", err
	}
	return GetAddressByTxidAndIndex(txid, index)
}
