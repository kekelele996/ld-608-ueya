package services

import (
	"groundTurn/src/constants"
	"groundTurn/src/models"
	"groundTurn/src/types"
)

type typesConflictItem = types.ConflictItem

func findTask(tasks []models.GroundTask, taskType string) *models.GroundTask {
	for i := range tasks {
		if tasks[i].TaskType == taskType {
			return &tasks[i]
		}
	}
	return nil
}

func (f *fixture) codeByResourceID(id uint) string {
	for code, rid := range f.resources {
		if rid == id {
			return code
		}
	}
	return ""
}

// silence unused constants in test builds if an assertion drops one.
var _ = constants.ConflictOffline
