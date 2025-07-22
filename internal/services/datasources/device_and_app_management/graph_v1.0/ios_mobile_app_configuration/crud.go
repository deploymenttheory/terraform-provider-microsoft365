package graphIOSMobileAppConfiguration

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-sdk-go/models"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	commonerrors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors"
)

// Define static errors
var (
	ErrInvalidResourceType = errors.New("resource is not an iOS mobile app configuration")
	ErrNilAPIResponse      = errors.New("received nil response from Microsoft Graph API")
	ErrResourceNotFound    = errors.New("iOS mobile app configuration not found")
	ErrMultipleResources   = errors.New("multiple iOS mobile app configurations found")
)

// Read retrieves information about an iOS mobile app configuration
func (d *IOSMobileAppConfigurationDataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	var data IOSMobileAppConfigurationDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Ensure either ID or DisplayName is provided
	if data.Id.IsNull() && data.DisplayName.IsNull() {
		resp.Diagnostics.AddError(
			"Invalid Configuration",
			"Either 'id' or 'display_name' must be provided to lookup the iOS mobile app configuration.",
		)
		return
	}

	// Set a default timeout
	ctx, cancel := crud.HandleTimeout(
		ctx,
		data.Timeouts.Read,
		time.Duration(ReadTimeout)*time.Second,
		&resp.Diagnostics,
	)
	if cancel == nil {
		return
	}
	defer cancel()

	// Perform the read operation
	if !data.Id.IsNull() {
		// Read by ID
		err := d.readByID(ctx, &data, resp)
		if err != nil {
			return
		}
	} else {
		// Read by display name
		err := d.readByDisplayName(ctx, &data, resp)
		if err != nil {
			return
		}
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// readByID retrieves the iOS mobile app configuration by ID
func (d *IOSMobileAppConfigurationDataSource) readByID(
	ctx context.Context,
	data *IOSMobileAppConfigurationDataSourceModel,
	resp *datasource.ReadResponse,
) error {
	tflog.Debug(ctx, "Reading iOS Mobile App Configuration by ID", map[string]interface{}{
		"id": data.Id.ValueString(),
	})

	// Get the mobile app configuration
	resource, err := d.client.DeviceAppManagement().
		MobileAppConfigurations().
		ByManagedDeviceMobileAppConfigurationId(data.Id.ValueString()).
		Get(ctx, nil)
	if err != nil {
		commonerrors.HandleGraphError(ctx, err, resp, "Read", d.ReadPermissions)
		return fmt.Errorf("failed to get mobile app configuration by ID: %w", err)
	}

	// Check if we got an iOS mobile app configuration
	iosConfig, ok := resource.(models.IosMobileAppConfigurationable)
	if !ok {
		resp.Diagnostics.AddError(
			"Invalid Resource Type",
			"Resource is not an iOS mobile app configuration",
		)
		return ErrInvalidResourceType
	}

	// Map the response to the data source model
	mapResourceToDataSourceState(ctx, iosConfig, data, &resp.Diagnostics)

	// Read assignments separately
	assignmentsResp, err := d.client.DeviceAppManagement().
		MobileAppConfigurations().
		ByManagedDeviceMobileAppConfigurationId(data.Id.ValueString()).
		Assignments().
		Get(ctx, nil)
	if err != nil {
		// Log warning but don't fail the read
		tflog.Warn(ctx, "Failed to read assignments", map[string]interface{}{"error": err.Error()})
	} else if assignmentsResp != nil && assignmentsResp.GetValue() != nil {
		data.Assignments = mapAssignmentsToDataSourceState(ctx, assignmentsResp.GetValue(), &resp.Diagnostics)
	}

	return nil
}

// readByDisplayName retrieves the iOS mobile app configuration by display name
func (d *IOSMobileAppConfigurationDataSource) readByDisplayName(
	ctx context.Context,
	data *IOSMobileAppConfigurationDataSourceModel,
	resp *datasource.ReadResponse,
) error {
	tflog.Debug(ctx, "Reading iOS Mobile App Configuration by Display Name", map[string]interface{}{
		"display_name": data.DisplayName.ValueString(),
	})

	// List all mobile app configurations (will filter manually)
	result, err := d.client.DeviceAppManagement().
		MobileAppConfigurations().
		Get(ctx, nil)
	if err != nil {
		commonerrors.HandleGraphError(ctx, err, resp, "Read", d.ReadPermissions)
		return fmt.Errorf("failed to list mobile app configurations: %w", err)
	}

	// Check if result is nil
	if result == nil {
		resp.Diagnostics.AddError(
			"API Error",
			"Received nil response from Microsoft Graph API",
		)
		return ErrNilAPIResponse
	}

	// Filter results manually since the API might not support OData filtering
	var matchingConfigs []models.ManagedDeviceMobileAppConfigurationable
	configurations := result.GetValue()

	for _, config := range configurations {
		if config == nil {
			continue
		}

		// Check if it's an iOS configuration
		iosConfig, ok := config.(models.IosMobileAppConfigurationable)
		if !ok {
			continue
		}

		// Check if display name matches
		if iosConfig.GetDisplayName() != nil &&
			*iosConfig.GetDisplayName() == data.DisplayName.ValueString() {
			matchingConfigs = append(matchingConfigs, iosConfig)
		}
	}

	// Check if we found any configurations
	if len(matchingConfigs) == 0 {
		resp.Diagnostics.AddError(
			"Resource Not Found",
			fmt.Sprintf(
				"No iOS mobile app configuration found with display name '%s'",
				data.DisplayName.ValueString(),
			),
		)
		return fmt.Errorf(
			"%w: display name '%s'",
			ErrResourceNotFound,
			data.DisplayName.ValueString(),
		)
	}

	// Check if we found multiple configurations
	if len(matchingConfigs) > 1 {
		resp.Diagnostics.AddError(
			"Multiple Resources Found",
			fmt.Sprintf(
				"Multiple iOS mobile app configurations found with display name '%s'",
				data.DisplayName.ValueString(),
			),
		)
		return fmt.Errorf(
			"%w: display name '%s'",
			ErrMultipleResources,
			data.DisplayName.ValueString(),
		)
	}

	// Get the first (and only) configuration
	iosConfig := matchingConfigs[0].(models.IosMobileAppConfigurationable)

	// Map the response to the data source model
	mapResourceToDataSourceState(ctx, iosConfig, data, &resp.Diagnostics)

	// Read assignments separately if we have an ID
	if iosConfig.GetId() != nil {
		assignmentsResp, err := d.client.DeviceAppManagement().
			MobileAppConfigurations().
			ByManagedDeviceMobileAppConfigurationId(*iosConfig.GetId()).
			Assignments().
			Get(ctx, nil)
		if err != nil {
			// Log warning but don't fail the read
			tflog.Warn(
				ctx,
				"Failed to read assignments",
				map[string]interface{}{"error": err.Error()},
			)
		} else if assignmentsResp != nil && assignmentsResp.GetValue() != nil {
			data.Assignments = mapAssignmentsToDataSourceState(ctx, assignmentsResp.GetValue(), &resp.Diagnostics)
		}
	}

	return nil
}
