package services

import (
	"errors"
	"fmt"
	"time"

	"groundTurn/src/constants"
	"groundTurn/src/constructors"
	"groundTurn/src/models"
	"groundTurn/src/repositories"
	"groundTurn/src/types"
	"groundTurn/src/utils"

	"gorm.io/gorm"
)

type TurnaroundService struct {
	svc *Services
	repo *repositories.FlightTurnaroundRepository
	taskRepo *repositories.GroundTaskRepository
	resRepo  *repositories.GroundResourceRepository
	bookRepo *repositories.ResourceBookingRepository
	delayRepo *repositories.DelayEventRepository
}

func NewTurnaroundService(svc *Services) *TurnaroundService {
	return &TurnaroundService{
		svc:       svc,
		repo:      repositories.NewFlightTurnaroundRepository(svc.DB),
		taskRepo:  repositories.NewGroundTaskRepository(svc.DB),
		resRepo:   repositories.NewGroundResourceRepository(svc.DB),
		bookRepo:  repositories.NewResourceBookingRepository(svc.DB),
		delayRepo: repositories.NewDelayEventRepository(svc.DB),
	}
}

func appError(status int, code, message string) *types.AppError {
	return &types.AppError{Code: code, Message: message, StatusCode: status}
}

func (s *TurnaroundService) List() ([]models.FlightTurnaround, error) {
	return s.repo.List()
}

func (s *TurnaroundService) Get(id uint) (*models.FlightTurnaround, error) {
	row, err := s.repo.Get(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, appError(404, constants.CodeNotFound,
			fmt.Sprintf(constants.MsgNotFound, "航班过站"))
	}
	return row, err
}

func (s *TurnaroundService) Detail(id uint) (*types.TurnaroundDetail, error) {
	if _, err := s.Get(id); err != nil {
		return nil, err
	}
	turnaround, _ := s.repo.Get(id)
	tasks, err := s.taskRepo.ListByTurnaround(s.svc.DB, id)
	if err != nil {
		return nil, err
	}
	bookings, err := s.bookRepo.ListByTurnaround(s.svc.DB, id)
	if err != nil {
		return nil, err
	}
	delays, err := s.delayRepo.ListByTurnaround(s.svc.DB, id)
	if err != nil {
		return nil, err
	}
	return &types.TurnaroundDetail{
		FlightTurnaround: *turnaround,
		Tasks:            tasks,
		Bookings:         bookings,
		Delays:           delays,
	}, nil
}

func (s *TurnaroundService) Create(actor Actor, req types.CreateTurnaroundRequest) (*models.FlightTurnaround, error) {
	if !req.DepartureTime.After(req.ArrivalTime) {
		return nil, appError(400, constants.CodeValidationFailed,
			fmt.Sprintf(constants.MsgValidationFailed, "离港时间必须晚于到站时间"))
	}
	row := constructors.NewFlightTurnaround(req)
	err := s.svc.DB.Transaction(func(tx *gorm.DB) error {
		if err := s.repo.Create(tx, row); err != nil {
			return err
		}
		s.svc.audit(tx, actor.Username, actor.Role, constants.ActionTurnaroundCreate,
			"FlightTurnaround", fmt.Sprintf("%d", row.ID),
			fmt.Sprintf(constants.LogTemplates[constants.ActionTurnaroundCreate],
				row.FlightNo, row.StandNo, utils.FormatTime(row.ArrivalTime)))
		return nil
	})
	if err != nil {
		return nil, appError(500, constants.CodeInternal, constants.MsgInternal)
	}
	utils.LogOperation(constants.ActionTurnaroundCreate, row.FlightNo, row.StandNo, utils.FormatTime(row.ArrivalTime))
	return row, nil
}

