package handlers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"trade/models"
	"trade/services/pool"
)

// query

func QueryPoolInfo(c *gin.Context) {
	_ = c.MustGet("username").(string)
	tokenA := c.Query("token_a")
	tokenB := c.Query("token_b")
	poolInfo, err := pool.QueryPoolInfo(tokenA, tokenB)
	if err != nil {

		if errors.Is(err, pool.PoolDoesNotExistErr) {
			c.JSON(http.StatusOK, Result2{
				Errno:  models.PoolDoesNotExistErr.Code(),
				ErrMsg: err.Error(),
				Data:   poolInfo,
			})
			return
		}

		c.JSON(http.StatusOK, Result2{
			Errno:  models.QueryPoolInfoErr.Code(),
			ErrMsg: err.Error(),
			Data:   poolInfo,
		})
		return
	}
	c.JSON(http.StatusOK, Result2{
		Errno:  0,
		ErrMsg: models.SUCCESS.Error(),
		Data:   poolInfo,
	})
}

func QueryShareRecordsCount(c *gin.Context) {
	tokenA := c.Query("token_a")
	tokenB := c.Query("token_b")
	username := c.Query("username")
	var count int64
	var err error
	if username == "" {
		count, err = pool.QueryShareRecordsCount(tokenA, tokenB)
		if err != nil {
			c.JSON(http.StatusOK, Result2{
				Errno:  models.QueryShareRecordsCountErr.Code(),
				ErrMsg: err.Error(),
				Data:   0,
			})
			return
		}
	} else {
		count, err = pool.QueryUserShareRecordsCount(tokenA, tokenB, username)
		if err != nil {
			c.JSON(http.StatusOK, Result2{
				Errno:  models.QueryUserShareRecordsCountErr.Code(),
				ErrMsg: err.Error(),
				Data:   0,
			})
			return
		}
	}
	c.JSON(http.StatusOK, Result2{
		Errno:  0,
		ErrMsg: models.SUCCESS.Error(),
		Data:   count,
	})
}

func QueryShareRecords(c *gin.Context) {
	tokenA := c.Query("token_a")
	tokenB := c.Query("token_b")
	limit := c.Query("limit")
	offset := c.Query("offset")
	username := c.Query("username")

	if limit == "" {
		err := errors.New("limit is empty")
		c.JSON(http.StatusOK, Result2{
			Errno:  models.LimitEmptyErr.Code(),
			ErrMsg: err.Error(),
			Data:   &[]pool.ShareRecordInfo{},
		})
		return
	}
	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		c.JSON(http.StatusOK, Result2{
			Errno:  models.AtoiErr.Code(),
			ErrMsg: err.Error(),
			Data:   &[]pool.ShareRecordInfo{},
		})
		return
	}
	if limitInt < 0 {
		err := errors.New("limit is less than 0")
		c.JSON(http.StatusOK, Result2{
			Errno:  models.LimitLessThanZeroErr.Code(),
			ErrMsg: err.Error(),
			Data:   &[]pool.ShareRecordInfo{},
		})
		return
	}
	if offset == "" {
		err := errors.New("offset is empty")
		c.JSON(http.StatusOK, Result2{
			Errno:  models.OffsetEmptyErr.Code(),
			ErrMsg: err.Error(),
			Data:   &[]pool.ShareRecordInfo{},
		})
		return
	}
	offsetInt, err := strconv.Atoi(offset)
	if err != nil {
		c.JSON(http.StatusOK, Result2{
			Errno:  models.AtoiErr.Code(),
			ErrMsg: err.Error(),
			Data:   &[]pool.ShareRecordInfo{},
		})
		return
	}
	if offsetInt < 0 {
		err := errors.New("offset is less than 0")
		c.JSON(http.StatusOK, Result2{
			Errno:  models.OffsetLessThanZeroErr.Code(),
			ErrMsg: err.Error(),
			Data:   &[]pool.ShareRecordInfo{},
		})
		return
	}
	var shareRecords *[]pool.ShareRecordInfo

	if username == "" {
		shareRecords, err = pool.QueryShareRecords(tokenA, tokenB, limitInt, offsetInt)
		if err != nil {
			c.JSON(http.StatusOK, Result2{
				Errno:  models.QueryShareRecordsErr.Code(),
				ErrMsg: err.Error(),
				Data:   &[]pool.ShareRecordInfo{},
			})
			return
		}
	} else {
		shareRecords, err = pool.QueryUserShareRecords(tokenA, tokenB, username, limitInt, offsetInt)
		if err != nil {
			c.JSON(http.StatusOK, Result2{
				Errno:  models.QueryUserShareRecordsErr.Code(),
				ErrMsg: err.Error(),
				Data:   &[]pool.ShareRecordInfo{},
			})
			return
		}
	}
	c.JSON(http.StatusOK, Result2{
		Errno:  0,
		ErrMsg: models.SUCCESS.Error(),
		Data:   shareRecords,
	})
}

func QueryUserShareBalance(c *gin.Context) {
	username := c.MustGet("username").(string)
	tokenA := c.Query("token_a")
	tokenB := c.Query("token_b")
	if username == "" {
		err := errors.New("username is empty")
		c.JSON(http.StatusOK, Result2{
			Errno:  models.UsernameEmptyErr.Code(),
			ErrMsg: err.Error(),
			Data:   &pool.PoolShareBalanceInfo{},
		})
		return
	}
	shareBalanceInfo, err := pool.QueryUserShareBalance(tokenA, tokenB, username)
	if err != nil {
		c.JSON(http.StatusOK, Result2{
			Errno:  models.QueryUserShareBalanceErr.Code(),
			ErrMsg: err.Error(),
			Data:   &pool.PoolShareBalanceInfo{},
		})
		return
	}
	c.JSON(http.StatusOK, Result2{
		Errno:  0,
		ErrMsg: models.SUCCESS.Error(),
		Data:   shareBalanceInfo,
	})
}

func QuerySwapRecordsCount(c *gin.Context) {
	tokenA := c.Query("token_a")
	tokenB := c.Query("token_b")
	username := c.Query("username")
	var count int64
	var err error
	if username == "" {
		count, err = pool.QuerySwapRecordsCount(tokenA, tokenB)
		if err != nil {
			c.JSON(http.StatusOK, Result2{
				Errno:  models.QuerySwapRecordsCountErr.Code(),
				ErrMsg: err.Error(),
				Data:   0,
			})
			return
		}
	} else {
		count, err = pool.QueryUserSwapRecordsCount(tokenA, tokenB, username)
		if err != nil {
			c.JSON(http.StatusOK, Result2{
				Errno:  models.QueryUserSwapRecordsCountErr.Code(),
				ErrMsg: err.Error(),
				Data:   0,
			})
			return
		}
	}
	c.JSON(http.StatusOK, Result2{
		Errno:  0,
		ErrMsg: models.SUCCESS.Error(),
		Data:   count,
	})

}

