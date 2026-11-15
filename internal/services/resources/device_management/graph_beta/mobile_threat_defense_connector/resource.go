package graphBetaMobileThreatDefenseConnector

import (
	"context"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/identityschema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
)

const (
	ResourceName  = "microsoft365_graph_beta_device_management_mobile_threat_defense_connector"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180

	// MicrosoftDefenderForEndpointPartnerId is the fixed connector ID of the Microsoft Defender for Endpoint partner.
	MicrosoftDefenderForEndpointPartnerId = "fc780465-2017-40d4-a0c5-307022471b92"
)

var (
	// Basic resource interface (CRUD operations)
	_ resource.Resource = &MobileThreatDefenseConnectorResource{}

	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &MobileThreatDefenseConnectorResource{}

	// Enables import functionality
	_ resource.ResourceWithImportState = &MobileThreatDefenseConnectorResource{}

	// Enables plan modification/diff suppression
	_ resource.ResourceWithModifyPlan = &MobileThreatDefenseConnectorResource{}

	// Enables partner-specific config validation
	_ resource.ResourceWithValidateConfig = &MobileThreatDefenseConnectorResource{}

	// Enables identity schema for list resource support
	_ resource.ResourceWithIdentity = &MobileThreatDefenseConnectorResource{}
)

// https://learn.microsoft.com/en-us/graph/api/intune-onboarding-mobilethreatdefenseconnector-create?view=graph-rest-beta
func NewMobileThreatDefenseConnectorResource() resource.Resource {
	return &MobileThreatDefenseConnectorResource{
		ReadPermissions: []string{
			"DeviceManagementServiceConfig.Read.All",
		},
		WritePermissions: []string{
			"DeviceManagementServiceConfig.ReadWrite.All",
		},
	}
}

type MobileThreatDefenseConnectorResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
}

// Metadata returns the resource type name.
func (r *MobileThreatDefenseConnectorResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = ResourceName
}

// Configure sets the client for the resource.
func (r *MobileThreatDefenseConnectorResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, ResourceName)
}

