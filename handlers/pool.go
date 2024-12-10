package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
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
			Data:   new(pool.PoolInfo),
		})
		return
	}
	c.JSON(http.StatusOK, Result2{
		Errno:  0,
		ErrMsg: models.SUCCESS.Error(),
		Data:   poolInfo,
	})
}

func QueryShareRecords(c *gin.Context) {
	//tokenA
	//tokenB
	//limit
	//offset
}

func QuerySwapRecords(c *gin.Context) {

}

func QueryWithdrawAwardRecords(c *gin.Context) {

}