func QuerySwapRecords(c *gin.Context) {
	tokenA := c.Query("token_a")
	tokenB := c.Query("token_b")
	limit := c.Query("limit")
	offset := c.Query("offset")
	username := c.Query("username")

	if limit == "" {
		err := errors.New("limit is empty")
		c.JSON(http.StatusOK, Result2{
			Errno:  models.LimitEmptyErr.Code(),
			ErrMsg: err.Error(),
			Data:   &[]pool.SwapRecordInfo{},
		})
		return
	}
	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		c.JSON(http.StatusOK, Result2{
			Errno:  models.AtoiErr.Code(),
			ErrMsg: err.Error(),
			Data:   &[]pool.SwapRecordInfo{},
		})
		return
	}
	if limitInt < 0 {
		err := errors.New("limit is less than 0")
		c.JSON(http.StatusOK, Result2{
			Errno:  models.LimitLessThanZeroErr.Code(),
			ErrMsg: err.Error(),
			Data:   &[]pool.SwapRecordInfo{},
		})
		return
	}

	if offset == "" {
		err := errors.New("offset is empty")
		c.JSON(http.StatusOK, Result2{
			Errno:  models.OffsetEmptyErr.Code(),
			ErrMsg: err.Error(),
			Data:   &[]pool.SwapRecordInfo{},
		})
		return
	}
	offsetInt, err := strconv.Atoi(offset)
	if err != nil {
		c.JSON(http.StatusOK, Result2{
			Errno:  models.AtoiErr.Code(),
			ErrMsg: err.Error(),
			Data:   &[]pool.SwapRecordInfo{},
		})
		return
	}
	if offsetInt < 0 {
		err := errors.New("offset is less than 0")
		c.JSON(http.StatusOK, Result2{
			Errno:  models.OffsetLessThanZeroErr.Code(),
			ErrMsg: err.Error(),
			Data:   &[]pool.SwapRecordInfo{},
		})
		return
	}

	var swapRecords *[]pool.SwapRecordInfo

	if username == "" {
		swapRecords, err = pool.QuerySwapRecords(tokenA, tokenB, limitInt, offsetInt)
		if err != nil {
			c.JSON(http.StatusOK, Result2{
				Errno:  models.QuerySwapRecordsErr.Code(),
				ErrMsg: err.Error(),
				Data:   &[]pool.SwapRecordInfo{},
			})
			return
		}
	} else {
		swapRecords, err = pool.QueryUserSwapRecords(tokenA, tokenB, username, limitInt, offsetInt)
		if err != nil {
			c.JSON(http.StatusOK, Result2{
				Errno:  models.QueryUserSwapRecordsErr.Code(),
				ErrMsg: err.Error(),
				Data:   &[]pool.SwapRecordInfo{},
			})
			return
		}
	}
	c.JSON(http.StatusOK, Result2{
		Errno:  0,
		ErrMsg: models.SUCCESS.Error(),
		Data:   swapRecords,
	})
}

func QueryUserLpAwardBalance(c *gin.Context) {
	username := c.MustGet("username").(string)
	if username == "" {
		err := errors.New("username is empty")
		c.JSON(http.StatusOK, Result2{
			Errno:  models.UsernameEmptyErr.Code(),
			ErrMsg: err.Error(),
			Data:   &pool.LpAwardBalanceInfo{},
		})
		return
	}
	lpAwardBalance, err := pool.QueryUserLpAwardBalance(username)
	if err != nil {
		c.JSON(http.StatusOK, Result2{
			Errno:  models.QueryUserLpAwardBalanceErr.Code(),
			ErrMsg: err.Error(),
			Data:   &pool.LpAwardBalanceInfo{},
		})
		return
	}
	c.JSON(http.StatusOK, Result2{
		Errno:  0,
		ErrMsg: models.SUCCESS.Error(),
		Data:   lpAwardBalance,
	})
}

func QueryWithdrawAwardRecordsCount(c *gin.Context) {
	username := c.MustGet("username").(string)
	if username == "" {
		err := errors.New("username is empty")
		c.JSON(http.StatusOK, Result2{
			Errno:  models.UsernameEmptyErr.Code(),
			ErrMsg: err.Error(),
			Data:   0,
		})
		return
	}
	count, err := pool.QueryUserWithdrawAwardRecordsCount(username)
	if err != nil {
		c.JSON(http.StatusOK, Result2{
			Errno:  models.QueryUserWithdrawAwardRecordsCountErr.Code(),
			ErrMsg: err.Error(),
			Data:   0,
		})
		return
	}
	c.JSON(http.StatusOK, Result2{
		Errno:  0,
		ErrMsg: models.SUCCESS.Error(),
		Data:   count,
	})
}

func QueryWithdrawAwardRecords(c *gin.Context) {
	username := c.MustGet("username").(string)
	if username == "" {
		err := errors.New("username is empty")
		c.JSON(http.StatusOK, Result2{
			Errno:  models.UsernameEmptyErr.Code(),
			ErrMsg: err.Error(),
			Data:   &[]pool.PoolWithdrawAwardRecord{},
		})
		return
	}
	limit := c.Query("limit")
	offset := c.Query("offset")

	if limit == "" {
		err := errors.New("limit is empty")
		c.JSON(http.StatusOK, Result2{
			Errno:  models.LimitEmptyErr.Code(),
			ErrMsg: err.Error(),
			Data:   &[]pool.WithdrawAwardRecordInfo{},
		})
		return
	}
	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		c.JSON(http.StatusOK, Result2{
			Errno:  models.AtoiErr.Code(),
			ErrMsg: err.Error(),
			Data:   &[]pool.WithdrawAwardRecordInfo{},
		})
		return
	}
	if limitInt < 0 {
		err := errors.New("limit is less than 0")
		c.JSON(http.StatusOK, Result2{
			Errno:  models.LimitLessThanZeroErr.Code(),
			ErrMsg: err.Error(),
			Data:   &[]pool.WithdrawAwardRecordInfo{},
		})
		return
	}

	if offset == "" {
		err := errors.New("offset is empty")
		c.JSON(http.StatusOK, Result2{
			Errno:  models.OffsetEmptyErr.Code(),
			ErrMsg: err.Error(),
			Data:   &[]pool.WithdrawAwardRecordInfo{},
		})
		return
	}
	offsetInt, err := strconv.Atoi(offset)
	if err != nil {
		c.JSON(http.StatusOK, Result2{
			Errno:  models.AtoiErr.Code(),
			ErrMsg: err.Error(),
			Data:   &[]pool.WithdrawAwardRecordInfo{},
		})
		return
	}
	if offsetInt < 0 {
		err := errors.New("offset is less than 0")
		c.JSON(http.StatusOK, Result2{
			Errno:  models.OffsetLessThanZeroErr.Code(),
			ErrMsg: err.Error(),
			Data:   &[]pool.WithdrawAwardRecordInfo{},
		})
		return
	}

	var poolWithdrawAwardRecords *[]pool.WithdrawAwardRecordInfo

	poolWithdrawAwardRecords, err = pool.QueryUserWithdrawAwardRecords(username, limitInt, offsetInt)
	if err != nil {
		c.JSON(http.StatusOK, Result2{
			Errno:  models.QueryUserWithdrawAwardRecordsErr.Code(),
			ErrMsg: err.Error(),
			Data:   &[]pool.WithdrawAwardRecordInfo{},
		})
		return
	}
	c.JSON(http.StatusOK, Result2{
		Errno:  0,
		ErrMsg: models.SUCCESS.Error(),
		Data:   poolWithdrawAwardRecords,
	})
}

// calc

