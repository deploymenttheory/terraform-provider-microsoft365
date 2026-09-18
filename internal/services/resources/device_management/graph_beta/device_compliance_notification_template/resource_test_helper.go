package graphBetaDeviceComplianceNotificationTemplate

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/exists"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

type DeviceComplianceNotificationTemplateTestResource struct{}

func (r DeviceComplianceNotificationTemplateTestResource) Exists(ctx context.Context, _ any, state *terraform.InstanceState) (*bool, error) {
	//nolint:wrapcheck // Direct pass-through to generic helper with contextual errors
	return exists.CheckResourceExists(ctx, state, func(client *msgraphbetasdk.GraphServiceClient, ctx context.Context, state *terraform.InstanceState) error {
		_, err := client.
			DeviceManagement().
			NotificationMessageTemplates().
			ByNotificationMessageTemplateId(state.ID).
			Get(ctx, nil)
		return err
	})
}
