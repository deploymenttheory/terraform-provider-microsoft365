package graphBetaApplicationsAgentIdentityBlueprintServicePrincipal

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	ResourceName  = "microsoft365_graph_beta_agents_agent_identity_blueprint_service_principal"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	// Basic resource interface (CRUD operations)
	_ resource.Resource = &AgentIdentityBlueprintServicePrincipalResource{}

	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &AgentIdentityBlueprintServicePrincipalResource{}

	// Enables import functionality
	_ resource.ResourceWithImportState = &AgentIdentityBlueprintServicePrincipalResource{}
)

func NewAgentIdentityBlueprintServicePrincipalResource() resource.Resource {
	return &AgentIdentityBlueprintServicePrincipalResource{
		ReadPermissions: []string{
			"AgentIdentityBlueprintPrincipal.Read.All",
		},
		WritePermissions: []string{
			"AgentIdentityBlueprintPrincipal.ReadWrite.All",
			"Directory.ReadWrite.All",
		},
		ResourcePath: "/servicePrincipals",
	}
}

type AgentIdentityBlueprintServicePrincipalResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

// Metadata returns the resource type name.
func (r *AgentIdentityBlueprintServicePrincipalResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = ResourceName
}

// Configure sets the client for the resource.
func (r *AgentIdentityBlueprintServicePrincipalResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, ResourceName)
}

// ImportState handles the import functionality.
func (r *AgentIdentityBlueprintServicePrincipalResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Schema returns the schema for the resource.
func (r *AgentIdentityBlueprintServicePrincipalResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a service principal for an Agent Identity Blueprint in Microsoft Entra ID using the `/servicePrincipals/{id}/microsoft.graph.agentIdentityBlueprintPrincipal` endpoint. " +
			"This resource creates an agentIdentityBlueprintPrincipal service principal for an existing agent identity blueprint.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique identifier for the service principal. Read-only.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"app_id": schema.StringAttribute{
				MarkdownDescription: "The application (client) ID of the agent identity blueprint for which to create the service principal. Required.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"timeouts": commonschema.Timeouts(ctx),
		},
	}
}
