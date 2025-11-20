package graphBetaWindowsDeviceComplianceScript

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	planmodifiers "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/plan_modifiers"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	ResourceName  = "microsoft365_graph_beta_device_management_windows_device_compliance_script"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	// Basic resource interface (CRUD operations)
	_ resource.Resource = &DeviceComplianceScriptResource{}

	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &DeviceComplianceScriptResource{}

	// Enables import functionality
	_ resource.ResourceWithImportState = &DeviceComplianceScriptResource{}

	// Enables plan modification/diff suppression
	_ resource.ResourceWithModifyPlan = &DeviceComplianceScriptResource{}
)

func NewDeviceComplianceScriptResource() resource.Resource {
	return &DeviceComplianceScriptResource{
		ReadPermissions: []string{
			"DeviceManagementConfiguration.Read.All",
		},
		WritePermissions: []string{
			"DeviceManagementConfiguration.ReadWrite.All",
		},
		ResourcePath: "/deviceManagement/deviceComplianceScripts",
	}
}

type DeviceComplianceScriptResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

// Metadata returns the resource type name.
func (r *DeviceComplianceScriptResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = ResourceName
}

// Configure sets the client for the resource.
func (r *DeviceComplianceScriptResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, ResourceName)
}

// ImportState imports the resource state.
func (r *DeviceComplianceScriptResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Schema defines the schema for the resource.
func (r *DeviceComplianceScriptResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages Windows device compliance scripts using the `/deviceManagement/deviceComplianceScripts` endpoint. Device compliance scripts enable running PowerShell detection scripts on enrolled Windows 10/11 Azure Active Directory joined devices to validate compliance requirements and provide custom compliance assessments.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
				MarkdownDescription: "Unique identifier for the device compliance script.",
			},
			"display_name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Name of the device compliance script.",
			},
			"description": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Description of the device compliance script.",
			},
			"publisher": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Name of the device compliance script publisher.",
			},
			"version": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Version of the device compliance script.",
			},
			"run_as_32_bit": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "Indicate whether PowerShell script should run as 32-bit. Default is false, which runs script in 64-bit PowerShell.",
			},
			"run_as_account": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Indicates the type of execution context. Possible values are: system, user.",
				Validators: []validator.String{
					stringvalidator.OneOf("system", "user"),
				},
			},
			"enforce_signature_check": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "Indicate whether the script signature needs be checked. Default is false, which does not check the script signature.",
			},
			"detection_script_content": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The entire content of the detection PowerShell script.",
			},
			"role_scope_tag_ids": schema.SetAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Set of scope tag IDs for this device compliance script.",
				PlanModifiers: []planmodifier.Set{
					planmodifiers.DefaultSetValue(
						[]attr.Value{types.StringValue("0")},
					),
				},
			},
			"created_date_time": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The timestamp of when the device compliance script was created. This property is read-only.",
			},
			"last_modified_date_time": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The timestamp of when the device compliance script was modified. This property is read-only.",
			},
			"timeouts": commonschema.Timeouts(ctx),
		},
	}
}
