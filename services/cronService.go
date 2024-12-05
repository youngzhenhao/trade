package services

import (
	"crypto/sha256"
	"encoding/hex"
	"strconv"
	"strings"
	"time"
	"trade/api"
	"trade/btlLog"
	"trade/config"
	"trade/middleware"
	"trade/models"
	"trade/services/satBackQueue"
	"trade/utils"
)

type CronService struct{}

// @dev: Auto
func CheckIfAutoUpdateScheduledTask() {
	if config.GetLoadConfig().IsAutoUpdateScheduledTask {
		err := CreateFairLaunchProcessions()
		if err != nil {
			btlLog.ScheduledTask.Info("%v", err)
		}
		err = CreateSnapshotProcessions()
		if err != nil {
			btlLog.ScheduledTask.Info("%v", err)
		}
		err = CreateSetTransfersAndReceives()
		if err != nil {
			btlLog.ScheduledTask.Info("%v", err)
		}
		err = CreateNftPresaleProcessions()
		if err != nil {
			btlLog.ScheduledTask.Info("%v", err)
		}
		err = CreatePushQueueProcessions()
		if err != nil {
			btlLog.ScheduledTask.Info("%v", err)
		}

	}
}

// CreateFairLaunchProcessions
// @dev: Use this to update scheduled task table
func CreateFairLaunchProcessions() (err error) {
	return CreateOrUpdateScheduledTasks(&[]models.ScheduledTask{
		{
			Name:           "ProcessFairLaunchNoPay",
			CronExpression: "*/20 * * * * *",
			FunctionName:   "ProcessFairLaunchNoPay",
			Package:        "services",
		}, {
			Name:           "ProcessFairLaunchPaidPending",
			CronExpression: "*/20 * * * * *",
			FunctionName:   "ProcessFairLaunchPaidPending",
			Package:        "services",
		}, {
			Name:           "ProcessFairLaunchPaidNoIssue",
			CronExpression: "0 */2 * * * *",
			FunctionName:   "ProcessFairLaunchPaidNoIssue",
			Package:        "services",
		}, {
			Name:           "ProcessFairLaunchIssuedPending",
			CronExpression: "*/20 * * * * *",
			FunctionName:   "ProcessFairLaunchIssuedPending",
			Package:        "services",
		}, {
			Name:           "ProcessFairLaunchReservedSentPending",
			CronExpression: "*/20 * * * * *",
			FunctionName:   "ProcessFairLaunchReservedSentPending",
			Package:        "services",
		}, {
			Name:           "FairLaunchMint",
			CronExpression: "*/20 * * * * *",
			FunctionName:   "FairLaunchMint",
			Package:        "services",
		}, {
			Name:           "SendFairLaunchAsset",
			CronExpression: "0 */5 * * * *",
			FunctionName:   "SendFairLaunchAsset",
			Package:        "services",
		},
		{
			Name:           "SnapshotToZipLast",
			CronExpression: "0 */5 * * * *",
			FunctionName:   "SnapshotToZipLast",
			Package:        "services",
		},
		{
			Name:           "UpdateFairLaunchIncomesSatAmountByTxids",
			CronExpression: "*/20 * * * * *",
			FunctionName:   "UpdateFairLaunchIncomesSatAmountByTxids",
			Package:        "services",
		},
	})
}

func TaskCountRecordByRedis(name string) error {
	var record string
	var count int
	var err error
	record, err = middleware.RedisGet(name)
	if err != nil {
		// @dev: no value has been set
		err = middleware.RedisSet(name, "1"+","+utils.GetTimeNow(), 6*time.Minute)
		if err != nil {
			return err
		}
		return nil
	}
	split := strings.Split(record, ",")
	count, err = strconv.Atoi(split[0])
	if err != nil {
		return err
	}
	err = middleware.RedisSet(name, strconv.Itoa(count+1)+","+utils.GetTimeNow(), 6*time.Minute)
	if err != nil {
		return err
	}
	return nil
}

func (cs *CronService) FairLaunchIssuance() {
	FairLaunchIssuance()
	err := TaskCountRecordByRedis("FairLaunchIssuance")
	if err != nil {
		return
	}
}

func (cs *CronService) SnapshotToZipLast() {
	SnapshotToZipLast()
	err := TaskCountRecordByRedis("SnapshotToZipLast")
	if err != nil {
		return
	}
}

func (cs *CronService) ProcessFairLaunchNoPay() {
	ProcessFairLaunchNoPay()
	err := TaskCountRecordByRedis("ProcessFairLaunchNoPay")
	if err != nil {
		return
	}
}

func (cs *CronService) ProcessFairLaunchPaidPending() {
	ProcessFairLaunchPaidPending()
	err := TaskCountRecordByRedis("ProcessFairLaunchPaidPending")
	if err != nil {
		return
	}
}

func (cs *CronService) ProcessFairLaunchPaidNoIssue() {
	ProcessFairLaunchPaidNoIssue()
	err := TaskCountRecordByRedis("ProcessFairLaunchPaidNoIssue")
	if err != nil {
		return
	}
}

func (cs *CronService) ProcessFairLaunchIssuedPending() {
	ProcessFairLaunchIssuedPending()
	err := TaskCountRecordByRedis("ProcessFairLaunchIssuedPending")
	if err != nil {
		return
	}
}

func (cs *CronService) ProcessFairLaunchReservedSentPending() {
	ProcessFairLaunchReservedSentPending()
	err := TaskCountRecordByRedis("ProcessFairLaunchReservedSentPending")
	if err != nil {
		return
	}
}

