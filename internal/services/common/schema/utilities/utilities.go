package schema // Or your schema utilities package

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// SchemaAttributesToAttrTypes converts a map of schema.Attribute to a map of attr.Type.
// This is used to provide the necessary type information for framework functions like
// types.ObjectValueFrom, types.SetValueFrom, value.As, etc., by deriving the
// attr.Type map directly from the resource's schema definition.
func SchemaAttributesToAttrTypes(schemaAttributes map[string]schema.Attribute) (map[string]attr.Type, error) {
	attrTypes := make(map[string]attr.Type, len(schemaAttributes))
	var err error

	for attrName, schemaAttr := range schemaAttributes {
		switch typedSchemaAttr := schemaAttr.(type) {
		case schema.BoolAttribute:
			attrTypes[attrName] = types.BoolType
		case schema.Float64Attribute:
			attrTypes[attrName] = types.Float64Type
		case schema.Int64Attribute:
			attrTypes[attrName] = types.Int64Type
		case schema.ListAttribute:
			attrTypes[attrName] = types.ListType{ElemType: typedSchemaAttr.ElementType}
		case schema.ListNestedAttribute:
			var nestedAttrTypes map[string]attr.Type
			nestedAttrTypes, err = SchemaAttributesToAttrTypes(typedSchemaAttr.NestedObject.Attributes)
			if err != nil {
				return nil, fmt.Errorf("failed to convert nested attributes for list attribute '%s': %w", attrName, err)
			}
			attrTypes[attrName] = types.ListType{ElemType: types.ObjectType{AttrTypes: nestedAttrTypes}}
		case schema.MapAttribute:
			attrTypes[attrName] = types.MapType{ElemType: typedSchemaAttr.ElementType}
		case schema.MapNestedAttribute:
			var nestedAttrTypes map[string]attr.Type
			nestedAttrTypes, err = SchemaAttributesToAttrTypes(typedSchemaAttr.NestedObject.Attributes)
			if err != nil {
				return nil, fmt.Errorf("failed to convert nested attributes for map attribute '%s': %w", attrName, err)
			}
			attrTypes[attrName] = types.MapType{ElemType: types.ObjectType{AttrTypes: nestedAttrTypes}}
		case schema.NumberAttribute:
			attrTypes[attrName] = types.NumberType
		case schema.ObjectAttribute:
			attrTypes[attrName] = types.ObjectType{AttrTypes: typedSchemaAttr.AttributeTypes}
		case schema.SetAttribute:
			attrTypes[attrName] = types.SetType{ElemType: typedSchemaAttr.ElementType}
		case schema.SetNestedAttribute:
			var nestedAttrTypes map[string]attr.Type
			nestedAttrTypes, err = SchemaAttributesToAttrTypes(typedSchemaAttr.NestedObject.Attributes)
			if err != nil {
				return nil, fmt.Errorf("failed to convert nested attributes for set attribute '%s': %w", attrName, err)
			}
			attrTypes[attrName] = types.SetType{ElemType: types.ObjectType{AttrTypes: nestedAttrTypes}}
		case schema.SingleNestedAttribute:
			var nestedAttrTypes map[string]attr.Type
			nestedAttrTypes, err = SchemaAttributesToAttrTypes(typedSchemaAttr.Attributes)
			if err != nil {
				return nil, fmt.Errorf("failed to convert nested attributes for single nested attribute '%s': %w", attrName, err)
			}
			attrTypes[attrName] = types.ObjectType{AttrTypes: nestedAttrTypes}
		case schema.StringAttribute:
			attrTypes[attrName] = types.StringType
		default:
			return nil, fmt.Errorf("unhandled schema.Attribute type for attribute '%s': %T", attrName, typedSchemaAttr)
		}
	}
	return attrTypes, nil
}
