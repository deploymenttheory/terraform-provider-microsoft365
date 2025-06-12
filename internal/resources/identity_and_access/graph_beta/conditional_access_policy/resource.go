package graphBetaConditionalAccessPolicy

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common"
	planmodifiers "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/plan_modifiers"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/schema"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	ResourceName  = "graph_beta_conditional_access_policy"
	CreateTimeout = 600
	UpdateTimeout = 600
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	// Basic resource interface (CRUD operations)
	_ resource.Resource = &ConditionalAccessPolicyResource{}

	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &ConditionalAccessPolicyResource{}

	// Enables import functionality
	_ resource.ResourceWithImportState = &ConditionalAccessPolicyResource{}

	// Enables plan modification/diff suppression
	_ resource.ResourceWithModifyPlan = &ConditionalAccessPolicyResource{}
)

func NewConditionalAccessPolicyResource() resource.Resource {
	return &ConditionalAccessPolicyResource{
		ReadPermissions: []string{
			"Policy.Read.All",
		},
		WritePermissions: []string{
			"Policy.ReadWrite.ConditionalAccess",
		},
		ResourcePath: "/identity/conditionalAccess/policies",
	}
}

type ConditionalAccessPolicyResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ProviderTypeName string
	TypeName         string
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

// Metadata returns the resource type name.
func (r *ConditionalAccessPolicyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	r.ProviderTypeName = req.ProviderTypeName
	r.TypeName = ResourceName
	resp.TypeName = r.FullTypeName()
}

// FullTypeName returns the full type name of the resource for logging purposes.
func (r *ConditionalAccessPolicyResource) FullTypeName() string {
	return r.ProviderTypeName + "_" + r.TypeName
}

// Configure sets the client for the resource.
func (r *ConditionalAccessPolicyResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = common.SetGraphBetaClientForResource(ctx, req, resp, constants.PROVIDER_NAME+"_"+ResourceName)
}

