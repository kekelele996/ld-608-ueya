package constants

import "time"

// TaskBlueprint defines how plan generation derives a GroundTask and its
// resource appointment window from the on-stand time.
type TaskBlueprint struct {
	Kind        GroundTaskType
	StartOffset time.Duration // relative to arrival_time (on-stand)
	Duration    time.Duration
	DeadlineGap time.Duration // allowed slack after planned_end before deadline
}

func minutes(n int64) time.Duration { return time.Duration(n) * time.Minute }

// TaskBlueprints is ordered like the turnaround SOP timeline.
var TaskBlueprints = []TaskBlueprint{
	{TaskBaggage, minutes(0), minutes(30), minutes(5)},
	{TaskCleaning, minutes(5), minutes(35), minutes(5)},
	{TaskCatering, minutes(10), minutes(30), minutes(5)},
	{TaskWaterService, minutes(10), minutes(20), minutes(5)},
	{TaskRefuel, minutes(15), minutes(30), minutes(5)},
	{TaskPushback, minutes(55), minutes(15), minutes(0)},
}

func BlueprintFor(kind GroundTaskType) TaskBlueprint {
	for _, blueprint := range TaskBlueprints {
		if blueprint.Kind == kind {
			return blueprint
		}
	}
	return TaskBlueprints[0]
}
