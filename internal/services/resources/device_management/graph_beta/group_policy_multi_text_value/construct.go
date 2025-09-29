package graphBetaGroupPolicyMultiTextValue

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/devicemanagement"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// constructResource constructs the updateDefinitionValues request for group policy text value operations
func constructResource(ctx context.Context, data *GroupPolicyMultiTextValueResourceModel, operation string) (*devicemanagement.GroupPolicyConfigurationsItemUpdateDefinitionValuesPostRequestBody, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing updateDefinitionValues request for %s", ResourceName))

	requestBody := devicemanagement.NewGroupPolicyConfigurationsItemUpdateDefinitionValuesPostRequestBody()

	definitionValue := models.NewGroupPolicyDefinitionValue()

	convert.FrameworkToGraphBool(data.Enabled, definitionValue.SetEnabled)

	definitionID := data.GroupPolicyDefinitionValueID.ValueString() // This contains the definition ID after resolveIDs
	definitionBindURL := fmt.Sprintf("https://graph.microsoft.com/beta/deviceManagement/groupPolicyDefinitions('%s')", definitionID)

	definitionValue.SetAdditionalData(map[string]any{
		"definition@odata.bind": definitionBindURL,
	})

	multiTextPresentationValue := models.NewGroupPolicyPresentationValueMultiText()
	odataType := "#microsoft.graph.groupPolicyPresentationValueMultiText"
	multiTextPresentationValue.SetOdataType(&odataType)

	convert.FrameworkToGraphStringSet(ctx, data.Values, multiTextPresentationValue.SetValues)

	presentationID := data.PresentationID.ValueString()
	presentationBindURL := fmt.Sprintf("https://graph.microsoft.com/beta/deviceManagement/groupPolicyDefinitions('%s')/presentations('%s')", definitionID, presentationID)

	multiTextPresentationValue.SetAdditionalData(map[string]any{
		"presentation@odata.bind": presentationBindURL,
	})

	presentationValues := []models.GroupPolicyPresentationValueable{multiTextPresentationValue}
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

		presentationValueInstanceID, ok := data.AdditionalData["presentationValueInstanceID"].(string)
		if !ok {
			return nil, fmt.Errorf("missing presentationValueInstanceID in AdditionalData for update operation")
		}

		definitionValue.SetId(&definitionValueInstanceID)

		multiTextPresentationValue.SetId(&presentationValueInstanceID)

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
