package pool

import (
	"errors"
	"gorm.io/gorm"
	"strconv"
	"trade/middleware"
	"trade/utils"
)

type PoolBatchType int64

// TODO
const (
	BatchCreated PoolBatchType = iota
	BatchPending
	BatchCompleted
	BatchFailed = -1
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
	RequestUser      string        `json:"request_user" gorm:"type:varchar(255);index"`
	Username         string        `json:"username" gorm:"type:varchar(255);index"`
	Amount           string        `json:"amount" gorm:"type:varchar(255);index"`
	ResultNewBalance string        `json:"result_new_balance" gorm:"type:varchar(255);index"`
	ResultErr        string        `json:"result_err" gorm:"type:varchar(255);index"`
	ProcessTimes     uint64        `json:"process_times" gorm:"index"`
	State            PoolBatchType `json:"state" gorm:"index"`
}

// Request

type PoolAddLiquidityBatchRequest struct {
	TokenA         string `json:"token_a" gorm:"type:varchar(255);index"`
	TokenB         string `json:"token_b" gorm:"type:varchar(255);index"`
	AmountADesired string `json:"amount_a_desired" gorm:"type:varchar(255);index"`
	AmountBDesired string `json:"amount_b_desired" gorm:"type:varchar(255);index"`
	AmountAMin     string `json:"amount_a_min" gorm:"type:varchar(255);index"`
	AmountBMin     string `json:"amount_b_min" gorm:"type:varchar(255);index"`
	Username       string `json:"username" gorm:"type:varchar(255);index"`
}

type PoolRemoveLiquidityBatchRequest struct {
	TokenA     string `json:"token_a" gorm:"type:varchar(255);index"`
	TokenB     string `json:"token_b" gorm:"type:varchar(255);index"`
	Liquidity  string `json:"liquidity" gorm:"type:varchar(255);index"`
	AmountAMin string `json:"amount_a_min" gorm:"type:varchar(255);index"`
	AmountBMin string `json:"amount_b_min" gorm:"type:varchar(255);index"`
	Username   string `json:"username" gorm:"type:varchar(255);index"`
	FeeK       uint16 `json:"fee_k" gorm:"index"`
}

type PoolSwapExactTokenForTokenNoPathBatchRequest struct {
	TokenIn          string `json:"token_in" gorm:"type:varchar(255);index"`
	TokenOut         string `json:"token_out" gorm:"type:varchar(255);index"`
	AmountIn         string `json:"amount_in" gorm:"type:varchar(255);index"`
	AmountOutMin     string `json:"amount_out_min" gorm:"type:varchar(255);index"`
	Username         string `json:"username" gorm:"type:varchar(255);index"`
	ProjectPartyFeeK uint16 `json:"project_party_fee_k" gorm:"index"`
	LpAwardFeeK      uint16 `json:"lp_award_fee_k" gorm:"index"`
}

type PoolSwapTokenForExactTokenNoPathBatchRequest struct {
	TokenIn          string `json:"token_in" gorm:"type:varchar(255);index"`
	TokenOut         string `json:"token_out" gorm:"type:varchar(255);index"`
	AmountOut        string `json:"amount_out" gorm:"type:varchar(255);index"`
	AmountInMax      string `json:"amount_in_max" gorm:"type:varchar(255);index"`
	Username         string `json:"username" gorm:"type:varchar(255);index"`
	ProjectPartyFeeK uint16 `json:"project_party_fee_k" gorm:"index"`
	LpAwardFeeK      uint16 `json:"lp_award_fee_k" gorm:"index"`
}

type PoolWithdrawAwardBatchRequest struct {
	Username string `json:"username" gorm:"type:varchar(255);index"`
	Amount   string `json:"amount" gorm:"type:varchar(255);index"`
}

// func

func Create(data any) (err error) {
	return middleware.DB.Create(data).Error
}

