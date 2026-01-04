package graphBetaApplyCloudPcProvisioningPolicy

import (
	"context"
	"regexp"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/action/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	ActionName    = "microsoft365_graph_beta_device_management_windows_365_apply_cloud_pc_provisioning_policy"
	InvokeTimeout = 60
)

var (
	_ action.Action                   = &ApplyCloudPcProvisioningPolicyAction{}
	_ action.ActionWithConfigure      = &ApplyCloudPcProvisioningPolicyAction{}
	_ action.ActionWithValidateConfig = &ApplyCloudPcProvisioningPolicyAction{}
)

func NewApplyCloudPcProvisioningPolicyAction() action.Action {
	return &ApplyCloudPcProvisioningPolicyAction{
		ReadPermissions: []string{
			"CloudPC.ReadWrite.All",
		},
		WritePermissions: []string{
			"CloudPC.ReadWrite.All",
		},
	}
}

type ApplyCloudPcProvisioningPolicyAction struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
}

func (a *ApplyCloudPcProvisioningPolicyAction) Metadata(ctx context.Context, req action.MetadataRequest, resp *action.MetadataResponse) {
	resp.TypeName = ActionName
}

func (a *ApplyCloudPcProvisioningPolicyAction) Configure(ctx context.Context, req action.ConfigureRequest, resp *action.ConfigureResponse) {
	a.client = client.SetGraphBetaClientForAction(ctx, req, resp, ActionName)
}

func (a *ApplyCloudPcProvisioningPolicyAction) Schema(ctx context.Context, req action.SchemaRequest, resp *action.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Applies configuration settings to a Cloud PC provisioning policy using the " +
			"`/deviceManagement/virtualEndpoint/provisioningPolicies/{id}/apply` endpoint. " +
			"This action allows applying region or single sign-on settings to existing Cloud PCs that were provisioned with this policy. " +
			"When you change the network, image, region, or single sign-on configuration in a provisioning policy, " +
			"the changes only apply to newly provisioned or reprovisioned Cloud PCs. Use this action to apply " +
			"region or single sign-on changes to previously provisioned Cloud PCs.",
		Attributes: map[string]schema.Attribute{
			"provisioning_policy_id": schema.StringAttribute{
				Required: true,
				MarkdownDescription: "The unique identifier of the Cloud PC provisioning policy to apply settings to. " +
					"This is the ID of the provisioning policy in Microsoft Intune.",
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(constants.GuidRegex),
						"must be a valid GUID format (e.g., 12345678-1234-1234-1234-123456789abc)",
					),
				},
			},
			"policy_settings": schema.StringAttribute{
				Optional: true,
				MarkdownDescription: "The target property of the apply action. " +
					"Valid values are:\n" +
					"- `region`: Apply region configuration changes to existing Cloud PCs (default)\n" +
					"- `singleSignOn`: Apply single sign-on configuration changes to existing Cloud PCs\n\n" +
					"The default value is `region`. This action applies region as a value if this parameter is null.\n" +
					"Note: Network and image changes cannot be applied retrospectively and require reprovisioning.",
				Validators: []validator.String{
					stringvalidator.OneOf("singleSignOn", "region"),
				},
			},
			"reserve_percentage": schema.Int32Attribute{
				Optional: true,
				MarkdownDescription: "For Frontline shared Cloud PCs only. The percentage of Cloud PCs to keep available. " +
					"Administrators can set this property to a value from 0 to 99. " +
					"Cloud PCs are reprovisioned only when there are no active and connected Cloud PC users. " +
					"This parameter is only applicable for Frontline shared provisioning policies.",
				Validators: []validator.Int32{
					int32validator.Between(0, 99),
				},
			},
			"validate_policy_exists": schema.BoolAttribute{
				Optional: true,
				MarkdownDescription: "When set to `true` (default), the action will validate that the provisioning policy exists " +
					"and check the provisioning type if reserve_percentage is specified before attempting to apply settings. " +
					"When `false`, policy validation is skipped and the action will attempt to apply settings directly. " +
					"Disabling validation can improve performance but may result in errors if the policy doesn't exist or is the wrong type.",
			},
			"timeouts": commonschema.ActionTimeouts(ctx),
		},
	}
}
