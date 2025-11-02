package graphBetaApplicationsIpApplicationSegment

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// IpApplicationSegmentResourceModel represents the Terraform resource model for On-Premises publishing IP application segments
type IpApplicationSegmentResourceModel struct {
	ID              types.String   `tfsdk:"id"`
	ApplicationID   types.String   `tfsdk:"application_id"`
	DestinationHost types.String   `tfsdk:"destination_host"`
	DestinationType types.String   `tfsdk:"destination_type"`
	Ports           types.Set      `tfsdk:"ports"`
	Protocol        types.String   `tfsdk:"protocol"`
	Timeouts        timeouts.Value `tfsdk:"timeouts"`
}
