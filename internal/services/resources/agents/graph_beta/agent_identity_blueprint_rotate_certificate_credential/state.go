package graphBetaAgentIdentityBlueprintKeyCredential

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteResourceStateToTerraform maps the key credentials from an application response to the Terraform state
// by finding the key credential with the matching keyId
func MapRemoteResourceStateToTerraform(ctx context.Context, data *AgentIdentityBlueprintKeyCredentialResourceModel, application graphmodels.Applicationable, keyID string) error {
	if application == nil {
		return fmt.Errorf("received nil application in MapRemoteResourceStateToTerraform")
	}

	tflog.Debug(ctx, fmt.Sprintf("Starting to map remote state to Terraform state for %s", ResourceName))

	keyCredentials := application.GetKeyCredentials()

	// Find the key credential with matching keyId
	var matchedCredential graphmodels.KeyCredentialable
	for _, cred := range keyCredentials {
		if cred.GetKeyId() != nil {
			credKeyID := cred.GetKeyId().String()
			tflog.Debug(ctx, fmt.Sprintf("Checking key credential with keyId: %s", credKeyID))
			if credKeyID == keyID {
				matchedCredential = cred
				break
			}
		}
	}

	if matchedCredential == nil {
		return fmt.Errorf("key credential with keyId %s not found in application", keyID)
	}

	tflog.Debug(ctx, fmt.Sprintf("Found matching key credential with keyId: %s", keyID))

	// Map the matched credential to state using convert helpers
	data.KeyID = convert.GraphToFrameworkUUID(matchedCredential.GetKeyId())
	data.DisplayName = convert.GraphToFrameworkString(matchedCredential.GetDisplayName())
	data.StartDateTime = convert.GraphToFrameworkTime(matchedCredential.GetStartDateTime())
	data.EndDateTime = convert.GraphToFrameworkTime(matchedCredential.GetEndDateTime())
	data.CustomKeyIdentifier = convert.GraphToFrameworkBytes(matchedCredential.GetCustomKeyIdentifier())

	// Map type and usage from the API response
	if matchedCredential.GetTypeEscaped() != nil {
		data.KeyType = types.StringValue(*matchedCredential.GetTypeEscaped())
	}
	if matchedCredential.GetUsage() != nil {
		data.Usage = types.StringValue(*matchedCredential.GetUsage())
	}

	// Note: Key, Proof, and PasswordSecretText are input-only values
	// preserved from the plan/state as they cannot be retrieved from the API
	// (key is only returned with $select, and proof/password are never returned)

	tflog.Debug(ctx, fmt.Sprintf("Finished mapping %s with key_id: %s", ResourceName, data.KeyID.ValueString()))

	return nil
}
