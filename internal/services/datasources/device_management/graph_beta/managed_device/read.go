// REF: https://learn.microsoft.com/en-us/graph/api/intune-devices-windowsmanageddevice-list?view=graph-rest-beta
package graphBetaManagedDevice

import (
	"context"
	"fmt"
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

// lookupMethod represents the different ways to look up managed devices
type lookupMethod int

const (
	lookupByODataQuery lookupMethod = iota
	lookupByDeviceId
	lookupByDeviceName
	lookupByOperatingSystem
	lookupByAzureADDeviceId
	lookupBySerialNumber
	lookupByUserPrincipalName
	lookupListAll
)

// Read handles the Read operation.
func (d *ManagedDeviceDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var object ManagedDeviceDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Starting Read method for datasource: %s", DataSourceName))

	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Read, ReadTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	var devices []graphmodels.ManagedDeviceable
	var err error

	method := determineLookupMethod(object)
	switch method {
	case lookupByODataQuery:
		devices, err = d.getDevicesByODataQuery(ctx, object)
	case lookupByDeviceId:
		devices, err = d.getDeviceById(ctx, object)
	case lookupByDeviceName:
		devices, err = d.getDevicesByDeviceName(ctx, object)
	case lookupByOperatingSystem:
		devices, err = d.getDevicesByOperatingSystem(ctx, object)
	case lookupByAzureADDeviceId:
		devices, err = d.getDevicesByAzureADDeviceId(ctx, object)
	case lookupBySerialNumber:
		devices, err = d.getDevicesBySerialNumber(ctx, object)
	case lookupByUserPrincipalName:
		devices, err = d.getDevicesByUserPrincipalName(ctx, object)
	case lookupListAll:
		devices, err = d.listAllDevices(ctx)
	}

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationRead, d.ReadPermissions)
		return
	}

	if len(devices) == 0 {
		resp.Diagnostics.AddWarning(
			"No Managed Devices Found",
			"The lookup did not return any managed devices matching the specified criteria.",
		)
	}

	tflog.Debug(ctx, fmt.Sprintf("Successfully found %d managed device(s)", len(devices)))

	// Map devices to items
	filteredItems := make([]ManagedDeviceDeviceDataItemModel, 0, len(devices))
	for _, device := range devices {
		mdItem := MapRemoteStateToDataSource(device)
		filteredItems = append(filteredItems, mdItem)
	}

	object.Items = filteredItems

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Datasource Read Method: %s", DataSourceName))
}

// determineLookupMethod determines which lookup method to use based on provided attributes
func determineLookupMethod(object ManagedDeviceDataSourceModel) lookupMethod {
	switch {
	case !object.ODataQuery.IsNull() && object.ODataQuery.ValueString() != "":
		return lookupByODataQuery
	case !object.DeviceId.IsNull() && object.DeviceId.ValueString() != "":
		return lookupByDeviceId
	case !object.DeviceName.IsNull() && object.DeviceName.ValueString() != "":
		return lookupByDeviceName
	case !object.OperatingSystem.IsNull() && object.OperatingSystem.ValueString() != "":
		return lookupByOperatingSystem
	case !object.AzureADDeviceId.IsNull() && object.AzureADDeviceId.ValueString() != "":
		return lookupByAzureADDeviceId
	case !object.SerialNumber.IsNull() && object.SerialNumber.ValueString() != "":
		return lookupBySerialNumber
	case !object.UserPrincipalName.IsNull() && object.UserPrincipalName.ValueString() != "":
		return lookupByUserPrincipalName
	case !object.ListAll.IsNull() && object.ListAll.ValueBool():
		return lookupListAll
	default:
		return lookupByDeviceId // This should never happen due to schema validators
	}
}

// getDeviceById retrieves a managed device by its device ID
func (d *ManagedDeviceDataSource) getDeviceById(ctx context.Context, object ManagedDeviceDataSourceModel) ([]graphmodels.ManagedDeviceable, error) {
	deviceId := object.DeviceId.ValueString()

	tflog.Debug(ctx, fmt.Sprintf("Looking up managed device by ID: %s", deviceId))

	device, err := d.client.
		DeviceManagement().
		ManagedDevices().
		ByManagedDeviceId(deviceId).
		Get(ctx, nil)

	if err != nil {
		return nil, err
	}

	if device == nil {
		return []graphmodels.ManagedDeviceable{}, nil
	}

	return []graphmodels.ManagedDeviceable{device}, nil
}

