package graphSubscribedSkus

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-sdk-go/models"
)

// Read handles the Read operation for Subscribed SKUs data source.
func (d *SubscribedSkusDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var object SubscribedSkusDataSourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting Read method for: %s", DataSourceName))

	resp.Diagnostics.Append(req.Config.Get(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Read, ReadTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	lookupMethod := d.determineLookupMethod(object)
	tflog.Debug(ctx, fmt.Sprintf("Determined lookup method: %s", lookupMethod))

	var skus []models.SubscribedSkuable
	var err error

	switch lookupMethod {
	case "sku_id":
		skus, err = d.getSkuById(ctx, object.SkuId.ValueString())
	case "sku_part_number":
		skus, err = d.getSkusByPartNumber(ctx, object.SkuPartNumber.ValueString())
	case "account_id":
		skus, err = d.getSkusByAccountId(ctx, object.AccountId.ValueString())
	case "account_name":
		skus, err = d.getSkusByAccountName(ctx, object.AccountName.ValueString())
	case "applies_to":
		skus, err = d.getSkusByAppliesTo(ctx, object.AppliesTo.ValueString())
	case "list_all":
		skus, err = d.listAllSkus(ctx)
	default:
		resp.Diagnostics.AddError(
			"Invalid Configuration",
			"No valid lookup method specified. Please provide one of: sku_id, sku_part_number, account_id, account_name, applies_to, or list_all.",
		)
		return
	}

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationRead, d.ReadPermissions)
		return
	}

	filteredItems := make([]SubscribedSkuModel, 0, len(skus))
	for _, sku := range skus {
		mappedItem := MapRemoteStateToDataSource(sku)
		filteredItems = append(filteredItems, mappedItem)
	}

	itemsList, diags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: getSubscribedSkuObjectType()}, filteredItems)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	object.Items = itemsList
	object.ID = types.StringValue(fmt.Sprintf("subscribed-skus-%d", time.Now().Unix()))

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Datasource Read Method: %s, found %d items", DataSourceName, len(filteredItems)))
}

// determineLookupMethod determines which lookup method to use based on the provided attributes
func (d *SubscribedSkusDataSource) determineLookupMethod(object SubscribedSkusDataSourceModel) string {
	if !object.SkuId.IsNull() && object.SkuId.ValueString() != "" {
		return "sku_id"
	}
	if !object.SkuPartNumber.IsNull() && object.SkuPartNumber.ValueString() != "" {
		return "sku_part_number"
	}
	if !object.AccountId.IsNull() && object.AccountId.ValueString() != "" {
		return "account_id"
	}
	if !object.AccountName.IsNull() && object.AccountName.ValueString() != "" {
		return "account_name"
	}
	if !object.AppliesTo.IsNull() && object.AppliesTo.ValueString() != "" {
		return "applies_to"
	}
	if !object.ListAll.IsNull() && object.ListAll.ValueBool() {
		return "list_all"
	}
	return ""
}

// getSkuById retrieves a specific SKU by its ID using the GET /subscribedSkus/{id} endpoint
func (d *SubscribedSkusDataSource) getSkuById(ctx context.Context, skuId string) ([]models.SubscribedSkuable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Fetching SKU by ID: %s", skuId))

	sku, err := d.client.
		SubscribedSkus().
		BySubscribedSkuId(skuId).
		Get(ctx, nil)

	if err != nil {
		return nil, err
	}

	return []models.SubscribedSkuable{sku}, nil
}

// listAllSkus retrieves all subscribed SKUs
func (d *SubscribedSkusDataSource) listAllSkus(ctx context.Context) ([]models.SubscribedSkuable, error) {
	tflog.Debug(ctx, "Fetching all subscribed SKUs")

	result, err := d.client.
		SubscribedSkus().
		Get(ctx, nil)

	if err != nil {
		return nil, err
	}

	return result.GetValue(), nil
}

// getSkusByPartNumber retrieves SKUs filtered by part number (local filtering)
func (d *SubscribedSkusDataSource) getSkusByPartNumber(ctx context.Context, partNumber string) ([]models.SubscribedSkuable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Fetching SKUs by part number: %s", partNumber))

	allSkus, err := d.listAllSkus(ctx)
	if err != nil {
		return nil, err
	}

	filtered := make([]models.SubscribedSkuable, 0)
	searchValue := strings.ToLower(partNumber)

	for _, sku := range allSkus {
		if sku.GetSkuPartNumber() != nil {
			skuPartNumber := strings.ToLower(*sku.GetSkuPartNumber())
			if strings.Contains(skuPartNumber, searchValue) {
				filtered = append(filtered, sku)
			}
		}
	}

	return filtered, nil
}

// getSkusByAccountId retrieves SKUs filtered by account ID (local filtering)
func (d *SubscribedSkusDataSource) getSkusByAccountId(ctx context.Context, accountId string) ([]models.SubscribedSkuable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Fetching SKUs by account ID: %s", accountId))

	allSkus, err := d.listAllSkus(ctx)
	if err != nil {
		return nil, err
	}

	filtered := make([]models.SubscribedSkuable, 0)
	searchValue := strings.ToLower(accountId)

	for _, sku := range allSkus {
		if sku.GetAccountId() != nil {
			skuAccountId := strings.ToLower(*sku.GetAccountId())
			if strings.Contains(skuAccountId, searchValue) {
				filtered = append(filtered, sku)
			}
		}
	}

	return filtered, nil
}

// getSkusByAccountName retrieves SKUs filtered by account name (local filtering)
func (d *SubscribedSkusDataSource) getSkusByAccountName(ctx context.Context, accountName string) ([]models.SubscribedSkuable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Fetching SKUs by account name: %s", accountName))

	allSkus, err := d.listAllSkus(ctx)
	if err != nil {
		return nil, err
	}

	filtered := make([]models.SubscribedSkuable, 0)
	searchValue := strings.ToLower(accountName)

	for _, sku := range allSkus {
		if sku.GetAccountName() != nil {
			skuAccountName := strings.ToLower(*sku.GetAccountName())
			if strings.Contains(skuAccountName, searchValue) {
				filtered = append(filtered, sku)
			}
		}
	}

	return filtered, nil
}

// getSkusByAppliesTo retrieves SKUs filtered by applies_to (local filtering)
func (d *SubscribedSkusDataSource) getSkusByAppliesTo(ctx context.Context, appliesTo string) ([]models.SubscribedSkuable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Fetching SKUs by applies_to: %s", appliesTo))

	allSkus, err := d.listAllSkus(ctx)
	if err != nil {
		return nil, err
	}

	filtered := make([]models.SubscribedSkuable, 0)
	searchValue := strings.ToLower(appliesTo)

	for _, sku := range allSkus {
		if sku.GetAppliesTo() != nil {
			skuAppliesTo := strings.ToLower(*sku.GetAppliesTo())
			if skuAppliesTo == searchValue {
				filtered = append(filtered, sku)
			}
		}
	}

	return filtered, nil
}
