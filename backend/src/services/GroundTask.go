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

// TaskService handles sign-off, completion and blocking.
type TaskService struct {
	DB       *gorm.DB
	Task     *repositories.TaskRepository
	Booking  *repositories.BookingRepository
	Turn     *repositories.TurnaroundRepository
	Audit    *repositories.AuditRepository
}

func NewTaskService(db *gorm.DB) *TaskService {
	return &TaskService{
		DB:      db,
		Task:    repositories.NewTaskRepository(db),
		Booking: repositories.NewBookingRepository(db),
		Turn:    repositories.NewTurnaroundRepository(db),
		Audit:   repositories.NewAuditRepository(db),
	}
}

// Accept signs a PLANNED task. Once signed, delay reschedules must not move
// its planned time point.
func (s *TaskService) Accept(taskID uint, actor, role, teamID string) (*models.GroundTask, error) {
	task, err := s.Task.Get(taskID)
	if err != nil {
		return nil, utils.NewAPIError(404, constants.CodeNotFound,
			fmt.Sprintf(constants.MsgNotFound, "地勤任务", taskID), nil)
	}
	if task.Status != constants.TaskPlanned {
		return nil, utils.NewAPIError(409, constants.CodeInvalidTransition,
			fmt.Sprintf(constants.MsgInvalidTransition, task.ID, task.Status, "签收"), nil)
	}
	now := time.Now()
	task.Status = constants.TaskAccepted
	task.SignedAt = &now
	if err := s.Task.Save(task); err != nil {
		return nil, err
	}
	turn, _ := s.Turn.Get(task.TurnaroundID)
	flightNo := ""
	if turn != nil {
		flightNo = turn.FlightNo
	}
	s.audit(actor, role, fmt.Sprintf(constants.LogTaskAccepted,
		task.ID, constants.GroundTaskTypeText[task.TaskType], task.TeamID,
		task.PlannedStart.Format("01-02 15:04")), "GroundTask", fmt.Sprintf("%d/%s", task.ID, flightNo))
	return task, nil
}

// Complete finishes a signed task and releases its resource booking.
func (s *TaskService) Complete(taskID uint, actor, role string) (*models.GroundTask, error) {
	task, err := s.Task.Get(taskID)
	if err != nil {
		return nil, utils.NewAPIError(404, constants.CodeNotFound,
			fmt.Sprintf(constants.MsgNotFound, "地勤任务", taskID), nil)
	}
	if task.Status != constants.TaskAccepted {
		return nil, utils.NewAPIError(409, constants.CodeInvalidTransition,
			fmt.Sprintf(constants.MsgInvalidTransition, task.ID, task.Status, "完成确认"), nil)
	}
	now := time.Now()
	task.Status = constants.TaskCompleted
	task.ActualFinish = &now
	err = s.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(task).Error; err != nil {
			return err
		}
		bookings, err := s.Booking.ListByTask(taskID)
		if err != nil {
			return err
		}
		for i := range bookings {
			if bookings[i].BookingStatus == constants.BookingReleased {
				continue
			}
			bookings[i].BookingStatus = constants.BookingReleased
			if err := tx.Save(&bookings[i]).Error; err != nil {
				return err
			}
			s.audit(actor, role, fmt.Sprintf(constants.LogBookingReleased, bookings[i].ID),
				"ResourceBooking", fmt.Sprintf("%d", bookings[i].ID))
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	s.audit(actor, role, fmt.Sprintf(constants.LogTaskCompleted,
		task.ID, constants.GroundTaskTypeText[task.TaskType],
		now.Format("15:04"), 0), "GroundTask", fmt.Sprintf("%d", task.ID))
	return task, nil
}

// Block marks a signed task blocked with a reason.
func (s *TaskService) Block(taskID uint, note, actor, role string) (*models.GroundTask, error) {
	task, err := s.Task.Get(taskID)
	if err != nil {
		return nil, utils.NewAPIError(404, constants.CodeNotFound,
			fmt.Sprintf(constants.MsgNotFound, "地勤任务", taskID), nil)
	}
	if task.Status != constants.TaskAccepted && task.Status != constants.TaskPlanned {
		return nil, utils.NewAPIError(409, constants.CodeInvalidTransition,
			fmt.Sprintf(constants.MsgInvalidTransition, task.ID, task.Status, "阻塞"), nil)
	}
	task.Status = constants.TaskBlocked
	task.BlockerNote = note
	if err := s.Task.Save(task); err != nil {
		return nil, err
	}
	s.audit(actor, role, fmt.Sprintf(constants.LogTaskBlocked,
		task.ID, constants.GroundTaskTypeText[task.TaskType], note),
		"GroundTask", fmt.Sprintf("%d", task.ID))
	return task, nil
}

func (s *TaskService) audit(actor, role, message, targetType, targetID string) {
	_ = s.Audit.Create(&models.AuditLog{
		Actor: actor, Role: role, Action: message, TargetType: targetType, TargetID: targetID,
	})
}
