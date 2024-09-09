package custodyAssets

import (
	"context"
	"github.com/lightninglabs/taproot-assets/taprpc"
	"strconv"
	"trade/btlLog"
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
	btlLog.CUST.Info("AddressServer start")
	for {
		event, err := stream.Recv()
		if err != nil {
			return
		}

		if event != nil {
			btlLog.CUST.Info("%v", event)
			switch event.Status {
			case taprpc.AddrEventStatus_ADDR_EVENT_STATUS_COMPLETED:

			default:

			}
		}
	}
}
