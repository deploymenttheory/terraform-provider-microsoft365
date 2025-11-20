package graphBetaDeviceComplianceNotificationTemplate

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	customrequests "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/custom_requests"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/devicemanagement"
)

// Create handles the resource creation
func (r *DeviceComplianceNotificationTemplateResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan DeviceComplianceNotificationTemplateResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting creation of resource: %s", ResourceName))

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := crud.HandleTimeout(ctx, plan.Timeouts.Create, CreateTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	// Step 1: Create base template without localized messages
	baseTemplate, err := constructResource(ctx, &plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing base template for create method",
			fmt.Sprintf("Could not construct base template: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	createdTemplate, err := r.client.
		DeviceManagement().
		NotificationMessageTemplates().
		Post(ctx, baseTemplate, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, "Create", r.WritePermissions)
		return
	}

	// Extract the created template ID for subsequent localized message creation
	templateId := createdTemplate.GetId()
	if templateId == nil {
		resp.Diagnostics.AddError(
			"Error getting template ID after creation",
			"Created template did not return a valid ID",
		)
		return
	}

	plan.ID = types.StringValue(*templateId)

	// Step 2: Create localized notification messages if provided
	if !plan.LocalizedNotificationMessages.IsNull() && !plan.LocalizedNotificationMessages.IsUnknown() {
		var localizedMessages []LocalizedNotificationMessageModel
		plan.LocalizedNotificationMessages.ElementsAs(ctx, &localizedMessages, false)

		// Sort messages so the one with isDefault=true is processed first
		// This ensures the defaultLocale is set correctly by the API
		var defaultMessage *LocalizedNotificationMessageModel
		var otherMessages []LocalizedNotificationMessageModel
		var defaultIndex int = -1

		for i, msg := range localizedMessages {
			if !msg.IsDefault.IsNull() && !msg.IsDefault.IsUnknown() && msg.IsDefault.ValueBool() {
				defaultMessage = &localizedMessages[i]
				defaultIndex = i
			} else {
				otherMessages = append(otherMessages, localizedMessages[i])
			}
		}

		// Process default message first if it exists
		if defaultMessage != nil {
			tflog.Debug(ctx, fmt.Sprintf("Processing default localized message first: %s", defaultMessage.Locale.ValueString()))
			localizedMessage, err := constructLocalizedMessage(ctx, defaultMessage, false)
			if err != nil {
				resp.Diagnostics.AddError(
					"Error constructing default localized message",
					fmt.Sprintf("Could not construct default localized message for locale %s: %s", defaultMessage.Locale.ValueString(), err.Error()),
				)
				return
			}

			createdLocalizedMessage, err := r.client.
				DeviceManagement().
				NotificationMessageTemplates().
				ByNotificationMessageTemplateId(*templateId).
				LocalizedNotificationMessages().
				Post(ctx, localizedMessage, nil)

			if err != nil {
				resp.Diagnostics.AddError(
					"Error creating default localized message",
					fmt.Sprintf("Could not create default localized message for locale %s: %s", defaultMessage.Locale.ValueString(), err.Error()),
				)
				return
			}

			// Update the plan with the created message ID
			if createdLocalizedMessage.GetId() != nil {
				localizedMessages[defaultIndex].ID = types.StringValue(*createdLocalizedMessage.GetId())
			}
		}

		// Process remaining messages
		for _, msg := range otherMessages {
			tflog.Debug(ctx, fmt.Sprintf("Processing additional localized message: %s", msg.Locale.ValueString()))
			localizedMessage, err := constructLocalizedMessage(ctx, &msg, false)
			if err != nil {
				resp.Diagnostics.AddError(
					"Error constructing localized message",
					fmt.Sprintf("Could not construct localized message for locale %s: %s", msg.Locale.ValueString(), err.Error()),
				)
				return
			}

			createdLocalizedMessage, err := r.client.
				DeviceManagement().
				NotificationMessageTemplates().
				ByNotificationMessageTemplateId(*templateId).
				LocalizedNotificationMessages().
				Post(ctx, localizedMessage, nil)

			if err != nil {
				resp.Diagnostics.AddError(
					"Error creating localized message",
					fmt.Sprintf("Could not create localized message for locale %s: %s", msg.Locale.ValueString(), err.Error()),
				)
				return
			}

			for j := range localizedMessages {
				if localizedMessages[j].Locale.ValueString() == msg.Locale.ValueString() && j != defaultIndex {
					if createdLocalizedMessage.GetId() != nil {
						localizedMessages[j].ID = types.StringValue(*createdLocalizedMessage.GetId())
					}
					break
				}
			}
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	readReq := resource.ReadRequest{State: resp.State, ProviderMeta: req.ProviderMeta}
	stateContainer := &crud.CreateResponseContainer{CreateResponse: resp}

	opts := crud.DefaultReadWithRetryOptions()
	opts.Operation = "Create"
	opts.ResourceTypeName = ResourceName

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

// Read handles the resource read
func (r *DeviceComplianceNotificationTemplateResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state DeviceComplianceNotificationTemplateResourceModel

	operation := "Read"
	if ctxOp := ctx.Value("retry_operation"); ctxOp != nil {
		if opStr, ok := ctxOp.(string); ok {
			operation = opStr
		}
	}

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Reading %s with ID: %s (operation: %s)", ResourceName, state.ID.ValueString(), operation))

	ctx, cancel := crud.HandleTimeout(ctx, state.Timeouts.Read, ReadTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "resource_id", state.ID.ValueString())
	tflog.Debug(ctx, "Starting read of notification message template")

	template, err := r.client.
		DeviceManagement().
		NotificationMessageTemplates().
		ByNotificationMessageTemplateId(state.ID.ValueString()).
		Get(ctx, &devicemanagement.NotificationMessageTemplatesNotificationMessageTemplateItemRequestBuilderGetRequestConfiguration{
			QueryParameters: &devicemanagement.NotificationMessageTemplatesNotificationMessageTemplateItemRequestBuilderGetQueryParameters{
				Expand: []string{"localizedNotificationMessages"},
			},
		})

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, operation, r.ReadPermissions)
		return
	}

	MapRemoteStateToTerraform(ctx, &state, template)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Read Method: %s", ResourceName))
}

// Update handles the resource update
func (r *DeviceComplianceNotificationTemplateResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan DeviceComplianceNotificationTemplateResourceModel
	var state DeviceComplianceNotificationTemplateResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting update method for: %s", ResourceName))

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := crud.HandleTimeout(ctx, plan.Timeouts.Update, UpdateTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	// Step 1: Update base template properties (without localized messages)
	baseTemplate, err := constructResource(ctx, &plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing base template for update method",
			fmt.Sprintf("Could not construct base template: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	_, err = r.client.
		DeviceManagement().
		NotificationMessageTemplates().
		ByNotificationMessageTemplateId(plan.ID.ValueString()).
		Patch(ctx, baseTemplate, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, "Update", r.WritePermissions)
		return
	}

	// Step 2: Handle localized notification messages with full CRUD operations
	templateId := plan.ID.ValueString()

	// Get current state localized messages to compare
	var currentMessages []LocalizedNotificationMessageModel
	if !state.LocalizedNotificationMessages.IsNull() && !state.LocalizedNotificationMessages.IsUnknown() {
		state.LocalizedNotificationMessages.ElementsAs(ctx, &currentMessages, false)
	}

	// Create maps for easy lookup
	currentMessagesByLocale := make(map[string]LocalizedNotificationMessageModel)
	for _, msg := range currentMessages {
		currentMessagesByLocale[msg.Locale.ValueString()] = msg
	}

	var planMessages []LocalizedNotificationMessageModel
	planMessagesByLocale := make(map[string]LocalizedNotificationMessageModel)
	if !plan.LocalizedNotificationMessages.IsNull() && !plan.LocalizedNotificationMessages.IsUnknown() {
		plan.LocalizedNotificationMessages.ElementsAs(ctx, &planMessages, false)
		for _, msg := range planMessages {
			planMessagesByLocale[msg.Locale.ValueString()] = msg
		}
	}

	// Step 2a: DELETE messages that exist in current state but not in plan
	for locale, currentMsg := range currentMessagesByLocale {
		if _, existsInPlan := planMessagesByLocale[locale]; !existsInPlan {
			tflog.Debug(ctx, fmt.Sprintf("Deleting localized message for locale: %s", locale))

			err := r.client.
				DeviceManagement().
				NotificationMessageTemplates().
				ByNotificationMessageTemplateId(templateId).
				LocalizedNotificationMessages().
				ByLocalizedNotificationMessageId(currentMsg.ID.ValueString()).
				Delete(ctx, nil)

			if err != nil {
				tflog.Warn(ctx, "Failed to delete localized message", map[string]any{
					"locale":    locale,
					"messageId": currentMsg.ID.ValueString(),
					"error":     err.Error(),
				})
				// Continue with other operations even if delete fails
			}
		}
	}

	// Step 2b: Process planned messages (CREATE new, UPDATE existing)
	var finalMessages []LocalizedNotificationMessageModel

	for _, msg := range planMessages {
		locale := msg.Locale.ValueString()

		if currentMsg, existsInCurrent := currentMessagesByLocale[locale]; existsInCurrent {
			// Message exists - PATCH it
			tflog.Debug(ctx, fmt.Sprintf("Updating existing localized message for locale: %s", locale))

			messageRequestBody, err := constructLocalizedMessage(ctx, &msg, true)
			if err != nil {
				resp.Diagnostics.AddError(
					"Error constructing localized notification message for update",
					fmt.Sprintf("Could not construct localized notification message: %s", err.Error()),
				)
				return
			}

			// Extract GUID part from the template ID
			parts := strings.SplitN(templateId, "_", 2)
			var guidPart string
			if len(parts) == 2 {
				guidPart = parts[1]
			} else {
				guidPart = templateId
			}

			messageId := guidPart + "_" + locale
			time.Sleep(2 * time.Second) // Delay to avoid rate limiting

			config := customrequests.PatchRequestConfig{
				APIVersion:        customrequests.GraphAPIBeta,
				Endpoint:          fmt.Sprintf("deviceManagement/notificationMessageTemplates/%s/localizedNotificationMessages", guidPart),
				ResourceID:        messageId,
				ResourceIDPattern: "/{id}",
				RequestBody:       messageRequestBody,
			}

			localizedMessageUrl := fmt.Sprintf("https://graph.microsoft.com/beta/deviceManagement/notificationMessageTemplates/%s/localizedNotificationMessages/%s", guidPart, messageId)
			tflog.Debug(ctx, fmt.Sprintf("Performing custom PATCH request to update localized message: %s", localizedMessageUrl))

			err = customrequests.PatchRequestByResourceId(ctx, r.client.GetAdapter(), config)
			if err != nil {
				tflog.Error(ctx, fmt.Sprintf("Failed custom PATCH request to: %s - Error: %s", localizedMessageUrl, err.Error()))
				errors.HandleKiotaGraphError(ctx, err, resp, "Update localized message", r.WritePermissions)
				return
			}

			// Keep the existing message ID and add to final messages
			updatedMsg := msg
			updatedMsg.ID = currentMsg.ID
			finalMessages = append(finalMessages, updatedMsg)

		} else {
			// Message is new - POST it
			tflog.Debug(ctx, fmt.Sprintf("Creating new localized message for locale: %s", locale))

			messageRequestBody, err := constructLocalizedMessage(ctx, &msg, false)
			if err != nil {
				resp.Diagnostics.AddError(
					"Error constructing new localized notification message",
					fmt.Sprintf("Could not construct new localized notification message: %s", err.Error()),
				)
				return
			}

			createdLocalizedMessage, err := r.client.
				DeviceManagement().
				NotificationMessageTemplates().
				ByNotificationMessageTemplateId(templateId).
				LocalizedNotificationMessages().
				Post(ctx, messageRequestBody, nil)

			if err != nil {
				resp.Diagnostics.AddError(
					"Error creating new localized message",
					fmt.Sprintf("Could not create new localized message for locale %s: %s", locale, err.Error()),
				)
				return
			}

			// Create final message with the created ID
			newMsg := msg
			if createdLocalizedMessage.GetId() != nil {
				newMsg.ID = types.StringValue(*createdLocalizedMessage.GetId())
			}
			finalMessages = append(finalMessages, newMsg)
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	readReq := resource.ReadRequest{State: resp.State, ProviderMeta: req.ProviderMeta}
	stateContainer := &crud.UpdateResponseContainer{UpdateResponse: resp}

	opts := crud.DefaultReadWithRetryOptions()
	opts.Operation = "Update"
	opts.ResourceTypeName = ResourceName

	err = crud.ReadWithRetry(ctx, r.Read, readReq, stateContainer, opts)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading resource state after update",
			fmt.Sprintf("Could not read resource state: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished updating %s with ID: %s", ResourceName, state.ID.ValueString()))
}

// Delete handles the resource deletion
func (r *DeviceComplianceNotificationTemplateResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state DeviceComplianceNotificationTemplateResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting deletion of resource: %s", ResourceName))

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := crud.HandleTimeout(ctx, state.Timeouts.Delete, DeleteTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	err := r.client.
		DeviceManagement().
		NotificationMessageTemplates().
		ByNotificationMessageTemplateId(state.ID.ValueString()).
		Delete(ctx, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, "Delete", r.WritePermissions)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Removing %s from Terraform state", ResourceName))

	resp.State.RemoveResource(ctx)

	tflog.Debug(ctx, fmt.Sprintf("Finished Delete Method: %s", ResourceName))
}
