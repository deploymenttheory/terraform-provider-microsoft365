package graphBetaMacOSLobApp

import (
	"context"
	"fmt"
	"path/filepath"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	construct "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors/graph_beta/device_and_app_management"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	helpers "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud/graph_beta/device_and_app_management"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/shared_models/graph_beta/device_and_app_management"
	sharedstater "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/state/graph_beta/device_and_app_management"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/deviceappmanagement"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// Create handles the Create operation for macOS LOB App resources.
//
// Operation: Creates a new macOS line-of-business application with content file upload workflow
// API Calls:
//   - POST /deviceAppManagement/mobileApps
//   - POST /deviceAppManagement/mobileApps/{mobileAppId}/microsoft.graph.macOSLobApp/contentVersions
//   - POST /deviceAppManagement/mobileApps/{mobileAppId}/microsoft.graph.macOSLobApp/contentVersions/{contentVersionId}/files
//   - POST /deviceAppManagement/mobileApps/{mobileAppId}/microsoft.graph.macOSLobApp/contentVersions/{contentVersionId}/files/{fileId}/commit
//   - PATCH /deviceAppManagement/mobileApps/{mobileAppId}
//
// Reference: https://learn.microsoft.com/en-us/graph/api/intune-apps-macoslobapp-create?view=graph-rest-beta
// Note: Includes encrypted file upload to Azure Storage and content version management
func (r *MacOSLobAppResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var object MacOSLobAppResourceModel

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
			fmt.Sprintf("Could not construct resource: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	// Step 4: Create the mobile app base resource in Graph API
	baseResource, err := r.client.
		DeviceAppManagement().
		MobileApps().
		Post(ctx, requestBody, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationCreate, r.WritePermissions)
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
			errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationCreate, r.WritePermissions)
			return
		}
	}

	// If a LOB app installer file is provided, process the content version and file upload
	if installerSourcePath != "" {

		// Step 6: Initialize content version
		tflog.Debug(ctx, "Initializing content version for file upload")

		content := graphmodels.NewMobileAppContent()

		contentBuilder := r.client.
			DeviceAppManagement().
			MobileApps().
			ByMobileAppId(object.ID.ValueString()).
			GraphMacOSLobApp().
			ContentVersions()

		contentVersion, err := contentBuilder.Post(ctx, content, nil)
		if err != nil {
			errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationCreate, r.WritePermissions)
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
			errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationCreate, r.WritePermissions)
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

			tflog.Debug(ctx, fmt.Sprintf("Waiting for Azure Storage URI, current state: %s", state.String()))
			return retry.RetryableError(fmt.Errorf("waiting for Azure Storage URI, current state: %s", state.String()))
		})
		if err != nil {
			errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationCreate, r.WritePermissions)
			return
		}

		// Step 10: Retrieve the Azure Storage URI and upload the encrypted file
		tflog.Debug(ctx, "Retrieving file status for Azure Storage URI")

		fileStatus, err := contentBuilder.
			ByMobileAppContentId(*contentVersion.GetId()).
			Files().
			ByMobileAppContentFileId(*createdFile.GetId()).
			Get(ctx, nil)

		if err != nil {
			errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationCreate, r.WritePermissions)
			return
		}

		if fileStatus.GetAzureStorageUri() == nil {
			tflog.Debug(ctx, "Azure Storage URI is nil in the retrieved file status")
			errors.HandleKiotaGraphError(ctx, fmt.Errorf("azure Storage URI is nil"), resp, constants.TfOperationCreate, r.WritePermissions)
			return
		}

		tflog.Debug(ctx, fmt.Sprintf("Retrieved Azure Storage URI: %s", *fileStatus.GetAzureStorageUri()))

		// IMPORTANT: Upload the encrypted file (.bin) not the original source file.
		encryptedFilePath := installerSourcePath + ".bin"

		tflog.Debug(ctx, fmt.Sprintf("Uploading encrypted file: %s", encryptedFilePath))

		err = construct.UploadToAzureStorage(ctx, *fileStatus.GetAzureStorageUri(), encryptedFilePath)
		if err != nil {
			tflog.Debug(ctx, fmt.Sprintf("Failed to upload to Azure Storage: %v", err))
			errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationCreate, r.WritePermissions)
			return
		}

		// Step 11: Commit the file with encryption metadata
		tflog.Debug(ctx, "Committing file with encryption metadata")
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

		// Step 12: Wait for commit to complete
		tflog.Debug(ctx, "Waiting for file commit to complete")

		err = WaitForFileCommitCompletion(
			ctx,
			contentBuilder,
			*contentVersion.GetId(),
			*createdFile.GetId(),
			encryptionInfo,
			resp,
			r.WritePermissions,
		)

		if err != nil {
			resp.Diagnostics.AddError(
				"Error waiting for file commit",
				err.Error(),
			)
			return
		}

		// Step 13: Update the App with the Committed Content Version
		updatePayload := graphmodels.NewMacOSLobApp()
		updatePayload.SetCommittedContentVersion(contentVersion.GetId())

		_, err = r.client.
			DeviceAppManagement().
			MobileApps().
			ByMobileAppId(object.ID.ValueString()).
			Patch(ctx, updatePayload, nil)

		if err != nil {
			errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationCreate, r.WritePermissions)
			return
		}

		tflog.Debug(ctx, "App updated successfully with committed content version")
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	readReq := resource.ReadRequest{State: resp.State}
	stateContainer := &crud.CreateResponseContainer{CreateResponse: resp}

	opts := crud.DefaultReadWithRetryOptions()
	opts.Operation = constants.TfOperationCreate
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

