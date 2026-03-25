package graphBetaAdministrativeUnitMembership

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// createReferenceRequest creates a reference request body for adding a member to an administrative unit
func createReferenceRequest(memberID string) graphmodels.ReferenceCreateable {
	reference := graphmodels.NewReferenceCreate()
	odataID := fmt.Sprintf("https://graph.microsoft.com/beta/directoryObjects/%s", memberID)
	reference.SetOdataId(&odataID)
	return reference
}

// constructAddMembersRequests creates a slice of reference requests for adding members
func constructAddMembersRequests(ctx context.Context, memberIDs []string) []graphmodels.ReferenceCreateable {
	tflog.Debug(ctx, fmt.Sprintf("Constructing add members requests for %d members", len(memberIDs)))

	requests := make([]graphmodels.ReferenceCreateable, 0, len(memberIDs))
	for _, id := range memberIDs {
		requests = append(requests, createReferenceRequest(id))
	}

	if len(requests) > 0 {
		if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("First member reference to be sent to Graph API for resource %s", ResourceName), requests[0]); err != nil {
			tflog.Error(ctx, "Failed to debug log object", map[string]any{
				"error": err.Error(),
			})
		}
	}

	return requests
}

// extractMemberIDsFromSet extracts string IDs from a Terraform set
func extractMemberIDsFromSet(memberSet types.Set) []string {
	if memberSet.IsNull() || memberSet.IsUnknown() {
		return []string{}
	}

	elements := memberSet.Elements()
	ids := make([]string, 0, len(elements))
	for _, elem := range elements {
		if strVal, ok := elem.(types.String); ok && !strVal.IsNull() {
			ids = append(ids, strVal.ValueString())
		}
	}
	return ids
}

// calculateMembershipChanges determines which members to add and remove
func calculateMembershipChanges(planMembers, stateMembers types.Set) (toAdd []string, toRemove []string) {
	planIDs := make(map[string]bool)
	for _, id := range extractMemberIDsFromSet(planMembers) {
		planIDs[id] = true
	}

	stateIDs := make(map[string]bool)
	for _, id := range extractMemberIDsFromSet(stateMembers) {
		stateIDs[id] = true
	}

	for id := range planIDs {
		if !stateIDs[id] {
			toAdd = append(toAdd, id)
		}
	}

	for id := range stateIDs {
		if !planIDs[id] {
			toRemove = append(toRemove, id)
		}
	}

	return toAdd, toRemove
}
