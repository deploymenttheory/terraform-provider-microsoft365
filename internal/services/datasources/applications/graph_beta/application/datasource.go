package graphBetaApplication

import (
	"context"
	"regexp"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	validate "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/validate/attribute"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	DataSourceName = "microsoft365_graph_beta_applications_application"
	ReadTimeout    = 180
)

var (
	_ datasource.DataSource              = &ApplicationDataSource{}
	_ datasource.DataSourceWithConfigure = &ApplicationDataSource{}
)

func NewApplicationDataSource() datasource.DataSource {
	return &ApplicationDataSource{
		ReadPermissions: []string{
			"Application.Read.All",
			"Directory.Read.All",
		},
	}
}

type ApplicationDataSource struct {
	client *msgraphbetasdk.GraphServiceClient

	ReadPermissions []string
}

func (d *ApplicationDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = DataSourceName
}

func (d *ApplicationDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.client = client.SetGraphBetaClientForDataSource(ctx, req, resp, DataSourceName)
}

func (d *ApplicationDataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Retrieves information about a Microsoft Entra ID (Azure AD) application using the `/applications` endpoint. This data source is used to query application details by ID, app ID, display name, or advanced OData filtering.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Unique identifier for the application object. This property is referred to as **Object ID** in the Microsoft Entra admin center. Key. Not nullable. Read-only. Supports `$filter` (`eq`, `ne`, `not`, `in`).",
			},
			"object_id": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The object ID of the application. One of `object_id`, `app_id`, `display_name`, or `odata_query` must be specified.",
				Validators: []validator.String{
					stringvalidator.AtLeastOneOf(
						path.MatchRoot("app_id"),
						path.MatchRoot("display_name"),
						path.MatchRoot("odata_query"),
					),
					stringvalidator.ConflictsWith(
						path.MatchRoot("app_id"),
						path.MatchRoot("display_name"),
						path.MatchRoot("odata_query"),
					),
					stringvalidator.RegexMatches(
						regexp.MustCompile(constants.GuidRegex),
						"must be a valid GUID in the format 00000000-0000-0000-0000-000000000000",
					),
				},
			},
			"app_id": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The unique identifier for the application that is assigned by Microsoft Entra ID. Not nullable. Read-only. Alternate key. Supports `$filter` (`eq`). One of `object_id`, `app_id`, `display_name`, or `odata_query` must be specified.",
				Validators: []validator.String{
					stringvalidator.AtLeastOneOf(
						path.MatchRoot("object_id"),
						path.MatchRoot("display_name"),
						path.MatchRoot("odata_query"),
					),
					stringvalidator.ConflictsWith(
						path.MatchRoot("object_id"),
						path.MatchRoot("display_name"),
						path.MatchRoot("odata_query"),
					),
					stringvalidator.RegexMatches(
						regexp.MustCompile(constants.GuidRegex),
						"must be a valid GUID in the format 00000000-0000-0000-0000-000000000000",
					),
				},
			},
			"display_name": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The display name for the application. Maximum length is 256 characters. Supports `$filter` (`eq`, `ne`, `not`, `ge`, `le`, `in`, `startsWith`, and `eq` on null values), `$search`, and `$orderby`. One of `object_id`, `app_id`, `display_name`, or `odata_query` must be specified.",
				Validators: []validator.String{
					stringvalidator.AtLeastOneOf(
						path.MatchRoot("object_id"),
						path.MatchRoot("app_id"),
						path.MatchRoot("odata_query"),
					),
					stringvalidator.ConflictsWith(
						path.MatchRoot("object_id"),
						path.MatchRoot("app_id"),
						path.MatchRoot("odata_query"),
					),
					validate.StringLengthAtMost(256),
				},
			},
			"odata_query": schema.StringAttribute{
				Optional: true,
				MarkdownDescription: "Custom OData filter query. " +
					"Use this for advanced filtering when the standard lookup attributes don't meet your needs. " +
					"Cannot be combined with `object_id`, `app_id`, or `display_name`. " +
					"Example: `displayName eq 'My Application' and signInAudience eq 'AzureADMyOrg'`",
				Validators: []validator.String{
					stringvalidator.AtLeastOneOf(
						path.MatchRoot("object_id"),
						path.MatchRoot("app_id"),
						path.MatchRoot("display_name"),
					),
					stringvalidator.ConflictsWith(
						path.MatchRoot("object_id"),
						path.MatchRoot("app_id"),
						path.MatchRoot("display_name"),
					),
				},
			},
			"description": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Free text field to provide a description of the application object to end users. The maximum allowed size is 1,024 characters. Returned by default. Supports `$filter` (`eq`, `ne`, `not`, `ge`, `le`, `startsWith`) and `$search`.",
			},
			"sign_in_audience": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Specifies the Microsoft accounts that are supported for the current application. The possible values are: `AzureADMyOrg` (default), `AzureADMultipleOrgs`, `AzureADandPersonalMicrosoftAccount`, and `PersonalMicrosoftAccount`. The value of this object also limits the number of permissions an app can request. For more information, see Limits on requested permissions per app. The value for this property has implications on other app object properties. As a result, if you change this property, you may need to change other properties first. For more information, see Validation differences for signInAudience. Supports `$filter` (`eq`, `ne`, `not`).",
			},
			"identifier_uris": schema.SetAttribute{
				Computed:            true,
				ElementType:         types.StringType,
				MarkdownDescription: "Also known as App ID URI, this value is set when an application is used as a resource app. The identifierUris acts as the prefix for the scopes you reference in your API's code, and it must be globally unique across Microsoft Entra ID. For more information on valid identifierUris patterns and best practices, see Microsoft Entra application registration security best practices. Not nullable. Supports `$filter` (`eq`, `ne`, `ge`, `le`, `startsWith`).",
			},
			"group_membership_claims": schema.SetAttribute{
				Computed:            true,
				ElementType:         types.StringType,
				MarkdownDescription: "Configures the groups claim issued in a user or OAuth 2.0 access token that the application expects. To set this attribute, use one of the following string values: `None`, `SecurityGroup` (for security groups and Microsoft Entra roles), `All` (this gets all security groups, distribution groups, and Microsoft Entra directory roles that the signed-in user is a member of).",
			},
			"notes": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Notes relevant for the management of the application.",
			},
			"is_device_only_auth_supported": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "Specifies whether this application supports device authentication without a user. The default is false.",
			},
			"is_fallback_public_client": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "Specifies the fallback application type as public client, such as an installed application running on a mobile device. The default value is false, which means the fallback application type is confidential client such as a web app. There are certain scenarios where Microsoft Entra ID can't determine the client application type. For example, the ROPC flow where the application is configured without specifying a redirect URI. In those cases Microsoft Entra ID interprets the application type based on the value of this property.",
			},
			"service_management_reference": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "References application or service contact information from a Service or Asset Management database. Nullable.",
			},
			"tags": schema.SetAttribute{
				Computed:            true,
				ElementType:         types.StringType,
				MarkdownDescription: "Custom strings that can be used to categorize and identify the application. Not nullable. Strings added here also appear in the tags property of any associated service principals. Supports `$filter` (`eq`, `not`, `ge`, `le`, `startsWith`) and `$search`.",
			},
			"disabled_by_microsoft_status": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Specifies whether Microsoft has disabled the registered application. The possible values are: null (default value), `NotDisabled`, and `DisabledDueToViolationOfServicesAgreement` (reasons may include suspicious, abusive, or malicious activity, or a violation of the Microsoft Services Agreement). Supports `$filter` (`eq`, `ne`, `not`). Read-only.",
			},
			"publisher_domain": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The verified publisher domain for the application. Read-only. Supports `$filter` (`eq`, `ne`, `ge`, `le`, `startsWith`).",
			},
			"created_date_time": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The date and time the application was registered. The DateTimeOffset type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z. Read-only. Supports `$filter` (`eq`, `ne`, `not`, `ge`, `le`, `in`, and `eq` on null values) and `$orderby`.",
			},
			"deleted_date_time": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The date and time the application was deleted. The DateTimeOffset type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z. Read-only.",
			},
			"api": schema.SingleNestedAttribute{
				Computed:            true,
				MarkdownDescription: "Specifies settings for an application that implements a web API.",
				Attributes: map[string]schema.Attribute{
					"accept_mapped_claims": schema.BoolAttribute{
						Computed:            true,
						MarkdownDescription: "Allows an application to use claims mapping without specifying a custom signing key.",
					},
					"known_client_applications": schema.SetAttribute{
						Computed:            true,
						ElementType:         types.StringType,
						MarkdownDescription: "Used for bundling consent if you have a solution that contains two parts: a client app and a custom web API app. If you set the appID of the client app to this value, the user only consents once to the client app. Microsoft Entra ID knows that consenting to the client means implicitly consenting to the web API and automatically provisions service principals for both APIs at the same time. Both the client and the web API app must be registered in the same tenant.",
					},
					"oauth2_permission_scopes": schema.SetNestedAttribute{
						Computed:            true,
						MarkdownDescription: "The definition of the delegated permissions exposed by the web API represented by this application registration. These delegated permissions may be requested by a client application, and may be granted by users or administrators during consent. Delegated permissions are sometimes referred to as OAuth 2.0 scopes.",
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"id": schema.StringAttribute{
									Computed:            true,
									MarkdownDescription: "Unique scope permission identifier inside the oauth2PermissionScopes collection. Required.",
								},
								"admin_consent_description": schema.StringAttribute{
									Computed:            true,
									MarkdownDescription: "A description of the delegated permissions, intended to be read by an administrator granting the permission on behalf of all users. This text appears in tenant-wide admin consent experiences.",
								},
								"admin_consent_display_name": schema.StringAttribute{
									Computed:            true,
									MarkdownDescription: "The permission's title, intended to be read by an administrator granting the permission on behalf of all users.",
								},
								"is_enabled": schema.BoolAttribute{
									Computed:            true,
									MarkdownDescription: "When you create or update a permission, this property must be set to true (which is the default). To delete a permission, this property must first be set to false. At that point, in a subsequent call, the permission may be removed.",
								},
								"type": schema.StringAttribute{
									Computed:            true,
									MarkdownDescription: "The possible values are: `User` and `Admin`. Specifies whether this delegated permission should be considered safe for non-admin users to consent to on behalf of themselves, or whether an administrator consent should be required for the permissions. While Microsoft Graph defines the default consent requirement for each permission, the tenant administrator may override the behavior in their organization (by allowing, restricting, or limiting user consent to this delegated permission). For more information, see Configure how users consent to applications.",
								},
								"user_consent_description": schema.StringAttribute{
									Computed:            true,
									MarkdownDescription: "A description of the delegated permissions, intended to be read by a user granting the permission on their own behalf. This text appears in consent experiences where the user is consenting only on behalf of themselves.",
								},
								"user_consent_display_name": schema.StringAttribute{
									Computed:            true,
									MarkdownDescription: "A title for the permission, intended to be read by a user granting the permission on their own behalf. This text appears in consent experiences where the user is consenting only on behalf of themselves.",
								},
								"value": schema.StringAttribute{
									Computed:            true,
									MarkdownDescription: "Specifies the value to include in the scp (scope) claim in access tokens. Must not exceed 120 characters in length. Allowed characters are : ! # $ % & ' ( ) * + , - . / : ;  =  ? @ [ ] ^ + _  {  } ~, and characters in the ranges 0-9, A-Z and a-z. Any other character, including the space character, aren't allowed. May not begin with `.`.",
								},
							},
						},
					},
					"pre_authorized_applications": schema.SetNestedAttribute{
						Computed:            true,
						MarkdownDescription: "Lists the client applications that are preauthorized with the specified delegated permissions to access this application's APIs. Users aren't required to consent to any preauthorized application (for the permissions specified). However, any other permissions not listed in preAuthorizedApplications (requested through incremental consent for example) will require user consent.",
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"app_id": schema.StringAttribute{
									Computed:            true,
									MarkdownDescription: "The unique identifier for the client application.",
								},
								"delegated_permission_ids": schema.SetAttribute{
									Computed:            true,
									ElementType:         types.StringType,
									MarkdownDescription: "The unique identifier for the scopes the client application is granted.",
								},
							},
						},
					},
					"requested_access_token_version": schema.Int32Attribute{
						Computed:            true,
						MarkdownDescription: "Specifies the access token version expected by this resource. This changes the version and format of the JWT produced independent of the endpoint or client used to request the access token. The endpoint used, v1.0 or v2.0, is chosen by the client and only impacts the version of id_tokens. Resources need to explicitly configure requestedAccessTokenVersion to indicate the supported access token format. Possible values for requestedAccessTokenVersion are 1, 2, or null. If the value is null, this defaults to 1, which corresponds to the v1.0 endpoint. If signInAudience on the application is configured as AzureADandPersonalMicrosoftAccount or PersonalMicrosoftAccount, the value for this property must be 2.",
					},
				},
			},
			"app_roles": schema.SetNestedAttribute{
				Computed:            true,
				MarkdownDescription: "The collection of roles defined for the application. With app role assignments, these roles can be assigned to users, groups, or service principals associated with other applications. Not nullable.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "Unique role identifier inside the appRoles collection. When creating a new app role, a new GUID identifier must be provided. Required.",
						},
						"allowed_member_types": schema.SetAttribute{
							Computed:            true,
							ElementType:         types.StringType,
							MarkdownDescription: "Specifies whether this app role can be assigned to users and groups (by setting to `['User']`), to other application's (by setting to `['Application']`, or both (by setting to `['User', 'Application']`). App roles supporting assignment to other applications' service principals are also known as application permissions. The 'Application' value is only supported for app roles defined on application entities. Required.",
						},
						"description": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The description for the app role. This is displayed when the app role is being assigned and, if the app role functions as an application permission, during consent experiences. Required.",
						},
						"display_name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "Display name for the permission that appears in the app role assignment and consent experiences. Required.",
						},
						"is_enabled": schema.BoolAttribute{
							Computed:            true,
							MarkdownDescription: "Defines whether the application's app role is enabled or disabled. Required.",
						},
						"origin": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "Specifies if the app role is defined on the application object or on the servicePrincipal entity. Must not be included in any POST or PATCH requests. Read-only.",
						},
						"value": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "Specifies the value to include in the roles claim in ID tokens and access tokens authenticating an assigned user or service principal. Must not exceed 120 characters in length. Allowed characters are : ! # $ % & ' ( ) * + , - . / : ;  =  ? @ [ ] ^ + _  {  } ~, and characters in the ranges 0-9, A-Z and a-z. Any other character, including the space character, aren't allowed. May not begin with `.`. Nullable.",
						},
					},
				},
			},
			"info": schema.SingleNestedAttribute{
				Computed:            true,
				MarkdownDescription: "Basic profile information of the application, such as it's marketing, support, terms of service, and privacy statement URLs. The terms of service and privacy statement are surfaced to users through the user consent experience. For more information, see How to: Add Terms of service and privacy statement for registered Microsoft Entra apps. Supports `$filter` (`eq`, `ne`, `not`, `ge`, `le`, and `eq` on null values).",
				Attributes: map[string]schema.Attribute{
					"logo_url": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "CDN URL to the application's logo. Read-only.",
					},
					"marketing_url": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "Link to the application's marketing page. For example, https://www.contoso.com/app/marketing.",
					},
					"privacy_statement_url": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "Link to the application's privacy statement. For example, https://www.contoso.com/app/privacy.",
					},
					"support_url": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "Link to the application's support page. For example, https://www.contoso.com/app/support.",
					},
					"terms_of_service_url": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "Link to the application's terms of service statement. For example, https://www.contoso.com/app/termsofservice.",
					},
				},
			},
			"key_credentials": schema.SetNestedAttribute{
				Computed:            true,
				MarkdownDescription: "The collection of key credentials associated with the application. This is a read-only attribute. To manage certificate credentials, use the `microsoft365_graph_beta_applications_application_certificate_credential` resource instead.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"custom_key_identifier": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "A 40-character binary type that can be used to identify the credential. Optional. When not provided in the payload, defaults to the thumbprint of the certificate.",
						},
						"display_name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "Friendly name for the key. Optional.",
						},
						"end_date_time": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The date and time at which the credential expires. The DateTimeOffset type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z.",
						},
						"key": schema.StringAttribute{
							Computed:            true,
							Sensitive:           true,
							MarkdownDescription: "Value for the key credential. Should be a Base64 encoded value. Returned only on $select for a single object, that is, GET applications/{applicationId}?$select=keyCredentials or GET servicePrincipals/{servicePrincipalId}?$select=keyCredentials; otherwise, it's always null. From a .cer certificate, you can read the key using the Convert.ToBase64String() method. For more information, see Get the certificate key.",
						},
						"key_id": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The unique identifier (GUID) for the key.",
						},
						"start_date_time": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The date and time at which the credential becomes valid. The DateTimeOffset type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z.",
						},
						"type": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The type of key credential; for example, `Symmetric`, `AsymmetricX509Cert`, or `X509CertAndPassword`.",
						},
						"usage": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "A string that describes the purpose for which the key can be used; for example, `None​`, `Verify`, `PairwiseIdentifier`, `Sign`.",
						},
					},
				},
			},
			"password_credentials": schema.SetNestedAttribute{
				Computed:            true,
				MarkdownDescription: "The collection of password credentials associated with the application. Not nullable.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"custom_key_identifier": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "Do not use.",
						},
						"display_name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "Friendly name for the password. Optional.",
						},
						"end_date_time": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The date and time at which the password expires represented using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z. Optional.",
						},
						"hint": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "Contains the first three characters of the password. Read-only.",
						},
						"key_id": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The unique identifier for the password. Required.",
						},
						"secret_text": schema.StringAttribute{
							Computed:            true,
							Sensitive:           true,
							MarkdownDescription: "Read-only; Contains the strong passwords generated by Microsoft Entra ID that are 16-64 characters in length. The generated password value is only returned during the initial POST request to addPassword. There is no way to retrieve this password in the future.",
						},
						"start_date_time": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The date and time at which the password becomes valid. The Timestamp type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z. Optional.",
						},
					},
				},
			},
			"optional_claims": schema.SingleNestedAttribute{
				Computed:            true,
				MarkdownDescription: "Application developers can configure optional claims in their Microsoft Entra applications to specify the claims that are sent to their application by the Microsoft security token service. For more information, see How to: Provide optional claims to your app.",
				Attributes: map[string]schema.Attribute{
					"access_token": schema.SetNestedAttribute{
						Computed:            true,
						MarkdownDescription: "The optional claims returned in the JWT access token.",
						NestedObject: schema.NestedAttributeObject{
							Attributes: optionalClaimAttributes(),
						},
					},
					"id_token": schema.SetNestedAttribute{
						Computed:            true,
						MarkdownDescription: "The optional claims returned in the JWT ID token.",
						NestedObject: schema.NestedAttributeObject{
							Attributes: optionalClaimAttributes(),
						},
					},
					"saml2_token": schema.SetNestedAttribute{
						Computed:            true,
						MarkdownDescription: "The optional claims returned in the SAML token.",
						NestedObject: schema.NestedAttributeObject{
							Attributes: optionalClaimAttributes(),
						},
					},
				},
			},
			"parental_control_settings": schema.SingleNestedAttribute{
				Computed:            true,
				MarkdownDescription: "Specifies parental control settings for an application.",
				Attributes: map[string]schema.Attribute{
					"countries_blocked_for_minors": schema.SetAttribute{
						Computed:            true,
						ElementType:         types.StringType,
						MarkdownDescription: "Specifies the two-letter ISO country codes. Access to the application will be blocked for minors from the countries specified in this list.",
					},
					"legal_age_group_rule": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "Specifies the legal age group rule that applies to users of the app. Can be set to one of the following values: `Allow`, `RequireConsentForPrivacyServices`, `RequireConsentForMinors`, `RequireConsentForKids`, `BlockMinors`.",
					},
				},
			},
			"public_client": schema.SingleNestedAttribute{
				Computed:            true,
				MarkdownDescription: "Specifies settings for installed clients such as desktop or mobile devices.",
				Attributes: map[string]schema.Attribute{
					"redirect_uris": schema.SetAttribute{
						Computed:            true,
						ElementType:         types.StringType,
						MarkdownDescription: "Specifies the URLs where user tokens are sent for sign-in, or the redirect URIs where OAuth 2.0 authorization codes and access tokens are sent.",
					},
				},
			},
			"required_resource_access": schema.SetNestedAttribute{
				Computed:            true,
				MarkdownDescription: "Specifies the resources that the application needs to access. This property also specifies the set of delegated permissions and application roles that it needs for each of those resources. This configuration of access to the required resources drives the consent experience. No more than 50 resource services (APIs) can be configured. Beginning mid-October 2021, the total number of required permissions must not exceed 400. For more information, see Limits on requested permissions per app. Not nullable. Supports `$filter` (`eq`, `not`, `ge`, `le`).",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"resource_app_id": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The unique identifier for the resource that the application requires access to. This should be equal to the appId declared on the target resource application. Required.",
						},
						"resource_access": schema.SetNestedAttribute{
							Computed:            true,
							MarkdownDescription: "The list of OAuth2.0 permission scopes and app roles that the application requires from the specified resource. Required.",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"id": schema.StringAttribute{
										Computed:            true,
										MarkdownDescription: "The unique identifier of an app role or delegated permission exposed by the resource application. For delegated permissions, this should match the id property of one of the delegated permissions in the oauth2PermissionScopes collection of the resource application's service principal. For app roles (application permissions), this should match the id property of an app role in the appRoles collection of the resource application's service principal. Required.",
									},
									"type": schema.StringAttribute{
										Computed:            true,
										MarkdownDescription: "Specifies whether the id property references a delegated permission or an app role (application permission). The possible values are: `Scope` (for delegated permissions) or `Role` (for app roles). Required.",
									},
								},
							},
						},
					},
				},
			},
			"spa": schema.SingleNestedAttribute{
				Computed:            true,
				MarkdownDescription: "Specifies settings for a single-page application, including sign out URLs and redirect URIs for authorization codes and access tokens.",
				Attributes: map[string]schema.Attribute{
					"redirect_uris": schema.SetAttribute{
						Computed:            true,
						ElementType:         types.StringType,
						MarkdownDescription: "Specifies the URLs where user tokens are sent for sign-in, or the redirect URIs where OAuth 2.0 authorization codes and access tokens are sent.",
					},
				},
			},
			"web": schema.SingleNestedAttribute{
				Computed:            true,
				MarkdownDescription: "Specifies settings for a web application.",
				Attributes: map[string]schema.Attribute{
					"home_page_url": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "Home page or landing page of the application.",
					},
					"logout_url": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "Specifies the URL that is used by Microsoft's authorization service to log out a user using front-channel, back-channel or SAML logout protocols.",
					},
					"redirect_uris": schema.SetAttribute{
						Computed:            true,
						ElementType:         types.StringType,
						MarkdownDescription: "Specifies the URLs where user tokens are sent for sign-in, or the redirect URIs where OAuth 2.0 authorization codes and access tokens are sent.",
					},
					"implicit_grant_settings": schema.SingleNestedAttribute{
						Computed:            true,
						MarkdownDescription: "Specifies whether this web application can request tokens using the OAuth 2.0 implicit flow.",
						Attributes: map[string]schema.Attribute{
							"enable_access_token_issuance": schema.BoolAttribute{
								Computed:            true,
								MarkdownDescription: "Specifies whether this web application can request an access token using the OAuth 2.0 implicit flow.",
							},
							"enable_id_token_issuance": schema.BoolAttribute{
								Computed:            true,
								MarkdownDescription: "Specifies whether this web application can request an ID token using the OAuth 2.0 implicit flow.",
							},
						},
					},
					"redirect_uri_settings": schema.SetNestedAttribute{
						Computed:            true,
						MarkdownDescription: "Specifies the index of the URLs where user tokens are sent for sign-in. This is only valid for applications using SAML. Note: If not specified, the API may auto-generate settings based on redirect_uris. To manage this field, you must provide at least one entry; empty arrays are not supported as the API auto-generates values.",
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"uri": schema.StringAttribute{
									Computed:            true,
									MarkdownDescription: "Specifies the URI that tokens are sent to.",
								},
								"index": schema.Int32Attribute{
									Computed:            true,
									MarkdownDescription: "Identifies the specific URI within the redirectURIs collection in SAML SSO flows. Defaults to null. The index is unique across all the redirectUris for the application.",
								},
							},
						},
					},
				},
			},
			"sign_in_audience_restrictions": schema.SingleNestedAttribute{
				Computed:            true,
				MarkdownDescription: "Specifies restrictions on the supported account types specified in signInAudience. The value type determines the restrictions that can be applied: unrestrictedAudience (There are no additional restrictions on the supported account types allowed by signInAudience) or allowedTenantsAudience (The application can only be used in the specified Entra tenants. Only supported when signInAudience is AzureADMultipleOrgs). Default is a value of type unrestrictedAudience. Returned only on `$select`.",
				Attributes: map[string]schema.Attribute{
					"odata_type": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "The OData type. Must be `#microsoft.graph.allowedTenantsAudience` or `#microsoft.graph.unrestrictedAudience`.",
					},
					"is_home_tenant_allowed": schema.BoolAttribute{
						Computed:            true,
						MarkdownDescription: "Indicates whether the home tenant is allowed. Only applicable when odata_type is `#microsoft.graph.allowedTenantsAudience`.",
					},
					"allowed_tenant_ids": schema.SetAttribute{
						Computed:            true,
						ElementType:         types.StringType,
						MarkdownDescription: "The list of allowed tenant IDs. Only applicable when odata_type is `#microsoft.graph.allowedTenantsAudience`.",
					},
				},
			},
			"owner_user_ids": schema.SetAttribute{
				Computed:            true,
				ElementType:         types.StringType,
				MarkdownDescription: "The user IDs of the owners for the application. At least one owner is typically required when creating an application. Owners are a set of non-admin users or service principals allowed to modify this object.",
			},
			"timeouts": commonschema.DatasourceTimeouts(ctx),
		},
	}
}

// optionalClaimAttributes returns the common attributes for optional claims
func optionalClaimAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"name": schema.StringAttribute{
			Computed:            true,
			MarkdownDescription: "The name of the optional claim. Required.",
		},
		"source": schema.StringAttribute{
			Computed:            true,
			MarkdownDescription: "The source (directory object) of the claim. There are predefined claims and user-defined claims from extension properties. If the source value is null, the claim is a predefined optional claim. If the source value is user, the value in the name property is the extension property from the user object.",
		},
		"essential": schema.BoolAttribute{
			Computed:            true,
			MarkdownDescription: "If the value is true, the claim specified by the client is necessary to ensure a smooth authorization experience for the specific task requested by the end user. The default value is false.",
		},
		"additional_properties": schema.SetAttribute{
			Computed:            true,
			ElementType:         types.StringType,
			MarkdownDescription: "Additional properties of the claim. If a property exists in this collection, it modifies the behavior of the optional claim specified in the name property.",
		},
	}
}
