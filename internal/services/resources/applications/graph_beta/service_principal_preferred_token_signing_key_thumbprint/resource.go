package graphBetaServicePrincipalPreferredTokenSigningKeyThumbprint

import (
	"context"
	"regexp"

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
	ResourceName  = "microsoft365_graph_beta_applications_service_principal_preferred_token_signing_key_thumbprint"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	// Basic resource interface (CRUD operations)
	_ resource.Resource = &ServicePrincipalPreferredTokenSigningKeyThumbprintResource{}

	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &ServicePrincipalPreferredTokenSigningKeyThumbprintResource{}

	// Enables import functionality
	_ resource.ResourceWithImportState = &ServicePrincipalPreferredTokenSigningKeyThumbprintResource{}

	// Enables identity schema for list resource support
	_ resource.ResourceWithIdentity = &ServicePrincipalPreferredTokenSigningKeyThumbprintResource{}
)

func NewServicePrincipalPreferredTokenSigningKeyThumbprintResource() resource.Resource {
	return &ServicePrincipalPreferredTokenSigningKeyThumbprintResource{
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

type ServicePrincipalPreferredTokenSigningKeyThumbprintResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

// Metadata returns the resource type name.
func (r *ServicePrincipalPreferredTokenSigningKeyThumbprintResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = ResourceName
}

// Configure sets the client for the resource.
func (r *ServicePrincipalPreferredTokenSigningKeyThumbprintResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, ResourceName)
}

// ImportState imports the resource state.
// Import format: service_principal_id (the thumbprint is populated by Read)
func (r *ServicePrincipalPreferredTokenSigningKeyThumbprintResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("service_principal_id"), req.ID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
}

// IdentitySchema defines the identity schema for this resource, used by list operations to uniquely identify instances
func (r *ServicePrincipalPreferredTokenSigningKeyThumbprintResource) IdentitySchema(ctx context.Context, req resource.IdentitySchemaRequest, resp *resource.IdentitySchemaResponse) {
	resp.IdentitySchema = identityschema.Schema{
		Attributes: map[string]identityschema.Attribute{
			"id": identityschema.StringAttribute{
				RequiredForImport: true,
			},
		},
	}
}

// Schema defines the schema for the resource.
func (r *ServicePrincipalPreferredTokenSigningKeyThumbprintResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages the preferred token signing key (`preferredTokenSigningKeyThumbprint`) of a Microsoft Entra Service Principal. " +
			"Setting this property activates a SAML token signing certificate for the service principal, which is required when configuring " +
			"SAML-based single sign-on. Destroying this resource clears the property, returning the service principal to Microsoft Entra's " +
			"automatic signing key selection; the certificate itself is not removed.\n\n" +
			"Note: Microsoft Graph does not validate that the thumbprint references an existing `keyCredential` on the service principal — " +
			"ensure the thumbprint comes from a token signing certificate that exists on the same service principal.\n\n" +
			"For more information, see the [Microsoft Graph API documentation](https://learn.microsoft.com/en-us/graph/api/serviceprincipal-update?view=graph-rest-beta).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique identifier for this resource, equal to `service_principal_id`.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"service_principal_id": schema.StringAttribute{
				MarkdownDescription: "The object ID (UUID) of the service principal.",
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
			"thumbprint": schema.StringAttribute{
				MarkdownDescription: "The SHA-1 thumbprint of the token signing certificate to activate, as a 40-character hexadecimal string " +
					"without colons or spaces. Changing this value updates the service principal in place, which supports certificate rotation " +
					"with `create_before_destroy` on the certificate resource.",
				Required: true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^[0-9a-fA-F]{40}$`),
						"must be a 40-character hexadecimal SHA-1 certificate thumbprint without colons or spaces",
					),
				},
			},
			"timeouts": commonschema.ResourceTimeouts(ctx),
		},
	}
}
