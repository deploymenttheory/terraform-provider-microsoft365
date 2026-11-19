package graphBetaServicePrincipalSynchronizationJob

import (
	"context"

	"github.com/hashicorp/terraform-plugin-testing/terraform"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/exists"
)

// ServicePrincipalSynchronizationJobTestResource provides resource existence checks for acceptance tests.
type ServicePrincipalSynchronizationJobTestResource struct{}

// Exists checks whether the resource exists in Microsoft Graph.
func (r ServicePrincipalSynchronizationJobTestResource) Exists(ctx context.Context, _ any, state *terraform.InstanceState) (*bool, error) {
	//nolint:wrapcheck // Preserve the shared acceptance helper's result.
	return exists.CheckResourceExists(ctx, state, func(client *msgraphbetasdk.GraphServiceClient, ctx context.Context, state *terraform.InstanceState) error {
		parent, job, err := parseResourceID(state.ID)
		if err != nil {
			return err
		}
		_, err = client.ServicePrincipals().ByServicePrincipalId(parent).Synchronization().Jobs().BySynchronizationJobId(job).Get(ctx, nil)
		return err
	})
}
