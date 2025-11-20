package graphBetaUpdateDeviceProperties

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
	ActionName = "microsoft365_graph_beta_device_management_windows_autopilot_device_identity_update_device_properties"
)

var (
	// Basic action interface
	_ action.Action = &UpdateDevicePropertiesAction{}

	// Allows the action to be configured with the provider client
	_ action.ActionWithConfigure = &UpdateDevicePropertiesAction{}

	// Allows the action to validate configuration
	_ action.ActionWithValidateConfig = &UpdateDevicePropertiesAction{}
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
	ReadPermissions  []string
	WritePermissions []string
}

// Metadata returns the action type name.
func (a *UpdateDevicePropertiesAction) Metadata(ctx context.Context, req action.MetadataRequest, resp *action.MetadataResponse) {
	resp.TypeName = ActionName
}

// Configure sets the client for the action.
func (a *UpdateDevicePropertiesAction) Configure(ctx context.Context, req action.ConfigureRequest, resp *action.ConfigureResponse) {
	a.client = client.SetGraphBetaClientForAction(ctx, req, resp, ActionName)
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
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(constants.GuidRegex),
						"must be a valid GUID format (e.g., 12345678-1234-1234-1234-123456789abc)",
					),
				},
			},
			"user_principal_name": schema.StringAttribute{
				Optional: true,
				MarkdownDescription: "The user principal name (UPN) of the user to assign to the device. " +
					"This is typically the user's email address in the format user@domain.com. " +
					"If not provided, the user assignment will not be updated.",
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(constants.EmailRegex),
						"must be a valid email address format (e.g., user@domain.com)",
					),
				},
			},
			"addressable_user_name": schema.StringAttribute{
				Optional: true,
				MarkdownDescription: "The addressable user name for the user being assigned to the device. " +
					"This is the display name or friendly name of the user. " +
					"If not provided, the addressable user name will not be updated.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"group_tag": schema.StringAttribute{
				Optional: true,
				MarkdownDescription: "The group tag to assign to the device. " +
					"Group tags are used for organizing and categorizing devices. " +
					"If not provided, the group tag will not be updated.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"display_name": schema.StringAttribute{
				Optional: true,
				MarkdownDescription: "The display name for the device. " +
					"This is a human-readable name for the device. " +
					"If not provided, the display name will not be updated.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"timeouts": commonschema.Timeouts(ctx),
		},
	}
}
