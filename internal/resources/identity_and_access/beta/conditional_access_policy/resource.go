package graphBetaConditionalAccessPolicy

import (
	"context"
	"regexp"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/schema"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/objectvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	ResourceName = "graph_beta_identity_and_access_conditional_access_policy"
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
	}
}

type ConditionalAccessPolicyResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ProviderTypeName string
	TypeName         string
	ReadPermissions  []string
	WritePermissions []string
}

// Metadata returns the resource type name.
func (r *ConditionalAccessPolicyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + ResourceName
}

// Configure sets the client for the resource.
func (r *ConditionalAccessPolicyResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = common.SetGraphBetaClientForResource(ctx, req, resp, r.TypeName)
}

// ImportState imports the resource state.
func (r *ConditionalAccessPolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Schema returns the schema for the resource.
func (r *ConditionalAccessPolicyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Specifies the identifier of a conditionalAccessPolicy object. Read-only.",
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "Not used.",
			},
			"display_name": schema.StringAttribute{
				Required:    true,
				Description: "Specifies a display name for the conditionalAccessPolicy object.",
			},
			"created_date_time": schema.StringAttribute{
				Computed:    true,
				Description: "The Timestamp type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is `2014-01-01T00:00:00Z`. Readonly.",
			},
			"modified_date_time": schema.StringAttribute{
				Computed:    true,
				Description: "The Timestamp type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is `2014-01-01T00:00:00Z`. Readonly.",
			},
			"state": schema.StringAttribute{
				Required:    true,
				Description: "Specifies the state of the conditionalAccessPolicy object. Possible values are: `enabled`, `disabled`, `enabledForReportingButNotEnforced`. Required.",
				Validators: []validator.String{
					stringvalidator.OneOf("enabled", "disabled", "enabledForReportingButNotEnforced"),
				},
			},
			"conditions": schema.SingleNestedAttribute{
				Required:    true,
				Description: "Specifies the rules that must be met for the policy to apply. Required.",
				Attributes:  r.conditionalAccessConditionsSchema(),
			},
			"grant_controls": schema.SingleNestedAttribute{
				Required:    true,
				Description: "Specifies the grant controls that must be fulfilled to pass the policy.",
				Attributes:  r.conditionalAccessGrantControlsSchema(),
			},
			"session_controls": schema.SingleNestedAttribute{
				Optional:    true,
				Description: "Specifies the session controls that are enforced after sign-in.",
				Attributes:  r.conditionalAccessSessionControlsSchema(),
			},
			"timeouts": commonschema.Timeouts(ctx),
		},
	}
}

