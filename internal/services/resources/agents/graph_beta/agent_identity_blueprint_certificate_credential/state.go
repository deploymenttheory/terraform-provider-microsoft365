package graphBetaAgentIdentityBlueprintCertificateCredential

import (
	"context"
	"encoding/hex"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteResourceStateToTerraform maps the key credentials response to the Terraform state
// by finding the key credential with the matching keyId
func MapRemoteResourceStateToTerraform(ctx context.Context, data *AgentIdentityBlueprintCertificateCredentialResourceModel, application graphmodels.Applicationable, keyID string) error {
	if application == nil {
		return fmt.Errorf("received nil application in MapRemoteResourceStateToTerraform")
	}

	tflog.Debug(ctx, fmt.Sprintf("Starting to map remote state to Terraform state for %s", ResourceName))

	keyCredentials := application.GetKeyCredentials()

	// Find the key credential with matching keyId
	var matchedCredential graphmodels.KeyCredentialable
	for _, cred := range keyCredentials {
		if cred.GetKeyId() != nil && cred.GetKeyId().String() == keyID {
			matchedCredential = cred
			break
		}
	}

	if matchedCredential == nil {
		return fmt.Errorf("certificate credential with keyId %s not found in application", keyID)
	}

	tflog.Debug(ctx, fmt.Sprintf("Found matching certificate credential with keyId: %s", matchedCredential.GetKeyId().String()))

	// Map the matched credential to state
	data.KeyID = convert.GraphToFrameworkUUID(matchedCredential.GetKeyId())
	data.DisplayName = convert.GraphToFrameworkString(matchedCredential.GetDisplayName())
	data.StartDateTime = convert.GraphToFrameworkTime(matchedCredential.GetStartDateTime())
	data.EndDateTime = convert.GraphToFrameworkTime(matchedCredential.GetEndDateTime())
	data.KeyType = convert.GraphToFrameworkString(matchedCredential.GetTypeEscaped())
	data.Usage = convert.GraphToFrameworkString(matchedCredential.GetUsage())

	// Map custom key identifier (bytes to hex string for display)
	if customKeyId := matchedCredential.GetCustomKeyIdentifier(); customKeyId != nil {
		data.CustomKeyIdentifier = types.StringValue(hex.EncodeToString(customKeyId))
	} else {
		data.CustomKeyIdentifier = types.StringNull()
	}

	// Map thumbprint from custom key identifier (they are often the same for certificates)
	if customKeyId := matchedCredential.GetCustomKeyIdentifier(); customKeyId != nil {
		data.Thumbprint = types.StringValue(hex.EncodeToString(customKeyId))
	} else {
		data.Thumbprint = types.StringNull()
	}

	// Note: Key is input-only and cannot be retrieved from the API

	tflog.Debug(ctx, fmt.Sprintf("Finished mapping %s with key_id: %s", ResourceName, data.KeyID.ValueString()))

	return nil
}
