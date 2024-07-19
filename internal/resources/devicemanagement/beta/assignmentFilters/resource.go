package assignmentFilter

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

var _ resource.Resource = &AssignmentFilterResource{}

type AssignmentFilterResource struct {
	client *msgraphbetasdk.GraphServiceClient
}

type AssignmentFilterResourceModel struct {
	ID          types.String `tfsdk:"id"`
	DisplayName types.String `tfsdk:"display_name"`
	Description types.String `tfsdk:"description"`
	Platform    types.String `tfsdk:"platform"`
	Rule        types.String `tfsdk:"rule"`
}

// Schema returns the schema for the resource.
func (r *AssignmentFilterResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The unique identifier of the assignment filter.",
			},
			"display_name": schema.StringAttribute{
				Required:    true,
				Description: "The display name of the assignment filter.",
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "The description of the assignment filter.",
			},
			"platform": schema.StringAttribute{
				Required:    true,
				Description: "The platform for the assignment filter.",
				Validators:  []validator.String{validatePlatform()},
			},
			"rule": schema.StringAttribute{
				Required:    true,
				Description: "The rule for the assignment filter.",
			},
		},
	}
}

// Create handles the Create operation.
func (r *AssignmentFilterResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data AssignmentFilterResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	requestBody := models.NewDeviceAndAppManagementAssignmentFilter()
	displayName := data.DisplayName.ValueString()
	requestBody.SetDisplayName(&displayName)

	description := data.Description.ValueString()
	requestBody.SetDescription(&description)

	platform := data.Platform.ValueString()
	requestBody.SetPlatform(&platform)

	rule := data.Rule.ValueString()
	requestBody.SetRule(&rule)

	roleScopeTags := []string{"0"} // Adjust if necessary
	requestBody.SetRoleScopeTags(roleScopeTags)

	assignmentFilter, err := r.client.DeviceManagement().AssignmentFilters().Post(ctx, requestBody, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating assignment filter",
			fmt.Sprintf("Could not create assignment filter: %s", err.Error()),
		)
		return
	}

	data.ID = types.StringValue(*assignmentFilter.GetId())
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Read handles the Read operation.
func (r *AssignmentFilterResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data AssignmentFilterResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	filter, err := r.client.DeviceManagement().AssignmentFilters().ByDeviceAndAppManagementAssignmentFilterId(data.ID.ValueString()).Get(ctx, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading assignment filter",
			fmt.Sprintf("Could not read assignment filter: %s", err.Error()),
		)
		return
	}

	data.DisplayName = types.StringValue(*filter.GetDisplayName())
	data.Description = types.StringValue(*filter.GetDescription())
	data.Platform = types.StringValue(*filter.GetPlatform())
	data.Rule = types.StringValue(*filter.GetRule())
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}


// Update handles the Update operation.
func (r *AssignmentFilterResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data AssignmentFilterResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	requestBody := models.NewDeviceAndAppManagementAssignmentFilter()
	displayName := data.DisplayName.ValueString()
	requestBody.SetDisplayName(&displayName)

	description := data.Description.ValueString()
	requestBody.SetDescription(&description)

	platform := data.Platform.ValueString()
	requestBody.SetPlatform(&platform)

	rule := data.Rule.ValueString()
	requestBody.SetRule(&rule)

	roleScopeTags := []string{"0"} // Adjust if necessary
	requestBody.SetRoleScopeTags(roleScopeTags)

	_, err := r.client.DeviceManagement().AssignmentFilters().ById(data.ID.ValueString()).Update(ctx, requestBody, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating assignment filter",
			fmt.Sprintf("Could not update assignment filter: %s", err.Error()),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Delete handles the Delete operation.
func (r *AssignmentFilterResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data AssignmentFilterResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeviceManagement().AssignmentFilters().ById(data.ID.ValueString()).Delete(ctx, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting assignment filter",
			fmt.Sprintf("Could not delete assignment filter: %s", err.Error()),
		)
		return
	}

	resp.State.RemoveResource(ctx)
}
