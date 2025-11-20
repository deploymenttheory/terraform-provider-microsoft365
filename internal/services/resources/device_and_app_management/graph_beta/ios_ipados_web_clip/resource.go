package graphBetaIOSiPadOSWebClip

import (
	"context"
	"regexp"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	planmodifiers "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/plan_modifiers"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	commonschemagraphbeta "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema/graph_beta/device_and_app_management"
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
	ResourceName  = "microsoft365_graph_beta_device_and_app_management_ios_ipados_web_clip"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	// Basic resource interface (CRUD operations)
	_ resource.Resource = &IOSiPadOSWebClipResource{}

	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &IOSiPadOSWebClipResource{}

	// Enables import functionality
	_ resource.ResourceWithImportState = &IOSiPadOSWebClipResource{}

	// Enables plan modification/diff suppression
	_ resource.ResourceWithModifyPlan = &IOSiPadOSWebClipResource{}
)

func NewIOSiPadOSWebClipResource() resource.Resource {
	return &IOSiPadOSWebClipResource{
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

type IOSiPadOSWebClipResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

// Metadata returns the resource type name.
func (r *IOSiPadOSWebClipResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = ResourceName
}

// Configure sets the client for the resource.
func (r *IOSiPadOSWebClipResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, ResourceName)
}

// ImportState imports the resource state.
func (r *IOSiPadOSWebClipResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Schema returns the schema for the resource.
func (r *IOSiPadOSWebClipResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages iOS/iPadOS Web Clip apps in Microsoft Intune using the `/deviceAppManagement/mobileApps` endpoint. iOS/iPadOS Web Clips are shortcuts to web pages that appear as icons on the home screen of iOS devices.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
				MarkdownDescription: "Key of the entity. This property is read-only.",
			},
			"display_name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The admin provided or imported title of the app.",
			},
			"description": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The description of the app.",
			},
			"publisher": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The publisher of the app.",
				Validators: []validator.String{
					stringvalidator.LengthBetween(2, 1024),
				},
			},
			"app_url": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Indicates iOS/iPadOS web clip app URL. Example: 'https://www.contoso.com'",
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(constants.HttpOrHttpsUrlRegex),
						"must be a valid URL starting with http:// or https://",
					),
				},
				PlanModifiers: []planmodifier.String{
					planmodifiers.RequiresReplaceString(),
				},
			},
			"use_managed_browser": schema.BoolAttribute{
				Optional:            true,
				MarkdownDescription: "Whether or not to use managed browser. When TRUE, the app will be required to be opened in Microsoft Edge. When FALSE, the app will not be required to be opened in Microsoft Edge. By default, this property is set to FALSE.",
			},
			"full_screen_enabled": schema.BoolAttribute{
				Optional:            true,
				MarkdownDescription: "Whether or not to open the web clip as a full-screen web app. Defaults to false. If TRUE, opens the web clip as a full-screen web app. If FALSE, the web clip opens inside of another app, such as Safari or the app specified with target_application_bundle_identifier.",
			},
			"target_application_bundle_identifier": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Specifies the application bundle identifier which opens the URL. Available in iOS 14 and later.",
			},
			"pre_composed_icon_enabled": schema.BoolAttribute{
				Optional:            true,
				MarkdownDescription: "Whether or not the icon for the app is precomosed. Defaults to false. If TRUE, prevents SpringBoard from adding \"shine\" to the icon. If FALSE, SpringBoard can add \"shine\".",
			},
			"ignore_manifest_scope": schema.BoolAttribute{
				Optional:            true,
				MarkdownDescription: "Whether or not a full screen web clip can navigate to an external web site without showing the Safari UI. Defaults to false. If FALSE, the Safari UI appears when navigating away. If TRUE, the Safari UI will not be shown.",
			},
			"is_featured": schema.BoolAttribute{
				Optional:            true,
				MarkdownDescription: "The value indicating whether the app is marked as featured by the admin.",
			},
			"privacy_information_url": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The privacy statement Url.",
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(constants.HttpOrHttpsUrlRegex),
						"must be a valid URL starting with http:// or https://",
					),
				},
			},
			"information_url": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The more information Url.",
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(constants.HttpOrHttpsUrlRegex),
						"must be a valid URL starting with http:// or https://",
					),
				},
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
			"role_scope_tag_ids": schema.SetAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "List of scope tag ids for this mobile app.",
				PlanModifiers: []planmodifier.Set{
					planmodifiers.DefaultSetValue(
						[]attr.Value{types.StringValue("0")},
					),
				},
			},
			"categories": schema.SetAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				MarkdownDescription: "The list of categories for this app. You can use either the predefined Intune category names like 'Business', 'Productivity', etc., or provide specific category UUIDs.",
				Validators: []validator.Set{
					setvalidator.ValueStringsAre(
						stringvalidator.RegexMatches(
							regexp.MustCompile(`^(Other apps|Books & Reference|Data management|Productivity|Business|Development & Design|Photos & Media|Collaboration & Social|Computer management|[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12})$`),
							"must be either a predefined category name or a valid GUID in the format xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
						),
					),
				},
			},
			"created_date_time": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The date and time the app was created. This property is read-only.",
			},
			"last_modified_date_time": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The date and time the app was last modified. This property is read-only.",
			},
			"publishing_state": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The publishing state for the app. The app cannot be assigned unless the app is published. This property is read-only. Possible values are: notPublished, processing, published.",
			},
			"is_assigned": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "The value indicating whether the app is assigned to at least one group. This property is read-only.",
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
			"upload_state": schema.Int32Attribute{
				Computed:            true,
				MarkdownDescription: "The upload state. Possible values are: 0 - Not Ready, 1 - Ready, 2 - Processing. This property is read-only.",
			},
			"app_icon": commonschemagraphbeta.MobileAppIconSchema(),
			"timeouts": commonschema.Timeouts(ctx),
		},
	}
}
