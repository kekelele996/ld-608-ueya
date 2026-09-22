package services

import (
	"sort"

	"groundTurn/src/constants"
	"groundTurn/src/models"
	"groundTurn/src/repositories"
	"groundTurn/src/types"
)

type DashboardStats struct {
	TurnaroundTotal int `json:"turnaround_total"`
	InService       int `json:"in_service"`
	Delayed         int `json:"delayed"`
	TaskTotal       int `json:"task_total"`
	TaskFinished    int `json:"task_finished"`
	TaskOverdue     int `json:"task_overdue"`
	TaskBlocked     int `json:"task_blocked"`
	PendingBookings int `json:"pending_bookings"`
	DelayMinutes    int `json:"delay_minutes"`
}

type ReportService struct {
	svc       *Services
	turnRepo  *repositories.FlightTurnaroundRepository
	taskRepo  *repositories.GroundTaskRepository
	bookRepo  *repositories.ResourceBookingRepository
	delayRepo *repositories.DelayEventRepository
	auditRepo *repositories.AuditLogRepository
}

func NewReportService(svc *Services) *ReportService {
	return &ReportService{
		svc:       svc,
		turnRepo:  repositories.NewFlightTurnaroundRepository(svc.DB),
		taskRepo:  repositories.NewGroundTaskRepository(svc.DB),
		bookRepo:  repositories.NewResourceBookingRepository(svc.DB),
		delayRepo: repositories.NewDelayEventRepository(svc.DB),
		auditRepo: repositories.NewAuditLogRepository(svc.DB),
	}
}

func (r *ReportService) Dashboard() (*DashboardStats, error) {
	turns, err := r.turnRepo.List()
	if err != nil {
		return nil, err
	}
	tasks, err := r.taskRepo.ListAll()
	if err != nil {
		return nil, err
	}
	delays, err := r.delayRepo.List()
	if err != nil {
		return nil, err
	}
	pending, err := r.bookRepo.ListPending()
	if err != nil {
		return nil, err
	}
	stats := &DashboardStats{
		TurnaroundTotal: len(turns),
		PendingBookings: len(pending),
	}
	for _, t := range turns {
		switch constants.TurnaroundStatus(t.TurnaroundStatus) {
		case constants.StatusInService:
			stats.InService++
		case constants.StatusDelayed:
			stats.Delayed++
		}
	}
	now := nowTime()
	for _, task := range tasks {
		stats.TaskTotal++
		if task.Status == string(constants.TaskFinished) {
			stats.TaskFinished++
		}
		if task.Status == string(constants.TaskBlocked) {
			stats.TaskBlocked++
		}
		if task.Status != string(constants.TaskFinished) && task.Deadline.Before(now) {
			stats.TaskOverdue++
		}
	}
	for _, d := range delays {
		stats.DelayMinutes += d.Minutes
	}
	return stats, nil
}

// DelayImpacts aggregates minutes by delay type for ChartPanel.
func (r *ReportService) DelayImpacts() ([]types.DelayImpact, error) {
	delays, err := r.delayRepo.List()
	if err != nil {
		return nil, err
	}
	agg := map[string]*types.DelayImpact{}
	for _, d := range delays {
		row, ok := agg[d.DelayType]
		if !ok {
			row = &types.DelayImpact{
				DelayType:          d.DelayType,
				ResponsibilityTeam: d.ResponsibilityTeam,
			}
			agg[d.DelayType] = row
		}
		row.EventCount++
		row.TotalMinutes += d.Minutes
	}
	out := make([]types.DelayImpact, 0, len(agg))
	for _, row := range agg {
		out = append(out, *row)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].TotalMinutes > out[j].TotalMinutes })
	return out, nil
}

func (r *ReportService) AuditLogs(limit int) ([]models.AuditLog, error) {
	return r.auditRepo.List(limit)
}
