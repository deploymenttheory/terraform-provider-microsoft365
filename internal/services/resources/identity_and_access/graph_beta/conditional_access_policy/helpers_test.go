package graphBetaConditionalAccessPolicy

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Mock enum type for testing
type MockEnum string

func (m MockEnum) String() string {
	return string(m)
}

const (
	MockEnumValue1 MockEnum = "value1"
	MockEnumValue2 MockEnum = "value2"
	MockEnumValue3 MockEnum = "value3"
)

func TestMapStringSliceToSetPreserveEmpty(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name     string
		input    []string
		wantNull bool
		wantLen  int
	}{
		{
			name:     "empty slice returns empty set",
			input:    []string{},
			wantNull: false,
			wantLen:  0,
		},
		{
			name:     "nil slice returns empty set",
			input:    nil,
			wantNull: false,
			wantLen:  0,
		},
		{
			name:     "single value",
			input:    []string{"value1"},
			wantNull: false,
			wantLen:  1,
		},
		{
			name:     "multiple values",
			input:    []string{"value1", "value2", "value3"},
			wantNull: false,
			wantLen:  3,
		},
		{
			name:     "values with special characters",
			input:    []string{"value-1", "value_2", "value.3"},
			wantNull: false,
			wantLen:  3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := mapStringSliceToSetPreserveEmpty(ctx, tt.input)

			if result.IsNull() != tt.wantNull {
				t.Errorf("IsNull() = %v, want %v", result.IsNull(), tt.wantNull)
			}

			if !result.IsNull() {
				elements := result.Elements()
				if len(elements) != tt.wantLen {
					t.Errorf("len(elements) = %d, want %d", len(elements), tt.wantLen)
				}

				// Verify values match input
				if tt.wantLen > 0 {
					for i, elem := range elements {
						strVal, ok := elem.(types.String)
						if !ok {
							t.Errorf("element %d is not types.String", i)
							continue
						}
						if strVal.ValueString() != tt.input[i] {
							t.Errorf("element %d = %q, want %q", i, strVal.ValueString(), tt.input[i])
						}
					}
				}
			}
		})
	}
}

func TestMapEnumCollectionToSet(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name      string
		input     []MockEnum
		fieldName string
		wantNull  bool
		wantLen   int
		wantVals  []string
	}{
		{
			name:      "empty slice returns empty set",
			input:     []MockEnum{},
			fieldName: "testField",
			wantNull:  false,
			wantLen:   0,
			wantVals:  []string{},
		},
		{
			name:      "nil slice returns empty set",
			input:     nil,
			fieldName: "testField",
			wantNull:  false,
			wantLen:   0,
			wantVals:  []string{},
		},
		{
			name:      "single enum value",
			input:     []MockEnum{MockEnumValue1},
			fieldName: "testField",
			wantNull:  false,
			wantLen:   1,
			wantVals:  []string{"value1"},
		},
		{
			name:      "multiple enum values",
			input:     []MockEnum{MockEnumValue1, MockEnumValue2, MockEnumValue3},
			fieldName: "testField",
			wantNull:  false,
			wantLen:   3,
			wantVals:  []string{"value1", "value2", "value3"},
		},
		{
			name:      "enum values in specific order",
			input:     []MockEnum{MockEnumValue3, MockEnumValue1, MockEnumValue2},
			fieldName: "testField",
			wantNull:  false,
			wantLen:   3,
			wantVals:  []string{"value3", "value1", "value2"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := mapEnumCollectionToSet(ctx, tt.input, tt.fieldName)

			if result.IsNull() != tt.wantNull {
				t.Errorf("IsNull() = %v, want %v", result.IsNull(), tt.wantNull)
			}

			if !result.IsNull() {
				elements := result.Elements()
				if len(elements) != tt.wantLen {
					t.Errorf("len(elements) = %d, want %d", len(elements), tt.wantLen)
				}

				// Verify values match expected
				for i, elem := range elements {
					strVal, ok := elem.(types.String)
					if !ok {
						t.Errorf("element %d is not types.String", i)
						continue
					}
					if i < len(tt.wantVals) && strVal.ValueString() != tt.wantVals[i] {
						t.Errorf("element %d = %q, want %q", i, strVal.ValueString(), tt.wantVals[i])
					}
				}
			}
		})
	}
}

