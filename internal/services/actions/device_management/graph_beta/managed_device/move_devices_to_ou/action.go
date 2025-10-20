package graphBetaMoveDevicesToOUManagedDevice

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
	"github.com/hashicorp/terraform-plugin-framework/types"

	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	ActionName = "graph_beta_device_management_managed_device_move_devices_to_ou"
)

var (
	_ action.Action                   = &MoveDevicesToOUManagedDeviceAction{}
	_ action.ActionWithConfigure      = &MoveDevicesToOUManagedDeviceAction{}
	_ action.ActionWithValidateConfig = &MoveDevicesToOUManagedDeviceAction{}
)

func NewMoveDevicesToOUManagedDeviceAction() action.Action {
	return &MoveDevicesToOUManagedDeviceAction{
		ReadPermissions: []string{
			"DeviceManagementConfiguration.Read.All",
			"DeviceManagementManagedDevices.Read.All",
		},
		WritePermissions: []string{
			"DeviceManagementConfiguration.Read.All",
			"DeviceManagementManagedDevices.Read.All",
		},
	}
}

type MoveDevicesToOUManagedDeviceAction struct {
	client           *msgraphbetasdk.GraphServiceClient
	ProviderTypeName string
	TypeName         string
	ReadPermissions  []string
	WritePermissions []string
}

func (a *MoveDevicesToOUManagedDeviceAction) Metadata(ctx context.Context, req action.MetadataRequest, resp *action.MetadataResponse) {
	a.ProviderTypeName = req.ProviderTypeName
	a.TypeName = ActionName
	resp.TypeName = a.FullTypeName()
}

func (a *MoveDevicesToOUManagedDeviceAction) FullTypeName() string {
	return a.ProviderTypeName + "_" + ActionName
}

func (a *MoveDevicesToOUManagedDeviceAction) Configure(ctx context.Context, req action.ConfigureRequest, resp *action.ConfigureResponse) {
	a.client = client.SetGraphBetaClientForAction(ctx, req, resp, constants.PROVIDER_NAME+"_"+ActionName)
}

func (a *MoveDevicesToOUManagedDeviceAction) Schema(ctx context.Context, req action.SchemaRequest, resp *action.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Moves hybrid Azure AD joined Windows devices to a specified Active Directory Organizational Unit (OU) using the " +
			"`/deviceManagement/managedDevices/moveDevicesToOU` and " +
			"`/deviceManagement/comanagedDevices/moveDevicesToOU` endpoints. " +
			"This action updates the organizational unit placement of devices in on-premises Active Directory for hybrid-joined devices. " +
			"The move operation is performed at the collection level, allowing multiple devices to be moved to the same OU in a single operation.\n\n" +
			"**Important Notes:**\n" +
			"- Only works on **Hybrid Azure AD joined** Windows devices\n" +
			"- Requires on-premises Active Directory connectivity\n" +
			"- Requires Azure AD Connect sync\n" +
			"- All devices are moved to the **same** OU path\n" +
			"- OU path must be valid in on-premises AD\n" +
			"- Changes reflect after next Azure AD Connect sync\n" +
			"- Does not affect cloud-only or Workplace-joined devices\n\n" +
			"**Use Cases:**\n" +
			"- Reorganizing device structure in Active Directory\n" +
			"- Applying different Group Policy Objects (GPOs)\n" +
			"- Moving devices between departments or locations\n" +
			"- Aligning device placement with organizational structure\n" +
			"- Consolidating devices for management purposes\n" +
			"- Preparing devices for different security policies\n\n" +
			"**Platform Support:**\n" +
			"- **Windows**: Hybrid Azure AD joined devices only\n" +
			"- **Other Platforms**: Not supported (cloud-only management)\n\n" +
			"**Reference:** [Microsoft Graph API - Move Devices to OU](https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-movedevicestoou?view=graph-rest-beta)",
		Attributes: map[string]schema.Attribute{
			"organizational_unit_path": schema.StringAttribute{
				Required: true,
				MarkdownDescription: "The full distinguished name path of the target Organizational Unit in Active Directory. " +
					"All specified devices will be moved to this OU.\n\n" +
					"**Format**: Must be a valid Active Directory OU distinguished name.\n\n" +
					"**Examples**:\n" +
					"- `\"OU=Workstations,OU=Computers,DC=contoso,DC=com\"`\n" +
					"- `\"OU=Marketing,OU=Departments,DC=example,DC=local\"`\n" +
					"- `\"OU=Laptops,OU=Mobile,OU=Devices,DC=corp,DC=acme,DC=com\"`\n\n" +
					"**Important**: The OU must exist in your on-premises Active Directory, and the Azure AD Connect sync account " +
					"must have permissions to move computer objects to this OU.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^(OU=|CN=)[^,]+(,(OU=|CN=|DC=)[^,]+)*$`),
						"organizational_unit_path must be a valid Active Directory distinguished name (e.g., OU=Computers,DC=contoso,DC=com)",
					),
				},
			},
			"managed_device_ids": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				MarkdownDescription: "List of managed device IDs (GUIDs) to move to the specified Organizational Unit. " +
					"These are devices fully managed by Intune that are also hybrid Azure AD joined.\n\n" +
					"**Note:** At least one of `managed_device_ids` or `comanaged_device_ids` must be provided. " +
					"All devices in this list will be moved to the same OU path specified in `organizational_unit_path`.\n\n" +
					"**Important:** Only hybrid Azure AD joined Windows devices can be moved. Cloud-only or workplace-joined devices will be ignored.\n\n" +
					"Example: `[\"12345678-1234-1234-1234-123456789abc\", \"87654321-4321-4321-4321-ba9876543210\"]`",
				Validators: []validator.List{
					listvalidator.ValueStringsAre(
						stringvalidator.RegexMatches(
							regexp.MustCompile(constants.GuidRegex),
							"each device ID must be a valid GUID format (e.g., 12345678-1234-1234-1234-123456789abc)",
						),
					),
				},
			},
			"comanaged_device_ids": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				MarkdownDescription: "List of co-managed device IDs (GUIDs) to move to the specified Organizational Unit. " +
					"These are devices managed by both Intune and Configuration Manager (SCCM) that are hybrid Azure AD joined.\n\n" +
					"**Note:** At least one of `managed_device_ids` or `comanaged_device_ids` must be provided. " +
					"All devices in this list will be moved to the same OU path.\n\n" +
					"Example: `[\"abcdef12-3456-7890-abcd-ef1234567890\"]`",
				Validators: []validator.List{
					listvalidator.ValueStringsAre(
						stringvalidator.RegexMatches(
							regexp.MustCompile(constants.GuidRegex),
							"each device ID must be a valid GUID format (e.g., 12345678-1234-1234-1234-123456789abc)",
						),
					),
				},
			},
			"timeouts": commonschema.Timeouts(ctx),
		},
	}
}
