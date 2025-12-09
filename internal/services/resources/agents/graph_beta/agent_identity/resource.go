package graphBetaAgentIdentity

import (
	"context"
	"fmt"
	"strings"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	ResourceName  = "microsoft365_graph_beta_agents_agent_identity"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	// Basic resource interface (CRUD operations)
	_ resource.Resource = &AgentIdentityResource{}

	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &AgentIdentityResource{}

	// Enables import functionality
	_ resource.ResourceWithImportState = &AgentIdentityResource{}
)

func NewAgentIdentityResource() resource.Resource {
	return &AgentIdentityResource{
		ReadPermissions: []string{
			"AgentInstance.Read.All",
			"Directory.Read.All",
		},
		WritePermissions: []string{
			"AgentInstance.ReadWrite.All",
			"Directory.ReadWrite.All",
			"AgentIdentity.DeleteRestore.All", // Needed for deletion
		},
		ResourcePath: "/servicePrincipals",
	}
}

type AgentIdentityResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

// Metadata returns the resource type name.
func (r *AgentIdentityResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = ResourceName
}

// Configure sets the client for the resource.
func (r *AgentIdentityResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, ResourceName)
}

// ImportState handles the import functionality.
// Import ID format: {agent_identity_id}/{agent_identity_blueprint_id}
func (r *AgentIdentityResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.Split(req.ID, "/")
	if len(parts) != 2 {
		resp.Diagnostics.AddError(
			"Invalid import ID format",
			fmt.Sprintf("Import ID must be in format: agent_identity_id/agent_identity_blueprint_id. Got: %s", req.ID),
		)
		return
	}

	agentIdentityID := parts[0]
	agentIdentityBlueprintID := parts[1]

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), agentIdentityID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("agent_identity_blueprint_id"), agentIdentityBlueprintID)...)
}

// Schema returns the schema for the resource.
func (r *AgentIdentityResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages an Agent Identity in Microsoft Entra ID using the `/servicePrincipals/microsoft.graph.agentIdentity` endpoint. " +
			"An agent identity is a service principal that represents an AI agent instance, created from an agent identity blueprint. " +
			"Agent identities inherit settings from their blueprint and can be assigned permissions and credentials.\n\n" +
			"For more information, see the [Agent Identity documentation](https://learn.microsoft.com/en-us/graph/api/resources/agentidentity?view=graph-rest-beta).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique identifier for the agent identity service principal. Read-only.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"display_name": schema.StringAttribute{
				MarkdownDescription: "The display name for the agent identity. Maximum length is 256 characters. Required.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 256),
				},
			},
			"agent_identity_blueprint_id": schema.StringAttribute{
				MarkdownDescription: "The application (client) ID of the agent identity blueprint from which this agent identity is created. Required. " +
					"This is the `app_id` of the `microsoft365_graph_beta_agents_agent_identity_blueprint` resource.",
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"account_enabled": schema.BoolAttribute{
				MarkdownDescription: "Set whether the agent identity is enabled. If `false`, the agent identity cannot authenticate or access resources.",
				Required:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"created_by_app_id": schema.StringAttribute{
				MarkdownDescription: "The application ID of the application that created this agent identity. Read-only.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"created_date_time": schema.StringAttribute{
				MarkdownDescription: "The date and time when the agent identity was created. Read-only.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"disabled_by_microsoft_status": schema.StringAttribute{
				MarkdownDescription: "Indicates whether Microsoft has disabled the agent identity. Possible values are: `null`, `NotDisabled`, `DisabledDueToViolationOfServicesAgreement`. Read-only.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"service_principal_type": schema.StringAttribute{
				MarkdownDescription: "The type of the service principal. For agent identities, this is always `ServiceIdentity`. Read-only.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"tags": schema.SetAttribute{
				MarkdownDescription: "Custom strings that can be used to categorize and identify the agent identity.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
			},
			"sponsor_ids": schema.SetAttribute{
				MarkdownDescription: "The user IDs of the sponsors for the agent identity. At least one sponsor is " +
					"required when creating an agent identity. Sponsors are users who can approve or oversee the agent identity.",
				Required:    true,
				ElementType: types.StringType,
			},
			"owner_ids": schema.SetAttribute{
				MarkdownDescription: "The user IDs of the owners for the agent identity. At least one owner is required " +
					"when creating an agent identity. Owners are users who have full control over the agent identity.",
				Required:    true,
				ElementType: types.StringType,
			},
			// TODO: Add custom security attributes
			// can use the same schema as the user resource
			// currently results in a 403. feature probably
			// not supported atm via api. will revisit. 09/12/2025
			"timeouts": commonschema.Timeouts(ctx),
		},
	}
}
