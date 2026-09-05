package graphBetaNetworkTLSInspectionPolicy

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
)

func MapRemoteStateToTerraform(
	_ context.Context,
	data *NetworkTLSInspectionPolicyResourceModel,
	remote *tlsInspectionPolicyResponse,
) error {
	if remote == nil || remote.id == nil || *remote.id == "" || remote.name == nil ||
		remote.defaultAction == nil ||
		remote.version == nil ||
		remote.lastModifiedDateTime == nil {
		return fmt.Errorf("%w: missing required fields", errInvalidResponse)
	}
	data.ID = convert.GraphToFrameworkString(remote.id)
	data.Name = convert.GraphToFrameworkString(remote.name)
	data.Description = convert.GraphToFrameworkString(remote.description)
	data.DefaultAction = convert.GraphToFrameworkString(remote.defaultAction)
	data.Version = convert.GraphToFrameworkString(remote.version)
	data.LastModifiedDateTime = convert.GraphToFrameworkString(remote.lastModifiedDateTime)
	return nil
}
