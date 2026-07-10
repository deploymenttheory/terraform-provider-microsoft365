package graphBetaNetworkForwardingProfilePolicyLink

import (
	"context"
	"regexp"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	DataSourceName = "microsoft365_graph_beta_identity_and_access_network_forwarding_profile_policy_link"
	ReadTimeout    = 180
)

var (
	_ datasource.DataSource              = &NetworkForwardingProfilePolicyLinkDataSource{}
	_ datasource.DataSourceWithConfigure = &NetworkForwardingProfilePolicyLinkDataSource{}
)

func NewNetworkForwardingProfilePolicyLinkDataSource() datasource.DataSource {
	return &NetworkForwardingProfilePolicyLinkDataSource{
		ReadPermissions: []string{
			"NetworkAccess.Read.All",
		},
	}
}

type NetworkForwardingProfilePolicyLinkDataSource struct {
	client          *msgraphbetasdk.GraphServiceClient
	ReadPermissions []string
}

func (d *NetworkForwardingProfilePolicyLinkDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = DataSourceName
}

func (d *NetworkForwardingProfilePolicyLinkDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.client = client.SetGraphBetaClientForDataSource(ctx, req, resp, DataSourceName)
}

func (d *NetworkForwardingProfilePolicyLinkDataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Retrieves a Microsoft Entra Global Secure Access forwarding profile policy link by forwarding profile selector and linked policy name.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The unique identifier for the data source operation.",
			},
			"forwarding_profile_id": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The forwarding profile ID. Conflicts with `forwarding_profile_name` and `traffic_forwarding_type` when used as a lookup attribute.",
				Validators: []validator.String{
					stringvalidator.RegexMatches(regexp.MustCompile(constants.GuidRegex), "must be a valid UUID"),
					stringvalidator.ConflictsWith(
						path.MatchRoot("forwarding_profile_name"),
						path.MatchRoot("traffic_forwarding_type"),
					),
					stringvalidator.AtLeastOneOf(
						path.MatchRoot("forwarding_profile_id"),
						path.MatchRoot("forwarding_profile_name"),
						path.MatchRoot("traffic_forwarding_type"),
					),
				},
			},
			"forwarding_profile_name": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The forwarding profile name to match exactly, ignoring case. Conflicts with `forwarding_profile_id` and `traffic_forwarding_type`.",
				Validators: []validator.String{
					stringvalidator.ConflictsWith(
						path.MatchRoot("forwarding_profile_id"),
						path.MatchRoot("traffic_forwarding_type"),
					),
				},
			},
			"traffic_forwarding_type": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The forwarding profile traffic type. Valid values are `internet`, `m365`, and `private`. Conflicts with `forwarding_profile_id` and `forwarding_profile_name`.",
				Validators: []validator.String{
					stringvalidator.OneOf("internet", "m365", "private"),
					stringvalidator.ConflictsWith(
						path.MatchRoot("forwarding_profile_id"),
						path.MatchRoot("forwarding_profile_name"),
					),
				},
			},
			"policy_name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The linked forwarding policy name to match exactly, ignoring case.",
			},
			"policy_link_id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The forwarding policy link ID used in `/networkAccess/forwardingProfiles/{forwardingProfileId}/policies/{policyLinkId}`.",
			},
			"priority": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "The forwarding policy link priority.",
			},
			"state": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The forwarding policy link state.",
			},
			"version": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The forwarding policy link version.",
			},
			"policy_id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The linked forwarding policy ID. Use this value as `forwarding_policy_id` when managing forwarding policy rules.",
			},
			"policy_description": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The linked forwarding policy description.",
			},
			"policy_version": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The linked forwarding policy version.",
			},
			"private_access_app_id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The linked private access app ID when returned by Graph.",
			},
			"timeouts": commonschema.ResourceTimeouts(ctx),
		},
	}
}
