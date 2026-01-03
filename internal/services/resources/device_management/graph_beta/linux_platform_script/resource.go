package graphBetaLinuxPlatformScript

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	planmodifiers "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/plan_modifiers"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	commonschemagraphbeta "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema/graph_beta/device_management"
	customValidator "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/validate/attribute"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	ResourceName  = "microsoft365_graph_beta_device_management_linux_platform_script"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	// Basic resource interface (CRUD operations)
	_ resource.Resource = &LinuxPlatformScriptResource{}

	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &LinuxPlatformScriptResource{}

	// Enables import functionality
	_ resource.ResourceWithImportState = &LinuxPlatformScriptResource{}

	// Enables plan modification/diff suppression
	_ resource.ResourceWithModifyPlan = &LinuxPlatformScriptResource{}
)

func NewLinuxPlatformScriptResource() resource.Resource {
	return &LinuxPlatformScriptResource{
		ReadPermissions: []string{
			"DeviceManagementConfiguration.Read.All",
		},
		WritePermissions: []string{
			"DeviceManagementConfiguration.ReadWrite.All",
		},
		ResourcePath: "/deviceManagement/configurationPolicies",
	}
}

type LinuxPlatformScriptResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

// Metadata returns the resource type name.
func (r *LinuxPlatformScriptResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = ResourceName
}

// Configure sets the client for the resource.
func (r *LinuxPlatformScriptResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, ResourceName)
}

// ImportState imports the resource state.
func (r *LinuxPlatformScriptResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Schema defines the resource schema.
func (r *LinuxPlatformScriptResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages Linux platform scripts using the `/deviceManagement/configurationPolicies` endpoint. Linux platform scripts enable automated deployment and execution of shell scripts on managed Linux devices with configurable execution contexts, retry logic, and scheduled execution frequencies for system administration and maintenance tasks.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
				MarkdownDescription: "The unique identifier of the linux platform script.",
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the linux device management script.",
				Required:            true,
			},
			"description": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				PlanModifiers:       []planmodifier.String{planmodifiers.DefaultValueString("")},
				MarkdownDescription: "Optional description for the linux device management script.",
			},
			"script_content": schema.StringAttribute{
				MarkdownDescription: "The linux script content. This will be base64 encoded as part of the request.",
				Required:            true,
				Sensitive:           true,
			},
			"role_scope_tag_ids": schema.SetAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Set of scope tag IDs for this Settings Catalog template profile.",
				PlanModifiers: []planmodifier.Set{
					planmodifiers.DefaultSetValue(
						[]attr.Value{types.StringValue("0")},
					),
				},
			},
			"platforms": schema.StringAttribute{
				Computed: true,
				MarkdownDescription: "Platform type for this linux platform script." +
					"Will always be set to ['linux']," +
					"unknownFutureValue, androidEnterprise, or aosp. Defaults to none.",
				Validators: []validator.String{
					stringvalidator.OneOf(
						"linux",
					),
				},
				PlanModifiers: []planmodifier.String{planmodifiers.DefaultValueString("linux")},
			},
			"technologies": schema.ListAttribute{
				ElementType:         types.StringType,
				Computed:            true,
				MarkdownDescription: "Describes the technologies this settings catalog setting can be deployed with. Defaults to ['linuxMdm'].",
				Validators: []validator.List{
					customValidator.StringListAllowedValues(
						"linuxMdm",
					),
				},
				PlanModifiers: []planmodifier.List{
					planmodifiers.DefaultListValue([]attr.Value{types.StringValue("linuxMdm")}),
				},
			},
			"execution_context": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Execution context for the linux platform script. Can be one of: user or root. Defaults to user.",
				Validators: []validator.String{
					stringvalidator.OneOf(
						"user",
						"root",
					),
				},
				PlanModifiers: []planmodifier.String{planmodifiers.DefaultValueString("user")},
			},
			"execution_frequency": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Execution frequency for the Linux platform script. Can be one of: `15minutes`, `30minutes`, `1hour`, `2hour`, `3hour`, `6hour`, `12hour`, `1day`, or `1week`. Defaults to `15minutes`.",
				Validators: []validator.String{
					stringvalidator.OneOf(
						"15minutes",
						"30minutes",
						"1hour",
						"2hour",
						"3hour",
						"6hour",
						"12hour",
						"1day",
						"1week",
					),
				},
				PlanModifiers: []planmodifier.String{planmodifiers.DefaultValueString("15minutes")},
			},
			"execution_retries": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Number of times the Linux platform script should be retried on failure. Can be one of: `1`, `2`, or `3`. Defaults to `1`.",
				Validators: []validator.String{
					stringvalidator.OneOf(
						"1",
						"2",
						"3",
					),
				},
				PlanModifiers: []planmodifier.String{planmodifiers.DefaultValueString("1")},
			},
			"assignments": commonschemagraphbeta.DeviceConfigurationWithAllGroupAssignmentsAndFilterSchema(),
			"timeouts":    commonschema.ResourceTimeouts(ctx),
		},
	}
}