func CalcAddLiquidity(c *gin.Context) {
	requestUser := c.MustGet("username").(string)
	var poolAddLiquidityBatchRequest pool.PoolAddLiquidityRequest
	err := c.ShouldBindJSON(&poolAddLiquidityBatchRequest)
	if err != nil {
		c.JSON(http.StatusOK, Result2{
			Errno:  models.ShouldBindJsonErr.Code(),
			ErrMsg: err.Error(),
			Data:   new(pool.CalcAddLiquidityResponse),
		})
		return
	}

	_, err = pool.ProcessPoolAddLiquidityBatchRequest(&poolAddLiquidityBatchRequest, requestUser)
	if err != nil {
		c.JSON(http.StatusOK, Result2{
			Errno:  models.ProcessPoolAddLiquidityBatchRequestErr.Code(),
			ErrMsg: err.Error(),
			Data:   new(pool.CalcAddLiquidityResponse),
		})
		return
	}

	if requestUser != poolAddLiquidityBatchRequest.Username {
		c.JSON(http.StatusOK, Result2{
			Errno:  models.UsernameNotMatchErr.Code(),
			ErrMsg: "username not match",
			Data:   new(pool.CalcAddLiquidityResponse),
		})
		return
	}

	var calcAddLiquidityResponse pool.CalcAddLiquidityResponse

	var (
		amountA     string
		amountB     string
		liquidity   string
		shareRecord *pool.PoolShareRecord
	)
	tokenA := poolAddLiquidityBatchRequest.TokenA
	tokenB := poolAddLiquidityBatchRequest.TokenB
	amountADesired := poolAddLiquidityBatchRequest.AmountADesired
	amountBDesired := poolAddLiquidityBatchRequest.AmountBDesired
	amountAMin := poolAddLiquidityBatchRequest.AmountAMin
	amountBMin := poolAddLiquidityBatchRequest.AmountBMin
	username := poolAddLiquidityBatchRequest.Username

	amountA, amountB, liquidity, shareRecord, err = pool.CalcAddLiquidity(tokenA, tokenB, amountADesired, amountBDesired, amountAMin, amountBMin, username)

	if err != nil {
		c.JSON(http.StatusOK, Result2{
			Errno:  models.CalcAddLiquidityErr.Code(),
			ErrMsg: err.Error(),
			Data:   new(pool.CalcAddLiquidityResponse),
		})
		return
	}
	calcAddLiquidityResponse = pool.CalcAddLiquidityResponse{
		AmountA:     amountA,
		AmountB:     amountB,
		Liquidity:   liquidity,
		ShareRecord: shareRecord.ToShareRecordInfo(),
	}

	c.JSON(http.StatusOK, Result2{
		Errno:  0,
		ErrMsg: models.SUCCESS.Error(),
		Data:   &calcAddLiquidityResponse,
	})
}

func CalcRemoveLiquidity(c *gin.Context) {
	requestUser := c.MustGet("username").(string)
	var poolRemoveLiquidityBatchRequest pool.PoolRemoveLiquidityRequest
	err := c.ShouldBindJSON(&poolRemoveLiquidityBatchRequest)
	if err != nil {
		c.JSON(http.StatusOK, Result2{
			Errno:  models.ShouldBindJsonErr.Code(),
			ErrMsg: err.Error(),
			Data:   new(pool.CalcRemoveLiquidityResponse),
		})
		return
	}

	_, err = pool.ProcessPoolRemoveLiquidityBatchRequest(&poolRemoveLiquidityBatchRequest, requestUser)
	if err != nil {
		c.JSON(http.StatusOK, Result2{
			Errno:  models.ProcessPoolRemoveLiquidityBatchRequestErr.Code(),
			ErrMsg: err.Error(),
			Data:   new(pool.CalcRemoveLiquidityResponse),
		})
		return
	}

	if requestUser != poolRemoveLiquidityBatchRequest.Username {
		c.JSON(http.StatusOK, Result2{
			Errno:  models.UsernameNotMatchErr.Code(),
			ErrMsg: "username not match",
			Data:   new(pool.CalcRemoveLiquidityResponse),
		})
		return
	}

	var calcRemoveLiquidityResponse pool.CalcRemoveLiquidityResponse

	var (
		amountA     string
		amountB     string
		shareRecord *pool.PoolShareRecord
	)

	tokenA := poolRemoveLiquidityBatchRequest.TokenA
	tokenB := poolRemoveLiquidityBatchRequest.TokenB
	liquidity := poolRemoveLiquidityBatchRequest.Liquidity
	amountAMin := poolRemoveLiquidityBatchRequest.AmountAMin
	amountBMin := poolRemoveLiquidityBatchRequest.AmountBMin
	username := poolRemoveLiquidityBatchRequest.Username
	feeK := poolRemoveLiquidityBatchRequest.FeeK

	amountA, amountB, shareRecord, err = pool.CalcRemoveLiquidity(tokenA, tokenB, liquidity, amountAMin, amountBMin, username, feeK)
	if err != nil {
		c.JSON(http.StatusOK, Result2{
			Errno:  models.CalcRemoveLiquidityErr.Code(),
			ErrMsg: err.Error(),
			Data:   new(pool.CalcRemoveLiquidityResponse),
		})
		return
	}

	calcRemoveLiquidityResponse = pool.CalcRemoveLiquidityResponse{
		AmountA:     amountA,
		AmountB:     amountB,
		ShareRecord: shareRecord.ToShareRecordInfo(),
	}

	c.JSON(http.StatusOK, Result2{
		Errno:  0,
		ErrMsg: models.SUCCESS.Error(),
		Data:   &calcRemoveLiquidityResponse,
	})
}

func CalcSwapExactTokenForTokenNoPath(c *gin.Context) {
	requestUser := c.MustGet("username").(string)
	var poolSwapExactTokenForTokenNoPathBatchRequest pool.PoolSwapExactTokenForTokenNoPathRequest
	err := c.ShouldBindJSON(&poolSwapExactTokenForTokenNoPathBatchRequest)
	if err != nil {
		c.JSON(http.StatusOK, Result2{
			Errno:  models.ShouldBindJsonErr.Code(),
			ErrMsg: err.Error(),
			Data:   new(pool.CalcSwapExactTokenForTokenNoPathResponse),
		})
		return
	}

	_, err = pool.ProcessPoolSwapExactTokenForTokenNoPathBatchRequest(&poolSwapExactTokenForTokenNoPathBatchRequest, requestUser)
	if err != nil {
		c.JSON(http.StatusOK, Result2{
			Errno:  models.ProcessPoolSwapExactTokenForTokenNoPathBatchRequestErr.Code(),
			ErrMsg: err.Error(),
			Data:   new(pool.CalcSwapExactTokenForTokenNoPathResponse),
		})
		return
	}

	if requestUser != poolSwapExactTokenForTokenNoPathBatchRequest.Username {
		c.JSON(http.StatusOK, Result2{
			Errno:  models.UsernameNotMatchErr.Code(),
			ErrMsg: "username not match",
			Data:   new(pool.CalcSwapExactTokenForTokenNoPathResponse),
		})
		return
	}

	var calcSwapExactTokenForTokenNoPathResponse pool.CalcSwapExactTokenForTokenNoPathResponse

	var (
		amountOut  string
		swapRecord *pool.PoolSwapRecord
	)

	tokenIn := poolSwapExactTokenForTokenNoPathBatchRequest.TokenIn
	tokenOut := poolSwapExactTokenForTokenNoPathBatchRequest.TokenOut
	amountIn := poolSwapExactTokenForTokenNoPathBatchRequest.AmountIn
	amountOutMin := poolSwapExactTokenForTokenNoPathBatchRequest.AmountOutMin
	username := poolSwapExactTokenForTokenNoPathBatchRequest.Username
	projectPartyFeeK := poolSwapExactTokenForTokenNoPathBatchRequest.ProjectPartyFeeK
	lpAwardFeeK := poolSwapExactTokenForTokenNoPathBatchRequest.LpAwardFeeK

	amountOut, swapRecord, err = pool.CalcSwapExactTokenForTokenNoPath(tokenIn, tokenOut, amountIn, amountOutMin, username, projectPartyFeeK, lpAwardFeeK)
	if err != nil {
		c.JSON(http.StatusOK, Result2{
			Errno:  models.CalcSwapExactTokenForTokenNoPathErr.Code(),
			ErrMsg: err.Error(),
			Data:   new(pool.CalcSwapExactTokenForTokenNoPathResponse),
		})
		return
	}

	calcSwapExactTokenForTokenNoPathResponse = pool.CalcSwapExactTokenForTokenNoPathResponse{
		AmountOut:  amountOut,
		SwapRecord: swapRecord.ToSwapRecordInfo(),
	}

	c.JSON(http.StatusOK, Result2{
		Errno:  0,
		ErrMsg: models.SUCCESS.Error(),
		Data:   &calcSwapExactTokenForTokenNoPathResponse,
	})
}

