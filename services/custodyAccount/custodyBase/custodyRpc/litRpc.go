package custodyRpc

import (
	"errors"
	"github.com/lightninglabs/lightning-terminal/litrpc"
	caccount "trade/services/custodyAccount/account"
	"trade/services/servicesrpc"
)

const (
	UpdateBalancePlus  string = "Plus"
	UpdateBalanceMinus string = "Minus"
)

var NotEnoughBalanceError = errors.New("not enough balance")

func GetAccountInfo(usr *caccount.UserInfo) (*litrpc.Account, error) {
	usr.RpcMux.Lock()
	defer usr.RpcMux.Unlock()
	return servicesrpc.AccountInfo(usr.Account.UserAccountCode)
}

func UpdateBalance(usr *caccount.UserInfo, dirt string, amount int64) (*litrpc.Account, error) {
	usr.RpcMux.Lock()
	defer usr.RpcMux.Unlock()

	account, err := servicesrpc.AccountInfo(usr.Account.UserAccountCode)
	if err != nil {
		return nil, err
	}
	if dirt == UpdateBalancePlus {
		account.CurrentBalance += amount
	} else if dirt == UpdateBalanceMinus {
		account.CurrentBalance -= amount
	}

	if account.CurrentBalance < 0 {
		return nil, NotEnoughBalanceError
	}

	result, err := servicesrpc.AccountUpdate(usr.Account.UserAccountCode, account.CurrentBalance, -1)
	if err != nil {
		return nil, err
	}
	return result, nil
}
