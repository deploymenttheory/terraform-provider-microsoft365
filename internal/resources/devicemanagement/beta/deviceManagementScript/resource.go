package deviceManagementScript

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
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
	client      *devicemanagement.DeviceManagementScriptsDeviceManagementScriptItemRequestBuilder
	assignments *devicemanagement.DeviceManagementScriptsItemAssignmentsRequestBuilder
}

// deviceManagementScriptData is the model for the device management script resource.
type deviceManagementScriptData struct {
	ID                          types.String `tfsdk:"id"`
	Name                        types.String `tfsdk:"name"`
	Description                 types.String `tfsdk:"description"`
	DetectionScriptContent      types.String `tfsdk:"detection_script_content"`
	RemediationScriptContent    types.String `tfsdk:"remediation_script_content"`
	CreatedDateTime             types.String `tfsdk:"created_date_time"`
	LastModifiedDateTime        types.String `tfsdk:"last_modified_date_time"`
	RunAsAccount                types.String `tfsdk:"run_as_account"`
	EnforceSignatureCheck       types.Bool   `tfsdk:"enforce_signature_check"`
	RunAs32Bit                  types.Bool   `tfsdk:"run_as_32_bit"`
	RoleScopeTagIds             types.List   `tfsdk:"role_scope_tag_ids"`
	IsGlobalScript              types.Bool   `tfsdk:"is_global_script"`
	HighestAvailableVersion     types.String `tfsdk:"highest_available_version"`
	DeviceHealthScriptType      types.String `tfsdk:"device_health_script_type"`
	DetectionScriptParameters   types.List   `tfsdk:"detection_script_parameters"`
	RemediationScriptParameters types.List   `tfsdk:"remediation_script_parameters"`
	Assignments                 types.List   `tfsdk:"assignments"`
	DeviceRunStates             types.List   `tfsdk:"device_run_states"`
	UserRunStates               types.List   `tfsdk:"user_run_states"`
	RunSummary                  types.Object `tfsdk:"run_summary"`
	Publisher                   types.String `tfsdk:"publisher"`
}

type deviceHealthScriptParameter struct {
	Name                             types.String `tfsdk:"name"`
	Description                      types.String `tfsdk:"description"`
	IsRequired                       types.Bool   `tfsdk:"is_required"`
	ApplyDefaultValueWhenNotAssigned types.Bool   `tfsdk:"apply_default_value_when_not_assigned"`
	DefaultValue                     types.String `tfsdk:"default_value"`
}

// deviceManagementScriptAssignmentData is the model for the device management script assignment resource.
type deviceManagementScriptAssignmentData struct {
	TargetGroupID        types.String `tfsdk:"target_group_id"`
	RunRemediationScript types.Bool   `tfsdk:"run_remediation_script"`
	RunSchedule          types.Object `tfsdk:"run_schedule"`
}

type runSchedule struct {
	Interval types.Int64  `tfsdk:"interval"`
	Time     types.String `tfsdk:"time"`
	UseUtc   types.Bool   `tfsdk:"use_utc"`
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
				Description:         "Description of the device health script.",
				MarkdownDescription: "Description of the device health script.",
				Optional:            true,
			},
			"detection_script_content": schema.StringAttribute{
				Description:         "The content of the detection PowerShell script.",
				MarkdownDescription: "The content of the detection PowerShell script.",
				Required:            true,
			},
			"remediation_script_content": schema.StringAttribute{
				Description:         "The content of the remediation PowerShell script.",
				MarkdownDescription: "The content of the remediation PowerShell script.",
				Required:            true,
			},
			"run_as_account": schema.StringAttribute{
				Description:         "Indicates the type of execution context the windows script runs in. Valid values are 'system' or 'user'.",
				MarkdownDescription: "Indicates the type of execution context the windows script runs in. Valid values are `system` or `user`.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("system", "user"),
				},
			},
			"run_as_32_bit": schema.BoolAttribute{
				Description:         "Indicate whether PowerShell script(s) should run as 32-bit.",
				MarkdownDescription: "Indicate whether PowerShell script(s) should run as 32-bit.",
				Optional:            true,
			},
			"enforce_signature_check": schema.BoolAttribute{
				Description:         "Indicate whether the script signature needs be checked.",
				MarkdownDescription: "Indicate whether the script signature needs be checked.",
				Optional:            true,
			},
			"role_scope_tag_ids": schema.ListAttribute{
				Description:         "List of Scope Tag IDs for the device health script.",
				MarkdownDescription: "List of Scope Tag IDs for the device health script.",
				Optional:            true,
				ElementType:         types.StringType,
			},
			"is_global_script": schema.BoolAttribute{
				Description:         "Indicates if the script is a global script.",
				MarkdownDescription: "Indicates if the script is a global script.",
				Optional:            true,
			},
			"highest_available_version": schema.StringAttribute{
				Description:         "The highest available version of the script.",
				MarkdownDescription: "The highest available version of the script.",
				Optional:            true,
			},
			"device_health_script_type": schema.StringAttribute{
				Description:         "The type of device health script.",
				MarkdownDescription: "The type of device health script.",
				Optional:            true,
			},
			"detection_script_parameters": schema.ListAttribute{
				Description:         "The parameters for the detection script.",
				MarkdownDescription: "The parameters for the detection script.",
				Optional:            true,
				ElementType: types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"name":                                  types.StringType,
						"description":                           types.StringType,
						"is_required":                           types.BoolType,
						"apply_default_value_when_not_assigned": types.BoolType,
						"default_value":                         types.StringType,
					},
				},
			},
			"remediation_script_parameters": schema.ListAttribute{
				Description:         "The parameters for the remediation script.",
				MarkdownDescription: "The parameters for the remediation script.",
				Optional:            true,
				ElementType: types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"name":                                  types.StringType,
						"description":                           types.StringType,
						"is_required":                           types.BoolType,
						"apply_default_value_when_not_assigned": types.BoolType,
						"default_value":                         types.StringType,
					},
				},
			},
			"assignments": schema.ListAttribute{
				Description:         "Assignments for the device management script.",
				MarkdownDescription: "Assignments for the device management script.",
				Optional:            true,
				ElementType: types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"target_group_id":        types.StringType,
						"run_remediation_script": types.BoolType,
						"run_schedule": types.ObjectType{
							AttrTypes: map[string]attr.Type{
								"interval": types.Int64Type,
								"time":     types.StringType,
								"use_utc":  types.BoolType,
							},
						},
					},
				},
			},
		},
	}
}

