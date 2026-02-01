package graphBetaAgentIdentityBlueprintPasswordCredential

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

// constructAddPasswordRequest constructs the request body for the addPassword API call
func constructAddPasswordRequest(ctx context.Context, data *AgentIdentityBlueprintPasswordCredentialResourceModel) (*applications.ItemAddPasswordPostRequestBody, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))

	if err := validateRequest(ctx, data); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	requestBody := applications.NewItemAddPasswordPostRequestBody()
	passwordCredential := graphmodels.NewPasswordCredential()

	convert.FrameworkToGraphString(data.DisplayName, passwordCredential.SetDisplayName)

	if err := convert.FrameworkToGraphTime(data.StartDateTime, passwordCredential.SetStartDateTime); err != nil {
		return nil, fmt.Errorf("failed to parse start_date_time: %w", err)
	}

	if err := convert.FrameworkToGraphTime(data.EndDateTime, passwordCredential.SetEndDateTime); err != nil {
		return nil, fmt.Errorf("failed to parse end_date_time: %w", err)
	}

	requestBody.SetPasswordCredential(passwordCredential)

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s resource", ResourceName))

	return requestBody, nil
}

// constructRemovePasswordRequest constructs the request body for the removePassword API call using Kiota SDK
func constructRemovePasswordRequest(ctx context.Context, keyID string) (*applications.ItemRemovePasswordPostRequestBody, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing removePassword request for key_id: %s", keyID))

	requestBody := applications.NewItemRemovePasswordPostRequestBody()

	if err := convert.FrameworkToGraphUUID(types.StringValue(keyID), requestBody.SetKeyId); err != nil {
		return nil, fmt.Errorf("failed to parse key_id as UUID: %w", err)
	}

	if err := constructors.DebugLogGraphObject(ctx, "RemovePassword request body", requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, "Finished constructing removePassword request")

	return requestBody, nil
}
