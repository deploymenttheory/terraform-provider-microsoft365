package graphBetaConditionalAccessTemplate

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	DataSourceName = "microsoft365_graph_beta_identity_and_access_conditional_access_template"
	ReadTimeout    = 180
)

var (
	_ datasource.DataSource              = &ConditionalAccessTemplateDataSource{}
	_ datasource.DataSourceWithConfigure = &ConditionalAccessTemplateDataSource{}
)

func NewConditionalAccessTemplateDataSource() datasource.DataSource {
	return &ConditionalAccessTemplateDataSource{
		ReadPermissions: []string{
			"Policy.Read.All",
		},
	}
}

type ConditionalAccessTemplateDataSource struct {
	client          *msgraphbetasdk.GraphServiceClient
	ReadPermissions []string
}

func (d *ConditionalAccessTemplateDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = DataSourceName
}

func (d *ConditionalAccessTemplateDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.client = client.SetGraphBetaClientForDataSource(ctx, req, resp, DataSourceName)
}

func (d *ConditionalAccessTemplateDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Retrieves conditional access policy templates from Microsoft Entra ID. " +
			"Templates provide pre-configured conditional access policies for common security scenarios. " +
			"You can query templates by ID or name to retrieve template details including conditions, grant controls, and session controls.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The unique identifier for this data source operation.",
			},
			"template_id": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The unique identifier (GUID) of the conditional access template.",
				Validators: []validator.String{
					stringvalidator.AtLeastOneOf(path.MatchRoot("name")),
					stringvalidator.ConflictsWith(path.MatchRoot("name")),
				},
			},
			"name": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The name of the conditional access template.",
				Validators: []validator.String{
					stringvalidator.AtLeastOneOf(path.MatchRoot("template_id")),
					stringvalidator.ConflictsWith(path.MatchRoot("template_id")),
				},
			},
			"description": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Description of what the conditional access template does.",
			},
			"scenarios": schema.SetAttribute{
				ElementType:         types.StringType,
				Computed:            true,
				MarkdownDescription: "Set of scenarios this template applies to (e.g., secureFoundation, zeroTrust, remoteWork, protectAdmins).",
			},
			"details": schema.SingleNestedAttribute{
				Computed:            true,
				MarkdownDescription: "The policy configuration details including conditions, grant controls, and session controls.",
				Attributes: map[string]schema.Attribute{
					"conditions": schema.SingleNestedAttribute{
						Computed:            true,
						MarkdownDescription: "The conditions that must be met for the policy to apply.",
						Attributes: map[string]schema.Attribute{
							"user_risk_levels": schema.ListAttribute{
								Computed:            true,
								ElementType:         types.StringType,
								MarkdownDescription: "User risk levels included in the policy.",
							},
							"sign_in_risk_levels": schema.ListAttribute{
								Computed:            true,
								ElementType:         types.StringType,
								MarkdownDescription: "Sign-in risk levels included in the policy.",
							},
							"client_app_types": schema.ListAttribute{
								Computed:            true,
								ElementType:         types.StringType,
								MarkdownDescription: "Client application types included in the policy.",
							},
							"service_principal_risk_levels": schema.ListAttribute{
								Computed:            true,
								ElementType:         types.StringType,
								MarkdownDescription: "Service principal risk levels included in the policy.",
							},
							"agent_id_risk_levels": schema.SetAttribute{
								Computed:            true,
								ElementType:         types.StringType,
								MarkdownDescription: "Agent identity risk levels included in the policy.",
							},
							"insider_risk_levels": schema.SetAttribute{
								Computed:            true,
								ElementType:         types.StringType,
								MarkdownDescription: "Insider risk levels included in the policy.",
							},
							"platforms": schema.SingleNestedAttribute{
								Computed:            true,
								MarkdownDescription: "Platform conditions.",
								Attributes: map[string]schema.Attribute{
									"include_platforms": schema.ListAttribute{
										Computed:            true,
										ElementType:         types.StringType,
										MarkdownDescription: "Platforms included in the policy.",
									},
									"exclude_platforms": schema.ListAttribute{
										Computed:            true,
										ElementType:         types.StringType,
										MarkdownDescription: "Platforms excluded from the policy.",
									},
								},
							},
							"locations": schema.SingleNestedAttribute{
								Computed:            true,
								MarkdownDescription: "Location conditions.",
								Attributes: map[string]schema.Attribute{
									"include_locations": schema.ListAttribute{
										Computed:            true,
										ElementType:         types.StringType,
										MarkdownDescription: "Locations included in the policy.",
									},
									"exclude_locations": schema.ListAttribute{
										Computed:            true,
										ElementType:         types.StringType,
										MarkdownDescription: "Locations excluded from the policy.",
									},
								},
							},
							"devices": schema.SingleNestedAttribute{
								Computed:            true,
								MarkdownDescription: "Device conditions.",
								Attributes: map[string]schema.Attribute{
									"include_device_states": schema.ListAttribute{
										Computed:            true,
										ElementType:         types.StringType,
										MarkdownDescription: "Device states included in the policy.",
									},
									"exclude_device_states": schema.ListAttribute{
										Computed:            true,
										ElementType:         types.StringType,
										MarkdownDescription: "Device states excluded from the policy.",
									},
									"include_devices": schema.ListAttribute{
										Computed:            true,
										ElementType:         types.StringType,
										MarkdownDescription: "Devices included in the policy.",
									},
									"exclude_devices": schema.ListAttribute{
										Computed:            true,
										ElementType:         types.StringType,
										MarkdownDescription: "Devices excluded from the policy.",
									},
									"device_filter": schema.SingleNestedAttribute{
										Computed:            true,
										MarkdownDescription: "Device filter configuration.",
										Attributes: map[string]schema.Attribute{
											"mode": schema.StringAttribute{
												Computed:            true,
												MarkdownDescription: "Filter mode (include or exclude).",
											},
											"rule": schema.StringAttribute{
												Computed:            true,
												MarkdownDescription: "Filter rule expression.",
											},
										},
									},
								},
							},
							"client_applications": schema.SingleNestedAttribute{
								Computed:            true,
								MarkdownDescription: "Client application conditions.",
								Attributes: map[string]schema.Attribute{
									"include_service_principals": schema.ListAttribute{
										Computed:            true,
										ElementType:         types.StringType,
										MarkdownDescription: "Service principals included in the policy.",
									},
									"include_agent_id_service_principals": schema.ListAttribute{
										Computed:            true,
										ElementType:         types.StringType,
										MarkdownDescription: "Agent identity service principals included in the policy.",
									},
									"exclude_service_principals": schema.ListAttribute{
										Computed:            true,
										ElementType:         types.StringType,
										MarkdownDescription: "Service principals excluded from the policy.",
									},
									"exclude_agent_id_service_principals": schema.ListAttribute{
										Computed:            true,
										ElementType:         types.StringType,
										MarkdownDescription: "Agent identity service principals excluded from the policy.",
									},
								},
							},
							"applications": schema.SingleNestedAttribute{
								Computed:            true,
								MarkdownDescription: "Application conditions.",
								Attributes: map[string]schema.Attribute{
									"include_applications": schema.ListAttribute{
										Computed:            true,
										ElementType:         types.StringType,
										MarkdownDescription: "Applications included in the policy.",
									},
									"exclude_applications": schema.ListAttribute{
										Computed:            true,
										ElementType:         types.StringType,
										MarkdownDescription: "Applications excluded from the policy.",
									},
									"include_user_actions": schema.ListAttribute{
										Computed:            true,
										ElementType:         types.StringType,
										MarkdownDescription: "User actions included in the policy.",
									},
									"include_authentication_context_class_references": schema.ListAttribute{
										Computed:            true,
										ElementType:         types.StringType,
										MarkdownDescription: "Authentication context class references included in the policy.",
									},
								},
							},
							"users": schema.SingleNestedAttribute{
								Computed:            true,
								MarkdownDescription: "User and group conditions.",
								Attributes: map[string]schema.Attribute{
									"include_users": schema.ListAttribute{
										Computed:            true,
										ElementType:         types.StringType,
										MarkdownDescription: "Users included in the policy.",
									},
									"exclude_users": schema.ListAttribute{
										Computed:            true,
										ElementType:         types.StringType,
										MarkdownDescription: "Users excluded from the policy.",
									},
									"include_groups": schema.ListAttribute{
										Computed:            true,
										ElementType:         types.StringType,
										MarkdownDescription: "Groups included in the policy.",
									},
									"exclude_groups": schema.ListAttribute{
										Computed:            true,
										ElementType:         types.StringType,
										MarkdownDescription: "Groups excluded from the policy.",
									},
									"include_roles": schema.ListAttribute{
										Computed:            true,
										ElementType:         types.StringType,
										MarkdownDescription: "Roles included in the policy.",
									},
									"exclude_roles": schema.ListAttribute{
										Computed:            true,
										ElementType:         types.StringType,
										MarkdownDescription: "Roles excluded from the policy.",
									},
									"include_guests_or_external_users": schema.SingleNestedAttribute{
										Computed:            true,
										MarkdownDescription: "Guest or external user inclusion conditions.",
										Attributes: map[string]schema.Attribute{
											"guest_or_external_user_types": schema.SetAttribute{
												Computed:            true,
												ElementType:         types.StringType,
												MarkdownDescription: "Types of guest or external users.",
											},
											"external_tenants": schema.SingleNestedAttribute{
												Computed:            true,
												MarkdownDescription: "External tenant configuration.",
												Attributes: map[string]schema.Attribute{
													"membership_kind": schema.StringAttribute{
														Computed:            true,
														MarkdownDescription: "Membership kind (e.g., all).",
													},
												},
											},
										},
									},
									"exclude_guests_or_external_users": schema.SingleNestedAttribute{
										Computed:            true,
										MarkdownDescription: "Guest or external user exclusion conditions.",
										Attributes: map[string]schema.Attribute{
											"guest_or_external_user_types": schema.SetAttribute{
												Computed:            true,
												ElementType:         types.StringType,
												MarkdownDescription: "Types of guest or external users.",
											},
											"external_tenants": schema.SingleNestedAttribute{
												Computed:            true,
												MarkdownDescription: "External tenant configuration.",
												Attributes: map[string]schema.Attribute{
													"membership_kind": schema.StringAttribute{
														Computed:            true,
														MarkdownDescription: "Membership kind (e.g., all).",
													},
												},
											},
										},
									},
								},
							},
						},
					},
					"grant_controls": schema.SingleNestedAttribute{
						Computed:            true,
						MarkdownDescription: "The grant controls applied by the policy.",
						Attributes: map[string]schema.Attribute{
							"operator": schema.StringAttribute{
								Computed:            true,
								MarkdownDescription: "Logical operator for grant controls (AND or OR).",
							},
							"built_in_controls": schema.ListAttribute{
								Computed:            true,
								ElementType:         types.StringType,
								MarkdownDescription: "Built-in grant controls (e.g., mfa, compliantDevice, domainJoinedDevice).",
							},
							"custom_authentication_factors": schema.ListAttribute{
								Computed:            true,
								ElementType:         types.StringType,
								MarkdownDescription: "Custom authentication factors.",
							},
							"terms_of_use": schema.ListAttribute{
								Computed:            true,
								ElementType:         types.StringType,
								MarkdownDescription: "Terms of use.",
							},
							"authentication_strength": schema.SingleNestedAttribute{
								Computed:            true,
								MarkdownDescription: "Authentication strength requirements.",
								Attributes: map[string]schema.Attribute{
									"id": schema.StringAttribute{
										Computed:            true,
										MarkdownDescription: "The unique identifier of the authentication strength policy.",
									},
									"created_date_time": schema.StringAttribute{
										Computed:            true,
										MarkdownDescription: "The date and time when the authentication strength policy was created.",
									},
									"modified_date_time": schema.StringAttribute{
										Computed:            true,
										MarkdownDescription: "The date and time when the authentication strength policy was last modified.",
									},
									"display_name": schema.StringAttribute{
										Computed:            true,
										MarkdownDescription: "The display name of the authentication strength policy.",
									},
									"description": schema.StringAttribute{
										Computed:            true,
										MarkdownDescription: "The description of the authentication strength policy.",
									},
									"policy_type": schema.StringAttribute{
										Computed:            true,
										MarkdownDescription: "The type of the authentication strength policy (e.g., builtIn).",
									},
									"requirements_satisfied": schema.StringAttribute{
										Computed:            true,
										MarkdownDescription: "The requirements satisfied by this authentication strength (e.g., mfa).",
									},
									"allowed_combinations": schema.ListAttribute{
										Computed:            true,
										ElementType:         types.StringType,
										MarkdownDescription: "The allowed authentication method combinations.",
									},
								},
							},
						},
					},
					"session_controls": schema.SingleNestedAttribute{
						Computed:            true,
						MarkdownDescription: "The session controls applied by the policy.",
						Attributes: map[string]schema.Attribute{
							"disable_resilience_defaults": schema.BoolAttribute{
								Computed:            true,
								MarkdownDescription: "Session control that determines whether it's acceptable for Microsoft Entra ID to extend existing sessions based on information collected prior to an outage or not.",
							},
							"application_enforced_restrictions": schema.SingleNestedAttribute{
								Computed:            true,
								MarkdownDescription: "Application enforced restrictions.",
								Attributes: map[string]schema.Attribute{
									"is_enabled": schema.BoolAttribute{
										Computed:            true,
										MarkdownDescription: "Whether application enforced restrictions are enabled.",
									},
								},
							},
							"cloud_app_security": schema.SingleNestedAttribute{
								Computed:            true,
								MarkdownDescription: "Session control to apply cloud app security.",
								Attributes: map[string]schema.Attribute{
									"is_enabled": schema.BoolAttribute{
										Computed:            true,
										MarkdownDescription: "Specifies whether the session control is enabled.",
									},
									"cloud_app_security_type": schema.StringAttribute{
										Computed:            true,
										MarkdownDescription: "The possible values are: mcasConfigured, monitorOnly, blockDownloads.",
									},
								},
							},
							"persistent_browser": schema.SingleNestedAttribute{
								Computed:            true,
								MarkdownDescription: "Persistent browser session settings.",
								Attributes: map[string]schema.Attribute{
									"mode": schema.StringAttribute{
										Computed:            true,
										MarkdownDescription: "Persistent browser mode (e.g., always, never).",
									},
									"is_enabled": schema.BoolAttribute{
										Computed:            true,
										MarkdownDescription: "Whether persistent browser session is enabled.",
									},
								},
							},
							"continuous_access_evaluation": schema.SingleNestedAttribute{
								Computed:            true,
								MarkdownDescription: "Session control for continuous access evaluation settings.",
								Attributes: map[string]schema.Attribute{
									"mode": schema.StringAttribute{
										Computed:            true,
										MarkdownDescription: "Specifies continuous access evaluation settings. The possible values are: strictEnforcement, disabled, unknownFutureValue, strictLocation.",
									},
								},
							},
							"secure_sign_in_session": schema.SingleNestedAttribute{
								Computed:            true,
								MarkdownDescription: "Session control to require sign in sessions to be bound to a device.",
								Attributes: map[string]schema.Attribute{
									"is_enabled": schema.BoolAttribute{
										Computed:            true,
										MarkdownDescription: "Specifies whether the session control is enabled.",
									},
								},
							},
							"global_secure_access_filtering_profile": schema.SingleNestedAttribute{
								Computed:            true,
								MarkdownDescription: "Session control to link to Global Secure Access security profiles or filtering profiles.",
								Attributes: map[string]schema.Attribute{
									"is_enabled": schema.BoolAttribute{
										Computed:            true,
										MarkdownDescription: "Specifies whether the session control is enabled.",
									},
									"profile_id": schema.StringAttribute{
										Computed:            true,
										MarkdownDescription: "Specifies the distinct identifier that is assigned to the security profile or filtering profile.",
									},
								},
							},
							"sign_in_frequency": schema.SingleNestedAttribute{
								Computed:            true,
								MarkdownDescription: "Sign-in frequency settings.",
								Attributes: map[string]schema.Attribute{
									"value": schema.Int64Attribute{
										Computed:            true,
										MarkdownDescription: "The frequency value.",
									},
									"type": schema.StringAttribute{
										Computed:            true,
										MarkdownDescription: "The frequency type (e.g., hours, days).",
									},
									"authentication_type": schema.StringAttribute{
										Computed:            true,
										MarkdownDescription: "The authentication type this frequency applies to.",
									},
									"frequency_interval": schema.StringAttribute{
										Computed:            true,
										MarkdownDescription: "The frequency interval (e.g., timeBased, everyTime).",
									},
									"is_enabled": schema.BoolAttribute{
										Computed:            true,
										MarkdownDescription: "Whether sign-in frequency is enabled.",
									},
								},
							},
						},
					},
				},
			},
			"timeouts": commonschema.ResourceTimeouts(ctx),
		},
	}
}
