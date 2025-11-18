package graphBetaAuthenticationStrengthPolicy

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// MapRemoteResourceStateToTerraform maps the remote authentication strength policy to Terraform state
func MapRemoteResourceStateToTerraform(ctx context.Context, data *AuthenticationStrengthPolicyResourceModel, remoteResource map[string]any) {
	if id, ok := remoteResource["id"].(string); ok {
		data.ID = types.StringValue(id)
	}

	data.DisplayName = convert.MapToFrameworkString(remoteResource, "displayName")
	data.Description = convert.MapToFrameworkString(remoteResource, "description")
	data.PolicyType = convert.MapToFrameworkString(remoteResource, "policyType")
	data.RequirementsSatisfied = convert.MapToFrameworkString(remoteResource, "requirementsSatisfied")
	data.CreatedDateTime = convert.MapToFrameworkString(remoteResource, "createdDateTime")
	data.ModifiedDateTime = convert.MapToFrameworkString(remoteResource, "modifiedDateTime")
	data.AllowedCombinations = convert.MapToFrameworkStringSet(ctx, remoteResource, "allowedCombinations")

	// Map combination configurations, preserving allowedIssuers from current state (since API never returns it)
	data.CombinationConfigurations = mapCombinationConfigurations(ctx, remoteResource, data.CombinationConfigurations)
}

// mapCombinationConfigurations maps the combinationConfigurations from the API response to Terraform state
// It preserves allowedIssuers from currentState when the API doesn't return it
func mapCombinationConfigurations(ctx context.Context, remoteResource map[string]any, currentState types.List) types.List {
	// Check if combinationConfigurations exists in the response
	combinationConfigsRaw, ok := remoteResource["combinationConfigurations"]
	if !ok {
		return types.ListNull(types.ObjectType{
			AttrTypes: getCombinationConfigurationAttrTypes(),
		})
	}

	combinationConfigsArray, ok := combinationConfigsRaw.([]any)
	if !ok || len(combinationConfigsArray) == 0 {
		return types.ListNull(types.ObjectType{
			AttrTypes: getCombinationConfigurationAttrTypes(),
		})
	}

	// Extract current state configurations to preserve allowedIssuers
	var currentStateConfigs []CombinationConfigurationModel
	if !currentState.IsNull() && !currentState.IsUnknown() {
		_ = currentState.ElementsAs(ctx, &currentStateConfigs, false)
	}

	// Build the list of combination configurations
	configModels := make([]CombinationConfigurationModel, 0, len(combinationConfigsArray))

	for i, configRaw := range combinationConfigsArray {
		configMap, ok := configRaw.(map[string]any)
		if !ok {
			continue
		}

		config := CombinationConfigurationModel{
			ID:                    convert.MapToFrameworkString(configMap, "id"),
			ODataType:             convert.MapToFrameworkString(configMap, "@odata.type"),
			AppliesToCombinations: convert.MapToFrameworkStringSet(ctx, configMap, "appliesToCombinations"),
		}

		// Map fields from API
		config.AllowedAAGUIDs = convert.MapToFrameworkStringSet(ctx, configMap, "allowedAAGUIDs")
		config.AllowedIssuerSkis = convert.MapToFrameworkStringSet(ctx, configMap, "allowedIssuerSkis")
		config.AllowedPolicyOIDs = convert.MapToFrameworkStringSet(ctx, configMap, "allowedPolicyOIDs")

		// Special handling for allowedIssuers: API accepts it but NEVER returns it
		// Always preserve from current state if it exists
		allowedIssuersFromAPI := convert.MapToFrameworkStringSet(ctx, configMap, "allowedIssuers")
		if allowedIssuersFromAPI.IsNull() && i < len(currentStateConfigs) && !currentStateConfigs[i].AllowedIssuers.IsNull() {
			config.AllowedIssuers = currentStateConfigs[i].AllowedIssuers
		} else {
			config.AllowedIssuers = allowedIssuersFromAPI
		}

		configModels = append(configModels, config)
	}

	// Convert to types.List
	listValue, diags := types.ListValueFrom(ctx, types.ObjectType{
		AttrTypes: getCombinationConfigurationAttrTypes(),
	}, configModels)

	if diags.HasError() {
		return types.ListNull(types.ObjectType{
			AttrTypes: getCombinationConfigurationAttrTypes(),
		})
	}

	return listValue
}

// getCombinationConfigurationAttrTypes returns the attribute types for combination configuration objects
func getCombinationConfigurationAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":                      types.StringType,
		"odata_type":              types.StringType,
		"applies_to_combinations": types.SetType{ElemType: types.StringType},
		"allowed_aaguids":         types.SetType{ElemType: types.StringType},
		"allowed_issuer_skis":     types.SetType{ElemType: types.StringType},
		"allowed_issuers":         types.SetType{ElemType: types.StringType},
		"allowed_policy_oids":     types.SetType{ElemType: types.StringType},
	}
}
