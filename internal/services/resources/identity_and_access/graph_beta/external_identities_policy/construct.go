package graphBetaExternalIdentitiesPolicy

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
)

var errInvalidPolicy = errors.New(
	"allow_external_identities_to_leave and allow_deleted_identities_data_removal must be known Booleans",
)

// constructResource sends only the writable external identities policy settings.
func constructResource(
	ctx context.Context,
	data *ExternalIdentitiesPolicyResourceModel,
) (graphmodels.ExternalIdentitiesPolicyable, error) {
	if data.AllowExternalIdentitiesToLeave.IsNull() ||
		data.AllowExternalIdentitiesToLeave.IsUnknown() || data.AllowDeletedIdentitiesDataRemoval.IsNull() || data.AllowDeletedIdentitiesDataRemoval.IsUnknown() {
		return nil, errInvalidPolicy
	}
	body := graphmodels.NewExternalIdentitiesPolicy()
	convert.FrameworkToGraphBool(
		data.AllowExternalIdentitiesToLeave,
		body.SetAllowExternalIdentitiesToLeave,
	)
	convert.FrameworkToGraphBool(
		data.AllowDeletedIdentitiesDataRemoval,
		body.SetAllowDeletedIdentitiesDataRemoval,
	)
	if err := constructors.DebugLogGraphObject(
		ctx,
		fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName),
		body,
	); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{"error": err.Error()})
	}
	return body, nil
}
