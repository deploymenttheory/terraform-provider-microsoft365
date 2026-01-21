package graphBetaUsersUser

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	attributevalidator "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/validate/attribute"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	ResourceName  = "microsoft365_graph_beta_users_user"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	// Basic resource interface (CRUD operations)
	_ resource.Resource = &UserResource{}

	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &UserResource{}

	// Enables import functionality
	_ resource.ResourceWithImportState = &UserResource{}

	// Enables plan modification/diff suppression
	_ resource.ResourceWithModifyPlan = &UserResource{}
)

func NewUserResource() resource.Resource {
	return &UserResource{
		ReadPermissions: []string{
			"User.Read.All",
			"Directory.Read.All",
			"CustomSecAttributeAssignment.Read.All",
		},
		WritePermissions: []string{
			"User.ReadWrite.All",
			"Directory.ReadWrite.All",
			"CustomSecAttributeAssignment.ReadWrite.All",
			"LifeCycleInfo.ReadWrite.All ",
			"User-PasswordProfile.ReadWrite.All",
		},
		ResourcePath: "/users",
	}
}

type UserResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

// Metadata returns the resource type name.
func (r *UserResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = ResourceName
}

// Configure sets the client for the resource.
func (r *UserResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, ResourceName)
}

// ImportState handles importing the resource with an extended ID format.
//
// Supported formats:
//   - Simple:   "resource_id" (hard_delete defaults to false)
//   - Extended: "resource_id:hard_delete=true" or "resource_id:hard_delete=false"
//
// Example:
//
//	terraform import microsoft365_graph_beta_users_user.example "12345678-1234-1234-1234-123456789012:hard_delete=true"
func (r *UserResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	idParts := strings.Split(req.ID, ":")
	resourceID := idParts[0]
	hardDelete := false // Default to soft delete for safety

	if len(idParts) > 1 {
		for _, part := range idParts[1:] {
			if strings.HasPrefix(part, "hard_delete=") {
				value := strings.TrimPrefix(part, "hard_delete=")
				switch strings.ToLower(value) {
				case "true":
					hardDelete = true
				case "false":
					hardDelete = false
				default:
					resp.Diagnostics.AddError(
						"Invalid Import ID",
						fmt.Sprintf("Invalid hard_delete value '%s'. Must be 'true' or 'false'.", value),
					)
					return
				}
			}
		}
	}

	tflog.Info(ctx, fmt.Sprintf("Importing %s with ID: %s, hard_delete: %t", ResourceName, resourceID, hardDelete))

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), resourceID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("hard_delete"), hardDelete)...)
}

