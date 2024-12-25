package services

import (
	"errors"
	"trade/middleware"
	"trade/models"
	"trade/utils"
)

func UnspentUtxosToBtcUtxos(username string, requests *[]models.UnspentUtxo, opExists map[string]bool) (btcUtxos *[]models.BtcUtxo, opDelete map[string]bool, err error) {
	if requests == nil {
		return nil, nil, errors.New("requests is nil")
	}
	var _btcUtxos []models.BtcUtxo

	requestOp := make(map[string]bool)

	for _, request := range *requests {
		requestOp[request.Outpoint] = true
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

	opDelete = make(map[string]bool)

	for opExist := range opExists {
		if !requestOp[opExist] {
			opDelete[opExist] = true
		}
	}

	return &_btcUtxos, opDelete, nil

}

func BtcUtxosToBtcUtxoHistories(btcUtxos *[]models.BtcUtxo) (btcUtxoHistories *[]models.BtcUtxoHistory) {
	if btcUtxos == nil {
		return nil
	}

	var _btcUtxoHistories []models.BtcUtxoHistory

	for _, btcUtxo := range *btcUtxos {
		btcUtxoHistory := models.BtcUtxoHistory{
			Username: btcUtxo.Username,
			UnspentUtxo: models.UnspentUtxo{
				AddressType:   btcUtxo.AddressType,
				Address:       btcUtxo.Address,
				AmountSat:     btcUtxo.AmountSat,
				PkScript:      btcUtxo.PkScript,
				Outpoint:      btcUtxo.Outpoint,
				Confirmations: btcUtxo.Confirmations,
			},
		}
		_btcUtxoHistories = append(_btcUtxoHistories, btcUtxoHistory)
	}

	return &_btcUtxoHistories

}

func SetBtcUtxo(username string, requests *[]models.UnspentUtxo) (err error) {

	if requests == nil {
		return errors.New("requests is nil")
	}

	var allOps []string
	for _, request := range *requests {
		allOps = append(allOps, request.Outpoint)
	}

	var dbExistOps []string

	err = middleware.DB.Table("btc_utxos").
		Where("username = ?", username).
		Pluck("outpoint", &dbExistOps).Error

	opExists := make(map[string]bool)
	for _, op := range dbExistOps {
		opExists[op] = true
	}

	btcUtxos, opDeleteMap, err := UnspentUtxosToBtcUtxos(username, requests, opExists)
	if err != nil {
		return utils.AppendErrorInfo(err, "UnspentUtxosToBtcUtxos")
	}

	//fmt.Printf("btcUtxos: %v\n", utils.ValueJsonString(btcUtxos))

	var opDelete []string
	for op := range opDeleteMap {
		opDelete = append(opDelete, op)
	}

	err = middleware.DB.Where("outpoint IN (?)", opDelete).Delete(&models.BtcUtxo{}).Error
	if err != nil {
		return utils.AppendErrorInfo(err, "Delete BtcUtxo")
	}

	if len(*btcUtxos) == 0 {
		return nil
	}

	//fmt.Printf("opDelete: %v\n", utils.ValueJsonString(opDelete))

	btcUtxoHistories := BtcUtxosToBtcUtxoHistories(btcUtxos)

	err = middleware.DB.Create(btcUtxos).Error
	if err != nil {
		return utils.AppendErrorInfo(err, "Create BtcUtxo")
	}

	err = middleware.DB.Create(btcUtxoHistories).Error
	if err != nil {
		return utils.AppendErrorInfo(err, "Create BtcUtxoHistory")
	}

	return nil
}