// Read handles the Read operation for macOS LOB App resources.
//
// Operation: Retrieves a macOS line-of-business application by ID
// API Calls:
//   - GET /deviceAppManagement/mobileApps/{mobileAppId}
//
// Reference: https://learn.microsoft.com/en-us/graph/api/intune-apps-macoslobapp-get?view=graph-rest-beta
func (r *MacOSLobAppResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var object MacOSLobAppResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting Read method for: %s", ResourceName))

	operation := constants.TfOperationRead
	if ctxOp := ctx.Value("retry_operation"); ctxOp != nil {
		if opStr, ok := ctxOp.(string); ok {
			operation = opStr
		}
	}
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

	// 1. get base resource with expanded query to return categories
	requestParameters := &deviceappmanagement.MobileAppsMobileAppItemRequestBuilderGetRequestConfiguration{
		QueryParameters: &deviceappmanagement.MobileAppsMobileAppItemRequestBuilderGetQueryParameters{
			Expand: []string{"categories"},
		},
	}

	respBaseResource, err := r.client.
		DeviceAppManagement().
		MobileApps().
		ByMobileAppId(object.ID.ValueString()).
		Get(ctx, requestParameters)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, operation, r.ReadPermissions)
		return
	}

	// This ensures type safety as the Graph API returns a base interface that needs
	// to be converted to the specific app type
	macOSLobApp, ok := respBaseResource.(graphmodels.MacOSLobAppable)
	if !ok {
		resp.Diagnostics.AddError(
			"Resource type mismatch",
			fmt.Sprintf("Expected resource of type MacOSLobAppable but got %T", respBaseResource),
		)
		return
	}

	MapRemoteResourceStateToTerraform(ctx, &object, macOSLobApp)

	// 2. get committed app content file versions
	committedVersionIdPtr := macOSLobApp.GetCommittedContentVersion()

	if committedVersionIdPtr != nil && *committedVersionIdPtr != "" {
		committedVersionId := *committedVersionIdPtr

		respFiles, err := r.client.DeviceAppManagement().
			MobileApps().
			ByMobileAppId(object.ID.ValueString()).
			GraphMacOSLobApp().
			ContentVersions().
			ByMobileAppContentId(committedVersionId).
			Files().
			Get(ctx, nil)

		if err != nil {
			errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationRead, r.ReadPermissions)
			return
		}

		var installerFileName string
		if macOSLobApp.GetFileName() != nil {
			installerFileName = *macOSLobApp.GetFileName()
		} else {
			var metadataObj sharedmodels.MobileAppMetaDataResourceModel
			if !object.AppInstaller.IsNull() {
				diags := object.AppInstaller.As(ctx, &metadataObj, basetypes.ObjectAsOptions{})
				if !diags.HasError() && !metadataObj.InstallerFilePathSource.IsNull() {
					filePath := metadataObj.InstallerFilePathSource.ValueString()
					if filePath != "" {
						installerFileName = filepath.Base(filePath)
					}
				}
			}
		}

		object.ContentVersion = sharedstater.MapCommittedContentVersionStateToTerraform(ctx, committedVersionId, respFiles, err, installerFileName)
	}

	// 3. Get app metadata by processing app installer file
	var existingMetadata sharedmodels.MobileAppMetaDataResourceModel
	if !req.State.Raw.IsNull() {
		diags := req.State.GetAttribute(ctx, path.Root("app_installer"), &existingMetadata)
		if diags.HasError() {
			resp.Diagnostics.Append(diags...)
			return
		}
		object.AppInstaller = sharedstater.MapAppMetadataStateToTerraform(ctx, &existingMetadata)
	}

	// 6. set final state
	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Read Method: %s", ResourceName))
}

