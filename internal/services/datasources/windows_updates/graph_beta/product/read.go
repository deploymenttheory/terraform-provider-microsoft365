package graphBetaWindowsUpdateProduct

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphadmin "github.com/microsoftgraph/msgraph-beta-sdk-go/admin"
)

func (d *WindowsUpdateProductDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var object WindowsUpdateProductDataSourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting Read method for: %s", DataSourceName))

	resp.Diagnostics.Append(req.Config.Get(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	searchType := object.SearchType.ValueString()
	searchValue := object.SearchValue.ValueString()

	tflog.Debug(ctx, fmt.Sprintf("Reading %s with search_type: %s, search_value: %s", DataSourceName, searchType, searchValue))

	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Read, ReadTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	var products []WindowsUpdateProduct

	switch searchType {
	case "catalog_id":
		requestParameters := &graphadmin.WindowsUpdatesProductsMicrosoftGraphWindowsUpdatesFindByCatalogIdWithCatalogIDRequestBuilderGetQueryParameters{
			Expand: []string{"revisions($expand=catalogEntry,knowledgeBaseArticle)", "knownIssues($expand=originatingKnowledgeBaseArticle,resolvingKnowledgeBaseArticle)"},
		}
		configuration := &graphadmin.WindowsUpdatesProductsMicrosoftGraphWindowsUpdatesFindByCatalogIdWithCatalogIDRequestBuilderGetRequestConfiguration{
			QueryParameters: requestParameters,
		}

		respList, err := d.client.
			Admin().
			Windows().
			Updates().
			Products().
			MicrosoftGraphWindowsUpdatesFindByCatalogIdWithCatalogID(&searchValue).
			GetAsFindByCatalogIdWithCatalogIDGetResponse(ctx, configuration)

		if err != nil {
			errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationRead, d.ReadPermissions)
			return
		}

		if respList != nil && respList.GetValue() != nil {
			for _, product := range respList.GetValue() {
				products = append(products, MapRemoteStateToDataSource(product))
			}
		}

	case "kb_number":
		kbNumber, err := strconv.ParseInt(searchValue, 10, 32)
		if err != nil {
			resp.Diagnostics.AddError(
				"Invalid KB Number",
				fmt.Sprintf("search_value must be a valid integer when search_type is 'kb_number'. Error: %s", err.Error()),
			)
			return
		}

		kbNumber32 := int32(kbNumber)

		requestParameters := &graphadmin.WindowsUpdatesProductsMicrosoftGraphWindowsUpdatesFindByKbNumberWithKbNumberRequestBuilderGetQueryParameters{
			Expand: []string{"revisions($expand=catalogEntry,knowledgeBaseArticle)", "knownIssues($expand=originatingKnowledgeBaseArticle,resolvingKnowledgeBaseArticle)"},
		}
		configuration := &graphadmin.WindowsUpdatesProductsMicrosoftGraphWindowsUpdatesFindByKbNumberWithKbNumberRequestBuilderGetRequestConfiguration{
			QueryParameters: requestParameters,
		}

		respList, err := d.client.
			Admin().
			Windows().
			Updates().
			Products().
			MicrosoftGraphWindowsUpdatesFindByKbNumberWithKbNumber(&kbNumber32).
			GetAsFindByKbNumberWithKbNumberGetResponse(ctx, configuration)

		if err != nil {
			errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationRead, d.ReadPermissions)
			return
		}

		if respList != nil && respList.GetValue() != nil {
			for _, product := range respList.GetValue() {
				products = append(products, MapRemoteStateToDataSource(product))
			}
		}

	default:
		resp.Diagnostics.AddError(
			"Invalid Search Type",
			fmt.Sprintf("Unsupported search_type: %s", searchType),
		)
		return
	}

	object.Products = products

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Datasource Read Method: %s, found %d products", DataSourceName, len(products)))
}
