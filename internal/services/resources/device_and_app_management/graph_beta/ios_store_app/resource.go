package graphBetaIOSStoreApp

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
	ResourceName  = "graph_beta_device_and_app_management_ios_store_app"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	// Basic resource interface (CRUD operations)
	_ resource.Resource = &IOSStoreAppResource{}

	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &IOSStoreAppResource{}

	// Enables import functionality
	_ resource.ResourceWithImportState = &IOSStoreAppResource{}

	// Enables plan modification/diff suppression
	_ resource.ResourceWithModifyPlan = &IOSStoreAppResource{}
)

func NewIOSStoreAppResource() resource.Resource {
	return &IOSStoreAppResource{
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

type IOSStoreAppResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ProviderTypeName string
	TypeName         string
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

// Metadata returns the resource type name.
func (r *IOSStoreAppResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	r.ProviderTypeName = req.ProviderTypeName
	r.TypeName = ResourceName
	resp.TypeName = r.FullTypeName()
}

// FullTypeName returns the full resource type name in the format "providername_resourcename".
func (r *IOSStoreAppResource) FullTypeName() string {
	return r.ProviderTypeName + "_" + r.TypeName
}

// Configure sets the client for the resource.
func (r *IOSStoreAppResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, constants.PROVIDER_NAME+"_"+ResourceName)
}

// ImportState imports the resource state.
func (r *IOSStoreAppResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Schema returns the schema for the resource.
func (r *IOSStoreAppResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages iOS Store apps in Microsoft Intune using the `/deviceAppManagement/mobileApps` endpoint. iOS Store apps are applications from the Apple App Store that can be managed through Intune.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
				MarkdownDescription: "The unique identifier of the iOS Store app.",
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
			"app_store_url": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The Apple AppStoreUrl.",
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^https://apps.apple.com/.*$`),
						"must be a valid Apple App Store URL",
					),
				},
				PlanModifiers: []planmodifier.String{
					planmodifiers.RequiresReplaceString(),
				},
			},
			"applicable_device_type": schema.SingleNestedAttribute{
				Required:            true,
				MarkdownDescription: "The iOS architecture for which this app can run on.",
				Attributes: map[string]schema.Attribute{
					"ipad": schema.BoolAttribute{
						Required:            true,
						MarkdownDescription: "Whether the app should run on iPads.",
					},
					"iphone_and_ipod": schema.BoolAttribute{
						Required:            true,
						MarkdownDescription: "Whether the app should run on iPhones and iPods.",
					},
				},
			},
			"minimum_supported_operating_system": schema.SingleNestedAttribute{
				Required:            true,
				MarkdownDescription: "The value for the minimum supported operating system.",
				Attributes: map[string]schema.Attribute{
					"v8_0": schema.BoolAttribute{
						Optional: true,
						Computed: true,
						MarkdownDescription: "Indicates the minimum iOS version support required for the managed device." +
							"When 'True', iOS with OS Version 8.0 or later is required to install the app. " +
							"If 'False', iOS Version 8.0 is not the minimum version. Default value is False." +
							"Exactly one of the minimum operating system boolean values will be TRUE.",
						PlanModifiers: []planmodifier.Bool{
							planmodifiers.BoolDefaultValue(false),
						},
					},
					"v9_0": schema.BoolAttribute{
						Optional: true,
						Computed: true,
						MarkdownDescription: "Indicates the minimum iOS version support required for the managed device." +
							"When 'True', iOS with OS Version 9.0 or later is required to install the app. " +
							"If 'False', iOS Version 9.0 is not the minimum version. Default value is False." +
							"Exactly one of the minimum operating system boolean values will be TRUE.",
						PlanModifiers: []planmodifier.Bool{
							planmodifiers.BoolDefaultValue(false),
						},
					},
					"v10_0": schema.BoolAttribute{
						Optional: true,
						Computed: true,
						MarkdownDescription: "Indicates the minimum iOS version support required for the managed device." +
							"When 'True', iOS with OS Version 10.0 or later is required to install the app. " +
							"If 'False', iOS Version 10.0 is not the minimum version. Default value is False." +
							"Exactly one of the minimum operating system boolean values will be TRUE.",
						PlanModifiers: []planmodifier.Bool{
							planmodifiers.BoolDefaultValue(false),
						},
					},
					"v11_0": schema.BoolAttribute{
						Optional: true,
						Computed: true,
						MarkdownDescription: "Indicates the minimum iOS version support required for the managed device." +
							"When 'True', iOS with OS Version 11.0 or later is required to install the app. " +
							"If 'False', iOS Version 11.0 is not the minimum version. Default value is False." +
							"Exactly one of the minimum operating system boolean values will be TRUE.",
						PlanModifiers: []planmodifier.Bool{
							planmodifiers.BoolDefaultValue(false),
						},
					},
					"v12_0": schema.BoolAttribute{
						Optional: true,
						Computed: true,
						MarkdownDescription: "Indicates the minimum iOS version support required for the managed device." +
							"When 'True', iOS with OS Version 12.0 or later is required to install the app. " +
							"If 'False', iOS Version 12.0 is not the minimum version. Default value is False." +
							"Exactly one of the minimum operating system boolean values will be TRUE.",
						PlanModifiers: []planmodifier.Bool{
							planmodifiers.BoolDefaultValue(false),
						},
					},
					"v13_0": schema.BoolAttribute{
						Optional: true,
						Computed: true,
						MarkdownDescription: "Indicates the minimum iOS version support required for the managed device." +
							"When 'True', iOS with OS Version 13.0 or later is required to install the app. " +
							"If 'False', iOS Version 13.0 is not the minimum version. Default value is False." +
							"Exactly one of the minimum operating system boolean values will be TRUE.",
						PlanModifiers: []planmodifier.Bool{
							planmodifiers.BoolDefaultValue(false),
						},
					},
					"v14_0": schema.BoolAttribute{
						Optional: true,
						Computed: true,
						MarkdownDescription: "Indicates the minimum iOS version support required for the managed device." +
							"When 'True', iOS with OS Version 14.0 or later is required to install the app. " +
							"If 'False', iOS Version 14.0 is not the minimum version. Default value is False." +
							"Exactly one of the minimum operating system boolean values will be TRUE.",
						PlanModifiers: []planmodifier.Bool{
							planmodifiers.BoolDefaultValue(false),
						},
					},
					"v15_0": schema.BoolAttribute{
						Optional: true,
						Computed: true,
						MarkdownDescription: "Indicates the minimum iOS version support required for the managed device." +
							"When 'True', iOS with OS Version 15.0 or later is required to install the app. " +
							"If 'False', iOS Version 15.0 is not the minimum version. Default value is False." +
							"Exactly one of the minimum operating system boolean values will be TRUE.",
						PlanModifiers: []planmodifier.Bool{
							planmodifiers.BoolDefaultValue(false),
						},
					},
					"v16_0": schema.BoolAttribute{
						Optional: true,
						Computed: true,
						MarkdownDescription: "Indicates the minimum iOS version support required for the managed device." +
							"When 'True', iOS with OS Version 16.0 or later is required to install the app. " +
							"If 'False', iOS Version 16.0 is not the minimum version. Default value is False." +
							"Exactly one of the minimum operating system boolean values will be TRUE.",
						PlanModifiers: []planmodifier.Bool{
							planmodifiers.BoolDefaultValue(false),
						},
					},
					"v17_0": schema.BoolAttribute{
						Optional: true,
						Computed: true,
						MarkdownDescription: "Indicates the minimum iOS version support required for the managed device." +
							"When 'True', iOS with OS Version 17.0 or later is required to install the app. " +
							"If 'False', iOS Version 17.0 is not the minimum version. Default value is False." +
							"Exactly one of the minimum operating system boolean values will be TRUE.",
						PlanModifiers: []planmodifier.Bool{
							planmodifiers.BoolDefaultValue(false),
						},
					},
					"v18_0": schema.BoolAttribute{
						Optional: true,
						Computed: true,
						MarkdownDescription: "Indicates the minimum iOS version support required for the managed device." +
							"When 'True', iOS with OS Version 18.0 or later is required to install the app. " +
							"If 'False', iOS Version 18.0 is not the minimum version. Default value is False." +
							"Exactly one of the minimum operating system boolean values will be TRUE.",
						PlanModifiers: []planmodifier.Bool{
							planmodifiers.BoolDefaultValue(false),
						},
					},
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
			"app_icon": commonschemagraphbeta.MobileAppIconSchema(),
			"timeouts": commonschema.Timeouts(ctx),
		},
	}
}
