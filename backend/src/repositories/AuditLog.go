package repositories

import (
	"groundTurn/src/models"

	"gorm.io/gorm"
)

type AuditLogRepository struct{ DB *gorm.DB }

func NewAuditLogRepository(db *gorm.DB) *AuditLogRepository {
	return &AuditLogRepository{DB: db}
}

func (r *AuditLogRepository) Create(tx *gorm.DB, row *models.AuditLog) error {
	return tx.Create(row).Error
}

func (r *AuditLogRepository) List(limit int) ([]models.AuditLog, error) {
	var rows []models.AuditLog
	if limit <= 0 {
		limit = 100
	}
	err := r.DB.Order("created_at DESC").Limit(limit).Find(&rows).Error
	return rows, err
}
