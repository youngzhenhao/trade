package config

import (
	"gopkg.in/yaml.v3"
	"os"
	"trade/utils"
)

type Config struct {
	NetWork   string `yaml:"network" json:"network"`
	GinConfig struct {
		Bind      string `yaml:"bind" json:"bind"`
		Port      string `yaml:"port" json:"port"`
		Mode      string `yaml:"mode" json:"mode"`
		LocalPort string `yaml:"local_port" json:"local_port"`
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
		Host                 string `yaml:"host" json:"host"`
		Port                 string `yaml:"port" json:"port"`
		Username             string `yaml:"username" json:"username"`
		Password             string `yaml:"password" json:"password"`
		DB                   int    `yaml:"db" json:"db"`
		ExpirationTimeMinute int    `yaml:"expiration_time_minute" json:"expiration_time_minute"`
		RedisSetTimeMinute   int    `yaml:"redis_set_time_minute" json:"redis_set_time_minute"`
	} `yaml:"redis" json:"redis"`
	RouterDisable struct {
		Login                 bool `yaml:"login" json:"login"`
		FairLaunch            bool `yaml:"fair_launch" json:"fair_launch"`
		Fee                   bool `yaml:"fee" json:"fee"`
		CustodyAccount        bool `yaml:"custody_account" json:"custody_account"`
		Ping                  bool `yaml:"ping" json:"ping"`
		Proof                 bool `yaml:"proof" json:"proof"`
		Ido                   bool `yaml:"ido" json:"ido"`
		Snapshot              bool `yaml:"snapshot" json:"snapshot"`
		BtcBalance            bool `yaml:"btc_balance" json:"btc_balance"`
		AssetTransfer         bool `yaml:"asset_transfer" json:"asset_transfer"`
		Bitcoind              bool `yaml:"bitcoind" json:"bitcoind"`
		Shell                 bool `yaml:"shell" json:"shell"`
		AddrReceive           bool `yaml:"addr_receive" json:"addr_receive"`
		BatchTransfer         bool `yaml:"batch_transfer" json:"batch_transfer"`
		AssetAddr             bool `yaml:"asset_addr" json:"asset_addr"`
		AssetLock             bool `yaml:"asset_lock" json:"asset_lock"`
		ValidateToken         bool `yaml:"validate_token" json:"validate_token"`
		AssetBalance          bool `yaml:"asset_balance" json:"asset_balance"`
		AssetBurn             bool `yaml:"asset_burn" json:"asset_burn"`
		AssetLocalMint        bool `yaml:"asset_local_mint" json:"asset_local_mint"`
		User                  bool `yaml:"user" json:"user"`
		AssetRecommend        bool `yaml:"asset_recommend" json:"asset_recommend"`
		FairLaunchFollow      bool `yaml:"fair_launch_follow" json:"fair_launch_follow"`
		AssetLocalMintHistory bool `yaml:"asset_local_mint_history" json:"asset_local_mint_history"`
		AssetManagedUtxo      bool `yaml:"asset_managed_utxo" json:"asset_managed_utxo"`
		LogFileUpload         bool `yaml:"log_file_upload" json:"log_file_upload"`
		AccountAsset          bool `yaml:"account_asset" json:"account_asset"`
		AssetGroup            bool `yaml:"asset_group" json:"asset_group"`
		NftTransfer           bool `yaml:"nft_transfer" json:"nft_transfer"`
		NftInfo               bool `yaml:"nft_info" json:"nft_info"`
	} `yaml:"router_disable" json:"router_disable"`
	ApiConfig struct {
		Lnd struct {
			Host         string `yaml:"host" json:"host"`
			Port         int    `yaml:"port" json:"port"`
			Dir          string `yaml:"dir" json:"dir"`
			TlsCertPath  string `yaml:"tls_cert_path" json:"tls_cert_path"`
			MacaroonPath string `yaml:"macaroon_path" json:"macaroon_path"`
		} `yaml:"lnd" json:"lnd"`
		Tapd struct {
			Host         string `yaml:"host" json:"host"`
			Port         int    `yaml:"port" json:"port"`
			Dir          string `yaml:"dir" json:"dir"`
			TlsCertPath  string `yaml:"tls_cert_path" json:"tls_cert_path"`
			MacaroonPath string `yaml:"macaroon_path" json:"macaroon_path"`
			UniverseHost string `yaml:"universe_host" json:"universe_host"`
		} `yaml:"tapd" json:"tapd"`
		Litd struct {
			Host         string `yaml:"host" json:"host"`
			Port         int    `yaml:"port" json:"port"`
			Dir          string `yaml:"dir" json:"dir"`
			TlsCertPath  string `yaml:"tls_cert_path" json:"tls_cert_path"`
			MacaroonPath string `yaml:"macaroon_path" json:"macaroon_path"`
		} `yaml:"litd" json:"litd"`
		Bitcoind struct {
			Mainnet struct {
				Ip           string `yaml:"ip" json:"ip"`
				Port         int    `yaml:"port" json:"port"`
				Wallet       string `yaml:"wallet" json:"wallet"`
				RpcUser      string `yaml:"rpc_user" json:"rpc_user"`
				RpcPasswd    string `yaml:"rpc_passwd" json:"rpc_passwd"`
				HttpPostMode bool   `yaml:"http_post_mode" json:"http_post_mode"`
				DisableTLS   bool   `yaml:"disable_tls" json:"disable_tls"`
			} `yaml:"mainnet" json:"mainnet"`
			Testnet struct {
				Ip           string `yaml:"ip" json:"ip"`
				Port         int    `yaml:"port" json:"port"`
				Wallet       string `yaml:"wallet" json:"wallet"`
				RpcUser      string `yaml:"rpc_user" json:"rpc_user"`
				RpcPasswd    string `yaml:"rpc_passwd" json:"rpc_passwd"`
				HttpPostMode bool   `yaml:"http_post_mode" json:"http_post_mode"`
				DisableTLS   bool   `yaml:"disable_tls" json:"disable_tls"`
			} `yaml:"testnet" json:"testnet"`
			Regtest struct {
				Ip           string `yaml:"ip" json:"ip"`
				Port         int    `yaml:"port" json:"port"`
				Wallet       string `yaml:"wallet" json:"wallet"`
				RpcUser      string `yaml:"rpc_user" json:"rpc_user"`
				RpcPasswd    string `yaml:"rpc_passwd" json:"rpc_passwd"`
				HttpPostMode bool   `yaml:"http_post_mode" json:"http_post_mode"`
				DisableTLS   bool   `yaml:"disable_tls" json:"disable_tls"`
			} `yaml:"regtest" json:"regtest"`
		} `yaml:"bitcoind" json:"bitcoind"`
		CustodyAccount struct {
			MacaroonDir string `yaml:"macaroon_dir" json:"macaroon_dir"`
		} `yaml:"custody_account" json:"custody_account"`
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