func CalcSwapTokenForExactTokenNoPath(c *gin.Context) {
	requestUser := c.MustGet("username").(string)
	var poolSwapTokenForExactTokenNoPathBatchRequest pool.PoolSwapTokenForExactTokenNoPathRequest
	err := c.ShouldBindJSON(&poolSwapTokenForExactTokenNoPathBatchRequest)
	if err != nil {
		c.JSON(http.StatusOK, Result2{
			Errno:  models.ShouldBindJsonErr.Code(),
			ErrMsg: err.Error(),
			Data:   new(pool.CalcSwapTokenForExactTokenNoPathResponse),
		})
		return
	}

	_, err = pool.ProcessPoolSwapTokenForExactTokenNoPathBatchRequest(&poolSwapTokenForExactTokenNoPathBatchRequest, requestUser)
	if err != nil {
		c.JSON(http.StatusOK, Result2{
			Errno:  models.ProcessPoolSwapTokenForExactTokenNoPathBatchRequestErr.Code(),
			ErrMsg: err.Error(),
			Data:   new(pool.CalcSwapTokenForExactTokenNoPathResponse),
		})
		return
	}

	if requestUser != poolSwapTokenForExactTokenNoPathBatchRequest.Username {
		c.JSON(http.StatusOK, Result2{
			Errno:  models.UsernameNotMatchErr.Code(),
			ErrMsg: "username not match",
			Data:   new(pool.CalcSwapTokenForExactTokenNoPathResponse),
		})
		return
	}

	var calcSwapTokenForExactTokenNoPathResponse pool.CalcSwapTokenForExactTokenNoPathResponse

	var (
		amountIn   string
		swapRecord *pool.PoolSwapRecord
	)

	tokenIn := poolSwapTokenForExactTokenNoPathBatchRequest.TokenIn
	tokenOut := poolSwapTokenForExactTokenNoPathBatchRequest.TokenOut
	amountOut := poolSwapTokenForExactTokenNoPathBatchRequest.AmountOut
	amountInMax := poolSwapTokenForExactTokenNoPathBatchRequest.AmountInMax
	username := poolSwapTokenForExactTokenNoPathBatchRequest.Username
	projectPartyFeeK := poolSwapTokenForExactTokenNoPathBatchRequest.ProjectPartyFeeK
	lpAwardFeeK := poolSwapTokenForExactTokenNoPathBatchRequest.LpAwardFeeK

	amountIn, swapRecord, err = pool.CalcSwapTokenForExactTokenNoPath(tokenIn, tokenOut, amountOut, amountInMax, username, projectPartyFeeK, lpAwardFeeK)
	if err != nil {
		c.JSON(http.StatusOK, Result2{
			Errno:  models.CalcSwapTokenForExactTokenNoPathErr.Code(),
			ErrMsg: err.Error(),
			Data:   new(pool.CalcSwapTokenForExactTokenNoPathResponse),
		})
		return
	}

	calcSwapTokenForExactTokenNoPathResponse = pool.CalcSwapTokenForExactTokenNoPathResponse{
		AmountIn:   amountIn,
		SwapRecord: swapRecord.ToSwapRecordInfo(),
	}

	c.JSON(http.StatusOK, Result2{
		Errno:  0,
		ErrMsg: models.SUCCESS.Error(),
		Data:   &calcSwapTokenForExactTokenNoPathResponse,
	})
}

func CalcQuote(c *gin.Context) {
	_ = c.MustGet("username").(string)
	tokenA := c.Query("token_a")
	tokenB := c.Query("token_b")
	amountA := c.Query("amount_a")
	if amountA == "" {
		err := errors.New("amount_a is empty")
		c.JSON(http.StatusOK, Result2{
			Errno:  models.QueryParamEmptyErr.Code(),
			ErrMsg: err.Error(),
			Data:   pool.ZeroValue,
		})
		return
	}

	amountB, err := pool.CalcQuote(tokenA, tokenB, amountA)
	if err != nil {
		c.JSON(http.StatusOK, Result2{
			Errno:  models.CalcQuoteErr.Code(),
			ErrMsg: err.Error(),
			Data:   pool.ZeroValue,
		})
		return
	}
	c.JSON(http.StatusOK, Result2{
		Errno:  0,
		ErrMsg: models.SUCCESS.Error(),
		Data:   amountB,
	})
}

func CalcBurnLiquidity(c *gin.Context) {
	username := c.MustGet("username").(string)
	tokenA := c.Query("token_a")
	tokenB := c.Query("token_b")
	liquidity := c.Query("liquidity")
	if liquidity == "" {
		err := errors.New("liquidity is empty")
		c.JSON(http.StatusOK, Result2{
			Errno:  models.QueryParamEmptyErr.Code(),
			ErrMsg: err.Error(),
			Data:   pool.CalcBurnLiquidityResponse{},
		})
		return
	}

	amountA, amountB, err := pool.CalcBurnLiquidity(tokenA, tokenB, liquidity, username, pool.RemoveLiquidityFeeK)
	if err != nil {
		c.JSON(http.StatusOK, Result2{
			Errno:  models.CalcBurnLiquidityErr.Code(),
			ErrMsg: err.Error(),
			Data:   pool.CalcBurnLiquidityResponse{},
		})
		return
	}
	c.JSON(http.StatusOK, Result2{
		Errno:  0,
		ErrMsg: models.SUCCESS.Error(),
		Data: pool.CalcBurnLiquidityResponse{
			AmountA: amountA,
			AmountB: amountB,
		},
	})
}

func CalcAmountOut(c *gin.Context) {
	_ = c.MustGet("username").(string)
	tokenIn := c.Query("token_in")
	tokenOut := c.Query("token_out")
	amountIn := c.Query("amount_in")
	if amountIn == "" {
		err := errors.New("amount_in is empty")
		c.JSON(http.StatusOK, Result2{
			Errno:  models.QueryParamEmptyErr.Code(),
			ErrMsg: err.Error(),
			Data:   pool.ZeroValue,
		})
		return
	}

	amountOut, err := pool.CalcAmountOut(tokenIn, tokenOut, amountIn)
	if err != nil {
		c.JSON(http.StatusOK, Result2{
			Errno:  models.CalcAmountOutErr.Code(),
			ErrMsg: err.Error(),
			Data:   pool.ZeroValue,
		})
		return
	}
	c.JSON(http.StatusOK, Result2{
		Errno:  0,
		ErrMsg: models.SUCCESS.Error(),
		Data:   amountOut,
	})
}

func CalcAmountIn(c *gin.Context) {
	_ = c.MustGet("username").(string)
	tokenIn := c.Query("token_in")
	tokenOut := c.Query("token_out")
	amountOut := c.Query("amount_out")
	if amountOut == "" {
		err := errors.New("amount_out is empty")
		c.JSON(http.StatusOK, Result2{
			Errno:  models.QueryParamEmptyErr.Code(),
			ErrMsg: err.Error(),
			Data:   pool.ZeroValue,
		})
		return
	}

	amountIn, err := pool.CalcAmountIn(tokenIn, tokenOut, amountOut)
	if err != nil {
		c.JSON(http.StatusOK, Result2{
			Errno:  models.CalcAmountInErr.Code(),
			ErrMsg: err.Error(),
			Data:   pool.ZeroValue,
		})
		return
	}
	c.JSON(http.StatusOK, Result2{
		Errno:  0,
		ErrMsg: models.SUCCESS.Error(),
		Data:   amountIn,
	})
}

