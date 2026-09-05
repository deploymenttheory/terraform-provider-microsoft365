package graphBetaNetworkTLSInspectionPolicyRule

import (
	"context"

	"github.com/hashicorp/terraform-plugin-testing/terraform"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/exists"
)

type NetworkTLSInspectionPolicyRuleTestResource struct{}

func (r NetworkTLSInspectionPolicyRuleTestResource) Exists(
	ctx context.Context,
	_ any,
	state *terraform.InstanceState,
) (*bool, error) {
	//nolint:wrapcheck // The generic existence helper already adds operation context.
	return exists.CheckResourceExists(
		ctx,
		state,
		func(client *msgraphbetasdk.GraphServiceClient, ctx context.Context, state *terraform.InstanceState) error {
			resource := &NetworkTLSInspectionPolicyRuleResource{client: client}
			_, err := resource.getTLSInspectionPolicyRule(
				ctx,
				state.Attributes["tls_inspection_policy_id"],
				state.ID,
			)
			return err
		},
	)
}