func (r *DeviceManagementScriptResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data deviceManagementScriptData

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Construct the device management script object
	script, err := objectConstruction(data)
	if err != nil {
		resp.Diagnostics.AddError(
			"Object Construction Error",
			fmt.Sprintf("Unable to construct device management script object, got error: %s", err),
		)
		return
	}

	// Instantiate the DeviceManagementScriptsRequestBuilder
	scriptsRequestBuilder := devicemanagement.NewDeviceManagementScriptsRequestBuilder(r.client.BaseRequestBuilder.UrlTemplate, r.client.BaseRequestBuilder.RequestAdapter)

	// Make the API call to create the device management script
	result, err := scriptsRequestBuilder.Post(ctx, script, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to create device management script, got error: %s", err),
		)
		return
	}

	tflog.Trace(ctx, "created a device management script")

	// Handle Assignments
	assignments, err := assignmentObjectConstruction(data)
	if err != nil {
		resp.Diagnostics.AddError(
			"Assignment Construction Error",
			fmt.Sprintf("Unable to construct assignment objects, got error: %s", err),
		)
		return
	}

	if assignments != nil {
		// Create AssignPostRequestBody and set assignments
		assignRequestBody := devicemanagement.NewAssignPostRequestBody()
		assignRequestBody.SetDeviceManagementScriptAssignments(assignments)

		// Make the API call to assign the device management script
		_, err = r.client.ByDeviceManagementScriptId(result.GetId()).Assign().Post(ctx, assignRequestBody, nil)
		if err != nil {
			resp.Diagnostics.AddError(
				"Assignment Error",
				fmt.Sprintf("Unable to assign device management script, got error: %s", err),
			)
			return
		}
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, resourcedeviceManagementScriptForState(result))...)
}

func (r *DeviceManagementScriptResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data deviceManagementScriptData

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	requestParameters := &devicemanagement.DeviceManagementScriptsDeviceManagementScriptItemRequestBuilderGetQueryParameters{
		Expand: []string{"assignments", "runSummary"},
	}
	requestConfiguration := &devicemanagement.DeviceManagementScriptsDeviceManagementScriptItemRequestBuilderGetRequestConfiguration{
		QueryParameters: requestParameters,
	}

	result, err := r.client.Get(ctx, requestConfiguration)
	if err != nil {
		resp.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to read device management script %s, got error: %s", data.ID.ValueString(), err),
		)
		return
	}

	tflog.Trace(ctx, "read a device management script")

	resp.Diagnostics.Append(resp.State.Set(ctx, resourcedeviceManagementScriptForState(result))...)
}

func (r *DeviceManagementScriptResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data deviceManagementScriptData

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
	resp.Diagnostics.Append(resp.State.Set(ctx, resourcedeviceManagementScriptForState(result))...)
}

func (r *DeviceManagementScriptResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data deviceManagementScriptData

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

func resourcedeviceManagementScriptForState(dms models.DeviceManagementScriptable) *deviceManagementScriptData {
	return &deviceManagementScriptData{
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
