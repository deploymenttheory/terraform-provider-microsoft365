package graphBetaNetworkCloudFirewallPolicyRule

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/identityschema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	kiotahttp "github.com/microsoft/kiota-http-go"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	ResourceName  = "microsoft365_graph_beta_identity_and_access_network_cloud_firewall_policy_rule"
	CreateTimeout = 180
	ReadTimeout   = 180
	UpdateTimeout = 180
	DeleteTimeout = 180
)

var (
	_ resource.Resource                   = &NetworkCloudFirewallPolicyRuleResource{}
	_ resource.ResourceWithConfigure      = &NetworkCloudFirewallPolicyRuleResource{}
	_ resource.ResourceWithImportState    = &NetworkCloudFirewallPolicyRuleResource{}
	_ resource.ResourceWithIdentity       = &NetworkCloudFirewallPolicyRuleResource{}
	_ resource.ResourceWithValidateConfig = &NetworkCloudFirewallPolicyRuleResource{}
)

type NetworkCloudFirewallPolicyRuleResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	retryOptions     kiotahttp.RetryHandlerOptions
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

func NewNetworkCloudFirewallPolicyRuleResource() resource.Resource {
	return &NetworkCloudFirewallPolicyRuleResource{ReadPermissions: []string{"NetworkAccess.Read.All"}, WritePermissions: []string{"NetworkAccess.ReadWrite.All"}, ResourcePath: "/networkaccess/cloudFirewallPolicies/{policyId}/policyRules"}
}

func (r *NetworkCloudFirewallPolicyRuleResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = ResourceName
}
func (r *NetworkCloudFirewallPolicyRuleResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, ResourceName)
	r.retryOptions = client.GraphRetryOptionsForResource(req.ProviderData)
}

func (r *NetworkCloudFirewallPolicyRuleResource) IdentitySchema(_ context.Context, _ resource.IdentitySchemaRequest, resp *resource.IdentitySchemaResponse) {
	resp.IdentitySchema = identityschema.Schema{Attributes: map[string]identityschema.Attribute{"policy_id": identityschema.StringAttribute{RequiredForImport: true}, "id": identityschema.StringAttribute{RequiredForImport: true}}}
}

func (r *NetworkCloudFirewallPolicyRuleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	var identity ruleIdentity
	switch {
	case req.ID != "":
		parts := strings.Split(req.ID, "/")
		if len(parts) != 2 {
			resp.Diagnostics.AddError("Invalid import ID", "Expected {policy_id}/{rule_id}.")
			return
		}
		identity = ruleIdentity{PolicyID: parts[0], ID: parts[1]}
	case req.Identity != nil:
		resp.Diagnostics.Append(req.Identity.Get(ctx, &identity)...)
	default:
		resp.Diagnostics.AddError("Missing import identity", "Provide a composite ID or policy_id and id identity attributes.")
		return
	}
	if resp.Diagnostics.HasError() {
		return
	}
	for _, id := range []string{identity.PolicyID, identity.ID} {
		if !regexp.MustCompile(constants.GuidRegex).MatchString(id) {
			resp.Diagnostics.AddError("Invalid import ID", fmt.Sprintf("Expected UUID identifiers, got %q.", id))
			return
		}
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("policy_id"), identity.PolicyID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), identity.ID)...)
	if resp.Identity != nil {
		resp.Diagnostics.Append(resp.Identity.Set(ctx, identity)...)
	}
}

