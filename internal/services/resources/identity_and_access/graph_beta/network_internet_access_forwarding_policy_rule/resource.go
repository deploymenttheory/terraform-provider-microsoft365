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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
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

	_ resource.ResourceWithConfigValidators = &NetworkInternetAccessForwardingPolicyRuleResource{}
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

func (r *NetworkInternetAccessForwardingPolicyRuleResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{destinationsValidator{}}
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
				MarkdownDescription: "The rule name. Graph does not allow updating this value after creation; changing it requires replacement.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.RequiresReplaceString(),
				},
			},
			"action": schema.StringAttribute{
				MarkdownDescription: "The forwarding action. Possible values are `forward` and `bypass`. Graph requires `forward` for acquire policies and `bypass` for bypass policies; changing this value requires replacement.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("forward", "bypass"),
				},
				PlanModifiers: []planmodifier.String{
					planmodifiers.RequiresReplaceString(),
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
				MarkdownDescription: "What the client does with matching traffic when it cannot reach Global Secure Access. Possible values are `default`, `bypass`, and `block`. When omitted on create, Graph sets `default`; removing it later keeps the existing value. Graph does not allow updating it after creation, so changing it requires replacement.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("default", "bypass", "block"),
				},
				// The framework's UseStateForUnknown leaves a value that is unknown in
				// config alone, so RequiresReplaceIfConfigured still sees the change.
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplaceIfConfigured(),
				},
			},
			"ports": schema.SetAttribute{
				MarkdownDescription: "Network ports matched by the rule. Each entry is a single port such as `443` or an inclusive range such as `1000-2000`, between 0 and 65535.",
				Required:            true,
				ElementType:         types.StringType,
				Validators: []validator.Set{
					setvalidator.SizeAtLeast(1),
					setvalidator.ValueStringsAre(portValidator{}),
				},
			},
			"protocol": schema.StringAttribute{
				MarkdownDescription: "Network protocol. Possible values are `tcp` and `udp`; a rule matches one protocol.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("tcp", "udp"),
				},
			},
			"destinations": schema.ListNestedAttribute{
				MarkdownDescription: "Destinations for the rule. Every destination must have the same `type` as `rule_type`. IP addresses, ranges and subnets must be IPv4; Graph rejects IPv6 for Internet Access rules.",
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
							MarkdownDescription: "Destination value for `fqdn`, `ip_address`, and `ip_subnet`. Use CIDR notation for `ip_subnet`. Not used for `ip_range`.",
							Optional:            true,
						},
						"begin_address": schema.StringAttribute{
							MarkdownDescription: "Beginning IPv4 address for `ip_range`.",
							Optional:            true,
						},
						"end_address": schema.StringAttribute{
							MarkdownDescription: "Ending IPv4 address for `ip_range`. Must not be lower than `begin_address`.",
							Optional:            true,
						},
					},
				},
			},
			"timeouts": commonschema.ResourceTimeouts(ctx),
		},
	}
}
