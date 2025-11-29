package utilityLicensingServicePlanReference

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// searchProductsByName searches for products by name (partial match, case-insensitive)
func searchProductsByName(ctx context.Context, data []LicenseData, searchTerm string, model *licensingServicePlanReferenceDataSourceModel, diags *diag.Diagnostics) {
	var matchingProducts []LicenseData

	for _, product := range data {
		if containsIgnoreCase(product.ProductName, searchTerm) {
			matchingProducts = append(matchingProducts, product)
		}
	}

	tflog.Debug(ctx, "Product name search completed", map[string]interface{}{
		"search_term":   searchTerm,
		"matches_found": len(matchingProducts),
	})

	if len(matchingProducts) == 0 {
		diags.AddWarning(
			"No Products Found",
			fmt.Sprintf("No products found matching the name: %s", searchTerm),
		)
	}

	model.MatchingProducts = mapProductsToFrameworkList(ctx, matchingProducts, diags)
	model.MatchingServicePlans = types.ListNull(types.ObjectType{AttrTypes: servicePlanAttrTypes()})
}

// searchProductByStringId searches for a product by string ID (exact match, case-insensitive)
func searchProductByStringId(ctx context.Context, data []LicenseData, stringId string, model *licensingServicePlanReferenceDataSourceModel, diags *diag.Diagnostics) {
	var matchingProducts []LicenseData

	for _, product := range data {
		if equalsIgnoreCase(product.StringId, stringId) {
			matchingProducts = append(matchingProducts, product)
		}
	}

	tflog.Debug(ctx, "Product string_id search completed", map[string]interface{}{
		"string_id":     stringId,
		"matches_found": len(matchingProducts),
	})

	if len(matchingProducts) == 0 {
		diags.AddError(
			"Product Not Found",
			fmt.Sprintf("No product found with string_id: %s", stringId),
		)
		return
	}

	model.MatchingProducts = mapProductsToFrameworkList(ctx, matchingProducts, diags)
	model.MatchingServicePlans = types.ListNull(types.ObjectType{AttrTypes: servicePlanAttrTypes()})
}

// searchProductByGuid searches for a product by GUID (exact match)
func searchProductByGuid(ctx context.Context, data []LicenseData, guid string, model *licensingServicePlanReferenceDataSourceModel, diags *diag.Diagnostics) {
	var matchingProducts []LicenseData

	for _, product := range data {
		if equalsIgnoreCase(product.Guid, guid) {
			matchingProducts = append(matchingProducts, product)
		}
	}

	tflog.Debug(ctx, "Product GUID search completed", map[string]interface{}{
		"guid":          guid,
		"matches_found": len(matchingProducts),
	})

	if len(matchingProducts) == 0 {
		diags.AddError(
			"Product Not Found",
			fmt.Sprintf("No product found with GUID: %s", guid),
		)
		return
	}

	model.MatchingProducts = mapProductsToFrameworkList(ctx, matchingProducts, diags)
	model.MatchingServicePlans = types.ListNull(types.ObjectType{AttrTypes: servicePlanAttrTypes()})
}

// searchServicePlanById searches for service plans by ID (partial match, case-insensitive)
func searchServicePlanById(ctx context.Context, data []LicenseData, searchTerm string, model *licensingServicePlanReferenceDataSourceModel, diags *diag.Diagnostics) {
	servicePlanMap := buildServicePlanMap(data)

	var matchingPlans []servicePlanWithSkus

	for _, plan := range servicePlanMap {
		if containsIgnoreCase(plan.Id, searchTerm) {
			matchingPlans = append(matchingPlans, plan)
		}
	}

	tflog.Debug(ctx, "Service plan ID search completed", map[string]interface{}{
		"search_term":   searchTerm,
		"matches_found": len(matchingPlans),
	})

	if len(matchingPlans) == 0 {
		diags.AddWarning(
			"No Service Plans Found",
			fmt.Sprintf("No service plans found matching the ID: %s", searchTerm),
		)
	}

	model.MatchingServicePlans = mapServicePlansToFrameworkList(ctx, matchingPlans, diags)
	model.MatchingProducts = types.ListNull(types.ObjectType{AttrTypes: productAttrTypes()})
}

