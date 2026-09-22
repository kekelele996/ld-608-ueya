package constructors

import (
	"time"

	"groundTurn/src/constants"
	"groundTurn/src/models"
	"groundTurn/src/types"
	"groundTurn/src/utils"
)

// NewTaskView builds the GroundTask response DTO. Display components read
// task_type_text / status_text from here instead of formatting locally.
func NewTaskView(t models.GroundTask) types.TaskView {
	return types.TaskView{
		ID:           t.ID,
		TurnaroundID: t.TurnaroundID,
		TaskType:     t.TaskType,
		TaskTypeText: constants.GroundTaskTypeText[t.TaskType],
		TeamID:       t.TeamID,
		PlannedStart: utils.FormatTime(t.PlannedStart),
		Deadline:     utils.FormatTime(t.Deadline),
		ActualFinish: utils.FormatTimePtr(t.ActualFinish),
		Status:       t.Status,
		StatusText:   constants.GroundTaskStatusText[t.Status],
		BlockerNote:  t.BlockerNote,
		SignedAt:     utils.FormatTimePtr(t.SignedAt),
	}
}

func NewTaskViews(rows []models.GroundTask) []types.TaskView {
	out := make([]types.TaskView, 0, len(rows))
	for _, row := range rows {
		out = append(out, NewTaskView(row))
	}
	return out
}

// NewBookingView assembles a booking with joined resource display fields.
func NewBookingView(b models.ResourceBooking, resource *models.GroundResource) types.BookingView {
	view := types.BookingView{
		ID:             b.ID,
		ResourceID:     b.ResourceID,
		TurnaroundID:   b.TurnaroundID,
		TaskID:         b.TaskID,
		StartTime:      utils.FormatTime(b.StartTime),
		EndTime:        utils.FormatTime(b.EndTime),
		BookingStatus:  b.BookingStatus,
		StatusText:     constants.BookingStatusText[b.BookingStatus],
		ConflictCode:   b.ConflictReason,
		ConflictReason: constants.ConflictReasonText[b.ConflictReason],
	}
	if resource != nil {
		view.ResourceCode = resource.ResourceCode
		view.ResourceType = resource.ResourceType
	}
	return view
}

// NewPlannedTask constructs one candidate GroundTask from the rule table.
// Used by the plan builder; never by request JSON binding.
func NewPlannedTask(turnaroundID uint, taskType, teamID string, start, deadline time.Time) models.GroundTask {
	return models.GroundTask{
		TurnaroundID: turnaroundID,
		TaskType:     taskType,
		TeamID:       teamID,
		PlannedStart: start,
		Deadline:     deadline,
		Status:       constants.TaskPlanned,
	}
}

// NewPlannedBooking constructs one candidate ResourceBooking.
func NewPlannedBooking(resourceID, turnaroundID, taskID uint, start, end time.Time) models.ResourceBooking {
	return models.ResourceBooking{
		ResourceID:    resourceID,
		TurnaroundID:  turnaroundID,
		TaskID:        taskID,
		StartTime:     start,
		EndTime:       end,
		BookingStatus: constants.BookingConfirmed,
	}
}
