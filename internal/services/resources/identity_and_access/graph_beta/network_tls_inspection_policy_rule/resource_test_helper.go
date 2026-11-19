package graphBetaNetworkTLSInspectionPolicyRule

import (
	"context"

	"github.com/hashicorp/terraform-plugin-testing/terraform"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
	models "github.com/microsoftgraph/msgraph-beta-sdk-go/models/networkaccess"

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
			readBuilder := resource.client.NetworkAccess().
				TlsInspectionPolicies().
				ByTlsInspectionPolicyId(state.Attributes["tls_inspection_policy_id"]).
				PolicyRules().
				ByPolicyRuleId(state.ID)
			readRequest, err := readBuilder.ToGetRequestInformation(ctx, nil)
			if err != nil {
				return err
			}
			readResult, err := resource.client.GetAdapter().Send(
				ctx,
				readRequest,
				models.CreateTlsInspectionRuleFromDiscriminatorValue,
				tlsInspectionRuleErrorMapping,
			)
			_, err = decodeGetTLSInspectionPolicyRuleResponse(readResult, err)
			return err
		},
	)
}
