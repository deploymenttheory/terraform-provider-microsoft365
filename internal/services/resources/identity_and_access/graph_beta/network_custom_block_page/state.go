package graphBetaNetworkCustomBlockPage

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	models "github.com/microsoftgraph/msgraph-beta-sdk-go/models/networkaccess"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
)

func mapRemoteState(
	data *NetworkCustomBlockPageResourceModel,
	remote models.CustomBlockPageable,
) error {
	if remote == nil || remote.GetState() == nil {
		return errMissingState
	}
	configuration := types.ObjectNull(configurationAttrTypes)
	if config := remote.GetConfiguration(); config != nil {
		markdown, ok := config.(models.MarkdownBlockMessageConfigurationable)
		if !ok {
			return errUnsupportedConfiguration
		}
		configuration = types.ObjectValueMust(configurationAttrTypes, map[string]attr.Value{
			"body": convert.GraphToFrameworkString(markdown.GetBody()),
		})
	}
	data.ID = types.StringValue(singletonID)
	data.State = convert.GraphToFrameworkEnum(remote.GetState())
	data.Configuration = configuration
	return nil
}
