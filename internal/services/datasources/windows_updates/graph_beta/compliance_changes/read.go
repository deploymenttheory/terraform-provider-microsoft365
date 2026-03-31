package graphBetaWindowsUpdatesComplianceChanges

import (
	"context"
	"fmt"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func (d *ComplianceChangesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config ComplianceChangesDataSourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting Read method for: %s", DataSourceName))

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var object ComplianceChangesDataSourceModel
	object.UpdatePolicyId = config.UpdatePolicyId
	object.Timeouts = config.Timeouts

	updatePolicyId := object.UpdatePolicyId.ValueString()
	tflog.Debug(ctx, fmt.Sprintf("Reading %s for update_policy_id: %s", DataSourceName, updatePolicyId))

	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Read, ReadTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	respList, err := d.client.
		Admin().
		Windows().
		Updates().
		UpdatePolicies().
		ByUpdatePolicyId(updatePolicyId).
		ComplianceChanges().
		Get(ctx, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationRead, d.ReadPermissions)
		return
	}

	var complianceChanges []ComplianceChange

	if respList != nil && respList.GetValue() != nil {
		for _, change := range respList.GetValue() {
			complianceChanges = append(complianceChanges, MapRemoteStateToDataSource(ctx, change))
		}
	}

	object.ComplianceChanges = complianceChanges

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Datasource Read Method: %s, found %d compliance changes", DataSourceName, len(complianceChanges)))
}
