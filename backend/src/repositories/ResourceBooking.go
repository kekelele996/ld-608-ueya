package repositories

import (
	"groundTurn/src/models"

	"gorm.io/gorm"
)

type ResourceBookingRepository struct{ DB *gorm.DB }

func NewResourceBookingRepository(db *gorm.DB) *ResourceBookingRepository {
	return &ResourceBookingRepository{DB: db}
}

func (r *ResourceBookingRepository) List() ([]models.ResourceBooking, error) {
	return r.ListAll(nil)
}

// ListAll reads every booking; pass the active tx inside a transaction so
// SQLite does not open a competing connection mid-write.
func (r *ResourceBookingRepository) ListAll(tx *gorm.DB) ([]models.ResourceBooking, error) {
	var rows []models.ResourceBooking
	q := r.DB
	if tx != nil {
		q = tx
	}
	err := q.Order("start_time DESC").Find(&rows).Error
	return rows, err
}

func (r *ResourceBookingRepository) ListByTurnaround(tx *gorm.DB, turnaroundID uint) ([]models.ResourceBooking, error) {
	var rows []models.ResourceBooking
	err := tx.Where("turnaround_id = ?", turnaroundID).
		Order("start_time ASC").Find(&rows).Error
	return rows, err
}

func (r *ResourceBookingRepository) ListPending() ([]models.ResourceBooking, error) {
	var rows []models.ResourceBooking
	err := r.DB.Where("booking_status = ?", "PENDING").
		Order("start_time ASC").Find(&rows).Error
	return rows, err
}

// ActiveForResource returns confirmed bookings of a resource, optionally
// excluding one booking id (used during adjustment).
func (r *ResourceBookingRepository) ActiveForResource(tx *gorm.DB, resourceID uint, excludeID uint) ([]models.ResourceBooking, error) {
	var rows []models.ResourceBooking
	q := tx.Where("resource_id = ? AND booking_status = ?", resourceID, "CONFIRMED")
	if excludeID > 0 {
		q = q.Where("id <> ?", excludeID)
	}
	err := q.Order("start_time ASC").Find(&rows).Error
	return rows, err
}

func (r *ResourceBookingRepository) Get(tx *gorm.DB, id uint) (*models.ResourceBooking, error) {
	var row models.ResourceBooking
	if err := tx.First(&row, id).Error; err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *ResourceBookingRepository) Create(tx *gorm.DB, row *models.ResourceBooking) error {
	return tx.Create(row).Error
}

func (r *ResourceBookingRepository) Update(tx *gorm.DB, row *models.ResourceBooking) error {
	return tx.Save(row).Error
}
