package graphBetaApplicationOwner

import (
	"context"
	"regexp"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	planmodifiers "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/plan_modifiers"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	ResourceName  = "microsoft365_graph_beta_applications_application_owner"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	// Basic resource interface (CRUD operations)
	_ resource.Resource = &ApplicationOwnerResource{}

	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &ApplicationOwnerResource{}

	// Enables import functionality
	_ resource.ResourceWithImportState = &ApplicationOwnerResource{}
)

func NewApplicationOwnerResource() resource.Resource {
	return &ApplicationOwnerResource{
		ReadPermissions: []string{
			"Application.Read.All",
			"Directory.Read.All",
		},
		WritePermissions: []string{
			"Application.ReadWrite.All",
			"Directory.ReadWrite.All",
		},
		ResourcePath: "/applications",
	}
}

type ApplicationOwnerResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

// Metadata returns the resource type name.
func (r *ApplicationOwnerResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = ResourceName
}

// Configure sets the client for the resource.
func (r *ApplicationOwnerResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, ResourceName)
}

// ImportState imports the resource state.
func (r *ApplicationOwnerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Schema defines the schema for the resource.
func (r *ApplicationOwnerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages Azure AD/Entra application owner assignments using the `/applications/{application-id}/owners` endpoint. This resource enables adding and removing users or service principals as owners of applications.\n\n**Owner Type Support:**\n- **Users**: Individual user accounts\n- **Service Principals**: Service principal objects\n\n**Important Notes:**\n- Owners can modify the application object\n- As a recommended best practice, apps should have at least two owners\n- The last owner (user object) of an application cannot be removed\n\n**Required Permissions by Owner Type:**\n- **Users**: `Application.ReadWrite.All` or `Directory.ReadWrite.All`\n- **Service Principals**: `Application.ReadWrite.All` or `Directory.ReadWrite.All`\n\n**Supported Microsoft Entra Roles:**\n- Application owners (can modify their own applications)\n- Application Developer (for applications they own)\n- Cloud Application Administrator\n- Application Administrator\n- Hybrid Identity Administrator.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
				MarkdownDescription: "The unique identifier for this application owner assignment. This is a composite ID formed by combining the application ID and owner ID.",
			},
			"application_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The unique identifier (UUID) for the application.",
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(constants.GuidRegex),
						"must be a valid GUID in the format 00000000-0000-0000-0000-000000000000",
					),
				},
			},
			"owner_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The unique identifier (UUID) for the owner to be added to the application. This can be a user or service principal.",
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(constants.GuidRegex),
						"must be a valid GUID in the format 00000000-0000-0000-0000-000000000000",
					),
				},
			},
			"owner_object_type": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The type of object being added as an owner. This determines the correct Microsoft Graph API endpoint to use. Valid values: 'User', 'ServicePrincipal'.",
				Validators: []validator.String{
					stringvalidator.OneOf("User", "ServicePrincipal"),
				},
			},
			"owner_type": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The type of the owner object as returned by Microsoft Graph (e.g., 'User', 'ServicePrincipal'). Read-only.",
			},
			"owner_display_name": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The display name of the owner. Read-only.",
			},
			"timeouts": commonschema.ResourceTimeouts(ctx),
		},
	}
}
