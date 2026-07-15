package graphBetaWindowsCustomConfiguration

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	planmodifiers "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/plan_modifiers"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	commonschemagraphbeta "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema/graph_beta/device_management"
	customvalidator "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/validate/attribute"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/identityschema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	ResourceName  = "microsoft365_graph_beta_device_management_windows_custom_configuration"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	// Basic resource interface (CRUD operations)
	_ resource.Resource = &WindowsCustomConfigurationResource{}

	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &WindowsCustomConfigurationResource{}

	// Enables import functionality
	_ resource.ResourceWithImportState = &WindowsCustomConfigurationResource{}

	// Enables plan modification/diff suppression
	_ resource.ResourceWithModifyPlan = &WindowsCustomConfigurationResource{}

	// Enables identity schema for list resource support
	_ resource.ResourceWithIdentity = &WindowsCustomConfigurationResource{}
)

func NewWindowsCustomConfigurationResource() resource.Resource {
	return &WindowsCustomConfigurationResource{
		ReadPermissions: []string{
			"DeviceManagementConfiguration.Read.All",
		},
		WritePermissions: []string{
			"DeviceManagementConfiguration.ReadWrite.All",
		},
		ResourcePath: "/deviceManagement/deviceConfigurations",
	}
}

type WindowsCustomConfigurationResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

// Metadata returns the resource type name.
func (r *WindowsCustomConfigurationResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = ResourceName
}

// Configure sets the client for the resource.
func (r *WindowsCustomConfigurationResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, ResourceName)
}

// ImportState imports the resource state.
func (r *WindowsCustomConfigurationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// IdentitySchema defines the identity schema for this resource, used by list operations to uniquely identify instances
func (r *WindowsCustomConfigurationResource) IdentitySchema(ctx context.Context, req resource.IdentitySchemaRequest, resp *resource.IdentitySchemaResponse) {
	resp.IdentitySchema = identityschema.Schema{
		Attributes: map[string]identityschema.Attribute{
			"id": identityschema.StringAttribute{
				RequiredForImport: true,
			},
		},
	}
}

// Schema defines the schema for the resource.
func (r *WindowsCustomConfigurationResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages Windows custom device configuration profiles (OMA-URI settings) in Microsoft Intune using the " +
			"`/deviceManagement/deviceConfigurations` endpoint with the `#microsoft.graph.windows10CustomConfiguration` OData type. " +
			"Custom profiles use OMA-URI (Open Mobile Alliance Uniform Resource Identifier) settings to configure features that " +
			"aren't built into Intune yet, such as ADMX policy ingestion for third party applications.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The unique identifier for the Windows custom configuration profile.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"display_name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The display name for the Windows custom configuration profile.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"description": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Optional description of the resource. Maximum length is 1500 characters.",
				Validators: []validator.String{
					stringvalidator.LengthAtMost(1500),
				},
			},
			"role_scope_tag_ids": schema.SetAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Set of scope tag IDs for this Windows custom configuration profile.",
				PlanModifiers: []planmodifier.Set{
					planmodifiers.DefaultSetValue(
						[]attr.Value{types.StringValue("0")},
					),
				},
			},
			"oma_settings": schema.ListNestedAttribute{
				Required: true,
				MarkdownDescription: "The list of OMA-URI settings deployed by this profile. Each setting targets a single OMA-URI " +
					"node in the Windows CSP (Configuration Service Provider) tree.",
				Validators: []validator.List{
					listvalidator.SizeAtLeast(1),
				},
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"odata_type": schema.StringAttribute{
							Required: true,
							MarkdownDescription: "The OData type of the OMA setting, which determines the data type of `value`. " +
								"Possible values are: `#microsoft.graph.omaSettingString`, `#microsoft.graph.omaSettingInteger`, " +
								"`#microsoft.graph.omaSettingBoolean`, `#microsoft.graph.omaSettingBase64`, " +
								"`#microsoft.graph.omaSettingDateTime`, `#microsoft.graph.omaSettingFloatingPoint`, " +
								"`#microsoft.graph.omaSettingStringXml`.",
							Validators: []validator.String{
								stringvalidator.OneOf(
									"#microsoft.graph.omaSettingString",
									"#microsoft.graph.omaSettingInteger",
									"#microsoft.graph.omaSettingBoolean",
									"#microsoft.graph.omaSettingBase64",
									"#microsoft.graph.omaSettingDateTime",
									"#microsoft.graph.omaSettingFloatingPoint",
									"#microsoft.graph.omaSettingStringXml",
								),
							},
						},
						"display_name": schema.StringAttribute{
							Required:            true,
							MarkdownDescription: "The display name of the OMA setting.",
							Validators: []validator.String{
								stringvalidator.LengthAtLeast(1),
							},
						},
						"description": schema.StringAttribute{
							Optional:            true,
							MarkdownDescription: "Optional description of the OMA setting.",
						},
						"oma_uri": schema.StringAttribute{
							Required: true,
							MarkdownDescription: "The OMA-URI of the setting, e.g. " +
								"`./Device/Vendor/MSFT/Policy/Config/VSCode~Policy~Application/UpdateMode`.",
							Validators: []validator.String{
								stringvalidator.LengthAtLeast(1),
							},
						},
						"value": schema.StringAttribute{
							Required: true,
							MarkdownDescription: "The value of the OMA setting, expressed as a string and converted according to " +
								"`odata_type`: string and stringXml values are used as-is, integer values must be whole numbers " +
								"without leading zeros, boolean values must be `true` or `false` (lowercase), base64 values must " +
								"be base64-encoded content, dateTime values must be UTC RFC3339 timestamps (e.g. " +
								"`2024-01-01T00:00:00Z`), and floatingPoint values must be decimal numbers without trailing " +
								"zeros. Values must be in this canonical form because it is what the Graph API returns on read; " +
								"non-canonical values are rejected at plan time.",
						},
						"file_name": schema.StringAttribute{
							Optional: true,
							MarkdownDescription: "The file name associated with the value, e.g. `policy.xml` or `logo.png`. " +
								"Only applicable when `odata_type` is `#microsoft.graph.omaSettingBase64` or " +
								"`#microsoft.graph.omaSettingStringXml`.",
						},
					},
				},
			},
			"assignments": windowsCustomConfigurationAssignmentsSchema(),
			"timeouts":    commonschema.ResourceTimeouts(ctx),
		},
	}
}

// windowsCustomConfigurationAssignmentsSchema reuses the common device configuration assignment
// schema and adds the conditional group_id validation required before this resource mutates the tenant.
func windowsCustomConfigurationAssignmentsSchema() schema.SetNestedAttribute {
	assignmentsSchema := commonschemagraphbeta.DeviceConfigurationWithAllGroupAssignmentsAndFilterSchema()
	groupIDAttribute := assignmentsSchema.NestedObject.Attributes["group_id"].(schema.StringAttribute)
	groupIDAttribute.Validators = append(
		groupIDAttribute.Validators,
		customvalidator.RequiredWhenEquals("type", types.StringValue("groupAssignmentTarget")),
		customvalidator.RequiredWhenEquals("type", types.StringValue("exclusionGroupAssignmentTarget")),
	)
	assignmentsSchema.NestedObject.Attributes["group_id"] = groupIDAttribute

	return assignmentsSchema
}
