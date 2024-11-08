package graphBetaSettingsCatalog

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common"
	planmodifiers "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/plan_modifiers"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/schema"
	customValidator "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/validators"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

var (
	// Basic resource interface (CRUD operations)
	_ resource.Resource = &SettingsCatalogResource{}

	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &SettingsCatalogResource{}

	// Enables import functionality
	_ resource.ResourceWithImportState = &SettingsCatalogResource{}

	// Enables plan modification/diff suppression
	_ resource.ResourceWithModifyPlan = &SettingsCatalogResource{}
)

const (
	// Define maximum allowed depth for nested settings catalog settings
	maxDepth = 5
)

func NewSettingsCatalogResource() resource.Resource {
	return &SettingsCatalogResource{
		ReadPermissions: []string{
			"DeviceManagementApps.Read.All",
			"DeviceManagementConfiguration.Read.All",
		},
		WritePermissions: []string{
			"DeviceManagementApps.ReadWrite.All",
			"DeviceManagementConfiguration.ReadWrite.All",
		},
	}
}

type SettingsCatalogResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ProviderTypeName string
	TypeName         string
	ReadPermissions  []string
	WritePermissions []string
}

// Metadata returns the resource type name.
func (r *SettingsCatalogResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_graph_beta_device_and_app_management_settings_catalog"
}

// Configure sets the client for the resource.
func (r *SettingsCatalogResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = common.SetGraphBetaClientForResource(ctx, req, resp, r.TypeName)
}

// ImportState imports the resource state.
func (r *SettingsCatalogResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Function to create the full device management configuration policy schema
func (r *SettingsCatalogResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Settings Catalog profile in Microsoft Intune.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				MarkdownDescription: "The unique identifier for this policy",
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Policy name",
			},
			"description": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				PlanModifiers:       []planmodifier.String{planmodifiers.DefaultValueString("")},
				MarkdownDescription: "Policy description",
			},
			"settings": schema.ListNestedAttribute{
				Required: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: settingInstance(0),
				},
				MarkdownDescription: "Policy settings configuration. Supports up to 5 levels of nesting. Starts at depth 0.",
			},
			"platforms": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					customValidator.EnumValues(
						"none", "android", "iOS", "macOS", "windows10X",
						"windows10", "linux", "unknownFutureValue",
						"androidEnterprise", "aosp",
					),
				},
				PlanModifiers:       []planmodifier.String{planmodifiers.DefaultValueString("windows10")},
				Computed:            true,
				MarkdownDescription: "Platforms for this policy",
			},
			"technologies": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					customValidator.EnumValues(
						"none", "mdm", "windows10XManagement", "configManager",
						"intuneManagementExtension", "thirdParty", "documentGateway",
						"appleRemoteManagement", "microsoftSense", "exchangeOnline",
						"mobileApplicationManagement", "linuxMdm", "enrollment",
						"endpointPrivilegeManagement", "unknownFutureValue",
						"windowsOsRecovery", "android",
					),
				},
				PlanModifiers:       []planmodifier.String{planmodifiers.DefaultValueString("mdm")},
				Computed:            true,
				MarkdownDescription: "Technologies for this policy",
			},
			"role_scope_tag_ids": schema.ListAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				MarkdownDescription: "List of scope tag IDs for this Windows Settings Catalog profile.",
				PlanModifiers: []planmodifier.List{
					planmodifiers.DefaultListValue(
						[]attr.Value{types.StringValue("0")},
					),
				},
			},

			"created_date_time": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				MarkdownDescription: "Creation date and time of the policy",
			},
			"last_modified_date_time": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				MarkdownDescription: "Last modification date and time of the policy",
			},
			"settings_count": schema.Int64Attribute{
				Computed: true,
				PlanModifiers: []planmodifier.Int64{
					planmodifiers.UseStateForUnknownInt64(),
				},
				MarkdownDescription: "Number of settings in the policy",
			},
			"is_assigned": schema.BoolAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.Bool{
					planmodifiers.UseStateForUnknownBool(),
				},
				MarkdownDescription: "Indicates if the policy is assigned to any scope",
			},
			"assignments": commonschema.AssignmentsSchema(),

			"timeouts": commonschema.Timeouts(ctx),
		},
	}
}