// Arrive confirms on-stand; the turnaround becomes eligible for plan generation.
func (s *TurnaroundService) Arrive(actor Actor, id uint, req types.ArriveRequest) (*models.FlightTurnaround, error) {
	row, err := s.Get(id)
	if err != nil {
		return nil, err
	}
	previous := row.TurnaroundStatus
	arrivedAt := time.Now()
	if req.ArrivalTime != nil {
		arrivedAt = *req.ArrivalTime
	}
	err = s.svc.DB.Transaction(func(tx *gorm.DB) error {
		row.ArrivalTime = arrivedAt
		if row.TurnaroundStatus == string(constants.StatusArriving) {
			row.TurnaroundStatus = string(constants.StatusOnStand)
		}
		if err := s.repo.Update(tx, row); err != nil {
			return err
		}
		s.svc.audit(tx, actor.Username, actor.Role, constants.ActionTurnaroundArrive,
			"FlightTurnaround", fmt.Sprintf("%d", row.ID),
			fmt.Sprintf(constants.LogTemplates[constants.ActionTurnaroundArrive],
				row.FlightNo, previous, row.TurnaroundStatus))
		return nil
	})
	if err != nil {
		return nil, appError(500, constants.CodeInternal, constants.MsgInternal)
	}
	utils.LogOperation(constants.ActionTurnaroundArrive, row.FlightNo, previous, row.TurnaroundStatus)
	return row, nil
}

// Release closes the turnaround once every task is finished.
func (s *TurnaroundService) Release(actor Actor, id uint) (*models.FlightTurnaround, error) {
	row, err := s.Get(id)
	if err != nil {
		return nil, err
	}
	tasks, err := s.taskRepo.ListByTurnaround(s.svc.DB, id)
	if err != nil {
		return nil, err
	}
	finished := 0
	for _, task := range tasks {
		if task.Status == string(constants.TaskFinished) {
			finished++
		}
	}
	if len(tasks) == 0 || finished != len(tasks) {
		return nil, appError(409, constants.CodeInvalidTransition,
			fmt.Sprintf(constants.MsgInvalidTransition,
				fmt.Sprintf("尚有 %d/%d 项任务未完成", len(tasks)-finished, len(tasks))))
	}
	err = s.svc.DB.Transaction(func(tx *gorm.DB) error {
		row.TurnaroundStatus = string(constants.StatusDeparted)
		if err := s.repo.Update(tx, row); err != nil {
			return err
		}
		s.svc.audit(tx, actor.Username, actor.Role, constants.ActionTurnaroundRelease,
			"FlightTurnaround", fmt.Sprintf("%d", row.ID),
			fmt.Sprintf(constants.LogTemplates[constants.ActionTurnaroundRelease],
				row.FlightNo, finished, len(tasks)))
		return nil
	})
	if err != nil {
		return nil, appError(500, constants.CodeInternal, constants.MsgInternal)
	}
	return row, nil
}

