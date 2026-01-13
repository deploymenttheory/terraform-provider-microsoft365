package graphBetaUnassignUserFromDevice

import (
	"context"
	"regexp"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/action/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	ActionName    = "microsoft365_graph_beta_device_management_windows_autopilot_device_identity_unassign_user_from_device"
	InvokeTimeout = 60
)

var (
	// Basic action interface
	_ action.Action = &UnassignUserFromDeviceAction{}

	// Allows the action to be configured with the provider client
	_ action.ActionWithConfigure = &UnassignUserFromDeviceAction{}
)

func NewUnassignUserFromDeviceAction() action.Action {
	return &UnassignUserFromDeviceAction{
		ReadPermissions: []string{
			"DeviceManagementServiceConfig.Read.All",
			"DeviceManagementServiceConfig.ReadWrite.All",
		},
		WritePermissions: []string{
			"DeviceManagementServiceConfig.ReadWrite.All",
		},
	}
}

type UnassignUserFromDeviceAction struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
}

// Metadata returns the action type name.
func (a *UnassignUserFromDeviceAction) Metadata(ctx context.Context, req action.MetadataRequest, resp *action.MetadataResponse) {
	resp.TypeName = ActionName
}

// Description returns the action description.
func (a *UnassignUserFromDeviceAction) Description(ctx context.Context) string {
	return "Unassigns a user from an Autopilot device in Microsoft Intune using the " +
		"`/deviceManagement/windowsAutopilotDeviceIdentities/{windowsAutopilotDeviceIdentityId}/unassignUserFromDevice` endpoint. " +
		"This action removes the user assignment from Autopilot devices."
}

// Configure sets the client for the action.
func (a *UnassignUserFromDeviceAction) Configure(ctx context.Context, req action.ConfigureRequest, resp *action.ConfigureResponse) {
	a.client = client.SetGraphBetaClientForAction(ctx, req, resp, ActionName)
}

// Schema returns the schema for the action.
func (a *UnassignUserFromDeviceAction) Schema(ctx context.Context, req action.SchemaRequest, resp *action.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Unassigns a user from an Autopilot device in Microsoft Intune using the " +
			"`/deviceManagement/windowsAutopilotDeviceIdentities/{windowsAutopilotDeviceIdentityId}/unassignUserFromDevice` endpoint. " +
			"This action removes the user assignment from Autopilot devices.",
		Attributes: map[string]schema.Attribute{
			"windows_autopilot_device_identity_id": schema.StringAttribute{
				Required: true,
				MarkdownDescription: "The unique identifier of the Windows Autopilot device identity to unassign the user from. " +
					"This is the ID of the Windows Autopilot device in Microsoft Intune.",
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(constants.GuidRegex),
						"must be a valid GUID format (e.g., 12345678-1234-1234-1234-123456789abc)",
					),
				},
			},
			"timeouts": commonschema.ActionTimeouts(ctx),
		},
	}
}
