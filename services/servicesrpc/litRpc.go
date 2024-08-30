package servicesrpc

import (
	"context"
	"github.com/lightninglabs/lightning-terminal/litrpc"
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