func settingInstance(depth int) map[string]schema.Attribute {
	if depth >= maxDepth {
		return map[string]schema.Attribute{}
	}

	return map[string]schema.Attribute{
		"odata_type": schema.StringAttribute{
			Computed:            true,
			MarkdownDescription: "The OData type of the setting instance. Always #microsoft.graph.deviceManagementConfigurationSetting",
		},
		"setting_instance": schema.SingleNestedAttribute{
			Optional:   true,
			Attributes: settingInstanceValueType(depth),
		},
	}
}

func settingInstanceValueType(depth int) map[string]schema.Attribute {
	if depth >= maxDepth {
		return map[string]schema.Attribute{}
	}

	return map[string]schema.Attribute{
		"odata_type": schema.StringAttribute{
			Required:            true,
			MarkdownDescription: "The OData type of the setting instance. This must be specified and is used to determine the specific setting instance type.",
			Validators: []validator.String{
				stringvalidator.OneOf(
					DeviceManagementConfigurationChoiceSettingInstance,
					DeviceManagementConfigurationChoiceSettingCollectionInstance,
					DeviceManagementConfigurationSimpleSettingInstance,
					DeviceManagementConfigurationSimpleSettingCollectionInstance,
					DeviceManagementConfigurationSettingGroupInstance,
					DeviceManagementConfigurationGroupSettingInstance,
					DeviceManagementConfigurationSettingGroupCollectionInstance,
					DeviceManagementConfigurationGroupSettingCollectionInstance,
				),
			},
		},
		"setting_definition_id": schema.StringAttribute{
			Required:            true,
			Description:         "settingDefinitionId",
			MarkdownDescription: "The settings catalog setting definition ID, e.g., `device_vendor_msft_bitlocker_removabledrivesexcludedfromencryption`.",
		},
		"choice": schema.SingleNestedAttribute{
			Optional:   true,
			Attributes: choiceSettingInstanceAttributes(depth + 1), // TODO "choice": GetChoiceSchema(),
			MarkdownDescription: "Choice setting instance with @odata.type: " +
				"#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance.\n\n" +
				"For details, see [Choice Setting Instance Documentation](https://learn.microsoft.com/en-us/graph/" +
				"api/resources/intune-deviceconfigv2-deviceManagementConfigurationChoiceSettingInstance?view=graph-rest-beta).",
		},
		"choice_collection": schema.SingleNestedAttribute{
			Optional:   true,
			Attributes: choiceSettingCollectionInstanceAttributes(depth + 1),
			MarkdownDescription: "Choice setting collection instance with @odata.type: " +
				"#microsoft.graph.deviceManagementConfigurationChoiceSettingCollectionInstance.\n\n" +
				"For details, see [Choice Setting Collection Documentation](https://learn.microsoft.com/en-us/graph/" +
				"api/resources/intune-deviceconfigv2-deviceManagementConfigurationChoiceSettingCollectionInstance?view=graph-rest-beta).",
		},
		"group": schema.SingleNestedAttribute{
			Optional:   true,
			Attributes: groupSettingInstanceAttributes(depth + 1),
			MarkdownDescription: "Group setting instance with @odata.type: " +
				"#microsoft.graph.deviceManagementConfigurationGroupSettingInstance.\n\n" +
				"For details, see [Group Setting Instance Documentation](https://learn.microsoft.com/en-us/graph/" +
				"api/resources/intune-deviceconfigv2-deviceManagementConfigurationGroupSettingInstance?view=graph-rest-beta).",
		},
		"group_collection": schema.SingleNestedAttribute{
			Optional:   true,
			Attributes: groupSettingCollectionInstanceAttributes(depth + 1),
			MarkdownDescription: "Group setting collection instance with @odata.type: " +
				"#microsoft.graph.deviceManagementConfigurationGroupSettingCollectionInstance.\n\n" +
				"For details, see [Group Setting Collection Documentation](https://learn.microsoft.com/en-us/graph/" +
				"api/resources/intune-deviceconfigv2-deviceManagementConfigurationGroupSettingCollectionInstance?view=graph-rest-beta).",
		},
		// "setting_group": schema.SingleNestedAttribute{
		// 	Optional: true,
		// 	Attributes: map[string]schema.Attribute{
		// 		"children": schema.ListNestedAttribute{
		// 			Required: true,
		// 			NestedObject: schema.NestedAttributeObject{
		// 				Attributes: settingInstance(depth + 1),
		// 			},
		// 			PlanModifiers: []planmodifier.List{
		// 				planmodifiers.DefaultListEmptyValue(),
		// 			},
		// 		},
		// 	},
		// 	MarkdownDescription: "Setting group instance with @odata.type: " +
		// 		"#microsoft.graph.deviceManagementConfigurationSettingGroupInstance.\n\n" +
		// 		"For details, see [Setting Group Instance Documentation](https://learn.microsoft.com/en-us/graph/" +
		// 		"api/resources/intune-deviceconfigv2-deviceManagementConfigurationSettingGroupInstance?view=graph-rest-beta).",
		// },
		// "setting_group_collection": schema.SingleNestedAttribute{
		// 	Optional: true,
		// 	Attributes: map[string]schema.Attribute{
		// 		"children": schema.ListNestedAttribute{
		// 			Required: true,
		// 			NestedObject: schema.NestedAttributeObject{
		// 				Attributes: settingInstance(depth + 1),
		// 			},
		// 		},
		// 	},
		// 	MarkdownDescription: "Setting group collection instance with @odata.type: " +
		// 		"#microsoft.graph.deviceManagementConfigurationSettingGroupCollectionInstance.\n\n" +
		// 		"For details, see [Setting Group Collection Documentation](https://learn.microsoft.com/en-us/graph/" +
		// 		"api/resources/intune-deviceconfigv2-deviceManagementConfigurationSettingGroupCollectionInstance?view=graph-rest-beta).",
		// },
		"simple": schema.SingleNestedAttribute{
			Optional:   true,
			Attributes: simpleSettingInstanceAttributes(),
			MarkdownDescription: "Simple setting instance with @odata.type: " +
				"#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance.\n\n" +
				"For details, see [Simple Setting Instance Documentation](https://learn.microsoft.com/en-us/graph/" +
				"api/resources/intune-deviceconfigv2-deviceManagementConfigurationSimpleSettingInstance?view=graph-rest-beta).",
		},
		"simple_collection": schema.SingleNestedAttribute{
			Optional:   true,
			Attributes: simpleSettingCollectionInstanceAttributes(),
			MarkdownDescription: "Simple setting collection instance with @odata.type: " +
				"#microsoft.graph.deviceManagementConfigurationSimpleSettingCollectionInstance.\n\n" +
				"For details, see [Simple Setting Collection Documentation](https://learn.microsoft.com/en-us/graph/" +
				"api/resources/intune-deviceconfigv2-deviceManagementConfigurationSimpleSettingCollectionInstance?view=graph-rest-beta).",
		},
	}
}

