package graphBetaGroupPolicyConfigurations

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/devicemanagement"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// constructResource converts the Terraform resource model to a plain map for JSON marshaling
// Returns a map[string]interface{} that can be directly JSON marshaled by the HTTP client
func constructResource(ctx context.Context, data *GroupPolicyConfigurationResourceModel) (map[string]interface{}, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))

	requestBody := make(map[string]interface{})

	// Basic properties using convert helpers
	convert.FrameworkToGraphString(data.DisplayName, func(val *string) {
		if val != nil {
			requestBody["displayName"] = *val
		}
	})

	convert.FrameworkToGraphString(data.Description, func(val *string) {
		if val != nil {
			requestBody["description"] = *val
		}
	})

	// Handle role scope tag IDs
	if err := convert.FrameworkToGraphStringSet(ctx, data.RoleScopeTagIds, func(values []string) {
		if len(values) > 0 {
			requestBody["roleScopeTagIds"] = values
		}
	}); err != nil {
		return nil, fmt.Errorf("failed to convert role scope tag IDs: %w", err)
	}

	// Debug logging using plain JSON marshal
	if debugJSON, err := json.MarshalIndent(requestBody, "", "    "); err == nil {
		tflog.Debug(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), map[string]interface{}{
			"json": "\n" + string(debugJSON),
		})
	} else {
		tflog.Error(ctx, "Failed to debug log object", map[string]interface{}{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s resource", ResourceName))

	return requestBody, nil
}

// constructUpdateDefinitionValuesRequest constructs the request for updating definition values
// This uses the special updateDefinitionValues endpoint that handles the complex nested structure
func constructUpdateDefinitionValuesRequest(ctx context.Context, data *GroupPolicyConfigurationResourceModel, lookupService *DefinitionLookupService) (*devicemanagement.GroupPolicyConfigurationsItemUpdateDefinitionValuesPostRequestBody, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing definition values update request for %s", ResourceName))

	// Create the SDK request body
	requestBody := devicemanagement.NewGroupPolicyConfigurationsItemUpdateDefinitionValuesPostRequestBody()

	// Initialize the arrays for added, updated, and deleted definition values
	addedValues := make([]models.GroupPolicyDefinitionValueable, 0)
	updatedValues := make([]models.GroupPolicyDefinitionValueable, 0)
	deletedIds := make([]string, 0)

	// Process definition values from the Terraform model
	if !data.DefinitionValues.IsNull() && !data.DefinitionValues.IsUnknown() {
		definitionValuesElements := make([]DefinitionValueModel, 0, len(data.DefinitionValues.Elements()))
		diag := data.DefinitionValues.ElementsAs(ctx, &definitionValuesElements, false)
		if diag.HasError() {
			return nil, fmt.Errorf("failed to convert definition values: %v", diag)
		}

		for _, defValue := range definitionValuesElements {
			// Create SDK definition value
			sdkDefValue := models.NewGroupPolicyDefinitionValue()

			// Set enabled status
			convert.FrameworkToGraphBool(defValue.Enabled, func(val *bool) {
				if val != nil {
					sdkDefValue.SetEnabled(val)
				}
			})

			// Lookup definition ID from display name, class type, and optional category path
			var displayName, classType, categoryPath string
			convert.FrameworkToGraphString(defValue.DisplayName, func(val *string) {
				if val != nil {
					displayName = *val
				}
			})
			convert.FrameworkToGraphString(defValue.ClassType, func(val *string) {
				if val != nil {
					classType = *val
				}
			})
			convert.FrameworkToGraphString(defValue.CategoryPath, func(val *string) {
				if val != nil {
					categoryPath = *val
				}
			})

			if displayName == "" || classType == "" {
				return nil, fmt.Errorf("display_name and class_type are required for constructing definition values")
			}

			// Lookup the definition ID
			definitionID, err := lookupService.LookupDefinitionID(ctx, displayName, classType, categoryPath)
			if err != nil {
				return nil, fmt.Errorf("failed to lookup definition ID for displayName='%s', classType='%s': %w", displayName, classType, err)
			}

			tflog.Debug(ctx, fmt.Sprintf("[CONSTRUCT] Resolved displayName='%s', classType='%s' to definitionID='%s'", displayName, classType, definitionID))

			// Set definition reference using odata.bind
			additionalData := map[string]interface{}{
				"definition@odata.bind": fmt.Sprintf("https://graph.microsoft.com/beta/deviceManagement/groupPolicyDefinitions('%s')", definitionID),
			}
			sdkDefValue.SetAdditionalData(additionalData)

			// Process presentation values
			if !defValue.PresentationValues.IsNull() && !defValue.PresentationValues.IsUnknown() {
				presentationValuesElements := make([]PresentationValueModel, 0, len(defValue.PresentationValues.Elements()))
				diag := defValue.PresentationValues.ElementsAs(ctx, &presentationValuesElements, false)
				if diag.HasError() {
					return nil, fmt.Errorf("failed to convert presentation values: %v", diag)
				}

				presentationValues := make([]models.GroupPolicyPresentationValueable, len(presentationValuesElements))
				for i, presValue := range presentationValuesElements {
					var presentationID string
					convert.FrameworkToGraphString(presValue.PresentationID, func(val *string) {
						if val != nil {
							presentationID = *val
						}
					})

					if presentationID == "" {
						return nil, fmt.Errorf("presentation ID is required for constructing presentation values")
					}

					// Create presentation value based on odata type
					var odataType string
					convert.FrameworkToGraphString(presValue.ODataType, func(val *string) {
						if val != nil {
							odataType = *val
						}
					})

					// Create presentation value with OData bind reference
					presBindURL := fmt.Sprintf("https://graph.microsoft.com/beta/deviceManagement/groupPolicyDefinitions('%s')/presentations('%s')", definitionID, presentationID)

					switch odataType {
					case "#microsoft.graph.groupPolicyPresentationValueBoolean":
						boolValue := models.NewGroupPolicyPresentationValueBoolean()
						convert.FrameworkToGraphBool(presValue.BooleanValue, func(val *bool) {
							if val != nil {
								boolValue.SetValue(val)
							}
						})
						// Use additionalData to set both OData type and bind reference
						additionalData := map[string]interface{}{
							"@odata.type":             odataType,
							"presentation@odata.bind": presBindURL,
						}
						boolValue.SetAdditionalData(additionalData)
						presentationValues[i] = boolValue
					case "#microsoft.graph.groupPolicyPresentationValueText":
						textValue := models.NewGroupPolicyPresentationValueText()
						convert.FrameworkToGraphString(presValue.TextValue, func(val *string) {
							if val != nil {
								textValue.SetValue(val)
							}
						})
						// Use additionalData to set both OData type and bind reference
						additionalData := map[string]interface{}{
							"@odata.type":             odataType,
							"presentation@odata.bind": presBindURL,
						}
						textValue.SetAdditionalData(additionalData)
						presentationValues[i] = textValue
					case "#microsoft.graph.groupPolicyPresentationValueDecimal":
						decimalValue := models.NewGroupPolicyPresentationValueDecimal()
						convert.FrameworkToGraphInt64(presValue.DecimalValue, func(val *int64) {
							if val != nil {
								decimalValue.SetValue(val)
							}
						})
						// Use additionalData to set both OData type and bind reference
						additionalData := map[string]interface{}{
							"@odata.type":             odataType,
							"presentation@odata.bind": presBindURL,
						}
						decimalValue.SetAdditionalData(additionalData)
						presentationValues[i] = decimalValue
					default:
						return nil, fmt.Errorf("unsupported odata type: %s", odataType)
					}
				}

				sdkDefValue.SetPresentationValues(presentationValues)
			}

			addedValues = append(addedValues, sdkDefValue)
		}
	}

	// Set the arrays on the request body
	requestBody.SetAdded(addedValues)
	requestBody.SetUpdated(updatedValues)
	requestBody.SetDeletedIds(deletedIds)

	tflog.Debug(ctx, fmt.Sprintf("Successfully constructed definition values update request for %s", ResourceName))

	// Debug log the request body structure
	tflog.Debug(ctx, fmt.Sprintf("[CONSTRUCT] Request body summary: Added=%d, Updated=%d, DeletedIds=%d",
		len(addedValues), len(updatedValues), len(deletedIds)))

	for i, addedValue := range addedValues {
		if addedValue != nil {
			tflog.Debug(ctx, fmt.Sprintf("[CONSTRUCT] Added[%d]: Enabled=%v, PresentationValues=%d",
				i, addedValue.GetEnabled(), len(addedValue.GetPresentationValues())))

			// Log the additional data (OData binds)
			if additionalData := addedValue.GetAdditionalData(); additionalData != nil {
				for key, value := range additionalData {
					tflog.Debug(ctx, fmt.Sprintf("[CONSTRUCT] Added[%d] AdditionalData[%s]: %v", i, key, value))
				}
			}

			// Log presentation values details
			for j, presValue := range addedValue.GetPresentationValues() {
				if presValue != nil {
					odataType := "nil"
					if presValue.GetOdataType() != nil {
						odataType = *presValue.GetOdataType()
					}
					tflog.Debug(ctx, fmt.Sprintf("[CONSTRUCT] Added[%d] PresentationValue[%d]: Type=%s", i, j, odataType))
					if presValue.GetAdditionalData() != nil {
						for key, value := range presValue.GetAdditionalData() {
							tflog.Debug(ctx, fmt.Sprintf("[CONSTRUCT] Added[%d] PresentationValue[%d] AdditionalData[%s]: %v", i, j, key, value))
						}
					}
				}
			}
		}
	}

	return requestBody, nil
}
