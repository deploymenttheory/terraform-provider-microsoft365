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

// GetTypeName returns the type name of the resource from the state model.
func (r *SettingsCatalogResource) GetTypeName() string {
	return r.TypeName
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
					Attributes: getSettingInstanceAttributes(),
				},
				MarkdownDescription: "Policy settings configuration",
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
			"template_reference": schema.SingleNestedAttribute{
				Optional:            true,
				Attributes:          getTemplateReferenceAttributes(),
				PlanModifiers:       []planmodifier.Object{planmodifiers.ObjectDefaultValueEmpty()},
				Computed:            true,
				MarkdownDescription: "Template reference information",
			},
			"role_scope_tag_ids": schema.SetAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					planmodifiers.DefaultValueSet(
						[]attr.Value{types.StringValue("0")},
					),
				},
				MarkdownDescription: "List of Scope Tags for this Entity instance",
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
			"setting_count": schema.Int64Attribute{
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

// Function to create setting instance attributes with recursion handling
func getSettingInstanceAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"definition_id": schema.StringAttribute{
			Required:            true,
			Description:         `settingDefinitionId`,
			MarkdownDescription: "This is the settings catalog setting definition ID. e.g `device_vendor_msft_bitlocker_removabledrivesexcludedfromencryption`",
		},
		"template_reference": schema.SingleNestedAttribute{
			Optional:            true,
			Attributes:          deviceManagementConfigurationPolicyDeviceManagementConfigurationSettingInstanceTemplateReferenceAttributes,
			Description:         `settingInstanceTemplateReference`,
			MarkdownDescription: "Setting Instance Template Reference",
		},
		"choice": generic.OdataDerivedTypeNestedAttributeRs{
			DerivedType: "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance",
			SingleNestedAttribute: schema.SingleNestedAttribute{
				Optional:            true,
				Attributes:          getChoiceSettingInstanceAttributes(),
				MarkdownDescription: "Choice setting instance",
			},
		},
		"choice_collection": generic.OdataDerivedTypeNestedAttributeRs{
			DerivedType: "#microsoft.graph.deviceManagementConfigurationChoiceSettingCollectionInstance",
			SingleNestedAttribute: schema.SingleNestedAttribute{
				Optional:            true,
				Attributes:          getChoiceSettingCollectionInstanceAttributes(),
				MarkdownDescription: "Choice setting collection instance",
			},
		},
		"group": generic.OdataDerivedTypeNestedAttributeRs{
			DerivedType: "#microsoft.graph.deviceManagementConfigurationGroupSettingInstance",
			SingleNestedAttribute: schema.SingleNestedAttribute{
				Optional:            true,
				Attributes:          getGroupSettingInstanceAttributes(),
				MarkdownDescription: "Group setting instance",
			},
		},
		"group_collection": generic.OdataDerivedTypeNestedAttributeRs{
			DerivedType: "#microsoft.graph.deviceManagementConfigurationGroupSettingCollectionInstance",
			SingleNestedAttribute: schema.SingleNestedAttribute{
				Optional:            true,
				Attributes:          getGroupSettingCollectionInstanceAttributes(),
				MarkdownDescription: "Group setting collection instance",
			},
		},
		"setting_group": generic.OdataDerivedTypeNestedAttributeRs{
			DerivedType: "#microsoft.graph.deviceManagementConfigurationSettingGroupInstance",
			SingleNestedAttribute: schema.SingleNestedAttribute{
				Optional:            true,
				Attributes:          map[string]schema.Attribute{},
				MarkdownDescription: "Setting group instance",
			},
		},
		"setting_group_collection": generic.OdataDerivedTypeNestedAttributeRs{
			DerivedType: "#microsoft.graph.deviceManagementConfigurationSettingGroupCollectionInstance",
			SingleNestedAttribute: schema.SingleNestedAttribute{
				Optional:            true,
				Attributes:          map[string]schema.Attribute{},
				MarkdownDescription: "Setting group collection instance",
			},
		},
		"simple": generic.OdataDerivedTypeNestedAttributeRs{
			DerivedType: "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance",
			SingleNestedAttribute: schema.SingleNestedAttribute{
				Optional:            true,
				Attributes:          getSimpleSettingInstanceAttributes(),
				MarkdownDescription: "Simple setting instance",
			},
		},
		"simple_collection": generic.OdataDerivedTypeNestedAttributeRs{
			DerivedType: "#microsoft.graph.deviceManagementConfigurationSimpleSettingCollectionInstance",
			SingleNestedAttribute: schema.SingleNestedAttribute{
				Optional:            true,
				Attributes:          getSimpleSettingCollectionInstanceAttributes(),
				MarkdownDescription: "Simple setting collection instance",
			},
		},
	}
}

