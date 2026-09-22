package repositories

import (
	"groundTurn/src/models"

	"gorm.io/gorm"
)

type DelayEventRepository struct{ DB *gorm.DB }

func NewDelayEventRepository(db *gorm.DB) *DelayEventRepository {
	return &DelayEventRepository{DB: db}
}

func (r *DelayEventRepository) List() ([]models.DelayEvent, error) {
	var rows []models.DelayEvent
	err := r.DB.Order("created_at DESC").Find(&rows).Error
	return rows, err
}

func (r *DelayEventRepository) ListByTurnaround(tx *gorm.DB, turnaroundID uint) ([]models.DelayEvent, error) {
	var rows []models.DelayEvent
	err := tx.Where("turnaround_id = ?", turnaroundID).
		Order("created_at ASC").Find(&rows).Error
	return rows, err
}

func (r *DelayEventRepository) Create(tx *gorm.DB, row *models.DelayEvent) error {
	return tx.Create(row).Error
}

func (r *DelayEventRepository) Update(tx *gorm.DB, row *models.DelayEvent) error {
	return tx.Save(row).Error
}
