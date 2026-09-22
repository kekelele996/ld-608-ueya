package constructors

import (
	"groundTurn/src/models"
	"groundTurn/src/types"
)

// NewDelayEvent builds an unresolved delay event from a register request.
func NewDelayEvent(turnaroundID uint, req types.RegisterDelayRequest) *models.DelayEvent {
	return &models.DelayEvent{
		TurnaroundID:       turnaroundID,
		DelayType:          req.DelayType,
		Minutes:            req.Minutes,
		RootCause:          req.RootCause,
		ResponsibilityTeam: req.ResponsibilityTeam,
	}
}