func (r *ConditionalAccessPolicyResource) conditionalAccessConditionsSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"applications": schema.SingleNestedAttribute{
			Required:    true,
			Description: "Applications and user actions included in and excluded from the policy. Required.",
			Attributes:  r.conditionalAccessApplicationsSchema(),
		},
		"authentication_flows": schema.SingleNestedAttribute{
			Optional:    true,
			Description: "Authentication flows included in the policy scope. For more information, see Conditional Access: Authentication flows.",
			Attributes:  r.conditionalAccessAuthenticationFlowsSchema(),
		},
		"users": schema.SingleNestedAttribute{
			Required:    true,
			Description: "Users, groups, and roles included in and excluded from the policy. Either users or clientApplications is required.",
			Attributes:  r.conditionalAccessUsersSchema(),
		},
		"client_applications": schema.SingleNestedAttribute{
			Optional:    true,
			Description: "Client applications (service principals and workload identities) included in and excluded from the policy. Either users or clientApplications is required.",
			Attributes:  r.conditionalAccessClientApplicationsSchema(),
		},
		"client_app_types": schema.ListAttribute{
			Required:    true,
			Description: "Client application types included in the policy. Possible values are: all, browser, mobileAppsAndDesktopClients, exchangeActiveSync, easSupported, other. Required. The easUnsupported enumeration member is deprecated in favor of exchangeActiveSync, which includes EAS supported and unsupported platforms.",
			ElementType: types.StringType,
			Validators: []validator.List{
				listvalidator.ValueStringsAre(
					stringvalidator.OneOf("all", "browser", "mobileAppsAndDesktopClients", "exchangeActiveSync", "easSupported", "other"),
				),
			},
		},
		"device_states": schema.SingleNestedAttribute{
			Optional:    true,
			Description: "Device states in the policy. To be deprecated and removed. Use the devices property instead.",
			Attributes:  r.conditionalAccessDeviceStatesSchema(),
		},
		"devices": schema.SingleNestedAttribute{
			Optional:    true,
			Description: "Devices in the policy.",
			Attributes:  r.conditionalAccessDevicesSchema(),
		},
		"locations": schema.SingleNestedAttribute{
			Optional:    true,
			Description: "Locations included in and excluded from the policy.",
			Attributes:  r.conditionalAccessLocationsSchema(),
		},
		"platforms": schema.SingleNestedAttribute{
			Optional:    true,
			Description: "Platforms included in and excluded from the policy.",
			Attributes:  r.conditionalAccessPlatformsSchema(),
		},
		"service_principal_risk_levels": schema.ListAttribute{
			Optional:    true,
			Description: "Service principal risk levels included in the policy. Possible values are: low, medium, high, none, unknownFutureValue.",
			ElementType: types.StringType,
			Validators: []validator.List{
				listvalidator.ValueStringsAre(
					stringvalidator.OneOf("low", "medium", "high", "none", "unknownFutureValue"),
				),
			},
		},
		"sign_in_risk_levels": schema.ListAttribute{
			Required:    true,
			Description: "Sign-in risk levels included in the policy. Possible values are: low, medium, high, hidden, none, unknownFutureValue. Required.",
			ElementType: types.StringType,
			Validators: []validator.List{
				listvalidator.ValueStringsAre(
					stringvalidator.OneOf("low", "medium", "high", "hidden", "none", "unknownFutureValue"),
				),
			},
		},
		"user_risk_levels": schema.ListAttribute{
			Required:    true,
			Description: "User risk levels included in the policy. Possible values are: low, medium, high, hidden, none, unknownFutureValue. Required.",
			ElementType: types.StringType,
			Validators: []validator.List{
				listvalidator.ValueStringsAre(
					stringvalidator.OneOf("low", "medium", "high", "hidden", "none", "unknownFutureValue"),
				),
			},
		},
		"insider_risk_levels": schema.StringAttribute{
			Optional:    true,
			Description: "Insider risk levels included in the policy. The possible values are: minor, moderate, elevated, unknownFutureValue.",
			Validators: []validator.String{
				stringvalidator.OneOf("minor", "moderate", "elevated", "unknownFutureValue"),
			},
		},
	}
}

func (r *ConditionalAccessPolicyResource) conditionalAccessApplicationsSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"include_applications": schema.ListAttribute{
			Optional:    true,
			Description: "List of application IDs the policy applies to, unless explicitly excluded. Can be one of the following: The list of client IDs (appId) the policy applies to, 'All', 'Office365' (For the list of apps included in Office365, see Apps included in Conditional Access Office 365 app suite), 'MicrosoftAdminPortals' (For more information, see Conditional Access Target resources: Microsoft Admin Portals).",
			ElementType: types.StringType,
			Validators: []validator.List{
				listvalidator.ValueStringsAre(
					stringvalidator.Any(
						stringvalidator.OneOf("All", "Office365", "MicrosoftAdminPortals"),
						stringvalidator.RegexMatches(regexp.MustCompile(`^[0-9a-fA-F]{8}-([0-9a-fA-F]{4}-){3}[0-9a-fA-F]{12}$`), "must be a valid UUID"),
					),
				),
			},
		},
		"exclude_applications": schema.ListAttribute{
			Optional:    true,
			Description: "List of application IDs explicitly excluded from the policy. Can be one of the following: The list of client IDs (appId) explicitly excluded from the policy, 'Office365' (For the list of apps included in Office365, see Apps included in Conditional Access Office 365 app suite), 'MicrosoftAdminPortals' (For more information, see Conditional Access Target resources: Microsoft Admin Portals).",
			ElementType: types.StringType,
			Validators: []validator.List{
				listvalidator.ValueStringsAre(
					stringvalidator.Any(
						stringvalidator.OneOf("Office365", "MicrosoftAdminPortals"),
						stringvalidator.RegexMatches(regexp.MustCompile(`^[0-9a-fA-F]{8}-([0-9a-fA-F]{4}-){3}[0-9a-fA-F]{12}$`), "must be a valid UUID"),
					),
				),
			},
		},
		"application_filter": filterSchema(
			"Filter that defines the dynamic-application-syntax rule to include/exclude cloud applications. A filter can use custom security attributes to include/exclude applications.",
		),
		"include_user_actions": schema.ListAttribute{
			Optional:    true,
			Description: "User actions to include. Supported values are 'urn:user:registersecurityinfo' and 'urn:user:registerdevice'.",
			ElementType: types.StringType,
			Validators: []validator.List{
				listvalidator.ValueStringsAre(
					stringvalidator.OneOf("urn:user:registersecurityinfo", "urn:user:registerdevice"),
				),
			},
		},
		"include_authentication_context_class_references": schema.ListAttribute{
			Optional:    true,
			Description: "Authentication context class references include. Supported values are 'c1' through 'c25'.",
			ElementType: types.StringType,
			Validators: []validator.List{
				listvalidator.ValueStringsAre(
					stringvalidator.RegexMatches(regexp.MustCompile(`^c([1-9]|1[0-9]|2[0-5])$`), "must be in the format 'c1' through 'c25'"),
				),
			},
		},
	}
}

