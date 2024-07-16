package deviceManagementScript

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/devicemanagement"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models"
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
			"created_date_time": schema.StringAttribute{
				Description:         "The date and time the device management script was created.",
				MarkdownDescription: "The date and time the device management script was created.",
				Computed:            true,
			},
			"last_modified_date_time": schema.StringAttribute{
				Description:         "The date and time the device management script was last modified.",
				MarkdownDescription: "The date and time the device management script was last modified.",
				Computed:            true,
			},
			"role_scope_tag_ids": schema.ListAttribute{
				Description:         "List of Scope Tag IDs for this PowerShellScript instance.",
				MarkdownDescription: "List of Scope Tag IDs for this PowerShellScript instance.",
				Optional:            true,
				ElementType:         types.StringType,
			},
			"run_as_32_bit": schema.BoolAttribute{
				Description:         "A value indicating whether the PowerShell script should run as 32-bit.",
				MarkdownDescription: "A value indicating whether the PowerShell script should run as 32-bit.",
				Optional:            true,
			},
			"enforce_signature_check": schema.BoolAttribute{
				Description:         "Indicate whether the script signature needs be checked.",
				MarkdownDescription: "Indicate whether the script signature needs be checked.",
				Optional:            true,
			},
			"file_name": schema.StringAttribute{
				Description:         "Script file name.",
				MarkdownDescription: "Script file name.",
				Optional:            true,
			},
			"assignments": schema.ListAttribute{
				Description:         "The list of group assignments for the device management script.",
				MarkdownDescription: "The list of group assignments for the device management script.",
				Optional:            true,
				ElementType: types.ObjectType{AttrTypes: map[string]schema.Attribute{
					"target_group_id": schema.StringAttribute{
						Description:         "The Id of the Azure Active Directory group we are targeting the script to.",
						MarkdownDescription: "The Id of the Azure Active Directory group we are targeting the script to.",
						Required:            true,
					},
				}},
			},
			"device_run_states": schema.ListAttribute{
				Description:         "List of run states for this script across all devices.",
				MarkdownDescription: "List of run states for this script across all devices.",
				Computed:            true,
				ElementType:         types.ObjectType{AttrTypes: map[string]schema.Attribute{
					// Add device run state attributes here
				}},
			},
			"user_run_states": schema.ListAttribute{
				Description:         "List of run states for this script across all users.",
				MarkdownDescription: "List of run states for this script across all users.",
				Computed:            true,
				ElementType:         types.ObjectType{AttrTypes: map[string]schema.Attribute{
					// Add user run state attributes here
				}},
			},
			"run_summary": schema.SingleNestedAttribute{
				Description:         "Run summary for device management script.",
				MarkdownDescription: "Run summary for device management script.",
				Computed:            true,
				Attributes:          map[string]schema.Attribute{
					// Add run summary attributes here
				},
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
	script.SetDisplayName(data.Name.ValueString())
	script.SetDescription(data.Description.ValueString())
	script.SetScriptContent([]byte(data.ScriptContent.ValueString()))
	if !data.RunAsAccount.IsNull() {
		script.SetRunAsAccount(models.RunAsAccountType(data.RunAsAccount.ValueString()))
	}
	if !data.RoleScopeTagIds.IsNull() {
		script.SetRoleScopeTagIds(expandStringList(data.RoleScopeTagIds))
	}
	if !data.RunAs32Bit.IsNull() {
		script.SetRunAs32Bit(data.RunAs32Bit.ValueBool())
	}
	if !data.EnforceSignatureCheck.IsNull() {
		script.SetEnforceSignatureCheck(data.EnforceSignatureCheck.ValueBool())
	}
	if !data.FileName.IsNull() {
		script.SetFileName(data.FileName.ValueString())
	}

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
	script.SetDisplayName(data.Name.ValueString())
	script.SetDescription(data.Description.ValueString())
	script.SetScriptContent([]byte(data.ScriptContent.ValueString()))
	if !data.RunAsAccount.IsNull() {
		script.SetRunAsAccount(models.RunAsAccountType(data.RunAsAccount.ValueString()))
	}
	if !data.RoleScopeTagIds.IsNull() {
		script.SetRoleScopeTagIds(expandStringList(data.RoleScopeTagIds))
	}
	if !data.RunAs32Bit.IsNull() {
		script.SetRunAs32Bit(data.RunAs32Bit.ValueBool())
	}
	if !data.EnforceSignatureCheck.IsNull() {
		script.SetEnforceSignatureCheck(data.EnforceSignatureCheck.ValueBool())
	}
	if !data.FileName.IsNull() {
		script.SetFileName(data.FileName.ValueString())
	}

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
		ID:                    types.StringValue(dms.GetId()),
		Name:                  types.StringValue(dms.GetDisplayName()),
		Description:           types.StringValue(dms.GetDescription()),
		ScriptContent:         types.StringValue(string(dms.GetScriptContent())),
		RunAsAccount:          types.StringValue(string(*dms.GetRunAsAccount())),
		CreatedDateTime:       types.StringValue(dms.GetCreatedDateTime().Format(time.RFC3339)),
		LastModifiedDateTime:  types.StringValue(dms.GetLastModifiedDateTime().Format(time.RFC3339)),
		RoleScopeTagIds:       flattenStringList(dms.GetRoleScopeTagIds()),
		RunAs32Bit:            types.BoolValue(dms.GetRunAs32Bit()),
		EnforceSignatureCheck: types.BoolValue(dms.GetEnforceSignatureCheck()),
		FileName:              types.StringValue(dms.GetFileName()),
		Assignments:           flattenAssignmentList(dms.GetAssignments()),
		DeviceRunStates:       flattenDeviceRunStateList(dms.GetDeviceRunStates()),
		UserRunStates:         flattenUserRunStateList(dms.GetUserRunStates()),
		RunSummary:            flattenRunSummary(dms.GetRunSummary()),
	}
}

