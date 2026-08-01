package graphBetaServicePrincipalTokenSigningCertificate

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
	"github.com/hashicorp/terraform-plugin-framework/resource/identityschema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	ResourceName  = "microsoft365_graph_beta_applications_service_principal_token_signing_certificate"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	// Basic resource interface (CRUD operations)
	_ resource.Resource = &ServicePrincipalTokenSigningCertificateResource{}

	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &ServicePrincipalTokenSigningCertificateResource{}

	// Enables import functionality
	_ resource.ResourceWithImportState = &ServicePrincipalTokenSigningCertificateResource{}

	// Enables identity schema for list resource support
	_ resource.ResourceWithIdentity = &ServicePrincipalTokenSigningCertificateResource{}
)

func NewServicePrincipalTokenSigningCertificateResource() resource.Resource {
	return &ServicePrincipalTokenSigningCertificateResource{
		ReadPermissions: []string{
			"Application.Read.All",
			"Directory.Read.All",
		},
		WritePermissions: []string{
			"Application.ReadWrite.All",
		},
		ResourcePath: "/servicePrincipals",
	}
}

type ServicePrincipalTokenSigningCertificateResource struct {
	client *msgraphbetasdk.GraphServiceClient
	// httpClient is used for the credential-removal PATCH on delete, which must
	// round-trip the raw credential JSON untouched (see constructDeleteBody).
	httpClient       *client.AuthenticatedHTTPClient
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

// Metadata returns the resource type name.
func (r *ServicePrincipalTokenSigningCertificateResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = ResourceName
}

// Configure sets the clients for the resource.
func (r *ServicePrincipalTokenSigningCertificateResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, ResourceName)
	r.httpClient = client.SetGraphBetaHTTPClientForResource(ctx, req, resp, ResourceName)
}

// ImportState imports the resource state.
// Import format: service_principal_id/key_id
func (r *ServicePrincipalTokenSigningCertificateResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.SplitN(req.ID, "/", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			fmt.Sprintf("Expected import ID in format 'service_principal_id/key_id', got: %s", req.ID),
		)
		return
	}

	// Graph emits lowercase UUIDs; normalize key_id so keyId lookups match exactly.
	// service_principal_id keeps its given casing: it is a RequiresReplace attribute,
	// so rewriting it would produce a spurious replacement plan for configurations
	// that reference the service principal with an uppercase GUID literal.
	servicePrincipalID := parts[0]
	keyId := strings.ToLower(parts[1])

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("service_principal_id"), servicePrincipalID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("key_id"), keyId)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), servicePrincipalID+"/"+keyId)...)
}

// IdentitySchema defines the identity schema for this resource, used by list operations to uniquely identify instances
func (r *ServicePrincipalTokenSigningCertificateResource) IdentitySchema(ctx context.Context, req resource.IdentitySchemaRequest, resp *resource.IdentitySchemaResponse) {
	resp.IdentitySchema = identityschema.Schema{
		Attributes: map[string]identityschema.Attribute{
			"id": identityschema.StringAttribute{
				RequiredForImport: true,
			},
		},
	}
}

// Schema defines the schema for the resource.
func (r *ServicePrincipalTokenSigningCertificateResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a token signing certificate for a Microsoft Entra Service Principal, used for SAML-based single sign-on. " +
			"Microsoft Entra generates a self-signed certificate and adds the corresponding signing and verification key credentials " +
			"(plus a password credential) to the service principal.\n\n" +
			"The certificate is immutable — changing any attribute forces recreation. For zero-downtime rotation, use " +
			"`create_before_destroy` together with the `microsoft365_graph_beta_applications_service_principal_preferred_token_signing_key_thumbprint` " +
			"resource, which activates the certificate and updates in place when the thumbprint changes.\n\n" +
			"For more information, see the [Microsoft Graph API documentation](https://learn.microsoft.com/en-us/graph/api/serviceprincipal-addtokensigningcertificate?view=graph-rest-beta).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique identifier for this resource, in the format `service_principal_id/key_id`.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"service_principal_id": schema.StringAttribute{
				MarkdownDescription: "The object ID (UUID) of the service principal for which to generate the token signing certificate.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(constants.GuidRegex),
						"must be a valid GUID in the format 00000000-0000-0000-0000-000000000000",
					),
				},
			},
			"display_name": schema.StringAttribute{
				MarkdownDescription: "The friendly name of the certificate, in the format `CN=<name>`. " +
					"Defaults to `CN=Microsoft Azure Federated SSO Certificate` when omitted. Changing this forces recreation.",
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplaceIfConfigured(),
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^CN=`),
						"must start with 'CN='",
					),
				},
			},
			"end_date_time": schema.StringAttribute{
				MarkdownDescription: "The end date and time of the certificate validity period, in RFC3339 / ISO 8601 UTC format " +
					"(e.g. `2028-01-01T14:59:59Z`). Defaults to three years from creation when omitted. Changing this forces recreation.",
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplaceIfConfigured(),
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(\.\d+)?(Z|[+-]\d{2}:\d{2})$`),
						"must be a valid RFC3339 date-time, e.g. 2028-01-01T14:59:59Z",
					),
				},
			},
			"key_id": schema.StringAttribute{
				MarkdownDescription: "The unique identifier (UUID) of the signing key credential. Read-only.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"start_date_time": schema.StringAttribute{
				MarkdownDescription: "The start date and time of the certificate validity period, in RFC3339 / ISO 8601 UTC format. Read-only.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"thumbprint": schema.StringAttribute{
				MarkdownDescription: "The SHA-1 thumbprint of the certificate. Reference this from the " +
					"`microsoft365_graph_beta_applications_service_principal_preferred_token_signing_key_thumbprint` resource to activate the certificate. Read-only.",
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"value": schema.StringAttribute{
				MarkdownDescription: "The base64-encoded public certificate, captured from the creation response. Read-only. " +
					"May be unset after import: Microsoft Graph does not return the signing key credential's key material on reads.",
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"timeouts": commonschema.ResourceTimeouts(ctx),
		},
	}
}
