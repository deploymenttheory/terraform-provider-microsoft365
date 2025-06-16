package crud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

// Mock structs for testing
type mockIdentifiable struct {
	id types.String
}

func (m mockIdentifiable) GetID() types.String {
	return m.id
}

type mockAssignment struct {
	mockIdentifiable
	target interface{}
}

func (m mockAssignment) GetTarget() interface{} {
	return m.target
}

func TestExistsInSlice(t *testing.T) {
	tests := []struct {
		name     string
		item     int
		slice    []int
		compare  CompareFunc
		expected bool
	}{
		{
			name:     "Item exists",
			item:     2,
			slice:    []int{1, 2, 3},
			compare:  func(a, b interface{}) bool { return a.(int) == b.(int) },
			expected: true,
		},
		{
			name:     "Item does not exist",
			item:     4,
			slice:    []int{1, 2, 3},
			compare:  func(a, b interface{}) bool { return a.(int) == b.(int) },
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ExistsInSlice(tt.item, tt.slice, tt.compare)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestDefaultAssignmentsEqual(t *testing.T) {
	tests := []struct {
		name     string
		a, b     Assignment
		expected bool
	}{
		{
			name:     "Equal IDs",
			a:        mockAssignment{mockIdentifiable{types.StringValue("1")}, "target"},
			b:        mockAssignment{mockIdentifiable{types.StringValue("1")}, "different"},
			expected: true,
		},
		{
			name:     "Different IDs",
			a:        mockAssignment{mockIdentifiable{types.StringValue("1")}, "target"},
			b:        mockAssignment{mockIdentifiable{types.StringValue("2")}, "target"},
			expected: false,
		},
		{
			name:     "Null IDs, equal targets",
			a:        mockAssignment{mockIdentifiable{types.StringNull()}, "target"},
			b:        mockAssignment{mockIdentifiable{types.StringNull()}, "target"},
			expected: true,
		},
		{
			name:     "Null IDs, different targets",
			a:        mockAssignment{mockIdentifiable{types.StringNull()}, "target1"},
			b:        mockAssignment{mockIdentifiable{types.StringNull()}, "target2"},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := DefaultAssignmentsEqual(tt.a, tt.b)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestAssignmentExistsInSlice(t *testing.T) {
	assignment1 := mockAssignment{mockIdentifiable{types.StringValue("1")}, "target1"}
	assignment2 := mockAssignment{mockIdentifiable{types.StringValue("2")}, "target2"}
	assignment3 := mockAssignment{mockIdentifiable{types.StringValue("3")}, "target3"}

	assignments := []mockAssignment{assignment1, assignment2}

	tests := []struct {
		name     string
		item     mockAssignment
		slice    []mockAssignment
		compare  CompareFunc
		expected bool
	}{
		{
			name:     "Assignment exists",
			item:     assignment1,
			slice:    assignments,
			compare:  nil,
			expected: true,
		},
		{
			name:     "Assignment does not exist",
			item:     assignment3,
			slice:    assignments,
			compare:  nil,
			expected: false,
		},
		{
			name:     "Custom compare function",
			item:     mockAssignment{mockIdentifiable{types.StringValue("1")}, "different"},
			slice:    assignments,
			compare:  func(a, b interface{}) bool { return a.(Assignment).GetID() == b.(Assignment).GetID() },
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := AssignmentExistsInSlice(tt.item, tt.slice, tt.compare)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIdentifiableExistsInSlice(t *testing.T) {
	item1 := mockIdentifiable{types.StringValue("1")}
	item2 := mockIdentifiable{types.StringValue("2")}
	item3 := mockIdentifiable{types.StringValue("3")}

	slice := []mockIdentifiable{item1, item2}

	tests := []struct {
		name     string
		item     mockIdentifiable
		slice    []mockIdentifiable
		expected bool
	}{
		{
			name:     "Item exists",
			item:     item1,
			slice:    slice,
			expected: true,
		},
		{
			name:     "Item does not exist",
			item:     item3,
			slice:    slice,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IdentifiableExistsInSlice(tt.item, tt.slice)
			assert.Equal(t, tt.expected, result)
		})
	}
}
