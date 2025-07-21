package graphIOSMobileAppConfiguration

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	planmodifiers "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/plan_modifiers"
)

const (
	ResourceName  = "graph_v1_device_and_app_management_ios_mobile_app_configuration"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	_ resource.Resource                = &IOSMobileAppConfigurationResource{}
	_ resource.ResourceWithConfigure   = &IOSMobileAppConfigurationResource{}
	_ resource.ResourceWithImportState = &IOSMobileAppConfigurationResource{}
	_ resource.ResourceWithModifyPlan  = &IOSMobileAppConfigurationResource{}
)

type IOSMobileAppConfigurationResource struct {
	client           *msgraphsdk.GraphServiceClient
	ProviderTypeName string
	TypeName         string
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

func NewIOSMobileAppConfigurationResource() resource.Resource {
	return &IOSMobileAppConfigurationResource{
		ReadPermissions: []string{
			"DeviceManagementApps.Read.All",
		},
		WritePermissions: []string{
			"DeviceManagementApps.ReadWrite.All",
		},
		ResourcePath: "/deviceAppManagement/mobileAppConfigurations",
	}
}

func (r *IOSMobileAppConfigurationResource) Metadata(
	ctx context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_" + ResourceName
}

func (r *IOSMobileAppConfigurationResource) Configure(
	ctx context.Context,
	req resource.ConfigureRequest,
	resp *resource.ConfigureResponse,
) {
	r.client = client.SetGraphStableClientForResource(
		ctx,
		req,
		resp,
		constants.PROVIDER_NAME+"_"+ResourceName,
	)
}

func (r *IOSMobileAppConfigurationResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *IOSMobileAppConfigurationResource) Schema(
	ctx context.Context,
	req resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages iOS mobile app configuration policies in Microsoft Intune. These policies allow you to configure app-specific settings for managed iOS applications.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique identifier for the iOS mobile app configuration. Read-only.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"display_name": schema.StringAttribute{
				MarkdownDescription: "The display name of the iOS mobile app configuration.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 1000),
				},
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The description of the iOS mobile app configuration.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(0, 1000),
				},
			},
			"targeted_mobile_apps": schema.ListAttribute{
				MarkdownDescription: "The list of targeted mobile app IDs.",
				Optional:            true,
				ElementType:         types.StringType,
				Validators: []validator.List{
					listvalidator.SizeAtLeast(1),
				},
			},
			"encoded_setting_xml": schema.StringAttribute{
				MarkdownDescription: "Base64 encoded configuration XML.",
				Optional:            true,
				Sensitive:           true,
			},
			"created_date_time": schema.StringAttribute{
				MarkdownDescription: "DateTime the object was created. Read-only.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
			},
			"last_modified_date_time": schema.StringAttribute{
				MarkdownDescription: "DateTime the object was last modified. Read-only.",
				Computed:            true,
			},
			"version": schema.Int32Attribute{
				MarkdownDescription: "Version of the device configuration. Read-only.",
				Computed:            true,
			},
			"timeouts": timeouts.Attributes(ctx, timeouts.Opts{
				Create: true,
				Read:   true,
				Update: true,
				Delete: true,
			}),
		},
		Blocks: map[string]schema.Block{
			"settings": schema.ListNestedBlock{
				MarkdownDescription: "iOS app configuration settings.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"app_config_key": schema.StringAttribute{
							MarkdownDescription: "The application configuration key.",
							Required:            true,
						},
						"app_config_key_type": schema.StringAttribute{
							MarkdownDescription: "The application configuration key type. Possible values are: stringType, integerType, realType, booleanType, tokenType.",
							Required:            true,
							Validators: []validator.String{
								stringvalidator.OneOf(
									"stringType",
									"integerType",
									"realType",
									"booleanType",
									"tokenType",
								),
							},
						},
						"app_config_key_value": schema.StringAttribute{
							MarkdownDescription: "The application configuration key value.",
							Required:            true,
						},
					},
				},
			},
			"assignments": schema.ListNestedBlock{
				MarkdownDescription: "The list of assignments for this iOS mobile app configuration.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							MarkdownDescription: "Key of the entity. Read-only.",
							Computed:            true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.UseStateForUnknown(),
							},
						},
					},
					Blocks: map[string]schema.Block{
						"target": schema.SingleNestedBlock{
							MarkdownDescription: "Target for the assignment.",
							Attributes: map[string]schema.Attribute{
								"odata_type": schema.StringAttribute{
									MarkdownDescription: "The type of assignment target. Possible values are: #microsoft.graph.allLicensedUsersAssignmentTarget, #microsoft.graph.allDevicesAssignmentTarget, #microsoft.graph.exclusionGroupAssignmentTarget, #microsoft.graph.groupAssignmentTarget.",
									Required:            true,
									Validators: []validator.String{
										stringvalidator.OneOf(
											"#microsoft.graph.allLicensedUsersAssignmentTarget",
											"#microsoft.graph.allDevicesAssignmentTarget",
											"#microsoft.graph.exclusionGroupAssignmentTarget",
											"#microsoft.graph.groupAssignmentTarget",
										),
									},
								},
								"group_id": schema.StringAttribute{
									MarkdownDescription: "The group Id that is the target of the assignment. Required when odata_type is groupAssignmentTarget or exclusionGroupAssignmentTarget.",
									Optional:            true,
									Validators: []validator.String{
										stringvalidator.LengthBetween(36, 36),
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

func (r *IOSMobileAppConfigurationResource) ModifyPlan(
	ctx context.Context,
	req resource.ModifyPlanRequest,
	resp *resource.ModifyPlanResponse,
) {
	if req.State.Raw.IsNull() || req.Plan.Raw.IsNull() {
		return
	}

	var state IOSMobileAppConfigurationResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var plan IOSMobileAppConfigurationResourceModel
	diags = req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if len(state.Assignments) > 0 && len(plan.Assignments) > 0 {
		for i := range plan.Assignments {
			if i < len(state.Assignments) && !state.Assignments[i].Id.IsNull() &&
				!state.Assignments[i].Id.IsUnknown() {
				plan.Assignments[i].Id = state.Assignments[i].Id
			}
		}
		resp.Plan.Set(ctx, &plan)
	}
}
