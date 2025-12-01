// REF: https://learn.microsoft.com/en-us/entra/identity/users/licensing-service-plan-reference
package utilityLicensingServicePlanReference

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/datasource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// licensingServicePlanReferenceDataSourceModel represents the Terraform data source model
type licensingServicePlanReferenceDataSourceModel struct {
	Id                   types.String   `tfsdk:"id"`
	ProductName          types.String   `tfsdk:"product_name"`
	StringId             types.String   `tfsdk:"string_id"`
	Guid                 types.String   `tfsdk:"guid"`
	ServicePlanId        types.String   `tfsdk:"service_plan_id"`
	ServicePlanName      types.String   `tfsdk:"service_plan_name"`
	ServicePlanGuid      types.String   `tfsdk:"service_plan_guid"`
	MatchingProducts     types.List     `tfsdk:"matching_products"`      // List of ProductModel
	MatchingServicePlans types.List     `tfsdk:"matching_service_plans"` // List of ServicePlanModel
	Timeouts             timeouts.Value `tfsdk:"timeouts"`
}

// ProductModel represents a license product/SKU
type ProductModel struct {
	ProductName          types.String `tfsdk:"product_name"`
	StringId             types.String `tfsdk:"string_id"`
	Guid                 types.String `tfsdk:"guid"`
	ServicePlansIncluded types.List   `tfsdk:"service_plans_included"` // List of ServicePlanIncludedModel
}

// ServicePlanIncludedModel represents a service plan included in a product
type ServicePlanIncludedModel struct {
	Id   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
	Guid types.String `tfsdk:"guid"`
}

// ServicePlanModel represents a unique service plan across all products
type ServicePlanModel struct {
	Id             types.String `tfsdk:"id"`
	Name           types.String `tfsdk:"name"`
	Guid           types.String `tfsdk:"guid"`
	IncludedInSkus types.List   `tfsdk:"included_in_skus"` // List of SkuReferenceModel
}

// SkuReferenceModel represents a reference to a SKU that includes a service plan
type SkuReferenceModel struct {
	ProductName types.String `tfsdk:"product_name"`
	StringId    types.String `tfsdk:"string_id"`
	Guid        types.String `tfsdk:"guid"`
}

// LicenseData represents the structure of the JSON data file
type LicenseData struct {
	ProductName          string        `json:"product_name"`
	StringId             string        `json:"string_id"`
	Guid                 string        `json:"guid"`
	ServicePlansIncluded []ServicePlan `json:"service_plans_included"`
}

// ServicePlan represents a service plan in the JSON data
type ServicePlan struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Guid string `json:"guid"`
}
