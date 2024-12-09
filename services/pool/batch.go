package pool

import "gorm.io/gorm"

type PoolBatchType int64

// TODO
const (
	BatchCreated PoolBatchType = iota
	BatchProcessing
	BatchCompleted
)

type PoolAddLiquidityBatch struct {
	gorm.Model
	RequestUser     string        `json:"request_user" gorm:"type:varchar(255);index"`
	TokenA          string        `json:"token_a" gorm:"type:varchar(255);index"`
	TokenB          string        `json:"token_b" gorm:"type:varchar(255);index"`
	AmountADesired  string        `json:"amount_a_desired" gorm:"type:varchar(255);index"`
	AmountBDesired  string        `json:"amount_b_desired" gorm:"type:varchar(255);index"`
	AmountAMin      string        `json:"amount_a_min" gorm:"type:varchar(255);index"`
	AmountBMin      string        `json:"amount_b_min" gorm:"type:varchar(255);index"`
	Username        string        `json:"username" gorm:"type:varchar(255);index"`
	ResultAmountA   string        `json:"result_amount_a" gorm:"type:varchar(255);index"`
	ResultAmountB   string        `json:"result_amount_b" gorm:"type:varchar(255);index"`
	ResultLiquidity string        `json:"result_liquidity" gorm:"type:varchar(255);index"`
	ResultErr       string        `json:"result_err" gorm:"type:varchar(255);index"`
	ProcessTimes    uint64        `json:"process_times" gorm:"index"`
	State           PoolBatchType `json:"state" gorm:"index"`
}

type PoolRemoveLiquidityBatch struct {
	gorm.Model
	RequestUser   string        `json:"request_user" gorm:"type:varchar(255);index"`
	TokenA        string        `json:"token_a" gorm:"type:varchar(255);index"`
	TokenB        string        `json:"token_b" gorm:"type:varchar(255);index"`
	Liquidity     string        `json:"liquidity" gorm:"type:varchar(255);index"`
	AmountAMin    string        `json:"amount_a_min" gorm:"type:varchar(255);index"`
	AmountBMin    string        `json:"amount_b_min" gorm:"type:varchar(255);index"`
	Username      string        `json:"username" gorm:"type:varchar(255);index"`
	FeeK          uint16        `json:"fee_k" gorm:"index"`
	ResultAmountA string        `json:"result_amount_a" gorm:"type:varchar(255);index"`
	ResultAmountB string        `json:"result_amount_b" gorm:"type:varchar(255);index"`
	ResultErr     string        `json:"result_err" gorm:"type:varchar(255);index"`
	ProcessTimes  uint64        `json:"process_times" gorm:"index"`
	State         PoolBatchType `json:"state" gorm:"index"`
}

type PoolSwapExactTokenForTokenNoPathBatch struct {
	gorm.Model
	RequestUser      string        `json:"request_user" gorm:"type:varchar(255);index"`
	TokenIn          string        `json:"token_in" gorm:"type:varchar(255);index"`
	TokenOut         string        `json:"token_out" gorm:"type:varchar(255);index"`
	AmountIn         string        `json:"amount_in" gorm:"type:varchar(255);index"`
	AmountOutMin     string        `json:"amount_out_min" gorm:"type:varchar(255);index"`
	Username         string        `json:"username" gorm:"type:varchar(255);index"`
	ProjectPartyFeeK uint16        `json:"project_party_fee_k" gorm:"index"`
	LpAwardFeeK      uint16        `json:"lp_award_fee_k" gorm:"index"`
	ResultAmountOut  string        `json:"result_amount_out" gorm:"type:varchar(255);index"`
	ResultErr        string        `json:"result_err" gorm:"type:varchar(255);index"`
	ProcessTimes     uint64        `json:"process_times" gorm:"index"`
	State            PoolBatchType `json:"state" gorm:"index"`
}

type PoolSwapTokenForExactTokenNoPathBatch struct {
	gorm.Model
	RequestUser      string        `json:"request_user" gorm:"type:varchar(255);index"`
	TokenIn          string        `json:"token_in" gorm:"type:varchar(255);index"`
	TokenOut         string        `json:"token_out" gorm:"type:varchar(255);index"`
	AmountOut        string        `json:"amount_out" gorm:"type:varchar(255);index"`
	AmountInMax      string        `json:"amount_in_max" gorm:"type:varchar(255);index"`
	Username         string        `json:"username" gorm:"type:varchar(255);index"`
	ProjectPartyFeeK uint16        `json:"project_party_fee_k" gorm:"index"`
	LpAwardFeeK      uint16        `json:"lp_award_fee_k" gorm:"index"`
	ResultAmountIn   string        `json:"result_amount_in" gorm:"type:varchar(255);index"`
	ResultErr        string        `json:"result_err" gorm:"type:varchar(255);index"`
	ProcessTimes     uint64        `json:"process_times" gorm:"index"`
	State            PoolBatchType `json:"state" gorm:"index"`
}

type PoolWithdrawAwardBatch struct {
	gorm.Model
	RequestUser  string        `json:"request_user" gorm:"type:varchar(255);index"`
	Username     string        `json:"username" gorm:"type:varchar(255);index"`
	Amount       string        `json:"amount" gorm:"type:varchar(255);index"`
	AwardBalance string        `json:"award_balance" gorm:"type:varchar(255);index"`
	ResultErr    string        `json:"result_err" gorm:"type:varchar(255);index"`
	ProcessTimes uint64        `json:"process_times" gorm:"index"`
	State        PoolBatchType `json:"state" gorm:"index"`
}
