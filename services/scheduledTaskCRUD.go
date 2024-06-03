package services

import (
	"gorm.io/gorm"
	"trade/models"
)

type ScheduledTaskStore struct {
	DB *gorm.DB
}

// ScheduledTask

func (s *ScheduledTaskStore) CreateScheduledTask(scheduledTask *models.ScheduledTask) error {
	return s.DB.Create(scheduledTask).Error
}

func (s *ScheduledTaskStore) ReadScheduledTask(id uint) (*models.ScheduledTask, error) {
	var scheduledTask models.ScheduledTask
	err := s.DB.First(&scheduledTask, id).Error
	return &scheduledTask, err
}

func (s *ScheduledTaskStore) UpdateScheduledTask(scheduledTask *models.ScheduledTask) error {
	return s.DB.Save(scheduledTask).Error
}

func (s *ScheduledTaskStore) DeleteScheduledTask(id uint) error {
	var scheduledTask models.ScheduledTask
	return s.DB.Delete(&scheduledTask, id).Error
}
