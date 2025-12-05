package graphBetaAgentIdentityBlueprintKeyCredential

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
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
	ResourceName  = "microsoft365_graph_beta_agents_agent_identity_blueprint_rotate_certificate_credential"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	// Basic resource interface (CRUD operations)
	_ resource.Resource = &AgentIdentityBlueprintKeyCredentialResource{}

	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &AgentIdentityBlueprintKeyCredentialResource{}

	// Note: ImportState is NOT supported for this resource because the key
	// cannot be retrieved from the API after creation.
)

func NewAgentIdentityBlueprintKeyCredentialResource() resource.Resource {
	return &AgentIdentityBlueprintKeyCredentialResource{
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

type AgentIdentityBlueprintKeyCredentialResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

// Metadata returns the resource type name.
func (r *AgentIdentityBlueprintKeyCredentialResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = ResourceName
}

// Configure sets the client for the resource.
func (r *AgentIdentityBlueprintKeyCredentialResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, ResourceName)
}

// Schema returns the schema for the resource.
func (r *AgentIdentityBlueprintKeyCredentialResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a key credential for an Agent Identity Blueprint in Microsoft Entra ID using the `/applications` endpoint. " +
			"This resource adds a key credential to an agentIdentityBlueprint using the " +
			"[addKey](https://learn.microsoft.com/en-us/graph/api/agentidentityblueprint-addkey?view=graph-rest-beta) API. " +
			"This method, along with removeKey, can be used to automate rolling its expiring keys.\n\n" +
			"**Important:** You should only provide the public key value when adding a certificate credential. " +
			"Adding a private key certificate risks compromising the application.\n\n" +
			"As part of the request validation for this method, a proof of possession of an existing key is verified before the action can be performed.\n\n" +
			"**Note:** Agent identity blueprints that don't have any existing valid certificates (no certificates have been added yet, " +
			"or all certificates have expired), won't be able to use this service action. Use the Update agent identity blueprint operation instead.\n\n" +
			"**Note:** This resource does not support import because the key value cannot be retrieved from the API after creation.",
		Attributes: map[string]schema.Attribute{
			"blueprint_id": schema.StringAttribute{
				MarkdownDescription: "The unique identifier (Object ID) of the agent identity blueprint to add the key credential to.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"key": schema.StringAttribute{
				MarkdownDescription: "The certificate's raw data in byte array converted to Base64 string. " +
					"From a .cer certificate, you can read the key using the Convert.ToBase64String() method. " +
					"For more information, see [Get the certificate key](https://learn.microsoft.com/en-us/graph/applications-how-to-add-certificate#get-the-certificate-key). " +
					"**Important:** Only provide the public key value. Adding a private key risks compromising the application.",
				Required:  true,
				Sensitive: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "The type of key credential. Supported key types are: " +
					"`AsymmetricX509Cert` (the usage must be `Verify`), " +
					"`X509CertAndPassword` (the usage must be `Sign`).",
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf("AsymmetricX509Cert", "X509CertAndPassword"),
				},
			},
			"usage": schema.StringAttribute{
				MarkdownDescription: "A string that describes the purpose for which the key can be used. " +
					"Possible values are: `None`, `Verify`, `PairwiseIdentifier`, `Delegation`, `Decrypt`, `Encrypt`, `HashedIdentifier`, `SelfSignedTls`, or `Sign`. " +
					"If usage is `Sign`, the type should be `X509CertAndPassword`, and the `password_secret_text` for signing should be defined.",
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf("None", "Verify", "PairwiseIdentifier", "Delegation", "Decrypt", "Encrypt", "HashedIdentifier", "SelfSignedTls", "Sign"),
				},
			},
			"proof": schema.StringAttribute{
				MarkdownDescription: "A self-signed JWT token used as a proof of possession of the existing keys. " +
					"This JWT token must be signed using the private key of one of the application's existing valid certificates. " +
					"The token should contain the following claims:\n" +
					"- **aud**: Audience needs to be `00000002-0000-0000-c000-000000000000`.\n" +
					"- **iss**: Issuer needs to be the ID of the application that initiates the request.\n" +
					"- **nbf**: Not before time.\n" +
					"- **exp**: Expiration time should be the value of `nbf` + 10 minutes.\n\n" +
					"For steps to generate this proof of possession token, see " +
					"[Generating proof of possession tokens for rolling keys](https://learn.microsoft.com/en-us/graph/application-rollkey-prooftoken).",
				Required:  true,
				Sensitive: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"display_name": schema.StringAttribute{
				MarkdownDescription: "Friendly name for the key. Maximum length is 90 characters.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(90),
				},
			},
			"start_date_time": schema.StringAttribute{
				MarkdownDescription: "The date and time at which the credential becomes valid. " +
					"The Timestamp type represents date and time information using ISO 8601 format and is always in UTC time. " +
					"For example, midnight UTC on Jan 1, 2014 is `2014-01-01T00:00:00Z`.",
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"end_date_time": schema.StringAttribute{
				MarkdownDescription: "The date and time at which the credential expires. " +
					"The Timestamp type represents date and time information using ISO 8601 format and is always in UTC time. " +
					"For example, midnight UTC on Jan 1, 2014 is `2014-01-01T00:00:00Z`.",
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"custom_key_identifier": schema.StringAttribute{
				MarkdownDescription: "A 40-character binary type that can be used to identify the credential. Optional. " +
					"When not provided in the payload, defaults to the thumbprint of the certificate.",
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"password_secret_text": schema.StringAttribute{
				MarkdownDescription: "The password for the key. Required only for keys of type `X509CertAndPassword`. " +
					"Set it to null otherwise.",
				Optional:  true,
				Sensitive: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"key_id": schema.StringAttribute{
				MarkdownDescription: "The unique identifier (GUID) for the key. Read-only.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"timeouts": commonschema.Timeouts(ctx),
		},
	}
}