func (cs *CronService) FairLaunchMint() {
	FairLaunchMint()
	err := TaskCountRecordByRedis("FairLaunchMint")
	if err != nil {
		return
	}
}

func (cs *CronService) SendFairLaunchAsset() {
	SendFairLaunchAsset()
	err := TaskCountRecordByRedis("SendFairLaunchAsset")
	if err != nil {
		return
	}
}

func (cs *CronService) RemoveMintedInventories() {
	RemoveMintedInventories()
	err := TaskCountRecordByRedis("RemoveMintedInventories")
	if err != nil {
		return
	}
}

func CreateSnapshotProcessions() (err error) {
	return CreateOrUpdateScheduledTasks(&[]models.ScheduledTask{
		{
			Name:           "SnapshotToZipLast",
			CronExpression: "0 0 */12 * * *",
			FunctionName:   "SnapshotToZipLast",
			Package:        "services",
		},
	})
}

func (cs *CronService) ListAndSetAssetTransfers() {
	network, err := api.NetworkStringToNetwork(config.GetLoadConfig().NetWork)
	if err != nil {
		return
	}
	userByte := sha256.Sum256([]byte(AdminUploadUserName))
	username := hex.EncodeToString(userByte[:])
	if len(username) < 16 {
		return
	}
	deviceId := username[:16]
	err = ListAndSetAssetTransfers(network, deviceId)
	if err != nil {
		return
	}
}

func (cs *CronService) GetAndSetAddrReceivesEvents() {
	userByte := sha256.Sum256([]byte(AdminUploadUserName))
	username := hex.EncodeToString(userByte[:])
	if len(username) < 16 {
		return
	}
	deviceId := username[:16]
	err := GetAndSetAddrReceivesEvents(deviceId)
	if err != nil {
		return
	}
}

func CreateSetTransfersAndReceives() (err error) {
	return CreateOrUpdateScheduledTasks(&[]models.ScheduledTask{
		{
			Name:           "ListAndSetAssetTransfers",
			CronExpression: "0 */3 * * * *",
			FunctionName:   "ListAndSetAssetTransfers",
			Package:        "services",
		},
		{
			Name:           "GetAndSetAddrReceivesEvents",
			CronExpression: "0 */3 * * * *",
			FunctionName:   "GetAndSetAddrReceivesEvents",
			Package:        "services",
		},
	})
}

func (cs *CronService) UpdateFairLaunchIncomesSatAmountByTxids() {
	network, err := api.NetworkStringToNetwork(config.GetLoadConfig().NetWork)
	if err != nil {
		return
	}
	err = UpdateFairLaunchIncomesSatAmountByTxids(network)
	if err != nil {
		return
	}
}

func (cs *CronService) ProcessNftPresaleBoughtNotPay() {
	ProcessNftPresaleBoughtNotPay()
	err := TaskCountRecordByRedis("ProcessNftPresaleBoughtNotPay")
	if err != nil {
		return
	}
}

func (cs *CronService) ProcessNftPresalePaidPending() {
	ProcessNftPresalePaidPending()
	err := TaskCountRecordByRedis("ProcessNftPresalePaidPending")
	if err != nil {
		return
	}
}

func (cs *CronService) ProcessNftPresalePaidNotSend() {
	ProcessNftPresalePaidNotSend()
	err := TaskCountRecordByRedis("ProcessNftPresalePaidNotSend")
	if err != nil {
		return
	}
}

func (cs *CronService) ProcessNftPresaleSentPending() {
	ProcessNftPresaleSentPending()
	err := TaskCountRecordByRedis("ProcessNftPresaleSentPending")
	if err != nil {
		return
	}
}

func CreateNftPresaleProcessions() (err error) {
	return CreateOrUpdateScheduledTasks(&[]models.ScheduledTask{
		{
			Name:           "ProcessNftPresaleBoughtNotPay",
			CronExpression: "*/20 * * * * *",
			FunctionName:   "ProcessNftPresaleBoughtNotPay",
			Package:        "services",
		},
		{
			Name:           "ProcessNftPresalePaidPending",
			CronExpression: "*/20 * * * * *",
			FunctionName:   "ProcessNftPresalePaidPending",
			Package:        "services",
		},
		{
			Name:           "ProcessNftPresalePaidNotSend",
			CronExpression: "0 */5 * * * *",
			FunctionName:   "ProcessNftPresalePaidNotSend",
			Package:        "services",
		},
		{
			Name:           "ProcessNftPresaleSentPending",
			CronExpression: "*/20 * * * * *",
			FunctionName:   "ProcessNftPresaleSentPending",
			Package:        "services",
		},
	})
}

func CreatePushQueueProcessions() (err error) {
	return CreateOrUpdateScheduledTasks(&[]models.ScheduledTask{
		{
			Name:           "GetAndPushClaimAsset",
			CronExpression: "*/30 * * * * *",
			FunctionName:   "GetAndPushClaimAsset",
			Package:        "services",
		},
		{
			Name:           "GetAndPushPurchasePresaleNFT",
			CronExpression: "*/30 * * * * *",
			FunctionName:   "GetAndPushPurchasePresaleNFT",
			Package:        "services",
		},
	})
}

func (cs *CronService) GetAndPushClaimAsset() {
	satBackQueue.GetAndPushClaimAsset()
}

func (cs *CronService) GetAndPushPurchasePresaleNFT() {
	satBackQueue.GetAndPushPurchasePresaleNFT()
}
