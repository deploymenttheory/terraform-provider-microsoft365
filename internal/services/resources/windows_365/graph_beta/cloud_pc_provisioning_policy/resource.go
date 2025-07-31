package graphBetaCloudPcProvisioningPolicy

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	planmodifiers "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/plan_modifiers"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	ResourceName  = "graph_beta_windows_365_cloud_pc_provisioning_policy"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	// Basic resource interface (CRUD operations)
	_ resource.Resource = &CloudPcProvisioningPolicyResource{}

	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &CloudPcProvisioningPolicyResource{}

	// Enables import functionality
	_ resource.ResourceWithImportState = &CloudPcProvisioningPolicyResource{}

	// Enables plan modification/diff suppression
	_ resource.ResourceWithModifyPlan = &CloudPcProvisioningPolicyResource{}
)

func NewCloudPcProvisioningPolicyResource() resource.Resource {
	return &CloudPcProvisioningPolicyResource{
		ReadPermissions: []string{
			"CloudPC.Read.All",
		},
		WritePermissions: []string{
			"CloudPC.ReadWrite.All",
		},
	}
}

type CloudPcProvisioningPolicyResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ProviderTypeName string
	TypeName         string
	ReadPermissions  []string
	WritePermissions []string
}

// Metadata returns the resource type name.
func (r *CloudPcProvisioningPolicyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	r.ProviderTypeName = req.ProviderTypeName
	r.TypeName = ResourceName
	resp.TypeName = r.FullTypeName()
}

// FullTypeName returns the full resource type name in the format "providername_resourcename".
func (r *CloudPcProvisioningPolicyResource) FullTypeName() string {
	return r.ProviderTypeName + "_" + r.TypeName
}

// Configure sets the client for the resource.
func (r *CloudPcProvisioningPolicyResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, constants.PROVIDER_NAME+"_"+ResourceName)
}

