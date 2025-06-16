// REF: https://learn.microsoft.com/en-us/graph/api/resources/browsersitelist?view=graph-rest-beta
package graphBetaBrowserSiteList

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type BrowserSiteListResourceModel struct {
	ID                   types.String   `tfsdk:"id"`
	Description          types.String   `tfsdk:"description"`
	DisplayName          types.String   `tfsdk:"display_name"`
	LastModifiedDateTime types.String   `tfsdk:"last_modified_date_time"`
	PublishedDateTime    types.String   `tfsdk:"published_date_time"`
	Revision             types.String   `tfsdk:"revision"`
	Status               types.String   `tfsdk:"status"`
	Timeouts             timeouts.Value `tfsdk:"timeouts"`
}
