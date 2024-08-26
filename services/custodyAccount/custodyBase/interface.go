package custodyBase

type CustodyEvent interface {
	GetBalance() (*[]Balance, error)
	ApplyInvoice()
	PayToOutside()
	PayToInside()
	GetTransactionHistory()
}
