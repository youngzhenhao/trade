package services

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"trade/api"
	"trade/models"
	"trade/utils"
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
	Txid    string `json:"txid"`
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
		txid, _ := utils.OutpointToTransactionAndIndex(addrReceiveEvent.Outpoint)
		assetReceives = append(assetReceives, AssetReceive{
			AssetId: addrReceiveEvent.AddrAssetID,
			Txid:    txid,
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

func AssetReceivesToAssetIdAndReceives(assetReceives *[]AssetReceive) *[]AssetIdAndReceive {
	assetIdMapAssetReceives := AssetReceivesToAssetIdMapAssetReceives(assetReceives)
	assetIdAndReceives := AssetIdMapAssetReceivesToAssetIdAndReceives(assetIdMapAssetReceives)
	return assetIdAndReceives
}

func AssetReceivesToUserMapAssetReceives(assetReceives *[]AssetReceive) *map[int]*[]AssetReceive {
	userMapAssetReceives := make(map[int]*[]AssetReceive)
	for _, assetReceive := range *assetReceives {
		balances, ok := userMapAssetReceives[assetReceive.UserId]
		if !ok {
			userMapAssetReceives[assetReceive.UserId] = &[]AssetReceive{assetReceive}
		} else {
			*balances = append(*balances, assetReceive)
		}
	}
	return &userMapAssetReceives
}

func UserMapAssetReceivesToUserAssetReceives(userMapAssetReceives *map[int]*[]AssetReceive) *[]UserAssetReceive {
	var userAssetReceives []UserAssetReceive
	for userId, assetReceives := range *userMapAssetReceives {
		userAssetReceives = append(userAssetReceives, UserAssetReceive{
			UserId:        userId,
			AssetReceives: assetReceives,
		})
	}
	return &userAssetReceives
}

func AssetReceivesToUserAssetReceives(assetReceives *[]AssetReceive) *[]UserAssetReceive {
	userMapAssetReceives := AssetReceivesToUserMapAssetReceives(assetReceives)
	userAssetReceives := UserMapAssetReceivesToUserAssetReceives(userMapAssetReceives)
	return userAssetReceives
}

// GetAllAssetIdAndUserAssetReceives
// @Description: Get all asset id and user asset receives
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

type AssetIdAndUserAssetReceiveAmount struct {
	AssetId                 string                    `json:"asset_id"`
	UserAssetReceiveAmounts *[]UserAssetReceiveAmount `json:"user_asset_receive_amounts"`
}

type AssetIdAndUserAssetReceiveAmountMap struct {
	AssetId                   string       `json:"asset_id"`
	UserAssetReceiveAmountMap *map[int]int `json:"user_asset_receive_amount_map"`
}

type UserAssetReceiveAmount struct {
	UserId             int `json:"user_id"`
	AssetReceiveAmount int `json:"asset_receive_amount"`
}

func AssetReceivesToUserMapAssetReceiveAmount(assetReceives *[]AssetReceive) *map[int]int {
	userMapAssetReceiveAmount := make(map[int]int)
	for _, assetReceive := range *assetReceives {
		balances, ok := userMapAssetReceiveAmount[assetReceive.UserId]
		if !ok || balances == 0 {
			userMapAssetReceiveAmount[assetReceive.UserId] = assetReceive.Amount
		} else {
			userMapAssetReceiveAmount[assetReceive.UserId] += assetReceive.Amount
		}
	}
	return &userMapAssetReceiveAmount
}

func UserMapAssetReceiveAmountToUserAssetReceiveAmount(userMapAssetReceiveAmount *map[int]int) *[]UserAssetReceiveAmount {
	var userAssetReceiveAmount []UserAssetReceiveAmount
	for userId, receiveAmount := range *userMapAssetReceiveAmount {
		userAssetReceiveAmount = append(userAssetReceiveAmount, UserAssetReceiveAmount{
			UserId:             userId,
			AssetReceiveAmount: receiveAmount,
		})
	}
	return &userAssetReceiveAmount
}

func AssetReceivesToUserAssetReceiveAmount(assetReceives *[]AssetReceive) *[]UserAssetReceiveAmount {
	userMapAssetReceives := AssetReceivesToUserMapAssetReceiveAmount(assetReceives)
	userAssetReceiveAmount := UserMapAssetReceiveAmountToUserAssetReceiveAmount(userMapAssetReceives)
	return userAssetReceiveAmount
}

// GetAllAssetIdAndUserAssetReceiveAmount
// @Description: get all asset id and user asset receive amount
func GetAllAssetIdAndUserAssetReceiveAmount() (*[]AssetIdAndUserAssetReceiveAmount, error) {
	var assetIdAndUserAssetReceiveAmount []AssetIdAndUserAssetReceiveAmount
	allAssetReceives, err := GetAllAssetReceives()
	if err != nil {
		return nil, err
	}
	assetIdAndReceives := AssetReceivesToAssetIdAndReceives(allAssetReceives)
	for _, assetIdAndReceive := range *assetIdAndReceives {
		userAssetReceiveAmount := AssetReceivesToUserAssetReceiveAmount(assetIdAndReceive.AssetReceives)
		assetIdAndUserAssetReceiveAmount = append(assetIdAndUserAssetReceiveAmount, AssetIdAndUserAssetReceiveAmount{
			AssetId:                 assetIdAndReceive.AssetId,
			UserAssetReceiveAmounts: userAssetReceiveAmount,
		})
	}
	return &assetIdAndUserAssetReceiveAmount, nil
}

func AssetReceivesToUserAssetReceiveAmountMap(assetReceives *[]AssetReceive) *map[int]int {
	userMapAssetReceives := AssetReceivesToUserMapAssetReceiveAmount(assetReceives)
	return userMapAssetReceives
}

// @dev: Use map
func GetAllAssetIdAndUserAssetReceiveAmountMap() (*[]AssetIdAndUserAssetReceiveAmountMap, error) {
	var assetIdAndUserAssetReceiveAmount []AssetIdAndUserAssetReceiveAmountMap
	allAssetReceives, err := GetAllAssetReceives()
	if err != nil {
		return nil, err
	}
	assetIdAndReceives := AssetReceivesToAssetIdAndReceives(allAssetReceives)
	for _, assetIdAndReceive := range *assetIdAndReceives {
		userAssetReceiveAmountMap := AssetReceivesToUserAssetReceiveAmountMap(assetIdAndReceive.AssetReceives)
		assetIdAndUserAssetReceiveAmount = append(assetIdAndUserAssetReceiveAmount, AssetIdAndUserAssetReceiveAmountMap{
			AssetId:                   assetIdAndReceive.AssetId,
			UserAssetReceiveAmountMap: userAssetReceiveAmountMap,
		})
	}
	return &assetIdAndUserAssetReceiveAmount, nil
}

func AssetReceivesToAddressAmountMap(assetReceiveEvents *[]models.AddrReceiveEvent, opMapAddress *map[string]string) *map[string]*AssetIdAndAmount {
	addressAmountMap := make(map[string]*AssetIdAndAmount)
	for _, assetReceive := range *assetReceiveEvents {
		op := assetReceive.Outpoint
		address, ok := (*opMapAddress)[op]
		if !ok {
			continue
		}
		_, ok = addressAmountMap[address]
		if !ok {
			addressAmountMap[address] = &AssetIdAndAmount{
				AssetId: assetReceive.AddrAssetID,
			}
		}
		if (*(addressAmountMap[address])).AssetId == assetReceive.AddrAssetID {
			(*(addressAmountMap[address])).Amount += assetReceive.AddrAmount
		}
	}
	return &addressAmountMap
}

func AssetReceiveEventsToOutpointSlice(addrReceiveEvents *[]models.AddrReceiveEvent) []string {
	var ops []string
	for _, event := range *addrReceiveEvents {
		ops = append(ops, event.Outpoint)
	}
	return ops
}

// AllAssetReceivesToAddressAmountMap
// @Description: All asset receives to address amount map
func AllAssetReceivesToAddressAmountMap(network models.Network) (*map[string]*AssetIdAndAmount, error) {
	allAssetReceiveEvents, err := GetAllAddrReceiveEvents()
	if err != nil {
		return nil, err
	}
	ops := AssetReceiveEventsToOutpointSlice(allAssetReceiveEvents)
	opMapAddress, err := api.GetAddressesByOutpointSlice(network, ops)
	if err != nil {
		return nil, err
	}
	addressAmountMap := AssetReceivesToAddressAmountMap(allAssetReceiveEvents, &opMapAddress)
	return addressAmountMap, nil
}

func SetAddrReceivesEvents(receives *[]models.AddrReceiveEventSetRequest) error {
	userByte := sha256.Sum256([]byte(AdminUploadUserName))
	username := hex.EncodeToString(userByte[:])
	userId, err := NameToId(username)
	if err != nil {
		// @dev: Admin upload user does not exist
		password, _ := hashPassword(username)
		if password == "" {
			password = username
		}
		err = CreateUser(&models.User{
			Username: username,
			Password: password,
		})
		if err != nil {
			return err
		}
		userId, err = NameToId(username)
		if err != nil {
			return err
		}
	}
	addrReceiveEvents := ProcessAddrReceiveEventsSetRequest(userId, receives)
	err = CreateOrUpdateAddrReceiveEvents(addrReceiveEvents)
	if err != nil {
		return err
	}
	return nil
}

// GetAndSetAddrReceivesEvents
// @Description: Get and set addr receives events
func GetAndSetAddrReceivesEvents(deviceId string) error {
	receives, err := api.AddrReceivesAndGetEventSetRequests(deviceId)
	if err != nil {
		return err
	}
	if receives == nil || len(*receives) == 0 {
		return nil
	}
	err = SetAddrReceivesEvents(receives)
	if err != nil {
		return nil
	}
	return nil
}
