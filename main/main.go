package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/robfig/cron/v3"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"trade/btlLog"
	"trade/config"
	"trade/dao"
	"trade/middleware"
	"trade/routers"
	"trade/routers/RouterSecond"
	"trade/services"
	"trade/services/custodyAccount"
	"trade/task"
	"trade/utils"
)

func main() {
	loadConfig, err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Printf("Failed to load config: %v", err)
		return
	}
	mode := loadConfig.GinConfig.Mode
	if !(mode == gin.DebugMode || mode == gin.ReleaseMode || mode == gin.TestMode) {
		mode = gin.DebugMode
	}
	gin.SetMode(mode)
	utils.PrintAsciiLogoAndInfo()
	if mode != gin.ReleaseMode {
		utils.PrintTitle(false, "Initialize")
	}
	// Initialize the database connection
	if err = middleware.InitMysql(); err != nil {
		log.Printf("Failed to initialize database: %v", err)
		return
	}
	if err = middleware.RedisConnect(); err != nil {
		log.Printf("Failed to initialize redis: %v", err)
		return
	}
	if config.GetLoadConfig().IsAutoMigrate {
		if err = dao.Migrate(); err != nil {
			utils.LogError("AutoMigrate error", err)
			return
		}
	}
	utils.PrintTitle(true, "Check Start")
	if !checkStart() {
		return
	}
	// Setup cron jobs
	services.CheckIfAutoUpdateScheduledTask()
	var jobs []task.Job
	if jobs, err = task.LoadJobs(); err != nil {
		log.Println(err)
		return
	}
	c := cron.New(cron.WithSeconds())
	for _, job := range jobs {
		// Schedule each job using cron
		_, err = c.AddFunc(job.CronExpression, func() {
			task.ExecuteWithLock(job.Name)
		})
		if err != nil {
			log.Printf("Error scheduling job %s: %v\n", job.Name, err)
			continue
		}
	}
	c.Start()
	defer c.Stop() // Ensure cron scheduler is stopped on shutdown
	// Setup HTTP server
	if mode != gin.ReleaseMode {
		utils.PrintTitle(true, "Setup Router")
	}
	go middleware.MonitorDatabaseConnections()

	r := routers.SetupRouter()
	bind := loadConfig.GinConfig.Bind
	port := loadConfig.GinConfig.Port
	if port == "" {
		port = "8080"
	}
	srv := &http.Server{
		Addr:    bind + ":" + port,
		Handler: r,
	}
	// Start HTTP server in a goroutine
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			fmt.Printf("listen: %s\n", err)
		}
	}()
	// Setup HTTP server2
	r2 := RouterSecond.SetupRouter()
	r2bind := "127.0.0.1"
	if config.GetConfig().NetWork == "regtest" {
		r2bind = "0.0.0.0"
	}
	localPort := loadConfig.GinConfig.LocalPort
	if localPort == "" {
		localPort = "10080"
	}
	srv2 := &http.Server{
		Addr:    r2bind + ":" + localPort,
		Handler: r2,
	}

	// Start HTTP server in a goroutine
	go func() {
		if err := srv2.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			fmt.Printf("listen: %s\n", err)
		}
	}()

	// Create a channel to listen to interrupt signals
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	// Block until a signal is received
	sig := <-signalChan
	log.Printf("Received signal: %s", sig)
	// Gracefully shutdown the server, waiting max 30 seconds for current operations to complete
	_, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	// Close Redis connection
	defer func(client *redis.Client) {
		if err = client.Close(); err != nil {
			log.Printf("Failed to close Redis connection: %v", err)
			return
		} else {
			log.Println("Redis connection closed successfully.")
		}
	}(middleware.Client)
	// Close database connection
	var db *sql.DB
	if db, err = middleware.DB.DB(); err != nil {
		log.Println(err)
		return
	}
	defer func(db *sql.DB) {
		if err = db.Close(); err != nil {
			log.Printf("Error closing database: %v", err)
			return
		} else {
			log.Println("Database connection closed successfully.")
		}
	}(db)
	// Perform any other shutdown tasks here
	log.Println("Shutting down the server...")
	os.Exit(0)
}

// Check config
func checkStart() bool {
	cfg := config.GetConfig()
	//检查区块网络配置
	switch cfg.NetWork {
	case "testnet":
		log.Println("Running on testnet")
	case "mainnet":
		log.Println("Running on mainnet")
	case "regtest":
		log.Println("Running on regtest")
	default:
		log.Println("NetWork need set testnet, mainnet or regtest")
		return false
	}
	//加载日志系统
	if err := btlLog.InitBtlLog(); err != nil {
		log.Printf("Failed to initialize btl log: %v", err)
		return false
	}
	//加载托管账户系统
	ctx := context.Background()
	if !custodyAccount.CustodyStart(ctx, cfg) {
		return false
	}
	return true
}
