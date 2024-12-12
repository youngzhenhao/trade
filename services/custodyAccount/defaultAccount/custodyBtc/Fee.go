package custodyBtc

import (
	"gorm.io/gorm"
	"trade/models/custodyModels"
	"trade/services/custodyAccount/account"
)

func PayFee(Db *gorm.DB, usr *account.UserInfo, amount float64, balanceId uint) error {
	_, err := LessBtcBalance(Db, usr, amount, balanceId, custodyModels.ChangeTypeBtcFee)
	return err
}
