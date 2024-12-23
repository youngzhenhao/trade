package servicesrpc

import (
	"context"
	"encoding/hex"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/lightningnetwork/lnd/lnrpc"
	"github.com/lightningnetwork/lnd/lnrpc/chainrpc"
	"github.com/lightningnetwork/lnd/lnrpc/invoicesrpc"
	"github.com/lightningnetwork/lnd/lnrpc/routerrpc"
	"github.com/lightningnetwork/lnd/lnwallet"
	"github.com/lightningnetwork/lnd/lnwire"
	"strconv"
	"time"
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

// 开具发票
func InvoiceCreate(amount int64, memo string) (*lnrpc.AddInvoiceResponse, error) {
	lndconf := config.GetConfig().ApiConfig.Lnd
	grpcHost := lndconf.Host + ":" + strconv.Itoa(lndconf.Port)
	tlsCertPath := lndconf.TlsCertPath
	macaroonPath := config.GetConfig().ApiConfig.Lnd.MacaroonPath

	conn, connClose := utils.GetConn(grpcHost, tlsCertPath, macaroonPath)
	defer connClose()

	request := &lnrpc.Invoice{
		Value: amount,
		Memo:  memo,
	}

	client := lnrpc.NewLightningClient(conn)
	response, err := client.AddInvoice(context.Background(), request)
	if err != nil {
		return nil, err
	}
	return response, err
}

// 解析发票
func InvoiceDecode(invoice string) (*lnrpc.PayReq, error) {
	lndconf := config.GetConfig().ApiConfig.Lnd

	grpcHost := lndconf.Host + ":" + strconv.Itoa(lndconf.Port)
	tlsCertPath := lndconf.TlsCertPath
	macaroonPath := lndconf.MacaroonPath

	conn, connClose := utils.GetConn(grpcHost, tlsCertPath, macaroonPath)
	defer connClose()

	request := &lnrpc.PayReqString{
		PayReq: invoice,
	}
	client := lnrpc.NewLightningClient(conn)
	response, err := client.DecodePayReq(context.Background(), request)
	if err != nil {
		return nil, err
	}
	return response, err
}

// 在节点上查询发票
func InvoiceFind(rHash []byte) (*lnrpc.Invoice, error) {
	lndconf := config.GetConfig().ApiConfig.Lnd

	grpcHost := lndconf.Host + ":" + strconv.Itoa(lndconf.Port)
	tlsCertPath := lndconf.TlsCertPath
	macaroonPath := lndconf.MacaroonPath

	conn, connClose := utils.GetConn(grpcHost, tlsCertPath, macaroonPath)
	defer connClose()

	request := &lnrpc.PaymentHash{
		RHash: rHash,
	}
	client := lnrpc.NewLightningClient(conn)
	response, err := client.LookupInvoice(context.Background(), request)
	if err != nil {
		return nil, err
	}
	return response, nil
}

// 支付非0发票
func InvoicePay(invoice string, amt, feeLimit int64) (*lnrpc.Payment, error) {
	lndconf := config.GetConfig().ApiConfig.Lnd
	macaroonFile := config.GetConfig().ApiConfig.Lnd.MacaroonPath
	grpcHost := lndconf.Host + ":" + strconv.Itoa(lndconf.Port)
	tlsCertPath := lndconf.TlsCertPath

	conn, connClose := utils.GetConn(grpcHost, tlsCertPath, macaroonFile)
	defer connClose()

	var paymentTimeout = time.Second * 60
	request := &routerrpc.SendPaymentRequest{
		PaymentRequest:    invoice,
		DestCustomRecords: make(map[uint64][]byte),
		//FeeLimitSat:    feeLimit,
		TimeoutSeconds: int32(paymentTimeout.Seconds()),
		MaxParts:       16,
	}
	if feeLimit > 1 {
		request.FeeLimitSat = feeLimit
	} else {
		amtMsat := lnwire.NewMSatFromSatoshis(btcutil.Amount(amt))
		request.FeeLimitSat = int64(lnwallet.DefaultRoutingFeeLimitForAmount(amtMsat).ToSatoshis())
	}
	client := routerrpc.NewRouterClient(conn)
	stream, err := client.SendPaymentV2(context.Background(), request)
	if err != nil {
		return nil, err
	}
	for {
		payment, err := stream.Recv()
		if err != nil {
			return nil, err
		}
		// Terminate loop if payments state is final.
		if payment.Status != lnrpc.Payment_IN_FLIGHT &&
			payment.Status != lnrpc.Payment_INITIATED {
			return payment, nil
		}
	}
}

func PaymentTrack(paymentHash string) (*lnrpc.Payment, error) {
	lndconf := config.GetConfig().ApiConfig.Lnd

	grpcHost := lndconf.Host + ":" + strconv.Itoa(lndconf.Port)
	tlsCertPath := lndconf.TlsCertPath
	macaroonPath := lndconf.MacaroonPath

	conn, connClose := utils.GetConn(grpcHost, tlsCertPath, macaroonPath)
	defer connClose()
	hash, _ := hex.DecodeString(paymentHash)
	request := &routerrpc.TrackPaymentRequest{
		PaymentHash: hash,
	}
	client := routerrpc.NewRouterClient(conn)
	stream, err := client.TrackPaymentV2(context.Background(), request)
	if err != nil {
		return nil, err
	}
	defer func(stream routerrpc.Router_TrackPaymentV2Client) {
		err := stream.CloseSend()
		if err != nil {

		}
	}(stream)
	for {
		payment, err := stream.Recv()
		if err != nil {
			return nil, err
		}
		if payment != nil {
			if payment.Status == lnrpc.Payment_SUCCEEDED || payment.Status == lnrpc.Payment_FAILED {
				return payment, nil
			}
		}
	}
}

func InvoiceCancel(hash []byte) error {
	lndconf := config.GetConfig().ApiConfig.Lnd

	grpcHost := lndconf.Host + ":" + strconv.Itoa(lndconf.Port)
	tlsCertPath := lndconf.TlsCertPath
	macaroonPath := lndconf.MacaroonPath

	conn, connClose := utils.GetConn(grpcHost, tlsCertPath, macaroonPath)
	defer connClose()
	client := invoicesrpc.NewInvoicesClient(conn)
	request := &invoicesrpc.CancelInvoiceMsg{
		PaymentHash: hash,
	}
	_, err := client.CancelInvoice(context.Background(), request)
	if err != nil {
		return err
	}
	return nil
}

func ListUnspent() (*lnrpc.ListUnspentResponse, error) {
	lndconf := config.GetConfig().ApiConfig.Lnd

	grpcHost := lndconf.Host + ":" + strconv.Itoa(lndconf.Port)
	tlsCertPath := lndconf.TlsCertPath
	macaroonPath := lndconf.MacaroonPath

	conn, connClose := utils.GetConn(grpcHost, tlsCertPath, macaroonPath)
	defer connClose()

	client := lnrpc.NewLightningClient(conn)
	request := &lnrpc.ListUnspentRequest{
		Account: "default",
	}
	response, err := client.ListUnspent(context.Background(), request)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func GetBalance() (*lnrpc.WalletBalanceResponse, error) {
	lndconf := config.GetConfig().ApiConfig.Lnd

	grpcHost := lndconf.Host + ":" + strconv.Itoa(lndconf.Port)
	tlsCertPath := lndconf.TlsCertPath
	macaroonPath := lndconf.MacaroonPath

	conn, connClose := utils.GetConn(grpcHost, tlsCertPath, macaroonPath)
	defer connClose()

	client := lnrpc.NewLightningClient(conn)
	request := &lnrpc.WalletBalanceRequest{}
	response, err := client.WalletBalance(context.Background(), request)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func GetChannelInfo() ([]*lnrpc.Channel, error) {
	lndconf := config.GetConfig().ApiConfig.Lnd

	grpcHost := lndconf.Host + ":" + strconv.Itoa(lndconf.Port)
	tlsCertPath := lndconf.TlsCertPath
	macaroonPath := lndconf.MacaroonPath

	conn, connClose := utils.GetConn(grpcHost, tlsCertPath, macaroonPath)
	defer connClose()

	client := lnrpc.NewLightningClient(conn)
	request := lnrpc.ListChannelsRequest{}
	response, err := client.ListChannels(context.Background(), &request)
	if err != nil {
		return nil, err
	}
	return response.Channels, nil
}
