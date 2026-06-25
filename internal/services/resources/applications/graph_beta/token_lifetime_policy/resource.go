package graphBetaApplicationsTokenLifetimePolicy

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/identityschema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	ResourceName  = "microsoft365_graph_beta_applications_token_lifetime_policy"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	_ resource.Resource                = &TokenLifetimePolicyResource{}
	_ resource.ResourceWithConfigure   = &TokenLifetimePolicyResource{}
	_ resource.ResourceWithImportState = &TokenLifetimePolicyResource{}
	_ resource.ResourceWithIdentity    = &TokenLifetimePolicyResource{}
)

func NewTokenLifetimePolicyResource() resource.Resource {
	return &TokenLifetimePolicyResource{
		ReadPermissions: []string{
			"Policy.Read.All",
		},
		WritePermissions: []string{
			"Policy.ReadWrite.ApplicationConfiguration",
		},
	}
}

type TokenLifetimePolicyResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
}

func (r *TokenLifetimePolicyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = ResourceName
}

func (r *TokenLifetimePolicyResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, ResourceName)
}

func (r *TokenLifetimePolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *TokenLifetimePolicyResource) IdentitySchema(ctx context.Context, req resource.IdentitySchemaRequest, resp *resource.IdentitySchemaResponse) {
	resp.IdentitySchema = identityschema.Schema{
		Attributes: map[string]identityschema.Attribute{
			"id": identityschema.StringAttribute{
				RequiredForImport: true,
			},
		},
	}
}

func (r *TokenLifetimePolicyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages token lifetime policies in Microsoft Entra ID using the `/policies/tokenLifetimePolicies` endpoint. " +
			"Token lifetime policies control how long access tokens, ID tokens, and SAML 1.1/2.0 tokens issued for applications are valid. " +
			"You can set token lifetimes for all apps in your organization, for a multi-tenant application, or for a specific service principal. " +
			"Only one token lifetime policy can be assigned to a service principal at a time.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique identifier of the token lifetime policy.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"display_name": schema.StringAttribute{
				MarkdownDescription: "The display name of the token lifetime policy.",
				Required:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The description of the token lifetime policy.",
				Optional:            true,
				Computed:            true,
			},
			"definition": schema.ListAttribute{
				MarkdownDescription: "A JSON string collection that defines the token lifetime policy rules and settings. " +
					"The collection must contain exactly one JSON string. " +
					"For details on the JSON structure and configurable properties, see the " +
					"[Configure token lifetimes](https://learn.microsoft.com/en-us/entra/identity-platform/configure-token-lifetimes) " +
					"documentation.",
				ElementType: types.StringType,
				Required:    true,
			},
			"is_organization_default": schema.BoolAttribute{
				MarkdownDescription: "If `true`, this is the default policy for the organization. There can only be one organization default token lifetime policy. " +
					"Defaults to `false`.",
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(false),
			},
			"deleted_date_time": schema.StringAttribute{
				MarkdownDescription: "The date and time the policy was deleted. Null if the policy has not been deleted.",
				Computed:            true,
			},
			"timeouts": commonschema.ResourceTimeouts(ctx),
		},
	}
}