// ImportState imports the resource state.
func (r *CloudPcProvisioningPolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Schema returns the schema for the resource.
func (r *CloudPcProvisioningPolicyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The unique identifier of the provisioning policy.",
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
			},
			"display_name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The display name for the provisioning policy.",
			},
			"alternate_resource_url": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The URL of the alternate resource that links to this provisioning policy. Read-only.",
			},
			"cloud_pc_group_display_name": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The display name of the Cloud PC group that the Cloud PCs reside in. Read-only.",
			},
			"cloud_pc_naming_template": schema.StringAttribute{
				Optional: true,
				MarkdownDescription: "Create unique names for your devices. Names must be between 5 and 15 characters, and can contain letters, numbers, and hyphens. Names cannot include a blank space. Use the %USERNAME:x% macro to add the first x letters of username. Use the %RAND:y% macro to add a random alphanumeric string of length y, y must be 5 or more. Names must contain a randomized string." +
					"For example, CPC-%USERNAME:4%-%RAND:5% means that the name of the Cloud PC starts with CPC-, followed by a four-character username, a - character, and then five random characters. The total length of the text generated by the template can't exceed 15 characters. Supports $filter, $select, and $orderby.",
			},
			"description": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The provisioning policy description. Supports $filter, $select, and $orderBy.",
			},
			"provisioning_type": schema.StringAttribute{
				Optional: true,
				Computed: true,
				MarkdownDescription: "Specifies the type of license used when provisioning Cloud PCs using this policy." +
					"By default, the license type is dedicated if the provisioningType isn't specified when you create the cloudPcProvisioningPolicy." +
					"Possible values are: dedicated, shared, sharedByUser, sharedByEntraGroup, unknownFutureValue." +
					"Changes to this attribute will force recreation of the resource." +
					"dedicated: (Enterprise) Each user will get their own Cloud PC without restrictions on when they can connect to it." +
					"shared: (Frontline - Dedicated) Recommended for users who need part time access to their Cloud PCs or follow a set schedule, such as shifts. A single license lets you provision up to three Cloud PCs that can be used non-concurrently, each assigned to a single user. Provides one concurrent session." +
					"sharedByEntraGroup: (Frontline - Shared) Recommended for users who use Cloud PC for a short period of time and do not require data to be preserved. A single license lets you provision one Cloud PC that can be shared non-concurrently among a group of users. Provides one concurrent session.",
				Default: stringdefault.StaticString("dedicated"),
				Validators: []validator.String{
					stringvalidator.OneOf("dedicated", "shared", "sharedByUser", "sharedByEntraGroup", "unknownFutureValue"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"domain_join_configurations": schema.ListNestedAttribute{
				Optional:            true,
				MarkdownDescription: "Specifies a list ordered by priority on how Cloud PCs join Microsoft Entra ID (Azure AD).",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"domain_join_type": schema.StringAttribute{
							Optional:            true,
							Computed:            true,
							MarkdownDescription: "Specifies the method by which the provisioned Cloud PC joins Microsoft Entra ID.",
							Default:             stringdefault.StaticString("azureADJoin"),
							Validators: []validator.String{
								stringvalidator.OneOf("azureADJoin", "hybridAzureADJoin", "unknownFutureValue"),
							},
						},
						"on_premises_connection_id": schema.StringAttribute{
							Optional:            true,
							MarkdownDescription: "The Azure network connection ID that matches the virtual network IT admins want the provisioning policy to use when they create Cloud PCs.",
						},
						"region_name": schema.StringAttribute{
							Optional: true,
							Computed: true,
							MarkdownDescription: "The supported Azure region where the IT admin wants the provisioning policy to create Cloud PCs. It is recommended using the Automatic option." +
								"The option allows Windows 365 to make the best selection which decreases the chance of provisioning failure.Must be one of: automatic, japaneast, eastasia.",
							Default: stringdefault.StaticString("automatic"),
							Validators: []validator.String{
								stringvalidator.OneOf("automatic", "japaneast", "eastasia", "southeastasia", "koreacentral"),
							},
						},
						"region_group": schema.StringAttribute{
							Required:            true,
							MarkdownDescription: "The Cloud PC region group. Must be one of: default, australia, canada, usCentral, usEast, usWest, france, germany, europeUnion, unitedKingdom, japan, asia, india, southAmerica, euap, usGovernment, usGovernmentDOD, unknownFutureValue, norway, switzerland, southKorea, middleEast, mexico, australasia, europe.",
							Validators: []validator.String{
								stringvalidator.OneOf(
									"default",
									"australia",
									"canada",
									"usCentral",
									"usEast",
									"usWest",
									"france",
									"germany",
									"europeUnion",
									"unitedKingdom",
									"japan",
									"asia",
									"india",
									"southAmerica",
									"euap",
									"usGovernment",
									"usGovernmentDOD",
									"unknownFutureValue",
									"norway",
									"switzerland",
									"southKorea",
									"middleEast",
									"mexico",
									"australasia",
									"europe",
								),
							},
						},
					},
				},
			},
			"enable_single_sign_on": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "True if the provisioned Cloud PC can be accessed by single sign-on. False indicates that the provisioned Cloud PC doesn't support this feature. The default value is false. ",
				Default:             booldefault.StaticBool(false),
			},
			"grace_period_in_hours": schema.Int32Attribute{
				Computed:            true,
				MarkdownDescription: "The number of hours to wait before reprovisioning/deprovisioning happens. Read-only.",
			},
			"image_display_name": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The display name of the operating system image that is used for provisioning. Supports $filter, $select, and $orderBy.",
			},
			"image_id": schema.StringAttribute{
				Required: true,
				MarkdownDescription: "The unique identifier that represents an operating system image used for provisioning new Cloud PCs. Must be one of:" +
					"'microsoftwindowsdesktop_windows-ent-cpc_win11-24H2-ent-cpc'," +
					"'microsoftwindowsdesktop_windows-ent-cpc_win11-24H2-ent-cpc-m365'," +
					"'microsoftwindowsdesktop_windows-ent-cpc_win11-23h2-ent-cpc-m365'," +
					"'microsoftwindowsdesktop_windows-ent-cpc_win11-23h2-ent-cpc'.",
				Validators: []validator.String{
					stringvalidator.OneOf(
						"microsoftwindowsdesktop_windows-ent-cpc_win11-24H2-ent-cpc",
						"microsoftwindowsdesktop_windows-ent-cpc_win11-24H2-ent-cpc-m365",
						"microsoftwindowsdesktop_windows-ent-cpc_win11-23h2-ent-cpc",
						"microsoftwindowsdesktop_windows-ent-cpc_win11-23h2-ent-cpc-m365",
					),
				},
			},
			"image_type": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The type of operating system image (custom or gallery) that is used for provisioning on Cloud PCs. Possible values are: gallery, custom. The default value is gallery. Supports $filter, $select, and $orderBy.",
				Default:             stringdefault.StaticString("gallery"),
				Validators: []validator.String{
					stringvalidator.OneOf("gallery", "custom"),
				},
			},
			"local_admin_enabled": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "When true, the local admin is enabled for Cloud PCs; false indicates that the local admin isn't enabled for Cloud PCs. The default value is false. Supports $filter, $select, and $orderBy.",
				Default:             booldefault.StaticBool(false),
			},
			"managed_by": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Specifies which service manages the Cloud PC provisioning policy. Possible values: windows365, devBox, unknownFutureValue, rpaBox. See Microsoft Graph documentation for details.",
				Default:             stringdefault.StaticString("windows365"),
				Validators: []validator.String{
					stringvalidator.OneOf("windows365", "devBox", "unknownFutureValue", "rpaBox"),
				},
			},
			"scope_ids": schema.SetAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Set of scope tag IDs for this resource.",
				PlanModifiers: []planmodifier.Set{
					planmodifiers.DefaultSetValue(
						[]attr.Value{types.StringValue("0")},
					),
				},
			},
			"autopatch": schema.SingleNestedAttribute{
				Optional:            true,
				MarkdownDescription: "Autopatch settings for the provisioning policy.",
				Attributes: map[string]schema.Attribute{
					"autopatch_group_id": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "The ID of the autopatch group.",
					},
				},
			},
			"autopilot_configuration": schema.SingleNestedAttribute{
				Optional:            true,
				MarkdownDescription: "Autopilot configuration for the provisioning policy.",
				Attributes: map[string]schema.Attribute{
					"device_preparation_profile_id": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "The ID of the device preparation profile.",
					},
					"application_timeout_in_minutes": schema.Int32Attribute{
						Optional:            true,
						MarkdownDescription: "The application timeout in minutes.",
					},
					"on_failure_device_access_denied": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Whether device access is denied on failure.",
					},
				},
			},
			"microsoft_managed_desktop": schema.SingleNestedAttribute{
				Required: true,
				MarkdownDescription: "This block is currently not supported in terraform. The registration to the Autopatch service currently only suports the intune gui," +
					"with no publically available api to call. This will return a 403 error if this block is used. Raise a ticket with Microsoft to make the Autopatch service available via api." +
					"The specific settings for Microsoft Managed Desktop that enables Microsoft Managed Desktop customers to get device managed experience for Cloud PC.",
				Attributes: map[string]schema.Attribute{
					"managed_type": schema.StringAttribute{
						Optional:            true,
						Computed:            true,
						MarkdownDescription: "Indicates the provisioning policy associated with Microsoft Managed Desktop settings.",
						Default:             stringdefault.StaticString("notManaged"),
						Validators: []validator.String{
							stringvalidator.OneOf("notManaged", "premiumManaged", "standardManaged", "starterManaged", "unknownFutureValue"),
						},
					},
					"profile": schema.StringAttribute{
						Optional: true,
						Computed: true,
						Default:  stringdefault.StaticString("4aa9b805-9494-4eed-a04b-ed51ec9e631e"),
						MarkdownDescription: "The name of the Microsoft Managed Desktop profile that the Windows 365 Cloud PC is associated with." +
							"'4aa9b805-9494-4eed-a04b-ed51ec9e631e' is the default Autopatch Group ID. Via 'https://mmdls.microsoft.com/device/v1/windows365/autopatchGroups'",
					},
				},
			},
			"windows_setting": schema.SingleNestedAttribute{
				Required:            true,
				MarkdownDescription: "Indicates a specific Windows setting to configure during the creation of Cloud PCs for this provisioning policy.",
				Attributes: map[string]schema.Attribute{
					"locale": schema.StringAttribute{
						Required:            true,
						MarkdownDescription: "The Windows language or region tag to use for language pack configuration and localization of the Cloud PC. The default value is en-US, which corresponds to English (United States).",
						Validators: []validator.String{
							stringvalidator.OneOf(
								"ar-SA",      // Arabic (Saudi Arabia)
								"bg-BG",      // Bulgarian (Bulgaria)
								"zh-CN",      // Chinese (Simplified, China)
								"zh-TW",      // Chinese (Traditional, Taiwan)
								"hr-HR",      // Croatian (Croatia)
								"cs-CZ",      // Czech (Czech Republic)
								"da-DK",      // Danish (Denmark)
								"nl-NL",      // Dutch (Netherlands)
								"en-AU",      // English (Australia)
								"en-IE",      // English (Ireland)
								"en-NZ",      // English (New Zealand)
								"en-GB",      // English (United Kingdom)
								"en-US",      // English (United States)
								"et-EE",      // Estonian (Estonia)
								"fi-FI",      // Finnish (Finland)
								"fr-CA",      // French (Canada)
								"fr-FR",      // French (France)
								"de-DE",      // German (Germany)
								"el-GR",      // Greek (Greece)
								"he-IL",      // Hebrew (Israel)
								"hu-HU",      // Hungarian (Hungary)
								"it-IT",      // Italian (Italy)
								"ja-JP",      // Japanese (Japan)
								"ko-KR",      // Korean (Korea)
								"lv-LV",      // Latvian (Latvia)
								"lt-LT",      // Lithuanian (Lithuania)
								"nb-NO",      // Norwegian Bokmål (Norway)
								"pl-PL",      // Polish (Poland)
								"pt-BR",      // Portuguese (Brazil)
								"pt-PT",      // Portuguese (Portugal)
								"ro-RO",      // Romanian (Romania)
								"ru-RU",      // Russian (Russia)
								"sr-Cyrl-CS", // Serbian (Cyrillic, Serbia)
								"sk-SK",      // Slovak (Slovakia)
								"sl-SI",      // Slovenian (Slovenia)
								"es-MX",      // Spanish (Mexico)
								"es-ES",      // Spanish (Spain)
								"sv-SE",      // Swedish (Sweden)
								"th-TH",      // Thai (Thailand)
								"tr-TR",      // Turkish (Turkey)
								"uk-UA",      // Ukrainian (Ukraine)
							),
						},
					},
				},
			},
			"apply_to_existing_cloud_pcs": schema.SingleNestedAttribute{
				Optional: true,
				MarkdownDescription: "If you change the network, image, region or single sign-on configuration in a provisioning policy," +
					"no change will occur for previously provisioned Cloud PCs. Newly provisioned or reprovisioned Cloud PCs will honor the" +
					"changes in your provisioning policy. To change the network or image of previously provisioned Cloud PCs to align with the changes," +
					"you must reprovision those Cloud PCs. To change the region or single sign-on of previously provisioned Cloud PCs to align with the changes," +
					"you must apply the changed configuration retrospectively using the apply_to_existing_cloud_pcs block.",
				Attributes: map[string]schema.Attribute{
					"microsoft_entra_single_sign_on_for_all_devices": schema.BoolAttribute{
						Computed:            true,
						Optional:            true,
						Default:             booldefault.StaticBool(false),
						MarkdownDescription: "When true, Microsoft Entra single sign-on is applied to all existing Cloud PCs. Applied only during resource updates and not during resource creation. Default is false.",
					},
					"region_or_azure_network_connection_for_all_devices": schema.BoolAttribute{
						Computed:            true,
						Optional:            true,
						Default:             booldefault.StaticBool(false),
						MarkdownDescription: "When true, region or Azure network connection settings are applied to all existing Cloud PCs. Applied only during resource updates and not during resource creation. Default is false.",
					},
					"region_or_azure_network_connection_for_select_devices": schema.BoolAttribute{
						Computed:            true,
						Optional:            true,
						Default:             booldefault.StaticBool(false),
						MarkdownDescription: "When true, region or Azure network connection settings are applied only to selected Cloud PCs. Applied only during resource updates and not during resource creation. Default is false.",
					},
				},
			},
			"assignments": Windows365ProvisioningPolicyAssignmentsSchema(),
			"timeouts":    commonschema.Timeouts(ctx),
		},
	}
}

