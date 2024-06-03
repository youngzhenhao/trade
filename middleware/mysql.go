package middleware

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
	"trade/config"
	"trade/utils"
)

var (
	DB *gorm.DB
)

func InitMysql() error {
	loadConfig, err := config.LoadConfig("config.yaml")
	if err != nil {
		panic("failed to load config: " + err.Error())
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		loadConfig.GormConfig.Mysql.Username,
		loadConfig.GormConfig.Mysql.Password,
		loadConfig.GormConfig.Mysql.Host,
		loadConfig.GormConfig.Mysql.Port,
		loadConfig.GormConfig.Mysql.DBName,
	)

	gormDB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		utils.LogError("failed to connect database", err)
		return err
	}

	// Get generic database object sql.DB to use its functions
	sqlDB, err := gormDB.DB()
	if err != nil {
		utils.LogError("failed to get generic database object", err)
		return err
	}

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(10)

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(100)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqlDB.SetConnMaxLifetime(time.Hour)

	DB = gormDB
	return nil
}
