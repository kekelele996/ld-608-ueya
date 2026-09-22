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

// PlanService owns the turnaround closed loop: arrival registration,
// whole-batch task+booking generation with all-or-nothing validation,
// delay-driven reschedule and task sign-off.
type PlanService struct {
	DB        *gorm.DB
	Turn      *repositories.TurnaroundRepository
	Task      *repositories.TaskRepository
	Resource  *repositories.ResourceRepository
	Booking   *repositories.BookingRepository
	Delay     *repositories.DelayRepository
	Audit     *repositories.AuditRepository
}

func NewPlanService(db *gorm.DB) *PlanService {
	return &PlanService{
		DB:       db,
		Turn:     repositories.NewTurnaroundRepository(db),
		Task:     repositories.NewTaskRepository(db),
		Resource: repositories.NewResourceRepository(db),
		Booking:  repositories.NewBookingRepository(db),
		Delay:    repositories.NewDelayRepository(db),
		Audit:    repositories.NewAuditRepository(db),
	}
}

type candidateTask struct {
	taskType string
	teamID   string
	start    time.Time
	end      time.Time
	deadline time.Time
}

// assignmentWindow is a simulated in-batch resource occupation used while
// validating a whole plan before anything is persisted.
type assignmentWindow struct {
	start, end time.Time
	taskType   string
}

// Register creates a turnaround after flight arrival.
func (s *PlanService) Register(flightNo, reg, stand string, arrival, departure time.Time, actor, role string) (*models.FlightTurnaround, error) {
	if !departure.After(arrival) {
		return nil, utils.NewAPIError(400, constants.CodeValidationFailed,
			fmt.Sprintf(constants.MsgValidationFailed, "离港时间必须晚于到站时间"), nil)
	}
	t := &models.FlightTurnaround{
		FlightNo:         flightNo,
		AircraftReg:      reg,
		StandNo:          stand,
		ArrivalTime:      arrival,
		DepartureTime:    departure,
		TurnaroundStatus: constants.StatusOnStand,
	}
	if err := s.Turn.Create(t); err != nil {
		return nil, utils.NewAPIError(500, constants.CodeInternal, constants.MsgInternal, err.Error())
	}
	s.audit(actor, role, fmt.Sprintf(constants.LogTurnaroundCreated,
		t.FlightNo, stand, constants.TurnaroundStatusText[t.TurnaroundStatus]),
		"FlightTurnaround", fmt.Sprintf("%d", t.ID))
	return t, nil
}

// GetDetail returns turnaround + tasks + bookings for pages and refreshes.
func (s *PlanService) GetDetail(id uint) (*models.FlightTurnaround, []models.GroundTask, []models.ResourceBooking, map[uint]models.GroundResource, error) {
	turn, err := s.Turn.Get(id)
	if err != nil {
		return nil, nil, nil, nil, utils.NewAPIError(404, constants.CodeNotFound,
			fmt.Sprintf(constants.MsgNotFound, "航班过站", id), nil)
	}
	tasks, err := s.Task.ListByTurnaround(id)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	bookings, err := s.Booking.ListByTurnaround(id)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	resourceMap := map[uint]models.GroundResource{}
	for _, b := range bookings {
		if _, ok := resourceMap[b.ResourceID]; ok {
			continue
		}
		if r, rerr := s.Resource.Get(b.ResourceID); rerr == nil {
			resourceMap[b.ResourceID] = *r
		}
	}
	return turn, tasks, bookings, resourceMap, nil
}

