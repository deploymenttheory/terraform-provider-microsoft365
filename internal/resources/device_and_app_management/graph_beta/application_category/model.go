// REF: https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-mobileappcategory?view=graph-rest-beta
package graphBetaApplicationCategory

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ApplicationCategoryResourceModel struct {
	ID                   types.String   `tfsdk:"id"`
	DisplayName          types.String   `tfsdk:"display_name"`
	LastModifiedDateTime types.String   `tfsdk:"last_modified_date_time"`
	Timeouts             timeouts.Value `tfsdk:"timeouts"`
}
