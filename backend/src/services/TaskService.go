package services

import (
	"errors"
	"fmt"
	"time"

	"groundTurn/src/constants"
	"groundTurn/src/models"
	"groundTurn/src/repositories"
	"groundTurn/src/types"

	"gorm.io/gorm"
)

type TaskService struct {
	svc      *Services
	repo     *repositories.GroundTaskRepository
	bookRepo *repositories.ResourceBookingRepository
}

func NewTaskService(svc *Services) *TaskService {
	return &TaskService{
		svc:      svc,
		repo:     repositories.NewGroundTaskRepository(svc.DB),
		bookRepo: repositories.NewResourceBookingRepository(svc.DB),
	}
}

func (s *TaskService) ListAll() ([]models.GroundTask, error) { return s.repo.ListAll() }

func (s *TaskService) ListByTurnaround(id uint) ([]models.GroundTask, error) {
	return s.repo.ListByTurnaround(s.svc.DB, id)
}

func (s *TaskService) Get(id uint) (*models.GroundTask, error) {
	row, err := s.repo.Get(s.svc.DB, id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, appError(404, constants.CodeNotFound,
			fmt.Sprintf(constants.MsgNotFound, "地勤任务"))
	}
	return row, err
}

// Sign marks dispatch acceptance by the crew. Acknowledged tasks are frozen
// during later delay replanning.
func (s *TaskService) Sign(actor Actor, id uint, req types.SignTaskRequest) (*models.GroundTask, error) {
	task, err := s.Get(id)
	if err != nil {
		return nil, err
	}
	if task.Status != string(constants.TaskPlanned) {
		return nil, appError(409, constants.CodeInvalidTransition,
			fmt.Sprintf(constants.MsgInvalidTransition,
				fmt.Sprintf("任务当前为「%s」，仅待签收任务可签收", constants.GroundTaskStatus(task.Status).Text())))
	}
	operator := actor.Label()
	if req.Operator != "" {
		operator = req.Operator
	}
	err = s.svc.DB.Transaction(func(tx *gorm.DB) error {
		task.Status = string(constants.TaskSigned)
		task.SignedBy = operator
		if err := s.repo.Update(tx, task); err != nil {
			return err
		}
		s.svc.audit(tx, actor.Username, actor.Role, constants.ActionTaskSign,
			"GroundTask", fmt.Sprintf("%d", task.ID),
			fmt.Sprintf(constants.LogTemplates[constants.ActionTaskSign],
				task.ID, constants.GroundTaskType(task.TaskType).Text(),
				task.TeamID, formatTimeL(task.PlannedStart)))
		return nil
	})
	if err != nil {
		return nil, appError(500, constants.CodeInternal, constants.MsgInternal)
	}
	return task, nil
}

// Finish records actual completion time.
func (s *TaskService) Finish(actor Actor, id uint, req types.FinishTaskRequest) (*models.GroundTask, error) {
	task, err := s.Get(id)
	if err != nil {
		return nil, err
	}
	if task.Status != string(constants.TaskSigned) && task.Status != string(constants.TaskBlocked) {
		return nil, appError(409, constants.CodeInvalidTransition,
			fmt.Sprintf(constants.MsgInvalidTransition, "仅已签收/阻塞任务可完成"))
	}
	now := time.Now()
	err = s.svc.DB.Transaction(func(tx *gorm.DB) error {
		task.Status = string(constants.TaskFinished)
		task.ActualFinish = &now
		task.BlockerNote = req.Note
		if err := s.repo.Update(tx, task); err != nil {
			return err
		}
		s.svc.audit(tx, actor.Username, actor.Role, constants.ActionTaskFinish,
			"GroundTask", fmt.Sprintf("%d", task.ID),
			fmt.Sprintf(constants.LogTemplates[constants.ActionTaskFinish],
				task.ID, constants.GroundTaskType(task.TaskType).Text(), now.Format(timeLayout)))
		return nil
	})
	if err != nil {
		return nil, appError(500, constants.CodeInternal, constants.MsgInternal)
	}
	return task, nil
}

// Block records why the task cannot proceed.
func (s *TaskService) Block(actor Actor, id uint, req types.BlockTaskRequest) (*models.GroundTask, error) {
	task, err := s.Get(id)
	if err != nil {
		return nil, err
	}
	if task.Status == string(constants.TaskFinished) {
		return nil, appError(409, constants.CodeInvalidTransition,
			fmt.Sprintf(constants.MsgInvalidTransition, "已完成任务不可阻塞"))
	}
	err = s.svc.DB.Transaction(func(tx *gorm.DB) error {
		task.Status = string(constants.TaskBlocked)
		task.BlockerNote = req.Note
		if err := s.repo.Update(tx, task); err != nil {
			return err
		}
		s.svc.audit(tx, actor.Username, actor.Role, constants.ActionTaskBlock,
			"GroundTask", fmt.Sprintf("%d", task.ID),
			fmt.Sprintf(constants.LogTemplates[constants.ActionTaskBlock],
				task.ID, constants.GroundTaskType(task.TaskType).Text(), req.Note))
		return nil
	})
	if err != nil {
		return nil, appError(500, constants.CodeInternal, constants.MsgInternal)
	}
	return task, nil
}

const timeLayout = "2006-01-02 15:04"

func formatTimeL(t time.Time) string { return t.Local().Format(timeLayout) }
