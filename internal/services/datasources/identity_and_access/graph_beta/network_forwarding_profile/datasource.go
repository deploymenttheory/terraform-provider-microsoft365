package graphBetaNetworkForwardingProfile

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
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
			"filter_type": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Type of filter to apply. Valid values are `all`, `id`, `name`, and `traffic_forwarding_type`.",
				Validators: []validator.String{
					stringvalidator.OneOf("all", "id", "name", "traffic_forwarding_type"),
				},
			},
			"filter_value": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Filter value. Required unless `filter_type` is `all`.",
			},
			"items": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "Forwarding profiles matching the filter.",
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
							MarkdownDescription: "Forwarding policy links associated with the profile. The policy link ID and policy ID are distinct Graph identifiers.",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"policy_link_id":          schema.StringAttribute{Computed: true, MarkdownDescription: "The forwarding policy link ID."},
									"priority":                schema.Int64Attribute{Computed: true, MarkdownDescription: "The forwarding policy link priority."},
									"state":                   schema.StringAttribute{Computed: true, MarkdownDescription: "The forwarding policy link state."},
									"version":                 schema.StringAttribute{Computed: true, MarkdownDescription: "The forwarding policy link version."},
									"policy_id":               schema.StringAttribute{Computed: true, MarkdownDescription: "The linked forwarding policy ID."},
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
