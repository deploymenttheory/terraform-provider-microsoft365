package graphBetaGroupPolicyBooleanValue

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/devicemanagement"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// constructResource constructs the updateDefinitionValues request for group policy boolean value operations
func constructResource(ctx context.Context, data *GroupPolicyBooleanValueResourceModel, operation string) (*devicemanagement.GroupPolicyConfigurationsItemUpdateDefinitionValuesPostRequestBody, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing updateDefinitionValues request for %s", ResourceName))

	requestBody := devicemanagement.NewGroupPolicyConfigurationsItemUpdateDefinitionValuesPostRequestBody()

	definitionValue := models.NewGroupPolicyDefinitionValue()

	convert.FrameworkToGraphBool(data.Enabled, definitionValue.SetEnabled)

	definitionID := data.GroupPolicyDefinitionValueID.ValueString() // This contains the definition ID after resolveIDs
	definitionBindURL := fmt.Sprintf("https://graph.microsoft.com/beta/deviceManagement/groupPolicyDefinitions('%s')", definitionID)

	definitionValue.SetAdditionalData(map[string]any{
		"definition@odata.bind": definitionBindURL,
	})

	// Convert the Terraform list to Go slice of BooleanPresentationValue structs
	var booleanValues []BooleanPresentationValue
	data.Values.ElementsAs(ctx, &booleanValues, false)

	// Get resolved presentations from AdditionalData
	resolvedPresentations, ok := data.AdditionalData["resolvedPresentations"].([]ResolvedPresentation)
	if !ok && operation != "delete" {
		return nil, fmt.Errorf("missing resolved presentations in AdditionalData")
	}

	var presentationValues []models.GroupPolicyPresentationValueable

	// Create presentation values for each boolean value
	for i, boolValue := range booleanValues {
		booleanPresentationValue := models.NewGroupPolicyPresentationValueBoolean()
		odataType := "#microsoft.graph.groupPolicyPresentationValueBoolean"
		booleanPresentationValue.SetOdataType(&odataType)

		convert.FrameworkToGraphBool(boolValue.Value, booleanPresentationValue.SetValue)

		// Get the appropriate presentation ID
		var presentationID string
		if operation == "delete" {
			// For delete, use the stored presentation ID from the value
			presentationID = boolValue.PresentationID.ValueString()
		} else if i < len(resolvedPresentations) {
			// Use resolved presentation ID
			presentationID = resolvedPresentations[i].TemplateID
		} else {
			return nil, fmt.Errorf("no presentation ID available for boolean value at index %d", i)
		}

		presentationBindURL := fmt.Sprintf("https://graph.microsoft.com/beta/deviceManagement/groupPolicyDefinitions('%s')/presentations('%s')", definitionID, presentationID)

		booleanPresentationValue.SetAdditionalData(map[string]any{
			"presentation@odata.bind": presentationBindURL,
		})

		presentationValues = append(presentationValues, booleanPresentationValue)
	}

	definitionValue.SetPresentationValues(presentationValues)

	// the request body supports add, update, and delete with distinct structures for each operation.
	switch operation {
	case "create":
		addedValues := []models.GroupPolicyDefinitionValueable{definitionValue}
		requestBody.SetAdded(addedValues)
		requestBody.SetUpdated([]models.GroupPolicyDefinitionValueable{})
		requestBody.SetDeletedIds([]string{})

	case "update":
		// For update, use template IDs for bindings and instance IDs for the update
		// Get instance IDs from AdditionalData (set by the ID resolver)
		definitionValueInstanceID, ok := data.AdditionalData["definitionValueInstanceID"].(string)
		if !ok {
			return nil, fmt.Errorf("missing definitionValueInstanceID in AdditionalData for update operation")
		}

		definitionValue.SetId(&definitionValueInstanceID)

		// Set instance IDs for each presentation value
		for i, presentationValue := range presentationValues {
			if i < len(resolvedPresentations) {
				instanceID := resolvedPresentations[i].InstanceID
				if boolValue, ok := presentationValue.(*models.GroupPolicyPresentationValueBoolean); ok {
					boolValue.SetId(&instanceID)
				}
			}
		}

		updatedValues := []models.GroupPolicyDefinitionValueable{definitionValue}
		requestBody.SetAdded([]models.GroupPolicyDefinitionValueable{})
		requestBody.SetUpdated(updatedValues)
		requestBody.SetDeletedIds([]string{})

	case "delete":
		definitionValueInstanceID := data.GroupPolicyDefinitionValueID.ValueString()
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
