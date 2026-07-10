package graphBetaNetworkForwardingProfilePolicyLink

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	planmodifiers "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/plan_modifiers"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/identityschema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	ResourceName  = "microsoft365_graph_beta_identity_and_access_network_forwarding_profile_policy_link"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180

	forwardingPolicyLinkODataType = "#microsoft.graph.networkaccess.forwardingPolicyLink"
)

var (
	_ resource.Resource                = &NetworkForwardingProfilePolicyLinkResource{}
	_ resource.ResourceWithConfigure   = &NetworkForwardingProfilePolicyLinkResource{}
	_ resource.ResourceWithImportState = &NetworkForwardingProfilePolicyLinkResource{}
	_ resource.ResourceWithIdentity    = &NetworkForwardingProfilePolicyLinkResource{}
)

func NewNetworkForwardingProfilePolicyLinkResource() resource.Resource {
	return &NetworkForwardingProfilePolicyLinkResource{
		ReadPermissions: []string{
			"NetworkAccess.Read.All",
		},
		WritePermissions: []string{
			"NetworkAccess.ReadWrite.All",
		},
		ResourcePath: "/networkAccess/forwardingProfiles/{forwardingProfileId}/policies/{policyLinkId}",
	}
}

type NetworkForwardingProfilePolicyLinkResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

func (r *NetworkForwardingProfilePolicyLinkResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = ResourceName
}

func (r *NetworkForwardingProfilePolicyLinkResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, ResourceName)
}

func (r *NetworkForwardingProfilePolicyLinkResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.Split(req.ID, "/")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		resp.Diagnostics.AddError("Invalid import ID", fmt.Sprintf("Expected import ID in the format {forwarding_profile_id}/{policy_link_id}, got %q.", req.ID))
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("forwarding_profile_id"), parts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("policy_link_id"), parts[1])...)
}

func (r *NetworkForwardingProfilePolicyLinkResource) IdentitySchema(ctx context.Context, req resource.IdentitySchemaRequest, resp *resource.IdentitySchemaResponse) {
	resp.IdentitySchema = identityschema.Schema{
		Attributes: map[string]identityschema.Attribute{
			"id": identityschema.StringAttribute{RequiredForImport: true},
		},
	}
}

func (r *NetworkForwardingProfilePolicyLinkResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages the state of a Microsoft-created Global Secure Access forwarding profile policy link using Microsoft Graph beta. Destroy does not delete the Microsoft-managed link; it only removes Terraform state.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Terraform identifier in `{forwarding_profile_id}/{policy_link_id}` format.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
			},
			"forwarding_profile_id": schema.StringAttribute{
				MarkdownDescription: "The forwarding profile ID.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(regexp.MustCompile(constants.GuidRegex), "must be a valid UUID"),
				},
				PlanModifiers: []planmodifier.String{
					planmodifiers.RequiresReplaceString(),
				},
			},
			"policy_link_id": schema.StringAttribute{
				MarkdownDescription: "The forwarding policy link ID. This is distinct from the linked forwarding policy ID.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(regexp.MustCompile(constants.GuidRegex), "must be a valid UUID"),
				},
				PlanModifiers: []planmodifier.String{
					planmodifiers.RequiresReplaceString(),
				},
			},
			"state": schema.StringAttribute{
				MarkdownDescription: "The forwarding policy link state. Possible values are `enabled` and `disabled`.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("enabled", "disabled"),
				},
			},
			"priority":                schema.Int64Attribute{Computed: true, MarkdownDescription: "The forwarding policy link priority."},
			"version":                 schema.StringAttribute{Computed: true, MarkdownDescription: "The forwarding policy link version."},
			"policy_id":               schema.StringAttribute{Computed: true, MarkdownDescription: "The linked forwarding policy ID."},
			"policy_name":             schema.StringAttribute{Computed: true, MarkdownDescription: "The linked forwarding policy name."},
			"policy_description":      schema.StringAttribute{Computed: true, MarkdownDescription: "The linked forwarding policy description."},
			"traffic_forwarding_type": schema.StringAttribute{Computed: true, MarkdownDescription: "The linked forwarding policy traffic forwarding type."},
			"timeouts":                commonschema.ResourceTimeouts(ctx),
		},
	}
}

func compositeID(forwardingProfileID, policyLinkID string) types.String {
	if forwardingProfileID == "" || policyLinkID == "" {
		return types.StringNull()
	}
	return types.StringValue(forwardingProfileID + "/" + policyLinkID)
}
