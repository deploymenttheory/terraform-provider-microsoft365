package graphBetaCloudPcs

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/validators"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	DataSourceName = "graph_beta_windows_365_cloud_pcs"
	ReadTimeout    = 180
)

var (
	_ datasource.DataSource              = &CloudPcsDataSource{}
	_ datasource.DataSourceWithConfigure = &CloudPcsDataSource{}
)

func NewCloudPcsDataSource() datasource.DataSource {
	return &CloudPcsDataSource{
		ReadPermissions: []string{
			"CloudPC.Read.All",
		},
	}
}

type CloudPcsDataSource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ProviderTypeName string
	TypeName         string
	ReadPermissions  []string
}

func (d *CloudPcsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + DataSourceName
}

func (d *CloudPcsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.client = client.SetGraphBetaClientForDataSource(ctx, req, resp, d.TypeName)
}

func (d *CloudPcsDataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Retrieves Cloud PCs from Microsoft Intune. Using the endpoint '/deviceManagement/virtualEndpoint/cloudPCs'. Supports filtering by all, id, display_name, user_principal_name, status, product_type, or odata.",
		Attributes: map[string]schema.Attribute{
			"filter_type": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Type of filter to apply. Valid values are: `all`, `id`, `display_name`, `odata`. Use 'all' to retrieve all Cloud PCs, 'id' to retrieve a specific Cloud PC by its unique identifier, 'display_name' to filter by name, or 'odata' to use advanced OData query parameters.",
				Validators: []validator.String{
					stringvalidator.OneOf("all", "id", "display_name", "odata"),
					validators.ODataParameterValidator("odata_filter", "odata_select", "odata_top", "odata_count"),
				},
			},
			"filter_value": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Value to filter by. Not required when filter_type is 'all'. For 'id', provide the Cloud PC ID. For other filters, provide the appropriate value to match.",
			},
			"odata_filter": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "OData $filter query parameter. Only applicable when filter_type is 'odata'. Example: \"status eq 'provisioned'\"",
			},
			"odata_select": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "OData $select query parameter to specify which fields to return. Only applicable when filter_type is 'odata'. Example: \"id,displayName,status\"",
			},
			"odata_top": schema.Int64Attribute{
				Optional:            true,
				MarkdownDescription: "OData $top query parameter to limit the number of items returned. Only applicable when filter_type is 'odata'.",
			},
			"odata_count": schema.BoolAttribute{
				Optional:            true,
				MarkdownDescription: "OData $count query parameter to include a count of items. Only applicable when filter_type is 'odata'.",
			},
			"items": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "The list of Cloud PCs that match the filter criteria.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The unique identifier for the Cloud PC.",
						},
						"display_name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The display name of the Cloud PC.",
						},
						"aad_device_id": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The Azure AD device ID associated with the Cloud PC.",
						},
						"image_display_name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The display name of the image used for the Cloud PC.",
						},
						"managed_device_id": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The managed device ID associated with the Cloud PC.",
						},
						"managed_device_name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The name of the managed device associated with the Cloud PC.",
						},
						"provisioning_policy_id": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The ID of the provisioning policy used for the Cloud PC.",
						},
						"provisioning_policy_name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The name of the provisioning policy used for the Cloud PC.",
						},
						"on_premises_connection_name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The name of the on-premises connection used for the Cloud PC.",
						},
						"service_plan_id": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The ID of the service plan associated with the Cloud PC.",
						},
						"service_plan_name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The name of the service plan associated with the Cloud PC.",
						},
						"service_plan_type": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The type of service plan associated with the Cloud PC.",
						},
						"status": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The current status of the Cloud PC (e.g., provisioned, provisioning, failed).",
						},
						"user_principal_name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The user principal name (UPN) of the user assigned to the Cloud PC.",
						},
						"last_modified_date_time": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The date and time when the Cloud PC was last modified.",
						},
						"status_detail_code": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The error/warning code associated with the Cloud PC status.",
						},
						"status_detail_message": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The status message associated with the error code.",
						},
						"grace_period_end_date_time": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The date and time when the grace period for the Cloud PC ends.",
						},
						"provisioning_type": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The type of provisioning used for the Cloud PC (e.g., dedicated).",
						},
						"device_region_name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The Azure region where the Cloud PC is deployed.",
						},
						"disk_encryption_state": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The disk encryption state of the Cloud PC (e.g., encryptedUsingPlatformManagedKey).",
						},
						"product_type": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The product type of the Cloud PC (e.g., enterprise).",
						},
						"user_account_type": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The account type of the user on provisioned Cloud PCs (e.g., standardUser, administrator).",
						},
						"enable_single_sign_on": schema.BoolAttribute{
							Computed:            true,
							MarkdownDescription: "Indicates whether single sign-on is enabled for the Cloud PC.",
						},
					},
				},
			},
			"timeouts": commonschema.Timeouts(ctx),
		},
	}
}
