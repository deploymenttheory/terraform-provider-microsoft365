package graphBetaGroupPolicyConfigurations

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteResourceStateToTerraform maps the remote group policy configuration to Terraform state
func MapRemoteResourceStateToTerraform(ctx context.Context, data *GroupPolicyConfigurationResourceModel, remoteResource graphmodels.GroupPolicyConfigurationable) error {
	tflog.Debug(ctx, fmt.Sprintf("Mapping remote resource state to Terraform for %s", ResourceName))

	if remoteResource == nil {
		tflog.Debug(ctx, "Remote resource is nil")
		return fmt.Errorf("remote resource is nil")
	}

	// Basic properties using SDK getters and convert helpers
	data.ID = convert.GraphToFrameworkString(remoteResource.GetId())
	data.DisplayName = convert.GraphToFrameworkString(remoteResource.GetDisplayName())
	data.Description = convert.GraphToFrameworkString(remoteResource.GetDescription())

	// Handle time.Time for dates
	if createdDateTime := remoteResource.GetCreatedDateTime(); createdDateTime != nil {
		data.CreatedDateTime = types.StringValue(createdDateTime.Format("2006-01-02T15:04:05.000Z"))
	} else {
		data.CreatedDateTime = types.StringNull()
	}

	if lastModifiedDateTime := remoteResource.GetLastModifiedDateTime(); lastModifiedDateTime != nil {
		data.LastModifiedDateTime = types.StringValue(lastModifiedDateTime.Format("2006-01-02T15:04:05.000Z"))
	} else {
		data.LastModifiedDateTime = types.StringNull()
	}

	// Handle role scope tag IDs
	if roleScopeTagIds := remoteResource.GetRoleScopeTagIds(); roleScopeTagIds != nil {
		data.RoleScopeTagIds = convert.GraphToFrameworkStringSet(ctx, roleScopeTagIds)
	} else {
		data.RoleScopeTagIds = types.SetNull(types.StringType)
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished mapping remote resource state to Terraform for %s", ResourceName))
	return nil
}

// MapRemoteDefinitionValuesToTerraform maps the remote definition values to Terraform state
// This is called separately because definition values are fetched from a different endpoint
func MapRemoteDefinitionValuesToTerraform(ctx context.Context, data *GroupPolicyConfigurationResourceModel, definitionValues []graphmodels.GroupPolicyDefinitionValueable, lookupService *DefinitionLookupService) error {
	tflog.Debug(ctx, fmt.Sprintf("Mapping remote definition values to Terraform for %s", ResourceName))

	if len(definitionValues) == 0 {
		data.DefinitionValues = types.SetNull(types.ObjectType{
			AttrTypes: getDefinitionValueAttrTypes(),
		})
		return nil
	}

	// Convert each definition value
	definitionValueObjects := make([]attr.Value, 0, len(definitionValues))

	for _, defValue := range definitionValues {
		if defValue == nil {
			continue
		}

		// Map basic definition value properties using SDK methods
		defValueAttrs := map[string]attr.Value{
			"id":                      convert.GraphToFrameworkString(defValue.GetId()),
			"enabled":                 convert.GraphToFrameworkBool(defValue.GetEnabled()),
			"created_date_time":       convert.GraphToFrameworkTime(defValue.GetCreatedDateTime()),
			"last_modified_date_time": convert.GraphToFrameworkTime(defValue.GetLastModifiedDateTime()),
		}

		// Map configuration type
		if configType := defValue.GetConfigurationType(); configType != nil {
			defValueAttrs["configuration_type"] = types.StringValue(configType.String())
		} else {
			defValueAttrs["configuration_type"] = types.StringNull()
		}

		// Extract definition ID and perform reverse lookup for display_name and class_type
		if definition := defValue.GetDefinition(); definition != nil {
			definitionID := definition.GetId()
			defValueAttrs["definition_id"] = convert.GraphToFrameworkString(definitionID)

			// Perform reverse lookup to get display_name, class_type, and category_path
			if definitionID != nil && *definitionID != "" {
				displayName, classType, categoryPath, err := lookupService.LookupDefinitionInfo(ctx, *definitionID)
				if err != nil {
					tflog.Warn(ctx, fmt.Sprintf("[STATE] Failed to reverse lookup definition info for ID='%s': %s", *definitionID, err.Error()))
					// Set to null values if lookup fails
					defValueAttrs["display_name"] = types.StringNull()
					defValueAttrs["class_type"] = types.StringNull()
					defValueAttrs["category_path"] = types.StringNull()
				} else {
					tflog.Debug(ctx, fmt.Sprintf("[STATE] Reverse lookup successful: ID='%s' -> DisplayName='%s', ClassType='%s', CategoryPath='%s'", *definitionID, displayName, classType, categoryPath))
					defValueAttrs["display_name"] = types.StringValue(displayName)
					defValueAttrs["class_type"] = types.StringValue(classType)
					if categoryPath != "" {
						defValueAttrs["category_path"] = types.StringValue(categoryPath)
					} else {
						defValueAttrs["category_path"] = types.StringNull()
					}
				}
			} else {
				defValueAttrs["display_name"] = types.StringNull()
				defValueAttrs["class_type"] = types.StringNull()
				defValueAttrs["category_path"] = types.StringNull()
			}
		} else {
			defValueAttrs["definition_id"] = types.StringNull()
			defValueAttrs["display_name"] = types.StringNull()
			defValueAttrs["class_type"] = types.StringNull()
			defValueAttrs["category_path"] = types.StringNull()
		}

		// Map presentation values from the SDK
		if presentationValues := defValue.GetPresentationValues(); presentationValues != nil && len(presentationValues) > 0 {
			presentationValuesSet, err := MapRemotePresentationValuesToTerraformSDK(ctx, presentationValues)
			if err != nil {
				return fmt.Errorf("failed to map presentation values: %v", err)
			}
			defValueAttrs["presentation_values"] = presentationValuesSet
		} else {
			// If no presentationValues in API response, set to empty set (not null)
			defValueAttrs["presentation_values"] = types.SetValueMust(types.ObjectType{
				AttrTypes: getPresentationValueAttrTypes(),
			}, []attr.Value{})
		}

		// Create the definition value object
		defValueObj, diag := types.ObjectValue(getDefinitionValueAttrTypes(), defValueAttrs)
		if diag.HasError() {
			return fmt.Errorf("failed to create definition value object: %v", diag)
		}

		definitionValueObjects = append(definitionValueObjects, defValueObj)
	}

	// Create the set of definition values
	definitionValuesSet, diag := types.SetValue(types.ObjectType{
		AttrTypes: getDefinitionValueAttrTypes(),
	}, definitionValueObjects)
	if diag.HasError() {
		return fmt.Errorf("failed to create definition values set: %v", diag)
	}

	data.DefinitionValues = definitionValuesSet

	tflog.Debug(ctx, fmt.Sprintf("Finished mapping remote definition values to Terraform for %s", ResourceName))
	return nil
}

// MapRemotePresentationValuesToTerraform maps the remote presentation values to a definition value
// This is called when we need to populate presentation values from the API response
func MapRemotePresentationValuesToTerraform(ctx context.Context, presentationValues []interface{}) (types.Set, error) {
	if len(presentationValues) == 0 {
		return types.SetValueMust(types.ObjectType{
			AttrTypes: getPresentationValueAttrTypes(),
		}, []attr.Value{}), nil
	}

	presentationValueObjects := make([]attr.Value, 0, len(presentationValues))

	for _, presValueInterface := range presentationValues {
		presValueMap, ok := presValueInterface.(map[string]interface{})
		if !ok {
			continue
		}

		// Map basic presentation value properties
		presValueAttrs := map[string]attr.Value{
			"id":                      convert.GraphToFrameworkString(getStringPtr(presValueMap, "id")),
			"created_date_time":       convert.GraphToFrameworkString(getStringPtr(presValueMap, "createdDateTime")),
			"last_modified_date_time": convert.GraphToFrameworkString(getStringPtr(presValueMap, "lastModifiedDateTime")),
			"odata_type":              convert.GraphToFrameworkString(getStringPtr(presValueMap, "@odata.type")),
		}

		// Extract presentation ID from the presentation object or odata.bind
		if presentation, ok := presValueMap["presentation"].(map[string]interface{}); ok {
			presValueAttrs["presentation_id"] = convert.GraphToFrameworkString(getStringPtr(presentation, "id"))
		} else {
			presValueAttrs["presentation_id"] = types.StringNull()
		}

		// Handle different value types based on odata type
		odataType := getStringPtr(presValueMap, "@odata.type")
		if odataType != nil {
			switch *odataType {
			case "#microsoft.graph.groupPolicyPresentationValueText":
				presValueAttrs["text_value"] = convert.GraphToFrameworkString(getStringPtr(presValueMap, "value"))
				presValueAttrs["value"] = convert.GraphToFrameworkString(getStringPtr(presValueMap, "value"))
				presValueAttrs["decimal_value"] = types.Int64Null()
				presValueAttrs["boolean_value"] = types.BoolNull()
				presValueAttrs["list_values"] = types.SetNull(types.StringType)
				presValueAttrs["multi_text_values"] = types.SetNull(types.StringType)

			case "#microsoft.graph.groupPolicyPresentationValueDecimal", "#microsoft.graph.groupPolicyPresentationValueLongDecimal":
				presValueAttrs["decimal_value"] = convert.GraphToFrameworkInt64(getInt64Ptr(presValueMap, "value"))
				presValueAttrs["text_value"] = types.StringNull()
				presValueAttrs["value"] = types.StringNull()
				presValueAttrs["boolean_value"] = types.BoolNull()
				presValueAttrs["list_values"] = types.SetNull(types.StringType)
				presValueAttrs["multi_text_values"] = types.SetNull(types.StringType)

			case "#microsoft.graph.groupPolicyPresentationValueBoolean":
				presValueAttrs["boolean_value"] = convert.GraphToFrameworkBool(getBoolPtr(presValueMap, "value"))
				presValueAttrs["text_value"] = types.StringNull()
				presValueAttrs["value"] = types.StringNull()
				presValueAttrs["decimal_value"] = types.Int64Null()
				presValueAttrs["list_values"] = types.SetNull(types.StringType)
				presValueAttrs["multi_text_values"] = types.SetNull(types.StringType)

			case "#microsoft.graph.groupPolicyPresentationValueList":
				// Handle list values
				if values, ok := presValueMap["values"].([]interface{}); ok {
					stringValues := make([]string, 0, len(values))
					for _, value := range values {
						if valueMap, ok := value.(map[string]interface{}); ok {
							if valueStr, ok := valueMap["value"].(string); ok {
								stringValues = append(stringValues, valueStr)
							}
						}
					}
					presValueAttrs["list_values"] = convert.GraphToFrameworkStringSet(ctx, stringValues)
				} else {
					presValueAttrs["list_values"] = types.SetNull(types.StringType)
				}
				presValueAttrs["text_value"] = types.StringNull()
				presValueAttrs["value"] = types.StringNull()
				presValueAttrs["decimal_value"] = types.Int64Null()
				presValueAttrs["boolean_value"] = types.BoolNull()
				presValueAttrs["multi_text_values"] = types.SetNull(types.StringType)

			case "#microsoft.graph.groupPolicyPresentationValueMultiText":
				// Handle multi-text values
				if values, ok := presValueMap["values"].([]interface{}); ok {
					stringValues := make([]string, 0, len(values))
					for _, value := range values {
						if valueStr, ok := value.(string); ok {
							stringValues = append(stringValues, valueStr)
						}
					}
					presValueAttrs["multi_text_values"] = convert.GraphToFrameworkStringSet(ctx, stringValues)
				} else {
					presValueAttrs["multi_text_values"] = types.SetNull(types.StringType)
				}
				presValueAttrs["text_value"] = types.StringNull()
				presValueAttrs["value"] = types.StringNull()
				presValueAttrs["decimal_value"] = types.Int64Null()
				presValueAttrs["boolean_value"] = types.BoolNull()
				presValueAttrs["list_values"] = types.SetNull(types.StringType)

			default:
				// For unknown types, try to use the generic value field
				presValueAttrs["value"] = convert.GraphToFrameworkString(getStringPtr(presValueMap, "value"))
				presValueAttrs["text_value"] = types.StringNull()
				presValueAttrs["decimal_value"] = types.Int64Null()
				presValueAttrs["boolean_value"] = types.BoolNull()
				presValueAttrs["list_values"] = types.SetNull(types.StringType)
				presValueAttrs["multi_text_values"] = types.SetNull(types.StringType)
			}
		} else {
			// No odata type, set all type-specific fields to null
			presValueAttrs["value"] = types.StringNull()
			presValueAttrs["text_value"] = types.StringNull()
			presValueAttrs["decimal_value"] = types.Int64Null()
			presValueAttrs["boolean_value"] = types.BoolNull()
			presValueAttrs["list_values"] = types.SetNull(types.StringType)
			presValueAttrs["multi_text_values"] = types.SetNull(types.StringType)
		}

		// Create the presentation value object
		presValueObj, diag := types.ObjectValue(getPresentationValueAttrTypes(), presValueAttrs)
		if diag.HasError() {
			return types.SetNull(types.ObjectType{
				AttrTypes: getPresentationValueAttrTypes(),
			}), fmt.Errorf("failed to create presentation value object: %v", diag)
		}

		presentationValueObjects = append(presentationValueObjects, presValueObj)
	}

	// Create the set of presentation values
	presentationValuesSet, diag := types.SetValue(types.ObjectType{
		AttrTypes: getPresentationValueAttrTypes(),
	}, presentationValueObjects)
	if diag.HasError() {
		return types.SetNull(types.ObjectType{
			AttrTypes: getPresentationValueAttrTypes(),
		}), fmt.Errorf("failed to create presentation values set: %v", diag)
	}

	return presentationValuesSet, nil
}

// Helper functions for attribute type definitions
func getDefinitionValueAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":                      types.StringType,
		"enabled":                 types.BoolType,
		"configuration_type":      types.StringType,
		"created_date_time":       types.StringType,
		"last_modified_date_time": types.StringType,
		"display_name":            types.StringType,
		"class_type":              types.StringType,
		"category_path":           types.StringType,
		"definition_id":           types.StringType,
		"presentation_values": types.SetType{
			ElemType: types.ObjectType{
				AttrTypes: getPresentationValueAttrTypes(),
			},
		},
	}
}

func getPresentationValueAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":                      types.StringType,
		"created_date_time":       types.StringType,
		"last_modified_date_time": types.StringType,
		"presentation_id":         types.StringType,
		"odata_type":              types.StringType,
		"value":                   types.StringType,
		"text_value":              types.StringType,
		"decimal_value":           types.Int64Type,
		"boolean_value":           types.BoolType,
		"list_values": types.SetType{
			ElemType: types.StringType,
		},
		"multi_text_values": types.SetType{
			ElemType: types.StringType,
		},
	}
}

// Helper function to get string pointer from map
func getStringPtr(data map[string]interface{}, key string) *string {
	if value, ok := data[key].(string); ok {
		return &value
	}
	return nil
}

// Helper function to get bool pointer from map
func getBoolPtr(data map[string]interface{}, key string) *bool {
	if value, ok := data[key].(bool); ok {
		return &value
	}
	return nil
}

// Helper function to get int64 pointer from map
func getInt64Ptr(data map[string]interface{}, key string) *int64 {
	if value, ok := data[key].(float64); ok {
		int64Val := int64(value)
		return &int64Val
	}
	if value, ok := data[key].(int64); ok {
		return &value
	}
	if value, ok := data[key].(int); ok {
		int64Val := int64(value)
		return &int64Val
	}
	return nil
}

