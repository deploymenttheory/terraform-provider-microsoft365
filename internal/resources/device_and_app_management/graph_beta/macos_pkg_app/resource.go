package graphBetaMacOSPKGApp

import (
	"context"
	"regexp"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common"
	planmodifiers "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/plan_modifiers"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/schema"
	commonschemagraphbeta "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/schema/graph_beta/device_and_app_management"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/validators"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	ResourceName  = "graph_beta_device_and_app_management_macos_pkg_app"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	// Basic resource interface (CRUD operations)
	_ resource.Resource = &MacOSPKGAppResource{}

	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &MacOSPKGAppResource{}

	// Enables import functionality
	_ resource.ResourceWithImportState = &MacOSPKGAppResource{}

	// Enables plan modification/diff suppression
	_ resource.ResourceWithModifyPlan = &MacOSPKGAppResource{}
)

func NewMacOSPKGAppResource() resource.Resource {
	return &MacOSPKGAppResource{
		ReadPermissions: []string{
			"DeviceManagementApps.Read.All",
			"DeviceManagementConfiguration.Read.All",
		},
		WritePermissions: []string{
			"DeviceManagementApps.ReadWrite.All",
			"DeviceManagementConfiguration.ReadWrite.All",
		},
		ResourcePath: "/deviceAppManagement/mobileApps",
	}
}

type MacOSPKGAppResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ProviderTypeName string
	TypeName         string
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

// Metadata returns the resource type name.
func (r *MacOSPKGAppResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + ResourceName
}

// Configure sets the client for the resource.
func (r *MacOSPKGAppResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = common.SetGraphBetaClientForResource(ctx, req, resp, r.TypeName)
}

