package graphBetaWin32App

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
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	ResourceName  = "microsoft365_graph_beta_device_and_app_management_win32_app"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	// Basic resource interface (CRUD operations)
	_ resource.Resource = &Win32LobAppResource{}

	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &Win32LobAppResource{}

	// Enables import functionality
	_ resource.ResourceWithImportState = &Win32LobAppResource{}

	// Enables plan modification/diff suppression
	_ resource.ResourceWithModifyPlan = &Win32LobAppResource{}
)

func NewWin32LobAppResource() resource.Resource {
	return &Win32LobAppResource{
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

type Win32LobAppResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

// Metadata returns the resource type name.
func (r *Win32LobAppResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = ResourceName
}

// Configure sets the client for the resource.
func (r *Win32LobAppResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, ResourceName)
}

// ImportState imports the resource state.
func (r *Win32LobAppResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Function to create the full device management win32 lob app schema
func (r *Win32LobAppResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages Win32 applications using the `/deviceAppManagement/mobileApps` endpoint. " +
			"Win apps enable deployment of custom Windows applications (.exe, .msi) with advanced installation logic, detection rules, " +
			"and dependency management for enterprise software distribution. They must be wrapped in the .intunewin file type." +
			"'https://learn.microsoft.com/en-us/intune/intune-service/apps/apps-win32-app-management'",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
				MarkdownDescription: "The unique identifier for this Intune win32 lob application",
			},
			"display_name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The admin provided or imported title of the app.",
			},
			"description": schema.StringAttribute{
				Required: true,
				MarkdownDescription: "A detailed description of the WinGet/ Microsoft Store for Business app." +
					"This field is automatically populated based on the package identifier when `automatically_generate_metadata` is set to true.",
				Validators: []validator.String{
					stringvalidator.LengthAtMost(10000),
				},
			},
			"publisher": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The publisher of the Intune macOS pkg application.",
				Validators: []validator.String{
					stringvalidator.LengthBetween(2, 1024),
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
			"created_date_time": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The date and time the app was created. This property is read-only.",
			},
			"last_modified_date_time": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The date and time the app was last modified. This property is read-only.",
			},
			"is_featured": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "The value indicating whether the app is marked as featured by the admin.",
			},
			"privacy_information_url": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The privacy statement Url.",
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(constants.HttpOrHttpsUrlRegex),
						"must be a valid URL starting with http:// or https://",
					),
				},
			},
			"information_url": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The more information Url.",
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(constants.HttpOrHttpsUrlRegex),
						"must be a valid URL starting with http:// or https://",
					),
				},
			},
			"owner": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The owner of the app.",
			},
			"developer": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The developer of the app.",
			},
			"notes": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Notes for the app.",
			},
			"upload_state": schema.Int32Attribute{
				Computed:            true,
				MarkdownDescription: "The upload state. Possible values are: 0 - Not Ready, 1 - Ready, 2 - Processing. This property is read-only.",
			},
			"publishing_state": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The publishing state for the app. The app cannot be assigned unless the app is published. This property is read-only. Possible values are: notPublished, processing, published.",
			},
			"is_assigned": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "The value indicating whether the app is assigned to at least one group. This property is read-only.",
			},
			"role_scope_tag_ids": schema.SetAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Set of scope tag IDs for this Settings Catalog template profile.",
				PlanModifiers: []planmodifier.Set{
					planmodifiers.DefaultSetValue(
						[]attr.Value{types.StringValue("0")},
					),
				},
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
			"committed_content_version": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The internal committed content version.",
			},
			"file_name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The name of the main Lob application file.",
			},
			"size": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "The total size, including all uploaded files. This property is read-only.",
			},
			"install_command_line": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The command line to install this app. Typically formatted as 'msiexec /i \"application_name.msi\" /qn'",
			},
			"uninstall_command_line": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The command line to uninstall this app. Typically formatted as 'msiexec /x {00000000-0000-0000-0000-000000000000} /qn'",
			},
			"allowed_architectures": schema.SetAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The Windows architecture(s) for which this app can run on. Possible values are: none, x64, x86, arm64.",
				Validators: []validator.Set{
					setvalidator.ValueStringsAre(
						stringvalidator.OneOf("none", "x64", "x86", "arm64"),
					),
				},
				PlanModifiers: []planmodifier.Set{
					planmodifiers.DefaultSetValue(
						[]attr.Value{types.StringValue("none")},
					),
				},
			},
			"minimum_free_disk_space_in_mb": schema.Int32Attribute{
				Optional:            true,
				MarkdownDescription: "The value for the minimum free disk space which is required to install this app.",
			},
			"minimum_memory_in_mb": schema.Int32Attribute{
				Optional:            true,
				MarkdownDescription: "The value for the minimum physical memory which is required to install this app.",
			},
			"minimum_number_of_processors": schema.Int32Attribute{
				Optional:            true,
				MarkdownDescription: "The value for the minimum number of processors which is required to install this app.",
			},
			"minimum_cpu_speed_in_mhz": schema.Int32Attribute{
				Optional:            true,
				MarkdownDescription: "The value for the minimum CPU speed which is required to install this app.",
			},
			"detection_rules": schema.ListNestedAttribute{
				Optional:            true,
				MarkdownDescription: "The detection rules to detect Win32 Line of Business (LoB) app.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						// Common attributes for all detection types
						"detection_type": schema.StringAttribute{
							Required:            true,
							MarkdownDescription: "The detection rule type. Possible values are: registry, msi_information, file_system, powershell_script.",
							Validators: []validator.String{
								stringvalidator.OneOf("registry", "msi_information", "file_system", "powershell_script"),
							},
						},
						"check_32_bit_on_64_system": schema.BoolAttribute{
							Optional:            true,
							MarkdownDescription: "Whether to check 32-bit registry or file system on 64-bit system. Applicable for registry, file_system, and PowerShell script detection.",
						},
						// Registry Detection specific attributes
						"registry_detection_type": schema.StringAttribute{
							Optional:            true,
							MarkdownDescription: "The registry detection type for registry detection. Possible values are: notConfigured, exists, doesNotExist, string, integer, version.",
							Validators: []validator.String{
								stringvalidator.OneOf("notConfigured", "exists", "doesNotExist", "string", "integer", "version"),
							},
						},
						"key_path": schema.StringAttribute{
							Optional:            true,
							MarkdownDescription: "The registry key path for registry detection.",
						},
						"value_name": schema.StringAttribute{
							Optional:            true,
							MarkdownDescription: "The registry value name for registry detection.",
						},
						"registry_detection_operator": schema.StringAttribute{
							Optional:            true,
							MarkdownDescription: "The registry detection operator for registry detection. Possible values are: notConfigured, equal, notEqual, greaterThan, greaterThanOrEqual, lessThan, lessThanOrEqual. Used for registry and file_system detection types.",
							Validators: []validator.String{
								stringvalidator.OneOf("notConfigured", "equal", "notEqual", "greaterThan", "greaterThanOrEqual", "lessThan", "lessThanOrEqual"),
							},
						},
						"detection_value": schema.StringAttribute{
							Optional:            true,
							MarkdownDescription: "The registry detection value for registry detection.",
						},
						// MSI Information Detection specific attributes
						"product_code": schema.StringAttribute{
							Optional:            true,
							MarkdownDescription: "The MSI product code for MSI detection.",
						},
						"product_version": schema.StringAttribute{
							Optional:            true,
							MarkdownDescription: "The MSI product version for MSI detection.",
						},
						"product_version_operator": schema.StringAttribute{
							Optional:            true,
							MarkdownDescription: "The MSI product version operator for MSI detection.Possible values are: notConfigured, equal, notEqual, greaterThan, greaterThanOrEqual, lessThan, lessThanOrEqual. Used for registry and file_system detection types.",
							Validators: []validator.String{
								stringvalidator.OneOf("notConfigured", "equal", "notEqual", "greaterThan", "greaterThanOrEqual", "lessThan", "lessThanOrEqual"),
							},
						},
						// File System Detection specific attributes
						"file_path": schema.StringAttribute{
							Optional:            true,
							MarkdownDescription: "The file path for file system detection.",
						},
						"file_or_folder_name": schema.StringAttribute{
							Optional:            true,
							MarkdownDescription: "The file name for file system detection.",
						},
						"file_system_detection_type": schema.StringAttribute{
							Optional: true,
							MarkdownDescription: "The comparison operator for detection. Possible values are: notConfigured, exists," +
								"modifiedDate, createdDate, version, sizeInMB, doesNotExist.",
							Validators: []validator.String{
								stringvalidator.OneOf(
									"notConfigured", "exists", "modifiedDate", "createdDate", "version", "sizeInMB", "doesNotExist",
								),
							},
						},
						"file_system_detection_operator": schema.StringAttribute{
							Optional:            true,
							MarkdownDescription: "The filesystem detection operator for filesystem detection. Possible values are: notConfigured, equal, notEqual, greaterThan, greaterThanOrEqual, lessThan, lessThanOrEqual. Used for registry and file_system detection types.",
							Validators: []validator.String{
								stringvalidator.OneOf("notConfigured", "equal", "notEqual", "greaterThan", "greaterThanOrEqual", "lessThan", "lessThanOrEqual"),
							},
						},
						// PowerShell Script Detection specific attributes
						"script_content": schema.StringAttribute{
							Optional: true,
							MarkdownDescription: "The PowerShell script content to run for script detection." +
								"This will be base64-encoded before being sent to the API. Supports PowerShell 5.1 and PowerShell 7.0.",
						},
						"enforce_signature_check": schema.BoolAttribute{
							Optional:            true,
							MarkdownDescription: "Whether to enforce signature checking for the PowerShell script.",
						},
						"run_as_32_bit": schema.BoolAttribute{
							Optional:            true,
							MarkdownDescription: "Whether to run the PowerShell script in 32-bit mode on 64-bit systems.",
						},
					},
				},
			},
			"requirement_rules": schema.ListNestedAttribute{
				Optional:            true,
				MarkdownDescription: "The requirement rules to detect Win32 Line of Business (LoB) app.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"requirement_type": schema.StringAttribute{
							Required:            true,
							MarkdownDescription: "The requirement rule type. Possible values are: registry, file, script.",
							Validators: []validator.String{
								stringvalidator.OneOf("registry", "file", "script"),
							},
						},
						"operator": schema.StringAttribute{
							Optional:            true,
							MarkdownDescription: "The operator for the requirement. Possible values are: notConfigured, equal, notEqual, greaterThan, greaterThanOrEqual, lessThan, lessThanOrEqual.",
						},
						"detection_value": schema.StringAttribute{
							Optional:            true,
							MarkdownDescription: "The value to check for.",
						},
						"check_32_bit_on_64_system": schema.BoolAttribute{
							Optional:            true,
							MarkdownDescription: "A value indicating whether to check 32-bit on 64-bit system.",
						},
						"key_path": schema.StringAttribute{
							Optional:            true,
							MarkdownDescription: "The registry key path for registry requirement.",
						},
						"value_name": schema.StringAttribute{
							Optional:            true,
							MarkdownDescription: "The registry value name for registry requirement.",
						},
						"detection_type": schema.StringAttribute{
							Optional:            true,
							MarkdownDescription: "The detection type for registry requirement.",
						},
						"file_or_folder_name": schema.StringAttribute{
							Optional:            true,
							MarkdownDescription: "The file or folder name to check for file requirement.",
						},
					},
				},
			},
			"rules": schema.ListNestedAttribute{
				Optional:            true,
				MarkdownDescription: "The detection and requirement rules for this app.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"rule_type": schema.StringAttribute{
							Required:            true,
							MarkdownDescription: "The rule type. Possible values are: detection, requirement.",
							Validators: []validator.String{
								stringvalidator.OneOf("detection", "requirement"),
							},
						},
						"rule_sub_type": schema.StringAttribute{
							Required:            true,
							MarkdownDescription: "The rule sub-type. Possible values are: registry, file_system, powershell_script.",
							Validators: []validator.String{
								stringvalidator.OneOf("registry", "file_system", "powershell_script"),
							},
						},
						"check_32_bit_on_64_system": schema.BoolAttribute{
							Optional:            true,
							MarkdownDescription: "A value indicating whether to check 32-bit on 64-bit system.",
						},
						// Common fields for all rule types
						"lob_app_rule_operator": schema.StringAttribute{
							Optional:            true,
							MarkdownDescription: "The operator for the rule. Possible values are: notConfigured, equal, notEqual, greaterThan, greaterThanOrEqual, lessThan, lessThanOrEqual.",
							Validators: []validator.String{
								stringvalidator.OneOf("notConfigured", "equal", "notEqual", "greaterThan", "greaterThanOrEqual", "lessThan", "lessThanOrEqual"),
							},
						},
						"comparison_value": schema.StringAttribute{
							Optional:            true,
							MarkdownDescription: "The value to compare against.",
						},
						// Registry rule specific fields
						"key_path": schema.StringAttribute{
							Optional:            true,
							MarkdownDescription: "The registry key path to detect or check. Required for registry rules.",
						},
						"value_name": schema.StringAttribute{
							Optional:            true,
							MarkdownDescription: "The registry value name to detect or check. Required for registry rules.",
						},
						"operation_type": schema.StringAttribute{
							Optional:            true,
							MarkdownDescription: "The operation type for registry rules. Possible values are: notConfigured, exists, doesNotExist, string, integer, version.",
							Validators: []validator.String{
								stringvalidator.OneOf("notConfigured", "exists", "doesNotExist", "string", "integer", "version"),
							},
						},
						// File system rule specific fields
						"path": schema.StringAttribute{
							Optional:            true,
							MarkdownDescription: "The path to the file or folder. Required for file_system rules.",
						},
						"file_or_folder_name": schema.StringAttribute{
							Optional:            true,
							MarkdownDescription: "The file or folder name to check. Required for file_system rules.",
						},
						"file_system_operation_type": schema.StringAttribute{
							Optional:            true,
							MarkdownDescription: "The operation type for file system rules. Possible values are: notConfigured, exists, modifiedDate, createdDate, version, sizeInMB, doesNotExist.",
							Validators: []validator.String{
								stringvalidator.OneOf("notConfigured", "exists", "modifiedDate", "createdDate", "version", "sizeInMB", "doesNotExist"),
							},
						},
						// PowerShell script rule specific fields
						"display_name": schema.StringAttribute{
							Optional:            true,
							MarkdownDescription: "The display name for PowerShell script rules.",
						},
						"enforce_signature_check": schema.BoolAttribute{
							Optional:            true,
							MarkdownDescription: "Whether to enforce signature checking for PowerShell script rules.",
						},
						"run_as_32_bit": schema.BoolAttribute{
							Optional:            true,
							MarkdownDescription: "Whether to run PowerShell scripts in 32-bit mode on 64-bit systems.",
						},
						"run_as_account": schema.StringAttribute{
							Optional:            true,
							MarkdownDescription: "The execution context for PowerShell scripts. Possible values are: system, user.",
							Validators: []validator.String{
								stringvalidator.OneOf("system", "user"),
							},
						},
						"script_content": schema.StringAttribute{
							Optional: true,
							MarkdownDescription: "The PowerShell script content to run for script detection." +
								"This will be base64-encoded before being sent to the API. Supports PowerShell 5.1 and PowerShell 7.0.",
						},
						"powershell_script_rule_operation_type": schema.StringAttribute{
							Optional:            true,
							MarkdownDescription: "The operation type for PowerShell script rules. Possible values are: notConfigured, string, dateTime, integer, float, version, boolean.",
							Validators: []validator.String{
								stringvalidator.OneOf("notConfigured", "string", "dateTime", "integer", "float", "version", "boolean"),
							},
						},
					},
				},
			},
			"install_experience": schema.SingleNestedAttribute{
				Optional:            true,
				MarkdownDescription: "The install experience for this app.",
				Attributes: map[string]schema.Attribute{
					"run_as_account": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "The execution context. Possible values are: system, user.",
						Validators: []validator.String{
							stringvalidator.OneOf("system", "user"),
						},
					},
					"device_restart_behavior": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "The device restart behavior. Possible values are: basedOnReturnCode, allow, suppress, force.",
						Validators: []validator.String{
							stringvalidator.OneOf("basedOnReturnCode", "allow", "suppress", "force"),
						},
					},
					"max_run_time_in_minutes": schema.Int32Attribute{
						Optional:            true,
						MarkdownDescription: "The maximum run time in minutes for the installation.",
					},
				},
			},
			"return_codes": schema.ListNestedAttribute{
				Optional:            true,
				MarkdownDescription: "The return codes for post installation behavior.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"return_code": schema.Int32Attribute{
							Required:            true,
							MarkdownDescription: "The return code.",
						},
						"type": schema.StringAttribute{
							Required:            true,
							MarkdownDescription: "The return code type. Possible values are: failed, success, softReboot, hardReboot, retry.",
							Validators: []validator.String{
								stringvalidator.OneOf("failed", "success", "softReboot", "hardReboot", "retry"),
							},
						},
					},
				},
			},
			"msi_information": schema.SingleNestedAttribute{
				Required:            true,
				MarkdownDescription: "The MSI details if this Win32 app is an MSI app.",
				Attributes: map[string]schema.Attribute{
					"package_type": schema.StringAttribute{
						Required:            true,
						MarkdownDescription: "The MSI package type. Possible values are: perMachine, perUser.",
						Validators: []validator.String{
							stringvalidator.OneOf("perMachine", "perUser"),
						},
					},
					"product_code": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "The MSI product code.",
					},
					"product_name": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "The MSI product name.",
					},
					"publisher": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "The MSI publisher.",
					},
					"product_version": schema.StringAttribute{
						Required:            true,
						MarkdownDescription: "The MSI product version.",
					},
					"upgrade_code": schema.StringAttribute{
						Required:            true,
						MarkdownDescription: "The MSI upgrade code.",
					},
					"requires_reboot": schema.BoolAttribute{
						Required:            true,
						MarkdownDescription: "A value indicating whether the MSI app requires a reboot.",
					},
				},
			},
			"setup_file_path": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The relative path of the setup file in the encrypted Win32LobApp package.",
			},
			"minimum_supported_windows_release": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The value for the minimum supported windows release.",
				Validators: []validator.String{
					stringvalidator.OneOf("Windows11_23H2", "Windows11_22H2", "Windows11_21H2", "Windows10_22H2", "Windows10_21H2", "21H1", "2H20", "2004", "1909", "1903", "1809", "1803", "1709", "1703", "1607"),
				},
			},
			"display_version": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The version displayed in the UX for this app.",
			},
			"allow_available_uninstall": schema.BoolAttribute{
				Required:            true,
				MarkdownDescription: "When TRUE, indicates that uninstall is supported from the company portal for the Windows app (Win32) with an Available assignment. When FALSE, indicates that uninstall is not supported for the Windows app (Win32) with an Available assignment. Default value is FALSE.",
			},
			"content_version": commonschemagraphbeta.MobileAppContentVersionSchema(),
			"app_installer":   commonschemagraphbeta.MobileAppWin32LobInstallerMetadataSchema(),
			"app_icon":        commonschemagraphbeta.MobileAppIconSchema(),
			"timeouts":        commonschema.ResourceTimeouts(ctx),
		},
	}
}
