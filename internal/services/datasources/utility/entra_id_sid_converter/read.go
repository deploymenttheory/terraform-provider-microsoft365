package entra_id_sid_converter

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func (d *entraIdSidConverterDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "Reading Entra ID SID Converter data source")

	var state EntraIdSidConverterDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

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
}
