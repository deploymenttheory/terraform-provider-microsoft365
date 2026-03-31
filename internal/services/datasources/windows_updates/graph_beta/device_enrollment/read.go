package graphBetaWindowsUpdatesDeviceEnrollment

import (
	"context"
	"fmt"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	abstractions "github.com/microsoft/kiota-abstractions-go"
	graphadmin "github.com/microsoftgraph/msgraph-beta-sdk-go/admin"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models/windowsupdates"
	graphcore "github.com/microsoftgraph/msgraph-sdk-go-core"
)

// lookupMethod represents the different ways to look up device enrollment
type lookupMethod int

const (
	lookupByEntraDeviceId lookupMethod = iota
	lookupByDeviceName
	lookupListAll
)

func (d *DeviceEnrollmentDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config DeviceEnrollmentDataSourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting Read method for: %s", DataSourceName))

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var object DeviceEnrollmentDataSourceModel
	object.EntraDeviceId = config.EntraDeviceId
	object.DeviceName = config.DeviceName
	object.ListAll = config.ListAll
	object.UpdateCategory = config.UpdateCategory
	object.ODataFilter = config.ODataFilter
	object.Timeouts = config.Timeouts

	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Read, ReadTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	var devices []EnrolledDevice
	var err error

	method := determineLookupMethod(object)
	switch method {
	case lookupByEntraDeviceId:
		devices, err = d.getDeviceByEntraId(ctx, object)
	case lookupByDeviceName:
		devices, err = d.getDeviceByName(ctx, object)
	case lookupListAll:
		devices, err = d.listAllDevices(ctx, object)
	}

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationRead, d.ReadPermissions)
		return
	}

	// Apply update_category filter if specified
	if !object.UpdateCategory.IsNull() && object.UpdateCategory.ValueString() != "" {
		devices = filterDevicesByUpdateCategory(devices, object.UpdateCategory.ValueString())
	}

	object.Devices = devices

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Datasource Read Method: %s, found %d device(s)", DataSourceName, len(devices)))
}

// determineLookupMethod determines which lookup method to use based on provided attributes
func determineLookupMethod(object DeviceEnrollmentDataSourceModel) lookupMethod {
	switch {
	case !object.EntraDeviceId.IsNull() && object.EntraDeviceId.ValueString() != "":
		return lookupByEntraDeviceId
	case !object.DeviceName.IsNull() && object.DeviceName.ValueString() != "":
		return lookupByDeviceName
	case !object.ListAll.IsNull() && object.ListAll.ValueBool():
		return lookupListAll
	default:
		return lookupByEntraDeviceId // This should never happen due to schema validators
	}
}

// getDeviceByEntraId retrieves a device by its Entra ID
func (d *DeviceEnrollmentDataSource) getDeviceByEntraId(ctx context.Context, object DeviceEnrollmentDataSourceModel) ([]EnrolledDevice, error) {
	deviceId := object.EntraDeviceId.ValueString()
	tflog.Debug(ctx, fmt.Sprintf("Looking up device by Entra ID: %s", deviceId))

	maxRetries := 6
	retryDelay := 10 * time.Second

	var device graphmodels.UpdatableAssetable
	var err error

	for attempt := 1; attempt <= maxRetries; attempt++ {
		device, err = d.client.
			Admin().
			Windows().
			Updates().
			UpdatableAssets().
			ByUpdatableAssetId(deviceId).
			Get(ctx, nil)

		if err == nil && device != nil && device.GetId() != nil {
			return []EnrolledDevice{mapDeviceToModel(device)}, nil
		}

		if err != nil {
			errorInfo := errors.GraphError(ctx, err)
			if errors.IsNonRetryableReadError(&errorInfo) {
				return nil, err
			}
		}

		if attempt < maxRetries {
			tflog.Debug(ctx, fmt.Sprintf("Device not found, retrying in %v (attempt %d/%d)", retryDelay, attempt, maxRetries))
			time.Sleep(retryDelay)
		}
	}

	if err != nil {
		return nil, fmt.Errorf("device not found after %d attempts: %w", maxRetries, err)
	}
	return nil, fmt.Errorf("device not found or invalid after %d attempts", maxRetries)
}

// getDeviceByName retrieves a device by resolving its name to an Entra ID first
func (d *DeviceEnrollmentDataSource) getDeviceByName(ctx context.Context, object DeviceEnrollmentDataSourceModel) ([]EnrolledDevice, error) {
	deviceName := object.DeviceName.ValueString()
	tflog.Debug(ctx, fmt.Sprintf("Looking up device by name: %s", deviceName))

	// First, we need to query the device management API to get the Entra device ID
	// This requires calling the managed devices endpoint and filtering by display name
	managedDevices, err := d.client.
		DeviceManagement().
		ManagedDevices().
		Get(ctx, nil)

	if err != nil {
		return nil, fmt.Errorf("failed to query managed devices: %w", err)
	}

	var entraDeviceId string
	if managedDevices != nil && managedDevices.GetValue() != nil {
		for _, device := range managedDevices.GetValue() {
			if device.GetDeviceName() != nil && *device.GetDeviceName() == deviceName {
				if device.GetAzureADDeviceId() != nil {
					entraDeviceId = *device.GetAzureADDeviceId()
					break
				}
			}
		}
	}

	if entraDeviceId == "" {
		return nil, fmt.Errorf("no device found with name: %s", deviceName)
	}

	tflog.Debug(ctx, fmt.Sprintf("Resolved device name '%s' to Entra ID: %s", deviceName, entraDeviceId))

	// Now fetch the enrollment status using the Entra device ID
	object.EntraDeviceId = types.StringValue(entraDeviceId)
	return d.getDeviceByEntraId(ctx, object)
}

