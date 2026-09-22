package services

import (
	"fmt"

	"groundTurn/src/constants"
	"groundTurn/src/constructors"
	"groundTurn/src/models"
	"groundTurn/src/repositories"
	"groundTurn/src/types"
	"groundTurn/src/utils"

	"gorm.io/gorm"
)

type DelayService struct {
	svc       *Services
	repo      *repositories.DelayEventRepository
	turnRepo  *repositories.FlightTurnaroundRepository
	taskRepo  *repositories.GroundTaskRepository
	bookRepo  *repositories.ResourceBookingRepository
	resRepo   *repositories.GroundResourceRepository
}

func NewDelayService(svc *Services) *DelayService {
	return &DelayService{
		svc:      svc,
		repo:     repositories.NewDelayEventRepository(svc.DB),
		turnRepo: repositories.NewFlightTurnaroundRepository(svc.DB),
		taskRepo: repositories.NewGroundTaskRepository(svc.DB),
		bookRepo: repositories.NewResourceBookingRepository(svc.DB),
		resRepo:  repositories.NewGroundResourceRepository(svc.DB),
	}
}

func (s *DelayService) List() ([]models.DelayEvent, error) { return s.repo.List() }

// Register records a delay in minutes and replans unfinished work inside one
// transaction:
//   - PLANNED/BLOCKED tasks and their appointments shift forward by `minutes`;
//   - SIGNED/FINISHED tasks keep their original time points (already 签收);
//   - appointments colliding with another flight after the shift are bumped
//     to PENDING with a conflict reason for manual triage.
func (s *DelayService) Register(actor Actor, turnaroundID uint, req types.RegisterDelayRequest) (*types.ReplanResult, error) {
	kind := constants.DelayType(req.DelayType)
	if !kind.Valid() {
		return nil, appError(400, constants.CodeValidationFailed,
			fmt.Sprintf(constants.MsgValidationFailed, "未知延误类型 "+req.DelayType))
	}
	turnaround, err := s.turnRepo.Get(turnaroundID)
	if err != nil {
		return nil, appError(404, constants.CodeNotFound,
			fmt.Sprintf(constants.MsgNotFound, "航班过站"))
	}

	event := constructors.NewDelayEvent(turnaroundID, req)
	result := &types.ReplanResult{
		ShiftMinutes:  req.Minutes,
		KeptTasks:     []types.KeptTaskView{},
		MovedTasks:    []models.GroundTask{},
		MovedBookings: []models.ResourceBooking{},
		Bumped:        []types.ConflictItem{},
	}
	shift := minutesDuration(req.Minutes)

	err = s.svc.DB.Transaction(func(tx *gorm.DB) error {
		if err := s.repo.Create(tx, event); err != nil {
			return err
		}
		result.DelayEventID = event.ID
		s.svc.audit(tx, actor.Username, actor.Role, constants.ActionDelayRegister,
			"DelayEvent", fmt.Sprintf("%d", event.ID),
			fmt.Sprintf(constants.LogTemplates[constants.ActionDelayRegister],
				turnaround.FlightNo, req.Minutes, kind.Text(), req.ResponsibilityTeam))

		// Shift the turnaround's schedule and mark DELAYED.
		turnaround.ArrivalTime = turnaround.ArrivalTime.Add(shift)
		turnaround.DepartureTime = turnaround.DepartureTime.Add(shift)
		turnaround.AccumulatedDelay += req.Minutes
		turnaround.DelayReason = kind.Text()
		if turnaround.TurnaroundStatus != string(constants.StatusDeparted) {
			turnaround.TurnaroundStatus = string(constants.StatusDelayed)
		}
		if err := s.turnRepo.Update(tx, turnaround); err != nil {
			return err
		}

		tasks, err := s.taskRepo.ListByTurnaround(tx, turnaroundID)
		if err != nil {
			return err
		}
		bookings, err := s.bookRepo.ListByTurnaround(tx, turnaroundID)
		if err != nil {
			return err
		}
		allBookings, err := s.bookRepo.ListAll(tx)
		if err != nil {
			return err
		}

		for _, task := range tasks {
			status := constants.GroundTaskStatus(task.Status)
			if status.Acknowledged() {
				result.KeptTasks = append(result.KeptTasks, types.KeptTaskView{
					TaskID:       task.ID,
					TaskType:     task.TaskType,
					Status:       task.Status,
					PlannedStart: task.PlannedStart,
					SignedBy:     task.SignedBy,
					Reason:       "已签收任务保持原计划时点，不参与延误重排",
				})
				continue
			}

			// Unfinished (PLANNED/BLOCKED): shift planned times and deadline.
			task.PlannedStart = task.PlannedStart.Add(shift)
			task.PlannedEnd = task.PlannedEnd.Add(shift)
			task.Deadline = task.Deadline.Add(shift)
			if err := s.taskRepo.Update(tx, &task); err != nil {
				return err
			}
			result.MovedTasks = append(result.MovedTasks, task)

			// Shift its resource appointment and test the new window.
			for i := range bookings {
				booking := &bookings[i]
				if booking.TaskID == nil || *booking.TaskID != task.ID {
					continue
				}
				if booking.BookingStatus == string(constants.BookingReleased) {
					continue
				}
				booking.StartTime = booking.StartTime.Add(shift)
				booking.EndTime = booking.EndTime.Add(shift)

				resource, rerr := s.resRepo.Get(tx, booking.ResourceID)
				if rerr == nil {
					if code, message := resourceCovers(resource, booking.StartTime, booking.EndTime); code != "" {
						bumpBooking(booking, message)
						result.Bumped = append(result.Bumped, types.ConflictItem{
							Code:         code,
							Message:      message,
							TaskType:     task.TaskType,
							TaskID:       ptrID(task.ID),
							BookingID:    ptrID(booking.ID),
							ResourceID:   ptrID(resource.ID),
							ResourceCode: resource.ResourceCode,
							PlannedStart: utils.FormatTime(booking.StartTime),
							PlannedEnd:   utils.FormatTime(booking.EndTime),
						})
					} else if hit := bookingOverlap(allBookings, resource.ID,
						booking.StartTime, booking.EndTime, booking.ID, turnaroundID); hit != nil {
						message := fmt.Sprintf(constants.ConflictMessages[constants.ItemBookingBumped],
							resource.ResourceCode, utils.FormatTime(booking.StartTime),
							utils.FormatTime(booking.EndTime))
						bumpBooking(booking, message)
						result.Bumped = append(result.Bumped, types.ConflictItem{
							Code:         constants.ItemResourceOverlap,
							Message:      message,
							TaskType:     task.TaskType,
							TaskID:       ptrID(task.ID),
							BookingID:    ptrID(booking.ID),
							ResourceID:   ptrID(resource.ID),
							ResourceCode: resource.ResourceCode,
							PlannedStart: utils.FormatTime(booking.StartTime),
							PlannedEnd:   utils.FormatTime(booking.EndTime),
						})
					}
				}

				if err := s.bookRepo.Update(tx, booking); err != nil {
					return err
				}
				result.MovedBookings = append(result.MovedBookings, *booking)

				action := constants.ActionBookingAdjust
				if booking.BookingStatus == string(constants.BookingPending) {
					action = constants.ActionBookingBump
				}
				s.svc.audit(tx, actor.Username, actor.Role, action,
					"ResourceBooking", fmt.Sprintf("%d", booking.ID),
					fmt.Sprintf(constants.LogTemplates[action],
						booking.ID, booking.ResourceID, booking.ConflictReason))
			}
		}

		s.svc.audit(tx, actor.Username, actor.Role, constants.ActionDelayReplan,
			"FlightTurnaround", fmt.Sprintf("%d", turnaroundID),
			fmt.Sprintf(constants.LogTemplates[constants.ActionDelayReplan],
				turnaround.FlightNo, req.Minutes,
				len(result.MovedTasks), len(result.Bumped)))
		return nil
	})
	if err != nil {
		return nil, appError(500, constants.CodeInternal, constants.MsgInternal)
	}

	utils.LogOperation(constants.ActionDelayRegister,
		turnaround.FlightNo, req.Minutes, kind.Text(), req.ResponsibilityTeam)
	utils.LogOperation(constants.ActionDelayReplan,
		turnaround.FlightNo, req.Minutes, len(result.MovedTasks), len(result.Bumped))
	return result, nil
}

func bumpBooking(booking *models.ResourceBooking, reason string) {
	booking.BookingStatus = string(constants.BookingPending)
	booking.ConflictReason = reason
}

// Resolve closes a delay event.
func (s *DelayService) Resolve(actor Actor, id uint) (*models.DelayEvent, error) {
	var event models.DelayEvent
	if err := s.svc.DB.First(&event, id).Error; err != nil {
		return nil, appError(404, constants.CodeNotFound,
			fmt.Sprintf(constants.MsgNotFound, "延误事件"))
	}
	err := s.svc.DB.Transaction(func(tx *gorm.DB) error {
		resolved := resolveNow()
		event.ResolvedAt = &resolved
		if err := s.repo.Update(tx, &event); err != nil {
			return err
		}
		s.svc.audit(tx, actor.Username, actor.Role, constants.ActionDelayResolve,
			"DelayEvent", fmt.Sprintf("%d", event.ID),
			fmt.Sprintf(constants.LogTemplates[constants.ActionDelayResolve],
				event.ID, event.TurnaroundID))
		return nil
	})
	if err != nil {
		return nil, appError(500, constants.CodeInternal, constants.MsgInternal)
	}
	return &event, nil
}
