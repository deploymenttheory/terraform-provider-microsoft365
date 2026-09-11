package graphBetaNetworkThreatIntelligencePolicy

import (
	"context"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/identityschema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	planmodifiers "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/plan_modifiers"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
)

const (
	ResourceName  = "microsoft365_graph_beta_identity_and_access_network_threat_intelligence_policy"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	_ resource.Resource                = &NetworkThreatIntelligencePolicyResource{}
	_ resource.ResourceWithConfigure   = &NetworkThreatIntelligencePolicyResource{}
	_ resource.ResourceWithImportState = &NetworkThreatIntelligencePolicyResource{}
	_ resource.ResourceWithIdentity    = &NetworkThreatIntelligencePolicyResource{}
)

func NewNetworkThreatIntelligencePolicyResource() resource.Resource {
	return &NetworkThreatIntelligencePolicyResource{
		ReadPermissions:  []string{"NetworkAccess.Read.All"},
		WritePermissions: []string{"NetworkAccess.ReadWrite.All"},
		ResourcePath:     "/networkaccess/threatIntelligencePolicies",
	}
}

type NetworkThreatIntelligencePolicyResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

func (r *NetworkThreatIntelligencePolicyResource) Metadata(
	ctx context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = ResourceName
}

func (r *NetworkThreatIntelligencePolicyResource) Configure(
	ctx context.Context,
	req resource.ConfigureRequest,
	resp *resource.ConfigureResponse,
) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, ResourceName)
}

func (r *NetworkThreatIntelligencePolicyResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	if req.ID == "" && req.Identity != nil {
		resp.Diagnostics.Append(req.Identity.GetAttribute(ctx, path.Root("id"), &req.ID)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}
	if !regexp.MustCompile(constants.GuidRegex).MatchString(req.ID) {
		resp.Diagnostics.AddError(
			"Invalid import ID",
			"Expected the threat intelligence policy UUID.",
		)
		return
	}
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *NetworkThreatIntelligencePolicyResource) IdentitySchema(
	ctx context.Context,
	req resource.IdentitySchemaRequest,
	resp *resource.IdentitySchemaResponse,
) {
	resp.IdentitySchema = identityschema.Schema{Attributes: map[string]identityschema.Attribute{
		"id": identityschema.StringAttribute{RequiredForImport: true},
	}}
}

func (r *NetworkThreatIntelligencePolicyResource) Schema(
	ctx context.Context,
	req resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages Microsoft Entra Global Secure Access threat intelligence policies using the Microsoft Graph beta `/networkaccess/threatIntelligencePolicies` endpoint. Rules are managed separately with `microsoft365_graph_beta_identity_and_access_network_threat_intelligence_policy_rule`. Microsoft automatically creates a default blocking rule. It cannot be deleted or disabled independently and is managed only through the policy lifecycle. Deleting a policy also deletes its rules. Creating a policy does not apply it to traffic. Traffic enforcement requires a filtering profile link and its applicable assignment. See [Configure threat intelligence](https://learn.microsoft.com/en-us/entra/global-secure-access/how-to-configure-threat-intelligence).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique identifier for the threat intelligence policy.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the threat intelligence policy.",
				Required:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The policy description. Omission clears the description to null; an empty string is preserved.",
				Optional:            true,
			},
			"default_action": schema.StringAttribute{
				MarkdownDescription: "The action when no rule matches: `allow`. Must be explicitly configured.",
				Required:            true,
				Validators:          []validator.String{stringvalidator.OneOf("allow")},
			},
			"version": schema.StringAttribute{
				MarkdownDescription: "The API-assigned policy version. Read-only.",
				Computed:            true,
			},
			"last_modified_date_time": schema.StringAttribute{
				MarkdownDescription: "The last modification timestamp. Changes to child rules can also change this value.",
				Computed:            true,
			},
			"timeouts": commonschema.ResourceTimeouts(ctx),
		},
	}
}