// Schema defines the schema for the resource.
func (r *UserResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages Microsoft 365 users using the `/users` endpoint. This resource is used to the user resource lets admins specify user preferences for languages and date/time formats for the user's primary Exchange mailboxes and Microsoft Entra profile. Permissions for this resource are complex and depend on the specific fields you wish tomanage. For more information, see the Microsoft Documentation. https://learn.microsoft.com/en-us/graph/api/user-update?view=graph-rest-beta&tabs=http#permissions-for-specific-scenarios.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique identifier for the user. Read-only.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"about_me": schema.StringAttribute{
				MarkdownDescription: "A freeform text entry field for users to describe themselves.",
				Optional:            true,
				Computed:            true,
			},
			"account_enabled": schema.BoolAttribute{
				MarkdownDescription: "Set to `true` if the account is enabled; otherwise, `false`. This property is required when a user is created.",
				Required:            true,
			},
			"age_group": schema.StringAttribute{
				MarkdownDescription: "Sets the age group of the user. Allowed values: `null`, `Minor`, `NotAdult`, `Adult`. Refer to the legal age group property definitions for further information.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("null", "Minor", "NotAdult", "Adult"),
				},
			},
			"business_phones": schema.SetAttribute{
				MarkdownDescription: "The telephone numbers for the user. NOTE: Although it is a string collection, only one number can be set for this property.",
				ElementType:         types.StringType,
				Optional:            true,
				Computed:            true,
				Validators: []validator.Set{
					setvalidator.SizeAtMost(1),
				},
			},
			"city": schema.StringAttribute{
				MarkdownDescription: "The city where the user is located. Maximum length is 128 characters.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(128),
				},
			},
			"company_name": schema.StringAttribute{
				MarkdownDescription: "The name of the company that the user is associated with. This property can be useful " +
					"for describing the company that an external user comes from. Maximum length is 64 characters.",
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(64),
				},
			},
			"consent_provided_for_minor": schema.StringAttribute{
				MarkdownDescription: "Sets whether consent was obtained for minors. Allowed values: `null`, `Granted`, `Denied`, `NotRequired`. " +
					"Refer to the legal age group property definitions for further information.",
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.OneOf("null", "Granted", "Denied", "NotRequired"),
				},
			},
			"country": schema.StringAttribute{
				MarkdownDescription: "The country or region where the user is located; for example, `US` or `UK`. Maximum length is 128 characters.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(128),
				},
			},
			"created_date_time": schema.StringAttribute{
				MarkdownDescription: "The date and time the user was created, in ISO 8601 format and UTC. Read-only.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"creation_type": schema.StringAttribute{
				MarkdownDescription: "Indicates whether the user account was created through one of the following methods: " +
					"As a regular school or work account (`null`), as an external account (`Invitation`), as a local account for " +
					"an Azure Active Directory B2C tenant (`LocalAccount`), through self-service sign-up by an internal user using " +
					"email verification (`EmailVerified`), or through self-service sign-up by an external user signing up through a " +
					"link that is part of a user flow (`SelfServiceSignUp`). Read-only.",
				Computed: true,
			},
			"custom_security_attributes": schema.SetNestedAttribute{
				MarkdownDescription: "An open complex type that holds the value of a custom security attribute that is assigned to " +
					"a directory object. Nullable. Returned only on `$select`. Supports `$filter` (eq, ne, not, startsWith). " +
					"The filter value is case-sensitive. To read this property, the calling app must be assigned the " +
					"`CustomSecAttributeAssignment.Read.All` permission. To write this property, the calling app must be assigned " +
					"the `CustomSecAttributeAssignment.ReadWrite.All` permission.",
				Optional: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"attribute_set": schema.StringAttribute{
							MarkdownDescription: "The name of the attribute set (e.g., `Engineering`, `Marketing`). " +
								"This groups related custom security attributes together.",
							Required: true,
						},
						"attributes": schema.SetNestedAttribute{
							MarkdownDescription: "The collection of custom security attributes within this attribute set.",
							Required:            true,
							NestedObject: schema.NestedAttributeObject{
								Validators: []validator.Object{
									attributevalidator.ExactlyOneOfMixedTypes(
										"string_value",
										"int_value",
										"bool_value",
										"string_values",
										"int_values",
									),
								},
								Attributes: map[string]schema.Attribute{
									"name": schema.StringAttribute{
										MarkdownDescription: "The name of the custom security attribute.",
										Required:            true,
									},
									"string_value": schema.StringAttribute{
										MarkdownDescription: "The value if the attribute is a single-valued string type. " +
											"Exactly one of `string_value`, `int_value`, `bool_value`, `string_values`, or `int_values` must be specified.",
										Optional: true,
									},
									"int_value": schema.Int32Attribute{
										MarkdownDescription: "The value if the attribute is a single-valued integer type. " +
											"Exactly one of `string_value`, `int_value`, `bool_value`, `string_values`, or `int_values` must be specified.",
										Optional: true,
									},
									"bool_value": schema.BoolAttribute{
										MarkdownDescription: "The value if the attribute is a boolean type. " +
											"Exactly one of `string_value`, `int_value`, `bool_value`, `string_values`, or `int_values` must be specified.",
										Optional: true,
									},
									"string_values": schema.SetAttribute{
										MarkdownDescription: "The values if the attribute is a multi-valued string type. " +
											"Exactly one of `string_value`, `int_value`, `bool_value`, `string_values`, or `int_values` must be specified.",
										Optional:    true,
										ElementType: types.StringType,
									},
									"int_values": schema.SetAttribute{
										MarkdownDescription: "The values if the attribute is a multi-valued integer type. " +
											"Exactly one of `string_value`, `int_value`, `bool_value`, `string_values`, or `int_values` must be specified.",
										Optional:    true,
										ElementType: types.Int32Type,
									},
								},
							},
						},
					},
				},
			},
			"deleted_date_time": schema.StringAttribute{
				MarkdownDescription: "The date and time the user was deleted. Read-only.",
				Computed:            true,
			},
			"department": schema.StringAttribute{
				MarkdownDescription: "The name of the department in which the user works. Maximum length is 64 characters.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(64),
				},
			},
			"display_name": schema.StringAttribute{
				MarkdownDescription: "The name displayed in the address book for the user. This is usually the combination " +
					"of the user's first name, middle initial, and last name. This property is required when a user is created " +
					"and it cannot be cleared during updates. Maximum length is 256 characters.",
				Required: true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 256),
				},
			},
			"employee_hire_date": schema.StringAttribute{
				MarkdownDescription: "The date and time when the user was hired or will start work in case of a future hire, " +
					"in ISO 8601 format and UTC.",
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(constants.ISO8601DateTimeRegex),
						"must be a valid ISO 8601 datetime format (e.g., 2023-05-01T13:45:30Z)",
					),
				},
			},
			"employee_id": schema.StringAttribute{
				MarkdownDescription: "The employee identifier assigned to the user by the organization. Maximum length is 16 characters.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(16),
				},
			},
			"employee_type": schema.StringAttribute{
				MarkdownDescription: "Captures enterprise worker type. For example, `Employee`, `Contractor`, `Consultant`, or `Vendor`.",
				Optional:            true,
				Computed:            true,
			},
			"external_user_state": schema.StringAttribute{
				MarkdownDescription: "For an external user invited to the tenant using the invitation API, this property represents the " +
					"invited user's invitation status. For invited users, the state can be `PendingAcceptance` or `Accepted`, or `null` for all other users. Read-only.",
				Computed: true,
				Validators: []validator.String{
					stringvalidator.OneOf("PendingAcceptance", "Accepted"),
				},
			},
			"external_user_state_change_date_time": schema.StringAttribute{
				MarkdownDescription: "Shows the timestamp for the latest change to the externalUserState property. Read-only.",
				Computed:            true,
			},
			"fax_number": schema.StringAttribute{
				MarkdownDescription: "The fax number of the user.",
				Optional:            true,
				Computed:            true,
			},
			"given_name": schema.StringAttribute{
				MarkdownDescription: "The given name (first name) of the user. Maximum length is 64 characters.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(64),
				},
			},
			"job_title": schema.StringAttribute{
				MarkdownDescription: "The user's job title. Maximum length is 128 characters.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(128),
				},
			},
			"mail": schema.StringAttribute{
				MarkdownDescription: "The SMTP address for the user, for example, `jeff@contoso.com`. Changes to this property also update the user's proxyAddresses collection to include the value as an SMTP address.",
				Optional:            true,
				Computed:            true,
			},
			"mail_nickname": schema.StringAttribute{
				MarkdownDescription: "The mail alias for the user. This property must be specified when a user is created. Maximum length is 64 characters.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 64),
				},
			},
			"manager_id": schema.StringAttribute{
				MarkdownDescription: "The user ID of the user's manager. Used to set the organizational hierarchy. " +
					"To update the manager, provide the user ID of the new manager. To remove the manager, set this to an empty string.",
				Optional: true,
				Computed: true,
			},
			"mobile_phone": schema.StringAttribute{
				MarkdownDescription: "The primary cellular telephone number for the user.",
				Optional:            true,
				Computed:            true,
			},
			"office_location": schema.StringAttribute{
				MarkdownDescription: "The office location in the user's place of business. Maximum length is 128 characters.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(128),
				},
			},
			"on_premises_immutable_id": schema.StringAttribute{
				MarkdownDescription: "This property is used to associate an on-premises Active Directory user account to their Microsoft Entra user object. This property must be specified when creating a new user account in the Graph if you're using a federated domain for the user's userPrincipalName (UPN) property.",
				Optional:            true,
				Computed:            true,
			},
			"other_mails": schema.SetAttribute{
				MarkdownDescription: "A list of additional email addresses for the user; for example: `[\"bob@contoso.com\", \"Robert@fabrikam.com\"]`. NOTE: This property can't contain accent characters. Maximum length per value is 250 characters.",
				ElementType:         types.StringType,
				Optional:            true,
				Computed:            true,
				Validators: []validator.Set{
					setvalidator.ValueStringsAre(
						stringvalidator.LengthAtMost(250),
					),
				},
			},
			"password_policies": schema.StringAttribute{
				MarkdownDescription: "Specifies password policies for the user. This value is an enumeration with one possible value being `DisableStrongPassword`, which allows weaker passwords than the default policy to be specified. `DisablePasswordExpiration` can also be specified. The two may be specified together; for example: `DisablePasswordExpiration, DisableStrongPassword`.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.OneOf(
						"DisableStrongPassword",
						"DisablePasswordExpiration",
						"DisablePasswordExpiration, DisableStrongPassword",
					),
				},
			},
			"password_profile": schema.SingleNestedAttribute{
				MarkdownDescription: "Specifies the password profile for the user. The profile contains the user's password. " +
					"This property is required when a user is created. These fields are write-only and used only for initial user provisioning. " +
					"Password management after user creation should be handled through proper identity management workflows, not Terraform.",
				Required: true,
				Attributes: map[string]schema.Attribute{
					"password": schema.StringAttribute{
						MarkdownDescription: "The password for the user. This property is required when a user is created. " +
							"This is a write-only field used only for initial provisioning - the API never returns password values.",
						Required:  true,
						WriteOnly: true,
						Sensitive: true,
					},
					"force_change_password_next_sign_in": schema.BoolAttribute{
						MarkdownDescription: "true if the user must change their password on the next login; otherwise false. " +
							"This is a write-only field used only for initial provisioning.",
						Required:  true,
						WriteOnly: true,
					},
					"force_change_password_next_sign_in_with_mfa": schema.BoolAttribute{
						MarkdownDescription: "If true, at next sign-in, the user must perform a multi-factor authentication (MFA) before " +
							"being forced to change their password. The behavior is identical to forceChangePasswordNextSignIn except that " +
							"the user is required to first perform a multi-factor authentication before password change. " +
							"This is a write-only field used only for initial provisioning. Defaults to false if not specified.",
						Optional:  true,
						WriteOnly: true,
					},
				},
			},
			"postal_code": schema.StringAttribute{
				MarkdownDescription: "The postal code for the user's postal address. The postal code is specific to the user's country/region. In the United States of America, this attribute contains the ZIP code. Maximum length is 40 characters.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(40),
				},
			},
			"preferred_language": schema.StringAttribute{
				MarkdownDescription: "The preferred language for the user. The preferred language format is based on RFC 4646. The name combines an ISO 639 two-letter lowercase culture code associated with the language and an ISO 3166 two-letter uppercase subculture code associated with the country or region. Example: `en-US`, or `es-ES`.",
				Optional:            true,
				Computed:            true,
			},
			"proxy_addresses": schema.SetAttribute{
				MarkdownDescription: "Email addresses that also represent the user for the same mailbox. For example: `[\"SMTP: bob@contoso.com\", \"smtp: bob@sales.contoso.com\"]`. Changes to the mail property also update this collection to include the value as an SMTP address. For more information, see mail and proxyAddresses properties. The proxy address prefixed with SMTP (capitalized) is the primary proxy address. This property can't contain accent characters. Read-only in Microsoft Graph; you can only update this property through the Microsoft 365 admin center.",
				ElementType:         types.StringType,
				Optional:            true,
				Computed:            true,
			},
			"show_in_address_list": schema.BoolAttribute{
				MarkdownDescription: "`true` if the Outlook global address list should contain this user, otherwise `false`. If not set, this will be treated as `true`. For users invited through the invitation manager, this property will be set to `false`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
			},
			"state": schema.StringAttribute{
				MarkdownDescription: "The state or province in the user's address. Maximum length is 128 characters.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(128),
				},
			},
			"street_address": schema.StringAttribute{
				MarkdownDescription: "The street address of the user's place of business. Maximum length is 1024 characters.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(1024),
				},
			},
			"surname": schema.StringAttribute{
				MarkdownDescription: "The user's surname (family name or last name). Maximum length is 64 characters.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(64),
				},
			},
			"usage_location": schema.StringAttribute{
				MarkdownDescription: "A two-letter country code (ISO standard 3166). Required for users that are assigned " +
					"licenses due to legal requirements to check for availability of services in countries. " +
					"Examples include: `US`, `JP`, and `GB`. Not nullable.",
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.OneOf(
						"AD", "AE", "AF", "AG", "AI", "AL", "AM", "AO", "AQ", "AR", "AS", "AT", "AU", "AW", "AX", "AZ",
						"BA", "BB", "BD", "BE", "BF", "BG", "BH", "BI", "BJ", "BL", "BM", "BN", "BO", "BQ", "BR", "BS", "BT", "BV", "BW", "BY", "BZ",
						"CA", "CC", "CD", "CF", "CG", "CH", "CI", "CK", "CL", "CM", "CN", "CO", "CR", "CU", "CV", "CW", "CX", "CY", "CZ",
						"DE", "DJ", "DK", "DM", "DO", "DZ",
						"EC", "EE", "EG", "EH", "ER", "ES", "ET",
						"FI", "FJ", "FK", "FM", "FO", "FR",
						"GA", "GB", "GD", "GE", "GF", "GG", "GH", "GI", "GL", "GM", "GN", "GP", "GQ", "GR", "GS", "GT", "GU", "GW", "GY",
						"HK", "HM", "HN", "HR", "HT", "HU",
						"ID", "IE", "IL", "IM", "IN", "IO", "IQ", "IR", "IS", "IT",
						"JE", "JM", "JO", "JP",
						"KE", "KG", "KH", "KI", "KM", "KN", "KP", "KR", "KW", "KY", "KZ",
						"LA", "LB", "LC", "LI", "LK", "LR", "LS", "LT", "LU", "LV", "LY",
						"MA", "MC", "MD", "ME", "MF", "MG", "MH", "MK", "ML", "MM", "MN", "MO", "MP", "MQ", "MR", "MS", "MT", "MU", "MV", "MW", "MX", "MY", "MZ",
						"NA", "NC", "NE", "NF", "NG", "NI", "NL", "NO", "NP", "NR", "NU", "NZ",
						"OM",
						"PA", "PE", "PF", "PG", "PH", "PK", "PL", "PM", "PN", "PR", "PS", "PT", "PW", "PY",
						"QA",
						"RE", "RO", "RS", "RU", "RW",
						"SA", "SB", "SC", "SD", "SE", "SG", "SH", "SI", "SJ", "SK", "SL", "SM", "SN", "SO", "SR", "SS", "ST", "SV", "SX", "SY", "SZ",
						"TC", "TD", "TF", "TG", "TH", "TJ", "TK", "TL", "TM", "TN", "TO", "TR", "TT", "TV", "TW", "TZ",
						"UA", "UG", "UM", "US", "UY", "UZ",
						"VA", "VC", "VE", "VG", "VI", "VN", "VU",
						"WF", "WS",
						"YE", "YT",
						"ZA", "ZM", "ZW",
					),
				},
			},
			"user_principal_name": schema.StringAttribute{
				MarkdownDescription: "The user principal name (UPN) of the user. The UPN is an Internet-style sign-in name for the user based on the Internet standard RFC 822. " +
					"By convention, this should map to the user's email name. The general format is alias@domain, where the domain must be present in the tenant's collection of verified domains. " +
					"This property is required when a user is created. The verified domains for the tenant can be accessed from the verifiedDomains property of organization. " +
					"NOTE: This property can't contain accent characters. Only the following characters are allowed: A-Z, a-z, 0-9, ' . - _ ! # ^ ~. For the complete list of allowed characters, see username policies.",
				Required: true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(constants.UserPrincipalNameRegex),
						"must be a valid user principal name in the format alias@domain. Only the following characters are allowed in the alias: A-Z, a-z, 0-9, ' . - _ ! # ^ ~",
					),
				},
			},
			"user_type": schema.StringAttribute{
				MarkdownDescription: "A string value that can be used to classify user types in your directory, such as `Member` and `Guest`.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("Member", "Guest"),
				},
			},
			"on_premises_distinguished_name": schema.StringAttribute{
				MarkdownDescription: "Contains the on-premises Active Directory `distinguished name` or `DN`. The property is only populated for customers who are synchronizing their on-premises directory to Microsoft Entra ID via Microsoft Entra Connect. Read-only.",
				Optional:            true,
				Computed:            true,
			},
			"on_premises_domain_name": schema.StringAttribute{
				MarkdownDescription: "Contains the on-premises `domainFQDN`, also called dnsDomainName synchronized from the on-premises directory. The property is only populated for customers who are synchronizing their on-premises directory to Microsoft Entra ID via Microsoft Entra Connect. Read-only.",
				Optional:            true,
				Computed:            true,
			},
			"on_premises_last_sync_date_time": schema.StringAttribute{
				MarkdownDescription: "Indicates the last time at which the object was synced with the on-premises directory; for example: `2013-02-16T03:04:54Z`. The Timestamp type represents date and time information using ISO 8601 format and is always in UTC. For example, midnight UTC on Jan 1, 2014 is `2014-01-01T00:00:00Z`. Read-only.",
				Optional:            true,
				Computed:            true,
			},
			"on_premises_sam_account_name": schema.StringAttribute{
				MarkdownDescription: "Contains the on-premises `sAMAccountName` synchronized from the on-premises directory. The property is only populated for customers who are synchronizing their on-premises directory to Microsoft Entra ID via Microsoft Entra Connect. Read-only.",
				Optional:            true,
				Computed:            true,
			},
			"on_premises_security_identifier": schema.StringAttribute{
				MarkdownDescription: "Contains the on-premises security identifier (SID) for the user that was synchronized from on-premises to the cloud. Read-only.",
				Optional:            true,
				Computed:            true,
			},
			"on_premises_sync_enabled": schema.BoolAttribute{
				MarkdownDescription: "`true` if this user object is currently being synced from an on-premises Active Directory (AD); otherwise, the user isn't being synced and can be managed in Microsoft Entra ID. Read-only. The value is `null` for cloud-only users.",
				Optional:            true,
				Computed:            true,
			},
			"on_premises_user_principal_name": schema.StringAttribute{
				MarkdownDescription: "Contains the on-premises `userPrincipalName` synchronized from the on-premises directory. The property is only populated for customers who are synchronizing their on-premises directory to Microsoft Entra ID via Microsoft Entra Connect. Read-only.",
				Optional:            true,
				Computed:            true,
			},
			"preferred_data_location": schema.StringAttribute{
				MarkdownDescription: "The preferred data location for the user. For more information, see OneDrive Online Multi-Geo.",
				Optional:            true,
				Computed:            true,
			},
			"preferred_name": schema.StringAttribute{
				MarkdownDescription: "The preferred name for the user. **Not Supported.** This attribute returns an empty string.",
				Optional:            true,
				Computed:            true,
			},
			"security_identifier": schema.StringAttribute{
				MarkdownDescription: "Security identifier (SID) of the user, used in Windows scenarios. Read-only.",
				Optional:            true,
				Computed:            true,
			},
			"sign_in_sessions_valid_from_date_time": schema.StringAttribute{
				MarkdownDescription: "Any refresh tokens or sessions tokens (session cookies) issued before this time are invalid, and applications get an error when using an invalid refresh or sessions token to acquire a delegated access token (to access APIs such as Microsoft Graph). If this happens, the application needs to acquire a new refresh token by making a request to the authorize endpoint. Read-only. Use revokeSignInSessions to reset.",
				Optional:            true,
				Computed:            true,
			},
			"hard_delete": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(false),
				MarkdownDescription: "When `true`, the user will be permanently deleted (hard delete) during destroy. " +
					"When `false` (default), the user will only be soft deleted and moved to the deleted items container where it can be restored within 30 days. " +
					"Note: This field defaults to `false` on import since the API does not return this value.",
			},
			"timeouts": commonschema.ResourceTimeouts(ctx),
		},
	}
}
