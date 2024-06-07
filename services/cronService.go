package services

import (
	"fmt"
	"strconv"
	"strings"
	"time"
	"trade/middleware"
	"trade/models"
	"trade/utils"
)

type CronService struct{}

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

func CreateScheduledTask(scheduledTask *models.ScheduledTask) (err error) {
	s := ScheduledTaskStore{DB: middleware.DB}
	return s.CreateScheduledTask(scheduledTask)
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
		CronExpression: "0 */1 * * * *",
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

func CreateFairLaunchScheduledTasks() {
	err := CreateFairLaunchIssuance()
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
