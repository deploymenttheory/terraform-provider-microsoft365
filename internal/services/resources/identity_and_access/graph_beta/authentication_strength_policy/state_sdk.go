package graphBetaAuthenticationStrengthPolicy

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteResourceStateToTerraformSDK maps the SDK response to Terraform state
func MapRemoteResourceStateToTerraformSDK(ctx context.Context, data *AuthenticationStrengthPolicyResourceModel, remoteResource graphmodels.AuthenticationStrengthPolicyable, currentState types.List) {
	data.ID = convert.GraphToFrameworkString(remoteResource.GetId())
	data.DisplayName = convert.GraphToFrameworkString(remoteResource.GetDisplayName())
	data.Description = convert.GraphToFrameworkString(remoteResource.GetDescription())
	data.PolicyType = convert.GraphToFrameworkEnum(remoteResource.GetPolicyType())
	data.RequirementsSatisfied = convert.GraphToFrameworkEnum(remoteResource.GetRequirementsSatisfied())
	data.CreatedDateTime = convert.GraphToFrameworkTime(remoteResource.GetCreatedDateTime())
	data.ModifiedDateTime = convert.GraphToFrameworkTime(remoteResource.GetModifiedDateTime())

	// Map allowedCombinations using enum slice helper
	if allowedCombos := remoteResource.GetAllowedCombinations(); len(allowedCombos) > 0 {
		enumStrings := convert.GraphToFrameworkEnumSlice(allowedCombos)
		attrValues := make([]attr.Value, len(enumStrings))
		for i, s := range enumStrings {
			attrValues[i] = s
		}
		data.AllowedCombinations = types.SetValueMust(types.StringType, attrValues)
	} else {
		data.AllowedCombinations = types.SetNull(types.StringType)
	}

	// Map combinationConfigurations
	data.CombinationConfigurations = mapCombinationConfigurationsFromSDK(ctx, remoteResource, currentState)
}

// mapCombinationConfigurationsFromSDK maps SDK combination configurations to Terraform state
func mapCombinationConfigurationsFromSDK(ctx context.Context, remoteResource graphmodels.AuthenticationStrengthPolicyable, currentState types.List) types.List {
	combinationConfigs := remoteResource.GetCombinationConfigurations()
	if len(combinationConfigs) == 0 {
		return types.ListNull(types.ObjectType{
			AttrTypes: getCombinationConfigurationAttrTypes(),
		})
	}

	// Extract current state configurations to preserve allowedIssuers
	var currentStateConfigs []CombinationConfigurationModel
	if !currentState.IsNull() && !currentState.IsUnknown() {
		_ = currentState.ElementsAs(ctx, &currentStateConfigs, false)
	}

	configModels := make([]CombinationConfigurationModel, 0, len(combinationConfigs))

	for i, config := range combinationConfigs {
		configModel := CombinationConfigurationModel{}

		// Get basic fields using helpers
		configModel.ID = convert.GraphToFrameworkString(config.GetId())
		configModel.ODataType = convert.GraphToFrameworkString(config.GetOdataType())

		// Get appliesToCombinations using enum slice helper
		if appliesToCombos := config.GetAppliesToCombinations(); len(appliesToCombos) > 0 {
			enumStrings := convert.GraphToFrameworkEnumSlice(appliesToCombos)
			attrValues := make([]attr.Value, len(enumStrings))
			for i, s := range enumStrings {
				attrValues[i] = s
			}
			configModel.AppliesToCombinations = types.SetValueMust(types.StringType, attrValues)
		} else {
			configModel.AppliesToCombinations = types.SetNull(types.StringType)
		}

		// Type-specific fields
		switch sdkConfig := config.(type) {
		case *graphmodels.Fido2CombinationConfiguration:
			configModel.AllowedAAGUIDs = convert.GraphToFrameworkStringSet(ctx, sdkConfig.GetAllowedAAGUIDs())
			configModel.AllowedIssuerSkis = types.SetNull(types.StringType)
			configModel.AllowedPolicyOIDs = types.SetNull(types.StringType)
			configModel.AllowedIssuers = types.SetNull(types.StringType)

		case *graphmodels.X509CertificateCombinationConfiguration:
			configModel.AllowedAAGUIDs = types.SetNull(types.StringType)
			configModel.AllowedIssuerSkis = convert.GraphToFrameworkStringSet(ctx, sdkConfig.GetAllowedIssuerSkis())
			configModel.AllowedPolicyOIDs = convert.GraphToFrameworkStringSet(ctx, sdkConfig.GetAllowedPolicyOIDs())

			// Special handling for allowedIssuers: API never returns it, preserve from current state
			if i < len(currentStateConfigs) && !currentStateConfigs[i].AllowedIssuers.IsNull() {
				configModel.AllowedIssuers = currentStateConfigs[i].AllowedIssuers
			} else {
				configModel.AllowedIssuers = types.SetNull(types.StringType)
			}
		}

		configModels = append(configModels, configModel)
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
