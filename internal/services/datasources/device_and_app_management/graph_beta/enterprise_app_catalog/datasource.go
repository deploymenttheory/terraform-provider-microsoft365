package enterpriseappcatalog

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	datasourceName = "graph_beta_device_and_app_management_enterprise_app_catalog"
	ReadTimeout    = 180
)

var (
	// Basic datasource interface (CRUD operations)
	_ datasource.DataSource = &EnterpriseAppCatalogDataSource{}

	// Allows the datasource to be configured with the provider client
	_ datasource.DataSourceWithConfigure = &EnterpriseAppCatalogDataSource{}
)

func NewEnterpriseAppCatalogDataSource() datasource.DataSource {
	return &EnterpriseAppCatalogDataSource{
		ReadPermissions: []string{
			"DeviceManagementApps.Read.All",
		},
	}
}

type EnterpriseAppCatalogDataSource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ProviderTypeName string
	TypeName         string
	ReadPermissions  []string
}

// Metadata returns the datasource type name.
func (d *EnterpriseAppCatalogDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + datasourceName
}

// Configure configures the data source with the provider client
func (d *EnterpriseAppCatalogDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.client = client.SetGraphBetaClientForDataSource(ctx, req, resp, d.TypeName)
}

// Schema defines the schema for the data source
func (d *EnterpriseAppCatalogDataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Retrieves Enterprise App catalog packages. Leverages the `/deviceAppManagement/mobileAppCatalogPackages` endpoint to return the product reference" +
			"and leverages the `/deviceAppManagement/mobileApps/convertFromMobileAppCatalogPackage(mobileAppCatalogPackageId='00000000-0000-0000-0000-000000000000')` endpoint to return the app configuration." +
			" This data source provides information about applications available in the Intune Enterprise App catalog." +
			"'https://learn.microsoft.com/en-us/intune/intune-service/apps/apps-enterprise-app-management'",
		Attributes: map[string]schema.Attribute{
			"filter_type": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Type of filter to apply. Valid values are: `all`, `id`, `display_name`, `odata`.",
				Validators: []validator.String{
					stringvalidator.OneOf("all", "id", "display_name", "odata"),
				},
			},
			"filter_value": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Value to filter by. Not required when filter_type is 'all'.",
			},
			"odata_filter": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "OData $filter parameter for filtering results. Only used when filter_type is 'odata'.",
			},
			"odata_top": schema.Int32Attribute{
				Optional:            true,
				MarkdownDescription: "OData $top parameter to limit the number of results. Only used when filter_type is 'odata'.",
			},
			"odata_skip": schema.Int32Attribute{
				Optional:            true,
				MarkdownDescription: "OData $skip parameter for pagination. Only used when filter_type is 'odata'.",
			},
			"odata_select": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "OData $select parameter to specify which fields to include. Only used when filter_type is 'odata'.",
			},
			"odata_orderby": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "OData $orderby parameter to sort results. Only used when filter_type is 'odata'.",
			},
			"include_app_config": schema.BoolAttribute{
				Optional:            true,
				MarkdownDescription: "Whether to include detailed app configuration by calling the convertFromMobileAppCatalogPackage endpoint. This adds significant processing time.",
			},
			"items": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "List of mobile app catalog packages retrieved.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The unique identifier for the package.",
						},
						"product_id": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The product identifier for the package.",
						},
						"product_display_name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The display name of the product.",
						},
						"publisher_display_name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The display name of the publisher.",
						},
						"version_display_name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The display name of the version.",
						},
						"branch_display_name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The display name of the branch.",
						},
						"applicable_architectures": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The architectures this package is applicable for (e.g., x86, x64, arm64).",
						},
						"locales": schema.ListAttribute{
							Computed:            true,
							MarkdownDescription: "List of locales supported by this package.",
							ElementType:         types.StringType,
						},
						"package_auto_update_capable": schema.BoolAttribute{
							Computed:            true,
							MarkdownDescription: "Whether the package is capable of auto-updating.",
						},
						"app_config": schema.SingleNestedAttribute{
							Computed:            true,
							MarkdownDescription: "Detailed app configuration from the convertFromMobileAppCatalogPackage endpoint. Only populated if include_app_config is true.",
							Attributes: map[string]schema.Attribute{
								"odata_type": schema.StringAttribute{
									Computed:            true,
									MarkdownDescription: "The OData type of the app (e.g., #microsoft.graph.win32CatalogApp).",
								},
								"display_name": schema.StringAttribute{
									Computed:            true,
									MarkdownDescription: "The display name of the app.",
								},
								"description": schema.StringAttribute{
									Computed:            true,
									MarkdownDescription: "The description of the app.",
								},
								"publisher": schema.StringAttribute{
									Computed:            true,
									MarkdownDescription: "The publisher of the app.",
								},
								"developer": schema.StringAttribute{
									Computed:            true,
									MarkdownDescription: "The developer of the app.",
								},
								"privacy_information_url": schema.StringAttribute{
									Computed:            true,
									MarkdownDescription: "The privacy information URL of the app.",
								},
								"information_url": schema.StringAttribute{
									Computed:            true,
									MarkdownDescription: "The information URL of the app.",
								},
								"file_name": schema.StringAttribute{
									Computed:            true,
									MarkdownDescription: "The file name of the app installer.",
								},
								"size": schema.Int64Attribute{
									Computed:            true,
									MarkdownDescription: "The size of the app installer in bytes.",
								},
								"install_command_line": schema.StringAttribute{
									Computed:            true,
									MarkdownDescription: "The command line to install the app.",
								},
								"uninstall_command_line": schema.StringAttribute{
									Computed:            true,
									MarkdownDescription: "The command line to uninstall the app.",
								},
								"applicable_architectures": schema.StringAttribute{
									Computed:            true,
									MarkdownDescription: "The architectures this app is applicable for.",
								},
								"allowed_architectures": schema.StringAttribute{
									Computed:            true,
									MarkdownDescription: "The architectures this app is allowed to run on.",
								},
								"setup_file_path": schema.StringAttribute{
									Computed:            true,
									MarkdownDescription: "The path to the setup file.",
								},
								"min_supported_windows_release": schema.StringAttribute{
									Computed:            true,
									MarkdownDescription: "The minimum supported Windows release.",
								},
								"display_version": schema.StringAttribute{
									Computed:            true,
									MarkdownDescription: "The display version of the app.",
								},
								"allow_available_uninstall": schema.BoolAttribute{
									Computed:            true,
									MarkdownDescription: "Whether the app can be uninstalled when available.",
								},
								"rules": schema.ListNestedAttribute{
									Computed:            true,
									MarkdownDescription: "The detection and requirement rules for the app.",
									NestedObject: schema.NestedAttributeObject{
										Attributes: map[string]schema.Attribute{
											"odata_type": schema.StringAttribute{
												Computed:            true,
												MarkdownDescription: "The OData type of the rule.",
											},
											"rule_type": schema.StringAttribute{
												Computed:            true,
												MarkdownDescription: "The type of rule (detection, requirement).",
											},
											"path": schema.StringAttribute{
												Computed:            true,
												MarkdownDescription: "The path for file system rules.",
											},
											"file_or_folder_name": schema.StringAttribute{
												Computed:            true,
												MarkdownDescription: "The file or folder name for file system rules.",
											},
											"check_32bit_on_64_system": schema.BoolAttribute{
												Computed:            true,
												MarkdownDescription: "Whether to check 32-bit paths on 64-bit systems.",
											},
											"operation_type": schema.StringAttribute{
												Computed:            true,
												MarkdownDescription: "The operation type (e.g., sizeInBytes, version, string).",
											},
											"operator": schema.StringAttribute{
												Computed:            true,
												MarkdownDescription: "The operator (e.g., equal, notEqual).",
											},
											"comparison_value": schema.StringAttribute{
												Computed:            true,
												MarkdownDescription: "The value to compare against.",
											},
											"key_path": schema.StringAttribute{
												Computed:            true,
												MarkdownDescription: "The registry key path for registry rules.",
											},
											"value_name": schema.StringAttribute{
												Computed:            true,
												MarkdownDescription: "The registry value name for registry rules.",
											},
										},
									},
								},
								"install_experience": schema.SingleNestedAttribute{
									Computed:            true,
									MarkdownDescription: "The install experience configuration.",
									Attributes: map[string]schema.Attribute{
										"run_as_account": schema.StringAttribute{
											Computed:            true,
											MarkdownDescription: "The account to run the installer as (e.g., system, user).",
										},
										"max_run_time_in_minutes": schema.Int64Attribute{
											Computed:            true,
											MarkdownDescription: "The maximum allowed run time in minutes.",
										},
										"device_restart_behavior": schema.StringAttribute{
											Computed:            true,
											MarkdownDescription: "The device restart behavior (e.g., basedOnReturnCode, allow, suppress).",
										},
									},
								},
								"return_codes": schema.ListNestedAttribute{
									Computed:            true,
									MarkdownDescription: "The return codes configuration.",
									NestedObject: schema.NestedAttributeObject{
										Attributes: map[string]schema.Attribute{
											"return_code": schema.Int64Attribute{
												Computed:            true,
												MarkdownDescription: "The return code value.",
											},
											"type": schema.StringAttribute{
												Computed:            true,
												MarkdownDescription: "The type of return code (e.g., success, softReboot, hardReboot, retry).",
											},
										},
									},
								},
							},
						},
					},
				},
			},
			"timeouts": commonschema.Timeouts(ctx),
		},
	}
}
