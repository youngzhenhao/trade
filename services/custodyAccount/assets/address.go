package assets

import (
	"context"
	"github.com/lightninglabs/taproot-assets/taprpc"
	"strconv"
	"trade/config"
	"trade/utils"
)

type SubscribeAddressServer struct {
}

var AddressServer SubscribeAddressServer

func (s *SubscribeAddressServer) Start() {
	go s.runServer()
}
func (s *SubscribeAddressServer) runServer() {
	tapdconf := config.GetConfig().ApiConfig.Tapd
	grpcHost := tapdconf.Host + ":" + strconv.Itoa(tapdconf.Port)
	tlsCertPath := tapdconf.TlsCertPath
	macaroonPath := tapdconf.MacaroonPath
	conn, connClose := utils.GetConn(grpcHost, tlsCertPath, macaroonPath)
	defer connClose()

	client := taprpc.NewTaprootAssetsClient(conn)
	request := &taprpc.SubscribeReceiveEventsRequest{}
	stream, err := client.SubscribeReceiveEvents(context.Background(), request)
	if err != nil {
		return
	}
	for {
		event, err := stream.Recv()
		if err != nil {
			return
		}
		if event.Status == taprpc.AddrEventStatus_ADDR_EVENT_STATUS_COMPLETED {

		}
	}
}
