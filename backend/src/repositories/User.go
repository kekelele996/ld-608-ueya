package repositories

import (
	"groundTurn/src/models"

	"gorm.io/gorm"
)

type UserRepository struct{ DB *gorm.DB }

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (r *UserRepository) FindByUsername(username string) (*models.User, error) {
	var row models.User
	if err := r.DB.Where("username = ?", username).First(&row).Error; err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *UserRepository) Create(tx *gorm.DB, row *models.User) error {
	return tx.Create(row).Error
}

func (r *UserRepository) Count(tx *gorm.DB) (int64, error) {
	var count int64
	err := tx.Model(&models.User{}).Count(&count).Error
	return count, err
}
