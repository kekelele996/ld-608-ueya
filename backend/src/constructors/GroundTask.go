package constructors

import (
	"time"

	"groundTurn/src/constants"
	"groundTurn/src/models"
)

// NewGroundTask constructs a planned task from a SOP blueprint projection.
func NewGroundTask(turnaroundID uint, kind constants.GroundTaskType, start time.Time) *models.GroundTask {
	bp := constants.BlueprintFor(kind)
	team := constants.TeamForTask(kind)
	end := start.Add(bp.Duration)
	return &models.GroundTask{
		TurnaroundID: turnaroundID,
		TaskType:     string(kind),
		TeamID:       team.ID,
		PlannedStart: start,
		PlannedEnd:   end,
		Deadline:     end.Add(bp.DeadlineGap),
		Status:       string(constants.TaskPlanned),
	}
}
