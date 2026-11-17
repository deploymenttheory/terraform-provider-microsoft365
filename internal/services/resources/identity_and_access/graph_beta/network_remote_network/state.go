package graphBetaIdentityAndAccessNetworkRemoteNetwork

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
)

func mapResponse(state *NetworkRemoteNetworkResourceModel, remote *remoteNetworkResponse) error {
	if remote == nil || remote.ID == nil || *remote.ID == "" || remote.Name == nil || remote.Region == nil {
		return errInvalidResponse
	}
	state.ID = convert.GraphToFrameworkString(remote.ID)
	state.Name = convert.GraphToFrameworkString(remote.Name)
	state.Region = convert.GraphToFrameworkString(remote.Region)
	state.Version = convert.GraphToFrameworkString(remote.Version)

	ids := make([]attr.Value, 0, len(remote.ForwardingProfiles))
	for _, profile := range remote.ForwardingProfiles {
		if profile.ID != nil && *profile.ID != "" {
			ids = append(ids, types.StringValue(*profile.ID))
		}
	}
	state.ForwardingProfileIDs = types.SetValueMust(types.StringType, ids)
	return nil
}
