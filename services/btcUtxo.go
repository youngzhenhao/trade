package services

import (
	"errors"
	"trade/middleware"
	"trade/models"
	"trade/utils"
)

func UnspentUtxosToBtcUtxos(username string, requests *[]models.UnspentUtxo, opExists map[string]bool) (btcUtxos *[]models.BtcUtxo, err error) {
	if requests == nil {
		return nil, errors.New("requests is nil")
	}
	var _btcUtxos []models.BtcUtxo

	for _, request := range *requests {
		if opExists[request.Outpoint] {
			continue
		}
		btcUtxo := models.BtcUtxo{
			Username: username,
			UnspentUtxo: models.UnspentUtxo{
				AddressType:   request.AddressType,
				Address:       request.Address,
				AmountSat:     request.AmountSat,
				PkScript:      request.PkScript,
				Outpoint:      request.Outpoint,
				Confirmations: request.Confirmations,
			},
		}
		_btcUtxos = append(_btcUtxos, btcUtxo)
	}

	return &_btcUtxos, nil

}

func SetBtcUtxo(username string, requests *[]models.UnspentUtxo) (err error) {

	if requests == nil {
		return errors.New("requests is nil")
	}

	var allOps []string
	for _, request := range *requests {
		allOps = append(allOps, request.Outpoint)
	}

	var existOps []string

	err = middleware.DB.Table("btc_utxos").
		Where("outpoint IN (?)", allOps).
		Pluck("outpoint", &existOps).Error

	opExists := make(map[string]bool)
	for _, op := range existOps {
		opExists[op] = true
	}

	btcUtxos, err := UnspentUtxosToBtcUtxos(username, requests, opExists)
	if err != nil {
		return utils.AppendErrorInfo(err, "UnspentUtxosToBtcUtxos")
	}

	//fmt.Println(utils.ValueJsonString(btcUtxos))

	if len(*btcUtxos) == 0 {
		return nil
	}

	err = middleware.DB.Create(btcUtxos).Error
	if err != nil {
		return utils.AppendErrorInfo(err, "Create BtcUtxo")
	}

	return nil
}
