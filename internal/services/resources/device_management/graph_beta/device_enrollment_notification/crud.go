package graphBetaDeviceEnrollmentNotification

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	customrequests "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/custom_requests"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

/*
Create Function API Call Summary

	Step 1: Create Base Enrollment Notification Configuration

	API Call: POST /deviceManagement/deviceEnrollmentConfigurations
	- Purpose: Creates the main device enrollment notification configuration resource
	- Request Body: Contains display name, description, platform type, default locale, role scope tag IDs, and transformed notification templates (e.g., "email" → "email_00000000-0000-0000-0000-000000000000")
	- Response: Returns the created configuration with an ID and notification template IDs (format: "Email_56475a42-fa91-4245-99ee-a19c2eabdb6c", "Push_17f215be-e998-4020-9498-de6f2fe0e15c")

	Step 2: Extract Notification Templates

	- No API Call: Extracts notification template IDs from the creation response
	- Processing: Logs and categorizes templates as email or push based on their ID prefixes

	Step 3: Create Localized Notification Messages (if specified)

	API Call: POST /deviceManagement/notificationMessageTemplates/{templateGuid}/localizedNotificationMessages
	- Purpose: Creates localized messages for each notification template type
	- Frequency: Called for each template ID × each localized message combination
	- Request Processing:
	  - Extracts GUID from template ID (e.g., "Email_56475a42-fa91-4245-99ee-a19c2eabdb6c" → "56475a42-fa91-4245-99ee-a19c2eabdb6c")
	  - Matches template type (email/push) with message template type
	  - 2-second delay between each message creation to avoid rate limiting
	- Request Body: Contains locale, subject, message template, is_default flag

	Step 3.5: Update Branding Options (if specified)

	API Call: PATCH /deviceManagement/notificationMessageTemplates/{templateGuid}
	- Purpose: Updates branding options for each notification template
	- Frequency: Called once for each notification template created
	- Request Processing:
	  - Extracts GUID from each template ID
	  - Converts set of branding options to comma-separated bitmask format
	  - 1-second delay before each branding update
	- Request Body: Contains branding options enum (e.g., "includeCompanyLogo,includeCompanyName")

	Step 4: Assign Configuration to Groups (if specified)

	API Call: POST /deviceManagement/deviceEnrollmentConfigurations/{configId}/assign
	- Purpose: Assigns the enrollment configuration to specified groups or all licensed users
	- Frequency: Called once with all assignments
	- Request Body: Contains array of enrollment configuration assignments with target types:
	  - #microsoft.graph.groupAssignmentTarget (with groupId)
	  - #microsoft.graph.allLicensedUsersAssignmentTarget (no groupId needed)

	Step 5: Final State Read

	API Call: GET /deviceManagement/deviceEnrollmentConfigurations/{configId} (via ReadWithRetry)
	- Purpose: Retrieves the final state of the created resource for Terraform state
	- Additional Call: GET /deviceManagement/deviceEnrollmentConfigurations/{configId}/assignments
	- Purpose: Retrieves assignments separately (as they don't support expand)

	API Call Flow Summary

	1. 1 POST - Create base configuration
	2. N × M POST - Create localized messages (N templates × M messages)
	3. N PATCH - Update branding options (N templates)
	4. 1 POST - Assign to groups/users
	5. 2 GET - Read final state (base + assignments)

	Total API Calls

	For a typical configuration with:
	- 2 notification templates (email + push)
	- 2 localized messages per template
	- Branding options enabled
	- 1 assignment

	Total: 9 API calls
	- 1 create + 4 localized messages + 2 branding updates + 1 assignment + 1 final read + 1 assignments read
*/
func (r *DeviceEnrollmentNotificationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var object DeviceEnrollmentNotificationResourceModel

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

	// Step 1: Create the base enrollment notification configuration
	requestBody, err := constructResource(ctx, &object, true)
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

	if baseResource.GetId() == nil {
		resp.Diagnostics.AddError(
			"Error creating resource",
			fmt.Sprintf("Could not create %s: ID was nil", ResourceName),
		)
		return
	}

	object.ID = types.StringValue(*baseResource.GetId())
	tflog.Debug(ctx, fmt.Sprintf("Created base resource with ID: %s", *baseResource.GetId()))

	// Step 2: Extract notification templates from the response
	var notificationTemplates []string
	if enrollmentNotificationConfig, ok := baseResource.(graphmodels.DeviceEnrollmentNotificationConfigurationable); ok {
		notificationTemplates = enrollmentNotificationConfig.GetNotificationTemplates()
		tflog.Debug(ctx, fmt.Sprintf("Retrieved notification templates: %v", notificationTemplates))

		// Log detailed information about each template ID
		for i, templateID := range notificationTemplates {
			tflog.Debug(ctx, fmt.Sprintf("Template[%d]: ID=%s", i, templateID))

			// Check if it's an email or push template
			if strings.Contains(strings.ToLower(templateID), "email") {
				tflog.Debug(ctx, fmt.Sprintf("Identified as EMAIL template: %s", templateID))
			} else if strings.Contains(strings.ToLower(templateID), "push") {
				tflog.Debug(ctx, fmt.Sprintf("Identified as PUSH template: %s", templateID))
			} else {
				tflog.Debug(ctx, fmt.Sprintf("Unknown template type for ID: %s", templateID))
			}
		}
	} else {
		resp.Diagnostics.AddError(
			"Error extracting notification templates",
			fmt.Sprintf("Could not extract notification templates from response for %s", ResourceName),
		)
		return
	}

	// Step 3: Create localized notification messages for each template type if specified
	if !object.LocalizedNotificationMessages.IsNull() && !object.LocalizedNotificationMessages.IsUnknown() {
		// Add a small delay to ensure templates are fully created before adding messages
		tflog.Debug(ctx, "Waiting 2 seconds before creating localized notification messages to ensure templates are ready")
		time.Sleep(2 * time.Second)

		localizedMessages := make([]LocalizedNotificationMessageModel, 0, len(object.LocalizedNotificationMessages.Elements()))
		diags := object.LocalizedNotificationMessages.ElementsAs(ctx, &localizedMessages, false)
		if diags.HasError() {
			resp.Diagnostics.Append(diags...)
			return
		}

		// Create localized messages for each template ID
		for _, templateID := range notificationTemplates {
			for _, msg := range localizedMessages {
				// Determine if this is an email or push template
				templateType := ""
				if strings.Contains(strings.ToLower(templateID), "email") {
					templateType = "email"
					tflog.Debug(ctx, fmt.Sprintf("Processing EMAIL template: %s", templateID))
				} else if strings.Contains(strings.ToLower(templateID), "push") {
					templateType = "push"
					tflog.Debug(ctx, fmt.Sprintf("Processing PUSH template: %s", templateID))
				} else {
					tflog.Debug(ctx, fmt.Sprintf("Skipping unknown template type for ID: %s", templateID))
					continue
				}

				// Only create messages for the matching template type
				if msg.TemplateType.ValueString() != templateType {
					tflog.Debug(ctx, fmt.Sprintf("Template type mismatch - message type: %s, template type: %s - skipping",
						msg.TemplateType.ValueString(), templateType))
					continue
				}

				tflog.Debug(ctx, fmt.Sprintf("Found matching template type: %s for message type: %s",
					templateType, msg.TemplateType.ValueString()))

				// Construct the localized notification message using the proper constructor
				messageRequestBody, err := constructLocalizedNotificationMessage(ctx, msg)
				if err != nil {
					resp.Diagnostics.AddError(
						"Error constructing localized notification message",
						fmt.Sprintf("Could not construct localized notification message: %s", err.Error()),
					)
					return
				}

				// Log the request payload for debugging
				tflog.Debug(ctx, fmt.Sprintf("Creating localized notification message for template %s with locale %s",
					templateID, msg.Locale.ValueString()))

				// Extract just the GUID part from the template ID
				// The format is "Email_GUID" or "Push_GUID", we need just the GUID
				parts := strings.SplitN(templateID, "_", 2)
				var guidPart string
				if len(parts) == 2 {
					guidPart = parts[1]
					tflog.Debug(ctx, fmt.Sprintf("Extracted GUID part from template ID: %s -> %s", templateID, guidPart))
				} else {
					// If we can't split it, use the original (shouldn't happen)
					guidPart = templateID
					tflog.Debug(ctx, fmt.Sprintf("Could not extract GUID part, using original: %s", templateID))
				}

				// Add a small delay between each message creation to avoid rate limiting
				time.Sleep(2 * time.Second)

				tflog.Debug(ctx, fmt.Sprintf("Using GUID part %s for template ID in API call", guidPart))

				result, err := r.client.
					DeviceManagement().
					NotificationMessageTemplates().
					ByNotificationMessageTemplateId(guidPart).
					LocalizedNotificationMessages().
					Post(ctx, messageRequestBody, nil)

				if err != nil {
					// Log more details about the error
					tflog.Error(ctx, fmt.Sprintf("Error creating localized message for template %s: %s", templateID, err.Error()))
					errors.HandleGraphError(ctx, err, resp, "Create localized message", r.WritePermissions)
					return
				}

				if result != nil && result.GetId() != nil {
					tflog.Debug(ctx, fmt.Sprintf("Successfully created localized notification message with ID: %s for template %s",
						*result.GetId(), templateID))
				} else {
					tflog.Debug(ctx, fmt.Sprintf("Created localized notification message for template %s (no ID returned)", templateID))
				}
			}
		}
	}

	// Step 3.5: Update branding options for the notification templates
	// Construct branding options
	brandingTemplate, err := constructBrandingOptions(ctx, &object)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing branding options",
			fmt.Sprintf("Could not construct branding options for %s: %s", ResourceName, err.Error()),
		)
		return
	}

	if brandingTemplate != nil && len(notificationTemplates) > 0 {
		// Update branding options for each template
		for _, templateID := range notificationTemplates {
			// Extract just the GUID part from the template ID
			parts := strings.SplitN(templateID, "_", 2)
			var guidPart string
			if len(parts) == 2 {
				guidPart = parts[1]
				tflog.Debug(ctx, fmt.Sprintf("Extracted GUID part for branding options: %s -> %s", templateID, guidPart))
			} else {
				guidPart = templateID
				tflog.Debug(ctx, fmt.Sprintf("Using original template ID for branding: %s", templateID))
			}

			tflog.Debug(ctx, fmt.Sprintf("Updating branding options for template ID: %s", guidPart))

			// Add a small delay before updating branding options
			time.Sleep(1 * time.Second)

			_, err = r.client.
				DeviceManagement().
				NotificationMessageTemplates().
				ByNotificationMessageTemplateId(guidPart).
				Patch(ctx, brandingTemplate, nil)

			if err != nil {
				errors.HandleGraphError(ctx, err, resp, "Update branding options", r.WritePermissions)
				return
			}

			tflog.Debug(ctx, fmt.Sprintf("Successfully updated branding options for template %s", guidPart))
		}
	}

	// Step 4: Assign the configuration to groups if assignments are specified
	if !object.Assignments.IsNull() && !object.Assignments.IsUnknown() {
		assignmentRequestBody, err := constructAssignments(ctx, &object)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error constructing assignments",
				fmt.Sprintf("Could not construct assignments for %s: %s", ResourceName, err.Error()),
			)
			return
		}

		err = r.client.
			DeviceManagement().
			DeviceEnrollmentConfigurations().
			ByDeviceEnrollmentConfigurationId(object.ID.ValueString()).
			Assign().
			Post(ctx, assignmentRequestBody, nil)

		if err != nil {
			errors.HandleGraphError(ctx, err, resp, "Create assignments", r.WritePermissions)
			return
		}

		tflog.Debug(ctx, fmt.Sprintf("Successfully assigned configuration %s", object.ID.ValueString()))
	}

	// Set the initial state and call Read to get the latest state from API
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

