package custodyBase

type CustodyEvent interface {
	GetBalance() ([]Balance, error)
	ApplyPayReq(PayReqApplyRequest) (PayReqApplyResponse, error)
	SendPayment(PayPacket) error
	GetTransactionHistory()
}
type PayReqApplyRequest interface {
	GetPayReqAmount() int64
}

type PayReqApplyResponse interface {
	GetAmount() int64
	GetPayReq() string
}

type PayPacket interface {
	VerifyPayReq(int64) error
}