// ImportState imports the resource state.
func (r *ConditionalAccessPolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *ConditionalAccessPolicyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a Conditional Access Policy using the `/identity/conditionalAccess/policies` endpoint. Represents a Microsoft Entra Conditional Access policy. Conditional access policies are custom rules that define an access scenario.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
				MarkdownDescription: "Specifies the identifier of a conditionalAccessPolicy object. Read-only.",
			},
			"display_name": schema.StringAttribute{
				MarkdownDescription: "Specifies a display name for the conditionalAccessPolicy object.",
				Required:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Not used.",
				Optional:            true,
			},
			"state": schema.StringAttribute{
				MarkdownDescription: "Specifies the state of the conditionalAccessPolicy object. Possible values are: `enabled`, `disabled`, `enabledForReportingButNotEnforced`. Required.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("enabled", "disabled", "enabledForReportingButNotEnforced"),
				},
			},
			"created_date_time": schema.StringAttribute{
				MarkdownDescription: "The Timestamp type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is `2014-01-01T00:00:00Z`. Readonly.",
				Computed:            true,
			},
			"modified_date_time": schema.StringAttribute{
				MarkdownDescription: "The Timestamp type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is `2014-01-01T00:00:00Z`. Readonly.",
				Computed:            true,
			},
			"conditions": schema.SingleNestedAttribute{
				MarkdownDescription: "Specifies the rules that must be met for the policy to apply. Required.",
				Required:            true,
				Attributes: map[string]schema.Attribute{
					"applications": schema.SingleNestedAttribute{
						MarkdownDescription: "Applications and user actions included in and excluded from the policy. Required.",
						Optional:            true,
						Attributes: map[string]schema.Attribute{
							"application_filter": schema.SingleNestedAttribute{
								MarkdownDescription: "Filter that defines the dynamic-application-syntax rule to include/exclude cloud applications. A filter can use custom security attributes to include/exclude applications.",
								Optional:            true,
								Attributes: map[string]schema.Attribute{
									"mode": schema.StringAttribute{
										MarkdownDescription: "Mode to use for the filter. Possible values are `include` or `exclude`.",
										Required:            true,
										Validators: []validator.String{
											stringvalidator.OneOf("include", "exclude"),
										},
									},
									"rule": schema.StringAttribute{
										MarkdownDescription: "Rule syntax is similar to that used for membership rules for groups in Microsoft Entra ID.",
										Required:            true,
									},
								},
							},
							"exclude_applications": schema.SetAttribute{
								ElementType: types.StringType,
								Optional:    true,
								MarkdownDescription: "Can be one of the following:\n" +
									"- The list of client IDs (**appId**) explicitly excluded from the policy.\n" +
									"- `Office365` - For the list of apps included in `Office365`, see Apps included in Conditional Access Office 365 app suite\n" +
									"- `MicrosoftAdminPortals` - For more information, see Conditional Access Target resources: Microsoft Admin Portals",
								Default:  setdefault.StaticValue(types.SetValueMust(types.StringType, []attr.Value{})),
								Computed: true,
							},
							"include_applications": schema.SetAttribute{
								ElementType: types.StringType,
								Optional:    true,
								MarkdownDescription: "Can be one of the following:\n" +
									"- The list of client IDs (**appId**) the policy applies to, unless explicitly excluded (in **excludeApplications**)\n" +
									"- `All`\n" +
									"- `Office365` - For the list of apps included in `Office365`, see Apps included in Conditional Access Office 365 app suite\n" +
									"- `MicrosoftAdminPortals` - For more information, see Conditional Access Target resources: Microsoft Admin Portals",
								Default:  setdefault.StaticValue(types.SetValueMust(types.StringType, []attr.Value{})),
								Computed: true,
							},
							"include_user_actions": schema.SetAttribute{
								ElementType:         types.StringType,
								Optional:            true,
								MarkdownDescription: "User actions to include. Supported values are `urn:user:registersecurityinfo` and `urn:user:registerdevice`",
								Default:             setdefault.StaticValue(types.SetValueMust(types.StringType, []attr.Value{})),
								Computed:            true,
							},
							"include_authentication_context_class_references": schema.SetAttribute{
								ElementType:         types.StringType,
								Optional:            true,
								MarkdownDescription: "Authentication context class references include. Supported values are `c1` through `c25`.",
								Default:             setdefault.StaticValue(types.SetValueMust(types.StringType, []attr.Value{})),
								Computed:            true,
							},
						},
					},
					"authentication_flows": schema.SingleNestedAttribute{
						MarkdownDescription: "Authentication flows included in the policy scope. For more information, see Conditional Access: Authentication flows.",
						Optional:            true,
						Attributes: map[string]schema.Attribute{
							"transfer_methods": schema.SetAttribute{
								ElementType:         types.StringType,
								Optional:            true,
								MarkdownDescription: "Represents the transfer methods in scope for the policy. The possible values are: `none`, `deviceCodeFlow`, `authenticationTransfer`, `unknownFutureValue`.",
								Default:             setdefault.StaticValue(types.SetValueMust(types.StringType, []attr.Value{})),
								Computed:            true,
							},
						},
					},
					"users": schema.SingleNestedAttribute{
						MarkdownDescription: "Users, groups, and roles included in and excluded from the policy. Either **users** or **clientApplications** is required.",
						Optional:            true,
						Attributes: map[string]schema.Attribute{
							"exclude_groups": schema.SetAttribute{
								ElementType:         types.StringType,
								Optional:            true,
								MarkdownDescription: "Group IDs excluded from scope of policy.",
								Default:             setdefault.StaticValue(types.SetValueMust(types.StringType, []attr.Value{})),
								Computed:            true,
							},
							"exclude_guests_or_external_users": schema.SingleNestedAttribute{
								MarkdownDescription: "Internal guests or external users excluded from the policy scope. Optionally populated.",
								Optional:            true,
								Attributes: map[string]schema.Attribute{
									"external_tenants": schema.SingleNestedAttribute{
										MarkdownDescription: "The tenant IDs of the selected types of external users. Either all B2B tenant or a collection of tenant IDs. External tenants can be specified only when the property **guestOrExternalUserTypes** isn't `null` or an empty String.",
										Optional:            true,
										Attributes: map[string]schema.Attribute{
											"membership_kind": schema.StringAttribute{
												MarkdownDescription: "The membership kind. Possible values are: `all`, `enumerated`, `unknownFutureValue`.",
												Required:            true,
												Validators: []validator.String{
													stringvalidator.OneOf("all", "enumerated", "unknownFutureValue"),
												},
											},
											"members": schema.SetAttribute{
												ElementType:         types.StringType,
												Optional:            true,
												MarkdownDescription: "The tenant IDs of the external tenants.",
												Default:             setdefault.StaticValue(types.SetValueMust(types.StringType, []attr.Value{})),
												Computed:            true,
											},
										},
									},
									"guest_or_external_user_types": schema.SetAttribute{
										ElementType:         types.StringType,
										Optional:            true,
										MarkdownDescription: "Indicates internal guests or external user types, and is a multi-valued property. Possible values are: `none`, `internalGuest`, `b2bCollaborationGuest`, `b2bCollaborationMember`, `b2bDirectConnectUser`, `otherExternalUser`, `serviceProvider`, `unknownFutureValue`.",
										Default:             setdefault.StaticValue(types.SetValueMust(types.StringType, []attr.Value{})),
										Computed:            true,
									},
								},
							},
							"exclude_roles": schema.SetAttribute{
								ElementType:         types.StringType,
								Optional:            true,
								MarkdownDescription: "Role IDs excluded from scope of policy.",
								Default:             setdefault.StaticValue(types.SetValueMust(types.StringType, []attr.Value{})),
								Computed:            true,
							},
							"exclude_users": schema.SetAttribute{
								ElementType:         types.StringType,
								Optional:            true,
								MarkdownDescription: "User IDs excluded from scope of policy and/or `GuestsOrExternalUsers`.",
								Default:             setdefault.StaticValue(types.SetValueMust(types.StringType, []attr.Value{})),
								Computed:            true,
							},
							"include_groups": schema.SetAttribute{
								ElementType:         types.StringType,
								Optional:            true,
								MarkdownDescription: "Group IDs in scope of policy unless explicitly excluded.",
								Default:             setdefault.StaticValue(types.SetValueMust(types.StringType, []attr.Value{})),
								Computed:            true,
							},
							"include_guests_or_external_users": schema.SingleNestedAttribute{
								MarkdownDescription: "Internal guests or external users included in the policy scope. Optionally populated.",
								Optional:            true,
								Attributes: map[string]schema.Attribute{
									"external_tenants": schema.SingleNestedAttribute{
										MarkdownDescription: "The tenant IDs of the selected types of external users. Either all B2B tenant or a collection of tenant IDs. External tenants can be specified only when the property **guestOrExternalUserTypes** isn't `null` or an empty String.",
										Optional:            true,
										Attributes: map[string]schema.Attribute{
											"membership_kind": schema.StringAttribute{
												MarkdownDescription: "The membership kind. Possible values are: `all`, `enumerated`, `unknownFutureValue`.",
												Required:            true,
												Validators: []validator.String{
													stringvalidator.OneOf("all", "enumerated", "unknownFutureValue"),
												},
											},
											"members": schema.SetAttribute{
												ElementType:         types.StringType,
												Optional:            true,
												MarkdownDescription: "The tenant IDs of the external tenants.",
												Default:             setdefault.StaticValue(types.SetValueMust(types.StringType, []attr.Value{})),
												Computed:            true,
											},
										},
									},
									"guest_or_external_user_types": schema.SetAttribute{
										ElementType:         types.StringType,
										Optional:            true,
										MarkdownDescription: "Indicates internal guests or external user types, and is a multi-valued property. Possible values are: `none`, `internalGuest`, `b2bCollaborationGuest`, `b2bCollaborationMember`, `b2bDirectConnectUser`, `otherExternalUser`, `serviceProvider`, `unknownFutureValue`.",
										Default:             setdefault.StaticValue(types.SetValueMust(types.StringType, []attr.Value{})),
										Computed:            true,
									},
								},
							},
							"include_roles": schema.SetAttribute{
								ElementType:         types.StringType,
								Optional:            true,
								MarkdownDescription: "Role IDs in scope of policy unless explicitly excluded.",
								Default:             setdefault.StaticValue(types.SetValueMust(types.StringType, []attr.Value{})),
								Computed:            true,
							},
							"include_users": schema.SetAttribute{
								ElementType:         types.StringType,
								Optional:            true,
								MarkdownDescription: "User IDs in scope of policy unless explicitly excluded, `None`, `All`, or `GuestsOrExternalUsers`.",
								Default:             setdefault.StaticValue(types.SetValueMust(types.StringType, []attr.Value{})),
								Computed:            true,
							},
						},
					},
					"client_applications": schema.SingleNestedAttribute{
						MarkdownDescription: "Client applications (service principals and workload identities) included in and excluded from the policy. Either **users** or **clientApplications** is required.",
						Optional:            true,
						Attributes: map[string]schema.Attribute{
							"exclude_service_principals": schema.SetAttribute{
								ElementType:         types.StringType,
								Optional:            true,
								MarkdownDescription: "Service principal IDs excluded from the policy scope.",
								Default:             setdefault.StaticValue(types.SetValueMust(types.StringType, []attr.Value{})),
								Computed:            true,
							},
							"include_service_principals": schema.SetAttribute{
								ElementType:         types.StringType,
								Optional:            true,
								MarkdownDescription: "Service principal IDs included in the policy scope, or `ServicePrincipalsInMyTenant`.",
								Default:             setdefault.StaticValue(types.SetValueMust(types.StringType, []attr.Value{})),
								Computed:            true,
							},
							"service_principal_filter": schema.SingleNestedAttribute{
								MarkdownDescription: "Filter that defines the dynamic-servicePrincipal-syntax rule to include/exclude service principals. A filter can use custom security attributes to include/exclude service principals.",
								Optional:            true,
								Attributes: map[string]schema.Attribute{
									"mode": schema.StringAttribute{
										MarkdownDescription: "Mode to use for the filter. Possible values are `include` or `exclude`.",
										Required:            true,
										Validators: []validator.String{
											stringvalidator.OneOf("include", "exclude"),
										},
									},
									"rule": schema.StringAttribute{
										MarkdownDescription: "Rule syntax is similar to that used for membership rules for groups in Microsoft Entra ID.",
										Required:            true,
									},
								},
							},
						},
					},
					"client_app_types": schema.SetAttribute{
						ElementType:         types.StringType,
						Optional:            true,
						MarkdownDescription: "Client application types included in the policy. Possible values are: `all`, `browser`, `mobileAppsAndDesktopClients`, `exchangeActiveSync`, `easSupported`, `other`. Required. The `easUnsupported` enumeration member is deprecated in favor of `exchangeActiveSync`, which includes EAS supported and unsupported platforms.",
						Default:             setdefault.StaticValue(types.SetValueMust(types.StringType, []attr.Value{})),
						Computed:            true,
					},
					"device_states": schema.SingleNestedAttribute{
						MarkdownDescription: "Device states in the policy. To be deprecated and removed. Use the **devices** property instead.",
						Optional:            true,
						Attributes: map[string]schema.Attribute{
							"include_states": schema.SetAttribute{
								ElementType:         types.StringType,
								Optional:            true,
								MarkdownDescription: "States in the scope of the policy. `All` is the only allowed value.",
								Default:             setdefault.StaticValue(types.SetValueMust(types.StringType, []attr.Value{})),
								Computed:            true,
							},
							"exclude_states": schema.SetAttribute{
								ElementType:         types.StringType,
								Optional:            true,
								MarkdownDescription: "States excluded from the scope of the policy. Possible values: `Compliant`, `DomainJoined`.",
								Default:             setdefault.StaticValue(types.SetValueMust(types.StringType, []attr.Value{})),
								Computed:            true,
							},
						},
					},
					"devices": schema.SingleNestedAttribute{
						MarkdownDescription: "Devices in the policy.",
						Optional:            true,
						Attributes: map[string]schema.Attribute{
							"include_devices": schema.SetAttribute{
								ElementType:         types.StringType,
								Optional:            true,
								MarkdownDescription: "States in the scope of the policy. `All` is the only allowed value. Cannot be set if **deviceFilter** is set.",
								Default:             setdefault.StaticValue(types.SetValueMust(types.StringType, []attr.Value{})),
								Computed:            true,
							},
							"exclude_devices": schema.SetAttribute{
								ElementType:         types.StringType,
								Optional:            true,
								MarkdownDescription: "States excluded from the scope of the policy. Possible values: `Compliant`, `DomainJoined`. Cannot be set if **deviceFIlter** is set.",
								Default:             setdefault.StaticValue(types.SetValueMust(types.StringType, []attr.Value{})),
								Computed:            true,
							},
							"device_filter": schema.SingleNestedAttribute{
								MarkdownDescription: "Filter that defines the dynamic-device-syntax rule to include/exclude devices. A filter can use device properties (such as extension attributes) to include/exclude them. Cannot be set if **includeDevices** or **excludeDevices** is set.",
								Optional:            true,
								Attributes: map[string]schema.Attribute{
									"mode": schema.StringAttribute{
										MarkdownDescription: "Mode to use for the filter. Possible values are `include` or `exclude`.",
										Required:            true,
										Validators: []validator.String{
											stringvalidator.OneOf("include", "exclude"),
										},
									},
									"rule": schema.StringAttribute{
										MarkdownDescription: "Rule syntax is similar to that used for membership rules for groups in Microsoft Entra ID.",
										Required:            true,
									},
								},
							},
							"include_device_states": schema.SetAttribute{
								ElementType:         types.StringType,
								Optional:            true,
								MarkdownDescription: "States in the scope of the policy. `All` is the only allowed value. (deprecated)",
								Default:             setdefault.StaticValue(types.SetValueMust(types.StringType, []attr.Value{})),
								Computed:            true,
							},
							"exclude_device_states": schema.SetAttribute{
								ElementType:         types.StringType,
								Optional:            true,
								MarkdownDescription: "States excluded from the scope of the policy. Possible values: `Compliant`, `DomainJoined`. (deprecated)",
								Default:             setdefault.StaticValue(types.SetValueMust(types.StringType, []attr.Value{})),
								Computed:            true,
							},
						},
					},
					"locations": schema.SingleNestedAttribute{
						MarkdownDescription: "Locations included in and excluded from the policy.",
						Optional:            true,
						Attributes: map[string]schema.Attribute{
							"include_locations": schema.SetAttribute{
								ElementType:         types.StringType,
								Optional:            true,
								MarkdownDescription: "Location IDs in scope of policy unless explicitly excluded, `All`, or `AllTrusted`.",
								Default:             setdefault.StaticValue(types.SetValueMust(types.StringType, []attr.Value{})),
								Computed:            true,
							},
							"exclude_locations": schema.SetAttribute{
								ElementType:         types.StringType,
								Optional:            true,
								MarkdownDescription: "Location IDs excluded from scope of policy.",
								Default:             setdefault.StaticValue(types.SetValueMust(types.StringType, []attr.Value{})),
								Computed:            true,
							},
						},
					},
					"platforms": schema.SingleNestedAttribute{
						MarkdownDescription: "Platforms included in and excluded from the policy scope.",
						Optional:            true,
						Attributes: map[string]schema.Attribute{
							"include_platforms": schema.SetAttribute{
								ElementType:         types.StringType,
								Optional:            true,
								MarkdownDescription: "Possible values are: `android`, `iOS`, `windows`, `windowsPhone`, `macOS`, `all`, `unknownFutureValue`, `linux`.",
								Default:             setdefault.StaticValue(types.SetValueMust(types.StringType, []attr.Value{})),
								Computed:            true,
							},
							"exclude_platforms": schema.SetAttribute{
								ElementType:         types.StringType,
								Optional:            true,
								MarkdownDescription: "Possible values are: `android`, `iOS`, `windows`, `windowsPhone`, `macOS`, `all`, `unknownFutureValue`, `linux`.",
								Default:             setdefault.StaticValue(types.SetValueMust(types.StringType, []attr.Value{})),
								Computed:            true,
							},
						},
					},
					"service_principal_risk_levels": schema.SetAttribute{
						ElementType:         types.StringType,
						Optional:            true,
						MarkdownDescription: "Service principal risk levels included in the policy. Possible values are: `low`, `medium`, `high`, `none`, `unknownFutureValue`.",
						Default:             setdefault.StaticValue(types.SetValueMust(types.StringType, []attr.Value{})),
						Computed:            true,
					},
					"sign_in_risk_levels": schema.SetAttribute{
						ElementType:         types.StringType,
						Optional:            true,
						MarkdownDescription: "Sign-in risk levels included in the policy. Possible values are: `low`, `medium`, `high`, `hidden`, `none`, `unknownFutureValue`. Required.",
						Default:             setdefault.StaticValue(types.SetValueMust(types.StringType, []attr.Value{})),
						Computed:            true,
					},
					"user_risk_levels": schema.SetAttribute{
						ElementType:         types.StringType,
						Optional:            true,
						MarkdownDescription: "User risk levels included in the policy. Possible values are: `low`, `medium`, `high`, `hidden`, `none`, `unknownFutureValue`. Required.",
						Default:             setdefault.StaticValue(types.SetValueMust(types.StringType, []attr.Value{})),
						Computed:            true,
					},
					"insider_risk_levels": schema.SetAttribute{
						ElementType:         types.StringType,
						Optional:            true,
						MarkdownDescription: "Insider risk levels included in the policy. The possible values are: `minor`, `moderate`, `elevated`, `unknownFutureValue`.",
						Default:             setdefault.StaticValue(types.SetValueMust(types.StringType, []attr.Value{})),
						Computed:            true,
					},
				},
			},
			"grant_controls": schema.SingleNestedAttribute{
				MarkdownDescription: "Specifies the grant controls that must be fulfilled to pass the policy.",
				Optional:            true,
				Attributes: map[string]schema.Attribute{
					"operator": schema.StringAttribute{
						MarkdownDescription: "Defines the relationship of the grant controls. Possible values: `AND`, `OR`.",
						Optional:            true,
						Default:             stringdefault.StaticString("OR"),
						Computed:            true,
						Validators: []validator.String{
							stringvalidator.OneOf("AND", "OR"),
						},
					},
					"built_in_controls": schema.SetAttribute{
						ElementType:         types.StringType,
						Optional:            true,
						MarkdownDescription: "List of values of built-in controls required by the policy. Possible values: `block`, `mfa`, `compliantDevice`, `domainJoinedDevice`, `approvedApplication`, `compliantApplication`, `passwordChange`, `unknownFutureValue`.",
						Default:             setdefault.StaticValue(types.SetValueMust(types.StringType, []attr.Value{})),
						Computed:            true,
					},
					"custom_authentication_factors": schema.SetAttribute{
						ElementType:         types.StringType,
						Optional:            true,
						MarkdownDescription: "List of custom controls IDs required by the policy. To learn more about custom control, see Custom controls (preview).",
						Default:             setdefault.StaticValue(types.SetValueMust(types.StringType, []attr.Value{})),
						Computed:            true,
					},
					"terms_of_use": schema.SetAttribute{
						ElementType:         types.StringType,
						Optional:            true,
						MarkdownDescription: "List of terms of use IDs required by the policy.",
						Default:             setdefault.StaticValue(types.SetValueMust(types.StringType, []attr.Value{})),
						Computed:            true,
					},
				},
			},
			"session_controls": schema.SingleNestedAttribute{
				MarkdownDescription: "Specifies the session controls that are enforced after sign-in.",
				Optional:            true,
				Attributes: map[string]schema.Attribute{
					"application_enforced_restrictions": schema.SingleNestedAttribute{
						MarkdownDescription: "Session control to enforce application restrictions.",
						Optional:            true,
						Attributes: map[string]schema.Attribute{
							"is_enabled": schema.BoolAttribute{
								MarkdownDescription: "Specifies whether the session control is enabled.",
								Required:            true,
							},
						},
					},
					"cloud_app_security": schema.SingleNestedAttribute{
						MarkdownDescription: "Session control to apply cloud app security.",
						Optional:            true,
						Attributes: map[string]schema.Attribute{
							"is_enabled": schema.BoolAttribute{
								MarkdownDescription: "Specifies whether the session control is enabled.",
								Required:            true,
							},
							"cloud_app_security_type": schema.StringAttribute{
								MarkdownDescription: "The cloud app security control type. Possible values: `mcasConfigured`, `monitorOnly`, `blockDownloads`, `unknownFutureValue`.",
								Optional:            true,
								Validators: []validator.String{
									stringvalidator.OneOf("mcasConfigured", "monitorOnly", "blockDownloads", "unknownFutureValue"),
								},
							},
						},
					},
					"sign_in_frequency": schema.SingleNestedAttribute{
						MarkdownDescription: "Session control to enforce sign-in frequency.",
						Optional:            true,
						Attributes: map[string]schema.Attribute{
							"is_enabled": schema.BoolAttribute{
								MarkdownDescription: "Specifies whether the session control is enabled.",
								Required:            true,
							},
							"type": schema.StringAttribute{
								MarkdownDescription: "The sign-in frequency type. Possible values: `days`, `hours`.",
								Optional:            true,
								Validators: []validator.String{
									stringvalidator.OneOf("days", "hours"),
								},
							},
							"value": schema.Int32Attribute{
								MarkdownDescription: "The sign-in frequency value.",
								Optional:            true,
							},
							"authentication_type": schema.StringAttribute{
								MarkdownDescription: "The authentication type for sign-in frequency. Possible values: `primaryAndSecondaryAuthentication`, `secondaryAuthentication`, `unknownFutureValue`.",
								Optional:            true,
								Validators: []validator.String{
									stringvalidator.OneOf("primaryAndSecondaryAuthentication", "secondaryAuthentication", "unknownFutureValue"),
								},
							},
							"frequency_interval": schema.StringAttribute{
								MarkdownDescription: "The frequency interval. Possible values: `timeBased`, `everyTime`, `unknownFutureValue`.",
								Optional:            true,
								Validators: []validator.String{
									stringvalidator.OneOf("timeBased", "everyTime", "unknownFutureValue"),
								},
							},
						},
					},
					"persistent_browser": schema.SingleNestedAttribute{
						MarkdownDescription: "Session control to define whether a session is persistent.",
						Optional:            true,
						Attributes: map[string]schema.Attribute{
							"is_enabled": schema.BoolAttribute{
								MarkdownDescription: "Specifies whether the session control is enabled.",
								Required:            true,
							},
							"mode": schema.StringAttribute{
								MarkdownDescription: "The persistent browser mode. Possible values: `always`, `never`, `unknownFutureValue`.",
								Optional:            true,
								Validators: []validator.String{
									stringvalidator.OneOf("always", "never", "unknownFutureValue"),
								},
							},
						},
					},
					"disable_resilience_defaults": schema.BoolAttribute{
						MarkdownDescription: "Session control that determines whether it is acceptable for Microsoft Entra ID to extend existing sessions based on information collected prior to an outage.",
						Optional:            true,
						Default:             booldefault.StaticBool(false),
						Computed:            true,
					},
				},
			},
			"timeouts": commonschema.Timeouts(ctx),
		},
	}
}
