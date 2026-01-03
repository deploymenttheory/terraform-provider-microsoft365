package graphBetaAutopatchGroups

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	planmodifiers "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/plan_modifiers"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const (
	ResourceName  = "microsoft365_graph_beta_device_management_autopatch_groups"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	// Basic resource interface (CRUD operations)
	_ resource.Resource = &AutopatchGroupsResource{}

	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &AutopatchGroupsResource{}

	// Enables import functionality
	_ resource.ResourceWithImportState = &AutopatchGroupsResource{}

	// Enables plan modification/diff suppression
	_ resource.ResourceWithModifyPlan = &AutopatchGroupsResource{}
)

func NewAutopatchGroupsResource() resource.Resource {
	return &AutopatchGroupsResource{
		ReadPermissions: []string{
			"WindowsUpdates.ReadWrite.All", // This does not work
		},
		WritePermissions: []string{
			"WindowsUpdates.ReadWrite.All", // This does not work
		},
		ResourcePath: "/device/v2/autopatchGroups",
		APIEndpoint:  "https://services.autopatch.microsoft.com",
	}
}

type AutopatchGroupsResource struct {
	httpClient       *client.AuthenticatedHTTPClient
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
	APIEndpoint      string
}

// Metadata returns the resource type name.
func (r *AutopatchGroupsResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = ResourceName
}

// Configure sets the client for the resource.
func (r *AutopatchGroupsResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.httpClient = client.SetGraphBetaHTTPClientForResource(ctx, req, resp, ResourceName)
}

