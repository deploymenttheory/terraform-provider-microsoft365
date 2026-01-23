package utilityGuidListSharder

import (
	"context"
	"regexp"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/validate/attribute"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	DataSourceName = "microsoft365_utility_guid_list_sharder"
	ReadTimeout    = 180
)

var (
	_ datasource.DataSource              = &guidListSharderDataSource{}
	_ datasource.DataSourceWithConfigure = &guidListSharderDataSource{}
)

func NewGuidListSharderDataSource() datasource.DataSource {
	return &guidListSharderDataSource{
		ReadPermissions: []string{
			"User.Read.All",
			"Group.Read.All",
			"Device.Read.All",
			"Directory.Read.All",
		},
	}
}

type guidListSharderDataSource struct {
	client *msgraphbetasdk.GraphServiceClient

	ReadPermissions []string
}

func (d *guidListSharderDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = DataSourceName
}

func (d *guidListSharderDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.client = client.SetGraphBetaClientForDataSource(ctx, req, resp, DataSourceName)
}

func (d *guidListSharderDataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Retrieves object IDs (GUIDs) from Microsoft Graph API and distributes them into configurable shards for progressive rollouts and phased deployments. " +
			"Queries `/users`, `/devices`, or `/groups/{id}/members` endpoints with optional OData filtering, then applies sharding strategies (random, sequential, or percentage-based) " +
			"to distribute results. Output shards are sets that can be directly used in conditional access policies, groups, and other resources requiring object ID collections.\n\n" +
			"**API Endpoints:** `GET /users`, `GET /devices`, `GET /groups/{id}/members` (with pagination and `ConsistencyLevel: eventual` header)\n\n" +
			"**Common Use Cases:** MFA rollouts, Windows Update rings, conditional access pilots, group splitting, A/B testing\n\n" +
			"For detailed examples and best practices, see the [Progressive Rollout with GUID List Sharder](https://registry.terraform.io/providers/deploymenttheory/microsoft365/latest/docs/guides/progressive_rollout_with_guid_list_sharder) guide.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The ID of this resource.",
			},
			"resource_type": schema.StringAttribute{
				Required: true,
				MarkdownDescription: "The type of Microsoft Graph resource to query and shard. " +
					"`users` queries `/users` for user-based policies (MFA, conditional access). " +
					"`devices` queries `/devices` for device policies (Windows Updates, compliance). " +
					"`group_members` queries `/groups/{id}/members` to split existing group membership (requires `group_id`).",
				Validators: []validator.String{
					stringvalidator.OneOf("users", "devices", "group_members"),
				},
			},
			"group_id": schema.StringAttribute{
				Optional: true,
				MarkdownDescription: "The object ID of the group to query members from. " +
					"Required when `resource_type = \"group_members\"`, ignored otherwise. " +
					"Use this to split an existing group's membership into multiple new groups for targeted policy application.",
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(constants.GuidRegex),
						"must be a valid GUID in the format 00000000-0000-0000-0000-000000000000",
					),
					attribute.RequiredWhenEquals("resource_type", types.StringValue("group_members")),
				},
			},
			"odata_query": schema.StringAttribute{
				Optional: true,
				MarkdownDescription: "Optional OData filter applied at the API level before sharding. " +
					"Common examples: `$filter=accountEnabled eq true` (active accounts only), " +
					"`$filter=operatingSystem eq 'Windows'` (Windows devices), " +
					"`$filter=userType eq 'Member'` (exclude guests). " +
					"Leave empty to query all resources without filtering.",
			},
			"shard_count": schema.Int64Attribute{
				Optional: true,
				MarkdownDescription: "Number of equally-sized shards to create (minimum 1). " +
					"Use with `round-robin` strategy. Conflicts with `shard_percentages`. " +
					"Creates shards named `shard_0`, `shard_1`, ..., `shard_N-1`. " +
					"For custom-sized shards (e.g., 10% pilot, 30% broader, 60% full), use `shard_percentages` with `percentage` strategy instead.",
				Validators: []validator.Int64{
					int64validator.AtLeast(1),
					int64validator.ExactlyOneOf(path.MatchRoot("shard_percentages")),
				},
			},
			"shard_percentages": schema.ListAttribute{
				ElementType: types.Int64Type,
				Optional:    true,
				MarkdownDescription: "List of percentages for custom-sized shards. Use with `percentage` strategy. Conflicts with `shard_count`. " +
					"Values must be non-negative integers that sum to exactly 100. " +
					"Example: `[10, 30, 60]` creates 10% pilot, 30% broader pilot, 60% full rollout. " +
					"Common patterns: `[5, 15, 80]` (Windows Update rings), `[33, 33, 34]` (A/B/C testing). " +
					"Last shard receives all remaining GUIDs to prevent loss.",
				Validators: []validator.List{
					listvalidator.SizeAtLeast(1),
					listvalidator.ValueInt64sAre(int64validator.AtLeast(0)),
					listvalidator.ExactlyOneOf(path.MatchRoot("shard_count")),
					attribute.Int64ListSumEquals(100),
				},
			},
			"strategy": schema.StringAttribute{
				Required: true,
				MarkdownDescription: "The distribution strategy for sharding GUIDs. " +
					"`round-robin` distributes in circular order (guarantees equal sizes, optional seed for reproducibility). " +
					"`percentage` distributes by specified percentages (requires `shard_percentages`, optional seed for reproducibility). " +
					"See the [guide](https://registry.terraform.io/providers/deploymenttheory/microsoft365/latest/docs/guides/progressive_rollout_with_guid_list_sharder) for detailed comparison.",
				Validators: []validator.String{
					stringvalidator.OneOf("round-robin", "percentage"),
				},
			},
			"seed": schema.StringAttribute{
				Optional: true,
				MarkdownDescription: "Optional seed value for deterministic distribution. When provided, makes results reproducible across Terraform runs. " +
					"**`round-robin` strategy**: No seed = uses API order (may change). With seed = shuffles deterministically first, then applies round-robin (reproducible). " +
					"**`percentage` strategy**: No seed = uses API order (may change). With seed = shuffles deterministically first, then applies percentage split (reproducible). " +
					"Use different seeds for different rollouts to distribute pilot burden: User X might be in shard_0 for MFA but shard_2 for Windows Updates.",
			},
			"shards": schema.MapAttribute{
				ElementType: types.SetType{ElemType: types.StringType},
				Computed:    true,
				MarkdownDescription: "Computed map of shard names (`shard_0`, `shard_1`, ...) to sets of GUIDs. " +
					"Each value is a `set(string)` type, directly compatible with resource attributes expecting object ID sets " +
					"(e.g., `conditions.users.include_users` in conditional access policies, `group_members` in groups). " +
					"Access with `data.example.shards[\"shard_0\"]`, check size with `length(data.example.shards[\"shard_0\"])`.",
			},
			"timeouts": commonschema.DatasourceTimeouts(ctx),
		},
	}
}
