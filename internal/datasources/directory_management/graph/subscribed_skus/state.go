package graphSubscribedSkus

import (
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/state"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	graphmodels "github.com/microsoftgraph/msgraph-sdk-go/models"
)

// MapRemoteStateToDataSource maps a SubscribedSku to a model
func MapRemoteStateToDataSource(data graphmodels.SubscribedSkuable) SubscribedSkuModel {
	model := SubscribedSkuModel{
		ID:               types.StringPointerValue(data.GetId()),
		AccountId:        types.StringPointerValue(data.GetAccountId()),
		AccountName:      types.StringPointerValue(data.GetAccountName()),
		AppliesTo:        types.StringPointerValue(data.GetAppliesTo()),
		CapabilityStatus: types.StringPointerValue(data.GetCapabilityStatus()),
		SkuId:            state.UUIDPtrToTypeString(data.GetSkuId()),
		SkuPartNumber:    types.StringPointerValue(data.GetSkuPartNumber()),
	}

	// Handle consumed units
	if consumedUnits := data.GetConsumedUnits(); consumedUnits != nil {
		model.ConsumedUnits = types.Int64Value(int64(*consumedUnits))
	} else {
		model.ConsumedUnits = types.Int64Null()
	}

	// Handle prepaid units
	if prepaidUnits := data.GetPrepaidUnits(); prepaidUnits != nil {
		prepaidUnitsAttrs := map[string]attr.Value{
			"enabled":    types.Int64Null(),
			"locked_out": types.Int64Null(),
			"suspended":  types.Int64Null(),
			"warning":    types.Int64Null(),
		}

		if enabled := prepaidUnits.GetEnabled(); enabled != nil {
			prepaidUnitsAttrs["enabled"] = types.Int64Value(int64(*enabled))
		}
		if lockedOut := prepaidUnits.GetLockedOut(); lockedOut != nil {
			prepaidUnitsAttrs["locked_out"] = types.Int64Value(int64(*lockedOut))
		}
		if suspended := prepaidUnits.GetSuspended(); suspended != nil {
			prepaidUnitsAttrs["suspended"] = types.Int64Value(int64(*suspended))
		}
		if warning := prepaidUnits.GetWarning(); warning != nil {
			prepaidUnitsAttrs["warning"] = types.Int64Value(int64(*warning))
		}

		prepaidUnitsObj, _ := types.ObjectValue(getPrepaidUnitsObjectType(), prepaidUnitsAttrs)
		model.PrepaidUnits = prepaidUnitsObj
	} else {
		model.PrepaidUnits = types.ObjectNull(getPrepaidUnitsObjectType())
	}

	// Handle service plans
	if servicePlans := data.GetServicePlans(); servicePlans != nil {
		var servicePlanItems []attr.Value
		for _, plan := range servicePlans {
			servicePlanAttrs := map[string]attr.Value{
				"service_plan_id":     state.UUIDPtrToTypeString(plan.GetServicePlanId()),
				"service_plan_name":   types.StringPointerValue(plan.GetServicePlanName()),
				"provisioning_status": types.StringPointerValue(plan.GetProvisioningStatus()),
				"applies_to":          types.StringPointerValue(plan.GetAppliesTo()),
			}
			servicePlanObj, _ := types.ObjectValue(getServicePlanObjectType(), servicePlanAttrs)
			servicePlanItems = append(servicePlanItems, servicePlanObj)
		}

		servicePlansList, _ := types.ListValue(types.ObjectType{AttrTypes: getServicePlanObjectType()}, servicePlanItems)
		model.ServicePlans = servicePlansList
	} else {
		model.ServicePlans = types.ListNull(types.ObjectType{AttrTypes: getServicePlanObjectType()})
	}

	// Handle subscription IDs
	if subscriptionIds := data.GetSubscriptionIds(); subscriptionIds != nil {
		var subscriptionItems []attr.Value
		for _, subId := range subscriptionIds {
			subscriptionItems = append(subscriptionItems, types.StringValue(subId))
		}

		subscriptionList, _ := types.ListValue(types.StringType, subscriptionItems)
		model.SubscriptionIds = subscriptionList
	} else {
		model.SubscriptionIds = types.ListNull(types.StringType)
	}

	return model
}

// Helper functions for object type definitions
func getSubscribedSkuObjectType() map[string]attr.Type {
	return map[string]attr.Type{
		"id":                types.StringType,
		"account_id":        types.StringType,
		"account_name":      types.StringType,
		"applies_to":        types.StringType,
		"capability_status": types.StringType,
		"consumed_units":    types.Int64Type,
		"sku_id":            types.StringType,
		"sku_part_number":   types.StringType,
		"prepaid_units":     types.ObjectType{AttrTypes: getPrepaidUnitsObjectType()},
		"service_plans":     types.ListType{ElemType: types.ObjectType{AttrTypes: getServicePlanObjectType()}},
		"subscription_ids":  types.ListType{ElemType: types.StringType},
	}
}

func getPrepaidUnitsObjectType() map[string]attr.Type {
	return map[string]attr.Type{
		"enabled":    types.Int64Type,
		"locked_out": types.Int64Type,
		"suspended":  types.Int64Type,
		"warning":    types.Int64Type,
	}
}

func getServicePlanObjectType() map[string]attr.Type {
	return map[string]attr.Type{
		"service_plan_id":     types.StringType,
		"service_plan_name":   types.StringType,
		"provisioning_status": types.StringType,
		"applies_to":          types.StringType,
	}
}
