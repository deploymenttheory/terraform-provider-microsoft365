package graphBetaDirectorySettingTemplates

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Read handles the Read operation for Directory Setting Templates data source.
func (d *DirectorySettingTemplatesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var object DirectorySettingTemplatesDataSourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting Read method for: %s", d.ProviderTypeName+"_"+d.TypeName))

	resp.Diagnostics.Append(req.Config.Get(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	filterType := object.FilterType.ValueString()
	tflog.Debug(ctx, fmt.Sprintf("Reading %s with filter_type: %s", d.ProviderTypeName+"_"+d.TypeName, filterType))

	if filterType != "all" && (object.FilterValue.IsNull() || object.FilterValue.ValueString() == "") {
		resp.Diagnostics.AddError(
			"Missing Required Parameter",
			fmt.Sprintf("filter_value must be provided when filter_type is '%s'", filterType),
		)
		return
	}

	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Read, ReadTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	var filteredItems []DirectorySettingTemplateModel
	filterValue := object.FilterValue.ValueString()

	// For ID filter, we can make a direct API call
	if filterType == "id" {
		respItem, err := d.client.
			DirectorySettingTemplates().
			ByDirectorySettingTemplateId(filterValue).
			Get(ctx, nil)

		if err != nil {
			errors.HandleGraphError(ctx, err, resp, "Read", d.ReadPermissions)
			return
		}

		filteredItems = append(filteredItems, MapRemoteStateToDataSource(ctx, respItem))
	} else {
		// For all other filters, we need to get all templates and filter locally
		respList, err := d.client.
			DirectorySettingTemplates().
			Get(ctx, nil)

		if err != nil {
			errors.HandleGraphError(ctx, err, resp, "Read", d.ReadPermissions)
			return
		}

		for _, item := range respList.GetValue() {
			switch filterType {
			case "all":
				filteredItems = append(filteredItems, MapRemoteStateToDataSource(ctx, item))

			case "display_name":
				if item.GetDisplayName() != nil && strings.Contains(
					strings.ToLower(*item.GetDisplayName()),
					strings.ToLower(filterValue)) {
					filteredItems = append(filteredItems, MapRemoteStateToDataSource(ctx, item))
				}
			}
		}
	}

	object.DirectorySettingTemplates = filteredItems
	object.ID = types.StringValue(fmt.Sprintf("directory-setting-templates-%d", time.Now().Unix()))

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Datasource Read Method: %s, found %d items", d.ProviderTypeName+"_"+d.TypeName, len(filteredItems)))
}
