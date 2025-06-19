package state

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

func TestBuildObjectSetFromSlice(t *testing.T) {
	ctx := context.Background()

	// Define test attribute types for objects
	attrTypes := map[string]attr.Type{
		"name":   types.StringType,
		"age":    types.Int32Type,
		"active": types.BoolType,
	}

	t.Run("Empty slice returns empty set", func(t *testing.T) {
		extract := func(i int) map[string]attr.Value {
			return map[string]attr.Value{}
		}

		result := BuildObjectSetFromSlice(ctx, attrTypes, extract, 0)

		assert.False(t, result.IsNull(), "Should not return null set for empty slice")
		assert.False(t, result.IsUnknown(), "Should not return unknown set for empty slice")

		elements := result.Elements()
		assert.Equal(t, 0, len(elements), "Should return empty set for zero length")
	})

	t.Run("Single element slice", func(t *testing.T) {
		extract := func(i int) map[string]attr.Value {
			return map[string]attr.Value{
				"name":   types.StringValue("John"),
				"age":    types.Int32Value(30),
				"active": types.BoolValue(true),
			}
		}

		result := BuildObjectSetFromSlice(ctx, attrTypes, extract, 1)

		assert.False(t, result.IsNull(), "Should not return null set")
		assert.False(t, result.IsUnknown(), "Should not return unknown set")

		elements := result.Elements()
		assert.Equal(t, 1, len(elements), "Should contain one element")

		// Verify the object structure
		obj := elements[0].(types.Object)
		objAttrs := obj.Attributes()
		assert.Equal(t, types.StringValue("John"), objAttrs["name"])
		assert.Equal(t, types.Int32Value(30), objAttrs["age"])
		assert.Equal(t, types.BoolValue(true), objAttrs["active"])
	})

	t.Run("Multiple elements slice", func(t *testing.T) {
		testData := []struct {
			name   string
			age    int32
			active bool
		}{
			{"Alice", 25, true},
			{"Bob", 35, false},
			{"Charlie", 40, true},
		}

		extract := func(i int) map[string]attr.Value {
			return map[string]attr.Value{
				"name":   types.StringValue(testData[i].name),
				"age":    types.Int32Value(int32(testData[i].age)),
				"active": types.BoolValue(testData[i].active),
			}
		}

		result := BuildObjectSetFromSlice(ctx, attrTypes, extract, len(testData))

		assert.False(t, result.IsNull(), "Should not return null set")
		assert.False(t, result.IsUnknown(), "Should not return unknown set")

		elements := result.Elements()
		assert.Equal(t, 3, len(elements), "Should contain three elements")

		for _, element := range elements {
			obj := element.(types.Object)
			objAttrs := obj.Attributes()

			// Note: Set elements might not be in the same order as input
			// So we need to find the matching element
			nameAttr := objAttrs["name"].(types.String)
			name := nameAttr.ValueString()

			var expectedData *struct {
				name   string
				age    int32
				active bool
			}

			for j := range testData {
				if testData[j].name == name {
					expectedData = &testData[j]
					break
				}
			}

			assert.NotNil(t, expectedData, "Should find matching test data for name: %s", name)
			assert.Equal(t, types.StringValue(expectedData.name), objAttrs["name"])
			assert.Equal(t, types.Int32Value(expectedData.age), objAttrs["age"])
			assert.Equal(t, types.BoolValue(expectedData.active), objAttrs["active"])
		}
	})

	t.Run("Extract function returns invalid values", func(t *testing.T) {
		// Extract function that returns invalid types
		extract := func(i int) map[string]attr.Value {
			return map[string]attr.Value{
				"name":   types.Int32Value(123),        // Wrong type - should be string
				"age":    types.StringValue("invalid"), // Wrong type - should be int64
				"active": types.StringValue("true"),    // Wrong type - should be bool
			}
		}

		result := BuildObjectSetFromSlice(ctx, attrTypes, extract, 1)

		// Should return empty set since object creation failed
		assert.False(t, result.IsNull(), "Should not return null set")
		elements := result.Elements()
		assert.Equal(t, 0, len(elements), "Should return empty set when object creation fails")
	})

	t.Run("Extract function returns missing required attributes", func(t *testing.T) {
		extract := func(i int) map[string]attr.Value {
			return map[string]attr.Value{
				"name": types.StringValue("John"),
				// Missing "age" and "active" attributes
			}
		}

		result := BuildObjectSetFromSlice(ctx, attrTypes, extract, 1)

		// Should return empty set since object creation failed
		assert.False(t, result.IsNull(), "Should not return null set")
		elements := result.Elements()
		assert.Equal(t, 0, len(elements), "Should return empty set when required attributes are missing")
	})

	t.Run("Mixed valid and invalid elements", func(t *testing.T) {
		extract := func(i int) map[string]attr.Value {
			if i == 1 {
				// Return invalid attributes for second element
				return map[string]attr.Value{
					"name":   types.Int32Value(123), // Wrong type
					"age":    types.StringValue("invalid"),
					"active": types.StringValue("true"),
				}
			}
			// Return valid attributes for other elements
			return map[string]attr.Value{
				"name":   types.StringValue("Valid"),
				"age":    types.Int32Value(25),
				"active": types.BoolValue(true),
			}
		}

		result := BuildObjectSetFromSlice(ctx, attrTypes, extract, 3)

		assert.False(t, result.IsNull(), "Should not return null set")
		elements := result.Elements()
		// Should contain only the valid elements (2 out of 3)
		assert.Equal(t, 2, len(elements), "Should contain only valid elements")
	})

	t.Run("Empty attribute types map", func(t *testing.T) {
		emptyAttrTypes := map[string]attr.Type{}

		extract := func(i int) map[string]attr.Value {
			return map[string]attr.Value{}
		}

		result := BuildObjectSetFromSlice(ctx, emptyAttrTypes, extract, 1)

		assert.False(t, result.IsNull(), "Should not return null set")
		elements := result.Elements()
		assert.Equal(t, 1, len(elements), "Should contain one element")

		obj := elements[0].(types.Object)
		objAttrs := obj.Attributes()
		assert.Equal(t, 0, len(objAttrs), "Object should have no attributes")
	})

	t.Run("Nil context handling", func(t *testing.T) {
		extract := func(i int) map[string]attr.Value {
			return map[string]attr.Value{
				"name":   types.StringValue("Test"),
				"age":    types.Int32Value(30),
				"active": types.BoolValue(true),
			}
		}

		// This should not panic even with nil context
		result := BuildObjectSetFromSlice(nil, attrTypes, extract, 1)

		assert.False(t, result.IsNull(), "Should not return null set")
		elements := result.Elements()
		assert.Equal(t, 1, len(elements), "Should contain one element")
	})
}
