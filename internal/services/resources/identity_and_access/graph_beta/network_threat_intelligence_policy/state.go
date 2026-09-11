package graphBetaNetworkThreatIntelligencePolicy

import (
	"context"
	"fmt"

	models "github.com/microsoftgraph/msgraph-beta-sdk-go/models/networkaccess"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
)

// MapRemoteStateToTerraform validates the response before replacing known state.
func MapRemoteStateToTerraform(
	_ context.Context,
	data *NetworkThreatIntelligencePolicyResourceModel,
	remote models.ThreatIntelligencePolicyable,
) error {
	if remote == nil || remote.GetId() == nil || *remote.GetId() == "" || remote.GetName() == nil ||
		remote.GetSettings() == nil ||
		remote.GetSettings().GetDefaultAction() == nil ||
		remote.GetVersion() == nil ||
		remote.GetLastModifiedDateTime() == nil {
		return fmt.Errorf(
			"%w: invalid threat intelligence policy response: missing required fields",
			errInvalidResponse,
		)
	}
	data.ID = convert.GraphToFrameworkString(remote.GetId())
	data.Name = convert.GraphToFrameworkString(remote.GetName())
	data.Description = convert.GraphToFrameworkString(remote.GetDescription())
	data.DefaultAction = convert.GraphToFrameworkEnum(remote.GetSettings().GetDefaultAction())
	data.Version = convert.GraphToFrameworkString(remote.GetVersion())
	data.LastModifiedDateTime = convert.GraphToFrameworkTime(remote.GetLastModifiedDateTime())
	return nil
}
