package graphSubscribedSkus

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	graphmodels "github.com/microsoftgraph/msgraph-sdk-go/models"
)

// MapRemoteStateToDataSource maps a SubscribedSku to a model
func MapRemoteStateToDataSource(data graphmodels.SubscribedSkuable) SubscribedSkuModel {
	model := SubscribedSkuModel{
		ID:               convert.GraphToFrameworkString(data.GetId()),
		AccountId:        convert.GraphToFrameworkString(data.GetAccountId()),
		AccountName:      convert.GraphToFrameworkString(data.GetAccountName()),
		AppliesTo:        convert.GraphToFrameworkString(data.GetAppliesTo()),
		CapabilityStatus: convert.GraphToFrameworkString(data.GetCapabilityStatus()),
		SkuId:            convert.GraphToFrameworkUUID(data.GetSkuId()),
		SkuPartNumber:    convert.GraphToFrameworkString(data.GetSkuPartNumber()),
		ConsumedUnits:    convert.GraphToFrameworkInt32(data.GetConsumedUnits()),
	}

	model.PrepaidUnits = mapPrepaidUnitsToState(data.GetPrepaidUnits())
	model.ServicePlans = mapServicePlansToState(data.GetServicePlans())
	model.SubscriptionIds = mapSubscriptionIdsToState(data.GetSubscriptionIds())

	return model
}

// mapPrepaidUnitsToState maps the prepaid units to state
func mapPrepaidUnitsToState(prepaidUnits graphmodels.LicenseUnitsDetailable) types.Object {
	if prepaidUnits == nil {
		return types.ObjectNull(map[string]attr.Type{
			"enabled":    types.Int32Type,
			"locked_out": types.Int32Type,
			"suspended":  types.Int32Type,
			"warning":    types.Int32Type,
		})
	}

	attrs := map[string]attr.Value{
		"enabled":    convert.GraphToFrameworkInt32(prepaidUnits.GetEnabled()),
		"locked_out": convert.GraphToFrameworkInt32(prepaidUnits.GetLockedOut()),
		"suspended":  convert.GraphToFrameworkInt32(prepaidUnits.GetSuspended()),
		"warning":    convert.GraphToFrameworkInt32(prepaidUnits.GetWarning()),
	}

	attrTypes := map[string]attr.Type{
		"enabled":    types.Int32Type,
		"locked_out": types.Int32Type,
		"suspended":  types.Int32Type,
		"warning":    types.Int32Type,
	}

	result, diags := types.ObjectValue(attrTypes, attrs)
	if diags.HasError() {
		return types.ObjectNull(attrTypes)
	}
	return result
}

// mapServicePlansToState maps service plans to state
func mapServicePlansToState(servicePlans []graphmodels.ServicePlanInfoable) types.List {
	if len(servicePlans) == 0 {
		return types.ListNull(types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"service_plan_id":     types.StringType,
				"service_plan_name":   types.StringType,
				"provisioning_status": types.StringType,
				"applies_to":          types.StringType,
			},
		})
	}

	attrTypes := map[string]attr.Type{
		"service_plan_id":     types.StringType,
		"service_plan_name":   types.StringType,
		"provisioning_status": types.StringType,
		"applies_to":          types.StringType,
	}

	items := make([]types.Object, 0, len(servicePlans))
	for _, plan := range servicePlans {
		if plan == nil {
			continue
		}

		attrs := map[string]attr.Value{
			"service_plan_id":     convert.GraphToFrameworkUUID(plan.GetServicePlanId()),
			"service_plan_name":   convert.GraphToFrameworkString(plan.GetServicePlanName()),
			"provisioning_status": convert.GraphToFrameworkString(plan.GetProvisioningStatus()),
			"applies_to":          convert.GraphToFrameworkString(plan.GetAppliesTo()),
		}

		obj, diags := types.ObjectValue(attrTypes, attrs)
		if diags.HasError() {
			continue
		}
		items = append(items, obj)
	}

	// Create the list from the objects
	result, diags := types.ListValueFrom(context.Background(), types.ObjectType{AttrTypes: attrTypes}, items)
	if diags.HasError() {
		return types.ListNull(types.ObjectType{AttrTypes: attrTypes})
	}
	return result
}

// mapSubscriptionIdsToState maps subscription IDs to state
func mapSubscriptionIdsToState(subscriptionIds []string) types.List {
	if len(subscriptionIds) == 0 {
		return types.ListNull(types.StringType)
	}

	items := make([]types.String, 0, len(subscriptionIds))
	for _, id := range subscriptionIds {
		items = append(items, types.StringValue(id))
	}

	// Create the list from the strings
	result, diags := types.ListValueFrom(context.Background(), types.StringType, items)
	if diags.HasError() {
		return types.ListNull(types.StringType)
	}
	return result
}

// getSubscribedSkuObjectType returns the object type definition for a SubscribedSku
func getSubscribedSkuObjectType() map[string]attr.Type {
	return map[string]attr.Type{
		"id":                types.StringType,
		"account_id":        types.StringType,
		"account_name":      types.StringType,
		"applies_to":        types.StringType,
		"capability_status": types.StringType,
		"consumed_units":    types.Int32Type,
		"sku_id":            types.StringType,
		"sku_part_number":   types.StringType,
		"prepaid_units": types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"enabled":    types.Int32Type,
				"locked_out": types.Int32Type,
				"suspended":  types.Int32Type,
				"warning":    types.Int32Type,
			},
		},
		"service_plans": types.ListType{
			ElemType: types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"service_plan_id":     types.StringType,
					"service_plan_name":   types.StringType,
					"provisioning_status": types.StringType,
					"applies_to":          types.StringType,
				},
			},
		},
		"subscription_ids": types.ListType{
			ElemType: types.StringType,
		},
	}
}
