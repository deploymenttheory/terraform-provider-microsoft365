package graphBetaUpdateDeviceProperties

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
	ActionName = "graph_beta_device_management_windows_autopilot_device_identity_update_device_properties"
)

var (
	// Basic action interface
	_ action.Action = &UpdateDevicePropertiesAction{}

	// Allows the action to be configured with the provider client
	_ action.ActionWithConfigure = &UpdateDevicePropertiesAction{}
)

func NewUpdateDevicePropertiesAction() action.Action {
	return &UpdateDevicePropertiesAction{
		ReadPermissions: []string{
			"DeviceManagementServiceConfig.ReadWrite.All",
		},
		WritePermissions: []string{
			"DeviceManagementServiceConfig.ReadWrite.All",
		},
	}
}

type UpdateDevicePropertiesAction struct {
	client           *msgraphbetasdk.GraphServiceClient
	ProviderTypeName string
	TypeName         string
	ReadPermissions  []string
	WritePermissions []string
}

// Metadata returns the action type name.
func (a *UpdateDevicePropertiesAction) Metadata(ctx context.Context, req action.MetadataRequest, resp *action.MetadataResponse) {
	a.ProviderTypeName = req.ProviderTypeName
	a.TypeName = ActionName
	resp.TypeName = a.FullTypeName()
}

// FullTypeName returns the full action type name in the format "providername_actionname".
func (a *UpdateDevicePropertiesAction) FullTypeName() string {
	return a.ProviderTypeName + "_" + a.TypeName
}

// Description returns the action description.
func (a *UpdateDevicePropertiesAction) Description(ctx context.Context) string {
	return "Updates properties on an Autopilot device in Microsoft Intune using the " +
		"`/deviceManagement/windowsAutopilotDeviceIdentities/{windowsAutopilotDeviceIdentityId}/updateDeviceProperties` endpoint. " +
		"This action allows updating various properties of Autopilot devices including user assignment, group tag, and display name."
}

// Configure sets the client for the action.
func (a *UpdateDevicePropertiesAction) Configure(ctx context.Context, req action.ConfigureRequest, resp *action.ConfigureResponse) {
	a.client = client.SetGraphBetaClientForAction(ctx, req, resp, constants.PROVIDER_NAME+"_"+ActionName)
}

// Schema returns the schema for the action.
func (a *UpdateDevicePropertiesAction) Schema(ctx context.Context, req action.SchemaRequest, resp *action.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Updates properties on an Autopilot device in Microsoft Intune using the " +
			"`/deviceManagement/windowsAutopilotDeviceIdentities/{windowsAutopilotDeviceIdentityId}/updateDeviceProperties` endpoint. " +
			"This action allows updating various properties of Autopilot devices including user assignment, group tag, and display name.",
		Attributes: map[string]schema.Attribute{
			"windows_autopilot_device_identity_id": schema.StringAttribute{
				Required: true,
				MarkdownDescription: "The unique identifier of the Windows Autopilot device identity to update properties for. " +
					"This is the ID of the Windows Autopilot device in Microsoft Intune.",
			},
			"user_principal_name": schema.StringAttribute{
				Optional: true,
				MarkdownDescription: "The user principal name (UPN) of the user to assign to the device. " +
					"This is typically the user's email address in the format user@domain.com. " +
					"If not provided, the user assignment will not be updated.",
			},
			"addressable_user_name": schema.StringAttribute{
				Optional: true,
				MarkdownDescription: "The addressable user name for the user being assigned to the device. " +
					"This is the display name or friendly name of the user. " +
					"If not provided, the addressable user name will not be updated.",
			},
			"group_tag": schema.StringAttribute{
				Optional: true,
				MarkdownDescription: "The group tag to assign to the device. " +
					"Group tags are used for organizing and categorizing devices. " +
					"If not provided, the group tag will not be updated.",
			},
			"display_name": schema.StringAttribute{
				Optional: true,
				MarkdownDescription: "The display name for the device. " +
					"This is a human-readable name for the device. " +
					"If not provided, the display name will not be updated.",
			},
			"timeouts": commonschema.Timeouts(ctx),
		},
	}
}
