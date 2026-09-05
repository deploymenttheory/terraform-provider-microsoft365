package graphBetaNetworkTLSInspectionPolicyRule

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
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

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	planmodifiers "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/plan_modifiers"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
)

const (
	ResourceName  = "microsoft365_graph_beta_identity_and_access_network_tls_inspection_policy_rule"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180

	destinationTypeWebCategory = "web_category"
	destinationTypeFQDN        = "fqdn"
)

var (
	_ resource.Resource                = &NetworkTLSInspectionPolicyRuleResource{}
	_ resource.ResourceWithConfigure   = &NetworkTLSInspectionPolicyRuleResource{}
	_ resource.ResourceWithImportState = &NetworkTLSInspectionPolicyRuleResource{}
	_ resource.ResourceWithIdentity    = &NetworkTLSInspectionPolicyRuleResource{}
)

func NewNetworkTLSInspectionPolicyRuleResource() resource.Resource {
	return &NetworkTLSInspectionPolicyRuleResource{
		ReadPermissions:  []string{"NetworkAccess.Read.All"},
		WritePermissions: []string{"NetworkAccess.ReadWrite.All"},
		ResourcePath:     "/networkaccess/tlsInspectionPolicies/{tlsInspectionPolicyId}/policyRules",
	}
}

type NetworkTLSInspectionPolicyRuleResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

func (r *NetworkTLSInspectionPolicyRuleResource) Metadata(
	ctx context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = ResourceName
}

func (r *NetworkTLSInspectionPolicyRuleResource) Configure(
	ctx context.Context,
	req resource.ConfigureRequest,
	resp *resource.ConfigureResponse,
) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, ResourceName)
}

func (r *NetworkTLSInspectionPolicyRuleResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	if req.Identity != nil && req.ID == "" {
		var identity TLSInspectionPolicyRuleIdentity
		resp.Diagnostics.Append(req.Identity.Get(ctx, &identity)...)
		if resp.Diagnostics.HasError() {
			return
		}
		req.ID = identity.TLSInspectionPolicyID + "/" + identity.ID
	}
	parts := strings.Split(req.ID, "/")
	if len(parts) != 2 || !regexp.MustCompile(constants.GuidRegex).MatchString(parts[0]) ||
		!regexp.MustCompile(constants.GuidRegex).MatchString(parts[1]) {
		resp.Diagnostics.AddError(
			"Invalid import ID",
			fmt.Sprintf(
				"Expected import ID in the format {tls_inspection_policy_id}/{rule_id}, got %q.",
				req.ID,
			),
		)
		return
	}
	resp.Diagnostics.Append(
		resp.State.SetAttribute(ctx, path.Root("tls_inspection_policy_id"), parts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), parts[1])...)
}

func (r *NetworkTLSInspectionPolicyRuleResource) IdentitySchema(
	ctx context.Context,
	req resource.IdentitySchemaRequest,
	resp *resource.IdentitySchemaResponse,
) {
	resp.IdentitySchema = identityschema.Schema{Attributes: map[string]identityschema.Attribute{
		"id":                       identityschema.StringAttribute{RequiredForImport: true},
		"tls_inspection_policy_id": identityschema.StringAttribute{RequiredForImport: true},
	}}
}

func (r *NetworkTLSInspectionPolicyRuleResource) Schema(
	ctx context.Context,
	req resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages rules within Microsoft Entra Global Secure Access TLS inspection policies using the Microsoft Graph beta `/networkaccess/tlsInspectionPolicies/{tlsInspectionPolicyId}/policyRules` endpoint. The parent policy is managed separately. Priority must be unique within the policy; use a spare priority when swapping rules. The automatically created system rule (priority 50) cannot be managed by this resource. The recommended bypass rule may be imported explicitly. This resource requires at least one nonempty destination group. Creating a rule does not apply the policy to traffic.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The TLS inspection policy rule ID.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
			},
			"tls_inspection_policy_id": schema.StringAttribute{
				MarkdownDescription: "The TLS inspection policy ID that owns this rule.",
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
				MarkdownDescription: "The action applied to matching destinations: `inspect` or `bypass`.",
				Required:            true,
				Validators:          []validator.String{stringvalidator.OneOf("inspect", "bypass")},
			},
			"priority": schema.Int32Attribute{
				MarkdownDescription: "The rule priority, from 100 through 2147483647. Lower numbers are evaluated first. Must be unique within the parent policy; 65000 is initially used by the recommended bypass rule.",
				Required:            true,
				Validators:          []validator.Int32{int32validator.AtLeast(100)},
			},
			"enabled": schema.BoolAttribute{
				MarkdownDescription: "Whether the rule is enabled. Must be explicitly configured.",
				Required:            true,
			},
			"status": schema.StringAttribute{
				MarkdownDescription: "The raw API rule status. Normally `enabled` or `disabled`; this is not a certificate lifecycle status.",
				Computed:            true,
			},
			"destinations": schema.ListNestedAttribute{
				MarkdownDescription: "Destination groups matched by the rule.",
				Required:            true,
				Validators:          []validator.List{listvalidator.SizeAtLeast(1)},
				NestedObject: schema.NestedAttributeObject{Attributes: map[string]schema.Attribute{
					"type": schema.StringAttribute{
						MarkdownDescription: "Destination type. Possible values are `web_category` and `fqdn`.",
						Required:            true,
						Validators: []validator.String{
							stringvalidator.OneOf(destinationTypeWebCategory, destinationTypeFQDN),
						},
					},
					"values": schema.SetAttribute{
						MarkdownDescription: "Destination values for this group.",
						ElementType:         types.StringType,
						Required:            true,
						Validators: []validator.Set{
							setvalidator.SizeAtLeast(1),
							setvalidator.ValueStringsAre(stringvalidator.LengthAtLeast(1)),
						},
					},
				}},
			},
			"timeouts": commonschema.ResourceTimeouts(ctx),
		},
	}
}
