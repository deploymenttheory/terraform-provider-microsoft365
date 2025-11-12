package graphBetaGroupPolicyTextValue

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/devicemanagement"
)

// populateMetadataFromAPI fetches policy metadata (policy_name, class_type, category_path)
// from the Graph API when these fields are missing (e.g., during import).
// This allows the resolver to function properly even when starting from just IDs.
func (r *GroupPolicyTextValueResource) populateMetadataFromAPI(
	ctx context.Context,
	data *GroupPolicyTextValueResourceModel,
) error {
	configID := data.GroupPolicyConfigurationID.ValueString()
	definitionValueID := data.GroupPolicyDefinitionValueID.ValueString()

	tflog.Debug(ctx, fmt.Sprintf("[TEXT_VALUE] Fetching metadata from API for configID=%s, definitionValueID=%s",
		configID, definitionValueID))

	// Fetch definition values with expansion to get the definition metadata
	// We fetch all definition values and find ours because single-item GET doesn't support $expand
	requestConfig := &devicemanagement.GroupPolicyConfigurationsItemDefinitionValuesRequestBuilderGetRequestConfiguration{
		QueryParameters: &devicemanagement.GroupPolicyConfigurationsItemDefinitionValuesRequestBuilderGetQueryParameters{
			Expand: []string{"definition($select=id,classType,displayName,categoryPath)"},
		},
	}

	definitionValues, err := r.client.
		DeviceManagement().
		GroupPolicyConfigurations().
		ByGroupPolicyConfigurationId(configID).
		DefinitionValues().
		Get(ctx, requestConfig)

	if err != nil {
		return fmt.Errorf("failed to fetch definition values: %w", err)
	}

	// Find our specific definition value
	var definitionValue graphmodels.GroupPolicyDefinitionValueable
	if definitionValues != nil && definitionValues.GetValue() != nil {
		for _, dv := range definitionValues.GetValue() {
			if dv != nil && dv.GetId() != nil && *dv.GetId() == definitionValueID {
				definitionValue = dv
				break
			}
		}
	}

	if definitionValue == nil {
		return fmt.Errorf("could not find definition value with ID %s", definitionValueID)
	}

	// Extract the definition (contains policy metadata)
	definition := definitionValue.GetDefinition()
	if definition == nil {
		return fmt.Errorf("definition value does not contain definition information")
	}

	// Also fetch the presentation value to get the presentation template ID
	presentationValueID := data.ID.ValueString()
	if presentationValueID == "" {
		return fmt.Errorf("presentation value ID is required but was empty")
	}

	presentationValue, err := r.client.
		DeviceManagement().
		GroupPolicyConfigurations().
		ByGroupPolicyConfigurationId(configID).
		DefinitionValues().
		ByGroupPolicyDefinitionValueId(definitionValueID).
		PresentationValues().
		ByGroupPolicyPresentationValueId(presentationValueID).
		Get(ctx, nil)

	if err != nil {
		return fmt.Errorf("failed to fetch presentation value: %w", err)
	}

	// Extract presentation template ID from the presentation value
	// Note: The presentation relationship is a navigation property that may or may not be populated
	// We'll check if it's available, but won't fail if it's not since the resolver can handle that
	var presentationTemplateID string
	if textValue, ok := presentationValue.(graphmodels.GroupPolicyPresentationValueTextable); ok {
		presentation := textValue.GetPresentation()
		if presentation != nil && presentation.GetId() != nil {
			presentationTemplateID = *presentation.GetId()
			data.PresentationID = types.StringValue(presentationTemplateID)
			tflog.Debug(ctx, fmt.Sprintf("[TEXT_VALUE] Extracted presentation template ID: %s", presentationTemplateID))
		}
	}

	// Get policy metadata from definition
	policyName := ""
	if definition.GetDisplayName() != nil {
		policyName = *definition.GetDisplayName()
	}

	classType := ""
	if definition.GetClassType() != nil {
		classType = definition.GetClassType().String()
	}

	categoryPath := ""
	if definition.GetCategoryPath() != nil {
		categoryPath = *definition.GetCategoryPath()
	}

	// Validate we got all required fields
	if policyName == "" || classType == "" || categoryPath == "" {
		return fmt.Errorf("could not extract required policy metadata from definition (policyName=%s, classType=%s, categoryPath=%s)",
			policyName, classType, categoryPath)
	}

	// Populate the metadata fields in the model
	data.PolicyName = types.StringValue(policyName)
	data.ClassType = types.StringValue(classType)
	data.CategoryPath = types.StringValue(categoryPath)

	return nil
}
