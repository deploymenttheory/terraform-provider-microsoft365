package graphBetaDeleteUserFromSharedAppleDevice

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
	ActionName = "microsoft365_graph_beta_device_management_managed_device_delete_user_from_shared_apple_device"
	InvokeTimeout = 60
)

var (
	_ action.Action                   = &DeleteUserFromSharedAppleDeviceAction{}
	_ action.ActionWithConfigure      = &DeleteUserFromSharedAppleDeviceAction{}
	_ action.ActionWithValidateConfig = &DeleteUserFromSharedAppleDeviceAction{}
)

func NewDeleteUserFromSharedAppleDeviceAction() action.Action {
	return &DeleteUserFromSharedAppleDeviceAction{
		ReadPermissions: []string{
			"DeviceManagementManagedDevices.PrivilegedOperations.All",
		},
		WritePermissions: []string{
			"DeviceManagementManagedDevices.PrivilegedOperations.All",
		},
	}
}

type DeleteUserFromSharedAppleDeviceAction struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
}

func (a *DeleteUserFromSharedAppleDeviceAction) Metadata(ctx context.Context, req action.MetadataRequest, resp *action.MetadataResponse) {
	resp.TypeName = ActionName
}

func (a *DeleteUserFromSharedAppleDeviceAction) Configure(ctx context.Context, req action.ConfigureRequest, resp *action.ConfigureResponse) {
	a.client = client.SetGraphBetaClientForAction(ctx, req, resp, ActionName)
}

func (a *DeleteUserFromSharedAppleDeviceAction) Schema(ctx context.Context, req action.SchemaRequest, resp *action.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Deletes a user and their cached data from Shared iPad devices using the " +
			"`/deviceManagement/managedDevices/{managedDeviceId}/deleteUserFromSharedAppleDevice` and " +
			"`/deviceManagement/comanagedDevices/{managedDeviceId}/deleteUserFromSharedAppleDevice` endpoints. " +
			"This action permanently removes the specified user's account and all associated cached data from " +
			"the Shared iPad, freeing up storage space for other users.\n\n" +
			"**What This Action Does:**\n" +
			"- Permanently deletes user from Shared iPad device roster\n" +
			"- Removes all cached user data (documents, photos, app data)\n" +
			"- Frees up device storage space\n" +
			"- Prevents user from logging back into that specific device\n" +
			"- Does not affect user's account or data in the cloud\n" +
			"- Cannot be undone (user must be re-added if needed)\n\n" +
			"**Platform Support:**\n" +
			"- **iPadOS**: Full support (Shared iPad mode only)\n" +
			"- **iOS**: Not supported (iPhones don't support Shared mode)\n" +
			"- **Other platforms**: Not supported\n\n" +
			"**Reference:** [Microsoft Graph API - Delete User From Shared Apple Device](https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-deleteuserfromsharedappledevice?view=graph-rest-beta)",
		Attributes: map[string]schema.Attribute{
			"managed_devices": schema.ListNestedAttribute{
				Optional:            true,
				MarkdownDescription: "List of managed device-user pairs. Managed devices are fully managed by Intune only.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"device_id": schema.StringAttribute{
							Required:            true,
							MarkdownDescription: "The managed device ID (GUID) of the Shared iPad.",
							Validators: []validator.String{
								stringvalidator.RegexMatches(
									regexp.MustCompile(constants.GuidRegex),
									"device_id must be a valid GUID format",
								),
							},
						},
						"user_principal_name": schema.StringAttribute{
							Required:            true,
							MarkdownDescription: "The user principal name (UPN) to delete from the device.",
						},
					},
				},
			},
			"comanaged_devices": schema.ListNestedAttribute{
				Optional:            true,
				MarkdownDescription: "List of co-managed device-user pairs. Co-managed devices are managed by both Intune and Configuration Manager (SCCM).",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"device_id": schema.StringAttribute{
							Required:            true,
							MarkdownDescription: "The co-managed device ID (GUID) of the Shared iPad.",
							Validators: []validator.String{
								stringvalidator.RegexMatches(
									regexp.MustCompile(constants.GuidRegex),
									"device_id must be a valid GUID format",
								),
							},
						},
						"user_principal_name": schema.StringAttribute{
							Required:            true,
							MarkdownDescription: "The user principal name (UPN) to delete from the device.",
						},
					},
				},
			},
			"timeouts": commonschema.ActionTimeouts(ctx),
		},
	}
}
