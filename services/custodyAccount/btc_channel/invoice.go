package btc_channel

type SubscribeInvoiceServer struct {
}

func NewSubscribeInvoiceServer() *SubscribeInvoiceServer {
	return &SubscribeInvoiceServer{}
}

func (s *SubscribeInvoiceServer) Start() {
	go s.runServer()
}
func (s *SubscribeInvoiceServer) runServer() {

}
