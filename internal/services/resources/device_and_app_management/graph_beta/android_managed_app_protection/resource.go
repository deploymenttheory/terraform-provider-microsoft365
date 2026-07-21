package graphBetaDeviceAndAppManagementAndroidManagedAppProtection

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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	ResourceName  = "microsoft365_graph_beta_device_and_app_management_android_managed_app_protection"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	_ resource.Resource                = &AndroidManagedAppProtectionResource{}
	_ resource.ResourceWithConfigure   = &AndroidManagedAppProtectionResource{}
	_ resource.ResourceWithImportState = &AndroidManagedAppProtectionResource{}
	_ resource.ResourceWithModifyPlan  = &AndroidManagedAppProtectionResource{}
	_ resource.ResourceWithIdentity    = &AndroidManagedAppProtectionResource{}
)

func NewAndroidManagedAppProtectionResource() resource.Resource {
	return &AndroidManagedAppProtectionResource{
		ReadPermissions: []string{
			"DeviceManagementApps.Read.All",
		},
		WritePermissions: []string{
			"DeviceManagementApps.ReadWrite.All",
		},
		ResourcePath: "/deviceAppManagement/androidManagedAppProtections",
	}
}

type AndroidManagedAppProtectionResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

func (r *AndroidManagedAppProtectionResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = ResourceName
}

func (r *AndroidManagedAppProtectionResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, ResourceName)
}

func (r *AndroidManagedAppProtectionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
}

func (r *AndroidManagedAppProtectionResource) IdentitySchema(ctx context.Context, req resource.IdentitySchemaRequest, resp *resource.IdentitySchemaResponse) {
	resp.IdentitySchema = identityschema.Schema{
		Attributes: map[string]identityschema.Attribute{
			"id": identityschema.StringAttribute{
				RequiredForImport: true,
			},
		},
	}
}

func (r *AndroidManagedAppProtectionResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages Android Mobile Application Management (MAM) app protection policies in Microsoft Intune. " +
			"These policies control how managed apps handle corporate data on Android devices, including data transfer restrictions, " +
			"PIN requirements, encryption, clipboard behaviour, and threat response actions. " +
			"Uses the `/beta/deviceAppManagement/androidManagedAppProtections` endpoint.",
		Attributes: map[string]schema.Attribute{

			// --- Computed-only ---
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
				MarkdownDescription: "The unique identifier for this Android managed app protection policy. Set by the API on creation.",
			},
			"created_date_time": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The date and time the policy was created. Set by the API, read-only.",
			},
			"last_modified_date_time": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The date and time the policy was last modified. Set by the API, read-only.",
			},
			"version": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Version of the entity. Set by the API, read-only.",
			},
			"is_assigned": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "Indicates if the policy is deployed to any inclusion groups. Set by the API, read-only.",
			},
			"deployed_app_count": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "Count of apps to which the current policy is deployed. Set by the API, read-only.",
			},

			// --- Required ---
			"display_name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Policy display name. Must be unique within your tenant.",
			},

			// --- Optional ---
			"description": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
				MarkdownDescription: "The policy's description.",
			},
			"period_offline_before_access_check": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("P30D"),
				MarkdownDescription: "The period after which access is checked when the device is not connected to the internet. ISO 8601 duration format. e.g. `PT720H` for 30 days.",
			},
			"period_online_before_access_check": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("PT30M"),
				MarkdownDescription: "The period after which access is checked when the device is connected to the internet. ISO 8601 duration format. e.g. `PT30M` for 30 minutes.",
			},
			"allowed_inbound_data_transfer_sources": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString("allApps"),
				MarkdownDescription: "Sources from which data is allowed to be transferred. " +
					"Possible values: `allApps`, `managedApps`, `none`.",
				Validators: []validator.String{
					stringvalidator.OneOf("allApps", "managedApps", "none"),
				},
			},
			"allowed_outbound_data_transfer_destinations": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString("allApps"),
				MarkdownDescription: "Destinations to which data is allowed to be transferred. " +
					"Possible values: `allApps`, `managedApps`, `none`.",
				Validators: []validator.String{
					stringvalidator.OneOf("allApps", "managedApps", "none"),
				},
			},
			"organizational_credentials_required": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "Indicates whether organizational credentials are required for app use.",
			},
			"allowed_outbound_clipboard_sharing_level": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString("allApps"),
				MarkdownDescription: "The level to which the clipboard may be shared between apps on the managed device. " +
					"Possible values: `allApps`, `managedAppsWithPasteIn`, `managedApps`, `blocked`.",
				Validators: []validator.String{
					stringvalidator.OneOf("allApps", "managedAppsWithPasteIn", "managedApps", "blocked"),
				},
			},
			"data_backup_blocked": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "Indicates whether the backup of a managed app's data is blocked.",
			},
			"device_compliance_required": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
				MarkdownDescription: "Indicates whether device compliance is required.",
			},
			"managed_browser_to_open_links_required": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "Indicates whether internet links should be opened in the managed browser app.",
			},
			"save_as_blocked": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "Indicates whether users may use the Save As menu item to save a copy of protected files.",
			},
			"period_offline_before_wipe_is_enforced": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("P90D"),
				MarkdownDescription: "The amount of time an app is allowed to remain disconnected from the internet before all managed data is wiped. ISO 8601 duration format. e.g. `P90D` for 90 days.",
			},
			"pin_required": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
				MarkdownDescription: "Indicates whether an app-level PIN is required.",
			},
			"maximum_pin_retries": schema.Int64Attribute{
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(5),
				MarkdownDescription: "Maximum number of incorrect PIN retry attempts before the managed app is either blocked or wiped. Valid values 1 to 65535.",
			},
			"simple_pin_blocked": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "Indicates whether simple PINs are blocked.",
			},
			"minimum_pin_length": schema.Int64Attribute{
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(4),
				MarkdownDescription: "Minimum PIN length required for an app-level PIN if PinRequired is set to true.",
			},
			"pin_character_set": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString("numeric"),
				MarkdownDescription: "Character set which may be used for an app-level PIN if PinRequired is set to true. " +
					"Possible values: `numeric`, `alphanumericAndSymbol`.",
				Validators: []validator.String{
					stringvalidator.OneOf("numeric", "alphanumericAndSymbol"),
				},
			},
			"period_before_pin_reset": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("P365D"),
				MarkdownDescription: "TimePeriod before the all-level PIN must be reset if PinRequired is set to true. ISO 8601 duration format.",
			},
			"allowed_data_storage_locations": schema.ListAttribute{
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
				MarkdownDescription: "Data storage locations where a user may store managed data. Possible values: `oneDriveForBusiness`, `sharePoint`, `box`, `localStorage`.",
			},
			"contact_sync_blocked": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "Indicates whether contacts can be synced to the user's device.",
			},
			"print_blocked": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "Indicates whether printing is allowed from managed apps.",
			},
			"fingerprint_blocked": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "Indicates whether use of the fingerprint reader is allowed in place of a PIN if PinRequired is set to true.",
			},
			"disable_app_pin_if_device_pin_is_set": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "Indicates whether use of the app PIN is required if the device PIN is set.",
			},
			"minimum_required_os_version": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
				MarkdownDescription: "Versions less than the specified version will block the managed app from accessing company data.",
			},
			"minimum_warning_os_version": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
				MarkdownDescription: "Versions less than the specified version will result in a warning message on the managed app.",
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
			"managed_browser": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString("notConfigured"),
				MarkdownDescription: "Indicates in which managed browser internet links should be opened. " +
					"Possible values: `notConfigured`, `microsoftEdge`.",
				Validators: []validator.String{
					stringvalidator.OneOf("notConfigured", "microsoftEdge"),
				},
			},
			"screen_capture_blocked": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "Indicates whether a managed user can take screen captures of managed apps.",
			},
			"disable_app_encryption_if_device_encryption_is_enabled": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "When enabled, app level encryption is disabled if device level encryption is enabled.",
			},
			"encrypt_app_data": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
				MarkdownDescription: "Indicates whether application data for managed apps should be encrypted.",
			},
			"minimum_required_patch_version": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
				MarkdownDescription: "Define the oldest required Android security patch level a user can have to gain secure access to the app.",
			},
			"minimum_warning_patch_version": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
				MarkdownDescription: "Define the oldest recommended Android security patch level a user can have for secure access to the app.",
			},
			"custom_browser_package_id": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
				MarkdownDescription: "Unique identifier of the preferred custom browser to open weblinks on Android. Requires managed_browser_to_open_links_required to be true.",
			},
			"custom_browser_display_name": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
				MarkdownDescription: "Friendly name of the preferred custom browser to open weblinks on Android. Requires managed_browser_to_open_links_required to be true.",
			},

			"timeouts": commonschema.ResourceTimeouts(ctx),
		},
	}
}
