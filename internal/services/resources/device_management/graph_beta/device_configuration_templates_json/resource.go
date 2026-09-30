package graphBetaDeviceConfigurationTemplatesJson

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/identityschema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	planmodifiers "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/plan_modifiers"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	assignmentschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema/graph_beta/device_management"
	customvalidator "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/validate/attribute"
)

const (
	ResourceName  = "microsoft365_graph_beta_device_management_device_configuration_templates_json"
	CreateTimeout = 180
	ReadTimeout   = 180
	UpdateTimeout = 180
	DeleteTimeout = 180
)

var (
	// Basic resource interface (CRUD operations)
	_ resource.Resource = &DeviceConfigurationTemplatesJsonResource{}

	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &DeviceConfigurationTemplatesJsonResource{}

	// Enables import functionality
	_ resource.ResourceWithImportState = &DeviceConfigurationTemplatesJsonResource{}

	// Enables plan modification/diff suppression
	_ resource.ResourceWithIdentity = &DeviceConfigurationTemplatesJsonResource{}

	// Enables identity schema for list resource support
	_ resource.ResourceWithValidateConfig = &DeviceConfigurationTemplatesJsonResource{}
)

func NewDeviceConfigurationTemplatesJsonResource() resource.Resource {
	return &DeviceConfigurationTemplatesJsonResource{
		ReadPermissions: []string{
			"DeviceManagementConfiguration.Read.All",
		},
		WritePermissions: []string{
			"DeviceManagementConfiguration.ReadWrite.All",
		},
		ResourcePath: "/deviceManagement/deviceConfigurations",
	}
}

type DeviceConfigurationTemplatesJsonResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

// Metadata returns the resource type name.
func (r *DeviceConfigurationTemplatesJsonResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = ResourceName
}

// Configure sets the client for the resource.
func (r *DeviceConfigurationTemplatesJsonResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, ResourceName)
}

// ImportState imports the resource state.
func (r *DeviceConfigurationTemplatesJsonResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// IdentitySchema defines the identity schema for this resource, used by list operations to uniquely identify instances
func (r *DeviceConfigurationTemplatesJsonResource) IdentitySchema(ctx context.Context, req resource.IdentitySchemaRequest, resp *resource.IdentitySchemaResponse) {
	resp.IdentitySchema = identityschema.Schema{
		Attributes: map[string]identityschema.Attribute{
			"id": identityschema.StringAttribute{
				RequiredForImport: true,
			},
		},
	}
}

func (r *DeviceConfigurationTemplatesJsonResource) Schema(
	ctx context.Context,
	_ resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	assignments := assignmentschema.DeviceConfigurationWithAllGroupAssignmentsAndFilterSchema()
	assignments.MarkdownDescription = "Assignments owned by this profile. Removing all entries unassigns the profile. Do not manage its assignments with another resource."
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages cross-platform Microsoft Intune device configuration templates through `/deviceManagement/deviceConfigurations`. " +
			"Profile metadata uses standard Terraform attributes. Template-specific settings, including the root `@odata.type`, use JSON. Assignments are managed separately in this resource. " +
			"This resource is separate from Settings Catalog (`/configurationPolicies`). Do not manage the same profile with both a typed resource and this JSON resource. " +
			"See [device configurations](https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-deviceconfiguration?view=graph-rest-beta).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
				MarkdownDescription: "The Intune device configuration ID.",
			},
			"display_name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Display name of the device configuration profile.",
				Validators:          []validator.String{stringvalidator.LengthAtLeast(1)},
			},
			"description": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Description of the device configuration profile. Maximum length is 1500 characters.",
				Validators: []validator.String{
					stringvalidator.LengthAtMost(1500),
				},
			},
			"role_scope_tag_ids": schema.SetAttribute{
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
				MarkdownDescription: "Set of Intune scope tag IDs. Defaults to the default scope tag, `0`.",
				Default: setdefault.StaticValue(
					types.SetValueMust(types.StringType, []attr.Value{types.StringValue("0")}),
				),
			},
			"settings": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Complete writable template settings as a JSON object, including the root `@odata.type`. Use `jsonencode()`, a JSON heredoc, or `file(\"profile.json\")`. Configure metadata through `display_name`, `description`, and `role_scope_tag_ids`; configure assignments through `assignments`. Exclude read-only response metadata (`id`, timestamps, `version`, `supportsScopeTags`) and OData response annotations. This endpoint uses a flat profile object. Unlike Settings Catalog, there is no nested `settings` collection. Include the complete writable settings returned by Graph, including default, null, and empty values. The provider reads the full remote settings; it does not project the response onto previously configured keys. Updates use PATCH, so use explicit API-supported reset values to clear settings. To change the root `@odata.type`, recreate the profile with Terraform’s `-replace` option. Nested OData discriminators and relationship bindings remain JSON. Use full Graph URLs for `@odata.bind` references; omit a binding to remove it. The provider reads certificate relationships separately and updates them through `$ref`. Windows encrypted OMA values are recovered through the plaintext endpoint; `omaSettingStringXml.value` accepts cleartext XML, which the provider Base64-encodes only when sending requests. Binary `omaSettingBase64` values and certificate content remain Base64-encoded. Wi-Fi pre-shared keys cannot be recovered on import and must be supplied again. Windows Update derived pause/rollback fields and OMA encryption metadata are read-only and must be excluded. Secret values are stored in Terraform state. The API determines which device configuration template types are supported; see the examples for tested iOS/iPadOS, macOS, Windows, and Android profiles.",
				Validators:          []validator.String{customvalidator.JSONSchemaValidator()},
			},
			"assignments": assignments,
			"timeouts":    commonschema.ResourceTimeouts(ctx),
		},
	}
}
