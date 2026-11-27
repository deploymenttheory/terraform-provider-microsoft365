package graphBetaTenantInformation

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
)

func (d *TenantInformationDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var object TenantInformationDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}
	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Read, ReadTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()
	var result graphmodels.TenantInformationable
	var err error
	switch {
	case !object.TenantID.IsNull() && !object.TenantID.IsUnknown():
		value := object.TenantID.ValueString()
		result, err = d.client.TenantRelationships().FindTenantInformationByTenantIdWithTenantId(&value).Get(ctx, nil)
	case !object.DomainName.IsNull() && !object.DomainName.IsUnknown():
		value := object.DomainName.ValueString()
		result, err = d.client.TenantRelationships().FindTenantInformationByDomainNameWithDomainName(&value).Get(ctx, nil)
	default:
		resp.Diagnostics.AddError("Missing tenant lookup", "Specify exactly one of tenant_id or domain_name.")
		return
	}
	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationRead, d.ReadPermissions)
		return
	}
	if result == nil || result.GetTenantId() == nil || *result.GetTenantId() == "" {
		resp.Diagnostics.AddError("Invalid tenant information response", "Microsoft Graph returned no tenant identifier for the requested lookup.")
		return
	}
	object.ID = types.StringValue(*result.GetTenantId())
	object.Items = []TenantInformationItem{{
		TenantID:            convert.GraphToFrameworkString(result.GetTenantId()),
		DisplayName:         convert.GraphToFrameworkString(result.GetDisplayName()),
		DefaultDomainName:   convert.GraphToFrameworkString(result.GetDefaultDomainName()),
		FederationBrandName: convert.GraphToFrameworkString(result.GetFederationBrandName()),
	}}
	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
}
