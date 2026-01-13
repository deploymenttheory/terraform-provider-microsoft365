// REF: https://learn.microsoft.com/en-us/graph/api/intune-devices-windowsmanageddevice-list?view=graph-rest-beta
package graphBetaManagedDevice

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	abstractions "github.com/microsoft/kiota-abstractions-go"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/devicemanagement"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
	graphcore "github.com/microsoftgraph/msgraph-sdk-go-core"
)

func (d *ManagedDeviceDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var object ManagedDeviceDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &object)...)
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

	filteredItems := []ManagedDeviceDeviceDataItemModel{}
	filterValue := object.FilterValue.ValueString()

	switch filterType {
	case "id":
		managedDevice, err := d.client.
			DeviceManagement().
			ManagedDevices().
			ByManagedDeviceId(filterValue).Get(ctx, nil)

		if err != nil {
			errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationRead, d.ReadPermissions)
			return
		}

		mdItem := MapRemoteStateToDataSource(managedDevice)
		filteredItems = append(filteredItems, mdItem)

	case "odata":
		headers := abstractions.NewRequestHeaders()
		headers.Add("ConsistencyLevel", "eventual")

		requestParameters := d.buildODataRequestParameters(ctx, &object, headers)

		tflog.Debug(ctx, "Using Microsoft Graph SDK PageIterator for managed devices")

		allManagedDevices, err := d.listAllManagedDevicesWithPageIterator(ctx, requestParameters)
		if err != nil {
			tflog.Error(ctx, fmt.Sprintf("Error in OData query with pagination: %v", err))
			errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationRead, d.ReadPermissions)
			return
		}

		tflog.Debug(ctx, fmt.Sprintf("PageIterator returned %d results", len(allManagedDevices)))

		for _, managedDevice := range allManagedDevices {
			mdItem := MapRemoteStateToDataSource(managedDevice)
			filteredItems = append(filteredItems, mdItem)
		}

	default:
		tflog.Debug(ctx, "Using Microsoft Graph SDK PageIterator for managed devices (all/basic filter)")

		allManagedDevices, err := d.listAllManagedDevicesWithPageIterator(ctx, nil)
		if err != nil {
			errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationRead, d.ReadPermissions)
			return
		}

		for _, managedDevice := range allManagedDevices {
			mdItem := MapRemoteStateToDataSource(managedDevice)

			switch filterType {
			case "all":
				filteredItems = append(filteredItems, mdItem)

			case "device_name":
				if managedDevice.GetDeviceName() != nil && strings.Contains(
					strings.ToLower(*managedDevice.GetDeviceName()),
					strings.ToLower(filterValue)) {
					filteredItems = append(filteredItems, mdItem)
				}
			case "serial_number":
				if managedDevice.GetSerialNumber() != nil && strings.Contains(
					strings.ToLower(*managedDevice.GetSerialNumber()),
					strings.ToLower(filterValue)) {
					filteredItems = append(filteredItems, mdItem)
				}
			case "user_id":
				if managedDevice.GetUserId() != nil && strings.Contains(
					strings.ToLower(*managedDevice.GetUserId()),
					strings.ToLower(filterValue)) {
					filteredItems = append(filteredItems, mdItem)
				}
			}
		}
	}

	object.Items = filteredItems

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Datasource Read Method: %s, found %d items", DataSourceName, len(filteredItems)))
}

