package teams_calling_policy

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/powershell"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// Create implements resource.Resource Create for TeamsCallingPolicy
func (r *TeamsCallingPolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data TeamsCallingPolicyResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	cmd := []string{"New-CsTeamsCallingPolicy"}
	cmd = append(cmd, fmt.Sprintf("-Identity '%s'", data.ID.ValueString()))
	appendAllCallingPolicyFieldsToCmd(ctx, &data, &cmd)
	psCmd := strings.Join(cmd, " ")
	_, err := powershell.RunPowerShell(psCmd)
	if err != nil {
		resp.Diagnostics.AddError("Error creating Teams Calling Policy", err.Error())
		return
	}
	resp.State.Set(ctx, &data)
}

func (r *TeamsCallingPolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data TeamsCallingPolicyResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	psCmd := fmt.Sprintf("Get-CsTeamsCallingPolicy -Identity '%s' | ConvertTo-Json", data.ID.ValueString())
	output, err := powershell.RunPowerShell(psCmd)
	if err != nil {
		resp.Diagnostics.AddError("Error reading Teams Calling Policy", err.Error())
		return
	}
	var p map[string]any
	err = json.Unmarshal([]byte(output), &p)
	if err != nil {
		resp.Diagnostics.AddError("Error parsing Teams Calling Policy JSON", err.Error())
		return
	}
	mapAllCallingPolicyFieldsFromJson(ctx, p, &data)
	resp.State.Set(ctx, &data)
}

func (r *TeamsCallingPolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data TeamsCallingPolicyResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	cmd := []string{"Set-CsTeamsCallingPolicy"}
	cmd = append(cmd, fmt.Sprintf("-Identity '%s'", data.ID.ValueString()))
	appendAllCallingPolicyFieldsToCmd(ctx, &data, &cmd)
	psCmd := strings.Join(cmd, " ")
	_, err := powershell.RunPowerShell(psCmd)
	if err != nil {
		resp.Diagnostics.AddError("Error updating Teams Calling Policy", err.Error())
		return
	}
	resp.State.Set(ctx, &data)
}

func (r *TeamsCallingPolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data TeamsCallingPolicyResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	psCmd := fmt.Sprintf("Remove-CsTeamsCallingPolicy -Identity '%s' -Confirm:$false", data.ID.ValueString())
	_, err := powershell.RunPowerShell(psCmd)
	if err != nil {
		resp.Diagnostics.AddError("Error deleting Teams Calling Policy", err.Error())
	}
}