// Read handles the Read operation for Android Enterprise Notification Configuration resources.
func (r *DeviceEnrollmentNotificationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var object DeviceEnrollmentNotificationResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Reading resource: %s with ID: %s", ResourceName, object.ID.ValueString()))

	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Read, ReadTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	// Step 1: Get the base resource
	resource, err := r.client.
		DeviceManagement().
		DeviceEnrollmentConfigurations().
		ByDeviceEnrollmentConfigurationId(object.ID.ValueString()).
		Get(ctx, nil)

	if err != nil {
		errors.HandleGraphError(ctx, err, resp, "Read", r.ReadPermissions)
		return
	}

	// Step 2: Get assignments separately (enrollment configurations don't support expand for assignments)
	assignments, err := r.client.
		DeviceManagement().
		DeviceEnrollmentConfigurations().
		ByDeviceEnrollmentConfigurationId(object.ID.ValueString()).
		Assignments().
		Get(ctx, nil)

	if err != nil {
		tflog.Warn(ctx, fmt.Sprintf("Could not fetch assignments for configuration %s: %s", object.ID.ValueString(), err.Error()))
		// Don't fail the read operation if assignments can't be fetched - set empty assignments
		assignments = nil
	}

	// Step 3: Map the base resource to Terraform state
	MapRemoteStateToTerraform(ctx, &object, resource)

	// Step 4: Map assignments to Terraform state
	if assignments != nil && assignments.GetValue() != nil {
		assignmentList := assignments.GetValue()
		tflog.Debug(ctx, fmt.Sprintf("Found %d assignments for configuration", len(assignmentList)))
		MapAssignmentsToTerraform(ctx, &object, assignmentList)
	} else {
		tflog.Debug(ctx, "No assignments found for configuration")
		object.Assignments = types.SetNull(AndroidEnterpriseNotificationsAssignmentType())
	}

	// Step 5: Get notification templates with localized messages
	if enrollmentNotificationConfig, ok := resource.(graphmodels.DeviceEnrollmentNotificationConfigurationable); ok {
		notificationTemplates := enrollmentNotificationConfig.GetNotificationTemplates()
		if len(notificationTemplates) > 0 {
			tflog.Debug(ctx, fmt.Sprintf("Found notification templates: %v", notificationTemplates))

			// Create a map to store templates and their types
			templates := []graphmodels.NotificationMessageTemplateable{}
			templateTypes := []string{}

			// Fetch each template with its localized messages
			for _, templateID := range notificationTemplates {
				// Determine template type from ID
				templateType := ""
				if strings.Contains(strings.ToLower(templateID), "email") {
					templateType = "email"
				} else if strings.Contains(strings.ToLower(templateID), "push") {
					templateType = "push"
				} else {
					tflog.Debug(ctx, fmt.Sprintf("Unknown template type for ID: %s", templateID))
					continue
				}

				// Extract just the GUID part from the template ID
				// The format is "Email_GUID" or "Push_GUID", we need just the GUID
				parts := strings.SplitN(templateID, "_", 2)
				var guidPart string
				if len(parts) == 2 {
					guidPart = parts[1]
					tflog.Debug(ctx, fmt.Sprintf("Extracted GUID part from template ID: %s -> %s", templateID, guidPart))
				} else {
					// If we can't split it, use the original (shouldn't happen)
					guidPart = templateID
					tflog.Debug(ctx, fmt.Sprintf("Could not extract GUID part, using original: %s", templateID))
				}

				// Get the template with localized messages using the GUID part
				// We need to make a separate call to get the localized messages since they don't auto-expand
				template, err := r.client.
					DeviceManagement().
					NotificationMessageTemplates().
					ByNotificationMessageTemplateId(guidPart).
					Get(ctx, nil)

				if err != nil {
					tflog.Warn(ctx, fmt.Sprintf("Error fetching template %s: %s", templateID, err.Error()))
					continue
				}

				// Separately fetch the localized messages for this template
				localizedMessages, err := r.client.
					DeviceManagement().
					NotificationMessageTemplates().
					ByNotificationMessageTemplateId(guidPart).
					LocalizedNotificationMessages().
					Get(ctx, nil)

				if err != nil {
					tflog.Warn(ctx, fmt.Sprintf("Error fetching localized messages for template %s: %s", templateID, err.Error()))
					// Continue with template but without localized messages
				} else if localizedMessages != nil && localizedMessages.GetValue() != nil {
					// Manually set the localized messages on the template
					template.SetLocalizedNotificationMessages(localizedMessages.GetValue())
					tflog.Debug(ctx, fmt.Sprintf("Fetched %d localized messages for template %s", len(localizedMessages.GetValue()), templateID))
				}

				templates = append(templates, template)
				templateTypes = append(templateTypes, templateType)
			}

			StateLocalizedNotificationMessages(ctx, &object, templates, templateTypes)
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	tflog.Debug(ctx, fmt.Sprintf("Finished reading resource: %s with ID: %s", ResourceName, object.ID.ValueString()))
}

// Update handles the Update operation for Android Enterprise Notification Configuration resources.
func (r *DeviceEnrollmentNotificationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan DeviceEnrollmentNotificationResourceModel
	var state DeviceEnrollmentNotificationResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Starting update of resource: %s with ID: %s", ResourceName, state.ID.ValueString()))

	ctx, cancel := crud.HandleTimeout(ctx, state.Timeouts.Update, UpdateTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	// Step 1: Get current resource to extract existing notification template GUIDs
	currentResource, err := r.client.
		DeviceManagement().
		DeviceEnrollmentConfigurations().
		ByDeviceEnrollmentConfigurationId(state.ID.ValueString()).
		Get(ctx, nil)

	if err != nil {
		errors.HandleGraphError(ctx, err, resp, "Read current resource for update", r.ReadPermissions)
		return
	}

	// Extract current notification templates to preserve actual GUIDs
	var currentNotificationTemplates []string
	if currentEnrollmentConfig, ok := currentResource.(graphmodels.DeviceEnrollmentNotificationConfigurationable); ok {
		currentNotificationTemplates = currentEnrollmentConfig.GetNotificationTemplates()
		tflog.Debug(ctx, fmt.Sprintf("Current notification templates: %v", currentNotificationTemplates))
	}

	// Step 2: Update the base enrollment notification configuration with actual template GUIDs
	requestBody, err := constructResource(ctx, &plan, false, currentNotificationTemplates)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing resource for update",
			fmt.Sprintf("Could not construct resource: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	baseUrl := fmt.Sprintf("https://graph.microsoft.com/beta/deviceManagement/deviceEnrollmentConfigurations/%s", state.ID.ValueString())
	tflog.Debug(ctx, fmt.Sprintf("Performing PATCH request to update base resource: %s", baseUrl))

	_, err = r.client.
		DeviceManagement().
		DeviceEnrollmentConfigurations().
		ByDeviceEnrollmentConfigurationId(state.ID.ValueString()).
		Patch(ctx, requestBody, nil)

	if err != nil {
		tflog.Error(ctx, fmt.Sprintf("Failed PATCH request to: %s - Error: %s", baseUrl, err.Error()))
		errors.HandleGraphError(ctx, err, resp, "Update", r.WritePermissions)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Completed PATCH request to update base resource: %s", baseUrl))

	// Step 3: Use the current notification templates (already retrieved in Step 1)
	notificationTemplates := currentNotificationTemplates
	tflog.Debug(ctx, fmt.Sprintf("Using current notification templates for localized message updates: %v", notificationTemplates))

	// Step 4: Update localized notification messages for each template type if specified
	if !plan.LocalizedNotificationMessages.IsNull() && !plan.LocalizedNotificationMessages.IsUnknown() {
		tflog.Debug(ctx, "Updating localized notification messages")

		localizedMessages := make([]LocalizedNotificationMessageModel, 0, len(plan.LocalizedNotificationMessages.Elements()))
		diags := plan.LocalizedNotificationMessages.ElementsAs(ctx, &localizedMessages, false)
		if diags.HasError() {
			resp.Diagnostics.Append(diags...)
			return
		}

		// Update/create localized messages for each template ID
		for _, templateID := range notificationTemplates {
			for _, msg := range localizedMessages {
				// Determine template type from ID
				templateType := ""
				if strings.Contains(strings.ToLower(templateID), "email") {
					templateType = "email"
				} else if strings.Contains(strings.ToLower(templateID), "push") {
					templateType = "push"
				} else {
					tflog.Debug(ctx, fmt.Sprintf("Skipping unknown template type for ID: %s", templateID))
					continue
				}

				// Only update messages for the matching template type
				if msg.TemplateType.ValueString() != templateType {
					continue
				}

				// Construct the localized notification message for update (exclude locale and isDefault)
				messageRequestBody, err := constructLocalizedNotificationMessage(ctx, msg, true)
				if err != nil {
					resp.Diagnostics.AddError(
						"Error constructing localized notification message",
						fmt.Sprintf("Could not construct localized notification message: %s", err.Error()),
					)
					return
				}

				// Extract GUID part from the template ID
				parts := strings.SplitN(templateID, "_", 2)
				var guidPart string
				if len(parts) == 2 {
					guidPart = parts[1]
				} else {
					guidPart = templateID
				}

				// Update the existing localized message using custom PATCH request with specific message ID
				// Message ID format: {templateGuid}_{locale}
				messageId := guidPart + "_" + msg.Locale.ValueString()

				time.Sleep(2 * time.Second) // Delay to avoid rate limiting

				// Use custom request to properly construct the URL path
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
					errors.HandleGraphError(ctx, err, resp, "Update localized message", r.WritePermissions)
					return
				}

				tflog.Debug(ctx, fmt.Sprintf("Completed custom PATCH request to update localized message: %s", localizedMessageUrl))
			}
		}
	}

	// Step 5: Update branding options for the notification templates
	brandingTemplate, err := constructBrandingOptions(ctx, &plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing branding options",
			fmt.Sprintf("Could not construct branding options for %s: %s", ResourceName, err.Error()),
		)
		return
	}

	if brandingTemplate != nil && len(notificationTemplates) > 0 {
		// Update branding options for each template
		for _, templateID := range notificationTemplates {
			// Extract GUID part from the template ID
			parts := strings.SplitN(templateID, "_", 2)
			var guidPart string
			if len(parts) == 2 {
				guidPart = parts[1]
			} else {
				guidPart = templateID
			}

			brandingUrl := fmt.Sprintf("https://graph.microsoft.com/beta/deviceManagement/notificationMessageTemplates/%s", guidPart)
			tflog.Debug(ctx, fmt.Sprintf("Performing PATCH request to update branding options: %s", brandingUrl))

			time.Sleep(1 * time.Second) // Delay to avoid rate limiting

			_, err = r.client.
				DeviceManagement().
				NotificationMessageTemplates().
				ByNotificationMessageTemplateId(guidPart).
				Patch(ctx, brandingTemplate, nil)

			if err != nil {
				tflog.Error(ctx, fmt.Sprintf("Failed PATCH request to: %s - Error: %s", brandingUrl, err.Error()))
				errors.HandleGraphError(ctx, err, resp, "Update branding options", r.WritePermissions)
				return
			}

			tflog.Debug(ctx, fmt.Sprintf("Completed PATCH request to update branding options: %s", brandingUrl))
		}
	}

	// Step 6: Update assignments if specified
	if !plan.Assignments.IsNull() && !plan.Assignments.IsUnknown() {
		assignmentRequestBody, err := constructAssignments(ctx, &plan)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error constructing assignments",
				fmt.Sprintf("Could not construct assignments for %s: %s", ResourceName, err.Error()),
			)
			return
		}

		assignmentUrl := fmt.Sprintf("https://graph.microsoft.com/beta/deviceManagement/deviceEnrollmentConfigurations/%s/assign", state.ID.ValueString())
		tflog.Debug(ctx, fmt.Sprintf("Performing POST request to update assignments: %s", assignmentUrl))

		err = r.client.
			DeviceManagement().
			DeviceEnrollmentConfigurations().
			ByDeviceEnrollmentConfigurationId(state.ID.ValueString()).
			Assign().
			Post(ctx, assignmentRequestBody, nil)

		if err != nil {
			tflog.Error(ctx, fmt.Sprintf("Failed POST request to: %s - Error: %s", assignmentUrl, err.Error()))
			errors.HandleGraphError(ctx, err, resp, "Update assignments", r.WritePermissions)
			return
		}

		tflog.Debug(ctx, fmt.Sprintf("Completed POST request to update assignments: %s", assignmentUrl))
	}

	// Set the plan state and call Read to get the latest state from API
	plan.ID = state.ID // Preserve the ID
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
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

	tflog.Debug(ctx, fmt.Sprintf("Finished updating %s with ID: %s", ResourceName, state.ID.ValueString()))
}

// Delete handles the Delete operation for Android Enterprise Notification Configuration resources.
func (r *DeviceEnrollmentNotificationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var object DeviceEnrollmentNotificationResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Starting deletion of resource: %s with ID: %s", ResourceName, object.ID.ValueString()))

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

	tflog.Debug(ctx, fmt.Sprintf("Finished deleting resource: %s with ID: %s", ResourceName, object.ID.ValueString()))
}
