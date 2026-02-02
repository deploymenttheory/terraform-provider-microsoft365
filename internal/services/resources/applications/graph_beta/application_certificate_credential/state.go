package graphBetaApplicationCertificateCredential

import (
	"context"
	"encoding/hex"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteResourceStateToTerraform maps the key credential to the Terraform state
func MapRemoteResourceStateToTerraform(ctx context.Context, data *ApplicationCertificateCredentialResourceModel, credential graphmodels.KeyCredentialable) {
	tflog.Debug(ctx, fmt.Sprintf("Mapping credential to state for %s", ResourceName))

	data.KeyID = convert.GraphToFrameworkUUID(credential.GetKeyId())
	data.DisplayName = convert.GraphToFrameworkString(credential.GetDisplayName())
	data.StartDateTime = convert.GraphToFrameworkTime(credential.GetStartDateTime())
	data.EndDateTime = convert.GraphToFrameworkTime(credential.GetEndDateTime())
	data.KeyType = convert.GraphToFrameworkString(credential.GetTypeEscaped())
	data.Usage = convert.GraphToFrameworkString(credential.GetUsage())

	// The custom_key_identifier is auto-generated as the certificate thumbprint
	if customKeyId := credential.GetCustomKeyIdentifier(); customKeyId != nil {
		data.CustomKeyIdentifier = types.StringValue(hex.EncodeToString(customKeyId))
		data.Thumbprint = types.StringValue(hex.EncodeToString(customKeyId))
	} else {
		data.CustomKeyIdentifier = types.StringNull()
		data.Thumbprint = types.StringNull()
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished mapping %s with key_id: %s", ResourceName, data.KeyID.ValueString()))
}
