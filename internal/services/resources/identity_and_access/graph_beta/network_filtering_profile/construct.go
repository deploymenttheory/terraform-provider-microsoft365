package graphBetaNetworkFilteringProfile

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	models "github.com/microsoftgraph/msgraph-beta-sdk-go/models/networkaccess"
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
			return filteringProfileRequest{}, fmt.Errorf("invalid filtering profile state: %s", err)
		}
		if parsed == nil {
			return filteringProfileRequest{}, fmt.Errorf("invalid filtering profile state: %s", data.State.ValueString())
		}
		requestBody.State = data.State.ValueStringPointer()
	}

	if requestBody.Priority == nil {
		return filteringProfileRequest{}, fmt.Errorf("filtering profile priority is required")
	}

	if requestBody.State == nil {
		return filteringProfileRequest{}, fmt.Errorf("filtering profile state is required")
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s resource", ResourceName))

	return requestBody, nil
}
