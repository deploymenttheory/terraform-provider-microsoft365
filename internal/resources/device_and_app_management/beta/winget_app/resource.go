package graphBetaWinGetApp

import (
	"context"
	"regexp"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common"
	planmodifiers "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/plan_modifiers"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/schema"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
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
	resp.TypeName = req.ProviderTypeName + "_" + ResourceName
}

// Configure sets the client for the resource.
func (r *WinGetAppResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = common.SetGraphBetaClientForResource(ctx, req, resp, r.TypeName)
}

// ImportState imports the resource state.
func (r *WinGetAppResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Schema returns the schema for the resource.
func (r *WinGetAppResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {

	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages an Intune Microsoft Store app (new) resource aka winget, using the mobileapps graph beta API.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				MarkdownDescription: "The unique graph guid that identifies this resource." +
					"Assigned at time of resource creation. This property is read-only.",
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
			"install_experience": schema.SingleNestedAttribute{
				Required: true,
				Attributes: map[string]schema.Attribute{
					"run_as_account": schema.StringAttribute{
						Required: true,
						MarkdownDescription: "The account type (System or User) that actions should be run as on target devices. " +
							"Required at creation time.",
						Validators: []validator.String{
							stringvalidator.OneOf("system", "user"),
						},
					},
				},
				MarkdownDescription: "The install experience settings associated with this application.",
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
			"upload_state": schema.Int64Attribute{
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
			"role_scope_tag_ids": schema.ListAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				MarkdownDescription: "List of scope tag ids for this mobile app.",
			},
			"dependent_app_count": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "The total number of dependencies the child app has. This property is read-only.",
			},
			"superseding_app_count": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "The total number of apps this app directly or indirectly supersedes. This property is read-only.",
			},
			"superseded_app_count": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "The total number of apps this app is directly or indirectly superseded by. This property is read-only.",
			},
			"manifest_hash": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Hash of package metadata properties used to validate that the application matches the metadata in the source repository.",
			},
			"assignments": WinGetAppAssignmentsSchema(),
			"timeouts":    commonschema.Timeouts(ctx),
		},
	}
}

func WinGetAppAssignmentsSchema() schema.ListNestedAttribute {
	return schema.ListNestedAttribute{
		Optional:            true,
		MarkdownDescription: "Manages the assignments for Intune Microsoft Store app (new) resource aka winget, using the mobileapps graph beta API.",
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"id": schema.StringAttribute{
					MarkdownDescription: "The ID of the winget app associated with this assignment.",
					Computed:            true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.UseStateForUnknown(),
					},
				},
				"intent": schema.StringAttribute{
							MarkdownDescription: "The install intent defined by the admin. Possible values are: available, required, uninstall, availableWithoutEnrollment.",
							Required:            true,
						},
				"source": schema.StringAttribute{
							MarkdownDescription: "The resource type which is the source for the assignment. Possible values are: direct, policySets. This property is read-only.",
							Required:            true,
						},
				"source_id": schema.StringAttribute{
							MarkdownDescription: "The identifier of the source of the assignment. This property is read-only.",
							Computed:            true,
						},
				"target": schema.SingleNestedAttribute{
					MarkdownDescription: "The target group assignment defined by the admin.",
					Required:            true,
					Attributes: map[string]schema.Attribute{
						"target_type": schema.StringAttribute{
							MarkdownDescription: "The type of target. Possible values: groupAssignmentTarget, allLicensedUsersAssignmentTarget, allDevicesAssignmentTarget, allLicensedUsers, allDevices",
							Required:            true,
						},
						"group_id": schema.StringAttribute{
							MarkdownDescription: "The ID of the target group.",
							Optional:            true,
						},
						"device_and_app_management_assignment_filter_id": schema.StringAttribute{
							MarkdownDescription: "The ID of the filter for the target assignment.",
							Optional:            true,
						},
						"device_and_app_management_assignment_filter_type": schema.StringAttribute{
							MarkdownDescription: "The type of filter for the target assignment. Possible values: include, exclude, none",
							Optional:            true,
						},
						"is_exclusion_group": schema.BoolAttribute{
							MarkdownDescription: "Indicates whether this is an exclusion group.",
							Optional:            true,
						},
					},
				},
				"settings": schema.SingleNestedAttribute{
					MarkdownDescription: "The settings for target assignment defined by the admin.",
					Optional:            true,
					Attributes: map[string]schema.Attribute{
						"notifications": schema.StringAttribute{
							MarkdownDescription: "The notification settings for the assignment. Possible values: showAll, showReboot, hideAll",
							Optional:            true,
						},
						"install_time_settings": schema.SingleNestedAttribute{
							MarkdownDescription: "Settings related to install time.",
							Optional:            true,
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
							MarkdownDescription: "Settings related to restarts after installation.",
							Optional:            true,
							Attributes: map[string]schema.Attribute{
								"grace_period_in_minutes": schema.Int64Attribute{
											MarkdownDescription: "The number of minutes to wait before restarting the device after an app installation.",
											Optional:            true,
										},
										"countdown_display_before_restart_in_minutes": schema.Int64Attribute{
											MarkdownDescription: "The number of minutes before the restart time to display the countdown dialog for pending restarts.",
											Optional:            true,
										},
										"restart_notification_snooze_duration_in_minutes": schema.Int64Attribute{
											MarkdownDescription: "The number of minutes to snooze the restart notification dialog when the snooze button is selected.",
											Optional:            true,
										},
							},
						},
					},
				},
			},
		},
	}
}
