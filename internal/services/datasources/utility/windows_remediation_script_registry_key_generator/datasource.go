package utilityWindowsRemediationScriptRegistryKeyGenerator

import (
	"context"

	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

const (
	DataSourceName = "windows_remediation_script_registry_key_generator"
	ReadTimeout    = 180
)

var (
	// Basic datasource interface (Read operations)
	_ datasource.DataSource = &WindowsRemediationScriptRegistryKeyGeneratorDataSource{}

	// Allows the datasource to be configured with the provider client
	_ datasource.DataSourceWithConfigure = &WindowsRemediationScriptRegistryKeyGeneratorDataSource{}
)

func NewWindowsRemediationScriptRegistryKeyGeneratorDataSource() datasource.DataSource {
	return &WindowsRemediationScriptRegistryKeyGeneratorDataSource{}
}

type WindowsRemediationScriptRegistryKeyGeneratorDataSource struct{}

func (d *WindowsRemediationScriptRegistryKeyGeneratorDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + DataSourceName
}

// For utility datasources that perform local computations. Required for interface compliance.
func (d *WindowsRemediationScriptRegistryKeyGeneratorDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
}

func (d *WindowsRemediationScriptRegistryKeyGeneratorDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Generates PowerShell detection and remediation scripts for Intune Proactive Remediations to manage Windows registry keys and values. " +
			"This utility helps create standardized scripts for setting registry values in either the current user's context (HKEY_CURRENT_USER) " +
			"or for all users on a device (HKEY_USERS). The generated scripts follow Microsoft's recommended patterns for Proactive Remediations, " +
			"including proper error handling, exit codes, and user context management.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Unique identifier for the data source (computed)",
				Computed:            true,
			},
			"context": schema.StringAttribute{
				MarkdownDescription: "The execution context for the registry operation. Valid values:\n" +
					"  - `current_user`: Applies to the currently logged-on user's registry hive (HKEY_CURRENT_USER)\n" +
					"  - `all_users`: Applies to all user profiles on the device (HKEY_USERS), excluding system accounts",
				Required: true,
				Validators: []validator.String{
					stringvalidator.OneOf("current_user", "all_users"),
				},
			},
			"registry_key_path": schema.StringAttribute{
				MarkdownDescription: "The registry key path relative to the user hive (without HKCU or HKU prefix). " +
					"Example: `Software\\Policies\\Microsoft\\WindowsStore\\` or `Software\\MyApp\\Settings`. " +
					"Use double backslashes (`\\\\`) to escape path separators.",
				Required: true,
			},
			"value_name": schema.StringAttribute{
				MarkdownDescription: "The name of the registry value to manage. Use `(Default)` for the default value of a key.",
				Required:            true,
			},
			"value_type": schema.StringAttribute{
				MarkdownDescription: "The registry value type. Valid values:\n" +
					"  - `REG_SZ`: String value\n" +
					"  - `REG_DWORD`: 32-bit integer value\n" +
					"  - `REG_QWORD`: 64-bit integer value\n" +
					"  - `REG_MULTI_SZ`: Multi-string value (separate strings with newlines)\n" +
					"  - `REG_EXPAND_SZ`: Expandable string value (can contain environment variables)\n" +
					"  - `REG_BINARY`: Binary data (provide as hex string, e.g., '01AF3C')",
				Required: true,
				Validators: []validator.String{
					stringvalidator.OneOf("REG_SZ", "REG_DWORD", "REG_QWORD", "REG_MULTI_SZ", "REG_EXPAND_SZ", "REG_BINARY"),
				},
			},
			"value_data": schema.StringAttribute{
				MarkdownDescription: "The desired value data. Format depends on value_type:\n" +
					"  - `REG_SZ`, `REG_EXPAND_SZ`: String (e.g., 'Enabled')\n" +
					"  - `REG_DWORD`: Decimal integer (e.g., '1', '0', '255')\n" +
					"  - `REG_QWORD`: Decimal integer (e.g., '1234567890')\n" +
					"  - `REG_MULTI_SZ`: Multiple strings separated by newlines\n" +
					"  - `REG_BINARY`: Hexadecimal string (e.g., '01AF3C')",
				Required: true,
			},
			"detection_script": schema.StringAttribute{
				MarkdownDescription: "Generated PowerShell detection script that checks if the registry value exists and matches the desired state. " +
					"Returns exit code 0 if compliant, 1 if remediation is needed. Use this as the detection script in Intune Proactive Remediations.",
				Computed: true,
			},
			"remediation_script": schema.StringAttribute{
				MarkdownDescription: "Generated PowerShell remediation script that creates or updates the registry key and value to the desired state. " +
					"Use this as the remediation script in Intune Proactive Remediations.",
				Computed: true,
			},
			"timeouts": commonschema.DatasourceTimeouts(ctx),
		},
	}
}
