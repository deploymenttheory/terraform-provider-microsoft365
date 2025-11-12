package utilityEntraIdSidConverter

import (
	"context"
	"fmt"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Read handles the Read operation for Entra ID SID Converter data source.
func (d *entraIdSidConverterDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state EntraIdSidConverterDataSourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting Read method for: %s", DataSourceName))

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := crud.HandleTimeout(ctx, state.Timeouts.Read, ReadTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	if !state.Sid.IsNull() && !state.Sid.IsUnknown() {
		sid := state.Sid.ValueString()

		tflog.Debug(ctx, fmt.Sprintf("Converting SID to Object ID: %s", sid))

		objectId, err := convertSidToObjectId(sid)
		if err != nil {
			resp.Diagnostics.AddAttributeError(
				path.Root("sid"),
				"SID Conversion Failed",
				fmt.Sprintf("Could not convert SID '%s' to Object ID: %s", sid, err),
			)
			return
		}

		state.ObjectId = types.StringValue(objectId)
		state.Id = types.StringValue(sid)

		tflog.Debug(ctx, fmt.Sprintf("Successfully converted SID to Object ID: %s", objectId))

	} else if !state.ObjectId.IsNull() && !state.ObjectId.IsUnknown() {
		objectId := state.ObjectId.ValueString()

		tflog.Debug(ctx, fmt.Sprintf("Converting Object ID to SID: %s", objectId))

		sid, err := convertObjectIdToSid(objectId)
		if err != nil {
			resp.Diagnostics.AddAttributeError(
				path.Root("object_id"),
				"Object ID Conversion Failed",
				fmt.Sprintf("Could not convert Object ID '%s' to SID: %s", objectId, err),
			)
			return
		}

		state.Sid = types.StringValue(sid)
		state.Id = types.StringValue(objectId)

		tflog.Debug(ctx, fmt.Sprintf("Successfully converted Object ID to SID: %s", sid))

	} else {
		resp.Diagnostics.AddError(
			"Missing Input",
			"Either 'sid' or 'object_id' must be provided.",
		)
		return
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)

	tflog.Debug(ctx, fmt.Sprintf("Finished Read Method: %s", DataSourceName))
}
