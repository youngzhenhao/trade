package services

import (
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
