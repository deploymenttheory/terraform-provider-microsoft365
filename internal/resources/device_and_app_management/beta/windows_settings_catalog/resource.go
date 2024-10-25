package graphBetaWindowsSettingsCatalog

import (
	"context"
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
var guidRegex *regexp.Regexp

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

// GuidRegex returns the regex pattern for a GUID.
func (r *WindowsSettingsCatalogResource) GuidRegex() string {
	guidRegex = regexp.MustCompile(helpers.GuidRegex)
	return helpers.GuidRegex
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
			"settings_count": schema.Int64Attribute{
				Computed:    true,
				Description: "The number of settings in this profile.",
			},
			"name": schema.StringAttribute{
				Computed:    true,
				Description: "The name of the Windows Settings Catalog profile.",
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
			"template_reference": schema.SingleNestedAttribute{
				Required:    true,
				Description: "The template reference for this Windows Settings Catalog profile.",
				Attributes: map[string]schema.Attribute{
					"template_id": schema.StringAttribute{
						Required:    true,
						Description: "The ID of the template used for this profile.",
					},
					"template_display_name": schema.StringAttribute{
						Computed:    true,
						Description: "The display name of the template used for this profile.",
					},
					"template_display_version": schema.StringAttribute{
						Computed:    true,
						Description: "The display version of the template used for this profile.",
					},
				},
			},

			// Settings attributes
			"settings": settingsSchema(),

			"assignments": schema.ListNestedAttribute{
				Optional:    true,
				Description: "The list of group assignments for this Windows Settings Catalog profile.",
				NestedObject: schema.NestedAttributeObject{
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
									"must be a valid GUID",
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
												"must be a valid GUID",
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
												"must be a valid GUID",
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
										"must be a valid GUID",
									),
								),
							},
						},
					},
				},
			},
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
									Description: "The value template reference for the choice setting.",
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
								"value": schema.StringAttribute{
									Required:    true,
									Description: "The value of the choice setting.",
								},
								"children": schema.ListNestedAttribute{
									Optional:    true,
									Description: "The child settings of this choice setting.",
									NestedObject: schema.NestedAttributeObject{
										Attributes: map[string]schema.Attribute{
											// This is a recursive reference to the setting_instance structure
											// We'll define it separately and reference it here
										},
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

// Define the recursive setting_instance structure
func settingInstanceSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
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
					Description: "The value template reference for the choice setting.",
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
				"value": schema.StringAttribute{
					Required:    true,
					Description: "The value of the choice setting.",
				},
				"children": schema.ListNestedAttribute{
					Optional:    true,
					Description: "The child settings of this choice setting.",
					NestedObject: schema.NestedAttributeObject{
						Attributes: settingInstanceSchema(),
					},
				},
			},
		},
	}
}
