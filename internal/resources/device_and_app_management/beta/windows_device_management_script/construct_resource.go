package graphBetaDeviceManagementScript

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/construct"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// Main entry point to construct the intune windows device management script resource for the Terraform provider.
func constructResource(ctx context.Context, data *DeviceManagementScriptResourceModel) (graphmodels.DeviceManagementScriptable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))

	requestBody := graphmodels.NewDeviceManagementScript()

	if !data.DisplayName.IsNull() {
		displayName := data.DisplayName.ValueString()
		requestBody.SetDisplayName(&displayName)
	}

	if !data.Description.IsNull() {
		description := data.Description.ValueString()
		requestBody.SetDescription(&description)
	}

	if !data.ScriptContent.IsNull() {
		encodedContent := base64.StdEncoding.EncodeToString([]byte(data.ScriptContent.ValueString()))
		scriptContent := []byte(encodedContent)
		requestBody.SetScriptContent(scriptContent)
	}

	if !data.RunAsAccount.IsNull() {
		runAsAccountStr := data.RunAsAccount.ValueString()
		var runAsAccount graphmodels.RunAsAccountType
		switch runAsAccountStr {
		case "system":
			runAsAccount = graphmodels.SYSTEM_RUNASACCOUNTTYPE
		case "user":
			runAsAccount = graphmodels.USER_RUNASACCOUNTTYPE
		}
		requestBody.SetRunAsAccount(&runAsAccount)
	}

	if !data.EnforceSignatureCheck.IsNull() {
		enforceSignatureCheck := data.EnforceSignatureCheck.ValueBool()
		requestBody.SetEnforceSignatureCheck(&enforceSignatureCheck)
	}

	if !data.FileName.IsNull() {
		fileName := data.FileName.ValueString()
		requestBody.SetFileName(&fileName)
	}

	if len(data.RoleScopeTagIds) > 0 {
		roleScopeTagIds := make([]string, 0, len(data.RoleScopeTagIds))
		for _, v := range data.RoleScopeTagIds {
			if !v.IsNull() && !v.IsUnknown() {
				roleScopeTagIds = append(roleScopeTagIds, v.ValueString())
			}
		}
		if len(roleScopeTagIds) > 0 {
			requestBody.SetRoleScopeTagIds(roleScopeTagIds)
		}
	}

	if !data.RunAs32Bit.IsNull() {
		runAs32Bit := data.RunAs32Bit.ValueBool()
		requestBody.SetRunAs32Bit(&runAs32Bit)
	}

	if err := construct.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]interface{}{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s resource", ResourceName))

	return requestBody, nil
}
