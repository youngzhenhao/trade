package Interface

type CustodyService interface {
	GetBalance()
	ApplyInvoice()
	PayToOutside()
	PayToInside()
	GetTransactionHistory()
}
