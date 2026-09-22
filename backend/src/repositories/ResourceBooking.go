package repositories

import (
	"time"

	"groundTurn/src/models"

	"gorm.io/gorm"
)

type BookingRepository struct{ DB *gorm.DB }

func NewBookingRepository(db *gorm.DB) *BookingRepository { return &BookingRepository{DB: db} }

func (r *BookingRepository) CreateBatch(bookings []models.ResourceBooking) error {
	if len(bookings) == 0 {
		return nil
	}
	return r.DB.Create(&bookings).Error
}

func (r *BookingRepository) Save(t *models.ResourceBooking) error { return r.DB.Save(t).Error }

func (r *BookingRepository) Get(id uint) (*models.ResourceBooking, error) {
	var t models.ResourceBooking
	if err := r.DB.First(&t, id).Error; err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *BookingRepository) ListByTurnaround(turnaroundID uint) ([]models.ResourceBooking, error) {
	var rows []models.ResourceBooking
	err := r.DB.Where("turnaround_id = ?", turnaroundID).Order("start_time ASC").Find(&rows).Error
	return rows, err
}

func (r *BookingRepository) ListByTask(taskID uint) ([]models.ResourceBooking, error) {
	var rows []models.ResourceBooking
	err := r.DB.Where("task_id = ?", taskID).Find(&rows).Error
	return rows, err
}

func (r *BookingRepository) ListPending() ([]models.ResourceBooking, error) {
	var rows []models.ResourceBooking
	err := r.DB.Where("booking_status = ?", "PENDING").Order("start_time ASC").Find(&rows).Error
	return rows, err
}

func (r *BookingRepository) List(page, size int, status string) ([]models.ResourceBooking, int64, error) {
	var rows []models.ResourceBooking
	var total int64
	q := r.DB.Model(&models.ResourceBooking{})
	if status != "" {
		q = q.Where("booking_status = ?", status)
	}
	q.Count(&total)
	err := q.Order("start_time ASC").Offset((page - 1) * size).Limit(size).Find(&rows).Error
	return rows, total, err
}

// Overlapping returns CONFIRMED/PENDING bookings of a resource whose window
// overlaps [start,end). excludeBookingID lets adjustment ignore itself.
func (r *BookingRepository) Overlapping(resourceID uint, start, end time.Time, excludeBookingID uint) ([]models.ResourceBooking, error) {
	var rows []models.ResourceBooking
	q := r.DB.Where("resource_id = ?", resourceID).
		Where("booking_status <> ?", "RELEASED").
		Where("start_time < ? AND end_time > ?", end, start)
	if excludeBookingID != 0 {
		q = q.Where("id <> ?", excludeBookingID)
	}
	err := q.Find(&rows).Error
	return rows, err
}

// FrozenBookingsBetween returns non-released bookings whose task has already
// been signed/completed: those windows cannot move during reschedule.
func (r *BookingRepository) FrozenBookingsBetween(resourceID uint, start, end time.Time) ([]models.ResourceBooking, error) {
	var rows []models.ResourceBooking
	err := r.DB.Where("resource_booking.resource_id = ?", resourceID).
		Where("resource_booking.booking_status = ?", "CONFIRMED").
		Where("resource_booking.start_time < ? AND resource_booking.end_time > ?", end, start).
		Joins("JOIN ground_task ON ground_task.id = resource_booking.task_id").
		Where("ground_task.status IN ?", []string{"ACCEPTED", "COMPLETED"}).
		Find(&rows).Error
	return rows, err
}

// ListInWindow supports the resource calendar view.
func (r *BookingRepository) ListInWindow(from, to time.Time) ([]models.ResourceBooking, error) {
	var rows []models.ResourceBooking
	err := r.DB.Where("start_time < ? AND end_time > ?", to, from).
		Order("start_time ASC").Find(&rows).Error
	return rows, err
}
