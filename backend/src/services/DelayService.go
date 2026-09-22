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

// DelayService registers delay minutes and reschedules open tasks.
type DelayService struct {
	DB      *gorm.DB
	Plan    *PlanService
	Turn    *repositories.TurnaroundRepository
	Task    *repositories.TaskRepository
	Booking *repositories.BookingRepository
	Delay   *repositories.DelayRepository
	Resource *repositories.ResourceRepository
	Audit   *repositories.AuditRepository
}

func NewDelayService(db *gorm.DB, plan *PlanService) *DelayService {
	return &DelayService{
		DB:       db,
		Plan:     plan,
		Turn:     repositories.NewTurnaroundRepository(db),
		Task:     repositories.NewTaskRepository(db),
		Booking:  repositories.NewBookingRepository(db),
		Delay:    repositories.NewDelayRepository(db),
		Resource: repositories.NewResourceRepository(db),
		Audit:    repositories.NewAuditRepository(db),
	}
}

// RegisterDelay appends minutes to the turnaround and reschedules every
// unfinished, unsigned task. Signed/completed tasks keep their original
// time points; bookings squeezed out become PENDING with a conflict reason.
func (s *DelayService) RegisterDelay(turnaroundID uint, delayType string, minutes int, rootCause, team, actor, role string) (*types.RescheduleResult, error) {
	if minutes <= 0 {
		return nil, utils.NewAPIError(400, constants.CodeValidationFailed,
			fmt.Sprintf(constants.MsgValidationFailed, "延误分钟数必须大于 0"), nil)
	}
	validType := false
	for _, t := range constants.DelayType {
		if t == delayType {
			validType = true
		}
	}
	if !validType {
		return nil, utils.NewAPIError(400, constants.CodeValidationFailed,
			fmt.Sprintf(constants.MsgValidationFailed, "未知延误类型 "+delayType), nil)
	}

	turn, err := s.Turn.Get(turnaroundID)
	if err != nil {
		return nil, utils.NewAPIError(404, constants.CodeNotFound,
			fmt.Sprintf(constants.MsgNotFound, "航班过站", turnaroundID), nil)
	}
	tasks, err := s.Task.ListByTurnaround(turnaroundID)
	if err != nil {
		return nil, err
	}
	if len(tasks) == 0 {
		return nil, utils.NewAPIError(409, constants.CodeValidationFailed,
			"该航班尚未生成保障计划，无法重排", nil)
	}
	bookings, err := s.Booking.ListByTurnaround(turnaroundID)
	if err != nil {
		return nil, err
	}
	bookingByTask := map[uint]models.ResourceBooking{}
	for _, b := range bookings {
		bookingByTask[b.TaskID] = b
	}

	// New time anchors.
	shift := time.Duration(minutes) * time.Minute
	newArrival := turn.ArrivalTime.Add(shift)
	newDeparture := turn.DepartureTime.Add(shift)

	moved := 0
	frozen := 0
	pending := 0
	conflicts := []types.ConflictItem{}

	err = s.DB.Transaction(func(tx *gorm.DB) error {
		// Persist the delay event first.
		event := &models.DelayEvent{
			TurnaroundID:       turnaroundID,
			DelayType:          delayType,
			Minutes:            minutes,
			RootCause:          rootCause,
			ResponsibilityTeam: team,
		}
		if err := tx.Create(event).Error; err != nil {
			return err
		}

		// Shift only open (PLANNED / BLOCKED) tasks. ACCEPTED/COMPLETED are
		// signed: their planned time points and bookings stay unchanged.
		for i := range tasks {
			task := &tasks[i]
			if task.Status == constants.TaskAccepted || task.Status == constants.TaskCompleted {
				frozen++
				continue
			}
			task.PlannedStart = task.PlannedStart.Add(shift)
			task.Deadline = task.Deadline.Add(shift)
			if task.Status == constants.TaskBlocked {
				// A previously blocked, unsigned task is replanned by the delay.
				task.Status = constants.TaskPlanned
				task.BlockerNote = ""
			}
			if err := tx.Save(task).Error; err != nil {
				return err
			}
			moved++

			booking, ok := bookingByTask[task.ID]
			if !ok {
				continue
			}
			booking.StartTime = booking.StartTime.Add(shift)
			booking.EndTime = booking.EndTime.Add(shift)

			// Check the shifted window against OTHER bookings on the resource
			// (excluding this turnaround's own moving bookings; signed tasks on
			// other flights remain frozen blockers). All reads/writes use tx so
			// rows saved earlier in this batch are visible.
			blocked, reason, blocker := s.findShiftBlocker(tx, booking, newDeparture)
			if blocked {
				booking.BookingStatus = constants.BookingPending
				booking.ConflictReason = reason
				pending++
				item := types.ConflictItem{
					Code:         reason,
					Message:      s.pendingMessage(reason, booking, blocker),
					TaskType:     task.TaskType,
					ResourceID:   booking.ResourceID,
					BookingID:    booking.ID,
					StartTime:    utils.FormatTime(booking.StartTime),
					EndTime:      utils.FormatTime(booking.EndTime),
				}
				conflicts = append(conflicts, item)
			} else {
				booking.BookingStatus = constants.BookingConfirmed
				booking.ConflictReason = ""
			}
			if err := tx.Save(&booking).Error; err != nil {
				return err
			}
		}

		// Update turnaround clock and status.
		turn.ArrivalTime = newArrival
		turn.DepartureTime = newDeparture
		turn.TotalDelayMinutes += minutes
		if rootCause != "" {
			turn.DelayReason = rootCause
		}
		turn.TurnaroundStatus = constants.StatusDelayed
		if err := tx.Save(turn).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		if apiErr, ok := err.(*utils.APIError); ok {
			return nil, apiErr
		}
		return nil, utils.NewAPIError(500, constants.CodeInternal, constants.MsgInternal, err.Error())
	}

	// Reload refreshed state — every page then sees the same result.
	_, reTasks, reBookings, resourceMap, _ := s.Plan.GetDetail(turnaroundID)
	result := &types.RescheduleResult{
		RescheduledTasks:  moved,
		FrozenTasks:       frozen,
		PendingBookings:   pending,
		AddedDelayMinutes: minutes,
	}
	result.Saved = true
	result.TurnaroundID = turnaroundID
	result.Tasks = taskViews(reTasks)
	result.Bookings = bookingViews(reBookings, resourceMap)
	result.Conflicts = conflicts

	s.audit(actor, role, fmt.Sprintf(constants.LogDelayRegistered,
		turn.FlightNo, minutes, constants.DelayTypeText[delayType], moved, pending),
		"DelayEvent", fmt.Sprintf("%d", turnaroundID))
	return result, nil
}

// findShiftBlocker checks a shifted booking (inside tx) against resource
// state and frozen/active bookings belonging to other turnarounds.
func (s *DelayService) findShiftBlocker(tx *gorm.DB, booking models.ResourceBooking, newDeparture time.Time) (bool, string, *models.ResourceBooking) {
	var resource models.GroundResource
	if err := tx.First(&resource, booking.ResourceID).Error; err == nil {
		if resource.AvailabilityStatus == constants.ResourceOffline {
			return true, constants.ConflictOffline, nil
		}
		if resource.AvailabilityStatus == constants.ResourceMaintenance ||
			(resource.MaintenanceDueAt != nil && resource.MaintenanceEnd != nil &&
				booking.StartTime.Before(*resource.MaintenanceEnd) && resource.MaintenanceDueAt.Before(booking.EndTime)) {
			return true, constants.ConflictMaintenance, nil
		}
	}
	if booking.EndTime.After(newDeparture) {
		return true, constants.ConflictDeadline, nil
	}
	// Other bookings on the same resource; ignore this booking itself,
	// RELEASED bookings and bookings of the same turnaround (they move too).
	var others []models.ResourceBooking
	tx.Where("resource_id = ?", booking.ResourceID).
		Where("booking_status <> ?", constants.BookingReleased).
		Where("start_time < ? AND end_time > ?", booking.EndTime, booking.StartTime).
		Where("id <> ?", booking.ID).
		Where("turnaround_id <> ?", booking.TurnaroundID).
		Find(&others)
	for i := range others {
		return true, constants.ConflictOverlap, &others[i]
	}
	return false, "", nil
}

func (s *DelayService) pendingMessage(reason string, booking models.ResourceBooking, blocker *models.ResourceBooking) string {
	switch reason {
	case constants.ConflictOffline:
		return fmt.Sprintf("资源 %d 已离线，预约 %d 待改派", booking.ResourceID, booking.ID)
	case constants.ConflictMaintenance:
		return fmt.Sprintf("资源 %d 处于维护期，预约 %d 待改派", booking.ResourceID, booking.ID)
	case constants.ConflictDeadline:
		return fmt.Sprintf("预约 %d 重排后超出新离港时间", booking.ID)
	default:
		blockerID := uint(0)
		if blocker != nil {
			blockerID = blocker.ID
		}
		return fmt.Sprintf(constants.MsgBookingConflict, fmt.Sprintf("#%d", booking.ResourceID),
			booking.StartTime.Format("15:04"), booking.EndTime.Format("15:04"), blockerID)
	}
}

func (s *DelayService) audit(actor, role, message, targetType, targetID string) {
	_ = s.Audit.Create(&models.AuditLog{
		Actor: actor, Role: role, Action: message, TargetType: targetType, TargetID: targetID,
	})
}

// Resolve closes a delay event after attribution.
func (s *DelayService) Resolve(eventID uint, actor, role string) (*models.DelayEvent, error) {
	event, err := s.Delay.Get(eventID)
	if err != nil {
		return nil, utils.NewAPIError(404, constants.CodeNotFound,
			fmt.Sprintf(constants.MsgNotFound, "延误事件", eventID), nil)
	}
	now := time.Now()
	event.ResolvedAt = &now
	if err := s.Delay.Save(event); err != nil {
		return nil, err
	}
	s.audit(actor, role, fmt.Sprintf(constants.LogDelayResolved, event.ID, constants.DelayTypeText[event.DelayType]),
		"DelayEvent", fmt.Sprintf("%d", event.ID))
	return event, nil
}
