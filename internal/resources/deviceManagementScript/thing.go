package deviceManagementScript

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
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
			"assignments": schema.ListAttribute{
				Description:         "Assignments for the device management script.",
				MarkdownDescription: "Assignments for the device management script.",
				Optional:            true,
				ElementType: types.ObjectType{
					AttrTypes: map[string]schema.Attribute{
						"target_group_id": schema.StringAttribute{
							Description:         "The target Azure AD group ID for the assignment.",
							MarkdownDescription: "The target Azure AD group ID for the assignment.",
							Required:            true,
						},
						"run_remediation_script": schema.BoolAttribute{
							Description:         "Whether to run the remediation script.",
							MarkdownDescription: "Whether to run the remediation script.",
							Required:            true,
						},
						"run_schedule": schema.ObjectAttribute{
							Description:         "The schedule for running the script.",
							MarkdownDescription: "The schedule for running the script.",
							Required:            true,
							AttrTypes: map[string]schema.Attribute{
								"interval": schema.IntAttribute{
									Description:         "The interval in days.",
									MarkdownDescription: "The interval in days.",
									Required:            true,
								},
								"time": schema.StringAttribute{
									Description:         "The time to run the script.",
									MarkdownDescription: "The time to run the script.",
									Required:            true,
								},
								"use_utc": schema.BoolAttribute{
									Description:         "Whether to use UTC time.",
									MarkdownDescription: "Whether to use UTC time.",
									Required:            true,
								},
							},
						},
					},
				},
			},
		},
	}
}

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
	Publisher             types.String `tfsdk:"publisher"`
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

func objectConstruction(data deviceManagementScript) (*models.DeviceHealthScript, error) {
	// Initialize a new device health script
	script := models.NewDeviceHealthScript()
	displayName := data.Name.ValueString()
	script.SetDisplayName(&displayName)

	if !data.Description.IsNull() {
		description := data.Description.ValueString()
		script.SetDescription(&description)
	}

	detectionScriptContent := []byte(data.ScriptContent.ValueString())
	script.SetDetectionScriptContent(detectionScriptContent)

	remediationScriptContent := []byte(data.ScriptContent.ValueString())
	script.SetRemediationScriptContent(remediationScriptContent)

	if !data.Publisher.IsNull() {
		publisher := data.Publisher.ValueString()
		script.SetPublisher(&publisher)
	}

	if !data.RunAsAccount.IsNull() {
		runAsAccount, err := models.ParseRunAsAccountType(data.RunAsAccount.ValueString())
		if err != nil || runAsAccount == nil {
			return nil, fmt.Errorf("invalid RunAsAccount value: got %q, should be one of %q or %q", data.RunAsAccount.ValueString(), models.SYSTEM_RUNASACCOUNTTYPE.String(), models.USER_RUNASACCOUNTTYPE.String())
		}
		script.SetRunAsAccount(runAsAccount.(*models.RunAsAccountType))
	}

	if !data.RoleScopeTagIds.IsNull() {
		script.SetRoleScopeTagIds(expandStringList(data.RoleScopeTagIds))
	}

	if !data.RunAs32Bit.IsNull() {
		runAs32Bit := data.RunAs32Bit.ValueBool()
		script.SetRunAs32Bit(&runAs32Bit)
	}

	if !data.EnforceSignatureCheck.IsNull() {
		enforceSignatureCheck := data.EnforceSignatureCheck.ValueBool()
		script.SetEnforceSignatureCheck(&enforceSignatureCheck)
	}

	return script, nil
}

func assignmentObjectConstruction(data deviceManagementScript) ([]models.DeviceHealthScriptAssignmentable, error) {
	if data.Assignments.IsNull() || len(data.Assignments.Elements()) == 0 {
		return nil, nil
	}

	var assignments []models.DeviceHealthScriptAssignmentable
	for _, assignment := range data.Assignments.Elements() {
		assignmentMap := assignment.(map[string]types.Value)
		deviceHealthScriptAssignment := models.NewDeviceHealthScriptAssignment()

		// Set target
		targetGroupID := assignmentMap["target_group_id"].(types.String)
		target := models.NewGroupAssignmentTarget()
		target.SetGroupId(targetGroupID.ValueString())
		deviceHealthScriptAssignment.SetTarget(target)

		// Set runRemediationScript
		runRemediationScript := assignmentMap["run_remediation_script"].(types.Bool)
		deviceHealthScriptAssignment.SetRunRemediationScript(&runRemediationScript.ValueBool())

		// Set runSchedule
		scheduleMap := assignmentMap["run_schedule"].(types.Object).Attrs
		runSchedule := models.NewDeviceHealthScriptDailySchedule()
		interval := scheduleMap["interval"].(types.Int).ValueInt()
		runSchedule.SetInterval(&interval)
		time := scheduleMap["time"].(types.String).ValueString()
		runSchedule.SetTime(&time)
		useUtc := scheduleMap["use_utc"].(types.Bool).ValueBool()
		runSchedule.SetUseUtc(&useUtc)
		deviceHealthScriptAssignment.SetRunSchedule(runSchedule)

		assignments = append(assignments, deviceHealthScriptAssignment)
	}

	return assignments, nil
}

func (r *DeviceManagementScriptResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data deviceManagementScript

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Construct the device health script object
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
		assignRequestBody := devicemanagement.NewAssignPostRequestBody()
		assignRequestBody.SetDeviceHealthScriptAssignments(assignments)
		_, err = r.client.ByDeviceHealthScriptId(result.GetId()).Assign().Post(ctx, assignRequestBody, nil)
		if err != nil {
			resp.Diagnostics.AddError(
				"Assignment Error",
				fmt.Sprintf("Unable to assign device management script, got error: %s", err),
			)
			return
		}
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, deviceManagementScriptForState(result))...)
}

func (r *DeviceManagementScriptResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data deviceManagementScript

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
