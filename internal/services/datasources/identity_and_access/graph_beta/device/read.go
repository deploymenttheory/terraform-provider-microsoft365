// REF: https://learn.microsoft.com/en-us/graph/api/device-list?view=graph-rest-beta
// REF: https://learn.microsoft.com/en-us/graph/api/device-get?view=graph-rest-beta
// REF: https://learn.microsoft.com/en-us/graph/api/device-list-memberof?view=graph-rest-beta
// REF: https://learn.microsoft.com/en-us/graph/api/device-list-registeredowners?view=graph-rest-beta
// REF: https://learn.microsoft.com/en-us/graph/api/device-list-registeredusers?view=graph-rest-beta
package graphBetaDevice

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
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
	graphdevices "github.com/microsoftgraph/msgraph-beta-sdk-go/devices"
	graphcore "github.com/microsoftgraph/msgraph-sdk-go-core"
)

type lookupMethod int

const (
	lookupByObjectId lookupMethod = iota
	lookupByDisplayName
	lookupByDeviceId
	lookupByODataQuery
	lookupListAll
)

// Read handles the Read operation.
func (d *DeviceDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var object DeviceDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Read, ReadTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	method := d.determineLookupMethod(object)
	var devices []graphmodels.Deviceable
	var err error

	switch method {
	case lookupByObjectId:
		devices, err = d.getDeviceByObjectId(ctx, object)
	case lookupByDisplayName:
		devices, err = d.getDevicesByDisplayName(ctx, object)
	case lookupByDeviceId:
		devices, err = d.getDevicesByDeviceId(ctx, object)
	case lookupByODataQuery:
		devices, err = d.getDevicesByODataQuery(ctx, object)
	case lookupListAll:
		devices, err = d.listAllDevices(ctx)
	}

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationRead, d.ReadPermissions)
		return
	}

	// Handle additional lookups for memberOf, registeredOwners, registeredUsers
	// These require an object ID - use the provided one or extract from first device
	var objectIdForRelationships string
	if !object.ObjectId.IsNull() && !object.ObjectId.IsUnknown() {
		objectIdForRelationships = object.ObjectId.ValueString()
	} else if len(devices) > 0 && devices[0].GetId() != nil {
		// Extract object ID from first device for relationship lookups
		objectIdForRelationships = *devices[0].GetId()
	}

	if objectIdForRelationships != "" {
		if !object.ListMemberOf.IsNull() && object.ListMemberOf.ValueBool() {
			memberOf, err := d.getDeviceMemberOf(ctx, objectIdForRelationships)
			if err != nil {
				errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationRead, d.ReadPermissions)
				return
			}
			object.MemberOf = ConstructDirectoryObjectItems(memberOf)
		}

		if !object.ListRegisteredOwners.IsNull() && object.ListRegisteredOwners.ValueBool() {
			registeredOwners, err := d.getDeviceRegisteredOwners(ctx, objectIdForRelationships)
			if err != nil {
				errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationRead, d.ReadPermissions)
				return
			}
			object.RegisteredOwners = ConstructDirectoryObjectItems(registeredOwners)
		}

		if !object.ListRegisteredUsers.IsNull() && object.ListRegisteredUsers.ValueBool() {
			registeredUsers, err := d.getDeviceRegisteredUsers(ctx, objectIdForRelationships)
			if err != nil {
				errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationRead, d.ReadPermissions)
				return
			}
			object.RegisteredUsers = ConstructDirectoryObjectItems(registeredUsers)
		}
	}

	object.Items = ConstructDeviceItems(devices)
	object.ID = types.StringValue(fmt.Sprintf("device-datasource-%d", time.Now().Unix()))

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
}

// determineLookupMethod determines which lookup method to use based on the provided attributes
func (d *DeviceDataSource) determineLookupMethod(object DeviceDataSourceModel) lookupMethod {
	switch {
	case !object.ObjectId.IsNull() && object.ObjectId.ValueString() != "":
		return lookupByObjectId
	case !object.DisplayName.IsNull() && object.DisplayName.ValueString() != "":
		return lookupByDisplayName
	case !object.DeviceId.IsNull() && object.DeviceId.ValueString() != "":
		return lookupByDeviceId
	case !object.ODataQuery.IsNull() && object.ODataQuery.ValueString() != "":
		return lookupByODataQuery
	case !object.ListAll.IsNull() && object.ListAll.ValueBool():
		return lookupListAll
	default:
		return lookupListAll
	}
}

