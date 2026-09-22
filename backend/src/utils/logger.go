package utils

import (
	"fmt"
	"log"

	"groundTurn/src/constants"
)

// LogOperation renders a centralized log template. Every write action calls
// this so field changes force edits to constants/logTemplates.go too.
func LogOperation(action string, args ...any) {
	template := constants.LogTemplates[action]
	if template == "" {
		template = action
	}
	log.Printf("audit|%s|%s", action, fmt.Sprintf(template, args...))
}
