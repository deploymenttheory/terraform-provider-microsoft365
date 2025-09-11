package graphBetaGroupPolicyMultiTextValue

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// constructResource maps the Terraform resource data to the Graph API request model
func constructResource(ctx context.Context, data *GroupPolicyMultiTextValueResourceModel, client *msgraphbetasdk.GraphServiceClient) (models.GroupPolicyPresentationValueMultiTextable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))

	// Resolve IDs if using simplified input
	err := resolveIDs(ctx, data, client)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve IDs: %w", err)
	}

	requestBody := models.NewGroupPolicyPresentationValueMultiText()

	// Set the OData type
	odataType := "#microsoft.graph.groupPolicyPresentationValueMultiText"
	requestBody.SetOdataType(&odataType)

	// Convert List to []string and set values
	if !data.Values.IsNull() && !data.Values.IsUnknown() {
		var stringValues []string
		diags := data.Values.ElementsAs(ctx, &stringValues, false)
		if diags.HasError() {
			return nil, fmt.Errorf("failed to convert values list to string slice")
		}
		requestBody.SetValues(stringValues)
	}

	// Set the presentation reference using resolved IDs
	if !data.PresentationID.IsNull() && !data.PresentationID.IsUnknown() {
		groupPolicyConfigurationID := data.GroupPolicyConfigurationID.ValueString()
		groupPolicyDefinitionValueID := data.GroupPolicyDefinitionValueID.ValueString()

		presentationBindURL := fmt.Sprintf(
			"https://graph.microsoft.com/beta/deviceManagement/groupPolicyConfigurations('%s')/definitionValues('%s')/presentation",
			groupPolicyConfigurationID,
			groupPolicyDefinitionValueID,
		)

		additionalData := map[string]interface{}{
			"presentation@odata.bind": presentationBindURL,
		}
		requestBody.SetAdditionalData(additionalData)

		tflog.Debug(ctx, fmt.Sprintf("Set presentation@odata.bind to: %s", presentationBindURL))
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s resource", ResourceName))
	return requestBody, nil
}

// resolveIDs resolves the definition value and presentation IDs from simplified input
func resolveIDs(ctx context.Context, data *GroupPolicyMultiTextValueResourceModel, client *msgraphbetasdk.GraphServiceClient) error {
	// Skip resolution if IDs are already set (backward compatibility)
	if !data.GroupPolicyDefinitionValueID.IsNull() && !data.GroupPolicyDefinitionValueID.IsUnknown() &&
		!data.PresentationID.IsNull() && !data.PresentationID.IsUnknown() {
		tflog.Debug(ctx, "IDs already resolved, skipping lookup")
		return nil
	}

	// Check if we have the required fields for lookup
	if data.PolicyName.IsNull() || data.PolicyName.IsUnknown() ||
		data.ClassType.IsNull() || data.ClassType.IsUnknown() {
		return fmt.Errorf("either provide group_policy_definition_value_id and presentation_id, or provide policy_name and class_type for auto-discovery")
	}

	// Get presentation index (default to 0)
	presentationIndex := int64(0)
	if !data.PresentationIndex.IsNull() && !data.PresentationIndex.IsUnknown() {
		presentationIndex = data.PresentationIndex.ValueInt64()
	}

	// Use lookup service to resolve IDs
	lookupService := NewLookupService(client)
	definitionValueID, presentationID, err := lookupService.ResolveDefinitionValueAndPresentation(
		ctx,
		data.GroupPolicyConfigurationID.ValueString(),
		data.PolicyName.ValueString(),
		data.ClassType.ValueString(),
		presentationIndex,
	)

	if err != nil {
		return err
	}

	// Set the resolved IDs in the model
	data.GroupPolicyDefinitionValueID = types.StringValue(definitionValueID)
	data.PresentationID = types.StringValue(presentationID)

	tflog.Debug(ctx, fmt.Sprintf("Resolved IDs - definitionValueID: %s, presentationID: %s", definitionValueID, presentationID))
	return nil
}
