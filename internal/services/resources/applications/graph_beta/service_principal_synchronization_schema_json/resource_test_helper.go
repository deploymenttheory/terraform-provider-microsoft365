package graphBetaServicePrincipalSynchronizationSchemaJson

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-testing/terraform"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/exists"
)

var errInvalidTestResourceID = errors.New("invalid synchronization schema ID")

// ServicePrincipalSynchronizationSchemaJsonTestResource implements the types.TestResource interface for synchronization schemas
type ServicePrincipalSynchronizationSchemaJsonTestResource struct{}

// Exists checks whether the synchronization job schema exists in Microsoft Graph
func (r ServicePrincipalSynchronizationSchemaJsonTestResource) Exists(ctx context.Context, _ any, state *terraform.InstanceState) (*bool, error) {
	//nolint:wrapcheck // Direct pass-through to generic helper with contextual errors
	return exists.CheckResourceExists(ctx, state, func(client *msgraphbetasdk.GraphServiceClient, ctx context.Context, state *terraform.InstanceState) error {
		servicePrincipalId, jobId, ok := strings.Cut(state.ID, "/")
		if !ok {
			return fmt.Errorf("%w %q", errInvalidTestResourceID, state.ID)
		}
		_, err := client.ServicePrincipals().ByServicePrincipalId(servicePrincipalId).Synchronization().Jobs().BySynchronizationJobId(jobId).Schema().Get(ctx, nil)
		return err
	})
}
