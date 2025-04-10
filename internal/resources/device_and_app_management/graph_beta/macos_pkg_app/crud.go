package graphBetaMacOSPKGApp

import (
	"context"
	"fmt"
	"os"
	"time"

	construct "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/constructors/graph_beta/device_and_app_management"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/crud"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/errors"
	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/shared_models/graph_beta/device_and_app_management"
	sharedstater "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/state/graph_beta/device_and_app_management"
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
	installerSourcePath, tempFileInfo, err := setInstallerSourcePath(ctx, object.MacOSPkgApp)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error determining installer file path",
			err.Error(),
		)
		return
	}

	// Ensure cleanup of temporary file when we're done
	if tempFileInfo.ShouldCleanup {
		defer cleanupTempFile(ctx, tempFileInfo)
	}

	// Get file size and store it in the model
	fileInfo, err := os.Stat(installerSourcePath)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting file information",
			err.Error(),
		)
		return
	}
	currentSize := fileInfo.Size()
	object.InstallerSizeInBytes = types.Int64Value(currentSize)

	// Step 2: Construct the base resource from the Terraform model
	createdResource, err := constructResource(ctx, &object, installerSourcePath)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing resource for Create method",
			fmt.Sprintf("Could not construct resource: %s_%s: %s", r.ProviderTypeName, r.TypeName, err.Error()),
		)
		return
	}

	// Step 3: Create the mobile app resource in Graph API
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

	// Step 4: Associate categories with the app if provided
	if !object.Categories.IsNull() {
		var categoryValues []string
		diags := object.Categories.ElementsAs(ctx, &categoryValues, false)
		if diags.HasError() {
			resp.Diagnostics.Append(diags...)
			return
		}

		err = associateAppWithCategories(ctx, r.client, object.ID.ValueString(), categoryValues, r.ReadPermissions)
		if err != nil {
			errors.HandleGraphError(ctx, err, resp, "Associate Categories", r.WritePermissions)
			return
		}
	}

	// If a package installer file is provided, process the content version and file upload
	if (!object.MacOSPkgApp.InstallerFilePathSource.IsNull() && object.MacOSPkgApp.InstallerFilePathSource.ValueString() != "") ||
		(!object.MacOSPkgApp.InstallerURLSource.IsNull() && object.MacOSPkgApp.InstallerURLSource.ValueString() != "") {

		// Step 5: Initialize content version
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

		// Step 6: Encrypt mobile app file and output file and encryption metadata
		tflog.Debug(ctx, "Encrypting installer file and constructing file metadata")
		contentFile, encryptionInfo, err := encryptMobileAppAndConstructFileContentMetadata(ctx, installerSourcePath)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error constructing and encrypting Intune Mobile App content file",
				err.Error(),
			)
			return
		}

		// Step 7: Create the content file resource in Graph API
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

		// Step 8: Wait for Graph API to generate a valid Azure Storage URI
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

		// Step 9: Retrieve the Azure Storage URI and upload the encrypted file
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

		// Step 10: Commit the file with encryption metadata
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

		// Step 11: Wait for commit to complete
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

		// Step 12: Update the App with the Committed Content Version
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

	// Create request configuration with expand query parameter
	requestParameters := &deviceappmanagement.MobileAppsMobileAppItemRequestBuilderGetRequestConfiguration{
		QueryParameters: &deviceappmanagement.MobileAppsMobileAppItemRequestBuilderGetQueryParameters{
			Expand: []string{"categories"},
		},
	}

	resource, err := r.client.
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
	macOSPkgApp, ok := resource.(graphmodels.MacOSPkgAppable)
	if !ok {
		resp.Diagnostics.AddError(
			"Resource type mismatch",
			fmt.Sprintf("Expected resource of type MacOSPkgAppable but got %T", resource),
		)
		return
	}

	MapRemoteResourceStateToTerraform(ctx, &object, macOSPkgApp)

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

	if respAssignments != nil && len(respAssignments.GetValue()) > 0 {
		object.Assignments = make([]sharedmodels.MobileAppAssignmentResourceModel, len(respAssignments.GetValue()))
		sharedstater.StateMobileAppAssignment(ctx, object.Assignments, respAssignments)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Read Method: %s_%s", r.ProviderTypeName, r.TypeName))
}

