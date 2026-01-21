package graphBetaGroup

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
	"github.com/hashicorp/terraform-plugin-framework/types"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	DataSourceName = "microsoft365_graph_beta_groups_group"
	ReadTimeout    = 180
)

var (
	_ datasource.DataSource              = &GroupDataSource{}
	_ datasource.DataSourceWithConfigure = &GroupDataSource{}
)

func NewGroupDataSource() datasource.DataSource {
	return &GroupDataSource{
		ReadPermissions: []string{
			"Group.Read.All",
			"Directory.Read.All",
		},
	}
}

type GroupDataSource struct {
	client *msgraphbetasdk.GraphServiceClient

	ReadPermissions []string
}

func (d *GroupDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = DataSourceName
}

func (d *GroupDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.client = client.SetGraphBetaClientForDataSource(ctx, req, resp, DataSourceName)
}

func (d *GroupDataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Gets information about a Microsoft Entra ID (Azure AD) group.\n\n" +
			"See the [Microsoft Graph API documentation](https://learn.microsoft.com/en-us/graph/api/group-get?view=graph-rest-beta) for more details.\n\n" +
			"## API Permissions\n\n" +
			"The following API permissions are required:\n\n" +
			"- `Group.Read.All` or `Directory.Read.All` when authenticated with a service principal.\n" +
			"- No additional roles required when authenticated with a user principal.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The unique identifier for the group.",
			},
			"display_name": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The display name for the group.",
				Validators: []validator.String{
					stringvalidator.AtLeastOneOf(
						path.MatchRoot("object_id"),
						path.MatchRoot("mail_nickname"),
						path.MatchRoot("odata_query"),
					),
					stringvalidator.ConflictsWith(
						path.MatchRoot("object_id"),
						path.MatchRoot("mail_nickname"),
						path.MatchRoot("odata_query"),
					),
				},
			},
			"object_id": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The object ID of the group. One of `object_id`, `display_name`, `mail_nickname`, or `odata_query` must be specified.",
				Validators: []validator.String{
					stringvalidator.AtLeastOneOf(
						path.MatchRoot("display_name"),
						path.MatchRoot("mail_nickname"),
						path.MatchRoot("odata_query"),
					),
					stringvalidator.ConflictsWith(
						path.MatchRoot("display_name"),
						path.MatchRoot("mail_nickname"),
						path.MatchRoot("odata_query"),
					),
					stringvalidator.RegexMatches(
						regexp.MustCompile(constants.GuidRegex),
						"must be a valid GUID in the format 00000000-0000-0000-0000-000000000000",
					),
				},
			},
			"mail_nickname": schema.StringAttribute{
				Optional: true,
				Computed: true,
				MarkdownDescription: "The mail alias for the group, unique in the organisation. " +
					"One of `object_id`, `display_name`, or `mail_nickname` must be specified.",
				Validators: []validator.String{
					stringvalidator.AtLeastOneOf(
						path.MatchRoot("object_id"),
						path.MatchRoot("display_name"),
						path.MatchRoot("odata_query"),
					),
					stringvalidator.ConflictsWith(
						path.MatchRoot("object_id"),
						path.MatchRoot("display_name"),
						path.MatchRoot("odata_query"),
					),
				},
			},
			"odata_query": schema.StringAttribute{
				Optional: true,
				MarkdownDescription: "Custom OData filter query. " +
					"Use this for advanced filtering when the standard lookup attributes don't meet your needs. " +
					"Cannot be combined with `object_id`, `display_name`, or `mail_nickname`. " +
					"Example: `displayName eq 'My Group' and securityEnabled eq true`",
				Validators: []validator.String{
					stringvalidator.AtLeastOneOf(
						path.MatchRoot("object_id"),
						path.MatchRoot("display_name"),
						path.MatchRoot("mail_nickname"),
					),
					stringvalidator.ConflictsWith(
						path.MatchRoot("object_id"),
						path.MatchRoot("display_name"),
						path.MatchRoot("mail_nickname"),
					),
				},
			},
			"mail_enabled": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				MarkdownDescription: "Whether the group is mail-enabled. " +
					"Can be used as an additional filter when combined with other lookup attributes.",
			},
			"security_enabled": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				MarkdownDescription: "Whether the group is a security group. " +
					"Can be used as an additional filter when combined with other lookup attributes.",
			},
			"description": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The optional description of the group.",
			},
			"classification": schema.StringAttribute{
				Computed: true,
				MarkdownDescription: "A classification for the group (such as low, medium or high business impact). " +
					"Valid values are defined by creating a ClassificationList setting value in the directory.",
			},
			"group_types": schema.SetAttribute{
				ElementType: types.StringType,
				Computed:    true,
				MarkdownDescription: "A list of group types configured for the group. Possible values include:\n" +
					"  - `DynamicMembership`: Denotes a group with dynamic membership\n" +
					"  - `Unified`: Specifies a Microsoft 365 group",
			},
			"visibility": schema.StringAttribute{
				Computed: true,
				MarkdownDescription: "The group join policy and group content visibility. Possible values are:\n" +
					"  - `Private`: Only members can view content\n" +
					"  - `Public`: Anyone can view content\n" +
					"  - `Hiddenmembership`: Only members can see membership (Microsoft 365 groups only)",
			},
			"assignable_to_role": schema.BoolAttribute{
				Computed: true,
				MarkdownDescription: "Indicates whether this group can be assigned to an Azure AD role. " +
					"Can only be set during group creation and cannot be changed afterwards.",
			},
			"dynamic_membership_enabled": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "Whether the group has dynamic membership enabled.",
			},
			"membership_rule": schema.StringAttribute{
				Computed: true,
				MarkdownDescription: "The rule that determines members for a dynamic membership group. " +
					"Only populated when `dynamic_membership_enabled` is `true`.",
			},
			"membership_rule_processing_state": schema.StringAttribute{
				Computed: true,
				MarkdownDescription: "Indicates whether the dynamic membership is processing. Possible values are:\n" +
					"  - `On`: Dynamic membership is active\n" +
					"  - `Paused`: Dynamic membership is paused",
			},
			"created_date_time": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The date and time when the group was created in RFC3339 format.",
			},
			"mail": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The SMTP address for the group.",
			},
			"proxy_addresses": schema.SetAttribute{
				ElementType:         types.StringType,
				Computed:            true,
				MarkdownDescription: "Email addresses for the group that direct to the same group mailbox.",
			},
			"assigned_licenses": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "The licenses that are assigned to the group for group-based licensing.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"sku_id": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The unique identifier (GUID) for the service SKU.",
						},
						"disabled_plans": schema.SetAttribute{
							ElementType:         types.StringType,
							Computed:            true,
							MarkdownDescription: "A collection of the unique identifiers for plans that have been disabled.",
						},
					},
				},
			},
			"has_members_with_license_errors": schema.BoolAttribute{
				Computed: true,
				MarkdownDescription: "Indicates whether there are members in this group that have license errors from group-based license assignment. " +
					"This property is never returned on a GET operation unless explicitly requested via $select.",
			},
			"hide_from_address_lists": schema.BoolAttribute{
				Computed: true,
				MarkdownDescription: "True if the group is not displayed in certain parts of the Outlook UI: the Address Book, " +
					"address lists for selecting message recipients, and the Browse Groups dialog for searching groups; " +
					"otherwise false. Default value is false.",
			},
			"hide_from_outlook_clients": schema.BoolAttribute{
				Computed: true,
				MarkdownDescription: "True if the group is not displayed in Outlook clients, such as Outlook for Windows and Outlook on the web; " +
					"otherwise false. Default value is false.",
			},
			"onpremises_sync_enabled": schema.BoolAttribute{
				Computed: true,
				MarkdownDescription: "Whether this group is synchronised from an on-premises directory. " +
					"Possible values are `true` (synced), `false` (no longer synced), or null (never synced).",
			},
			"onpremises_last_sync_date_time": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The last time the group was synced from the on-premises directory in RFC3339 format.",
			},
			"onpremises_sam_account_name": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The on-premises SAM account name, synchronised from the on-premises directory when Azure AD Connect is used.",
			},
			"onpremises_domain_name": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The on-premises FQDN (dnsDomainName), synchronised from the on-premises directory when Azure AD Connect is used.",
			},
			"onpremises_netbios_name": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The on-premises NetBIOS name, synchronised from the on-premises directory when Azure AD Connect is used.",
			},
			"onpremises_security_identifier": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The on-premises security identifier (SID), synchronised from the on-premises directory when Azure AD Connect is used.",
			},
			"members": schema.SetAttribute{
				ElementType:         types.StringType,
				Computed:            true,
				MarkdownDescription: "List of object IDs of the group members.",
			},
			"owners": schema.SetAttribute{
				ElementType:         types.StringType,
				Computed:            true,
				MarkdownDescription: "List of object IDs of the group owners.",
			},
			"timeouts": commonschema.DatasourceTimeouts(ctx),
		},
	}
}
