package graphBetaAgentIdentityBlueprintIdentifierUri

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	ResourceName  = "microsoft365_graph_beta_agents_agent_identity_blueprint_identifier_uri"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	// Basic resource interface (CRUD operations)
	_ resource.Resource = &AgentIdentityBlueprintIdentifierUriResource{}

	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &AgentIdentityBlueprintIdentifierUriResource{}

	// Allows the resource to be imported
	_ resource.ResourceWithImportState = &AgentIdentityBlueprintIdentifierUriResource{}
)

func NewAgentIdentityBlueprintIdentifierUriResource() resource.Resource {
	return &AgentIdentityBlueprintIdentifierUriResource{
		ReadPermissions: []string{
			"AgentIdentityBlueprint.Read.All",
			"Application.Read.All",
			"Directory.Read.All",
		},
		WritePermissions: []string{
			"AgentIdentityBlueprint.ReadWrite.All",
			"Directory.ReadWrite.All",
		},
		ResourcePath: "/applications",
	}
}

type AgentIdentityBlueprintIdentifierUriResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

// Metadata returns the resource type name.
func (r *AgentIdentityBlueprintIdentifierUriResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = ResourceName
}

// Configure sets the client for the resource.
func (r *AgentIdentityBlueprintIdentifierUriResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, ResourceName)
}

// ImportState handles the import of the resource.
// Import format: blueprint_id/identifier_uri
func (r *AgentIdentityBlueprintIdentifierUriResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Use SplitN with limit 2 to handle identifier URIs that contain "/" (e.g., api://...)
	parts := strings.SplitN(req.ID, "/", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			fmt.Sprintf("Expected import ID in format 'blueprint_id/identifier_uri', got: %s", req.ID),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("blueprint_id"), parts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("identifier_uri"), parts[1])...)
}

// Schema returns the schema for the resource.
func (r *AgentIdentityBlueprintIdentifierUriResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages an identifier URI and optional OAuth2 permission scope for an Agent Identity Blueprint in Microsoft Entra ID using the `/applications` endpoint. " +
			"This resource configures the identifier URI and optional scope using a " +
			"[PATCH](https://learn.microsoft.com/en-us/graph/api/application-update?view=graph-rest-beta) to the application endpoint.\n\n" +
			"The identifier URI is used to uniquely identify the agent identity blueprint and is required for " +
			"receiving incoming requests from users and other agents.\n\n" +
			"**Note:** This resource manages a single identifier URI. To manage multiple URIs, create multiple resource instances.",
		Attributes: map[string]schema.Attribute{
			"blueprint_id": schema.StringAttribute{
				MarkdownDescription: "The unique identifier (Object ID) of the agent identity blueprint to configure. Required.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"identifier_uri": schema.StringAttribute{
				MarkdownDescription: "The identifier URI for the agent identity blueprint. Valid formats include " +
					"`api://<guid>`, `api://<domain>/<path>`, `https://<domain>/<path>`, or `urn:<namespace>:<identifier>`. Required.",
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					stringvalidator.RegexMatches(
						regexp.MustCompile(constants.IdentifierUriRegex),
						"must be a valid identifier URI (api://, https://, or urn: prefix)",
					),
				},
			},
			"scope": schema.SingleNestedAttribute{
				MarkdownDescription: "Optional OAuth2 permission scope configuration. Defines the scope that allows " +
					"applications to access the agent on behalf of the signed-in user.",
				Optional: true,
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						MarkdownDescription: "The unique identifier for the OAuth2 permission scope. If not specified, a UUID will be generated.",
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"admin_consent_description": schema.StringAttribute{
						MarkdownDescription: "A description of the delegated permission, intended to be read by an administrator granting the permission. " +
							"Default: `Allow the application to access the agent on behalf of the signed-in user.`",
						Optional: true,
						Computed: true,
						Default:  stringdefault.StaticString("Allow the application to access the agent on behalf of the signed-in user."),
					},
					"admin_consent_display_name": schema.StringAttribute{
						MarkdownDescription: "The display name for the permission shown in the admin consent experience. " +
							"Default: `Access agent`",
						Optional: true,
						Computed: true,
						Default:  stringdefault.StaticString("Access agent"),
					},
					"is_enabled": schema.BoolAttribute{
						MarkdownDescription: "Whether the permission scope is enabled. Default: `true`",
						Optional:            true,
						Computed:            true,
						Default:             booldefault.StaticBool(true),
					},
					"type": schema.StringAttribute{
						MarkdownDescription: "The type of permission. Valid values are `User` or `Admin`. Default: `User`",
						Optional:            true,
						Computed:            true,
						Default:             stringdefault.StaticString("User"),
						Validators: []validator.String{
							stringvalidator.OneOf("User", "Admin"),
						},
					},
					"value": schema.StringAttribute{
						MarkdownDescription: "The value of the scope claim that the resource application should expect in the OAuth 2.0 access token. " +
							"Default: `access_agent`",
						Optional: true,
						Computed: true,
						Default:  stringdefault.StaticString("access_agent"),
					},
				},
			},
			"timeouts": commonschema.Timeouts(ctx),
		},
	}
}
