package graphBetaWindowsSettingsCatalog

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/schema"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
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
	_ resource.Resource = &WindowsSettingsCatalogResource{}

	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &WindowsSettingsCatalogResource{}

	// Enables import functionality
	_ resource.ResourceWithImportState = &WindowsSettingsCatalogResource{}

	// Enables plan modification/diff suppression
	_ resource.ResourceWithModifyPlan = &WindowsSettingsCatalogResource{}
)

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
			// Policy base attributes
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The unique identifier for this Windows Settings Catalog profile.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
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
			"settings": settingsSchema(),

			// Assignments
			"assignments": commonschema.AssignmentsSchema(),

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
						"choice_setting_value": schema.SingleNestedAttribute{
							Optional:    true,
							Description: "The choice setting value.",
							Attributes: map[string]schema.Attribute{
								"odata_type": schema.StringAttribute{
									Required:    true,
									Description: "The OData type of the choice setting value.",
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
										Attributes: settingSchema(0),
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
func settingSchema(depth int) map[string]schema.Attribute {
	if depth >= 10 {
		return nil
	}

	childAttributes := map[string]schema.Attribute{
		"odata_type": schema.StringAttribute{
			Optional:    true,
			Description: fmt.Sprintf("The OData type for child setting at depth %d.", depth+1),
		},
		"setting_definition_id": schema.StringAttribute{
			Optional:    true,
			Description: fmt.Sprintf("The setting definition ID for child setting at depth %d.", depth+1),
		},
		"choice_setting_value": schema.SingleNestedAttribute{
			Optional:    true,
			Description: fmt.Sprintf("The choice setting value for child level %d.", depth+1),
			Attributes: map[string]schema.Attribute{
				"odata_type": schema.StringAttribute{
					Optional: true,
				},
				"int_value": schema.Int32Attribute{
					Optional: true,
				},
				"string_value": schema.StringAttribute{
					Optional: true,
				},
				"children": schema.ListNestedAttribute{
					Optional:    true,
					Description: fmt.Sprintf("The child settings of this choice setting (level %d).", depth+2),
					NestedObject: schema.NestedAttributeObject{
						Attributes: settingSchema(depth + 1),
					},
				},
			},
		},
	}

	return childAttributes
}