// ImportState imports the resource state.
func (r *AutopatchGroupsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Schema returns the schema for the resource.
func (r *AutopatchGroupsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages Windows Autopatch groups using the `https://services.autopatch.microsoft.com/device/v2/autopatchGroups` endpoint." +
			" Autopatch groups help organize devices into logical groups for automated Windows Update deployment with customizable deployment rings and policy settings." +
			"This resource is not documented in the Microsoft Graph API documentation. This resource is experimental and may change in the future." +
			"There's currently 401 errors when trying using this resource. There appears to be a seperate unobservable authentication step between intune and " +
			" autopatch micro service that cannot be replicated in the terraform provider. Entra ID client id / secret are not sufficient to authenticate.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
				MarkdownDescription: "The unique identifier for this Autopatch group",
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The name of the Autopatch group",
			},
			"description": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The description of the Autopatch group",
			},
			"tenant_id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The tenant ID associated with this Autopatch group",
			},
			"type": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The type of the Autopatch group (Default, User)",
			},
			"status": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The status of the Autopatch group (Active, Creating, etc.)",
			},
			"is_locked_by_policy": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "Whether the group is locked by policy",
			},
			"distribution_type": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The distribution type (Mixed, AdminAssigned)",
			},
			"read_only": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "Whether the group is read-only",
			},
			"number_of_registered_devices": schema.Int32Attribute{
				Computed:            true,
				MarkdownDescription: "The number of registered devices in the group",
			},
			"user_has_all_scope_tag": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "Whether the user has all scope tags",
			},
			"flow_id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The flow ID for the operation",
			},
			"flow_type": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The flow type for the operation",
			},
			"flow_status": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The flow status for the operation",
			},
			"umbrella_group_id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The umbrella group ID",
			},
			"enable_driver_update": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Whether driver updates are enabled",
			},
			"enabled_content_types": schema.Int32Attribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Enabled content types bitmask",
			},
			"global_user_managed_aad_groups": schema.SetNestedAttribute{
				Optional:            true,
				MarkdownDescription: "Global user-managed Azure AD groups",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Required:            true,
							MarkdownDescription: "The ID of the Azure AD group",
						},
						"type": schema.StringAttribute{
							Required:            true,
							MarkdownDescription: "The type of the group (Device, User)",
						},
					},
				},
			},
			"deployment_groups": schema.ListNestedAttribute{
				Optional:            true,
				MarkdownDescription: "The deployment groups (rings) within this Autopatch group",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"aad_id": schema.StringAttribute{
							Optional:            true,
							Computed:            true,
							MarkdownDescription: "The Azure AD group ID for this deployment group",
							Default:             stringdefault.StaticString("00000000-0000-0000-0000-000000000000"),
						},
						"name": schema.StringAttribute{
							Required:            true,
							MarkdownDescription: "The name of the deployment group",
						},
						"distribution": schema.Int32Attribute{
							Optional:            true,
							MarkdownDescription: "Distribution percentage for this deployment group",
						},
						"failed_prerequisite_check_count": schema.Int32Attribute{
							Computed:            true,
							MarkdownDescription: "Number of failed prerequisite checks",
						},
						"user_managed_aad_groups": schema.SetNestedAttribute{
							Optional:            true,
							MarkdownDescription: "User-managed Azure AD groups for this deployment group",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"id": schema.StringAttribute{
										Required:            true,
										MarkdownDescription: "The ID of the Azure AD group",
									},
									"name": schema.StringAttribute{
										Optional:            true,
										MarkdownDescription: "The name of the Azure AD group",
									},
									"type": schema.StringAttribute{
										Optional:            true,
										Computed:            true,
										MarkdownDescription: "The type of the group (Device, None)",
										Default:             stringdefault.StaticString("Device"),
									},
								},
							},
						},
						"deployment_group_policy_settings": schema.SingleNestedAttribute{
							Optional:            true,
							MarkdownDescription: "Policy settings for this deployment group",
							Attributes: map[string]schema.Attribute{
								"aad_group_name": schema.StringAttribute{
									Optional:            true,
									MarkdownDescription: "The Azure AD group name",
								},
								"is_update_settings_modified": schema.BoolAttribute{
									Optional:            true,
									MarkdownDescription: "Whether update settings are modified",
								},
								"device_configuration_setting": schema.SingleNestedAttribute{
									Optional:            true,
									MarkdownDescription: "Device configuration settings",
									Attributes: map[string]schema.Attribute{
										"policy_id": schema.StringAttribute{
											Optional:            true,
											Computed:            true,
											MarkdownDescription: "The policy ID",
										},
										"update_behavior": schema.StringAttribute{
											Optional:            true,
											MarkdownDescription: "Update behavior setting",
										},
										"notification_setting": schema.StringAttribute{
											Optional:            true,
											MarkdownDescription: "Notification setting",
										},
										"quality_deployment_settings": schema.SingleNestedAttribute{
											Optional:            true,
											MarkdownDescription: "Quality update deployment settings",
											Attributes: map[string]schema.Attribute{
												"deadline": schema.Int32Attribute{
													Optional:            true,
													MarkdownDescription: "Deadline in days",
												},
												"deferral": schema.Int32Attribute{
													Optional:            true,
													MarkdownDescription: "Deferral in days",
												},
												"grace_period": schema.Int32Attribute{
													Optional:            true,
													MarkdownDescription: "Grace period in days",
												},
											},
										},
										"feature_deployment_settings": schema.SingleNestedAttribute{
											Optional:            true,
											MarkdownDescription: "Feature update deployment settings",
											Attributes: map[string]schema.Attribute{
												"deadline": schema.Int32Attribute{
													Optional:            true,
													MarkdownDescription: "Deadline in days",
												},
												"deferral": schema.Int32Attribute{
													Optional:            true,
													MarkdownDescription: "Deferral in days",
												},
											},
										},
									},
								},
								"dnf_update_cloud_setting": schema.SingleNestedAttribute{
									Optional:            true,
									MarkdownDescription: "Driver and firmware update cloud settings",
									Attributes: map[string]schema.Attribute{
										"policy_id": schema.StringAttribute{
											Optional:            true,
											Computed:            true,
											MarkdownDescription: "The policy ID",
										},
										"approval_type": schema.StringAttribute{
											Optional:            true,
											MarkdownDescription: "Approval type (Automatic or Manual)",
										},
										"deployment_deferral_in_days": schema.Int32Attribute{
											Optional:            true,
											MarkdownDescription: "Deployment deferral in days",
										},
									},
								},
								"office_dcv2_setting": schema.SingleNestedAttribute{
									Optional:            true,
									MarkdownDescription: "Office update delivery optimization v2 settings",
									Attributes: map[string]schema.Attribute{
										"policy_id": schema.StringAttribute{
											Optional:            true,
											Computed:            true,
											MarkdownDescription: "The policy ID",
										},
										"deadline": schema.Int32Attribute{
											Optional:            true,
											MarkdownDescription: "Deadline in days",
										},
										"deferral": schema.Int32Attribute{
											Optional:            true,
											MarkdownDescription: "Deferral in days",
										},
										"hide_update_notifications": schema.BoolAttribute{
											Optional:            true,
											MarkdownDescription: "Whether to hide update notifications",
										},
										"target_channel": schema.StringAttribute{
											Optional:            true,
											MarkdownDescription: "Target channel for Office updates (e.g., MonthlyEnterprise)",
										},
										"enable_automatic_update": schema.BoolAttribute{
											Optional:            true,
											Computed:            true,
											MarkdownDescription: "Whether to enable automatic updates",
										},
										"hide_enable_disable_update": schema.BoolAttribute{
											Optional:            true,
											Computed:            true,
											MarkdownDescription: "Whether to hide enable/disable update option",
										},
										"enable_office_mgmt": schema.BoolAttribute{
											Optional:            true,
											Computed:            true,
											MarkdownDescription: "Whether to enable Office management",
										},
										"update_path": schema.StringAttribute{
											Optional:            true,
											Computed:            true,
											MarkdownDescription: "The update path URL for Office updates",
										},
									},
								},
								"edge_dcv2_setting": schema.SingleNestedAttribute{
									Optional:            true,
									MarkdownDescription: "Edge update delivery optimization v2 settings",
									Attributes: map[string]schema.Attribute{
										"policy_id": schema.StringAttribute{
											Optional:            true,
											Computed:            true,
											MarkdownDescription: "The policy ID",
										},
										"target_channel": schema.StringAttribute{
											Optional:            true,
											MarkdownDescription: "Target channel for Edge updates (e.g., Stable, Beta)",
										},
									},
								},
								"feature_update_anchor_cloud_setting": schema.SingleNestedAttribute{
									Optional:            true,
									MarkdownDescription: "Feature update anchor cloud settings",
									Attributes: map[string]schema.Attribute{
										"target_os_version": schema.StringAttribute{
											Optional:            true,
											MarkdownDescription: "Target OS version (e.g., 'Windows 11, version 25H2')",
										},
										"install_latest_windows10_on_windows11_ineligible_device": schema.BoolAttribute{
											Optional:            true,
											MarkdownDescription: "Whether to install the latest Windows 10 on Windows 11 ineligible devices",
										},
										"policy_id": schema.StringAttribute{
											Optional:            true,
											Computed:            true,
											MarkdownDescription: "The policy ID",
										},
									},
								},
							},
						},
					},
				},
			},
			"scope_tags": schema.SetAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Set of scope tag IDs for this Settings Catalog template profile.",
				PlanModifiers: []planmodifier.Set{
					planmodifiers.DefaultSetValue(
						[]attr.Value{types.StringValue("0")},
					),
				},
			},
			"timeouts": commonschema.ResourceTimeouts(ctx),
		},
	}
}
