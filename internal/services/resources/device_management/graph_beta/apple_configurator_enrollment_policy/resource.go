package graphBetaAppleConfiguratorEnrollmentPolicy

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	planmodifiers "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/plan_modifiers"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/validators"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	ResourceName  = "graph_beta_device_management_apple_configurator_enrollment_policy"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	// Basic resource interface (CRUD operations)
	_ resource.Resource = &AppleConfiguratorEnrollmentPolicyResource{}

	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &AppleConfiguratorEnrollmentPolicyResource{}

	// Enables import functionality
	_ resource.ResourceWithImportState = &AppleConfiguratorEnrollmentPolicyResource{}

	// Enables plan modification/diff suppression
	_ resource.ResourceWithModifyPlan = &AppleConfiguratorEnrollmentPolicyResource{}
)

func NewAppleConfiguratorEnrollmentPolicyResource() resource.Resource {
	return &AppleConfiguratorEnrollmentPolicyResource{
		ReadPermissions: []string{
			"DeviceManagementServiceConfig.Read.All",
		},
		WritePermissions: []string{
			"DeviceManagementServiceConfig.ReadWrite.All",
		},
		ResourcePath: "deviceManagement/depOnboardingSettings/{depOnboardingSettingsId}/enrollmentProfiles",
	}
}

type AppleConfiguratorEnrollmentPolicyResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ProviderTypeName string
	TypeName         string
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

func (r *AppleConfiguratorEnrollmentPolicyResource) FullTypeName() string {
	return r.ProviderTypeName + "_" + ResourceName
}

func (r *AppleConfiguratorEnrollmentPolicyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	r.ProviderTypeName = req.ProviderTypeName
	r.TypeName = ResourceName
	resp.TypeName = r.FullTypeName()
}

func (r *AppleConfiguratorEnrollmentPolicyResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, constants.PROVIDER_NAME+"_"+ResourceName)
}

func (r *AppleConfiguratorEnrollmentPolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *AppleConfiguratorEnrollmentPolicyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages apple configurator enrollment policies in Microsoft Intune using the `/deviceManagement/depOnboardingSettings/{depOnboardingSettingsId}/enrollmentProfiles` endpoint." +
			"Apple configurator enrollment policy defines the enrollment policy for Apple devices when using the Apple Configurator enrollment type.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
				MarkdownDescription: "The unique identifier of the enrollment profile.",
			},
			"dep_onboarding_settings_id": schema.StringAttribute{
				Computed: true,
				MarkdownDescription: "Identifier of the parent depOnboardingSetting that contains this apple business manager enrollment profile. " +
					"This is resolved from the /deviceManagement endpoint and correlated to the intuneAccountId.",
			},
			"display_name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The name of the apple business manager enrollment profile displayed in Intune.",
			},
			"description": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Description of the apple business manager enrollment profile.",
			},
			"requires_user_authentication": schema.BoolAttribute{
				Required:            true,
				MarkdownDescription: "Indicates whether the user must authenticate during Apple device setup.",
			},
			"enable_authentication_via_company_portal": schema.BoolAttribute{
				Required:            true,
				MarkdownDescription: "How users first sign in to authenticate with Intune.If your organization uses multi-factor authentication, set this to true; the app will then automatically install on devices at time of enrollment. ",
				Validators: []validator.Bool{
					validators.MutuallyExclusiveBool("require_company_portal_on_setup_assistant_enrolled_devices", "enable_authentication_via_company_portal and require_company_portal_on_setup_assistant_enrolled_devices cannot both be set to true"),
				},
			},
			"require_company_portal_on_setup_assistant_enrolled_devices": schema.BoolAttribute{
				Required:            true,
				MarkdownDescription: "How users first sign in to authenticate with Intune. Setup assistant as a standalone authentication method has been superseded by setup assistant with modern authentication. The modern authentication method prompts users to authenticate based on the settings you've configured in Microsoft Entra ID.",
				Validators: []validator.Bool{
					validators.MutuallyExclusiveBool("enable_authentication_via_company_portal", "enable_authentication_via_company_portal and require_company_portal_on_setup_assistant_enrolled_devices cannot both be set to true"),
				},
			},
			"configuration_endpoint_url": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Apple Configurator enrollment configuration endpoint URL generated by Intune.",
			},
			"timeouts": commonschema.Timeouts(ctx),
		},
	}
}
