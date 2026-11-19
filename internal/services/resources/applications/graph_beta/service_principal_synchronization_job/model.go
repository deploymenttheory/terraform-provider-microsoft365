package graphBetaServicePrincipalSynchronizationJob

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ServicePrincipalSynchronizationJobResourceModel struct {
	ID                 types.String   `tfsdk:"id"`
	ServicePrincipalID types.String   `tfsdk:"service_principal_id"`
	JobID              types.String   `tfsdk:"job_id"`
	TemplateID         types.String   `tfsdk:"template_id"`
	Enabled            types.Bool     `tfsdk:"enabled"`
	Schedule           types.Object   `tfsdk:"schedule"`
	StatusCode         types.String   `tfsdk:"status_code"`
	Timeouts           timeouts.Value `tfsdk:"timeouts"`
}

var scheduleAttributeTypes = map[string]attr.Type{"expiration": types.StringType, "interval": types.StringType, "state": types.StringType}
