package repositories

import (
	"time"

	"groundTurn/src/models"

	"gorm.io/gorm"
)

type ResourceRepository struct{ DB *gorm.DB }

func NewResourceRepository(db *gorm.DB) *ResourceRepository { return &ResourceRepository{DB: db} }

func (r *ResourceRepository) CreateBatch(resources []models.GroundResource) error {
	if len(resources) == 0 {
		return nil
	}
	return r.DB.Create(&resources).Error
}

func (r *ResourceRepository) Save(t *models.GroundResource) error { return r.DB.Save(t).Error }

func (r *ResourceRepository) Get(id uint) (*models.GroundResource, error) {
	var t models.GroundResource
	if err := r.DB.First(&t, id).Error; err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *ResourceRepository) List(page, size int, status, resourceType string) ([]models.GroundResource, int64, error) {
	var rows []models.GroundResource
	var total int64
	q := r.DB.Model(&models.GroundResource{})
	if status != "" {
		q = q.Where("availability_status = ?", status)
	}
	if resourceType != "" {
		q = q.Where("resource_type = ?", resourceType)
	}
	q.Count(&total)
	err := q.Order("resource_code ASC").Offset((page - 1) * size).Limit(size).Find(&rows).Error
	return rows, total, err
}

func (r *ResourceRepository) ListByType(resourceType string) ([]models.GroundResource, error) {
	var rows []models.GroundResource
	err := r.DB.Where("resource_type = ?", resourceType).Order("resource_code ASC").Find(&rows).Error
	return rows, err
}

// MaintenanceOverlap returns resources whose maintenance window intersects.
func (r *ResourceRepository) MaintenanceOverlap(start, end time.Time) ([]models.GroundResource, error) {
	var rows []models.GroundResource
	err := r.DB.Where("maintenance_due_at IS NOT NULL AND maintenance_end IS NOT NULL").
		Where("maintenance_due_at < ? AND maintenance_end > ?", end, start).
		Find(&rows).Error
	return rows, err
}
