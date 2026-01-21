// REF: https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-activatedeviceesim?view=graph-rest-beta
package graphBetaActivateDeviceEsimManagedDevice

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
	ActionName    = "microsoft365_graph_beta_device_management_managed_device_activate_device_esim"
	InvokeTimeout = 120
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ action.Action                   = &ActivateDeviceEsimManagedDeviceAction{}
	_ action.ActionWithConfigure      = &ActivateDeviceEsimManagedDeviceAction{}
	_ action.ActionWithValidateConfig = &ActivateDeviceEsimManagedDeviceAction{}
)

func NewActivateDeviceEsimManagedDeviceAction() action.Action {
	return &ActivateDeviceEsimManagedDeviceAction{}
}

type ActivateDeviceEsimManagedDeviceAction struct {
	client *msgraphbetasdk.GraphServiceClient
}

// Metadata returns the action type name.
func (a *ActivateDeviceEsimManagedDeviceAction) Metadata(ctx context.Context, req action.MetadataRequest, resp *action.MetadataResponse) {
	resp.TypeName = ActionName
}

// Configure adds the provider configured client to the action.
func (a *ActivateDeviceEsimManagedDeviceAction) Configure(ctx context.Context, req action.ConfigureRequest, resp *action.ConfigureResponse) {
	a.client = client.SetGraphBetaClientForAction(ctx, req, resp, ActionName)
}

// Schema defines the schema for the action.
func (a *ActivateDeviceEsimManagedDeviceAction) Schema(ctx context.Context, req action.SchemaRequest, resp *action.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Activates eSIM cellular data plans on iOS and iPadOS devices in Microsoft Intune using the " +
			"`/deviceManagement/managedDevices/{managedDeviceId}/activateDeviceEsim` and " +
			"`/deviceManagement/comanagedDevices/{managedDeviceId}/activateDeviceEsim` endpoints. " +
			"This action is used to remotely activate eSIM cellular plans without physical SIM cards, making it easier to manage connectivity for users.",
		Attributes: map[string]schema.Attribute{
			"managed_devices": schema.ListNestedAttribute{
				Optional: true,
				MarkdownDescription: "List of iOS/iPadOS managed devices to activate eSIM on. These are devices fully managed by Intune only. " +
					"Devices must have eSIM hardware capability (iPhone XS and later, cellular iPads with eSIM support).",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"device_id": schema.StringAttribute{
							Required: true,
							MarkdownDescription: "The unique identifier (GUID) of the iOS/iPadOS device to activate eSIM on. " +
								"Device must have eSIM hardware capability (iPhone XS+, cellular iPad with eSIM). " +
								"Example: `\"12345678-1234-1234-1234-123456789abc\"`",
							Validators: []validator.String{
								stringvalidator.RegexMatches(
									regexp.MustCompile(constants.GuidRegex),
									"device_id must be a valid GUID format (e.g., 12345678-1234-1234-1234-123456789abc)",
								),
							},
						},
						"carrier_url": schema.StringAttribute{
							Required: true,
							MarkdownDescription: "The activation server URL provided by your mobile carrier for eSIM activation. " +
								"This URL is carrier-specific and contains the activation profile. " +
								"Example: `\"https://carrier.example.com/esim/activate?token=abc123\"`",
							Validators: []validator.String{
								stringvalidator.RegexMatches(
									regexp.MustCompile(constants.HttpOrHttpsUrlRegex),
									"must be a valid HTTP or HTTPS URL (e.g., https://carrier.example.com/esim/activate)",
								),
							},
						},
					},
				},
			},
			"comanaged_devices": schema.ListNestedAttribute{
				Optional: true,
				MarkdownDescription: "List of iOS/iPadOS co-managed devices to activate eSIM on. These are devices managed by both Intune and " +
					"Configuration Manager (SCCM). Devices must have eSIM hardware capability (iPhone XS and later, cellular iPads with eSIM support).",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"device_id": schema.StringAttribute{
							Required: true,
							MarkdownDescription: "The unique identifier (GUID) of the iOS/iPadOS co-managed device to activate eSIM on. " +
								"Device must have eSIM hardware capability (iPhone XS+, cellular iPad with eSIM). " +
								"Example: `\"12345678-1234-1234-1234-123456789abc\"`",
							Validators: []validator.String{
								stringvalidator.RegexMatches(
									regexp.MustCompile(constants.GuidRegex),
									"device_id must be a valid GUID format (e.g., 12345678-1234-1234-1234-123456789abc)",
								),
							},
						},
						"carrier_url": schema.StringAttribute{
							Required: true,
							MarkdownDescription: "The activation server URL provided by your mobile carrier for eSIM activation. " +
								"This URL is carrier-specific and contains the activation profile. " +
								"Example: `\"https://carrier.example.com/esim/activate?code=xyz789\"`",
							Validators: []validator.String{
								stringvalidator.RegexMatches(
									regexp.MustCompile(constants.HttpOrHttpsUrlRegex),
									"must be a valid HTTP or HTTPS URL (e.g., https://carrier.example.com/esim/activate)",
								),
							},
						},
					},
				},
			},
			"ignore_partial_failures": schema.BoolAttribute{
				Optional: true,
				MarkdownDescription: "If set to `true`, the action will succeed even if some devices fail eSIM activation. " +
					"Failed devices will be reported as warnings instead of errors. " +
					"Default: `false` (action fails if any device fails).",
			},
			"validate_device_exists": schema.BoolAttribute{
				Optional: true,
				MarkdownDescription: "Whether to validate that devices exist before attempting activation. " +
					"Disabling this can speed up planning but may result in runtime errors for non-existent devices. " +
					"Default: `true`.",
			},
			"timeouts": commonschema.ActionTimeouts(ctx),
		},
	}
}
