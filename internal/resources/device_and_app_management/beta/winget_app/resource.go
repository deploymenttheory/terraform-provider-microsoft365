package graphBetaWinGetApp

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/schema"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

var _ resource.Resource = &WinGetAppResource{}
var _ resource.ResourceWithConfigure = &WinGetAppResource{}
var _ resource.ResourceWithImportState = &WinGetAppResource{}

func NewWinGetAppResource() resource.Resource {
	return &WinGetAppResource{}
}

type WinGetAppResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ProviderTypeName string
	TypeName         string
}

// GetID returns the ID of a resource from the state model.
func (s *WinGetAppResourceModel) GetID() string {
	return s.ID.ValueString()
}

// GetTypeName returns the type name of the resource from the state model.
func (r *WinGetAppResource) GetTypeName() string {
	return r.TypeName
}

// Metadata returns the resource type name.
func (r *WinGetAppResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_graph_beta_device_and_app_management_win_get_app"
}

// Configure sets the client for the resource.
func (r *WinGetAppResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = common.SetGraphBetaClientForResource(ctx, req, resp, r.TypeName)
}

// ImportState imports the resource state.
func (r *WinGetAppResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *WinGetAppResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a WinGet application in Microsoft Intune.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Key of the entity. This property is read-only.",
			},
			"display_name": schema.StringAttribute{
				Required:    true,
				Description: "The admin provided or imported title of the app.",
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "The description of the app.",
			},
			"publisher": schema.StringAttribute{
				Optional:    true,
				Description: "The publisher of the app.",
			},
			"large_icon": schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{
					"type": schema.StringAttribute{
						Required:    true,
						Description: "The MIME type of the icon.",
					},
					"value": schema.StringAttribute{
						Required:    true,
						Description: "The icon data, request for mat should be in raw bytes.",
					},
				},
				Description: "The large icon, to be displayed in the app details and used for upload of the icon.",
			},
			"created_date_time": schema.StringAttribute{
				Computed:    true,
				Description: "The date and time the app was created. This property is read-only.",
			},
			"last_modified_date_time": schema.StringAttribute{
				Computed:    true,
				Description: "The date and time the app was last modified. This property is read-only.",
			},
			"is_featured": schema.BoolAttribute{
				Optional:    true,
				Description: "The value indicating whether the app is marked as featured by the admin.",
			},
			"privacy_information_url": schema.StringAttribute{
				Optional:    true,
				Description: "The privacy statement Url.",
			},
			"information_url": schema.StringAttribute{
				Optional:    true,
				Description: "The more information Url.",
			},
			"owner": schema.StringAttribute{
				Optional:    true,
				Description: "The owner of the app.",
			},
			"developer": schema.StringAttribute{
				Optional:    true,
				Description: "The developer of the app.",
			},
			"notes": schema.StringAttribute{
				Optional:    true,
				Description: "Notes for the app.",
			},
			"upload_state": schema.Int64Attribute{
				Computed:    true,
				Description: "The upload state. Possible values are: 0 - Not Ready, 1 - Ready, 2 - Processing. This property is read-only.",
			},
			"publishing_state": schema.StringAttribute{
				Computed: true,
				Description: "The publishing state for the app. The app cannot be assigned unless the app is published. " +
					"Possible values are: notPublished, processing, published.",
			},
			"is_assigned": schema.BoolAttribute{
				Computed:    true,
				Description: "The value indicating whether the app is assigned to at least one group. This property is read-only.",
			},
			"role_scope_tag_ids": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "List of scope tag ids for this mobile app.",
			},
			"dependent_app_count": schema.Int64Attribute{
				Computed:    true,
				Description: "The total number of dependencies the child app has. This property is read-only.",
			},
			"superseding_app_count": schema.Int64Attribute{
				Computed:    true,
				Description: "The total number of apps this app directly or indirectly supersedes. This property is read-only.",
			},
			"superseded_app_count": schema.Int64Attribute{
				Computed:    true,
				Description: "The total number of apps this app is directly or indirectly superseded by. This property is read-only.",
			},
			"manifest_hash": schema.StringAttribute{
				Computed:    true,
				Description: "Hash of package metadata properties used to validate that the application matches the metadata in the source repository.",
			},
			"package_identifier": schema.StringAttribute{
				Required:    true,
				Description: "The PackageIdentifier from the WinGet source repository REST API. This also maps to the Id when using the WinGet client command line application. Required at creation time, cannot be modified on existing objects.",
			},
			"install_experience": schema.SingleNestedAttribute{
				Required: true,
				Attributes: map[string]schema.Attribute{
					"run_as_account": schema.StringAttribute{
						Required: true,
						Description: "The account type (System or User) that actions should be run as on target devices. " +
							"Required at creation time.",
						Validators: []validator.String{
							stringvalidator.OneOf("system", "user"),
						},
					},
				},
				Description: "The install experience settings associated with this application.",
			},
			"assignments": commonschema.DeviceAndAppManagementAssignments(),
			"timeouts":    commonschema.Timeouts(ctx),
		},
	}
}
