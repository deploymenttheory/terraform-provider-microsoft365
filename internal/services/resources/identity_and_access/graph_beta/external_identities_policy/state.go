package graphBetaExternalIdentitiesPolicy

import (
	"context"
	"errors"

	"github.com/hashicorp/terraform-plugin-framework/types"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
)

var errInvalidPolicyResponse = errors.New(
	"response must contain the externalIdentityPolicy singleton and both policy Booleans",
)

// mapRemoteState validates and maps the singleton object returned by the beta GET endpoint.
func mapRemoteState(
	ctx context.Context,
	data *ExternalIdentitiesPolicyResourceModel,
	remote graphmodels.ExternalIdentitiesPolicyable,
) error {
	if remote == nil || remote.GetId() == nil || *remote.GetId() != singletonID ||
		remote.GetAllowExternalIdentitiesToLeave() == nil || remote.GetAllowDeletedIdentitiesDataRemoval() == nil {
		return errInvalidPolicyResponse
	}
	data.ID = types.StringValue(singletonID)
	data.DisplayName = convert.GraphToFrameworkString(remote.GetDisplayName())
	data.Description = convert.GraphToFrameworkString(remote.GetDescription())
	data.AllowExternalIdentitiesToLeave = convert.GraphToFrameworkBool(
		remote.GetAllowExternalIdentitiesToLeave(),
	)
	data.AllowDeletedIdentitiesDataRemoval = convert.GraphToFrameworkBool(
		remote.GetAllowDeletedIdentitiesDataRemoval(),
	)
	return nil
}
