package graphBetaGroupLifecycleExpirationPolicyAssignment

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	graphgrouplifecyclepolicies "github.com/microsoftgraph/msgraph-beta-sdk-go/grouplifecyclepolicies"
)

// constructAddGroupRequest constructs the request body for adding a group to the lifecycle policy.
func (r *GroupLifecycleExpirationPolicyAssignmentResource) constructAddGroupRequest(
	ctx context.Context,
	groupID string,
	diagnostics *diag.Diagnostics,
) (policyID string, requestBody *graphgrouplifecyclepolicies.ItemAddGroupPostRequestBody, err error) {

	policyID, err = r.validateRequest(ctx, groupID, diagnostics)
	if err != nil || policyID == "" {
		return "", nil, err
	}

	requestBody = graphgrouplifecyclepolicies.NewItemAddGroupPostRequestBody()
	requestBody.SetGroupId(&groupID)

	return policyID, requestBody, nil
}

// constructRemoveGroupRequest constructs the request body for removing a group from the lifecycle policy.
func (r *GroupLifecycleExpirationPolicyAssignmentResource) constructRemoveGroupRequest(
	ctx context.Context,
	groupID string,
	diagnostics *diag.Diagnostics,
) (policyID string, requestBody *graphgrouplifecyclepolicies.ItemRemoveGroupPostRequestBody, err error) {

	policyID, err = r.validateRequest(ctx, groupID, diagnostics)
	if err != nil || policyID == "" {
		return "", nil, err
	}

	requestBody = graphgrouplifecyclepolicies.NewItemRemoveGroupPostRequestBody()
	requestBody.SetGroupId(&groupID)

	return policyID, requestBody, nil
}