// Common template reference attributes
var deviceManagementConfigurationPolicyDeviceManagementConfigurationSettingValueTemplateReferenceAttributes = map[string]schema.Attribute{
	"default_id": schema.StringAttribute{
		Optional:            true,
		MarkdownDescription: "Default Setting Definition ID to be applied when the template is instantiated.",
	},
	"children": schema.ListNestedAttribute{
		Optional:            true,
		NestedObject:        schema.NestedAttributeObject{Attributes: getTemplateReferenceChildAttributes()},
		MarkdownDescription: "Collection of child setting instance template references.",
	},
}

func getTemplateReferenceChildAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"default_id": schema.StringAttribute{
			Optional:            true,
			MarkdownDescription: "Default Setting Definition ID to be applied when the template is instantiated.",
		},
	}
}

// Simple setting value attributes
var deviceManagementConfigurationPolicyDeviceManagementConfigurationSimpleSettingValueAttributes = map[string]schema.Attribute{
	"string_value": schema.StringAttribute{
		Optional:            true,
		MarkdownDescription: "Simple setting value in string form.",
	},
	"boolean_value": schema.BoolAttribute{
		Optional:            true,
		MarkdownDescription: "Simple setting value in boolean form.",
	},
	"int_value": schema.Int64Attribute{
		Optional:            true,
		MarkdownDescription: "Simple setting value in integer form.",
	},
}

// Template reference instance attributes
var deviceManagementConfigurationPolicyDeviceManagementConfigurationSettingInstanceTemplateReferenceAttributes = map[string]schema.Attribute{
	"default_id": schema.StringAttribute{
		Optional:            true,
		MarkdownDescription: "Default Setting Definition ID to be applied when the template is instantiated.",
	},
	"children": schema.ListNestedAttribute{
		Optional:            true,
		NestedObject:        schema.NestedAttributeObject{Attributes: getTemplateReferenceChildAttributes()},
		MarkdownDescription: "Collection of child setting instance template references.",
	},
}

// Function to create simple setting collection instance attributes
func getSimpleSettingCollectionInstanceAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"values": schema.ListNestedAttribute{
			Required: true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: deviceManagementConfigurationPolicyDeviceManagementConfigurationSimpleSettingValueAttributes,
			},
			Description:         `simpleSettingCollectionValue`,
			MarkdownDescription: "Simple setting collection instance value",
		},
	}
}

// Function to create simple setting instance attributes
func getSimpleSettingInstanceAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"value": schema.SingleNestedAttribute{
			Required:            true,
			Attributes:          deviceManagementConfigurationPolicyDeviceManagementConfigurationSimpleSettingValueAttributes,
			Description:         `simpleSettingValue`,
			MarkdownDescription: "Simple setting instance value",
		},
	}
}

