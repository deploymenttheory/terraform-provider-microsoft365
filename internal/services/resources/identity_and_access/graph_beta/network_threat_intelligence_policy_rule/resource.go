package graphBetaNetworkThreatIntelligencePolicyRule

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
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
	ResourceName  = "microsoft365_graph_beta_identity_and_access_network_threat_intelligence_policy_rule"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180

	destinationTypeFQDN = "fqdn"
)

var (
	_ resource.Resource                = &NetworkThreatIntelligencePolicyRuleResource{}
	_ resource.ResourceWithConfigure   = &NetworkThreatIntelligencePolicyRuleResource{}
	_ resource.ResourceWithImportState = &NetworkThreatIntelligencePolicyRuleResource{}
	_ resource.ResourceWithIdentity    = &NetworkThreatIntelligencePolicyRuleResource{}
)

func NewNetworkThreatIntelligencePolicyRuleResource() resource.Resource {
	return &NetworkThreatIntelligencePolicyRuleResource{
		ReadPermissions:  []string{"NetworkAccess.Read.All"},
		WritePermissions: []string{"NetworkAccess.ReadWrite.All"},
		ResourcePath:     "/networkaccess/threatIntelligencePolicies/{threatIntelligencePolicyId}/policyRules",
	}
}

type NetworkThreatIntelligencePolicyRuleResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

func (r *NetworkThreatIntelligencePolicyRuleResource) Metadata(
	ctx context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = ResourceName
}

func (r *NetworkThreatIntelligencePolicyRuleResource) Configure(
	ctx context.Context,
	req resource.ConfigureRequest,
	resp *resource.ConfigureResponse,
) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, ResourceName)
}

func (r *NetworkThreatIntelligencePolicyRuleResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	if req.Identity != nil && req.ID == "" {
		var identity ThreatIntelligencePolicyRuleIdentity
		resp.Diagnostics.Append(req.Identity.Get(ctx, &identity)...)
		if resp.Diagnostics.HasError() {
			return
		}
		req.ID = identity.ThreatIntelligencePolicyID + "/" + identity.ID
	}
	parts := strings.Split(req.ID, "/")
	if len(parts) != 2 || !regexp.MustCompile(constants.GuidRegex).MatchString(parts[0]) ||
		!regexp.MustCompile(constants.GuidRegex).MatchString(parts[1]) {
		resp.Diagnostics.AddError(
			"Invalid import ID",
			fmt.Sprintf(
				"Expected import ID in the format {threat_intelligence_policy_id}/{rule_id}, got %q.",
				req.ID,
			),
		)
		return
	}
	resp.Diagnostics.Append(
		resp.State.SetAttribute(ctx, path.Root("threat_intelligence_policy_id"), parts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), parts[1])...)
}

func (r *NetworkThreatIntelligencePolicyRuleResource) IdentitySchema(
	ctx context.Context,
	req resource.IdentitySchemaRequest,
	resp *resource.IdentitySchemaResponse,
) {
	resp.IdentitySchema = identityschema.Schema{Attributes: map[string]identityschema.Attribute{
		"id":                            identityschema.StringAttribute{RequiredForImport: true},
		"threat_intelligence_policy_id": identityschema.StringAttribute{RequiredForImport: true},
	}}
}

func (r *NetworkThreatIntelligencePolicyRuleResource) Schema(
	ctx context.Context,
	req resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages rules within Microsoft Entra Global Secure Access threat intelligence policies using the Microsoft Graph beta `/networkaccess/threatIntelligencePolicies/{threatIntelligencePolicyId}/policyRules` endpoint. The parent policy is managed separately. Priority must be unique within the policy; use a spare priority when swapping rules. The automatically created default rule (priority 65000) cannot be managed or imported by this resource. Creating a rule does not apply the policy to traffic.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The threat intelligence policy rule ID.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
			},
			"threat_intelligence_policy_id": schema.StringAttribute{
				MarkdownDescription: "The threat intelligence policy ID that owns this rule.",
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
				MarkdownDescription: "The rule name. Must be unique within the parent policy.",
				Required:            true,
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
				MarkdownDescription: "The rule priority, from 100 through 2147483647. Lower numbers are evaluated first. Must be unique within the parent policy; 65000 is reserved for the automatic default rule.",
				Required:            true,
				Validators: []validator.Int32{
					int32validator.AtLeast(100),
					int32validator.NoneOf(65000),
				},
			},
			"severity": schema.StringAttribute{
				MarkdownDescription: "Threat severity. Explicitly configure `high`; the service default differs from its accepted write values. Import preserves an existing `low` value, but the next apply must reconcile it to `high`.",
				Required:            true,
				Validators:          []validator.String{stringvalidator.OneOf("high")},
			},
			"enabled": schema.BoolAttribute{
				MarkdownDescription: "Whether the rule is enabled. Must be explicitly configured.",
				Required:            true,
			},
			"status": schema.StringAttribute{
				MarkdownDescription: "The raw API rule status. Normally `enabled` or `disabled`; unknown statuses produce a diagnostic without being interpreted as disabled.",
				Computed:            true,
			},
			"destinations": schema.ListNestedAttribute{
				MarkdownDescription: "Destination groups matched by the rule. An empty list is accepted; at most one FQDN group is supported.",
				Required:            true,
				Validators:          []validator.List{listvalidator.SizeAtMost(1)},
				NestedObject: schema.NestedAttributeObject{Attributes: map[string]schema.Attribute{
					"type": schema.StringAttribute{
						MarkdownDescription: "Destination type. Only `fqdn` is supported.",
						Required:            true,
						Validators: []validator.String{
							stringvalidator.OneOf(destinationTypeFQDN),
						},
					},
					"values": schema.ListAttribute{
						MarkdownDescription: "FQDN values. Order, duplicates, and case are preserved as stored by the API. URLs and trailing slashes are rejected by the service.",
						ElementType:         types.StringType,
						Required:            true,
						Validators: []validator.List{
							listvalidator.SizeAtLeast(1),
							listvalidator.ValueStringsAre(stringvalidator.LengthAtLeast(1)),
						},
					},
				}},
			},
			"timeouts": commonschema.ResourceTimeouts(ctx),
		},
	}
}