// getDeviceByObjectId retrieves a device by its object ID
func (d *DeviceDataSource) getDeviceByObjectId(ctx context.Context, object DeviceDataSourceModel) ([]graphmodels.Deviceable, error) {
	objectId := object.ObjectId.ValueString()

	tflog.Debug(ctx, fmt.Sprintf("Looking up device by object ID: %s", objectId))

	device, err := d.client.Devices().ByDeviceId(objectId).Get(ctx, nil)
	if err != nil {
		return nil, err
	}

	if device == nil {
		return []graphmodels.Deviceable{}, nil
	}

	return []graphmodels.Deviceable{device}, nil
}

// getDevicesByDisplayName retrieves devices filtered by display name
func (d *DeviceDataSource) getDevicesByDisplayName(ctx context.Context, object DeviceDataSourceModel) ([]graphmodels.Deviceable, error) {
	displayName := object.DisplayName.ValueString()

	tflog.Debug(ctx, fmt.Sprintf("Looking up devices by display name: %s", displayName))

	filter := fmt.Sprintf("displayName eq '%s'", displayName)

	headers := abstractions.NewRequestHeaders()
	headers.Add("ConsistencyLevel", "eventual")

	requestConfig := &graphdevices.DevicesRequestBuilderGetRequestConfiguration{
		Headers: headers,
		QueryParameters: &graphdevices.DevicesRequestBuilderGetQueryParameters{
			Filter: &filter,
		},
	}

	return d.listAllDevicesWithPageIterator(ctx, requestConfig)
}

// getDevicesByDeviceId retrieves devices filtered by device ID
func (d *DeviceDataSource) getDevicesByDeviceId(ctx context.Context, object DeviceDataSourceModel) ([]graphmodels.Deviceable, error) {
	deviceId := object.DeviceId.ValueString()

	tflog.Debug(ctx, fmt.Sprintf("Looking up devices by device ID: %s", deviceId))

	filter := fmt.Sprintf("deviceId eq '%s'", deviceId)

	headers := abstractions.NewRequestHeaders()
	headers.Add("ConsistencyLevel", "eventual")

	requestConfig := &graphdevices.DevicesRequestBuilderGetRequestConfiguration{
		Headers: headers,
		QueryParameters: &graphdevices.DevicesRequestBuilderGetQueryParameters{
			Filter: &filter,
		},
	}

	return d.listAllDevicesWithPageIterator(ctx, requestConfig)
}

// getDevicesByODataQuery retrieves devices using a custom OData query
func (d *DeviceDataSource) getDevicesByODataQuery(ctx context.Context, object DeviceDataSourceModel) ([]graphmodels.Deviceable, error) {
	filter := object.ODataQuery.ValueString()

	tflog.Debug(ctx, fmt.Sprintf("Looking up devices with OData query: %s", filter))

	headers := abstractions.NewRequestHeaders()
	headers.Add("ConsistencyLevel", "eventual")

	requestConfig := &graphdevices.DevicesRequestBuilderGetRequestConfiguration{
		Headers: headers,
		QueryParameters: &graphdevices.DevicesRequestBuilderGetQueryParameters{
			Filter: &filter,
		},
	}

	return d.listAllDevicesWithPageIterator(ctx, requestConfig)
}

// listAllDevices retrieves all devices in the tenant
func (d *DeviceDataSource) listAllDevices(ctx context.Context) ([]graphmodels.Deviceable, error) {
	tflog.Debug(ctx, "Listing all devices")

	return d.listAllDevicesWithPageIterator(ctx, nil)
}

