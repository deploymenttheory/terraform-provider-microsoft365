package graphBetaNetworkContentPolicyRule

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	ResourceName  = "microsoft365_graph_beta_identity_and_access_network_content_policy_rule"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180

	destinationTypeWebCategory = "web_category"
	destinationTypeFQDN        = "fqdn"
	destinationTypeURL         = "url"
)

var (
	_ resource.Resource                     = &NetworkContentPolicyRuleResource{}
	_ resource.ResourceWithConfigure        = &NetworkContentPolicyRuleResource{}
	_ resource.ResourceWithImportState      = &NetworkContentPolicyRuleResource{}
	_ resource.ResourceWithIdentity         = &NetworkContentPolicyRuleResource{}
	_ resource.ResourceWithConfigValidators = &NetworkContentPolicyRuleResource{}
)

func NewNetworkContentPolicyRuleResource() resource.Resource {
	return &NetworkContentPolicyRuleResource{
		ReadPermissions:  []string{"NetworkAccess.Read.All"},
		WritePermissions: []string{"NetworkAccess.ReadWrite.All"},
		ResourcePath:     "/networkaccess/filePolicies/{contentPolicyId}/policyRules",
	}
}

type NetworkContentPolicyRuleResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

func (r *NetworkContentPolicyRuleResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = ResourceName
}

func (r *NetworkContentPolicyRuleResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, ResourceName)
}

func (r *NetworkContentPolicyRuleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.Split(req.ID, "/")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		resp.Diagnostics.AddError("Invalid import ID", fmt.Sprintf("Expected import ID in the format {content_policy_id}/{rule_id}, got %q.", req.ID))
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("content_policy_id"), parts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), parts[1])...)
}

func (r *NetworkContentPolicyRuleResource) IdentitySchema(ctx context.Context, req resource.IdentitySchemaRequest, resp *resource.IdentitySchemaResponse) {
	resp.IdentitySchema = identityschema.Schema{Attributes: map[string]identityschema.Attribute{
		"id": identityschema.StringAttribute{RequiredForImport: true},
	}}
}

func (r *NetworkContentPolicyRuleResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{
		resourcevalidator.AtLeastOneOf(path.MatchRoot("content_types"), path.MatchRoot("text_content_types")),
	}
}

func (r *NetworkContentPolicyRuleResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages rules within Microsoft Entra Global Secure Access content policies using the Microsoft Graph beta `/networkaccess/filePolicies/{contentPolicyId}/policyRules` endpoint.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The content policy rule ID.",
				Computed:            true,
				PlanModifiers:       []planmodifier.String{planmodifiers.UseStateForUnknownString()},
			},
			"content_policy_id": schema.StringAttribute{
				MarkdownDescription: "The content policy ID that owns this rule.",
				Required:            true,
				Validators:          []validator.String{stringvalidator.RegexMatches(regexp.MustCompile(constants.GuidRegex), "must be a valid UUID")},
				PlanModifiers:       []planmodifier.String{planmodifiers.RequiresReplaceString()},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The rule name.",
				Required:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The rule description. Defaults to an empty string.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
				Validators:          []validator.String{stringvalidator.LengthAtMost(8192)},
			},
			"action": schema.StringAttribute{
				MarkdownDescription: "The action applied to matching files. Possible values are `allow`, `block`, and `scanPurview`.",
				Required:            true,
				Validators:          []validator.String{stringvalidator.OneOf("allow", "block", "scanPurview")},
			},
			"priority": schema.Int64Attribute{
				MarkdownDescription: "The rule priority. Lower numbers are evaluated first.",
				Required:            true,
				Validators:          []validator.Int64{int64validator.Between(100, 65000)},
			},
			"status": schema.StringAttribute{
				MarkdownDescription: "The rule status. Possible values are `enabled` and `disabled`.",
				Required:            true,
				Validators:          []validator.String{stringvalidator.OneOf("enabled", "disabled")},
			},
			"activities": schema.SetAttribute{
				MarkdownDescription: "File activities matched by the rule. Possible values are `download` and `upload`.",
				ElementType:         types.StringType,
				Required:            true,
				Validators: []validator.Set{
					setvalidator.SizeAtLeast(1),
					setvalidator.ValueStringsAre(stringvalidator.OneOf("download", "upload")),
				},
			},
			"content_types": schema.SetAttribute{
				MarkdownDescription: "MIME content types matched by the rule.",
				ElementType:         types.StringType,
				Optional:            true,
				Validators: []validator.Set{
					setvalidator.SizeAtLeast(1),
					setvalidator.ValueStringsAre(stringvalidator.LengthAtLeast(1)),
				},
			},
			"text_content_types": schema.SetAttribute{
				MarkdownDescription: "Text content types matched by the rule. Possible values are `json`, `plain`, `html`, and `xml`.",
				ElementType:         types.StringType,
				Optional:            true,
				Validators: []validator.Set{
					setvalidator.SizeAtLeast(1),
					setvalidator.ValueStringsAre(stringvalidator.OneOf("json", "plain", "html", "xml")),
				},
			},
			"destinations": schema.ListNestedAttribute{
				MarkdownDescription: "Destination groups matched by the rule.",
				Required:            true,
				Validators:          []validator.List{listvalidator.SizeAtLeast(1)},
				NestedObject: schema.NestedAttributeObject{Attributes: map[string]schema.Attribute{
					"type": schema.StringAttribute{
						MarkdownDescription: "Destination type. Possible values are `web_category`, `fqdn`, and `url`.",
						Required:            true,
						Validators:          []validator.String{stringvalidator.OneOf(destinationTypeWebCategory, destinationTypeFQDN, destinationTypeURL)},
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
			"session_types": schema.SetAttribute{
				MarkdownDescription: "Session types matched by the rule. Possible values are `user` and `agent`.",
				ElementType:         types.StringType,
				Required:            true,
				Validators: []validator.Set{
					setvalidator.SizeAtLeast(1),
					setvalidator.ValueStringsAre(stringvalidator.OneOf("user", "agent")),
				},
			},
			"timeouts": commonschema.ResourceTimeouts(ctx),
		},
	}
}