// GeneratePlan builds all tasks and bookings for a turnaround in one batch.
// ANY conflict rejects the entire batch (transaction rollback) and every
// conflict reason is returned item by item.
func (s *PlanService) GeneratePlan(turnaroundID uint, taskTypes []string, actor, role string) (*types.PlanResult, error) {
	turn, err := s.Turn.Get(turnaroundID)
	if err != nil {
		return nil, utils.NewAPIError(404, constants.CodeNotFound,
			fmt.Sprintf(constants.MsgNotFound, "航班过站", turnaroundID), nil)
	}
	if existing, _ := s.Task.CountByTurnaround(turnaroundID); existing > 0 {
		return nil, utils.NewAPIError(409, constants.CodeConflict,
			fmt.Sprintf(constants.MsgAlreadyExists, turn.FlightNo), nil)
	}
	if len(taskTypes) == 0 {
		taskTypes = constants.GroundTaskType
	}

	candidates := make([]candidateTask, 0, len(taskTypes))
	for _, taskType := range taskTypes {
		if !constants.IsValidTaskType(taskType) {
			return nil, utils.NewAPIError(400, constants.CodeValidationFailed,
				fmt.Sprintf(constants.MsgValidationFailed, "未知任务类型 "+taskType), nil)
		}
		spec := constants.TaskOffsetMinutes[taskType]
		start := turn.ArrivalTime.Add(time.Duration(spec[0]) * time.Minute)
		end := start.Add(time.Duration(spec[1]) * time.Minute)
		candidates = append(candidates, candidateTask{
			taskType: taskType,
			teamID:   constants.TaskDefaultTeam[taskType],
			start:    start,
			end:      end,
			deadline: end,
		})
	}

	conflicts := s.validateCandidates(candidates, turn, nil)

	result := &types.PlanResult{TurnaroundID: turnaroundID, Conflicts: conflicts}
	if len(conflicts) > 0 {
		s.audit(actor, role, fmt.Sprintf(constants.LogTaskPlanRejected, turn.FlightNo, len(conflicts)),
			"GroundTask", fmt.Sprintf("%d", turnaroundID))
		result.Saved = false
		return result, utils.NewAPIError(409, constants.CodeConflict,
			fmt.Sprintf(constants.MsgPlanConflict, len(conflicts)), conflicts)
	}

	// Persist the whole batch atomically.
	err = s.DB.Transaction(func(tx *gorm.DB) error {
		tasks := make([]models.GroundTask, 0, len(candidates))
		for _, cand := range candidates {
			tasks = append(tasks, models.GroundTask{
				TurnaroundID: turnaroundID,
				TaskType:     cand.taskType,
				TeamID:       cand.teamID,
				PlannedStart: cand.start,
				Deadline:     cand.deadline,
				Status:       constants.TaskPlanned,
			})
		}
		if err := tx.Create(&tasks).Error; err != nil {
			return err
		}
		bookings := make([]models.ResourceBooking, 0, len(tasks))
		// Mirror the validator's in-batch assignments so two new tasks never
		// grab the same resource window.
		persistedAssign := map[uint][]assignmentWindow{}
		for i, cand := range candidates {
			chosen := s.chooseResource(cand, nil, persistedAssign)
			if chosen == 0 {
				return utils.NewAPIError(409, constants.CodeConflict, constants.MsgPlanConflict, nil)
			}
			persistedAssign[chosen] = append(persistedAssign[chosen],
				assignmentWindow{start: cand.start, end: cand.end, taskType: cand.taskType})
			bookings = append(bookings, models.ResourceBooking{
				ResourceID:    chosen,
				TurnaroundID:  turnaroundID,
				TaskID:        tasks[i].ID,
				StartTime:     cand.start,
				EndTime:       cand.end,
				BookingStatus: constants.BookingConfirmed,
			})
		}
		if err := tx.Create(&bookings).Error; err != nil {
			return err
		}
		if turn.TurnaroundStatus == constants.StatusOnStand {
			turn.TurnaroundStatus = constants.StatusInService
			if err := tx.Save(turn).Error; err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		if apiErr, ok := err.(*utils.APIError); ok {
			return nil, apiErr
		}
		return nil, utils.NewAPIError(500, constants.CodeInternal, constants.MsgInternal, err.Error())
	}

	// Reload for the response.
	_, tasks, bookings, resourceMap, _ := s.GetDetail(turnaroundID)
	result.Saved = true
	result.Tasks = taskViews(tasks)
	result.Bookings = bookingViews(bookings, resourceMap)
	s.audit(actor, role, fmt.Sprintf(constants.LogTaskPlanGenerated,
		turn.FlightNo, len(tasks), len(bookings)),
		"GroundTask", fmt.Sprintf("%d", turnaroundID))
	return result, nil
}

// validateCandidates checks each candidate against deadline rules and all
// resources of its required type. proposed (when non-nil) represents the
// in-flight batch so two new tasks cannot grab the same resource window.
func (s *PlanService) validateCandidates(candidates []candidateTask, turn *models.FlightTurnaround, proposed map[string][]models.ResourceBooking) []types.ConflictItem {
	conflicts := []types.ConflictItem{}
	// simulated assignments within this batch: resourceID -> windows
	localAssign := map[uint][]assignmentWindow{}

	for _, cand := range candidates {
		if !cand.deadline.Before(turn.DepartureTime) && !cand.deadline.Equal(turn.DepartureTime) {
			conflicts = append(conflicts, types.ConflictItem{
				Code:      constants.ConflictDeadline,
				Message:   fmt.Sprintf(constants.MsgDeadlineConflict, constants.GroundTaskTypeText[cand.taskType], cand.deadline.Format("15:04"), turn.DepartureTime.Format("15:04")),
				TaskType:  cand.taskType,
				StartTime: utils.FormatTime(cand.start),
				EndTime:   utils.FormatTime(cand.end),
			})
			continue
		}
		if cand.start.Before(turn.ArrivalTime) {
			conflicts = append(conflicts, types.ConflictItem{
				Code:      constants.ConflictBeforeArrival,
				Message:   fmt.Sprintf(constants.MsgTaskOrder, constants.GroundTaskTypeText[cand.taskType], cand.start.Format("15:04"), turn.ArrivalTime.Format("15:04")),
				TaskType:  cand.taskType,
				StartTime: utils.FormatTime(cand.start),
				EndTime:   utils.FormatTime(cand.end),
			})
			continue
		}

		resourceType := constants.TaskResourceType[cand.taskType]
		resources, err := s.Resource.ListByType(resourceType)
		if err != nil || len(resources) == 0 {
			conflicts = append(conflicts, types.ConflictItem{
				Code:     constants.ConflictNoResource,
				Message:  fmt.Sprintf(constants.MsgNoResource, constants.GroundTaskTypeText[cand.taskType], resourceType),
				TaskType: cand.taskType,
			})
			continue
		}

		bestID, perResource := s.pickResource(resources, cand, proposed, localAssign)
		if bestID == 0 {
			// Surface the most specific reason per resource, itemized.
			for _, item := range perResource {
				conflicts = append(conflicts, item)
			}
			if len(perResource) == 0 {
				conflicts = append(conflicts, types.ConflictItem{
					Code:     constants.ConflictNoResource,
					Message:  fmt.Sprintf(constants.MsgNoResource, constants.GroundTaskTypeText[cand.taskType], resourceType),
					TaskType: cand.taskType,
				})
			}
			continue
		}
		localAssign[bestID] = append(localAssign[bestID], assignmentWindow{cand.start, cand.end, cand.taskType})
	}
	return conflicts
}

// pickResource returns the first usable resource id, collecting itemized
// reasons explaining why every other resource was rejected.
func (s *PlanService) pickResource(resources []models.GroundResource, cand candidateTask, proposed map[string][]models.ResourceBooking, localAssign map[uint][]assignmentWindow) (uint, []types.ConflictItem) {
	items := []types.ConflictItem{}
	usable := uint(0)
	for _, resource := range resources {
		reason := s.resourceBlockReason(resource, cand, proposed, localAssign)
		if reason == nil {
			if usable == 0 {
				usable = resource.ID
			}
			continue
		}
		reason.TaskType = cand.taskType
		reason.ResourceID = resource.ID
		reason.ResourceCode = resource.ResourceCode
		reason.StartTime = utils.FormatTime(cand.start)
		reason.EndTime = utils.FormatTime(cand.end)
		items = append(items, *reason)
	}
	return usable, items
}

// resourceBlockReason returns a ConflictItem when this resource cannot host
// the candidate window.
func (s *PlanService) resourceBlockReason(resource models.GroundResource, cand candidateTask, proposed map[string][]models.ResourceBooking, localAssign map[uint][]assignmentWindow) *types.ConflictItem {
	switch resource.AvailabilityStatus {
	case constants.ResourceOffline:
		return &types.ConflictItem{
			Code:    constants.ConflictOffline,
			Message: fmt.Sprintf(constants.MsgResourceOffline, resource.ResourceCode),
		}
	case constants.ResourceMaintenance:
		return &types.ConflictItem{
			Code:    constants.ConflictMaintenance,
			Message: fmt.Sprintf("资源 %s 标记为维护中", resource.ResourceCode),
		}
	}
	if resource.MaintenanceDueAt != nil && resource.MaintenanceEnd != nil &&
		cand.start.Before(*resource.MaintenanceEnd) && resource.MaintenanceDueAt.Before(cand.end) {
		return &types.ConflictItem{
			Code:    constants.ConflictMaintenance,
			Message: fmt.Sprintf(constants.MsgResourceMaint, resource.ResourceCode,
				resource.MaintenanceDueAt.Format("01-02 15:04"), cand.start.Format("15:04"), cand.end.Format("15:04")),
		}
	}
	// Existing persisted bookings (other flights) — frozen and active.
	existing, err := s.Booking.Overlapping(resource.ID, cand.start, cand.end, 0)
	if err == nil && len(existing) > 0 {
		first := existing[0]
		return &types.ConflictItem{
			Code:    constants.ConflictOverlap,
			Message: fmt.Sprintf(constants.MsgBookingConflict, resource.ResourceCode,
				cand.start.Format("15:04"), cand.end.Format("15:04"), first.ID),
			BookingID: first.ID,
		}
	}
	// In-flight batch double-booking.
	for _, win := range localAssign[resource.ID] {
		if cand.start.Before(win.end) && win.start.Before(cand.end) {
			return &types.ConflictItem{
				Code:    constants.ConflictOverlap,
				Message: fmt.Sprintf("资源 %s 时段与同批 %s 任务重叠", resource.ResourceCode, constants.GroundTaskTypeText[win.taskType]),
			}
		}
	}
	return nil
}

// chooseResource repeats selection during persistence (batch already valid).
func (s *PlanService) chooseResource(cand candidateTask, proposed map[string][]models.ResourceBooking, localAssign map[uint][]assignmentWindow) uint {
	resources, err := s.Resource.ListByType(constants.TaskResourceType[cand.taskType])
	if err != nil {
		return 0
	}
	id, _ := s.pickResource(resources, cand, proposed, localAssign)
	return id
}

func (s *PlanService) audit(actor, role, message, targetType, targetID string) {
	_ = s.Audit.Create(&models.AuditLog{
		Actor:      actor,
		Role:       role,
		Action:     message,
		TargetType: targetType,
		TargetID:   targetID,
	})
}
