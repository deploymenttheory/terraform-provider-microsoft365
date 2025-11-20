package graphSubscribedSkus

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Read handles the Read operation for Subscribed SKUs data source.
func (d *SubscribedSkusDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var object SubscribedSkusDataSourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting Read method for: %s", DataSourceName))

	resp.Diagnostics.Append(req.Config.Get(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Reading %s with filters - sku_id: %s, sku_part_number: %s, applies_to: %s",
		DataSourceName,
		object.SkuId.ValueString(),
		object.SkuPartNumber.ValueString(),
		object.AppliesTo.ValueString()))

	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Read, ReadTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	var filteredItems []SubscribedSkuModel

	// If we have a specific SKU ID, try to get it directly
	if !object.SkuId.IsNull() && object.SkuId.ValueString() != "" {
		skuId := object.SkuId.ValueString()

		respItem, err := d.client.
			SubscribedSkus().
			BySubscribedSkuId(skuId).
			Get(ctx, nil)

		if err != nil {
			errors.HandleKiotaGraphError(ctx, err, resp, "Read", d.ReadPermissions)
			return
		}

		mappedItem := MapRemoteStateToDataSource(respItem)

		// Apply additional filters if specified
		if shouldIncludeItem(mappedItem, object) {
			filteredItems = append(filteredItems, mappedItem)
		}
	} else {
		// Get all subscribed SKUs and filter locally
		respList, err := d.client.
			SubscribedSkus().
			Get(ctx, nil)

		if err != nil {
			errors.HandleKiotaGraphError(ctx, err, resp, "Read", d.ReadPermissions)
			return
		}

		for _, item := range respList.GetValue() {
			mappedItem := MapRemoteStateToDataSource(item)

			if shouldIncludeItem(mappedItem, object) {
				filteredItems = append(filteredItems, mappedItem)
			}
		}
	}

	// Convert the filtered items to a Terraform list
	subscribedSkusList, diags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: getSubscribedSkuObjectType()}, filteredItems)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	object.SubscribedSkus = subscribedSkusList
	object.ID = types.StringValue(fmt.Sprintf("subscribed-skus-%d", time.Now().Unix()))

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Datasource Read Method: %s, found %d items", DataSourceName, len(filteredItems)))
}

// shouldIncludeItem determines whether an item should be included based on the filters
func shouldIncludeItem(item SubscribedSkuModel, filters SubscribedSkusDataSourceModel) bool {
	// Filter by SKU part number if specified
	if !filters.SkuPartNumber.IsNull() && filters.SkuPartNumber.ValueString() != "" {
		filterValue := strings.ToLower(filters.SkuPartNumber.ValueString())
		itemValue := strings.ToLower(item.SkuPartNumber.ValueString())
		if !strings.Contains(itemValue, filterValue) {
			return false
		}
	}

	// Filter by applies_to if specified
	if !filters.AppliesTo.IsNull() && filters.AppliesTo.ValueString() != "" {
		filterValue := strings.ToLower(filters.AppliesTo.ValueString())
		itemValue := strings.ToLower(item.AppliesTo.ValueString())
		if itemValue != filterValue {
			return false
		}
	}

	return true
}
