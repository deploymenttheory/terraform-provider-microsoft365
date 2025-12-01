package utilityLicensingServicePlanReference

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

//go:embed data/licensing_service_plan_reference.json
var licensingDataJSON []byte

var cachedLicensingData []LicenseData

// loadLicensingData loads and caches the licensing data
func loadLicensingData() ([]LicenseData, error) {
	if cachedLicensingData != nil {
		return cachedLicensingData, nil
	}

	var data []LicenseData
	if err := json.Unmarshal(licensingDataJSON, &data); err != nil {
		return nil, fmt.Errorf("failed to parse licensing data: %w", err)
	}

	cachedLicensingData = data
	return data, nil
}

// Read implements the datasource.DataSource interface
func (d *licensingServicePlanReferenceDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config licensingServicePlanReferenceDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	licensingData, err := loadLicensingData()
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to Load Licensing Data",
			fmt.Sprintf("Could not load embedded licensing data: %s", err.Error()),
		)
		return
	}

	tflog.Debug(ctx, "Loaded licensing data", map[string]any{
		"product_count": len(licensingData),
	})

	// Determine search mode and execute appropriate search
	if !config.ProductName.IsNull() {
		searchProductsByName(ctx, licensingData, config.ProductName.ValueString(), &config, &resp.Diagnostics)
	} else if !config.StringId.IsNull() {
		searchProductByStringId(ctx, licensingData, config.StringId.ValueString(), &config, &resp.Diagnostics)
	} else if !config.Guid.IsNull() {
		searchProductByGuid(ctx, licensingData, config.Guid.ValueString(), &config, &resp.Diagnostics)
	} else if !config.ServicePlanId.IsNull() {
		searchServicePlanById(ctx, licensingData, config.ServicePlanId.ValueString(), &config, &resp.Diagnostics)
	} else if !config.ServicePlanName.IsNull() {
		searchServicePlanByName(ctx, licensingData, config.ServicePlanName.ValueString(), &config, &resp.Diagnostics)
	} else if !config.ServicePlanGuid.IsNull() {
		searchServicePlanByGuid(ctx, licensingData, config.ServicePlanGuid.ValueString(), &config, &resp.Diagnostics)
	} else {
		resp.Diagnostics.AddError(
			"Missing Search Parameter",
			"At least one search parameter must be specified: product_name, string_id, guid, service_plan_id, service_plan_name, or service_plan_guid",
		)
		return
	}

	if resp.Diagnostics.HasError() {
		return
	}

	config.Id = config.generateId()

	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}

// generateId creates a unique ID for the data source operation
func (m *licensingServicePlanReferenceDataSourceModel) generateId() types.String {
	if !m.ProductName.IsNull() {
		return types.StringValue(fmt.Sprintf("product_name:%s", m.ProductName.ValueString()))
	} else if !m.StringId.IsNull() {
		return types.StringValue(fmt.Sprintf("string_id:%s", m.StringId.ValueString()))
	} else if !m.Guid.IsNull() {
		return types.StringValue(fmt.Sprintf("guid:%s", m.Guid.ValueString()))
	} else if !m.ServicePlanId.IsNull() {
		return types.StringValue(fmt.Sprintf("service_plan_id:%s", m.ServicePlanId.ValueString()))
	} else if !m.ServicePlanName.IsNull() {
		return types.StringValue(fmt.Sprintf("service_plan_name:%s", m.ServicePlanName.ValueString()))
	} else if !m.ServicePlanGuid.IsNull() {
		return types.StringValue(fmt.Sprintf("service_plan_guid:%s", m.ServicePlanGuid.ValueString()))
	}
	return types.StringValue("unknown")
}

// containsIgnoreCase performs case-insensitive substring search
func containsIgnoreCase(haystack, needle string) bool {
	return strings.Contains(strings.ToLower(haystack), strings.ToLower(needle))
}

// equalsIgnoreCase performs case-insensitive equality check
func equalsIgnoreCase(a, b string) bool {
	return strings.EqualFold(a, b)
}
