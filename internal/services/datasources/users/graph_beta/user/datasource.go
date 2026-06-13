package graphBetaUser

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	"github.com/hashicorp/terraform-plugin-framework-validators/boolvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	DataSourceName = "microsoft365_graph_beta_users_user"
	ReadTimeout    = 180
)

var (
	_ datasource.DataSource              = &UserDataSource{}
	_ datasource.DataSourceWithConfigure = &UserDataSource{}
)

func NewUserDataSource() datasource.DataSource {
	return &UserDataSource{
		ReadPermissions: []string{
			"User.Read.All",
			"Directory.Read.All",
		},
	}
}

type UserDataSource struct {
	client *msgraphbetasdk.GraphServiceClient

	ReadPermissions []string
}

func (d *UserDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = DataSourceName
}

func (d *UserDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.client = client.SetGraphBetaClientForDataSource(ctx, req, resp, DataSourceName)
}

func (d *UserDataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Retrieves Microsoft Entra Users using the `/users` endpoint. " +
			"Supports flexible lookup by object ID, display name, employee ID, given name, " +
			"user principal name, on-premises immutable ID, on-premises distinguished name, " +
			"or a custom OData query. Can also list all users in the tenant.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The unique identifier for the data source. This is a placeholder attribute required by Terraform.",
			},
			"object_id": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The unique object identifier of the user in Microsoft Entra ID. Conflicts with other lookup attributes.",
				Validators: []validator.String{
					stringvalidator.ConflictsWith(
						path.MatchRoot("display_name"),
						path.MatchRoot("employee_id"),
						path.MatchRoot("given_name"),
						path.MatchRoot("user_principal_name"),
						path.MatchRoot("on_premises_immutable_id"),
						path.MatchRoot("on_premises_distinguished_name"),
						path.MatchRoot("list_all"),
						path.MatchRoot("odata_query"),
					),
					stringvalidator.AtLeastOneOf(
						path.MatchRoot("object_id"),
						path.MatchRoot("display_name"),
						path.MatchRoot("employee_id"),
						path.MatchRoot("given_name"),
						path.MatchRoot("user_principal_name"),
						path.MatchRoot("on_premises_immutable_id"),
						path.MatchRoot("on_premises_distinguished_name"),
						path.MatchRoot("list_all"),
						path.MatchRoot("odata_query"),
					),
				},
			},
			"display_name": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The display name of the user. Conflicts with other lookup attributes.",
				Validators: []validator.String{
					stringvalidator.ConflictsWith(
						path.MatchRoot("object_id"),
						path.MatchRoot("list_all"),
						path.MatchRoot("odata_query"),
					),
				},
			},
			"employee_id": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The employee identifier assigned to the user by the organization. Conflicts with other lookup attributes.",
				Validators: []validator.String{
					stringvalidator.ConflictsWith(
						path.MatchRoot("object_id"),
						path.MatchRoot("list_all"),
						path.MatchRoot("odata_query"),
					),
				},
			},
			"given_name": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The given name (first name) of the user. Conflicts with other lookup attributes.",
				Validators: []validator.String{
					stringvalidator.ConflictsWith(
						path.MatchRoot("object_id"),
						path.MatchRoot("list_all"),
						path.MatchRoot("odata_query"),
					),
				},
			},
			"user_principal_name": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The user principal name (UPN) of the user. Conflicts with other lookup attributes.",
				Validators: []validator.String{
					stringvalidator.ConflictsWith(
						path.MatchRoot("object_id"),
						path.MatchRoot("list_all"),
						path.MatchRoot("odata_query"),
					),
				},
			},
			"on_premises_immutable_id": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The on-premises immutable ID (sourceAnchor) used to associate an on-premises Active Directory user account. Conflicts with other lookup attributes.",
				Validators: []validator.String{
					stringvalidator.ConflictsWith(
						path.MatchRoot("object_id"),
						path.MatchRoot("list_all"),
						path.MatchRoot("odata_query"),
					),
				},
			},
			"on_premises_distinguished_name": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The on-premises Active Directory distinguished name (DN) of the user. Conflicts with other lookup attributes.",
				Validators: []validator.String{
					stringvalidator.ConflictsWith(
						path.MatchRoot("object_id"),
						path.MatchRoot("list_all"),
						path.MatchRoot("odata_query"),
					),
				},
			},
			"odata_query": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Custom OData filter expression for advanced queries (e.g., `accountEnabled eq true and userType eq 'Member'`). Conflicts with specific lookup attributes.",
				Validators: []validator.String{
					stringvalidator.ConflictsWith(
						path.MatchRoot("object_id"),
						path.MatchRoot("display_name"),
						path.MatchRoot("employee_id"),
						path.MatchRoot("given_name"),
						path.MatchRoot("user_principal_name"),
						path.MatchRoot("on_premises_immutable_id"),
						path.MatchRoot("on_premises_distinguished_name"),
						path.MatchRoot("list_all"),
					),
				},
			},
			"list_all": schema.BoolAttribute{
				Optional:            true,
				MarkdownDescription: "Retrieve all users in the tenant. Conflicts with specific lookup attributes.",
				Validators: []validator.Bool{
					boolvalidator.ConflictsWith(
						path.MatchRoot("object_id"),
						path.MatchRoot("display_name"),
						path.MatchRoot("employee_id"),
						path.MatchRoot("given_name"),
						path.MatchRoot("user_principal_name"),
						path.MatchRoot("on_premises_immutable_id"),
						path.MatchRoot("on_premises_distinguished_name"),
						path.MatchRoot("odata_query"),
					),
				},
			},
			"items": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "List of users matching the query criteria.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The unique identifier for the user object.",
						},
						"about_me": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "A freeform text entry field for the user to describe themselves.",
						},
						"account_enabled": schema.BoolAttribute{
							Computed:            true,
							MarkdownDescription: "true if the account is enabled; otherwise, false.",
						},
						"age_group": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "Sets the age group of the user (null, Minor, NotAdult, Adult).",
						},
						"business_phones": schema.ListAttribute{
							ElementType:         types.StringType,
							Computed:            true,
							MarkdownDescription: "The telephone numbers for the user.",
						},
						"city": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The city where the user is located.",
						},
						"company_name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The company name which the user is associated with.",
						},
						"consent_provided_for_minor": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "Sets whether consent was obtained for minors (null, Granted, Denied, NotRequired).",
						},
						"country": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The country/region where the user is located.",
						},
						"created_date_time": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The date and time the user was created.",
						},
						"creation_type": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "Indicates whether the user account was created as a regular school or work account, an external account, etc.",
						},
						"deleted_date_time": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The date and time the user was deleted.",
						},
						"department": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The name of the department in which the user works.",
						},
						"display_name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The name displayed in the address book for the user.",
						},
						"employee_hire_date": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The date and time when the user was hired or will start work.",
						},
						"employee_id": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The employee identifier assigned to the user by the organization.",
						},
						"employee_type": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "Captures enterprise worker type (Employee, Contractor, Consultant, Vendor, etc.).",
						},
						"external_user_state": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "For an external user invited to the tenant, this represents the invitation status (PendingAcceptance, Accepted).",
						},
						"external_user_state_change_date_time": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "Shows the timestamp for the latest change to the external_user_state property.",
						},
						"fax_number": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The fax number of the user.",
						},
						"given_name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The given name (first name) of the user.",
						},
						"job_title": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The user's job title.",
						},
						"mail": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The SMTP address for the user.",
						},
						"mail_nickname": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The mail alias for the user.",
						},
						"mobile_phone": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The primary cellular telephone number for the user.",
						},
						"office_location": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The office location in the user's place of business.",
						},
						"on_premises_distinguished_name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "Contains the on-premises Active Directory distinguished name (DN).",
						},
						"on_premises_domain_name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "Contains the on-premises domainFQDN, also called dnsDomainName synchronized from the on-premises directory.",
						},
						"on_premises_immutable_id": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The on-premises immutable ID (sourceAnchor) used to associate an on-premises Active Directory user account.",
						},
						"on_premises_last_sync_date_time": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "Indicates the last time at which the object was synced with the on-premises directory.",
						},
						"on_premises_sam_account_name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "Contains the on-premises samAccountName synchronized from the on-premises directory.",
						},
						"on_premises_security_identifier": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "Contains the on-premises security identifier (SID) for the user that was synchronized from on-premises to the cloud.",
						},
						"on_premises_sync_enabled": schema.BoolAttribute{
							Computed:            true,
							MarkdownDescription: "true if this user object is currently being synced from an on-premises Active Directory.",
						},
						"on_premises_user_principal_name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "Contains the on-premises userPrincipalName synchronized from the on-premises directory.",
						},
						"other_mails": schema.ListAttribute{
							ElementType:         types.StringType,
							Computed:            true,
							MarkdownDescription: "A list of additional email addresses for the user.",
						},
						"password_policies": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "Specifies password policies for the user (DisableStrongPassword, DisablePasswordExpiration).",
						},
						"postal_code": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The postal code for the user's postal address.",
						},
						"preferred_data_location": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The preferred data location for the user.",
						},
						"preferred_language": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The preferred language for the user, in ISO 639-1 format.",
						},
						"preferred_name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The preferred name for the user.",
						},
						"proxy_addresses": schema.ListAttribute{
							ElementType:         types.StringType,
							Computed:            true,
							MarkdownDescription: "For example: [\"SMTP: bob@contoso.com\", \"smtp: bob@sales.contoso.com\"].",
						},
						"security_identifier": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "Security identifier (SID) of the user, used in Windows scenarios.",
						},
						"show_in_address_list": schema.BoolAttribute{
							Computed:            true,
							MarkdownDescription: "Do not use in Microsoft Graph. Manage this property through the Microsoft 365 admin center instead.",
						},
						"sign_in_sessions_valid_from_date_time": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "Any refresh tokens or session tokens issued before this time are invalid.",
						},
						"state": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The state or province in the user's address.",
						},
						"street_address": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The street address of the user's place of business.",
						},
						"surname": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The user's surname (family name or last name).",
						},
						"usage_location": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "A two letter country code (ISO standard 3166), required for users that are assigned licenses.",
						},
						"user_principal_name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The user principal name (UPN) of the user.",
						},
						"user_type": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "A string value that can be used to classify user types in your directory (Member, Guest).",
						},
					},
				},
			},
			"timeouts": commonschema.DatasourceTimeouts(ctx),
		},
	}
}
