package repositories

import (
	"groundTurn/src/models"

	"gorm.io/gorm"
)

type DelayRepository struct{ DB *gorm.DB }

func NewDelayRepository(db *gorm.DB) *DelayRepository { return &DelayRepository{DB: db} }

func (r *DelayRepository) Create(t *models.DelayEvent) error { return r.DB.Create(t).Error }

func (r *DelayRepository) Save(t *models.DelayEvent) error { return r.DB.Save(t).Error }

func (r *DelayRepository) Get(id uint) (*models.DelayEvent, error) {
	var t models.DelayEvent
	if err := r.DB.First(&t, id).Error; err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *DelayRepository) ListByTurnaround(turnaroundID uint) ([]models.DelayEvent, error) {
	var rows []models.DelayEvent
	err := r.DB.Where("turnaround_id = ?", turnaroundID).Order("created_at DESC").Find(&rows).Error
	return rows, err
}

func (r *DelayRepository) List(page, size int) ([]models.DelayEvent, int64, error) {
	var rows []models.DelayEvent
	var total int64
	r.DB.Model(&models.DelayEvent{}).Count(&total)
	err := r.DB.Order("created_at DESC").Offset((page - 1) * size).Limit(size).Find(&rows).Error
	return rows, total, err
}
