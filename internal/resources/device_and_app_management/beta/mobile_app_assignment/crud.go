package graphBetaMobileAppAssignment

import (
	"context"
	"fmt"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/crud"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/deviceappmanagement"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// Create handles the Create operation.
func (r *MobileAppAssignmentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan MobileAppAssignmentResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting creation of resource: %s_%s", r.ProviderTypeName, r.TypeName))

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := crud.HandleTimeout(ctx, plan.Timeouts.Create, 30*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	mobileAppAssignment, err := constructResource(ctx, &plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing resource",
			fmt.Sprintf("Could not construct %s_%s: %s", r.ProviderTypeName, r.TypeName, err.Error()),
		)
		return
	}

	requestBody := deviceappmanagement.NewMobileAppsItemAssignPostRequestBody()
	requestBody.SetMobileAppAssignments([]models.MobileAppAssignmentable{mobileAppAssignment})

	err = r.client.DeviceAppManagement().
		MobileApps().
		ByMobileAppId(plan.SourceID.ValueString()).
		Assign().
		Post(ctx, requestBody, nil)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating resource",
			fmt.Sprintf("Could not create %s_%s: %s", r.ProviderTypeName, r.TypeName, err.Error()),
		)
		return
	}

	plan.ID = types.StringValue(plan.SourceID.ValueString())

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)

	tflog.Debug(ctx, fmt.Sprintf("Finished Create Method: %s_%s", r.ProviderTypeName, r.TypeName))
}

// Read handles the Read operation.
func (r *MobileAppAssignmentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state MobileAppAssignmentResourceModel
	tflog.Debug(ctx, fmt.Sprintf("Starting Read method for: %s_%s", r.ProviderTypeName, r.TypeName))

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Reading %s_%s with ID: %s", r.ProviderTypeName, r.TypeName, state.ID.ValueString()))

	ctx, cancel := crud.HandleTimeout(ctx, state.Timeouts.Read, 30*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	assignmentsResponse, err := r.client.DeviceAppManagement().MobileApps().ByMobileAppId(state.SourceID.ValueString()).Assignments().Get(ctx, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading mobile app assignments",
			fmt.Sprintf("Could not read assignments for mobile app %s: %s", state.SourceID.ValueString(), err.Error()),
		)
		return
	}

	// Get the actual assignments from the response
	assignments := assignmentsResponse.GetValue()

	// Find the specific assignment we're interested in
	var targetAssignment models.MobileAppAssignmentable
	for _, assignment := range assignments {
		if assignment.GetId() != nil && *assignment.GetId() == state.ID.ValueString() {
			targetAssignment = assignment
			break
		}
	}

	if targetAssignment == nil {
		resp.Diagnostics.AddError(
			"Error reading mobile app assignment",
			fmt.Sprintf("Could not find assignment with ID %s for mobile app %s", state.ID.ValueString(), state.SourceID.ValueString()),
		)
		return
	}

	MapRemoteStateToTerraform(ctx, &state, targetAssignment)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	tflog.Debug(ctx, fmt.Sprintf("Finished Read Method: %s_%s", r.ProviderTypeName, r.TypeName))
}

// Update handles the Update operation.
func (r *MobileAppAssignmentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan MobileAppAssignmentResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting Update of resource: %s_%s", r.ProviderTypeName, r.TypeName))

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := crud.HandleTimeout(ctx, plan.Timeouts.Update, 30*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	mobileAppAssignment, err := constructResource(ctx, &plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing resource for update method",
			fmt.Sprintf("Could not construct resource: %s_%s: %s", r.ProviderTypeName, r.TypeName, err.Error()),
		)
		return
	}

	requestBody := deviceappmanagement.NewMobileAppsItemAssignPostRequestBody()
	requestBody.SetMobileAppAssignments([]models.MobileAppAssignmentable{mobileAppAssignment})

	err = r.client.DeviceAppManagement().
		MobileApps().
		ByMobileAppId(plan.ID.ValueString()).
		Assign().
		Post(ctx, requestBody, nil)

	if err != nil {
		crud.HandleUpdateErrorIfNotFound(ctx, resp, r, &plan, err)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)

	tflog.Debug(ctx, fmt.Sprintf("Finished Update Method: %s_%s", r.ProviderTypeName, r.TypeName))
}

// Delete handles the Delete operation.
func (r *MobileAppAssignmentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data MobileAppAssignmentResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting deletion of resource: %s_%s", r.ProviderTypeName, r.TypeName))

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := crud.HandleTimeout(ctx, data.Timeouts.Delete, 30*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	err := r.client.DeviceAppManagement().MobileApps().ByMobileAppId(data.ID.ValueString()).Delete(ctx, nil)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when deleting %s_%s", r.ProviderTypeName, r.TypeName), err.Error())
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Delete Method: %s_%s", r.ProviderTypeName, r.TypeName))

	resp.State.RemoveResource(ctx)
}
