package api

import (
	"github.com/lightningnetwork/lnd/lnrpc"
)

func ListChainTxnsAndGetResponse() (*lnrpc.TransactionDetails, error) {
	return listChainTxns()
}

func GetListChainTransactions() (*[]ChainTransaction, error) {
	response, err := listChainTxns()
	if err != nil {
		//utils.LogError("", err)
		return nil, err
	}
	result := processChainTransactions(response)
	return result, nil
}
