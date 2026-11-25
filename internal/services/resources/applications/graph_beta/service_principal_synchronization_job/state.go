package graphBetaServicePrincipalSynchronizationJob

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
)

// MapRemoteStateToTerraform maps the remote Microsoft Graph response to Terraform state.
func MapRemoteStateToTerraform(data *ServicePrincipalSynchronizationJobResourceModel, remote graphmodels.SynchronizationJobable) error {
	if remote == nil || remote.GetId() == nil || *remote.GetId() == "" || remote.GetTemplateId() == nil || remote.GetSchedule() == nil || remote.GetSchedule().GetState() == nil {
		return errJobResponse
	}
	data.JobID = convert.GraphToFrameworkString(remote.GetId())
	data.ID = types.StringValue(data.ServicePrincipalID.ValueString() + "/" + data.JobID.ValueString())
	data.TemplateID = convert.GraphToFrameworkString(remote.GetTemplateId())
	data.Enabled = types.BoolValue(*remote.GetSchedule().GetState() == graphmodels.ACTIVE_SYNCHRONIZATIONSCHEDULESTATE)
	data.Schedule = types.ObjectValueMust(scheduleAttributeTypes, map[string]attr.Value{
		"expiration": convert.GraphToFrameworkTime(remote.GetSchedule().GetExpiration()),
		"interval":   convert.GraphToFrameworkISODuration(remote.GetSchedule().GetInterval()),
		"state":      convert.GraphToFrameworkEnum(remote.GetSchedule().GetState()),
	})
	data.StatusCode = types.StringNull()
	if remote.GetStatus() != nil {
		data.StatusCode = convert.GraphToFrameworkEnum(remote.GetStatus().GetCode())
	}
	return nil
}
