package graphBetaWindowsSettingsCatalog

import (
	"context"
	"fmt"
	"regexp"

	helpers "github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/schema"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

var _ resource.Resource = &WindowsSettingsCatalogResource{}
var _ resource.ResourceWithConfigure = &WindowsSettingsCatalogResource{}
var _ resource.ResourceWithImportState = &WindowsSettingsCatalogResource{}
var guidRegex = regexp.MustCompile(helpers.GuidRegex)

func NewWindowsSettingsCatalogResource() resource.Resource {
	return &WindowsSettingsCatalogResource{
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

type WindowsSettingsCatalogResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ProviderTypeName string
	TypeName         string
	ReadPermissions  []string
	WritePermissions []string
}

// GetTypeName returns the type name of the resource from the state model.
func (r *WindowsSettingsCatalogResource) GetTypeName() string {
	return r.TypeName
}

// Metadata returns the resource type name.
func (r *WindowsSettingsCatalogResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_graph_beta_device_and_app_management_windows_settings_catalog"
}

// Configure sets the client for the resource.
func (r *WindowsSettingsCatalogResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = common.SetGraphBetaClientForResource(ctx, req, resp, r.TypeName)
}

// ImportState imports the resource state.
func (r *WindowsSettingsCatalogResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *WindowsSettingsCatalogResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Windows Settings Catalog profile in Microsoft Intune.",
		Attributes: map[string]schema.Attribute{
			// Profile attributes
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The unique identifier for this Windows Settings Catalog profile.",
			},
			"display_name": schema.StringAttribute{
				Required:    true,
				Description: "The display name of the Windows Settings Catalog profile.",
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "The description of the Windows Settings Catalog profile.",
			},
			"platforms": schema.StringAttribute{
				Computed:    true,
				Description: "The platforms this profile supports.",
			},
			"technologies": schema.StringAttribute{
				Computed:    true,
				Description: "The technologies this profile uses.",
			},
			"settings_count": schema.Int32Attribute{
				Computed:    true,
				Description: "The number of settings in this profile.",
			},
			"creation_source": schema.StringAttribute{
				Computed:    true,
				Description: "The source of creation for this profile.",
			},
			"role_scope_tag_ids": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "List of scope tag IDs for this Windows Settings Catalog profile.",
			},
			"is_assigned": schema.BoolAttribute{
				Computed:    true,
				Description: "Indicates whether this profile is assigned to any groups.",
			},
			"last_modified_date_time": schema.StringAttribute{
				Computed:    true,
				Description: "The date and time when this profile was last modified.",
			},
			"created_date_time": schema.StringAttribute{
				Computed:    true,
				Description: "The date and time when this profile was created.",
			},

			// Settings attributes
			"settings": settingsSchema(),

			// Assignments
			"assignments": assignmentsSchema(),

			"timeouts": commonschema.Timeouts(ctx),
		},
	}
}
func settingsSchema() schema.ListNestedAttribute {
	return schema.ListNestedAttribute{
		Required:    true,
		Description: "The list of settings for this Windows Settings Catalog profile.",
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"odata_type": schema.StringAttribute{
					Required:    true,
					Description: "The OData type of the setting.",
				},
				"id": schema.StringAttribute{
					Computed:    true,
					Description: "The ID of the setting.",
				},
				"setting_instance": schema.SingleNestedAttribute{
					Required:    true,
					Description: "The instance of the setting.",
					Attributes: map[string]schema.Attribute{
						"odata_type": schema.StringAttribute{
							Required:    true,
							Description: "The OData type of the setting instance.",
						},
						"setting_definition_id": schema.StringAttribute{
							Required:    true,
							Description: "The ID of the setting definition.",
						},
						"setting_instance_template_reference": schema.SingleNestedAttribute{
							Optional:    true,
							Description: "The template reference for this setting instance.",
							Attributes: map[string]schema.Attribute{
								"odata_type": schema.StringAttribute{
									Required:    true,
									Description: "The OData type of the setting instance template reference.",
								},
								"setting_instance_template_id": schema.StringAttribute{
									Required:    true,
									Description: "The ID of the setting instance template.",
								},
							},
						},
						"choice_setting_value": schema.SingleNestedAttribute{
							Optional:    true,
							Description: "The choice setting value.",
							Attributes: map[string]schema.Attribute{
								"odata_type": schema.StringAttribute{
									Required:    true,
									Description: "The OData type of the choice setting value.",
								},
								"setting_value_template_reference": schema.SingleNestedAttribute{
									Optional:    true,
									Description: "The template reference for this setting value.",
									Attributes: map[string]schema.Attribute{
										"odata_type": schema.StringAttribute{
											Required:    true,
											Description: "The OData type of the setting value template reference.",
										},
										"setting_value_template_id": schema.StringAttribute{
											Required:    true,
											Description: "The ID of the setting value template.",
										},
										"use_template_default": schema.BoolAttribute{
											Optional:    true,
											Description: "Whether to use the template default value.",
										},
									},
								},
								"int_value": schema.Int32Attribute{
									Optional:    true,
									Description: "The integer value of the choice setting.",
								},
								"string_value": schema.StringAttribute{
									Optional:    true,
									Description: "The string value of the choice setting.",
								},
								"children": schema.ListNestedAttribute{
									Optional:    true,
									Description: "The child settings of this choice setting.",
									NestedObject: schema.NestedAttributeObject{
										Attributes: settingInstanceSchema(0),
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

// settingInstanceSchema returns the schema for a setting instance, used recursively for children
func settingInstanceSchema(depth int) map[string]schema.Attribute {
	if depth >= 10 {
		return nil
	}

	childAttributes := map[string]schema.Attribute{
		"odata_type": schema.StringAttribute{
			Optional:    true, // Change to Optional for children
			Description: fmt.Sprintf("The OData type for child setting at depth %d.", depth+1),
		},
		"setting_definition_id": schema.StringAttribute{
			Optional:    true, // Change to Optional for children
			Description: fmt.Sprintf("The setting definition ID for child setting at depth %d.", depth+1),
		},
		"choice_setting_value": schema.SingleNestedAttribute{
			Optional:    true,
			Description: fmt.Sprintf("The choice setting value for child level %d.", depth+1),
			Attributes: map[string]schema.Attribute{
				"odata_type": schema.StringAttribute{
					Optional:    true, // Make Optional for children within choice_setting_value
					Description: "The OData type of the choice setting value.",
				},
				"setting_value_template_reference": schema.SingleNestedAttribute{
					Optional:    true,
					Description: "The template reference for this setting value.",
					Attributes: map[string]schema.Attribute{
						"odata_type": schema.StringAttribute{
							Optional:    true, // Make Optional for nested template reference
							Description: "The OData type of the setting value template reference.",
						},
						"setting_value_template_id": schema.StringAttribute{
							Optional:    true, // Make Optional for nested template reference
							Description: "The ID of the setting value template.",
						},
						"use_template_default": schema.BoolAttribute{
							Optional:    true,
							Description: "Whether to use the template default value.",
						},
					},
				},
				"int_value": schema.Int32Attribute{
					Optional:    true,
					Description: "The integer value of the choice setting.",
				},
				"string_value": schema.StringAttribute{
					Optional:    true,
					Description: "The string value of the choice setting.",
				},
			},
		},
	}

	// Add children for the next level if not at max depth
	if depth < 9 {
		childAttributes["children"] = schema.ListNestedAttribute{
			Optional:    true,
			Description: fmt.Sprintf("The child settings of this choice setting (level %d).", depth+2),
			NestedObject: schema.NestedAttributeObject{
				Attributes: settingInstanceSchema(depth + 1),
			},
		}
	}

	return childAttributes
}

func assignmentsSchema() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Optional:    true,
		Description: "The assignment configuration for this Windows Settings Catalog profile.",
		Attributes: map[string]schema.Attribute{
			"all_devices": schema.BoolAttribute{
				Optional: true,
				MarkdownDescription: "Specifies whether this assignment applies to all devices. " +
					"When set to `true`, the assignment targets all devices in the organization." +
					"Can be used in conjuction with `all_devices_filter_type` or `all_devices_filter_id`." +
					"Can be used as an alternative to `include_groups`." +
					"Can be used in conjuction with `all_users` and `all_users_filter_type` or `all_users_filter_id`.",
			},
			"all_devices_filter_type": schema.StringAttribute{
				Optional: true,
				MarkdownDescription: "The filter type for all devices assignment. " +
					"Valid values are:\n" +
					"- `include`: Apply the assignment to devices that match the filter.\n" +
					"- `exclude`: Do not apply the assignment to devices that match the filter.\n" +
					"- `none`: No filter applied.",
				Validators: []validator.String{
					stringvalidator.OneOf("include", "exclude", "none"),
				},
			},
			"all_devices_filter_id": schema.StringAttribute{
				Optional: true,
				MarkdownDescription: "The ID of the device group filter to apply when `all_devices` is set to `true`. " +
					"This should be a valid GUID of an existing device group filter.",
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						guidRegex,
						"must be a valid GUID",
					),
				},
			},
			"all_users": schema.BoolAttribute{
				Optional: true,
				MarkdownDescription: "Specifies whether this assignment applies to all users. " +
					"When set to `true`, the assignment targets all licensed users within the organization." +
					"Can be used in conjuction with `all_users_filter_type` or `all_users_filter_id`." +
					"Can be used as an alternative to `include_groups`." +
					"Can be used in conjuction with `all_devices` and `all_devices_filter_type` or `all_devices_filter_id`.",
			},
			"all_users_filter_type": schema.StringAttribute{
				Optional: true,
				MarkdownDescription: "The filter type for all users assignment. " +
					"Valid values are:\n" +
					"- `include`: Apply the assignment to users that match the filter.\n" +
					"- `exclude`: Do not apply the assignment to users that match the filter.\n" +
					"- `none`: No filter applied.",
				Validators: []validator.String{
					stringvalidator.OneOf("include", "exclude", "none"),
				},
			},
			"all_users_filter_id": schema.StringAttribute{
				Optional: true,
				MarkdownDescription: "The ID of the filter to apply when `all_users` is set to `true`. " +
					"This should be a valid GUID of an existing filter.",
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						guidRegex,
						"assignment filer for all_users must be a valid GUID",
					),
				},
			},
			"include_groups": schema.ListNestedAttribute{
				Optional: true,
				MarkdownDescription: "A list of entra id group Id's to include in the assignment. " +
					"Each group can have its own filter type and filter ID.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"group_id": schema.StringAttribute{
							Required: true,
							MarkdownDescription: "The entra ID group ID of the group to include in the assignment. " +
								"This should be a valid GUID of an existing group.",
							Validators: []validator.String{
								stringvalidator.RegexMatches(
									guidRegex,
									"assignment include group(s) must be a valid GUID",
								),
							},
						},
						"include_groups_filter_type": schema.StringAttribute{
							Optional: true,
							MarkdownDescription: "The device group filter type for the included group. " +
								"Valid values are:\n" +
								"- `include`: Apply the assignment to group members that match the filter.\n" +
								"- `exclude`: Do not apply the assignment to group members that match the filter.\n" +
								"- `none`: No filter applied.",
							Validators: []validator.String{
								stringvalidator.OneOf("include", "exclude", "none"),
							},
						},
						"include_groups_filter_id": schema.StringAttribute{
							Optional: true,
							MarkdownDescription: "The Entra ID Group ID of the filter to apply to the included group. " +
								"This should be a valid GUID of an existing filter.",
							Validators: []validator.String{
								stringvalidator.RegexMatches(
									guidRegex,
									"assignment group filter id must be a valid GUID",
								),
							},
						},
					},
				},
			},
			"exclude_group_ids": schema.ListAttribute{
				Optional:    true,
				ElementType: types.StringType,
				MarkdownDescription: "A list of group IDs to exclude from the assignment. " +
					"These groups will not receive the assignment, even if they match other inclusion criteria.",
				Validators: []validator.List{
					listvalidator.ValueStringsAre(
						stringvalidator.RegexMatches(
							guidRegex,
							"assignment exclude group id must be a valid GUID",
						),
					),
				},
			},
		},
	}
}
