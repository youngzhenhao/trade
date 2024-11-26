package pool

const (
	ZeroValue            = "0"
	TokenSatTag   string = "sat"
	MinLiquidity  uint   = 1e3
	MinSatFee     uint   = 20
	AssetIdLength        = 64
)

// F is 1000 * Fee Rate
// i.e. Fee Rate = F / 1000

const (
	AddLiquidityF    uint16 = 0
	RemoveLiquidityF uint16 = 3
)

const (
	ProjectPartyF = 3
	LpAwardF      = 3
)

const (
	SwapF uint16 = ProjectPartyF + LpAwardF
)
