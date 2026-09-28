package graphBetaIdentityAndAccessB2bManagementPolicyAssignment

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/exists"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

// B2bManagementPolicyAssignmentTestResource implements the types.TestResource interface
type B2bManagementPolicyAssignmentTestResource struct{}

// Exists checks whether the B2B management policy applies to the directory object in Microsoft Graph
func (r B2bManagementPolicyAssignmentTestResource) Exists(ctx context.Context, _ any, state *terraform.InstanceState) (*bool, error) {
	return exists.CheckResourceExists(ctx, state, func(client *msgraphbetasdk.GraphServiceClient, ctx context.Context, state *terraform.InstanceState) error {
		policyID := state.Attributes["b2b_management_policy_id"]
		directoryObjectID := state.Attributes["directory_object_id"]

		if policyID == "" || directoryObjectID == "" {
			return fmt.Errorf("missing b2b_management_policy_id or directory_object_id in state")
		}

		directoryObjects, err := client.Policies().B2bManagementPolicies().ByB2bManagementPolicyId(policyID).AppliesTo().Get(ctx, nil)
		if err != nil {
			return err
		}

		for _, directoryObject := range directoryObjects.GetValue() {
			if directoryObject.GetId() != nil && *directoryObject.GetId() == directoryObjectID {
				return nil
			}
		}

		return fmt.Errorf("B2B management policy %s does not apply to directory object %s", policyID, directoryObjectID)
	})
}
