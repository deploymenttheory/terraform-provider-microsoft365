package graphBetaNetworkTLSInspectionPolicy

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
	ResourceName  = "microsoft365_graph_beta_identity_and_access_network_tls_inspection_policy"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	_ resource.Resource                = &NetworkTLSInspectionPolicyResource{}
	_ resource.ResourceWithConfigure   = &NetworkTLSInspectionPolicyResource{}
	_ resource.ResourceWithImportState = &NetworkTLSInspectionPolicyResource{}
	_ resource.ResourceWithIdentity    = &NetworkTLSInspectionPolicyResource{}
)

func NewNetworkTLSInspectionPolicyResource() resource.Resource {
	return &NetworkTLSInspectionPolicyResource{
		ReadPermissions:  []string{"NetworkAccess.Read.All"},
		WritePermissions: []string{"NetworkAccess.ReadWrite.All"},
		ResourcePath:     "/networkaccess/tlsInspectionPolicies",
	}
}

type NetworkTLSInspectionPolicyResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

func (r *NetworkTLSInspectionPolicyResource) Metadata(
	ctx context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = ResourceName
}

func (r *NetworkTLSInspectionPolicyResource) Configure(
	ctx context.Context,
	req resource.ConfigureRequest,
	resp *resource.ConfigureResponse,
) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, ResourceName)
}

func (r *NetworkTLSInspectionPolicyResource) ImportState(
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
		resp.Diagnostics.AddError("Invalid import ID", "Expected the TLS inspection policy UUID.")
		return
	}
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *NetworkTLSInspectionPolicyResource) IdentitySchema(
	ctx context.Context,
	req resource.IdentitySchemaRequest,
	resp *resource.IdentitySchemaResponse,
) {
	resp.IdentitySchema = identityschema.Schema{Attributes: map[string]identityschema.Attribute{
		"id": identityschema.StringAttribute{RequiredForImport: true},
	}}
}

func (r *NetworkTLSInspectionPolicyResource) Schema(
	ctx context.Context,
	req resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages Microsoft Entra Global Secure Access TLS inspection policies using the Microsoft Graph beta `/networkaccess/tlsInspectionPolicies` endpoint. Rules are managed separately with `microsoft365_graph_beta_identity_and_access_network_tls_inspection_policy_rule`. Microsoft creates system and recommended bypass rules automatically; this resource does not adopt them. Deleting a policy also deletes its rules. Creating a policy does not apply it to traffic. TLS inspection requires a configured certificate authority and a filtering profile link. See [TLS inspection](https://learn.microsoft.com/en-us/entra/global-secure-access/how-to-transport-layer-security).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique identifier for the TLS inspection policy.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The nonempty name of the TLS inspection policy.",
				Required:            true,
				Validators:          []validator.String{stringvalidator.LengthAtLeast(1)},
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The policy description. Omission clears the description to null; an empty string is preserved.",
				Optional:            true,
			},
			"default_action": schema.StringAttribute{
				MarkdownDescription: "The action when no rule matches: `inspect` or `bypass`. Must be explicitly configured.",
				Required:            true,
				Validators:          []validator.String{stringvalidator.OneOf("inspect", "bypass")},
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
