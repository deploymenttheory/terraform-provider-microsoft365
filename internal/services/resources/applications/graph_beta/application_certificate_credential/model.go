// REF: https://learn.microsoft.com/en-us/graph/api/application-addkey?view=graph-rest-beta&tabs=http
// REF: https://learn.microsoft.com/en-us/graph/api/application-removekey?view=graph-rest-beta&tabs=http
package graphBetaApplicationCertificateCredential

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ApplicationCertificateCredentialResourceModel describes the resource data model.
type ApplicationCertificateCredentialResourceModel struct {
	// Required inputs
	ApplicationID types.String `tfsdk:"application_id"`
	Key           types.String `tfsdk:"key"`
	Encoding      types.String `tfsdk:"encoding"`
	Usage         types.String `tfsdk:"usage"`
	KeyType       types.String `tfsdk:"type"`

	// Optional inputs
	DisplayName                 types.String `tfsdk:"display_name"`
	StartDateTime               types.String `tfsdk:"start_date_time"`
	EndDateTime                 types.String `tfsdk:"end_date_time"`
	ReplaceExistingCertificates types.Bool   `tfsdk:"replace_existing_certificates"`

	// Computed outputs
	KeyID               types.String `tfsdk:"key_id"`
	CustomKeyIdentifier types.String `tfsdk:"custom_key_identifier"`
	Thumbprint          types.String `tfsdk:"thumbprint"`

	Timeouts timeouts.Value `tfsdk:"timeouts"`
}
