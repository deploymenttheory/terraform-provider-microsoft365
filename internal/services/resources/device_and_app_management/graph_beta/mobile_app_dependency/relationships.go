package graphBetaMobileAppDependency

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
		"a relationship already exists for the target app; import an existing dependency instead of creating it",
	)
	errMissingDependency       = errors.New("app dependency not found (404)")
	errConflictingRelationship = errors.New(
		"the target app relationship is not the managed dependency",
	)
	errMissingDependencyDeletion = errors.New("app dependency is still visible after deletion")
)

// listRelationships retrieves every page of the source app's relationship collection.
func (r *MobileAppDependencyResource) listRelationships(
	ctx context.Context,
	sourceID string,
) ([]graphmodels.MobileAppRelationshipable, error) {
	relationships, err := apprelationships.List(ctx, r.client, sourceID)
	if err != nil {
		return nil, fmt.Errorf("read app dependency relationships: %w", err)
	}
	return relationships, nil
}

// isIncomingRelationship identifies parent relationships excluded from outgoing updates.
func isIncomingRelationship(relationship graphmodels.MobileAppRelationshipable) bool {
	return apprelationships.Incoming(relationship)
}

// changeDependency applies one dependency change while retaining other outgoing relationships.
func (r *MobileAppDependencyResource) changeDependency(
	ctx context.Context,
	data *MobileAppDependencyResourceModel,
	operation string,
) (bool, error) {
	var replacement graphmodels.MobileAppDependencyable
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
		dependency, ok := relationship.(graphmodels.MobileAppDependencyable)
		if !ok {
			return false, errConflictingRelationship
		}
		found = true
		if operation == constants.TfOperationUpdate {
			dependency.SetDependencyType(replacement.GetDependencyType())
			desired = append(desired, dependency)
		}
	}
	if operation == constants.TfOperationCreate {
		desired = append(desired, replacement)
	} else if !found {
		if operation == constants.TfOperationUpdate {
			return false, errMissingDependency
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
