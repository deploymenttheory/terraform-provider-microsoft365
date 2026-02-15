package graphBetaWindowsPlatformScript

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/devicemanagement"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models"
	graphcore "github.com/microsoftgraph/msgraph-sdk-go-core"
)

// listAllResourcesWithPageIterator retrieves all Windows platform scripts using the PageIterator
// This handles pagination automatically and returns ALL scripts across all pages
func (r *WindowsPlatformScriptListResource) listAllResourcesWithPageIterator(
	ctx context.Context,
	requestConfig *devicemanagement.DeviceManagementScriptsRequestBuilderGetRequestConfiguration,
) ([]models.DeviceManagementScriptable, error) {
	var allScripts []models.DeviceManagementScriptable

	tflog.Debug(ctx, "Fetching first page of Windows platform scripts")

	scriptsResponse, err := r.client.
		DeviceManagement().
		DeviceManagementScripts().
		Get(ctx, requestConfig)

	if err != nil {
		return nil, fmt.Errorf("failed to get scripts: %w", err)
	}

	pageIterator, err := graphcore.NewPageIterator[models.DeviceManagementScriptable](
		scriptsResponse,
		r.client.GetAdapter(),
		models.CreateDeviceManagementScriptCollectionResponseFromDiscriminatorValue,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create page iterator: %w", err)
	}

	pageCount := 0
	err = pageIterator.Iterate(ctx, func(item models.DeviceManagementScriptable) bool {
		if item != nil {
			allScripts = append(allScripts, item)

			if len(allScripts)%100 == 0 {
				pageCount++
				tflog.Debug(ctx, fmt.Sprintf("PageIterator: collected %d scripts (estimated page %d)", len(allScripts), pageCount))
			}
		}
		return true
	})

	if err != nil {
		return nil, fmt.Errorf("failed to iterate pages: %w", err)
	}

	tflog.Debug(ctx, fmt.Sprintf("PageIterator complete: collected %d total scripts", len(allScripts)))

	return allScripts, nil
}
