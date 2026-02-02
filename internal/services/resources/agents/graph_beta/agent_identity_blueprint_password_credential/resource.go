package graphBetaAgentIdentityBlueprintPasswordCredential

import (
	"context"
	"regexp"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	ResourceName  = "microsoft365_graph_beta_agents_agent_identity_blueprint_password_credential"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	// Basic resource interface (CRUD operations)
	_ resource.Resource = &AgentIdentityBlueprintPasswordCredentialResource{}

	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &AgentIdentityBlueprintPasswordCredentialResource{}

	// Note: ImportState is NOT supported for this resource because the secret_text
	// is only available at creation time and cannot be retrieved from the API.
)

func NewAgentIdentityBlueprintPasswordCredentialResource() resource.Resource {
	return &AgentIdentityBlueprintPasswordCredentialResource{
		ReadPermissions: []string{
			"AgentIdentityBlueprint.Read.All",
			"Application.Read.All",
			"Directory.Read.All",
		},
		WritePermissions: []string{
			"AgentIdentityBlueprint.AddRemoveCreds.All",
			"AgentIdentityBlueprint.ReadWrite.All",
			"Directory.ReadWrite.All",
		},
		ResourcePath: "/applications",
	}
}

type AgentIdentityBlueprintPasswordCredentialResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

// Metadata returns the resource type name.
func (r *AgentIdentityBlueprintPasswordCredentialResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = ResourceName
}

// Configure sets the client for the resource.
func (r *AgentIdentityBlueprintPasswordCredentialResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, ResourceName)
}

// Schema returns the schema for the resource.
func (r *AgentIdentityBlueprintPasswordCredentialResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a password credential for an Agent Identity Blueprint using the `/applications` endpoint. This resource is used to adds a strong password to an agentIdentityBlueprint using the [addPassword](https://learn.microsoft.com/en-us/graph/api/agentidentityblueprint-addpassword?view=graph-rest-beta) API.\n\n**Important:** The `secret_text` attribute contains the password and is only available at creation time. It cannot be retrieved after initial creation. Store this value securely.\n\n**Note:** This resource does not support import because the password cannot be retrieved from the API after creation..",
		Attributes: map[string]schema.Attribute{
			"blueprint_id": schema.StringAttribute{
				MarkdownDescription: "The unique identifier (Object ID) of the agent identity blueprint to add the password credential to. Required.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"display_name": schema.StringAttribute{
				MarkdownDescription: "Friendly name for the password credential. This helps identify the purpose of the credential. Optional but recommended.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(256),
				},
			},
			"start_date_time": schema.StringAttribute{
				MarkdownDescription: "The date and time at which the password becomes valid. The Timestamp type represents date and time information " +
					"using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z. " +
					"Optional. The default value is 'now'.",
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(constants.ISO8601DateTimeRegex),
						"must be a valid ISO 8601 datetime format (e.g., 2026-01-01T00:00:00Z)",
					),
				},
			},
			"end_date_time": schema.StringAttribute{
				MarkdownDescription: "The date and time at which the password expires represented using ISO 8601 format and is always in UTC time. " +
					"For example, midnight UTC on Jan 1, 2026 is 2026-01-01T00:00:00Z. Optional. The default value is `startDateTime + 2 years`.",
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(constants.ISO8601DateTimeRegex),
						"must be a valid ISO 8601 datetime format (e.g., 2026-01-01T00:00:00Z)",
					),
				},
			},
			"key_id": schema.StringAttribute{
				MarkdownDescription: "The unique identifier for the password credential. This is generated by Microsoft Entra ID and used to reference the credential for removal. Read-only.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"secret_text": schema.StringAttribute{
				MarkdownDescription: "The strong password generated by Microsoft Entra ID (16-64 characters). " +
					"**Important:** This value is only available at creation time and cannot be retrieved later. Store it securely.",
				Computed:  true,
				Sensitive: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"hint": schema.StringAttribute{
				MarkdownDescription: "A hint for the password (typically the first few characters). Read-only.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"custom_key_identifier": schema.StringAttribute{
				MarkdownDescription: "A custom key identifier for the password credential. Read-only.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"timeouts": commonschema.ResourceTimeouts(ctx),
		},
	}
}
