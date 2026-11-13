package graphBetaMobileAppSupersedence

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/deviceappmanagement"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	apprelationships "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud/graph_beta/device_and_app_management/app_relationships"
	kiotaerrors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
)

// updateRetryInterval is the wait between attempts after Intune rejects updateRelationships
// with 412 Precondition Failed.
const updateRetryInterval = 5 * time.Second

var (
	errExistingRelationship = errors.New(
		"a relationship already exists for the target app; import an existing supersedence instead of creating it",
	)
	errMissingSupersedence     = errors.New("app supersedence not found (404)")
	errConflictingRelationship = errors.New(
		"the target app relationship is not the managed supersedence",
	)
	errMissingSupersedenceDeletion = errors.New("app supersedence is still visible after deletion")
)

// listRelationships retrieves every page of the source app's relationship collection.
func (r *MobileAppSupersedenceResource) listRelationships(
	ctx context.Context,
	sourceID string,
) ([]graphmodels.MobileAppRelationshipable, error) {
	relationships, err := apprelationships.List(ctx, r.client, sourceID)
	if err != nil {
		return nil, fmt.Errorf("read app supersedence relationships: %w", err)
	}
	return relationships, nil
}

// isIncomingRelationship identifies parent relationships excluded from outgoing updates.
func isIncomingRelationship(relationship graphmodels.MobileAppRelationshipable) bool {
	return apprelationships.Incoming(relationship)
}

// changeSupersedence applies one supersedence change while retaining other outgoing relationships.
//
// Intune rejects updateRelationships with 412 Precondition Failed while it is still
// processing an earlier relationship change for the app (for example a delete straight
// after an update). Each attempt re-reads the collection and rebuilds the request, so a
// retry never sends a stale collection, and a write that was applied despite the 412 is
// recognised instead of reported as a conflict. Other errors fail immediately.
func (r *MobileAppSupersedenceResource) changeSupersedence(
	ctx context.Context,
	data *MobileAppSupersedenceResourceModel,
	operation string,
) (bool, error) {
	var changed bool
	attempt := 0
	err := crud.PollUntil(ctx, updateRetryInterval, func(ctx context.Context) (bool, error) {
		attempt++
		var err error
		changed, err = r.applySupersedence(ctx, data, operation, attempt > 1)
		if err == nil {
			return true, nil
		}
		if kiotaerrors.GraphError(ctx, err).StatusCode == http.StatusPreconditionFailed {
			tflog.Debug(ctx, fmt.Sprintf("updateRelationships 412, attempt %d", attempt))
			return false, err
		}
		return false, &crud.FatalPollError{Err: err}
	})
	if err != nil {
		return false, err //nolint:wrapcheck // already wrapped by applySupersedence
	}
	return changed, nil
}

// applySupersedence makes one read-modify-write attempt against the relationship collection.
// retrying is set after a 412, when the previous attempt may already have been applied.
func (r *MobileAppSupersedenceResource) applySupersedence(
	ctx context.Context,
	data *MobileAppSupersedenceResourceModel,
	operation string,
	retrying bool,
) (bool, error) {
	var replacement graphmodels.MobileAppSupersedenceable
	if operation != constants.TfOperationDelete {
		var err error
		replacement, err = constructResource(ctx, data)
		if err != nil {
			return false, err
		}
	}
	relationships, err := r.listRelationships(ctx, data.SourceID.ValueString())
	if err != nil {
		return false, err
	}
	desired := make([]graphmodels.MobileAppRelationshipable, 0, len(relationships)+1)
	found := false
	for _, relationship := range relationships {
		if isIncomingRelationship(relationship) {
			continue
		}
		if relationship.GetTargetId() == nil ||
			!strings.EqualFold(*relationship.GetTargetId(), data.TargetID.ValueString()) {
			desired = append(desired, relationship)
			continue
		}
		if operation == constants.TfOperationCreate {
			if retrying && supersedenceType(relationship) == data.SupersedenceType.ValueString() {
				return true, nil
			}
			return false, errExistingRelationship
		}
		if relationship.GetId() == nil || *relationship.GetId() != data.ID.ValueString() {
			desired = append(desired, relationship)
			continue
		}
		supersedence, ok := relationship.(graphmodels.MobileAppSupersedenceable)
		if !ok {
			return false, errConflictingRelationship
		}
		found = true
		if operation == constants.TfOperationUpdate {
			supersedence.SetSupersedenceType(replacement.GetSupersedenceType())
			desired = append(desired, supersedence)
		}
	}
	if operation == constants.TfOperationCreate {
		desired = append(desired, replacement)
	} else if !found {
		if operation == constants.TfOperationUpdate {
			return false, errMissingSupersedence
		}
		return false, nil
	}
	requestBody := deviceappmanagement.NewMobileAppsItemUpdateRelationshipsPostRequestBody()
	requestBody.SetRelationships(desired)
	err = r.client.
		DeviceAppManagement().
		MobileApps().
		ByMobileAppId(data.SourceID.ValueString()).
		UpdateRelationships().
		Post(ctx, requestBody, nil)
	if err != nil {
		return false, fmt.Errorf("update app relationships: %w", err)
	}
	return true, nil
}

// supersedenceType returns the relationship's supersedence type, or "" for other relationships.
func supersedenceType(relationship graphmodels.MobileAppRelationshipable) string {
	supersedence, ok := relationship.(graphmodels.MobileAppSupersedenceable)
	if !ok || supersedence.GetSupersedenceType() == nil {
		return ""
	}
	return supersedence.GetSupersedenceType().String()
}