func TestMapEnumCollectionToSetNullIfEmpty(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name      string
		input     []MockEnum
		fieldName string
		wantNull  bool
		wantLen   int
		wantVals  []string
	}{
		{
			name:      "empty slice returns null",
			input:     []MockEnum{},
			fieldName: "testField",
			wantNull:  true,
			wantLen:   0,
		},
		{
			name:      "nil slice returns null",
			input:     nil,
			fieldName: "testField",
			wantNull:  true,
			wantLen:   0,
		},
		{
			name:      "single enum value",
			input:     []MockEnum{MockEnumValue1},
			fieldName: "testField",
			wantNull:  false,
			wantLen:   1,
			wantVals:  []string{"value1"},
		},
		{
			name:      "multiple enum values",
			input:     []MockEnum{MockEnumValue1, MockEnumValue2, MockEnumValue3},
			fieldName: "testField",
			wantNull:  false,
			wantLen:   3,
			wantVals:  []string{"value1", "value2", "value3"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := mapEnumCollectionToSetNullIfEmpty(ctx, tt.input, tt.fieldName)

			if result.IsNull() != tt.wantNull {
				t.Errorf("IsNull() = %v, want %v", result.IsNull(), tt.wantNull)
			}

			if !result.IsNull() {
				elements := result.Elements()
				if len(elements) != tt.wantLen {
					t.Errorf("len(elements) = %d, want %d", len(elements), tt.wantLen)
				}

				// Verify values match expected
				for i, elem := range elements {
					strVal, ok := elem.(types.String)
					if !ok {
						t.Errorf("element %d is not types.String", i)
						continue
					}
					if strVal.ValueString() != tt.wantVals[i] {
						t.Errorf("element %d = %q, want %q", i, strVal.ValueString(), tt.wantVals[i])
					}
				}
			}
		})
	}
}

func TestMapStringSliceToSetNullIfEmpty(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name     string
		input    []string
		wantNull bool
		wantLen  int
	}{
		{
			name:     "empty slice returns null",
			input:    []string{},
			wantNull: true,
			wantLen:  0,
		},
		{
			name:     "nil slice returns null",
			input:    nil,
			wantNull: true,
			wantLen:  0,
		},
		{
			name:     "single value",
			input:    []string{"value1"},
			wantNull: false,
			wantLen:  1,
		},
		{
			name:     "multiple values",
			input:    []string{"value1", "value2", "value3"},
			wantNull: false,
			wantLen:  3,
		},
		{
			name:     "values with GUIDs",
			input:    []string{"12345678-1234-1234-1234-123456789012", "87654321-4321-4321-4321-210987654321"},
			wantNull: false,
			wantLen:  2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := mapStringSliceToSetNullIfEmpty(ctx, tt.input)

			if result.IsNull() != tt.wantNull {
				t.Errorf("IsNull() = %v, want %v", result.IsNull(), tt.wantNull)
			}

			if !result.IsNull() {
				elements := result.Elements()
				if len(elements) != tt.wantLen {
					t.Errorf("len(elements) = %d, want %d", len(elements), tt.wantLen)
				}

				// Verify values match input
				if tt.wantLen > 0 {
					for i, elem := range elements {
						strVal, ok := elem.(types.String)
						if !ok {
							t.Errorf("element %d is not types.String", i)
							continue
						}
						if strVal.ValueString() != tt.input[i] {
							t.Errorf("element %d = %q, want %q", i, strVal.ValueString(), tt.input[i])
						}
					}
				}
			}
		})
	}
}

// TestMapStringSliceToSetPreserveEmpty_VerifyEmptyVsNull ensures the function correctly
// distinguishes between empty sets and null for state consistency
func TestMapStringSliceToSetPreserveEmpty_VerifyEmptyVsNull(t *testing.T) {
	ctx := context.Background()

	emptySlice := []string{}
	result := mapStringSliceToSetPreserveEmpty(ctx, emptySlice)

	if result.IsNull() {
		t.Error("Expected empty set, got null - this breaks state consistency for required empty sets")
	}

	if result.IsUnknown() {
		t.Error("Expected empty set, got unknown")
	}

	elements := result.Elements()
	if len(elements) != 0 {
		t.Errorf("Expected empty set with 0 elements, got %d elements", len(elements))
	}

	// Verify it's a valid empty set, not null
	expectedEmptySet := types.SetValueMust(types.StringType, []attr.Value{})
	if result.Equal(expectedEmptySet) != true {
		t.Error("Empty set does not match expected empty set value")
	}
}

