package graphBetaWinGetApp

import (
	"context"
	"regexp"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common"
	planmodifiers "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/plan_modifiers"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	ResourceName  = "graph_beta_device_and_app_management_win_get_app"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	// Basic resource interface (CRUD operations)
	_ resource.Resource = &WinGetAppResource{}

	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &WinGetAppResource{}

	// Enables import functionality
	_ resource.ResourceWithImportState = &WinGetAppResource{}

	// Enables plan modification/diff suppression
	_ resource.ResourceWithModifyPlan = &WinGetAppResource{}
)

func NewWinGetAppResource() resource.Resource {
	return &WinGetAppResource{
		ReadPermissions: []string{
			"DeviceManagementApps.Read.All",
			"DeviceManagementConfiguration.Read.All",
		},
		WritePermissions: []string{
			"DeviceManagementApps.ReadWrite.All",
			"DeviceManagementConfiguration.ReadWrite.All",
		},
		ResourcePath: "/deviceAppManagement/mobileApps",
	}
}

type WinGetAppResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ProviderTypeName string
	TypeName         string
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

// Metadata returns the resource type name.
func (r *WinGetAppResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	r.ProviderTypeName = req.ProviderTypeName
	r.TypeName = ResourceName
	resp.TypeName = r.FullTypeName()
}

// FullTypeName returns the full resource type name in the format "providername_resourcename".
func (r *WinGetAppResource) FullTypeName() string {
	return r.ProviderTypeName + "_" + r.TypeName
}

// Configure sets the client for the resource.
func (r *WinGetAppResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = common.SetGraphBetaClientForResource(ctx, req, resp, constants.PROVIDER_NAME+"_"+ResourceName)
}

