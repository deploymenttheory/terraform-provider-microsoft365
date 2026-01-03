package graphBetaLinuxDeviceComplianceScript

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	planmodifiers "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/plan_modifiers"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	ResourceName  = "microsoft365_graph_beta_device_management_linux_device_compliance_script"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	// Basic resource interface (CRUD operations)
	_ resource.Resource = &LinuxDeviceComplianceScriptResource{}

	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &LinuxDeviceComplianceScriptResource{}

	// Enables import functionality
	_ resource.ResourceWithImportState = &LinuxDeviceComplianceScriptResource{}

	// Enables plan modification/diff suppression
	_ resource.ResourceWithModifyPlan = &LinuxDeviceComplianceScriptResource{}
)

func NewLinuxDeviceComplianceScriptResource() resource.Resource {
	return &LinuxDeviceComplianceScriptResource{
		ReadPermissions: []string{
			"DeviceManagementConfiguration.Read.All",
		},
		WritePermissions: []string{
			"DeviceManagementConfiguration.ReadWrite.All",
		},
		ResourcePath: "/deviceManagement/reusablePolicySettings",
	}
}

type LinuxDeviceComplianceScriptResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

// Metadata returns the resource type name.
func (r *LinuxDeviceComplianceScriptResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = ResourceName
}

// Configure sets the client for the resource.
func (r *LinuxDeviceComplianceScriptResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, ResourceName)
}

// ImportState imports the resource state.
func (r *LinuxDeviceComplianceScriptResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Schema defines the schema for the resource.
func (r *LinuxDeviceComplianceScriptResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages Linux device compliance scripts using the `/deviceManagement/reusablePolicySettings` endpoint. Linux device compliance scripts enable running shell scripts on enrolled Linux devices to validate compliance requirements and provide custom compliance assessments.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
				MarkdownDescription: "Unique identifier for the Linux device compliance script.",
			},
			"display_name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Name of the Linux device compliance script.",
			},
			"description": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Description of the Linux device compliance script.",
			},
			"detection_script_content": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The entire content of the detection shell script for Linux compliance checking.",
			},
			"setting_definition_id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The setting definition ID for Linux custom compliance discovery script.",
			},
			"version": schema.Int32Attribute{
				Computed:            true,
				MarkdownDescription: "Version of the Linux device compliance script.",
			},
			"last_modified_date_time": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The timestamp of when the Linux device compliance script was modified. This property is read-only.",
			},
			"timeouts": commonschema.ResourceTimeouts(ctx),
		},
	}
}