func (r *ConditionalAccessPolicyResource) conditionalAccessAuthenticationFlowsSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"transfer_methods": schema.StringAttribute{
			Optional:    true,
			Description: "Represents the transfer methods in scope for the policy. The possible values are: none, deviceCodeFlow, authenticationTransfer, unknownFutureValue.",
			Validators: []validator.String{
				stringvalidator.OneOf(
					"none",
					"deviceCodeFlow",
					"authenticationTransfer",
					"unknownFutureValue",
				),
			},
		},
	}
}

func (r *ConditionalAccessPolicyResource) conditionalAccessUsersSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"exclude_groups": schema.ListAttribute{
			Optional:    true,
			Description: "Group IDs excluded from scope of policy.",
			ElementType: types.StringType,
			Validators: []validator.List{
				listvalidator.ValueStringsAre(
					stringvalidator.RegexMatches(regexp.MustCompile(`^[0-9a-fA-F]{8}-([0-9a-fA-F]{4}-){3}[0-9a-fA-F]{12}$`), "must be a valid UUID"),
				),
			},
		},
		"exclude_guests_or_external_users": schema.SingleNestedAttribute{
			Optional:    true,
			Description: "Internal guests or external users excluded from the policy scope. Optionally populated.",
			Attributes:  r.conditionalAccessGuestsOrExternalUsersSchema(),
		},
		"exclude_roles": schema.ListAttribute{
			Optional:    true,
			Description: "Role IDs excluded from scope of policy.",
			ElementType: types.StringType,
			Validators: []validator.List{
				listvalidator.ValueStringsAre(
					stringvalidator.RegexMatches(regexp.MustCompile(`^[0-9a-fA-F]{8}-([0-9a-fA-F]{4}-){3}[0-9a-fA-F]{12}$`), "must be a valid UUID"),
				),
			},
		},
		"exclude_users": schema.ListAttribute{
			Optional:    true,
			Description: "User IDs excluded from scope of policy and/or GuestsOrExternalUsers.",
			ElementType: types.StringType,
			Validators: []validator.List{
				listvalidator.ValueStringsAre(
					stringvalidator.Any(
						stringvalidator.OneOf("GuestsOrExternalUsers"),
						stringvalidator.RegexMatches(regexp.MustCompile(`^[0-9a-fA-F]{8}-([0-9a-fA-F]{4}-){3}[0-9a-fA-F]{12}$`), "must be a valid UUID"),
					),
				),
			},
		},
		"include_groups": schema.ListAttribute{
			Optional:    true,
			Description: "Group IDs in scope of policy unless explicitly excluded.",
			ElementType: types.StringType,
			Validators: []validator.List{
				listvalidator.ValueStringsAre(
					stringvalidator.RegexMatches(regexp.MustCompile(`^[0-9a-fA-F]{8}-([0-9a-fA-F]{4}-){3}[0-9a-fA-F]{12}$`), "must be a valid UUID"),
				),
			},
		},
		"include_guests_or_external_users": schema.SingleNestedAttribute{
			Optional:    true,
			Description: "Internal guests or external users included in the policy scope. Optionally populated.",
			Attributes:  r.conditionalAccessGuestsOrExternalUsersSchema(),
		},
		"include_roles": schema.ListAttribute{
			Optional:    true,
			Description: "Role IDs in scope of policy unless explicitly excluded.",
			ElementType: types.StringType,
			Validators: []validator.List{
				listvalidator.ValueStringsAre(
					stringvalidator.RegexMatches(regexp.MustCompile(`^[0-9a-fA-F]{8}-([0-9a-fA-F]{4}-){3}[0-9a-fA-F]{12}$`), "must be a valid UUID"),
				),
			},
		},
		"include_users": schema.ListAttribute{
			Optional:    true,
			Description: "User IDs in scope of policy unless explicitly excluded, None, All, or GuestsOrExternalUsers.",
			ElementType: types.StringType,
			Validators: []validator.List{
				listvalidator.ValueStringsAre(
					stringvalidator.Any(
						stringvalidator.OneOf("None", "All", "GuestsOrExternalUsers"),
						stringvalidator.RegexMatches(regexp.MustCompile(`^[0-9a-fA-F]{8}-([0-9a-fA-F]{4}-){3}[0-9a-fA-F]{12}$`), "must be a valid UUID"),
					),
				),
			},
		},
	}
}