// ImportState imports the resource state.
func (r *MacOSPKGAppResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Schema returns the schema for the resource.
func (r *MacOSPKGAppResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages an Intune macOS app (PKG), using the mobileapps graph beta API. Apps are deployed using the Microsoft Intune management agent for macOS.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
				MarkdownDescription: "The unique identifier of the macOS PKG.",
			},
			// "application_type": schema.StringAttribute{
			// 	Required: true,
			// 	MarkdownDescription: "The type of Intune application to deploy. Possible values are:" +
			// 		"`AndroidForWorkApp`, `AndroidLobApp`, `AndroidManagedStoreApp`, `AndroidManagedStoreWebApp`," +
			// 		"`AndroidStoreApp`, `IosiPadOSWebClip`, `IosLobApp`, `IosStoreApp`, `IosVppApp`, `MacOSDmgApp`," +
			// 		"`MacOSLobApp`, `MacOSMicrosoftDefenderApp`, `MacOSMicrosoftEdgeApp`, `MacOSOfficeSuiteApp`," +
			// 		"`MacOSPkgApp`, `MacOsVppApp`, `MacOSWebClip`, `ManagedAndroidLobApp`, `ManagedAndroidStoreApp`," +
			// 		"`ManagedApp`, `ManagedIOSLobApp`, `ManagedIOSStoreApp`, `ManagedMobileLobApp`, `MicrosoftStoreForBusinessApp`," +
			// 		"`MobileLobApp`, `OfficeSuiteApp`, `WebApp`, `Win32CatalogApp`, `Win32LobApp`, `WindowsAppX`," +
			// 		"`WindowsMicrosoftEdgeApp`, `WindowsMobileMSI`, `WindowsPhone81AppX`, `WindowsPhone81AppXBundle`," +
			// 		"`WindowsPhone81StoreApp`, `WindowsPhoneXAP`, `WindowsStoreApp`, `WindowsUniversalAppX`, `WindowsWebApp`, `MacOSPKGAppResource`",
			// 	Validators: []validator.String{
			// 		stringvalidator.OneOf(
			// 			"AndroidForWorkApp",
			// 			"AndroidLobApp",
			// 			"AndroidManagedStoreApp",
			// 			"AndroidManagedStoreWebApp",
			// 			"AndroidStoreApp",
			// 			"IosiPadOSWebClip",
			// 			"IosLobApp",
			// 			"IosStoreApp",
			// 			"IosVppApp",
			// 			"MacOSDmgApp",
			// 			"MacOSLobApp",
			// 			"MacOSMicrosoftDefenderApp",
			// 			"MacOSMicrosoftEdgeApp",
			// 			"MacOSOfficeSuiteApp",
			// 			"MacOSPkgApp",
			// 			"MacOsVppApp",
			// 			"MacOSWebClip",
			// 			"ManagedAndroidLobApp",
			// 			"ManagedAndroidStoreApp",
			// 			"ManagedApp",
			// 			"ManagedIOSLobApp",
			// 			"ManagedIOSStoreApp",
			// 			"ManagedMobileLobApp",
			// 			"MicrosoftStoreForBusinessApp",
			// 			"MobileLobApp",
			// 			"OfficeSuiteApp",
			// 			"WebApp",
			// 			"Win32CatalogApp",
			// 			"Win32LobApp",
			// 			"WindowsAppX",
			// 			"WindowsMicrosoftEdgeApp",
			// 			"WindowsMobileMSI",
			// 			"WindowsPhone81AppX",
			// 			"WindowsPhone81AppXBundle",
			// 			"WindowsPhone81StoreApp",
			// 			"WindowsPhoneXAP",
			// 			"WindowsStoreApp",
			// 			"WindowsUniversalAppX",
			// 			"WindowsWebApp",
			// 			"MacOSPKGAppResource",
			// 		),
			// 	},
			// },
			"is_featured": schema.BoolAttribute{
				Optional:            true,
				MarkdownDescription: "The value indicating whether the app is marked as featured by the admin.",
			},
			"privacy_information_url": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The privacy statement Url.",
			},
			"information_url": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The more information Url.",
			},
			"owner": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The owner of the app.",
			},
			"developer": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The developer of the app.",
			},
			"notes": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Notes for the app.",
			},
			"display_name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The title of the Intune macOS pkg application.",
			},
			"description": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "A detailed description of the Intune macOS pkg application.",
			},
			"publisher": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The publisher of the Intune macOS pkg application.",
			},
			"app_icon": schema.SingleNestedAttribute{
				MarkdownDescription: "The path to the icon file to be uploaded. Resource supports both local file sources and url based sources.",
				Optional:            true,
				Attributes: map[string]schema.Attribute{
					"icon_file_path_source": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "The file path to the icon file (PNG) to be uploaded.",
						Validators: []validator.String{
							stringvalidator.RegexMatches(
								regexp.MustCompile(`\.png$`),
								"must end with .png file extension",
							),
						},
					},
					"icon_url_source": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "The web location of the icon file, can be a http(s) URL.",
						Validators: []validator.String{
							stringvalidator.RegexMatches(
								regexp.MustCompile(`^(http|https|file)://.*$|^(/|./|../).*$`),
								"Must be a valid URL.",
							),
							stringvalidator.RegexMatches(
								regexp.MustCompile(`\.png$`),
								"must end with .png file extension",
							),
						},
						PlanModifiers: []planmodifier.String{
							planmodifiers.UseStateForUnknownString(),
						},
					},
				},
			},
			"categories": schema.SetAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				MarkdownDescription: "Set of category names to associate with this application. You can use either thebpredefined Intune category names like 'Business', 'Productivity', etc., or provide specific category UUIDs. Predefined values include: 'Other apps', 'Books & Reference', 'Data management', 'Productivity', 'Business', 'Development & Design', 'Photos & Media', 'Collaboration & Social', 'Computer management'.",
				Validators: []validator.Set{
					setvalidator.ValueStringsAre(
						stringvalidator.RegexMatches(
							regexp.MustCompile(`^(Other apps|Books & Reference|Data management|Productivity|Business|Development & Design|Photos & Media|Collaboration & Social|Computer management|[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12})$`),
							"must be either a predefined category name or a valid GUID in the format xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
						),
					),
				},
			},
			"relationships": schema.ListNestedAttribute{
				Optional:            true,
				MarkdownDescription: "List of relationships associated with this application.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The unique identifier for the relationship.",
						},
						"source_display_name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The display name of the source application in the relationship.",
						},
						"source_display_version": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The display version of the source application in the relationship.",
						},
						"source_id": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The unique identifier of the source application in the relationship.",
						},
						"source_publisher_display_name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The display name of the source application's publisher in the relationship.",
						},
						"target_display_name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The display name of the target application in the relationship.",
						},
						"target_display_version": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The display version of the target application in the relationship.",
						},
						"target_id": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The unique identifier of the target application in the relationship.",
						},
						"target_publisher": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The publisher of the target application in the relationship.",
						},
						"target_publisher_display_name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The display name of the target application's publisher in the relationship.",
						},
						"target_type": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The type of the target in the relationship.",
						},
					},
				},
			},
			"created_date_time": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
				MarkdownDescription: "The date and time the app was created. This property is read-only.",
			},
			"last_modified_date_time": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The date and time the app was last modified. This property is read-only.",
			},
			"upload_state": schema.Int64Attribute{
				Computed: true,
				PlanModifiers: []planmodifier.Int64{
					planmodifiers.UseStateForUnknownInt64(),
				},
				MarkdownDescription: "The upload state. Possible values are: 0 - Not Ready, 1 - Ready, 2 - Processing. This property is read-only.",
			},
			"publishing_state": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
				MarkdownDescription: "The publishing state for the app. The app cannot be assigned unless the app is published. " +
					"Possible values are: notPublished, processing, published.",
			},
			"is_assigned": schema.BoolAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.Bool{
					planmodifiers.UseStateForUnknownBool(),
				},
				MarkdownDescription: "The value indicating whether the app is assigned to at least one group. This property is read-only.",
			},
			"role_scope_tag_ids": schema.SetAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				MarkdownDescription: "Set of scope tag ids for this mobile app.",
			},
			"dependent_app_count": schema.Int64Attribute{
				Computed: true,
				PlanModifiers: []planmodifier.Int64{
					planmodifiers.UseStateForUnknownInt64(),
				},
				MarkdownDescription: "The total number of dependencies the child app has. This property is read-only.",
			},
			"superseding_app_count": schema.Int64Attribute{
				Computed: true,
				PlanModifiers: []planmodifier.Int64{
					planmodifiers.UseStateForUnknownInt64(),
				},
				MarkdownDescription: "The total number of apps this app directly or indirectly supersedes. This property is read-only.",
			},
			"superseded_app_count": schema.Int64Attribute{
				Computed: true,
				PlanModifiers: []planmodifier.Int64{
					planmodifiers.UseStateForUnknownInt64(),
				},
				MarkdownDescription: "The total number of apps this app is directly or indirectly superseded by. This property is read-only.",
			},
			"macos_pkg_app": schema.SingleNestedAttribute{
				MarkdownDescription: "The path to the PKG file to be uploaded. Resource supports both local file sources and url based sources.",
				Optional:            true,
				Validators: []validator.Object{
					validators.ExactlyOneOf("installer_file_path_source", "installer_url_source"),
				},
				Attributes: map[string]schema.Attribute{
					"installer_file_path_source": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "The path to the PKG file to be uploaded. The file must be a valid `.pkg` file.",
						Validators: []validator.String{
							stringvalidator.RegexMatches(
								regexp.MustCompile(`.*\.pkg$`),
								"File path must point to a valid .pkg file.",
							),
						},
						// PlanModifiers: []planmodifier.String{
						// 	planmodifiers.UseStateForUnknownString(),
						// },
					},
					"installer_url_source": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "The web location of the PKG file, can be a http(s) URL. The file must be a valid `.pkg` file.",
						Validators: []validator.String{
							stringvalidator.RegexMatches(
								regexp.MustCompile(`^(http|https|file)://.*$|^(/|./|../).*$`),
								"Must be a valid URL.",
							),
						},
						PlanModifiers: []planmodifier.String{
							planmodifiers.UseStateForUnknownString(),
						},
					},

					"ignore_version_detection": schema.BoolAttribute{
						Required:            true,
						MarkdownDescription: "Select 'true' for apps that are automatically updated by app developer or to only check for app bundleID before installation. Select 'false' to check for app bundleID and version number before installation.",
					},
					"included_apps": schema.SetNestedAttribute{
						Optional: true,
						Computed: true,
						PlanModifiers: []planmodifier.Set{
							planmodifiers.UseStateForUnknownSet(),
						},
						MarkdownDescription: "Define the app bundle identifiers and version numbers to be used to detect the presence of the macOS app installation. This list is dynamically populated based on the PKG metadata, and users can also append additional entries. Maximum of 500 apps. +\n" +
							"\n" +
							"### Notes: +\n" +
							"- Included app bundle IDs (`CFBundleIdentifier`) and build numbers (`CFBundleShortVersionString`) are used for detecting and monitoring app installation status of the uploaded file. +\n" +
							"- The list should **only** contain the application(s) installed by the uploaded file in the `/Applications` folder on macOS. +\n" +
							"- Any other type of file that is not an application or is not installed in the `/Applications` folder should **not** be included. +\n" +
							"- If the list contains files that are not applications or none of the listed apps are installed, app installation status will **not** report success. +\n" +
							"- When multiple apps are present in the PKG, the **first app** in the list is used to identify the application. +\n" +
							"\n" +
							"### Example: +\n" +
							"To retrieve the `CFBundleIdentifier` and `CFBundleShortVersionString` of an installed application, you can use the macOS Terminal: +\n" +
							"\n" +
							"```bash +\n" +
							"# Retrieve the Bundle Identifier +\n" +
							"defaults read /Applications/Company\\ Portal.app/Contents/Info CFBundleIdentifier +\n" +
							"\n" +
							"# Retrieve the Short Version String +\n" +
							"defaults read /Applications/Company\\ Portal.app/Contents/Info CFBundleShortVersionString +\n" +
							"``` +\n" +
							"\n" +
							"Alternatively, these values can also be located in the `<app_name>.app/Contents/Info.plist` file inside the mounted PKG or DMG. +\n" +
							"\n" +
							"For apps added to Intune, the Intune admin center can also provide the app bundle ID. +\n",
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"bundle_id": schema.StringAttribute{
									Required:            true,
									MarkdownDescription: "The `CFBundleIdentifier` of the app as defined in the PKG metadata or appended manually.",
								},
								"bundle_version": schema.StringAttribute{
									Required:            true,
									MarkdownDescription: "The `CFBundleShortVersionString` of the app as defined in the PKG metadata or appended manually.",
								},
							},
						},
					},
					"minimum_supported_operating_system": schema.SingleNestedAttribute{
						Required:            true,
						MarkdownDescription: "Specifies the minimum macOS version required for the application. Each field indicates whether the version is supported.",
						Attributes: map[string]schema.Attribute{
							"v10_7": schema.BoolAttribute{
								Optional:            true,
								Computed:            true,
								Default:             booldefault.StaticBool(false),
								MarkdownDescription: "Application supports macOS 10.7 or later. Defaults to `false`.",
							},
							"v10_8": schema.BoolAttribute{
								Optional:            true,
								Computed:            true,
								Default:             booldefault.StaticBool(false),
								MarkdownDescription: "Application supports macOS 10.8 or later. Defaults to `false`.",
							},
							"v10_9": schema.BoolAttribute{
								Optional:            true,
								Computed:            true,
								Default:             booldefault.StaticBool(false),
								MarkdownDescription: "Application supports macOS 10.9 or later. Defaults to `false`.",
							},
							"v10_10": schema.BoolAttribute{
								Optional:            true,
								Computed:            true,
								Default:             booldefault.StaticBool(false),
								MarkdownDescription: "Application supports macOS 10.10 or later. Defaults to `false`.",
							},
							"v10_11": schema.BoolAttribute{
								Optional:            true,
								Computed:            true,
								Default:             booldefault.StaticBool(false),
								MarkdownDescription: "Application supports macOS 10.11 or later. Defaults to `false`.",
							},
							"v10_12": schema.BoolAttribute{
								Optional:            true,
								Computed:            true,
								Default:             booldefault.StaticBool(false),
								MarkdownDescription: "Application supports macOS 10.12 or later. Defaults to `false`.",
							},
							"v10_13": schema.BoolAttribute{
								Optional:            true,
								Computed:            true,
								Default:             booldefault.StaticBool(false),
								MarkdownDescription: "Application supports macOS 10.13 or later. Defaults to `false`.",
							},
							"v10_14": schema.BoolAttribute{
								Optional:            true,
								Computed:            true,
								Default:             booldefault.StaticBool(false),
								MarkdownDescription: "Application supports macOS 10.14 or later. Defaults to `false`.",
							},
							"v10_15": schema.BoolAttribute{
								Optional:            true,
								Computed:            true,
								Default:             booldefault.StaticBool(false),
								MarkdownDescription: "Application supports macOS 10.15 or later. Defaults to `false`.",
							},
							"v11_0": schema.BoolAttribute{
								Optional:            true,
								Computed:            true,
								Default:             booldefault.StaticBool(false),
								MarkdownDescription: "Application supports macOS 11.0 or later. Defaults to `false`.",
							},
							"v12_0": schema.BoolAttribute{
								Optional:            true,
								Computed:            true,
								Default:             booldefault.StaticBool(false),
								MarkdownDescription: "Application supports macOS 12.0 or later. Defaults to `false`.",
							},
							"v13_0": schema.BoolAttribute{
								Optional:            true,
								Computed:            true,
								Default:             booldefault.StaticBool(false),
								MarkdownDescription: "Application supports macOS 13.0 or later. Defaults to `false`.",
							},
							"v14_0": schema.BoolAttribute{
								Optional:            true,
								Computed:            true,
								Default:             booldefault.StaticBool(false),
								MarkdownDescription: "Application supports macOS 14.0 or later. Defaults to `false`.",
							},
						},
					},
					"pre_install_script": schema.SingleNestedAttribute{
						Optional: true,
						Attributes: map[string]schema.Attribute{
							"script_content": schema.StringAttribute{
								Required:            true,
								MarkdownDescription: "Base64 encoded shell script to execute on macOS device before app installation. Requires base64encode()",
							},
						},
					},
					"post_install_script": schema.SingleNestedAttribute{
						Optional: true,
						Attributes: map[string]schema.Attribute{
							"script_content": schema.StringAttribute{
								Required:            true,
								MarkdownDescription: "Base64 encoded shell script to execute on macOS device after app installation. Requires base64encode()",
							},
						},
					},
					"primary_bundle_id": schema.StringAttribute{
						Computed: true,
						// PlanModifiers: []planmodifier.String{
						// 	planmodifiers.UseStateForUnknownString(),
						// },
						MarkdownDescription: "The bundleId of the primary app in the PKG. Maps to CFBundleIdentifier in the app's bundle configuration.",
					},
					"primary_bundle_version": schema.StringAttribute{
						Computed: true,
						// PlanModifiers: []planmodifier.String{
						// 	planmodifiers.UseStateForUnknownString(),
						// },
						MarkdownDescription: "The version of the primary app in the PKG. Maps to CFBundleShortVersion in the app's bundle configuration.",
					},
				},
			},
			"installer_size_in_bytes": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "The size of the installer file in bytes. Used to detect changes in content.",
				PlanModifiers: []planmodifier.Int64{
					planmodifiers.UseStateForUnknownInt64(),
				},
			},
			"assignments": commonschemagraphbeta.MobileAppAssignmentSchema(),
			"timeouts":    commonschema.Timeouts(ctx),
		},
	}
}
