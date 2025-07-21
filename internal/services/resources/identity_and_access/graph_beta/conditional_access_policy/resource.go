package graphBetaConditionalAccessPolicy

import (
	"context"
	"regexp"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const (
	ResourceName  = "graph_beta_identity_and_access_conditional_access_policy"
	CreateTimeout = 180
	UpdateTimeout = 180
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
			"Policy.Read.ConditionalAccess",
		},
		WritePermissions: []string{
			"Policy.ReadWrite.ConditionalAccess",
		},
		ResourcePath: "/identity/conditionalAccess/policies",
	}
}

type ConditionalAccessPolicyResource struct {
	httpClient       *client.AuthenticatedHTTPClient
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

// FullTypeName returns the full resource type name in the format "providername_resourcename".
func (r *ConditionalAccessPolicyResource) FullTypeName() string {
	return r.ProviderTypeName + "_" + r.TypeName
}

// Configure sets the client for the resource.
func (r *ConditionalAccessPolicyResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.httpClient = client.SetGraphBetaHTTPClientForResource(ctx, req, resp, constants.PROVIDER_NAME+"_"+ResourceName)
}

// ImportState imports the resource state.
func (r *ConditionalAccessPolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Schema defines the schema for the resource.
func (r *ConditionalAccessPolicyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages Microsoft 365 Conditional Access Policies using the `/identity/conditionalAccess/policies` endpoint. Conditional Access policies define the conditions under which users can access cloud apps.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "String (identifier)",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(constants.GuidRegex),
						"must be a valid GUID",
					),
				},
			},
			"display_name": schema.StringAttribute{
				MarkdownDescription: "The display name for the Conditional Access policy.",
				Required:            true,
			},
			"state": schema.StringAttribute{
				MarkdownDescription: "Specifies the state of the policy. Possible values are: enabled, disabled, enabledForReportingButNotEnforced.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("enabled", "disabled", "enabledForReportingButNotEnforced"),
				},
			},
			"created_date_time": schema.StringAttribute{
				MarkdownDescription: "The creation date and time of the policy.",
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(constants.TimeFormatRFC3339Regex),
						"must be a valid RFC3339 date-time string",
					),
				},
			},
			"modified_date_time": schema.StringAttribute{
				MarkdownDescription: "The last modified date and time of the policy.",
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(constants.TimeFormatRFC3339Regex),
						"must be a valid RFC3339 date-time string",
					),
				},
			},
			"deleted_date_time": schema.StringAttribute{
				MarkdownDescription: "The deletion date and time of the policy, if applicable.",
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(constants.TimeFormatRFC3339Regex),
						"must be a valid RFC3339 date-time string",
					),
				},
			},
			"conditions": schema.SingleNestedAttribute{
				MarkdownDescription: "Conditions that must be met for the policy to apply.",
				Required:            true,
				Attributes: map[string]schema.Attribute{
					"client_app_types": schema.SetAttribute{
						MarkdownDescription: "Client application types included in the policy. Possible values are: all, browser, mobileAppsAndDesktopClients, exchangeActiveSync, other.",
						ElementType:         types.StringType,
						Required:            true,
						Validators: []validator.Set{
							setvalidator.ValueStringsAre(
								stringvalidator.OneOf(
									"all",
									"browser",
									"mobileAppsAndDesktopClients",
									"exchangeActiveSync",
									"other",
								),
							),
						},
					},
					"applications": schema.SingleNestedAttribute{
						MarkdownDescription: "Applications and user actions included in and excluded from the policy.",
						Required:            true,
						Attributes: map[string]schema.Attribute{
							"include_applications": schema.SetAttribute{
								MarkdownDescription: "Applications to include in the policy. Can use the special value 'All' to include all applications.",
								ElementType:         types.StringType,
								Required:            true,
								Validators: []validator.Set{
									setvalidator.ValueStringsAre(
										stringvalidator.Any(
											stringvalidator.OneOf("All", "Office365"),
											stringvalidator.RegexMatches(
												regexp.MustCompile(constants.GuidRegex),
												"must be a valid GUID or one of the special values: All, Office365",
											),
										),
									),
								},
							},
							"exclude_applications": schema.SetAttribute{
								MarkdownDescription: "Applications to exclude from the policy.",
								ElementType:         types.StringType,
								Optional:            true,
								Validators: []validator.Set{
									setvalidator.ValueStringsAre(
										stringvalidator.RegexMatches(
											regexp.MustCompile(constants.GuidRegex),
											"must be a valid GUID",
										),
									),
								},
							},
							"include_user_actions": schema.SetAttribute{
								MarkdownDescription: "User actions to include in the policy.",
								ElementType:         types.StringType,
								Optional:            true,
								Validators: []validator.Set{
									setvalidator.ValueStringsAre(
										stringvalidator.OneOf(
											"urn:user:registersecurityinfo",
										),
									),
								},
							},
							"include_authentication_context_class_references": schema.SetAttribute{
								MarkdownDescription: "Authentication context class references to include in the policy.",
								ElementType:         types.StringType,
								Optional:            true,
								Validators: []validator.Set{
									setvalidator.ValueStringsAre(
										stringvalidator.RegexMatches(
											regexp.MustCompile(`^c[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`),
											"must be in the format 'c' followed by a GUID",
										),
									),
								},
							},
							"application_filter": schema.SingleNestedAttribute{
								MarkdownDescription: "Filter that defines the applications the policy applies to.",
								Optional:            true,
								Attributes: map[string]schema.Attribute{
									"mode": schema.StringAttribute{
										MarkdownDescription: "Mode of the filter. Possible values are: include, exclude.",
										Required:            true,
										Validators: []validator.String{
											stringvalidator.OneOf("include", "exclude"),
										},
									},
									"rule": schema.StringAttribute{
										MarkdownDescription: "Rule syntax for the filter.",
										Required:            true,
									},
								},
							},
						},
					},
					"users": schema.SingleNestedAttribute{
						MarkdownDescription: "Users, groups, and roles included in and excluded from the policy.",
						Required:            true,
						Attributes: map[string]schema.Attribute{
							"include_users": schema.SetAttribute{
								MarkdownDescription: "Users to include in the policy. Can use special values like 'All', 'None', or 'GuestsOrExternalUsers'.",
								ElementType:         types.StringType,
								Optional:            true,
								Validators: []validator.Set{
									setvalidator.ValueStringsAre(
										stringvalidator.Any(
											stringvalidator.OneOf("All", "None", "GuestsOrExternalUsers"),
											stringvalidator.RegexMatches(
												regexp.MustCompile(constants.GuidRegex),
												"must be a valid GUID or one of the special values: All, None, GuestsOrExternalUsers",
											),
										),
									),
								},
							},
							"exclude_users": schema.SetAttribute{
								MarkdownDescription: "Users to exclude from the policy. Can use special values like 'GuestsOrExternalUsers'.",
								ElementType:         types.StringType,
								Optional:            true,
								Validators: []validator.Set{
									setvalidator.ValueStringsAre(
										stringvalidator.Any(
											stringvalidator.OneOf("GuestsOrExternalUsers"),
											stringvalidator.RegexMatches(
												regexp.MustCompile(constants.GuidRegex),
												"must be a valid GUID or the special value: GuestsOrExternalUsers",
											),
										),
									),
								},
							},
							"include_groups": schema.SetAttribute{
								MarkdownDescription: "Groups to include in the policy.",
								ElementType:         types.StringType,
								Optional:            true,
								Validators: []validator.Set{
									setvalidator.ValueStringsAre(
										stringvalidator.RegexMatches(
											regexp.MustCompile(constants.GuidRegex),
											"must be a valid GUID",
										),
									),
								},
							},
							"exclude_groups": schema.SetAttribute{
								MarkdownDescription: "Groups to exclude from the policy.",
								ElementType:         types.StringType,
								Optional:            true,
								Validators: []validator.Set{
									setvalidator.ValueStringsAre(
										stringvalidator.RegexMatches(
											regexp.MustCompile(constants.GuidRegex),
											"must be a valid GUID",
										),
									),
								},
							},
							"include_roles": schema.SetAttribute{
								MarkdownDescription: "Roles to include in the policy.",
								ElementType:         types.StringType,
								Optional:            true,
								Validators: []validator.Set{
									setvalidator.ValueStringsAre(
										stringvalidator.RegexMatches(
											regexp.MustCompile(constants.GuidRegex),
											"must be a valid GUID",
										),
									),
								},
							},
							"exclude_roles": schema.SetAttribute{
								MarkdownDescription: "Roles to exclude from the policy.",
								ElementType:         types.StringType,
								Optional:            true,
								Validators: []validator.Set{
									setvalidator.ValueStringsAre(
										stringvalidator.RegexMatches(
											regexp.MustCompile(constants.GuidRegex),
											"must be a valid GUID",
										),
									),
								},
							},
							"include_guests_or_external_users": schema.SingleNestedAttribute{
								MarkdownDescription: "Configuration for including guests or external users.",
								Optional:            true,
								Attributes: map[string]schema.Attribute{
									"guest_or_external_user_types": schema.StringAttribute{
										MarkdownDescription: "Types of guests or external users to include. Possible values are: internalGuest, b2bCollaborationGuest, b2bCollaborationMember, b2bDirectConnectUser, otherExternalUser, serviceProvider.",
										Required:            true,
									},
									"external_tenants": schema.SingleNestedAttribute{
										MarkdownDescription: "Configuration for external tenants.",
										Required:            true,
										Attributes: map[string]schema.Attribute{
											"membership_kind": schema.StringAttribute{
												MarkdownDescription: "Kind of membership. Possible values are: all, enumerated, unknownFutureValue.",
												Required:            true,
												Validators: []validator.String{
													stringvalidator.OneOf("all", "enumerated", "unknownFutureValue"),
												},
											},
											"members": schema.SetAttribute{
												MarkdownDescription: "The list of tenant IDs for external tenants.",
												ElementType:         types.StringType,
												Required:            true,
												Validators: []validator.Set{
													setvalidator.ValueStringsAre(
														stringvalidator.RegexMatches(
															regexp.MustCompile(constants.GuidRegex),
															"must be a valid GUID",
														),
													),
												},
											},
										},
									},
								},
							},
							"exclude_guests_or_external_users": schema.SingleNestedAttribute{
								MarkdownDescription: "Configuration for excluding guests or external users.",
								Optional:            true,
								Attributes: map[string]schema.Attribute{
									"guest_or_external_user_types": schema.StringAttribute{
										MarkdownDescription: "Types of guests or external users to exclude. Possible values are: internalGuest, b2bCollaborationGuest, b2bCollaborationMember, b2bDirectConnectUser, otherExternalUser, serviceProvider.",
										Required:            true,
									},
									"external_tenants": schema.SingleNestedAttribute{
										MarkdownDescription: "Configuration for external tenants.",
										Required:            true,
										Attributes: map[string]schema.Attribute{
											"membership_kind": schema.StringAttribute{
												MarkdownDescription: "Kind of membership. Possible values are: all, enumerated, unknownFutureValue.",
												Required:            true,
												Validators: []validator.String{
													stringvalidator.OneOf("all", "enumerated", "unknownFutureValue"),
												},
											},
											"members": schema.SetAttribute{
												MarkdownDescription: "The list of tenant IDs for external tenants.",
												ElementType:         types.StringType,
												Required:            true,
												Validators: []validator.Set{
													setvalidator.ValueStringsAre(
														stringvalidator.RegexMatches(
															regexp.MustCompile(constants.GuidRegex),
															"must be a valid GUID",
														),
													),
												},
											},
										},
									},
								},
							},
						},
					},
					"platforms": schema.SingleNestedAttribute{
						MarkdownDescription: "Platforms included in and excluded from the policy.",
						Optional:            true,
						Attributes: map[string]schema.Attribute{
							"include_platforms": schema.SetAttribute{
								MarkdownDescription: "Platforms to include in the policy.",
								ElementType:         types.StringType,
								Required:            true,
								Validators: []validator.Set{
									setvalidator.ValueStringsAre(
										stringvalidator.OneOf(
											"all",
											"android",
											"iOS",
											"windows",
											"windowsPhone",
											"macOS",
											"linux",
											"unknownFutureValue",
										),
									),
								},
							},
							"exclude_platforms": schema.SetAttribute{
								MarkdownDescription: "Platforms to exclude from the policy.",
								ElementType:         types.StringType,
								Optional:            true,
								Validators: []validator.Set{
									setvalidator.ValueStringsAre(
										stringvalidator.OneOf(
											"all",
											"android",
											"iOS",
											"windows",
											"windowsPhone",
											"macOS",
											"linux",
											"unknownFutureValue",
										),
									),
								},
							},
						},
					},
					"locations": schema.SingleNestedAttribute{
						MarkdownDescription: "Locations included in and excluded from the policy.",
						Optional:            true,
						Attributes: map[string]schema.Attribute{
							"include_locations": schema.SetAttribute{
								MarkdownDescription: "Locations to include in the policy. Can use special values like 'All' or 'AllTrusted'.",
								ElementType:         types.StringType,
								Required:            true,
								Validators: []validator.Set{
									setvalidator.ValueStringsAre(
										stringvalidator.Any(
											stringvalidator.OneOf("All", "AllTrusted"),
											stringvalidator.RegexMatches(
												regexp.MustCompile(constants.GuidRegex),
												"must be a valid GUID or one of the special values: All, AllTrusted",
											),
										),
									),
								},
							},
							"exclude_locations": schema.SetAttribute{
								MarkdownDescription: "Locations to exclude from the policy. Can use special values like 'AllTrusted'.",
								ElementType:         types.StringType,
								Optional:            true,
								Validators: []validator.Set{
									setvalidator.ValueStringsAre(
										stringvalidator.Any(
											stringvalidator.OneOf("AllTrusted"),
											stringvalidator.RegexMatches(
												regexp.MustCompile(constants.GuidRegex),
												"must be a valid GUID or the special value: AllTrusted",
											),
										),
									),
								},
							},
						},
					},
					"devices": schema.SingleNestedAttribute{
						MarkdownDescription: "Devices included in and excluded from the policy.",
						Optional:            true,
						Attributes: map[string]schema.Attribute{
							"include_devices": schema.SetAttribute{
								MarkdownDescription: "Devices to include in the policy.",
								ElementType:         types.StringType,
								Optional:            true,
							},
							"exclude_devices": schema.SetAttribute{
								MarkdownDescription: "Devices to exclude from the policy.",
								ElementType:         types.StringType,
								Optional:            true,
							},
							"include_device_states": schema.SetAttribute{
								MarkdownDescription: "Device states to include in the policy.",
								ElementType:         types.StringType,
								Optional:            true,
							},
							"exclude_device_states": schema.SetAttribute{
								MarkdownDescription: "Device states to exclude from the policy.",
								ElementType:         types.StringType,
								Optional:            true,
							},
							"device_filter": schema.SingleNestedAttribute{
								MarkdownDescription: "Filter that defines the devices the policy applies to.",
								Optional:            true,
								Attributes: map[string]schema.Attribute{
									"mode": schema.StringAttribute{
										MarkdownDescription: "Mode of the filter. Possible values are: include, exclude.",
										Required:            true,
										Validators: []validator.String{
											stringvalidator.OneOf("include", "exclude"),
										},
									},
									"rule": schema.StringAttribute{
										MarkdownDescription: "Rule syntax for the filter.",
										Required:            true,
									},
								},
							},
						},
					},
					"sign_in_risk_levels": schema.SetAttribute{
						MarkdownDescription: "Sign-in risk levels included in the policy. Possible values are: low, medium, high, hidden, none, unknownFutureValue.",
						ElementType:         types.StringType,
						Required:            true,
						Validators: []validator.Set{
							setvalidator.ValueStringsAre(
								stringvalidator.OneOf(
									"low",
									"medium",
									"high",
									"hidden",
									"none",
									"unknownFutureValue",
								),
							),
						},
					},
					"user_risk_levels": schema.SetAttribute{
						MarkdownDescription: "User risk levels included in the policy. Possible values are: low, medium, high, hidden, none, unknownFutureValue.",
						ElementType:         types.StringType,
						Optional:            true,
						Validators: []validator.Set{
							setvalidator.ValueStringsAre(
								stringvalidator.OneOf(
									"low",
									"medium",
									"high",
									"hidden",
									"none",
									"unknownFutureValue",
								),
							),
						},
					},
					"service_principal_risk_levels": schema.SetAttribute{
						MarkdownDescription: "Service principal risk levels included in the policy. Possible values are: low, medium, high, hidden, none, unknownFutureValue.",
						ElementType:         types.StringType,
						Optional:            true,
						Validators: []validator.Set{
							setvalidator.ValueStringsAre(
								stringvalidator.OneOf(
									"low",
									"medium",
									"high",
									"hidden",
									"none",
									"unknownFutureValue",
								),
							),
						},
					},
					"times": schema.SingleNestedAttribute{
						MarkdownDescription: "Times and days when the policy applies.",
						Optional:            true,
						Attributes: map[string]schema.Attribute{
							"included_ranges": schema.SetAttribute{
								MarkdownDescription: "Time ranges when the policy applies.",
								ElementType:         types.StringType,
								Optional:            true,
							},
							"excluded_ranges": schema.SetAttribute{
								MarkdownDescription: "Time ranges when the policy does not apply.",
								ElementType:         types.StringType,
								Optional:            true,
							},
							"all_day": schema.BoolAttribute{
								MarkdownDescription: "Whether the policy applies all day.",
								Optional:            true,
							},
							"start_time": schema.StringAttribute{
								MarkdownDescription: "Start time for the policy.",
								Optional:            true,
								Validators: []validator.String{
									stringvalidator.RegexMatches(
										regexp.MustCompile(constants.TimeFormatHHMMSSRegex),
										"Time must be in the format 'HH:MM:SS' (24-hour format)",
									),
								},
							},
							"end_time": schema.StringAttribute{
								MarkdownDescription: "End time for the policy.",
								Optional:            true,
								Validators: []validator.String{
									stringvalidator.RegexMatches(
										regexp.MustCompile(constants.TimeFormatHHMMSSRegex),
										"Time must be in the format 'HH:MM:SS' (24-hour format)",
									),
								},
							},
							"time_zone": schema.StringAttribute{
								MarkdownDescription: "Time zone for the policy times.",
								Optional:            true,
							},
						},
					},
					"device_states": schema.SingleNestedAttribute{
						MarkdownDescription: "Device states included in and excluded from the policy.",
						Optional:            true,
						Attributes: map[string]schema.Attribute{
							"include_states": schema.SetAttribute{
								MarkdownDescription: "Device states to include in the policy.",
								ElementType:         types.StringType,
								Optional:            true,
							},
							"exclude_states": schema.SetAttribute{
								MarkdownDescription: "Device states to exclude from the policy.",
								ElementType:         types.StringType,
								Optional:            true,
							},
						},
					},
					"client_applications": schema.SingleNestedAttribute{
						MarkdownDescription: "Client applications configuration for the conditional access policy.",
						Optional:            true,
						Attributes: map[string]schema.Attribute{
							"include_service_principals": schema.SetAttribute{
								MarkdownDescription: "Service principals to include in the policy. Can use the special value 'ServicePrincipalsInMyTenant' to include all service principals.",
								ElementType:         types.StringType,
								Required:            true,
							},
							"exclude_service_principals": schema.SetAttribute{
								MarkdownDescription: "Service principals to exclude from the policy.",
								ElementType:         types.StringType,
								Optional:            true,
							},
						},
					},
				},
			},
			"grant_controls": schema.SingleNestedAttribute{
				MarkdownDescription: "Controls for granting access.",
				Required:            true,
				Attributes: map[string]schema.Attribute{
					"operator": schema.StringAttribute{
						MarkdownDescription: "Operator to apply to the controls. Possible values are: AND, OR.",
						Required:            true,
						Validators: []validator.String{
							stringvalidator.OneOf("AND", "OR"),
						},
					},
					"built_in_controls": schema.SetAttribute{
						MarkdownDescription: "List of built-in controls required by the policy. Possible values are: block, mfa, compliantDevice, domainJoinedDevice, approvedApplication, compliantApplication, passwordChange, unknownFutureValue.",
						ElementType:         types.StringType,
						Optional:            true,
						Validators: []validator.Set{
							setvalidator.ValueStringsAre(
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
					"custom_authentication_factors": schema.SetAttribute{
						MarkdownDescription: "Custom authentication factors for granting access.",
						ElementType:         types.StringType,
						Optional:            true,
					},
					"terms_of_use": schema.SetAttribute{
						MarkdownDescription: "Terms of use required for granting access.",
						ElementType:         types.StringType,
						Optional:            true,
					},
					"authentication_strength": schema.SingleNestedAttribute{
						MarkdownDescription: "Authentication strength required for granting access.",
						Optional:            true,
						Attributes: map[string]schema.Attribute{
							"id": schema.StringAttribute{
								MarkdownDescription: "ID of the authentication strength policy.",
								Required:            true,
								Validators: []validator.String{
									stringvalidator.RegexMatches(
										regexp.MustCompile(constants.GuidRegex),
										"must be a valid GUID",
									),
								},
							},
							"display_name": schema.StringAttribute{
								MarkdownDescription: "Display name of the authentication strength policy.",
								Optional:            true,
								Computed:            true,
							},
							"description": schema.StringAttribute{
								MarkdownDescription: "Description of the authentication strength policy.",
								Optional:            true,
								Computed:            true,
							},
							"policy_type": schema.StringAttribute{
								MarkdownDescription: "Type of the policy. Possible values are: builtIn, custom.",
								Optional:            true,
								Computed:            true,
								Validators: []validator.String{
									stringvalidator.OneOf("builtIn", "custom"),
								},
							},
							"requirements_satisfied": schema.StringAttribute{
								MarkdownDescription: "Requirements satisfied by the policy.",
								Optional:            true,
								Computed:            true,
							},
							"allowed_combinations": schema.SetAttribute{
								MarkdownDescription: "The allowed authentication method combinations that satisfy the authentication strength policy.",
								ElementType:         types.StringType,
								Optional:            true,
								Validators: []validator.Set{
									setvalidator.ValueStringsAre(
										stringvalidator.OneOf(
											"windowsHelloForBusiness",
											"fido2",
											"x509CertificateMultiFactor",
											"deviceBasedPush",
											"temporaryAccessPassOneTime",
											"temporaryAccessPassMultiUse",
											"password,microsoftAuthenticatorPush",
											"password,softwareOath",
											"password,hardwareOath",
											"password,sms",
											"password,voice",
											"federatedMultiFactor",
											"microsoftAuthenticatorPush,federatedSingleFactor",
											"softwareOath,federatedSingleFactor",
											"hardwareOath,federatedSingleFactor",
											"sms,federatedSingleFactor",
											"voice,federatedSingleFactor",
										),
									),
								},
							},
							"created_date_time": schema.StringAttribute{
								MarkdownDescription: "Creation date and time of the authentication strength policy.",
								Computed:            true,
								Validators: []validator.String{
									stringvalidator.RegexMatches(
										regexp.MustCompile(constants.TimeFormatRFC3339Regex),
										"must be a valid RFC3339 date-time string",
									),
								},
							},
							"modified_date_time": schema.StringAttribute{
								MarkdownDescription: "Last modified date and time of the authentication strength policy.",
								Computed:            true,
								Validators: []validator.String{
									stringvalidator.RegexMatches(
										regexp.MustCompile(constants.TimeFormatRFC3339Regex),
										"must be a valid RFC3339 date-time string",
									),
								},
							},
						},
					},
				},
			},
			"session_controls": schema.SingleNestedAttribute{
				MarkdownDescription: "Controls for managing user sessions.",
				Optional:            true,
				Attributes: map[string]schema.Attribute{
					"application_enforced_restrictions": schema.SingleNestedAttribute{
						MarkdownDescription: "Application enforced restrictions for the session.",
						Optional:            true,
						Attributes: map[string]schema.Attribute{
							"is_enabled": schema.BoolAttribute{
								MarkdownDescription: "Whether application enforced restrictions are enabled.",
								Required:            true,
							},
						},
					},
					"cloud_app_security": schema.SingleNestedAttribute{
						MarkdownDescription: "Cloud app security controls for the session.",
						Optional:            true,
						Attributes: map[string]schema.Attribute{
							"is_enabled": schema.BoolAttribute{
								MarkdownDescription: "Whether cloud app security controls are enabled.",
								Required:            true,
							},
							"cloud_app_security_type": schema.StringAttribute{
								MarkdownDescription: "Type of cloud app security control. Possible values are: blockDownloads, mcasConfigured, monitorOnly, unknownFutureValue.",
								Required:            true,
								Validators: []validator.String{
									stringvalidator.OneOf("blockDownloads", "mcasConfigured", "monitorOnly", "unknownFutureValue"),
								},
							},
						},
					},
					"sign_in_frequency": schema.SingleNestedAttribute{
						MarkdownDescription: "Sign-in frequency controls for the session.",
						Optional:            true,
						Attributes: map[string]schema.Attribute{
							"is_enabled": schema.BoolAttribute{
								MarkdownDescription: "Whether sign-in frequency controls are enabled.",
								Required:            true,
							},
							"type": schema.StringAttribute{
								MarkdownDescription: "Type of sign-in frequency control. Possible values are: days, hours.",
								Required:            true,
								Validators: []validator.String{
									stringvalidator.OneOf("days", "hours"),
								},
							},
							"value": schema.Int64Attribute{
								MarkdownDescription: "Value for the sign-in frequency.",
								Required:            true,
							},
							"authentication_type": schema.StringAttribute{
								MarkdownDescription: "Authentication type for sign-in frequency. Possible values are: primaryAndSecondaryAuthentication, secondaryAuthentication.",
								Optional:            true,
								Computed:            true,
								Validators: []validator.String{
									stringvalidator.OneOf("primaryAndSecondaryAuthentication", "secondaryAuthentication"),
								},
							},
							"frequency_interval": schema.StringAttribute{
								MarkdownDescription: "Frequency interval for sign-in frequency. Possible values are: timeBased, everyTime.",
								Optional:            true,
								Computed:            true,
								Validators: []validator.String{
									stringvalidator.OneOf("timeBased", "everyTime"),
								},
							},
						},
					},
					"persistent_browser": schema.SingleNestedAttribute{
						MarkdownDescription: "Persistent browser controls for the session.",
						Optional:            true,
						Attributes: map[string]schema.Attribute{
							"is_enabled": schema.BoolAttribute{
								MarkdownDescription: "Whether persistent browser controls are enabled.",
								Required:            true,
							},
							"mode": schema.StringAttribute{
								MarkdownDescription: "Mode for persistent browser. Possible values are: always, never.",
								Required:            true,
								Validators: []validator.String{
									stringvalidator.OneOf("always", "never"),
								},
							},
						},
					},
					"disable_resilience_defaults": schema.BoolAttribute{
						MarkdownDescription: "Whether to disable resilience defaults.",
						Optional:            true,
					},
					"continuous_access_evaluation": schema.SingleNestedAttribute{
						MarkdownDescription: "Continuous access evaluation controls for the session.",
						Optional:            true,
						Attributes: map[string]schema.Attribute{
							"mode": schema.StringAttribute{
								MarkdownDescription: "Mode for continuous access evaluation. Possible values are: disabled, basic, strict.",
								Required:            true,
								Validators: []validator.String{
									stringvalidator.OneOf("disabled", "basic", "strict"),
								},
							},
						},
					},
					"secure_sign_in_session": schema.SingleNestedAttribute{
						MarkdownDescription: "Secure sign-in session controls.",
						Optional:            true,
						Attributes: map[string]schema.Attribute{
							"is_enabled": schema.BoolAttribute{
								MarkdownDescription: "Whether secure sign-in session controls are enabled.",
								Required:            true,
							},
						},
					},
				},
			},
			"template_id": schema.StringAttribute{
				MarkdownDescription: "ID of the template this policy is derived from.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(constants.GuidOrEmptyValueRegex),
						"must be a valid GUID or empty",
					),
				},
			},
			"partial_enablement_strategy": schema.StringAttribute{
				MarkdownDescription: "Strategy for partial enablement of the policy.",
				Optional:            true,
				Computed:            true,
			},
			"timeouts": commonschema.Timeouts(ctx),
		},
	}
}
