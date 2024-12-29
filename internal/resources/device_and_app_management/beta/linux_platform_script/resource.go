package graphBetaLinuxPlatformScript

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	ResourceName  = "graph_beta_device_and_app_management_linux_platform_script"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

// var (
// 	// Basic resource interface (CRUD operations)
// 	_ resource.Resource = &LinuxPlatformScriptResource{}

// 	// Allows the resource to be configured with the provider client
// 	_ resource.ResourceWithConfigure = &LinuxPlatformScriptResource{}

// 	// Enables import functionality
// 	_ resource.ResourceWithImportState = &LinuxPlatformScriptResource{}

// 	// Enables plan modification/diff suppression
// 	_ resource.ResourceWithModifyPlan = &LinuxPlatformScriptResource{}
// )

// func NewLinuxPlatformScriptResource() resource.Resource {
// 	return &LinuxPlatformScriptResource{
// 		ReadPermissions: []string{
// 			"DeviceManagementConfiguration.Read.All",
// 		},
// 		WritePermissions: []string{
// 			"DeviceManagementConfiguration.ReadWrite.All",
// 		},
// 		ResourcePath: "/deviceManagement/configurationPolicies",
// 	}
// }

type LinuxPlatformScriptResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ProviderTypeName string
	TypeName         string
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

// Metadata returns the resource type name.
func (r *LinuxPlatformScriptResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + ResourceName
}

// Configure sets the client for the resource.
func (r *LinuxPlatformScriptResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = common.SetGraphBetaClientForResource(ctx, req, resp, r.TypeName)
}

// ImportState imports the resource state.
func (r *LinuxPlatformScriptResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Schema defines the resource schema.
func (r *LinuxPlatformScriptResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages an Intune Linux platform script using the 'configurationPolicies' Graph Beta API.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Unique Identifier for the device management script.",
				Computed:    true,
			},
			"display_name": schema.StringAttribute{
				Description: "Name of the device management script.",
				Required:    true,
			},
			"description": schema.StringAttribute{
				Description: "Optional description for the device management script.",
				Optional:    true,
			},
			"script_content": schema.StringAttribute{
				Description: "The script content, base64 encoded.",
				Required:    true,
				Sensitive:   true,
			},
			"role_scope_tag_ids": schema.ListAttribute{
				Description: "List of Scope Tag IDs for this script.",
				Optional:    true,
				ElementType: types.StringType,
			},
			"platforms": schema.StringAttribute{
				Description: "Platform type for the script, always `LINUX`.",
				Computed:    true,
			},
			"technologies": schema.StringAttribute{
				Description: "Technology type for the linux platform script, always `LINUXMDM`.",
				Computed:    true,
			},
			"settings": schema.ListNestedAttribute{
				Description: "List of configuration settings.",
				Optional:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"setting_definition_id": schema.StringAttribute{
							Description: "The ID of the setting definition.",
							Required:    true,
						},
						"setting_value": schema.StringAttribute{
							Description: "The value for the setting.",
							Required:    true,
						},
						"setting_value_template_reference": schema.StringAttribute{
							Description: "The template reference ID for the setting value.",
							Optional:    true,
						},
						"setting_instance_template_id": schema.StringAttribute{
							Description: "The instance template ID for the setting.",
							Optional:    true,
						},
						"children": schema.ListNestedAttribute{
							Description: "Nested children configuration settings.",
							Optional:    true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"setting_definition_id": schema.StringAttribute{
										Description: "The ID of the nested setting definition.",
										Required:    true,
									},
									"setting_value": schema.StringAttribute{
										Description: "The value for the nested setting.",
										Required:    true,
									},
									"setting_value_template_reference": schema.StringAttribute{
										Description: "The template reference ID for the nested setting value.",
										Optional:    true,
									},
									"setting_instance_template_id": schema.StringAttribute{
										Description: "The instance template ID for the nested setting.",
										Optional:    true,
									},
								},
							},
						},
					},
				},
			},
			"template_reference_id": schema.StringAttribute{
				Description: "The ID of the configuration policy template reference.",
				Required:    true,
			},
			"assignments": commonschema.IntuneScriptAssignmentsSchema(),
			"timeouts":    commonschema.Timeouts(ctx),
		},
	}
}
