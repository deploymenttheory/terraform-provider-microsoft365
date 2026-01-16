package utilityDeploymentScheduler

import (
	"fmt"
	"strings"
)

// StatusMessageBuilder constructs structured, human-readable status messages
type StatusMessageBuilder struct {
	gateOpen         bool
	manualOverride   bool
	conditions       []conditionStatus
	dependencyStatus *dependencyConditionStatus
}

type conditionStatus struct {
	name   string
	met    bool
	detail string
}

type dependencyConditionStatus struct {
	met    bool
	detail string
}

// newStatusMessageBuilder creates a new status message builder
func newStatusMessageBuilder(gateOpen bool) *StatusMessageBuilder {
	return &StatusMessageBuilder{
		gateOpen:   gateOpen,
		conditions: []conditionStatus{},
	}
}

// setManualOverride marks the gate as manually overridden
func (smb *StatusMessageBuilder) setManualOverride() {
	smb.manualOverride = true
}

// addTimeCondition adds time condition status
func (smb *StatusMessageBuilder) addTimeCondition(met bool, detail string) {
	smb.conditions = append(smb.conditions, conditionStatus{
		name:   "Time",
		met:    met,
		detail: detail,
	})
}

// addInclusionWindow adds inclusion window status
func (smb *StatusMessageBuilder) addInclusionWindow(met bool, detail string) {
	smb.conditions = append(smb.conditions, conditionStatus{
		name:   "Inclusion Window",
		met:    met,
		detail: detail,
	})
}

// addExclusionWindow adds exclusion window status
func (smb *StatusMessageBuilder) addExclusionWindow(active bool, detail string) {
	smb.conditions = append(smb.conditions, conditionStatus{
		name:   "Exclusion Window",
		met:    !active, // Invert because active = bad
		detail: detail,
	})
}

// setDependency sets dependency gate status
func (smb *StatusMessageBuilder) setDependency(met bool, detail string) {
	smb.dependencyStatus = &dependencyConditionStatus{
		met:    met,
		detail: detail,
	}
}

// build constructs the final status message
func (smb *StatusMessageBuilder) build() string {
	if smb.manualOverride {
		return "[GATE OPEN] Manual override enabled - all conditions bypassed"
	}

	var parts []string

	// Header
	if smb.gateOpen {
		parts = append(parts, "[GATE OPEN]")
	} else {
		parts = append(parts, "[GATE CLOSED]")
	}

	// Conditions
	if len(smb.conditions) > 0 {
		var conditionParts []string
		for _, cond := range smb.conditions {
			status := "PASS"
			if !cond.met {
				status = "FAIL"
			}
			conditionParts = append(conditionParts, fmt.Sprintf("%s - %s: %s", status, cond.name, cond.detail))
		}
		parts = append(parts, strings.Join(conditionParts, " | "))
	}

	// Dependency
	if smb.dependencyStatus != nil {
		status := "PASS"
		if !smb.dependencyStatus.met {
			status = "FAIL"
		}
		parts = append(parts, fmt.Sprintf("%s - Dependency: %s", status, smb.dependencyStatus.detail))
	}

	return strings.Join(parts, " | ")
}
