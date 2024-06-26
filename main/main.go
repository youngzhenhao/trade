package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
	"trade/config"
	"trade/dao"
	"trade/middleware"
	"trade/routers"
	"trade/services"
	"trade/task"
	"trade/utils"
)

func main() {
	loadConfig, err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
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
	if err := middleware.InitMysql(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	if err := middleware.RedisConnect(); err != nil {
		log.Fatalf("Failed to initialize redis: %v", err)
	}
	if config.GetLoadConfig().IsAutoMigrate {
		err = dao.Migrate()
	}
	if err != nil {
		utils.LogError("AutoMigrate error", err)
		return
	}
	utils.PrintTitle(true, "Check Start")
	if !checkStart() {
		return
	}
	// Setup cron jobs
	services.CheckIfAutoUpdateScheduledTask()
	jobs, err := task.LoadJobs()
	if err != nil {
		log.Fatal(err)
	}
	c := cron.New(cron.WithSeconds())
	for _, job := range jobs {
		// Schedule each job using cron
		_, err := c.AddFunc(job.CronExpression, func() {
			task.ExecuteWithLock(job.Name)
		})
		if err != nil {
			log.Printf("Error scheduling job %s: %v\n", job.Name, err)
		}
	}
	c.Start()
	defer c.Stop() // Ensure cron scheduler is stopped on shutdown
	// Setup HTTP server
	if mode != gin.ReleaseMode {
		utils.PrintTitle(true, "Setup Router")
	}
	r := routers.SetupRouter()
	/**
	for _, routeInfo := range r.Routes() {
		fmt.Println(routeInfo.Method, routeInfo.Path)
	}
	*/
	bind := loadConfig.GinConfig.Bind
	port := loadConfig.GinConfig.Port
	if port == "" {
		port = "8080"
	}
	utils.PrintTitle(true, "Run Router")
	err = r.Run(fmt.Sprintf("%s:%s", bind, port))
	if err != nil {
		return
	}
	// Create a channel to listen to interrupt signals
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	// Start HTTP server in a goroutine
	go func() {
		if err := r.Run(fmt.Sprintf(":%s", port)); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()
	// Block until a signal is received
	sig := <-signalChan
	log.Printf("Received signal: %s", sig)
	// Gracefully shutdown the server, waiting max 30 seconds for current operations to complete
	_, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	// Close Redis connection
	if err := middleware.Client.Close(); err != nil {
		log.Printf("Failed to close Redis connection: %v", err)
	} else {
		log.Println("Redis connection closed successfully.")
	}
	// Close database connection
	if db, err := middleware.DB.DB(); err == nil {
		if err := db.Close(); err != nil {
			log.Printf("Error closing database: %v", err)
		} else {
			log.Println("Database connection closed successfully.")
		}
	}
	// Perform any other shutdown tasks here
	log.Println("Shutting down the server...")
	os.Exit(0)
}

// check config
func checkStart() bool {
	cfg := config.GetConfig()
	if cfg.ApiConfig.CustodyAccount.MacaroonDir == "" {
		log.Println("Custody account MacaroonDir is not set")
		return false
	}
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

	fmt.Println("Custody account MacaroonDir is set:", cfg.ApiConfig.CustodyAccount.MacaroonDir)
	// 检测admin账户
	if !services.CheckAdminAccount() {
		log.Println("Admin account is not set")
		return false
	}
	return true
}
