package constructors

import (
	"groundTurn/src/constants"
	"groundTurn/src/models"
	"groundTurn/src/types"
)

// NewFlightTurnaround builds a model from a create request (factory file
// deliberately separated from service/controller).
func NewFlightTurnaround(req types.CreateTurnaroundRequest) *models.FlightTurnaround {
	return &models.FlightTurnaround{
		FlightNo:         req.FlightNo,
		AircraftReg:      req.AircraftReg,
		StandNo:          req.StandNo,
		ArrivalTime:      req.ArrivalTime,
		DepartureTime:    req.DepartureTime,
		TurnaroundStatus: string(constants.StatusArriving),
	}
}
