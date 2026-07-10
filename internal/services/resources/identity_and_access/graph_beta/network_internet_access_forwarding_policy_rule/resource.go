package graphBetaNetworkInternetAccessForwardingPolicyRule

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	planmodifiers "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/plan_modifiers"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
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
	ResourceName  = "microsoft365_graph_beta_identity_and_access_network_internet_access_forwarding_policy_rule"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180

	ruleTypeFQDN      = "fqdn"
	ruleTypeIPAddress = "ip_address"
	ruleTypeIPRange   = "ip_range"
	ruleTypeIPSubnet  = "ip_subnet"
)

var (
	_ resource.Resource                = &NetworkInternetAccessForwardingPolicyRuleResource{}
	_ resource.ResourceWithConfigure   = &NetworkInternetAccessForwardingPolicyRuleResource{}
	_ resource.ResourceWithImportState = &NetworkInternetAccessForwardingPolicyRuleResource{}
	_ resource.ResourceWithIdentity    = &NetworkInternetAccessForwardingPolicyRuleResource{}
)

func NewNetworkInternetAccessForwardingPolicyRuleResource() resource.Resource {
	return &NetworkInternetAccessForwardingPolicyRuleResource{
		ReadPermissions: []string{
			"NetworkAccess.Read.All",
		},
		WritePermissions: []string{
			"NetworkAccess.ReadWrite.All",
		},
		ResourcePath: "/networkAccess/forwardingPolicies/{forwardingPolicyId}/policyRules",
	}
}

type NetworkInternetAccessForwardingPolicyRuleResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

func (r *NetworkInternetAccessForwardingPolicyRuleResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = ResourceName
}

func (r *NetworkInternetAccessForwardingPolicyRuleResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, ResourceName)
}

func (r *NetworkInternetAccessForwardingPolicyRuleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.Split(req.ID, "/")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		resp.Diagnostics.AddError("Invalid import ID", fmt.Sprintf("Expected import ID in the format {forwarding_policy_id}/{rule_id}, got %q.", req.ID))
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("forwarding_policy_id"), parts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), parts[1])...)
}

func (r *NetworkInternetAccessForwardingPolicyRuleResource) IdentitySchema(ctx context.Context, req resource.IdentitySchemaRequest, resp *resource.IdentitySchemaResponse) {
	resp.IdentitySchema = identityschema.Schema{
		Attributes: map[string]identityschema.Attribute{
			"id": identityschema.StringAttribute{RequiredForImport: true},
		},
	}
}

func (r *NetworkInternetAccessForwardingPolicyRuleResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages Microsoft Entra Global Secure Access Internet Access forwarding policy rules using Microsoft Graph beta `/networkAccess/forwardingPolicies/{forwardingPolicyId}/policyRules`.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The forwarding policy rule ID.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
			},
			"forwarding_policy_id": schema.StringAttribute{
				MarkdownDescription: "The forwarding policy ID that owns this rule.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(regexp.MustCompile(constants.GuidRegex), "must be a valid UUID"),
				},
				PlanModifiers: []planmodifier.String{
					planmodifiers.RequiresReplaceString(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The rule name.",
				Required:            true,
			},
			"action": schema.StringAttribute{
				MarkdownDescription: "The forwarding action. Possible values are `forward` and `bypass`.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("forward", "bypass"),
				},
			},
			"rule_type": schema.StringAttribute{
				MarkdownDescription: "The destination rule type. Possible values are `fqdn`, `ip_address`, `ip_range`, and `ip_subnet`.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf(ruleTypeFQDN, ruleTypeIPAddress, ruleTypeIPRange, ruleTypeIPSubnet),
				},
			},
			"client_fallback_action": schema.StringAttribute{
				MarkdownDescription: "The Graph-computed client fallback action.",
				Computed:            true,
			},
			"ports": schema.SetAttribute{
				MarkdownDescription: "Network ports matched by the rule. Observed Internet Access rules use values such as `80` and `443`.",
				Required:            true,
				ElementType:         types.StringType,
				Validators: []validator.Set{
					setvalidator.SizeAtLeast(1),
				},
			},
			"protocol": schema.StringAttribute{
				MarkdownDescription: "Network protocol. Validated values are `tcp` and `udp` until additional Graph contract probing confirms more values.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("tcp", "udp"),
				},
			},
			"destinations": schema.ListNestedAttribute{
				MarkdownDescription: "Destinations for the rule. FQDN, IP address, IP range, and CIDR/IP subnet shapes are supported.",
				Required:            true,
				Validators: []validator.List{
					listvalidator.SizeAtLeast(1),
				},
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"type": schema.StringAttribute{
							MarkdownDescription: "Destination type. Possible values are `fqdn`, `ip_address`, `ip_range`, and `ip_subnet`.",
							Required:            true,
							Validators: []validator.String{
								stringvalidator.OneOf(ruleTypeFQDN, ruleTypeIPAddress, ruleTypeIPRange, ruleTypeIPSubnet),
							},
						},
						"value": schema.StringAttribute{
							MarkdownDescription: "Destination value for `fqdn`, `ip_address`, and `ip_subnet`. Use CIDR notation for `ip_subnet`.",
							Optional:            true,
						},
						"begin_address": schema.StringAttribute{
							MarkdownDescription: "Beginning IP address for `ip_range`.",
							Optional:            true,
						},
						"end_address": schema.StringAttribute{
							MarkdownDescription: "Ending IP address for `ip_range`.",
							Optional:            true,
						},
					},
				},
			},
			"timeouts": commonschema.ResourceTimeouts(ctx),
		},
	}
}