// request

func RequestAddLiquidity(c *gin.Context) {
	requestUser := c.MustGet("username").(string)
	var poolAddLiquidityBatchRequest pool.PoolAddLiquidityRequest
	err := c.ShouldBindJSON(&poolAddLiquidityBatchRequest)
	if err != nil {
		c.JSON(http.StatusOK, Result2{
			Errno:  models.ShouldBindJsonErr.Code(),
			ErrMsg: err.Error(),
			Data:   nil,
		})
		return
	}
	poolAddLiquidityBatch, err := pool.ProcessPoolAddLiquidityBatchRequest(&poolAddLiquidityBatchRequest, requestUser)
	if err != nil {
		c.JSON(http.StatusOK, Result2{
			Errno:  models.ProcessPoolAddLiquidityBatchRequestErr.Code(),
			ErrMsg: err.Error(),
			Data:   nil,
		})
		return
	}
	err = pool.Create(poolAddLiquidityBatch)
	if err != nil {
		c.JSON(http.StatusOK, Result2{
			Errno:  models.PoolCreateErr.Code(),
			ErrMsg: err.Error(),
			Data:   nil,
		})
		return
	}
	// success
	c.JSON(http.StatusOK, Result2{
		Errno:  0,
		ErrMsg: models.SUCCESS.Error(),
		Data:   nil,
	})
}

func RequestRemoveLiquidity(c *gin.Context) {
	requestUser := c.MustGet("username").(string)
	var poolRemoveLiquidityBatchRequest pool.PoolRemoveLiquidityRequest
	err := c.ShouldBindJSON(&poolRemoveLiquidityBatchRequest)
	if err != nil {
		c.JSON(http.StatusOK, Result2{
			Errno:  models.ShouldBindJsonErr.Code(),
			ErrMsg: err.Error(),
			Data:   nil,
		})
		return
	}
	poolRemoveLiquidityBatch, err := pool.ProcessPoolRemoveLiquidityBatchRequest(&poolRemoveLiquidityBatchRequest, requestUser)
	if err != nil {
		c.JSON(http.StatusOK, Result2{
			Errno:  models.ProcessPoolRemoveLiquidityBatchRequestErr.Code(),
			ErrMsg: err.Error(),
			Data:   nil,
		})
		return
	}
	err = pool.Create(poolRemoveLiquidityBatch)
	if err != nil {
		c.JSON(http.StatusOK, Result2{
			Errno:  models.PoolCreateErr.Code(),
			ErrMsg: err.Error(),
			Data:   nil,
		})
		return
	}
	c.JSON(http.StatusOK, Result2{
		Errno:  0,
		ErrMsg: models.SUCCESS.Error(),
		Data:   nil,
	})
}

func RequestSwapExactTokenForTokenNoPath(c *gin.Context) {
	requestUser := c.MustGet("username").(string)
	var poolSwapExactTokenForTokenNoPathBatchRequest pool.PoolSwapExactTokenForTokenNoPathRequest
	err := c.ShouldBindJSON(&poolSwapExactTokenForTokenNoPathBatchRequest)
	if err != nil {
		c.JSON(http.StatusOK, Result2{
			Errno:  models.ShouldBindJsonErr.Code(),
			ErrMsg: err.Error(),
			Data:   nil,
		})
		return
	}
	poolSwapExactTokenForTokenNoPathBatch, err := pool.ProcessPoolSwapExactTokenForTokenNoPathBatchRequest(&poolSwapExactTokenForTokenNoPathBatchRequest, requestUser)
	if err != nil {
		c.JSON(http.StatusOK, Result2{
			Errno:  models.ProcessPoolSwapExactTokenForTokenNoPathBatchRequestErr.Code(),
			ErrMsg: err.Error(),
			Data:   nil,
		})
		return
	}
	err = pool.Create(poolSwapExactTokenForTokenNoPathBatch)
	if err != nil {
		c.JSON(http.StatusOK, Result2{
			Errno:  models.PoolCreateErr.Code(),
			ErrMsg: err.Error(),
			Data:   nil,
		})
		return
	}
	c.JSON(http.StatusOK, Result2{
		Errno:  0,
		ErrMsg: models.SUCCESS.Error(),
		Data:   nil,
	})
}

func RequestSwapTokenForExactTokenNoPath(c *gin.Context) {
	requestUser := c.MustGet("username").(string)
	var poolSwapTokenForExactTokenNoPathBatchRequest pool.PoolSwapTokenForExactTokenNoPathRequest
	err := c.ShouldBindJSON(&poolSwapTokenForExactTokenNoPathBatchRequest)
	if err != nil {
		c.JSON(http.StatusOK, Result2{
			Errno:  models.ShouldBindJsonErr.Code(),
			ErrMsg: err.Error(),
			Data:   nil,
		})
		return
	}
	poolSwapTokenForExactTokenNoPathBatch, err := pool.ProcessPoolSwapTokenForExactTokenNoPathBatchRequest(&poolSwapTokenForExactTokenNoPathBatchRequest, requestUser)
	if err != nil {
		c.JSON(http.StatusOK, Result2{
			Errno:  models.ProcessPoolSwapTokenForExactTokenNoPathBatchRequestErr.Code(),
			ErrMsg: err.Error(),
			Data:   nil,
		})
		return
	}
	err = pool.Create(poolSwapTokenForExactTokenNoPathBatch)
	if err != nil {
		c.JSON(http.StatusOK, Result2{
			Errno:  models.PoolCreateErr.Code(),
			ErrMsg: err.Error(),
			Data:   nil,
		})
		return
	}
	c.JSON(http.StatusOK, Result2{
		Errno:  0,
		ErrMsg: models.SUCCESS.Error(),
		Data:   nil,
	})
}

func RequestWithdrawAward(c *gin.Context) {
	requestUser := c.MustGet("username").(string)
	var poolWithdrawAwardBatchRequest pool.PoolWithdrawAwardRequest
	err := c.ShouldBindJSON(&poolWithdrawAwardBatchRequest)
	if err != nil {
		c.JSON(http.StatusOK, Result2{
			Errno:  models.ShouldBindJsonErr.Code(),
			ErrMsg: err.Error(),
			Data:   nil,
		})
		return
	}
	poolWithdrawAwardBatch, err := pool.ProcessPoolWithdrawAwardBatchRequest(&poolWithdrawAwardBatchRequest, requestUser)
	if err != nil {
		c.JSON(http.StatusOK, Result2{
			Errno:  models.ProcessPoolSwapExactTokenForTokenNoPathBatchRequestErr.Code(),
			ErrMsg: err.Error(),
			Data:   nil,
		})
		return
	}
	err = pool.Create(poolWithdrawAwardBatch)
	if err != nil {
		c.JSON(http.StatusOK, Result2{
			Errno:  models.PoolCreateErr.Code(),
			ErrMsg: err.Error(),
			Data:   nil,
		})
		return
	}
	c.JSON(http.StatusOK, Result2{
		Errno:  0,
		ErrMsg: models.SUCCESS.Error(),
		Data:   nil,
	})
}

// batch

func QueryAddLiquidityBatchCount(c *gin.Context) {
	username := c.MustGet("username").(string)
	var count int64
	var err error

	count, err = pool.QueryAddLiquidityBatchCount(username)
	if err != nil {
		c.JSON(http.StatusOK, Result2{
			Errno:  models.QueryAddLiquidityBatchCountErr.Code(),
			ErrMsg: err.Error(),
			Data:   0,
		})
		return
	}

	c.JSON(http.StatusOK, Result2{
		Errno:  0,
		ErrMsg: models.SUCCESS.Error(),
		Data:   count,
	})
}

