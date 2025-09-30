package graphBetaAssignUserToDevice

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/action/schema"

	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	ActionName = "graph_beta_device_management_windows_autopilot_device_identity_assign_user_to_device"
)

var (
	// Basic action interface
	_ action.Action = &AssignUserToDeviceAction{}

	// Allows the action to be configured with the provider client
	_ action.ActionWithConfigure = &AssignUserToDeviceAction{}
)

func NewAssignUserToDeviceAction() action.Action {
	return &AssignUserToDeviceAction{
		ReadPermissions: []string{
			"DeviceManagementServiceConfig.ReadWrite.All",
		},
		WritePermissions: []string{
			"DeviceManagementServiceConfig.ReadWrite.All",
		},
	}
}

type AssignUserToDeviceAction struct {
	client           *msgraphbetasdk.GraphServiceClient
	ProviderTypeName string
	TypeName         string
	ReadPermissions  []string
	WritePermissions []string
}

// Metadata returns the action type name.
func (a *AssignUserToDeviceAction) Metadata(ctx context.Context, req action.MetadataRequest, resp *action.MetadataResponse) {
	a.ProviderTypeName = req.ProviderTypeName
	a.TypeName = ActionName
	resp.TypeName = a.FullTypeName()
}

// FullTypeName returns the full action type name in the format "providername_actionname".
func (a *AssignUserToDeviceAction) FullTypeName() string {
	return a.ProviderTypeName + "_" + a.TypeName
}

// Configure sets the client for the action.
func (a *AssignUserToDeviceAction) Configure(ctx context.Context, req action.ConfigureRequest, resp *action.ConfigureResponse) {
	a.client = client.SetGraphBetaClientForAction(ctx, req, resp, constants.PROVIDER_NAME+"_"+ActionName)
}

// Schema returns the schema for the action.
func (a *AssignUserToDeviceAction) Schema(ctx context.Context, req action.SchemaRequest, resp *action.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Assigns a user to an Autopilot device in Microsoft Intune using the " +
			"`/deviceManagement/windowsAutopilotDeviceIdentities/{windowsAutopilotDeviceIdentityId}/assignUserToDevice` endpoint. " +
			"This action assigns user to Autopilot devices for streamlined device setup and management.",
		Attributes: map[string]schema.Attribute{
			"windows_autopilot_device_identity_id": schema.StringAttribute{
				Required: true,
				MarkdownDescription: "The unique identifier of the Windows Autopilot device identity to assign the user to. " +
					"This is the ID of the Windows Autopilot device in Microsoft Intune.",
			},
			"user_principal_name": schema.StringAttribute{
				Required: true,
				MarkdownDescription: "The user principal name (UPN) of the user to assign to the device. " +
					"This is typically the user's email address in the format user@domain.com.",
			},
			"addressable_user_name": schema.StringAttribute{
				Required: true,
				MarkdownDescription: "The addressable user name for the user being assigned to the device. " +
					"This is the display name or friendly name of the user.",
			},
			"timeouts": commonschema.Timeouts(ctx),
		},
	}
}