// ImportState imports the resource state.
func (r *WinGetAppResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Schema returns the schema for the resource.
func (r *WinGetAppResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages WinGet applications from the Microsoft Store using the `/deviceAppManagement/mobileApps` endpoint. WinGet apps enable deployment of Microsoft Store applications with automatic metadata population, streamlined package management, and integration with the Windows Package Manager for efficient app distribution to managed devices.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
				MarkdownDescription: "The unique identifier for this Intune Microsoft Store app",
			},
			"display_name": schema.StringAttribute{
				Computed: true,
				Optional: true,
				MarkdownDescription: "The title of the WinGet app imported from the Microsoft Store for Business." +
					"This field value must match the expected title of the app in the Microsoft Store for Business associated with the `package_identifier`." +
					"This field is automatically populated based on the package identifier when `automatically_generate_metadata` is set to true.",
			},
			"description": schema.StringAttribute{
				Computed: true,
				Optional: true,
				MarkdownDescription: "A detailed description of the WinGet/ Microsoft Store for Business app." +
					"This field is automatically populated based on the package identifier when `automatically_generate_metadata` is set to true.",
			},
			"publisher": schema.StringAttribute{
				Computed: true,
				Optional: true,
				MarkdownDescription: "The publisher of the WinGet/ Microsoft Store for Business app." +
					"This field is automatically populated based on the package identifier when `automatically_generate_metadata` is set to true.",
			},
			"categories": schema.SetAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				MarkdownDescription: "Set of category names to associate with this application. You can use either thebpredefined Intune category names like 'Business', 'Productivity', etc., or provide specific category UUIDs. Predefined values include: 'Other apps', 'Books & Reference', 'Data management', 'Productivity', 'Business', 'Development & Design', 'Photos & Media', 'Collaboration & Social', 'Computer management'.",
				Validators: []validator.Set{
					setvalidator.ValueStringsAre(
						stringvalidator.RegexMatches(
							regexp.MustCompile(`^(Other apps|Books & Reference|Data management|Productivity|Business|Development & Design|Photos & Media|Collaboration & Social|Computer management|[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12})$`),
							"must be either a predefined category name or a valid GUID in the format xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
						),
					),
				},
			},
			"package_identifier": schema.StringAttribute{
				Required: true,
				MarkdownDescription: "The **unique package identifier** for the WinGet/Microsoft Store app from the storefront.\n\n" +
					"For example, for the app Microsoft Edge Browser URL [https://apps.microsoft.com/detail/xpfftq037jwmhs?hl=en-us&gl=US](https://apps.microsoft.com/detail/xpfftq037jwmhs?hl=en-us&gl=US), " +
					"the package identifier is `xpfftq037jwmhs`.\n\n" +
					"**Important notes:**\n" +
					"- This identifier is **required** at creation time.\n" +
					"- It **cannot be modified** for existing Terraform-deployed WinGet/Microsoft Store apps.\n\n" +
					"attempting to modify this value will result in a failed request.",
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^[A-Za-z0-9]+$`),
						"package_identifier value must contain only uppercase or lowercase letters and numbers.",
					),
				},
				PlanModifiers: []planmodifier.String{
					planmodifiers.CaseInsensitiveString(),
					stringplanmodifier.RequiresReplace(),
				},
			},
			"is_featured": schema.BoolAttribute{
				Optional:            true,
				MarkdownDescription: "The value indicating whether the app is marked as featured by the admin.",
			},
			"privacy_information_url": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The privacy statement Url.",
			},
			"information_url": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The more information Url.",
			},
			"owner": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The owner of the app.",
			},
			"developer": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The developer of the app.",
			},
			"notes": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Notes for the app.",
			},
			"automatically_generate_metadata": schema.BoolAttribute{
				Required: true,
				MarkdownDescription: "When set to `true`, the provider will automatically fetch metadata from the Microsoft Store for Business " +
					"using the package identifier. This will populate the `display_name`, `description`, `publisher`, and 'icon' fields.",
			},
			"install_experience": schema.SingleNestedAttribute{
				Required:            true,
				MarkdownDescription: "The install experience settings associated with this application.the value is idempotent and any changes to this field will trigger a recreation of the application.",
				Attributes: map[string]schema.Attribute{
					"run_as_account": schema.StringAttribute{
						Required: true,
						MarkdownDescription: "The account type (System or User) that actions should be run as on target devices.  " +
							"Required at creation time.",
						Validators: []validator.String{
							stringvalidator.OneOf("system", "user"),
						},
					},
				},
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.RequiresReplace(),
				},
			},
			"large_icon": schema.SingleNestedAttribute{
				Computed: true,
				Optional: true,
				Attributes: map[string]schema.Attribute{
					"type": schema.StringAttribute{
						Computed:            true,
						Optional:            true,
						MarkdownDescription: "The MIME type of the app's large icon, automatically populated based on the `package_identifier` when `automatically_generate_metadata` is true. Example: `image/png`",
					},
					"value": schema.StringAttribute{
						Computed:            true,
						Optional:            true,
						Sensitive:           true, // not sensitive in a true sense, but we don't want to show the icon base64 encode in the plan.
						MarkdownDescription: "The icon image to use for the winget app. This field is automatically populated based on the `package_identifier` when `automatically_generate_metadata` is set to true.",
					},
				},
				MarkdownDescription: "The large icon for the WinGet app, automatically downloaded and set from the Microsoft Store for Business.",
			},
			"created_date_time": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The date and time the app was created. This property is read-only.",
			},
			"last_modified_date_time": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The date and time the app was last modified. This property is read-only.",
			},
			"upload_state": schema.Int32Attribute{
				Computed:            true,
				MarkdownDescription: "The upload state. Possible values are: 0 - Not Ready, 1 - Ready, 2 - Processing. This property is read-only.",
			},
			"publishing_state": schema.StringAttribute{
				Computed: true,
				MarkdownDescription: "The publishing state for the app. The app cannot be assigned unless the app is published. " +
					"Possible values are: notPublished, processing, published.",
			},
			"is_assigned": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "The value indicating whether the app is assigned to at least one group. This property is read-only.",
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
			"dependent_app_count": schema.Int32Attribute{
				Computed:            true,
				MarkdownDescription: "The total number of dependencies the child app has. This property is read-only.",
			},
			"superseding_app_count": schema.Int32Attribute{
				Computed:            true,
				MarkdownDescription: "The total number of apps this app directly or indirectly supersedes. This property is read-only.",
			},
			"superseded_app_count": schema.Int32Attribute{
				Computed:            true,
				MarkdownDescription: "The total number of apps this app is directly or indirectly superseded by. This property is read-only.",
			},
			"manifest_hash": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Hash of package metadata properties used to validate that the application matches the metadata in the source repository.",
			},
			//"assignments": wingetAppAssignmentsSchema(),
			"timeouts": commonschema.Timeouts(ctx),
		},
	}
}

// wingetAppAssignmentsSchema returns the nested schema for the assignments wrapper,
// which contains three sets: required, available, and uninstall.
// wingetAppAssignmentsSchema returns the schema for app assignments
func wingetAppAssignmentsSchema() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Optional: true,
		Computed: true,
		MarkdownDescription: "Configures application deployment assignments for this Microsoft Store (WinGet) application in Intune.\n\n" +
			"App assignments define which users or devices receive the application and how it's delivered. Assignments support:\n\n" +
			"* Targeting groups by entra group id, all devices, or all licensed users\n" +
			"* Assignment filters to refine targeting with dynamic criteria\n" +
			"* Installation deadlines and restart behaviors\n" +
			"* Notification controls for end-users\n\n" +
			"Each assignment must have an intent (required, available, or uninstall) which determines the application installation intent delivered to assigned devices.",
		Attributes: map[string]schema.Attribute{
			// Required assignments set
			"required": schema.SetNestedAttribute{
				Optional: true,
				Computed: true,
				MarkdownDescription: "A set of assignments with **required** intent, where the app will be automatically installed " +
					"on assigned devices without user intervention. Required installations can include deadline times and " +
					"restart settings. Commonly used for critical applications that must be deployed to specific devices or users.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: commonAssignmentAttrs(),
				},
			},
			// Available assignments set
			"available": schema.SetNestedAttribute{
				Optional: true,
				Computed: true,
				MarkdownDescription: "A set of assignments with **available** intent, where the app appears in the Company Portal " +
					"for users to install at their discretion. The app is not automatically installed, giving users control " +
					"over when or whether to install it. Ideal for optional applications or applications where user timing " +
					"preferences should be respected.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: commonAssignmentAttrs(),
				},
			},
			// Uninstall assignments set
			"uninstall": schema.SetNestedAttribute{
				Optional: true,
				Computed: true,
				MarkdownDescription: "A set of assignments with **uninstall** intent, where the app will be automatically removed " +
					"from assigned devices. Useful for retiring applications, removing unapproved software, or preparing " +
					"for a replacement application. Uninstallation can be targeted to specific groups, with appropriate " +
					"notification settings to alert users.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: commonAssignmentAttrs(),
				},
			},
		},
	}
}

// commonAssignmentAttrs defines the shared attributes for each assignment in any set.
func commonAssignmentAttrs() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			MarkdownDescription: "The ID of the app assignment associated with the Intune application.",
			Computed:            true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"source": schema.StringAttribute{
			MarkdownDescription: "The resource type which is the source for the assignment. Possible values are: direct, policySets.",
			Optional:            true,
			Computed:            true,
			Default:             stringdefault.StaticString("direct"),
		},
		"source_id": schema.StringAttribute{
			MarkdownDescription: "The identifier of the source of the assignment.",
			Optional:            true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"target": schema.SingleNestedAttribute{
			Required:            true,
			MarkdownDescription: "Target for this assignment.",
			Attributes: map[string]schema.Attribute{
				"target_type": schema.StringAttribute{
					Required: true,
					MarkdownDescription: "The target group type for the application assignment. Possible values are:\n\n" +
						"- **allDevices**: Target all devices in the tenant\n" +
						"- **allLicensedUsers**: Target all licensed users in the tenant\n" +
						"- **configurationManagerCollection**: Target System Centre Configuration Manager collection\n" +
						"- **exclusionGroupAssignment**: Target a specific Entra ID group for exclusion\n" +
						"- **groupAssignment**: Target a specific Entra ID group",
					Validators: []validator.String{
						stringvalidator.OneOf(
							"allDevices",
							"allLicensedUsers",
							"androidFotaDeployment",
							"configurationManagerCollection",
							"exclusionGroupAssignment",
							"groupAssignment",
						),
					},
				},
				"group_id": schema.StringAttribute{
					MarkdownDescription: "The entra ID group ID for the application assignment target. Required when target_type is 'groupAssignment', 'exclusionGroupAssignment', or 'androidFotaDeployment'.",
					Optional:            true,
					Validators: []validator.String{
						stringvalidator.RegexMatches(
							regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`),
							"Must be a valid GUID format (xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx)",
						),
					},
				},
				"collection_id": schema.StringAttribute{
					MarkdownDescription: "The SCCM group collection ID for the application assignment target. Default collections start with 'SMS', while custom collections start with your site code (e.g., 'MEM').",
					Optional:            true,
					Validators: []validator.String{
						stringvalidator.RegexMatches(
							regexp.MustCompile(`^[A-Za-z]{2,8}[0-9A-Za-z]{8}$`),
							"Must be a valid SCCM collection ID format. Default collections start with 'SMS' followed by an alphanumeric ID. Custom collections start with your site code (e.g., 'MEM') followed by an alphanumeric ID.",
						),
					},
				},
				"assignment_filter_id": schema.StringAttribute{
					MarkdownDescription: "The Id of the scope filter applied to the target assignment.",
					Optional:            true,
					Validators: []validator.String{
						stringvalidator.RegexMatches(
							regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`),
							"Must be a valid GUID format (xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx)",
						),
					},
				},
				"assignment_filter_type": schema.StringAttribute{
					Optional: true,
					Computed: true,
					MarkdownDescription: "The type of scope filter for the target assignment. Possible values are:\n\n" +
						"- **include**: Only include devices or users matching the filter\n" +
						"- **exclude**: Exclude devices or users matching the filter\n" +
						"- **none**: No assignment filter applied",
					Validators: []validator.String{
						stringvalidator.OneOf(
							"include",
							"exclude",
							"none",
						),
					},
				},
			},
		},
		"settings": schema.SingleNestedAttribute{
			Optional:            true,
			Computed:            true,
			MarkdownDescription: "Assignment-specific settings.",
			Attributes: map[string]schema.Attribute{
				"notifications": schema.StringAttribute{
					MarkdownDescription: "The notification settings for the assignment. The supported values are 'showAll', 'showReboot', 'hideAll'.",
					Optional:            true,
					Computed:            true,
					Default:             stringdefault.StaticString("showAll"),
					Validators: []validator.String{
						stringvalidator.OneOf(
							"showAll",
							"showReboot",
							"hideAll",
						),
					},
				},
				"install_time_settings": schema.SingleNestedAttribute{
					Optional: true,
					Attributes: map[string]schema.Attribute{
						"use_local_time": schema.BoolAttribute{
							MarkdownDescription: "Whether the local device time or UTC time should be used when determining the deadline times.",
							Optional:            true,
						},
						"deadline_date_time": schema.StringAttribute{
							MarkdownDescription: "The time at which the app should be installed.",
							Optional:            true,
						},
					},
				},
				"restart_settings": schema.SingleNestedAttribute{
					Optional: true,
					Attributes: map[string]schema.Attribute{
						"grace_period_in_minutes": schema.Int32Attribute{
							MarkdownDescription: "The number of minutes to wait before restarting the device after an app installation.",
							Optional:            true,
						},
						"countdown_display_before_restart_in_minutes": schema.Int32Attribute{
							MarkdownDescription: "The number of minutes before the restart time to display the countdown dialog for pending restarts.",
							Optional:            true,
						},
						"restart_notification_snooze_duration_in_minutes": schema.Int32Attribute{
							MarkdownDescription: "The number of minutes to snooze the restart notification dialog when the snooze button is selected.",
							Optional:            true,
						},
					},
				},
			},
		},
	}
}
