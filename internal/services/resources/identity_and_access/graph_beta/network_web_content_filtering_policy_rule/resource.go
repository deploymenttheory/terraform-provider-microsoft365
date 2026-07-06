package graphBetaNetworkWebContentFilteringPolicyRule

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	planmodifiers "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/plan_modifiers"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/resourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/identityschema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	ResourceName  = "microsoft365_graph_beta_identity_and_access_network_web_content_filtering_policy_rule"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	_ resource.Resource                     = &NetworkWebContentFilteringPolicyRuleResource{}
	_ resource.ResourceWithConfigure        = &NetworkWebContentFilteringPolicyRuleResource{}
	_ resource.ResourceWithImportState      = &NetworkWebContentFilteringPolicyRuleResource{}
	_ resource.ResourceWithIdentity         = &NetworkWebContentFilteringPolicyRuleResource{}
	_ resource.ResourceWithConfigValidators = &NetworkWebContentFilteringPolicyRuleResource{}
)

func NewNetworkWebContentFilteringPolicyRuleResource() resource.Resource {
	return &NetworkWebContentFilteringPolicyRuleResource{
		ReadPermissions: []string{
			"NetworkAccess.Read.All",
		},
		WritePermissions: []string{
			"NetworkAccess.ReadWrite.All",
		},
		ResourcePath: "/networkaccess/webFilteringPolicies/{webContentFilteringPolicyId}/policyRules",
	}
}

type NetworkWebContentFilteringPolicyRuleResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

func (r *NetworkWebContentFilteringPolicyRuleResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = ResourceName
}

func (r *NetworkWebContentFilteringPolicyRuleResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, ResourceName)
}

func (r *NetworkWebContentFilteringPolicyRuleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	idParts := strings.Split(req.ID, "/")
	if len(idParts) != 2 {
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			fmt.Sprintf("Expected import ID in format 'webContentFilteringPolicyId/ruleId', got: %s", req.ID),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("web_content_filtering_policy_id"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), idParts[1])...)
}

func (r *NetworkWebContentFilteringPolicyRuleResource) IdentitySchema(ctx context.Context, req resource.IdentitySchemaRequest, resp *resource.IdentitySchemaResponse) {
	resp.IdentitySchema = identityschema.Schema{
		Attributes: map[string]identityschema.Attribute{
			"id": identityschema.StringAttribute{
				RequiredForImport: true,
			},
		},
	}
}

func (r *NetworkWebContentFilteringPolicyRuleResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{
		resourcevalidator.AtLeastOneOf(
			path.MatchRoot("urls_or_fqdns"),
			path.MatchRoot("web_categories"),
		),
	}
}

func (r *NetworkWebContentFilteringPolicyRuleResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages rules for Microsoft Entra Global Secure Access web content filtering policies using the Microsoft Graph beta `/networkaccess/webFilteringPolicies/{id}/policyRules` endpoint observed from the Entra portal.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique identifier for the web content filtering policy rule.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
			},
			"web_content_filtering_policy_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the web content filtering policy that owns this rule.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.RequiresReplaceString(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the web content filtering rule.",
				Required:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Optional description of the web content filtering rule. Maximum length is 8192 characters.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(8192),
				},
			},
			"priority": schema.Int64Attribute{
				MarkdownDescription: "The rule priority. Lower numbers are evaluated before higher numbers. The Entra portal accepts values from 100 to 65000.",
				Required:            true,
				Validators: []validator.Int64{
					int64validator.Between(100, 65000),
				},
			},
			"action": schema.StringAttribute{
				MarkdownDescription: "The action for matching traffic. Possible values are `allow` and `block`.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("allow", "block"),
				},
			},
			"status": schema.StringAttribute{
				MarkdownDescription: "The rule status. Possible values are `enabled` and `disabled`.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("enabled", "disabled"),
				},
			},
			"urls_or_fqdns": schema.SetAttribute{
				MarkdownDescription: "URL or FQDN destination patterns for the rule, for example `www.MySite.com`, `www.MySite.com/a/b/c`, `www.MySite.com/a/*`, or `*.mysite.com`. Use `*` to match any URL or FQDN. At least one of `urls_or_fqdns` or `web_categories` must be specified; both can be set on the same rule. If set, this attribute must contain at least one value. The Entra portal shows URL/FQDN destinations as comma-separated text, while Microsoft Graph stores them as a values array; Terraform follows the Graph shape with one set element per destination.",
				ElementType:         types.StringType,
				Optional:            true,
				Validators: []validator.Set{
					setvalidator.SizeAtLeast(1),
				},
			},
			"web_categories": schema.SetAttribute{
				MarkdownDescription: "Web category IDs for the rule, for example `AlcoholAndTobacco`. At least one of `urls_or_fqdns` or `web_categories` must be specified; both can be set on the same rule. If set, this attribute must contain at least one value. Category IDs are passed through to Microsoft Graph unchanged.",
				ElementType:         types.StringType,
				Optional:            true,
				Validators: []validator.Set{
					setvalidator.SizeAtLeast(1),
				},
			},
			"http_methods": schema.SetAttribute{
				MarkdownDescription: "HTTP methods that must match the rule. The Entra portal sends these as comma-separated lowercase values.",
				ElementType:         types.StringType,
				Optional:            true,
				Validators: []validator.Set{
					setvalidator.ValueStringsAre(stringvalidator.OneOf("get", "post", "put", "patch", "delete")),
				},
			},
			"session_types": schema.SetAttribute{
				MarkdownDescription: "Session types that must match the rule. Possible values are `user` and `agent`.",
				ElementType:         types.StringType,
				Optional:            true,
				Validators: []validator.Set{
					setvalidator.ValueStringsAre(stringvalidator.OneOf("user", "agent")),
				},
			},
			"custom_headers": schema.ListNestedAttribute{
				MarkdownDescription: "Custom request headers to add for allow rules. Microsoft Graph accepts these only when `action` is `allow`; the Entra portal serializes them as `action.headerSettings.modifications`.",
				Optional:            true,
				Validators: []validator.List{
					listvalidator.SizeAtMost(10),
				},
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"header_name": schema.StringAttribute{
							MarkdownDescription: "The custom header name.",
							Required:            true,
							Validators: []validator.String{
								stringvalidator.LengthAtMost(128),
								stringvalidator.RegexMatches(regexp.MustCompile(`^[!#$%&'*+\-.^_`+"`"+`|~0-9A-Za-z]+$`), "must be a valid HTTP header name"),
							},
						},
						"header_value": schema.StringAttribute{
							MarkdownDescription: "The custom header value.",
							Required:            true,
							Validators: []validator.String{
								stringvalidator.LengthAtMost(2048),
								stringvalidator.RegexMatches(regexp.MustCompile(`^[\x20-\x7E]*$`), "must contain only printable ASCII characters"),
							},
						},
					},
				},
			},
			"timeouts": commonschema.ResourceTimeouts(ctx),
		},
	}
}
