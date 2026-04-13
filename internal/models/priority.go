package models

import "fmt"

// Priority represents the desire level of a gift (1=lowest, 5=highest)
type Priority int

const (
	// PriorityLow is the lowest desire level
	PriorityLow Priority = 1
	// PriorityMedium is medium desire level
	PriorityMedium Priority = 2
	// PriorityHigh is high desire level
	PriorityHigh Priority = 3
	// PriorityUrgent is urgent desire level
	PriorityUrgent Priority = 4
	// PriorityMust is the highest desire level
	PriorityMust Priority = 5
)

// Valid checks if priority is within allowed range
func (p Priority) Valid() bool {
	return p >= PriorityLow && p <= PriorityMust
}

// String returns string representation
func (p Priority) String() string {
	switch p {
	case PriorityLow:
		return "low"
	case PriorityMedium:
		return "medium"
	case PriorityHigh:
		return "high"
	case PriorityUrgent:
		return "urgent"
	case PriorityMust:
		return "must"
	default:
		return fmt.Sprintf("Priority(%d)", p)
	}
}
