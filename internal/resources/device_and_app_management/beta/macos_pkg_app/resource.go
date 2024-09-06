package graphbetamacospkgapp

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ resource.Resource = &MacOSPkgAppResource{}
var _ resource.ResourceWithConfigure = &MacOSPkgAppResource{}
var _ resource.ResourceWithImportState = &MacOSPkgAppResource{}

func NewMacOSPkgAppResource() resource.Resource {
	return &MacOSPkgAppResource{}
}

type MacOSPkgAppResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ProviderTypeName string
	TypeName         string
}

// GetID returns the ID of a resource from the state model.
func (s *MacOSPkgAppResourceModel) GetID() string {
	return s.ID.ValueString()
}

// GetTypeName returns the type name of the resource from the state model.
func (r *MacOSPkgAppResource) GetTypeName() string {
	return r.TypeName
}

// Metadata returns the resource type name.
func (r *MacOSPkgAppResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_graph_beta_device_and_app_management_macos_pkg_app"
}

// Configure sets the client for the resource.
func (r *MacOSPkgAppResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = common.SetGraphBetaClientForResource(ctx, req, resp, r.TypeName)
}

// ImportState imports the resource state.
func (r *MacOSPkgAppResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *MacOSPkgAppResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "The resource `macos_pkg_app` manages a macOS PKG app",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Key of the entity. This property is read-only.",
				Computed:    true,
			},
			"display_name": schema.StringAttribute{
				Description: "The admin provided or imported title of the app.",
				Required:    true,
			},
			"description": schema.StringAttribute{
				Description: "The description of the app.",
				Optional:    true,
			},
			"publisher": schema.StringAttribute{
				Description: "The publisher of the app.",
				Optional:    true,
			},
			"large_icon": schema.SingleNestedAttribute{
				Description: "The large icon, to be displayed in the app details and used for upload of the icon.",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"type": schema.StringAttribute{
						Description: "The MIME type of the icon.",
						Required:    true,
					},
					"value": schema.StringAttribute{
						Description: "The base64-encoded icon data.",
						Required:    true,
					},
				},
			},
			"created_date_time": schema.StringAttribute{
				Description: "The date and time the app was created. This property is read-only.",
				Computed:    true,
			},
			"last_modified_date_time": schema.StringAttribute{
				Description: "The date and time the app was last modified. This property is read-only.",
				Computed:    true,
			},
			"is_featured": schema.BoolAttribute{
				Description: "The value indicating whether the app is marked as featured by the admin.",
				Optional:    true,
			},
			"privacy_information_url": schema.StringAttribute{
				Description: "The privacy statement Url.",
				Optional:    true,
			},
			"information_url": schema.StringAttribute{
				Description: "The more information Url.",
				Optional:    true,
			},
			"owner": schema.StringAttribute{
				Description: "The owner of the app.",
				Optional:    true,
			},
			"developer": schema.StringAttribute{
				Description: "The developer of the app.",
				Optional:    true,
			},
			"notes": schema.StringAttribute{
				Description: "Notes for the app.",
				Optional:    true,
			},
			"upload_state": schema.Int64Attribute{
				Description: "The upload state. Possible values are: 0 - `Not Ready`, 1 - `Ready`, 2 - `Processing`. This property is read-only.",
				Computed:    true,
			},
			"publishing_state": schema.StringAttribute{
				Description: "The publishing state for the app. Possible values are: `notPublished`, `processing`, `published`.",
				Computed:    true,
			},
			"is_assigned": schema.BoolAttribute{
				Description: "The value indicating whether the app is assigned to at least one group. This property is read-only.",
				Computed:    true,
			},
			"role_scope_tag_ids": schema.ListAttribute{
				Description: "List of scope tag ids for this mobile app.",
				Optional:    true,
				ElementType: types.StringType,
			},
			"dependent_app_count": schema.Int64Attribute{
				Description: "The total number of dependencies the child app has. This property is read-only.",
				Computed:    true,
			},
			"superseding_app_count": schema.Int64Attribute{
				Description: "The total number of apps this app directly or indirectly supersedes. This property is read-only.",
				Computed:    true,
			},
			"superseded_app_count": schema.Int64Attribute{
				Description: "The total number of apps this app is directly or indirectly superseded by. This property is read-only.",
				Computed:    true,
			},
			"committed_content_version": schema.StringAttribute{
				Description: "The internal committed content version.",
				Computed:    true,
			},
			"file_name": schema.StringAttribute{
				Description: "The name of the main Lob application file.",
				Required:    true,
			},
			"size": schema.Int64Attribute{
				Description: "The total size, including all uploaded files. This property is read-only.",
				Computed:    true,
			},
			"primary_bundle_id": schema.StringAttribute{
				Description: "The bundleId of the primary app in the PKG. This maps to the CFBundleIdentifier in the app's bundle configuration.",
				Required:    true,
			},
			"primary_bundle_version": schema.StringAttribute{
				Description: "The version of the primary app in the PKG. This maps to the CFBundleShortVersion in the app's bundle configuration.",
				Required:    true,
			},
			"included_apps": schema.ListNestedAttribute{
				Description: "The list of apps expected to be installed by the PKG. This collection can contain a maximum of 500 elements.",
				Optional:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"bundle_id": schema.StringAttribute{
							Description: "The bundle ID of the included app.",
							Required:    true,
						},
						"bundle_version": schema.StringAttribute{
							Description: "The bundle version of the included app.",
							Required:    true,
						},
					},
				},
			},
			"ignore_version_detection": schema.BoolAttribute{
				Description: "When TRUE, indicates that the app's version will NOT be used to detect if the app is installed on a device. When FALSE, indicates that the app's version will be used to detect if the app is installed on a device. Set this to true for apps that use a self update feature. The default value is FALSE.",
				Optional:    true,
			},
			"minimum_supported_operating_system": schema.SingleNestedAttribute{
				Description: "The minimum operating system applicable for the application.",
				Required:    true,
				Attributes: map[string]schema.Attribute{
					"v10_7": schema.BoolAttribute{
						Description: "Mac OS 10.7 or later.",
						Optional:    true,
					},
					"v10_8": schema.BoolAttribute{
						Description: "Mac OS 10.8 or later.",
						Optional:    true,
					},
					"v10_9": schema.BoolAttribute{
						Description: "Mac OS 10.9 or later.",
						Optional:    true,
					},
					"v10_10": schema.BoolAttribute{
						Description: "Mac OS 10.10 or later.",
						Optional:    true,
					},
					"v10_11": schema.BoolAttribute{
						Description: "Mac OS 10.11 or later.",
						Optional:    true,
					},
					"v10_12": schema.BoolAttribute{
						Description: "Mac OS 10.12 or later.",
						Optional:    true,
					},
					"v10_13": schema.BoolAttribute{
						Description: "Mac OS 10.13 or later.",
						Optional:    true,
					},
					"v10_14": schema.BoolAttribute{
						Description: "Mac OS 10.14 or later.",
						Optional:    true,
					},
					"v10_15": schema.BoolAttribute{
						Description: "Mac OS 10.15 or later.",
						Optional:    true,
					},
					"v11_0": schema.BoolAttribute{
						Description: "Mac OS 11.0 or later.",
						Optional:    true,
					},
					"v12_0": schema.BoolAttribute{
						Description: "Mac OS 12.0 or later.",
						Optional:    true,
					},
					"v13_0": schema.BoolAttribute{
						Description: "Mac OS 13.0 or later.",
						Optional:    true,
					},
					"v14_0": schema.BoolAttribute{
						Description: "Mac OS 14.0 or later.",
						Optional:    true,
					},
				},
			},
			"pre_install_script": schema.SingleNestedAttribute{
				Description: "The pre-install script for the app. This will execute on the macOS device before the app is installed.",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"script_content": schema.StringAttribute{
						Description: "The script content.",
						Required:    true,
					},
				},
			},
			"post_install_script": schema.SingleNestedAttribute{
				Description: "The post-install script for the app. This will execute on the macOS device after the app is installed.",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"script_content": schema.StringAttribute{
						Description: "The script content.",
						Required:    true,
					},
				},
			},
			"timeouts": commonschema.Timeouts(ctx),
		},
	}
}
