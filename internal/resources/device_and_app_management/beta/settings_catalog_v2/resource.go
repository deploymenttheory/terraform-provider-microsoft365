// Base type definitions and interfaces
package graphBetaSettingsCatalog

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/schema"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
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

// Main schema definition incorporating all components
func (r *SettingsCatalogResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Windows Settings Catalog profile in Microsoft Intune.",
		Attributes: map[string]schema.Attribute{
			// Policy base attributes
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				MarkdownDescription: "The unique identifier for this policy",
			},
			"display_name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The policy name",
			},
			"description": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The policy description",
			},
			"platforms": schema.StringAttribute{
				Required:    true,
				Description: "The platforms this settings catalog policy supports. Valid values are: android, androidEnterprise, aosp, iOS, linux, macOS, windows10, windows10X",
				Validators: []validator.String{
					stringvalidator.OneOf(
						"android",
						"androidEnterprise",
						"aosp",
						"iOS",
						"linux",
						"macOS",
						"windows10",
						"windows10X",
					),
				},
			},
			"technologies": schema.StringAttribute{
				Computed:    true,
				Description: "The technologies this profile uses.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"settings_count": schema.Int32Attribute{
				Computed:    true,
				Description: "The number of settings in this profile.",
			},
			"role_scope_tag_ids": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "List of scope tag IDs for this Windows Settings Catalog profile.",
			},
			"template_reference": schema.SingleNestedAttribute{
				Optional:            true,
				Attributes:          deviceManagementConfigurationPolicyDeviceManagementConfigurationSettingValueTemplateReferenceAttributes(),
				MarkdownDescription: "The policy template reference",
			},
			"last_modified_date_time": schema.StringAttribute{
				Computed:    true,
				Description: "The date and time when this profile was last modified.",
			},
			"created_date_time": schema.StringAttribute{
				Computed:    true,
				Description: "The date and time when this profile was created.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			// Settings attributes
			"settings": schema.ListNestedAttribute{
				Required: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: createSettingInstanceSchema(0),
				},
				MarkdownDescription: "The policy settings",
			},
			// Assignments
			"assignments": commonschema.AssignmentsSchema(),

			"timeouts": commonschema.Timeouts(ctx),
		},
	}
}

// Main function to create setting instance schema with recursion
func createSettingInstanceSchema(depth int) map[string]schema.Attribute {
	if depth >= 5 {
		return nil
	}

	baseAttrs := createBaseSettingInstance()

	settingTypes := map[string]schema.Attribute{
		"choice": schema.SingleNestedAttribute{
			Optional:            true,
			Attributes:          createChoiceSettingSchema(depth),
			MarkdownDescription: "Choice setting instance",
		},
		"choice_collection": schema.SingleNestedAttribute{
			Optional:            true,
			Attributes:          createChoiceCollectionSettingSchema(depth),
			MarkdownDescription: "Choice setting collection instance",
		},
		"group": schema.SingleNestedAttribute{
			Optional:            true,
			Attributes:          createGroupSettingSchema(depth),
			MarkdownDescription: "Group setting instance",
		},
		"group_collection": schema.SingleNestedAttribute{
			Optional:            true,
			Attributes:          createGroupCollectionSettingSchema(depth),
			MarkdownDescription: "Group setting collection instance",
		},
		"simple": schema.SingleNestedAttribute{
			Optional:            true,
			Attributes:          createSimpleSettingSchema(),
			MarkdownDescription: "Simple setting instance",
		},
		"simple_collection": schema.SingleNestedAttribute{
			Optional:            true,
			Attributes:          createSimpleCollectionSettingSchema(),
			MarkdownDescription: "Simple setting collection instance",
		},
	}

	for k, v := range settingTypes {
		baseAttrs[k] = v
	}

	return baseAttrs
}

// Helper function to create base setting instance attributes
func createBaseSettingInstance() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"setting_definition_id": schema.StringAttribute{
			Required:            true,
			Description:         `settingDefinitionId`,
			MarkdownDescription: "This is the settings catalog setting definition ID. e.g `device_vendor_msft_bitlocker_removabledrivesexcludedfromencryption`",
		},
		"template_reference": schema.SingleNestedAttribute{
			Optional:            true,
			Attributes:          deviceManagementConfigurationPolicyDeviceManagementConfigurationSettingValueTemplateReferenceAttributes(),
			Description:         `settingInstanceTemplateReference`,
			MarkdownDescription: "Setting Instance Template Reference",
		},
	}
}

// Template reference attributes used throughout the schema
// Creates template reference attributes used throughout the schema
func deviceManagementConfigurationPolicyDeviceManagementConfigurationSettingValueTemplateReferenceAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"display_name": schema.StringAttribute{
			Computed:            true,
			MarkdownDescription: "Template display name",
		},
		"display_version": schema.StringAttribute{
			Computed:            true,
			MarkdownDescription: "Template display version",
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
			MarkdownDescription: "Template family type. Ref https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfigv2-deviceManagementConfigurationPolicyTemplateReference?view=graph-rest-beta",
		},
		"id": schema.StringAttribute{
			Optional:            true,
			MarkdownDescription: "Template ID",
		},
	}
}