func QueryAddLiquidityBatch(c *gin.Context) {
	username := c.MustGet("username").(string)
	limit := c.Query("limit")
	offset := c.Query("offset")

	if limit == "" {
		err := errors.New("limit is empty")
		c.JSON(http.StatusOK, Result2{
			Errno:  models.LimitEmptyErr.Code(),
			ErrMsg: err.Error(),
			Data:   &[]pool.PoolAddLiquidityBatchInfo{},
		})
		return
	}
	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		c.JSON(http.StatusOK, Result2{
			Errno:  models.AtoiErr.Code(),
			ErrMsg: err.Error(),
			Data:   &[]pool.PoolAddLiquidityBatchInfo{},
		})
		return
	}
	if limitInt < 0 {
		err := errors.New("limit is less than 0")
		c.JSON(http.StatusOK, Result2{
			Errno:  models.LimitLessThanZeroErr.Code(),
			ErrMsg: err.Error(),
			Data:   &[]pool.PoolAddLiquidityBatchInfo{},
		})
		return
	}
	if offset == "" {
		err := errors.New("offset is empty")
		c.JSON(http.StatusOK, Result2{
			Errno:  models.OffsetEmptyErr.Code(),
			ErrMsg: err.Error(),
			Data:   &[]pool.PoolAddLiquidityBatchInfo{},
		})
		return
	}
	offsetInt, err := strconv.Atoi(offset)
	if err != nil {
		c.JSON(http.StatusOK, Result2{
			Errno:  models.AtoiErr.Code(),
			ErrMsg: err.Error(),
			Data:   &[]pool.PoolAddLiquidityBatchInfo{},
		})
		return
	}
	if offsetInt < 0 {
		err := errors.New("offset is less than 0")
		c.JSON(http.StatusOK, Result2{
			Errno:  models.OffsetLessThanZeroErr.Code(),
			ErrMsg: err.Error(),
			Data:   &[]pool.PoolAddLiquidityBatchInfo{},
		})
		return
	}
	var records *[]pool.PoolAddLiquidityBatchInfo

	records, err = pool.QueryAddLiquidityBatch(username, limitInt, offsetInt)
	if err != nil {
		c.JSON(http.StatusOK, Result2{
			Errno:  models.QueryAddLiquidityBatchErr.Code(),
			ErrMsg: err.Error(),
			Data:   &[]pool.PoolAddLiquidityBatchInfo{},
		})
		return
	}

	c.JSON(http.StatusOK, Result2{
		Errno:  0,
		ErrMsg: models.SUCCESS.Error(),
		Data:   records,
	})
}

func QueryRemoveLiquidityBatchCount(c *gin.Context) {
	username := c.MustGet("username").(string)
	var count int64
	var err error

	count, err = pool.QueryRemoveLiquidityBatchCount(username)
	if err != nil {
		c.JSON(http.StatusOK, Result2{
			Errno:  models.QueryRemoveLiquidityBatchCountErr.Code(),
			ErrMsg: err.Error(),
			Data:   0,
		})
		return
	}

	c.JSON(http.StatusOK, Result2{
		Errno:  0,
		ErrMsg: models.SUCCESS.Error(),
		Data:   count,
	})
}

func QueryRemoveLiquidityBatch(c *gin.Context) {
	username := c.MustGet("username").(string)
	limit := c.Query("limit")
	offset := c.Query("offset")

	if limit == "" {
		err := errors.New("limit is empty")
		c.JSON(http.StatusOK, Result2{
			Errno:  models.LimitEmptyErr.Code(),
			ErrMsg: err.Error(),
			Data:   &[]pool.PoolRemoveLiquidityBatchInfo{},
		})
		return
	}
	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		c.JSON(http.StatusOK, Result2{
			Errno:  models.AtoiErr.Code(),
			ErrMsg: err.Error(),
			Data:   &[]pool.PoolRemoveLiquidityBatchInfo{},
		})
		return
	}
	if limitInt < 0 {
		err := errors.New("limit is less than 0")
		c.JSON(http.StatusOK, Result2{
			Errno:  models.LimitLessThanZeroErr.Code(),
			ErrMsg: err.Error(),
			Data:   &[]pool.PoolRemoveLiquidityBatchInfo{},
		})
		return
	}
	if offset == "" {
		err := errors.New("offset is empty")
		c.JSON(http.StatusOK, Result2{
			Errno:  models.OffsetEmptyErr.Code(),
			ErrMsg: err.Error(),
			Data:   &[]pool.PoolRemoveLiquidityBatchInfo{},
		})
		return
	}
	offsetInt, err := strconv.Atoi(offset)
	if err != nil {
		c.JSON(http.StatusOK, Result2{
			Errno:  models.AtoiErr.Code(),
			ErrMsg: err.Error(),
			Data:   &[]pool.PoolRemoveLiquidityBatchInfo{},
		})
		return
	}
	if offsetInt < 0 {
		err := errors.New("offset is less than 0")
		c.JSON(http.StatusOK, Result2{
			Errno:  models.OffsetLessThanZeroErr.Code(),
			ErrMsg: err.Error(),
			Data:   &[]pool.PoolRemoveLiquidityBatchInfo{},
		})
		return
	}
	var records *[]pool.PoolRemoveLiquidityBatchInfo

	records, err = pool.QueryRemoveLiquidityBatch(username, limitInt, offsetInt)
	if err != nil {
		c.JSON(http.StatusOK, Result2{
			Errno:  models.QueryRemoveLiquidityBatchErr.Code(),
			ErrMsg: err.Error(),
			Data:   &[]pool.PoolRemoveLiquidityBatchInfo{},
		})
		return
	}

	c.JSON(http.StatusOK, Result2{
		Errno:  0,
		ErrMsg: models.SUCCESS.Error(),
		Data:   records,
	})
}

func QuerySwapExactTokenForTokenNoPathBatchCount(c *gin.Context) {
	username := c.MustGet("username").(string)
	var count int64
	var err error

	count, err = pool.QuerySwapExactTokenForTokenNoPathBatchCount(username)
	if err != nil {
		c.JSON(http.StatusOK, Result2{
			Errno:  models.QuerySwapExactTokenForTokenNoPathBatchCountErr.Code(),
			ErrMsg: err.Error(),
			Data:   0,
		})
		return
	}

	c.JSON(http.StatusOK, Result2{
		Errno:  0,
		ErrMsg: models.SUCCESS.Error(),
		Data:   count,
	})
}

