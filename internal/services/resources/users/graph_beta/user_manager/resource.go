package graphBetaUsersUserManager

import (
	"context"
	"regexp"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
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
	ResourceName        = "microsoft365_graph_beta_users_user_manager"
	CreateTimeout       = 180
	UpdateTimeout       = 180
	ReadTimeout         = 180
	DeleteTimeout       = 180
	ResourceDocMarkdown = "Manages the manager relationship for a user in Microsoft Entra ID using the Microsoft Graph Beta API."
)

var (
	// Basic resource interface (CRUD operations)
	_ resource.Resource = &UserManagerResource{}

	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &UserManagerResource{}

	// Allows the resource to define a custom import state logic
	_ resource.ResourceWithImportState = &UserManagerResource{}
)

// UserManagerResource defines the resource implementation.
type UserManagerResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ProviderTypeName string
	TypeName         string
	ReadPermissions  []string
	WritePermissions []string
}

// NewUserManagerResource returns a new instance of the resource.
func NewUserManagerResource() resource.Resource {
	return &UserManagerResource{
		ReadPermissions: []string{
			"User.Read.All",
		},
		WritePermissions: []string{
			"User.ReadWrite.All",
			"AgentIdUser.ReadWrite.IdentityParentedBy,",
			"AgentIdUser.ReadWrite.All",
		},
	}
}

// Metadata returns the resource type name.
func (r *UserManagerResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = ResourceName
}

// Configure sets the client for the resource.
func (r *UserManagerResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, ResourceName)
}

// ImportState handles the import of the resource.
// The import ID is the user_id (the user whose manager relationship is being imported).
func (r *UserManagerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Set both id and user_id from the import ID since they are the same value
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("user_id"), req.ID)...)
}

// Schema defines the schema for the resource.
func (r *UserManagerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: ResourceDocMarkdown,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The unique identifier for this user manager relationship. This is the same as user_id.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"user_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The unique identifier (GUID) of the user whose manager is being managed.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(constants.GuidRegex),
						"must be a valid GUID in the format xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
					),
				},
			},
			"manager_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The unique identifier (GUID) of the user or organizational contact to assign as the manager.",
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(constants.GuidRegex),
						"must be a valid GUID in the format xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
					),
				},
			},
			"timeouts": commonschema.ResourceTimeouts(ctx),
		},
	}
}
