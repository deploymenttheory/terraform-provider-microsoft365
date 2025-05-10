package graphBetaWindowsRemediationScript

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/constructors"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// constructResource creates a new DeviceHealthScript based on the resource model
func constructResource(ctx context.Context, data *DeviceHealthScriptResourceModel) (graphmodels.DeviceHealthScriptable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))

	requestBody := graphmodels.NewDeviceHealthScript()

	constructors.SetStringProperty(data.DisplayName, requestBody.SetDisplayName)
	constructors.SetStringProperty(data.Description, requestBody.SetDescription)
	constructors.SetStringProperty(data.Publisher, requestBody.SetPublisher)
	constructors.SetBoolProperty(data.RunAs32Bit, requestBody.SetRunAs32Bit)
	constructors.SetBoolProperty(data.EnforceSignatureCheck, requestBody.SetEnforceSignatureCheck)

	if err := constructors.SetEnumProperty(data.RunAsAccount, graphmodels.ParseRunAsAccountType, requestBody.SetRunAsAccount); err != nil {
		return nil, fmt.Errorf("invalid run as account type: %s", err)
	}

	if data.DeviceHealthScriptType.ValueString() != "" {
		if err := constructors.SetEnumProperty(data.DeviceHealthScriptType, graphmodels.ParseDeviceHealthScriptType, requestBody.SetDeviceHealthScriptType); err != nil {
			return nil, fmt.Errorf("invalid device health script type: %s", err)
		}
	}

	constructors.SetBytesProperty(data.DetectionScriptContent, requestBody.SetDetectionScriptContent)
	constructors.SetBytesProperty(data.RemediationScriptContent, requestBody.SetRemediationScriptContent)

	if err := constructors.SetStringSet(ctx, data.RoleScopeTagIds, requestBody.SetRoleScopeTagIds); err != nil {
		return nil, fmt.Errorf("failed to set role scope tags: %s", err)
	}

	var detParams []DeviceHealthScriptParameterModel

	if diags := data.DetectionScriptParameters.ElementsAs(ctx, &detParams, false); diags.HasError() {
		return nil, fmt.Errorf("unable to read detectionScriptParameters: %s", diags)
	}

	if len(detParams) > 0 {
		var graphDet []graphmodels.DeviceHealthScriptParameterable
		for _, p := range detParams {
			gp := graphmodels.NewDeviceHealthScriptParameter()
			constructors.SetStringProperty(p.Name, gp.SetName)
			constructors.SetStringProperty(p.Description, gp.SetDescription)
			constructors.SetBoolProperty(p.IsRequired, gp.SetIsRequired)
			constructors.SetBoolProperty(p.ApplyDefaultValueWhenNotAssigned, gp.SetApplyDefaultValueWhenNotAssigned)
			graphDet = append(graphDet, gp)
		}
		requestBody.SetDetectionScriptParameters(graphDet)
	}

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]interface{}{"error": err.Error()})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s resource", ResourceName))
	return requestBody, nil
}
