package graphBetaNetworkMCPPolicyRule

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
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
	planmodifiers "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/plan_modifiers"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
)

const (
	ResourceName  = "microsoft365_graph_beta_identity_and_access_network_mcp_policy_rule"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	_ resource.Resource                = &NetworkMCPPolicyRuleResource{}
	_ resource.ResourceWithConfigure   = &NetworkMCPPolicyRuleResource{}
	_ resource.ResourceWithImportState = &NetworkMCPPolicyRuleResource{}
	_ resource.ResourceWithIdentity    = &NetworkMCPPolicyRuleResource{}
)

// mcpConditionKeys maps each matching_conditions attribute to its Graph property under matchingConditions.destinations.
var mcpConditionKeys = map[string]string{
	"server_urls":         "serverUrls",
	"protocol_versions":   "protocolVersions",
	"insecure_connection": "insecureConnection",
	"missing_prm":         "missingPrm",
	"tool_matching":       "toolMatching",
	"resource_matching":   "resourceMatching",
	"prompt_matching":     "promptMatching",
}

// mcpMatchConditionAttributes is a string matching condition. Live POST/PATCH followed by GET preserved both order and
// duplicates for all five values arrays. Lists round-trip that configuration; sets would normalize it.
var mcpMatchConditionAttributes = map[string]schema.Attribute{
	"values": schema.ListAttribute{
		MarkdownDescription: "Values in API order. Empty arrays are preserved and are distinct from an omitted condition. Pattern evaluation is performed by the service; regex or wildcard behavior is not guaranteed by this resource.",
		Required:            true,
		ElementType:         types.StringType,
	},
	"match_type": schema.StringAttribute{
		MarkdownDescription: "Match operator: exactMatch, contains, notExactMatch, or notContains. Must be explicitly configured.",
		Required:            true,
		Validators:          []validator.String{stringvalidator.OneOf("exactMatch", "contains", "notExactMatch", "notContains")},
	},
}

// mcpPrimitiveConditionAttributes is a tool, resource or prompt matching condition.
var mcpPrimitiveConditionAttributes = map[string]schema.Attribute{
	"names": schema.SingleNestedAttribute{
		MarkdownDescription: "Primitive name matching.",
		Optional:            true,
		Attributes:          mcpMatchConditionAttributes,
	},
	"methods": schema.StringAttribute{
		MarkdownDescription: "API method flags, for example call or list,call for tools, read for resources, or get for prompts. Values are sent unchanged; unsupported combinations return API diagnostics.",
		Optional:            true,
		Validators:          []validator.String{stringvalidator.LengthAtLeast(1)},
	},
}

// mcpPropertyConditionAttribute is the insecure_connection and missing_prm condition.
var mcpPropertyConditionAttribute = schema.StringAttribute{
	MarkdownDescription: "Whether this property is required or excluded by the condition. Omission removes this constraint.",
	Optional:            true,
	Validators:          []validator.String{stringvalidator.OneOf("required", "excluded")},
}

// mcpPrimitiveConditionAttribute is the tool_matching, resource_matching and prompt_matching condition.
var mcpPrimitiveConditionAttribute = schema.SingleNestedAttribute{
	MarkdownDescription: "Primitive conditions applied to the rule. Multiple primitive types may be stored by the API; their combined traffic evaluation is not verified.",
	Optional:            true,
	Attributes:          mcpPrimitiveConditionAttributes,
}

// mcpStringConditionAttribute is the server_urls and protocol_versions condition.
var mcpStringConditionAttribute = schema.SingleNestedAttribute{
	MarkdownDescription: "String matching condition. Values are preserved without URL or case normalization.",
	Optional:            true,
	Attributes:          mcpMatchConditionAttributes,
}

// mcpConditionAttributes are the attributes of matching_conditions, one per key of mcpConditionKeys.
var mcpConditionAttributes = map[string]schema.Attribute{
	"server_urls":         mcpStringConditionAttribute,
	"protocol_versions":   mcpStringConditionAttribute,
	"insecure_connection": mcpPropertyConditionAttribute,
	"missing_prm":         mcpPropertyConditionAttribute,
	"tool_matching":       mcpPrimitiveConditionAttribute,
	"resource_matching":   mcpPrimitiveConditionAttribute,
	"prompt_matching":     mcpPrimitiveConditionAttribute,
}

// mcpConditionObjectType is the object type of matching_conditions.
var mcpConditionObjectType = schema.SingleNestedAttribute{Attributes: mcpConditionAttributes}.GetType().(types.ObjectType)

