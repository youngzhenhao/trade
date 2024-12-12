package pool

const (
	ZeroValue          = "0"
	TokenSatTag string = "sat"

	MinLiquidity uint = 1e1

	MinAddLiquiditySat uint = 1e3
	// TODO: ?
	MinRemoveLiquiditySat uint = 1e3

	MinWithdrawAwardSat uint = 1e2

	MinSwapSatFee uint = 20

	AssetIdLength = 64
)

// FeeK is 1000 * Fee Rate
// i.e. Fee Rate = FeeK / 1000

const (
	AddLiquidityFeeK    uint16 = 0
	RemoveLiquidityFeeK uint16 = 3
)

const (
	ProjectPartyFeeK uint16 = 3
	LpAwardFeeK      uint16 = 3
)

const (
	SwapFeeK uint16 = ProjectPartyFeeK + LpAwardFeeK
)

// 3334
var MinSwapSat = func() uint {
	minSwapSat := MinSwapSatFee * 1000 / uint(SwapFeeK)
	if MinSwapSatFee*1000%uint(SwapFeeK) != 0 {
		minSwapSat++
	}
	return minSwapSat
}()
