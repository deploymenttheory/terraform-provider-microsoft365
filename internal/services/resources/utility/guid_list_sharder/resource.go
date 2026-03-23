package utilityGuidListSharder

import (
	"context"
	"regexp"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	planmodifiers "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/plan_modifiers"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/validate/attribute"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	ResourceName  = "microsoft365_utility_guid_list_sharder"
	CreateTimeout = 180
	ReadTimeout   = 180
	UpdateTimeout = 180
	DeleteTimeout = 180
)

var (
	_ resource.Resource              = &GuidListSharderResource{}
	_ resource.ResourceWithConfigure = &GuidListSharderResource{}
)

func NewGuidListSharderResource() resource.Resource {
	return &GuidListSharderResource{
		ReadPermissions: []string{
			"Application.Read.All",
			"Device.Read.All",
			"Directory.Read.All",
			"Group.Read.All",
			"User.Read.All",
			"User.ReadBasic.All",
		},
	}
}

type GuidListSharderResource struct {
	client          *msgraphbetasdk.GraphServiceClient
	ReadPermissions []string
}

func (r *GuidListSharderResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = ResourceName
}

func (r *GuidListSharderResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, ResourceName)
}

func (r *GuidListSharderResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Retrieves object IDs (GUIDs) from Microsoft Graph API and distributes them into configurable shards for progressive rollouts and phased deployments. " +
			"Queries `/users`, `/devices`, `/applications`, or `/groups/{id}/members` endpoints with optional OData filtering, then applies sharding strategies (random, sequential, or percentage-based) " +
			"to distribute results. Output shards are sets that can be directly used in conditional access policies, groups, and other resources requiring object ID collections.\n\n" +
			"Unlike a datasource, this resource stores shard assignments in Terraform state. When `recalculate_on_next_run = false`, " +
			"the stored assignments are returned unchanged on every plan and apply — preventing membership churn from causing reassignments. " +
			"Set `recalculate_on_next_run = true` and run `terraform apply` to recompute shards from the current Graph API member list.\n\n" +
			"**API Endpoints:** `GET /users`, `GET /devices`, `GET /applications`, `GET /groups/{id}/members` (with pagination and `ConsistencyLevel: eventual` header)\n\n" +
			"**Common Use Cases:** MFA rollouts, Windows Update rings, conditional access pilots, application-based policies, group splitting, A/B testing",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
				MarkdownDescription: "The ID of this resource.",
			},
			"recalculate_on_next_run": schema.BoolAttribute{
				Required: true,
				MarkdownDescription: "Controls whether shard assignments are recomputed during Terraform plan/refresh and on `terraform apply` when configuration changes.\n\n" +
					"**`false` (recommended default):** Shard assignments are locked in state. No Graph API call is made during plan or apply, " +
					"and membership changes in your tenant (users added, removed) do not cause reassignments. " +
					"On the very first apply, when no prior state exists, assignments are always computed regardless of this value — " +
					"so you can safely set `false` from the outset without a two-step toggle.\n\n" +
					"**`true`:** Re-queries the Graph API and reruns the sharding algorithm on every plan refresh and every apply. " +
					"Use this only when you explicitly want to rebalance — for example after a large onboarding wave, a policy restructure, or a change to `shard_count` or `strategy`.\n\n" +
					"**Recommended workflow:** Set to `false` from day one. Initial assignments are computed automatically on the first apply. " +
					"Change to `true` (and run `terraform apply`) only when you intentionally want to reshard. Set back to `false` afterwards to re-lock assignments.",
			},
			"resource_type": schema.StringAttribute{
				Required: true,
				MarkdownDescription: "The type of Microsoft Graph resource to query and shard. " +
					"`users` queries `/users` for user-based policies (MFA, conditional access). " +
					"`devices` queries `/devices` for device policies (Windows Updates, compliance). " +
					"`applications` queries `/applications` for app registrations. " +
					"`service_principals` queries `/servicePrincipals` (enterprise apps) for application-based conditional access policies. " +
					"`group_members` queries `/groups/{id}/members` to split existing group membership (requires `group_id`).",
				Validators: []validator.String{
					stringvalidator.OneOf("users", "devices", "applications", "service_principals", "group_members"),
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
			"odata_filter": schema.StringAttribute{
				Optional: true,
				MarkdownDescription: "Optional OData filter expression applied at the API level before sharding. " +
					"**Users:** `accountEnabled eq true` (active accounts only), `userType eq 'Member'` (exclude guests). " +
					"**Devices:** `operatingSystem eq 'Windows'` (Windows devices only). " +
					"**Service Principals:** `startswith(displayName, 'Microsoft')` (Microsoft apps), `appId eq 'guid'` (specific app). " +
					"Leave empty to query all resources without filtering.",
			},
			"shard_count": schema.Int64Attribute{
				Optional: true,
				MarkdownDescription: "Number of equally-sized shards to create (minimum 1). " +
					"Use with `round-robin` strategy. Conflicts with `shard_percentages` and `shard_sizes`. " +
					"Creates shards named `shard_0`, `shard_1`, ..., `shard_N-1`.",
				Validators: []validator.Int64{
					int64validator.AtLeast(1),
					int64validator.ExactlyOneOf(path.MatchRoot("shard_percentages"), path.MatchRoot("shard_sizes")),
				},
			},
			"shard_percentages": schema.ListAttribute{
				ElementType: types.Int64Type,
				Optional:    true,
				MarkdownDescription: "List of percentages for custom-sized shards. Use with `percentage` strategy. Conflicts with `shard_count` and `shard_sizes`. " +
					"Values must be non-negative integers that sum to exactly 100. " +
					"Example: `[10, 30, 60]` creates 10% pilot, 30% broader pilot, 60% full rollout.",
				Validators: []validator.List{
					listvalidator.SizeAtLeast(1),
					listvalidator.ValueInt64sAre(int64validator.AtLeast(0)),
					listvalidator.ExactlyOneOf(path.MatchRoot("shard_count"), path.MatchRoot("shard_sizes")),
					attribute.Int64ListSumEquals(100),
				},
			},
			"shard_sizes": schema.ListAttribute{
				ElementType: types.Int64Type,
				Optional:    true,
				MarkdownDescription: "List of absolute shard sizes (exact number of GUIDs per shard). Use with `size` strategy. Conflicts with `shard_count` and `shard_percentages`. " +
					"Values must be positive integers or -1 (which means 'all remaining'). Only the last element can be -1. " +
					"Example: `[50, 200, -1]` creates 50 pilot users, 200 broader rollout, remainder for full deployment.",
				Validators: []validator.List{
					listvalidator.SizeAtLeast(1),
					listvalidator.ValueInt64sAre(int64validator.Any(
						int64validator.AtLeast(1),
						int64validator.OneOf(-1),
					)),
					listvalidator.ExactlyOneOf(path.MatchRoot("shard_count"), path.MatchRoot("shard_percentages")),
				},
			},
			"strategy": schema.StringAttribute{
				Required: true,
				MarkdownDescription: "The distribution strategy for sharding GUIDs. " +
					"`round-robin` distributes in circular order (guarantees equal sizes, optional seed for reproducibility). " +
					"`percentage` distributes by specified percentages (requires `shard_percentages`, optional seed for reproducibility). " +
					"`size` distributes by absolute sizes (requires `shard_sizes`, optional seed for reproducibility). " +
					"`rendezvous` uses Highest Random Weight algorithm (always deterministic, minimal disruption when shard count changes).",
				Validators: []validator.String{
					stringvalidator.OneOf("round-robin", "percentage", "size", "rendezvous"),
				},
			},
			"seed": schema.StringAttribute{
				Optional: true,
				MarkdownDescription: "Optional seed value for deterministic distribution. When provided, makes results reproducible across Terraform runs for the same input set. " +
					"**`round-robin`**: No seed = uses API order (may change). With seed = shuffles deterministically first, then round-robin. " +
					"**`percentage`/`size`**: Same shuffle behaviour as round-robin. " +
					"**`rendezvous`**: Always deterministic. Seed affects which shard each GUID is assigned to. " +
					"Use different seeds for different rollouts to vary pilot burden distribution.",
			},
			"shards": schema.MapAttribute{
				ElementType: types.SetType{ElemType: types.StringType},
				Computed:    true,
				MarkdownDescription: "Computed map of shard names (`shard_0`, `shard_1`, ...) to sets of GUIDs. " +
					"Each value is a `set(string)` type, directly compatible with resource attributes expecting object ID sets. " +
					"Access with `resource.example.shards[\"shard_0\"]`, check size with `length(resource.example.shards[\"shard_0\"])`.",
			},
			"timeouts": commonschema.ResourceTimeouts(ctx),
		},
	}
}