// searchServicePlanByName searches for service plans by name (partial match, case-insensitive)
func searchServicePlanByName(ctx context.Context, data []LicenseData, searchTerm string, model *licensingServicePlanReferenceDataSourceModel, diags *diag.Diagnostics) {
	servicePlanMap := buildServicePlanMap(data)

	var matchingPlans []servicePlanWithSkus

	for _, plan := range servicePlanMap {
		if containsIgnoreCase(plan.Name, searchTerm) {
			matchingPlans = append(matchingPlans, plan)
		}
	}

	tflog.Debug(ctx, "Service plan name search completed", map[string]interface{}{
		"search_term":   searchTerm,
		"matches_found": len(matchingPlans),
	})

	if len(matchingPlans) == 0 {
		diags.AddWarning(
			"No Service Plans Found",
			fmt.Sprintf("No service plans found matching the name: %s", searchTerm),
		)
	}

	model.MatchingServicePlans = mapServicePlansToFrameworkList(ctx, matchingPlans, diags)
	model.MatchingProducts = types.ListNull(types.ObjectType{AttrTypes: productAttrTypes()})
}

// searchServicePlanByGuid searches for a service plan by GUID (exact match)
func searchServicePlanByGuid(ctx context.Context, data []LicenseData, guid string, model *licensingServicePlanReferenceDataSourceModel, diags *diag.Diagnostics) {
	servicePlanMap := buildServicePlanMap(data)

	var matchingPlans []servicePlanWithSkus

	for _, plan := range servicePlanMap {
		if equalsIgnoreCase(plan.Guid, guid) {
			matchingPlans = append(matchingPlans, plan)
		}
	}

	tflog.Debug(ctx, "Service plan GUID search completed", map[string]interface{}{
		"guid":          guid,
		"matches_found": len(matchingPlans),
	})

	if len(matchingPlans) == 0 {
		diags.AddError(
			"Service Plan Not Found",
			fmt.Sprintf("No service plan found with GUID: %s", guid),
		)
		return
	}

	model.MatchingServicePlans = mapServicePlansToFrameworkList(ctx, matchingPlans, diags)
	model.MatchingProducts = types.ListNull(types.ObjectType{AttrTypes: productAttrTypes()})
}

// servicePlanWithSkus represents a service plan with its associated SKUs
type servicePlanWithSkus struct {
	Id             string
	Name           string
	Guid           string
	IncludedInSkus []skuReference
}

// skuReference represents a reference to a SKU
type skuReference struct {
	ProductName string
	StringId    string
	Guid        string
}

// buildServicePlanMap creates a map of unique service plans with their associated SKUs
func buildServicePlanMap(data []LicenseData) map[string]servicePlanWithSkus {
	planMap := make(map[string]servicePlanWithSkus)

	for _, product := range data {
		for _, plan := range product.ServicePlansIncluded {
			if plan.Guid == "" {
				continue
			}

			existing, found := planMap[plan.Guid]
			if !found {
				existing = servicePlanWithSkus{
					Id:             plan.Id,
					Name:           plan.Name,
					Guid:           plan.Guid,
					IncludedInSkus: []skuReference{},
				}
			}

			// Add this SKU to the list
			existing.IncludedInSkus = append(existing.IncludedInSkus, skuReference{
				ProductName: product.ProductName,
				StringId:    product.StringId,
				Guid:        product.Guid,
			})

			planMap[plan.Guid] = existing
		}
	}

	return planMap
}

