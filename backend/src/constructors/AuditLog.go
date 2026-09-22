package constructors

import (
	"time"

	"groundTurn/src/models"
)

// NewAuditLog builds the operation log row written by services/middleware.
func NewAuditLog(actor, role, action, targetType, targetID, detail string) *models.AuditLog {
	return &models.AuditLog{
		Actor:      actor,
		Role:       role,
		Action:     action,
		TargetType: targetType,
		TargetID:   targetID,
		Detail:     detail,
		CreatedAt:  time.Now(),
	}
}
