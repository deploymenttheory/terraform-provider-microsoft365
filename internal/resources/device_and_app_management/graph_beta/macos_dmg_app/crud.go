package graphBetaMacOSDmgApp

import (
	"context"
	"fmt"
	"time"

	construct "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/constructors/graph_beta/device_and_app_management"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/crud"
	helpers "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/crud/graph_beta/device_and_app_management"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/errors"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// Create handles the complete creation workflow for a macOS DMG app resource in Intune.
func (r *MacOSDmgAppResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var object MacOSDmgAppResourceModel

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

	deadline, _ := ctx.Deadline()
	retryTimeout := time.Until(deadline) - time.Second

	// Step 1: Determine installer source path (local or download via URL)
	installerSourcePath, tempFileInfo, err := helpers.SetInstallerSourcePath(ctx, object.AppInstaller)
	if err != nil {
		resp.Diagnostics.AddError("Error determining installer file path", err.Error())
		return
	}

	// Ensure cleanup of temporary file occurs post state read
	if tempFileInfo.ShouldCleanup {
		defer helpers.CleanupTempFile(ctx, tempFileInfo)
	}

	// Step 3: Construct the base resource from the Terraform model
	requestBody, err := constructResource(ctx, &object, installerSourcePath)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing resource for Create method",
			fmt.Sprintf("Could not construct resource: %s_%s: %s", r.ProviderTypeName, r.TypeName, err.Error()),
		)
		return
	}

	// Step 4: Create the mobile app base resource in Graph API
	baseResource, err := r.client.
		DeviceAppManagement().
		MobileApps().
		Post(ctx, requestBody, nil)

	if err != nil {
		errors.HandleGraphError(ctx, err, resp, "Create", r.WritePermissions)
		return
	}

	object.ID = types.StringValue(*baseResource.GetId())
	tflog.Debug(ctx, fmt.Sprintf("Base resource created with ID: %s", object.ID.ValueString()))

	// Step 5: Associate categories with the app if provided
	if !object.Categories.IsNull() {
		var categoryValues []string
		diags := object.Categories.ElementsAs(ctx, &categoryValues, false)
		if diags.HasError() {
			resp.Diagnostics.Append(diags...)
			return
		}

		err = construct.AssignMobileAppCategories(ctx, r.client, object.ID.ValueString(), categoryValues, r.ReadPermissions)
		if err != nil {
			errors.HandleGraphError(ctx, err, resp, "Create", r.WritePermissions)
			return
		}
	}

	// If a DMG installer file is provided, process the content version and file upload
	if installerSourcePath != "" {

		// Step 6: Initialize content version
		tflog.Debug(ctx, "Initializing content version for file upload")

		content := graphmodels.NewMobileAppContent()

		contentBuilder := r.client.
			DeviceAppManagement().
			MobileApps().
			ByMobileAppId(object.ID.ValueString()).
			GraphMacOSDmgApp().
			ContentVersions()

		contentVersion, err := contentBuilder.Post(ctx, content, nil)
		if err != nil {
			errors.HandleGraphError(ctx, err, resp, "Create", r.WritePermissions)
			return
		}
		tflog.Debug(ctx, fmt.Sprintf("Content version created with ID: %s", *contentVersion.GetId()))

		// Step 7: Encrypt mobile app file and output file and encryption metadata
		tflog.Debug(ctx, "Encrypting installer file and constructing file metadata")
		contentFile, encryptionInfo, err := construct.EncryptMobileAppAndConstructFileContentMetadata(ctx, installerSourcePath)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error constructing and encrypting Intune Mobile App content file",
				err.Error(),
			)
			return
		}

		// Step 8: Create the content file resource in Graph API
		tflog.Debug(ctx, "Creating content file resource in Graph API")

		createdFile, err := contentBuilder.
			ByMobileAppContentId(*contentVersion.GetId()).
			Files().
			Post(ctx, contentFile, nil)

		if err != nil {
			errors.HandleGraphError(ctx, err, resp, "Create", r.WritePermissions)
			return
		}
		tflog.Debug(ctx, fmt.Sprintf("Content file resource created with ID: %s", *createdFile.GetId()))

		// Step 9: Wait for Graph API to generate a valid Azure Storage URI
		tflog.Debug(ctx, "Waiting for Graph API to generate a valid Azure Storage URI")
		err = retry.RetryContext(ctx, time.Until(deadline), func() *retry.RetryError {
			file, err := contentBuilder.
				ByMobileAppContentId(*contentVersion.GetId()).
				Files().
				ByMobileAppContentFileId(*createdFile.GetId()).
				Get(ctx, nil)

			if err != nil {
				tflog.Debug(ctx, fmt.Sprintf("Failed to get file status: %v", err))
				return retry.RetryableError(fmt.Errorf("failed to get file status: %v", err))
			}

			if file.GetUploadState() == nil {
				tflog.Debug(ctx, "Upload state is nil; retrying until a state is returned")
				return retry.RetryableError(fmt.Errorf("upload state is nil"))
			}

			state := *file.GetUploadState()
			tflog.Debug(ctx, fmt.Sprintf("Current upload state: %s", state.String()))

			if state == graphmodels.AZURESTORAGEURIREQUESTSUCCESS_MOBILEAPPCONTENTFILEUPLOADSTATE {
				tflog.Debug(ctx, "Azure Storage URI request successful")
				return nil
			}

			if state == graphmodels.AZURESTORAGEURIREQUESTFAILED_MOBILEAPPCONTENTFILEUPLOADSTATE {
				tflog.Debug(ctx, "Azure Storage URI request failed")
				return retry.NonRetryableError(fmt.Errorf("azure storage URI request failed"))
			}

			return retry.RetryableError(fmt.Errorf("upload state %s is not ready", state.String()))
		})

		if err != nil {
			resp.Diagnostics.AddError(
				"Error waiting for Azure Storage URI",
				fmt.Sprintf("Failed to get Azure Storage URI: %s", err.Error()),
			)
			return
		}

		// Step 10: Get the file with Azure Storage URI
		file, err := contentBuilder.
			ByMobileAppContentId(*contentVersion.GetId()).
			Files().
			ByMobileAppContentFileId(*createdFile.GetId()).
			Get(ctx, nil)

		if err != nil {
			errors.HandleGraphError(ctx, err, resp, "Create", r.WritePermissions)
			return
		}

		// Step 11: Upload encrypted file to Azure Storage
		tflog.Debug(ctx, "Uploading encrypted file to Azure Storage")
		encryptedFilePath := installerSourcePath + ".bin"
		defer helpers.CleanupTempFile(ctx, helpers.TempFileInfo{
			FilePath:      encryptedFilePath,
			ShouldCleanup: true,
		})

		err = construct.UploadToAzureStorage(ctx, *file.GetAzureStorageUri(), encryptedFilePath)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error uploading file to Azure Storage",
				fmt.Sprintf("Failed to upload encrypted file: %s", err.Error()),
			)
			return
		}

		// Step 12: Commit the file
		tflog.Debug(ctx, "Committing file")
		err = retry.RetryContext(ctx, time.Until(deadline), func() *retry.RetryError {
			commitBody, err := construct.CommitUploadedMobileAppWithEncryptionMetadata(encryptionInfo)
			if err != nil {
				return retry.NonRetryableError(fmt.Errorf("failed to construct commit request: %v", err))
			}

			err = contentBuilder.
				ByMobileAppContentId(*contentVersion.GetId()).
				Files().
				ByMobileAppContentFileId(*createdFile.GetId()).
				Commit().
				Post(ctx, commitBody, nil)

			if err != nil {
				tflog.Debug(ctx, fmt.Sprintf("Failed to commit file: %v", err))
				return retry.RetryableError(fmt.Errorf("failed to commit file: %v", err))
			}
			tflog.Debug(ctx, "File commit request successful")
			return nil
		})

		if err != nil {
			resp.Diagnostics.AddError(
				"Error committing file",
				err.Error(),
			)
			return
		}

		// Step 13: Wait for commit completion
		tflog.Debug(ctx, "Waiting for commit completion")
		err = waitForCommitCompletion(ctx, r.client, object.ID.ValueString(), *contentVersion.GetId(), retryTimeout)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error waiting for commit completion",
				fmt.Sprintf("Failed waiting for commit completion: %s", err.Error()),
			)
			return
		}

		// Step 14: Update the mobile app with the committed content version
		tflog.Debug(ctx, "Updating mobile app with committed content version")
		updatePayload := graphmodels.NewMacOSDmgApp()
		updatePayload.SetCommittedContentVersion(contentVersion.GetId())

		_, err = r.client.
			DeviceAppManagement().
			MobileApps().
			ByMobileAppId(object.ID.ValueString()).
			Patch(ctx, updatePayload, nil)

		if err != nil {
			errors.HandleGraphError(ctx, err, resp, "Create", r.WritePermissions)
			return
		}

		tflog.Debug(ctx, "App updated successfully with committed content version")
	}

	// Step 15: Read the resource to get the final state
	readReq := resource.ReadRequest{
		State: resp.State,
	}
	readResp := resource.ReadResponse{
		State: resp.State,
	}

	r.Read(ctx, readReq, &readResp)
	resp.State = readResp.State
	resp.Diagnostics.Append(readResp.Diagnostics...)
}

