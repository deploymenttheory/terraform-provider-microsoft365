package graphBetaDeviceEnrollmentNotificationConfiguration

import (
	"context"
	"fmt"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/crud"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/errors"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Create handles the Create operation for Device Enrollment Notification Configuration resources.
func (r *DeviceEnrollmentNotificationConfigurationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var object DeviceEnrollmentNotificationConfigurationResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting creation of resource: %s", ResourceName))

	resp.Diagnostics.Append(req.Plan.Get(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Create, CreateTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	requestBody, err := constructResource(ctx, &object)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing resource",
			fmt.Sprintf("Could not construct resource: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	baseResource, err := r.client.
		DeviceManagement().
		DeviceEnrollmentConfigurations().
		Post(ctx, requestBody, nil)

	if err != nil {
		errors.HandleGraphError(ctx, err, resp, "Create", r.WritePermissions)
		return
	}

	object.ID = types.StringValue(*baseResource.GetId())

	// Create localized messages if specified
	if object.PushLocalizedMessage != nil {
		templateGUID, err := r.resolveNotificationTemplateID(ctx, object.ID.ValueString(), "push")
		if err != nil {
			resp.Diagnostics.AddError(
				"Error resolving push template ID",
				fmt.Sprintf("Could not resolve push template ID: %s", err.Error()),
			)
			return
		}

		requestBody := constructLocalizedMessage(ctx, object.PushLocalizedMessage)
		if requestBody != nil {
			_, err = r.client.
				DeviceManagement().
				NotificationMessageTemplates().
				ByNotificationMessageTemplateId(templateGUID).
				LocalizedNotificationMessages().
				Post(ctx, requestBody, nil)

			if err != nil {
				errors.HandleGraphError(ctx, err, resp, "Create", r.WritePermissions)
				return
			}

			tflog.Debug(ctx, "Successfully created push localized message")
		}
	}

	if object.EmailLocalizedMessage != nil {
		templateGUID, err := r.resolveNotificationTemplateID(ctx, object.ID.ValueString(), "email")
		if err != nil {
			resp.Diagnostics.AddError(
				"Error resolving email template ID",
				fmt.Sprintf("Could not resolve email template ID: %s", err.Error()),
			)
			return
		}

		requestBody := constructLocalizedMessage(ctx, object.EmailLocalizedMessage)
		if requestBody != nil {
			_, err = r.client.
				DeviceManagement().
				NotificationMessageTemplates().
				ByNotificationMessageTemplateId(templateGUID).
				LocalizedNotificationMessages().
				Post(ctx, requestBody, nil)

			if err != nil {
				errors.HandleGraphError(ctx, err, resp, "Create", r.WritePermissions)
				return
			}

			tflog.Debug(ctx, "Successfully created email localized message")
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	readReq := resource.ReadRequest{State: resp.State, ProviderMeta: req.ProviderMeta}
	stateContainer := &crud.CreateResponseContainer{CreateResponse: resp}

	opts := crud.DefaultReadWithRetryOptions()
	opts.Operation = "Create"
	opts.ResourceTypeName = constants.PROVIDER_NAME + "_" + ResourceName

	err = crud.ReadWithRetry(ctx, r.Read, readReq, stateContainer, opts)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading resource state after create",
			fmt.Sprintf("Could not read resource state: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Create Method: %s", ResourceName))
}

// Read handles the Read operation for Device Enrollment Notification Configuration resources.
func (r *DeviceEnrollmentNotificationConfigurationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var object DeviceEnrollmentNotificationConfigurationResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting Read method for: %s", ResourceName))

	resp.Diagnostics.Append(req.State.Get(ctx, &object)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Reading %s with ID: %s", ResourceName, object.ID.ValueString()))

	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Read, ReadTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	resource, err := r.client.
		DeviceManagement().
		DeviceEnrollmentConfigurations().
		ByDeviceEnrollmentConfigurationId(object.ID.ValueString()).
		Get(ctx, nil)

	if err != nil {
		errors.HandleGraphError(ctx, err, resp, "Read", r.ReadPermissions)
		return
	}

	mapRemoteStateToTerraform(ctx, &object, resource)

	// Read localized messages if template types are configured
	if !object.TemplateTypes.IsNull() && !object.TemplateTypes.IsUnknown() {
		templateTypes := make([]string, 0, len(object.TemplateTypes.Elements()))
		for _, element := range object.TemplateTypes.Elements() {
			if stringVal, ok := element.(types.String); ok && !stringVal.IsNull() {
				templateTypes = append(templateTypes, stringVal.ValueString())
			}
		}

		for _, templateType := range templateTypes {
			templateGUID, err := r.resolveNotificationTemplateID(ctx, object.ID.ValueString(), templateType)
			if err != nil {
				tflog.Warn(ctx, fmt.Sprintf("Failed to resolve template ID for type %s: %s", templateType, err.Error()))
				continue
			}

			messages, err := r.client.
				DeviceManagement().
				NotificationMessageTemplates().
				ByNotificationMessageTemplateId(templateGUID).
				LocalizedNotificationMessages().
				Get(ctx, nil)

			if err != nil {
				tflog.Warn(ctx, fmt.Sprintf("Failed to get localized messages for template %s: %s", templateGUID, err.Error()))
				continue
			}

			if messagesCollection := messages.GetValue(); len(messagesCollection) > 0 {
				// Get the first (and typically only) localized message
				message := messagesCollection[0]
				localizedModel := &LocalizedNotificationMessageModel{
					Locale:          types.StringPointerValue(message.GetLocale()),
					Subject:         types.StringPointerValue(message.GetSubject()),
					MessageTemplate: types.StringPointerValue(message.GetMessageTemplate()),
					IsDefault:       types.BoolPointerValue(message.GetIsDefault()),
				}

				// Assign to the appropriate field based on template type
				switch templateType {
				case "push":
					object.PushLocalizedMessage = localizedModel
				case "email":
					object.EmailLocalizedMessage = localizedModel
				}
				tflog.Debug(ctx, fmt.Sprintf("Successfully read localized message for template type %s", templateType))
			}
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Read Method: %s", ResourceName))
}

// Update handles the Update operation.
func (r *DeviceEnrollmentNotificationConfigurationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var object DeviceEnrollmentNotificationConfigurationResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting Update of resource: %s", ResourceName))

	resp.Diagnostics.Append(req.Plan.Get(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Update, UpdateTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	requestBody, err := constructResource(ctx, &object)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing resource for update method",
			fmt.Sprintf("Could not construct resource: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	_, err = r.client.
		DeviceManagement().
		DeviceEnrollmentConfigurations().
		ByDeviceEnrollmentConfigurationId(object.ID.ValueString()).
		Patch(ctx, requestBody, nil)

	if err != nil {
		errors.HandleGraphError(ctx, err, resp, "Update", r.WritePermissions)
		return
	}

	readReq := resource.ReadRequest{State: resp.State, ProviderMeta: req.ProviderMeta}
	stateContainer := &crud.UpdateResponseContainer{UpdateResponse: resp}

	opts := crud.DefaultReadWithRetryOptions()
	opts.Operation = "Update"
	opts.ResourceTypeName = constants.PROVIDER_NAME + "_" + ResourceName

	err = crud.ReadWithRetry(ctx, r.Read, readReq, stateContainer, opts)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading resource state after update",
			fmt.Sprintf("Could not read resource state: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Update Method: %s", ResourceName))
}

// Delete handles the Delete operation for Device Enrollment Notification Configuration resources.
func (r *DeviceEnrollmentNotificationConfigurationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var object DeviceEnrollmentNotificationConfigurationResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting deletion of resource: %s", ResourceName))

	resp.Diagnostics.Append(req.State.Get(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Delete, DeleteTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	err := r.client.
		DeviceManagement().
		DeviceEnrollmentConfigurations().
		ByDeviceEnrollmentConfigurationId(object.ID.ValueString()).
		Delete(ctx, nil)

	if err != nil {
		errors.HandleGraphError(ctx, err, resp, "Delete", r.WritePermissions)
		return
	}

	resp.State.RemoveResource(ctx)

	tflog.Debug(ctx, fmt.Sprintf("Finished Delete Method: %s", ResourceName))
}
