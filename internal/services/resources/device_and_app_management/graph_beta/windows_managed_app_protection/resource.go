package graphBetaDeviceAndAppManagementWindowsManagedAppProtection

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	planmodifiers "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/plan_modifiers"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/identityschema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	ResourceName  = "microsoft365_graph_beta_device_and_app_management_windows_managed_app_protection"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	_ resource.Resource                = &WindowsManagedAppProtectionResource{}
	_ resource.ResourceWithConfigure   = &WindowsManagedAppProtectionResource{}
	_ resource.ResourceWithImportState = &WindowsManagedAppProtectionResource{}
	_ resource.ResourceWithModifyPlan  = &WindowsManagedAppProtectionResource{}
	_ resource.ResourceWithIdentity    = &WindowsManagedAppProtectionResource{}
)

func NewWindowsManagedAppProtectionResource() resource.Resource {
	return &WindowsManagedAppProtectionResource{
		ReadPermissions: []string{
			"DeviceManagementApps.Read.All",
			"DeviceManagementConfiguration.Read.All",
		},
		WritePermissions: []string{
			"DeviceManagementApps.ReadWrite.All",
			"DeviceManagementConfiguration.ReadWrite.All",
		},
		ResourcePath: "/deviceAppManagement/windowsManagedAppProtections",
	}
}

type WindowsManagedAppProtectionResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

func (r *WindowsManagedAppProtectionResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = ResourceName
}

func (r *WindowsManagedAppProtectionResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, ResourceName)
}

func (r *WindowsManagedAppProtectionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
}

func (r *WindowsManagedAppProtectionResource) IdentitySchema(ctx context.Context, req resource.IdentitySchemaRequest, resp *resource.IdentitySchemaResponse) {
	resp.IdentitySchema = identityschema.Schema{
		Attributes: map[string]identityschema.Attribute{
			"id": identityschema.StringAttribute{
				RequiredForImport: true,
			},
		},
	}
}

func (r *WindowsManagedAppProtectionResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages Windows Mobile Application Management (MAM) app protection policies in Microsoft Intune. " +
			"These policies control how managed apps handle corporate data on Windows devices, including data transfer restrictions, " +
			"clipboard behaviour, and threat response actions. Uses the `/beta/deviceAppManagement/windowsManagedAppProtections` endpoint.",
		Attributes: map[string]schema.Attribute{

			// --- Computed-only (set by API, never writable, will never cause drift) ---
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
				MarkdownDescription: "The unique identifier for this Windows managed app protection policy. Set by the API on creation.",
			},
			"created_date_time": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
				MarkdownDescription: "The date and time the policy was created. Set by the API, read-only.",
			},
			"last_modified_date_time": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
				MarkdownDescription: "The date and time the policy was last modified. Set by the API, read-only.",
			},
			"version": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
				MarkdownDescription: "Version of the entity. Set by the API, read-only.",
			},
			"is_assigned": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "When TRUE, indicates that the policy is deployed to some inclusion groups. Set by the API, read-only.",
			},
			"deployed_app_count": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "Indicates the total number of applications for which the current policy is deployed. Set by the API, read-only.",
			},

			// --- Required ---
			"display_name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Policy display name. Must be unique within your tenant.",
			},

			// --- Optional with defaults (no plan modifier needed — default covers the unknown) ---
			"description": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
				MarkdownDescription: "The policy's description.",
			},
			"role_scope_tag_ids": schema.ListAttribute{
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
				MarkdownDescription: "List of scope tag IDs for this entity instance.",
			},
			"print_blocked": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "When TRUE, printing is blocked from managed apps. When FALSE, printing is allowed. Default value is FALSE.",
			},
			"allowed_inbound_data_transfer_sources": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString("allApps"),
				MarkdownDescription: "Indicates the sources from which data is allowed to be transferred into managed apps. " +
					"Possible values: `allApps`, `none`.",
				Validators: []validator.String{
					stringvalidator.OneOf("allApps", "none"),
				},
			},
			"allowed_outbound_clipboard_sharing_level": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString("anyDestinationAnySource"),
				MarkdownDescription: "Indicates the level to which the clipboard may be shared across org and non-org resources. " +
					"Possible values: `anyDestinationAnySource`, `none`, `orgDestinationAnySource`, `orgDestinationOrgSource`.",
				Validators: []validator.String{
					stringvalidator.OneOf("anyDestinationAnySource", "none", "orgDestinationAnySource", "orgDestinationOrgSource"),
				},
			},
			"allowed_outbound_data_transfer_destinations": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString("allApps"),
				MarkdownDescription: "Indicates the destinations to which data is allowed to be transferred from managed apps. " +
					"Possible values: `allApps`, `none`.",
				Validators: []validator.String{
					stringvalidator.OneOf("allApps", "none"),
				},
			},
			"maximum_allowed_device_threat_level": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString("notConfigured"),
				MarkdownDescription: "Maximum allowed device threat level as reported by the Mobile Threat Defense app. " +
					"Possible values: `notConfigured`, `secured`, `low`, `medium`, `high`.",
				Validators: []validator.String{
					stringvalidator.OneOf("notConfigured", "secured", "low", "medium", "high"),
				},
			},
			"mobile_threat_defense_remediation_action": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString("block"),
				MarkdownDescription: "Action to take if the mobile threat defense threat threshold is not met. " +
					"Possible values: `block`, `wipe`. Note: `warn` is not supported for this property.",
				Validators: []validator.String{
					stringvalidator.OneOf("block", "wipe"),
				},
			},
			"period_offline_before_wipe_is_enforced": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString("P90D"),
				MarkdownDescription: "The amount of time an app is allowed to remain disconnected from the internet before all managed data is wiped. " +
					"ISO 8601 duration format. e.g. `P5D` for 5 days, `PT0S` to never wipe.",
			},
			"period_offline_before_access_check": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString("PT720H"),
				MarkdownDescription: "The period after which access is checked when the device is not connected to the internet. " +
					"ISO 8601 duration format. e.g. `PT5M` for 5 minutes, `PT0S` to block immediately.",
			},

			// --- Optional without defaults (plan modifier required to suppress unknown after apply) ---
			"app_action_if_unable_to_authenticate_user": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
				MarkdownDescription: "Action to take when the user is unable to check in because their authentication token is invalid. " +
					"Possible values: `block`, `wipe`, `warn`, `blockWhenSettingIsSupported`. If not set, no action is taken.",
				Validators: []validator.String{
					stringvalidator.OneOf("block", "wipe", "warn", "blockWhenSettingIsSupported"),
				},
			},
			"minimum_required_sdk_version": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
				MarkdownDescription: "Versions less than the specified version will block the managed app from accessing company data. e.g. `8.1.0`.",
			},
			"minimum_wipe_sdk_version": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
				MarkdownDescription: "Versions less than the specified version will wipe the managed app and associated company data. e.g. `8.1.0`.",
			},
			"minimum_required_os_version": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
				MarkdownDescription: "Versions less than the specified version will block the managed app from accessing company data. e.g. `10.0.19041`.",
			},
			"minimum_warning_os_version": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
				MarkdownDescription: "Versions less than the specified version will result in a warning message on the managed app. e.g. `10.0.19041`.",
			},
			"minimum_wipe_os_version": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
				MarkdownDescription: "Versions less than the specified version will wipe the managed app and associated company data. e.g. `10.0.19041`.",
			},
			"minimum_required_app_version": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
				MarkdownDescription: "Versions less than the specified version will block the managed app from accessing company data.",
			},
			"minimum_warning_app_version": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
				MarkdownDescription: "Versions less than the specified version will result in a warning message on the managed app.",
			},
			"minimum_wipe_app_version": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
				MarkdownDescription: "Versions less than the specified version will wipe the managed app and associated company data.",
			},
			"maximum_required_os_version": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
				MarkdownDescription: "Versions greater than the specified version will block the managed app from accessing company data.",
			},
			"maximum_warning_os_version": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
				MarkdownDescription: "Versions greater than the specified version will result in a warning message on the managed app.",
			},
			"maximum_wipe_os_version": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
				MarkdownDescription: "Versions greater than the specified version will wipe the managed app and associated company data.",
			},

			"timeouts": commonschema.ResourceTimeouts(ctx),
		},
	}
}
