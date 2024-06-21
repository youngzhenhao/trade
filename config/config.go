package config

import (
	"gopkg.in/yaml.v3"
	"os"
	"trade/utils"
)

type Config struct {
	GinConfig struct {
		Bind string `yaml:"bind" json:"bind"`
		Port string `yaml:"port" json:"port"`
		Mode string `yaml:"mode" json:"mode"`
	} `yaml:"gin_config" json:"gin_config"`
	GormConfig struct {
		Mysql struct {
			Host     string `yaml:"host" json:"host"`
			Port     string `yaml:"port" json:"port"`
			Username string `yaml:"username" json:"username"`
			Password string `yaml:"password" json:"password"`
			DBName   string `yaml:"dbname" json:"dbname"`
		} `yaml:"mysql" json:"mysql"`
	} `yaml:"gorm_config" json:"gorm_config"`
	Redis struct {
		Host     string `yaml:"host" json:"host"`
		Port     string `yaml:"port" json:"port"`
		Username string `yaml:"username" json:"username"`
		Password string `yaml:"password" json:"password"`
		DB       int    `yaml:"db" json:"db"`
	} `yaml:"redis" json:"redis"`
	RouterDisable struct {
		Login          bool `yaml:"login" json:"login"`
		FairLaunch     bool `yaml:"fair_launch" json:"fair_launch"`
		Fee            bool `yaml:"fee" json:"fee"`
		CustodyAccount bool `yaml:"custody_account" json:"custody_account"`
		Ping           bool `yaml:"ping" json:"ping"`
		Proof          bool `yaml:"proof" json:"proof"`
		Snapshot       bool `yaml:"snapshot" json:"snapshot"`
	} `yaml:"router_disable" json:"router_disable"`
	ApiConfig struct {
		Lnd struct {
			Host         string `yaml:"host" json:"host"`
			Port         int    `yaml:"port" json:"port"`
			Dir          string `yaml:"dir" json:"dir"`
			TlsCertPath  string `yaml:"tls_cert_path" json:"tlsCertPath"`
			MacaroonPath string `yaml:"macaroon_path" json:"macaroonPath"`
		} `yaml:"lnd" json:"lnd"`
		Tapd struct {
			Host         string `yaml:"host" json:"host"`
			Port         int    `yaml:"port" json:"port"`
			Dir          string `yaml:"dir" json:"dir"`
			TlsCertPath  string `yaml:"tls_cert_path" json:"tlsCertPath"`
			MacaroonPath string `yaml:"macaroon_path" json:"macaroonPath"`
		} `yaml:"tapd" json:"tapd"`
		Litd struct {
			Host         string `yaml:"host" json:"host"`
			Port         int    `yaml:"port" json:"port"`
			Dir          string `yaml:"dir" json:"dir"`
			TlsCertPath  string `yaml:"tls_cert_path" json:"tlsCertPath"`
			MacaroonPath string `yaml:"macaroon_path" json:"macaroonPath"`
		} `yaml:"litd" json:"litd"`
		Bitcoind struct {
			Host         string `yaml:"host" json:"host"`
			Port         int    `yaml:"port" json:"port"`
			RpcUser      string `yaml:"rpcuser" json:"rpcUser"`
			RpcPasswd    string `yaml:"rpcpasswd" json:"rpcPasswd"`
			HTTPPostMode bool   `yaml:"http_post_mode" json:"HTTPPostMode"`
			DisableTLS   bool   `yaml:"disable_tls" json:"disableTLS"`
		} `yaml:"bitcoind" json:"bitcoind"`
		CustodyAccount struct {
			MacaroonDir string `yaml:"macaroon_dir" json:"macaroon_dir"`
		} `yaml:"custody_account" json:"custodyAccount"`
	} `yaml:"api_config" json:"api_config"`
	FairLaunchConfig struct {
		EstimateSmartFeeRateBlocks int  `yaml:"estimate_smart_fee_rate_blocks" json:"estimate_smart_fee_rate_blocks"`
		IsAutoUpdateFeeRate        bool `yaml:"is_auto_update_fee_rate" json:"is_auto_update_fee_rate"`
		MaxNumberOfMint            int  `yaml:"max_number_of_mint" json:"max_number_of_mint"`
	} `yaml:"fair_launch_config" json:"fair_launch_config"`
	AdminUser                 BasicAuth `yaml:"admin_user" json:"admin_user"`
	FrpsServer                string    `yaml:"frps_server" json:"frps_server"`
	IsAutoMigrate             bool      `yaml:"is_auto_migrate" json:"is_auto_migrate"`
	IsAutoUpdateScheduledTask bool      `yaml:"is_auto_update_scheduled_task" json:"is_auto_update_scheduled_task"`
}

type BasicAuth struct {
	Username string `yaml:"username" json:"username"`
	Password string `yaml:"password" json:"password"`
}

var (
	config Config
)

func GetConfig() *Config {
	return &config
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func GetLoadConfig() *Config {
	loadConfig, err := LoadConfig("config.yaml")
	if err != nil {
		utils.LogError("[ERROR] Failed to load config", err)
		return &Config{}
	}
	return loadConfig
}
