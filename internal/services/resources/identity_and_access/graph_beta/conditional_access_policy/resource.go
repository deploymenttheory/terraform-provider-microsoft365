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

// https://learn.microsoft.com/en-us/graph/custom-security-attributes-examples?tabs=http#prerequisites
func NewConditionalAccessPolicyResource() resource.Resource {
	return &ConditionalAccessPolicyResource{
		ReadPermissions: []string{
			"Policy.Read.All",
			"Policy.Read.ConditionalAccess",
			"Directory.Read.All",                    // for validation of roles
			"CustomSecAttributeAssignment.Read.All", // for custom security attributes
			"Application.Read.All",                  // for custom security attributes
		},
		WritePermissions: []string{
			"Policy.ReadWrite.ConditionalAccess",
			"CustomSecAttributeAssignment.ReadWrite.All", // Read and write custom security attribute assignments
			"Application.Read.All",                       // needs read permissions for write operations
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
						"must be a valid GUID in the format '00000000-0000-0000-0000-000000000000'",
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
											stringvalidator.OneOf("All", "None", "MicrosoftAdminPortals", "Office365"),
											stringvalidator.RegexMatches(
												regexp.MustCompile(constants.GuidRegex),
												"must be a valid GUID or one of the special values: All, None, MicrosoftAdminPortals, Office365",
											),
										),
									),
								},
							},
							"exclude_applications": schema.SetAttribute{
								MarkdownDescription: "Applications to exclude from the policy. For empty requests, use []",
								ElementType:         types.StringType,
								Required:            true,
								Validators: []validator.Set{
									setvalidator.ValueStringsAre(
										stringvalidator.Any(
											stringvalidator.OneOf("All", "MicrosoftAdminPortals", "Office365"),
											stringvalidator.RegexMatches(
												regexp.MustCompile(constants.GuidRegex),
												"must be a valid GUID or one of the special values: All, MicrosoftAdminPortals, Office365",
											),
										),
									),
								},
							},
							"include_user_actions": schema.SetAttribute{
								MarkdownDescription: "User actions to include in the policy.",
								ElementType:         types.StringType,
								Required:            true,
								Validators: []validator.Set{
									setvalidator.ValueStringsAre(
										stringvalidator.OneOf(
											"urn:user:registersecurityinfo",
											"urn:user:registerdevice",
										),
									),
								},
							},
							"include_authentication_context_class_references": schema.SetAttribute{
								MarkdownDescription: "Authentication context secures data and actions in applications, including custom applications, line-of-business (LOB) " +
									"applications, SharePoint, and applications protected by Microsoft Defender for Cloud Apps. Can be predefined builtin contexts: " +
									"`require_trusted_device` (or c1), `require_terms_of_use` (or c2), `require_trusted_location` (or c3), `require_strong_authentication` (or c4), " +
									"`required_trust_type:azure_ad_joined` (or c5), `require_access_from_an_approved_app` (or c6), `required_trust_type:hybrid_azure_ad_joined` (or c7) " +
									"or custom authentication context class references in the format 'c' followed by a number from 8 through to 99 " +
									"(e.g., c1, c8, c10, c25, c99). Learn more here 'https://learn.microsoft.com/en-us/entra/identity/conditional-access/concept-conditional-access-cloud-apps#authentication-context'.",
								ElementType: types.StringType,
								Required:    true,
								Validators: []validator.Set{
									setvalidator.ValueStringsAre(
										stringvalidator.Any(
											stringvalidator.OneOf(
												"require_trusted_device",
												"require_terms_of_use",
												"require_trusted_location",
												"require_strong_authentication",
												"required_trust_type:azure_ad_joined",
												"require_access_from_an_approved_app",
												"required_trust_type:hybrid_azure_ad_joined",
											),
											stringvalidator.RegexMatches(
												regexp.MustCompile(`^c([1-9]|[1-9][0-9])$`),
												"must be in the format 'c' followed by a number from 1 to 99",
											),
										),
									),
								},
							},
							"application_filter": schema.SingleNestedAttribute{
								MarkdownDescription: "Configure app filters you want to policy to apply to. Using custom security attributes you can use the rule builder " +
									"or rule syntax text box to create or edit the filter rules. this feature is currently in preview, only attributes of type String are supported. " +
									"Attributes of type Integer or Boolean are not currently supported. Learn more here 'https://learn.microsoft.com/en-us/entra/identity/conditional-access/concept-filter-for-applications'.",
								Optional: true,
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
							"global_secure_access": schema.SingleNestedAttribute{
								MarkdownDescription: "Global Secure Access settings for the conditional access policy.",
								Optional:            true,
								Computed:            true,
								Attributes:          map[string]schema.Attribute{
									// Note: This field appears in API responses but is typically null
									// Adding minimal structure based on API observations
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
								Required:            true,
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
								Required:            true,
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
								Required:            true,
								Validators: []validator.Set{
									setvalidator.ValueStringsAre(
										stringvalidator.RegexMatches(
											regexp.MustCompile(constants.GuidRegex),
											"must be a valid GUID in the format '00000000-0000-0000-0000-000000000000'",
										),
									),
								},
							},
							"exclude_groups": schema.SetAttribute{
								MarkdownDescription: "Groups to exclude from the policy.",
								ElementType:         types.StringType,
								Required:            true,
								Validators: []validator.Set{
									setvalidator.ValueStringsAre(
										stringvalidator.RegexMatches(
											regexp.MustCompile(constants.GuidRegex),
											"must be a valid GUID in the format '00000000-0000-0000-0000-000000000000'",
										),
									),
								},
							},
							"include_roles": schema.SetAttribute{
								MarkdownDescription: "Roles to include in the policy.",
								ElementType:         types.StringType,
								Required:            true,
								Validators: []validator.Set{
									setvalidator.ValueStringsAre(
										stringvalidator.RegexMatches(
											regexp.MustCompile(constants.GuidRegex),
											"must be a valid GUID in the format '00000000-0000-0000-0000-000000000000'",
										),
									),
								},
							},
							"exclude_roles": schema.SetAttribute{
								MarkdownDescription: "Microsoft Entra tenant roles to exclude from the policy.",
								ElementType:         types.StringType,
								Required:            true,
								Validators: []validator.Set{
									setvalidator.ValueStringsAre(
										stringvalidator.RegexMatches(
											regexp.MustCompile(constants.GuidRegex),
											"must be a valid GUID in the format '00000000-0000-0000-0000-000000000000'",
										),
									),
								},
							},
							"include_guests_or_external_users": schema.SingleNestedAttribute{
								MarkdownDescription: "Configuration for including guests or external users.",
								Optional:            true,
								Computed:            true,
								Attributes: map[string]schema.Attribute{
									"guest_or_external_user_types": schema.SetAttribute{
										MarkdownDescription: "Types of guests or external users to include. Possible values are: InternalGuest, B2bCollaborationGuest, B2bCollaborationMember, B2bDirectConnectUser, OtherExternalUser, ServiceProvider.",
										ElementType:         types.StringType,
										Optional:            true,
										Computed:            true,
										Validators: []validator.Set{
											setvalidator.ValueStringsAre(
												stringvalidator.OneOf(
													"b2bCollaborationGuest",
													"b2bCollaborationMember",
													"b2bDirectConnectUser",
													"internalGuest",
													"serviceProvider",
													"otherExternalUser",
												),
											),
										},
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
												MarkdownDescription: "The list of Microsoft Entra organization tenant IDs for external tenants to exclude from the CA policy.",
												ElementType:         types.StringType,
												Optional:            true,
												Computed:            true,
												Validators: []validator.Set{
													setvalidator.ValueStringsAre(
														stringvalidator.RegexMatches(
															regexp.MustCompile(constants.GuidRegex),
															"must be a valid GUID in the format '00000000-0000-0000-0000-000000000000'",
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
								Computed:            true,
								Attributes: map[string]schema.Attribute{
									"guest_or_external_user_types": schema.SetAttribute{
										MarkdownDescription: "Types of guests or external users to exclude. Possible values are: InternalGuest, B2bCollaborationGuest, B2bCollaborationMember, B2bDirectConnectUser, OtherExternalUser, ServiceProvider.",
										ElementType:         types.StringType,
										Optional:            true,
										Computed:            true,
										Validators: []validator.Set{
											setvalidator.ValueStringsAre(
												stringvalidator.OneOf(
													"b2bCollaborationGuest",
													"b2bCollaborationMember",
													"b2bDirectConnectUser",
													"internalGuest",
													"serviceProvider",
													"otherExternalUser",
												),
											),
										},
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
												Optional:            true,
												Computed:            true,
												Validators: []validator.Set{
													setvalidator.ValueStringsAre(
														stringvalidator.RegexMatches(
															regexp.MustCompile(constants.GuidRegex),
															"must be a valid GUID in the format '00000000-0000-0000-0000-000000000000'",
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
								MarkdownDescription: "Named locations to include in the policy. Can use special values like 'All' or 'AllTrusted' 'or provide guid's" +
									"of named locations.",
								ElementType: types.StringType,
								Required:    true,
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
								MarkdownDescription: "Named locations to exclude from the policy. Can use special values like 'AllTrusted' or provide guid's" +
									"of named locations.",
								ElementType: types.StringType,
								Required:    true,
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
						MarkdownDescription: "Operator to apply to the controls. Possible values are: AND, OR. When setting a singular operator, use 'OR'.",
						Required:            true,
						Validators: []validator.String{
							stringvalidator.OneOf("AND", "OR"),
						},
					},
					"built_in_controls": schema.SetAttribute{
						MarkdownDescription: "List of built-in controls required by the policy. Possible values are: block, mfa, compliantDevice, domainJoinedDevice, approvedApplication, compliantApplication, passwordChange, unknownFutureValue.",
						ElementType:         types.StringType,
						Required:            true,
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
								),
							),
						},
					},
					"custom_authentication_factors": schema.SetAttribute{
						MarkdownDescription: "Custom authentication factors for granting access.",
						ElementType:         types.StringType,
						Required:            true,
					},
					"terms_of_use": schema.SetAttribute{
						MarkdownDescription: "Terms of use required for granting access.",
						ElementType:         types.StringType,
						Optional:            true,
					},
					"authentication_strength": schema.SingleNestedAttribute{
						MarkdownDescription: "Authentication strength is a Conditional Access control that specifies which combinations of authentication " +
							"methods can be used to access a resource. Users can satisfy the strength requirements by authenticating with any of the allowed " +
							"combinations. read more here 'https://learn.microsoft.com/en-us/entra/identity/authentication/concept-authentication-strengths'.",
						Optional: true,
						Attributes: map[string]schema.Attribute{
							"id": schema.StringAttribute{
								MarkdownDescription: "ID of the authentication strength policy. Can be a GUID or predefined built-in values: 'multifactor_authentication' (maps to '00000000-0000-0000-0000-000000000002'), 'passwordless_mfa' (maps to '00000000-0000-0000-0000-000000000003'), or 'phishing_resistant_mfa' (maps to '00000000-0000-0000-0000-000000000004').",
								Required:            true,
								Validators: []validator.String{
									stringvalidator.Any(
										stringvalidator.RegexMatches(
											regexp.MustCompile(constants.GuidRegex),
											"must be a valid GUID in the format '00000000-0000-0000-0000-000000000000'",
										),
										stringvalidator.OneOf(
											"multifactor_authentication",
											"passwordless_mfa",
											"phishing_resistant_mfa"),
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
								Computed:            true,
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
								MarkdownDescription: "Type of sign-in frequency control. Possible values are: days, hours. Not used when frequency_interval is everyTime.",
								Optional:            true,
								Validators: []validator.String{
									stringvalidator.OneOf("days", "hours"),
								},
							},
							"value": schema.Int64Attribute{
								MarkdownDescription: "Value for the sign-in frequency. Not used when frequency_interval is everyTime.",
								Optional:            true,
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
					"global_secure_access_filtering_profile": schema.SingleNestedAttribute{
						MarkdownDescription: "Global Secure Access filtering profile for the session.",
						Optional:            true,
						Attributes: map[string]schema.Attribute{
							"is_enabled": schema.BoolAttribute{
								MarkdownDescription: "Whether global secure access filtering controls are enabled.",
								Required:            true,
							},
							"profile_id": schema.StringAttribute{
								MarkdownDescription: "ID of the global secure access filtering profile.",
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
