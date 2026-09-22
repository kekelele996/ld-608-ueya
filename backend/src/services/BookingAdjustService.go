package services

import (
	"fmt"
	"time"

	"groundTurn/src/constants"
	"groundTurn/src/models"
	"groundTurn/src/repositories"
	"groundTurn/src/types"
	"groundTurn/src/utils"

	"gorm.io/gorm"
)

// BookingService resolves PENDING bookings by adjusting their time window.
type BookingService struct {
	DB       *gorm.DB
	Booking  *repositories.BookingRepository
	Resource *repositories.ResourceRepository
	Task     *repositories.TaskRepository
	Turn     *repositories.TurnaroundRepository
	Audit    *repositories.AuditRepository
}

func NewBookingService(db *gorm.DB) *BookingService {
	return &BookingService{
		DB:       db,
		Booking:  repositories.NewBookingRepository(db),
		Resource: repositories.NewResourceRepository(db),
		Task:     repositories.NewTaskRepository(db),
		Turn:     repositories.NewTurnaroundRepository(db),
		Audit:    repositories.NewAuditRepository(db),
	}
}

// Adjust moves a non-released booking to a new window after re-running
// overlap / maintenance / deadline checks. If the new window is still
// unusable, the booking stays PENDING and the per-item reasons come back.
func (s *BookingService) Adjust(bookingID uint, start, end time.Time, actor, role string) (*models.ResourceBooking, []types.ConflictItem, error) {
	if !end.After(start) {
		return nil, nil, utils.NewAPIError(400, constants.CodeValidationFailed,
			fmt.Sprintf(constants.MsgValidationFailed, "结束时间必须晚于开始时间"), nil)
	}
	booking, err := s.Booking.Get(bookingID)
	if err != nil {
		return nil, nil, utils.NewAPIError(404, constants.CodeNotFound,
			fmt.Sprintf(constants.MsgNotFound, "资源预约", bookingID), nil)
	}
	if booking.BookingStatus == constants.BookingReleased {
		return nil, nil, utils.NewAPIError(409, constants.CodeInvalidTransition,
			fmt.Sprintf("预约 %d 已释放，不能再调整", bookingID), nil)
	}
	resource, err := s.Resource.Get(booking.ResourceID)
	if err != nil {
		return nil, nil, utils.NewAPIError(404, constants.CodeNotFound,
			fmt.Sprintf(constants.MsgNotFound, "保障资源", booking.ResourceID), nil)
	}
	turn, err := s.Turn.Get(booking.TurnaroundID)
	if err != nil {
		return nil, nil, err
	}

	conflicts := []types.ConflictItem{}
	switch resource.AvailabilityStatus {
	case constants.ResourceOffline:
		conflicts = append(conflicts, types.ConflictItem{
			Code: constants.ConflictOffline,
			Message: fmt.Sprintf(constants.MsgResourceOffline, resource.ResourceCode),
			ResourceID: resource.ID, ResourceCode: resource.ResourceCode,
			BookingID: booking.ID, StartTime: utils.FormatTime(start), EndTime: utils.FormatTime(end),
		})
	case constants.ResourceMaintenance:
		conflicts = append(conflicts, types.ConflictItem{
			Code: constants.ConflictMaintenance,
			Message: fmt.Sprintf("资源 %s 标记为维护中", resource.ResourceCode),
			ResourceID: resource.ID, ResourceCode: resource.ResourceCode,
			BookingID: booking.ID, StartTime: utils.FormatTime(start), EndTime: utils.FormatTime(end),
		})
	}
	if resource.MaintenanceDueAt != nil && resource.MaintenanceEnd != nil &&
		start.Before(*resource.MaintenanceEnd) && resource.MaintenanceDueAt.Before(end) {
		conflicts = append(conflicts, types.ConflictItem{
			Code: constants.ConflictMaintenance,
			Message: fmt.Sprintf(constants.MsgResourceMaint, resource.ResourceCode,
				resource.MaintenanceDueAt.Format("01-02 15:04"), start.Format("15:04"), end.Format("15:04")),
			ResourceID: resource.ID, ResourceCode: resource.ResourceCode,
			BookingID: booking.ID, StartTime: utils.FormatTime(start), EndTime: utils.FormatTime(end),
		})
	}
	if end.After(turn.DepartureTime) {
		conflicts = append(conflicts, types.ConflictItem{
			Code: constants.ConflictDeadline,
			Message: fmt.Sprintf("预约结束 %s 晚于离港时间 %s",
				end.Format("15:04"), turn.DepartureTime.Format("15:04")),
			ResourceID: resource.ID, ResourceCode: resource.ResourceCode,
			BookingID: booking.ID, StartTime: utils.FormatTime(start), EndTime: utils.FormatTime(end),
		})
	}
	// Overlaps against other bookings (self excluded); other bookings of the
	// same turnaround are moving peers but are still distinct windows.
	others, err := s.Booking.Overlapping(resource.ID, start, end, booking.ID)
	if err == nil {
		for _, other := range others {
			conflicts = append(conflicts, types.ConflictItem{
				Code: constants.ConflictOverlap,
				Message: fmt.Sprintf(constants.MsgBookingConflict, resource.ResourceCode,
					start.Format("15:04"), end.Format("15:04"), other.ID),
				ResourceID: resource.ID, ResourceCode: resource.ResourceCode,
				BookingID: other.ID, StartTime: utils.FormatTime(start), EndTime: utils.FormatTime(end),
			})
		}
	}

	if len(conflicts) > 0 {
		// Keep it pending with the freshest primary reason; nothing else moves.
		booking.StartTime = start
		booking.EndTime = end
		booking.BookingStatus = constants.BookingPending
		booking.ConflictReason = conflicts[0].Code
		if err := s.Booking.Save(booking); err != nil {
			return nil, nil, err
		}
		return booking, conflicts, utils.NewAPIError(409, constants.CodeBookingConflict,
			fmt.Sprintf("预约 %d 调整时段仍有 %d 项冲突，保持待处理", booking.ID, len(conflicts)), conflicts)
	}

	booking.StartTime = start
	booking.EndTime = end
	booking.BookingStatus = constants.BookingConfirmed
	booking.ConflictReason = ""
	if err := s.Booking.Save(booking); err != nil {
		return nil, nil, err
	}
	// Align the owning task's plan window for tasks still unsigned.
	if task, terr := s.Task.Get(booking.TaskID); terr == nil &&
		(task.Status == constants.TaskPlanned || task.Status == constants.TaskBlocked) {
		task.PlannedStart = start
		task.Deadline = end
		task.Status = constants.TaskPlanned
		task.BlockerNote = ""
		_ = s.Task.Save(task)
	}

	s.audit(actor, role, fmt.Sprintf(constants.LogBookingAdjusted,
		booking.ID, resource.ResourceCode,
		start.Format("01-02 15:04"), end.Format("01-02 15:04"),
		constants.BookingStatusText[booking.BookingStatus]),
		"ResourceBooking", fmt.Sprintf("%d", booking.ID))
	return booking, nil, nil
}

// Release frees a booking manually.
func (s *BookingService) Release(bookingID uint, actor, role string) (*models.ResourceBooking, error) {
	booking, err := s.Booking.Get(bookingID)
	if err != nil {
		return nil, utils.NewAPIError(404, constants.CodeNotFound,
			fmt.Sprintf(constants.MsgNotFound, "资源预约", bookingID), nil)
	}
	booking.BookingStatus = constants.BookingReleased
	if err := s.Booking.Save(booking); err != nil {
		return nil, err
	}
	s.audit(actor, role, fmt.Sprintf(constants.LogBookingReleased, booking.ID),
		"ResourceBooking", fmt.Sprintf("%d", booking.ID))
	return booking, nil
}

func (s *BookingService) audit(actor, role, message, targetType, targetID string) {
	_ = s.Audit.Create(&models.AuditLog{
		Actor: actor, Role: role, Action: message, TargetType: targetType, TargetID: targetID,
	})
}
