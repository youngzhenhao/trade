package api

import (
	"context"
	"github.com/lightningnetwork/lnd/lnrpc"
	"strconv"
	"trade/config"
	"trade/utils"
)

type ClientType int

var (
	ClientTypeLnd  ClientType = 1
	ClientTypeTapd ClientType = 2
	ClientTypeLitd ClientType = 3
)

type ConnConfiguration struct {
	GrpcHost     string
	TlsCertPath  string
	MacaroonPath string
}

func GetConnConfiguration(clientType ClientType) *ConnConfiguration {
	var connConfiguration ConnConfiguration
	if clientType == ClientTypeLnd {
		connConfiguration.GrpcHost = config.GetLoadConfig().ApiConfig.Lnd.Host + ":" + strconv.Itoa(config.GetLoadConfig().ApiConfig.Lnd.Port)
		connConfiguration.TlsCertPath = config.GetLoadConfig().ApiConfig.Lnd.TlsCertPath
		connConfiguration.MacaroonPath = config.GetLoadConfig().ApiConfig.Lnd.MacaroonPath
	} else if clientType == ClientTypeTapd {
		connConfiguration.GrpcHost = config.GetLoadConfig().ApiConfig.Tapd.Host + ":" + strconv.Itoa(config.GetLoadConfig().ApiConfig.Tapd.Port)
		connConfiguration.TlsCertPath = config.GetLoadConfig().ApiConfig.Tapd.TlsCertPath
		connConfiguration.MacaroonPath = config.GetLoadConfig().ApiConfig.Tapd.MacaroonPath
	} else if clientType == ClientTypeLitd {
		connConfiguration.GrpcHost = config.GetLoadConfig().ApiConfig.Litd.Host + ":" + strconv.Itoa(config.GetLoadConfig().ApiConfig.Litd.Port)
		connConfiguration.TlsCertPath = config.GetLoadConfig().ApiConfig.Litd.TlsCertPath
		connConfiguration.MacaroonPath = config.GetLoadConfig().ApiConfig.Litd.MacaroonPath
	} else {
		return nil
	}
	return &connConfiguration
}

func listChainTxns() (*lnrpc.TransactionDetails, error) {
	connConfiguration := GetConnConfiguration(ClientTypeLnd)
	conn, connClose := utils.GetConn(connConfiguration.GrpcHost, connConfiguration.TlsCertPath, connConfiguration.MacaroonPath)
	defer connClose()
	client := lnrpc.NewLightningClient(conn)
	request := &lnrpc.GetTransactionsRequest{
		//StartHeight: 0,
		//EndHeight:   0,
		//Account:     "",
	}
	response, err := client.GetTransactions(context.Background(), request)
	if err != nil {
		return nil, utils.AppendErrorInfo(err, "GetTransactions")
	}
	return response, nil
}

type ChainTransaction struct {
	TxHash            string             `json:"tx_hash"`
	Amount            int                `json:"amount"`
	NumConfirmations  int                `json:"num_confirmations"`
	BlockHash         string             `json:"block_hash"`
	BlockHeight       int                `json:"block_height"`
	TimeStamp         int                `json:"time_stamp"`
	TotalFees         int                `json:"total_fees"`
	DestAddresses     []string           `json:"dest_addresses"`
	OutputDetails     []OutputDetail     `json:"output_details"`
	RawTxHex          string             `json:"raw_tx_hex"`
	Label             string             `json:"label"`
	PreviousOutpoints []PreviousOutpoint `json:"previous_outpoints"`
}

type OutputDetail struct {
	OutputType   string `json:"output_type"`
	Address      string `json:"address"`
	PkScript     string `json:"pk_script"`
	OutputIndex  int    `json:"output_index"`
	Amount       int    `json:"amount"`
	IsOurAddress bool   `json:"is_our_address"`
}

type PreviousOutpoint struct {
	Outpoint    string `json:"outpoint"`
	IsOurOutput bool   `json:"is_our_output"`
}

func processChainTransactions(response *lnrpc.TransactionDetails) *[]ChainTransaction {
	var chainTransactions []ChainTransaction
	for _, transaction := range response.Transactions {
		var outputDetails []OutputDetail
		for _, outputDetail := range transaction.OutputDetails {
			outputDetails = append(outputDetails, OutputDetail{
				OutputType:   outputDetail.OutputType.String(),
				Address:      outputDetail.Address,
				PkScript:     outputDetail.PkScript,
				OutputIndex:  int(outputDetail.OutputIndex),
				Amount:       int(outputDetail.Amount),
				IsOurAddress: outputDetail.IsOurAddress,
			})
		}
		var previousOutpoints []PreviousOutpoint
		for _, previousOutpoint := range transaction.PreviousOutpoints {
			previousOutpoints = append(previousOutpoints, PreviousOutpoint{
				Outpoint:    previousOutpoint.Outpoint,
				IsOurOutput: previousOutpoint.IsOurOutput,
			})
		}
		chainTransactions = append(chainTransactions, ChainTransaction{
			TxHash:            transaction.TxHash,
			Amount:            int(transaction.Amount),
			NumConfirmations:  int(transaction.NumConfirmations),
			BlockHash:         transaction.BlockHash,
			BlockHeight:       int(transaction.BlockHeight),
			TimeStamp:         int(transaction.TimeStamp),
			TotalFees:         int(transaction.TotalFees),
			DestAddresses:     transaction.GetDestAddresses(),
			OutputDetails:     outputDetails,
			RawTxHex:          transaction.RawTxHex,
			Label:             transaction.Label,
			PreviousOutpoints: previousOutpoints,
		})
	}
	return &chainTransactions
}