// mapProductsToFrameworkList converts products to a Framework list
func mapProductsToFrameworkList(ctx context.Context, products []LicenseData, diags *diag.Diagnostics) types.List {
	if len(products) == 0 {
		return types.ListNull(types.ObjectType{AttrTypes: productAttrTypes()})
	}

	var productObjects []attr.Value

	for _, product := range products {
		// Map service plans included
		var servicePlanObjects []attr.Value
		for _, plan := range product.ServicePlansIncluded {
			planObj, diag := types.ObjectValue(
				servicePlanIncludedAttrTypes(),
				map[string]attr.Value{
					"id":   types.StringValue(plan.Id),
					"name": types.StringValue(plan.Name),
					"guid": types.StringValue(plan.Guid),
				},
			)
			diags.Append(diag...)
			servicePlanObjects = append(servicePlanObjects, planObj)
		}

		servicePlansList, diag := types.ListValue(
			types.ObjectType{AttrTypes: servicePlanIncludedAttrTypes()},
			servicePlanObjects,
		)
		diags.Append(diag...)

		productObj, diag := types.ObjectValue(
			productAttrTypes(),
			map[string]attr.Value{
				"product_name":           types.StringValue(product.ProductName),
				"string_id":              types.StringValue(product.StringId),
				"guid":                   types.StringValue(product.Guid),
				"service_plans_included": servicePlansList,
			},
		)
		diags.Append(diag...)
		productObjects = append(productObjects, productObj)
	}

	list, diag := types.ListValue(
		types.ObjectType{AttrTypes: productAttrTypes()},
		productObjects,
	)
	diags.Append(diag...)

	return list
}

// mapServicePlansToFrameworkList converts service plans to a Framework list
func mapServicePlansToFrameworkList(ctx context.Context, plans []servicePlanWithSkus, diags *diag.Diagnostics) types.List {
	if len(plans) == 0 {
		return types.ListNull(types.ObjectType{AttrTypes: servicePlanAttrTypes()})
	}

	var planObjects []attr.Value

	for _, plan := range plans {
		// Map included SKUs
		var skuObjects []attr.Value
		for _, sku := range plan.IncludedInSkus {
			skuObj, diag := types.ObjectValue(
				skuReferenceAttrTypes(),
				map[string]attr.Value{
					"product_name": types.StringValue(sku.ProductName),
					"string_id":    types.StringValue(sku.StringId),
					"guid":         types.StringValue(sku.Guid),
				},
			)
			diags.Append(diag...)
			skuObjects = append(skuObjects, skuObj)
		}

		skusList, diag := types.ListValue(
			types.ObjectType{AttrTypes: skuReferenceAttrTypes()},
			skuObjects,
		)
		diags.Append(diag...)

		planObj, diag := types.ObjectValue(
			servicePlanAttrTypes(),
			map[string]attr.Value{
				"id":               types.StringValue(plan.Id),
				"name":             types.StringValue(plan.Name),
				"guid":             types.StringValue(plan.Guid),
				"included_in_skus": skusList,
			},
		)
		diags.Append(diag...)
		planObjects = append(planObjects, planObj)
	}

	list, diag := types.ListValue(
		types.ObjectType{AttrTypes: servicePlanAttrTypes()},
		planObjects,
	)
	diags.Append(diag...)

	return list
}

// productAttrTypes returns attribute types for ProductModel
func productAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"product_name":           types.StringType,
		"string_id":              types.StringType,
		"guid":                   types.StringType,
		"service_plans_included": types.ListType{ElemType: types.ObjectType{AttrTypes: servicePlanIncludedAttrTypes()}},
	}
}

// servicePlanIncludedAttrTypes returns attribute types for ServicePlanIncludedModel
func servicePlanIncludedAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":   types.StringType,
		"name": types.StringType,
		"guid": types.StringType,
	}
}

// servicePlanAttrTypes returns attribute types for ServicePlanModel
func servicePlanAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":               types.StringType,
		"name":             types.StringType,
		"guid":             types.StringType,
		"included_in_skus": types.ListType{ElemType: types.ObjectType{AttrTypes: skuReferenceAttrTypes()}},
	}
}

// skuReferenceAttrTypes returns attribute types for SkuReferenceModel
func skuReferenceAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"product_name": types.StringType,
		"string_id":    types.StringType,
		"guid":         types.StringType,
	}
}
