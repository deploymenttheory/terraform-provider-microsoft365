package graphBetaNetworkMCPPolicy

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
)

func MapRemoteStateToTerraform(
	_ context.Context,
	data *NetworkMCPPolicyResourceModel,
	remote *mcpPolicyResponse,
) error {
	if remote == nil || remote.ID == nil || *remote.ID == "" || remote.Name == nil ||
		remote.Settings == nil || remote.Settings.DefaultAction == nil ||
		remote.Version == nil ||
		remote.LastModifiedDateTime == nil {
		return fmt.Errorf("%w: missing required fields", errInvalidResponse)
	}
	if *remote.Settings.DefaultAction != "allow" && *remote.Settings.DefaultAction != "block" {
		return fmt.Errorf(
			"%w: unsupported default action %q",
			errInvalidResponse,
			*remote.Settings.DefaultAction,
		)
	}
	data.ID = convert.GraphToFrameworkString(remote.ID)
	data.Name = convert.GraphToFrameworkString(remote.Name)
	data.Description = convert.GraphToFrameworkString(remote.Description)
	data.DefaultAction = convert.GraphToFrameworkString(remote.Settings.DefaultAction)
	data.Version = convert.GraphToFrameworkString(remote.Version)
	data.LastModifiedDateTime = convert.GraphToFrameworkString(remote.LastModifiedDateTime)
	return nil
}
