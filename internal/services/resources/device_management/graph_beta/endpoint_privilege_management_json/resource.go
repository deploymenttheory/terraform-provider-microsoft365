package graphBetaEndpointPrivilegeManagement

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	planmodifiers "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/plan_modifiers"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	commonschemagraphbeta "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema/graph_beta/device_management"
	customValidator "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/validate/attribute"
	sharedValidators "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/validate/graph_beta/device_management"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	ResourceName  = "microsoft365_graph_beta_device_management_endpoint_privilege_management_json"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	// Basic resource interface (CRUD operations)
	_ resource.Resource = &EndpointPrivilegeManagementResource{}

	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &EndpointPrivilegeManagementResource{}

	// Enables import functionality
	_ resource.ResourceWithImportState = &EndpointPrivilegeManagementResource{}

	// Enables plan modification/diff suppression
	_ resource.ResourceWithModifyPlan = &EndpointPrivilegeManagementResource{}
)

func NewEndpointPrivilegeManagementResource() resource.Resource {
	return &EndpointPrivilegeManagementResource{
		ReadPermissions: []string{
			"DeviceManagementConfiguration.Read.All",
		},
		WritePermissions: []string{
			"DeviceManagementConfiguration.ReadWrite.All",
		},
		ResourcePath: "/deviceManagement/configurationPolicies",
	}
}

type EndpointPrivilegeManagementResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

// Metadata returns the resource type name.
func (r *EndpointPrivilegeManagementResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = ResourceName
}

// Configure sets the client for the resource.
func (r *EndpointPrivilegeManagementResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, ResourceName)
}

// ImportState imports the resource state.
func (r *EndpointPrivilegeManagementResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Function to create the full device management win32 lob app schema
func (r *EndpointPrivilegeManagementResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages Endpoint Privilege Management policies using the `/deviceManagement/configurationPolicies` endpoint. " +
			"EPM policies control elevation settings and rules for Windows devices, allowing administrators to grant temporary administrative " +
			"privileges to standard users for specific applications or processes without compromising overall security posture.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
				MarkdownDescription: "The unique identifier for this Endpoint Privilege Management Policy",
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Policy name",
			},
			"description": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
				MarkdownDescription: "Optional description of the resource. Maximum length is 1500 characters.",
				Validators: []validator.String{
					stringvalidator.LengthAtMost(1500),
				},
			},
			"settings_catalog_template_type": schema.StringAttribute{
				Required: true,
				MarkdownDescription: "Defines which Endpoint Privilege Management Policy type with settings catalog setting will be deployed. " +
					"Options available are `elevation_settings_policy` or `elevation_rules_policy`.",
				Validators: []validator.String{
					stringvalidator.OneOf(
						"elevation_settings_policy",
						"elevation_rules_policy",
					),
				},
			},
			"settings": schema.StringAttribute{
				Required: true,
				MarkdownDescription: "Endpoint Privilege Management Policy settings defined as a JSON string. Please provide a valid JSON-encoded settings structure. " +
					"This can either be extracted from an existing policy using the Intune gui `export JSON` functionality if supported, via a script such as this powershell script." +
					" [Export-IntuneSettingsCatalogConfigurationById](https://github.com/deploymenttheory/terraform-provider-microsoft365/blob/main/scripts/device_management/Export-IntuneSettingsCatalogConfigurationById.ps1) " +
					"or created from scratch. The JSON structure should match the graph schema of the settings catalog. Please look at the " +
					"terraform documentation for the Endpoint Privilege Management Policy for examples and how to correctly format the HCL.\n\n" +
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
					"```\n\n" +
					"Alternatively, if you prefer a fully hcl based approach, you can use the resource `microsoft365_graph_beta_device_management_settings_catalog_configuration_policy` " +
					"by creating the epm policy via thr gui and then exporting the settings catalog configuration to hcl using the [Export-IntuneSettingsCatalogConfigurationToHCL.ps1] (https://github.com/deploymenttheory/terraform-provider-microsoft365/blob/main/scripts/device_management/Export-IntuneSettingsCatalogConfigurationToHCL.ps1) script." +
					"This will create a fully hcl based epm policy with all the settings catalog settings configured.",
				Validators: []validator.String{
					customValidator.JSONSchemaValidator(),
					sharedValidators.SettingsCatalogJSONValidator(),
				},
				PlanModifiers: []planmodifier.String{
					planmodifiers.NormalizeJSONPlanModifier{},
				},
			},
			"platforms": schema.StringAttribute{
				Computed: true,
				MarkdownDescription: "Platform type for this Endpoint Privilege Management Policy." +
					"Will always be set to `windows10`, as EPM currently only supports windows device types." +
					"Defaults to windows10.",
				Validators: []validator.String{
					stringvalidator.OneOf(
						"windows10",
					),
				},
				PlanModifiers: []planmodifier.String{planmodifiers.DefaultValueString("windows10")},
			},
			"technologies": schema.ListAttribute{
				ElementType: types.StringType,
				Computed:    true,
				MarkdownDescription: "Describes a list of technologies this Endpoint Privilege Management Policy with settings catalog setting will be deployed with." +
					"Defaults to `mdm`, `endpointPrivilegeManagement`.",
				Validators: []validator.List{
					customValidator.StringListAllowedValues(
						"mdm", "endpointPrivilegeManagement",
					),
				},
			},
			"role_scope_tag_ids": schema.SetAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Set of scope tag IDs for this Settings Catalog template profile.",
				PlanModifiers: []planmodifier.Set{
					planmodifiers.DefaultSetValue(
						[]attr.Value{types.StringValue("0")},
					),
				},
			},
			"created_date_time": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				MarkdownDescription: "Creation date and time of the settings catalog policy",
			},
			"last_modified_date_time": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Last modification date and time of the settings catalog policy",
			},
			"settings_count": schema.Int32Attribute{
				Computed:            true,
				MarkdownDescription: "Number of settings catalog settings with the policy. This will change over time as the resource is updated.",
			},
			"is_assigned": schema.BoolAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.Bool{
					planmodifiers.UseStateForUnknownBool(),
				},
				MarkdownDescription: "Indicates if the policy is assigned to any scope",
			},
			"assignments": commonschemagraphbeta.DeviceConfigurationWithAllGroupAssignmentsAndFilterSchema(),
			"timeouts":    commonschema.ResourceTimeouts(ctx),
		},
	}
}
