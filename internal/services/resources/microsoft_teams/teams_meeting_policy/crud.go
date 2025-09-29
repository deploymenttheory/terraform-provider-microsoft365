package teamsMeetingPolicy

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/powershell"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// Create implements resource.Resource Create for TeamsMeetingPolicy
func (r *TeamsMeetingPolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data TeamsMeetingPolicyResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	cmd := []string{"Set-CsTeamsMeetingPolicy"}
	cmd = append(cmd, fmt.Sprintf("-Identity '%s'", data.XdsIdentity.ValueString()))
	appendAllPolicyFieldsToCmd(ctx, &data, &cmd)
	psCmd := strings.Join(cmd, " ")
	_, err := powershell.RunPowerShell(psCmd)
	if err != nil {
		resp.Diagnostics.AddError("Error creating Teams Meeting Policy", err.Error())
		return
	}
	resp.State.Set(ctx, &data)
}

func (r *TeamsMeetingPolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data TeamsMeetingPolicyResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	psCmd := fmt.Sprintf("Get-CsTeamsMeetingPolicy -Identity '%s' | ConvertTo-Json", data.XdsIdentity.ValueString())
	output, err := powershell.RunPowerShell(psCmd)
	if err != nil {
		resp.Diagnostics.AddError("Error reading Teams Meeting Policy", err.Error())
		return
	}
	var p map[string]any
	err = json.Unmarshal([]byte(output), &p)
	if err != nil {
		resp.Diagnostics.AddError("Error parsing Teams Meeting Policy JSON", err.Error())
		return
	}
	mapAllPolicyFieldsFromJson(ctx, p, &data)
	resp.State.Set(ctx, &data)
}

func (r *TeamsMeetingPolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data TeamsMeetingPolicyResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	cmd := []string{"Set-CsTeamsMeetingPolicy"}
	cmd = append(cmd, fmt.Sprintf("-Identity '%s'", data.XdsIdentity.ValueString()))
	appendAllPolicyFieldsToCmd(ctx, &data, &cmd)
	psCmd := strings.Join(cmd, " ")
	_, err := powershell.RunPowerShell(psCmd)
	if err != nil {
		resp.Diagnostics.AddError("Error updating Teams Meeting Policy", err.Error())
		return
	}
	resp.State.Set(ctx, &data)
}

func (r *TeamsMeetingPolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data TeamsMeetingPolicyResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	psCmd := fmt.Sprintf("Remove-CsTeamsMeetingPolicy -Identity '%s' -Confirm:$false", data.XdsIdentity.ValueString())
	_, err := powershell.RunPowerShell(psCmd)
	if err != nil {
		resp.Diagnostics.AddError("Error deleting Teams Meeting Policy", err.Error())
	}
}
