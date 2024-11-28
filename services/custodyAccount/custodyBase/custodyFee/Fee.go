package custodyFee

var (
	ChannelBtcInsideServiceFee = uint64(10)
	ChannelBtcServiceFee       = uint64(100)
	AssetInsideFee             = uint64(10)
	AssetOutsideFee            = uint64(2500)
)

func SetMemoSign() string {
	return "internal transfer"
}