// listAllDevices retrieves all enrolled devices with optional filtering
func (d *DeviceEnrollmentDataSource) listAllDevices(ctx context.Context, object DeviceEnrollmentDataSourceModel) ([]EnrolledDevice, error) {
	tflog.Debug(ctx, "Listing all enrolled devices")

	headers := abstractions.NewRequestHeaders()
	headers.Add("ConsistencyLevel", "eventual")

	requestConfig := &graphadmin.WindowsUpdatesUpdatableAssetsRequestBuilderGetRequestConfiguration{
		Headers: headers,
		QueryParameters: &graphadmin.WindowsUpdatesUpdatableAssetsRequestBuilderGetQueryParameters{
			Filter: nil,
		},
	}

	// Apply custom OData filter if provided
	if !object.ODataFilter.IsNull() && object.ODataFilter.ValueString() != "" {
		filter := object.ODataFilter.ValueString()
		requestConfig.QueryParameters.Filter = &filter
	}

	maxRetries := 6
	retryDelay := 10 * time.Second

	for attempt := 1; attempt <= maxRetries; attempt++ {
		devicesResponse, err := d.client.
			Admin().
			Windows().
			Updates().
			UpdatableAssets().
			Get(ctx, requestConfig)

		if err != nil {
			errorInfo := errors.GraphError(ctx, err)
			if errors.IsNonRetryableReadError(&errorInfo) {
				return nil, err
			}
			if attempt < maxRetries {
				tflog.Debug(ctx, fmt.Sprintf("Error listing devices, retrying in %v (attempt %d/%d)", retryDelay, attempt, maxRetries))
				time.Sleep(retryDelay)
				continue
			}
			return nil, err
		}

		var devices []EnrolledDevice

		pageIterator, err := graphcore.NewPageIterator[graphmodels.UpdatableAssetable](
			devicesResponse,
			d.client.GetAdapter(),
			graphmodels.CreateUpdatableAssetCollectionResponseFromDiscriminatorValue,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to create page iterator: %w", err)
		}

		err = pageIterator.Iterate(ctx, func(item graphmodels.UpdatableAssetable) bool {
			if item != nil && item.GetId() != nil {
				devices = append(devices, mapDeviceToModel(item))
			}
			return true
		})

		if err != nil {
			return nil, fmt.Errorf("failed to iterate device pages: %w", err)
		}

		return devices, nil
	}

	return nil, fmt.Errorf("failed to list devices after %d attempts", maxRetries)
}

// filterDevicesByUpdateCategory filters devices by a specific update category
func filterDevicesByUpdateCategory(devices []EnrolledDevice, category string) []EnrolledDevice {
	var filtered []EnrolledDevice

	for _, device := range devices {
		hasCategory := false
		for _, enrollment := range device.Enrollments {
			if enrollment.UpdateCategory.ValueString() == category {
				hasCategory = true
				break
			}
		}
		if hasCategory {
			filtered = append(filtered, device)
		}
	}

	return filtered
}

// mapDeviceToModel maps a Graph API device to the Terraform model
func mapDeviceToModel(device graphmodels.UpdatableAssetable) EnrolledDevice {
	enrolledDevice := EnrolledDevice{
		ID:          types.StringValue(*device.GetId()),
		Enrollments: []UpdateManagementEnrollment{},
		Errors:      []UpdatableAssetError{},
	}

	if azureADDevice, ok := device.(graphmodels.AzureADDeviceable); ok {
		if enrollment := azureADDevice.GetEnrollment(); enrollment != nil {
			// Check for driver updateenrollment
			if driverEnrollment := enrollment.GetDriver(); driverEnrollment != nil {
				enrolledDevice.Enrollments = append(enrolledDevice.Enrollments, UpdateManagementEnrollment{
					UpdateCategory: types.StringValue("driver"),
				})
			}
			// Check for feature update enrollment
			if featureEnrollment := enrollment.GetFeature(); featureEnrollment != nil {
				enrolledDevice.Enrollments = append(enrolledDevice.Enrollments, UpdateManagementEnrollment{
					UpdateCategory: types.StringValue("feature"),
				})
			}
			// Check for quality update enrollment
			if qualityEnrollment := enrollment.GetQuality(); qualityEnrollment != nil {
				enrolledDevice.Enrollments = append(enrolledDevice.Enrollments, UpdateManagementEnrollment{
					UpdateCategory: types.StringValue("quality"),
				})
			}
		}

		if deviceErrors := azureADDevice.GetErrors(); deviceErrors != nil {
			for _, deviceError := range deviceErrors {
				if _, ok := deviceError.(graphmodels.AzureADDeviceRegistrationErrorable); ok {
					enrolledDevice.Errors = append(enrolledDevice.Errors, UpdatableAssetError{
						ErrorCode:    types.StringValue("AzureADDeviceRegistrationError"),
						ErrorMessage: types.StringValue("Device registration error occurred"),
					})
				}
			}
		}
	}

	return enrolledDevice
}
