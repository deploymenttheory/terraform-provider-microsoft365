package graphBetaNetworkFilteringProfilePolicyLink

import (
	"context"
	"fmt"
	"strings"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	planmodifiers "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/plan_modifiers"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
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
	ResourceName  = "microsoft365_graph_beta_identity_and_access_network_filtering_profile_policy_link"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180

	policyTypeFiltering          = "filtering_policy"
	policyTypeWebFiltering       = "web_filtering_policy"
	policyTypeCloudFirewall      = "cloud_firewall_policy"
	policyTypeThreatIntelligence = "threat_intelligence_policy"
	policyTypeTlsInspection      = "tls_inspection_policy"
)

var (
	_ resource.Resource                = &NetworkFilteringProfilePolicyLinkResource{}
	_ resource.ResourceWithConfigure   = &NetworkFilteringProfilePolicyLinkResource{}
	_ resource.ResourceWithImportState = &NetworkFilteringProfilePolicyLinkResource{}
	_ resource.ResourceWithIdentity    = &NetworkFilteringProfilePolicyLinkResource{}
)

func NewNetworkFilteringProfilePolicyLinkResource() resource.Resource {
	return &NetworkFilteringProfilePolicyLinkResource{
		ReadPermissions: []string{
			"NetworkAccess.Read.All",
		},
		WritePermissions: []string{
			"NetworkAccess.ReadWrite.All",
		},
		ResourcePath: "/networkAccess/filteringProfiles/{filteringProfileId}/policies",
	}
}

type NetworkFilteringProfilePolicyLinkResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

func (r *NetworkFilteringProfilePolicyLinkResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = ResourceName
}

func (r *NetworkFilteringProfilePolicyLinkResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, ResourceName)
}

func (r *NetworkFilteringProfilePolicyLinkResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.Split(req.ID, "/")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		resp.Diagnostics.AddError(
			"Invalid import ID",
			fmt.Sprintf("Expected import ID in the format {filtering_profile_id}/{policy_link_id}, got %q.", req.ID),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("filtering_profile_id"), parts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("policy_link_id"), parts[1])...)
}

func (r *NetworkFilteringProfilePolicyLinkResource) IdentitySchema(ctx context.Context, req resource.IdentitySchemaRequest, resp *resource.IdentitySchemaResponse) {
	resp.IdentitySchema = identityschema.Schema{
		Attributes: map[string]identityschema.Attribute{
			"id": identityschema.StringAttribute{
				RequiredForImport: true,
			},
		},
	}
}

func (r *NetworkFilteringProfilePolicyLinkResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a Microsoft Entra Global Secure Access filtering profile policy link using the Microsoft Graph beta `/networkAccess/filteringProfiles/{filteringProfileId}/policies` endpoint.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The Terraform identifier for the link in `{filtering_profile_id}/{policy_link_id}` format.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
			},
			"filtering_profile_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the filtering profile, shown as a security profile in the Microsoft Entra admin center.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.RequiresReplaceString(),
				},
			},
			"policy_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the existing policy to link to the filtering profile.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.RequiresReplaceString(),
				},
			},
			"policy_link_id": schema.StringAttribute{
				MarkdownDescription: "The server-generated ID of the policy link.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
			},
			"policy_type": schema.StringAttribute{
				MarkdownDescription: "The policy type to link. Known values are `filtering_policy`, `web_filtering_policy`, `cloud_firewall_policy`, `threat_intelligence_policy`, and `tls_inspection_policy`.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf(
						policyTypeFiltering,
						policyTypeWebFiltering,
						policyTypeCloudFirewall,
						policyTypeThreatIntelligence,
						policyTypeTlsInspection,
					),
				},
				PlanModifiers: []planmodifier.String{
					planmodifiers.RequiresReplaceString(),
				},
			},
			"state": schema.StringAttribute{
				MarkdownDescription: "The state of the policy link. Possible values are `enabled` and `disabled`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("enabled"),
				Validators: []validator.String{
					stringvalidator.OneOf("enabled", "disabled"),
				},
			},
			"priority": schema.Int64Attribute{
				MarkdownDescription: "The priority for a `filtering_policy` link. The Entra admin center sends this for legacy filtering policy links and does not send it for `web_filtering_policy` links.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					planmodifiers.UseStateForUnknownInt64(),
					planmodifiers.RequiresReplaceInt64(),
				},
			},
			"logging_state": schema.StringAttribute{
				MarkdownDescription: "The logging state for a `filtering_policy` link. The Entra admin center sends this for legacy filtering policy links and does not send it for `web_filtering_policy` links.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("enabled", "disabled"),
				},
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
					planmodifiers.RequiresReplaceString(),
				},
			},
			"created_date_time": schema.StringAttribute{
				MarkdownDescription: "The date and time when the policy link was created.",
				Computed:            true,
			},
			"last_modified_date_time": schema.StringAttribute{
				MarkdownDescription: "The date and time when the policy link was last modified.",
				Computed:            true,
			},
			"version": schema.StringAttribute{
				MarkdownDescription: "The Graph version string for the policy link.",
				Computed:            true,
			},
			"timeouts": commonschema.ResourceTimeouts(ctx),
		},
	}
}

func compositeID(filteringProfileID, policyLinkID string) types.String {
	if filteringProfileID == "" || policyLinkID == "" {
		return types.StringNull()
	}
	return types.StringValue(filteringProfileID + "/" + policyLinkID)
}
