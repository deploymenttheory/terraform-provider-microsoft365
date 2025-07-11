package graphBetaMobileAppSupersedence

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	planmodifiers "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/plan_modifiers"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	ResourceName  = "graph_beta_device_and_app_management_mobile_app_supersedence"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	// Basic resource interface (CRUD operations)
	_ resource.Resource = &MobileAppSupersedenceResource{}

	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &MobileAppSupersedenceResource{}

	// Enables import functionality
	_ resource.ResourceWithImportState = &MobileAppSupersedenceResource{}

	// Enables plan modification/diff suppression
	_ resource.ResourceWithModifyPlan = &MobileAppSupersedenceResource{}
)

func NewMobileAppSupersedenceResource() resource.Resource {
	return &MobileAppSupersedenceResource{
		ReadPermissions: []string{
			"DeviceManagementApps.Read.All",
			"DeviceManagementConfiguration.Read.All",
		},
		WritePermissions: []string{
			"DeviceManagementApps.ReadWrite.All",
			"DeviceManagementConfiguration.ReadWrite.All",
		},
		ResourcePath: "/deviceAppManagement/mobileApps/{mobileAppId}/relationships",
	}
}

type MobileAppSupersedenceResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ProviderTypeName string
	TypeName         string
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

// Metadata returns the resource type name.
func (r *MobileAppSupersedenceResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	r.ProviderTypeName = req.ProviderTypeName
	r.TypeName = ResourceName
	resp.TypeName = r.FullTypeName()
}

// FullTypeName returns the full resource type name in the format "providername_resourcename".
func (r *MobileAppSupersedenceResource) FullTypeName() string {
	return r.ProviderTypeName + "_" + r.TypeName
}

// Configure sets the client for the resource.
func (r *MobileAppSupersedenceResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, constants.PROVIDER_NAME+"_"+ResourceName)
}

// ImportState imports the resource state.
func (r *MobileAppSupersedenceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Schema returns the schema for the resource.
func (r *MobileAppSupersedenceResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages mobile app supersedence relationships in Microsoft Intune using the /deviceAppManagement/mobileApps/{mobileAppId}/relationships endpoint. Supersedence enables admins to upgrade or replace existing apps with newer versions in a controlled manner.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
				MarkdownDescription: "The unique identifier for the mobile app relationship entity. This is assigned at MobileAppRelationship entity creation. For example: 2dbc75b9-e993-4e4d-a071-91ac5a218672_43aaaf35-ce51-4695-9447-5eac6df31161. Read-Only.",
			},
			"target_id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				MarkdownDescription: "The unique app identifier of the target of the mobile app relationship entity. For example: 2dbc75b9-e993-4e4d-a071-91ac5a218672.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"target_display_name": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The display name of the app that is the target of the mobile app relationship entity. For example: Firefox Setup 52.0.2 32bit.intunewin. Maximum length is 500 characters. Read-Only.",
			},
			"target_display_version": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The display version of the app that is the target of the mobile app relationship entity. For example 1.0 or 1.2203.156. Read-Only.",
			},
			"target_publisher": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The publisher of the app that is the target of the mobile app relationship entity. For example: Fabrikam. Maximum length is 500 characters. Read-Only.",
			},
			"target_publisher_display_name": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The publisher display name of the app that is the target of the mobile app relationship entity. For example: Fabrikam. Maximum length is 500 characters. Read-Only.",
			},
			"source_id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				MarkdownDescription: "The unique app identifier of the source of the mobile app relationship entity. For example: 2dbc75b9-e993-4e4d-a071-91ac5a218672. If null during relationship creation, then it will be populated with parent Id. Read-Only.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"source_display_name": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The display name of the app that is the source of the mobile app relationship entity. For example: Orca. Maximum length is 500 characters. Read-Only.",
			},
			"source_display_version": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The display version of the app that is the source of the mobile app relationship entity. For example 1.0.12 or 1.2203.156 or 3. Read-Only.",
			},
			"source_publisher_display_name": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The publisher display name of the app that is the source of the mobile app relationship entity. For example: Fabrikam. Maximum length is 500 characters. Read-Only.",
			},
			"target_type": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The type of relationship indicating whether the target application of a relationship is a parent or child in the relationship. Possible values are: parent, child. Read-Only.",
			},
			"app_relationship_type": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The type of relationship indicating whether the target application of a relationship is a parent or child in the relationship. Possible values are: parent, child. Read-Only.",
			},
			"supersedence_type": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The supersedence relationship type between the parent and child apps. Possible values are: update, replace.",
				Validators: []validator.String{
					stringvalidator.OneOf("update", "replace"),
				},
			},
			"superseded_app_count": schema.Int32Attribute{
				Computed:            true,
				MarkdownDescription: "The total number of apps directly or indirectly superseded by the child app. Read-Only.",
			},
			"superseding_app_count": schema.Int32Attribute{
				Computed:            true,
				MarkdownDescription: "The total number of apps directly or indirectly superseding the parent app. Read-Only.",
			},
			"timeouts": commonschema.Timeouts(ctx),
		},
	}
}
