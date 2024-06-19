package routers

import (
	"github.com/gin-gonic/gin"
	"trade/config"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	if !config.GetLoadConfig().RouterDisable.Login {
		SetupLoginRouter(r)
	}
	if !config.GetLoadConfig().RouterDisable.FairLaunch {
		SetupFairLaunchRouter(r)
	}
	if !config.GetLoadConfig().RouterDisable.Fee {
		SetupFeeRouter(r)
	}
	if !config.GetLoadConfig().RouterDisable.CustodyAccount {
		SetupCustodyAccountRouter(r)
	}
	if !config.GetLoadConfig().RouterDisable.Proof {
		SetupProofRouter(r)
	}
	if !config.GetLoadConfig().RouterDisable.Ido {
		SetupIdoRouter(r)
	}
	if !config.GetLoadConfig().RouterDisable.Snapshot {
		SetupSnapshotRouter(r)
	}
	return r
}
