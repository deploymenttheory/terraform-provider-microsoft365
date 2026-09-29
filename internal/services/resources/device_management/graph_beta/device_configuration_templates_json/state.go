package graphBetaDeviceConfigurationTemplatesJson

import (
	"context"
	"errors"
	"fmt"

	"github.com/microsoft/kiota-abstractions-go/serialization"
	jsonserialization "github.com/microsoft/kiota-serialization-json-go"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
)

var errInvalidProfileResponse = errors.New("invalid device configuration response")

// MapRemoteResourceStateToTerraform maps the base profile metadata using the SDK model.
func MapRemoteResourceStateToTerraform(
	ctx context.Context,
	data *DeviceConfigurationTemplatesJsonResourceModel,
	response []byte,
) error {
	node, err := jsonserialization.NewJsonParseNode(response)
	if err != nil {
		return fmt.Errorf("parse device configuration response: %w", err)
	}
	// The base model reads metadata without discarding unknown derived-type settings.
	parsed, err := node.GetObjectValue(
		func(_ serialization.ParseNode) (serialization.Parsable, error) {
			return graphmodels.NewDeviceConfiguration(), nil
		},
	)
	if err != nil {
		return fmt.Errorf("deserialize device configuration metadata: %w", err)
	}
	remote, ok := parsed.(graphmodels.DeviceConfigurationable)
	if !ok || remote == nil || remote.GetId() == nil || *remote.GetId() == "" ||
		remote.GetOdataType() == nil ||
		remote.GetDisplayName() == nil {
		return fmt.Errorf(
			"%w: Graph must return id, displayName and @odata.type",
			errInvalidProfileResponse,
		)
	}
	data.ID = convert.GraphToFrameworkString(remote.GetId())
	data.DisplayName = convert.GraphToFrameworkString(remote.GetDisplayName())
	data.Description = convert.GraphToFrameworkString(remote.GetDescription())
	data.RoleScopeTagIds = convert.GraphToFrameworkStringSet(ctx, remote.GetRoleScopeTagIds())
	return nil
}
