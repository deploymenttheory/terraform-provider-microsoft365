// REF: https://learn.microsoft.com/en-us/graph/api/resources/ipapplicationsegment?view=graph-rest-beta
package graphBetaApplicationsOnPremisesIpApplicationSegment

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// OnPremisesIpApplicationSegmentResourceModel represents the Terraform resource model for On-Premises publishing IP application segments
type OnPremisesIpApplicationSegmentResourceModel struct {
	ID                  types.String   `tfsdk:"id"`
	ApplicationObjectID types.String   `tfsdk:"application_object_id"`
	DestinationHost     types.String   `tfsdk:"destination_host"`
	DestinationType     types.String   `tfsdk:"destination_type"`
	Ports               types.Set      `tfsdk:"ports"`
	Protocol            types.String   `tfsdk:"protocol"`
	Timeouts            timeouts.Value `tfsdk:"timeouts"`
}
