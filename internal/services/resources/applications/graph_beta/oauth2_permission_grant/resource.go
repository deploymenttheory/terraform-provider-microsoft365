package graphBetaOAuth2PermissionGrant

import (
	"context"
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/identityschema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	validate "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/validate/attribute"
)

const (
	ResourceName   = "microsoft365_graph_beta_applications_oauth2_permission_grant"
	CreateTimeout  = 180
	ReadTimeout    = 180
	UpdateTimeout  = 180
	DeleteTimeout  = 180
	maxScopeLength = 3850
)

var (
	_ resource.Resource                   = &OAuth2PermissionGrantResource{}
	_ resource.ResourceWithConfigure      = &OAuth2PermissionGrantResource{}
	_ resource.ResourceWithImportState    = &OAuth2PermissionGrantResource{}
	_ resource.ResourceWithIdentity       = &OAuth2PermissionGrantResource{}
	_ resource.ResourceWithValidateConfig = &OAuth2PermissionGrantResource{}
)

type OAuth2PermissionGrantResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
}

func NewOAuth2PermissionGrantResource() resource.Resource {
	return &OAuth2PermissionGrantResource{
		ReadPermissions:  []string{"Directory.Read.All"},
		WritePermissions: []string{"DelegatedPermissionGrant.ReadWrite.All"},
	}
}

func (r *OAuth2PermissionGrantResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = ResourceName
}

func (r *OAuth2PermissionGrantResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, ResourceName)
}

func (r *OAuth2PermissionGrantResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	if strings.TrimSpace(req.ID) == "" || strings.Contains(req.ID, "/") {
		resp.Diagnostics.AddError("Invalid import ID", "Supply the opaque Microsoft Graph delegated permission grant ID, without an API path.")
		return
	}
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *OAuth2PermissionGrantResource) IdentitySchema(_ context.Context, _ resource.IdentitySchemaRequest, resp *resource.IdentitySchemaResponse) {
	resp.IdentitySchema = identityschema.Schema{Attributes: map[string]identityschema.Attribute{
		"id": identityschema.StringAttribute{RequiredForImport: true},
	}}
}

func (r *OAuth2PermissionGrantResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	guid := stringvalidator.RegexMatches(regexp.MustCompile(constants.GuidRegex), "must be a valid UUID")
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages one delegated permission grant using `/oauth2PermissionGrants` on Microsoft Graph beta. A grant owns the complete set of scopes for a client service principal, resource service principal and consent identity. Import existing consent before managing it; do not manage individual scopes of the same grant in separate resources. Changing either service principal or the consent identity replaces the grant. Graph beta requires legacy creation timestamps but ignores their values; these do not provide grant-expiration controls.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The opaque Microsoft Graph grant ID. This is not a UUID.",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"client_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Object ID of the client service principal receiving delegated access. Use its service principal ID, not its application/client ID.",
				Validators:          []validator.String{guid},
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"resource_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Object ID of the resource service principal exposing the API.",
				Validators:          []validator.String{guid},
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"consent_type": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "`AllPrincipals` grants administrator consent for all users; `Principal` grants consent for the user identified by `principal_id`.",
				Validators:          []validator.String{stringvalidator.OneOf("AllPrincipals", "Principal")},
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"principal_id": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "User object ID. Required for `Principal` consent and omitted for `AllPrincipals`.",
				Validators: []validator.String{
					guid,
					validate.RequiredWhenEquals("consent_type", types.StringValue("Principal")),
					validate.StringMustBeEmptyWhenStringEquals("consent_type", "AllPrincipals", "principal_id must be omitted for AllPrincipals consent"),
				},
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"scopes": schema.SetAttribute{
				Required:            true,
				ElementType:         types.StringType,
				MarkdownDescription: "The complete set of delegated permission values, such as `User.Read`. Each element is one scope token, not a permission GUID. Order is ignored; the combined space-separated representation must not exceed 3,850 characters.",
				Validators: []validator.Set{
					setvalidator.ValueStringsAre(stringvalidator.RegexMatches(regexp.MustCompile(`^[^\s]+$`), "must contain one nonempty scope token without whitespace")),
				},
			},
			"timeouts": commonschema.ResourceTimeouts(ctx),
		},
	}
}

func (r *OAuth2PermissionGrantResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var scopes types.Set
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("scopes"), &scopes)...)
	if resp.Diagnostics.HasError() || scopes.IsNull() || scopes.IsUnknown() {
		return
	}
	for _, value := range scopes.Elements() {
		if value.IsUnknown() {
			return
		}
	}
	var values []string
	resp.Diagnostics.Append(scopes.ElementsAs(ctx, &values, false)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if utf8.RuneCountInString(strings.Join(values, " ")) > maxScopeLength {
		resp.Diagnostics.AddAttributeError(path.Root("scopes"), "Scope string exceeds Graph limit", "The complete space-separated scope string must not exceed 3,850 characters.")
	}
}
