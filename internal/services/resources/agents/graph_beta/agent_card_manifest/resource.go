package graphBetaAgentCardManifest

import (
	"context"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/identityschema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
)

const (
	ResourceName  = "microsoft365_graph_beta_agents_agent_card_manifest"
	CreateTimeout = 180
	ReadTimeout   = 180
	UpdateTimeout = 180
	DeleteTimeout = 180
)

var (
	_ resource.Resource                = &AgentCardManifestResource{}
	_ resource.ResourceWithConfigure   = &AgentCardManifestResource{}
	_ resource.ResourceWithImportState = &AgentCardManifestResource{}
	_ resource.ResourceWithModifyPlan  = &AgentCardManifestResource{}
	_ resource.ResourceWithIdentity    = &AgentCardManifestResource{}
)

type AgentCardManifestResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
}

func NewAgentCardManifestResource() resource.Resource {
	return &AgentCardManifestResource{
		ReadPermissions:  []string{"AgentCardManifest.Read.All", "AgentCardManifest.ReadWrite.All", "AgentCardManifest.ReadWrite.ManagedBy"},
		WritePermissions: []string{"AgentCardManifest.ReadWrite.All", "AgentCardManifest.ReadWrite.ManagedBy"},
	}
}

func (r *AgentCardManifestResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = ResourceName
}

func (r *AgentCardManifestResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, ResourceName)
}

func (r *AgentCardManifestResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	if strings.TrimSpace(req.ID) == "" {
		resp.Diagnostics.AddError("Invalid manifest ID", "Import requires the ID of an existing agent card manifest.")
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("agent_card_manifest_id"), req.ID)...)
}

func (r *AgentCardManifestResource) IdentitySchema(_ context.Context, _ resource.IdentitySchemaRequest, resp *resource.IdentitySchemaResponse) {
	resp.IdentitySchema = identityschema.Schema{Attributes: map[string]identityschema.Attribute{"id": identityschema.StringAttribute{RequiredForImport: true}}}
}

