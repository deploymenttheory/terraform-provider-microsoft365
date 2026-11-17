package graphBetaNetworkForwardingProfilePrivateAccessPolicyLink

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
)

const (
	ResourceName  = "microsoft365_graph_beta_identity_and_access_network_forwarding_profile_private_access_policy_link"
	CreateTimeout = 300
	ReadTimeout   = 180
	UpdateTimeout = 180
	DeleteTimeout = 300

	profileTypeUser   = "user"
	profileTypeDevice = "device"
	stateEnabled      = "enabled"
	stateDisabled     = "disabled"
)

var (
	_ resource.Resource                = &NetworkForwardingProfilePrivateAccessPolicyLinkResource{}
	_ resource.ResourceWithConfigure   = &NetworkForwardingProfilePrivateAccessPolicyLinkResource{}
	_ resource.ResourceWithImportState = &NetworkForwardingProfilePrivateAccessPolicyLinkResource{}

	guidPattern = regexp.MustCompile(constants.GuidRegex)
)

// NetworkForwardingProfilePrivateAccessPolicyLinkResource links a Private Access
// forwarding policy to a private forwarding profile.
type NetworkForwardingProfilePrivateAccessPolicyLinkResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

func NewNetworkForwardingProfilePrivateAccessPolicyLinkResource() resource.Resource {
	return &NetworkForwardingProfilePrivateAccessPolicyLinkResource{
		ReadPermissions:  []string{"NetworkAccess.Read.All"},
		WritePermissions: []string{"NetworkAccess.ReadWrite.All"},
		ResourcePath:     "/networkAccess/forwardingProfiles/{forwardingProfileId}/policies",
	}
}

func (r *NetworkForwardingProfilePrivateAccessPolicyLinkResource) Metadata(
	ctx context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = ResourceName
}

func (r *NetworkForwardingProfilePrivateAccessPolicyLinkResource) Configure(
	ctx context.Context,
	req resource.ConfigureRequest,
	resp *resource.ConfigureResponse,
) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, ResourceName)
}

// ImportState accepts {profile_type}/{forwarding_profile_id}/{policy_link_id} or
// {forwarding_profile_id}/{policy_link_id} (user profile).
func (r *NetworkForwardingProfilePrivateAccessPolicyLinkResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	profileType, profileID, linkID, ok := parseImportID(req.ID)
	if !ok {
		resp.Diagnostics.AddError(
			"Invalid import ID",
			fmt.Sprintf(
				"Expected {forwarding_profile_id}/{policy_link_id} or {user|device}/{forwarding_profile_id}/{policy_link_id} with GUID IDs, got %q.",
				req.ID,
			),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), compositeID(profileID, linkID))...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("profile_type"), profileType)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("forwarding_profile_id"), profileID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("policy_link_id"), linkID)...)
}

func parseImportID(id string) (profileType, profileID, linkID string, ok bool) {
	parts := strings.Split(id, "/")
	switch len(parts) {
	case 2:
		profileType, profileID, linkID = profileTypeUser, parts[0], parts[1]
	case 3:
		profileType, profileID, linkID = parts[0], parts[1], parts[2]
	default:
		return "", "", "", false
	}
	if profileType != profileTypeUser && profileType != profileTypeDevice {
		return "", "", "", false
	}
	if !guidPattern.MatchString(profileID) || !guidPattern.MatchString(linkID) {
		return "", "", "", false
	}
	// Graph IDs are lowercase; an uppercase import would otherwise force a replacement.
	return profileType, strings.ToLower(profileID), strings.ToLower(linkID), true
}

func compositeID(profileID, linkID string) types.String {
	if profileID == "" || linkID == "" {
		return types.StringNull()
	}
	return types.StringValue(profileID + "/" + linkID)
}

func (r *NetworkForwardingProfilePrivateAccessPolicyLinkResource) Schema(
	ctx context.Context,
	req resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	guid := stringvalidator.RegexMatches(guidPattern, "must be a valid UUID")
	resp.Schema = schema.Schema{
		MarkdownDescription: "Links one Global Secure Access Private Access forwarding policy (app or DNS policy) to a private " +
			"forwarding profile using the `policies@delta` collection of " +
			"`/networkAccess/forwardingProfiles/{id}` or `/networkAccess/deviceForwardingProfiles/{id}`. " +
			"Graph only supports adding and removing links, so every argument forces replacement.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Terraform identifier in `{forwarding_profile_id}/{policy_link_id}` format.",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"forwarding_profile_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "ID of the private forwarding profile to link the policy to.",
				Validators:          []validator.String{guid},
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"profile_type": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString(profileTypeUser),
				MarkdownDescription: "Forwarding profile kind. `user` (default) targets `/networkAccess/forwardingProfiles`; " +
					"`device` targets the pre-logon `/networkAccess/deviceForwardingProfiles`.",
				Validators:    []validator.String{stringvalidator.OneOf(profileTypeUser, profileTypeDevice)},
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"forwarding_policy_id": schema.StringAttribute{
				Required: true,
				MarkdownDescription: "ID of the Private Access forwarding policy to link. This is the forwardingPolicy ID, " +
					"not the application or service principal ID.",
				Validators:    []validator.String{guid},
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"config_type": schema.StringAttribute{
				Required: true,
				MarkdownDescription: "Link configuration type. Possible values are `appPolicy` and `dnsPolicy`. " +
					"`dnsPolicy` uses the same payload as `appPolicy` but has not been verified against a live tenant.",
				Validators:    []validator.String{stringvalidator.OneOf("appPolicy", "dnsPolicy")},
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"priority": schema.Int64Attribute{
				Optional: true,
				Computed: true,
				MarkdownDescription: "Link priority, 1-65000 (the Entra admin center offers 100-64999). If omitted, an " +
					"adopted automatic link keeps `65000` and a new link has none.",
				Validators: []validator.Int64{int64validator.Between(1, 65000)},
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
					int64planmodifier.RequiresReplace(),
				},
			},
			"state": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(stateEnabled),
				MarkdownDescription: "Link state. Possible values are `enabled` (default) and `disabled`.",
				Validators:          []validator.String{stringvalidator.OneOf(stateEnabled, stateDisabled)},
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"policy_link_id": schema.StringAttribute{
				Computed: true,
				MarkdownDescription: "ID of the forwardingPolicyLink object. Graph issues a new ID whenever the link is " +
					"removed and re-added.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"timeouts": commonschema.ResourceTimeouts(ctx),
		},
	}
}