// Update handles the Update operation for macOS LOB App resources.
//
// Operation: Updates an existing macOS line-of-business application
// API Calls:
//   - PATCH /deviceAppManagement/mobileApps/{mobileAppId}
//
// Reference: https://learn.microsoft.com/en-us/graph/api/intune-apps-macoslobapp-update?view=graph-rest-beta
func (r *MacOSLobAppResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state MacOSLobAppResourceModel
	var installerSourcePath string
	var tempFileInfo helpers.TempFileInfo
	var err error

	tflog.Debug(ctx, fmt.Sprintf("Updating %s with ID: %s", ResourceName, state.ID.ValueString()))

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

	// Ensure cleanup of temporary file when we're done
	if tempFileInfo.ShouldCleanup {
		defer helpers.CleanupTempFile(ctx, tempFileInfo)
	}

	installerSourcePath, tempFileInfo, err = helpers.SetInstallerSourcePath(ctx, plan.AppInstaller)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error determining installer file path",
			err.Error(),
		)
		return
	}

	// Step 1: Update the base mobile app resource
	requestBody, err := constructResource(ctx, &plan, installerSourcePath)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing resource for update method",
			fmt.Sprintf("Could not construct resource: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	_, err = r.client.
		DeviceAppManagement().
		MobileApps().
		ByMobileAppId(plan.ID.ValueString()).
		Patch(ctx, requestBody, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationUpdate, r.WritePermissions)
		return
	}

	// Step 2: Updated Categories
	if !plan.Categories.Equal(state.Categories) {
		tflog.Debug(ctx, "Categories have changed — updating categories")

		var categoryValues []string
		diags := plan.Categories.ElementsAs(ctx, &categoryValues, false)
		if diags.HasError() {
			resp.Diagnostics.Append(diags...)
			return
		}

		err = construct.AssignMobileAppCategories(ctx, r.client, plan.ID.ValueString(), categoryValues, r.ReadPermissions)
		if err != nil {
			errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationUpdate, r.WritePermissions)
			return
		}
	}

	// Read updated resource state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	readReq := resource.ReadRequest{State: resp.State}
	stateContainer := &crud.UpdateResponseContainer{UpdateResponse: resp}

	opts := crud.DefaultReadWithRetryOptions()
	opts.Operation = constants.TfOperationUpdate
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

// Delete handles the Delete operation for macOS LOB App resources.
//
// Operation: Deletes a macOS line-of-business application
// API Calls:
//   - DELETE /deviceAppManagement/mobileApps/{mobileAppId}
//
// Reference: https://learn.microsoft.com/en-us/graph/api/intune-apps-macoslobapp-delete?view=graph-rest-beta
func (r *MacOSLobAppResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var object MacOSLobAppResourceModel

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
		DeviceAppManagement().
		MobileApps().
		ByMobileAppId(object.ID.ValueString()).
		Delete(ctx, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationDelete, r.WritePermissions)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Removing %s from Terraform state", ResourceName))

	resp.State.RemoveResource(ctx)

	tflog.Debug(ctx, fmt.Sprintf("Finished Delete Method: %s", ResourceName))
}
