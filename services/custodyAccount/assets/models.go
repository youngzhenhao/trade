package assets

import (
	"github.com/lightninglabs/taproot-assets/taprpc"
)

// AssetAddressApplyRequest 资产接收地址申请请求结构体
type AssetAddressApplyRequest struct {
	Amount int64
}

func (req *AssetAddressApplyRequest) GetPayReqAmount() int64 {
	return req.Amount
}

// AssetApplyAddress 资产接收地址申请请求的结构体
type AssetApplyAddress struct {
	Addr   *taprpc.Addr
	Amount int64
}

func (a *AssetApplyAddress) GetAmount() int64 {
	return a.Amount
}
func (a *AssetApplyAddress) GetPayReq() string {
	return a.Addr.Encoded
}