// listAllDevicesWithPageIterator handles pagination for device list requests
func (d *DeviceDataSource) listAllDevicesWithPageIterator(ctx context.Context, requestConfig *graphdevices.DevicesRequestBuilderGetRequestConfiguration) ([]graphmodels.Deviceable, error) {
	var allDevices []graphmodels.Deviceable

	result, err := d.client.Devices().Get(ctx, requestConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to get initial page of devices: %w", err)
	}

	if result == nil {
		return allDevices, nil
	}

	pageIterator, err := graphcore.NewPageIterator[graphmodels.Deviceable](
		result,
		d.client.GetAdapter(),
		graphmodels.CreateDeviceCollectionResponseFromDiscriminatorValue,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create page iterator: %w", err)
	}

	err = pageIterator.Iterate(ctx, func(device graphmodels.Deviceable) bool {
		if device != nil {
			allDevices = append(allDevices, device)
		}
		return true
	})

	if err != nil {
		return nil, fmt.Errorf("error during pagination: %w", err)
	}

	return allDevices, nil
}

// getDeviceMemberOf retrieves the groups and administrative units that the device is a member of
func (d *DeviceDataSource) getDeviceMemberOf(ctx context.Context, objectId string) ([]graphmodels.DirectoryObjectable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Retrieving memberOf for device: %s", objectId))

	var allMemberOf []graphmodels.DirectoryObjectable

	result, err := d.client.Devices().ByDeviceId(objectId).MemberOf().Get(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get memberOf for device: %w", err)
	}

	if result == nil {
		return allMemberOf, nil
	}

	pageIterator, err := graphcore.NewPageIterator[graphmodels.DirectoryObjectable](
		result,
		d.client.GetAdapter(),
		graphmodels.CreateDirectoryObjectCollectionResponseFromDiscriminatorValue,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create memberOf page iterator: %w", err)
	}

	err = pageIterator.Iterate(ctx, func(item graphmodels.DirectoryObjectable) bool {
		if item != nil {
			allMemberOf = append(allMemberOf, item)
		}
		return true
	})

	if err != nil {
		return nil, fmt.Errorf("error during memberOf pagination: %w", err)
	}

	return allMemberOf, nil
}

// getDeviceRegisteredOwners retrieves the registered owners of the device
func (d *DeviceDataSource) getDeviceRegisteredOwners(ctx context.Context, objectId string) ([]graphmodels.DirectoryObjectable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Retrieving registered owners for device: %s", objectId))

	var allOwners []graphmodels.DirectoryObjectable

	result, err := d.client.Devices().ByDeviceId(objectId).RegisteredOwners().Get(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get registered owners for device: %w", err)
	}

	if result == nil {
		return allOwners, nil
	}

	pageIterator, err := graphcore.NewPageIterator[graphmodels.DirectoryObjectable](
		result,
		d.client.GetAdapter(),
		graphmodels.CreateDirectoryObjectCollectionResponseFromDiscriminatorValue,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create registered owners page iterator: %w", err)
	}

	err = pageIterator.Iterate(ctx, func(item graphmodels.DirectoryObjectable) bool {
		if item != nil {
			allOwners = append(allOwners, item)
		}
		return true
	})

	if err != nil {
		return nil, fmt.Errorf("error during registered owners pagination: %w", err)
	}

	return allOwners, nil
}

// getDeviceRegisteredUsers retrieves the registered users of the device
func (d *DeviceDataSource) getDeviceRegisteredUsers(ctx context.Context, objectId string) ([]graphmodels.DirectoryObjectable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Retrieving registered users for device: %s", objectId))

	var allUsers []graphmodels.DirectoryObjectable

	result, err := d.client.Devices().ByDeviceId(objectId).RegisteredUsers().Get(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get registered users for device: %w", err)
	}

	if result == nil {
		return allUsers, nil
	}

	pageIterator, err := graphcore.NewPageIterator[graphmodels.DirectoryObjectable](
		result,
		d.client.GetAdapter(),
		graphmodels.CreateDirectoryObjectCollectionResponseFromDiscriminatorValue,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create registered users page iterator: %w", err)
	}

	err = pageIterator.Iterate(ctx, func(item graphmodels.DirectoryObjectable) bool {
		if item != nil {
			allUsers = append(allUsers, item)
		}
		return true
	})

	if err != nil {
		return nil, fmt.Errorf("error during registered users pagination: %w", err)
	}

	return allUsers, nil
}
