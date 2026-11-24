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
	Protocol            types.Set      `tfsdk:"protocol"`
	Timeouts            timeouts.Value `tfsdk:"timeouts"`
}

// ipApplicationSegmentRequest is the ipApplicationSegment create and update body, in the wire shape the
// application-scoped beta endpoint accepts.
type ipApplicationSegmentRequest struct {
	DestinationHost string `json:"destinationHost"`
	DestinationType string `json:"destinationType"`
	// Learn marks "port" as deprecated / DO NOT USE, but its create example
	// still includes "port": 0. The real beta endpoint also returns "port": 0
	// together with "ports": ["443-443"] for segments created by Terraform.
	// Include the scalar field to match the endpoint's accepted wire shape while
	// using "ports" for the actual range values in Terraform state.
	Port     int32    `json:"port"`
	Ports    []string `json:"ports"`
	Protocol string   `json:"protocol"`
}

// ipApplicationSegmentResponse holds the ipApplicationSegment fields this resource manages. destinationType is a
// string, not the SDK enum, because the endpoint returns "ip" for a single IP address segment.
//
// Observed response after creating a segment with Terraform:
//
//	{
//	  "@odata.context": ".../applications('{applicationObjectId}')/onPremisesPublishing/segmentsConfiguration/microsoft.graph.ipSegmentConfiguration/applicationSegments/$entity",
//	  "action": "tunnel",
//	  "destinationHost": "10.10.10.10",
//	  "destinationType": "ip",
//	  "exclusions": null,
//	  "id": "{ipApplicationSegmentId}",
//	  "inclusions": null,
//	  "port": 0,
//	  "ports": ["443-443"],
//	  "protocol": "tcp"
//	}
type ipApplicationSegmentResponse struct {
	ID              *string  `json:"id"`
	DestinationHost *string  `json:"destinationHost"`
	DestinationType *string  `json:"destinationType"`
	Ports           []string `json:"ports"`
	Protocol        *string  `json:"protocol"`
}
