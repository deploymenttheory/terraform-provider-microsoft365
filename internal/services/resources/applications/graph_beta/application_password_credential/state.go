package graphBetaApplicationPasswordCredential

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteResourceStateToTerraform maps the password resource response to the Terraform state
func MapRemoteResourceStateToTerraform(ctx context.Context, data *ApplicationPasswordCredentialResourceModel, resource graphmodels.PasswordCredentialable) {
	if resource == nil {
		tflog.Warn(ctx, "Received nil resource in MapRemoteResourceStateToTerraform")
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Starting to map remote state to Terraform state for %s", ResourceName))

	data.KeyID = convert.GraphToFrameworkUUID(resource.GetKeyId())
	data.SecretText = convert.GraphToFrameworkString(resource.GetSecretText())
	data.Hint = convert.GraphToFrameworkString(resource.GetHint())
	data.DisplayName = convert.GraphToFrameworkString(resource.GetDisplayName())
	data.StartDateTime = convert.GraphToFrameworkTime(resource.GetStartDateTime())
	data.EndDateTime = convert.GraphToFrameworkTime(resource.GetEndDateTime())
	data.CustomKeyIdentifier = convert.GraphToFrameworkBytes(resource.GetCustomKeyIdentifier())

	tflog.Debug(ctx, fmt.Sprintf("Finished mapping %s with key_id: %s", ResourceName, data.KeyID.ValueString()))
}
