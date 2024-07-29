package servicesrpc

import (
	"context"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/lightningnetwork/lnd/lnrpc/chainrpc"
	"strconv"
	"trade/config"
	"trade/utils"
)

func GetBlockInfo(hash string) (*chainrpc.GetBlockResponse, error) {

	blockHash, err := chainhash.NewHashFromStr(hash)
	if err != nil {
		return nil, err
	}
	request := &chainrpc.GetBlockRequest{BlockHash: blockHash.CloneBytes()}
	response, err := getBlockInfo(request)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func getBlockInfo(request *chainrpc.GetBlockRequest) (*chainrpc.GetBlockResponse, error) {
	lndconf := config.GetConfig().ApiConfig.Lnd

	grpcHost := lndconf.Host + ":" + strconv.Itoa(lndconf.Port)
	tlsCertPath := lndconf.TlsCertPath
	macaroonPath := lndconf.MacaroonPath

	conn, connClose := utils.GetConn(grpcHost, tlsCertPath, macaroonPath)
	defer connClose()

	client := chainrpc.NewChainKitClient(conn)
	response, err := client.GetBlock(context.Background(), request)
	return response, err
}
