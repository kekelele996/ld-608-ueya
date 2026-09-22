package constructors

import (
	"groundTurn/src/constants"
	"groundTurn/src/models"
	"groundTurn/src/utils"
)

// DelayView is the DelayEvent response DTO.
type DelayView struct {
	ID                 uint   `json:"id"`
	TurnaroundID       uint   `json:"turnaround_id"`
	FlightNo           string `json:"flight_no"`
	DelayType          string `json:"delay_type"`
	DelayTypeText      string `json:"delay_type_text"`
	Minutes            int    `json:"minutes"`
	RootCause          string `json:"root_cause"`
	ResponsibilityTeam string `json:"responsibility_team"`
	ResolvedAt         string `json:"resolved_at"`
	CreatedAt          string `json:"created_at"`
}

func NewDelayView(d models.DelayEvent, flightNo string) DelayView {
	return DelayView{
		ID:                 d.ID,
		TurnaroundID:       d.TurnaroundID,
		FlightNo:           flightNo,
		DelayType:          d.DelayType,
		DelayTypeText:      constants.DelayTypeText[d.DelayType],
		Minutes:            d.Minutes,
		RootCause:          d.RootCause,
		ResponsibilityTeam: d.ResponsibilityTeam,
		ResolvedAt:         utils.FormatTimePtr(d.ResolvedAt),
		CreatedAt:          utils.FormatTime(d.CreatedAt),
	}
}

func NewDelayViews(rows []models.DelayEvent, flightMap map[uint]string) []DelayView {
	out := make([]DelayView, 0, len(rows))
	for _, row := range rows {
		out = append(out, NewDelayView(row, flightMap[row.TurnaroundID]))
	}
	return out
}

// AuditView is the operation-log response DTO.
type AuditView struct {
	ID         uint   `json:"id"`
	Actor      string `json:"actor"`
	Role       string `json:"role"`
	Action     string `json:"action"`
	TargetType string `json:"target_type"`
	TargetID   string `json:"target_id"`
	CreatedAt  string `json:"created_at"`
}

func NewAuditView(a models.AuditLog) AuditView {
	return AuditView{
		ID:         a.ID,
		Actor:      a.Actor,
		Role:       a.Role,
		Action:     a.Action,
		TargetType: a.TargetType,
		TargetID:   a.TargetID,
		CreatedAt:  utils.FormatTime(a.CreatedAt),
	}
}
