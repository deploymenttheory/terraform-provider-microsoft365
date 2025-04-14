package graphBetaMacOSPKGApp

import (
	"context"
	"fmt"
	"time"

	construct "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/constructors/graph_beta/device_and_app_management"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/crud"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/errors"
	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/shared_models/graph_beta/device_and_app_management"
	sharedstater "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/state/graph_beta/device_and_app_management"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
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
	installerSourcePath, tempFileInfo, err := setInstallerSourcePath(ctx, object.AppMetadata)
	if err != nil {
		resp.Diagnostics.AddError("Error determining installer file path", err.Error())
		return
	}

	// Ensure cleanup of temporary file occurs post state read
	if tempFileInfo.ShouldCleanup {
		defer cleanupTempFile(ctx, tempFileInfo)
	}

	// Step 3: Construct the base resource from the Terraform model
	createdResource, err := constructResource(ctx, &object, installerSourcePath)
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
		Post(ctx, createdResource, nil)
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

		err = assignCategoriesToMobileApplication(ctx, r.client, object.ID.ValueString(), categoryValues, r.ReadPermissions)
		if err != nil {
			errors.HandleGraphError(ctx, err, resp, "Associate Categories", r.WritePermissions)
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
		contentFile, encryptionInfo, err := encryptMobileAppAndConstructFileContentMetadata(ctx, installerSourcePath)
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
		err = uploadToAzureStorage(ctx, *fileStatus.GetAzureStorageUri(), encryptedFilePath)
		if err != nil {
			tflog.Debug(ctx, fmt.Sprintf("Failed to upload to Azure Storage: %v", err))
			errors.HandleGraphError(ctx, err, resp, "Create", r.WritePermissions)
			return
		}

		// Step 11: Commit the file with encryption metadata
		tflog.Debug(ctx, "Committing file with encryption metadata")
		err = retry.RetryContext(ctx, time.Until(deadline), func() *retry.RetryError {
			commitBody, err := CommitUploadedMobileAppWithEncryptionMetadata(encryptionInfo)
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
		maxRetries := 10
		for i := 0; i < maxRetries; i++ {
			file, err := contentBuilder.
				ByMobileAppContentId(*contentVersion.GetId()).
				Files().
				ByMobileAppContentFileId(*createdFile.GetId()).
				Get(ctx, nil)
			if err != nil {
				errors.HandleGraphError(ctx, err, resp, "Create", r.WritePermissions)
				return
			}
			state := *file.GetUploadState()
			tflog.Debug(ctx, fmt.Sprintf("Commit status check %d: state=%s", i+1, state.String()))
			if state == graphmodels.COMMITFILESUCCESS_MOBILEAPPCONTENTFILEUPLOADSTATE {
				tflog.Debug(ctx, "File commit completed successfully")
				break
			}
			if state == graphmodels.COMMITFILEFAILED_MOBILEAPPCONTENTFILEUPLOADSTATE {
				tflog.Debug(ctx, "File commit failed; retrying commit request")
				commitBody, err := CommitUploadedMobileAppWithEncryptionMetadata(encryptionInfo)
				if err != nil {
					tflog.Debug(ctx, fmt.Sprintf("Error constructing commit request during retry: %v", err))
					continue
				}
				err = contentBuilder.
					ByMobileAppContentId(*contentVersion.GetId()).
					Files().
					ByMobileAppContentFileId(*createdFile.GetId()).
					Commit().
					Post(ctx, commitBody, nil)
				if err != nil {
					tflog.Debug(ctx, fmt.Sprintf("Error during commit retry: %v", err))
					continue
				}
			}
			if i == maxRetries-1 {
				resp.Diagnostics.AddError(
					"Error waiting for file commit",
					fmt.Sprintf("File commit did not complete after %d attempts. Last state: %s", maxRetries, state.String()),
				)
				return
			}
			time.Sleep(10 * time.Second)
		}

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
	if object.Assignments != nil {
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

	// 2. content versions and it's files
	respContentVersions, err := r.client.DeviceAppManagement().
		MobileApps().
		ByMobileAppId(object.ID.ValueString()).
		GraphMacOSPkgApp().
		ContentVersions().
		Get(ctx, nil)

	if err != nil {
		errors.HandleGraphError(ctx, err, resp, "Read", r.ReadPermissions)
		return
	}

	var allContentVersionFiles = make(map[string][]graphmodels.MobileAppContentFileable)

	for _, version := range respContentVersions.GetValue() {
		if version == nil || version.GetId() == nil {
			continue
		}

		respFiles, err := r.client.DeviceAppManagement().
			MobileApps().
			ByMobileAppId(object.ID.ValueString()).
			GraphMacOSPkgApp().
			ContentVersions().
			ByMobileAppContentId(*version.GetId()).
			Files().
			Get(ctx, nil)

		if err != nil {
			tflog.Warn(ctx, fmt.Sprintf("Failed to retrieve files for content version %s: %v", *version.GetId(), err))
			continue
		}

		allContentVersionFiles[*version.GetId()] = respFiles.GetValue()
	}

	object.ContentVersion = MapContentVersionsStateToTerraform(ctx, respContentVersions.GetValue(), allContentVersionFiles)

	// 3. app assignments
	respAssignments, err := r.client.
		DeviceAppManagement().
		MobileApps().
		ByMobileAppId(object.ID.ValueString()).
		Assignments().
		Get(ctx, nil)
	if err != nil {
		errors.HandleGraphError(ctx, err, resp, "Read Assignments", r.ReadPermissions)
		return
	}

	object.Assignments = sharedstater.StateMobileAppAssignment(ctx, nil, respAssignments)

	// 4. Get app metadata by processing app installer file
	installerPath, tempInfo, err := setInstallerSourcePath(ctx, object.AppMetadata)
	if err != nil {
		resp.Diagnostics.AddError("Error determining installer path", err.Error())
		return
	}
	if tempInfo.ShouldCleanup {
		defer cleanupTempFile(ctx, tempInfo)
	}

	var existingMetadata sharedmodels.MobileAppMetaDataResourceModel
	if !req.State.Raw.IsNull() {
		diags := req.State.GetAttribute(ctx, path.Root("app_metadata"), &existingMetadata)
		if diags.HasError() {
			resp.Diagnostics.Append(diags...)
			return
		}
	}

	metadata, err := GetAppMetadata(ctx, installerPath, &existingMetadata)
	if err != nil {
		resp.Diagnostics.AddError("Error capturing installer metadata", err.Error())
		return
	}

	object.AppMetadata = MapAppMetadataStateToTerraform(ctx, metadata)

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
	var tempFileInfo TempFileInfo
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
		defer cleanupTempFile(ctx, tempFileInfo)
	}

	installerSourcePath, tempFileInfo, err = setInstallerSourcePath(ctx, object.AppMetadata)
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

		err = assignCategoriesToMobileApplication(ctx, r.client, object.ID.ValueString(), categoryValues, r.ReadPermissions)
		if err != nil {
			errors.HandleGraphError(ctx, err, resp, "Update Categories", r.WritePermissions)
			return
		}
	}

	// Step 3:  Updated Assignments
	if object.Assignments != nil && !state.ID.IsNull() {
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

	// Step 4:  Updated application content version and files
	if installerSourcePath != "" {

		// Evaluate if content update is needed
		contentUpdateNeeded, existingContentVersion, err := evaluateIfContentVersionUpdateRequired(ctx, &object, &state, r.client)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error evaluating if content update is needed",
				err.Error(),
			)
			return
		}

		// --- Begin File Processing Logic ---
		// Process content upload only if content has changed or we need a new version
		if contentUpdateNeeded {
			tflog.Debug(ctx, "Content has changed — deleting existing content versions before upload")

			// Step 4b
			tflog.Debug(ctx, "Initializing new content version for file upload")
			content := graphmodels.NewMobileAppContent()
			contentBuilder := r.client.
				DeviceAppManagement().
				MobileApps().
				ByMobileAppId(object.ID.ValueString()).
				GraphMacOSPkgApp().
				ContentVersions()

			contentVersion, err := contentBuilder.Post(ctx, content, nil)
			if err != nil {
				errors.HandleGraphError(ctx, err, resp, "Update", r.WritePermissions)
				return
			}
			tflog.Debug(ctx, fmt.Sprintf("New content version created with ID: %s", *contentVersion.GetId()))

			// Step 4c: Encrypt installer file and construct file metadata
			tflog.Debug(ctx, "Encrypting installer file and constructing file metadata")
			contentFile, encryptionInfo, err := encryptMobileAppAndConstructFileContentMetadata(ctx, installerSourcePath)
			if err != nil {
				resp.Diagnostics.AddError(
					"Error constructing and encrypting Intune Mobile App content file",
					err.Error(),
				)
				return
			}

			// Step 4d: Create the content file resource in Graph API
			tflog.Debug(ctx, "Creating content file resource in Graph API")
			createdFile, err := contentBuilder.
				ByMobileAppContentId(*contentVersion.GetId()).
				Files().
				Post(ctx, contentFile, nil)
			if err != nil {
				errors.HandleGraphError(ctx, err, resp, "Update", r.WritePermissions)
				return
			}
			tflog.Debug(ctx, fmt.Sprintf("Content file resource created with ID: %s", *createdFile.GetId()))

			// Step 5d: Wait for Graph API to generate a valid Azure Storage URI
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

				stateVal := *file.GetUploadState()
				tflog.Debug(ctx, fmt.Sprintf("Current upload state: %s", stateVal.String()))

				if stateVal == graphmodels.AZURESTORAGEURIREQUESTSUCCESS_MOBILEAPPCONTENTFILEUPLOADSTATE {
					tflog.Debug(ctx, "Azure Storage URI request successful")
					return nil
				}

				if stateVal == graphmodels.AZURESTORAGEURIREQUESTFAILED_MOBILEAPPCONTENTFILEUPLOADSTATE {
					tflog.Debug(ctx, "Azure Storage URI request failed")
					return retry.NonRetryableError(fmt.Errorf("azure storage URI request failed"))
				}

				tflog.Debug(ctx, fmt.Sprintf("Waiting for Azure Storage URI, current state: %s", stateVal.String()))
				return retry.RetryableError(fmt.Errorf("waiting for Azure Storage URI, current state: %s", stateVal.String()))
			})
			if err != nil {
				errors.HandleGraphError(ctx, err, resp, "Update", r.WritePermissions)
				return
			}

			// Step 5e: Retrieve the Azure Storage URI and upload the encrypted file
			tflog.Debug(ctx, "Retrieving file status for Azure Storage URI")
			fileStatus, err := contentBuilder.
				ByMobileAppContentId(*contentVersion.GetId()).
				Files().
				ByMobileAppContentFileId(*createdFile.GetId()).
				Get(ctx, nil)
			if err != nil {
				errors.HandleGraphError(ctx, err, resp, "Update", r.WritePermissions)
				return
			}
			if fileStatus.GetAzureStorageUri() == nil {
				tflog.Debug(ctx, "Azure Storage URI is nil in the retrieved file status")
				errors.HandleGraphError(ctx, fmt.Errorf("azure Storage URI is nil"), resp, "Update", r.WritePermissions)
				return
			}
			tflog.Debug(ctx, fmt.Sprintf("Retrieved Azure Storage URI: %s", *fileStatus.GetAzureStorageUri()))

			// IMPORTANT: Upload the encrypted file (.bin) not the original source file.
			encryptedFilePath := installerSourcePath + ".bin"
			tflog.Debug(ctx, fmt.Sprintf("Uploading encrypted file: %s", encryptedFilePath))
			err = uploadToAzureStorage(ctx, *fileStatus.GetAzureStorageUri(), encryptedFilePath)
			if err != nil {
				tflog.Debug(ctx, fmt.Sprintf("Failed to upload to Azure Storage: %v", err))
				errors.HandleGraphError(ctx, err, resp, "Update", r.WritePermissions)
				return
			}

			// Step 5f: Commit the file with encryption metadata
			tflog.Debug(ctx, "Committing file with encryption metadata")
			err = retry.RetryContext(ctx, time.Until(deadline), func() *retry.RetryError {
				commitBody, err := CommitUploadedMobileAppWithEncryptionMetadata(encryptionInfo)
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

			// Step 5g: Wait for commit to complete
			tflog.Debug(ctx, "Waiting for file commit to complete")
			maxRetries := 10
			for i := 0; i < maxRetries; i++ {
				file, err := contentBuilder.
					ByMobileAppContentId(*contentVersion.GetId()).
					Files().
					ByMobileAppContentFileId(*createdFile.GetId()).
					Get(ctx, nil)
				if err != nil {
					errors.HandleGraphError(ctx, err, resp, "Update", r.WritePermissions)
					return
				}
				stateVal := *file.GetUploadState()
				tflog.Debug(ctx, fmt.Sprintf("Commit status check %d: state=%s", i+1, stateVal.String()))
				if stateVal == graphmodels.COMMITFILESUCCESS_MOBILEAPPCONTENTFILEUPLOADSTATE {
					tflog.Debug(ctx, "File commit completed successfully")
					break
				}
				if stateVal == graphmodels.COMMITFILEFAILED_MOBILEAPPCONTENTFILEUPLOADSTATE {
					tflog.Debug(ctx, "File commit failed; retrying commit request")
					commitBody, err := CommitUploadedMobileAppWithEncryptionMetadata(encryptionInfo)
					if err != nil {
						tflog.Debug(ctx, fmt.Sprintf("Error constructing commit request during retry: %v", err))
						continue
					}
					err = contentBuilder.
						ByMobileAppContentId(*contentVersion.GetId()).
						Files().
						ByMobileAppContentFileId(*createdFile.GetId()).
						Commit().
						Post(ctx, commitBody, nil)
					if err != nil {
						tflog.Debug(ctx, fmt.Sprintf("Error during commit retry: %v", err))
						continue
					}
				}
				if i == maxRetries-1 {
					resp.Diagnostics.AddError(
						"Error waiting for file commit",
						fmt.Sprintf("File commit did not complete after %d attempts. Last state: %s", maxRetries, stateVal.String()),
					)
					return
				}
				time.Sleep(10 * time.Second)
			}

			// Step 5h: Update the App with the new Committed Content Version
			updatePayload := graphmodels.NewMacOSPkgApp()
			updatePayload.SetCommittedContentVersion(contentVersion.GetId())
			_, err = r.client.
				DeviceAppManagement().
				MobileApps().
				ByMobileAppId(object.ID.ValueString()).
				Patch(ctx, updatePayload, nil)
			if err != nil {
				errors.HandleGraphError(ctx, err, resp, "Update", r.WritePermissions)
				return
			}
			tflog.Debug(ctx, fmt.Sprintf("App updated with new committed content version: %s", *contentVersion.GetId()))
		} else if existingContentVersion != "" {
			// Content hasn't changed, preserve the existing content version
			tflog.Debug(ctx, fmt.Sprintf(
				"File content unchanged, preserving existing content version: %s",
				existingContentVersion))

			// Ensure that the committed content version is set in the app
			// This handles edge cases where another process might have cleared it
			updatePayload := graphmodels.NewMacOSPkgApp()
			existingVersionPtr := &existingContentVersion
			updatePayload.SetCommittedContentVersion(existingVersionPtr)
			_, err = r.client.
				DeviceAppManagement().
				MobileApps().
				ByMobileAppId(object.ID.ValueString()).
				Patch(ctx, updatePayload, nil)
			if err != nil {
				errors.HandleGraphError(ctx, err, resp, "Update", r.WritePermissions)
				return
			}
			tflog.Debug(ctx, "Ensured existing content version is still committed")
		}
		// --- End File Processing Logic ---
	}

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