func NewNetworkMCPPolicyRuleResource() resource.Resource {
	return &NetworkMCPPolicyRuleResource{
		ReadPermissions:  []string{"NetworkAccess.Read.All"},
		WritePermissions: []string{"NetworkAccess.ReadWrite.All"},
		ResourcePath:     "/networkaccess/mcpPolicies/{mcpPolicyId}/policyRules",
	}
}

type NetworkMCPPolicyRuleResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

func (r *NetworkMCPPolicyRuleResource) Metadata(
	ctx context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = ResourceName
}

func (r *NetworkMCPPolicyRuleResource) Configure(
	ctx context.Context,
	req resource.ConfigureRequest,
	resp *resource.ConfigureResponse,
) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, ResourceName)
}

func (r *NetworkMCPPolicyRuleResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	if req.Identity != nil && req.ID == "" {
		var identity MCPPolicyRuleIdentity
		resp.Diagnostics.Append(req.Identity.Get(ctx, &identity)...)
		if resp.Diagnostics.HasError() {
			return
		}
		req.ID = identity.MCPPolicyID + "/" + identity.ID
	}
	parts := strings.Split(req.ID, "/")
	if len(parts) != 2 || !regexp.MustCompile(constants.GuidRegex).MatchString(parts[0]) ||
		!regexp.MustCompile(constants.GuidRegex).MatchString(parts[1]) {
		resp.Diagnostics.AddError(
			"Invalid import ID",
			fmt.Sprintf(
				"Expected import ID in the format {mcp_policy_id}/{rule_id}, got %q.",
				req.ID,
			),
		)
		return
	}
	resp.Diagnostics.Append(
		resp.State.SetAttribute(ctx, path.Root("mcp_policy_id"), parts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), parts[1])...)
}

func (r *NetworkMCPPolicyRuleResource) IdentitySchema(
	ctx context.Context,
	req resource.IdentitySchemaRequest,
	resp *resource.IdentitySchemaResponse,
) {
	resp.IdentitySchema = identityschema.Schema{Attributes: map[string]identityschema.Attribute{
		"id":            identityschema.StringAttribute{RequiredForImport: true},
		"mcp_policy_id": identityschema.StringAttribute{RequiredForImport: true},
	}}
}

func (r *NetworkMCPPolicyRuleResource) Schema(
	ctx context.Context,
	req resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages independent rules under Microsoft Entra Global Secure Access MCP policies using the portal-backed Graph beta endpoint. Changing mcp_policy_id replaces the rule. Names and priorities must be unique within the parent. The API accepts rules without conditions. Creating a rule does not assign its policy to traffic. Condition persistence has been verified; traffic evaluation, AND/OR behavior and pattern semantics have not been verified. All CRUD operation timeouts default to 3 minutes.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The MCP policy rule ID.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
			},
			"mcp_policy_id": schema.StringAttribute{
				MarkdownDescription: "The MCP policy ID that owns this rule.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(constants.GuidRegex),
						"must be a valid UUID",
					),
				},
				// The framework modifier also replaces when a newly created parent ID is unknown.
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The nonempty rule name.",
				Required:            true,
				Validators:          []validator.String{stringvalidator.LengthAtLeast(1)},
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The rule description. Omission clears it to null; an empty string is preserved.",
				Optional:            true,
			},
			"action": schema.StringAttribute{
				MarkdownDescription: "The action applied to matching destinations: `allow` or `block`.",
				Required:            true,
				Validators:          []validator.String{stringvalidator.OneOf("allow", "block")},
			},
			"priority": schema.Int32Attribute{
				MarkdownDescription: "The rule priority, from 100 through 2147483647. Must be unique within the parent policy. Use a spare priority when swapping rules.",
				Required:            true,
				Validators:          []validator.Int32{int32validator.AtLeast(100)},
			},
			"enabled": schema.BoolAttribute{
				MarkdownDescription: "Whether the rule is enabled. Must be explicitly configured.",
				Required:            true,
			},
			"status": schema.StringAttribute{
				MarkdownDescription: "The raw API rule status. Normally `enabled` or `disabled`; unknown values are diagnosed rather than converted to enabled.",
				Computed:            true,
			},
			"matching_conditions": schema.SingleNestedAttribute{
				MarkdownDescription: "MCP destination conditions. Omission clears all conditions. Arrays preserve input order; the provider does not normalize URLs or infer traffic matching semantics.",
				Optional:            true,
				Attributes:          mcpConditionAttributes,
			},
			"timeouts": commonschema.ResourceTimeouts(ctx),
		},
	}
}
