package utilityWindowsMSIAppMetadata

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
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
	DataSourceName = "microsoft365_utility_windows_msi_app_metadata"
	ReadTimeout    = 300 // Extended timeout for file processing
)

var (
	// Basic datasource interface (CRUD operations)
	_ datasource.DataSource = &WindowsMSIAppMetadataDataSource{}

	// Allows the datasource to be configured with the provider client
	_ datasource.DataSourceWithConfigure = &WindowsMSIAppMetadataDataSource{}
)

func NewWindowsMSIAppMetadataDataSource() datasource.DataSource {
	return &WindowsMSIAppMetadataDataSource{}
}

type WindowsMSIAppMetadataDataSource struct {
	client *msgraphbetasdk.GraphServiceClient
}

// Metadata returns the datasource type name.
func (r *WindowsMSIAppMetadataDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = DataSourceName
}

// For utility datasources that perform local computations. Required for interface compliance.
func (d *WindowsMSIAppMetadataDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.client = client.SetGraphBetaClientForDataSource(ctx, req, resp, DataSourceName)
}

// Schema defines the schema for the data source
func (d *WindowsMSIAppMetadataDataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Extracts comprehensive metadata from Windows MSI installer files. Supports both local file paths and remote URLs.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Unique identifier for the data source instance.",
			},
			"installer_file_path_source": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Path to a local MSI file. Either this or installer_url_source must be provided.",
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.MatchRoot("installer_url_source")),
				},
			},
			"installer_url_source": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "URL to an MSI file. Either this or installer_file_path_source must be provided.",
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.MatchRoot("installer_file_path_source")),
				},
			},
			"metadata": schema.SingleNestedAttribute{
				Computed:            true,
				MarkdownDescription: "Extracted metadata from the MSI file.",
				Attributes: map[string]schema.Attribute{
					"product_code": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "MSI Product Code (GUID) that uniquely identifies the product.",
					},
					"product_version": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "Version number of the product.",
					},
					"product_name": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "Display name of the product.",
					},
					"publisher": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "Manufacturer/publisher of the product.",
					},
					"upgrade_code": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "MSI Upgrade Code (GUID) used for product upgrades.",
					},
					"language": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "Product language code.",
					},
					"package_type": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "Type of MSI package (e.g., Application, Patch).",
					},
					"install_location": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "Default installation directory.",
					},
					"install_command": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "Command to install the MSI package.",
					},
					"uninstall_command": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "Command to uninstall the MSI package.",
					},
					"transform_paths": schema.ListAttribute{
						Computed:            true,
						ElementType:         types.StringType,
						MarkdownDescription: "List of MST transform file paths.",
					},
					"size_mb": schema.Float64Attribute{
						Computed:            true,
						MarkdownDescription: "Size of the MSI file in megabytes.",
					},
					"sha256_checksum": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "SHA256 hash of the MSI file.",
					},
					"md5_checksum": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "MD5 hash of the MSI file.",
					},
					"properties": schema.MapAttribute{
						Computed:            true,
						ElementType:         types.StringType,
						MarkdownDescription: "All MSI properties as key-value pairs.",
					},
					"required_features": schema.ListAttribute{
						Computed:            true,
						ElementType:         types.StringType,
						MarkdownDescription: "List of required features for installation.",
					},
					"files": schema.ListAttribute{
						Computed:            true,
						ElementType:         types.StringType,
						MarkdownDescription: "List of files included in the MSI package.",
					},
					"min_os_version": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "Minimum operating system version required.",
					},
					"architecture": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "Target architecture (x86, x64, etc.).",
					},
				},
			},
			"timeouts": commonschema.DatasourceTimeouts(ctx),
		},
	}
}
