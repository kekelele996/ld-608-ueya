package repositories

import (
	"groundTurn/src/models"

	"gorm.io/gorm"
)

type GroundTaskRepository struct{ DB *gorm.DB }

func NewGroundTaskRepository(db *gorm.DB) *GroundTaskRepository {
	return &GroundTaskRepository{DB: db}
}

func (r *GroundTaskRepository) ListByTurnaround(tx *gorm.DB, turnaroundID uint) ([]models.GroundTask, error) {
	var rows []models.GroundTask
	err := tx.Where("turnaround_id = ?", turnaroundID).
		Order("planned_start ASC").Find(&rows).Error
	return rows, err
}

func (r *GroundTaskRepository) ListAll() ([]models.GroundTask, error) {
	var rows []models.GroundTask
	err := r.DB.Order("planned_start ASC").Find(&rows).Error
	return rows, err
}

func (r *GroundTaskRepository) Get(tx *gorm.DB, id uint) (*models.GroundTask, error) {
	var row models.GroundTask
	if err := tx.First(&row, id).Error; err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *GroundTaskRepository) Create(tx *gorm.DB, row *models.GroundTask) error {
	return tx.Create(row).Error
}

func (r *GroundTaskRepository) Update(tx *gorm.DB, row *models.GroundTask) error {
	return tx.Save(row).Error
}

func (r *GroundTaskRepository) CountByTurnaround(tx *gorm.DB, turnaroundID uint) (int64, error) {
	var count int64
	err := tx.Model(&models.GroundTask{}).
		Where("turnaround_id = ?", turnaroundID).Count(&count).Error
	return count, err
}
