package graphBetaAndroidManagedDeviceAppConfigurationPolicy

import (
	"context"
	"regexp"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	planmodifiers "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/plan_modifiers"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
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
	ResourceName  = "microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	// Basic resource interface (CRUD operations)
	_ resource.Resource = &AndroidManagedDeviceAppConfigurationPolicyResource{}

	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &AndroidManagedDeviceAppConfigurationPolicyResource{}

	// Enables import functionality
	_ resource.ResourceWithImportState = &AndroidManagedDeviceAppConfigurationPolicyResource{}

	// Enables plan modification/diff suppression
	_ resource.ResourceWithModifyPlan = &AndroidManagedDeviceAppConfigurationPolicyResource{}

	// Enables resource-level configuration validation
	_ resource.ResourceWithConfigValidators = &AndroidManagedDeviceAppConfigurationPolicyResource{}
)

func NewAndroidManagedDeviceAppConfigurationPolicyResource() resource.Resource {
	return &AndroidManagedDeviceAppConfigurationPolicyResource{
		ReadPermissions: []string{
			"DeviceManagementConfiguration.Read.All",
		},
		WritePermissions: []string{
			"DeviceManagementApps.ReadWrite.All",
		},
		ResourcePath: "/deviceAppManagement/mobileAppConfigurations",
	}
}

type AndroidManagedDeviceAppConfigurationPolicyResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

// Metadata returns the resource type name.
func (r *AndroidManagedDeviceAppConfigurationPolicyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = ResourceName
}

// Configure sets the client for the resource.
func (r *AndroidManagedDeviceAppConfigurationPolicyResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, ResourceName)
}

// ImportState imports the resource state.
func (r *AndroidManagedDeviceAppConfigurationPolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Schema returns the schema for the resource.
func (r *AndroidManagedDeviceAppConfigurationPolicyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages Android managed store app configurations using the `/deviceAppManagement/mobileAppConfigurations` endpoint. This resource is used to use app configuration policies in Microsoft Intune to provide custom configuration settings for Android apps from the managed Google Play store. These configuration settings allow an app to be customized based on the app supplier's direction using Android Enterprise managed configurations. Learn more here: https://learn.microsoft.com/en-us/mem/intune/apps/app-configuration-policies-use-android.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
				MarkdownDescription: "The unique identifier for this Android mobile app configuration",
			},
			"display_name": schema.StringAttribute{
				Required:    true,
				Description: "The display name of the Android mobile app configuration",
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The optional description of the Android mobile app configuration",
			},
			"targeted_mobile_apps": schema.SetAttribute{
				ElementType:         types.StringType,
				Required:            true,
				MarkdownDescription: "Set of Android mobile app IDs that this configuration targets.",
				Validators: []validator.Set{
					setvalidator.ValueStringsAre(
						stringvalidator.RegexMatches(
							regexp.MustCompile(constants.GuidRegex),
							"Must be a valid GUID format (xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx)",
						),
					),
				},
			},
			"role_scope_tag_ids": schema.SetAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Set of scope tag IDs for this Android mobile app configuration.",
				PlanModifiers: []planmodifier.Set{
					planmodifiers.DefaultSetValue(
						[]attr.Value{types.StringValue("0")},
					),
				},
			},
			"version": schema.Int32Attribute{
				Computed:            true,
				MarkdownDescription: "Version of the Android mobile app configuration.",
			},
			"package_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The package ID of the Android app (e.g., `app:com.microsoft.office.officehubrow`).",
			},
			"payload_json": schema.StringAttribute{
				Required:            true,
				Sensitive:           true,
				MarkdownDescription: "The Android Enterprise managed configuration in Base64 encoded JSON format.",
			},
			"profile_applicability": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The profile applicability for this configuration. Possible values: `default`, `androidWorkProfile`, `androidDeviceOwner`. Defaults to `default`.",
				Validators: []validator.String{
					stringvalidator.OneOf(
						"default",
						"androidWorkProfile",
						"androidDeviceOwner",
					),
				},
			},
			"connected_apps_enabled": schema.BoolAttribute{
				Optional:            true,
				MarkdownDescription: "Whether connected apps are enabled for this configuration.",
			},
			"app_supports_oem_config": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "Whether the app supports OEM configuration. This is a computed value from the API.",
			},
			"permission_actions": schema.SetNestedAttribute{
				Optional: true,
				MarkdownDescription: "List of Android permissions and their corresponding actions.Specify permissions you want to override." +
					"If they are not chosen/specified explicitly, then the default behavior will apply. Learn more here: https://learn.microsoft.com/en-us/intune/intune-service/apps/app-configuration-policies-use-android",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"permission": schema.StringAttribute{
							Required:    true,
							Description: "The Android permission (e.g., `android.permission.CAMERA`)",
							Validators: []validator.String{
								stringvalidator.OneOf(
									"android.permission-group.NEARBY_DEVICES",
									"android.permission.NEARBY_WIFI_DEVICES",
									"android.permission.BLUETOOTH_CONNECT",
									"android.permission.READ_MEDIA_AUDIO",
									"android.permission.READ_MEDIA_IMAGES",
									"android.permission.READ_MEDIA_VIDEO",
									"android.permission.POST_NOTIFICATIONS",
									"android.permission.WRITE_EXTERNAL_STORAGE",
									"android.permission.READ_EXTERNAL_STORAGE",
									"android.permission.RECEIVE_MMS",
									"android.permission.RECEIVE_WAP_PUSH",
									"android.permission.READ_SMS",
									"android.permission.RECEIVE_SMS",
									"android.permission.SEND_SMS",
									"android.permission.BODY_SENSORS_BACKGROUND",
									"android.permission.BODY_SENSORS",
									"android.permission.PROCESS_OUTGOING_CALLS",
									"android.permission.USE_SIP",
									"android.permission.ADD_VOICEMAIL",
									"android.permission.WRITE_CALL_LOG",
									"android.permission.READ_CALL_LOG",
									"android.permission.CALL_PHONE",
									"android.permission.READ_PHONE_STATE",
									"android.permission.RECORD_AUDIO",
									"android.permission.ACCESS_BACKGROUND_LOCATION",
									"android.permission.ACCESS_COARSE_LOCATION",
									"android.permission.ACCESS_FINE_LOCATION",
									"android.permission.GET_ACCOUNTS",
									"android.permission.WRITE_CONTACTS",
									"android.permission.READ_CONTACTS",
									"android.permission.CAMERA",
									"android.permission.WRITE_CALENDAR",
									"android.permission.READ_CALENDAR",
								),
							},
						},
						"action": schema.StringAttribute{
							Required:    true,
							Description: "The action for this permission. Possible values: `prompt`, `autoGrant`, `autoDeny`",
							Validators: []validator.String{
								stringvalidator.OneOf(
									"prompt",
									"autoGrant",
									"autoDeny",
								),
							},
						},
					},
				},
			},
			"timeouts": commonschema.ResourceTimeouts(ctx),
		},
	}
}
