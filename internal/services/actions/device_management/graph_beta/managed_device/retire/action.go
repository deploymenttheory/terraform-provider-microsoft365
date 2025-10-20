package graphBetaRetireManagedDevice

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
	ActionName = "graph_beta_device_management_managed_device_retire"
)

var (
	_ action.Action                   = &RetireManagedDeviceAction{}
	_ action.ActionWithConfigure      = &RetireManagedDeviceAction{}
	_ action.ActionWithValidateConfig = &RetireManagedDeviceAction{}
)

func NewRetireManagedDeviceAction() action.Action {
	return &RetireManagedDeviceAction{
		ReadPermissions: []string{
			"DeviceManagementManagedDevices.PrivilegedOperations.All",
		},
		WritePermissions: []string{
			"DeviceManagementManagedDevices.PrivilegedOperations.All",
		},
	}
}

type RetireManagedDeviceAction struct {
	client           *msgraphbetasdk.GraphServiceClient
	ProviderTypeName string
	TypeName         string
	ReadPermissions  []string
	WritePermissions []string
}

func (a *RetireManagedDeviceAction) Metadata(ctx context.Context, req action.MetadataRequest, resp *action.MetadataResponse) {
	a.ProviderTypeName = req.ProviderTypeName
	a.TypeName = ActionName
	resp.TypeName = a.FullTypeName()
}

func (a *RetireManagedDeviceAction) FullTypeName() string {
	return a.ProviderTypeName + "_" + ActionName
}

func (a *RetireManagedDeviceAction) Configure(ctx context.Context, req action.ConfigureRequest, resp *action.ConfigureResponse) {
	a.client = client.SetGraphBetaClientForAction(ctx, req, resp, constants.PROVIDER_NAME+"_"+ActionName)
}

func (a *RetireManagedDeviceAction) Schema(ctx context.Context, req action.SchemaRequest, resp *action.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Retires managed devices from Microsoft Intune using the " +
			"`/deviceManagement/managedDevices/{managedDeviceId}/retire` endpoint. " +
			"This action removes company data and managed apps from the device, while leaving personal data intact. " +
			"The device is removed from Intune management and can no longer access company resources. " +
			"This action supports retiring multiple devices in a single operation.\n\n" +
			"**Important Notes:**\n" +
			"- For iOS/iPadOS devices, all data is removed except when enrolled via Device Enrollment Program (DEP) with User Affinity\n" +
			"- For Windows devices, company data under `%PROGRAMDATA%\\Microsoft\\MDM` is removed\n" +
			"- For Android devices, company data is removed and managed apps are uninstalled\n" +
			"- This action cannot be reversed - devices must be re-enrolled to be managed again\n\n" +
			"**Reference:** [Microsoft Graph API - Retire Managed Device](https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-retire?view=graph-rest-beta)",
		Attributes: map[string]schema.Attribute{
			"device_ids": schema.ListAttribute{
				ElementType: types.StringType,
				Required:    true,
				MarkdownDescription: "List of managed device IDs to retire from Intune management. " +
					"Each ID must be a valid GUID format. Multiple devices can be retired in a single action. " +
					"Example: `[\"12345678-1234-1234-1234-123456789abc\", \"87654321-4321-4321-4321-ba9876543210\"]`",
				Validators: []validator.List{
					listvalidator.SizeAtLeast(1),
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
