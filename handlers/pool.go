package handlers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"trade/models"
	"trade/services/pool"
)

// request

func AddLiquidity(c *gin.Context) {
	requestUser := c.MustGet("username").(string)
	var poolAddLiquidityBatchRequest pool.PoolAddLiquidityBatchRequest
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

func RemoveLiquidity(c *gin.Context) {
	requestUser := c.MustGet("username").(string)
	var poolRemoveLiquidityBatchRequest pool.PoolRemoveLiquidityBatchRequest
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

func SwapExactTokenForTokenNoPath(c *gin.Context) {
	requestUser := c.MustGet("username").(string)
	var poolSwapExactTokenForTokenNoPathBatchRequest pool.PoolSwapExactTokenForTokenNoPathBatchRequest
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

func SwapTokenForExactTokenNoPath(c *gin.Context) {
	requestUser := c.MustGet("username").(string)
	var poolSwapTokenForExactTokenNoPathBatchRequest pool.PoolSwapTokenForExactTokenNoPathBatchRequest
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
			Errno:  models.ProcessPoolSwapExactTokenForTokenNoPathBatchRequestErr.Code(),
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

func WithdrawAward(c *gin.Context) {
	requestUser := c.MustGet("username").(string)
	var poolWithdrawAwardBatchRequest pool.PoolWithdrawAwardBatchRequest
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

// calc

func CalcAddLiquidity(c *gin.Context) {

}

func CalcRemoveLiquidity(c *gin.Context) {

}

func CalcSwapExactTokenForTokenNoPath(c *gin.Context) {

}

func CalcSwapTokenForExactTokenNoPath(c *gin.Context) {

}

func CalcWithdrawAward(c *gin.Context) {

}

// query

func QueryPoolInfo(c *gin.Context) {
	_ = c.MustGet("username").(string)
	tokenA := c.Query("token_a")
	tokenB := c.Query("token_b")
	poolInfo, err := pool.QueryPoolInfo(tokenA, tokenB)
	if err != nil {
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
