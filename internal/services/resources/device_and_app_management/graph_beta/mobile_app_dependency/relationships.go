package graphBetaMobileAppDependency

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/microsoftgraph/msgraph-beta-sdk-go/deviceappmanagement"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
)

var (
	errMissingRelationships = errors.New(
		"graph returned an incomplete relationship collection",
	)
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
	builder := r.client.DeviceAppManagement().MobileApps().ByMobileAppId(sourceID).Relationships()
	page, err := builder.Get(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("list app relationships: %w", err)
	}
	var relationships []graphmodels.MobileAppRelationshipable
	for {
		if page == nil || page.GetValue() == nil {
			return nil, errMissingRelationships
		}
		for _, relationship := range page.GetValue() {
			if relationship == nil {
				return nil, errMissingRelationships
			}
			relationships = append(relationships, relationship)
		}
		if nextLink := page.GetOdataNextLink(); nextLink != nil && *nextLink != "" {
			page, err = builder.WithUrl(*nextLink).Get(ctx, nil)
			if err != nil {
				return nil, fmt.Errorf("read next app relationship page: %w", err)
			}
			continue
		}
		return relationships, nil
	}
}

// isIncomingRelationship identifies parent relationships excluded from outgoing updates.
func isIncomingRelationship(relationship graphmodels.MobileAppRelationshipable) bool {
	return relationship.GetTargetType() != nil &&
		*relationship.GetTargetType() == graphmodels.PARENT_MOBILEAPPRELATIONSHIPTYPE
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
