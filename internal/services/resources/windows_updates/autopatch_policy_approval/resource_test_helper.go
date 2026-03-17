package graphBetaWindowsUpdatesAutopatchPolicyApproval

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/exists"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

type WindowsUpdatesAutopatchPolicyApprovalTestResource struct{}

func (r WindowsUpdatesAutopatchPolicyApprovalTestResource) Exists(ctx context.Context, _ any, state *terraform.InstanceState) (*bool, error) {
	//nolint:wrapcheck
	return exists.CheckResourceExists(ctx, state, func(client *msgraphbetasdk.GraphServiceClient, ctx context.Context, state *terraform.InstanceState) error {
		policyId := state.Attributes["policy_id"]
		approvalId := state.ID
		_, err := client.Admin().Windows().Updates().Policies().ByPolicyId(policyId).Approvals().ByPolicyApprovalId(approvalId).Get(ctx, nil)
		return err
	})
}