// Create choice setting schema, supporting recursive children as per the JSON structure.
func createChoiceSettingSchema(depth int) map[string]schema.Attribute {
	if depth >= 5 {
		return nil // Avoid excessive recursion for practical limits.
	}

	return map[string]schema.Attribute{
		"value": schema.SingleNestedAttribute{
			Required: true,
			Attributes: map[string]schema.Attribute{
				"odata_type": schema.StringAttribute{
					Computed:            true,
					MarkdownDescription: "OData type for choice setting value",
				},
				"template_reference": createSettingValueTemplateReference(),
				"children": schema.ListNestedAttribute{
					Optional: true,
					NestedObject: schema.NestedAttributeObject{
						Attributes: createChoiceSettingSchema(depth + 1),
					},
					MarkdownDescription: "Recursive child settings",
				},
				"value": schema.StringAttribute{
					Optional:            true,
					MarkdownDescription: "Choice setting value: an OptionDefinition ItemId",
				},
			},
			Description:         `choiceSettingValue`,
			MarkdownDescription: "Choice setting value with recursive children",
		},
	}
}

// Helper function to create simple setting value attributes
func createSimpleSettingValueAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"odata_type": schema.StringAttribute{
			Computed:            true,
			MarkdownDescription: "OData type for simple setting value",
		},
		"string_value": schema.StringAttribute{
			Optional:            true,
			MarkdownDescription: "Simple string value",
		},
		"boolean_value": schema.BoolAttribute{
			Optional:            true,
			MarkdownDescription: "Simple boolean value",
		},
		"int32_value": schema.Int32Attribute{
			Optional:            true,
			MarkdownDescription: "Simple integer value",
		},
	}
}

// Creates choice collection setting schema
func createChoiceCollectionSettingSchema(depth int) map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"values": schema.ListNestedAttribute{
			Required: true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"odata_type": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "OData type for choice collection value",
					},
					"template_reference": createSettingValueTemplateReference(),
					"children": schema.ListNestedAttribute{
						Optional: true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: createSettingInstanceSchema(depth + 1),
						},

						Computed:            true,
						MarkdownDescription: "Child settings",
					},
					"value": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "Choice setting value",
					},
				},
			},
			Description:         `choiceSettingCollectionValue`,
			MarkdownDescription: "Choice setting collection values",
		},
	}
}

// Creates group setting schema
func createGroupSettingSchema(depth int) map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"value": schema.SingleNestedAttribute{
			Required: true,
			Attributes: map[string]schema.Attribute{
				"odata_type": schema.StringAttribute{
					Computed:            true,
					MarkdownDescription: "OData type for group setting value",
				},
				"template_reference": createSettingValueTemplateReference(),
				"children": schema.ListNestedAttribute{
					Optional: true,
					NestedObject: schema.NestedAttributeObject{
						Attributes: createSettingInstanceSchema(depth + 1),
					},

					Computed:            true,
					MarkdownDescription: "Child settings",
				},
			},
			Description:         `groupSettingValue`,
			MarkdownDescription: "Group setting value",
		},
	}
}

// Creates group collection setting schema
func createGroupCollectionSettingSchema(depth int) map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"values": schema.ListNestedAttribute{
			Required: true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					// Template reference still needed for MS Graph API
					// "template_reference": schema.SingleNestedAttribute{
					// 	Optional: true,
					// 	Computed: true,
					// 	PlanModifiers: []planmodifier.Object{
					// 		objectplanmodifier.UseStateForUnknown(),
					// 	},
					// 	Attributes: createSettingValueTemplateReference(),
					// },
					"children": schema.ListNestedAttribute{
						Optional: true,
						Computed: true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: createSettingInstanceSchema(depth + 1),
						},
						PlanModifiers: []planmodifier.List{
							listplanmodifier.UseStateForUnknown(),
							listplanmodifier.RequiresReplace(),
						},
						MarkdownDescription: "Child settings",
					},
				},
			},
			PlanModifiers: []planmodifier.List{
				listplanmodifier.UseStateForUnknown(),
			},
			Description:         `groupSettingCollectionValue`,
			MarkdownDescription: "Group setting collection values",
		},
	}
}

// Creates simple setting schema
func createSimpleSettingSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"value": schema.SingleNestedAttribute{
			Required: true,
			Attributes: map[string]schema.Attribute{
				"odata_type": schema.StringAttribute{
					Computed:            true,
					MarkdownDescription: "OData type for simple setting value",
				},
				"value": schema.SingleNestedAttribute{
					Required:            true,
					Attributes:          createSimpleSettingValueAttributes(),
					Description:         `simpleSettingValue`,
					MarkdownDescription: "Simple setting value",
				},
			},
		},
	}
}

// Creates simple collection setting schema
func createSimpleCollectionSettingSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"values": schema.ListNestedAttribute{
			Required: true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: createSimpleSettingValueAttributes(),
			},
			Description:         `simpleSettingCollectionValue`,
			MarkdownDescription: "Simple setting collection values",
		},
	}
}

// Helper function to create template reference structure
func createSettingValueTemplateReference() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Optional:            true,
		Attributes:          deviceManagementConfigurationPolicyDeviceManagementConfigurationSettingValueTemplateReferenceAttributes(),
		Description:         `settingValueTemplateReference`,
		MarkdownDescription: "Setting value template reference information",
	}
}
