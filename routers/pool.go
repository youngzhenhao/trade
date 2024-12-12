package routers

import (
	"github.com/gin-gonic/gin"
	"trade/handlers"
	"trade/middleware"
)

func SetupPoolRouter(router *gin.Engine) *gin.Engine {
	pool := router.Group("/pool")
	pool.Use(middleware.AuthMiddleware())
	query := pool.Group("/query")
	{
		query.GET("/pool_info", handlers.QueryPoolInfo)
		query.GET("/share_records/count", handlers.QueryShareRecordsCount)
		query.GET("/share_records", handlers.QueryShareRecords)
		query.GET("/swap_records/count", handlers.QuerySwapRecordsCount)
		query.GET("/swap_records", handlers.QuerySwapRecords)
		query.GET("lp_award_balance", handlers.QueryUserLpAwardBalance)
		query.GET("/withdraw_award_records/count", handlers.QueryWithdrawAwardRecordsCount)
		query.GET("/withdraw_award_records", handlers.QueryWithdrawAwardRecords)
		query.GET("/quote", handlers.QueryQuote)
	}
	calc := pool.Group("/calc")
	{
		calc.POST("/add_liquidity", handlers.CalcAddLiquidity)
		calc.POST("/remove_liquidity", handlers.CalcRemoveLiquidity)
		calc.POST("/swap_exact_token_for_token_no_path", handlers.CalcSwapExactTokenForTokenNoPath)
		calc.POST("/swap_token_for_exact_token_no_path", handlers.CalcSwapTokenForExactTokenNoPath)
	}
	request := pool.Group("/request")
	{
		request.POST("/add_liquidity", handlers.AddLiquidity)
		request.POST("/remove_liquidity", handlers.RemoveLiquidity)
		request.POST("/swap_exact_token_for_token_no_path", handlers.SwapExactTokenForTokenNoPath)
		request.POST("/swap_token_for_exact_token_no_path", handlers.SwapTokenForExactTokenNoPath)
		request.POST("/withdraw_award", handlers.WithdrawAward)
	}
	batch := pool.Group("/batch")
	{
		batch.GET("/add_liquidity/count", handlers.QueryAddLiquidityBatchCount)
		batch.GET("/add_liquidity", handlers.QueryAddLiquidityBatch)
		batch.GET("/remove_liquidity/count", handlers.QueryRemoveLiquidityBatchCount)
		batch.GET("/remove_liquidity", handlers.QueryRemoveLiquidityBatch)
		batch.GET("/swap_exact_token_for_token_no_path/count", handlers.QuerySwapExactTokenForTokenNoPathBatchCount)
		batch.GET("/swap_exact_token_for_token_no_path", handlers.QuerySwapExactTokenForTokenNoPathBatch)
		batch.GET("/swap_token_for_exact_token_no_path/count", handlers.QuerySwapTokenForExactTokenNoPathBatchCount)
		batch.GET("/swap_token_for_exact_token_no_path", handlers.QuerySwapTokenForExactTokenNoPathBatch)
		batch.GET("/withdraw_award/count", handlers.QueryWithdrawAwardBatchCount)
		batch.GET("/withdraw_award", handlers.QueryWithdrawAwardBatch)
	}
	return router
}
