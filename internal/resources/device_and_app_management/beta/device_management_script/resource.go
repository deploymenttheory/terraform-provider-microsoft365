package graphBetaDeviceManagementScript

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/schema"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	ResourceName = "graph_beta_device_and_app_management_device_management_script"
)

var (
	// Basic resource interface (CRUD operations)
	_ resource.Resource = &DeviceManagementScriptResource{}

	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &DeviceManagementScriptResource{}

	// Enables import functionality
	_ resource.ResourceWithImportState = &DeviceManagementScriptResource{}

	// Enables plan modification/diff suppression
	_ resource.ResourceWithModifyPlan = &DeviceManagementScriptResource{}
)

func NewDeviceManagementScriptResource() resource.Resource {
	return &DeviceManagementScriptResource{
		ReadPermissions: []string{
			"DeviceManagementConfiguration.Read.All",
		},
		WritePermissions: []string{
			"DeviceManagementConfiguration.ReadWrite.All",
		},
	}
}

type DeviceManagementScriptResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ProviderTypeName string
	TypeName         string
	ReadPermissions  []string
	WritePermissions []string
}

// GetID returns the ID of a resource from the state model.
func (s *DeviceManagementScriptResourceModel) GetID() string {
	return s.ID.ValueString()
}

// GetTypeName returns the type name of the resource from the state model.
func (r *DeviceManagementScriptResource) GetTypeName() string {
	return r.TypeName
}

// Metadata returns the resource type name.
func (r *DeviceManagementScriptResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + ResourceName
}

// Configure sets the client for the resource.
func (r *DeviceManagementScriptResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = common.SetGraphBetaClientForResource(ctx, req, resp, r.TypeName)
}

// ImportState imports the resource state.
func (r *DeviceManagementScriptResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *DeviceManagementScriptResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "The resource `device_management_script` manages a device management script",
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
				Description: "The script content.",
				Required:    true,
				Sensitive:   true,
			},
			"created_date_time": schema.StringAttribute{
				Description: "The date and time the device management script was created. This property is read-only.",
				Computed:    true,
			},
			"last_modified_date_time": schema.StringAttribute{
				Description: "The date and time the device management script was last modified. This property is read-only.",
				Computed:    true,
			},
			"run_as_account": schema.StringAttribute{
				Description: "Indicates the type of execution context. Possible values are: `system`, `user`.",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("system", "user"),
				},
			},
			"enforce_signature_check": schema.BoolAttribute{
				Description: "Indicate whether the script signature needs be checked.",
				Optional:    true,
			},
			"file_name": schema.StringAttribute{
				Description: "Script file name.",
				Required:    true,
			},
			"role_scope_tag_ids": schema.ListAttribute{
				Description: "List of Scope Tag IDs for this PowerShellScript instance.",
				Optional:    true,
				ElementType: types.StringType,
			},
			"run_as_32_bit": schema.BoolAttribute{
				Description: "A value indicating whether the PowerShell script should run as 32-bit.",
				Optional:    true,
			},
			"assignments": schema.ListNestedAttribute{
				Description: "The assignments of the device management script.",
				Optional:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "Key of the device management script assignment entity. This property is read-only.",
							Computed:    true,
						},
						"target": schema.SingleNestedAttribute{
							Description: "The target of the assignment.",
							Required:    true,
							Attributes: map[string]schema.Attribute{
								"device_and_app_management_assignment_filter_id": schema.StringAttribute{
									Description: "The Id of the filter for the target assignment.",
									Optional:    true,
								},
								"device_and_app_management_assignment_filter_type": schema.StringAttribute{
									Description: "The type of filter of the target assignment i.e. Exclude or Include. Possible values are: `none`, `include`, `exclude`.",
									Optional:    true,
									Validators: []validator.String{
										stringvalidator.OneOf("none", "include", "exclude"),
									},
								},
								"target_type": schema.StringAttribute{
									Description: "The target type of the assignment.",
									Required:    true,
									Validators: []validator.String{
										stringvalidator.OneOf("user"),
									},
								},
								"entra_object_id": schema.StringAttribute{
									Description: "The ID of the Azure Active Directory object.",
									Optional:    true,
								},
							},
						},
					},
				},
			},
			"group_assignments": schema.ListNestedAttribute{
				Description: "The group assignments of the device management script.",
				Optional:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "Key of the device management script group assignment entity. This property is read-only.",
							Computed:    true,
						},
						"target_group_id": schema.StringAttribute{
							Description: "The Id of the Azure Active Directory group we are targeting the script to.",
							Required:    true,
						},
					},
				},
			},
			"timeouts": commonschema.Timeouts(ctx),
		},
	}
}
