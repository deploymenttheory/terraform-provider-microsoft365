package graphBetaWindowsAutopatchDeviceRegistration

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/exists"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/admin"
)

type WindowsAutopatchDeviceRegistrationTestResource struct{}

func (r WindowsAutopatchDeviceRegistrationTestResource) Exists(ctx context.Context, _ any, state *terraform.InstanceState) (*bool, error) {
	return exists.CheckResourceExists(ctx, state, func(client *msgraphbetasdk.GraphServiceClient, ctx context.Context, state *terraform.InstanceState) error {
		updateCategory := state.Attributes["update_category"]
		filterQuery := "$filter=isof('microsoft.graph.windowsUpdates.azureADDevice')"

		assetsResp, err := client.
			Admin().
			Windows().
			Updates().
			UpdatableAssets().
			Get(ctx, &admin.WindowsUpdatesUpdatableAssetsRequestBuilderGetRequestConfiguration{
				QueryParameters: &admin.WindowsUpdatesUpdatableAssetsRequestBuilderGetQueryParameters{
					Select: []string{"id", "enrollments"},
					Filter: &filterQuery,
				},
			})

		if err != nil {
			return err
		}

		devices := assetsResp.GetValue()
		if len(devices) == 0 {
			return fmt.Errorf("no devices enrolled for update category: %s", updateCategory)
		}

		return nil
	})
}
