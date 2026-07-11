package graphBetaNetworkForwardingProfile

import (
	"context"
	"regexp"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	"github.com/hashicorp/terraform-plugin-framework-validators/boolvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	DataSourceName = "microsoft365_graph_beta_identity_and_access_network_forwarding_profile"
	ReadTimeout    = 180
)

var (
	_ datasource.DataSource              = &NetworkForwardingProfileDataSource{}
	_ datasource.DataSourceWithConfigure = &NetworkForwardingProfileDataSource{}
)

func NewNetworkForwardingProfileDataSource() datasource.DataSource {
	return &NetworkForwardingProfileDataSource{
		ReadPermissions: []string{
			"NetworkAccess.Read.All",
		},
	}
}

type NetworkForwardingProfileDataSource struct {
	client          *msgraphbetasdk.GraphServiceClient
	ReadPermissions []string
}

func (d *NetworkForwardingProfileDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = DataSourceName
}

func (d *NetworkForwardingProfileDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.client = client.SetGraphBetaClientForDataSource(ctx, req, resp, DataSourceName)
}

func (d *NetworkForwardingProfileDataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Retrieves Microsoft Entra Global Secure Access forwarding profiles using Microsoft Graph beta `/networkAccess/forwardingProfiles` and expands associated forwarding policy links.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The unique identifier for the data source operation.",
			},
			"forwarding_profile_id": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The forwarding profile ID. Conflicts with `name`, `traffic_forwarding_type`, and `list_all`.",
				Validators: []validator.String{
					stringvalidator.RegexMatches(regexp.MustCompile(constants.GuidRegex), "must be a valid UUID"),
					stringvalidator.ConflictsWith(
						path.MatchRoot("name"),
						path.MatchRoot("traffic_forwarding_type"),
						path.MatchRoot("list_all"),
					),
					stringvalidator.AtLeastOneOf(
						path.MatchRoot("forwarding_profile_id"),
						path.MatchRoot("name"),
						path.MatchRoot("traffic_forwarding_type"),
						path.MatchRoot("list_all"),
					),
				},
			},
			"name": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The forwarding profile name to match exactly, ignoring case. Conflicts with `forwarding_profile_id`, `traffic_forwarding_type`, and `list_all`.",
				Validators: []validator.String{
					stringvalidator.ConflictsWith(
						path.MatchRoot("forwarding_profile_id"),
						path.MatchRoot("traffic_forwarding_type"),
						path.MatchRoot("list_all"),
					),
				},
			},
			"traffic_forwarding_type": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The forwarding profile traffic type. Valid values are `internet`, `m365`, and `private`. Conflicts with `forwarding_profile_id`, `name`, and `list_all`.",
				Validators: []validator.String{
					stringvalidator.OneOf("internet", "m365", "private"),
					stringvalidator.ConflictsWith(
						path.MatchRoot("forwarding_profile_id"),
						path.MatchRoot("name"),
						path.MatchRoot("list_all"),
					),
				},
			},
			"list_all": schema.BoolAttribute{
				Optional:            true,
				MarkdownDescription: "Retrieve all forwarding profiles. Conflicts with `forwarding_profile_id`, `name`, and `traffic_forwarding_type`.",
				Validators: []validator.Bool{
					boolvalidator.ConflictsWith(
						path.MatchRoot("forwarding_profile_id"),
						path.MatchRoot("name"),
						path.MatchRoot("traffic_forwarding_type"),
					),
				},
			},
			"items": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "Forwarding profiles matching the query criteria.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id":                       schema.StringAttribute{Computed: true, MarkdownDescription: "The forwarding profile ID."},
						"name":                     schema.StringAttribute{Computed: true, MarkdownDescription: "The forwarding profile name."},
						"description":              schema.StringAttribute{Computed: true, MarkdownDescription: "The forwarding profile description."},
						"state":                    schema.StringAttribute{Computed: true, MarkdownDescription: "The forwarding profile state."},
						"version":                  schema.StringAttribute{Computed: true, MarkdownDescription: "The forwarding profile version."},
						"last_modified_date_time":  schema.StringAttribute{Computed: true, MarkdownDescription: "The last modified timestamp."},
						"traffic_forwarding_type":  schema.StringAttribute{Computed: true, MarkdownDescription: "The traffic forwarding type, such as `internet`, `m365`, or `private`."},
						"priority":                 schema.Int32Attribute{Computed: true, MarkdownDescription: "The forwarding profile priority."},
						"is_custom_profile":        schema.BoolAttribute{Computed: true, MarkdownDescription: "Whether this is a custom forwarding profile."},
						"client_fallback_action":   schema.StringAttribute{Computed: true, MarkdownDescription: "The client fallback action."},
						"service_principal_app_id": schema.StringAttribute{Computed: true, MarkdownDescription: "The associated service principal application ID."},
						"service_principal_id":     schema.StringAttribute{Computed: true, MarkdownDescription: "The associated service principal object ID."},
						"policies": schema.ListNestedAttribute{
							Computed:            true,
							MarkdownDescription: "Forwarding policy links associated with the profile. `policy_link_id` is the Graph policyLink object ID used by the forwarding profile policies endpoint; `policy_id` is the linked forwarding policy ID.",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"policy_link_id":          schema.StringAttribute{Computed: true, MarkdownDescription: "The forwarding policy link ID used in `/networkAccess/forwardingProfiles/{forwardingProfileId}/policies/{policyLinkId}`."},
									"priority":                schema.Int64Attribute{Computed: true, MarkdownDescription: "The forwarding policy link priority."},
									"state":                   schema.StringAttribute{Computed: true, MarkdownDescription: "The forwarding policy link state."},
									"version":                 schema.StringAttribute{Computed: true, MarkdownDescription: "The forwarding policy link version."},
									"policy_id":               schema.StringAttribute{Computed: true, MarkdownDescription: "The linked forwarding policy ID. Use this value as `forwarding_policy_id` when managing forwarding policy rules."},
									"policy_name":             schema.StringAttribute{Computed: true, MarkdownDescription: "The linked forwarding policy name."},
									"policy_description":      schema.StringAttribute{Computed: true, MarkdownDescription: "The linked forwarding policy description."},
									"policy_version":          schema.StringAttribute{Computed: true, MarkdownDescription: "The linked forwarding policy version."},
									"traffic_forwarding_type": schema.StringAttribute{Computed: true, MarkdownDescription: "The linked forwarding policy traffic forwarding type."},
									"private_access_app_id":   schema.StringAttribute{Computed: true, MarkdownDescription: "The linked private access app ID when returned by Graph."},
								},
							},
						},
					},
				},
			},
			"timeouts": commonschema.ResourceTimeouts(ctx),
		},
	}
}
