// REF: https://learn.microsoft.com/en-us/graph/api/resources/user?view=graph-rest-beta
package graphBetaUsersUser

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// UserResourceModel represents the schema for the User resource
type UserResourceModel struct {
	ID                              types.String                 `tfsdk:"id"`
	AboutMe                         types.String                 `tfsdk:"about_me"`
	AccountEnabled                  types.Bool                   `tfsdk:"account_enabled"`
	AgeGroup                        types.String                 `tfsdk:"age_group"`
	BusinessPhones                  types.Set                    `tfsdk:"business_phones"`
	City                            types.String                 `tfsdk:"city"`
	CompanyName                     types.String                 `tfsdk:"company_name"`
	ConsentProvidedForMinor         types.String                 `tfsdk:"consent_provided_for_minor"`
	Country                         types.String                 `tfsdk:"country"`
	CreatedDateTime                 types.String                 `tfsdk:"created_date_time"`
	CreationType                    types.String                 `tfsdk:"creation_type"`
	CustomSecurityAttributes        []CustomSecurityAttributeSet `tfsdk:"custom_security_attributes"`
	DeletedDateTime                 types.String                 `tfsdk:"deleted_date_time"`
	Department                      types.String                 `tfsdk:"department"`
	DisplayName                     types.String                 `tfsdk:"display_name"`
	EmployeeHireDate                types.String                 `tfsdk:"employee_hire_date"`
	EmployeeId                      types.String                 `tfsdk:"employee_id"`
	EmployeeType                    types.String                 `tfsdk:"employee_type"`
	ExternalUserState               types.String                 `tfsdk:"external_user_state"`
	ExternalUserStateChangeDateTime types.String                 `tfsdk:"external_user_state_change_date_time"`
	FaxNumber                       types.String                 `tfsdk:"fax_number"`
	GivenName                       types.String                 `tfsdk:"given_name"`
	JobTitle                        types.String                 `tfsdk:"job_title"`
	Mail                            types.String                 `tfsdk:"mail"`
	MailNickname                    types.String                 `tfsdk:"mail_nickname"`
	ManagerId                       types.String                 `tfsdk:"manager_id"`
	MobilePhone                     types.String                 `tfsdk:"mobile_phone"`
	OfficeLocation                  types.String                 `tfsdk:"office_location"`
	OnPremisesDistinguishedName     types.String                 `tfsdk:"on_premises_distinguished_name"`
	OnPremisesDomainName            types.String                 `tfsdk:"on_premises_domain_name"`
	OnPremisesImmutableId           types.String                 `tfsdk:"on_premises_immutable_id"`
	OnPremisesLastSyncDateTime      types.String                 `tfsdk:"on_premises_last_sync_date_time"`
	OnPremisesSamAccountName        types.String                 `tfsdk:"on_premises_sam_account_name"`
	OnPremisesSecurityIdentifier    types.String                 `tfsdk:"on_premises_security_identifier"`
	OnPremisesSyncEnabled           types.Bool                   `tfsdk:"on_premises_sync_enabled"`
	OnPremisesUserPrincipalName     types.String                 `tfsdk:"on_premises_user_principal_name"`
	OtherMails                      types.Set                    `tfsdk:"other_mails"`
	PasswordPolicies                types.String                 `tfsdk:"password_policies"`
	PasswordProfile                 *PasswordProfile             `tfsdk:"password_profile"`
	PostalCode                      types.String                 `tfsdk:"postal_code"`
	PreferredDataLocation           types.String                 `tfsdk:"preferred_data_location"`
	PreferredLanguage               types.String                 `tfsdk:"preferred_language"`
	PreferredName                   types.String                 `tfsdk:"preferred_name"`
	ProxyAddresses                  types.Set                    `tfsdk:"proxy_addresses"`
	SecurityIdentifier              types.String                 `tfsdk:"security_identifier"`
	ShowInAddressList               types.Bool                   `tfsdk:"show_in_address_list"`
	SignInSessionsValidFromDateTime types.String                 `tfsdk:"sign_in_sessions_valid_from_date_time"`
	State                           types.String                 `tfsdk:"state"`
	StreetAddress                   types.String                 `tfsdk:"street_address"`
	Surname                         types.String                 `tfsdk:"surname"`
	UsageLocation                   types.String                 `tfsdk:"usage_location"`
	UserPrincipalName               types.String                 `tfsdk:"user_principal_name"`
	UserType                        types.String                 `tfsdk:"user_type"`
	HardDelete                      types.Bool                   `tfsdk:"hard_delete"`
	Timeouts                        timeouts.Value               `tfsdk:"timeouts"`
}

// PasswordProfile represents the password profile for a user
type PasswordProfile struct {
	Password                             types.String `tfsdk:"password"`
	ForceChangePasswordNextSignIn        types.Bool   `tfsdk:"force_change_password_next_sign_in"`
	ForceChangePasswordNextSignInWithMfa types.Bool   `tfsdk:"force_change_password_next_sign_in_with_mfa"`
}

// CustomSecurityAttributeSet represents a set of custom security attributes
type CustomSecurityAttributeSet struct {
	AttributeSet types.String                  `tfsdk:"attribute_set"`
	Attributes   []CustomSecurityAttributeItem `tfsdk:"attributes"`
}

// CustomSecurityAttributeItem represents a single custom security attribute within a set
type CustomSecurityAttributeItem struct {
	Name         types.String `tfsdk:"name"`
	StringValue  types.String `tfsdk:"string_value"`
	IntValue     types.Int32  `tfsdk:"int_value"`
	BoolValue    types.Bool   `tfsdk:"bool_value"`
	StringValues types.Set    `tfsdk:"string_values"`
	IntValues    types.Set    `tfsdk:"int_values"`
}
