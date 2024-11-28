package custodyRpc

import (
	"github.com/lightningnetwork/lnd/lnrpc"
	caccount "trade/services/custodyAccount/account"
	"trade/services/servicesrpc"
)

func PayBtcInvoice(usr *caccount.UserInfo, invoice string, amt, feeLimit int64) (*lnrpc.Payment, error) {
	usr.RpcMux.Lock()
	defer usr.RpcMux.Unlock()
	return servicesrpc.InvoicePay(invoice, amt, feeLimit)
}