func (r *AgentCardManifestResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages an existing Microsoft Graph beta agent card manifest using GET and PATCH /agentRegistry/agentCardManifests/{id}. Create adopts and patches the existing manifest; it does not create an agent instance or a manifest. Destroy removes Terraform state only because Graph documents no standalone manifest DELETE operation. Create the manifest through an agent instance first. Unconfigured optional properties retain remote values; explicitly empty collections are sent as empty arrays. Do not configure the same manifest fields in another Terraform resource. Security dictionaries and extension parameters use JSON to preserve their extensible shape.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Managed manifest ID.",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"agent_card_manifest_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "ID of the existing manifest. Non-GUID IDs are supported. Changing the target adopts a different manifest.",
				Validators:          []validator.String{stringvalidator.LengthAtLeast(1)},
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"display_name": schema.StringAttribute{
				Required:            true,
				Validators:          []validator.String{stringvalidator.LengthAtLeast(1)},
				MarkdownDescription: "Human-readable name of the agent or skill.",
			},
			"description": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
				MarkdownDescription: "Description of the agent, skill, or extension.",
			},
			"icon_url": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
				MarkdownDescription: "URL of the agent icon.",
			},
			"protocol_version": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
				MarkdownDescription: "Agent communication protocol version supported by this manifest.",
			},
			"version": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
				MarkdownDescription: "Version of the agent card manifest.",
			},
			"documentation_url": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
				MarkdownDescription: "URL of the agent documentation.",
			},
			"managed_by": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
				MarkdownDescription: "Application client ID managing the manifest.",
			},
			"originating_store": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
				MarkdownDescription: "Store or registry where the agent originated.",
			},
			"owner_ids": schema.SetAttribute{
				Optional:            true,
				Computed:            true,
				PlanModifiers:       []planmodifier.Set{setplanmodifier.UseStateForUnknown()},
				ElementType:         types.StringType,
				MarkdownDescription: "Object IDs of the manifest owners.",
			},
			"default_input_modes": schema.SetAttribute{
				Optional:            true,
				Computed:            true,
				PlanModifiers:       []planmodifier.Set{setplanmodifier.UseStateForUnknown()},
				ElementType:         types.StringType,
				MarkdownDescription: "Content types accepted by default.",
			},
			"default_output_modes": schema.SetAttribute{
				Optional: true,
				Computed: true, PlanModifiers: []planmodifier.Set{setplanmodifier.UseStateForUnknown()},
				ElementType:         types.StringType,
				MarkdownDescription: "Content types produced by default.",
			},
			"supports_authenticated_extended_card": schema.BoolAttribute{
				Optional: true,
				Computed: true, PlanModifiers: []planmodifier.Bool{boolplanmodifier.UseStateForUnknown()},
				MarkdownDescription: "Whether an authenticated extended agent card is supported.",
			},
			"agent_provider": schema.SingleNestedAttribute{
				Optional: true,
				Computed: true, PlanModifiers: []planmodifier.Object{objectplanmodifier.UseStateForUnknown()},
				MarkdownDescription: "Agent provider details (the Graph provider property).",
				Attributes: map[string]schema.Attribute{
					"organization": schema.StringAttribute{
						Optional:            true,
						Computed:            true,
						PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
						MarkdownDescription: "Name of the organization providing the agent.",
					},
					"url": schema.StringAttribute{
						Optional:            true,
						Computed:            true,
						PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
						MarkdownDescription: "URL of the agent provider.",
					},
				},
			},
			"capabilities": schema.SingleNestedAttribute{
				Optional: true,
				Computed: true, PlanModifiers: []planmodifier.Object{objectplanmodifier.UseStateForUnknown()},
				MarkdownDescription: "Capabilities supported by the agent.",
				Attributes: map[string]schema.Attribute{
					"streaming": schema.BoolAttribute{
						Optional:            true,
						Computed:            true,
						PlanModifiers:       []planmodifier.Bool{boolplanmodifier.UseStateForUnknown()},
						MarkdownDescription: "Whether the agent supports streaming responses.",
					},
					"push_notifications": schema.BoolAttribute{
						Optional:            true,
						Computed:            true,
						PlanModifiers:       []planmodifier.Bool{boolplanmodifier.UseStateForUnknown()},
						MarkdownDescription: "Whether the agent supports push notifications.",
					},
					"state_transition_history": schema.BoolAttribute{
						Optional:            true,
						Computed:            true,
						PlanModifiers:       []planmodifier.Bool{boolplanmodifier.UseStateForUnknown()},
						MarkdownDescription: "Whether the agent exposes state transition history.",
					},
					"extensions": schema.ListNestedAttribute{
						Optional: true,
						Computed: true, PlanModifiers: []planmodifier.List{listplanmodifier.UseStateForUnknown()},
						MarkdownDescription: "Protocol extensions supported by the agent.",
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"uri": schema.StringAttribute{
									Required:            true,
									Validators:          []validator.String{stringvalidator.LengthAtLeast(1)},
									MarkdownDescription: "Unique URI identifying the extension.",
								},
								"description": schema.StringAttribute{
									Optional:            true,
									Computed:            true,
									PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
									MarkdownDescription: "Description of the agent, skill, or extension.",
								},
								"required": schema.BoolAttribute{
									Optional:            true,
									Computed:            true,
									PlanModifiers:       []planmodifier.Bool{boolplanmodifier.UseStateForUnknown()},
									MarkdownDescription: "Whether clients must support this extension.",
								},
								"params_json": schema.StringAttribute{
									Optional: true,
									Computed: true, PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
									CustomType:          jsontypes.NormalizedType{},
									Sensitive:           true,
									MarkdownDescription: "JSON object containing arbitrary extension parameters.",
									Validators:          []validator.String{stringvalidator.RegexMatches(regexp.MustCompile(`^\s*\{`), "must be a JSON object")},
								},
							},
						},
					},
				},
			},
			"skills": schema.ListNestedAttribute{
				Optional: true,
				Computed: true, PlanModifiers: []planmodifier.List{listplanmodifier.UseStateForUnknown()},
				MarkdownDescription: "Skills advertised by the agent.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Required:            true,
							Validators:          []validator.String{stringvalidator.LengthAtLeast(1)},
							MarkdownDescription: "Unique identifier for the skill.",
						},
						"display_name": schema.StringAttribute{
							Required:            true,
							Validators:          []validator.String{stringvalidator.LengthAtLeast(1)},
							MarkdownDescription: "Human-readable name of the agent or skill.",
						},
						"description": schema.StringAttribute{
							Optional:            true,
							Computed:            true,
							PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
							MarkdownDescription: "Description of the agent, skill, or extension.",
						},
						"tags": schema.SetAttribute{
							Optional:            true,
							Computed:            true,
							PlanModifiers:       []planmodifier.Set{setplanmodifier.UseStateForUnknown()},
							ElementType:         types.StringType,
							MarkdownDescription: "Tags describing the skill.",
						},
						"examples": schema.SetAttribute{
							Optional: true,
							Computed: true, PlanModifiers: []planmodifier.Set{setplanmodifier.UseStateForUnknown()},
							ElementType:         types.StringType,
							MarkdownDescription: "Example prompts or uses of the skill.",
						},
						"input_modes": schema.SetAttribute{
							Optional: true,
							Computed: true, PlanModifiers: []planmodifier.Set{setplanmodifier.UseStateForUnknown()},
							ElementType:         types.StringType,
							MarkdownDescription: "Content types accepted by this skill.",
						},
						"output_modes": schema.SetAttribute{
							Optional: true,
							Computed: true, PlanModifiers: []planmodifier.Set{setplanmodifier.UseStateForUnknown()},
							ElementType:         types.StringType,
							MarkdownDescription: "Content types produced by this skill.",
						},
					},
				},
			},
			"security_schemes_json": schema.StringAttribute{
				Optional: true,
				Computed: true, PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
				CustomType:          jsontypes.NormalizedType{},
				Sensitive:           false,
				MarkdownDescription: "JSON object mapping arbitrary security scheme names to definitions.",
				Validators:          []validator.String{stringvalidator.RegexMatches(regexp.MustCompile(`^\s*\{`), "must be a JSON object")},
			},
			"security_json": schema.StringAttribute{
				Optional: true,
				Computed: true, PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
				CustomType:          jsontypes.NormalizedType{},
				Sensitive:           false,
				MarkdownDescription: "JSON array of security requirement objects.",
				Validators:          []validator.String{stringvalidator.RegexMatches(regexp.MustCompile(`^\s*\[`), "must be a JSON array")},
			},
			"created_by":              schema.StringAttribute{Computed: true, MarkdownDescription: "Identity that created the manifest, as returned by Graph."},
			"created_date_time":       schema.StringAttribute{Computed: true, MarkdownDescription: "Creation time returned by Graph."},
			"last_modified_date_time": schema.StringAttribute{Computed: true, MarkdownDescription: "Last modification time returned by Graph."},
			"timeouts":                commonschema.ResourceTimeouts(ctx),
		},
	}
}
