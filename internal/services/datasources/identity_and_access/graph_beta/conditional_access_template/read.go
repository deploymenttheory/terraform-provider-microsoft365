package graphBetaConditionalAccessTemplate

import (
	"context"
	"fmt"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func (d *ConditionalAccessTemplateDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var object ConditionalAccessTemplateDataSourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting Read method for: %s", DataSourceName))

	resp.Diagnostics.Append(req.Config.Get(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	name := object.Name.ValueString()
	templateID := object.TemplateID.ValueString()

	tflog.Debug(ctx, fmt.Sprintf("Reading %s with name: %s, template_id: %s", DataSourceName, name, templateID))

	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Read, ReadTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	template, ok := d.validateRequest(ctx, templateID, name, resp)
	if !ok {
		return
	}

	matchedTemplate := MapRemoteStateToDataSource(ctx, template)

	object.ID = types.StringValue(fmt.Sprintf("conditional-access-template-%s-%d", matchedTemplate.TemplateID.ValueString(), time.Now().Unix()))
	object.TemplateID = matchedTemplate.TemplateID
	object.Name = matchedTemplate.Name
	object.Description = matchedTemplate.Description
	object.Scenarios = matchedTemplate.Scenarios
	object.Details = matchedTemplate.Details

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Datasource Read Method: %s", DataSourceName))
}