// getDevicesByDeviceName retrieves managed devices by device name using OData filter
func (d *ManagedDeviceDataSource) getDevicesByDeviceName(ctx context.Context, object ManagedDeviceDataSourceModel) ([]graphmodels.ManagedDeviceable, error) {
	deviceName := object.DeviceName.ValueString()

	tflog.Debug(ctx, fmt.Sprintf("Looking up managed devices by device name: %s", deviceName))

	// Build OData filter for device name
	filter := fmt.Sprintf("deviceName eq '%s'", deviceName)

	headers := abstractions.NewRequestHeaders()
	headers.Add("ConsistencyLevel", "eventual")

	requestConfig := &devicemanagement.ManagedDevicesRequestBuilderGetRequestConfiguration{
		Headers: headers,
		QueryParameters: &devicemanagement.ManagedDevicesRequestBuilderGetQueryParameters{
			Filter: &filter,
		},
	}

	return d.listAllManagedDevicesWithPageIterator(ctx, requestConfig)
}

// getDevicesByOperatingSystem retrieves managed devices by operating system and optionally OS version
func (d *ManagedDeviceDataSource) getDevicesByOperatingSystem(ctx context.Context, object ManagedDeviceDataSourceModel) ([]graphmodels.ManagedDeviceable, error) {
	os := object.OperatingSystem.ValueString()
	osVersion := object.OsVersion.ValueString()

	tflog.Debug(ctx, fmt.Sprintf("Looking up managed devices by OS: %s, version: %s", os, osVersion))

	// Build OData filter
	var filter string
	if osVersion != "" {
		filter = fmt.Sprintf("operatingSystem eq '%s' and osVersion eq '%s'", os, osVersion)
	} else {
		filter = fmt.Sprintf("operatingSystem eq '%s'", os)
	}

	headers := abstractions.NewRequestHeaders()
	headers.Add("ConsistencyLevel", "eventual")

	requestConfig := &devicemanagement.ManagedDevicesRequestBuilderGetRequestConfiguration{
		Headers: headers,
		QueryParameters: &devicemanagement.ManagedDevicesRequestBuilderGetQueryParameters{
			Filter: &filter,
		},
	}

	return d.listAllManagedDevicesWithPageIterator(ctx, requestConfig)
}

// getDevicesByAzureADDeviceId retrieves managed devices by Azure AD device ID
func (d *ManagedDeviceDataSource) getDevicesByAzureADDeviceId(ctx context.Context, object ManagedDeviceDataSourceModel) ([]graphmodels.ManagedDeviceable, error) {
	azureADDeviceId := object.AzureADDeviceId.ValueString()

	tflog.Debug(ctx, fmt.Sprintf("Looking up managed devices by Azure AD device ID: %s", azureADDeviceId))

	// Build OData filter for Azure AD device ID
	filter := fmt.Sprintf("azureADDeviceId eq '%s'", azureADDeviceId)

	headers := abstractions.NewRequestHeaders()
	headers.Add("ConsistencyLevel", "eventual")

	requestConfig := &devicemanagement.ManagedDevicesRequestBuilderGetRequestConfiguration{
		Headers: headers,
		QueryParameters: &devicemanagement.ManagedDevicesRequestBuilderGetQueryParameters{
			Filter: &filter,
		},
	}

	return d.listAllManagedDevicesWithPageIterator(ctx, requestConfig)
}

// getDevicesBySerialNumber retrieves managed devices by serial number
func (d *ManagedDeviceDataSource) getDevicesBySerialNumber(ctx context.Context, object ManagedDeviceDataSourceModel) ([]graphmodels.ManagedDeviceable, error) {
	serialNumber := object.SerialNumber.ValueString()

	tflog.Debug(ctx, fmt.Sprintf("Looking up managed devices by serial number: %s", serialNumber))

	// Build OData filter for serial number
	filter := fmt.Sprintf("serialNumber eq '%s'", serialNumber)

	headers := abstractions.NewRequestHeaders()
	headers.Add("ConsistencyLevel", "eventual")

	requestConfig := &devicemanagement.ManagedDevicesRequestBuilderGetRequestConfiguration{
		Headers: headers,
		QueryParameters: &devicemanagement.ManagedDevicesRequestBuilderGetQueryParameters{
			Filter: &filter,
		},
	}

	return d.listAllManagedDevicesWithPageIterator(ctx, requestConfig)
}

