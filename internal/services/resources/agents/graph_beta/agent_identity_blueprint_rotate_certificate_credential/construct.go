package graphBetaAgentIdentityBlueprintKeyCredential

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	applications "github.com/microsoftgraph/msgraph-beta-sdk-go/applications"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// constructAddKeyRequest constructs the request body for the addKey API call
func constructAddKeyRequest(ctx context.Context, data *AgentIdentityBlueprintKeyCredentialResourceModel) (*applications.ItemAddKeyPostRequestBody, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))

	requestBody := applications.NewItemAddKeyPostRequestBody()
	keyCredential := graphmodels.NewKeyCredential()

	convert.FrameworkToGraphBytes(data.Key, keyCredential.SetKey)
	convert.FrameworkToGraphString(data.KeyType, keyCredential.SetTypeEscaped)
	convert.FrameworkToGraphString(data.Usage, keyCredential.SetUsage)
	convert.FrameworkToGraphString(data.DisplayName, keyCredential.SetDisplayName)
	convert.FrameworkToGraphBytes(data.CustomKeyIdentifier, keyCredential.SetCustomKeyIdentifier)

	if err := convert.FrameworkToGraphTime(data.StartDateTime, keyCredential.SetStartDateTime); err != nil {
		return nil, fmt.Errorf("failed to parse start_date_time: %w", err)
	}

	if err := convert.FrameworkToGraphTime(data.EndDateTime, keyCredential.SetEndDateTime); err != nil {
		return nil, fmt.Errorf("failed to parse end_date_time: %w", err)
	}

	requestBody.SetKeyCredential(keyCredential)

	// Set passwordCredential if provided (required for X509CertAndPassword type)
	if !data.PasswordSecretText.IsNull() && !data.PasswordSecretText.IsUnknown() {
		passwordCredential := graphmodels.NewPasswordCredential()
		convert.FrameworkToGraphString(data.PasswordSecretText, passwordCredential.SetSecretText)
		requestBody.SetPasswordCredential(passwordCredential)
	}

	// Set required proof
	convert.FrameworkToGraphString(data.Proof, requestBody.SetProof)

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s resource", ResourceName))

	return requestBody, nil
}

// constructRemoveKeyRequest constructs the request body for the removeKey API call using Kiota SDK
func constructRemoveKeyRequest(ctx context.Context, keyID string, proof string) (*applications.ItemRemoveKeyPostRequestBody, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing removeKey request for key_id: %s", keyID))

	requestBody := applications.NewItemRemoveKeyPostRequestBody()

	if err := convert.FrameworkToGraphUUID(types.StringValue(keyID), requestBody.SetKeyId); err != nil {
		return nil, fmt.Errorf("failed to parse key_id as UUID: %w", err)
	}

	requestBody.SetProof(&proof)

	if err := constructors.DebugLogGraphObject(ctx, "RemoveKey request body", requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, "Finished constructing removeKey request")

	return requestBody, nil
}
