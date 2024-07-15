package DeviceManagementScript

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models/odataerrors"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/devicemanagement"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ resource.Resource = &DeviceManagementScriptResource{}
var _ resource.ResourceWithImportState = &DeviceManagementScriptResource{}

func NewDeviceManagementScriptResource() resource.Resource {
	return &DeviceManagementScriptResource{}
}

// DeviceManagementScriptResource defines the resource implementation.
type DeviceManagementScriptResource struct {
	client *devicemanagement.DeviceManagementScriptsDeviceManagementScriptItemRequestBuilder
}

func (r *DeviceManagementScriptResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_device_management_script"
}

func (r *DeviceManagementScriptResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Manages a device management script.",
		MarkdownDescription: "The resource `device_management_script` manages a device management script.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description:         "ID of the device management script.",
				MarkdownDescription: "`ID` of the device management script.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description:         "Name of the device management script.",
				MarkdownDescription: "Name of the device management script.",
				Required:            true,
			},
			"description": schema.StringAttribute{
				Description:         "Description of the device management script.",
				MarkdownDescription: "Description of the device management script.",
				Optional:            true,
			},
			"script_content": schema.StringAttribute{
				Description:         "The content of the device management script.",
				MarkdownDescription: "The content of the device management script.",
				Required:            true,
			},
			"run_as_account": schema.StringAttribute{
				Description:         "The account to run the script as.",
				MarkdownDescription: "The account to run the script as.",
				Optional:            true,
			},
		},
	}
}

func (r *DeviceManagementScriptResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*devicemanagement.DeviceManagementScriptsDeviceManagementScriptItemRequestBuilder)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *devicemanagement.DeviceManagementScriptsDeviceManagementScriptItemRequestBuilder, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *DeviceManagementScriptResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data deviceManagementScript

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	script := models.NewDeviceManagementScript()
	script.SetName(data.Name.ValueString())
	script.SetDescription(data.Description.ValueString())
	script.SetScriptContent(data.ScriptContent.ValueString())
	script.SetRunAsAccount(data.RunAsAccount.ValueString())

	result, err := r.client.Create(ctx, script)
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to create device management script, got error: %s", err),
		)
		return
	}

	tflog.Trace(ctx, "created a device management script")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, deviceManagementScriptForState(result))...)
}

func (r *DeviceManagementScriptResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data deviceManagementScript

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	result, err := r.client.Get(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to read device management script %s, got error: %s", data.ID.ValueString(), err),
		)
		return
	}

	tflog.Trace(ctx, "read a device management script")

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, deviceManagementScriptForState(result))...)
}

func (r *DeviceManagementScriptResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data deviceManagementScript

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	script := models.NewDeviceManagementScript()
	script.SetName(data.Name.ValueString())
	script.SetDescription(data.Description.ValueString())
	script.SetScriptContent(data.ScriptContent.ValueString())
	script.SetRunAsAccount(data.RunAsAccount.ValueString())

	result, err := r.client.Update(ctx, data.ID.ValueString(), script)
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to update device management script %s, got error: %s", data.ID.ValueString(), err),
		)
		return
	}

	tflog.Trace(ctx, "updated a device management script")

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, deviceManagementScriptForState(result))...)
}

func (r *DeviceManagementScriptResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data deviceManagementScript

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.Delete(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to delete device management script %s, got error: %s", data.ID.ValueString(), err),
		)
		return
	}

	tflog.Trace(ctx, "deleted a device management script")
}

func (r *DeviceManagementScriptResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resourceImportStatePassthroughID(ctx, "device_management_script", req, resp)
}

func deviceManagementScriptForState(dms models.DeviceManagementScriptable) *deviceManagementScript {
	return &deviceManagementScript{
		ID:              types.StringValue(dms.GetId()),
		Name:            types.StringValue(dms.GetName()),
		Description:     types.StringValue(dms.GetDescription()),
		ScriptContent:   types.StringValue(dms.GetScriptContent()),
		RunAsAccount:    types.StringValue(dms.GetRunAsAccount()),
	}
}

type deviceManagementScript struct {
	ID              types.String `tfsdk:"id"`
	Name            types.String `tfsdk:"name"`
	Description     types.String `tfsdk:"description"`
	ScriptContent   types.String `tfsdk:"script_content"`
	RunAsAccount    types.String `tfsdk:"run_as_account"`
}
