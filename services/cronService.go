package services

import (
	"fmt"
	"trade/middleware"
	"trade/models"
)

type CronService struct{}

func (cs *CronService) FairLaunchIssuance() {
	//FairLaunchDebugLogger.Info("start cron job: FairLaunchIssuance")
	FairLaunchIssuance()
}

func (cs *CronService) FairLaunchMint() {
	//FairLaunchDebugLogger.Info("start cron job: FairLaunchMint")
	FairLaunchMint()
}

func (cs *CronService) SendFairLaunchAsset() {
	//FairLaunchDebugLogger.Info("start cron job: FairLaunchMint")
	SendFairLaunchAsset()
}

func (cs *CronService) UpdateFeeRateWeek() {
	//FairLaunchDebugLogger.Info("start cron job: UpdateFeeRateWeek")
	_ = CheckIfUpdateFeeRateInfoByBlockOfWeek()
}

func CreateScheduledTask(scheduledTask *models.ScheduledTask) (err error) {
	s := ScheduledTaskStore{DB: middleware.DB}
	return s.CreateScheduledTask(scheduledTask)
}

func CreateFairLaunchIssuance() (err error) {
	return CreateScheduledTask(&models.ScheduledTask{
		Name:           "FairLaunchIssuance",
		CronExpression: "* */1 * * * *",
		FunctionName:   "FairLaunchIssuance",
		Package:        "services",
	})
}

func CreateFairLaunchMint() (err error) {
	return CreateScheduledTask(&models.ScheduledTask{
		Name:           "FairLaunchMint",
		CronExpression: "* */1 * * * *",
		FunctionName:   "FairLaunchMint",
		Package:        "services",
	})
}

func CreateSendFairLaunchAsset() (err error) {
	return CreateScheduledTask(&models.ScheduledTask{
		Name:           "SendFairLaunchAsset",
		CronExpression: "* */1 * * * *",
		FunctionName:   "SendFairLaunchAsset",
		Package:        "services",
	})
}
func CreateUpdateFeeRateWeek() (err error) {
	return CreateScheduledTask(&models.ScheduledTask{
		Name:           "UpdateFeeRateWeek",
		CronExpression: "* */20 * * * *",
		FunctionName:   "UpdateFeeRateWeek",
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
	//if err != nil {
	//	FairLaunchDebugLogger.Error("", err)
	//}
	//err = CreateUpdateFeeRateWeek()
	if err != nil {
		FairLaunchDebugLogger.Error("", err)
	}
	fmt.Println("Create FairLaunch ScheduledTasks Procession finished!")
}

func (sm *CronService) PollPaymentCron() {
	pollPayment()
}

func (sm *CronService) PollInvoiceCron() {
	pollInvoice()
}

func (sm *CronService) PollPayInsideMission() {
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
