package api

import (
	"fmt"
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/rpcclient"
	"trade/config"
)

func getBitcoinConnConfig() *rpcclient.ConnConfig {
	return &rpcclient.ConnConfig{
		Host:         fmt.Sprintf("%s:%d", config.GetLoadConfig().ApiConfig.Bitcoind.Host, config.GetLoadConfig().ApiConfig.Bitcoind.Port),
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
		fmt.Println("Error connecting to the RPC server:", err)
		return
	}
	defer client.Shutdown()
	feeResult, err = client.EstimateSmartFee(confTarget, mode)
	if err != nil {
		fmt.Println("Error calling EstimateSmartFeeAndGetResult:", err)
		return
	}
	return feeResult, nil
}
