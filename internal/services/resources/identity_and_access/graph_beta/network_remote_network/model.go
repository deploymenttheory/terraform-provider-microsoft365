// REF: https://learn.microsoft.com/en-us/graph/api/resources/networkaccess-remotenetwork?view=graph-rest-beta
package graphBetaIdentityAndAccessNetworkRemoteNetwork

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// NetworkRemoteNetworkResourceModel is the Terraform model of a remote network.
type NetworkRemoteNetworkResourceModel struct {
	ID                   types.String   `tfsdk:"id"`
	Name                 types.String   `tfsdk:"name"`
	Region               types.String   `tfsdk:"region"`
	ForwardingProfileIDs types.Set      `tfsdk:"forwarding_profile_ids"`
	Version              types.String   `tfsdk:"version"`
	Timeouts             timeouts.Value `tfsdk:"timeouts"`
}

type remoteNetworkResponse struct {
	ID                 *string                     `json:"id"`
	Name               *string                     `json:"name"`
	Region             *string                     `json:"region"`
	Version            *string                     `json:"version"`
	ForwardingProfiles []forwardingProfileResponse `json:"forwardingProfiles"`
}

type forwardingProfileResponse struct {
	ID *string `json:"id"`
}