// choiceSettingInstanceAttributes dynamically adds children to handle recursion
func choiceSettingInstanceAttributes(depth int) map[string]schema.Attribute {
	if depth >= maxDepth {
		return map[string]schema.Attribute{}
	}

	attributes := map[string]schema.Attribute{
		"odata_type": schema.StringAttribute{
			Computed:            true,
			MarkdownDescription: "The OData type of the setting instance. This is automatically set by the graph SDK during request construction.",
		},
		"integer_value": schema.Int32Attribute{
			Optional: true,
			MarkdownDescription: "Value of the integer setting with @odata.type: #microsoft.graph.deviceManagementConfigurationIntegerSettingValue.\n\n" +
				"For more details, see [Intune Integer Setting Value Documentation](https://learn.microsoft.com/en-us/graph/" +
				"api/resources/intune-deviceconfigv2-deviceManagementConfigurationIntegerSettingValue?view=graph-rest-beta).",
		},
		"string_value": schema.StringAttribute{
			Optional: true,
			MarkdownDescription: "Value of the string setting with @odata.type: #microsoft.graph.deviceManagementConfigurationStringSettingValue.\n\n" +
				"For more details, see [String Setting Value Documentation](https://learn.microsoft.com/en-us/graph/" +
				"api/resources/intune-deviceconfigv2-deviceManagementConfigurationStringSettingValue?view=graph-rest-beta).",
		},
		"children": schema.ListNestedAttribute{
			Optional:    true,
			Description: "The child settings of this choice setting.",
			NestedObject: schema.NestedAttributeObject{
				Attributes: settingInstanceValueType(depth + 1),
			},
			MarkdownDescription: "Nested child settings instances, allowing recursive configurations.",
			PlanModifiers:       []planmodifier.List{planmodifiers.DefaultListEmptyValue()},
		},
	}

	return attributes
}

