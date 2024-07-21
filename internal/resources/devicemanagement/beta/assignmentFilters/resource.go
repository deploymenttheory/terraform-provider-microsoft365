package assignmentFilter

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

var _ resource.Resource = &AssignmentFilterResource{}
var _ resource.ResourceWithImportState = &AssignmentFilterResource{}

func NewUserResource() resource.Resource {
	return &AssignmentFilterResource{
		ProviderTypeName: "microsoft365",
		TypeName:         "_device_and_app_management_assignment_filter",
	}
}

type AssignmentFilterResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ProviderTypeName string
	TypeName         string
}

type AssignmentFilterResourceModel struct {
	ID          types.String `tfsdk:"id"`
	DisplayName types.String `tfsdk:"display_name"`
	Description types.String `tfsdk:"description"`
	Platform    types.String `tfsdk:"platform"`
	Rule        types.String `tfsdk:"rule"`
}

// Metadata returns the resource type name.
func (r *AssignmentFilterResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_device_and_app_management_assignment_filter"
}

// ImportState imports the resource state.
func (r *AssignmentFilterResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
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
				Validators: []validator.String{
					platformValidator{},
				},
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
	ctx, cancel := context.WithTimeout(ctx, 2*time.Minute)
	defer cancel()

	var data AssignmentFilterResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting creation of resource: %s_%s", r.ProviderTypeName, r.TypeName))

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	requestBody := models.NewDeviceAndAppManagementAssignmentFilter()
	displayName := data.DisplayName.ValueString()
	requestBody.SetDisplayName(&displayName)

	description := data.Description.ValueString()
	requestBody.SetDescription(&description)

	platformStr := data.Platform.ValueString()
	platform, err := StringToDevicePlatformType(platformStr, supportedPlatformTypes)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating assignment filter",
			fmt.Sprintf("Invalid platform: %s", err.Error()),
		)
		return
	}
	requestBody.SetPlatform(platform)

	rule := data.Rule.ValueString()
	requestBody.SetRule(&rule)

	roleScopeTags := []string{"0"}
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

	tflog.Debug(ctx, fmt.Sprintf("Finished creation of resource: %s_%s", r.ProviderTypeName, r.TypeName))
}

// Read handles the Read operation.
func (r *AssignmentFilterResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Minute)
	defer cancel()

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
	platformStr, err := DevicePlatformTypeToString(filter.GetPlatform())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading assignment filter",
			fmt.Sprintf("Could not convert platform: %s", err.Error()),
		)
		return
	}
	data.Platform = types.StringValue(platformStr)
	data.Rule = types.StringValue(*filter.GetRule())
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Update handles the Update operation.
func (r *AssignmentFilterResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Minute)
	defer cancel()

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

	platformStr := data.Platform.ValueString()
	platform, err := StringToDevicePlatformType(platformStr, supportedPlatformTypes)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating assignment filter",
			fmt.Sprintf("Invalid platform: %s", err.Error()),
		)
		return
	}
	requestBody.SetPlatform(platform)

	rule := data.Rule.ValueString()
	requestBody.SetRule(&rule)

	roleScopeTags := []string{"0"} // Adjust if necessary
	requestBody.SetRoleScopeTags(roleScopeTags)

	_, err = r.client.DeviceManagement().AssignmentFilters().ByDeviceAndAppManagementAssignmentFilterId(data.ID.ValueString()).Patch(ctx, requestBody, nil)
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
	ctx, cancel := context.WithTimeout(ctx, 1*time.Minute)
	defer cancel()

	var data AssignmentFilterResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting deletion of resource: %s_%s", r.ProviderTypeName, r.TypeName))

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeviceManagement().AssignmentFilters().ByDeviceAndAppManagementAssignmentFilterId(data.ID.ValueString()).Delete(ctx, nil)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when deleting %s_%s", r.ProviderTypeName, r.TypeName), err.Error())
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Completed deletion of resource: %s_%s", r.ProviderTypeName, r.TypeName))

	resp.State.RemoveResource(ctx)
}
