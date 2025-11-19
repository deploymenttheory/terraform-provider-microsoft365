package graphBetaAuthenticationStrengthPolicy

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteResourceStateToTerraform maps the SDK response to Terraform state
func MapRemoteResourceStateToTerraform(ctx context.Context, data *AuthenticationStrengthPolicyResourceModel, remoteResource graphmodels.AuthenticationStrengthPolicyable) {
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
	data.CombinationConfigurations = mapCombinationConfigurationsFromSDK(ctx, remoteResource)
}

// mapCombinationConfigurationsFromSDK maps SDK combination configurations to Terraform state
func mapCombinationConfigurationsFromSDK(ctx context.Context, remoteResource graphmodels.AuthenticationStrengthPolicyable) types.List {
	combinationConfigs := remoteResource.GetCombinationConfigurations()
	if len(combinationConfigs) == 0 {
		return types.ListNull(types.ObjectType{
			AttrTypes: getCombinationConfigurationAttrTypes(),
		})
	}

	configModels := make([]CombinationConfigurationModel, 0, len(combinationConfigs))

	for _, config := range combinationConfigs {
		configModel := CombinationConfigurationModel{}

		// Get basic fields using helpers
		configModel.ID = convert.GraphToFrameworkString(config.GetId())
		configModel.ODataType = convert.GraphToFrameworkString(config.GetOdataType())

		// Get appliesToCombinations - use the first element since it's now a single string
		if appliesToCombos := config.GetAppliesToCombinations(); len(appliesToCombos) > 0 {
			enumStrings := convert.GraphToFrameworkEnumSlice(appliesToCombos)
			if len(enumStrings) > 0 {
				configModel.AppliesToCombinations = enumStrings[0]
			} else {
				configModel.AppliesToCombinations = types.StringNull()
			}
		} else {
			configModel.AppliesToCombinations = types.StringNull()
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

			// Compute allowedIssuers from allowedIssuerSkis by adding "CUSTOMIDENTIFIER:" prefix
			if skis := sdkConfig.GetAllowedIssuerSkis(); len(skis) > 0 {
				issuers := make([]attr.Value, len(skis))
				for i, ski := range skis {
					issuers[i] = types.StringValue("CUSTOMIDENTIFIER:" + ski)
				}
				configModel.AllowedIssuers = types.SetValueMust(types.StringType, issuers)
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

// getCombinationConfigurationAttrTypes returns the attribute types for combination configuration objects
func getCombinationConfigurationAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":                      types.StringType,
		"odata_type":              types.StringType,
		"applies_to_combinations": types.StringType,
		"allowed_aaguids":         types.SetType{ElemType: types.StringType},
		"allowed_issuer_skis":     types.SetType{ElemType: types.StringType},
		"allowed_issuers":         types.SetType{ElemType: types.StringType},
		"allowed_policy_oids":     types.SetType{ElemType: types.StringType},
	}
}
