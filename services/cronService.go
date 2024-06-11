package services

import (
	"fmt"
	"strconv"
	"strings"
	"time"
	"trade/config"
	"trade/middleware"
	"trade/models"
	"trade/utils"
)

type CronService struct{}

// @dev: Auto
func CheckIfAutoUpdateScheduledTask() {
	if config.GetLoadConfig().IsAutoUpdateScheduledTask {
		err := CreateFairLaunchProcessions()
		if err != nil {
			ScheduledTask.Info("%v", err)
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

func CreateFairLaunchIssuance() (err error) {
	return CreateScheduledTask(&models.ScheduledTask{
		Name:           "FairLaunchIssuance",
		CronExpression: "0 */1 * * * *",
		FunctionName:   "FairLaunchIssuance",
		Package:        "services",
	})
}

func CreateFairLaunchMint() (err error) {
	return CreateScheduledTask(&models.ScheduledTask{
		Name:           "FairLaunchMint",
		CronExpression: "*/20 * * * * *",
		FunctionName:   "FairLaunchMint",
		Package:        "services",
	})
}

func CreateSendFairLaunchAsset() (err error) {
	return CreateScheduledTask(&models.ScheduledTask{
		Name:           "SendFairLaunchAsset",
		CronExpression: "0 */5 * * * *",
		FunctionName:   "SendFairLaunchAsset",
		Package:        "services",
	})
}

// Deprecated: Use CreateFairLaunchProcessions instead
func CreateFairLaunchScheduledTasks() {
	err := CreateFairLaunchProcessions()
	if err != nil {
		FairLaunchDebugLogger.Error("", err)
	}
	err = CreateFairLaunchMint()
	if err != nil {
		FairLaunchDebugLogger.Error("", err)
	}
	err = CreateSendFairLaunchAsset()
	if err != nil {
		FairLaunchDebugLogger.Error("", err)
	}
	fmt.Println("Create FairLaunch ScheduledTasks Procession finished!")
}

func (cs *CronService) PollPaymentCron() {
	pollPayment()
}

func (cs *CronService) PollInvoiceCron() {
	pollInvoice()
}

func (cs *CronService) PollPayInsideMission() {
	pollPayInsideMission()
}

func CreatePollPaymentCron() (err error) {
	return CreateScheduledTask(&models.ScheduledTask{
		Name:           "PollPaymentCron",
		CronExpression: "*/25 * * * * *",
		FunctionName:   "PollPaymentCron",
		Package:        "services",
	})
}

func CreatePollInvoiceCron() (err error) {
	return CreateScheduledTask(&models.ScheduledTask{
		Name:           "PollInvoiceCron",
		CronExpression: "*/25 * * * * *",
		FunctionName:   "PollInvoiceCron",
		Package:        "services",
	})
}

func CreatePollPayInvoiceMission() (err error) {
	return CreateScheduledTask(&models.ScheduledTask{
		Name:           "PollPayInsideMission",
		CronExpression: "*/25 * * * * *",
		FunctionName:   "PollPayInsideMission",
		Package:        "services",
	})
}

func CreatePAYTasks() {
	err := CreatePollPaymentCron()
	if err != nil {
		CUST.Error("", err)
	}
	err = CreatePollInvoiceCron()
	if err != nil {
		CUST.Error("", err)
	}
	err = CreatePollPayInvoiceMission()
	if err != nil {
		CUST.Error("", err)
	}
	fmt.Println("Create FairLaunch ScheduledTasks Procession finished!")
}
