package services

import (
	"errors"
	"trade/models"
)

func GetAssetLocksByUserId(userId int) (*[]models.AssetLock, error) {
	return ReadAssetLocksByUserId(userId)
}

func ProcessAssetLockSetRequest(userId int, assetLockSetRequest *models.AssetLockSetRequest) *models.AssetLock {
	var assetLock models.AssetLock
	assetLock = models.AssetLock{
		AssetId:          assetLockSetRequest.AssetId,
		AssetName:        assetLockSetRequest.AssetName,
		AssetType:        assetLockSetRequest.AssetType,
		LockAmount:       assetLockSetRequest.LockAmount,
		LockTime:         assetLockSetRequest.LockTime,
		RelativeLockTime: assetLockSetRequest.RelativeLockTime,
		HashLock:         assetLockSetRequest.HashLock,
		Invoice:          assetLockSetRequest.Invoice,
		DeviceId:         assetLockSetRequest.DeviceId,
		UserId:           userId,
	}
	return &assetLock
}

func IsAssetLockChanged(assetLockByInvoice *models.AssetLock, old *models.AssetLock) bool {
	if assetLockByInvoice == nil || old == nil {
		return true
	}
	if assetLockByInvoice.AssetId != old.AssetId {
		return true
	}
	if assetLockByInvoice.AssetName != old.AssetName {
		return true
	}
	if assetLockByInvoice.AssetType != old.AssetType {
		return true
	}
	if assetLockByInvoice.LockAmount != old.LockAmount {
		return true
	}
	if assetLockByInvoice.LockTime != old.LockTime {
		return true
	}
	if assetLockByInvoice.RelativeLockTime != old.RelativeLockTime {
		return true
	}
	if assetLockByInvoice.HashLock != old.HashLock {
		return true
	}
	if assetLockByInvoice.Invoice != old.Invoice {
		return true
	}
	if assetLockByInvoice.DeviceId != old.DeviceId {
		return true
	}
	if assetLockByInvoice.UserId != old.UserId {
		return true
	}
	return false
}

func CheckAssetLockIfUpdate(assetLock *models.AssetLock) (*models.AssetLock, error) {
	if assetLock == nil {
		return nil, errors.New("nil asset lock")
	}
	assetLockByInvoice, err := ReadAssetLockByInvoice(assetLock.Invoice)
	if err != nil {
		return assetLock, nil
	}
	if !IsAssetLockChanged(assetLockByInvoice, assetLock) {
		return assetLockByInvoice, nil
	}
	assetLockByInvoice.AssetId = assetLock.AssetId
	assetLockByInvoice.AssetName = assetLock.AssetName
	assetLockByInvoice.AssetType = assetLock.AssetType
	assetLockByInvoice.LockAmount = assetLock.LockAmount
	assetLockByInvoice.LockTime = assetLock.LockTime
	assetLockByInvoice.RelativeLockTime = assetLock.RelativeLockTime
	assetLockByInvoice.HashLock = assetLock.HashLock
	assetLockByInvoice.Invoice = assetLock.Invoice
	assetLockByInvoice.DeviceId = assetLock.DeviceId
	assetLockByInvoice.UserId = assetLock.UserId
	return assetLockByInvoice, nil
}

func CreateOrUpdateAssetLock(lock *models.AssetLock) (err error) {
	var assetLock *models.AssetLock
	assetLock, err = CheckAssetLockIfUpdate(lock)
	return UpdateAssetLock(assetLock)
}
