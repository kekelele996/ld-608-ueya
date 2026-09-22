package services

import (
	"errors"
	"fmt"

	"groundTurn/src/constants"
	"groundTurn/src/models"
	"groundTurn/src/repositories"
	"groundTurn/src/types"
	"groundTurn/src/utils"

	"gorm.io/gorm"
)

type ResourceService struct {
	svc      *Services
	repo     *repositories.GroundResourceRepository
	bookRepo *repositories.ResourceBookingRepository
}

func NewResourceService(svc *Services) *ResourceService {
	return &ResourceService{
		svc:      svc,
		repo:     repositories.NewGroundResourceRepository(svc.DB),
		bookRepo: repositories.NewResourceBookingRepository(svc.DB),
	}
}

func (s *ResourceService) List() ([]models.GroundResource, error) { return s.repo.List() }

func (s *ResourceService) Get(id uint) (*models.GroundResource, error) {
	row, err := s.repo.Get(s.svc.DB, id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, appError(404, constants.CodeNotFound,
			fmt.Sprintf(constants.MsgNotFound, "保障资源"))
	}
	return row, err
}

// SetMaintenance flips MAINTENANCE / OFFLINE / AVAILABLE and records the due time.
func (s *ResourceService) SetMaintenance(actor Actor, id uint, req types.MaintenanceRequest) (*models.GroundResource, error) {
	row, err := s.Get(id)
	if err != nil {
		return nil, err
	}
	target := constants.ResourceStatus(req.Status)
	if !target.Valid() {
		return nil, appError(400, constants.CodeValidationFailed,
			fmt.Sprintf(constants.MsgValidationFailed, "未知资源状态 "+req.Status))
	}
	err = s.svc.DB.Transaction(func(tx *gorm.DB) error {
		row.AvailabilityStatus = string(target)
		row.MaintenanceDueAt = req.MaintenanceDueAt
		if err := s.repo.Update(tx, row); err != nil {
			return err
		}
		s.svc.audit(tx, actor.Username, actor.Role, constants.ActionResourceMaintenance,
			"GroundResource", fmt.Sprintf("%d", row.ID),
			fmt.Sprintf(constants.LogTemplates[constants.ActionResourceMaintenance],
				row.ResourceCode, target.Text(), utils.FormatTimePtr(req.MaintenanceDueAt)))
		return nil
	})
	if err != nil {
		return nil, appError(500, constants.CodeInternal, constants.MsgInternal)
	}
	return row, nil
}

// PendingBookings lists appointments bumped into 待处理 for manual triage.
func (s *ResourceService) PendingBookings() ([]models.ResourceBooking, error) {
	return s.bookRepo.ListPending()
}

func (s *ResourceService) AllBookings() ([]models.ResourceBooking, error) {
	return s.bookRepo.List()
}
