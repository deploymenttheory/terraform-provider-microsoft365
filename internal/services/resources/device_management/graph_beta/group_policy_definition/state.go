package graphBetaGroupPolicyDefinition

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteStateToTerraform maps the remote state from Graph API to the Terraform resource model
func MapRemoteStateToTerraform(ctx context.Context, data *GroupPolicyDefinitionResourceModel, presentationValues []graphmodels.GroupPolicyPresentationValueable, definitionValue graphmodels.GroupPolicyDefinitionValueable) {
	tflog.Debug(ctx, "Starting to map remote state to Terraform state", map[string]any{
		"presentationValueCount": len(presentationValues),
	})

	// Map definition value fields
	if definitionValue != nil {
		data.Enabled = convert.GraphToFrameworkBool(definitionValue.GetEnabled())
		data.CreatedDateTime = convert.GraphToFrameworkTime(definitionValue.GetCreatedDateTime())
		data.LastModifiedDateTime = convert.GraphToFrameworkTime(definitionValue.GetLastModifiedDateTime())

		// ID is already set by the resolver in composite format: configID/definitionValueID
		// Do not overwrite it here
	}

	// Get resolved presentations from AdditionalData
	resolvedPresentations, ok := data.AdditionalData["resolvedPresentations"].([]ResolvedPresentation)
	if !ok {
		tflog.Error(ctx, "Missing resolved presentations in AdditionalData")
		// Keep existing values if we can't resolve
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("[STATE] Resolving %d presentations", len(resolvedPresentations)))
	tflog.Debug(ctx, fmt.Sprintf("[STATE] Got %d presentation values from API", len(presentationValues)))

	// Map presentation values
	var mappedValues []PresentationValue

	for i, resolved := range resolvedPresentations {
		tflog.Debug(ctx, fmt.Sprintf("[STATE] [%d] Looking for instance ID: %s, template ID: %s, label: %s",
			i, resolved.InstanceID, resolved.TemplateID, resolved.Label))

		// Find the matching presentation value instance
		found := false
		for j, presValue := range presentationValues {
			if presValue == nil {
				tflog.Debug(ctx, fmt.Sprintf("[STATE]   [%d] presValue is nil", j))
				continue
			}

			presValueID := presValue.GetId()
			if presValueID != nil {
				tflog.Debug(ctx, fmt.Sprintf("[STATE]   [%d] Checking presentation value ID: %s", j, *presValueID))
				if *presValueID == resolved.InstanceID {
					// Extract the value as a string based on type
					valueStr := ExtractValueFromPresentation(presValue)

					mappedValue := PresentationValue{
						ID:    types.StringValue(resolved.TemplateID),
						Label: types.StringValue(resolved.Label),
						Value: types.StringValue(valueStr),
					}
					mappedValues = append(mappedValues, mappedValue)

					tflog.Debug(ctx, fmt.Sprintf("[STATE] ✓ Mapped value [%d]: id='%s', label='%s', type='%s', value='%s'",
						i, resolved.TemplateID, resolved.Label, resolved.Type, valueStr))
					found = true
					break
				}
			} else {
				tflog.Debug(ctx, fmt.Sprintf("[STATE]   [%d] presValueID is nil", j))
			}
		}

		if !found {
			tflog.Warn(ctx, fmt.Sprintf("[STATE] ✗ Could not find presentation value for label='%s', instanceID='%s'", resolved.Label, resolved.InstanceID))
		}
	}

	// Only update values if we found some
	if len(mappedValues) > 0 {
		// Convert to Terraform Set
		valuesSet, diags := types.SetValueFrom(ctx, types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"id":    types.StringType,
				"label": types.StringType,
				"value": types.StringType,
			},
		}, mappedValues)

		if diags.HasError() {
			tflog.Error(ctx, "Failed to convert values to set", map[string]any{
				"errors": diags.Errors(),
			})
			return
		}

		data.Values = valuesSet
		tflog.Debug(ctx, fmt.Sprintf("[STATE] ✓ Finished mapping resource %s with %d values", ResourceName, len(mappedValues)))
	} else {
		tflog.Warn(ctx, "[STATE] No values were mapped - keeping existing values")
	}
}
