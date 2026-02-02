// REF: https://learn.microsoft.com/en-us/graph/api/application-addpassword?view=graph-rest-beta
// REF: https://learn.microsoft.com/en-us/graph/api/application-removepassword?view=graph-rest-beta
package graphBetaApplicationPasswordCredential

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ApplicationPasswordCredentialResourceModel describes the resource data model.
type ApplicationPasswordCredentialResourceModel struct {
	// Required inputs
	ApplicationID types.String `tfsdk:"application_id"`
	DisplayName   types.String `tfsdk:"display_name"`

	// Optional inputs
	StartDateTime types.String `tfsdk:"start_date_time"`
	EndDateTime   types.String `tfsdk:"end_date_time"`

	// Computed outputs (from API response)
	KeyID               types.String `tfsdk:"key_id"`
	SecretText          types.String `tfsdk:"secret_text"`
	Hint                types.String `tfsdk:"hint"`
	CustomKeyIdentifier types.String `tfsdk:"custom_key_identifier"`

	Timeouts timeouts.Value `tfsdk:"timeouts"`
}
