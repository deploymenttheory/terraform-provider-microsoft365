package utilityMacOSPKGAppMetadata

import (
	"context"
	"regexp"

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
	DataSourceName = "microsoft365_utility_macos_pkg_app_metadata"
	ReadTimeout    = 240 // Extended timeout for metadata extraction
)

var (
	// Basic datasource interface (CRUD operations)
	_ datasource.DataSource = &MacOSPKGAppMetadataDataSource{}

	// Allows the datasource to be configured with the provider client
	_ datasource.DataSourceWithConfigure = &MacOSPKGAppMetadataDataSource{}
)

func NewMacOSPKGAppMetadataDataSource() datasource.DataSource {
	return &MacOSPKGAppMetadataDataSource{}
}

type MacOSPKGAppMetadataDataSource struct {
	client *msgraphbetasdk.GraphServiceClient
}

// Metadata returns the datasource type name.
func (r *MacOSPKGAppMetadataDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = DataSourceName
}

// Configure configures the data source with the provider client
func (d *MacOSPKGAppMetadataDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Even though we don't need the Graph client for local file operations,
	// we'll set it up in case we need it for any future URL downloads through the Microsoft API
	d.client = client.SetGraphBetaClientForDataSource(ctx, req, resp, DataSourceName)
}

// Schema defines the schema for the data source
func (d *MacOSPKGAppMetadataDataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Extracts metadata from macOS PKG installer files locally or from URLs. This data source is used to retrieve bundle identifiers, versions, and package details for macOS app deployment.",
		Attributes: map[string]schema.Attribute{
			"installer_file_path_source": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The path to the PKG file to be uploaded. The file must be a valid `.pkg` file. Value is not returned by API call.",
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(`.*\.pkg$`),
						"File path must point to a valid .pkg file.",
					),
				},
			},
			"installer_url_source": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The web location of the PKG file, can be a http(s) URL. The file must be a valid `.pkg` file. Value is not returned by API call.",
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^(http|https|file)://.*$|^(/|./|../).*$`),
						"Must be a valid URL.",
					),
				},
			},
			"metadata": schema.SingleNestedAttribute{
				Computed:            true,
				MarkdownDescription: "Extracted metadata from the macOS PKG file.",
				Attributes: map[string]schema.Attribute{
					"cf_bundle_identifier": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "The bundle identifier (CFBundleIdentifier) extracted from the PKG file.",
					},
					"cf_bundle_short_version_string": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "The version (CFBundleShortVersionString) extracted from the PKG file.",
					},
					"name": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "The app name extracted from the PKG file.",
					},
					"package_ids": schema.ListAttribute{
						Computed:            true,
						ElementType:         types.StringType,
						MarkdownDescription: "The list of package IDs found in the PKG file.",
					},
					"install_location": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "The installation location specified in the PKG file.",
					},
					"app_paths": schema.ListAttribute{
						Computed:            true,
						ElementType:         types.StringType,
						MarkdownDescription: "List of application paths found in the PKG file.",
					},
					"min_os_version": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "The minimum OS version required by the PKG app.",
					},
					"size_mb": schema.Float64Attribute{
						Computed:            true,
						MarkdownDescription: "The size of the PKG file in megabytes (MB).",
					},
					"sha256_checksum": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "SHA256 hash of the PKG file content as a hexadecimal string.",
					},
					"md5_checksum": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "MD5 hash of the PKG file content as a hexadecimal string.",
					},
					"included_bundles": schema.ListNestedAttribute{
						Computed:            true,
						MarkdownDescription: "List of additional bundles included in the PKG file.",
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"bundle_id": schema.StringAttribute{
									Computed:            true,
									MarkdownDescription: "The bundle identifier of the included bundle.",
								},
								"version": schema.StringAttribute{
									Computed:            true,
									MarkdownDescription: "The version of the included bundle.",
								},
								"path": schema.StringAttribute{
									Computed:            true,
									MarkdownDescription: "The path of the included bundle within the PKG file.",
								},
								"cf_bundle_version": schema.StringAttribute{
									Computed:            true,
									MarkdownDescription: "The CFBundleVersion of the included bundle (may differ from version string).",
								},
							},
						},
					},
				},
			},
			"timeouts": commonschema.DatasourceTimeouts(ctx),
		},
	}
}