// Function to create simple setting collection instance attributes
func simpleSettingCollectionInstanceAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"values": schema.ListNestedAttribute{
			Required: true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: deviceManagementConfigurationSimpleSettingValueAttributes,
			},
			Description:         `simpleSettingCollectionValue`,
			MarkdownDescription: "Simple setting collection instance value",
		},
	}
}

// Function to create simple setting instance attributes
func simpleSettingInstanceAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"value": schema.SingleNestedAttribute{
			Required:            true,
			Attributes:          deviceManagementConfigurationSimpleSettingValueAttributes,
			Description:         `simpleSettingValue`,
			MarkdownDescription: "Simple setting instance value",
		},
	}
}

// Function to create choice setting value attributes
func choiceSettingValueAttributes(depth int) map[string]schema.Attribute {
	if depth >= maxDepth {
		return map[string]schema.Attribute{}
	}

	return map[string]schema.Attribute{
		// "children": schema.ListNestedAttribute{
		// 	Optional:            true,
		// 	NestedObject:        schema.NestedAttributeObject{Attributes: settingInstance(depth + 1)},
		// 	PlanModifiers:       []planmodifier.List{planmodifiers.DefaultListEmptyValue()},
		// 	Computed:            true,
		// 	MarkdownDescription: "Child settings.",
		// },
		"value": schema.StringAttribute{
			Optional:            true,
			MarkdownDescription: "Choice setting value: an OptionDefinition ItemId.",
		},
	}
}

// Function to create choice setting collection instance attributes
func choiceSettingCollectionInstanceAttributes(depth int) map[string]schema.Attribute {
	if depth >= maxDepth {
		return map[string]schema.Attribute{}
	}

	return map[string]schema.Attribute{
		"values": schema.ListNestedAttribute{
			Required: true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: choiceSettingValueAttributes(depth + 1),
			},
			Description:         `choiceSettingCollectionValue`,
			MarkdownDescription: "Choice setting collection value",
		},
	}
}

// Function to create group setting value attributes
func groupSettingValueAttributes(depth int) map[string]schema.Attribute {
	if depth >= maxDepth {
		return map[string]schema.Attribute{}
	}

	return map[string]schema.Attribute{
		"children": schema.ListNestedAttribute{
			Optional:            true,
			NestedObject:        schema.NestedAttributeObject{Attributes: settingInstance(depth + 1)},
			PlanModifiers:       []planmodifier.List{planmodifiers.DefaultListEmptyValue()},
			Computed:            true,
			MarkdownDescription: "Collection of child setting instances",
		},
	}
}

// Function to create group setting collection instance attributes
func groupSettingCollectionInstanceAttributes(depth int) map[string]schema.Attribute {
	if depth >= maxDepth {
		return map[string]schema.Attribute{}
	}

	return map[string]schema.Attribute{
		"values": schema.ListNestedAttribute{
			Required: true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: groupSettingValueAttributes(depth + 1),
			},
			Description:         `groupSettingCollectionValue`,
			MarkdownDescription: "Collection of group setting values",
		},
	}
}

