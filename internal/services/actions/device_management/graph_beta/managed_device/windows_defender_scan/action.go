package graphBetaWindowsDefenderScan

import (
	"context"
	"regexp"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/action/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	ActionName = "graph_beta_device_management_managed_device_windows_defender_scan"
)

var (
	_ action.Action                   = &WindowsDefenderScanAction{}
	_ action.ActionWithConfigure      = &WindowsDefenderScanAction{}
	_ action.ActionWithValidateConfig = &WindowsDefenderScanAction{}
)

func NewWindowsDefenderScanAction() action.Action {
	return &WindowsDefenderScanAction{
		ReadPermissions: []string{
			"DeviceManagementManagedDevices.PrivilegedOperations.All",
		},
		WritePermissions: []string{
			"DeviceManagementManagedDevices.PrivilegedOperations.All",
		},
	}
}

type WindowsDefenderScanAction struct {
	client           *msgraphbetasdk.GraphServiceClient
	ProviderTypeName string
	TypeName         string
	ReadPermissions  []string
	WritePermissions []string
}

func (a *WindowsDefenderScanAction) Metadata(ctx context.Context, req action.MetadataRequest, resp *action.MetadataResponse) {
	a.ProviderTypeName = req.ProviderTypeName
	a.TypeName = ActionName
	resp.TypeName = a.FullTypeName()
}

func (a *WindowsDefenderScanAction) FullTypeName() string {
	return a.ProviderTypeName + "_" + ActionName
}

func (a *WindowsDefenderScanAction) Configure(ctx context.Context, req action.ConfigureRequest, resp *action.ConfigureResponse) {
	a.client = client.SetGraphBetaClientForAction(ctx, req, resp, constants.PROVIDER_NAME+"_"+ActionName)
}

