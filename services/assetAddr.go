package services

import (
	"errors"
	"trade/models"
)

func GetAssetAddrsByUserId(userId int) (*[]models.AssetAddr, error) {
	return ReadAssetAddrsByUserId(userId)
}

func ProcessAssetAddrSetRequest(userId int, assetAddrSetRequest *models.AssetAddrSetRequest) *models.AssetAddr {
	var assetAddr models.AssetAddr
	assetAddr = models.AssetAddr{
		Encoded:          assetAddrSetRequest.Encoded,
		AssetId:          assetAddrSetRequest.AssetId,
		AssetType:        assetAddrSetRequest.AssetType,
		Amount:           assetAddrSetRequest.Amount,
		GroupKey:         assetAddrSetRequest.GroupKey,
		ScriptKey:        assetAddrSetRequest.ScriptKey,
		InternalKey:      assetAddrSetRequest.InternalKey,
		TapscriptSibling: assetAddrSetRequest.TapscriptSibling,
		TaprootOutputKey: assetAddrSetRequest.TaprootOutputKey,
		ProofCourierAddr: assetAddrSetRequest.ProofCourierAddr,
		AssetVersion:     assetAddrSetRequest.AssetVersion,
		DeviceID:         assetAddrSetRequest.DeviceID,
		UserId:           userId,
	}
	return &assetAddr
}

func IsAssetAddrChanged(assetAddrByAddrEncoded *models.AssetAddr, old *models.AssetAddr) bool {
	if assetAddrByAddrEncoded == nil || old == nil {
		return true
	}
	if assetAddrByAddrEncoded.Encoded != old.Encoded {
		return true
	}
	if assetAddrByAddrEncoded.AssetId != old.AssetId {
		return true
	}
	if assetAddrByAddrEncoded.AssetType != old.AssetType {
		return true
	}
	if assetAddrByAddrEncoded.Amount != old.Amount {
		return true
	}
	if assetAddrByAddrEncoded.GroupKey != old.GroupKey {
		return true
	}
	if assetAddrByAddrEncoded.ScriptKey != old.ScriptKey {
		return true
	}
	if assetAddrByAddrEncoded.InternalKey != old.InternalKey {
		return true
	}
	if assetAddrByAddrEncoded.TapscriptSibling != old.TapscriptSibling {
		return true
	}
	if assetAddrByAddrEncoded.TaprootOutputKey != old.TaprootOutputKey {
		return true
	}
	if assetAddrByAddrEncoded.ProofCourierAddr != old.ProofCourierAddr {
		return true
	}
	if assetAddrByAddrEncoded.AssetVersion != old.AssetVersion {
		return true
	}
	if assetAddrByAddrEncoded.DeviceID != old.DeviceID {
		return true
	}
	if assetAddrByAddrEncoded.UserId != old.UserId {
		return true
	}
	return false
}

func CheckAssetAddrIfUpdate(assetAddr *models.AssetAddr) (*models.AssetAddr, error) {
	if assetAddr == nil {
		return nil, errors.New("nil asset addr")
	}
	assetAddrByAddrEncoded, err := ReadAssetAddrByAddrEncoded(assetAddr.Encoded)
	if err != nil {
		return assetAddr, nil
	}
	if !IsAssetAddrChanged(assetAddrByAddrEncoded, assetAddr) {
		return assetAddrByAddrEncoded, nil
	}
	assetAddrByAddrEncoded.Encoded = assetAddr.Encoded
	assetAddrByAddrEncoded.AssetId = assetAddr.AssetId
	assetAddrByAddrEncoded.AssetType = assetAddr.AssetType
	assetAddrByAddrEncoded.Amount = assetAddr.Amount
	assetAddrByAddrEncoded.GroupKey = assetAddr.GroupKey
	assetAddrByAddrEncoded.ScriptKey = assetAddr.ScriptKey
	assetAddrByAddrEncoded.InternalKey = assetAddr.InternalKey
	assetAddrByAddrEncoded.TapscriptSibling = assetAddr.TapscriptSibling
	assetAddrByAddrEncoded.TaprootOutputKey = assetAddr.TaprootOutputKey
	assetAddrByAddrEncoded.ProofCourierAddr = assetAddr.ProofCourierAddr
	assetAddrByAddrEncoded.AssetVersion = assetAddr.AssetVersion
	assetAddrByAddrEncoded.DeviceID = assetAddr.DeviceID
	assetAddrByAddrEncoded.UserId = assetAddr.UserId
	return assetAddrByAddrEncoded, nil
}

func CreateOrUpdateAssetAddr(addr *models.AssetAddr) (err error) {
	var assetAddr *models.AssetAddr
	assetAddr, err = CheckAssetAddrIfUpdate(addr)
	return UpdateAssetAddr(assetAddr)
}

func GetAssetAddrsByScriptKey(scriptKey string) (*[]models.AssetAddr, error) {
	return ReadAssetAddrsByScriptKey(scriptKey)
}