// Windows365ProvisioningPolicyAssignmentsSchema returns the schema for the assignments attribute
func Windows365ProvisioningPolicyAssignmentsSchema() schema.SetNestedAttribute {
	return schema.SetNestedAttribute{
		Optional:            true,
		MarkdownDescription: "Assignments of the Cloud PC provisioning policy to groups. Only Microsoft 365 groups and security groups in Microsoft Entra ID are currently supported.",
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"type": schema.StringAttribute{
					Required:            true,
					MarkdownDescription: "The type of assignment target. Valid values are 'groupAssignmentTarget' only.",
					Validators: []validator.String{
						stringvalidator.OneOf(
							"groupAssignmentTarget",
						),
					},
				},
				"group_id": schema.StringAttribute{
					Required:            true,
					MarkdownDescription: "The ID of the Microsoft 365 group or security group in Microsoft Entra ID to assign the policy to.",
				},
				"service_plan_id": schema.StringAttribute{
					Optional: true,
					MarkdownDescription: "The ID of the frontlineservice plan. Required when provisioning_type is 'shared', 'sharedByUser', or 'sharedByEntraGroup'." +
						"This value can be obtained from the 'microsoft365_graph_beta_windows_365_cloud_pc_frontline_service_plan' data source.",
				},
				"allotment_license_count": schema.Int32Attribute{
					Optional: true,
					MarkdownDescription: "The number of licenses to allot. Required when provisioning_type is 'shared', 'sharedByUser', or 'sharedByEntraGroup'." +
						"The number must be between 0 and 900 and can't be more than the number of shared Cloud PC licenses available.",
				},
				"allotment_display_name": schema.StringAttribute{
					Optional:            true,
					MarkdownDescription: "A display name for the allotment. Required when provisioning_type is 'shared', 'sharedByUser', or 'sharedByEntraGroup'.",
				},
			},
		},
	}
}