// GeneratePlan projects tasks per type/team/blueprint and resource appointments.
// ANY conflict (overlap, maintenance/offline resource, availability window or
// deadline) aborts the whole transaction: nothing is saved, every reason is
// returned item by item.
func (s *TurnaroundService) GeneratePlan(actor Actor, id uint) (*types.ConflictResult, error) {
	row, err := s.Get(id)
	if err != nil {
		return nil, err
	}
	if row.PlanGenerated {
		return nil, appError(409, constants.CodeInvalidTransition,
			fmt.Sprintf(constants.MsgInvalidTransition, "保障计划已生成，不能重复派发"))
	}

	result := &types.ConflictResult{Saved: false, Items: []types.ConflictItem{}}

	err = s.svc.DB.Transaction(func(tx *gorm.DB) error {
		// Existing confirmed bookings across the whole airport (other flights).
		existingBookings, err := s.bookRepo.ListAll(tx)
		if err != nil {
			return err
		}

		type projection struct {
			task     *models.GroundTask
			resource *models.GroundResource
		}
		projections := make([]projection, 0, len(constants.TaskBlueprints))

		for _, bp := range constants.TaskBlueprints {
			kind := bp.Kind
			task := constructors.NewGroundTask(row.ID, kind, row.ArrivalTime.Add(bp.StartOffset))

			// 1) Deadline conflict: projected end beyond turnaround departure/deadline.
			limit := row.DepartureTime
			if task.Deadline.After(limit) {
				result.Items = append(result.Items, types.ConflictItem{
					Code:         constants.ItemDeadlineConflict,
					Message: fmt.Sprintf(constants.ConflictMessages[constants.ItemDeadlineConflict],
						kind.Text(), utils.FormatTime(task.PlannedEnd), utils.FormatTime(limit)),
					TaskType:     string(kind),
					TeamID:       task.TeamID,
					PlannedStart: utils.FormatTime(task.PlannedStart),
					PlannedEnd:   utils.FormatTime(task.PlannedEnd),
					Deadline:     utils.FormatTime(limit),
				})
				continue
			}

			// 2) Resource selection by task type, team ownership and window.
			candidates, err := s.resRepo.ListByType(tx, constants.ResourceTypeForTask(kind))
			if err != nil {
				return err
			}
			var chosen *models.GroundResource
			var reasons []types.ConflictItem
			for i := range candidates {
				candidate := &candidates[i]
				code, message := resourceCovers(candidate, task.PlannedStart, task.PlannedEnd)
				if code != "" {
					reasons = append(reasons, types.ConflictItem{
						Code:         code,
						Message:      message,
						TaskType:     string(kind),
						TeamID:       task.TeamID,
						ResourceID:   ptrID(candidate.ID),
						ResourceCode: candidate.ResourceCode,
						PlannedStart: utils.FormatTime(task.PlannedStart),
						PlannedEnd:   utils.FormatTime(task.PlannedEnd),
					})
					continue
				}
				if hit := bookingOverlap(existingBookings, candidate.ID, task.PlannedStart, task.PlannedEnd, 0, 0); hit != nil {
					reasons = append(reasons, types.ConflictItem{
						Code: constants.ItemResourceOverlap,
						Message: fmt.Sprintf(constants.ConflictMessages[constants.ItemResourceOverlap],
							candidate.ResourceCode, utils.FormatTime(task.PlannedStart),
							utils.FormatTime(task.PlannedEnd), hit.TaskIDValue()),
						TaskType:     string(kind),
						TeamID:       task.TeamID,
						ResourceID:   ptrID(candidate.ID),
						ResourceCode: candidate.ResourceCode,
						PlannedStart: utils.FormatTime(task.PlannedStart),
						PlannedEnd:   utils.FormatTime(task.PlannedEnd),
					})
					continue
				}
				chosen = candidate
				break
			}

			if chosen == nil {
				if len(reasons) == 0 {
					result.Items = append(result.Items, types.ConflictItem{
						Code:     constants.ItemNoResource,
						Message:  fmt.Sprintf(constants.ConflictMessages[constants.ItemNoResource], kind.Text(), kind),
						TaskType: string(kind),
						TeamID:   task.TeamID,
					})
				} else {
					result.Items = append(result.Items, reasons...)
				}
				continue
			}
			projections = append(projections, projection{task: task, resource: chosen})
		}

		if len(result.Items) > 0 {
			// Returning an error rolls back: the entire batch is not saved.
			s.svc.audit(tx, actor.Username, actor.Role, constants.ActionPlanGenerate,
				"FlightTurnaround", fmt.Sprintf("%d", row.ID),
				fmt.Sprintf(constants.MsgPlanConflict, len(result.Items)))
			return errBatchRejected
		}

		for _, p := range projections {
			if err := s.taskRepo.Create(tx, p.task); err != nil {
				return err
			}
			booking := constructors.NewResourceBooking(p.resource.ID, row.ID, p.task.ID,
				p.task.PlannedStart, p.task.PlannedEnd)
			if err := s.bookRepo.Create(tx, booking); err != nil {
				return err
			}
			s.svc.audit(tx, actor.Username, actor.Role, constants.ActionBookingCreate,
				"ResourceBooking", fmt.Sprintf("%d", booking.ID),
				fmt.Sprintf(constants.LogTemplates[constants.ActionBookingCreate],
					p.resource.ResourceCode, utils.FormatTime(booking.StartTime),
					utils.FormatTime(booking.EndTime), p.task.ID))
		}

		row.PlanGenerated = true
		row.TurnaroundStatus = string(constants.StatusInService)
		if err := s.repo.Update(tx, row); err != nil {
			return err
		}
		s.svc.audit(tx, actor.Username, actor.Role, constants.ActionTaskDispatch,
			"FlightTurnaround", fmt.Sprintf("%d", row.ID),
			fmt.Sprintf(constants.LogTemplates[constants.ActionTaskDispatch],
				row.FlightNo, len(projections)))
		result.Saved = true
		result.TaskCount = len(projections)
		return nil
	})

	if errors.Is(err, errBatchRejected) {
		return result, nil // 409-ish payload handled by controller; items explain everything
	}
	if err != nil {
		return nil, appError(500, constants.CodeInternal, constants.MsgInternal)
	}
	utils.LogOperation(constants.ActionPlanGenerate, row.FlightNo, result.TaskCount)
	return result, nil
}

var errBatchRejected = errors.New("batch rejected due to conflicts")

func ptrID(id uint) *uint { return &id }
