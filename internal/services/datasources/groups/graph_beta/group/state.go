package graphBetaGroup

import (
	"context"
	"slices"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteStateToDataSource maps a Group Graph API object to the Terraform data source model
func MapRemoteStateToDataSource(ctx context.Context, data graphmodels.Groupable, members []string, owners []string, config GroupDataSourceModel) GroupDataSourceModel {

	dynamicMembershipEnabled := data.GetGroupTypes() != nil &&
		slices.Contains(data.GetGroupTypes(), "DynamicMembership")

	assignedLicenses := mapAssignedLicenses(ctx, data.GetAssignedLicenses())

	return GroupDataSourceModel{
		ID:                            convert.GraphToFrameworkString(data.GetId()),
		DisplayName:                   convert.GraphToFrameworkString(data.GetDisplayName()),
		ObjectId:                      convert.GraphToFrameworkString(data.GetId()),
		Description:                   convert.GraphToFrameworkString(data.GetDescription()),
		Classification:                convert.GraphToFrameworkString(data.GetClassification()),
		MailNickname:                  convert.GraphToFrameworkString(data.GetMailNickname()),
		MailEnabled:                   convert.GraphToFrameworkBool(data.GetMailEnabled()),
		SecurityEnabled:               convert.GraphToFrameworkBool(data.GetSecurityEnabled()),
		GroupTypes:                    convert.GraphToFrameworkStringSet(ctx, data.GetGroupTypes()),
		Visibility:                    convert.GraphToFrameworkString(data.GetVisibility()),
		AssignableToRole:              convert.GraphToFrameworkBool(data.GetIsAssignableToRole()),
		DynamicMembershipEnabled:      types.BoolValue(dynamicMembershipEnabled),
		MembershipRule:                convert.GraphToFrameworkString(data.GetMembershipRule()),
		MembershipRuleProcessingState: convert.GraphToFrameworkString(data.GetMembershipRuleProcessingState()),
		CreatedDateTime:               convert.GraphToFrameworkTime(data.GetCreatedDateTime()),
		Mail:                          convert.GraphToFrameworkString(data.GetMail()),
		ProxyAddresses:                convert.GraphToFrameworkStringSet(ctx, data.GetProxyAddresses()),
		AssignedLicenses:              assignedLicenses,
		HasMembersWithLicenseErrors:   convert.GraphToFrameworkBool(data.GetHasMembersWithLicenseErrors()),
		HideFromAddressLists:          convert.GraphToFrameworkBool(data.GetHideFromAddressLists()),
		HideFromOutlookClients:        convert.GraphToFrameworkBool(data.GetHideFromOutlookClients()),
		OnPremisesSyncEnabled:         convert.GraphToFrameworkBool(data.GetOnPremisesSyncEnabled()),
		OnPremisesLastSyncDateTime:    convert.GraphToFrameworkTime(data.GetOnPremisesLastSyncDateTime()),
		OnPremisesSamAccountName:      convert.GraphToFrameworkString(data.GetOnPremisesSamAccountName()),
		OnPremisesDomainName:          convert.GraphToFrameworkString(data.GetOnPremisesDomainName()),
		OnPremisesNetBiosName:         convert.GraphToFrameworkString(data.GetOnPremisesNetBiosName()),
		OnPremisesSecurityIdentifier:  convert.GraphToFrameworkString(data.GetOnPremisesSecurityIdentifier()),
		Members:                       convert.GraphToFrameworkStringSet(ctx, members),
		Owners:                        convert.GraphToFrameworkStringSet(ctx, owners),
		ODataQuery:                    config.ODataQuery,
		Timeouts:                      config.Timeouts,
	}
}

// mapAssignedLicenses converts Graph API AssignedLicense objects to Terraform framework types
func mapAssignedLicenses(ctx context.Context, licenses []graphmodels.AssignedLicenseable) types.List {
	if licenses == nil {
		return types.ListNull(types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"sku_id":         types.StringType,
				"disabled_plans": types.SetType{ElemType: types.StringType},
			},
		})
	}

	var licenseList []attr.Value
	for _, license := range licenses {
		if license == nil {
			continue
		}

		disabledPlans := types.SetNull(types.StringType)
		if license.GetDisabledPlans() != nil && len(license.GetDisabledPlans()) > 0 {
			var plans []string
			for _, plan := range license.GetDisabledPlans() {
				plans = append(plans, plan.String())
			}
			disabledPlans = convert.GraphToFrameworkStringSet(ctx, plans)
		}

		licenseObj, _ := types.ObjectValue(
			map[string]attr.Type{
				"sku_id":         types.StringType,
				"disabled_plans": types.SetType{ElemType: types.StringType},
			},
			map[string]attr.Value{
				"sku_id":         convert.GraphToFrameworkUUID(license.GetSkuId()),
				"disabled_plans": disabledPlans,
			},
		)
		licenseList = append(licenseList, licenseObj)
	}

	if len(licenseList) == 0 {
		return types.ListNull(types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"sku_id":         types.StringType,
				"disabled_plans": types.SetType{ElemType: types.StringType},
			},
		})
	}

	result, _ := types.ListValue(
		types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"sku_id":         types.StringType,
				"disabled_plans": types.SetType{ElemType: types.StringType},
			},
		},
		licenseList,
	)

	return result
}
