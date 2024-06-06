package api

import (
	"github.com/lightningnetwork/lnd/lnrpc"
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
