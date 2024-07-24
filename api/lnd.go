package api

import (
	"errors"
	"github.com/lightningnetwork/lnd/lnrpc"
	"strconv"
	"trade/utils"
)

func ListChainTxnsAndGetResponse() (*lnrpc.TransactionDetails, error) {
	return listChainTxns()
}

func GetListChainTransactions() (*[]ChainTransaction, error) {
	response, err := listChainTxns()
	if err != nil {
		return nil, utils.AppendErrorInfo(err, "listChainTxns")
	}
	result := processChainTransactions(response)
	return result, nil
}

func WalletBalanceAndGetResponse() (*lnrpc.WalletBalanceResponse, error) {
	return walletBalance()
}

func GetListChainTransactionsOutpointAddress(outpoint string) (address string, err error) {
	response, err := GetListChainTransactions()
	if err != nil {
		return "", utils.AppendErrorInfo(err, "GetListChainTransactions")
	}
	tx, indexStr := utils.GetTransactionAndIndexByOutpoint(outpoint)
	index, err := strconv.Atoi(indexStr)
	if err != nil {
		return "", utils.AppendErrorInfo(err, "GetTransactionAndIndexByOutpoint")
	}
	for _, transaction := range *response {
		if transaction.TxHash == tx {
			return transaction.DestAddresses[index], nil
		}
	}
	err = errors.New("did not match transaction outpoint")
	return "", err
}
