package repositories

import (
	"groundTurn/src/models"

	"gorm.io/gorm"
)

type AuditRepository struct{ DB *gorm.DB }

func NewAuditRepository(db *gorm.DB) *AuditRepository { return &AuditRepository{DB: db} }

func (r *AuditRepository) Create(entry *models.AuditLog) error { return r.DB.Create(entry).Error }

func (r *AuditRepository) List(page, size int) ([]models.AuditLog, int64, error) {
	var rows []models.AuditLog
	var total int64
	r.DB.Model(&models.AuditLog{}).Count(&total)
	err := r.DB.Order("created_at DESC").Offset((page - 1) * size).Limit(size).Find(&rows).Error
	return rows, total, err
}

type UserRepository struct{ DB *gorm.DB }

func NewUserRepository(db *gorm.DB) *UserRepository { return &UserRepository{DB: db} }

func (r *UserRepository) GetByUsername(username string) (*models.User, error) {
	var u models.User
	if err := r.DB.Where("username = ?", username).First(&u).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) List() ([]models.User, error) {
	var rows []models.User
	err := r.DB.Order("id ASC").Find(&rows).Error
	return rows, err
}
