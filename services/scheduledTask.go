package services

import (
	"trade/middleware"
	"trade/models"
)

func CreateScheduledTask(scheduledTask *models.ScheduledTask) (err error) {
	s := ScheduledTaskStore{DB: middleware.DB}
	return s.CreateScheduledTask(scheduledTask)
}

func CreateScheduledTasks(scheduledTasks *[]models.ScheduledTask) (err error) {
	s := ScheduledTaskStore{DB: middleware.DB}
	return s.CreateScheduledTasks(scheduledTasks)
}

func ReadScheduledTaskByName(name string) (*models.ScheduledTask, error) {
	var scheduledTask models.ScheduledTask
	err := middleware.DB.Where("name = ?", name).First(&scheduledTask).Error
	return &scheduledTask, err
}

func IsScheduledTaskChanged(scheduledTask *models.ScheduledTask, old *models.ScheduledTask) bool {
	if scheduledTask == nil || old == nil {
		return true
	}
	if old.CronExpression != scheduledTask.CronExpression || old.FunctionName != scheduledTask.FunctionName || old.Package != scheduledTask.Package {
		return true
	}
	return false
}

func CreateScheduledTaskIfNotExistOrUpdate(scheduledTask *models.ScheduledTask) (err error) {
	scheduledTaskByName, err := ReadScheduledTaskByName(scheduledTask.Name)
	if err != nil {
		err = CreateScheduledTask(scheduledTask)
		if err != nil {
			return err
		}
		return nil
	}
	if !IsScheduledTaskChanged(scheduledTask, scheduledTaskByName) {
		return nil
	}
	scheduledTaskByName.CronExpression = scheduledTask.CronExpression
	scheduledTaskByName.FunctionName = scheduledTask.FunctionName
	scheduledTaskByName.Package = scheduledTask.Package
	s := ScheduledTaskStore{DB: middleware.DB}
	return s.UpdateScheduledTask(scheduledTaskByName)
}

func CreateOrUpdateScheduledTasks(scheduledTasks *[]models.ScheduledTask) (err error) {
	for _, scheduledTask := range *scheduledTasks {
		err = CreateScheduledTaskIfNotExistOrUpdate(&scheduledTask)
		if err != nil {
			return err
		}
	}
	return nil
}