// Function to create choice setting value attributes
func getChoiceSettingValueAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"template_reference": schema.SingleNestedAttribute{
			Optional:            true,
			Attributes:          deviceManagementConfigurationPolicyDeviceManagementConfigurationSettingValueTemplateReferenceAttributes,
			Description:         `settingValueTemplateReference`,
			MarkdownDescription: "Setting value template reference",
		},
		"children": schema.ListNestedAttribute{
			Optional:            true,
			NestedObject:        schema.NestedAttributeObject{Attributes: getSettingInstanceAttributes()},
			PlanModifiers:       []planmodifier.List{planmodifiers.ListDefaultValueEmpty()},
			Computed:            true,
			MarkdownDescription: "Child settings.",
		},
		"value": schema.StringAttribute{
			Optional:            true,
			MarkdownDescription: "Choice setting value: an OptionDefinition ItemId.",
		},
	}
}

// Function to create choice setting collection instance attributes
func getChoiceSettingCollectionInstanceAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"values": schema.ListNestedAttribute{
			Required: true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: getChoiceSettingValueAttributes(),
			},
			Description:         `choiceSettingCollectionValue`,
			MarkdownDescription: "Choice setting collection value",
		},
	}
}

// Function to create group setting value attributes
func getGroupSettingValueAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"template_reference": schema.SingleNestedAttribute{
			Optional:            true,
			Attributes:          deviceManagementConfigurationPolicyDeviceManagementConfigurationSettingValueTemplateReferenceAttributes,
			Description:         `settingValueTemplateReference`,
			MarkdownDescription: "Setting value template reference",
		},
		"children": schema.ListNestedAttribute{
			Optional:            true,
			NestedObject:        schema.NestedAttributeObject{Attributes: getSettingInstanceAttributes()},
			PlanModifiers:       []planmodifier.List{planmodifiers.ListDefaultValueEmpty()},
			Computed:            true,
			MarkdownDescription: "Collection of child setting instances",
		},
	}
}

// Function to create group setting collection instance attributes
func getGroupSettingCollectionInstanceAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"values": schema.ListNestedAttribute{
			Required: true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: getGroupSettingValueAttributes(),
			},
			Description:         `groupSettingCollectionValue`,
			MarkdownDescription: "Collection of group setting values",
		},
	}
}

// Function to create choice setting instance attributes
func getChoiceSettingInstanceAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"value": schema.SingleNestedAttribute{
			Required:            true,
			Attributes:          getChoiceSettingValueAttributes(),
			Description:         `choiceSettingValue`,
			MarkdownDescription: "Choice setting value",
		},
	}
}

// Function to create group setting instance attributes
func getGroupSettingInstanceAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"value": schema.SingleNestedAttribute{
			Required:            true,
			Attributes:          getGroupSettingValueAttributes(),
			Description:         `groupSettingValue`,
			MarkdownDescription: "Group setting value",
		},
	}
}

// Function to create template reference attributes
func getTemplateReferenceAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"display_name": schema.StringAttribute{
			Computed: true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
			Description: `templateDisplayName`,
		},
		"display_version": schema.StringAttribute{
			Computed: true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
			Description: `templateDisplayVersion`,
		},
		"family": schema.StringAttribute{
			Computed: true,
			Validators: []validator.String{
				stringvalidator.OneOf(
					"none",
					"endpointSecurityAntivirus",
					"endpointSecurityDiskEncryption",
					"endpointSecurityFirewall",
					"endpointSecurityEndpointDetectionAndResponse",
					"endpointSecurityAttackSurfaceReduction",
					"endpointSecurityAccountProtection",
					"endpointSecurityApplicationControl",
					"endpointSecurityEndpointPrivilegeManagement",
					"enrollmentConfiguration",
					"appQuietTime",
					"baseline",
					"unknownFutureValue",
					"deviceConfigurationScripts",
					"deviceConfigurationPolicies",
					"windowsOsRecoveryPolicies",
					"companyPortal",
				),
			},
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
			Description:         `templateFamily`,
			MarkdownDescription: "Describes the TemplateFamily for the Template entity",
		},
		"id": schema.StringAttribute{
			Optional: true,
			PlanModifiers: []planmodifier.String{
				planmodifiers.DefaultValueString(""),
				stringplanmodifier.RequiresReplace(),
			},
			Computed:            true,
			Description:         `templateId`,
			MarkdownDescription: "Template id",
		},
	}
}