func QuerySwapExactTokenForTokenNoPathBatch(c *gin.Context) {
	username := c.MustGet("username").(string)
	limit := c.Query("limit")
	offset := c.Query("offset")

	if limit == "" {
		err := errors.New("limit is empty")
		c.JSON(http.StatusOK, Result2{
			Errno:  models.LimitEmptyErr.Code(),
			ErrMsg: err.Error(),
			Data:   &[]pool.PoolSwapExactTokenForTokenNoPathBatchInfo{},
		})
		return
	}
	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		c.JSON(http.StatusOK, Result2{
			Errno:  models.AtoiErr.Code(),
			ErrMsg: err.Error(),
			Data:   &[]pool.PoolSwapExactTokenForTokenNoPathBatchInfo{},
		})
		return
	}
	if limitInt < 0 {
		err := errors.New("limit is less than 0")
		c.JSON(http.StatusOK, Result2{
			Errno:  models.LimitLessThanZeroErr.Code(),
			ErrMsg: err.Error(),
			Data:   &[]pool.PoolSwapExactTokenForTokenNoPathBatchInfo{},
		})
		return
	}
	if offset == "" {
		err := errors.New("offset is empty")
		c.JSON(http.StatusOK, Result2{
			Errno:  models.OffsetEmptyErr.Code(),
			ErrMsg: err.Error(),
			Data:   &[]pool.PoolSwapExactTokenForTokenNoPathBatchInfo{},
		})
		return
	}
	offsetInt, err := strconv.Atoi(offset)
	if err != nil {
		c.JSON(http.StatusOK, Result2{
			Errno:  models.AtoiErr.Code(),
			ErrMsg: err.Error(),
			Data:   &[]pool.PoolSwapExactTokenForTokenNoPathBatchInfo{},
		})
		return
	}
	if offsetInt < 0 {
		err := errors.New("offset is less than 0")
		c.JSON(http.StatusOK, Result2{
			Errno:  models.OffsetLessThanZeroErr.Code(),
			ErrMsg: err.Error(),
			Data:   &[]pool.PoolSwapExactTokenForTokenNoPathBatchInfo{},
		})
		return
	}
	var records *[]pool.PoolSwapExactTokenForTokenNoPathBatchInfo

	records, err = pool.QuerySwapExactTokenForTokenNoPathBatch(username, limitInt, offsetInt)
	if err != nil {
		c.JSON(http.StatusOK, Result2{
			Errno:  models.QuerySwapExactTokenForTokenNoPathBatchErr.Code(),
			ErrMsg: err.Error(),
			Data:   &[]pool.PoolSwapExactTokenForTokenNoPathBatchInfo{},
		})
		return
	}

	c.JSON(http.StatusOK, Result2{
		Errno:  0,
		ErrMsg: models.SUCCESS.Error(),
		Data:   records,
	})
}

func QuerySwapTokenForExactTokenNoPathBatchCount(c *gin.Context) {
	username := c.MustGet("username").(string)
	var count int64
	var err error

	count, err = pool.QuerySwapTokenForExactTokenNoPathBatchCount(username)
	if err != nil {
		c.JSON(http.StatusOK, Result2{
			Errno:  models.QuerySwapTokenForExactTokenNoPathBatchCountErr.Code(),
			ErrMsg: err.Error(),
			Data:   0,
		})
		return
	}

	c.JSON(http.StatusOK, Result2{
		Errno:  0,
		ErrMsg: models.SUCCESS.Error(),
		Data:   count,
	})
}

func QuerySwapTokenForExactTokenNoPathBatch(c *gin.Context) {
	username := c.MustGet("username").(string)
	limit := c.Query("limit")
	offset := c.Query("offset")

	if limit == "" {
		err := errors.New("limit is empty")
		c.JSON(http.StatusOK, Result2{
			Errno:  models.LimitEmptyErr.Code(),
			ErrMsg: err.Error(),
			Data:   &[]pool.PoolSwapTokenForExactTokenNoPathBatchInfo{},
		})
		return
	}
	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		c.JSON(http.StatusOK, Result2{
			Errno:  models.AtoiErr.Code(),
			ErrMsg: err.Error(),
			Data:   &[]pool.PoolSwapTokenForExactTokenNoPathBatchInfo{},
		})
		return
	}
	if limitInt < 0 {
		err := errors.New("limit is less than 0")
		c.JSON(http.StatusOK, Result2{
			Errno:  models.LimitLessThanZeroErr.Code(),
			ErrMsg: err.Error(),
			Data:   &[]pool.PoolSwapTokenForExactTokenNoPathBatchInfo{},
		})
		return
	}
	if offset == "" {
		err := errors.New("offset is empty")
		c.JSON(http.StatusOK, Result2{
			Errno:  models.OffsetEmptyErr.Code(),
			ErrMsg: err.Error(),
			Data:   &[]pool.PoolSwapTokenForExactTokenNoPathBatchInfo{},
		})
		return
	}
	offsetInt, err := strconv.Atoi(offset)
	if err != nil {
		c.JSON(http.StatusOK, Result2{
			Errno:  models.AtoiErr.Code(),
			ErrMsg: err.Error(),
			Data:   &[]pool.PoolSwapTokenForExactTokenNoPathBatchInfo{},
		})
		return
	}
	if offsetInt < 0 {
		err := errors.New("offset is less than 0")
		c.JSON(http.StatusOK, Result2{
			Errno:  models.OffsetLessThanZeroErr.Code(),
			ErrMsg: err.Error(),
			Data:   &[]pool.PoolSwapTokenForExactTokenNoPathBatchInfo{},
		})
		return
	}
	var records *[]pool.PoolSwapTokenForExactTokenNoPathBatchInfo

	records, err = pool.QuerySwapTokenForExactTokenNoPathBatch(username, limitInt, offsetInt)
	if err != nil {
		c.JSON(http.StatusOK, Result2{
			Errno:  models.QuerySwapTokenForExactTokenNoPathBatchErr.Code(),
			ErrMsg: err.Error(),
			Data:   &[]pool.PoolSwapTokenForExactTokenNoPathBatchInfo{},
		})
		return
	}

	c.JSON(http.StatusOK, Result2{
		Errno:  0,
		ErrMsg: models.SUCCESS.Error(),
		Data:   records,
	})
}

func QueryWithdrawAwardBatchCount(c *gin.Context) {
	username := c.MustGet("username").(string)
	var count int64
	var err error

	count, err = pool.QueryWithdrawAwardBatchCount(username)
	if err != nil {
		c.JSON(http.StatusOK, Result2{
			Errno:  models.QueryWithdrawAwardBatchCountErr.Code(),
			ErrMsg: err.Error(),
			Data:   0,
		})
		return
	}

	c.JSON(http.StatusOK, Result2{
		Errno:  0,
		ErrMsg: models.SUCCESS.Error(),
		Data:   count,
	})
}

func QueryWithdrawAwardBatch(c *gin.Context) {
	username := c.MustGet("username").(string)
	limit := c.Query("limit")
	offset := c.Query("offset")

	if limit == "" {
		err := errors.New("limit is empty")
		c.JSON(http.StatusOK, Result2{
			Errno:  models.LimitEmptyErr.Code(),
			ErrMsg: err.Error(),
			Data:   &[]pool.PoolWithdrawAwardBatchInfo{},
		})
		return
	}
	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		c.JSON(http.StatusOK, Result2{
			Errno:  models.AtoiErr.Code(),
			ErrMsg: err.Error(),
			Data:   &[]pool.PoolWithdrawAwardBatchInfo{},
		})
		return
	}
	if limitInt < 0 {
		err := errors.New("limit is less than 0")
		c.JSON(http.StatusOK, Result2{
			Errno:  models.LimitLessThanZeroErr.Code(),
			ErrMsg: err.Error(),
			Data:   &[]pool.PoolWithdrawAwardBatchInfo{},
		})
		return
	}
	if offset == "" {
		err := errors.New("offset is empty")
		c.JSON(http.StatusOK, Result2{
			Errno:  models.OffsetEmptyErr.Code(),
			ErrMsg: err.Error(),
			Data:   &[]pool.PoolWithdrawAwardBatchInfo{},
		})
		return
	}
	offsetInt, err := strconv.Atoi(offset)
	if err != nil {
		c.JSON(http.StatusOK, Result2{
			Errno:  models.AtoiErr.Code(),
			ErrMsg: err.Error(),
			Data:   &[]pool.PoolWithdrawAwardBatchInfo{},
		})
		return
	}
	if offsetInt < 0 {
		err := errors.New("offset is less than 0")
		c.JSON(http.StatusOK, Result2{
			Errno:  models.OffsetLessThanZeroErr.Code(),
			ErrMsg: err.Error(),
			Data:   &[]pool.PoolWithdrawAwardBatchInfo{},
		})
		return
	}
	var records *[]pool.PoolWithdrawAwardBatchInfo

	records, err = pool.QueryWithdrawAwardBatch(username, limitInt, offsetInt)
	if err != nil {
		c.JSON(http.StatusOK, Result2{
			Errno:  models.QueryWithdrawAwardBatchErr.Code(),
			ErrMsg: err.Error(),
			Data:   &[]pool.PoolWithdrawAwardBatchInfo{},
		})
		return
	}

	c.JSON(http.StatusOK, Result2{
		Errno:  0,
		ErrMsg: models.SUCCESS.Error(),
		Data:   records,
	})
}

