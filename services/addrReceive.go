package services

import (
	"errors"
	"trade/models"
)

func ProcessAddrReceiveEventsSetRequest(userId int, addrReceiveEventsSetRequest *[]models.AddrReceiveEventSetRequest) *[]models.AddrReceiveEvent {
	var addrReceiveEvents []models.AddrReceiveEvent
	for _, event := range *addrReceiveEventsSetRequest {
		addrReceiveEvents = append(addrReceiveEvents, models.AddrReceiveEvent{
			CreationTimeUnixSeconds: event.CreationTimeUnixSeconds,
			AddrEncoded:             event.Addr.Encoded,
			AddrAssetID:             event.Addr.AssetID,
			AddrAmount:              event.Addr.Amount,
			AddrScriptKey:           event.Addr.ScriptKey,
			AddrInternalKey:         event.Addr.InternalKey,
			AddrTaprootOutputKey:    event.Addr.TaprootOutputKey,
			AddrProofCourierAddr:    event.Addr.ProofCourierAddr,
			EventStatus:             event.Status,
			Outpoint:                event.Outpoint,
			UtxoAmtSat:              event.UtxoAmtSat,
			ConfirmationHeight:      event.ConfirmationHeight,
			HasProof:                event.HasProof,
			DeviceID:                event.DeviceID,
			UserID:                  userId,
		})
	}
	return &addrReceiveEvents
}

func GetAddrReceiveEventsByUserId(userId int) (*[]models.AddrReceiveEvent, error) {
	return ReadAddrReceiveEventsByUserId(userId)
}

func ProcessGetAddrReceiveEvents(addrReceiveEvents *[]models.AddrReceiveEvent) *[]models.AddrReceiveEventSetRequest {
	var addrReceiveEventsSetRequest []models.AddrReceiveEventSetRequest
	for _, event := range *addrReceiveEvents {
		addrReceiveEventsSetRequest = append(addrReceiveEventsSetRequest, models.AddrReceiveEventSetRequest{
			CreationTimeUnixSeconds: event.CreationTimeUnixSeconds,
			Addr: models.AddrReceiveEventSetRequestAddr{
				Encoded:          event.AddrEncoded,
				AssetID:          event.AddrAssetID,
				Amount:           event.AddrAmount,
				ScriptKey:        event.AddrScriptKey,
				InternalKey:      event.AddrInternalKey,
				TaprootOutputKey: event.AddrTaprootOutputKey,
				ProofCourierAddr: event.AddrProofCourierAddr,
			},
			Status:             event.EventStatus,
			Outpoint:           event.Outpoint,
			UtxoAmtSat:         event.UtxoAmtSat,
			ConfirmationHeight: event.ConfirmationHeight,
			HasProof:           event.HasProof,
			DeviceID:           event.DeviceID,
		})
	}
	return &addrReceiveEventsSetRequest
}

func GetAddrReceiveEventsProcessedOriginByUserId(userId int) (*[]models.AddrReceiveEventSetRequest, error) {
	addrReceiveEvents, err := ReadAddrReceiveEventsByUserId(userId)
	if err != nil {
		return nil, err
	}
	return ProcessGetAddrReceiveEvents(addrReceiveEvents), nil
}

