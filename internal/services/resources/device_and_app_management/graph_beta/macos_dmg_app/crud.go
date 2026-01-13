package graphBetaMacOSDmgApp

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

// Create handles the complete creation workflow for a macOS DMG app resource in Intune.
//
// The function performs the following steps:
//
// 1. Reads the planned resource state from Terraform.
// 2. Constructs and creates the base resource from the Terraform model.
//
// If a package installer file is provided, the workflow continues as follows:
//
//  3. Initializes a new content version.
//  4. Encrypts the installer file locally (producing a .bin file) and constructs the file metadata,
//     including file size, encrypted file size, and encryption metadata (keys, digest, IV, MAC, etc.).
//  5. Creates a content file resource (with the metadata) in Graph API under the new content version.
//  6. Waits (via a retry loop using GET) for the Graph API to generate a valid Azure Storage SAS URI for the content file.
//  7. Retrieves the SAS URI and uploads the encrypted file (.bin) directly to Azure Blob Storage in chunks.
//  8. Commits the file by sending a commit request (including the encryption metadata) to Graph API,
//     and waits until the commit is confirmed.
//  9. Updates the mobile app resource (via a PATCH call) to set its committedContentVersion, so that
//     Intune uses the newly committed content file.
//  10. Assigns any mobile app assignments to the mobile app
//
// Finally, if app assignments are provided, the function creates the assignments, and then
// performs a final read of the resource state to verify successful creation.
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

		//
		err = construct.AssignMobileAppCategories(ctx, r.client, object.ID.ValueString(), categoryValues, r.ReadPermissions)
		//

		if err != nil {
			errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationCreate, r.WritePermissions)
			return
		}
	}

	// If a package installer file is provided, process the content version and file upload
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
		updatePayload := graphmodels.NewMacOSDmgApp()
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

// Read retrieves the current state of a macOS DMG app resource from Intune.
func (r *MacOSDmgAppResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var object MacOSDmgAppResourceModel

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
	macOSDmgApp, ok := respBaseResource.(graphmodels.MacOSDmgAppable)
	if !ok {
		resp.Diagnostics.AddError(
			"Resource type mismatch",
			fmt.Sprintf("Expected resource of type MacOSDmgAppable but got %T", respBaseResource),
		)
		return
	}

	MapRemoteResourceStateToTerraform(ctx, &object, macOSDmgApp)

	// 2. get committed app content file versions
	committedVersionIdPtr := macOSDmgApp.GetCommittedContentVersion()

	if committedVersionIdPtr != nil && *committedVersionIdPtr != "" {
		committedVersionId := *committedVersionIdPtr

		respFiles, err := r.client.DeviceAppManagement().
			MobileApps().
			ByMobileAppId(object.ID.ValueString()).
			GraphMacOSDmgApp().
			ContentVersions().
			ByMobileAppContentId(committedVersionId).
			Files().
			Get(ctx, nil)

		if err != nil {
			errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationRead, r.ReadPermissions)
			return
		}

		var installerFileName string
		if macOSDmgApp.GetFileName() != nil {
			installerFileName = *macOSDmgApp.GetFileName()
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

// Update handles updates to a macOS DMG app resource in Intune.
func (r *MacOSDmgAppResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state MacOSDmgAppResourceModel
	var installerSourcePath string
	var tempFileInfo helpers.TempFileInfo
	var err error

	tflog.Debug(ctx, fmt.Sprintf("Starting to update resource: %s", ResourceName))

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

	// No update logic is defined here intentionally for the reupload of mobile apps as any changes to either
	// installer_file_path_source orinstaller_url_source triggers a destory and redeploy. Implementation attempts were
	// becoming grossly complex and overly difficult for not much benefit.
	// Rationale
	// Intune's default behaviour is to append new mobileAppContentFiles to a new mobileContainedApp versions. Incrementing by 1
	// and references only the latest version.
	// The Api docs suggest you can delete existing mobileContainedApp versions but you cannot if they are commited. Meaning
	// there are disparities between the api docs and what was possible with the go sdk.
	// https://learn.microsoft.com/en-us/graph/api/intune-apps-mobileappcontent-delete?view=graph-rest-1.0 <- doesnt work ?
	//
	// This causes terraform stating issues as we are not interesting in tracking previous versions, we only want the latest
	// mobileContainedApp version and it's mobileAppContentFiles within.
	// with no delete option this would cause stating issues.
	//
	// Since in real world scenarios' updating the mobile app content files would mean a new app deployment in real terms,
	// as the only time you'd do this is if the app has installation bugs or unintended Ux. I believe force replacing
	// the mobile app acheices the same outcome from a sys admin perspective.
	//
	// This approach ensures that the logic used to auto generate the detection logic for included_apps. Is always run and keeps
	// the code more concise.
	//

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
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationDelete, r.WritePermissions)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Successfully deleted resource: %s with ID: %s", ResourceName, appId))

	resp.State.RemoveResource(ctx)
}
