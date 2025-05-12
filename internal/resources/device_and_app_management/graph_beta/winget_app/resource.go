package graphBetaWinGetApp

import (
	"context"
	"regexp"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common"
	planmodifiers "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/plan_modifiers"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/schema"
	commonschemagraphbeta "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/schema/graph_beta/device_and_app_management"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
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
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
				MarkdownDescription: "The unique identifier for this Intune Microsoft Store app",
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
			"assignments": commonschemagraphbeta.MobileAppAssignmentSchema(),
			"timeouts":    commonschema.Timeouts(ctx),
		},
	}
}
