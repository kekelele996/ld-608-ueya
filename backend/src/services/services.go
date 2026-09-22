package services

import (
	"fmt"

	"groundTurn/src/config"
	"groundTurn/src/constructors"

	"gorm.io/gorm"
)

// Services bundles repositories so every entity service shares transactions
// and audit writes — the deliberate cross-file coupling of this codebase.
type Services struct {
	DB           *gorm.DB
	Cfg          *config.Config
	Turnarounds  *TurnaroundService
	Tasks        *TaskService
	Resources    *ResourceService
	Bookings     *BookingService
	Delays       *DelayService
	Auth         *AuthService
	Reports      *ReportService
}

func NewServices(db *gorm.DB, cfg *config.Config) *Services {
	s := &Services{DB: db, Cfg: cfg}
	s.Turnarounds = NewTurnaroundService(s)
	s.Tasks = NewTaskService(s)
	s.Resources = NewResourceService(s)
	s.Bookings = NewBookingService(s)
	s.Delays = NewDelayService(s)
	s.Auth = NewAuthService(s)
	s.Reports = NewReportService(s)
	return s
}

func (s *Services) audit(tx *gorm.DB, actor, role, action, targetType, targetID, detail string) {
	row := constructors.NewAuditLog(actor, role, action, targetType, targetID, detail)
	if err := tx.Create(row).Error; err != nil {
		fmt.Printf("audit write failed: %v\n", err)
	}
}

// Actor identifies the caller propagated from authMiddleware.
type Actor struct {
	Username string
	Role     string
	TeamID   string
	Name     string
}

func (a Actor) Label() string {
	if a.Name != "" {
		return a.Name
	}
	return a.Username
}

func (a Actor) IsZero() bool { return a.Username == "" }

func systemActor() Actor {
	return Actor{Username: "system", Role: "SYSTEM", Name: "系统"}
}

func (s *Services) actorOf(a Actor) Actor {
	if a.IsZero() {
		return systemActor()
	}
	return a
}
