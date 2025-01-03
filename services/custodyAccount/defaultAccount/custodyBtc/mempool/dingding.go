package mempool

import (
	"fmt"
	"github.com/CatchZeng/dingtalk/pkg/dingtalk"
	"github.com/gookit/color"
	"strings"
	"time"
)

var accessToken = "1999fdf9b8f932ca9295edd44d329f4c9bd3a32b265aa899833bf10411380aa6"
var secret = "SEC83ef77bc6f056ebe0a2c93469dc2d7edade724c1bcef9243cb38da9e8a00a42b"

type Dingding struct {
}
type Balance struct {
	Name  string  `json:"name"`
	Value float64 `json:"value"`
}

// 发送消息
func NewDingding() Dingding {
	return Dingding{}
}

var dingdingMsgOut = `
### 比特币提现通知  
数额:{amount}
时间: {time}
余额信息: 
{balances}`

func (d Dingding) SendBtcPayOutChange(Amount float64, balances []Balance) error {

	client := dingtalk.NewClient(accessToken, secret)

	// 获取设备信息
	var msg string

	msg = strings.Replace(dingdingMsgOut, "{amount}", fmt.Sprintf("%.2f", Amount), 1)
	msg = strings.Replace(msg, "{time}", time.Now().Format("2006-01-02 15:04:05"), 1)

	// 生成余额信息部分
	var balancesInfo string
	for _, balance := range balances {
		balancesInfo += fmt.Sprintf("%s: %.2f\n", balance.Name, balance.Value)
	}
	msg = strings.Replace(msg, "{balances}", balancesInfo, 1)

	color.Infoln("msg:", msg)

	mk := dingtalk.NewMarkdownMessage()
	mk.SetMarkdown("bitlong", msg)

	_, _, err := client.Send(mk)
	if err != nil {
		return err
	}
	return nil
}

var dingdingMsgIn = `
### 比特币充值通知  
数额:{amount}
时间: {time}
余额信息: 
{balances}`

func (d Dingding) ReceiveBtcChannel(Amount float64, balances []Balance) error {
	client := dingtalk.NewClient(accessToken, secret)

	// 获取设备信息
	var msg string

	msg = strings.Replace(dingdingMsgIn, "{amount}", fmt.Sprintf("%.2f", Amount), 1)
	msg = strings.Replace(msg, "{time}", time.Now().Format("2006-01-02 15:04:05"), 1)

	// 生成余额信息部分
	var balancesInfo string
	for _, balance := range balances {
		balancesInfo += fmt.Sprintf("%s: %.2f\n", balance.Name, balance.Value)
	}
	msg = strings.Replace(msg, "{balances}", balancesInfo, 1)

	color.Infoln("msg:", msg)

	mk := dingtalk.NewMarkdownMessage()
	mk.SetMarkdown("bitlong", msg)

	_, _, err := client.Send(mk)
	if err != nil {
		return err
	}
	return nil
}
