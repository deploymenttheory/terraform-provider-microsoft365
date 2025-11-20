package graphBetaUsersUser

import (
	"context"
	"regexp"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
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
		},
		WritePermissions: []string{
			"User.ReadWrite.All",
			"Directory.ReadWrite.All",
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

// ImportState imports the resource state.
func (r *UserResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Schema defines the schema for the resource.
func (r *UserResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages Microsoft 365 users using the `/users` endpoint. The user resource lets admins specify user preferences for languages and date/time formats for the user's primary Exchange mailboxes and Microsoft Entra profile.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "String (identifier)",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"about_me": schema.StringAttribute{
				MarkdownDescription: "A freeform text entry field for the user to describe themselves.",
				Optional:            true,
				Computed:            true,
			},
			"account_enabled": schema.BoolAttribute{
				MarkdownDescription: "true if the account is enabled; otherwise, false.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
			},
			"age_group": schema.StringAttribute{
				MarkdownDescription: "Sets the age group of the user. Allowed values: null, minor, notAdult, adult.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("null", "minor", "notAdult", "adult"),
				},
			},
			"business_phones": schema.SetAttribute{
				MarkdownDescription: "The telephone numbers for the user.",
				ElementType:         types.StringType,
				Optional:            true,
				Computed:            true,
			},
			"city": schema.StringAttribute{
				MarkdownDescription: "The city in which the user is located.",
				Optional:            true,
				Computed:            true,
			},
			"company_name": schema.StringAttribute{
				MarkdownDescription: "The company name which the user is associated.",
				Optional:            true,
				Computed:            true,
			},
			"consent_provided_for_minor": schema.StringAttribute{
				MarkdownDescription: "Sets whether consent has been obtained for minors.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("null", "granted", "denied", "notRequired"),
				},
			},
			"country": schema.StringAttribute{
				MarkdownDescription: "The country/region in which the user is located; for example, 'US' or 'UK'.",
				Optional:            true,
				Computed:            true,
			},
			"created_date_time": schema.StringAttribute{
				MarkdownDescription: "The created date of the user object.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"creation_type": schema.StringAttribute{
				MarkdownDescription: "Indicates whether the user account was created as a regular school or work account (null), an external account (Invitation), a local account for an Azure Active Directory B2C tenant (LocalAccount) or self-service sign-up using email verification (EmailVerified).",
				Computed:            true,
			},
			"deleted_date_time": schema.StringAttribute{
				MarkdownDescription: "The date and time when the user was deleted.",
				Computed:            true,
			},
			"department": schema.StringAttribute{
				MarkdownDescription: "The name for the department in which the user works.",
				Optional:            true,
				Computed:            true,
			},
			"display_name": schema.StringAttribute{
				MarkdownDescription: "The name displayed in the address book for the user. This is usually the combination of the user's first name, middle initial and last name.",
				Required:            true,
			},
			"employee_hire_date": schema.StringAttribute{
				MarkdownDescription: "The date and time when the user was hired or will start work in case of a future hire.",
				Optional:            true,
				Computed:            true,
			},
			"employee_id": schema.StringAttribute{
				MarkdownDescription: "The employee identifier assigned to the user by the organization.",
				Optional:            true,
				Computed:            true,
			},
			"employee_type": schema.StringAttribute{
				MarkdownDescription: "Captures enterprise worker type. For example, Employee, Contractor, Consultant, or Vendor.",
				Optional:            true,
				Computed:            true,
			},
			"external_user_state": schema.StringAttribute{
				MarkdownDescription: "For an external user invited to the tenant, this property represents the invited user's invitation status. Possible values: PendingAcceptance, Accepted.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("PendingAcceptance", "Accepted"),
				},
			},
			"external_user_state_change_date_time": schema.StringAttribute{
				MarkdownDescription: "Shows the timestamp for the latest change to the externalUserState property.",
				Computed:            true,
			},
			"fax_number": schema.StringAttribute{
				MarkdownDescription: "The fax number of the user.",
				Optional:            true,
				Computed:            true,
			},
			"given_name": schema.StringAttribute{
				MarkdownDescription: "The given name (first name) of the user.",
				Optional:            true,
				Computed:            true,
			},
			"job_title": schema.StringAttribute{
				MarkdownDescription: "The user's job title.",
				Optional:            true,
				Computed:            true,
			},
			"mail": schema.StringAttribute{
				MarkdownDescription: "The SMTP address for the user.",
				Optional:            true,
				Computed:            true,
			},
			"mail_nickname": schema.StringAttribute{
				MarkdownDescription: "The mail alias for the user.",
				Required:            true,
			},
			"mobile_phone": schema.StringAttribute{
				MarkdownDescription: "The primary cellular telephone number for the user.",
				Optional:            true,
				Computed:            true,
			},
			"office_location": schema.StringAttribute{
				MarkdownDescription: "The office location in the user's place of business.",
				Optional:            true,
				Computed:            true,
			},
			"on_premises_immutable_id": schema.StringAttribute{
				MarkdownDescription: "This property is used to associate an on-premises Active Directory user account to their Microsoft Entra ID user object. This property must be specified when creating a new user account in a federated domain if you are not using the userPrincipalName property.",
				Optional:            true,
				Computed:            true,
			},
			"other_mails": schema.SetAttribute{
				MarkdownDescription: "Additional email addresses for the user.",
				ElementType:         types.StringType,
				Optional:            true,
				Computed:            true,
			},
			"password_policies": schema.StringAttribute{
				MarkdownDescription: "Specifies password policies for the user. This value is an enumeration with one possible value being 'DisableStrongPassword', which allows weaker passwords than the default policy to be specified. 'DisablePasswordExpiration' can also be specified. The two may be specified together; for example: 'DisablePasswordExpiration, DisableStrongPassword'.",
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
				MarkdownDescription: "Specifies the password profile for the user. The profile contains the user's password. This property is required when a user is created.",
				Required:            true,
				Attributes: map[string]schema.Attribute{
					"password": schema.StringAttribute{
						MarkdownDescription: "The password for the user. This property is required when a user is created. It can be updated, but the user will be required to change the password on the next login.",
						Required:            true,
						Sensitive:           true,
					},
					"force_change_password_next_sign_in": schema.BoolAttribute{
						MarkdownDescription: "true if the user must change their password on the next login; otherwise false.",
						Optional:            true,
						Computed:            true,
						Default:             booldefault.StaticBool(false),
					},
					"force_change_password_next_sign_in_with_mfa": schema.BoolAttribute{
						MarkdownDescription: "If true, at next sign-in, the user must perform a multi-factor authentication (MFA) before being forced to change their password. The behavior is identical to forceChangePasswordNextSignIn except that the user is required to first perform a multi-factor authentication before password change.",
						Optional:            true,
						Computed:            true,
						Default:             booldefault.StaticBool(false),
					},
				},
			},
			"postal_code": schema.StringAttribute{
				MarkdownDescription: "The postal code for the user's postal address. The postal code is specific to the user's country/region. In the United States of America, this attribute contains the ZIP code.",
				Optional:            true,
				Computed:            true,
			},
			"preferred_language": schema.StringAttribute{
				MarkdownDescription: "The preferred language for the user. Should follow ISO 639-1 Code; for example 'en-US'.",
				Optional:            true,
				Computed:            true,
			},
			"proxy_addresses": schema.SetAttribute{
				MarkdownDescription: "For example: ['SMTP: bob@contoso.com', 'smtp: bob@sales.contoso.com']. Changes to the mail property will also update this collection to include the value as an SMTP address.",
				ElementType:         types.StringType,
				Optional:            true,
				Computed:            true,
			},
			"show_in_address_list": schema.BoolAttribute{
				MarkdownDescription: "true if the Outlook global address list should contain this user, otherwise false.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
			},
			"state": schema.StringAttribute{
				MarkdownDescription: "The state or province in the user's address.",
				Optional:            true,
				Computed:            true,
			},
			"street_address": schema.StringAttribute{
				MarkdownDescription: "The street address of the user's place of business.",
				Optional:            true,
				Computed:            true,
			},
			"surname": schema.StringAttribute{
				MarkdownDescription: "The user's surname (family name or last name).",
				Optional:            true,
				Computed:            true,
			},
			"usage_location": schema.StringAttribute{
				MarkdownDescription: "A two letter country code (ISO standard 3166). Required for users that will be assigned licenses due to legal requirement to check for availability of services in countries.",
				Optional:            true,
				Computed:            true,
			},
			"user_principal_name": schema.StringAttribute{
				MarkdownDescription: "The user principal name (someuser@contoso.com). It's an Internet-style login name for the user based on the Internet standard RFC 822. " +
					"By convention, this should map to the user's email name. The general format is alias@domain, where domain must be present in the tenant's collection of verified domains. " +
					"The verified domains for the tenant can be accessed from the verifiedDomains property of organization. " +
					"NOTE: This property cannot contain accent characters. Only the following characters are allowed A - Z, a - z, 0 - 9, ' . - _ ! # ^ ~. For the complete list of allowed characters, see https://learn.microsoft.com/en-us/entra/identity/authentication/concept-sspr-policy?tabs=ms-powershell#userprincipalname-policies-that-apply-to-all-user-accounts.",
				Required: true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(constants.UserPrincipalNameRegex),
						"must be a valid user principal name in the format alias@domain. Only the following characters are allowed in the alias: A-Z, a-z, 0-9, ' . - _ ! # ^ ~",
					),
				},
			},
			"user_type": schema.StringAttribute{
				MarkdownDescription: "A string value that can be used to classify user types in your directory, such as 'Member' and 'Guest'.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("Member", "Guest"),
				},
			},
			"identities": schema.SetNestedAttribute{
				MarkdownDescription: "Identities that can be used to sign in to this user account. An identity can be provided by Microsoft " +
					"(also known as a local account), by organizations, or by social identity providers such as Facebook, Google, and Microsoft, and tied to a user account.",
				Optional: true,
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"sign_in_type": schema.StringAttribute{
							MarkdownDescription: "The type of sign-in used by the identity. The possible values are: emailAddress, userName, federated, or userPrincipalName.",
							Optional:            true,
							Computed:            true,
							Validators: []validator.String{
								stringvalidator.OneOf(
									"emailAddress",
									"userName",
									"federated",
									"userPrincipalName",
								),
							},
						},
						"issuer": schema.StringAttribute{
							MarkdownDescription: "The name of the identity provider. Typically, this is the tenant domain name, such as 'contoso.onmicrosoft.com'.",
							Optional:            true,
							Computed:            true,
						},
						"issuer_assigned_id": schema.StringAttribute{
							MarkdownDescription: "The unique identifier assigned to the user by the issuer. Typically, this is the user's email address, such as 'jane.smith@contoso.com'.",
							Optional:            true,
							Computed:            true,
						},
					},
				},
			},
			"im_addresses": schema.SetAttribute{
				MarkdownDescription: "The instant message voice over IP (VOIP) session initiation protocol (SIP) addresses for the user.",
				ElementType:         types.StringType,
				Optional:            true,
				Computed:            true,
			},
			"on_premises_distinguished_name": schema.StringAttribute{
				MarkdownDescription: "Contains the on-premises Active Directory distinguished name or DN. The property is only populated for customers who are synchronizing their on-premises directory to Microsoft Entra ID via Microsoft Entra Connect.",
				Optional:            true,
				Computed:            true,
			},
			"on_premises_domain_name": schema.StringAttribute{
				MarkdownDescription: "Contains the on-premises domainFQDN, also called dnsDomainName synchronized from the on-premises directory.",
				Optional:            true,
				Computed:            true,
			},
			"on_premises_last_sync_date_time": schema.StringAttribute{
				MarkdownDescription: "Indicates the last time at which the object was synchronized with the on-premises directory; the property is only populated for customers who are synchronizing their on-premises directory to Microsoft Entra ID via Microsoft Entra Connect.",
				Optional:            true,
				Computed:            true,
			},
			"on_premises_sam_account_name": schema.StringAttribute{
				MarkdownDescription: "Contains the on-premises sAMAccountName synchronized from the on-premises directory.",
				Optional:            true,
				Computed:            true,
			},
			"on_premises_security_identifier": schema.StringAttribute{
				MarkdownDescription: "Contains the on-premises security identifier (SID) for the user that was synchronized from on-premises to the cloud.",
				Optional:            true,
				Computed:            true,
			},
			"on_premises_sync_enabled": schema.BoolAttribute{
				MarkdownDescription: "true if this object is synchronized from an on-premises directory; false if this object was originally synchronized from an on-premises directory but is no longer synchronized; null if this object has never been synchronized from an on-premises directory (default).",
				Optional:            true,
				Computed:            true,
			},
			"on_premises_user_principal_name": schema.StringAttribute{
				MarkdownDescription: "Contains the on-premises userPrincipalName synchronized from the on-premises directory.",
				Optional:            true,
				Computed:            true,
			},
			"preferred_data_location": schema.StringAttribute{
				MarkdownDescription: "The preferred data location for the user.",
				Optional:            true,
				Computed:            true,
			},
			"preferred_name": schema.StringAttribute{
				MarkdownDescription: "The preferred name for the user.",
				Optional:            true,
				Computed:            true,
			},
			"security_identifier": schema.StringAttribute{
				MarkdownDescription: "Security identifier (SID) of the user, used for Azure AD authentication.",
				Optional:            true,
				Computed:            true,
			},
			"sign_in_sessions_valid_from_date_time": schema.StringAttribute{
				MarkdownDescription: "Any refresh tokens or sessions tokens (session cookies) issued before this time are invalid, and applications will get an error when using an invalid refresh or sessions token to acquire a delegated access token (to access APIs such as Microsoft Graph).",
				Optional:            true,
				Computed:            true,
			},
			"timeouts": commonschema.Timeouts(ctx),
		},
	}
}
