package graphBetaNetworkFilteringProfile

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	models "github.com/microsoftgraph/msgraph-beta-sdk-go/models/networkaccess"
)

var (
	errInvalidState     = errors.New("invalid filtering profile state")
	errPriorityRequired = errors.New("filtering profile priority is required")
	errStateRequired    = errors.New("filtering profile state is required")
)

// constructResource maps Terraform schema values to the minimal Graph payload: name, description, priority and
// state. See the deviation block in crud.go for why this is raw JSON rather than the SDK FilteringProfile model.
//
// Microsoft Learn documents filteringProfile and the supported read/update
// surface here:
// https://learn.microsoft.com/en-us/graph/api/resources/networkaccess-filteringprofile?view=graph-rest-beta
// https://learn.microsoft.com/en-us/graph/api/networkaccess-filteringprofile-update?view=graph-rest-beta
func constructResource(ctx context.Context, data *NetworkFilteringProfileResourceModel) (filteringProfileRequest, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))

	var requestBody filteringProfileRequest
	if !data.Name.IsNull() && !data.Name.IsUnknown() {
		requestBody.Name = data.Name.ValueStringPointer()
	}
	if !data.Description.IsNull() && !data.Description.IsUnknown() {
		requestBody.Description = data.Description.ValueStringPointer()
	}
	if !data.Priority.IsNull() && !data.Priority.IsUnknown() {
		requestBody.Priority = data.Priority.ValueInt64Pointer()
	}

	if !data.State.IsNull() && !data.State.IsUnknown() {
		parsed, err := models.ParseStatus(data.State.ValueString())
		if err != nil {
			return filteringProfileRequest{}, fmt.Errorf("%w: %w", errInvalidState, err)
		}
		if parsed == nil {
			return filteringProfileRequest{}, fmt.Errorf("%w: %s", errInvalidState, data.State.ValueString())
		}
		requestBody.State = data.State.ValueStringPointer()
	}

	if requestBody.Priority == nil {
		return filteringProfileRequest{}, errPriorityRequired
	}

	if requestBody.State == nil {
		return filteringProfileRequest{}, errStateRequired
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s resource", ResourceName))

	return requestBody, nil
}
