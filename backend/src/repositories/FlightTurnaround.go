package repositories

import (
	"groundTurn/src/models"

	"gorm.io/gorm"
)

type TurnaroundRepository struct{ DB *gorm.DB }

func NewTurnaroundRepository(db *gorm.DB) *TurnaroundRepository {
	return &TurnaroundRepository{DB: db}
}

func (r *TurnaroundRepository) Create(t *models.FlightTurnaround) error {
	return r.DB.Create(t).Error
}

func (r *TurnaroundRepository) Save(t *models.FlightTurnaround) error {
	return r.DB.Save(t).Error
}

func (r *TurnaroundRepository) Get(id uint) (*models.FlightTurnaround, error) {
	var t models.FlightTurnaround
	if err := r.DB.First(&t, id).Error; err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *TurnaroundRepository) List(page, size int) ([]models.FlightTurnaround, int64, error) {
	var rows []models.FlightTurnaround
	var total int64
	r.DB.Model(&models.FlightTurnaround{}).Count(&total)
	err := r.DB.Order("arrival_time DESC").Offset((page - 1) * size).Limit(size).Find(&rows).Error
	return rows, total, err
}
