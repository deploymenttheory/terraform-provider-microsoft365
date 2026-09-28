package graphBetaDeviceConfigurationTemplatesJson

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"

	sharedconstructors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors/graph_beta/device_management"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
)

func constructResource(
	ctx context.Context,
	data *DeviceConfigurationTemplatesJsonResourceModel,
) (graphmodels.DeviceConfigurationable, error) {
	requestBody, err := sharedconstructors.ConstructDeviceConfigurationSettings(data.Settings)
	if err != nil {
		return nil, fmt.Errorf("construct template settings: %w", err)
	}
	convert.FrameworkToGraphString(data.DisplayName, requestBody.SetDisplayName)
	convert.FrameworkToGraphString(data.Description, requestBody.SetDescription)
	if err := convert.FrameworkToGraphStringSet(
		ctx,
		data.RoleScopeTagIds,
		requestBody.SetRoleScopeTagIds,
	); err != nil {
		return nil, fmt.Errorf("set role scope tags: %w", err)
	}
	tflog.Debug(ctx, "Constructed device configuration template request")
	return requestBody, nil
}
