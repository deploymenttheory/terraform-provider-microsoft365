package graphBetaMobileApp

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
	datasourceName = "graph_beta_device_and_app_management_mobile_app"
	ReadTimeout    = 180
)

var (
	// Basic datasource interface (CRUD operations)
	_ datasource.DataSource = &MobileAppDataSource{}

	// Allows the datasource to be configured with the provider client
	_ datasource.DataSourceWithConfigure = &MobileAppDataSource{}
)

func NewMobileAppDataSource() datasource.DataSource {
	return &MobileAppDataSource{
		ReadPermissions: []string{
			"DeviceManagementApps.Read.All",
		},
	}
}

type MobileAppDataSource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ProviderTypeName string
	TypeName         string
	ReadPermissions  []string
}

// Metadata returns the datasource type name.
func (r *MobileAppDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + datasourceName
}

// Configure sets the client for the data source
func (d *MobileAppDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.client = client.SetGraphBetaClientForDataSource(ctx, req, resp, d.TypeName)
}

// Schema defines the schema for the data source
func (d *MobileAppDataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Retrieves mobile applications from Microsoft Intune using the `/deviceAppManagement/mobileApps` endpoint. This data source enables querying all mobile app types including Win32 LOB apps, PKG/DMG apps, store apps, and web apps with advanced filtering capabilities for application discovery and configuration planning.",
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
			"app_type_filter": schema.StringAttribute{
				Optional: true,
				MarkdownDescription: "Optional filter that filters returned apps by the application type. Supported values are: " +
					"`macOSPkgApp`, `macOSDmgApp`, `macOSLobApp`, `macOSMicrosoftDefenderApp`, `macOSMicrosoftEdgeApp`, " +
					"`macOSOfficeSuiteApp`, `macOsVppApp`, `macOSWebClip`, `androidForWorkApp`, `androidLobApp`, " +
					"`androidManagedStoreApp`, `androidManagedStoreWebApp`, `androidStoreApp`, `managedAndroidLobApp`, " +
					"`managedAndroidStoreApp`, `iosiPadOSWebClip`, `iosLobApp`, `iosStoreApp`, `iosVppApp`, `managedIOSLobApp`, " +
					"`managedIOSStoreApp`, `windowsAppX`, `windowsMicrosoftEdgeApp`, `windowsMobileMSI`, `windowsPhone81AppX`, " +
					"`windowsPhone81AppXBundle`, `windowsPhone81StoreApp`, `windowsPhoneXAP`, `windowsStoreApp`, " +
					"`windowsUniversalAppX`, `windowsWebApp`, `winGetApp`, `webApp`, `microsoftStoreForBusinessApp`, " +
					"`officeSuiteApp`, `win32CatalogApp`, `win32LobApp`, `managedApp`, `managedMobileLobApp`, `mobileLobApp`.",
				Validators: []validator.String{
					stringvalidator.OneOf(
						// MacOS Apps
						"macOSPkgApp", "macOSDmgApp", "macOSLobApp", "macOSMicrosoftDefenderApp",
						"macOSMicrosoftEdgeApp", "macOSOfficeSuiteApp", "macOsVppApp", "macOSWebClip",
						// Android Apps
						"androidForWorkApp", "androidLobApp", "androidManagedStoreApp", "androidManagedStoreWebApp",
						"androidStoreApp", "managedAndroidLobApp", "managedAndroidStoreApp",
						// iOS Apps
						"iosiPadOSWebClip", "iosLobApp", "iosStoreApp", "iosVppApp",
						"managedIOSLobApp", "managedIOSStoreApp",
						// Windows Apps
						"windowsAppX", "windowsMicrosoftEdgeApp", "windowsMobileMSI", "windowsPhone81AppX",
						"windowsPhone81AppXBundle", "windowsPhone81StoreApp", "windowsPhoneXAP",
						"windowsStoreApp", "windowsUniversalAppX", "windowsWebApp", "winGetApp",
						// Web Apps
						"webApp",
						// Microsoft Store Apps
						"microsoftStoreForBusinessApp",
						// Office Apps
						"officeSuiteApp",
						// Win32 Apps
						"win32CatalogApp", "win32LobApp",
						// Other App Types
						"managedApp", "managedMobileLobApp", "mobileLobApp",
					),
				},
			},
			"items": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "The list of mobile apps that match the filter criteria.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "Key of the entity.",
						},
						"display_name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The admin provided or imported title of the app.",
						},
						"description": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The description of the app.",
						},
						"publisher": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The publisher of the app.",
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
							MarkdownDescription: "The value indicating whether the app is marked as featured by the admin.",
						},
						"privacy_information_url": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The privacy statement Url.",
						},
						"information_url": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The more information Url.",
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
							MarkdownDescription: "The publishing state for the app. The app cannot be assigned unless the app is published. Possible values are: `notPublished`, `processing`, `published`.",
						},
						"is_assigned": schema.BoolAttribute{
							Computed:            true,
							MarkdownDescription: "The value indicating whether the app is assigned to at least one group.",
						},
						"role_scope_tag_ids": schema.ListAttribute{
							Computed:            true,
							ElementType:         types.StringType,
							MarkdownDescription: "List of scope tag ids for this mobile app.",
						},
						"dependent_app_count": schema.Int32Attribute{
							Computed:            true,
							MarkdownDescription: "The total number of dependencies the child app has.",
						},
						"superseded_app_count": schema.Int32Attribute{
							Computed:            true,
							MarkdownDescription: "The total number of apps this app is directly or indirectly superseded by.",
						},
						"superseding_app_count": schema.Int32Attribute{
							Computed:            true,
							MarkdownDescription: "The total number of apps this app directly or indirectly supersedes.",
						},
						"categories": schema.ListAttribute{
							Computed:            true,
							ElementType:         types.StringType,
							MarkdownDescription: "The list of categories for this app.",
						},
					},
				},
			},
			"timeouts": commonschema.Timeouts(ctx),
		},
	}
}
