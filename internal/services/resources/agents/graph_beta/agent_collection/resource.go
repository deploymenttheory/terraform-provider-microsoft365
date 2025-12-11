package graphBetaAgentCollection

import (
	"context"
	"regexp"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	ResourceName  = "microsoft365_graph_beta_agents_agent_collection"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	_ resource.Resource                = &AgentCollectionResource{}
	_ resource.ResourceWithConfigure   = &AgentCollectionResource{}
	_ resource.ResourceWithImportState = &AgentCollectionResource{}
)

func NewAgentCollectionResource() resource.Resource {
	return &AgentCollectionResource{
		ReadPermissions: []string{
			"AgentCollection.Read.All",
		},
		WritePermissions: []string{
			"AgentCollection.ReadWrite.All",
			"AgentCollection.ReadWrite.ManagedBy",
			"AgentCollection.ReadWrite.Global",
			"AgentCollection.ReadWrite.Quarantined",
		},
		ResourcePath: "/agentRegistry/agentCollections",
	}
}

type AgentCollectionResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

// Metadata returns the resource type name.
func (r *AgentCollectionResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = ResourceName
}

// Configure sets the client for the resource.
func (r *AgentCollectionResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, ResourceName)
}

// ImportState handles importing the resource.
func (r *AgentCollectionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Schema returns the schema for the resource.
func (r *AgentCollectionResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages an Agent Collection in the Microsoft Entra Agent Registry using the `/agentRegistry/agentCollections` endpoint. " +
			"An agent collection represents a grouping of agent instances for organizational and access control purposes.\n\n" +
			"**Reserved Collections**: Two system-reserved collections are always available per tenant:\n" +
			"- **Global** (ID: `00000000-0000-0000-0000-000000000001`): Tenant-wide pool of generally available agents\n" +
			"- **Quarantined** (ID: `00000000-0000-0000-0000-000000000002`): Holding area for blocked/review-pending agents\n\n" +
			"Reserved collections cannot be updated or deleted. Attempting to create a collection with a reserved name returns a 409 Conflict error.\n\n" +
			"For more information, see the [agentCollection resource type](https://learn.microsoft.com/en-us/graph/api/resources/agentcollection?view=graph-rest-beta).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Unique identifier for the collection. Key. Inherited from entity.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"display_name": schema.StringAttribute{
				MarkdownDescription: "Friendly name of the collection.",
				Required:            true,
			},
			"owner_ids": schema.SetAttribute{
				MarkdownDescription: "List of object IDs for the owners of the agent collection.",
				Required:            true,
				ElementType:         types.StringType,
				Validators: []validator.Set{
					setvalidator.ValueStringsAre(
						stringvalidator.RegexMatches(
							regexp.MustCompile(constants.GuidRegex),
							"must be a valid GUID in the format xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
						),
					),
				},
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Description / purpose of the collection.",
				Optional:            true,
			},
			"managed_by": schema.StringAttribute{
				MarkdownDescription: "**appId** (referred to as **Application (client) ID** on the Microsoft Entra admin center) of the service principal managing this agent collection.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(constants.GuidRegex),
						"must be a valid GUID in the format xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
					),
				},
			},
			"originating_store": schema.StringAttribute{
				MarkdownDescription: "Source system/store where the collection originated. For example Copilot Studio. Changing this value will trigger resource recreation.",
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"created_by": schema.StringAttribute{
				MarkdownDescription: "Object ID of the user or app that created the agent collection. Read-only.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"created_date_time": schema.StringAttribute{
				MarkdownDescription: "Timestamp when agent collection was created. Read-only.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"last_modified_date_time": schema.StringAttribute{
				MarkdownDescription: "Timestamp of last modification. Read-only.",
				Computed:            true,
			},
			"timeouts": commonschema.Timeouts(ctx),
		},
	}
}
