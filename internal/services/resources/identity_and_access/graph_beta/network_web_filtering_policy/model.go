package graphBetaNetworkWebFilteringPolicy

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// NetworkWebFilteringPolicyResourceModel represents the Terraform schema
// for Global Secure Access web filtering policies.
type NetworkWebFilteringPolicyResourceModel struct {
	ID            types.String   `tfsdk:"id"`
	Name          types.String   `tfsdk:"name"`
	Description   types.String   `tfsdk:"description"`
	DefaultAction types.String   `tfsdk:"default_action"`
	Timeouts      timeouts.Value `tfsdk:"timeouts"`
}

// webFilteringPolicy is the raw JSON web filtering policy sent to and returned by /networkaccess/webFilteringPolicies.
type webFilteringPolicy struct {
	ID          *string                     `json:"id,omitempty"`
	Name        *string                     `json:"name,omitempty"`
	Description *string                     `json:"description,omitempty"`
	Settings    *webFilteringPolicySettings `json:"settings,omitempty"`
	// PolicyRules is sent as an empty array on create only; the policy rule resource manages the rules.
	PolicyRules *[]struct{} `json:"policyRules,omitempty"`
}

// webFilteringPolicySettings holds the policy's default action.
type webFilteringPolicySettings struct {
	DefaultAction *webFilteringDefaultAction `json:"defaultAction,omitempty"`
}

// webFilteringDefaultAction is the default action, identified only by its OData type.
type webFilteringDefaultAction struct {
	ODataType string `json:"@odata.type"` //nolint:tagliatelle // Graph requires the literal OData discriminator key.
}
