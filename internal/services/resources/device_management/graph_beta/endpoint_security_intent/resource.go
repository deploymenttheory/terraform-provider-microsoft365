package graphBetaEndpointSecurityIntent

import (
	"context"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/identityschema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	commonschemagraphbeta "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema/graph_beta/device_management"
)

const (
	ResourceName  = "microsoft365_graph_beta_device_management_endpoint_security_intent"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	// Basic resource interface (CRUD operations)
	_ resource.Resource = &EndpointSecurityIntentResource{}

	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &EndpointSecurityIntentResource{}

	// Enables import functionality
	_ resource.ResourceWithImportState = &EndpointSecurityIntentResource{}

	// Enables plan modification/diff suppression
	_ resource.ResourceWithModifyPlan = &EndpointSecurityIntentResource{}

	// Enables settings and assignment config validation
	_ resource.ResourceWithValidateConfig = &EndpointSecurityIntentResource{}

	// Enables identity schema for list resource support
	_ resource.ResourceWithIdentity = &EndpointSecurityIntentResource{}
)

// https://learn.microsoft.com/en-us/graph/api/intune-deviceintent-devicemanagementintent-create?view=graph-rest-beta
func NewEndpointSecurityIntentResource() resource.Resource {
	return &EndpointSecurityIntentResource{
		ReadPermissions: []string{
			"DeviceManagementConfiguration.Read.All",
		},
		WritePermissions: []string{
			"DeviceManagementConfiguration.ReadWrite.All",
		},
	}
}

type EndpointSecurityIntentResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
}

// Metadata returns the resource type name.
func (r *EndpointSecurityIntentResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = ResourceName
}

// Configure sets the client for the resource.
func (r *EndpointSecurityIntentResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, ResourceName)
}

// ImportState imports the resource state.
// The import ID is the intent ID, optionally followed by a slash and a comma-separated list of the setting definition IDs to manage,
// for example `<intent_id>/<definition_id>,<definition_id>`. Read then fills in the type and value of each listed setting.
func (r *EndpointSecurityIntentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	intentId, definitionIds, err := parseImportId(req.ID)
	if err != nil {
		resp.Diagnostics.AddError("Invalid import ID", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), intentId)...)
	if len(definitionIds) == 0 {
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("settings"), importedSettings(definitionIds))...)
}

// IdentitySchema defines the identity schema for this resource, used by list operations to uniquely identify instances
func (r *EndpointSecurityIntentResource) IdentitySchema(ctx context.Context, req resource.IdentitySchemaRequest, resp *resource.IdentitySchemaResponse) {
	resp.IdentitySchema = identityschema.Schema{
		Attributes: map[string]identityschema.Attribute{
			"id": identityschema.StringAttribute{
				RequiredForImport: true,
			},
		},
	}
}

// Schema defines the schema for the resource.
func (r *EndpointSecurityIntentResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a template-based Intune endpoint security policy (intent) using the `/deviceManagement/intents` endpoint. " +
			"An intent is created from a security template with `/deviceManagement/templates/{templateId}/createInstance`, its settings are changed with `/deviceManagement/intents/{id}/updateSettings` " +
			"and its assignments are replaced with `/deviceManagement/intents/{id}/assign`. Microsoft is moving these templates to settings catalog configuration policies, " +
			"so this resource is mainly for managing existing estates; new policies should use the settings catalog where the template has been migrated.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique identifier of the intent.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"template_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the security template the intent is created from, for example `d1174162-1dd2-4976-affc-6667049ab0ae` for the BitLocker disk encryption template. " +
					"List the available templates with `GET /deviceManagement/templates`. Changing this value forces a new intent to be created.",
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.RegexMatches(regexp.MustCompile(constants.GuidRegex), "must be a valid GUID in the format '00000000-0000-0000-0000-000000000000'"),
				},
			},
			"display_name": schema.StringAttribute{
				MarkdownDescription: "The display name of the intent.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The description of the intent. Defaults to an empty string.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"role_scope_tag_ids": schema.SetAttribute{
				MarkdownDescription: "The role scope tag IDs for the intent. Defaults to `[\"0\"]` (the Default scope tag).",
				ElementType:         types.StringType,
				Optional:            true,
				Computed:            true,
				Default:             setdefault.StaticValue(types.SetValueMust(types.StringType, []attr.Value{types.StringValue("0")})),
				Validators: []validator.Set{
					setvalidator.SizeAtLeast(1),
				},
			},
			"settings": schema.SetNestedAttribute{
				MarkdownDescription: "The template settings managed by this resource. Only the settings listed here are written and tracked; every other template setting keeps the value Intune gives it. " +
					"Settings are sent to `createInstance` as the settings delta on create and to `updateSettings` on update, so changing them never recreates the intent. " +
					"Removing a setting from configuration stops managing it but leaves its current value in Intune.",
				Optional: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"definition_id": schema.StringAttribute{
							MarkdownDescription: "The setting definition ID, for example `deviceConfiguration--windows10EndpointProtectionConfiguration_bitLockerEncryptDevice`. " +
								"List the definitions of an existing intent with `GET /deviceManagement/intents/{id}/settings`.",
							Required: true,
							Validators: []validator.String{
								stringvalidator.LengthAtLeast(1),
							},
						},
						"setting_type": schema.StringAttribute{
							MarkdownDescription: "The setting instance type, which Graph needs to accept the value. Must be one of: `boolean`, `integer`, `string`, `complex`, `collection`. " +
								"It matches the `@odata.type` of the setting returned by Graph, for example `boolean` for `#microsoft.graph.deviceManagementBooleanSettingInstance`.",
							Required: true,
							Validators: []validator.String{
								stringvalidator.OneOf(settingTypeBoolean, settingTypeInteger, settingTypeString, settingTypeComplex, settingTypeCollection),
							},
						},
						"value_json": schema.StringAttribute{
							MarkdownDescription: "The setting value as JSON, for example `jsonencode(true)` or `jsonencode(\"enabledForAzureAd\")`. Values that differ from Graph only in JSON formatting are not reported as drift.",
							Required:            true,
						},
					},
				},
			},
			"assignments": commonschemagraphbeta.DeviceConfigurationWithAllGroupAssignmentsAndFilterSchema(),
			"is_assigned": schema.BoolAttribute{
				MarkdownDescription: "Whether the intent is assigned to any groups.",
				Computed:            true,
			},
			"is_migrating_to_configuration_policy": schema.BoolAttribute{
				MarkdownDescription: "Whether Intune is migrating the intent to a settings catalog configuration policy.",
				Computed:            true,
			},
			"last_modified_date_time": schema.StringAttribute{
				MarkdownDescription: "The date and time the intent was last modified.",
				Computed:            true,
			},
			"timeouts": commonschema.ResourceTimeouts(ctx),
		},
	}
}
