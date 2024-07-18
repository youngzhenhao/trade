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
		return nil, errors.New("nil addr receive")
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

func GetAllAddrReceiveEvents() (*[]models.AddrReceiveEvent, error) {
	return ReadAllAddrReceiveEvents()
}

func GetAddrReceiveEventsByAssetId(assetId string) (*[]models.AddrReceiveEvent, error) {
	return ReadAddrReceiveEventsByAssetId(assetId)
}

type AssetReceive struct {
	AssetId string `json:"asset_id"`
	Encoded string `json:"encoded"`
	Amount  int    `json:"amount"`
	UserId  int    `json:"user_id"`
}

type UserAssetReceive struct {
	UserId        int             `json:"user_id"`
	AssetReceives *[]AssetReceive `json:"asset_receives"`
}

type AssetIdAndReceive struct {
	AssetId       string          `json:"asset_id"`
	AssetReceives *[]AssetReceive `json:"asset_receives"`
}

type AssetIdAndUserAssetReceive struct {
	AssetId           string              `json:"asset_id"`
	UserAssetReceives *[]UserAssetReceive `json:"user_asset_receives"`
}

func GetAllAssetReceives() (*[]AssetReceive, error) {
	allAddrReceiveEvents, err := GetAllAddrReceiveEvents()
	if err != nil {
		return nil, err
	}
	assetReceives := AddrReceiveEventsToAssetReceives(allAddrReceiveEvents)
	return assetReceives, nil
}

func AddrReceiveEventsToAssetReceives(allAddrReceiveEvents *[]models.AddrReceiveEvent) *[]AssetReceive {
	var assetReceives []AssetReceive
	for _, addrReceiveEvent := range *allAddrReceiveEvents {
		assetReceives = append(assetReceives, AssetReceive{
			AssetId: addrReceiveEvent.AddrAssetID,
			Encoded: addrReceiveEvent.AddrEncoded,
			Amount:  addrReceiveEvent.AddrAmount,
			UserId:  addrReceiveEvent.UserID,
		})
	}
	return &assetReceives
}

func AssetReceivesToAssetIdMapAssetReceives(assetReceives *[]AssetReceive) *map[string]*[]AssetReceive {
	AssetIdMapAssetReceives := make(map[string]*[]AssetReceive)
	for _, assetReceive := range *assetReceives {
		receives, ok := AssetIdMapAssetReceives[assetReceive.AssetId]
		if !ok {
			AssetIdMapAssetReceives[assetReceive.AssetId] = &[]AssetReceive{assetReceive}
		} else {
			*receives = append(*receives, assetReceive)
		}
	}
	return &AssetIdMapAssetReceives
}

func AssetIdMapAssetReceivesToAssetIdAndReceives(AssetIdMapAssetReceives *map[string]*[]AssetReceive) *[]AssetIdAndReceive {
	var assetIdAndReceives []AssetIdAndReceive
	for assetId, assetReceives := range *AssetIdMapAssetReceives {
		assetIdAndReceives = append(assetIdAndReceives, AssetIdAndReceive{
			AssetId:       assetId,
			AssetReceives: assetReceives,
		})
	}
	return &assetIdAndReceives
}

// TODO: Test
func AssetReceivesToAssetIdAndReceives(assetReceives *[]AssetReceive) *[]AssetIdAndReceive {
	assetIdMapAssetReceives := AssetReceivesToAssetIdMapAssetReceives(assetReceives)
	assetIdAndReceives := AssetIdMapAssetReceivesToAssetIdAndReceives(assetIdMapAssetReceives)
	return assetIdAndReceives
}

func AssetReceivesToUserAssetReceives(userAssetReceives *[]AssetReceive) *[]UserAssetReceive {
	// TODO
	return nil
}

func GetAllAssetIdAndUserAssetReceives() (*[]AssetIdAndUserAssetReceive, error) {
	var assetIdAndUserAssetReceives []AssetIdAndUserAssetReceive
	allAssetReceives, err := GetAllAssetReceives()
	if err != nil {
		return nil, err
	}
	assetIdAndReceives := AssetReceivesToAssetIdAndReceives(allAssetReceives)
	for _, assetIdAndReceive := range *assetIdAndReceives {
		userAssetReceives := AssetReceivesToUserAssetReceives(assetIdAndReceive.AssetReceives)
		assetIdAndUserAssetReceives = append(assetIdAndUserAssetReceives, AssetIdAndUserAssetReceive{
			AssetId:           assetIdAndReceive.AssetId,
			UserAssetReceives: userAssetReceives,
		})
	}
	return &assetIdAndUserAssetReceives, nil
}
