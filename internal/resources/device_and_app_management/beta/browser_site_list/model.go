// REF: https://learn.microsoft.com/en-us/graph/api/resources/browsersitelist?view=graph-rest-beta
package graphbetabrowsersite

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type BrowserSiteListResourceModel struct {
	Description types.String `tfsdk:"description"`
	DisplayName types.String `tfsdk:"display_name"`
	ID          types.String `tfsdk:"id"`
	//LastModifiedBy       sharedmodels.IdentitySetModel `tfsdk:"last_modified_by"`
	LastModifiedDateTime types.String `tfsdk:"last_modified_date_time"`
	//PublishedBy          sharedmodels.IdentitySetModel `tfsdk:"published_by"`
	PublishedDateTime types.String   `tfsdk:"published_date_time"`
	Revision          types.String   `tfsdk:"revision"`
	Status            types.String   `tfsdk:"status"`
	Timeouts          timeouts.Value `tfsdk:"timeouts"`
}
