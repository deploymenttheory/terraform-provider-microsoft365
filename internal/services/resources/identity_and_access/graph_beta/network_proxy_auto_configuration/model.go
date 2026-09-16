package graphBetaNetworkProxyAutoConfiguration

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// NetworkProxyAutoConfigurationResourceModel holds the managed configuration.
type NetworkProxyAutoConfigurationResourceModel struct {
	ID                   types.String   `tfsdk:"id"`
	Name                 types.String   `tfsdk:"name"`
	Content              types.String   `tfsdk:"content"`
	IsEnabled            types.Bool     `tfsdk:"is_enabled"`
	CreatedDateTime      types.String   `tfsdk:"created_date_time"`
	LastModifiedDateTime types.String   `tfsdk:"last_modified_date_time"`
	Timeouts             timeouts.Value `tfsdk:"timeouts"`
}

type pacResponse struct {
	ID                   *string `json:"id"`
	Name                 *string `json:"name"`
	Content              *string `json:"content"`
	IsEnabled            *bool   `json:"isEnabled"`
	CreatedDateTime      *string `json:"createdDateTime"`
	LastModifiedDateTime *string `json:"lastModifiedDateTime"`
}