func (r *ConditionalAccessPolicyResource) conditionalAccessClientApplicationsSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"exclude_service_principals": schema.ListAttribute{
			Optional:    true,
			Description: "Service principal IDs excluded from the policy scope.",
			ElementType: types.StringType,
			Validators: []validator.List{
				listvalidator.ValueStringsAre(
					stringvalidator.RegexMatches(regexp.MustCompile(`^[0-9a-fA-F]{8}-([0-9a-fA-F]{4}-){3}[0-9a-fA-F]{12}$`), "must be a valid UUID"),
				),
			},
		},
		"include_service_principals": schema.ListAttribute{
			Optional:    true,
			Description: "Service principal IDs included in the policy scope, or 'ServicePrincipalsInMyTenant'.",
			ElementType: types.StringType,
			Validators: []validator.List{
				listvalidator.ValueStringsAre(
					stringvalidator.Any(
						stringvalidator.OneOf("ServicePrincipalsInMyTenant"),
						stringvalidator.RegexMatches(regexp.MustCompile(`^[0-9a-fA-F]{8}-([0-9a-fA-F]{4}-){3}[0-9a-fA-F]{12}$`), "must be a valid UUID"),
					),
				),
			},
		},
		"service_principal_filter": filterSchema(
			"Filter that defines the dynamic-servicePrincipal-syntax rule to include/exclude service principals. A filter can use custom security attributes to include/exclude service principals.",
		),
	}
}

func (r *ConditionalAccessPolicyResource) conditionalAccessDeviceStatesSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"include_states": schema.ListAttribute{
			Optional:    true,
			Description: "States in the scope of the policy. 'All' is the only allowed value.",
			ElementType: types.StringType,
			Validators: []validator.List{
				listvalidator.SizeAtMost(1),
				listvalidator.ValueStringsAre(
					stringvalidator.OneOf("All"),
				),
			},
		},
		"exclude_states": schema.ListAttribute{
			Optional:    true,
			Description: "States excluded from the scope of the policy. Possible values: 'Compliant', 'DomainJoined'.",
			ElementType: types.StringType,
			Validators: []validator.List{
				listvalidator.ValueStringsAre(
					stringvalidator.OneOf("Compliant", "DomainJoined"),
				),
			},
		},
	}
}

