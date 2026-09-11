package graphBetaNetworkCloudFirewallPolicy

import (
	"context"
	"errors"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	graph "github.com/microsoftgraph/msgraph-beta-sdk-go/models/networkaccess"
)

var errInvalidResponse = errors.New("invalid cloud firewall response")

// MapRemoteStateToTerraform validates the response before replacing known state.
func MapRemoteStateToTerraform(_ context.Context, data *NetworkCloudFirewallPolicyResourceModel, remote graph.CloudFirewallPolicyable) error {
	if remote == nil || remote.GetId() == nil || *remote.GetId() == "" || remote.GetName() == nil || remote.GetSettings() == nil || remote.GetSettings().GetDefaultAction() == nil || remote.GetVersion() == nil || remote.GetLastModifiedDateTime() == nil {
		return fmt.Errorf("%w: cloud firewall policy response is missing required properties", errInvalidResponse)
	}
	if remote.GetSettings().GetDefaultAction().String() != "allow" {
		return fmt.Errorf("%w: unsupported cloud firewall default action %q", errInvalidResponse, remote.GetSettings().GetDefaultAction().String())
	}
	if !data.ID.IsNull() && !data.ID.IsUnknown() && data.ID.ValueString() != *remote.GetId() {
		return fmt.Errorf("%w: cloud firewall policy response ID does not match requested ID", errInvalidResponse)
	}
	data.ID = convert.GraphToFrameworkString(remote.GetId())
	data.Name = convert.GraphToFrameworkString(remote.GetName())
	data.Description = convert.GraphToFrameworkString(remote.GetDescription())
	data.DefaultAction = convert.GraphToFrameworkEnum(remote.GetSettings().GetDefaultAction())
	data.Version = convert.GraphToFrameworkString(remote.GetVersion())
	data.LastModifiedDateTime = convert.GraphToFrameworkTime(remote.GetLastModifiedDateTime())
	return nil
}
