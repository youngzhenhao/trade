package middleware

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"sync"
	"time"
	"trade/config"
	"trade/utils"
)

var (
	DB   *gorm.DB
	once sync.Once // 使用 sync.Once 确保只初始化一次
)

// InitMysql 初始化数据库连接
func InitMysql() error {
	var err error

	// 使用 sync.Once 确保只初始化一次，避免竞态条件
	once.Do(func() {
		loadConfig, loadErr := config.LoadConfig("config.yaml")
		if loadErr != nil {
			panic("failed to load config: " + loadErr.Error())
		}

		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			loadConfig.GormConfig.Mysql.Username,
			loadConfig.GormConfig.Mysql.Password,
			loadConfig.GormConfig.Mysql.Host,
			loadConfig.GormConfig.Mysql.Port,
			loadConfig.GormConfig.Mysql.DBName,
		)

		gormDB, openErr := gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if openErr != nil {
			utils.LogError("failed to connect database", openErr)
			err = openErr
			return
		}

		sqlDB, dbErr := gormDB.DB()
		if dbErr != nil {
			utils.LogError("failed to get generic database object", dbErr)
			err = dbErr
			return
		}

		// 配置连接池
		sqlDB.SetMaxIdleConns(10)  // 最大空闲连接数
		sqlDB.SetMaxOpenConns(300) // 最大打开连接数
		sqlDB.SetConnMaxIdleTime(5 * time.Minute)
		sqlDB.SetConnMaxLifetime(5 * time.Minute) // 连接的最大存活时间

		DB = gormDB
	})

	return err
}

func MonitorDatabaseConnections() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		sqlDB, _ := DB.DB()
		if err := sqlDB.Ping(); err != nil {
			log.Printf("Database ping failed: %v", err)
		}
	}
}