// Function to create group setting instance attributes
func groupSettingInstanceAttributes(depth int) map[string]schema.Attribute {
	if depth >= maxDepth {
		return map[string]schema.Attribute{}
	}

	return map[string]schema.Attribute{
		"value": schema.SingleNestedAttribute{
			Required:            true,
			Attributes:          groupSettingValueAttributes(depth + 1),
			Description:         `groupSettingValue`,
			MarkdownDescription: "Group setting value",
		},
	}
}

// Simple setting value attributes
var deviceManagementConfigurationSimpleSettingValueAttributes = map[string]schema.Attribute{
	"odata_type": schema.StringAttribute{
		Computed:            true,
		MarkdownDescription: "The OData type of the setting instance. This is automatically set by the graph SDK during request construction.",
	},
	"integer_value": schema.Int64Attribute{
		Optional: true,
		MarkdownDescription: "Value of the integer setting with @odata.type: #microsoft.graph.deviceManagementConfigurationIntegerSettingValue.\n\n" +
			"For more details, see [Intune Integer Setting Value Documentation](https://learn.microsoft.com/en-us/graph/" +
			"api/resources/intune-deviceconfigv2-deviceManagementConfigurationIntegerSettingValue?view=graph-rest-beta).",
	},
	"reference": schema.SingleNestedAttribute{
		Optional: true,
		Attributes: map[string]schema.Attribute{
			"value": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Value of the reference setting.",
			},
			"note": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "A note for contextual information provided by the admin.",
			},
		},
		MarkdownDescription: "Model for ReferenceSettingValue with @odata.type: " +
			"#microsoft.graph.deviceManagementConfigurationReferenceSettingValue.\n\n" +
			"For more details, see [Reference Setting Value Documentation](https://learn.microsoft.com/en-us/graph/" +
			"api/resources/intune-deviceconfigv2-deviceManagementConfigurationReferenceSettingValue?view=graph-rest-beta).",
	},
	"secret": schema.SingleNestedAttribute{
		Optional: true,
		Attributes: map[string]schema.Attribute{
			"secret_value": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Value of the secret setting.",
			},
			"state": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.OneOf("invalid", "notEncrypted", "encryptedValueToken"),
				},
				Description: `valueState`,
				MarkdownDescription: "Indicates the encryption state of the Value property, with possible values:\n\n" +
					"- **invalid**: Default invalid value.\n" +
					"- **notEncrypted**: Secret value is not encrypted.\n" +
					"- **encryptedValueToken**: A token for the encrypted value is returned by the service.\n\n" +
					"Model for SecretSettingValue with @odata.type: " +
					"#microsoft.graph.deviceManagementConfigurationSecretSettingValue.\n\n" +
					"For more details, see [Secret Setting Value Documentation](https://learn.microsoft.com/en-us/graph/" +
					"api/resources/intune-deviceconfigv2-deviceManagementConfigurationSecretSettingValue?view=graph-rest-beta).",
			},
		},
		MarkdownDescription: "Graph model for a secret setting value with @odata.type: " +
			"#microsoft.graph.deviceManagementConfigurationSecretSettingValue.\n\n" +
			"For more details, see [Secret Setting Value Documentation](https://learn.microsoft.com/en-us/graph/" +
			"api/resources/intune-deviceconfigv2-deviceManagementConfigurationSecretSettingValue?view=graph-rest-beta).",
	},
	"string_value": schema.StringAttribute{
		Optional: true,
		MarkdownDescription: "Value of the string setting with @odata.type: #microsoft.graph.deviceManagementConfigurationStringSettingValue.\n\n" +
			"For more details, see [String Setting Value Documentation](https://learn.microsoft.com/en-us/graph/" +
			"api/resources/intune-deviceconfigv2-deviceManagementConfigurationStringSettingValue?view=graph-rest-beta).",
	},
}