func (r *ConditionalAccessPolicyResource) conditionalAccessDevicesSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"include_devices": schema.ListAttribute{
			Optional:    true,
			Description: "States in the scope of the policy. 'All' is the only allowed value. Cannot be set if deviceFilter is set.",
			ElementType: types.StringType,
			Validators: []validator.List{
				listvalidator.SizeAtMost(1),
				listvalidator.ValueStringsAre(
					stringvalidator.OneOf("All"),
				),
				listvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("device_filter")),
			},
		},
		"exclude_devices": schema.ListAttribute{
			Optional:    true,
			Description: "States excluded from the scope of the policy. Possible values: 'Compliant', 'DomainJoined'. Cannot be set if deviceFilter is set.",
			ElementType: types.StringType,
			Validators: []validator.List{
				listvalidator.ValueStringsAre(
					stringvalidator.OneOf("Compliant", "DomainJoined"),
				),
				listvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("device_filter")),
			},
		},
		"device_filter": filterSchema(
			"Filter that defines the dynamic-device-syntax rule to include/exclude devices. A filter can use device properties (such as extension attributes) to include/exclude them. Cannot be set if includeDevices or excludeDevices is set.",
		),
		"include_device_states": schema.ListAttribute{
			Optional:           true,
			Description:        "(Deprecated) States in the scope of the policy. 'All' is the only allowed value.",
			DeprecationMessage: "This field is deprecated. Use include_devices instead.",
			ElementType:        types.StringType,
			Validators: []validator.List{
				listvalidator.SizeAtMost(1),
				listvalidator.ValueStringsAre(
					stringvalidator.OneOf("All"),
				),
			},
		},
		"exclude_device_states": schema.ListAttribute{
			Optional:           true,
			Description:        "(Deprecated) States excluded from the scope of the policy. Possible values: 'Compliant', 'DomainJoined'.",
			DeprecationMessage: "This field is deprecated. Use exclude_devices instead.",
			ElementType:        types.StringType,
			Validators: []validator.List{
				listvalidator.ValueStringsAre(
					stringvalidator.OneOf("Compliant", "DomainJoined"),
				),
			},
		},
	}
}

func (r *ConditionalAccessPolicyResource) conditionalAccessLocationsSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"include_locations": schema.ListAttribute{
			Required:    true,
			Description: "Location IDs in scope of policy unless explicitly excluded, 'All', or 'AllTrusted'.",
			ElementType: types.StringType,
			Validators: []validator.List{
				listvalidator.ValueStringsAre(
					stringvalidator.Any(
						stringvalidator.OneOf("All", "AllTrusted"),
						stringvalidator.RegexMatches(regexp.MustCompile(`^[0-9a-fA-F]{8}-([0-9a-fA-F]{4}-){3}[0-9a-fA-F]{12}$`), "must be a valid UUID"),
					),
				),
			},
		},
		"exclude_locations": schema.ListAttribute{
			Optional:    true,
			Description: "Location IDs excluded from scope of policy.",
			ElementType: types.StringType,
			Validators: []validator.List{
				listvalidator.ValueStringsAre(
					stringvalidator.RegexMatches(regexp.MustCompile(`^[0-9a-fA-F]{8}-([0-9a-fA-F]{4}-){3}[0-9a-fA-F]{12}$`), "must be a valid UUID"),
				),
			},
		},
	}
}

func (r *ConditionalAccessPolicyResource) conditionalAccessPlatformsSchema() map[string]schema.Attribute {
	platformValues := []string{
		"android",
		"iOS",
		"windows",
		"windowsPhone",
		"macOS",
		"all",
		"unknownFutureValue",
		"linux",
	}

	return map[string]schema.Attribute{
		"include_platforms": schema.ListAttribute{
			Optional:    true,
			Description: "Device platforms included in the policy scope.",
			ElementType: types.StringType,
			Validators: []validator.List{
				listvalidator.ValueStringsAre(
					stringvalidator.OneOf(platformValues...),
				),
			},
		},
		"exclude_platforms": schema.ListAttribute{
			Optional:    true,
			Description: "Device platforms excluded from the policy scope.",
			ElementType: types.StringType,
			Validators: []validator.List{
				listvalidator.ValueStringsAre(
					stringvalidator.OneOf(platformValues...),
				),
			},
		},
	}
}

