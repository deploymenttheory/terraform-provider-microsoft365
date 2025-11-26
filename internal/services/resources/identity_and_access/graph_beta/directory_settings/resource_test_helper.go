package graphBetaDirectorySettings

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// DirectorySettingsTestResource implements the types.TestResource interface for Directory Settings
type DirectorySettingsTestResource struct{}

// Exists checks whether the directory settings exist in Microsoft Graph
// For a singleton resource, this checks if the Group.Unified settings object exists
func (r DirectorySettingsTestResource) Exists(ctx context.Context, _ any, state *terraform.InstanceState) (*bool, error) {
	graphClient, err := acceptance.TestGraphClient()
	if err != nil {
		return nil, err
	}

	settingsID := state.ID

	settings, err := graphClient.
		Settings().
		ByDirectorySettingId(settingsID).
		Get(ctx, nil)

	if err != nil {
		errorInfo := errors.GraphError(ctx, err)
		// 404 means it doesn't exist
		if errorInfo.StatusCode == 404 ||
			errorInfo.ErrorCode == "ResourceNotFound" ||
			errorInfo.ErrorCode == "Request_ResourceNotFound" ||
			errorInfo.ErrorCode == "ItemNotFound" {
			exists := false
			return &exists, nil
		}
		return nil, err
	}

	if settings == nil {
		exists := false
		return &exists, nil
	}

	exists := true
	return &exists, nil
}
