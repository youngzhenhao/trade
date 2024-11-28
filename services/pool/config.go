package pool

const (
	ZeroValue          = "0"
	TokenSatTag string = "sat"
	// TODO: Set this
	MinLiquidity  uint = 1e2
	MinSatFee     uint = 20
	AssetIdLength      = 64
)

// FeeK is 1000 * Fee Rate
// i.e. Fee Rate = FeeK / 1000

const (
	AddLiquidityFeeK    uint16 = 0
	RemoveLiquidityFeeK uint16 = 3
)

const (
	ProjectPartyFeeK = 3
	LpAwardFeeK      = 3
)

const (
	SwapFeeK uint16 = ProjectPartyFeeK + LpAwardFeeK
)
