package graphBetaReuseablePolicySettings

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	planmodifiers "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/plan_modifiers"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	customValidator "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/validators"
	sharedValidators "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/validators/graph_beta/device_management"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	ResourceName  = "graph_beta_device_management_reuseable_policy_setting"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	// Basic resource interface (CRUD operations)
	_ resource.Resource = &ReuseablePolicySettingsResource{}

	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &ReuseablePolicySettingsResource{}

	// Enables import functionality
	_ resource.ResourceWithImportState = &ReuseablePolicySettingsResource{}

	// Enables plan modification/diff suppression
	_ resource.ResourceWithModifyPlan = &ReuseablePolicySettingsResource{}
)

func NewReuseablePolicySettingsResource() resource.Resource {
	return &ReuseablePolicySettingsResource{
		ReadPermissions: []string{
			"DeviceManagementConfiguration.Read.All",
		},
		WritePermissions: []string{
			"DeviceManagementConfiguration.ReadWrite.All",
		},
		ResourcePath: "/deviceManagement/reusablePolicySettings",
	}
}

type ReuseablePolicySettingsResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ProviderTypeName string
	TypeName         string
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

// Metadata returns the resource type name.
func (r *ReuseablePolicySettingsResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	r.ProviderTypeName = req.ProviderTypeName
	r.TypeName = ResourceName
	resp.TypeName = r.FullTypeName()
}

// FullTypeName returns the full type name of the resource for logging purposes.
func (r *ReuseablePolicySettingsResource) FullTypeName() string {
	return r.ProviderTypeName + "_" + r.TypeName
}

// Configure sets the client for the resource.
func (r *ReuseablePolicySettingsResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, constants.PROVIDER_NAME+"_"+ResourceName)
}

// ImportState imports the resource state.
func (r *ReuseablePolicySettingsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Function to create the full device management win32 lob app schema
func (r *ReuseablePolicySettingsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages reusable policy settings using the `/deviceManagement/reusablePolicySettings` endpoint. Reusable settings enable centralized management of configuration templates that can be referenced across multiple policies, reducing duplication and ensuring consistent settings application across Endpoint Privilege Management and other device configuration policies.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
				MarkdownDescription: "The unique identifier for this Reuseable Settings Policy",
			},
			"display_name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The reusable setting display name supplied by user.",
			},
			"description": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				PlanModifiers:       []planmodifier.String{planmodifiers.DefaultValueString("")},
				MarkdownDescription: "Reuseable Settings Policy description",
			},
			"settings": schema.StringAttribute{
				Required: true,
				MarkdownDescription: "Reuseable Settings Policy with settings catalog settings defined as a valid JSON string. Please provide a valid JSON-encoded settings structure. " +
					"This can either be extracted from an existing policy using the Intune gui `export JSON` functionality if supported, via a script such as this powershell script." +
					" [ExportReuseablePolicySettingsById](https://github.com/deploymenttheory/terraform-provider-microsoft365/blob/main/scripts/Export-IntuneReuseablePolicySettingsById.ps1) " +
					"or created from scratch. The JSON structure should match the graph schema of the settings catalog. Please look at the " +
					"terraform documentation for the settings catalog template for examples and how to correctly format the HCL.\n\n" +
					"A correctly formatted field in the HCL should begin and end like this:\n" +
					"```hcl\n" +
					"settings = jsonencode({\n" +
					"  \"settings\": [\n" +
					"    {\n" +
					"      \"id\": \"0\",\n" +
					"      \"settingInstance\": {\n" +
					"      }\n" +
					"    }\n" +
					"  ]\n" +
					"})\n" +
					"```\n\n" +
					"**Note:** Settings must always be provided as an array within the settings field, even when configuring a single setting." +
					"This is required because the Microsoft Graph SDK for Go always returns settings in an array format\n\n" +
					"**Note:** When configuring secret values (identified by @odata.type: \"#microsoft.graph.deviceManagementConfigurationSecretSettingValue\") " +
					"ensure the valueState is set to \"notEncrypted\". The value \"encryptedValueToken\" is reserved for server" +
					"responses and should not be used when creating or updating settings.\n\n" +
					"```hcl\n" +
					"settings = jsonencode({\n" +
					"  \"settings\": [\n" +
					"    {\n" +
					"      \"id\": \"0\",\n" +
					"      \"settingInstance\": {\n" +
					"        \"@odata.type\": \"#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance\",\n" +
					"        \"settingDefinitionId\": \"com.apple.loginwindow_autologinpassword\",\n" +
					"        \"settingInstanceTemplateReference\": null,\n" +
					"        \"simpleSettingValue\": {\n" +
					"          \"@odata.type\": \"#microsoft.graph.deviceManagementConfigurationSecretSettingValue\",\n" +
					"          \"valueState\": \"notEncrypted\",\n" +
					"          \"value\": \"your_secret_value\",\n" +
					"          \"settingValueTemplateReference\": null\n" +
					"        }\n" +
					"      }\n" +
					"    }\n" +
					"  ]\n" +
					"})\n" +
					"```\n\n",
				Validators: []validator.String{
					customValidator.JSONSchemaValidator(),
					sharedValidators.SettingsCatalogValidator(),
				},
				PlanModifiers: []planmodifier.String{
					planmodifiers.NormalizeJSONPlanModifier{},
				},
			},
			"created_date_time": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Creation date and time of the settings catalog policy",
			},
			"last_modified_date_time": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Last modification date and time of the settings catalog policy",
			},
			"version": schema.Int32Attribute{
				Computed:            true,
				MarkdownDescription: "Version of the policy",
			},
			"referencing_configuration_policies": schema.ListAttribute{
				ElementType: types.StringType,
				Computed:    true,
				PlanModifiers: []planmodifier.List{
					planmodifiers.UseStateForUnknownList(),
				},
				MarkdownDescription: "List of configuration policies referencing this reuseable policy",
			},
			"referencing_configuration_policy_count": schema.Int32Attribute{
				Computed:            true,
				MarkdownDescription: "Number of configuration policies referencing this reuseable policy",
			},
			"timeouts": commonschema.Timeouts(ctx),
		},
	}
}
