package graphBetaNetworkContentPolicy

import (
	"errors"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// contentPoliciesEndpoint is the portal-backed Graph collection of content (file) policies.
const contentPoliciesEndpoint = "networkaccess/filePolicies"

// errEmptyUpdateResponse reports a PATCH that succeeded without returning the updated policy.
var errEmptyUpdateResponse = errors.New("update content policy returned nil response")

// NetworkContentPolicyResourceModel represents a Global Secure Access content policy.
type NetworkContentPolicyResourceModel struct {
	ID            types.String   `tfsdk:"id"`
	Name          types.String   `tfsdk:"name"`
	Description   types.String   `tfsdk:"description"`
	DefaultAction types.String   `tfsdk:"default_action"`
	Timeouts      timeouts.Value `tfsdk:"timeouts"`
}

// contentPolicyResponse is the JSON shape of a content (file) policy returned by Graph. The SDK has no model for it (see
// the SDK deviation in crud.go).
type contentPolicyResponse struct {
	ID          *string `json:"id"`
	Name        *string `json:"name"`
	Description *string `json:"description"`
	Settings    *struct {
		DefaultAction *string `json:"defaultAction"`
	} `json:"settings"`
}
