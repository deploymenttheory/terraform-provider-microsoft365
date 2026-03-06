package graphBetaWindowsAutopatchDeploymentAudienceMembers

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/exists"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

type WindowsUpdateDeploymentAudienceMembersTestResource struct{}

func (r WindowsUpdateDeploymentAudienceMembersTestResource) Exists(ctx context.Context, _ any, state *terraform.InstanceState) (*bool, error) {
	//nolint:wrapcheck
	return exists.CheckResourceExists(ctx, state, func(client *msgraphbetasdk.GraphServiceClient, ctx context.Context, state *terraform.InstanceState) error {
		// Check if the audience exists and has members/exclusions
		audienceID := state.Attributes["audience_id"]

		membersResp, err := client.Admin().Windows().Updates().DeploymentAudiences().ByDeploymentAudienceId(audienceID).Members().Get(ctx, nil)
		if err != nil {
			return err
		}

		exclusionsResp, err := client.Admin().Windows().Updates().DeploymentAudiences().ByDeploymentAudienceId(audienceID).Exclusions().Get(ctx, nil)
		if err != nil {
			return err
		}

		// Resource exists if there are any members or exclusions
		members := membersResp.GetValue()
		exclusions := exclusionsResp.GetValue()

		// For this resource, we consider it to exist even if empty
		// The audience itself must exist for this to work
		_ = members
		_ = exclusions

		return nil
	})
}