// ImportState imports the resource state.
func (r *MobileThreatDefenseConnectorResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// IdentitySchema defines the identity schema for this resource, used by list operations to uniquely identify instances
func (r *MobileThreatDefenseConnectorResource) IdentitySchema(ctx context.Context, req resource.IdentitySchemaRequest, resp *resource.IdentitySchemaResponse) {
	resp.IdentitySchema = identityschema.Schema{
		Attributes: map[string]identityschema.Attribute{
			"id": identityschema.StringAttribute{
				RequiredForImport: true,
			},
		},
	}
}

// Schema defines the schema for the resource.
func (r *MobileThreatDefenseConnectorResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages an Intune Mobile Threat Defense partner connector using the `/deviceManagement/mobileThreatDefenseConnectors` endpoint. " +
			"The connector ID is the fixed partner ID of the Mobile Threat Defense partner (for example `" + MicrosoftDefenderForEndpointPartnerId + "` for Microsoft Defender for Endpoint). " +
			"Creating the resource posts the full connector configuration for that partner, which replaces any settings already present for the partner, and destroying it deletes the connector.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The partner ID that identifies the Mobile Threat Defense connector. Each partner has a fixed ID that is the same in every tenant, " +
					"for example `" + MicrosoftDefenderForEndpointPartnerId + "` for Microsoft Defender for Endpoint. Graph rejects IDs that do not belong to a supported partner. " +
					"Changing this value forces a new connector to be created.",
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.RegexMatches(regexp.MustCompile(constants.GuidRegex), "must be a valid GUID in the format '00000000-0000-0000-0000-000000000000'"),
				},
			},
			"android_enabled": schema.BoolAttribute{
				MarkdownDescription: "When `true`, data from the Mobile Threat Defense partner is used during compliance evaluations for Android devices. Defaults to `false`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"ios_enabled": schema.BoolAttribute{
				MarkdownDescription: "When `true`, data from the Mobile Threat Defense partner is used during compliance evaluations for iOS/iPadOS devices. Defaults to `false`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"mac_enabled": schema.BoolAttribute{
				MarkdownDescription: "When `true`, data from the Mobile Threat Defense partner is used during compliance evaluations for macOS devices. Defaults to `false`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"windows_enabled": schema.BoolAttribute{
				MarkdownDescription: "When `true`, data from the Mobile Threat Defense partner is used during compliance evaluations for Windows devices. Defaults to `false`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"android_mobile_application_management_enabled": schema.BoolAttribute{
				MarkdownDescription: "When `true`, data from the Mobile Threat Defense partner can be used during Mobile Application Management (MAM) evaluations for Android devices. " +
					"Only one partner per platform may be enabled for MAM evaluation. Defaults to `false`.",
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(false),
			},
			"ios_mobile_application_management_enabled": schema.BoolAttribute{
				MarkdownDescription: "When `true`, data from the Mobile Threat Defense partner can be used during Mobile Application Management (MAM) evaluations for iOS/iPadOS devices. " +
					"Only one partner per platform may be enabled for MAM evaluation. Defaults to `false`.",
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(false),
			},
			"windows_mobile_application_management_enabled": schema.BoolAttribute{
				MarkdownDescription: "When `true`, data from the Mobile Threat Defense partner can be used during Mobile Application Management (MAM) evaluations for Windows devices. " +
					"Only one partner per platform may be enabled for MAM evaluation. Defaults to `false`.",
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(false),
			},
			"partner_unsupported_os_version_blocked": schema.BoolAttribute{
				MarkdownDescription: "When `true`, Intune marks devices noncompliant on enabled platforms that do not meet the minimum version requirements of the Mobile Threat Defense partner. Defaults to `false`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"partner_unresponsiveness_threshold_in_days": schema.Int32Attribute{
				MarkdownDescription: "The number of days without a heartbeat from the Mobile Threat Defense partner before the partner is marked as unresponsive. " +
					"Intune then ignores the partner's data for the next compliance calculation. Defaults to `7`.",
				Optional: true,
				Computed: true,
				Default:  int32default.StaticInt32(7),
				Validators: []validator.Int32{
					int32validator.AtLeast(0),
				},
			},
			"allow_partner_to_collect_ios_application_metadata": schema.BoolAttribute{
				MarkdownDescription: "When `true`, the Mobile Threat Defense partner may collect metadata about installed applications from Intune for iOS/iPadOS devices. Defaults to `false`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"allow_partner_to_collect_ios_personal_application_metadata": schema.BoolAttribute{
				MarkdownDescription: "When `true`, the Mobile Threat Defense partner may collect metadata about personally installed applications from Intune for iOS/iPadOS devices. Defaults to `false`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"allow_partner_to_collect_ios_certificate_metadata": schema.BoolAttribute{
				MarkdownDescription: "When `true`, the Mobile Threat Defense partner may request the list of installed certificates on iOS/iPadOS devices from Intune, " +
					"including unmanaged certificates. Defaults to `false`.",
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(false),
			},
			"allow_partner_to_collect_ios_personal_certificate_metadata": schema.BoolAttribute{
				MarkdownDescription: "When `true`, the Mobile Threat Defense partner may request the list of installed certificates on personally owned iOS/iPadOS devices from Intune. Defaults to `false`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"microsoft_defender_for_endpoint_attach_enabled": schema.BoolAttribute{
				MarkdownDescription: "When `true`, configuration profile management through Microsoft Defender for Endpoint is enabled. " +
					"Graph only honours this setting for the Microsoft Defender for Endpoint connector and returns `false` for other partners. Defaults to `false`.",
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(false),
			},
			"grant_mobile_threat_defense_partner_role": schema.BoolAttribute{
				MarkdownDescription: "When `true`, the Mobile Threat Defense partner is granted its Intune role. Defaults to `false`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"launch_mobile_threat_defense_partner_on_setup_enabled": schema.BoolAttribute{
				MarkdownDescription: "When `true`, the Mobile Threat Defense partner app is launched during device setup. " +
					"Graph ignores this setting for partners that do not support it and returns `false`. Defaults to `false`.",
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(false),
			},
			"android_device_blocked_on_missing_partner_data": schema.BoolAttribute{
				MarkdownDescription: "When `true`, Intune must receive data from the partner before marking an Android device compliant. Graph derives this value from `android_enabled` and ignores writes to it.",
				Computed:            true,
			},
			"ios_device_blocked_on_missing_partner_data": schema.BoolAttribute{
				MarkdownDescription: "When `true`, Intune must receive data from the partner before marking an iOS/iPadOS device compliant. Graph derives this value from `ios_enabled` and ignores writes to it.",
				Computed:            true,
			},
			"mac_device_blocked_on_missing_partner_data": schema.BoolAttribute{
				MarkdownDescription: "When `true`, Intune must receive data from the partner before marking a macOS device compliant. Graph derives this value from `mac_enabled` and ignores writes to it.",
				Computed:            true,
			},
			"windows_device_blocked_on_missing_partner_data": schema.BoolAttribute{
				MarkdownDescription: "When `true`, Intune must receive data from the partner before marking a Windows device compliant. Graph derives this value from `windows_enabled` and ignores writes to it.",
				Computed:            true,
			},
			"partner_state": schema.StringAttribute{
				MarkdownDescription: "The partner state for this tenant, set by the service. Possible values are: `unavailable`, `available`, `enabled`, `unresponsive`, `notSetUp`, `error`, `neverConfigured`, `unknownFutureValue`.",
				Computed:            true,
			},
			"last_heartbeat_date_time": schema.StringAttribute{
				MarkdownDescription: "The date and time of the last heartbeat received from the Mobile Threat Defense partner. Null when the partner has never sent a heartbeat.",
				Computed:            true,
			},
			"partner_display_name": schema.StringAttribute{
				MarkdownDescription: "The display name of the Mobile Threat Defense partner, when the service reports one.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"partner_supports_android_mam": schema.BoolAttribute{
				MarkdownDescription: "Whether the partner supports Mobile Application Management (MAM) for Android.",
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"partner_supports_ios_mam": schema.BoolAttribute{
				MarkdownDescription: "Whether the partner supports Mobile Application Management (MAM) for iOS/iPadOS.",
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"partner_supports_windows_mam": schema.BoolAttribute{
				MarkdownDescription: "Whether the partner supports Mobile Application Management (MAM) for Windows.",
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"partner_supports_ios_certificate_sync": schema.BoolAttribute{
				MarkdownDescription: "Whether the partner supports iOS/iPadOS certificate metadata sync.",
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"partner_minimum_supported_android_version": schema.StringAttribute{
				MarkdownDescription: "The minimum Android version the partner supports, when the service reports one.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"partner_minimum_supported_ios_version": schema.StringAttribute{
				MarkdownDescription: "The minimum iOS/iPadOS version the partner supports, when the service reports one.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"partner_minimum_supported_mac_version": schema.StringAttribute{
				MarkdownDescription: "The minimum macOS version the partner supports, when the service reports one.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"partner_minimum_supported_windows_version": schema.StringAttribute{
				MarkdownDescription: "The minimum Windows version the partner supports, when the service reports one.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"timeouts": commonschema.ResourceTimeouts(ctx),
		},
	}
}