func ProcessPoolAddLiquidityBatchRequest(request *PoolAddLiquidityBatchRequest, requestUser string) (poolAddLiquidityBatch *PoolAddLiquidityBatch, err error) {
	if request == nil {
		err = errors.New("request is nil")
		return new(PoolAddLiquidityBatch), err
	}
	_, _, err = sortTokens(request.TokenA, request.TokenB)
	if err != nil {
		return new(PoolAddLiquidityBatch), utils.AppendErrorInfo(err, "sortTokens")
	}
	if request.AmountADesired == "" {
		err = errors.New("amount_a_desired is empty")
		return new(PoolAddLiquidityBatch), err
	}
	if request.AmountBDesired == "" {
		err = errors.New("amount_b_desired is empty")
		return new(PoolAddLiquidityBatch), err
	}
	if request.AmountAMin == "" {
		err = errors.New("amount_a_min is empty")
		return new(PoolAddLiquidityBatch), err
	}
	if request.AmountBMin == "" {
		err = errors.New("amount_b_min is empty")
		return new(PoolAddLiquidityBatch), err
	}
	if request.Username == "" {
		err = errors.New("username is empty")
		return new(PoolAddLiquidityBatch), err
	}
	return &PoolAddLiquidityBatch{
		RequestUser:    requestUser,
		TokenA:         request.TokenA,
		TokenB:         request.TokenB,
		AmountADesired: request.AmountADesired,
		AmountBDesired: request.AmountBDesired,
		AmountAMin:     request.AmountAMin,
		AmountBMin:     request.AmountBMin,
		Username:       request.Username,
	}, nil
}

func ProcessPoolRemoveLiquidityBatchRequest(request *PoolRemoveLiquidityBatchRequest, requestUser string) (poolRemoveLiquidityBatch *PoolRemoveLiquidityBatch, err error) {
	if request == nil {
		err = errors.New("request is nil")
		return new(PoolRemoveLiquidityBatch), err
	}
	_, _, err = sortTokens(request.TokenA, request.TokenB)
	if err != nil {
		return new(PoolRemoveLiquidityBatch), utils.AppendErrorInfo(err, "sortTokens")
	}
	if request.Liquidity == "" {
		err = errors.New("liquidity is empty")
		return new(PoolRemoveLiquidityBatch), err
	}
	if request.AmountAMin == "" {
		err = errors.New("amount_a_min is empty")
		return new(PoolRemoveLiquidityBatch), err
	}
	if request.AmountBMin == "" {
		err = errors.New("amount_b_min is empty")
		return new(PoolRemoveLiquidityBatch), err
	}
	if request.Username == "" {
		err = errors.New("username is empty")
		return new(PoolRemoveLiquidityBatch), err
	}
	if request.FeeK != RemoveLiquidityFeeK {
		err = errors.New("invalid fee_k(" + strconv.FormatUint(uint64(request.FeeK), 10) + ")")
		return new(PoolRemoveLiquidityBatch), err
	}
	return &PoolRemoveLiquidityBatch{
		RequestUser: requestUser,
		TokenA:      request.TokenA,
		TokenB:      request.TokenB,
		Liquidity:   request.Liquidity,
		AmountAMin:  request.AmountAMin,
		AmountBMin:  request.AmountBMin,
		Username:    request.Username,
		FeeK:        request.FeeK,
	}, nil
}

