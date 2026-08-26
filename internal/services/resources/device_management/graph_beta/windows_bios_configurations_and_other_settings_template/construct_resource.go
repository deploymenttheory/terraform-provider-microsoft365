package graphBetaWindowsBiosConfigurationsAndOtherSettingsTemplate

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
)

var (
	errDecodeConfigurationFileContent = errors.New(
		"failed to decode configuration_file_content, expected a base64 encoded string",
	)
	errSetHardwareConfigurationFormat = errors.New("failed to set hardwareConfigurationFormat")
	errSetRoleScopeTagIds             = errors.New("failed to set role scope tags")
)

// Main entry point to construct the intune hardware configuration resource for the Terraform provider.
func constructResource(
	ctx context.Context,
	data *WindowsBiosConfigurationsAndOtherSettingsTemplateResourceModel,
) (graphmodels.HardwareConfigurationable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))

	requestBody := graphmodels.NewHardwareConfiguration()

	convert.FrameworkToGraphString(data.DisplayName, requestBody.SetDisplayName)
	convert.FrameworkToGraphString(data.Description, requestBody.SetDescription)
	convert.FrameworkToGraphString(data.FileName, requestBody.SetFileName)

	// The schema holds base64, and Kiota base64 encodes []byte during serialization. Decoding here
	// keeps the value that reaches the Graph API byte for byte identical to what the practitioner supplied.
	if !data.ConfigurationFileContent.IsNull() && !data.ConfigurationFileContent.IsUnknown() {
		decoded, err := base64.StdEncoding.DecodeString(data.ConfigurationFileContent.ValueString())
		if err != nil {
			return nil, fmt.Errorf("%w: %w", errDecodeConfigurationFileContent, err)
		}
		requestBody.SetConfigurationFileContent(decoded)
	}

	if err := convert.FrameworkToGraphEnum(
		data.HardwareConfigurationFormat,
		graphmodels.ParseHardwareConfigurationFormat,
		requestBody.SetHardwareConfigurationFormat,
	); err != nil {
		return nil, fmt.Errorf("%w: %w", errSetHardwareConfigurationFormat, err)
	}

	convert.FrameworkToGraphBool(
		data.PerDevicePasswordDisabled,
		requestBody.SetPerDevicePasswordDisabled,
	)

	if err := convert.FrameworkToGraphStringSet(
		ctx,
		data.RoleScopeTagIds,
		requestBody.SetRoleScopeTagIds,
	); err != nil {
		return nil, fmt.Errorf("%w: %w", errSetRoleScopeTagIds, err)
	}

	if err := constructors.DebugLogGraphObject(
		ctx,
		fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName),
		requestBody,
	); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s resource", ResourceName))

	return requestBody, nil
}