// MapRemotePresentationValuesToTerraformSDK maps SDK presentation values to Terraform state
func MapRemotePresentationValuesToTerraformSDK(ctx context.Context, presentationValues []graphmodels.GroupPolicyPresentationValueable) (types.Set, error) {
	tflog.Debug(ctx, fmt.Sprintf("[STATE] MapRemotePresentationValuesToTerraformSDK: Processing %d presentation values", len(presentationValues)))

	if len(presentationValues) == 0 {
		tflog.Debug(ctx, "[STATE] No presentation values to map, returning empty set")
		return types.SetValueMust(types.ObjectType{
			AttrTypes: getPresentationValueAttrTypes(),
		}, []attr.Value{}), nil
	}

	presentationValueObjects := make([]attr.Value, 0, len(presentationValues))

	for i, presValue := range presentationValues {
		tflog.Debug(ctx, fmt.Sprintf("[STATE] Processing presentation value %d: ID=%s, Type=%T",
			i, stringPtrToString(presValue.GetId()), presValue))
		if presValue == nil {
			continue
		}

		// Map basic presentation value properties using SDK methods
		presValueAttrs := map[string]attr.Value{
			"id":                      convert.GraphToFrameworkString(presValue.GetId()),
			"created_date_time":       convert.GraphToFrameworkTime(presValue.GetCreatedDateTime()),
			"last_modified_date_time": convert.GraphToFrameworkTime(presValue.GetLastModifiedDateTime()),
		}

		// Get presentation ID from the presentation object
		if presentation := presValue.GetPresentation(); presentation != nil {
			presValueAttrs["presentation_id"] = convert.GraphToFrameworkString(presentation.GetId())
		} else {
			presValueAttrs["presentation_id"] = types.StringNull()
		}

		// Get definition value ID from the definition value object
		if definitionValue := presValue.GetDefinitionValue(); definitionValue != nil {
			presValueAttrs["definition_value_id"] = convert.GraphToFrameworkString(definitionValue.GetId())
		} else {
			presValueAttrs["definition_value_id"] = types.StringNull()
		}

		// Map OData type and specific value based on the concrete type
		switch v := presValue.(type) {
		case graphmodels.GroupPolicyPresentationValueTextable:
			tflog.Debug(ctx, fmt.Sprintf("[STATE] Mapping as Text presentation value: ID=%s, Value=%s",
				stringPtrToString(presValue.GetId()), stringPtrToString(v.GetValue())))
			presValueAttrs["odata_type"] = types.StringValue("#microsoft.graph.groupPolicyPresentationValueText")
			presValueAttrs["value"] = convert.GraphToFrameworkString(v.GetValue())
			presValueAttrs["text_value"] = convert.GraphToFrameworkString(v.GetValue())
			presValueAttrs["boolean_value"] = types.BoolNull()
			presValueAttrs["decimal_value"] = types.Int64Null()
			presValueAttrs["list_values"] = types.SetNull(types.StringType)
			presValueAttrs["multi_text_values"] = types.SetNull(types.StringType)
		case graphmodels.GroupPolicyPresentationValueDecimalable:
			tflog.Debug(ctx, fmt.Sprintf("[STATE] Mapping as Decimal presentation value: ID=%s, Value=%v",
				stringPtrToString(presValue.GetId()), v.GetValue()))
			presValueAttrs["odata_type"] = types.StringValue("#microsoft.graph.groupPolicyPresentationValueDecimal")
			if decimalValue := v.GetValue(); decimalValue != nil {
				presValueAttrs["value"] = types.StringValue(fmt.Sprintf("%d", *decimalValue))
				presValueAttrs["decimal_value"] = types.Int64Value(*decimalValue)
			} else {
				presValueAttrs["value"] = types.StringNull()
				presValueAttrs["decimal_value"] = types.Int64Null()
			}
			presValueAttrs["text_value"] = types.StringNull()
			presValueAttrs["boolean_value"] = types.BoolNull()
			presValueAttrs["list_values"] = types.SetNull(types.StringType)
			presValueAttrs["multi_text_values"] = types.SetNull(types.StringType)
		case graphmodels.GroupPolicyPresentationValueBooleanable:
			tflog.Debug(ctx, fmt.Sprintf("[STATE] Mapping as Boolean presentation value: ID=%s, Value=%v",
				stringPtrToString(presValue.GetId()), v.GetValue()))
			presValueAttrs["odata_type"] = types.StringValue("#microsoft.graph.groupPolicyPresentationValueBoolean")
			presValueAttrs["value"] = convert.GraphToFrameworkString(boolToStringPtr(v.GetValue()))
			presValueAttrs["boolean_value"] = convert.GraphToFrameworkBool(v.GetValue())
			presValueAttrs["text_value"] = types.StringNull()
			presValueAttrs["decimal_value"] = types.Int64Null()
			presValueAttrs["list_values"] = types.SetNull(types.StringType)
			presValueAttrs["multi_text_values"] = types.SetNull(types.StringType)
		default:
			tflog.Debug(ctx, fmt.Sprintf("[STATE] Mapping as Unknown presentation value type: ID=%s, Type=%T",
				stringPtrToString(presValue.GetId()), presValue))
			presValueAttrs["odata_type"] = types.StringValue("#microsoft.graph.groupPolicyPresentationValue")
			presValueAttrs["value"] = types.StringNull()
			presValueAttrs["text_value"] = types.StringNull()
			presValueAttrs["boolean_value"] = types.BoolNull()
			presValueAttrs["decimal_value"] = types.Int64Null()
			presValueAttrs["list_values"] = types.SetNull(types.StringType)
			presValueAttrs["multi_text_values"] = types.SetNull(types.StringType)
		}

		presValueObj, diags := types.ObjectValue(getPresentationValueAttrTypes(), presValueAttrs)
		if diags.HasError() {
			return types.SetNull(types.ObjectType{AttrTypes: getPresentationValueAttrTypes()}), fmt.Errorf("failed to create presentation value object: %v", diags.Errors())
		}

		presentationValueObjects = append(presentationValueObjects, presValueObj)
	}

	return types.SetValueMust(types.ObjectType{
		AttrTypes: getPresentationValueAttrTypes(),
	}, presentationValueObjects), nil
}

// Helper function to convert bool pointer to string pointer
func boolToStringPtr(value *bool) *string {
	if value == nil {
		return nil
	}
	result := fmt.Sprintf("%t", *value)
	return &result
}
