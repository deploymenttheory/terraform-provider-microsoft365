package graphBetaMobileAppRelationship

import (
	"context"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
)

const (
	DataSourceName = "microsoft365_graph_beta_device_and_app_management_mobile_app_relationship"
	ReadTimeout    = 180
)

var (
	_ datasource.DataSource              = &MobileAppRelationshipDataSource{}
	_ datasource.DataSourceWithConfigure = &MobileAppRelationshipDataSource{}
)

func NewMobileAppRelationshipDataSource() datasource.DataSource {
	return &MobileAppRelationshipDataSource{
		ReadPermissions: []string{
			"DeviceManagementApps.Read.All",
		},
	}
}

type MobileAppRelationshipDataSource struct {
	client          *msgraphbetasdk.GraphServiceClient
	ReadPermissions []string
}

func (d *MobileAppRelationshipDataSource) Metadata(
	ctx context.Context,
	req datasource.MetadataRequest,
	resp *datasource.MetadataResponse,
) {
	resp.TypeName = DataSourceName
}

func (d *MobileAppRelationshipDataSource) Configure(
	ctx context.Context,
	req datasource.ConfigureRequest,
	resp *datasource.ConfigureResponse,
) {
	d.client = client.SetGraphBetaClientForDataSource(ctx, req, resp, DataSourceName)
}

// lookupAttributes are mutually exclusive; exactly one selects the relationship.
var lookupAttributes = []string{
	"relationship_id",
	"target_id",
	"target_display_name",
	"odata_query",
}

// lookupValidators returns the exactly-one-of validators for a lookup attribute.
func lookupValidators(self string) []validator.String {
	others := make([]path.Expression, 0, len(lookupAttributes)-1)
	for _, name := range lookupAttributes {
		if name != self {
			others = append(others, path.MatchRoot(name))
		}
	}
	return []validator.String{
		stringvalidator.AtLeastOneOf(others...),
		stringvalidator.ConflictsWith(others...),
	}
}

func guidValidator() validator.String {
	return stringvalidator.RegexMatches(
		regexp.MustCompile(constants.GuidRegex),
		"must be a valid GUID in the format 00000000-0000-0000-0000-000000000000",
	)
}

func (d *MobileAppRelationshipDataSource) Schema(
	ctx context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Retrieves a single Intune mobile app relationship (dependency or supersedence) using the " +
			"`/deviceAppManagement/mobileApps/{id}/relationships` endpoint. Look the relationship up by `relationship_id`, " +
			"or search the relationships of the app named in `source_id` by `target_id`, `target_display_name`, or a custom `odata_query`. " +
			"The lookup must match exactly one relationship.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The unique identifier of the relationship, in the format `{source_id}_{target_id}`.",
			},
			"source_id": schema.StringAttribute{
				Optional: true,
				Computed: true,
				MarkdownDescription: "The ID of the app whose relationships are searched. Required with `target_id`, " +
					"`target_display_name`, or `odata_query`. Optional with `relationship_id`, where it is derived from the ID.",
				Validators: []validator.String{
					guidValidator(),
					stringvalidator.AtLeastOneOf(
						path.MatchRoot("relationship_id"),
						path.MatchRoot("target_id"),
						path.MatchRoot("target_display_name"),
						path.MatchRoot("odata_query"),
					),
				},
			},
			"relationship_id": schema.StringAttribute{
				Optional: true,
				Computed: true,
				MarkdownDescription: "The ID of the relationship, in the format `{source_id}_{target_id}`. " +
					"One of `relationship_id`, `target_id`, `target_display_name`, or `odata_query` must be specified.",
				Validators: lookupValidators("relationship_id"),
			},
			"target_id": schema.StringAttribute{
				Optional: true,
				Computed: true,
				MarkdownDescription: "The ID of the related app. Requires `source_id`. " +
					"One of `relationship_id`, `target_id`, `target_display_name`, or `odata_query` must be specified.",
				Validators: append(lookupValidators("target_id"),
					guidValidator(),
					stringvalidator.AlsoRequires(path.MatchRoot("source_id")),
				),
			},
			"target_display_name": schema.StringAttribute{
				Optional: true,
				Computed: true,
				MarkdownDescription: "The display name of the related app (case-insensitive). Requires `source_id`. " +
					"One of `relationship_id`, `target_id`, `target_display_name`, or `odata_query` must be specified.",
				Validators: append(lookupValidators("target_display_name"),
					stringvalidator.AlsoRequires(path.MatchRoot("source_id")),
				),
			},
			"odata_query": schema.StringAttribute{
				Optional: true,
				MarkdownDescription: "Custom OData `$filter` applied to the relationships of the `source_id` app. Requires `source_id`. " +
					"Use this for advanced filtering when the standard lookup attributes don't meet your needs. " +
					"Example: `targetType eq 'child'`",
				Validators: append(lookupValidators("odata_query"),
					stringvalidator.AlsoRequires(path.MatchRoot("source_id")),
				),
			},
			"target_type": schema.StringAttribute{
				Optional: true,
				Computed: true,
				MarkdownDescription: "Whether the target app is a `child` (the source app depends on or supersedes it) or a " +
					"`parent` (the target app depends on or supersedes the source app). " +
					"Can be used as an additional filter when combined with the lookup attributes.",
				Validators: []validator.String{
					stringvalidator.OneOf("child", "parent"),
				},
			},
			"relationship_type": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The kind of relationship: `dependency` or `supersedence`.",
			},
			"target_display_version": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The display version of the target app.",
			},
			"target_publisher": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The publisher of the target app.",
			},
			"target_publisher_display_name": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The publisher display name of the target app.",
			},
			"source_display_name": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The display name of the source app.",
			},
			"source_display_version": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The display version of the source app.",
			},
			"source_publisher_display_name": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The publisher display name of the source app.",
			},
			"dependency_type": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "For dependencies, how the dependency is handled: `detect` or `autoInstall`.",
			},
			"dependent_app_count": schema.Int32Attribute{
				Computed:            true,
				MarkdownDescription: "For dependencies, the total number of apps that directly or indirectly depend on the parent app.",
			},
			"depends_on_app_count": schema.Int32Attribute{
				Computed:            true,
				MarkdownDescription: "For dependencies, the total number of apps the child app directly or indirectly depends on.",
			},
			"supersedence_type": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "For supersedences, whether the superseded app is `update`d or `replace`d.",
			},
			"superseded_app_count": schema.Int32Attribute{
				Computed:            true,
				MarkdownDescription: "For supersedences, the total number of apps directly or indirectly superseded by the child app.",
			},
			"superseding_app_count": schema.Int32Attribute{
				Computed:            true,
				MarkdownDescription: "For supersedences, the total number of apps directly or indirectly superseding the parent app.",
			},
			"timeouts": commonschema.DatasourceTimeouts(ctx),
		},
	}
}