// getDevicesByUserPrincipalName retrieves managed devices by user principal name
func (d *ManagedDeviceDataSource) getDevicesByUserPrincipalName(ctx context.Context, object ManagedDeviceDataSourceModel) ([]graphmodels.ManagedDeviceable, error) {
	userPrincipalName := object.UserPrincipalName.ValueString()

	tflog.Debug(ctx, fmt.Sprintf("Looking up managed devices by user principal name: %s", userPrincipalName))

	filter := fmt.Sprintf("userPrincipalName eq '%s'", userPrincipalName)

	headers := abstractions.NewRequestHeaders()
	headers.Add("ConsistencyLevel", "eventual")

	requestConfig := &devicemanagement.ManagedDevicesRequestBuilderGetRequestConfiguration{
		Headers: headers,
		QueryParameters: &devicemanagement.ManagedDevicesRequestBuilderGetQueryParameters{
			Filter: &filter,
		},
	}

	return d.listAllManagedDevicesWithPageIterator(ctx, requestConfig)
}

// getDevicesByODataQuery retrieves managed devices using a custom OData query
func (d *ManagedDeviceDataSource) getDevicesByODataQuery(ctx context.Context, object ManagedDeviceDataSourceModel) ([]graphmodels.ManagedDeviceable, error) {
	filter := object.ODataQuery.ValueString()

	tflog.Debug(ctx, fmt.Sprintf("Looking up managed devices with OData query: %s", filter))

	headers := abstractions.NewRequestHeaders()
	headers.Add("ConsistencyLevel", "eventual")

	requestConfig := &devicemanagement.ManagedDevicesRequestBuilderGetRequestConfiguration{
		Headers: headers,
		QueryParameters: &devicemanagement.ManagedDevicesRequestBuilderGetQueryParameters{
			Filter: &filter,
		},
	}

	return d.listAllManagedDevicesWithPageIterator(ctx, requestConfig)
}

// listAllDevices retrieves all managed devices
func (d *ManagedDeviceDataSource) listAllDevices(ctx context.Context) ([]graphmodels.ManagedDeviceable, error) {
	tflog.Debug(ctx, "Listing all managed devices")

	headers := abstractions.NewRequestHeaders()
	headers.Add("ConsistencyLevel", "eventual")

	requestConfig := &devicemanagement.ManagedDevicesRequestBuilderGetRequestConfiguration{
		Headers: headers,
	}

	return d.listAllManagedDevicesWithPageIterator(ctx, requestConfig)
}

// listAllManagedDevicesWithPageIterator uses Microsoft Graph SDK's PageIterator to handle pagination
func (d *ManagedDeviceDataSource) listAllManagedDevicesWithPageIterator(
	ctx context.Context,
	requestConfig *devicemanagement.ManagedDevicesRequestBuilderGetRequestConfiguration,
) ([]graphmodels.ManagedDeviceable, error) {
	var allDevices []graphmodels.ManagedDeviceable

	resp, err := d.client.
		DeviceManagement().
		ManagedDevices().
		Get(ctx, requestConfig)

	if err != nil {
		return nil, fmt.Errorf("failed to get initial page of managed devices: %w", err)
	}

	pageIterator, err := graphcore.NewPageIterator[graphmodels.ManagedDeviceable](
		resp,
		d.client.GetAdapter(),
		graphmodels.CreateManagedDeviceCollectionResponseFromDiscriminatorValue,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create page iterator: %w", err)
	}

	err = pageIterator.Iterate(ctx, func(device graphmodels.ManagedDeviceable) bool {
		allDevices = append(allDevices, device)
		return true
	})

	if err != nil {
		return nil, fmt.Errorf("error during pagination: %w", err)
	}

	return allDevices, nil
}
