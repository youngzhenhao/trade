package servicesrpc

import (
	"context"
	"encoding/hex"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/lightninglabs/lightning-terminal/litrpc"
	"github.com/lightningnetwork/lnd/lnrpc"
	"github.com/lightningnetwork/lnd/lnrpc/invoicesrpc"
	"github.com/lightningnetwork/lnd/lnrpc/routerrpc"
	"github.com/lightningnetwork/lnd/lnwallet"
	"github.com/lightningnetwork/lnd/lnwire"
	"strconv"
	"trade/config"
	"trade/utils"
)

func AccountCreate(balance uint64, expirationDate int64) (*litrpc.Account, []byte, error) {
	litdconf := config.GetConfig().ApiConfig.Litd

	grpcHost := litdconf.Host + ":" + strconv.Itoa(litdconf.Port)
	tlsCertPath := litdconf.TlsCertPath
	macaroonPath := litdconf.MacaroonPath

	conn, connClose := utils.GetConn(grpcHost, tlsCertPath, macaroonPath)
	defer connClose()

	request := &litrpc.CreateAccountRequest{
		AccountBalance: balance,
		ExpirationDate: expirationDate,
	}
	client := litrpc.NewAccountsClient(conn)
	response, err := client.CreateAccount(context.Background(), request)
	if err != nil {
		return nil, nil, err
	}

	return response.Account, response.Macaroon, nil
}

func AccountInfo(id string) (*litrpc.Account, error) {
	litdconf := config.GetConfig().ApiConfig.Litd

	grpcHost := litdconf.Host + ":" + strconv.Itoa(litdconf.Port)
	tlsCertPath := litdconf.TlsCertPath
	macaroonPath := litdconf.MacaroonPath

	conn, connClose := utils.GetConn(grpcHost, tlsCertPath, macaroonPath)
	defer connClose()

	request := &litrpc.AccountInfoRequest{
		Id: id,
	}
	client := litrpc.NewAccountsClient(conn)
	response, err := client.AccountInfo(context.Background(), request)
	if err != nil {
		return nil, err
	}
	return response, err
}

func acountList() ([]*litrpc.Account, error) {
	litdconf := config.GetConfig().ApiConfig.Litd

	grpcHost := litdconf.Host + ":" + strconv.Itoa(litdconf.Port)
	tlsCertPath := litdconf.TlsCertPath
	macaroonPath := litdconf.MacaroonPath

	conn, connClose := utils.GetConn(grpcHost, tlsCertPath, macaroonPath)
	defer connClose()

	request := &litrpc.ListAccountsRequest{}
	client := litrpc.NewAccountsClient(conn)
	response, err := client.ListAccounts(context.Background(), request)
	if err != nil {
		return nil, err
	}
	return response.Accounts, err
}

func accountQueryId(label string) (string, error) {
	litdconf := config.GetConfig().ApiConfig.Litd

	grpcHost := litdconf.Host + ":" + strconv.Itoa(litdconf.Port)
	tlsCertPath := litdconf.TlsCertPath
	macaroonPath := litdconf.MacaroonPath

	conn, connClose := utils.GetConn(grpcHost, tlsCertPath, macaroonPath)
	defer connClose()

	request := &litrpc.AccountInfoRequest{
		Label: label,
	}
	client := litrpc.NewAccountsClient(conn)
	response, err := client.AccountInfo(context.Background(), request)
	if err != nil {
		return "", err
	}
	return response.Id, err
}

func AccountRemove(id string) error {
	litdconf := config.GetConfig().ApiConfig.Litd

	grpcHost := litdconf.Host + ":" + strconv.Itoa(litdconf.Port)
	tlsCertPath := litdconf.TlsCertPath
	macaroonPath := litdconf.MacaroonPath

	conn, connClose := utils.GetConn(grpcHost, tlsCertPath, macaroonPath)
	defer connClose()

	request := &litrpc.RemoveAccountRequest{
		Id: id,
	}
	client := litrpc.NewAccountsClient(conn)
	_, err := client.RemoveAccount(context.Background(), request)
	return err
}

func AccountUpdate(id string, balance int64, expirationDate int64) (*litrpc.Account, error) {
	litdconf := config.GetConfig().ApiConfig.Litd

	grpcHost := litdconf.Host + ":" + strconv.Itoa(litdconf.Port)
	tlsCertPath := litdconf.TlsCertPath
	macaroonPath := litdconf.MacaroonPath

	conn, connClose := utils.GetConn(grpcHost, tlsCertPath, macaroonPath)
	defer connClose()

	request := &litrpc.UpdateAccountRequest{
		Id:             id,
		AccountBalance: balance,
		ExpirationDate: expirationDate,
	}
	client := litrpc.NewAccountsClient(conn)
	response, err := client.UpdateAccount(context.Background(), request)
	if err != nil {
		return nil, err
	}
	return response, err
}

// TODO: 开通通道
func channelOpen() {}

// TODO: 关闭通道
func channelClose() {}

func LitdInfo() string {
	litdconf := config.GetConfig().ApiConfig.Litd

	grpcHost := litdconf.Host + ":" + strconv.Itoa(litdconf.Port)
	tlsCertPath := litdconf.TlsCertPath
	macaroonPath := litdconf.MacaroonPath

	conn, connClose := utils.GetConn(grpcHost, tlsCertPath, macaroonPath)
	defer connClose()

	request := &litrpc.GetInfoRequest{}

	client := litrpc.NewProxyClient(conn)
	response, err := client.GetInfo(context.Background(), request)
	if err != nil {
		return "Error: " + err.Error()
	}
	return response.String()
}

func LitdStatus() string {
	litdconf := config.GetConfig().ApiConfig.Litd

	grpcHost := litdconf.Host + ":" + strconv.Itoa(litdconf.Port)
	tlsCertPath := litdconf.TlsCertPath
	macaroonPath := litdconf.MacaroonPath

	conn, connClose := utils.GetConn(grpcHost, tlsCertPath, macaroonPath)
	defer connClose()

	request := &litrpc.SubServerStatusReq{}
	client := litrpc.NewStatusClient(conn)
	response, err := client.SubServerStatus(context.Background(), request)
	if err != nil {
		return "Error: " + err.Error()
	}
	return response.String()
}

// 开具发票
func InvoiceCreate(amount int64, memo string, macaroonPath string) (*lnrpc.AddInvoiceResponse, error) {
	lndconf := config.GetConfig().ApiConfig.Lnd

	grpcHost := lndconf.Host + ":" + strconv.Itoa(lndconf.Port)
	tlsCertPath := lndconf.TlsCertPath

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
func InvoicePay(macaroonPath string, invoice string, amt, feeLimit int64) (*lnrpc.Payment, error) {
	lndconf := config.GetConfig().ApiConfig.Lnd

	grpcHost := lndconf.Host + ":" + strconv.Itoa(lndconf.Port)
	tlsCertPath := lndconf.TlsCertPath

	conn, connClose := utils.GetConn(grpcHost, tlsCertPath, macaroonPath)
	defer connClose()

	request := &routerrpc.SendPaymentRequest{
		PaymentRequest: invoice,
		//FeeLimitSat:    feeLimit,
		TimeoutSeconds: 10,
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
		if payment != nil {
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
			return payment, nil
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
