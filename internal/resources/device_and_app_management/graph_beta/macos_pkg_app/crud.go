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
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// Create handles the Create operation.
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

	createdResource, err := constructResource(ctx, &object)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing resource for Create method",
			fmt.Sprintf("Could not construct resource: %s_%s: %s", r.ProviderTypeName, r.TypeName, err.Error()),
		)
		return
	}

	baseResource, err := r.client.
		DeviceAppManagement().
		MobileApps().
		Post(ctx, createdResource, nil)
	if err != nil {
		errors.HandleGraphError(ctx, err, resp, "Create", r.WritePermissions)
		return
	}

	object.ID = types.StringValue(*baseResource.GetId())

	// Initialize content and upload file if source path is provided
	if !object.MacOSPkgApp.PackageInstallerFileSource.IsNull() {
		// Initialize content version
		content := graphmodels.NewMobileAppContent()
		contentBuilder := r.client.DeviceAppManagement().
			MobileApps().
			ByMobileAppId(object.ID.ValueString()).
			GraphMacOSPkgApp().
			ContentVersions()

		contentVersion, err := contentBuilder.Post(ctx, content, nil)
		if err != nil {
			errors.HandleGraphError(ctx, err, resp, "Create", r.WritePermissions)
			return
		}

		// Create and encrypt content file
		contentFile, encryptionInfo, err := constructMobileAppContentFile(ctx, object.MacOSPkgApp.PackageInstallerFileSource.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Error constructing and encrypting Intune Mobile App content file",
				err.Error(),
			)
			return
		}

		// Create the content file in Graph API
		createdFile, err := contentBuilder.
			ByMobileAppContentId(*contentVersion.GetId()).
			Files().
			Post(ctx, contentFile, nil)

		if err != nil {
			errors.HandleGraphError(ctx, err, resp, "Create", r.WritePermissions)
			return
		}

		// Wait for Azure Storage URI
		deadline, _ := ctx.Deadline()
		err = retry.RetryContext(ctx, time.Until(deadline), func() *retry.RetryError {
			file, err := contentBuilder.
				ByMobileAppContentId(*contentVersion.GetId()).
				Files().
				ByMobileAppContentFileId(*createdFile.GetId()).
				Get(ctx, nil)
			if err != nil {
				return retry.RetryableError(fmt.Errorf("failed to get file status: %v", err))
			}

			if file.GetUploadState() == nil {
				return retry.RetryableError(fmt.Errorf("upload state is nil"))
			}

			state := *file.GetUploadState()
			if state == graphmodels.AZURESTORAGEURIREQUESTSUCCESS_MOBILEAPPCONTENTFILEUPLOADSTATE {
				return nil
			}

			if state == graphmodels.AZURESTORAGEURIREQUESTFAILED_MOBILEAPPCONTENTFILEUPLOADSTATE {
				return retry.NonRetryableError(fmt.Errorf("azure storage URI request failed"))
			}

			return retry.RetryableError(fmt.Errorf("waiting for Azure Storage URI, current state: %s", state.String()))
		})

		if err != nil {
			resp.Diagnostics.AddError(
				"Error waiting for Azure Storage URI",
				err.Error(),
			)
			return
		}

		// Get the file status to retrieve Azure Storage URI
		fileStatus, err := contentBuilder.
			ByMobileAppContentId(*contentVersion.GetId()).
			Files().
			ByMobileAppContentFileId(*createdFile.GetId()).
			Get(ctx, nil)

		if err != nil {
			errors.HandleGraphError(ctx, err, resp, "Create", r.WritePermissions)
			return
		}

		// Upload encrypted file to Azure Storage using the SAS URI
		err = uploadToAzureStorage(ctx, *fileStatus.GetAzureStorageUri(), object.MacOSPkgApp.PackageInstallerFileSource.ValueString()+".bin")
		if err != nil {
			resp.Diagnostics.AddError(
				"Error uploading file to Azure Storage",
				err.Error(),
			)
			return
		}

		// Commit the file with encryption info
		err = retry.RetryContext(ctx, time.Until(deadline), func() *retry.RetryError {
			requestBody, err := NewEncryptedIntuneMobileAppFileUpload(encryptionInfo)
			if err != nil {
				return retry.NonRetryableError(fmt.Errorf("failed to construct commit request: %v", err))
			}

			err = contentBuilder.
				ByMobileAppContentId(*contentVersion.GetId()).
				Files().
				ByMobileAppContentFileId(*createdFile.GetId()).
				Commit().
				Post(ctx, requestBody, nil)
			if err != nil {
				return retry.RetryableError(fmt.Errorf("failed to commit file: %v", err))
			}
			return nil
		})

		if err != nil {
			resp.Diagnostics.AddError(
				"Error committing file",
				err.Error(),
			)
			return
		}

		// Wait for commit to complete
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
			if state == graphmodels.COMMITFILESUCCESS_MOBILEAPPCONTENTFILEUPLOADSTATE {
				break
			}

			if state == graphmodels.COMMITFILEFAILED_MOBILEAPPCONTENTFILEUPLOADSTATE {
				// Retry the upload
				requestBody, err := NewEncryptedIntuneMobileAppFileUpload(encryptionInfo)
				if err != nil {
					continue
				}

				err = contentBuilder.
					ByMobileAppContentId(*contentVersion.GetId()).
					Files().
					ByMobileAppContentFileId(*createdFile.GetId()).
					Commit().
					Post(ctx, requestBody, nil)
				if err != nil {
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

	// resource, err := r.client.
	// 	DeviceAppManagement().
	// 	MobileApps().
	// 	ByMobileAppId(object.ID.ValueString()).
	// 	Get(ctx, nil)

	// if err != nil {
	// 	errors.HandleGraphError(ctx, err, resp, "Read", r.ReadPermissions)
	// 	return
	// }

	// This ensures type safety as the Graph API returns a base interface that needs
	// to be converted to the specific app type to access WinGetApp-specific fields.
	// winGetApp, ok := resource.(graphmodels.WinGetAppable)
	// if !ok {
	// 	resp.Diagnostics.AddError(
	// 		"Resource type mismatch",
	// 		fmt.Sprintf("Expected resource of type WinGetAppable but got %T", resource),
	// 	)
	// 	return
	// }

	//MapRemoteResourceStateToTerraform(ctx, &object, winGetApp)

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

	// Only map assignments if there are any assignments returned and the object has an assignments block
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
func (r *MacOSPKGAppResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var object, state MacOSPKGAppResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting Update of resource: %s_%s", r.ProviderTypeName, r.TypeName))

	resp.Diagnostics.Append(req.Plan.Get(ctx, &object)...)
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

	requestBody, err := constructResource(ctx, &object)
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

	if object.Assignments != nil && state.Assignments != nil && !state.ID.IsNull() {
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
			errors.HandleGraphError(ctx, err, resp, "Create", r.WritePermissions)
			return
		}
	}

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
