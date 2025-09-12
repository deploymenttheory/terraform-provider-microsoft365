package graphBetaGroupPolicyBooleanValue

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteStateToTerraform maps multiple GroupPolicyPresentationValueBoolean instances and GroupPolicyDefinitionValue to the Terraform resource model
func MapRemoteStateToTerraform(ctx context.Context, data *GroupPolicyBooleanValueResourceModel, presentationValues []graphmodels.GroupPolicyPresentationValueable, definitionValue graphmodels.GroupPolicyDefinitionValueable) {
	if len(presentationValues) == 0 {
		tflog.Debug(ctx, "No presentation values provided")
		return
	}

	tflog.Debug(ctx, "Starting to map remote state to Terraform state", map[string]any{
		"presentationValueCount": len(presentationValues),
	})

	// Map the enabled state and timestamps from the definition value
	if definitionValue != nil {
		data.Enabled = convert.GraphToFrameworkBool(definitionValue.GetEnabled())
		data.CreatedDateTime = convert.GraphToFrameworkTime(definitionValue.GetCreatedDateTime())
		data.LastModifiedDateTime = convert.GraphToFrameworkTime(definitionValue.GetLastModifiedDateTime())
	}

	// Get resolved presentations from AdditionalData to maintain order
	resolvedPresentations, ok := data.AdditionalData["resolvedPresentations"].([]ResolvedPresentation)
	if !ok {
		tflog.Error(ctx, "Missing resolved presentations in AdditionalData")
		return
	}

	// Create the values list
	var booleanValues []BooleanPresentationValue

	// Map each resolved presentation to its corresponding presentation value
	for _, resolved := range resolvedPresentations {
		// Find the matching presentation value by instance ID
		for _, presValue := range presentationValues {
			if presValue == nil {
				continue
			}

			presValueID := presValue.GetId()
			if presValueID != nil && *presValueID == resolved.InstanceID {
				// Cast to boolean presentation value
				if boolValue, ok := presValue.(graphmodels.GroupPolicyPresentationValueBooleanable); ok {
					booleanValues = append(booleanValues, BooleanPresentationValue{
						PresentationID: types.StringValue(resolved.TemplateID),
						Value:          convert.GraphToFrameworkBool(boolValue.GetValue()),
					})
					break
				}
			}
		}
	}

	// Convert to Terraform List
	valuesList, diags := types.ListValueFrom(ctx, types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"presentation_id": types.StringType,
			"value":           types.BoolType,
		},
	}, booleanValues)

	if diags.HasError() {
		tflog.Error(ctx, "Failed to convert boolean values to list", map[string]any{
			"errors": diags.Errors(),
		})
		return
	}

	data.Values = valuesList

	tflog.Debug(ctx, fmt.Sprintf("Finished mapping resource %s with %d boolean values", ResourceName, len(booleanValues)))
}
