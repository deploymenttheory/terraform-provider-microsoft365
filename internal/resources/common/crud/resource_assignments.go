package crud

import (
	"reflect"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Identifiable is an interface for types that have an ID
type Identifiable interface {
	GetID() types.String
}

// Assignment is a generic interface for types that represent assignments
type Assignment interface {
	Identifiable
	GetTarget() interface{}
}

// CompareFunc is a type for custom comparison functions
type CompareFunc func(a, b interface{}) bool

// ExistsInSlice checks if a given item exists in a slice based on a comparison function
func ExistsInSlice[T any](item T, slice []T, compareFunc CompareFunc) bool {
	for _, sliceItem := range slice {
		if compareFunc(item, sliceItem) {
			return true
		}
	}
	return false
}

// DefaultAssignmentsEqual compares two Assignment instances for equality
func DefaultAssignmentsEqual(a, b Assignment) bool {
	// Compare IDs if they are not null or unknown
	if !a.GetID().IsNull() && !a.GetID().IsUnknown() && !b.GetID().IsNull() && !b.GetID().IsUnknown() {
		return a.GetID() == b.GetID()
	}

	// If IDs are not available or not comparable, compare targets
	return reflect.DeepEqual(a.GetTarget(), b.GetTarget())
}

// AssignmentExistsInSlice checks if a given assignment exists in a slice of assignments
func AssignmentExistsInSlice[T Assignment](assignment T, assignments []T, compareFunc CompareFunc) bool {
	if compareFunc == nil {
		compareFunc = func(a, b interface{}) bool {
			return DefaultAssignmentsEqual(a.(Assignment), b.(Assignment))
		}
	}
	return ExistsInSlice(assignment, assignments, compareFunc)
}

// IdentifiableExistsInSlice checks if a given identifiable item exists in a slice based on ID
func IdentifiableExistsInSlice[T Identifiable](item T, slice []T) bool {
	return ExistsInSlice(item, slice, func(a, b interface{}) bool {
		return a.(Identifiable).GetID() == b.(Identifiable).GetID()
	})
}
