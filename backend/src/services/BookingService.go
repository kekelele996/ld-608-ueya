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

type BookingService struct {
	svc     *Services
	repo    *repositories.ResourceBookingRepository
	resRepo *repositories.GroundResourceRepository
}

func NewBookingService(svc *Services) *BookingService {
	return &BookingService{
		svc:     svc,
		repo:    repositories.NewResourceBookingRepository(svc.DB),
		resRepo: repositories.NewGroundResourceRepository(svc.DB),
	}
}

func (s *BookingService) List() ([]models.ResourceBooking, error) { return s.repo.List() }

func (s *BookingService) Get(id uint) (*models.ResourceBooking, error) {
	row, err := s.repo.Get(s.svc.DB, id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, appError(404, constants.CodeNotFound,
			fmt.Sprintf(constants.MsgNotFound, "资源预约"))
	}
	return row, err
}

// Adjust moves an appointment to a new resource/time window. The same
// validation rules as plan generation apply; a failed adjustment leaves the
// booking in PENDING with the reason, a successful one becomes CONFIRMED.
func (s *BookingService) Adjust(actor Actor, id uint, req types.AdjustBookingRequest) (*models.ResourceBooking, error) {
	booking, err := s.Get(id)
	if err != nil {
		return nil, err
	}
	if !req.EndTime.After(req.StartTime) {
		return nil, appError(400, constants.CodeValidationFailed,
			fmt.Sprintf(constants.MsgValidationFailed, "结束时间必须晚于开始时间"))
	}
	resource, err := s.resRepo.Get(s.svc.DB, req.ResourceID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, appError(404, constants.CodeNotFound,
			fmt.Sprintf(constants.MsgNotFound, "保障资源"))
	}
	if err != nil {
		return nil, err
	}

	reason := ""
	status := string(constants.BookingConfirmed)

	if code, message := resourceCovers(resource, req.StartTime, req.EndTime); code != "" {
		reason = message
		status = string(constants.BookingPending)
	} else {
		others, err := s.repo.ActiveForResource(s.svc.DB, resource.ID, booking.ID)
		if err != nil {
			return nil, appError(500, constants.CodeInternal, constants.MsgInternal)
		}
		if hit := bookingOverlap(others, resource.ID, req.StartTime, req.EndTime, booking.ID, 0); hit != nil {
			reason = fmt.Sprintf(constants.ConflictMessages[constants.ItemResourceOverlap],
				resource.ResourceCode, utils.FormatTime(req.StartTime),
				utils.FormatTime(req.EndTime), hit.TaskIDValue())
			status = string(constants.BookingPending)
		}
	}

	err = s.svc.DB.Transaction(func(tx *gorm.DB) error {
		booking.ResourceID = resource.ID
		booking.StartTime = req.StartTime
		booking.EndTime = req.EndTime
		booking.BookingStatus = status
		booking.ConflictReason = reason
		if err := s.repo.Update(tx, booking); err != nil {
			return err
		}
		action := constants.ActionBookingAdjust
		s.svc.audit(tx, actor.Username, actor.Role, action,
			"ResourceBooking", fmt.Sprintf("%d", booking.ID),
			fmt.Sprintf(constants.LogTemplates[constants.ActionBookingAdjust],
				booking.ID, resource.ResourceCode,
				utils.FormatTime(req.StartTime), utils.FormatTime(req.EndTime),
				constants.BookingStatus(status).Text()))
		return nil
	})
	if err != nil {
		return nil, appError(500, constants.CodeInternal, constants.MsgInternal)
	}
	utils.LogOperation(constants.ActionBookingAdjust, booking.ID, resource.ResourceCode,
		utils.FormatTime(req.StartTime), utils.FormatTime(req.EndTime),
		constants.BookingStatus(status).Text())
	return booking, nil
}

// Release frees a pending/unneeded appointment.
func (s *BookingService) Release(actor Actor, id uint) (*models.ResourceBooking, error) {
	booking, err := s.Get(id)
	if err != nil {
		return nil, err
	}
	err = s.svc.DB.Transaction(func(tx *gorm.DB) error {
		booking.BookingStatus = string(constants.BookingReleased)
		if err := s.repo.Update(tx, booking); err != nil {
			return err
		}
		s.svc.audit(tx, actor.Username, actor.Role, constants.ActionBookingRelease,
			"ResourceBooking", fmt.Sprintf("%d", booking.ID),
			fmt.Sprintf(constants.LogTemplates[constants.ActionBookingRelease],
				booking.ID, booking.ResourceID))
		return nil
	})
	if err != nil {
		return nil, appError(500, constants.CodeInternal, constants.MsgInternal)
	}
	return booking, nil
}
