package attr_test

import (
	"context"
	"testing"

	localattr "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/attr"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

func TestObjectValue(t *testing.T) {
	attrTypes := map[string]attr.Type{
		"name": types.StringType,
		"age":  types.Int64Type,
	}

	values := map[string]attr.Value{
		"name": types.StringValue("test"),
		"age":  types.Int64Value(42),
	}

	obj := localattr.ObjectValue(attrTypes, values)
	assert.False(t, obj.IsNull())
	assert.False(t, obj.IsUnknown())

	attrs := obj.Attributes()
	assert.Equal(t, types.StringValue("test"), attrs["name"])
	assert.Equal(t, types.Int64Value(42), attrs["age"])
}

func TestObjectNullIfEmpty(t *testing.T) {
	attrTypes := map[string]attr.Type{
		"name": types.StringType,
		"age":  types.Int64Type,
	}

	// Test with empty values
	emptyValues := map[string]attr.Value{}
	obj := localattr.ObjectNullIfEmpty(attrTypes, emptyValues)
	assert.True(t, obj.IsNull())

	// Test with non-empty values
	values := map[string]attr.Value{
		"name": types.StringValue("test"),
		"age":  types.Int64Value(42),
	}
	obj = localattr.ObjectNullIfEmpty(attrTypes, values)
	assert.False(t, obj.IsNull())
}

func TestGetObjectAttr(t *testing.T) {
	attrTypes := map[string]attr.Type{
		"name": types.StringType,
		"age":  types.Int64Type,
	}

	values := map[string]attr.Value{
		"name": types.StringValue("test"),
		"age":  types.Int64Value(42),
	}

	obj, _ := types.ObjectValue(attrTypes, values)

	// Test existing attribute
	nameAttr := localattr.GetObjectAttr(obj, "name", types.StringNull())
	assert.Equal(t, types.StringValue("test"), nameAttr)

	// Test non-existing attribute
	missingAttr := localattr.GetObjectAttr(obj, "missing", types.StringValue("default"))
	assert.Equal(t, types.StringValue("default"), missingAttr)

	// Test with null object
	nullObj := types.ObjectNull(attrTypes)
	nullAttr := localattr.GetObjectAttr(nullObj, "name", types.StringValue("default"))
	assert.Equal(t, types.StringValue("default"), nullAttr)
}

func TestListValue(t *testing.T) {
	elements := []attr.Value{
		types.StringValue("one"),
		types.StringValue("two"),
		types.StringValue("three"),
	}

	list := localattr.ListValue(types.StringType, elements)
	assert.False(t, list.IsNull())
	assert.Equal(t, 3, len(list.Elements()))
}

func TestListNullIfEmpty(t *testing.T) {
	// Test with empty values
	emptyValues := []attr.Value{}
	list := localattr.ListNullIfEmpty(types.StringType, emptyValues)
	assert.True(t, list.IsNull())

	// Test with non-empty values
	values := []attr.Value{
		types.StringValue("one"),
		types.StringValue("two"),
	}
	list = localattr.ListNullIfEmpty(types.StringType, values)
	assert.False(t, list.IsNull())
	assert.Equal(t, 2, len(list.Elements()))
}

func TestSetValue(t *testing.T) {
	elements := []attr.Value{
		types.StringValue("one"),
		types.StringValue("two"),
		types.StringValue("three"),
	}

	set := localattr.SetValue(types.StringType, elements)
	assert.False(t, set.IsNull())
	assert.Equal(t, 3, len(set.Elements()))
}

func TestSetNullIfEmpty(t *testing.T) {
	ctx := context.Background()

	// Test with empty values
	emptyValues := []attr.Value{}
	set := localattr.SetNullIfEmpty(ctx, types.StringType, emptyValues)
	assert.True(t, set.IsNull())

	// Test with non-empty values
	values := []attr.Value{
		types.StringValue("one"),
		types.StringValue("two"),
	}
	set = localattr.SetNullIfEmpty(ctx, types.StringType, values)
	assert.False(t, set.IsNull())
	assert.Equal(t, 2, len(set.Elements()))
}

func TestMapValue(t *testing.T) {
	elements := map[string]attr.Value{
		"key1": types.StringValue("value1"),
		"key2": types.StringValue("value2"),
	}

	m := localattr.MapValue(types.StringType, elements)
	assert.False(t, m.IsNull())
	assert.Equal(t, 2, len(m.Elements()))
}

func TestMapNullIfEmpty(t *testing.T) {
	// Test with empty values
	emptyValues := map[string]attr.Value{}
	m := localattr.MapNullIfEmpty(types.StringType, emptyValues)
	assert.True(t, m.IsNull())

	// Test with non-empty values
	values := map[string]attr.Value{
		"key1": types.StringValue("value1"),
		"key2": types.StringValue("value2"),
	}
	m = localattr.MapNullIfEmpty(types.StringType, values)
	assert.False(t, m.IsNull())
	assert.Equal(t, 2, len(m.Elements()))
}

func TestStringListElements(t *testing.T) {
	elements := []attr.Value{
		types.StringValue("one"),
		types.StringValue("two"),
		types.StringNull(),
		types.StringValue("three"),
	}

	list, _ := types.ListValue(types.StringType, elements)
	result := localattr.StringListElements(list)
	assert.Equal(t, []string{"one", "two", "three"}, result)

	// Test with null list
	nullList := types.ListNull(types.StringType)
	result = localattr.StringListElements(nullList)
	assert.Nil(t, result)
}

func TestStringSetElements(t *testing.T) {
	elements := []attr.Value{
		types.StringValue("one"),
		types.StringValue("two"),
		types.StringNull(),
		types.StringValue("three"),
	}

	set, _ := types.SetValue(types.StringType, elements)
	result := localattr.StringSetElements(set)
	// Order is not guaranteed in a set
	assert.ElementsMatch(t, []string{"one", "two", "three"}, result)

	// Test with null set
	nullSet := types.SetNull(types.StringType)
	result = localattr.StringSetElements(nullSet)
	assert.Nil(t, result)
}

func TestStringMapElements(t *testing.T) {
	elements := map[string]attr.Value{
		"key1": types.StringValue("value1"),
		"key2": types.StringNull(),
		"key3": types.StringValue("value3"),
	}

	m, _ := types.MapValue(types.StringType, elements)
	result := localattr.StringMapElements(m)
	expected := map[string]string{
		"key1": "value1",
		"key3": "value3",
	}
	assert.Equal(t, expected, result)

	// Test with null map
	nullMap := types.MapNull(types.StringType)
	result = localattr.StringMapElements(nullMap)
	assert.Nil(t, result)
}

func TestObjectSetElements(t *testing.T) {
	attrTypes := map[string]attr.Type{
		"name": types.StringType,
		"age":  types.Int64Type,
	}

	obj1, _ := types.ObjectValue(attrTypes, map[string]attr.Value{
		"name": types.StringValue("Alice"),
		"age":  types.Int64Value(30),
	})

	obj2, _ := types.ObjectValue(attrTypes, map[string]attr.Value{
		"name": types.StringValue("Bob"),
		"age":  types.Int64Value(25),
	})

	elements := []attr.Value{obj1, obj2}
	set, _ := types.SetValue(types.ObjectType{AttrTypes: attrTypes}, elements)

	result := localattr.ObjectSetElements(set)
	assert.Equal(t, 2, len(result))

	// Order is not guaranteed in a set, so we need to check both possible orders
	if result[0]["name"].(types.String).ValueString() == "Alice" {
		assert.Equal(t, "Alice", result[0]["name"].(types.String).ValueString())
		assert.Equal(t, int64(30), result[0]["age"].(types.Int64).ValueInt64())
		assert.Equal(t, "Bob", result[1]["name"].(types.String).ValueString())
		assert.Equal(t, int64(25), result[1]["age"].(types.Int64).ValueInt64())
	} else {
		assert.Equal(t, "Bob", result[0]["name"].(types.String).ValueString())
		assert.Equal(t, int64(25), result[0]["age"].(types.Int64).ValueInt64())
		assert.Equal(t, "Alice", result[1]["name"].(types.String).ValueString())
		assert.Equal(t, int64(30), result[1]["age"].(types.Int64).ValueInt64())
	}
}

func TestObjectListElements(t *testing.T) {
	attrTypes := map[string]attr.Type{
		"name": types.StringType,
		"age":  types.Int64Type,
	}

	obj1, _ := types.ObjectValue(attrTypes, map[string]attr.Value{
		"name": types.StringValue("Alice"),
		"age":  types.Int64Value(30),
	})

	obj2, _ := types.ObjectValue(attrTypes, map[string]attr.Value{
		"name": types.StringValue("Bob"),
		"age":  types.Int64Value(25),
	})

	elements := []attr.Value{obj1, obj2}
	list, _ := types.ListValue(types.ObjectType{AttrTypes: attrTypes}, elements)

	result := localattr.ObjectListElements(list)
	assert.Equal(t, 2, len(result))
	assert.Equal(t, "Alice", result[0]["name"].(types.String).ValueString())
	assert.Equal(t, int64(30), result[0]["age"].(types.Int64).ValueInt64())
	assert.Equal(t, "Bob", result[1]["name"].(types.String).ValueString())
	assert.Equal(t, int64(25), result[1]["age"].(types.Int64).ValueInt64())
}

func TestObjectMapElements(t *testing.T) {
	attrTypes := map[string]attr.Type{
		"name": types.StringType,
		"age":  types.Int64Type,
	}

	obj1, _ := types.ObjectValue(attrTypes, map[string]attr.Value{
		"name": types.StringValue("Alice"),
		"age":  types.Int64Value(30),
	})

	obj2, _ := types.ObjectValue(attrTypes, map[string]attr.Value{
		"name": types.StringValue("Bob"),
		"age":  types.Int64Value(25),
	})

	elements := map[string]attr.Value{
		"person1": obj1,
		"person2": obj2,
	}
	m, _ := types.MapValue(types.ObjectType{AttrTypes: attrTypes}, elements)

	result := localattr.ObjectMapElements(m)
	assert.Equal(t, 2, len(result))
	assert.Equal(t, "Alice", result["person1"]["name"].(types.String).ValueString())
	assert.Equal(t, int64(30), result["person1"]["age"].(types.Int64).ValueInt64())
	assert.Equal(t, "Bob", result["person2"]["name"].(types.String).ValueString())
	assert.Equal(t, int64(25), result["person2"]["age"].(types.Int64).ValueInt64())
}
