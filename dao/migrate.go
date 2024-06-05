package dao

import (
	"trade/middleware"
	"trade/models"
)

func Migrate() error {
	var err error
	if err = middleware.DB.AutoMigrate(&models.Account{}); err != nil {
		return err
	}
	if err = middleware.DB.AutoMigrate(&models.Balance{}); err != nil {
		return err
	}
	if err = middleware.DB.AutoMigrate(&models.BalanceExt{}); err != nil {
		return err
	}
	if err = middleware.DB.AutoMigrate(&models.User{}); err != nil {
		return err
	}
	if err = middleware.DB.AutoMigrate(&models.ScheduledTask{}); err != nil {
		return err
	}
	if err = middleware.DB.AutoMigrate(&models.Invoice{}); err != nil {
		return err
	}
	if err = middleware.DB.AutoMigrate(&models.FairLaunchInfo{}); err != nil {
		return err
	}
	if err = middleware.DB.AutoMigrate(&models.FairLaunchMintedInfo{}); err != nil {
		return err
	}
	if err = middleware.DB.AutoMigrate(&models.FairLaunchMintedUserInfo{}); err != nil {
		return err
	}
	if err = middleware.DB.AutoMigrate(&models.FairLaunchInventoryInfo{}); err != nil {
		return err
	}
	if err = middleware.DB.AutoMigrate(&models.FeeRateInfo{}); err != nil {
		return err
	}
	if err = middleware.DB.AutoMigrate(&models.AssetIssuance{}); err != nil {
		return err
	}
	if err = middleware.DB.AutoMigrate(&models.PayInside{}); err != nil {
		return err
	}
	return err
}
