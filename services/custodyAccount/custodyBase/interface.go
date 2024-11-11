package custodyBase

import (
	caccount "trade/services/custodyAccount/account"
)

type CustodyEvent interface {
	GetBalance() ([]Balance, error)
	ApplyPayReq(PayReqApplyRequest) (PayReqApplyResponse, error)
	SendPayment(PayPacket) error
	GetTransactionHistory(int, int) (*PaymentList, error)
}
type PayReqApplyRequest interface {
	GetPayReqAmount() int64
}

type PayReqApplyResponse interface {
	GetAmount() int64
	GetPayReq() string
}

type PayPacket interface {
	VerifyPayReq(*caccount.UserInfo) error
}