func (r *NetworkCloudFirewallPolicyRuleResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{MarkdownDescription: "Manages an independent rule in a [Microsoft Entra Global Secure Access cloud firewall policy](https://learn.microsoft.com/en-us/graph/api/resources/networkaccess-cloudfirewallrule?view=graph-rest-beta). Uses `/beta/networkaccess/cloudFirewallPolicies/{policy_id}/policyRules`. Creating a rule does not link its policy to traffic. Matching collections are sets: ordering and duplicates do not produce Terraform differences. Rules containing branch conditions are not currently supported.", Attributes: map[string]schema.Attribute{
		"id":          schema.StringAttribute{MarkdownDescription: "The rule ID.", Computed: true, PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
		"policy_id":   schema.StringAttribute{MarkdownDescription: "The parent cloud firewall policy ID. Changing the parent replaces this rule, including when the new parent's ID is not yet known.", Required: true, Validators: []validator.String{stringvalidator.RegexMatches(regexp.MustCompile(constants.GuidRegex), "must be a UUID")}, PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
		"name":        schema.StringAttribute{MarkdownDescription: "The rule name, unique within its parent policy. Updated in place.", Required: true},
		"description": schema.StringAttribute{MarkdownDescription: "The rule description. Omitting it clears the description; an empty string is preserved.", Optional: true},
		"priority":    schema.Int32Attribute{MarkdownDescription: "Priority, unique within the parent policy. Lower values are evaluated first. Microsoft recommends values of at least 100; the configuration API accepts signed 32-bit values. To exchange priorities, move one rule to an unused priority in a separate apply first.", Required: true},
		"action":      schema.StringAttribute{MarkdownDescription: "The action for matching traffic: `allow` or `block`.", Required: true, Validators: []validator.String{stringvalidator.OneOf("allow", "block")}},
		"enabled":     schema.BoolAttribute{MarkdownDescription: "Whether the rule is enabled. Terraform defaults to false and explicitly sends disabled; the API requires a status on creation.", Optional: true, Computed: true, Default: booldefault.StaticBool(false)},
		"status":      schema.StringAttribute{MarkdownDescription: "The observed service status, enabled or disabled. This does not confirm propagation to traffic enforcement.", Computed: true},
		"sources":     matchingSchema(false), "destinations": matchingSchema(true), "timeouts": commonschema.ResourceTimeouts(ctx),
	}}
}

func matchingSchema(destination bool) schema.SingleNestedAttribute {
	addressTypes := []string{"ip"}
	description := "Source conditions. Omit this object to leave source matching unspecified."
	if destination {
		addressTypes = append(addressTypes, "fqdn")
		description = "Destination conditions. Omit this object to leave destination matching unspecified. FQDN conditions can be stored but do not work for remote-network Internet Access traffic."
	}
	attrs := map[string]schema.Attribute{
		"addresses": schema.SetNestedAttribute{MarkdownDescription: "Address groups, at most one per type. Each group's values are matched with OR. Terraform defaults to an empty set.", Optional: true, Computed: true, Default: setdefault.StaticValue(types.SetValueMust(addressObjectType(), nil)), NestedObject: schema.NestedAttributeObject{Attributes: map[string]schema.Attribute{
			"type":   schema.StringAttribute{MarkdownDescription: "Address type: `ip` for IPv4 addresses, CIDRs or ranges; destinations also support `fqdn`.", Required: true, Validators: []validator.String{stringvalidator.OneOf(addressTypes...)}},
			"values": schema.SetAttribute{MarkdownDescription: "Address values. IPv4 examples: `192.0.2.1`, `192.0.2.0/24`, `192.0.2.1-192.0.2.10`. FQDN example: `api.example.com`; a standalone `*` is accepted, but `*.example.com` is not. Case and CIDR host bits are preserved.", Required: true, ElementType: types.StringType},
		}}},
		"ports": schema.SetAttribute{MarkdownDescription: "Ports as decimal strings in 0–65535 or inclusive ascending ranges such as `1000-1002`. Leading zeros and wildcard strings are not accepted. Terraform defaults to an empty set.", Optional: true, Computed: true, ElementType: types.StringType, Default: setdefault.StaticValue(types.SetValueMust(types.StringType, nil)), Validators: []validator.Set{setvalidator.ValueStringsAre(portValidator{})}},
	}
	if destination {
		attrs["protocols"] = schema.SetAttribute{MarkdownDescription: "Protocols to match: `tcp`, `udp`, or both. Required whenever destinations is specified; null is not sent to the API.", Required: true, ElementType: types.StringType, Validators: []validator.Set{setvalidator.SizeAtLeast(1), setvalidator.ValueStringsAre(stringvalidator.OneOf("tcp", "udp"))}}
	}
	return schema.SingleNestedAttribute{MarkdownDescription: description, Optional: true, Attributes: attrs}
}
