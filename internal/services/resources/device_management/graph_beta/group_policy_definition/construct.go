package graphBetaGroupPolicyDefinition

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/devicemanagement"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// constructResource constructs the updateDefinitionValues request for group policy definition operations
func constructResource(ctx context.Context, data *GroupPolicyDefinitionResourceModel, operation string) (*devicemanagement.GroupPolicyConfigurationsItemUpdateDefinitionValuesPostRequestBody, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing updateDefinitionValues request for %s", ResourceName))

	requestBody := devicemanagement.NewGroupPolicyConfigurationsItemUpdateDefinitionValuesPostRequestBody()

	definitionValue := models.NewGroupPolicyDefinitionValue()

	convert.FrameworkToGraphBool(data.Enabled, definitionValue.SetEnabled)

	definitionTemplateID := data.AdditionalData["definitionTemplateID"].(string)
	definitionBindURL := fmt.Sprintf("https://graph.microsoft.com/beta/deviceManagement/groupPolicyDefinitions('%s')", definitionTemplateID)

	definitionValue.SetAdditionalData(map[string]any{
		"definition@odata.bind": definitionBindURL,
	})

	var presentationValues []models.GroupPolicyPresentationValueable

	// For delete operations, we don't need to construct presentation values
	if operation != constants.TfTfOperationDelete {
		// Get resolved presentations from AdditionalData
		resolvedPresentations, ok := data.AdditionalData["resolvedPresentations"].([]ResolvedPresentation)
		if !ok {
			return nil, fmt.Errorf("missing resolved presentations in AdditionalData")
		}

		// Build a map of label -> resolved presentation for quick lookup
		labelToPresentation := make(map[string]ResolvedPresentation)
		for _, resolved := range resolvedPresentations {
			if resolved.Label != "" {
				labelToPresentation[resolved.Label] = resolved
			}
		}

		// Get user-provided values
		var userValues []PresentationValue
		if !data.Values.IsNull() && !data.Values.IsUnknown() {
			data.Values.ElementsAs(ctx, &userValues, false)
		}

		for _, userValue := range userValues {
			label := userValue.Label.ValueString()
			value := userValue.Value.ValueString()

			// Find the matching presentation
			resolved, found := labelToPresentation[label]
			if !found {
				return nil, fmt.Errorf("no presentation found with label '%s'", label)
			}

			// Convert the value based on presentation type
			presValue, err := ConvertValueForType(ctx, label, value, resolved.Type)
			if err != nil {
				return nil, fmt.Errorf("failed to convert value for label '%s': %w", label, err)
			}

			// Set the presentation binding
			presentationBindURL := fmt.Sprintf("https://graph.microsoft.com/beta/deviceManagement/groupPolicyDefinitions('%s')/presentations('%s')",
				definitionTemplateID, resolved.TemplateID)

			presValue.SetAdditionalData(map[string]any{
				"presentation@odata.bind": presentationBindURL,
			})

			// For update operations, set the instance ID
			if operation == constants.TfOperationUpdate && resolved.InstanceID != "" {
				presValue.SetId(&resolved.InstanceID)
			}

			presentationValues = append(presentationValues, presValue)

			tflog.Debug(ctx, fmt.Sprintf("[CONSTRUCT] Added presentation value: label='%s', type='%s', value='%s'", label, resolved.Type, value))
		}
	}

	definitionValue.SetPresentationValues(presentationValues)

	switch operation {
	case constants.TfOperationCreate:
		addedValues := []models.GroupPolicyDefinitionValueable{definitionValue}
		requestBody.SetAdded(addedValues)
		requestBody.SetUpdated([]models.GroupPolicyDefinitionValueable{})
		requestBody.SetDeletedIds([]string{})

	case constants.TfOperationUpdate:
		// For update, use the instance ID
		definitionValueInstanceID, ok := data.AdditionalData["definitionValueInstanceID"].(string)
		if !ok {
			return nil, fmt.Errorf("missing definitionValueInstanceID in AdditionalData for update operation")
		}

		definitionValue.SetId(&definitionValueInstanceID)

		updatedValues := []models.GroupPolicyDefinitionValueable{definitionValue}
		requestBody.SetAdded([]models.GroupPolicyDefinitionValueable{})
		requestBody.SetUpdated(updatedValues)
		requestBody.SetDeletedIds([]string{})

	case constants.TfTfOperationDelete:
		definitionValueInstanceID := data.AdditionalData["definitionValueInstanceID"].(string)
		requestBody.SetAdded([]models.GroupPolicyDefinitionValueable{})
		requestBody.SetUpdated([]models.GroupPolicyDefinitionValueable{})
		requestBody.SetDeletedIds([]string{definitionValueInstanceID})
	}

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final updateDefinitionValues JSON for resource %s", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing updateDefinitionValues request for %s", ResourceName))
	return requestBody, nil
}
