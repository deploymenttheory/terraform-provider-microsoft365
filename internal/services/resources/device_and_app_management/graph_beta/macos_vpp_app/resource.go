package graphBetaMacOSVppApp

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
	ResourceName  = "microsoft365_graph_beta_device_and_app_management_macos_vpp_app"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	// Basic resource interface (CRUD operations)
	_ resource.Resource = &MacOSVppAppResource{}

	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &MacOSVppAppResource{}

	// Enables import functionality
	_ resource.ResourceWithImportState = &MacOSVppAppResource{}

	// Enables plan modification/diff suppression
	_ resource.ResourceWithModifyPlan = &MacOSVppAppResource{}
)

func NewMacOSVppAppResource() resource.Resource {
	return &MacOSVppAppResource{
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

type MacOSVppAppResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

// Metadata returns the resource type name.
func (r *MacOSVppAppResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = ResourceName
}

// Configure sets the client for the resource.
func (r *MacOSVppAppResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, ResourceName)
}

// ImportState imports the resource state.
func (r *MacOSVppAppResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Schema returns the schema for the resource.
func (r *MacOSVppAppResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages macOS Volume-Purchased Program (VPP) apps in Microsoft Intune using the `/deviceAppManagement/mobileApps` endpoint. VPP apps are applications purchased through Apple's Volume Purchase Program that can be distributed to managed macOS devices through the Intune management agent.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
				MarkdownDescription: "The unique identifier of the macOS VPP app.",
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
			"relationships": schema.ListNestedAttribute{
				Optional:            true,
				MarkdownDescription: "The set of direct relationships for this app.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The unique identifier for the relationship.",
						},
						"source_display_name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The display name of the source application in the relationship.",
						},
						"source_display_version": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The display version of the source application in the relationship.",
						},
						"source_id": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The unique identifier of the source application in the relationship.",
						},
						"source_publisher_display_name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The display name of the source application's publisher in the relationship.",
						},
						"target_display_name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The display name of the target application in the relationship.",
						},
						"target_display_version": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The display version of the target application in the relationship.",
						},
						"target_id": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The unique identifier of the target application in the relationship.",
						},
						"target_publisher": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The publisher of the target application in the relationship.",
						},
						"target_publisher_display_name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The display name of the target application's publisher in the relationship.",
						},
						"target_type": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The type of the target in the relationship.",
						},
					},
				},
			},
			"created_date_time": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
				MarkdownDescription: "The date and time the app was created. This property is read-only.",
			},
			"last_modified_date_time": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
				MarkdownDescription: "The date and time the app was last modified. This property is read-only.",
			},
			"upload_state": schema.Int32Attribute{
				Computed: true,
				PlanModifiers: []planmodifier.Int32{
					planmodifiers.UseStateForUnknownInt32(),
				},
				MarkdownDescription: "The upload state. Possible values are: 0 - Not Ready, 1 - Ready, 2 - Processing. This property is read-only.",
			},
			"publishing_state": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
				MarkdownDescription: "The publishing state for the app. The app cannot be assigned unless the app is published. " +
					"Possible values are: notPublished, processing, published.",
			},
			"is_assigned": schema.BoolAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.Bool{
					planmodifiers.UseStateForUnknownBool(),
				},
				MarkdownDescription: "The value indicating whether the app is assigned to at least one group. This property is read-only.",
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
			"dependent_app_count": schema.Int32Attribute{
				Computed: true,
				PlanModifiers: []planmodifier.Int32{
					planmodifiers.UseStateForUnknownInt32(),
				},
				MarkdownDescription: "The total number of dependencies the child app has. This property is read-only.",
			},
			"superseding_app_count": schema.Int32Attribute{
				Computed: true,
				PlanModifiers: []planmodifier.Int32{
					planmodifiers.UseStateForUnknownInt32(),
				},
				MarkdownDescription: "The total number of apps this app directly or indirectly supersedes. This property is read-only.",
			},
			"superseded_app_count": schema.Int32Attribute{
				Computed: true,
				PlanModifiers: []planmodifier.Int32{
					planmodifiers.UseStateForUnknownInt32(),
				},
				MarkdownDescription: "The total number of apps this app is directly or indirectly superseded by. This property is read-only.",
			},
			"used_license_count": schema.Int32Attribute{
				Computed: true,
				PlanModifiers: []planmodifier.Int32{
					planmodifiers.UseStateForUnknownInt32(),
				},
				MarkdownDescription: "The number of VPP licenses in use.",
			},
			"total_license_count": schema.Int32Attribute{
				Computed: true,
				PlanModifiers: []planmodifier.Int32{
					planmodifiers.UseStateForUnknownInt32(),
				},
				MarkdownDescription: "The total number of VPP licenses.",
			},
			"release_date_time": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
				MarkdownDescription: "The VPP application release date and time.",
			},
			"app_store_url": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
				MarkdownDescription: "The store URL.",
			},
			"licensing_type": schema.SingleNestedAttribute{
				Optional:            true,
				MarkdownDescription: "The supported License Type.",
				Attributes: map[string]schema.Attribute{
					"support_user_licensing": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Whether user licensing is supported.",
					},
					"support_device_licensing": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Whether device licensing is supported.",
					},
					"supports_user_licensing": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Whether user licensing is supported (alternative property name).",
					},
					"supports_device_licensing": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Whether device licensing is supported (alternative property name).",
					},
				},
			},
			"vpp_token_organization_name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The organization associated with the Apple Volume Purchase Program Token.",
			},
			"vpp_token_account_type": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The type of volume purchase program which the given Apple Volume Purchase Program Token is associated with. Possible values are: business, education.",
				Validators: []validator.String{
					stringvalidator.OneOf("business", "education"),
				},
			},
			"vpp_token_apple_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The Apple Id associated with the given Apple Volume Purchase Program Token.",
			},
			"bundle_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The Identity Name.",
			},
			"vpp_token_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Identifier of the VPP token associated with this app.",
			},
			"vpp_token_display_name": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
				MarkdownDescription: "Display name of the VPP token associated with this app.",
			},
			"assigned_licenses": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "The licenses assigned to this app.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"user_id": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The user identifier.",
						},
						"user_email_address": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The email address of the user.",
						},
						"user_name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The name of the user.",
						},
						"user_principal_name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The principal name of the user.",
						},
					},
				},
			},
			"revoke_license_action_results": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "Results of revoke license actions on this app.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"user_id": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The user identifier.",
						},
						"managed_device_id": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The managed device identifier.",
						},
						"total_licenses_count": schema.Int32Attribute{
							Computed:            true,
							MarkdownDescription: "The total number of licenses.",
						},
						"failed_licenses_count": schema.Int32Attribute{
							Computed:            true,
							MarkdownDescription: "The number of failed licenses.",
						},
						"action_failure_reason": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The reason for action failure.",
						},
						"action_name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The name of the action.",
						},
						"action_state": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The state of the action.",
						},
						"start_date_time": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The start date and time of the action.",
						},
						"last_updated_date_time": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The last updated date and time of the action.",
						},
					},
				},
			},
			"app_icon": commonschemagraphbeta.MobileAppIconSchema(),
			"timeouts": commonschema.Timeouts(ctx),
		},
	}
}
