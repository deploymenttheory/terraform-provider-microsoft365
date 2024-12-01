package graphBetaM365AppsInstallationOptions

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/schema"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	ResourceName  = "graph_beta_device_and_app_management_m365_apps_installation_options"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	// Basic resource interface (CRUD operations)
	_ resource.Resource = &M365AppsInstallationOptionsResource{}

	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &M365AppsInstallationOptionsResource{}

	// Enables import functionality
	_ resource.ResourceWithImportState = &M365AppsInstallationOptionsResource{}

	// Enables plan modification/diff suppression
	_ resource.ResourceWithModifyPlan = &M365AppsInstallationOptionsResource{}
)

func NewM365AppsInstallationOptionsResource() resource.Resource {
	return &M365AppsInstallationOptionsResource{
		ReadPermissions: []string{
			"OrgSettings-Microsoft365Install.Read.All",
		},
		WritePermissions: []string{
			"OrgSettings-Microsoft365Install.ReadWrite.All",
		},
		ResourcePath: "/admin/microsoft365Apps/installationOptions",
	}
}

type M365AppsInstallationOptionsResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ProviderTypeName string
	TypeName         string
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

// Metadata returns the resource type name.
func (r *M365AppsInstallationOptionsResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + ResourceName
}

// Configure sets the client for the resource.
func (r *M365AppsInstallationOptionsResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = common.SetGraphBetaClientForResource(ctx, req, resp, r.TypeName)
}

// ImportState imports the resource state.
func (r *M365AppsInstallationOptionsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *M365AppsInstallationOptionsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The unique identifier for the M365AppsInstallationOptions.",
			},
			"update_channel": schema.StringAttribute{
				Required: true,
				Description: "Specifies how often users get feature updates for Microsoft 365 apps installed on devices running Windows. " +
					"The possible values are: `current`, `monthlyEnterprise`, or `semiAnnual`, with corresponding update frequencies of " +
					"`As soon as they're ready`, `Once a month`, and `Every six months`.",
				Validators: []validator.String{
					stringvalidator.OneOf("current", "monthlyEnterprise", "semiAnnual"),
				},
			},
			"apps_for_windows": schema.SingleNestedAttribute{
				Required: true,
				Attributes: map[string]schema.Attribute{
					"is_microsoft_365_apps_enabled": schema.BoolAttribute{
						Required:    true,
						Description: "Specifies whether users can install Microsoft 365 apps, including Skype for Business, on their Windows devices. The default value is `true`.",
					},
					"is_project_enabled": schema.BoolAttribute{
						Required:    true,
						Description: "Specifies whether users can install Microsoft Project on their Windows devices. The default value is `true`.",
					},
					"is_skype_for_business_enabled": schema.BoolAttribute{
						Required:    true,
						Description: "Specifies whether users can install Skype for Business (standalone) on their Windows devices. The default value is `true`.",
					},
					"is_visio_enabled": schema.BoolAttribute{
						Required:    true,
						Description: "Specifies whether users can install Visio on their Windows devices. The default value is `true`.",
					},
				},
				Description: "The Microsoft 365 apps installation options container object for a Windows platform.",
			},
			"apps_for_mac": schema.SingleNestedAttribute{
				Required: true,
				Attributes: map[string]schema.Attribute{
					"is_microsoft_365_apps_enabled": schema.BoolAttribute{
						Required:    true,
						Description: "Specifies whether users can install Microsoft 365 apps on their MAC devices. The default value is `true`.",
					},
					"is_skype_for_business_enabled": schema.BoolAttribute{
						Required:    true,
						Description: "Specifies whether users can install Skype for Business on their MAC devices running OS X El Capitan 10.11 or later. The default value is `true`.",
					},
				},
				Description: "The Microsoft 365 apps installation options container object for a MAC platform.",
			},
			"timeouts": commonschema.Timeouts(ctx),
		},
	}
}
