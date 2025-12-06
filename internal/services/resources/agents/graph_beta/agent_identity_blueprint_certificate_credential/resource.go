package graphBetaAgentIdentityBlueprintCertificateCredential

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	attribute "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/validate/attribute"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
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
	ResourceName  = "microsoft365_graph_beta_agents_agent_identity_blueprint_certificate_credential"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	// Basic resource interface (CRUD operations)
	_ resource.Resource = &AgentIdentityBlueprintCertificateCredentialResource{}

	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &AgentIdentityBlueprintCertificateCredentialResource{}
)

// AgentIdentityBlueprintCertificateCredentialResource defines the resource implementation.
type AgentIdentityBlueprintCertificateCredentialResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ProviderTypeName string
	TypeName         string
	ReadPermissions  []string
	WritePermissions []string
}

// Metadata returns the resource type name.
func (r *AgentIdentityBlueprintCertificateCredentialResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = ResourceName
}

// Configure sets the client for the resource.
func (r *AgentIdentityBlueprintCertificateCredentialResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, ResourceName)
}

// ImportState imports the resource state - not supported for certificate credentials
// as the key value cannot be retrieved after creation.
func (r *AgentIdentityBlueprintCertificateCredentialResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resp.Diagnostics.AddError(
		"Import not supported",
		"Certificate credentials cannot be imported as the key value cannot be retrieved after creation.",
	)
}

// Schema defines the schema for the resource.
func (r *AgentIdentityBlueprintCertificateCredentialResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a certificate credential for an Agent Identity Blueprint application using the Microsoft Graph Beta API. " +
			"This resource uses PATCH on the application's keyCredentials property with OData type cast to microsoft.graph.agentIdentityBlueprint.",
		Attributes: map[string]schema.Attribute{
			"blueprint_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The unique identifier of the agent identity blueprint application.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"key": schema.StringAttribute{
				Required:            true,
				Sensitive:           true,
				MarkdownDescription: "The certificate data. The format depends on the `encoding` attribute: for `pem`, use `file(\"path/to/cert.pem\")`; for `base64`, use `filebase64(\"path/to/cert.der\")` or a base64 string; for `hex`, provide a hex-encoded string.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"encoding": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("pem"),
				MarkdownDescription: "Specifies the encoding used for the `key` value. Possible values are `pem`, `base64`, or `hex`. Defaults to `pem`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf("pem", "base64", "hex"),
				},
			},
			"type": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("AsymmetricX509Cert"),
				MarkdownDescription: "The type of key credential. Unlike the standard application keyCredentials endpoint which supports both `AsymmetricX509Cert` and `Symmetric`, the Agent Identity Blueprint API only supports `AsymmetricX509Cert` (X.509 certificates).",
				Validators: []validator.String{
					stringvalidator.OneOf("AsymmetricX509Cert"),
				},
			},
			"usage": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("Verify"),
				MarkdownDescription: "A string that describes the purpose for which the key can be used. Must be `Verify`.",
				Validators: []validator.String{
					stringvalidator.OneOf("Verify"),
				},
			},
			"display_name": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Friendly name for the certificate. Optional.",
			},
			"start_date_time": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The date and time at which the credential becomes valid. The timestamp type represents date and time information using ISO 8601 format and is always in UTC. For example, midnight UTC on Jan 1, 2024 is 2024-01-01T00:00:00Z. If not specified, defaults to the current time.",
			},
			"end_date_time": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The date and time at which the credential expires. The timestamp type represents date and time information using ISO 8601 format and is always in UTC. For example, midnight UTC on Jan 1, 2025 is 2025-01-01T00:00:00Z. Must be in the future.",
				Validators: []validator.String{
					attribute.RolloutDateTime(0, 3650), // Must be in the future, up to 10 years
				},
			},
			"replace_existing_certificates": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "When `true`, replaces all existing certificates on the application. When `false` (default), preserves existing certificates and adds the new one.",
			},
			"key_id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The unique identifier (GUID) for the key credential.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"custom_key_identifier": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "A 40-character binary type that can be used to identify the credential. Optional. When not provided in the payload, defaults to the thumbprint of the certificate.",
			},
			"thumbprint": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The thumbprint (SHA-1 hash) of the certificate.",
			},
			"timeouts": commonschema.Timeouts(ctx),
		},
	}
}

func NewAgentIdentityBlueprintCertificateCredentialResource() resource.Resource {
	return &AgentIdentityBlueprintCertificateCredentialResource{
		ProviderTypeName: "microsoft365",
		TypeName:         ResourceName,
		ReadPermissions: []string{
			"Application.Read.All",
		},
		WritePermissions: []string{
			"Application.ReadWrite.All",
		},
	}
}
