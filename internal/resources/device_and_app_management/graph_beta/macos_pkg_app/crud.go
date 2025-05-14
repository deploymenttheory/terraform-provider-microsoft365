package graphBetaMacOSPKGApp

import (
	"context"
	"fmt"
	"path/filepath"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	construct "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/constructors/graph_beta/device_and_app_management"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/crud"
	helpers "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/crud/graph_beta/device_and_app_management"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/errors"
	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/shared_models/graph_beta/device_and_app_management"
	sharedstater "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/state/graph_beta/device_and_app_management"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/deviceappmanagement"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// Create handles the complete creation workflow for a macOS PKG app resource in Intune.
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
func (r *MacOSPKGAppResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var object MacOSPKGAppResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting creation of resource: %s_%s", r.ProviderTypeName, r.TypeName))

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

		//constants.GraphSDKMutex.Lock()
		err = construct.AssignMobileAppCategories(ctx, r.client, object.ID.ValueString(), categoryValues, r.ReadPermissions)
		//constants.GraphSDKMutex.Unlock()

		if err != nil {
			errors.HandleGraphError(ctx, err, resp, "Create", r.WritePermissions)
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
			GraphMacOSPkgApp().
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

		constants.GraphSDKMutex.Lock()
		createdFile, err := contentBuilder.
			ByMobileAppContentId(*contentVersion.GetId()).
			Files().
			Post(ctx, contentFile, nil)
		constants.GraphSDKMutex.Unlock()

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

			tflog.Debug(ctx, fmt.Sprintf("Waiting for Azure Storage URI, current state: %s", state.String()))
			return retry.RetryableError(fmt.Errorf("waiting for Azure Storage URI, current state: %s", state.String()))
		})
		if err != nil {
			errors.HandleGraphError(ctx, err, resp, "Create", r.WritePermissions)
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
			errors.HandleGraphError(ctx, err, resp, "Create", r.WritePermissions)
			return
		}

		if fileStatus.GetAzureStorageUri() == nil {
			tflog.Debug(ctx, "Azure Storage URI is nil in the retrieved file status")
			errors.HandleGraphError(ctx, fmt.Errorf("azure Storage URI is nil"), resp, "Create", r.WritePermissions)
			return
		}

		tflog.Debug(ctx, fmt.Sprintf("Retrieved Azure Storage URI: %s", *fileStatus.GetAzureStorageUri()))

		// IMPORTANT: Upload the encrypted file (.bin) not the original source file.
		encryptedFilePath := installerSourcePath + ".bin"

		tflog.Debug(ctx, fmt.Sprintf("Uploading encrypted file: %s", encryptedFilePath))

		err = construct.UploadToAzureStorage(ctx, *fileStatus.GetAzureStorageUri(), encryptedFilePath)
		if err != nil {
			tflog.Debug(ctx, fmt.Sprintf("Failed to upload to Azure Storage: %v", err))
			errors.HandleGraphError(ctx, err, resp, "Create", r.WritePermissions)
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

		// maxRetries := 10
		// for i := 0; i < maxRetries; i++ {

		// 	file, err := contentBuilder.
		// 		ByMobileAppContentId(*contentVersion.GetId()).
		// 		Files().
		// 		ByMobileAppContentFileId(*createdFile.GetId()).
		// 		Get(ctx, nil)

		// 	if err != nil {
		// 		errors.HandleGraphError(ctx, err, resp, "Create", r.WritePermissions)
		// 		return
		// 	}

		// 	state := *file.GetUploadState()

		// 	tflog.Debug(ctx, fmt.Sprintf("Commit status check %d: state=%s", i+1, state.String()))
		// 	if state == graphmodels.COMMITFILESUCCESS_MOBILEAPPCONTENTFILEUPLOADSTATE {
		// 		tflog.Debug(ctx, "File commit completed successfully")
		// 		break
		// 	}

		// 	if state == graphmodels.COMMITFILEFAILED_MOBILEAPPCONTENTFILEUPLOADSTATE {
		// 		tflog.Debug(ctx, "File commit failed; retrying commit request")
		// 		commitBody, err := construct.CommitUploadedMobileAppWithEncryptionMetadata(encryptionInfo)
		// 		if err != nil {
		// 			tflog.Debug(ctx, fmt.Sprintf("Error constructing commit request during retry: %v", err))
		// 			continue
		// 		}

		// 		err = contentBuilder.
		// 			ByMobileAppContentId(*contentVersion.GetId()).
		// 			Files().
		// 			ByMobileAppContentFileId(*createdFile.GetId()).
		// 			Commit().
		// 			Post(ctx, commitBody, nil)

		// 		if err != nil {
		// 			tflog.Debug(ctx, fmt.Sprintf("Error during commit retry: %v", err))
		// 			continue
		// 		}
		// 	}

		// 	if i == maxRetries-1 {
		// 		resp.Diagnostics.AddError(
		// 			"Error waiting for file commit",
		// 			fmt.Sprintf("File commit did not complete after %d attempts. Last state: %s", maxRetries, state.String()),
		// 		)
		// 		return
		// 	}
		// 	time.Sleep(10 * time.Second)
		// }

		// Step 13: Update the App with the Committed Content Version
		updatePayload := graphmodels.NewMacOSPkgApp()
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

	// Step 14: Apply Assignments
	if len(object.Assignments) > 0 {
		requestAssignment, err := construct.ConstructMobileAppAssignment(ctx, object.Assignments)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error constructing assignment for Create Method",
				fmt.Sprintf("Could not construct assignment: %s_%s: %s", r.ProviderTypeName, r.TypeName, err.Error()),
			)
			return
		}

		deadline, _ := ctx.Deadline()
		retryTimeout := time.Until(deadline) - time.Second

		err = retry.RetryContext(ctx, retryTimeout, func() *retry.RetryError {
			err := r.client.
				DeviceAppManagement().
				MobileApps().
				ByMobileAppId(object.ID.ValueString()).
				Assign().
				Post(ctx, requestAssignment, nil)

			if err != nil {
				return retry.RetryableError(fmt.Errorf("failed to create assignment: %s", err))
			}
			return nil
		})

		if err != nil {
			errors.HandleGraphError(ctx, err, resp, "Create", r.WritePermissions)
			return
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err = retry.RetryContext(ctx, retryTimeout, func() *retry.RetryError {
		readResp := &resource.ReadResponse{State: resp.State}
		r.Read(ctx, resource.ReadRequest{
			State:        resp.State,
			ProviderMeta: req.ProviderMeta,
		}, readResp)

		if readResp.Diagnostics.HasError() {
			return retry.NonRetryableError(fmt.Errorf("error reading resource state after Create Method: %s", readResp.Diagnostics.Errors()))
		}

		resp.State = readResp.State
		return nil
	})

	if err != nil {
		resp.Diagnostics.AddError(
			"Error waiting for resource creation",
			fmt.Sprintf("Failed to verify resource creation: %s", err),
		)
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("Finished Create Method: %s_%s", r.ProviderTypeName, r.TypeName))
}

// Read handles the Read operation.
func (r *MacOSPKGAppResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var object MacOSPKGAppResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting Read method for: %s_%s", r.ProviderTypeName, r.TypeName))

	resp.Diagnostics.Append(req.State.Get(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Reading %s_%s with ID: %s", r.ProviderTypeName, r.TypeName, object.ID.ValueString()))

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
		errors.HandleGraphError(ctx, err, resp, "Read", r.ReadPermissions)
		return
	}

	// This ensures type safety as the Graph API returns a base interface that needs
	// to be converted to the specific app type
	macOSPkgApp, ok := respBaseResource.(graphmodels.MacOSPkgAppable)
	if !ok {
		resp.Diagnostics.AddError(
			"Resource type mismatch",
			fmt.Sprintf("Expected resource of type MacOSPkgAppable but got %T", respBaseResource),
		)
		return
	}

	MapRemoteResourceStateToTerraform(ctx, &object, macOSPkgApp)

	// 2. get committed app content file versions
	committedVersionIdPtr := macOSPkgApp.GetCommittedContentVersion()

	if committedVersionIdPtr != nil && *committedVersionIdPtr != "" {
		committedVersionId := *committedVersionIdPtr

		respFiles, err := r.client.DeviceAppManagement().
			MobileApps().
			ByMobileAppId(object.ID.ValueString()).
			GraphMacOSPkgApp().
			ContentVersions().
			ByMobileAppContentId(committedVersionId).
			Files().
			Get(ctx, nil)

		if err != nil {
			errors.HandleGraphError(ctx, err, resp, "Read", r.ReadPermissions)
			return
		}

		var installerFileName string
		if macOSPkgApp.GetFileName() != nil {
			installerFileName = *macOSPkgApp.GetFileName()
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

	// 3. app assignments
	respAssignments, err := r.client.
		DeviceAppManagement().
		MobileApps().
		ByMobileAppId(object.ID.ValueString()).
		Assignments().
		Get(ctx, nil)

	if err != nil {
		errors.HandleGraphError(ctx, err, resp, "Read", r.ReadPermissions)
		return
	}

	object.Assignments = sharedstater.StateMobileAppAssignment(ctx, nil, respAssignments)

	// 4. Get app metadata by processing app installer file
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

	tflog.Debug(ctx, fmt.Sprintf("Finished Read Method: %s_%s", r.ProviderTypeName, r.TypeName))
}

// Update handles the Update operation.
func (r *MacOSPKGAppResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var object, state MacOSPKGAppResourceModel
	var installerSourcePath string
	var tempFileInfo helpers.TempFileInfo
	var err error

	tflog.Debug(ctx, fmt.Sprintf("Starting Update of resource: %s_%s", r.ProviderTypeName, r.TypeName))

	resp.Diagnostics.Append(req.Plan.Get(ctx, &object)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Update, UpdateTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	deadline, _ := ctx.Deadline()
	retryTimeout := time.Until(deadline) - time.Second
	// Ensure cleanup of temporary file when we're done
	if tempFileInfo.ShouldCleanup {
		defer helpers.CleanupTempFile(ctx, tempFileInfo)
	}
	installerSourcePath, tempFileInfo, err = helpers.SetInstallerSourcePath(ctx, object.AppInstaller)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error determining installer file path",
			err.Error(),
		)
		return
	}

	// Step 1: Update the base mobile app resource
	requestBody, err := constructResource(ctx, &object, installerSourcePath)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing resource for update method",
			fmt.Sprintf("Could not construct resource: %s_%s: %s", r.ProviderTypeName, r.TypeName, err.Error()),
		)
		return
	}

	_, err = r.client.
		DeviceAppManagement().
		MobileApps().
		ByMobileAppId(object.ID.ValueString()).
		Patch(ctx, requestBody, nil)

	if err != nil {
		errors.HandleGraphError(ctx, err, resp, "Update", r.WritePermissions)
		return
	}

	// Step 2: Updated Categories
	if !object.Categories.Equal(state.Categories) {
		tflog.Debug(ctx, "Categories have changed — updating categories")

		var categoryValues []string
		diags := object.Categories.ElementsAs(ctx, &categoryValues, false)
		if diags.HasError() {
			resp.Diagnostics.Append(diags...)
			return
		}

		err = construct.AssignMobileAppCategories(ctx, r.client, object.ID.ValueString(), categoryValues, r.ReadPermissions)
		if err != nil {
			errors.HandleGraphError(ctx, err, resp, "Update", r.WritePermissions)
			return
		}
	}

	// In the Update method in crud.go
	// Step 3: Updated Assignments
	if !state.ID.IsNull() {
		if len(object.Assignments) == 0 {
			tflog.Debug(ctx, "Empty assignments array detected. Removing all existing assignments individually")

			respAssignments, err := r.client.
				DeviceAppManagement().
				MobileApps().
				ByMobileAppId(object.ID.ValueString()).
				Assignments().
				Get(ctx, nil)

			if err != nil {
				errors.HandleGraphError(ctx, err, resp, "Update", r.WritePermissions)
				return
			}

			assignments := respAssignments.GetValue()

			if assignments == nil {
				tflog.Debug(ctx, "No assignments found to remove")
			} else {
				for _, assignment := range assignments {
					if assignment.GetId() == nil {
						continue // Skip assignments without an ID
					}

					assignmentId := *assignment.GetId()
					tflog.Debug(ctx, fmt.Sprintf("Deleting assignment with ID: %s", assignmentId))

					err := r.client.
						DeviceAppManagement().
						MobileApps().
						ByMobileAppId(object.ID.ValueString()).
						Assignments().
						ByMobileAppAssignmentId(assignmentId).
						Delete(ctx, nil)

					if err != nil {
						errors.HandleGraphError(ctx, err, resp, "Update", r.WritePermissions)
						return
					}

					tflog.Debug(ctx, fmt.Sprintf("Successfully deleted assignment with ID: %s", assignmentId))
				}

				tflog.Debug(ctx, "All assignments have been removed successfully")
			}
		} else {
			// Handle normal assignment update (non-empty assignments)
			requestAssignment, err := construct.ConstructMobileAppAssignment(ctx, object.Assignments)
			if err != nil {
				resp.Diagnostics.AddError(
					"Error constructing assignment for Update Method",
					fmt.Sprintf("Could not construct assignment: %s_%s: %s", r.ProviderTypeName, r.TypeName, err.Error()),
				)
				return
			}

			err = retry.RetryContext(ctx, retryTimeout, func() *retry.RetryError {
				err := r.client.
					DeviceAppManagement().
					MobileApps().
					ByMobileAppId(object.ID.ValueString()).
					Assign().
					Post(ctx, requestAssignment, nil)

				if err != nil {
					return retry.RetryableError(fmt.Errorf("failed to create assignment: %s", err))
				}
				return nil
			})

			if err != nil {
				errors.HandleGraphError(ctx, err, resp, "Update", r.WritePermissions)
				return
			}
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

	// Read updated resource state
	err = retry.RetryContext(ctx, retryTimeout, func() *retry.RetryError {
		readResp := &resource.ReadResponse{State: resp.State}
		r.Read(ctx, resource.ReadRequest{
			State:        resp.State,
			ProviderMeta: req.ProviderMeta,
		}, readResp)

		if readResp.Diagnostics.HasError() {
			return retry.NonRetryableError(fmt.Errorf("error reading resource state after Update Method: %s", readResp.Diagnostics.Errors()))
		}

		resp.State = readResp.State
		return nil
	})

	if err != nil {
		resp.Diagnostics.AddError(
			"Error waiting for resource update",
			fmt.Sprintf("Failed to verify resource update: %s", err),
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Update Method: %s_%s", r.ProviderTypeName, r.TypeName))
}

// Delete handles the Delete operation.
func (r *MacOSPKGAppResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var object MacOSPKGAppResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting deletion of resource: %s_%s", r.ProviderTypeName, r.TypeName))

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
		errors.HandleGraphError(ctx, err, resp, "Delete", r.WritePermissions)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Delete Method: %s_%s", r.ProviderTypeName, r.TypeName))

	resp.State.RemoveResource(ctx)
}