func (a *WindowsDefenderScanAction) Schema(ctx context.Context, req action.SchemaRequest, resp *action.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Triggers an antivirus scan on Windows devices using Windows Defender (Microsoft Defender Antivirus) via the " +
			"`/deviceManagement/managedDevices/{managedDeviceId}/windowsDefenderScan` and " +
			"`/deviceManagement/comanagedDevices/{managedDeviceId}/windowsDefenderScan` endpoints. " +
			"This action initiates either a quick scan or full scan remotely on Windows devices managed by Intune.\n\n" +
			"**What This Action Does:**\n" +
			"- Triggers immediate Windows Defender scan\n" +
			"- Supports both quick and full scan types\n" +
			"- Scans for viruses, malware, and threats\n" +
			"- Updates threat definitions before scanning\n" +
			"- Reports results to Intune\n" +
			"- Can be used for threat remediation\n" +
			"- Works on managed and co-managed devices\n\n" +
			"**Scan Types:**\n" +
			"- **Quick Scan**: Scans common threat locations (5-15 minutes)\n" +
			"  - System folders and registry keys\n" +
			"  - Active memory processes\n" +
			"  - Startup locations\n" +
			"  - Recommended for routine scans\n" +
			"- **Full Scan**: Comprehensive scan of entire system (30+ minutes to hours)\n" +
			"  - All files and folders\n" +
			"  - All drives and partitions\n" +
			"  - Archive files\n" +
			"  - Recommended when threat detected or troubleshooting\n\n" +
			"**Platform Support:**\n" +
			"- **Windows 10/11**: Full support (managed and co-managed)\n" +
			"- **Windows Server**: Full support (if Defender enabled)\n" +
			"- **Other platforms**: Not supported (Windows Defender only)\n\n" +
			"**Common Use Cases:**\n" +
			"- Security incident response\n" +
			"- Threat detection and remediation\n" +
			"- Compliance verification\n" +
			"- Post-malware cleanup\n" +
			"- Routine security checks\n" +
			"- After suspicious activity\n" +
			"- Emergency threat scanning\n\n" +
			"**Important Considerations:**\n" +
			"- Device must be online\n" +
			"- Full scans can impact performance\n" +
			"- Scans run in background\n" +
			"- Results reported to Intune\n" +
			"- May require user notification\n" +
			"- Can be resource-intensive\n\n" +
			"**Reference:** [Microsoft Graph API - Windows Defender Scan](https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-windowsdefenderscan?view=graph-rest-beta)",
		Attributes: map[string]schema.Attribute{
			"managed_devices": schema.ListNestedAttribute{
				Optional: true,
				MarkdownDescription: "List of managed Windows devices to scan with individual scan type configuration. " +
					"Each entry specifies a device ID and whether to perform a quick scan or full scan. " +
					"These are devices fully managed by Intune only.\n\n" +
					"Example:\n" +
					"```hcl\n" +
					"managed_devices = [\n" +
					"  {\n" +
					"    device_id  = \"12345678-1234-1234-1234-123456789abc\"\n" +
					"    quick_scan = true  # Quick scan (5-15 min)\n" +
					"  },\n" +
					"  {\n" +
					"    device_id  = \"87654321-4321-4321-4321-ba9876543210\"\n" +
					"    quick_scan = false # Full scan (30+ min)\n" +
					"  }\n" +
					"]\n" +
					"```\n\n" +
					"**Note:** At least one of `managed_devices` or `comanaged_devices` must be provided.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"device_id": schema.StringAttribute{
							Required: true,
							MarkdownDescription: "The managed device ID (GUID) of the Windows device to scan. " +
								"Example: `12345678-1234-1234-1234-123456789abc`",
							Validators: []validator.String{
								stringvalidator.RegexMatches(
									regexp.MustCompile(constants.GuidRegex),
									"device_id must be a valid GUID format (e.g., 12345678-1234-1234-1234-123456789abc)",
								),
							},
						},
						"quick_scan": schema.BoolAttribute{
							Required: true,
							MarkdownDescription: "Whether to perform a quick scan (`true`) or full scan (`false`).\n\n" +
								"- **Quick Scan (`true`)**: Fast scan of common threat locations (5-15 minutes)\n" +
								"  - Scans system folders, registry, memory, startup locations\n" +
								"  - Minimal impact on device performance\n" +
								"  - Recommended for routine/scheduled scans\n" +
								"  - Good for rapid security checks\n\n" +
								"- **Full Scan (`false`)**: Comprehensive scan of entire system (30+ minutes to hours)\n" +
								"  - Scans all files, folders, drives, archives\n" +
								"  - Higher impact on device performance\n" +
								"  - Recommended when threat detected\n" +
								"  - Thorough investigation of suspicious activity\n" +
								"  - Post-incident verification",
						},
					},
				},
				Validators: []validator.List{
					listvalidator.SizeAtLeast(1),
				},
			},
			"comanaged_devices": schema.ListNestedAttribute{
				Optional: true,
				MarkdownDescription: "List of co-managed Windows devices to scan with individual scan type configuration. " +
					"These are devices managed by both Intune and Configuration Manager (SCCM). " +
					"Each entry specifies a device ID and scan type.\n\n" +
					"**Co-Management Context:**\n" +
					"- Devices managed by both Intune and Configuration Manager\n" +
					"- Typically Windows 10/11 enterprise devices\n" +
					"- This action triggers Defender scan via Intune endpoint\n" +
					"- ConfigMgr can also trigger scans independently\n\n" +
					"**Note:** At least one of `managed_devices` or `comanaged_devices` must be provided.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"device_id": schema.StringAttribute{
							Required: true,
							MarkdownDescription: "The co-managed device ID (GUID) of the Windows device to scan. " +
								"Example: `12345678-1234-1234-1234-123456789abc`",
							Validators: []validator.String{
								stringvalidator.RegexMatches(
									regexp.MustCompile(constants.GuidRegex),
									"device_id must be a valid GUID format (e.g., 12345678-1234-1234-1234-123456789abc)",
								),
							},
						},
						"quick_scan": schema.BoolAttribute{
							Required: true,
							MarkdownDescription: "Whether to perform a quick scan (`true`) or full scan (`false`). " +
								"See managed_devices.quick_scan for detailed explanation of scan types.",
						},
					},
				},
				Validators: []validator.List{
					listvalidator.SizeAtLeast(1),
				},
			},
			"timeouts": commonschema.Timeouts(ctx),
		},
	}
}
