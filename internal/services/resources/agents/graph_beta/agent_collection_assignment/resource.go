package graphBetaAgentsAgentCollectionAssignment

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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	ResourceName  = "microsoft365_graph_beta_agents_agent_collection_assignment"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	_ resource.Resource                = &AgentCollectionAssignmentResource{}
	_ resource.ResourceWithConfigure   = &AgentCollectionAssignmentResource{}
	_ resource.ResourceWithImportState = &AgentCollectionAssignmentResource{}
)

func NewAgentCollectionAssignmentResource() resource.Resource {
	return &AgentCollectionAssignmentResource{
		ReadPermissions: []string{
			"AgentCollection.Read.All",
			"AgentInstance.Read.All",
		},
		WritePermissions: []string{
			"AgentCollection.ReadWrite.All",
			"AgentCollection.ReadWrite.ManagedBy",
			"AgentCollection.ReadWrite.Global",
			"AgentCollection.ReadWrite.Quarantined",
			"AgentInstance.Read.All",
		},
		ResourcePath: "/agentRegistry/agentInstances/{agentInstanceId}/collections/{agentCollectionId}/members",
	}
}

type AgentCollectionAssignmentResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

// Metadata returns the resource type name.
func (r *AgentCollectionAssignmentResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = ResourceName
}

// Configure sets the client for the resource.
func (r *AgentCollectionAssignmentResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, ResourceName)
}

// ImportState handles importing the resource.
// Import ID format: {agent_instance_id}/{agent_collection_id}
func (r *AgentCollectionAssignmentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.Split(req.ID, "/")
	if len(parts) != 2 {
		resp.Diagnostics.AddError(
			"Invalid import ID format",
			fmt.Sprintf("Import ID must be in format: agent_instance_id/agent_collection_id. Got: %s", req.ID),
		)
		return
	}

	agentInstanceID := parts[0]
	agentCollectionID := parts[1]

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("agent_instance_id"), agentInstanceID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("agent_collection_id"), agentCollectionID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
}

// Schema returns the schema for the resource.
func (r *AgentCollectionAssignmentResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages the assignment of an agent instance to an agent collection in the Microsoft Entra Agent Registry. " +
			"This resource adds an agent instance as a member of an agent collection.\n\n" +
			"Use this resource to control which agent instances belong to which collections for organizational and access control purposes.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique identifier for the assignment. Format: {agent_instance_id}/{agent_collection_id}.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"agent_instance_id": schema.StringAttribute{
				MarkdownDescription: "The unique identifier of the agent instance to add as a member of the collection.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(constants.GuidRegex),
						"must be a valid GUID in the format xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
					),
				},
			},
			"agent_collection_id": schema.StringAttribute{
				MarkdownDescription: "The unique identifier of the agent collection to add the agent instance to.",
				Required:            true,
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