// TestMapStringSliceToSetNullIfEmpty_VerifyEmptyVsNull ensures the function correctly
// returns null for empty slices (for optional fields)
func TestMapStringSliceToSetNullIfEmpty_VerifyEmptyVsNull(t *testing.T) {
	ctx := context.Background()

	emptySlice := []string{}
	result := mapStringSliceToSetNullIfEmpty(ctx, emptySlice)

	if !result.IsNull() {
		t.Error("Expected null, got non-null - this breaks state consistency for optional empty sets")
	}

	if result.IsUnknown() {
		t.Error("Expected null, got unknown")
	}

	// Verify it's null, not an empty set
	expectedNull := types.SetNull(types.StringType)
	if result.Equal(expectedNull) != true {
		t.Error("Null set does not match expected null set value")
	}
}

// TestMapEnumCollectionToSet_VerifyEmptyVsNull ensures the function correctly
// preserves empty sets (for required enum fields)
func TestMapEnumCollectionToSet_VerifyEmptyVsNull(t *testing.T) {
	ctx := context.Background()

	emptySlice := []MockEnum{}
	result := mapEnumCollectionToSet(ctx, emptySlice, "testField")

	if result.IsNull() {
		t.Error("Expected empty set, got null - this breaks state consistency for required empty enum sets")
	}

	if result.IsUnknown() {
		t.Error("Expected empty set, got unknown")
	}

	elements := result.Elements()
	if len(elements) != 0 {
		t.Errorf("Expected empty set with 0 elements, got %d elements", len(elements))
	}

	// Verify it's a valid empty set, not null
	expectedEmptySet := types.SetValueMust(types.StringType, []attr.Value{})
	if result.Equal(expectedEmptySet) != true {
		t.Error("Empty enum set does not match expected empty set value")
	}
}

// TestMapEnumCollectionToSetNullIfEmpty_VerifyEmptyVsNull ensures the function correctly
// returns null for empty enum slices
func TestMapEnumCollectionToSetNullIfEmpty_VerifyEmptyVsNull(t *testing.T) {
	ctx := context.Background()

	emptySlice := []MockEnum{}
	result := mapEnumCollectionToSetNullIfEmpty(ctx, emptySlice, "testField")

	if !result.IsNull() {
		t.Error("Expected null, got non-null - this breaks state consistency for optional enum sets")
	}

	if result.IsUnknown() {
		t.Error("Expected null, got unknown")
	}

	// Verify it's null, not an empty set
	expectedNull := types.SetNull(types.StringType)
	if result.Equal(expectedNull) != true {
		t.Error("Null set does not match expected null set value")
	}
}

// TestEnumCollectionBehaviorComparison validates the critical difference between
// mapEnumCollectionToSet and mapEnumCollectionToSetNullIfEmpty
func TestEnumCollectionBehaviorComparison(t *testing.T) {
	ctx := context.Background()

	t.Run("empty slice behavior difference", func(t *testing.T) {
		emptySlice := []MockEnum{}

		// mapEnumCollectionToSet should preserve empty set
		preserveResult := mapEnumCollectionToSet(ctx, emptySlice, "testField")
		if preserveResult.IsNull() {
			t.Error("mapEnumCollectionToSet returned null for empty slice, expected empty set")
		}

		// mapEnumCollectionToSetNullIfEmpty should return null
		nullResult := mapEnumCollectionToSetNullIfEmpty(ctx, emptySlice, "testField")
		if !nullResult.IsNull() {
			t.Error("mapEnumCollectionToSetNullIfEmpty returned non-null for empty slice, expected null")
		}

		// They should not be equal
		if preserveResult.Equal(nullResult) {
			t.Error("Empty set and null should not be equal")
		}
	})

	t.Run("populated slice behavior same", func(t *testing.T) {
		populatedSlice := []MockEnum{MockEnumValue1, MockEnumValue2}

		// Both should return non-null sets with same values
		preserveResult := mapEnumCollectionToSet(ctx, populatedSlice, "testField")
		nullResult := mapEnumCollectionToSetNullIfEmpty(ctx, populatedSlice, "testField")

		if preserveResult.IsNull() {
			t.Error("mapEnumCollectionToSet returned null for populated slice")
		}

		if nullResult.IsNull() {
			t.Error("mapEnumCollectionToSetNullIfEmpty returned null for populated slice")
		}

		// Both should have same values
		if len(preserveResult.Elements()) != len(nullResult.Elements()) {
			t.Errorf("Different element counts: preserve=%d, null=%d",
				len(preserveResult.Elements()), len(nullResult.Elements()))
		}

		// They should be equal for populated slices
		if !preserveResult.Equal(nullResult) {
			t.Error("Both functions should return equal sets for populated slices")
		}
	})
}
