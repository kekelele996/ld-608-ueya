package constructors

import (
	"groundTurn/src/constants"
	"groundTurn/src/models"
)

// NewGroundResource builds a resource ledger row during seeding/admin create.
func NewGroundResource(code, resourceType, location, teamID string, status constants.ResourceStatus) *models.GroundResource {
	return &models.GroundResource{
		ResourceCode:       code,
		ResourceType:       resourceType,
		Location:           location,
		AvailabilityStatus: string(status),
		OwnerTeam:          teamID,
	}
}
