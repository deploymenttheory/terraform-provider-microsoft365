package graphBetaNetworkCrossTenantAccessSettings

import (
	"errors"
	"fmt"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-framework/types"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models/networkaccess"
)

var (
	errMissingStatus = errors.New("response is missing networkPacketTaggingStatus")
	errInvalidStatus = errors.New("unsupported networkPacketTaggingStatus")
)

func mapRemoteState(data *NetworkCrossTenantAccessSettingsResourceModel, remote graphmodels.CrossTenantAccessSettingsable) error {
	if remote == nil || remote.GetNetworkPacketTaggingStatus() == nil {
		return errMissingStatus
	}
	status := convert.GraphToFrameworkEnum(remote.GetNetworkPacketTaggingStatus())
	if status.ValueString() != "enabled" && status.ValueString() != "disabled" {
		return fmt.Errorf("%w: %q", errInvalidStatus, status.ValueString())
	}
	data.ID = types.StringValue(singletonID)
	data.NetworkPacketTaggingStatus = status
	return nil
}
