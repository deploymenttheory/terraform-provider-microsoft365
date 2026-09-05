package graphBetaNetworkPromptPolicy

import (
	"context"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/identityschema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	planmodifiers "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/plan_modifiers"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
)

const (
	ResourceName  = "microsoft365_graph_beta_identity_and_access_network_prompt_policy"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	_ resource.Resource                = &NetworkPromptPolicyResource{}
	_ resource.ResourceWithConfigure   = &NetworkPromptPolicyResource{}
	_ resource.ResourceWithImportState = &NetworkPromptPolicyResource{}
	_ resource.ResourceWithIdentity    = &NetworkPromptPolicyResource{}
)

func NewNetworkPromptPolicyResource() resource.Resource {
	return &NetworkPromptPolicyResource{
		ReadPermissions:  []string{"NetworkAccess.Read.All"},
		WritePermissions: []string{"NetworkAccess.ReadWrite.All"},
		ResourcePath:     "/networkaccess/promptPolicies",
	}
}

type NetworkPromptPolicyResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

func (r *NetworkPromptPolicyResource) Metadata(
	ctx context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = ResourceName
}

func (r *NetworkPromptPolicyResource) Configure(
	ctx context.Context,
	req resource.ConfigureRequest,
	resp *resource.ConfigureResponse,
) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, ResourceName)
}

func (r *NetworkPromptPolicyResource) ImportState(
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
		resp.Diagnostics.AddError("Invalid import ID", "Expected the prompt policy UUID.")
		return
	}
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *NetworkPromptPolicyResource) IdentitySchema(
	ctx context.Context,
	req resource.IdentitySchemaRequest,
	resp *resource.IdentitySchemaResponse,
) {
	resp.IdentitySchema = identityschema.Schema{Attributes: map[string]identityschema.Attribute{
		"id": identityschema.StringAttribute{RequiredForImport: true},
	}}
}

func (r *NetworkPromptPolicyResource) Schema(
	ctx context.Context,
	req resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages Microsoft Entra Global Secure Access prompt policies using the portal-backed Microsoft Graph beta `/networkaccess/promptPolicies` endpoint. Rules are managed separately with `microsoft365_graph_beta_identity_and_access_network_prompt_policy_rule`. Deleting a policy also deletes all its rules, including rules outside Terraform. Creating a policy does not apply it to traffic; use a separate security profile policy link. See [prompt injection protection](https://learn.microsoft.com/en-us/entra/global-secure-access/how-to-ai-prompt-injection-protection).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique identifier for the prompt policy.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The nonempty name of the prompt policy.",
				Required:            true,
				Validators:          []validator.String{stringvalidator.LengthAtLeast(1)},
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The policy description. Omission clears the description to null; an empty string is preserved.",
				Optional:            true,
			},
			"default_action": schema.StringAttribute{
				MarkdownDescription: "The action when no rule matches. Only `allow` is supported. Defaults to `allow`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("allow"),
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
