package graphBetaTenantInformation

import (
	"context"
	"fmt"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func (d *TenantInformationDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var object TenantInformationDataSourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting Read method for: %s", DataSourceName))

	resp.Diagnostics.Append(req.Config.Get(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	filterType := object.FilterType.ValueString()
	filterValue := object.FilterValue.ValueString()

	tflog.Debug(ctx, fmt.Sprintf("Reading %s with filter_type: %s, filter_value: %s", DataSourceName, filterType, filterValue))

	if filterValue == "" {
		resp.Diagnostics.AddError(
			"Missing Required Parameter",
			"filter_value must be provided",
		)
		return
	}

	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Read, ReadTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	var tenantInfo TenantInformationDataSourceModel

	switch filterType {
	case "tenant_id":
		respItem, err := d.client.
			TenantRelationships().
			FindTenantInformationByTenantIdWithTenantId(&filterValue).
			Get(ctx, nil)

		if err != nil {
			errors.HandleKiotaGraphError(ctx, err, resp, "Read", d.ReadPermissions)
			return
		}

		tenantInfo = MapRemoteStateToDataSource(ctx, respItem)

	case "domain_name":
		respItem, err := d.client.
			TenantRelationships().
			FindTenantInformationByDomainNameWithDomainName(&filterValue).
			Get(ctx, nil)

		if err != nil {
			errors.HandleKiotaGraphError(ctx, err, resp, "Read", d.ReadPermissions)
			return
		}

		tenantInfo = MapRemoteStateToDataSource(ctx, respItem)

	default:
		resp.Diagnostics.AddError(
			"Invalid Filter Type",
			fmt.Sprintf("filter_type must be either 'tenant_id' or 'domain_name', got: %s", filterType),
		)
		return
	}

	object.ID = types.StringValue(fmt.Sprintf("tenant-information-%s-%d", filterValue, time.Now().Unix()))
	object.TenantID = tenantInfo.TenantID
	object.DisplayName = tenantInfo.DisplayName
	object.DefaultDomainName = tenantInfo.DefaultDomainName
	object.FederationBrandName = tenantInfo.FederationBrandName

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Datasource Read Method: %s", DataSourceName))
}
