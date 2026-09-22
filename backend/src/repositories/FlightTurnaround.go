package repositories

import (
	"groundTurn/src/models"

	"gorm.io/gorm"
)

type FlightTurnaroundRepository struct{ DB *gorm.DB }

func NewFlightTurnaroundRepository(db *gorm.DB) *FlightTurnaroundRepository {
	return &FlightTurnaroundRepository{DB: db}
}

func (r *FlightTurnaroundRepository) List() ([]models.FlightTurnaround, error) {
	var rows []models.FlightTurnaround
	err := r.DB.Order("arrival_time DESC").Find(&rows).Error
	return rows, err
}

func (r *FlightTurnaroundRepository) Get(id uint) (*models.FlightTurnaround, error) {
	var row models.FlightTurnaround
	if err := r.DB.First(&row, id).Error; err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *FlightTurnaroundRepository) Create(tx *gorm.DB, row *models.FlightTurnaround) error {
	return tx.Create(row).Error
}

func (r *FlightTurnaroundRepository) Update(tx *gorm.DB, row *models.FlightTurnaround) error {
	return tx.Save(row).Error
}
