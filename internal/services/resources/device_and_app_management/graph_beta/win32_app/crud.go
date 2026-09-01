package graphBetaWin32App

import (
	"context"
	"fmt"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	construct "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors/graph_beta/device_and_app_management"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	identitymodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/shared_models/graph_beta"
	sharedstater "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/state/graph_beta/device_and_app_management"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/deviceappmanagement"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// Create handles the Create operation for Win32 App resources.
//
// Operation: Creates a new Win32 application with content file upload workflow
// API Calls:
//   - POST /deviceAppManagement/mobileApps
//   - POST /deviceAppManagement/mobileApps/{mobileAppId}/microsoft.graph.win32LobApp/contentVersions
//   - POST /deviceAppManagement/mobileApps/{mobileAppId}/microsoft.graph.win32LobApp/contentVersions/{contentVersionId}/files
//   - GET /deviceAppManagement/mobileApps/{mobileAppId}/microsoft.graph.win32LobApp/contentVersions/{contentVersionId}/files/{fileId}
//   - POST /deviceAppManagement/mobileApps/{mobileAppId}/microsoft.graph.win32LobApp/contentVersions/{contentVersionId}/files/{fileId}/commit
//   - PATCH /deviceAppManagement/mobileApps/{mobileAppId}
//
// Reference: https://learn.microsoft.com/en-us/graph/api/intune-apps-win32lobapp-create?view=graph-rest-beta
// Note: Includes encrypted file upload to Azure Storage and content version management
func (r *Win32LobAppResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var object Win32LobAppResourceModel

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

	installerPackage, err := prepareInstaller(ctx, &object)
	if err != nil {
		resp.Diagnostics.AddError("Error preparing Win32 application content", err.Error())
		return
	}
	defer cleanupWin32Content(ctx, installerPackage)

	// Step 3: Construct the base resource from the Terraform model
	requestBody, err := constructResource(ctx, &object, installerPackage.fileName)
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

	if !r.publishContentVersion(ctx, object.ID.ValueString(), installerPackage, resp, constants.TfOperationCreate) {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	readReq := resource.ReadRequest{State: resp.State, ProviderMeta: req.ProviderMeta}
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

// Read handles the Read operation for Win32 App resources.
//
// Operation: Retrieves a Win32 application by ID
// API Calls:
//   - GET /deviceAppManagement/mobileApps/{mobileAppId}
//
// Reference: https://learn.microsoft.com/en-us/graph/api/intune-apps-win32lobapp-get?view=graph-rest-beta
func (r *Win32LobAppResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var object Win32LobAppResourceModel
	var identity identitymodels.ResourceIdentity

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

	identity.ID = object.ID.ValueString()

	if resp.Identity != nil {
		resp.Diagnostics.Append(resp.Identity.Set(ctx, identity)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

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
	win32LobApp, ok := respBaseResource.(graphmodels.Win32LobAppable)
	if !ok {
		resp.Diagnostics.AddError(
			"Resource type mismatch",
			fmt.Sprintf("Expected resource of type MacOSPkgAppable but got %T", respBaseResource),
		)
		return
	}

	MapRemoteResourceStateToTerraform(ctx, &object, win32LobApp)

	// 2. get committed app content file versions
	committedVersionIdPtr := win32LobApp.GetCommittedContentVersion()

	if committedVersionIdPtr != nil && *committedVersionIdPtr != "" {
		committedVersionId := *committedVersionIdPtr

		respFiles, err := r.client.DeviceAppManagement().
			MobileApps().
			ByMobileAppId(object.ID.ValueString()).
			GraphWin32LobApp().
			ContentVersions().
			ByMobileAppContentId(committedVersionId).
			Files().
			Get(ctx, nil)

		if err != nil {
			errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationRead, r.ReadPermissions)
			return
		}

		var installerFileName string
		if win32LobApp.GetFileName() != nil {
			installerFileName = *win32LobApp.GetFileName()
		}

		object.ContentVersion = sharedstater.MapCommittedContentVersionStateToTerraform(ctx, committedVersionId, respFiles, err, installerFileName)
	}

	// Installer sources are configuration-only values. Preserve both blocks from
	// the existing state; Graph cannot reconstruct local paths or download URLs.

	// 6. set final state
	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Read Method: %s", ResourceName))
}

// Update handles the Update operation for Win32 App resources.
//
// Operation: Updates an existing Win32 application
// API Calls:
//   - PATCH /deviceAppManagement/mobileApps/{mobileAppId}
//
// Reference: https://learn.microsoft.com/en-us/graph/api/intune-apps-win32lobapp-update?view=graph-rest-beta
func (r *Win32LobAppResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state Win32LobAppResourceModel

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

	var installerPackage *preparedWin32Content
	installerFileName := ""
	if installerSourceChanged(&plan, &state) && hasInstallerSource(&plan) {
		var err error
		installerPackage, err = prepareInstaller(ctx, &plan)
		if err != nil {
			resp.Diagnostics.AddError("Error preparing Win32 application content", err.Error())
			return
		}
		defer cleanupWin32Content(ctx, installerPackage)
		installerFileName = installerPackage.fileName
	}

	// Step 1: Update the base mobile app resource
	requestBody, err := constructResource(ctx, &plan, installerFileName)
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

	if installerPackage != nil {
		tflog.Debug(ctx, "Installer source changed; publishing a new content version on the existing application")
		if !r.publishContentVersion(ctx, state.ID.ValueString(), installerPackage, resp, constants.TfOperationUpdate) {
			return
		}
	}

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

// Delete handles the Delete operation for Win32 App resources.
//
// Operation: Deletes a Win32 application
// API Calls:
//   - DELETE /deviceAppManagement/mobileApps/{mobileAppId}
//
// Reference: https://learn.microsoft.com/en-us/graph/api/intune-apps-win32lobapp-delete?view=graph-rest-beta
func (r *Win32LobAppResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var object Win32LobAppResourceModel

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
