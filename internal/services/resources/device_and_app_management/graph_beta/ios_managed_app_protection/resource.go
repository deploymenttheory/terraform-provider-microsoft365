package graphBetaDeviceAndAppManagementIosManagedAppProtection

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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	ResourceName  = "microsoft365_graph_beta_device_and_app_management_ios_managed_app_protection"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	_ resource.Resource                = &IosManagedAppProtectionResource{}
	_ resource.ResourceWithConfigure   = &IosManagedAppProtectionResource{}
	_ resource.ResourceWithImportState = &IosManagedAppProtectionResource{}
	_ resource.ResourceWithModifyPlan  = &IosManagedAppProtectionResource{}
	_ resource.ResourceWithIdentity    = &IosManagedAppProtectionResource{}
)

func NewIosManagedAppProtectionResource() resource.Resource {
	return &IosManagedAppProtectionResource{
		ReadPermissions: []string{
			"DeviceManagementApps.Read.All",
		},
		WritePermissions: []string{
			"DeviceManagementApps.ReadWrite.All",
		},
		ResourcePath: "/deviceAppManagement/iosManagedAppProtections",
	}
}

type IosManagedAppProtectionResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

func (r *IosManagedAppProtectionResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = ResourceName
}

func (r *IosManagedAppProtectionResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, ResourceName)
}

func (r *IosManagedAppProtectionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
}

func (r *IosManagedAppProtectionResource) IdentitySchema(ctx context.Context, req resource.IdentitySchemaRequest, resp *resource.IdentitySchemaResponse) {
	resp.IdentitySchema = identityschema.Schema{
		Attributes: map[string]identityschema.Attribute{
			"id": identityschema.StringAttribute{
				RequiredForImport: true,
			},
		},
	}
}

func (r *IosManagedAppProtectionResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages iOS Mobile Application Management (MAM) app protection policies in Microsoft Intune. " +
			"These policies control how managed apps handle corporate data on iOS devices, including data transfer restrictions, " +
			"PIN/Face ID requirements, Open-In filtering, and app data encryption. " +
			"Uses the `/beta/deviceAppManagement/iosManagedAppProtections` endpoint.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
				MarkdownDescription: "The unique identifier for this iOS managed app protection policy. Set by the API on creation.",
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
			"deployed_app_count": schema.Int32Attribute{
				Computed:            true,
				MarkdownDescription: "Count of apps to which the current policy is deployed. Set by the API, read-only.",
			},
			"display_name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Policy display name. Must be unique within your tenant.",
			},
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
				MarkdownDescription: "The period after which access is checked when the device is not connected to the internet. ISO 8601 duration format.",
			},
			"period_online_before_access_check": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("PT30M"),
				MarkdownDescription: "The period after which access is checked when the device is connected to the internet. ISO 8601 duration format.",
			},
			"allowed_inbound_data_transfer_sources": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("allApps"),
				MarkdownDescription: "Sources from which data is allowed to be transferred. Possible values: `allApps`, `managedApps`, `none`.",
				Validators: []validator.String{
					stringvalidator.OneOf("allApps", "managedApps", "none"),
				},
			},
			"allowed_outbound_data_transfer_destinations": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("allApps"),
				MarkdownDescription: "Destinations to which data is allowed to be transferred. Possible values: `allApps`, `managedApps`, `none`.",
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
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("allApps"),
				MarkdownDescription: "The level to which the clipboard may be shared between apps. Possible values: `allApps`, `managedAppsWithPasteIn`, `managedApps`, `blocked`.",
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
				MarkdownDescription: "The amount of time an app is allowed to remain disconnected from the internet before all managed data is wiped. ISO 8601 duration format.",
			},
			"pin_required": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
				MarkdownDescription: "Indicates whether an app-level PIN is required.",
			},
			"maximum_pin_retries": schema.Int32Attribute{
				Optional:            true,
				Computed:            true,
				Default:             int32default.StaticInt32(5),
				MarkdownDescription: "Maximum number of incorrect PIN retry attempts before the managed app is blocked or wiped.",
			},
			"simple_pin_blocked": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "Indicates whether simple PINs are blocked.",
			},
			"minimum_pin_length": schema.Int32Attribute{
				Optional:            true,
				Computed:            true,
				Default:             int32default.StaticInt32(4),
				MarkdownDescription: "Minimum PIN length required for an app-level PIN if pin_required is true.",
			},
			"pin_character_set": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("numeric"),
				MarkdownDescription: "Character set which may be used for an app-level PIN. Possible values: `numeric`, `alphanumericAndSymbol`.",
				Validators: []validator.String{
					stringvalidator.OneOf("numeric", "alphanumericAndSymbol"),
				},
			},
			"period_before_pin_reset": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("P365D"),
				MarkdownDescription: "Time period before the app-level PIN must be reset if pin_required is true. ISO 8601 duration format.",
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
				MarkdownDescription: "Indicates whether use of Touch ID is allowed in place of a PIN if pin_required is true.",
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
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("notConfigured"),
				MarkdownDescription: "Indicates in which managed browser internet links should be opened. Possible values: `notConfigured`, `microsoftEdge`.",
				Validators: []validator.String{
					stringvalidator.OneOf("notConfigured", "microsoftEdge"),
				},
			},
			"face_id_blocked": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "Indicates whether use of Face ID is blocked as an alternative to an app PIN.",
			},
			"custom_browser_protocol": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
				MarkdownDescription: "Custom browser protocol used to open weblinks on iOS. Requires managed_browser_to_open_links_required to be true.",
			},
			"third_party_keyboards_blocked": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "Indicates whether third-party keyboards are blocked while accessing a managed app.",
			},
			"filter_open_in_to_only_managed_apps": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "Indicates whether Open-In (file sharing) is restricted to only other managed apps.",
			},
			"disable_protection_of_managed_outbound_open_in_data": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "Disables protection of managed data sent to other apps via iOS's Open-In feature.",
			},
			"app_data_encryption_type": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString("whenDeviceLocked"),
				MarkdownDescription: "Type of encryption applied to app data. Possible values: `useDeviceSettings`, `afterDeviceRestart`, " +
					"`whenDeviceLockedExceptOpenFiles`, `whenDeviceLocked`. Verify these against your SDK's generated enum if the build fails here.",
				Validators: []validator.String{
					stringvalidator.OneOf("useDeviceSettings", "afterDeviceRestart", "whenDeviceLockedExceptOpenFiles", "whenDeviceLocked"),
				},
			},
			"timeouts": commonschema.ResourceTimeouts(ctx),
		},
	}
}