func expandStringList(list types.List) []string {
	if list.IsNull() {
		return nil
	}

	strList := make([]string, len(list.Elems))
	for i, v := range list.Elems {
		strList[i] = v.ValueString()
	}

	return strList
}

func flattenStringList(list []string) types.List {
	if list == nil {
		return types.ListNull(types.StringType)
	}

	strList := make([]types.String, len(list))
	for i, v := range list {
		strList[i] = types.StringValue(v)
	}

	return types.List{
		ElemType: types.StringType,
		Elems:    strList,
	}
}

func flattenAssignmentList(assignments []models.DeviceManagementScriptAssignmentable) types.List {
	if assignments == nil {
		return types.ListNull(types.ObjectType{AttrTypes: map[string]schema.Attribute{
			"target_group_id": types.StringType,
		}})
	}

	elems := make([]types.Object, len(assignments))
	for i, assignment := range assignments {
		elems[i] = types.ObjectValue(map[string]types.Value{
			"target_group_id": types.StringValue(assignment.GetTargetGroupId()),
		})
	}

	return types.List{
		ElemType: types.ObjectType{AttrTypes: map[string]schema.Attribute{
			"target_group_id": types.StringType,
		}},
		Elems: elems,
	}
}

// Add flattenDeviceRunStateList and flattenUserRunStateList functions as needed

type deviceManagementScript struct {
	ID                    types.String `tfsdk:"id"`
	Name                  types.String `tfsdk:"name"`
	Description           types.String `tfsdk:"description"`
	ScriptContent         types.String `tfsdk:"script_content"`
	RunAsAccount          types.String `tfsdk:"run_as_account"`
	CreatedDateTime       types.String `tfsdk:"created_date_time"`
	LastModifiedDateTime  types.String `tfsdk:"last_modified_date_time"`
	RoleScopeTagIds       types.List   `tfsdk:"role_scope_tag_ids"`
	RunAs32Bit            types.Bool   `tfsdk:"run_as_32_bit"`
	EnforceSignatureCheck types.Bool   `tfsdk:"enforce_signature_check"`
	FileName              types.String `tfsdk:"file_name"`
	Assignments           types.List   `tfsdk:"assignments"`
	DeviceRunStates       types.List   `tfsdk:"device_run_states"`
	UserRunStates         types.List   `tfsdk:"user_run_states"`
	RunSummary            types.Object `tfsdk:"run_summary"`
}