// Read retrieves the current state of a macOS DMG app resource from Intune.
func (r *MacOSDmgAppResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var object MacOSDmgAppResourceModel
	tflog.Debug(ctx, fmt.Sprintf("Starting to read resource: %s", ResourceName))

	resp.Diagnostics.Append(req.State.Get(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Read, ReadTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	appId := object.ID.ValueString()
	tflog.Debug(ctx, fmt.Sprintf("Reading DMG app with ID: %s", appId))

	mobileApp, err := r.client.
		DeviceAppManagement().
		MobileApps().
		ByMobileAppId(appId).
		GraphMacOSDmgApp().
		Get(ctx, nil)

	if err != nil {
		errors.HandleGraphError(ctx, err, resp, "Read", r.ReadPermissions)
		return
	}

	MapRemoteStateToTerraform(ctx, &object, mobileApp)

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	tflog.Debug(ctx, fmt.Sprintf("Finished reading resource: %s with ID: %s", ResourceName, object.ID.ValueString()))
}

// Update handles updates to a macOS DMG app resource in Intune.
func (r *MacOSDmgAppResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var object MacOSDmgAppResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting to update resource: %s", ResourceName))

	resp.Diagnostics.Append(req.Plan.Get(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Update, UpdateTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	appId := object.ID.ValueString()
	tflog.Debug(ctx, fmt.Sprintf("Updating DMG app with ID: %s", appId))

	// Construct the update request
	requestBody, err := constructResource(ctx, &object, "")
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing resource for Update method",
			fmt.Sprintf("Could not construct resource: %s_%s: %s", r.ProviderTypeName, r.TypeName, err.Error()),
		)
		return
	}

	// Update the resource
	_, err = r.client.
		DeviceAppManagement().
		MobileApps().
		ByMobileAppId(appId).
		Patch(ctx, requestBody, nil)

	if err != nil {
		errors.HandleGraphError(ctx, err, resp, "Update", r.WritePermissions)
		return
	}

	// Handle category updates
	if !object.Categories.IsNull() {
		var categoryValues []string
		diags := object.Categories.ElementsAs(ctx, &categoryValues, false)
		if diags.HasError() {
			resp.Diagnostics.Append(diags...)
			return
		}

		err = construct.AssignMobileAppCategories(ctx, r.client, appId, categoryValues, r.ReadPermissions)
		if err != nil {
			errors.HandleGraphError(ctx, err, resp, "Update", r.WritePermissions)
			return
		}
	}

	// Read the updated resource
	readReq := resource.ReadRequest{
		State: resp.State,
	}
	readResp := resource.ReadResponse{
		State: resp.State,
	}

	r.Read(ctx, readReq, &readResp)
	resp.State = readResp.State
	resp.Diagnostics.Append(readResp.Diagnostics...)
}

// Delete removes a macOS DMG app resource from Intune.
func (r *MacOSDmgAppResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var object MacOSDmgAppResourceModel
	tflog.Debug(ctx, fmt.Sprintf("Starting to delete resource: %s", ResourceName))

	resp.Diagnostics.Append(req.State.Get(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Delete, DeleteTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	appId := object.ID.ValueString()
	tflog.Debug(ctx, fmt.Sprintf("Deleting DMG app with ID: %s", appId))

	err := r.client.
		DeviceAppManagement().
		MobileApps().
		ByMobileAppId(appId).
		Delete(ctx, nil)

	if err != nil {
		errors.HandleGraphError(ctx, err, resp, "Delete", r.WritePermissions)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Successfully deleted resource: %s with ID: %s", ResourceName, appId))

	resp.State.RemoveResource(ctx)
}
