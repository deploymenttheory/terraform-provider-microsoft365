package graphBetaAppControlForBusinessBuiltInControls

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	planmodifiers "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/plan_modifiers"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	commonschemagraphbeta "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema/graph_beta/device_management"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
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
	ResourceName  = "microsoft365_graph_beta_device_management_app_control_for_business_built_in_controls"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	// Basic resource interface (CRUD operations)
	_ resource.Resource = &AppControlForBusinessResourceBuiltInControls{}

	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &AppControlForBusinessResourceBuiltInControls{}

	// Enables import functionality
	_ resource.ResourceWithImportState = &AppControlForBusinessResourceBuiltInControls{}

	// Enables plan modification/diff suppression
	_ resource.ResourceWithModifyPlan = &AppControlForBusinessResourceBuiltInControls{}
)

func NewAppControlForBusinessResourceBuiltInControlsResource() resource.Resource {
	return &AppControlForBusinessResourceBuiltInControls{
		ReadPermissions: []string{
			"DeviceManagementConfiguration.Read.All",
		},
		WritePermissions: []string{
			"DeviceManagementConfiguration.ReadWrite.All",
		},
		ResourcePath: "/deviceManagement/configurationPolicies",
	}
}

type AppControlForBusinessResourceBuiltInControls struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

// Metadata returns the resource type name.
func (r *AppControlForBusinessResourceBuiltInControls) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = ResourceName
}

// Configure sets the client for the resource.
func (r *AppControlForBusinessResourceBuiltInControls) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, ResourceName)
}

// ImportState imports the resource state.
func (r *AppControlForBusinessResourceBuiltInControls) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Schema defines the resource schema.
func (r *AppControlForBusinessResourceBuiltInControls) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages App Control for Business configuration policies using the `/deviceManagement/configurationPolicies` endpoint. App Control for Business policies enable application control and trust settings on Windows devices with configurable enforcement modes and trusted application sources.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
				MarkdownDescription: "The unique identifier of the app control for business status.",
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the App Control for Business policy.",
				Required:            true,
			},
			"description": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Optional description of the resource. Maximum length is 1500 characters.",
				Validators: []validator.String{
					stringvalidator.LengthAtMost(1500),
				},
			},
			"enable_app_control": schema.StringAttribute{
				Required: true,
				MarkdownDescription: "Select Audit only to log all events in local client logs but not block any apps from running " +
					"or select Enforce to actively block apps from running in a deployed App Control for Business base policy. " +
					"App Control for Business policies created in either Audit only or Enforce mode will be deployed as rebootless base policies to all devices targeted. " +
					"By default, any devices targeted with this App Control for Business policy will have the setting to Trust Windows components and Store apps enabled, " +
					"in either audit or enforce mode based on your selection.",
				Validators: []validator.String{
					stringvalidator.OneOf(
						"enforce",
						"audit",
					),
				},
			},
			"additional_rules_for_trusting_apps": schema.SetAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
				MarkdownDescription: "By default, any devices targeted with this App Control for Business policy will have the setting to Trust Windows components and " +
					"Store apps enabled, in either audit or enforce mode based on your selection.Further, you can optionally add some additional rules to your " +
					"policy, such as selecting Trust apps with good reputation to allow reputable apps as defined by the Microsoft Intelligent Security Graph to run." +
					"Select Trust apps from managed installers to allow apps deployed via authorized sources of application deployment (managed installers). " +
					"The Intune management extension will be considered a managed installer if it has been set as such within your organization. Any apps not marked as " +
					"coming from a managed installer will not be allowed to run.All other apps and files not specified by the rules in this App Control for Business policy " +
					"will be audited only in local client logs (if Audit only is selected), or blocked (if Enforce is selected) from running on devices. " +
					"Can include: 'trust_apps_with_good_reputation' - allows reputable apps as defined by the Microsoft Intelligent Security Graph to run; " +
					"'trust_apps_from_managed_installers' - allows apps deployed via authorized sources of application deployment (managed installers).",
				Validators: []validator.Set{
					setvalidator.ValueStringsAre(
						stringvalidator.OneOf(
							"trust_apps_with_good_reputation",
							"trust_apps_from_managed_installers",
						),
					),
				},
			},
			"role_scope_tag_ids": schema.SetAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Set of scope tag IDs for this App Control for Business policy.",
				PlanModifiers: []planmodifier.Set{
					planmodifiers.DefaultSetValue(
						[]attr.Value{types.StringValue("0")},
					),
				},
			},
			"assignments": commonschemagraphbeta.DeviceConfigurationWithAllGroupAssignmentsAndFilterSchema(),
			"timeouts":    commonschema.ResourceTimeouts(ctx),
		},
	}
}
