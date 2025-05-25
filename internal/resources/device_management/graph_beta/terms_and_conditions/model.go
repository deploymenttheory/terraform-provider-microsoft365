// resource REF: https://learn.microsoft.com/en-us/graph/api/resources/intune-companyterms-termsandconditions?view=graph-rest-beta
package graphBetaTermsAndConditions

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type TermsAndConditionsResourceModel struct {
	ID                  types.String   `tfsdk:"id"`
	DisplayName         types.String   `tfsdk:"display_name"`
	Description         types.String   `tfsdk:"description"`
	Title               types.String   `tfsdk:"title"`
	BodyText            types.String   `tfsdk:"body_text"`
	AcceptanceStatement types.String   `tfsdk:"acceptance_statement"`
	Version             types.Int32    `tfsdk:"version"`
	RoleScopeTagIds     types.Set      `tfsdk:"role_scope_tag_ids"`
	CreatedDateTime     types.String   `tfsdk:"created_date_time"`
	ModifiedDateTime    types.String   `tfsdk:"modified_date_time"`
	Timeouts            timeouts.Value `tfsdk:"timeouts"`
}
