package graphBetaNetworkFilteringProfile

import (
	"context"
	stderrors "errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	models "github.com/microsoftgraph/msgraph-beta-sdk-go/models/networkaccess"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
)

// MapRemoteStateToTerraform maps the base properties of a NetworkFilteringProfileResourceModel to a Terraform state.
func MapRemoteStateToTerraform(ctx context.Context, data *NetworkFilteringProfileResourceModel, remoteResource models.FilteringProfileable) {
	if remoteResource == nil {
		tflog.Debug(ctx, "Remote resource is nil")
		return
	}

	tflog.Debug(ctx, "Starting to map remote state to Terraform state", map[string]any{
		"resourceId": convert.GraphToFrameworkString(remoteResource.GetId()),
	})

	data.ID = convert.GraphToFrameworkString(remoteResource.GetId())
	data.Name = convert.GraphToFrameworkString(remoteResource.GetName())
	data.Description = convert.GraphToFrameworkString(remoteResource.GetDescription())
	data.Priority = convert.GraphToFrameworkInt64(remoteResource.GetPriority())
	data.State = convert.GraphToFrameworkEnum(remoteResource.GetState())
	data.CreatedDateTime = convert.GraphToFrameworkTime(remoteResource.GetCreatedDateTime())
	data.LastModifiedDateTime = convert.GraphToFrameworkTime(remoteResource.GetLastModifiedDateTime())
	data.Version = convert.GraphToFrameworkString(remoteResource.GetVersion())

	tflog.Debug(ctx, fmt.Sprintf("Finished stating resource %s with id %s", ResourceName, data.ID.ValueString()))
}

// decodeCreateFilteringProfileResponse validates and decodes the Microsoft Graph response.
func decodeCreateFilteringProfileResponse(result any, err error) (models.FilteringProfileable, error) {
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, errCreateResponseNil
	}
	profile, ok := result.(models.FilteringProfileable)
	if !ok {
		return nil, fmt.Errorf("create filtering profile returned %T, %w", result, errUnexpectedResponseType)
	}
	return profile, nil
}

var (
	errCreateResponseNil      = stderrors.New("create filtering profile returned nil response")
	errUnexpectedResponseType = stderrors.New("expected FilteringProfileable")
)