func (r *ConditionalAccessPolicyResource) conditionalAccessGuestsOrExternalUsersSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"external_tenants": schema.SingleNestedAttribute{
			Optional:    true,
			Description: "The tenant IDs of the selected types of external users.",
			Attributes: map[string]schema.Attribute{
				"membership_kind": schema.StringAttribute{
					Required:    true,
					Description: "Indicates the kind of membership. Possible values are: all, enumerated, unknownFutureValue.",
					Validators: []validator.String{
						stringvalidator.OneOf("all", "enumerated", "unknownFutureValue"),
					},
				},
			},
		},
		"guest_or_external_user_types": schema.StringAttribute{
			Required:    true,
			Description: "Indicates internal guests or external user types. Possible values are: none, internalGuest, b2bCollaborationGuest, b2bCollaborationMember, b2bDirectConnectUser, otherExternalUser, serviceProvider, unknownFutureValue.",
			Validators: []validator.String{
				stringvalidator.OneOf(
					"none",
					"internalGuest",
					"b2bCollaborationGuest",
					"b2bCollaborationMember",
					"b2bDirectConnectUser",
					"otherExternalUser",
					"serviceProvider",
					"unknownFutureValue",
				),
			},
		},
	}
}

func (r *ConditionalAccessPolicyResource) conditionalAccessGrantControlsSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"operator": schema.StringAttribute{
			Required:    true,
			Description: "Defines the relationship of the grant controls. Possible values: AND, OR.",
			Validators: []validator.String{
				stringvalidator.OneOf("AND", "OR"),
			},
		},
		"built_in_controls": schema.ListAttribute{
			Optional:    true,
			Description: "List of values of built-in controls required by the policy. Possible values: block, mfa, compliantDevice, domainJoinedDevice, approvedApplication, compliantApplication, passwordChange, unknownFutureValue.",
			ElementType: types.StringType,
			Validators: []validator.List{
				listvalidator.ValueStringsAre(
					stringvalidator.OneOf(
						"block",
						"mfa",
						"compliantDevice",
						"domainJoinedDevice",
						"approvedApplication",
						"compliantApplication",
						"passwordChange",
						"unknownFutureValue",
					),
				),
			},
		},
		"custom_authentication_factors": schema.ListAttribute{
			Optional:    true,
			Description: "List of custom controls IDs required by the policy. To learn more about custom control, see Custom controls (preview).",
			ElementType: types.StringType,
		},
		"terms_of_use": schema.ListAttribute{
			Optional:    true,
			Description: "List of terms of use IDs required by the policy.",
			ElementType: types.StringType,
		},
		"authentication_strength": schema.SingleNestedAttribute{
			Optional:    true,
			Description: "The authentication strength policy to be used in the grant controls.",
			Attributes: map[string]schema.Attribute{
				"id": schema.StringAttribute{
					Optional:    true,
					Description: "The identifier of the authentication strength policy.",
				},
				"display_name": schema.StringAttribute{
					Optional:    true,
					Description: "The display name of the authentication strength policy.",
				},
				"description": schema.StringAttribute{
					Optional:    true,
					Description: "The description of the authentication strength policy.",
				},
				"created_date_time": schema.StringAttribute{
					Computed:    true,
					Description: "The creation date and time of the authentication strength policy.",
				},
				"modified_date_time": schema.StringAttribute{
					Computed:    true,
					Description: "The last modified date and time of the authentication strength policy.",
				},
				"policy_type": schema.StringAttribute{
					Optional:    true,
					Description: "The type of the authentication strength policy.",
				},
				"requirements_satisfied": schema.StringAttribute{
					Optional:    true,
					Description: "The requirements satisfied by the authentication strength policy.",
				},
				"allowed_combinations": schema.ListAttribute{
					Optional:    true,
					Description: "The allowed combinations for the authentication strength policy.",
					ElementType: types.StringType,
				},
			},
		},
	}
}

