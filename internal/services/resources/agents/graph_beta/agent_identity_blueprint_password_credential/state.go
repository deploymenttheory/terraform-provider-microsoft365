package graphBetaAgentIdentityBlueprintPasswordCredential

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteResourceStateToTerraform maps the password credential response to the Terraform state
func MapRemoteResourceStateToTerraform(ctx context.Context, data *AgentIdentityBlueprintPasswordCredentialResourceModel, credential graphmodels.PasswordCredentialable) {
	if credential == nil {
		tflog.Warn(ctx, "Received nil credential in MapRemoteResourceStateToTerraform")
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Starting to map remote state to Terraform state for %s", ResourceName))

	data.KeyID = convert.GraphToFrameworkUUID(credential.GetKeyId())
	data.SecretText = convert.GraphToFrameworkString(credential.GetSecretText())
	data.Hint = convert.GraphToFrameworkString(credential.GetHint())
	data.DisplayName = convert.GraphToFrameworkString(credential.GetDisplayName())
	data.StartDateTime = convert.GraphToFrameworkTime(credential.GetStartDateTime())
	data.EndDateTime = convert.GraphToFrameworkTime(credential.GetEndDateTime())
	data.CustomKeyIdentifier = convert.GraphToFrameworkBytes(credential.GetCustomKeyIdentifier())

	tflog.Debug(ctx, fmt.Sprintf("Finished mapping %s with key_id: %s", ResourceName, data.KeyID.ValueString()))
}