// Sync

func AddLiquidity(c *gin.Context) {
	requestUser := c.MustGet("username").(string)
	var poolAddLiquidityBatchRequest pool.PoolAddLiquidityRequest
	err := c.ShouldBindJSON(&poolAddLiquidityBatchRequest)
	if err != nil {
		c.JSON(http.StatusOK, Result2{
			Errno:  models.ShouldBindJsonErr.Code(),
			ErrMsg: err.Error(),
			Data:   new(pool.AddLiquidityResult),
		})
		return
	}

	if requestUser != poolAddLiquidityBatchRequest.Username {
		c.JSON(http.StatusOK, Result2{
			Errno:  models.UsernameNotMatchErr.Code(),
			ErrMsg: "username not match",
			Data:   new(pool.AddLiquidityResult),
		})
		return
	}

	result, err := pool.AddLiquidity(&poolAddLiquidityBatchRequest)
	if err != nil {
		c.JSON(http.StatusOK, Result2{
			Errno:  models.AddLiquidityErr.Code(),
			ErrMsg: err.Error(),
			Data:   new(pool.AddLiquidityResult),
		})
		return
	}

	c.JSON(http.StatusOK, Result2{
		Errno:  0,
		ErrMsg: models.SUCCESS.Error(),
		Data:   result,
	})
}

func RemoveLiquidity(c *gin.Context) {
	requestUser := c.MustGet("username").(string)
	var poolRemoveLiquidityBatchRequest pool.PoolRemoveLiquidityRequest
	err := c.ShouldBindJSON(&poolRemoveLiquidityBatchRequest)
	if err != nil {
		c.JSON(http.StatusOK, Result2{
			Errno:  models.ShouldBindJsonErr.Code(),
			ErrMsg: err.Error(),
			Data:   new(pool.RemoveLiquidityResult),
		})
		return
	}

	if requestUser != poolRemoveLiquidityBatchRequest.Username {
		c.JSON(http.StatusOK, Result2{
			Errno:  models.UsernameNotMatchErr.Code(),
			ErrMsg: "username not match",
			Data:   new(pool.RemoveLiquidityResult),
		})
		return
	}

	result, err := pool.RemoveLiquidity(&poolRemoveLiquidityBatchRequest)
	if err != nil {
		c.JSON(http.StatusOK, Result2{
			Errno:  models.RemoveLiquidityErr.Code(),
			ErrMsg: err.Error(),
			Data:   new(pool.RemoveLiquidityResult),
		})
		return
	}

	c.JSON(http.StatusOK, Result2{
		Errno:  0,
		ErrMsg: models.SUCCESS.Error(),
		Data:   result,
	})
}

func SwapExactTokenForTokenNoPath(c *gin.Context) {
	requestUser := c.MustGet("username").(string)
	var poolSwapExactTokenForTokenNoPathBatchRequest pool.PoolSwapExactTokenForTokenNoPathRequest
	err := c.ShouldBindJSON(&poolSwapExactTokenForTokenNoPathBatchRequest)
	if err != nil {
		c.JSON(http.StatusOK, Result2{
			Errno:  models.ShouldBindJsonErr.Code(),
			ErrMsg: err.Error(),
			Data:   new(pool.SwapExactTokenForTokenNoPathResult),
		})
		return
	}

	if requestUser != poolSwapExactTokenForTokenNoPathBatchRequest.Username {
		c.JSON(http.StatusOK, Result2{
			Errno:  models.UsernameNotMatchErr.Code(),
			ErrMsg: "username not match",
			Data:   new(pool.SwapExactTokenForTokenNoPathResult),
		})
		return
	}

	result, err := pool.SwapExactTokenForTokenNoPath(&poolSwapExactTokenForTokenNoPathBatchRequest)
	if err != nil {
		c.JSON(http.StatusOK, Result2{
			Errno:  models.SwapExactTokenForTokenNoPathErr.Code(),
			ErrMsg: err.Error(),
			Data:   new(pool.SwapExactTokenForTokenNoPathResult),
		})
		return
	}

	c.JSON(http.StatusOK, Result2{
		Errno:  0,
		ErrMsg: models.SUCCESS.Error(),
		Data:   result,
	})
}

func SwapTokenForExactTokenNoPath(c *gin.Context) {
	requestUser := c.MustGet("username").(string)
	var poolSwapTokenForExactTokenNoPathBatchRequest pool.PoolSwapTokenForExactTokenNoPathRequest
	err := c.ShouldBindJSON(&poolSwapTokenForExactTokenNoPathBatchRequest)
	if err != nil {
		c.JSON(http.StatusOK, Result2{
			Errno:  models.ShouldBindJsonErr.Code(),
			ErrMsg: err.Error(),
			Data:   new(pool.SwapTokenForExactTokenNoPathResult),
		})
		return
	}

	if requestUser != poolSwapTokenForExactTokenNoPathBatchRequest.Username {
		c.JSON(http.StatusOK, Result2{
			Errno:  models.UsernameNotMatchErr.Code(),
			ErrMsg: "username not match",
			Data:   new(pool.SwapTokenForExactTokenNoPathResult),
		})
		return
	}

	result, err := pool.SwapTokenForExactTokenNoPath(&poolSwapTokenForExactTokenNoPathBatchRequest)
	if err != nil {
		c.JSON(http.StatusOK, Result2{
			Errno:  models.SwapTokenForExactTokenNoPathErr.Code(),
			ErrMsg: err.Error(),
			Data:   new(pool.SwapTokenForExactTokenNoPathResult),
		})
		return
	}

	c.JSON(http.StatusOK, Result2{
		Errno:  0,
		ErrMsg: models.SUCCESS.Error(),
		Data:   result,
	})
}

func WithdrawAward(c *gin.Context) {
	requestUser := c.MustGet("username").(string)
	var poolWithdrawAwardBatchRequest pool.PoolWithdrawAwardRequest
	err := c.ShouldBindJSON(&poolWithdrawAwardBatchRequest)
	if err != nil {
		c.JSON(http.StatusOK, Result2{
			Errno:  models.ShouldBindJsonErr.Code(),
			ErrMsg: err.Error(),
			Data:   new(pool.WithdrawAwardResult),
		})
		return
	}

	if requestUser != poolWithdrawAwardBatchRequest.Username {
		c.JSON(http.StatusOK, Result2{
			Errno:  models.UsernameNotMatchErr.Code(),
			ErrMsg: "username not match",
			Data:   new(pool.WithdrawAwardResult),
		})
		return
	}

	result, err := pool.WithdrawAward(&poolWithdrawAwardBatchRequest)
	if err != nil {
		c.JSON(http.StatusOK, Result2{
			Errno:  models.WithdrawAwardErr.Code(),
			ErrMsg: err.Error(),
			Data:   new(pool.WithdrawAwardResult),
		})
		return
	}

	c.JSON(http.StatusOK, Result2{
		Errno:  0,
		ErrMsg: models.SUCCESS.Error(),
		Data:   result,
	})
}
