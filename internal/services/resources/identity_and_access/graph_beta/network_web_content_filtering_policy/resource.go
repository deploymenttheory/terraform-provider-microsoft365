package graphBetaNetworkWebContentFilteringPolicy

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	planmodifiers "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/plan_modifiers"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/identityschema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	ResourceName  = "microsoft365_graph_beta_identity_and_access_network_web_content_filtering_policy"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	_ resource.Resource                = &NetworkWebContentFilteringPolicyResource{}
	_ resource.ResourceWithConfigure   = &NetworkWebContentFilteringPolicyResource{}
	_ resource.ResourceWithImportState = &NetworkWebContentFilteringPolicyResource{}
	_ resource.ResourceWithIdentity    = &NetworkWebContentFilteringPolicyResource{}
)

func NewNetworkWebContentFilteringPolicyResource() resource.Resource {
	return &NetworkWebContentFilteringPolicyResource{
		ReadPermissions: []string{
			"NetworkAccess.Read.All",
		},
		WritePermissions: []string{
			"NetworkAccess.ReadWrite.All",
		},
		ResourcePath: "/networkaccess/webFilteringPolicies",
	}
}

type NetworkWebContentFilteringPolicyResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

func (r *NetworkWebContentFilteringPolicyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = ResourceName
}

func (r *NetworkWebContentFilteringPolicyResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, ResourceName)
}

func (r *NetworkWebContentFilteringPolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *NetworkWebContentFilteringPolicyResource) IdentitySchema(ctx context.Context, req resource.IdentitySchemaRequest, resp *resource.IdentitySchemaResponse) {
	resp.IdentitySchema = identityschema.Schema{
		Attributes: map[string]identityschema.Attribute{
			"id": identityschema.StringAttribute{
				RequiredForImport: true,
			},
		},
	}
}

func (r *NetworkWebContentFilteringPolicyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages Microsoft Entra Global Secure Access web content filtering policies using the Microsoft Graph beta `/networkaccess/webFilteringPolicies` endpoint observed from the Entra portal.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique identifier for the web content filtering policy.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the web content filtering policy.",
				Required:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Optional description of the web content filtering policy. Maximum length is 8192 characters.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(8192),
				},
			},
			"default_action": schema.StringAttribute{
				MarkdownDescription: "The default action applied by the policy. Possible values are `allow` and `block`.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("allow", "block"),
				},
			},
			"timeouts": commonschema.ResourceTimeouts(ctx),
		},
	}
}
