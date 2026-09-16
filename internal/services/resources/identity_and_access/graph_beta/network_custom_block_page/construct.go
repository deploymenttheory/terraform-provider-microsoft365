package graphBetaNetworkCustomBlockPage

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"
	models "github.com/microsoftgraph/msgraph-beta-sdk-go/models/networkaccess"
)

func constructResource(
	data *NetworkCustomBlockPageResourceModel,
) (models.CustomBlockPageable, error) {
	state, err := models.ParseStatus(data.State.ValueString())
	if err != nil || state == nil {
		return nil, fmt.Errorf("%w: %q", errInvalidState, data.State.ValueString())
	}
	body := models.NewCustomBlockPage()
	// Send writable fields only. The singleton's synthetic Terraform ID is never sent.
	body.SetOdataType(nil)
	body.SetState(state.(*models.Status))
	if !data.Configuration.IsNull() && !data.Configuration.IsUnknown() {
		attrs := data.Configuration.Attributes()
		message := attrs["body"].(types.String)
		if message.IsNull() || message.IsUnknown() {
			return nil, errUnknownBody
		}
		if err := validateBlockMessage(message.ValueString()); err != nil {
			return nil, err
		}
		configuration := models.NewMarkdownBlockMessageConfiguration()
		configuration.SetBody(message.ValueStringPointer())
		body.SetConfiguration(configuration)
	}

	return body, nil
}
