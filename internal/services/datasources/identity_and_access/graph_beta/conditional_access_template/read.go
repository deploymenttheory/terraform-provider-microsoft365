// Package graphBetaConditionalAccessTemplate reads Microsoft Entra conditional access policy templates.
//
// SDK deviation: unknown enum value
//
// Evidence:
//   - SDK: github.com/microsoftgraph/msgraph-beta-sdk-go@v0.168.0, models.ConditionalAccessGrantControl and
//     models.TemplateScenarios. ParseConditionalAccessGrantControl knows block, mfa, compliantDevice, domainJoinedDevice,
//     approvedApplication, compliantApplication, passwordChange, unknownFutureValue and riskRemediation, and returns nil
//     for anything else, so verifiedID comes back as block. ParseTemplateScenarios returns nil for the whole flags value
//     when one flag is unknown, so a template with aiAgents loses its scenarios.
//   - Learn: https://learn.microsoft.com/en-us/graph/api/resources/conditionalaccesstemplate?view=graph-rest-beta
//   - Observed: 2026-11-24, GET /beta/identity/conditionalAccess/templates returns the "Secure account recovery with
//     identity verification (Preview)" template with builtInControls ["verifiedID"] and the "Block high risk agent
//     identities from accessing resources" template with scenarios "aiAgents" (acceptance tests 20 and 21).
//
// The template list is read through custom_requests in raw JSON mode, and templateResponse keeps scenarios and
// builtInControls as strings while the SDK model parses every other field.
//
// Removal condition: models.ConditionalAccessGrantControl includes verifiedID and models.TemplateScenarios includes
// aiAgents, or the SDK keeps unknown enum members instead of dropping them.
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
