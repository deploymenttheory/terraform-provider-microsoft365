package graphBetaWin32LobApp

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common"
	planmodifiers "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/plan_modifiers"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	ResourceName  = "graph_beta_device_and_app_management_win32_lob_app"
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
		},
		WritePermissions: []string{
			"DeviceManagementApps.ReadWrite.All",
		},
		ResourcePath: "/deviceAppManagement/mobileApps",
	}
}

type Win32LobAppResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ProviderTypeName string
	TypeName         string
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

// Metadata returns the resource type name.
func (r *Win32LobAppResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	r.ProviderTypeName = req.ProviderTypeName
	r.TypeName = ResourceName
	resp.TypeName = r.FullTypeName()
}

// FullTypeName returns the full resource type name in the format "providername_resourcename".
func (r *Win32LobAppResource) FullTypeName() string {
	return r.ProviderTypeName + "_" + r.TypeName
}

// Configure sets the client for the resource.
func (r *Win32LobAppResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = common.SetGraphBetaClientForResource(ctx, req, resp, constants.PROVIDER_NAME+"_"+ResourceName)
}

// ImportState imports the resource state.
func (r *Win32LobAppResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Function to create the full device management win32 lob app schema
func (r *Win32LobAppResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages Win32 line-of-business applications using the `/deviceAppManagement/mobileApps` endpoint. Win32 LOB apps enable deployment of custom Windows applications (.exe, .msi) with advanced installation logic, detection rules, and dependency management for enterprise software distribution.",
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
				Optional:            true,
				MarkdownDescription: "The description of the app.",
			},
			"publisher": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The publisher of the app.",
			},
			"large_icon": schema.SingleNestedAttribute{
				Optional:            true,
				MarkdownDescription: "The large icon, to be displayed in the app details and used for upload of the icon.",
				Attributes: map[string]schema.Attribute{
					"type": schema.StringAttribute{
						Required:            true,
						MarkdownDescription: "The MIME type of the icon.",
					},
					"value": schema.StringAttribute{
						Required:            true,
						MarkdownDescription: "The base64-encoded icon data.",
					},
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
			"upload_state": schema.Int64Attribute{
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
			"dependent_app_count": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "The total number of dependencies the child app has. This property is read-only.",
			},
			"superseding_app_count": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "The total number of apps this app directly or indirectly supersedes. This property is read-only.",
			},
			"superseded_app_count": schema.Int64Attribute{
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
				MarkdownDescription: "The command line to install this app",
			},
			"uninstall_command_line": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The command line to uninstall this app",
			},
			"applicable_architectures": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The Windows architecture(s) for which this app can run on. Possible values are: none, x86, x64, arm, neutral, arm64.",
				Validators: []validator.String{
					stringvalidator.OneOf("none", "x86", "x64", "arm", "neutral", "arm64"),
				},
			},
			"minimum_supported_operating_system": schema.SingleNestedAttribute{
				Optional:            true,
				MarkdownDescription: "The value for the minimum applicable operating system.",
				Attributes: map[string]schema.Attribute{
					"v8_0": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Windows 8.0 or later.",
					},
					"v8_1": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Windows 8.1 or later.",
					},
					"v10_0": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Windows 10.0 or later.",
					},
					"v10_1607": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Windows 10 1607 or later.",
					},
					"v10_1703": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Windows 10 1703 or later.",
					},
					"v10_1709": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Windows 10 1709 or later.",
					},
					"v10_1803": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Windows 10 1803 or later.",
					},
					"v10_1809": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Windows 10 1809 or later.",
					},
					"v10_1903": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Windows 10 1903 or later.",
					},
					"v10_1909": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Windows 10 1909 or later.",
					},
					"v10_2004": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Windows 10 2004 or later.",
					},
					"v10_2h20": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Windows 10 20H2 or later.",
					},
					"v10_21h1": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Windows 10 21H1 or later.",
					},
				},
			},
			"minimum_free_disk_space_in_mb": schema.Int64Attribute{
				Optional:            true,
				MarkdownDescription: "The value for the minimum free disk space which is required to install this app.",
			},
			"minimum_memory_in_mb": schema.Int64Attribute{
				Optional:            true,
				MarkdownDescription: "The value for the minimum physical memory which is required to install this app.",
			},
			"minimum_number_of_processors": schema.Int64Attribute{
				Optional:            true,
				MarkdownDescription: "The value for the minimum number of processors which is required to install this app.",
			},
			"minimum_cpu_speed_in_mhz": schema.Int64Attribute{
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
						"key_path": schema.StringAttribute{
							Optional:            true,
							MarkdownDescription: "The registry key path for registry detection.",
						},
						"value_name": schema.StringAttribute{
							Optional:            true,
							MarkdownDescription: "The registry value name for registry detection.",
						},
						"registry_detection_type": schema.StringAttribute{
							Optional:            true,
							MarkdownDescription: "The comparison operator for detection. Possible values are: notConfigured, exists, doesNotExist, string, integer, version. Applicable for registry detection.",
							Validators: []validator.String{
								stringvalidator.OneOf(
									"notConfigured", "exists", "doesNotExist", "string", "integer", "version",
								),
							},
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
						"filesystem_detection_type": schema.StringAttribute{
							Optional: true,
							MarkdownDescription: "The comparison operator for detection. Possible values are: notConfigured, exists," +
								"modifiedDate, createdDate, version, sizeInMB, doesNotExist.",
							Validators: []validator.String{
								stringvalidator.OneOf(
									"notConfigured", "exists", "modifiedDate", "createdDate", "version", "sizeInMB", "doesNotExist",
								),
							},
						},
						"filesystem_detection_operator": schema.StringAttribute{
							Optional:            true,
							MarkdownDescription: "The filesystem detection operator for filesystem detection. Possible values are: notConfigured, equal, notEqual, greaterThan, greaterThanOrEqual, lessThan, lessThanOrEqual. Used for registry and file_system detection types.",
							Validators: []validator.String{
								stringvalidator.OneOf("notConfigured", "equal", "notEqual", "greaterThan", "greaterThanOrEqual", "lessThan", "lessThanOrEqual"),
							},
						},
						// PowerShell Script Detection specific attributes
						"script_content": schema.StringAttribute{
							Optional:            true,
							MarkdownDescription: "The PowerShell script content to run for script detection.",
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
						"key_path": schema.StringAttribute{
							Optional:            true,
							MarkdownDescription: "The path to check for file or registry type.",
						},
						"file_or_folder_name": schema.StringAttribute{
							Optional:            true,
							MarkdownDescription: "The file or folder name to check for.",
						},
						"check_32_bit_on_64_system": schema.BoolAttribute{
							Optional:            true,
							MarkdownDescription: "A value indicating whether to check 32-bit on 64-bit system.",
						},
						"operator": schema.StringAttribute{
							Optional:            true,
							MarkdownDescription: "The operator for the requirement. Possible values are: notConfigured, equal, notEqual, greaterThan, greaterThanOrEqual, lessThan, lessThanOrEqual.",
						},
						"detection_value": schema.StringAttribute{
							Optional:            true,
							MarkdownDescription: "The value to check for.",
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
							MarkdownDescription: "The rule type.",
						},
						"check_32_bit_on_64_system": schema.BoolAttribute{
							Optional:            true,
							MarkdownDescription: "A value indicating whether to check 32-bit on 64-bit system.",
						},
						"key_path": schema.StringAttribute{
							Optional:            true,
							MarkdownDescription: "The registry key path to detect or check.",
						},
						"value_name": schema.StringAttribute{
							Optional:            true,
							MarkdownDescription: "The registry value name to detect or check.",
						},
						"operation_type": schema.StringAttribute{
							Optional:            true,
							MarkdownDescription: "The operation type for the rule.",
						},
						"operator": schema.StringAttribute{
							Optional:            true,
							MarkdownDescription: "The operator for the rule.",
						},
						"comparison_value": schema.StringAttribute{
							Optional:            true,
							MarkdownDescription: "The value to compare against.",
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
				},
			},
			"return_codes": schema.ListNestedAttribute{
				Optional:            true,
				MarkdownDescription: "The return codes for post installation behavior.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"return_code": schema.Int64Attribute{
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
				Optional:            true,
				MarkdownDescription: "The MSI details if this Win32 app is an MSI app.",
				Attributes: map[string]schema.Attribute{
					"product_code": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "The MSI product code.",
					},
					"product_version": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "The MSI product version.",
					},
					"upgrade_code": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "The MSI upgrade code.",
					},
					"requires_reboot": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "A value indicating whether the MSI app requires a reboot.",
					},
					"package_type": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "The MSI package type. Possible values are: perMachine, perUser.",
						Validators: []validator.String{
							stringvalidator.OneOf("perMachine", "perUser"),
						},
					},
				},
			},
			"setup_file_path": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The relative path of the setup file in the encrypted Win32LobApp package.",
			},
			"minimum_supported_windows_release": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The value for the minimum supported windows release.",
			},
			"display_version": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The version displayed in the UX for this app.",
			},
			"allow_available_uninstall": schema.BoolAttribute{
				Optional:            true,
				MarkdownDescription: "When TRUE, indicates that uninstall is supported from the company portal for the Windows app (Win32) with an Available assignment. When FALSE, indicates that uninstall is not supported for the Windows app (Win32) with an Available assignment. Default value is FALSE.",
			},
			"timeouts": commonschema.Timeouts(ctx),
		},
	}
}
