package graphBetaNetworkTLSInspectionPolicy

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// NetworkTLSInspectionPolicyResourceModel represents a Global Secure Access TLS inspection policy.
type NetworkTLSInspectionPolicyResourceModel struct {
	ID                   types.String   `tfsdk:"id"`
	Name                 types.String   `tfsdk:"name"`
	Description          types.String   `tfsdk:"description"`
	DefaultAction        types.String   `tfsdk:"default_action"`
	Version              types.String   `tfsdk:"version"`
	LastModifiedDateTime types.String   `tfsdk:"last_modified_date_time"`
	Timeouts             timeouts.Value `tfsdk:"timeouts"`
}