// Update handles the Update operation.
// Update handles the Update operation.
func (r *MacOSPKGAppResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var object, state MacOSPKGAppResourceModel

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

	// Step 1: Determine installer source path (local or download via URL)
	var installerSourcePath string
	var tempFileInfo TempFileInfo
	var err error
	var needsContentUpload bool = false

	if (!object.MacOSPkgApp.InstallerFilePathSource.IsNull() && object.MacOSPkgApp.InstallerFilePathSource.ValueString() != "") ||
		(!object.MacOSPkgApp.InstallerURLSource.IsNull() && object.MacOSPkgApp.InstallerURLSource.ValueString() != "") {

		installerSourcePath, tempFileInfo, err = setInstallerSourcePath(ctx, object.MacOSPkgApp)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error determining installer file path",
				err.Error(),
			)
			return
		}

		// Ensure cleanup of temporary file when we're done
		if tempFileInfo.ShouldCleanup {
			defer cleanupTempFile(ctx, tempFileInfo)
		}

		// Get file size and store it in the model
		fileInfo, err := os.Stat(installerSourcePath)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error getting file information",
				err.Error(),
			)
			return
		}
		currentSize := fileInfo.Size()
		object.InstallerSizeInBytes = types.Int64Value(currentSize)

		// Determine if we need to process a content upload based on file size
		needsContentUpload = true

		// Check if we have a previous size in state to compare against
		if !state.InstallerSizeInBytes.IsNull() {
			previousSize := state.InstallerSizeInBytes.ValueInt64()

			// If size is the same, check if there's a committed content version already
			if previousSize == currentSize {
				// Retrieve the resource to check for committed content version
				resource, err := r.client.
					DeviceAppManagement().
					MobileApps().
					ByMobileAppId(object.ID.ValueString()).
					Get(ctx, nil)

				if err == nil {
					// This ensures type safety as the Graph API returns a base interface
					macOSPkgApp, ok := resource.(graphmodels.MacOSPkgAppable)
					if ok && macOSPkgApp.GetCommittedContentVersion() != nil &&
						*macOSPkgApp.GetCommittedContentVersion() != "" {
						// If size is unchanged and we have a committed version, skip upload
						needsContentUpload = false
						tflog.Debug(ctx, fmt.Sprintf(
							"File size unchanged (%d bytes) and committed version exists, skipping content upload",
							currentSize))
					}
				}
			} else {
				tflog.Debug(ctx, fmt.Sprintf(
					"File size changed (previous: %d, current: %d bytes), will process content upload",
					previousSize, currentSize))
			}
		} else {
			tflog.Debug(ctx, "No previous file size in state, will process content upload")
		}
	}

	// Step 2: Construct the base resource with the resolved installer path
	requestBody, err := constructResource(ctx, &object, installerSourcePath)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing resource for update method",
			fmt.Sprintf("Could not construct resource: %s_%s: %s", r.ProviderTypeName, r.TypeName, err.Error()),
		)
		return
	}

	// Step 3: Update the base mobile app resource
	_, err = r.client.
		DeviceAppManagement().
		MobileApps().
		ByMobileAppId(object.ID.ValueString()).
		Patch(ctx, requestBody, nil)

	if err != nil {
		errors.HandleGraphError(ctx, err, resp, "Update", r.WritePermissions)
		return
	}

	// --- Begin File Processing Logic ---
	// If a package installer file is provided and content upload is needed, process the file upload
	if installerSourcePath != "" && needsContentUpload {
		// Step 4: Initialize content version for file upload
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
			errors.HandleGraphError(ctx, err, resp, "Update", r.WritePermissions)
			return
		}
		tflog.Debug(ctx, fmt.Sprintf("Content version created with ID: %s", *contentVersion.GetId()))

		// Step 5: Encrypt installer file and construct file metadata
		tflog.Debug(ctx, "Encrypting installer file and constructing file metadata")
		contentFile, encryptionInfo, err := encryptMobileAppAndConstructFileContentMetadata(ctx, installerSourcePath)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error constructing and encrypting Intune Mobile App content file",
				err.Error(),
			)
			return
		}

		// Step 6: Create the content file resource in Graph API
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

		// Step 7: Wait for Graph API to generate a valid Azure Storage URI
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

		// Step 8: Retrieve the Azure Storage URI and upload the encrypted file
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

		// Step 9: Commit the file with encryption metadata
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

		// Step 10: Wait for commit to complete
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

		// Step 11: Update the App with the Committed Content Version
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
		tflog.Debug(ctx, "App updated successfully with committed content version")
	} else if installerSourcePath != "" && !needsContentUpload {
		tflog.Debug(ctx, "Skipping content upload as file content has not changed")
	}
	// --- End File Processing Logic ---

	// Process assignments if provided
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
