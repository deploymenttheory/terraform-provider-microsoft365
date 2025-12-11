package graphBetaAgentInstance

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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	ResourceName  = "microsoft365_graph_beta_agents_agent_instance"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	_ resource.Resource                = &AgentInstanceResource{}
	_ resource.ResourceWithConfigure   = &AgentInstanceResource{}
	_ resource.ResourceWithImportState = &AgentInstanceResource{}
)

func NewAgentInstanceResource() resource.Resource {
	return &AgentInstanceResource{
		ReadPermissions: []string{
			"AgentInstance.Read.All",
			"AgentCardManifest.Read.All",
		},
		WritePermissions: []string{
			"AgentInstance.ReadWrite.All",
			"AgentInstance.ReadWrite.ManagedBy",
			"AgentCardManifest.ReadWrite.All",
			"AgentCardManifest.ReadWrite.ManagedBy",
		},
		ResourcePath: "/agentRegistry/agentInstances",
	}
}

type AgentInstanceResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

// Metadata returns the resource type name.
func (r *AgentInstanceResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = ResourceName
}

// Configure sets the client for the resource.
func (r *AgentInstanceResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, ResourceName)
}

// ImportState handles importing the resource.
func (r *AgentInstanceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Schema returns the schema for the resource.
func (r *AgentInstanceResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages an Agent Instance in the Microsoft Entra Agent Registry using the `/agentRegistry/agentInstances` endpoint. " +
			"An agent instance represents a specific deployed instance of an AI agent. Agent instances can be associated with an " +
			"agentCardManifest that defines its capabilities, skills, and metadata.\n\n" +
			"For more information, see the [agentInstance resource type](https://learn.microsoft.com/en-us/graph/api/resources/agentinstance?view=graph-rest-beta).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Unique identifier for the agent instance. Key. Inherited from entity.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"display_name": schema.StringAttribute{
				MarkdownDescription: "Display name for the agent instance.",
				Required:            true,
			},
			"owner_ids": schema.SetAttribute{
				MarkdownDescription: "List of object IDs for the owners of the agent instance.",
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
			"managed_by": schema.StringAttribute{
				MarkdownDescription: "**appId** (referred to as **Application (client) ID** on the Microsoft Entra admin center) of the application managing this agent.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(constants.GuidRegex),
						"must be a valid GUID in the format xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
					),
				},
			},
			"originating_store": schema.StringAttribute{
				MarkdownDescription: "Name of the store/system where agent originated. For example Copilot Studio, or Microsoft Security Copilot etc. " +
					"Changing this value will force resource recreation.",
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"created_by": schema.StringAttribute{
				MarkdownDescription: "Object ID of the user or application that created the agent instance. Read-only.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"source_agent_id": schema.StringAttribute{
				MarkdownDescription: "Identifier of the agent in the original source system.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(constants.GuidRegex),
						"must be a valid GUID in the format xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
					),
				},
			},
			"agent_identity_blueprint_id": schema.StringAttribute{
				MarkdownDescription: "Object ID of the agentIdentityBlueprint object.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(constants.GuidRegex),
						"must be a valid GUID in the format xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
					),
				},
			},
			"agent_identity_id": schema.StringAttribute{
				MarkdownDescription: "Object ID of the agentIdentity object.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(constants.GuidRegex),
						"must be a valid GUID in the format xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
					),
				},
			},
			"agent_user_id": schema.StringAttribute{
				MarkdownDescription: "Object ID of the agentUser associated with the agent. Read-only.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"created_date_time": schema.StringAttribute{
				MarkdownDescription: "Timestamp when agent instance was created. Read-only.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"last_modified_date_time": schema.StringAttribute{
				MarkdownDescription: "Timestamp of last modification.",
				Computed:            true,
			},
			"url": schema.StringAttribute{
				MarkdownDescription: "Endpoint URL for the agent instance.",
				Optional:            true,
			},
			"preferred_transport": schema.StringAttribute{
				MarkdownDescription: "Preferred transport protocol. The possible values are `JSONRPC`, `GRPC`, and `HTTP+JSON`.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("JSONRPC", "GRPC", "HTTP+JSON"),
				},
			},
			"timeouts": commonschema.Timeouts(ctx),
			"additional_interfaces": schema.ListNestedAttribute{
				MarkdownDescription: "Additional interfaces/transports supported by the agent.",
				Optional:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"url": schema.StringAttribute{
							MarkdownDescription: "URL for the interface.",
							Required:            true,
						},
						"transport": schema.StringAttribute{
							MarkdownDescription: "Transport protocol. The possible values are `JSONRPC`, `GRPC`, and `HTTP+JSON`.",
							Required:            true,
							Validators: []validator.String{
								stringvalidator.OneOf("JSONRPC", "GRPC", "HTTP+JSON"),
							},
						},
					},
				},
			},
			"agent_card_manifest": schema.SingleNestedAttribute{
				MarkdownDescription: "The agent card manifest of the agent instance.",
				Required:            true,
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						MarkdownDescription: "Unique identifier for the agent card manifest.",
						Computed:            true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"owner_ids": schema.SetAttribute{
						MarkdownDescription: "List of owner identifiers for the agent card manifest.",
						Optional:            true,
						ElementType:         types.StringType,
					},
					"originating_store": schema.StringAttribute{
						MarkdownDescription: "Name of the store/system where the manifest originated.",
						Optional:            true,
					},
					"display_name": schema.StringAttribute{
						MarkdownDescription: "Display name for the agent card manifest.",
						Required:            true,
					},
					"description": schema.StringAttribute{
						MarkdownDescription: "Description of the agent card manifest.",
						Required:            true,
					},
					"icon_url": schema.StringAttribute{
						MarkdownDescription: "URL to the icon for the agent. Changing or removingthis value requires resource recreation.",
						Optional:            true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
					},
					"protocol_version": schema.StringAttribute{
						MarkdownDescription: "Protocol version for the agent card. Must be in either the Major.Minor versioning format X.Y (e.g., 1.0, 2.1) " +
							"or the semantic versioning format X.Y.Z (e.g., 1.0.0, 2.1.3)",
						Required: true,
						Validators: []validator.String{
							stringvalidator.Any(
								stringvalidator.RegexMatches(
									regexp.MustCompile(constants.MajorMinorVersionRegex),
									"must be a valid version in the Major.Minor versioning format X.Y (e.g., 1.0, 2.1)",
								),
								stringvalidator.RegexMatches(
									regexp.MustCompile(constants.SemVerRegex),
									"must be a valid semantic version in the format X.Y.Z (e.g., 1.0.0, 2.1.3)",
								),
							),
						},
					},
					"version": schema.StringAttribute{
						MarkdownDescription: "Version of the agent card manifest. Must be in the semantic versioning format X.Y.Z (e.g., 1.0.0, 2.1.3)",
						Required:            true,
						Validators: []validator.String{
							stringvalidator.RegexMatches(
								regexp.MustCompile(constants.SemVerRegex),
								"must be a valid semantic version in the format X.Y.Z (e.g., 1.0.0, 2.1.3)",
							),
						},
					},
					"documentation_url": schema.StringAttribute{
						MarkdownDescription: "URL to the documentation for the agent. Changing or removing this value requires resource recreation.",
						Optional:            true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
					},
					"default_input_modes": schema.SetAttribute{
						MarkdownDescription: "Default input modes supported by the agent. Changing or removing this value requires resource recreation.",
						Optional:            true,
						ElementType:         types.StringType,
						PlanModifiers: []planmodifier.Set{
							setplanmodifier.RequiresReplace(),
						},
					},
					"default_output_modes": schema.SetAttribute{
						MarkdownDescription: "Default output modes supported by the agent. Changing or removingthis value requires resource recreation.",
						Optional:            true,
						ElementType:         types.StringType,
						PlanModifiers: []planmodifier.Set{
							setplanmodifier.RequiresReplace(),
						},
					},
					"supports_authenticated_extended_card": schema.BoolAttribute{
						MarkdownDescription: "Whether the agent supports authenticated extended card.",
						Required:            true,
					},
					"capabilities": schema.SingleNestedAttribute{
						MarkdownDescription: "Capabilities of the agent.",
						Required:            true,
						Attributes: map[string]schema.Attribute{
							"streaming": schema.BoolAttribute{
								MarkdownDescription: "Whether the agent supports streaming.",
								Required:            true,
							},
							"push_notifications": schema.BoolAttribute{
								MarkdownDescription: "Whether the agent supports push notifications.",
								Required:            true,
							},
							"state_transition_history": schema.BoolAttribute{
								MarkdownDescription: "Whether the agent supports state transition history.",
								Required:            true,
							},
							"extensions": schema.ListNestedAttribute{
								MarkdownDescription: "Capability extensions for the agent.",
								Optional:            true,
								Computed:            true,
								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										"uri": schema.StringAttribute{
											MarkdownDescription: "URI of the extension.",
											Required:            true,
										},
										"description": schema.StringAttribute{
											MarkdownDescription: "Description of the extension.",
											Optional:            true,
										},
										"required": schema.BoolAttribute{
											MarkdownDescription: "Whether the extension is required.",
											Optional:            true,
										},
										"params": schema.MapAttribute{
											MarkdownDescription: "Parameters for the extension.",
											Optional:            true,
											ElementType:         types.StringType,
										},
									},
								},
							},
						},
					},
					"provider": schema.SingleNestedAttribute{
						MarkdownDescription: "Provider information for the agent card. Changing this value requires resource recreation.",
						Optional:            true,
						PlanModifiers: []planmodifier.Object{
							objectplanmodifier.RequiresReplace(),
						},
						Attributes: map[string]schema.Attribute{
							"organization": schema.StringAttribute{
								MarkdownDescription: "Organization name of the provider.",
								Optional:            true,
							},
							"url": schema.StringAttribute{
								MarkdownDescription: "URL of the provider.",
								Optional:            true,
							},
						},
					},
					"skills": schema.ListNestedAttribute{
						MarkdownDescription: "Skills defined in the agent card manifest.",
						Optional:            true,
						Computed:            true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"id": schema.StringAttribute{
									MarkdownDescription: "Unique identifier for the skill.",
									Required:            true,
								},
								"display_name": schema.StringAttribute{
									MarkdownDescription: "Display name for the skill.",
									Required:            true,
								},
								"description": schema.StringAttribute{
									MarkdownDescription: "Description of the skill.",
									Optional:            true,
								},
								"tags": schema.SetAttribute{
									MarkdownDescription: "Tags associated with the skill.",
									Optional:            true,
									ElementType:         types.StringType,
								},
								"examples": schema.SetAttribute{
									MarkdownDescription: "Example prompts for the skill.",
									Optional:            true,
									ElementType:         types.StringType,
								},
								"input_modes": schema.SetAttribute{
									MarkdownDescription: "Input modes supported by the skill.",
									Optional:            true,
									ElementType:         types.StringType,
								},
								"output_modes": schema.SetAttribute{
									MarkdownDescription: "Output modes supported by the skill.",
									Optional:            true,
									ElementType:         types.StringType,
								},
							},
						},
					},
				},
			},
		},
	}
}
