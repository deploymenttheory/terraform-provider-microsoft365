package graphBetaCloudPcAuditEvent

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/devicemanagement"
)

func (d *CloudPcAuditEventDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var object CloudPcAuditEventDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &object)...) // Get config into object
	if resp.Diagnostics.HasError() {
		return
	}

	filterType := object.FilterType.ValueString()
	tflog.Debug(ctx, fmt.Sprintf("Starting Read method for datasource: %s with filter_type: %s", DataSourceName, filterType))

	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Read, ReadTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	var filteredItems []CloudPcAuditEventItemModel
	filterValue := object.FilterValue.ValueString()

	if filterType == "id" {
		requestParameters := &devicemanagement.VirtualEndpointAuditEventsCloudPcAuditEventItemRequestBuilderGetRequestConfiguration{}

		auditEvent, err := d.client.
			DeviceManagement().
			VirtualEndpoint().
			AuditEvents().
			ByCloudPcAuditEventId(filterValue).
			Get(ctx, requestParameters)

		if err != nil {
			errors.HandleGraphError(ctx, err, resp, "Read", d.ReadPermissions)
			return
		}

		item := MapRemoteStateToDataSource(ctx, auditEvent)
		filteredItems = append(filteredItems, item)

	} else {
		respList, err := d.client.
			DeviceManagement().
			VirtualEndpoint().
			AuditEvents().
			Get(ctx, nil)

		if err != nil {
			errors.HandleGraphError(ctx, err, resp, "Read", d.ReadPermissions)
			return
		}

		for _, auditEvent := range respList.GetValue() {
			item := MapRemoteStateToDataSource(ctx, auditEvent)
			switch filterType {
			case "all":
				filteredItems = append(filteredItems, item)
			case "display_name":
				if auditEvent.GetDisplayName() != nil && strings.Contains(
					strings.ToLower(*auditEvent.GetDisplayName()),
					strings.ToLower(filterValue)) {
					filteredItems = append(filteredItems, item)
				}
			}
		}
	}

	object.Items = filteredItems

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...) // Set state
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("Finished Datasource Read Method: %s, found %d items", DataSourceName, len(filteredItems)))
}
