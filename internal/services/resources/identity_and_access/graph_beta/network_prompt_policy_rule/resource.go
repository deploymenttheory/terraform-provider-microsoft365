package graphBetaNetworkPromptPolicyRule

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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	planmodifiers "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/plan_modifiers"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
)

const (
	ResourceName  = "microsoft365_graph_beta_identity_and_access_network_prompt_policy_rule"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180

	schemeTypeCustom     = "custom"
	schemeTypePredefined = "predefined"
)

var (
	_ resource.Resource                   = &NetworkPromptPolicyRuleResource{}
	_ resource.ResourceWithConfigure      = &NetworkPromptPolicyRuleResource{}
	_ resource.ResourceWithImportState    = &NetworkPromptPolicyRuleResource{}
	_ resource.ResourceWithValidateConfig = &NetworkPromptPolicyRuleResource{}
	_ resource.ResourceWithIdentity       = &NetworkPromptPolicyRuleResource{}
)

func NewNetworkPromptPolicyRuleResource() resource.Resource {
	return &NetworkPromptPolicyRuleResource{
		ReadPermissions:  []string{"NetworkAccess.Read.All"},
		WritePermissions: []string{"NetworkAccess.ReadWrite.All"},
		ResourcePath:     "/networkaccess/promptPolicies/{promptPolicyId}/policyRules",
	}
}

type NetworkPromptPolicyRuleResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

func (r *NetworkPromptPolicyRuleResource) Metadata(
	ctx context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = ResourceName
}

func (r *NetworkPromptPolicyRuleResource) Configure(
	ctx context.Context,
	req resource.ConfigureRequest,
	resp *resource.ConfigureResponse,
) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, ResourceName)
}

func (r *NetworkPromptPolicyRuleResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	if req.Identity != nil && req.ID == "" {
		var identity PromptPolicyRuleIdentity
		resp.Diagnostics.Append(req.Identity.Get(ctx, &identity)...)
		if resp.Diagnostics.HasError() {
			return
		}
		req.ID = identity.PromptPolicyID + "/" + identity.ID
	}
	parts := strings.Split(req.ID, "/")
	if len(parts) != 2 || !regexp.MustCompile(constants.GuidRegex).MatchString(parts[0]) ||
		!regexp.MustCompile(constants.GuidRegex).MatchString(parts[1]) {
		resp.Diagnostics.AddError(
			"Invalid import ID",
			fmt.Sprintf(
				"Expected import ID in the format {prompt_policy_id}/{rule_id}, got %q.",
				req.ID,
			),
		)
		return
	}
	resp.Diagnostics.Append(
		resp.State.SetAttribute(ctx, path.Root("prompt_policy_id"), parts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), parts[1])...)
}

func (r *NetworkPromptPolicyRuleResource) IdentitySchema(
	ctx context.Context,
	req resource.IdentitySchemaRequest,
	resp *resource.IdentitySchemaResponse,
) {
	resp.IdentitySchema = identityschema.Schema{Attributes: map[string]identityschema.Attribute{
		"id":               identityschema.StringAttribute{RequiredForImport: true},
		"prompt_policy_id": identityschema.StringAttribute{RequiredForImport: true},
	}}
}

func (r *NetworkPromptPolicyRuleResource) Schema(
	ctx context.Context,
	req resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages independent rules within Microsoft Entra Global Secure Access prompt policies using the Microsoft Graph beta `/networkaccess/promptPolicies/{promptPolicyId}/policyRules` endpoint. The parent policy is managed separately. Priority must be unique within a policy; use a spare priority when swapping rules. Creating a rule does not apply the policy to traffic.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The prompt policy rule ID.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
			},
			"prompt_policy_id": schema.StringAttribute{
				MarkdownDescription: "The prompt policy ID that owns this rule.",
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
				MarkdownDescription: "The action applied to matching prompts: `allow` or `block`. Must be explicitly configured.",
				Required:            true,
				Validators:          []validator.String{stringvalidator.OneOf("allow", "block")},
			},
			"priority": schema.Int32Attribute{
				MarkdownDescription: "The rule priority, from 100 through 2147483647. Lower numbers are evaluated first. Must be unique within the parent policy.",
				Required:            true,
				Validators:          []validator.Int32{int32validator.AtLeast(100)},
			},
			"enabled": schema.BoolAttribute{
				MarkdownDescription: "Whether the rule is enabled. Must be explicitly configured.",
				Required:            true,
			},
			"status": schema.StringAttribute{
				MarkdownDescription: "The raw API rule status. Normally `enabled` or `disabled`.",
				Computed:            true,
			},
			"prompt_logging": schema.StringAttribute{
				MarkdownDescription: "When to log prompts: `never`, `always`, or `onBlock`. Defaults to `never`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("never"),
				Validators: []validator.String{
					stringvalidator.OneOf("never", "always", "onBlock"),
				},
			},
			"scan_result": schema.StringAttribute{
				MarkdownDescription: "The scan result to match. Defaults to `maliciousPromptDetected`, the currently supported value.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("maliciousPromptDetected"),
				Validators: []validator.String{
					stringvalidator.OneOf("maliciousPromptDetected"),
				},
			},
			"conversation_schemes": schema.ListNestedAttribute{
				MarkdownDescription: "Conversation schemes matched by the rule. The entire list is managed together. An empty list is accepted by the API; its effect on traffic has not been verified.",
				Required:            true,
				NestedObject: schema.NestedAttributeObject{Attributes: map[string]schema.Attribute{
					"type": schema.StringAttribute{
						MarkdownDescription: "Scheme type: `custom` or `predefined`.",
						Required:            true,
						Validators: []validator.String{
							stringvalidator.OneOf(schemeTypeCustom, schemeTypePredefined),
						},
					},
					"url": schema.StringAttribute{
						MarkdownDescription: "The HTTP or HTTPS endpoint URL for a custom scheme. Required for custom schemes and forbidden for predefined schemes. Include a path, such as `/chat` or `/`.",
						Optional:            true,
						Validators:          []validator.String{stringvalidator.LengthAtLeast(1)},
					},
					"json_path": schema.StringAttribute{
						MarkdownDescription: "JSONPath locating the prompt in a custom request. Only valid for custom schemes. When omitted the API uses an empty string; an explicitly configured empty string is preserved.",
						Optional:            true,
					},
					"scheme_name": schema.StringAttribute{
						MarkdownDescription: "The predefined scheme name. Required for predefined schemes and forbidden for custom schemes.",
						Optional:            true,
						Validators: []validator.String{
							stringvalidator.OneOf(
								"chatGpt",
								"claude",
								"cohere",
								"deepseek",
								"gemini",
								"grok",
								"mistral",
								"perplexity",
								"pi",
								"qwen",
							),
						},
					},
				}},
			},
			"timeouts": commonschema.ResourceTimeouts(ctx),
		},
	}
}
