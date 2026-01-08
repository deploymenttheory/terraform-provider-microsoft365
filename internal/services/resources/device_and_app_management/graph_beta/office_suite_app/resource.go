package graphBetaOfficeSuiteApp

import (
	"context"
	"regexp"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	planmodifiers "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/plan_modifiers"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	commonschemagraphbeta "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema/graph_beta/device_and_app_management"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	ResourceName  = "microsoft365_graph_beta_device_and_app_management_office_suite_app"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	// Basic resource interface (CRUD operations)
	_ resource.Resource = &OfficeSuiteAppResource{}

	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &OfficeSuiteAppResource{}

	// Enables import functionality
	_ resource.ResourceWithImportState = &OfficeSuiteAppResource{}

	// Enables plan modification/diff suppression
	_ resource.ResourceWithModifyPlan = &OfficeSuiteAppResource{}
)

func NewOfficeSuiteAppResource() resource.Resource {
	return &OfficeSuiteAppResource{
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

type OfficeSuiteAppResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

// Metadata returns the resource type name.
func (r *OfficeSuiteAppResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = ResourceName
}

// Configure sets the client for the resource.
func (r *OfficeSuiteAppResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, ResourceName)
}

// ImportState imports the resource state.
func (r *OfficeSuiteAppResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Schema returns the schema for the resource.
func (r *OfficeSuiteAppResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages Microsoft 365 Apps (Office Suite) applications using the `/deviceAppManagement/mobileApps` endpoint." +
			"Office Suite Apps enable deployment of Microsoft 365 office applications with configuration options including app exclusions, update channels, " +
			"localization settings, and shared computer activation for enterprise environments. Learn more here 'https://learn.microsoft.com/en-us/intune/intune-service/apps/apps-add-office365'",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
				MarkdownDescription: "The unique identifier for this Microsoft 365 Apps application",
			},
			"display_name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The title of the Microsoft 365 Apps application.",
			},
			"description": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Required. The description of the resource. Maximum length is 1500 characters.",
				Validators: []validator.String{
					stringvalidator.LengthAtMost(10000),
				},
			},
			"publisher": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The publisher of the Microsoft 365 Apps application. Typically 'Microsoft'.",
			},
			"role_scope_tag_ids": schema.SetAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Set of scope tag IDs for this Office Suite app.",
				Default: setdefault.StaticValue(
					types.SetValueMust(
						types.StringType,
						[]attr.Value{types.StringValue("0")},
					),
				),
			},
			"categories": schema.SetAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				MarkdownDescription: "Set of category names to associate with this application. You can use either the predefined Intune category names like 'Business', 'Productivity', etc., or provide specific category UUIDs. Predefined values include: 'Other apps', 'Books & Reference', 'Data management', 'Productivity', 'Business', 'Development & Design', 'Photos & Media', 'Collaboration & Social', 'Computer management'.",
				Validators: []validator.Set{
					setvalidator.ValueStringsAre(
						stringvalidator.RegexMatches(
							regexp.MustCompile(`^(Other apps|Books & Reference|Data management|Productivity|Business|Development & Design|Photos & Media|Collaboration & Social|Computer management|[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12})$`),
							"must be either a predefined category name or a valid GUID in the format xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
						),
					),
				},
			},
			"is_featured": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "The value indicating whether the app is marked as featured by the admin. Default is false.",
			},
			"privacy_information_url": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The privacy statement URL. This is automatically set to Microsoft's privacy statement URL.",
			},
			"information_url": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The more information URL.",
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(constants.HttpOrHttpsUrlRegex),
						"must be a valid URL starting with http:// or https://",
					),
				},
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
				Optional:            true,
				MarkdownDescription: "Notes for the app.",
			},
			"app_icon": commonschemagraphbeta.MobileAppIconSchema(),

			// Office Suite App configuration blocks (mutually exclusive)
			"configuration_designer": schema.SingleNestedAttribute{
				Optional:            true,
				MarkdownDescription: "Configuration Designer block for Office Suite App. Use this to configure Office applications using individual settings.",
				Attributes: map[string]schema.Attribute{
					"auto_accept_eula": schema.BoolAttribute{
						Optional:            true,
						Computed:            true,
						Default:             booldefault.StaticBool(true),
						MarkdownDescription: "The value to accept the EULA automatically on the enduser's device. Default is true.",
					},
					"excluded_apps": schema.SingleNestedAttribute{
						Optional:            true,
						MarkdownDescription: "The Office applications to exclude from the installation.",
						Attributes: map[string]schema.Attribute{
							"access": schema.BoolAttribute{
								Optional:            true,
								Computed:            true,
								Default:             booldefault.StaticBool(false),
								MarkdownDescription: "The value for if MS Office Access should be excluded or not. Default is false.",
							},
							"bing": schema.BoolAttribute{
								Optional:            true,
								Computed:            true,
								Default:             booldefault.StaticBool(false),
								MarkdownDescription: "The value for if Microsoft Search as default in Bing should be excluded or not. Default is false.",
							},
							"excel": schema.BoolAttribute{
								Optional:            true,
								Computed:            true,
								Default:             booldefault.StaticBool(false),
								MarkdownDescription: "The value for if MS Office Excel should be excluded or not. Default is false.",
							},
							"groove": schema.BoolAttribute{
								Optional:            true,
								Computed:            true,
								Default:             booldefault.StaticBool(false),
								MarkdownDescription: "The value for if MS Office OneDrive for Business – Groove should be excluded or not. Default is false.",
							},
							"info_path": schema.BoolAttribute{
								Optional:            true,
								Computed:            true,
								Default:             booldefault.StaticBool(false),
								MarkdownDescription: "The value for if MS Office InfoPath should be excluded or not. Default is false.",
							},
							"lync": schema.BoolAttribute{
								Optional:            true,
								Computed:            true,
								Default:             booldefault.StaticBool(false),
								MarkdownDescription: "The value for if MS Office Skype for Business – Lync should be excluded or not. Default is false.",
							},
							"one_drive": schema.BoolAttribute{
								Optional:            true,
								Computed:            true,
								Default:             booldefault.StaticBool(false),
								MarkdownDescription: "The value for if MS Office OneDrive should be excluded or not. Default is false.",
							},
							"one_note": schema.BoolAttribute{
								Optional:            true,
								Computed:            true,
								Default:             booldefault.StaticBool(false),
								MarkdownDescription: "The value for if MS Office OneNote should be excluded or not. Default is false.",
							},
							"outlook": schema.BoolAttribute{
								Optional:            true,
								Computed:            true,
								Default:             booldefault.StaticBool(false),
								MarkdownDescription: "The value for if MS Office Outlook should be excluded or not. Default is false.",
							},
							"power_point": schema.BoolAttribute{
								Optional:            true,
								Computed:            true,
								Default:             booldefault.StaticBool(false),
								MarkdownDescription: "The value for if MS Office PowerPoint should be excluded or not. Default is false.",
							},
							"publisher": schema.BoolAttribute{
								Optional:            true,
								Computed:            true,
								Default:             booldefault.StaticBool(false),
								MarkdownDescription: "The value for if MS Office Publisher should be excluded or not. Default is false.",
							},
							"share_point_designer": schema.BoolAttribute{
								Optional:            true,
								Computed:            true,
								Default:             booldefault.StaticBool(false),
								MarkdownDescription: "The value for if MS Office SharePoint Designer should be excluded or not. Default is false.",
							},
							"teams": schema.BoolAttribute{
								Optional:            true,
								Computed:            true,
								Default:             booldefault.StaticBool(false),
								MarkdownDescription: "The value for if MS Office Teams should be excluded or not. Default is false.",
							},
							"visio": schema.BoolAttribute{
								Optional:            true,
								Computed:            true,
								Default:             booldefault.StaticBool(false),
								MarkdownDescription: "The value for if MS Office Visio should be excluded or not. Default is false.",
							},
							"word": schema.BoolAttribute{
								Optional:            true,
								Computed:            true,
								Default:             booldefault.StaticBool(false),
								MarkdownDescription: "The value for if MS Office Word should be excluded or not. Default is false.",
							},
						},
					},
					"locales_to_install": schema.SetAttribute{
						ElementType: types.StringType,
						Optional:    true,
						MarkdownDescription: "By default, Intune will install Office with the default language of the operating system. Choose any additional languages that you want to install." +
							"Must be one of the supported Office locale codes in the format 'xx-xx' (e.g., 'en-us', 'ja-jp').",
						Validators: []validator.Set{
							setvalidator.ValueStringsAre(
								stringvalidator.OneOf(
									"ar-sa", "bg-bg", "zh-cn", "zh-tw", "hr-hr", "cs-cz", "da-dk", "nl-nl", "en-us", "et-ee",
									"fi-fi", "fr-fr", "de-de", "el-gr", "he-il", "hi-in", "hu-hu", "id-id", "it-it", "ja-jp",
									"kk-kz", "ko-kr", "lv-lv", "lt-lt", "ms-my", "nb-no", "pl-pl", "pt-br", "pt-pt", "ro-ro",
									"ru-ru", "sr-latn-rs", "sk-sk", "sl-si", "es-es", "sv-se", "th-th", "tr-tr", "uk-ua", "vi-vn",
									"af-za", "sq-al", "am-et", "hy-am", "as-in", "az-latn-az", "bn-bd", "bn-in", "eu-es", "be-by",
									"bs-latn-ba", "ca-es", "prs-af", "fil-ph", "gl-es", "ka-ge", "gu-in", "is-is", "ga-ie", "kn-in",
									"km-kh", "sw-ke", "kok-in", "ky-kg", "lb-lu", "mk-mk", "ml-in", "mt-mt", "mi-nz", "mr-in",
									"mn-mn", "ne-np", "nn-no", "or-in", "fa-ir", "pa-in", "quz-pe", "gd-gb", "sr-cyrl-ba", "sr-cyrl-rs",
									"sd-arab-pk", "si-lk", "ta-in", "tt-ru", "te-in", "tk-tm", "ur-pk", "ug-cn", "uz-latn-uz",
									"ca-es-valencia", "cy-gb", "ha-latn-ng", "ig-ng", "xh-za", "zu-za", "rw-rw", "ps-af", "rm-ch",
									"nso-za", "tn-za", "wo-sn", "yo-ng",
								),
							),
						},
					},
					"office_platform_architecture": schema.StringAttribute{
						Optional:            true,
						Computed:            true,
						MarkdownDescription: "The architecture for which to install Office. Possible values are: 'x86', 'x64'. Default is 'x64'. Changing this forces a new resource to be created.",
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
						Validators: []validator.String{
							stringvalidator.OneOf("x86", "x64"),
						},
					},
					"office_suite_app_default_file_format": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "The default file format for Office applications. Possible values are: 'officeOpenXMLFormat', 'officeOpenDocumentFormat'.",
						Validators: []validator.String{
							stringvalidator.OneOf("officeOpenXMLFormat", "officeOpenDocumentFormat"),
						},
					},
					"product_ids": schema.SetAttribute{
						ElementType:         types.StringType,
						Optional:            true,
						MarkdownDescription: "The Product IDs that represent the Office suite app. Example values: 'o365ProPlusRetail', 'o365BusinessRetail','projectProRetail', 'visioProRetail'.",
						Validators: []validator.Set{
							setvalidator.ValueStringsAre(
								stringvalidator.OneOf(
									"o365ProPlusRetail", "o365BusinessRetail", "visioProRetail", "projectProRetail",
								),
							),
						},
					},
					"should_uninstall_older_versions_of_office": schema.BoolAttribute{
						Optional:            true,
						Computed:            true,
						Default:             booldefault.StaticBool(false),
						MarkdownDescription: "The value to uninstall any existing MSI versions of Office. Default is false.",
					},
					"target_version": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "The specific version of Office to install. Example: '16.0.19029.20244'.",
					},
					"update_version": schema.StringAttribute{
						Optional:            true,
						Computed:            true,
						Default:             stringdefault.StaticString(""),
						MarkdownDescription: "The specific update version for the Office installation. Example: '2507'. For latest version, use empty string (default).",
					},
					"update_channel": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "The Office update channel. Possible values are: 'current', 'deferred', 'firstReleaseCurrent', 'firstReleaseDeferred', 'monthlyEnterprise'.",
						Validators: []validator.String{
							stringvalidator.OneOf("current", "deferred", "firstReleaseCurrent", "firstReleaseDeferred", "monthlyEnterprise"),
						},
					},
					"use_shared_computer_activation": schema.BoolAttribute{
						Optional:            true,
						Computed:            true,
						Default:             booldefault.StaticBool(false),
						MarkdownDescription: "The value to enable shared computer activation for Office. Shared computer activation lets you deploy Microsoft 365 Apps to computers that are used by multiple users. Normally, users can only install and activate Microsoft 365 Apps on a limited number of devices, such as 5 PCs. Using Microsoft 365 Apps with shared computer activation doesn't count against that limit. Default is false.",
					},
				},
			},
			"xml_configuration": schema.SingleNestedAttribute{
				Optional: true,
				MarkdownDescription: "XML Configuration block for Office Suite App. Use this to configure Office applications using XML configuration. " +
					"Learn more here'https://learn.microsoft.com/en-us/microsoft-365-apps/deploy/office-deployment-tool-configuration-options'.",
				Attributes: map[string]schema.Attribute{
					"office_configuration_xml": schema.StringAttribute{
						Required:            true,
						MarkdownDescription: "The XML configuration file for Office deployment. This is base64 encoded XML content that defines the Office installation configuration.",
					},
				},
			},
			// Common computed fields
			"created_date_time": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The date and time the app was created. This property is read-only.",
			},
			"last_modified_date_time": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The date and time the app was last modified. This property is read-only.",
			},
			"upload_state": schema.Int32Attribute{
				Computed:            true,
				MarkdownDescription: "The upload state. Possible values are: 0 - Not Ready, 1 - Ready, 2 - Processing. This property is read-only.",
			},
			"publishing_state": schema.StringAttribute{
				Computed: true,
				MarkdownDescription: "The publishing state for the app. The app cannot be assigned unless the app is published. " +
					"Possible values are: notPublished, processing, published.",
			},
			"is_assigned": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "The value indicating whether the app is assigned to at least one group. This property is read-only.",
			},
			"dependent_app_count": schema.Int32Attribute{
				Computed:            true,
				MarkdownDescription: "The total number of dependencies the child app has. This property is read-only.",
			},
			"superseding_app_count": schema.Int32Attribute{
				Computed:            true,
				MarkdownDescription: "The total number of apps this app directly or indirectly supersedes. This property is read-only.",
			},
			"superseded_app_count": schema.Int32Attribute{
				Computed:            true,
				MarkdownDescription: "The total number of apps this app is directly or indirectly superseded by. This property is read-only.",
			},
			"timeouts": commonschema.ResourceTimeouts(ctx),
		},
	}
}