func (r *ConditionalAccessPolicyResource) conditionalAccessSessionControlsSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"application_enforced_restrictions": schema.SingleNestedAttribute{
			Optional:    true,
			Description: "Session control to enforce application restrictions. Only Exchange Online and Sharepoint Online support this session control.",
			Attributes: map[string]schema.Attribute{
				"is_enabled": schema.BoolAttribute{
					Required:    true,
					Description: "Specifies whether the session control is enabled.",
				},
			},
		},
		"cloud_app_security": schema.SingleNestedAttribute{
			Optional:    true,
			Description: "Session control to apply cloud app security.",
			Attributes: map[string]schema.Attribute{
				"is_enabled": schema.BoolAttribute{
					Required:    true,
					Description: "Specifies whether the session control is enabled.",
				},
				"cloud_app_security_type": schema.StringAttribute{
					Required:    true,
					Description: "Possible values are: mcasConfigured, monitorOnly, blockDownloads, unknownFutureValue.",
					Validators: []validator.String{
						stringvalidator.OneOf("mcasConfigured", "monitorOnly", "blockDownloads", "unknownFutureValue"),
					},
				},
			},
		},
		"continuous_access_evaluation": schema.SingleNestedAttribute{
			Optional:    true,
			Description: "Session control for continuous access evaluation settings.",
			Attributes: map[string]schema.Attribute{
				"mode": schema.StringAttribute{
					Required:    true,
					Description: "Specifies continuous access evaluation settings. The possible values are: strictEnforcement, disabled, unknownFutureValue, strictLocation. Note that you must use the Prefer: include-unknown-enum-members request header to get the strictLocation value in this evolvable enum.",
					Validators: []validator.String{
						stringvalidator.OneOf("strictEnforcement", "disabled", "unknownFutureValue", "strictLocation"),
					},
				},
			},
		},
		"persistent_browser": schema.SingleNestedAttribute{
			Optional:    true,
			Description: "Session control to define whether to persist cookies or not. All apps should be selected for this session control to work correctly.",
			Attributes: map[string]schema.Attribute{
				"is_enabled": schema.BoolAttribute{
					Required:    true,
					Description: "Specifies whether the session control is enabled.",
				},
				"mode": schema.StringAttribute{
					Required:    true,
					Description: "Possible values are: always, never.",
					Validators: []validator.String{
						stringvalidator.OneOf("always", "never"),
					},
				},
			},
		},
		"sign_in_frequency": schema.SingleNestedAttribute{
			Optional:    true,
			Description: "Session control to enforce signin frequency.",
			Attributes: map[string]schema.Attribute{
				"is_enabled": schema.BoolAttribute{
					Required:    true,
					Description: "Specifies whether the session control is enabled.",
				},
				"type": schema.StringAttribute{
					Required:    true,
					Description: "Possible values are: days, hours.",
					Validators: []validator.String{
						stringvalidator.OneOf("days", "hours"),
					},
				},
				"value": schema.Int64Attribute{
					Required:    true,
					Description: "The number of days or hours.",
				},
				"authentication_type": schema.StringAttribute{
					Required:    true,
					Description: "The possible values are: primaryAndSecondaryAuthentication, secondaryAuthentication, unknownFutureValue.",
					Validators: []validator.String{
						stringvalidator.OneOf("primaryAndSecondaryAuthentication", "secondaryAuthentication", "unknownFutureValue"),
					},
				},
				"frequency_interval": schema.StringAttribute{
					Required:    true,
					Description: "The possible values are: timeBased, everyTime, unknownFutureValue.",
					Validators: []validator.String{
						stringvalidator.OneOf("timeBased", "everyTime", "unknownFutureValue"),
					},
				},
			},
		},
		"disable_resilience_defaults": schema.BoolAttribute{
			Optional:    true,
			Description: "Session control that determines whether it's acceptable for Microsoft Entra ID to extend existing sessions based on information collected prior to an outage or not.",
		},
		"secure_sign_in_session": schema.SingleNestedAttribute{
			Optional:    true,
			Description: "Session control to require sign in sessions to be bound to a device.",
			Attributes: map[string]schema.Attribute{
				"is_enabled": schema.BoolAttribute{
					Required:    true,
					Description: "Specifies whether the session control is enabled.",
				},
			},
		},
	}
}

func filterSchema(description string) schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Optional:    true,
		Description: description,
		Attributes: map[string]schema.Attribute{
			"mode": schema.StringAttribute{
				Required:    true,
				Description: "Mode to use for the filter. Possible values are include or exclude.",
				Validators: []validator.String{
					stringvalidator.OneOf("include", "exclude"),
				},
			},
			"rule": schema.StringAttribute{
				Required:    true,
				Description: "Rule syntax is similar to that used for membership rules for groups in Microsoft Entra ID.",
			},
		},
		Validators: []validator.Object{
			objectvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("include_devices")),
			objectvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("exclude_devices")),
		},
	}
}
