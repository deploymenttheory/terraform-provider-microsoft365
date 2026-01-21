package utilityMicrosoftStorePackageManifest

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
	DataSourceName = "microsoft365_utility_microsoft_store_package_manifest_metadata"
	ReadTimeout    = 240 // Extended timeout for API calls
)

var (
	// Basic datasource interface (CRUD operations)
	_ datasource.DataSource = &MicrosoftStorePackageManifestDataSource{}

	// Allows the datasource to be configured with the provider client
	_ datasource.DataSourceWithConfigure = &MicrosoftStorePackageManifestDataSource{}
)

func NewMicrosoftStorePackageManifestDataSource() datasource.DataSource {
	return &MicrosoftStorePackageManifestDataSource{}
}

type MicrosoftStorePackageManifestDataSource struct {
	client *msgraphbetasdk.GraphServiceClient
}

// Metadata returns the datasource type name.
func (r *MicrosoftStorePackageManifestDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = DataSourceName
}

// Configure configures the data source with the provider client
func (d *MicrosoftStorePackageManifestDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.client = client.SetGraphBetaClientForDataSource(ctx, req, resp, DataSourceName)
}

// Schema defines the schema for the data source
func (d *MicrosoftStorePackageManifestDataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Queries winget package manifests from the Microsoft Store API using package identifiers or search terms. This data source is used to discover installation metadata for Windows Package Manager applications.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The ID of this resource.",
			},
			"package_identifier": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The specific package identifier to retrieve manifest for. Either this or search_term must be provided.",
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.MatchRoot("search_term")),
				},
			},
			"search_term": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Search term to find packages. Either this or package_identifier must be provided.",
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.MatchRoot("package_identifier")),
				},
			},
			"manifests": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "List of package manifests retrieved.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"type": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The type identifier for the manifest.",
						},
						"package_identifier": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The package identifier.",
						},
						"versions": schema.ListNestedAttribute{
							Computed:            true,
							MarkdownDescription: "List of package versions.",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"type": schema.StringAttribute{
										Computed:            true,
										MarkdownDescription: "The type identifier for the version.",
									},
									"package_version": schema.StringAttribute{
										Computed:            true,
										MarkdownDescription: "The package version number.",
									},
									"default_locale": schema.SingleNestedAttribute{
										Computed:            true,
										MarkdownDescription: "Default locale information for the package.",
										Attributes: map[string]schema.Attribute{
											"type": schema.StringAttribute{
												Computed:            true,
												MarkdownDescription: "The type identifier for the locale.",
											},
											"package_locale": schema.StringAttribute{
												Computed:            true,
												MarkdownDescription: "The locale code (e.g., en-us).",
											},
											"publisher": schema.StringAttribute{
												Computed:            true,
												MarkdownDescription: "The publisher name.",
											},
											"publisher_url": schema.StringAttribute{
												Computed:            true,
												MarkdownDescription: "The publisher website URL.",
											},
											"privacy_url": schema.StringAttribute{
												Computed:            true,
												MarkdownDescription: "The privacy policy URL.",
											},
											"publisher_support_url": schema.StringAttribute{
												Computed:            true,
												MarkdownDescription: "The publisher support URL.",
											},
											"package_name": schema.StringAttribute{
												Computed:            true,
												MarkdownDescription: "The package name.",
											},
											"license": schema.StringAttribute{
												Computed:            true,
												MarkdownDescription: "The license information.",
											},
											"copyright": schema.StringAttribute{
												Computed:            true,
												MarkdownDescription: "The copyright information.",
											},
											"short_description": schema.StringAttribute{
												Computed:            true,
												MarkdownDescription: "Short description of the package.",
											},
											"description": schema.StringAttribute{
												Computed:            true,
												MarkdownDescription: "Detailed description of the package.",
											},
											"tags": schema.ListAttribute{
												Computed:            true,
												MarkdownDescription: "List of tags associated with the package.",
												ElementType:         types.StringType,
											},
											"agreements": schema.ListNestedAttribute{
												Computed:            true,
												MarkdownDescription: "List of agreements for the package.",
												NestedObject: schema.NestedAttributeObject{
													Attributes: map[string]schema.Attribute{
														"type": schema.StringAttribute{
															Computed:            true,
															MarkdownDescription: "The type identifier for the agreement.",
														},
														"agreement_label": schema.StringAttribute{
															Computed:            true,
															MarkdownDescription: "The agreement label.",
														},
														"agreement": schema.StringAttribute{
															Computed:            true,
															MarkdownDescription: "The agreement text.",
														},
														"agreement_url": schema.StringAttribute{
															Computed:            true,
															MarkdownDescription: "The agreement URL.",
														},
													},
												},
											},
										},
									},
									"locales": schema.ListNestedAttribute{
										Computed:            true,
										MarkdownDescription: "List of locale-specific information.",
										NestedObject: schema.NestedAttributeObject{
											Attributes: map[string]schema.Attribute{
												"type": schema.StringAttribute{
													Computed:            true,
													MarkdownDescription: "The type identifier for the locale.",
												},
												"package_locale": schema.StringAttribute{
													Computed:            true,
													MarkdownDescription: "The locale code (e.g., en-us).",
												},
												"publisher": schema.StringAttribute{
													Computed:            true,
													MarkdownDescription: "The publisher name.",
												},
												"publisher_url": schema.StringAttribute{
													Computed:            true,
													MarkdownDescription: "The publisher website URL.",
												},
												"privacy_url": schema.StringAttribute{
													Computed:            true,
													MarkdownDescription: "The privacy policy URL.",
												},
												"publisher_support_url": schema.StringAttribute{
													Computed:            true,
													MarkdownDescription: "The publisher support URL.",
												},
												"package_name": schema.StringAttribute{
													Computed:            true,
													MarkdownDescription: "The package name.",
												},
												"license": schema.StringAttribute{
													Computed:            true,
													MarkdownDescription: "The license information.",
												},
												"copyright": schema.StringAttribute{
													Computed:            true,
													MarkdownDescription: "The copyright information.",
												},
												"short_description": schema.StringAttribute{
													Computed:            true,
													MarkdownDescription: "Short description of the package.",
												},
												"description": schema.StringAttribute{
													Computed:            true,
													MarkdownDescription: "Detailed description of the package.",
												},
												"tags": schema.ListAttribute{
													Computed:            true,
													MarkdownDescription: "List of tags associated with the package.",
													ElementType:         types.StringType,
												},
											},
										},
									},
									"installers": schema.ListNestedAttribute{
										Computed:            true,
										MarkdownDescription: "List of installer information.",
										NestedObject: schema.NestedAttributeObject{
											Attributes: map[string]schema.Attribute{
												"type": schema.StringAttribute{
													Computed:            true,
													MarkdownDescription: "The type identifier for the installer.",
												},
												"ms_store_product_identifier": schema.StringAttribute{
													Computed:            true,
													MarkdownDescription: "Microsoft Store product identifier.",
												},
												"architecture": schema.StringAttribute{
													Computed:            true,
													MarkdownDescription: "The target architecture (e.g., x86, x64, arm64).",
												},
												"installer_type": schema.StringAttribute{
													Computed:            true,
													MarkdownDescription: "The installer type (e.g., msstore, exe).",
												},
												"package_family_name": schema.StringAttribute{
													Computed:            true,
													MarkdownDescription: "The package family name.",
												},
												"scope": schema.StringAttribute{
													Computed:            true,
													MarkdownDescription: "The installation scope (user or machine).",
												},
												"download_command_prohibited": schema.BoolAttribute{
													Computed:            true,
													MarkdownDescription: "Whether download command is prohibited.",
												},
												"installer_sha256": schema.StringAttribute{
													Computed:            true,
													MarkdownDescription: "SHA256 hash of the installer.",
												},
												"installer_url": schema.StringAttribute{
													Computed:            true,
													MarkdownDescription: "URL to download the installer.",
												},
												"installer_locale": schema.StringAttribute{
													Computed:            true,
													MarkdownDescription: "Locale for the installer.",
												},
												"minimum_os_version": schema.StringAttribute{
													Computed:            true,
													MarkdownDescription: "Minimum OS version required.",
												},
												"installer_success_codes": schema.ListAttribute{
													Computed:            true,
													MarkdownDescription: "List of success codes for the installer.",
													ElementType:         types.Int64Type,
												},
												"markets": schema.SingleNestedAttribute{
													Computed:            true,
													MarkdownDescription: "Market information for the installer.",
													Attributes: map[string]schema.Attribute{
														"type": schema.StringAttribute{
															Computed:            true,
															MarkdownDescription: "The type identifier for markets.",
														},
														"allowed_markets": schema.ListAttribute{
															Computed:            true,
															MarkdownDescription: "List of allowed markets.",
															ElementType:         types.StringType,
														},
													},
												},
												"installer_switches": schema.SingleNestedAttribute{
													Computed:            true,
													MarkdownDescription: "Installer switches information.",
													Attributes: map[string]schema.Attribute{
														"type": schema.StringAttribute{
															Computed:            true,
															MarkdownDescription: "The type identifier for installer switches.",
														},
														"silent": schema.StringAttribute{
															Computed:            true,
															MarkdownDescription: "Silent installation switch.",
														},
													},
												},
												"expected_return_codes": schema.ListNestedAttribute{
													Computed:            true,
													MarkdownDescription: "List of expected return codes.",
													NestedObject: schema.NestedAttributeObject{
														Attributes: map[string]schema.Attribute{
															"type": schema.StringAttribute{
																Computed:            true,
																MarkdownDescription: "The type identifier for return code.",
															},
															"installer_return_code": schema.Int64Attribute{
																Computed:            true,
																MarkdownDescription: "The return code.",
															},
															"return_response": schema.StringAttribute{
																Computed:            true,
																MarkdownDescription: "The response description for the return code.",
															},
														},
													},
												},
												"apps_and_features_entries": schema.ListNestedAttribute{
													Computed:            true,
													MarkdownDescription: "List of Apps and Features entries.",
													NestedObject: schema.NestedAttributeObject{
														Attributes: map[string]schema.Attribute{
															"type": schema.StringAttribute{
																Computed:            true,
																MarkdownDescription: "The type identifier for the entry.",
															},
															"display_name": schema.StringAttribute{
																Computed:            true,
																MarkdownDescription: "Display name in Apps and Features.",
															},
															"publisher": schema.StringAttribute{
																Computed:            true,
																MarkdownDescription: "Publisher name in Apps and Features.",
															},
															"display_version": schema.StringAttribute{
																Computed:            true,
																MarkdownDescription: "Display version in Apps and Features.",
															},
															"product_code": schema.StringAttribute{
																Computed:            true,
																MarkdownDescription: "Product code.",
															},
															"installer_type": schema.StringAttribute{
																Computed:            true,
																MarkdownDescription: "Installer type.",
															},
														},
													},
												},
											},
										},
									},
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
