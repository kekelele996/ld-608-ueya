package services

import (
	"time"

	"groundTurn/src/models"
	"groundTurn/src/repositories"

	"gorm.io/gorm"
)

// QueryService is the read side consumed by dashboard/list pages.
type QueryService struct {
	Turn     *repositories.TurnaroundRepository
	Task     *repositories.TaskRepository
	Resource *repositories.ResourceRepository
	Booking  *repositories.BookingRepository
	Delay    *repositories.DelayRepository
	Audit    *repositories.AuditRepository
	User     *repositories.UserRepository
	DB       *gorm.DB
}

func NewQueryService(db *gorm.DB) *QueryService {
	return &QueryService{
		Turn:     repositories.NewTurnaroundRepository(db),
		Task:     repositories.NewTaskRepository(db),
		Resource: repositories.NewResourceRepository(db),
		Booking:  repositories.NewBookingRepository(db),
		Delay:    repositories.NewDelayRepository(db),
		Audit:    repositories.NewAuditRepository(db),
		User:     repositories.NewUserRepository(db),
		DB:       db,
	}
}

// DashboardStats feeds StatCard row on the operations dashboard.
type DashboardStats struct {
	ActiveTurnarounds int `json:"active_turnarounds"`
	OpenTasks         int `json:"open_tasks"`
	PendingBookings   int `json:"pending_bookings"`
	TotalDelayMinutes int `json:"total_delay_minutes"`
}

func (q *QueryService) DashboardStats() (DashboardStats, error) {
	var raw struct {
		ActiveTurnarounds int64
		OpenTasks         int64
		PendingBookings   int64
		TotalDelayMinutes int64
	}
	q.DB.Model(&models.FlightTurnaround{}).
		Where("turnaround_status IN ?", []string{"ON_STAND", "IN_SERVICE", "DELAYED"}).
		Count(&raw.ActiveTurnarounds)
	q.DB.Model(&models.GroundTask{}).
		Where("status IN ?", []string{"PLANNED", "ACCEPTED", "BLOCKED"}).
		Count(&raw.OpenTasks)
	q.DB.Model(&models.ResourceBooking{}).
		Where("booking_status = ?", "PENDING").
		Count(&raw.PendingBookings)
	q.DB.Model(&models.FlightTurnaround{}).
		Select("COALESCE(SUM(total_delay_minutes),0)").Scan(&raw.TotalDelayMinutes)
	return DashboardStats{
		ActiveTurnarounds: int(raw.ActiveTurnarounds),
		OpenTasks:         int(raw.OpenTasks),
		PendingBookings:   int(raw.PendingBookings),
		TotalDelayMinutes: int(raw.TotalDelayMinutes),
	}, nil
}

// CalendarWindow returns bookings for the resource calendar.
func (q *QueryService) CalendarWindow(from, to time.Time) ([]models.ResourceBooking, error) {
	return q.Booking.ListInWindow(from, to)
}

func (q *QueryService) PendingBookings() ([]models.ResourceBooking, error) {
	return q.Booking.ListPending()
}

func (q *QueryService) AllResources() ([]models.GroundResource, error) {
	rows, _, err := q.Resource.List(1, 200, "", "")
	return rows, err
}