func IsAddrReceiveEventChanged(addrReceiveEventByAddrEncoded *models.AddrReceiveEvent, old *models.AddrReceiveEvent) bool {
	if addrReceiveEventByAddrEncoded == nil || old == nil {
		return true
	}
	if addrReceiveEventByAddrEncoded.CreationTimeUnixSeconds != old.CreationTimeUnixSeconds {
		return true
	}
	if addrReceiveEventByAddrEncoded.AddrEncoded != old.AddrEncoded {
		return true
	}
	if addrReceiveEventByAddrEncoded.AddrAssetID != old.AddrAssetID {
		return true
	}
	if addrReceiveEventByAddrEncoded.AddrAmount != old.AddrAmount {
		return true
	}
	if addrReceiveEventByAddrEncoded.AddrScriptKey != old.AddrScriptKey {
		return true
	}
	if addrReceiveEventByAddrEncoded.AddrInternalKey != old.AddrInternalKey {
		return true
	}
	if addrReceiveEventByAddrEncoded.AddrTaprootOutputKey != old.AddrTaprootOutputKey {
		return true
	}
	if addrReceiveEventByAddrEncoded.AddrProofCourierAddr != old.AddrProofCourierAddr {
		return true
	}
	if addrReceiveEventByAddrEncoded.EventStatus != old.EventStatus {
		return true
	}
	if addrReceiveEventByAddrEncoded.Outpoint != old.Outpoint {
		return true
	}
	if addrReceiveEventByAddrEncoded.UtxoAmtSat != old.UtxoAmtSat {
		return true
	}
	if addrReceiveEventByAddrEncoded.ConfirmationHeight != old.ConfirmationHeight {
		return true
	}
	if addrReceiveEventByAddrEncoded.HasProof != old.HasProof {
		return true
	}
	if addrReceiveEventByAddrEncoded.DeviceID != old.DeviceID {
		return true
	}
	if addrReceiveEventByAddrEncoded.UserID != old.UserID {
		return true
	}
	return false
}

func CheckAddrReceiveEventIfUpdate(addrReceiveEvent *models.AddrReceiveEvent) (*models.AddrReceiveEvent, error) {
	if addrReceiveEvent == nil {
		return nil, errors.New("nil addr receive event")
	}
	addrReceiveEventByAddrEncoded, err := ReadAddrReceiveEventByAddrEncoded(addrReceiveEvent.AddrEncoded)
	if err != nil {
		return addrReceiveEvent, nil
	}
	if !IsAddrReceiveEventChanged(addrReceiveEventByAddrEncoded, addrReceiveEvent) {
		return addrReceiveEventByAddrEncoded, nil
	}
	addrReceiveEventByAddrEncoded.CreationTimeUnixSeconds = addrReceiveEvent.CreationTimeUnixSeconds
	addrReceiveEventByAddrEncoded.AddrEncoded = addrReceiveEvent.AddrEncoded
	addrReceiveEventByAddrEncoded.AddrAssetID = addrReceiveEvent.AddrAssetID
	addrReceiveEventByAddrEncoded.AddrAmount = addrReceiveEvent.AddrAmount
	addrReceiveEventByAddrEncoded.AddrScriptKey = addrReceiveEvent.AddrScriptKey
	addrReceiveEventByAddrEncoded.AddrInternalKey = addrReceiveEvent.AddrInternalKey
	addrReceiveEventByAddrEncoded.AddrTaprootOutputKey = addrReceiveEvent.AddrTaprootOutputKey
	addrReceiveEventByAddrEncoded.AddrProofCourierAddr = addrReceiveEvent.AddrProofCourierAddr
	addrReceiveEventByAddrEncoded.EventStatus = addrReceiveEvent.EventStatus
	addrReceiveEventByAddrEncoded.Outpoint = addrReceiveEvent.Outpoint
	addrReceiveEventByAddrEncoded.UtxoAmtSat = addrReceiveEvent.UtxoAmtSat
	addrReceiveEventByAddrEncoded.ConfirmationHeight = addrReceiveEvent.ConfirmationHeight
	addrReceiveEventByAddrEncoded.HasProof = addrReceiveEvent.HasProof
	addrReceiveEventByAddrEncoded.DeviceID = addrReceiveEvent.DeviceID
	addrReceiveEventByAddrEncoded.UserID = addrReceiveEvent.UserID
	return addrReceiveEventByAddrEncoded, nil
}

func CreateOrUpdateAddrReceiveEvents(addrReceiveEvents *[]models.AddrReceiveEvent) (err error) {
	var addrReceives []models.AddrReceiveEvent
	var addrReceive *models.AddrReceiveEvent
	for _, addrReceiveEvent := range *addrReceiveEvents {
		addrReceive, err = CheckAddrReceiveEventIfUpdate(&addrReceiveEvent)
		if err != nil {
			return err
		}
		addrReceives = append(addrReceives, *addrReceive)
	}
	return UpdateAddrReceiveEvents(&addrReceives)
}
