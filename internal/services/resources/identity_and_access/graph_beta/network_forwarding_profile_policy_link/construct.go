package graphBetaNetworkForwardingProfilePolicyLink

import (
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	models "github.com/microsoftgraph/msgraph-beta-sdk-go/models/networkaccess"
)

// constructResource builds the forwarding policy link PATCH body, which carries only the link state.
func constructResource(data *NetworkForwardingProfilePolicyLinkResourceModel) (models.ForwardingPolicyLinkable, error) {
	requestBody := models.NewForwardingPolicyLink()
	if err := convert.FrameworkToGraphEnum(data.State, models.ParseStatus, requestBody.SetState); err != nil {
		return nil, fmt.Errorf("failed to set state: %w", err)
	}
	return requestBody, nil
}
