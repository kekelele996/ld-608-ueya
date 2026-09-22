package constructors

import (
	"groundTurn/src/constants"
	"groundTurn/src/models"
	"groundTurn/src/utils"
)

// ResourceView is the GroundResource list/detail response DTO.
type ResourceView struct {
	ID                uint   `json:"id"`
	ResourceCode      string `json:"resource_code"`
	ResourceType      string `json:"resource_type"`
	Location          string `json:"location"`
	AvailabilityStatus string `json:"availability_status"`
	StatusText        string `json:"status_text"`
	MaintenanceDueAt  string `json:"maintenance_due_at"`
	MaintenanceEnd    string `json:"maintenance_end"`
	OwnerTeam         string `json:"owner_team"`
}

func NewResourceView(r models.GroundResource) ResourceView {
	return ResourceView{
		ID:                 r.ID,
		ResourceCode:       r.ResourceCode,
		ResourceType:       r.ResourceType,
		Location:           r.Location,
		AvailabilityStatus: r.AvailabilityStatus,
		StatusText:         constants.ResourceStatusText[r.AvailabilityStatus],
		MaintenanceDueAt:   utils.FormatTimePtr(r.MaintenanceDueAt),
		MaintenanceEnd:     utils.FormatTimePtr(r.MaintenanceEnd),
		OwnerTeam:          r.OwnerTeam,
	}
}

func NewResourceViews(rows []models.GroundResource) []ResourceView {
	out := make([]ResourceView, 0, len(rows))
	for _, row := range rows {
		out = append(out, NewResourceView(row))
	}
	return out
}

// NewResourceFromRegister builds a GroundResource for the register form.
func NewResourceFromRegister(code, resourceType, location, status, team string) models.GroundResource {
	return models.GroundResource{
		ResourceCode:       code,
		ResourceType:       resourceType,
		Location:           location,
		AvailabilityStatus: status,
		OwnerTeam:          team,
	}
}
