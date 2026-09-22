package repositories

import (
	"groundTurn/src/models"

	"gorm.io/gorm"
)

type GroundResourceRepository struct{ DB *gorm.DB }

func NewGroundResourceRepository(db *gorm.DB) *GroundResourceRepository {
	return &GroundResourceRepository{DB: db}
}

func (r *GroundResourceRepository) List() ([]models.GroundResource, error) {
	var rows []models.GroundResource
	err := r.DB.Order("resource_type ASC, resource_code ASC").Find(&rows).Error
	return rows, err
}

func (r *GroundResourceRepository) ListByType(tx *gorm.DB, resourceType string) ([]models.GroundResource, error) {
	var rows []models.GroundResource
	err := tx.Where("resource_type = ?", resourceType).
		Order("resource_code ASC").Find(&rows).Error
	return rows, err
}

func (r *GroundResourceRepository) Get(tx *gorm.DB, id uint) (*models.GroundResource, error) {
	var row models.GroundResource
	if err := tx.First(&row, id).Error; err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *GroundResourceRepository) Create(tx *gorm.DB, row *models.GroundResource) error {
	return tx.Create(row).Error
}

func (r *GroundResourceRepository) Update(tx *gorm.DB, row *models.GroundResource) error {
	return tx.Save(row).Error
}