func ProcessPoolSwapExactTokenForTokenNoPathBatchRequest(request *PoolSwapExactTokenForTokenNoPathBatchRequest, requestUser string) (poolSwapExactTokenForTokenNoPathBatch *PoolSwapExactTokenForTokenNoPathBatch, err error) {
	if request == nil {
		err = errors.New("request is nil")
		return new(PoolSwapExactTokenForTokenNoPathBatch), err
	}
	_, _, err = sortTokens(request.TokenIn, request.TokenOut)
	if err != nil {
		return new(PoolSwapExactTokenForTokenNoPathBatch), utils.AppendErrorInfo(err, "sortTokens")
	}
	if request.AmountIn == "" {
		err = errors.New("amount_in is empty")
		return new(PoolSwapExactTokenForTokenNoPathBatch), err
	}
	if request.AmountOutMin == "" {
		err = errors.New("amount_out_min is empty")
		return new(PoolSwapExactTokenForTokenNoPathBatch), err
	}
	if request.Username == "" {
		err = errors.New("username is empty")
		return new(PoolSwapExactTokenForTokenNoPathBatch), err
	}
	if request.ProjectPartyFeeK != ProjectPartyFeeK {
		err = errors.New("invalid project_party_fee_k(" + strconv.FormatUint(uint64(request.ProjectPartyFeeK), 10) + ")")
		return new(PoolSwapExactTokenForTokenNoPathBatch), err
	}
	if request.LpAwardFeeK != LpAwardFeeK {
		err = errors.New("invalid lp_award_fee_k(" + strconv.FormatUint(uint64(request.LpAwardFeeK), 10) + ")")
		return new(PoolSwapExactTokenForTokenNoPathBatch), err
	}
	return &PoolSwapExactTokenForTokenNoPathBatch{
		RequestUser:      requestUser,
		TokenIn:          request.TokenIn,
		TokenOut:         request.TokenOut,
		AmountIn:         request.AmountIn,
		AmountOutMin:     request.AmountOutMin,
		Username:         request.Username,
		ProjectPartyFeeK: request.ProjectPartyFeeK,
		LpAwardFeeK:      request.LpAwardFeeK,
	}, nil
}

func ProcessPoolSwapTokenForExactTokenNoPathBatchRequest(request *PoolSwapTokenForExactTokenNoPathBatchRequest, requestUser string) (poolSwapTokenForExactTokenNoPathBatch *PoolSwapTokenForExactTokenNoPathBatch, err error) {
	if request == nil {
		err = errors.New("request is nil")
		return new(PoolSwapTokenForExactTokenNoPathBatch), err
	}
	_, _, err = sortTokens(request.TokenIn, request.TokenOut)
	if err != nil {
		return new(PoolSwapTokenForExactTokenNoPathBatch), utils.AppendErrorInfo(err, "sortTokens")
	}
	if request.AmountOut == "" {
		err = errors.New("amount_out is empty")
		return new(PoolSwapTokenForExactTokenNoPathBatch), err
	}
	if request.AmountInMax == "" {
		err = errors.New("amount_in_max is empty")
		return new(PoolSwapTokenForExactTokenNoPathBatch), err
	}
	if request.Username == "" {
		err = errors.New("username is empty")
		return new(PoolSwapTokenForExactTokenNoPathBatch), err
	}
	if request.ProjectPartyFeeK != ProjectPartyFeeK {
		err = errors.New("invalid project_party_fee_k(" + strconv.FormatUint(uint64(request.ProjectPartyFeeK), 10) + ")")
		return new(PoolSwapTokenForExactTokenNoPathBatch), err
	}
	if request.LpAwardFeeK != LpAwardFeeK {
		err = errors.New("invalid lp_award_fee_k(" + strconv.FormatUint(uint64(request.LpAwardFeeK), 10) + ")")
		return new(PoolSwapTokenForExactTokenNoPathBatch), err
	}
	return &PoolSwapTokenForExactTokenNoPathBatch{
		RequestUser:      requestUser,
		TokenIn:          request.TokenIn,
		TokenOut:         request.TokenOut,
		AmountOut:        request.AmountOut,
		AmountInMax:      request.AmountInMax,
		Username:         request.Username,
		ProjectPartyFeeK: request.ProjectPartyFeeK,
		LpAwardFeeK:      request.LpAwardFeeK,
	}, nil
}

func ProcessPoolWithdrawAwardBatchRequest(request *PoolWithdrawAwardBatchRequest, requestUser string) (poolWithdrawAwardBatch *PoolWithdrawAwardBatch, err error) {
	if request == nil {
		err = errors.New("request is nil")
		return new(PoolWithdrawAwardBatch), err
	}
	if request.Username == "" {
		err = errors.New("username is empty")
		return new(PoolWithdrawAwardBatch), err
	}
	if request.Amount == "" {
		err = errors.New("amount is empty")
		return new(PoolWithdrawAwardBatch), err
	}
	return &PoolWithdrawAwardBatch{
		RequestUser: requestUser,
		Username:    request.Username,
		Amount:      request.Amount,
	}, nil
}
