package utilityWindowsRemediationScriptRegistryKeyGenerator

import (
	"context"
	"fmt"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func (d *WindowsRemediationScriptRegistryKeyGeneratorDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, fmt.Sprintf("Starting Read method for: %s", DataSourceName))

	var data WindowsRemediationScriptRegistryKeyGeneratorDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := crud.HandleTimeout(ctx, data.Timeouts.Read, ReadTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	tflog.Debug(ctx, fmt.Sprintf("Generating scripts for context: %s, registry path: %s",
		data.Context.ValueString(), data.RegistryKeyPath.ValueString()))

	detectionScript, remediationScript, err := GenerateScripts(
		data.Context.ValueString(),
		data.RegistryKeyPath.ValueString(),
		data.ValueName.ValueString(),
		data.ValueType.ValueString(),
		data.ValueData.ValueString(),
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Script Generation Failed",
			fmt.Sprintf("Unable to generate registry scripts: %s", err.Error()),
		)
		return
	}

	data.ID = data.Context
	data.DetectionScript = detectionScript
	data.RemediationScript = remediationScript

	tflog.Debug(ctx, "Successfully generated detection and remediation scripts")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	tflog.Debug(ctx, fmt.Sprintf("Finished Read Method: %s", DataSourceName))
}