func (d *ManagedDeviceDataSource) buildODataRequestParameters(ctx context.Context, object *ManagedDeviceDataSourceModel, headers *abstractions.RequestHeaders) *devicemanagement.ManagedDevicesRequestBuilderGetRequestConfiguration {
	requestParameters := &devicemanagement.ManagedDevicesRequestBuilderGetRequestConfiguration{
		Headers:         headers,
		QueryParameters: &devicemanagement.ManagedDevicesRequestBuilderGetQueryParameters{},
	}

	if !object.ODataFilter.IsNull() && object.ODataFilter.ValueString() != "" {
		filter := object.ODataFilter.ValueString()
		requestParameters.QueryParameters.Filter = &filter
		tflog.Debug(ctx, fmt.Sprintf("Setting OData filter: %s", filter))
	}

	if !object.ODataTop.IsNull() {
		topValue := object.ODataTop.ValueInt32()
		requestParameters.QueryParameters.Top = &topValue
		tflog.Debug(ctx, fmt.Sprintf("Setting OData top: %d", topValue))
	}

	if !object.ODataSkip.IsNull() {
		skipValue := object.ODataSkip.ValueInt32()
		requestParameters.QueryParameters.Skip = &skipValue
		tflog.Debug(ctx, fmt.Sprintf("Setting OData skip: %d", skipValue))
	}

	if !object.ODataSelect.IsNull() && object.ODataSelect.ValueString() != "" {
		selectFields := strings.Split(object.ODataSelect.ValueString(), ",")
		for i, field := range selectFields {
			selectFields[i] = strings.TrimSpace(field)
		}
		requestParameters.QueryParameters.Select = selectFields
		tflog.Debug(ctx, fmt.Sprintf("Setting OData select: %v", selectFields))
	}

	if !object.ODataOrderBy.IsNull() && object.ODataOrderBy.ValueString() != "" {
		orderbyFields := strings.Split(object.ODataOrderBy.ValueString(), ",")
		for i, field := range orderbyFields {
			orderbyFields[i] = strings.TrimSpace(field)
		}
		requestParameters.QueryParameters.Orderby = orderbyFields
		tflog.Debug(ctx, fmt.Sprintf("Setting OData orderby: %v", orderbyFields))
	}

	if !object.ODataCount.IsNull() && object.ODataCount.ValueBool() {
		count := true
		requestParameters.QueryParameters.Count = &count
		tflog.Debug(ctx, "Setting OData count: true")
	}

	if !object.ODataSearch.IsNull() && object.ODataSearch.ValueString() != "" {
		search := object.ODataSearch.ValueString()
		requestParameters.QueryParameters.Search = &search
		tflog.Debug(ctx, fmt.Sprintf("Setting OData search: %s", search))
	}

	if !object.ODataExpand.IsNull() && object.ODataExpand.ValueString() != "" {
		expandFields := strings.Split(object.ODataExpand.ValueString(), ",")
		for i, field := range expandFields {
			expandFields[i] = strings.TrimSpace(field)
		}
		requestParameters.QueryParameters.Expand = expandFields
		tflog.Debug(ctx, fmt.Sprintf("Setting OData expand: %v", expandFields))
	}

	return requestParameters
}

func (d *ManagedDeviceDataSource) listAllManagedDevicesWithPageIterator(ctx context.Context, requestParameters *devicemanagement.ManagedDevicesRequestBuilderGetRequestConfiguration) ([]graphmodels.ManagedDeviceable, error) {
	var allManagedDevices []graphmodels.ManagedDeviceable

	managedDevicesResponse, err := d.client.
		DeviceManagement().
		ManagedDevices().
		Get(ctx, requestParameters)
	if err != nil {
		return nil, err
	}

	pageIterator, err := graphcore.NewPageIterator[graphmodels.ManagedDeviceable](
		managedDevicesResponse,
		d.client.GetAdapter(),
		graphmodels.CreateManagedDeviceCollectionResponseFromDiscriminatorValue,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create page iterator: %w", err)
	}

	pageCount := 0
	err = pageIterator.Iterate(ctx, func(item graphmodels.ManagedDeviceable) bool {
		if item != nil {
			allManagedDevices = append(allManagedDevices, item)

			if len(allManagedDevices)%25 == 0 {
				pageCount++
				tflog.Debug(ctx, fmt.Sprintf("PageIterator: collected %d managed devices (estimated page %d)", len(allManagedDevices), pageCount))
			}
		}
		return true
	})

	if err != nil {
		return nil, fmt.Errorf("failed to iterate pages: %w", err)
	}

	tflog.Debug(ctx, fmt.Sprintf("PageIterator complete: collected %d total managed devices", len(allManagedDevices)))

	return allManagedDevices, nil
}
