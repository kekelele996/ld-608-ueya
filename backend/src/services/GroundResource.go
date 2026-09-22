package services

import (
	"fmt"
	"time"

	"groundTurn/src/constants"
	"groundTurn/src/models"
	"groundTurn/src/repositories"
	"groundTurn/src/utils"

	"gorm.io/gorm"
)

// ResourceService maintains the resource ledger and availability windows.
type ResourceService struct {
	DB       *gorm.DB
	Resource *repositories.ResourceRepository
	Booking  *repositories.BookingRepository
	Audit    *repositories.AuditRepository
}

func NewResourceService(db *gorm.DB) *ResourceService {
	return &ResourceService{
		DB:       db,
		Resource: repositories.NewResourceRepository(db),
		Booking:  repositories.NewBookingRepository(db),
		Audit:    repositories.NewAuditRepository(db),
	}
}

// Register adds a new resource to the ledger.
func (s *ResourceService) Register(code, resourceType, location, team string, actor, role string) (*models.GroundResource, error) {
	r := &models.GroundResource{
		ResourceCode:       code,
		ResourceType:       resourceType,
		Location:           location,
		AvailabilityStatus: constants.ResourceAvailable,
		OwnerTeam:          team,
	}
	if err := s.Resource.Save(r); err != nil {
		return nil, utils.NewAPIError(409, constants.CodeAlreadyExists,
			fmt.Sprintf("资源编码 %s 已存在", code), nil)
	}
	s.audit(actor, role, fmt.Sprintf(constants.LogResourceCreated, code, resourceType, team),
		"GroundResource", fmt.Sprintf("%d", r.ID))
	return r, nil
}

// UpdateStatus changes availability and optionally registers a maintenance
// window. Existing bookings are not rewritten here; conflict detection runs
// on the next plan/adjust action.
func (s *ResourceService) UpdateStatus(resourceID uint, status string, maintenanceStart, maintenanceEnd *time.Time, actor, role string) (*models.GroundResource, error) {
	r, err := s.Resource.Get(resourceID)
	if err != nil {
		return nil, utils.NewAPIError(404, constants.CodeNotFound,
			fmt.Sprintf(constants.MsgNotFound, "保障资源", resourceID), nil)
	}
	if status != "" && !constants.IsValidResourceStatus(status) {
		return nil, utils.NewAPIError(400, constants.CodeValidationFailed,
			fmt.Sprintf(constants.MsgValidationFailed, "未知资源状态 "+status), nil)
	}
	previous := r.AvailabilityStatus
	if status != "" {
		r.AvailabilityStatus = status
	}
	if status == constants.ResourceMaintenance {
		r.MaintenanceDueAt = maintenanceStart
		r.MaintenanceEnd = maintenanceEnd
	} else if maintenanceStart != nil && maintenanceEnd != nil {
		r.MaintenanceDueAt = maintenanceStart
		r.MaintenanceEnd = maintenanceEnd
	}
	if err := s.Resource.Save(r); err != nil {
		return nil, err
	}
	if previous != r.AvailabilityStatus {
		s.audit(actor, role, fmt.Sprintf(constants.LogResourceStatus,
			r.ResourceCode, constants.ResourceStatusText[previous],
			constants.ResourceStatusText[r.AvailabilityStatus]),
			"GroundResource", fmt.Sprintf("%d", r.ID))
	}
	if r.MaintenanceDueAt != nil && r.MaintenanceEnd != nil {
		s.audit(actor, role, fmt.Sprintf(constants.LogResourceMaintenance,
			r.ResourceCode, r.MaintenanceDueAt.Format("01-02 15:04"),
			r.MaintenanceEnd.Format("01-02 15:04")),
			"GroundResource", fmt.Sprintf("%d", r.ID))
	}
	return r, nil
}

func (s *ResourceService) audit(actor, role, message, targetType, targetID string) {
	_ = s.Audit.Create(&models.AuditLog{
		Actor: actor, Role: role, Action: message, TargetType: targetType, TargetID: targetID,
	})
}
