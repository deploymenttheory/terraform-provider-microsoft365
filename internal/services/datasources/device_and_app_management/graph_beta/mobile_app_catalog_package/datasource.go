package graphBetaMobileAppCatalogPackage

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
	DataSourceName = "microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package"
	ReadTimeout    = 180
)

var (
	// Basic datasource interface (CRUD operations)
	_ datasource.DataSource = &MobileAppCatalogPackageDataSource{}

	// Allows the datasource to be configured with the provider client
	_ datasource.DataSourceWithConfigure = &MobileAppCatalogPackageDataSource{}
)

func NewMobileAppCatalogPackageDataSource() datasource.DataSource {
	return &MobileAppCatalogPackageDataSource{
		ReadPermissions: []string{
			"DeviceManagementApps.Read.All",
		},
	}
}

type MobileAppCatalogPackageDataSource struct {
	client          *msgraphbetasdk.GraphServiceClient
	ReadPermissions []string
}

// Metadata returns the datasource type name.
func (r *MobileAppCatalogPackageDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = DataSourceName
}

// Configure sets the client for the data source
func (d *MobileAppCatalogPackageDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.client = client.SetGraphBetaClientForDataSource(ctx, req, resp, DataSourceName)
}

// Schema defines the schema for the data source
func (d *MobileAppCatalogPackageDataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Retrieves mobile app catalog packages from Microsoft Intune using the `/deviceAppManagement/MobileAppCatalogPackage` endpoint. This data source enables querying mobile app catalog packages with advanced filtering capabilities including OData queries for filtering by product name, publisher, and other properties.",
		Attributes: map[string]schema.Attribute{
			"filter_type": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Type of filter to apply. Valid values are: `all`, `id`, `product_name`, `publisher_name`, `odata`.",
				Validators: []validator.String{
					stringvalidator.OneOf("all", "id", "product_name", "publisher_name", "odata"),
				},
			},
			"filter_value": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Value to filter by. Not required when filter_type is 'all' or 'odata'.",
			},
			"odata_filter": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "OData $filter parameter for filtering results. Only used when filter_type is 'odata'. Example: productDisplayName eq 'Microsoft Office'.",
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
				MarkdownDescription: "OData $orderby parameter to sort results. Only used when filter_type is 'odata'. Example: productDisplayName.",
			},
			"odata_count": schema.BoolAttribute{
				Optional:            true,
				MarkdownDescription: "OData $count parameter to include count of total results. Only used when filter_type is 'odata'.",
			},
			"odata_search": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "OData $search parameter for full-text search. Only used when filter_type is 'odata'.",
			},
			"odata_expand": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "OData $expand parameter to include related entities. Only used when filter_type is 'odata'.",
			},
			"items": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "The list of win32 catalog applications with full details that match the filter criteria.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The unique identifier for the application.",
						},
						"display_name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The display name of the application.",
						},
						"description": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The description of the application.",
						},
						"publisher": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The publisher of the application.",
						},
						"created_date_time": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The date and time the app was created.",
						},
						"last_modified_date_time": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The date and time the app was last modified.",
						},
						"is_featured": schema.BoolAttribute{
							Computed:            true,
							MarkdownDescription: "Indicates whether the app is marked as featured by the admin.",
						},
						"privacy_information_url": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The privacy statement URL.",
						},
						"information_url": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The more information URL.",
						},
						"owner": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The owner of the app.",
						},
						"developer": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The developer of the app.",
						},
						"notes": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "Notes for the app.",
						},
						"upload_state": schema.Int32Attribute{
							Computed:            true,
							MarkdownDescription: "The upload state.",
						},
						"publishing_state": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The publishing state for the app.",
						},
						"is_assigned": schema.BoolAttribute{
							Computed:            true,
							MarkdownDescription: "Indicates whether the app is assigned to at least one group.",
						},
						"role_scope_tag_ids": schema.ListAttribute{
							Computed:            true,
							ElementType:         types.StringType,
							MarkdownDescription: "List of scope tag IDs for this app.",
						},
						"dependent_app_count": schema.Int32Attribute{
							Computed:            true,
							MarkdownDescription: "The total number of dependencies the child app has.",
						},
						"superseding_app_count": schema.Int32Attribute{
							Computed:            true,
							MarkdownDescription: "The total number of apps this app directly or indirectly supersedes.",
						},
						"superseded_app_count": schema.Int32Attribute{
							Computed:            true,
							MarkdownDescription: "The total number of apps this app is directly or indirectly superseded by.",
						},
						"committed_content_version": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The internal committed content version.",
						},
						"file_name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The name of the main installation file.",
						},
						"size": schema.Int64Attribute{
							Computed:            true,
							MarkdownDescription: "The total size of the application in bytes.",
						},
						"install_command_line": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The command line to install this app.",
						},
						"uninstall_command_line": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The command line to uninstall this app.",
						},
						"applicable_architectures": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The Windows architecture(s) on which this app can run.",
						},
						"allowed_architectures": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The allowed target architectures for this app.",
						},
						"minimum_free_disk_space_in_mb": schema.Int32Attribute{
							Computed:            true,
							MarkdownDescription: "The minimum free disk space required to install this app.",
						},
						"minimum_memory_in_mb": schema.Int32Attribute{
							Computed:            true,
							MarkdownDescription: "The minimum memory required to install this app.",
						},
						"minimum_number_of_processors": schema.Int32Attribute{
							Computed:            true,
							MarkdownDescription: "The minimum number of processors required to install this app.",
						},
						"minimum_cpu_speed_in_mhz": schema.Int32Attribute{
							Computed:            true,
							MarkdownDescription: "The minimum CPU speed required to install this app.",
						},
						"setup_file_path": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The relative path of the setup file in the app package.",
						},
						"minimum_supported_windows_release": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The minimum supported Windows release version.",
						},
						"display_version": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The display version of the application.",
						},
						"allow_available_uninstall": schema.BoolAttribute{
							Computed:            true,
							MarkdownDescription: "Indicates whether the app can be uninstalled from the available context.",
						},
						"mobile_app_catalog_package_id": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The mobile app catalog package ID.",
						},
						"rules": schema.ListNestedAttribute{
							Computed:            true,
							MarkdownDescription: "The detection and requirement rules for this app.",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"odata_type": schema.StringAttribute{
										Computed:            true,
										MarkdownDescription: "The OData type of the rule.",
									},
									"rule_type": schema.StringAttribute{
										Computed:            true,
										MarkdownDescription: "The type of rule (detection or requirement).",
									},
									"path": schema.StringAttribute{
										Computed:            true,
										MarkdownDescription: "The file or folder path for file system rules.",
									},
									"file_or_folder_name": schema.StringAttribute{
										Computed:            true,
										MarkdownDescription: "The file or folder name to detect.",
									},
									"check_32bit_on_64system": schema.BoolAttribute{
										Computed:            true,
										MarkdownDescription: "Indicates whether to check 32-bit on a 64-bit system.",
									},
									"operation_type": schema.StringAttribute{
										Computed:            true,
										MarkdownDescription: "The operation type for the rule.",
									},
									"operator": schema.StringAttribute{
										Computed:            true,
										MarkdownDescription: "The operator for the rule.",
									},
									"comparison_value": schema.StringAttribute{
										Computed:            true,
										MarkdownDescription: "The comparison value for the rule.",
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
							MarkdownDescription: "The install experience for this app.",
							Attributes: map[string]schema.Attribute{
								"run_as_account": schema.StringAttribute{
									Computed:            true,
									MarkdownDescription: "Indicates the account context to execute the app.",
								},
								"max_run_time_in_minutes": schema.Int32Attribute{
									Computed:            true,
									MarkdownDescription: "The maximum run time in minutes.",
								},
								"device_restart_behavior": schema.StringAttribute{
									Computed:            true,
									MarkdownDescription: "Device restart behavior.",
								},
							},
						},
						"return_codes": schema.ListNestedAttribute{
							Computed:            true,
							MarkdownDescription: "The return codes for post installation behavior.",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"return_code": schema.Int32Attribute{
										Computed:            true,
										MarkdownDescription: "The return code.",
									},
									"type": schema.StringAttribute{
										Computed:            true,
										MarkdownDescription: "The type of return code.",
									},
								},
							},
						},
						"msi_information": schema.SingleNestedAttribute{
							Computed:            true,
							MarkdownDescription: "The MSI information for MSI-based apps.",
							Attributes: map[string]schema.Attribute{
								"product_code": schema.StringAttribute{
									Computed:            true,
									MarkdownDescription: "The MSI product code.",
								},
								"product_version": schema.StringAttribute{
									Computed:            true,
									MarkdownDescription: "The MSI product version.",
								},
								"upgrade_code": schema.StringAttribute{
									Computed:            true,
									MarkdownDescription: "The MSI upgrade code.",
								},
								"requires_reboot": schema.BoolAttribute{
									Computed:            true,
									MarkdownDescription: "Whether the MSI app requires reboot.",
								},
								"package_type": schema.StringAttribute{
									Computed:            true,
									MarkdownDescription: "The MSI package type.",
								},
								"product_name": schema.StringAttribute{
									Computed:            true,
									MarkdownDescription: "The MSI product name.",
								},
								"publisher": schema.StringAttribute{
									Computed:            true,
									MarkdownDescription: "The MSI publisher.",
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
