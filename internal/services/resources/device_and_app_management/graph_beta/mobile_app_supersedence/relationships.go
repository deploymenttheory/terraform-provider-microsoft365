package graphBetaMobileAppSupersedence

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/microsoftgraph/msgraph-beta-sdk-go/deviceappmanagement"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	apprelationships "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud/graph_beta/device_and_app_management/app_relationships"
)

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
func (r *MobileAppSupersedenceResource) changeSupersedence(
	ctx context.Context,
	data *MobileAppSupersedenceResourceModel,
	operation string,
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
