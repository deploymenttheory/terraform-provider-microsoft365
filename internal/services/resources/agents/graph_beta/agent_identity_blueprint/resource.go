package graphBetaAgentsAgentIdentityBlueprint

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

var _ resource.Resource = &AgentIdentityBlueprintResource{}
var _ resource.ResourceWithConfigure = &AgentIdentityBlueprintResource{}
var _ resource.ResourceWithImportState = &AgentIdentityBlueprintResource{}

const (
	ResourceName = "microsoft365_graph_beta_agents_agent_identity_blueprint"
)

func NewAgentIdentityBlueprintResource() resource.Resource {
	return &AgentIdentityBlueprintResource{
		ReadPermissions: []string{
			"AgentIdentityBlueprint.Read.All",
			"Application.Read.All",
			"Directory.Read.All",
		},
		WritePermissions: []string{
			"AgentIdentityBlueprint.AddRemoveCreds.All",
			"AgentIdentityBlueprint.UpdateBranding.All",
			"AgentIdentityBlueprint.ReadWrite.All",
			"Directory.ReadWrite.All",
		},
		ResourcePath: "/applications",
	}
}

type AgentIdentityBlueprintResource struct {
	httpClient       *client.AuthenticatedHTTPClient
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

// Metadata returns the resource type name.
func (r *AgentIdentityBlueprintResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = ResourceName
}

// Configure sets the client for the resource.
func (r *AgentIdentityBlueprintResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.httpClient = client.SetGraphBetaHTTPClientForResource(ctx, req, resp, ResourceName)
}

// ImportState imports the resource state.
func (r *AgentIdentityBlueprintResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Schema defines the schema for the resource.
func (r *AgentIdentityBlueprintResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages Microsoft Entra ID Agent Identity Blueprints using the `/applications` endpoint with OData type casting. " +
			"An agent identity blueprint serves as a template for creating agent identities within the Microsoft Entra ID ecosystem.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Unique identifier for the agent identity blueprint object. This property is referred to as **Object ID** in the Microsoft Entra admin center. " +
					"Key. Not nullable. Read-only. Inherited from [directoryObject](https://learn.microsoft.com/en-us/graph/api/resources/directoryobject?view=graph-rest-beta).",
				Computed: true,
			},
			"app_id": schema.StringAttribute{
				MarkdownDescription: "The unique identifier for the agent identity blueprint assigned by Microsoft Entra ID. " +
					"Also known as **Application (client) ID**. Read-only. Inherited from [application](https://learn.microsoft.com/en-us/graph/api/resources/application?view=graph-rest-beta).",
				Computed: true,
			},
			"display_name": schema.StringAttribute{
				MarkdownDescription: "The display name for the agent identity blueprint. Maximum length is 256 characters. " +
					"Inherited from [application](https://learn.microsoft.com/en-us/graph/api/resources/application?view=graph-rest-beta).",
				Required: true,
			},
			"unique_name": schema.StringAttribute{
				MarkdownDescription: "The unique identifier that can be assigned to an agent identity blueprint and used as an alternate key. Immutable. Read-only. " +
					"Inherited from [application](https://learn.microsoft.com/en-us/graph/api/resources/application?view=graph-rest-beta).",
				Computed: true,
			},
			"created_by_app_id": schema.StringAttribute{
				MarkdownDescription: "The **appId** (called **Application (client) ID** on the Microsoft Entra admin center) of the application that created this agent identity blueprint. " +
					"Set internally by Microsoft Entra ID. Read-only. Inherited from [application](https://learn.microsoft.com/en-us/graph/api/resources/application?view=graph-rest-beta).",
				Computed: true,
			},
			"created_date_time": schema.StringAttribute{
				MarkdownDescription: "The date and time the agent identity blueprint was registered. The DateTimeOffset type represents date and time information using ISO 8601 format and is always in UTC time. " +
					"Read-only. Inherited from [application](https://learn.microsoft.com/en-us/graph/api/resources/application?view=graph-rest-beta).",
				Computed: true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Free text field to provide a description of the agent identity blueprint to end users. The maximum allowed size is 1,024 characters. " +
					"Inherited from [application](https://learn.microsoft.com/en-us/graph/api/resources/application?view=graph-rest-beta).",
				Optional: true,
			},
			"publisher_domain": schema.StringAttribute{
				MarkdownDescription: "The verified publisher domain for the agent identity blueprint. Read-only. " +
					"Inherited from [application](https://learn.microsoft.com/en-us/graph/api/resources/application?view=graph-rest-beta).",
				Computed: true,
			},
			"sign_in_audience": schema.StringAttribute{
				MarkdownDescription: "Specifies the Microsoft accounts that are supported for the current agent identity blueprint. " +
					"The possible values are: `AzureADMyOrg` (default), `AzureADMultipleOrgs`, `AzureADandPersonalMicrosoftAccount`, and `PersonalMicrosoftAccount`. " +
					"Inherited from [application](https://learn.microsoft.com/en-us/graph/api/resources/application?view=graph-rest-beta).",
				Optional: true,
				Computed: true,
			},
			"group_membership_claims": schema.StringAttribute{
				MarkdownDescription: "Configures the groups claim issued in a user or OAuth 2.0 access token that the agent identity blueprint expects. " +
					"To set this attribute, use one of the following string values: `None`, `SecurityGroup` (for security groups and Microsoft Entra roles), " +
					"`All` (this gets all security groups, distribution groups, and Microsoft Entra directory roles that the signed-in user is a member of). " +
					"Inherited from [application](https://learn.microsoft.com/en-us/graph/api/resources/application?view=graph-rest-beta).",
				Optional: true,
			},
			"disabled_by_microsoft_status": schema.StringAttribute{
				MarkdownDescription: "Specifies whether Microsoft has disabled the registered agent identity blueprint. " +
					"Possible values are: null (default value), `NotDisabled`, and `DisabledDueToViolationOfServicesAgreement` (reasons may include suspicious, abusive, or malicious activity, or a violation of the Microsoft Services Agreement). " +
					"Inherited from [application](https://learn.microsoft.com/en-us/graph/api/resources/application?view=graph-rest-beta).",
				Computed: true,
			},
			"service_management_reference": schema.StringAttribute{
				MarkdownDescription: "References application or service contact information from a Service or Asset Management database. Nullable. " +
					"Inherited from [application](https://learn.microsoft.com/en-us/graph/api/resources/application?view=graph-rest-beta).",
				Optional: true,
			},
			"token_encryption_key_id": schema.StringAttribute{
				MarkdownDescription: "Specifies the keyId of a public key from the keyCredentials collection. " +
					"When configured, Microsoft Entra ID encrypts all the tokens it emits by using the key this property points to. " +
					"Inherited from [application](https://learn.microsoft.com/en-us/graph/api/resources/application?view=graph-rest-beta).",
				Optional: true,
			},
			"identifier_uris": schema.SetAttribute{
				MarkdownDescription: "Also known as App ID URI, this value is set when an agent identity blueprint is used as a resource app. " +
					"The identifierUris acts as the prefix for the scopes you reference in your API's code, and it must be globally unique across Microsoft Entra ID. " +
					"Not nullable. Inherited from [application](https://learn.microsoft.com/en-us/graph/api/resources/application?view=graph-rest-beta).",
				ElementType: schema.StringAttribute{}.GetType(),
				Optional:    true,
			},
			"tags": schema.SetAttribute{
				MarkdownDescription: "Custom strings that can be used to categorize and identify the agent identity blueprint. Not nullable. " +
					"Inherited from [application](https://learn.microsoft.com/en-us/graph/api/resources/application?view=graph-rest-beta).",
				ElementType: schema.StringAttribute{}.GetType(),
				Optional:    true,
			},
			"info": schema.SingleNestedAttribute{
				MarkdownDescription: "Basic profile information of the agent identity blueprint, such as its marketing, support, terms of service, and privacy statement URLs. " +
					"The terms of service and privacy statement are surfaced to users through the user consent experience. " +
					"Inherited from [application](https://learn.microsoft.com/en-us/graph/api/resources/application?view=graph-rest-beta). " +
					"See [informationalUrl](https://learn.microsoft.com/en-us/graph/api/resources/informationalurl?view=graph-rest-beta).",
				Optional: true,
				Attributes: map[string]schema.Attribute{
					"logo_url": schema.StringAttribute{
						MarkdownDescription: "CDN URL to the application's logo, Read-only.",
						Computed:            true,
					},
					"marketing_url": schema.StringAttribute{
						MarkdownDescription: "Link to the application's marketing page. For example, `https://www.contoso.com/app/marketing`.",
						Optional:            true,
					},
					"privacy_statement_url": schema.StringAttribute{
						MarkdownDescription: "Link to the application's privacy statement. For example, `https://www.contoso.com/app/privacy`.",
						Optional:            true,
					},
					"support_url": schema.StringAttribute{
						MarkdownDescription: "Link to the application's support page. For example, `https://www.contoso.com/app/support`.",
						Optional:            true,
					},
					"terms_of_service_url": schema.StringAttribute{
						MarkdownDescription: "Link to the application's terms of service statement. For example, `https://www.contoso.com/app/termsofservice`.",
						Optional:            true,
					},
				},
			},
			"key_credentials": schema.ListNestedAttribute{
				MarkdownDescription: "The collection of key credentials associated with the agent identity blueprint. Not nullable. " +
					"Inherited from [application](https://learn.microsoft.com/en-us/graph/api/resources/application?view=graph-rest-beta). " +
					"See [keyCredential](https://learn.microsoft.com/en-us/graph/api/resources/keycredential?view=graph-rest-beta).",
				Optional: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"custom_key_identifier": schema.StringAttribute{
							MarkdownDescription: "A 40-character binary type that can be used to identify the credential. Optional. When not provided in the payload, defaults to the thumbprint of the certificate.",
							Optional:            true,
						},
						"display_name": schema.StringAttribute{
							MarkdownDescription: "The friendly name for the key, with a maximum length of 90 characters. Longer values are accepted but shortened. Optional.",
							Optional:            true,
						},
						"end_date_time": schema.StringAttribute{
							MarkdownDescription: "The date and time at which the credential expires. The DateTimeOffset type represents date and time information using ISO 8601 format and is always in UTC time. " +
								"For example, midnight UTC on Jan 1, 2014 is `2014-01-01T00:00:00Z`.",
							Optional: true,
						},
						"key": schema.StringAttribute{
							MarkdownDescription: "Value for the key credential. Should be a Base64 encoded value. Returned only on `$select` for a single object, that is, " +
								"`GET applications/{applicationId}?$select=keyCredentials` or `GET servicePrincipals/{servicePrincipalId}?$select=keyCredentials`; otherwise, it's always null. " +
								"From a `.cer` certificate, you can read the key using the `Convert.ToBase64String()` method.",
							Optional: true,
						},
						"key_id": schema.StringAttribute{
							MarkdownDescription: "The unique identifier for the key.",
							Optional:            true,
						},
						"start_date_time": schema.StringAttribute{
							MarkdownDescription: "The date and time at which the credential becomes valid. The Timestamp type represents date and time information using ISO 8601 format and is always in UTC time. " +
								"For example, midnight UTC on Jan 1, 2014 is `2014-01-01T00:00:00Z`.",
							Optional: true,
						},
						"type": schema.StringAttribute{
							MarkdownDescription: "The type of key credential; for example, `Symmetric`, `AsymmetricX509Cert`, or `X509CertAndPassword`.",
							Optional:            true,
						},
						"usage": schema.StringAttribute{
							MarkdownDescription: "A string that describes the purpose for which the key can be used; for example, `None`, `Verify`, `PairwiseIdentifier`, `Delegation`, `Decrypt`, `Encrypt`, `HashedIdentifier`, `SelfSignedTls`, or `Sign`. " +
								"If `usage` is `Sign`, the `type` should be `X509CertAndPassword`, and the `passwordCredentials` for signing should be defined.",
							Optional: true,
						},
					},
				},
			},
			"password_credentials": schema.ListNestedAttribute{
				MarkdownDescription: "The collection of password credentials associated with the agent identity blueprint. Not nullable. " +
					"Inherited from [application](https://learn.microsoft.com/en-us/graph/api/resources/application?view=graph-rest-beta). " +
					"You can also add passwords after creating the agent identity blueprint by calling the [Add password](https://learn.microsoft.com/en-us/graph/api/agentidentityblueprint-addpassword?view=graph-rest-beta) API. " +
					"See [passwordCredential](https://learn.microsoft.com/en-us/graph/api/resources/passwordcredential?view=graph-rest-beta).",
				Optional: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"custom_key_identifier": schema.StringAttribute{
							MarkdownDescription: "Do not use.",
							Optional:            true,
						},
						"display_name": schema.StringAttribute{
							MarkdownDescription: "Friendly name for the password. Optional.",
							Optional:            true,
						},
						"end_date_time": schema.StringAttribute{
							MarkdownDescription: "The date and time at which the password expires represented using ISO 8601 format and is always in UTC time. " +
								"For example, midnight UTC on Jan 1, 2014 is `2014-01-01T00:00:00Z`. Optional.",
							Optional: true,
						},
						"hint": schema.StringAttribute{
							MarkdownDescription: "Contains the first three characters of the password. Read-only.",
							Computed:            true,
						},
						"key_id": schema.StringAttribute{
							MarkdownDescription: "The unique identifier for the password.",
							Optional:            true,
						},
						"secret_text": schema.StringAttribute{
							MarkdownDescription: "Read-only; Contains the strong passwords generated by Microsoft Entra ID that are 16-64 characters in length. " +
								"The generated password value is only returned during the initial POST request to [addPassword](https://learn.microsoft.com/en-us/graph/api/agentidentityblueprint-addpassword?view=graph-rest-beta). " +
								"There is no way to retrieve this password in the future.",
							Computed:  true,
							Sensitive: true,
						},
						"start_date_time": schema.StringAttribute{
							MarkdownDescription: "The date and time at which the password becomes valid. The Timestamp type represents date and time information using ISO 8601 format and is always in UTC time. " +
								"For example, midnight UTC on Jan 1, 2014 is `2014-01-01T00:00:00Z`. Optional.",
							Optional: true,
						},
					},
				},
			},
			"app_roles": schema.ListNestedAttribute{
				MarkdownDescription: "The collection of roles assigned to the agent identity blueprint. " +
					"With app role assignments, these roles can be assigned to users, groups, or service principals associated with other applications. Not nullable. " +
					"Inherited from [application](https://learn.microsoft.com/en-us/graph/api/resources/application?view=graph-rest-beta). " +
					"See [appRole](https://learn.microsoft.com/en-us/graph/api/resources/approle?view=graph-rest-beta).",
				Optional: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"allowed_member_types": schema.SetAttribute{
							MarkdownDescription: "Specifies whether this app role can be assigned to users and groups (by setting to `[\"User\"]`), to other application's (by setting to `[\"Application\"]`, or both (by setting to `[\"User\", \"Application\"]`). " +
								"App roles supporting assignment to other applications' service principals are also known as [application permissions](https://learn.microsoft.com/en-us/graph/auth/auth-concepts#microsoft-graph-permissions). " +
								"The \"Application\" value is only supported for app roles defined on **application** entities.",
							ElementType: schema.StringAttribute{}.GetType(),
							Optional:    true,
						},
						"description": schema.StringAttribute{
							MarkdownDescription: "The description for the app role. This is displayed when the app role is being assigned and, if the app role functions as an application permission, during consent experiences.",
							Optional:            true,
						},
						"display_name": schema.StringAttribute{
							MarkdownDescription: "Display name for the permission that appears in the app role assignment and consent experiences.",
							Optional:            true,
						},
						"id": schema.StringAttribute{
							MarkdownDescription: "Unique role identifier inside the **appRoles** collection. You must specify a new GUID identifier when you create a new app role.",
							Optional:            true,
							Computed:            true,
						},
						"is_enabled": schema.BoolAttribute{
							MarkdownDescription: "When you create or updating an app role, this value must be **true**. To delete a role, this must first be set to **false**. " +
								"At that point, in a subsequent call, this role might be removed. Default value is **true**.",
							Optional: true,
						},
						"origin": schema.StringAttribute{
							MarkdownDescription: "Specifies if the app role is defined on the [application](https://learn.microsoft.com/en-us/graph/api/resources/application?view=graph-rest-beta) object or on the " +
								"[servicePrincipal](https://learn.microsoft.com/en-us/graph/api/resources/serviceprincipal?view=graph-rest-beta) entity. Must _not_ be included in any POST or PATCH requests. Read-only.",
							Computed: true,
						},
						"value": schema.StringAttribute{
							MarkdownDescription: "Specifies the value to include in the roles claim in ID tokens and access tokens authenticating an assigned user or service principal. " +
								"Must not exceed 120 characters in length. Allowed characters are `: ! # $ % & ' ( ) * + , - . / : ; < = > ? @ [ ] ^ + _ \\` { | } ~`, and characters in the ranges `0-9`, `A-Z`, and `a-z`. " +
								"Any other character, including the space character, aren't allowed. May not begin with `.`.",
							Optional: true,
						},
					},
				},
			},
			"api": schema.SingleNestedAttribute{
				MarkdownDescription: "Specifies settings for an application that implements a web API. " +
					"Inherited from [application](https://learn.microsoft.com/en-us/graph/api/resources/application?view=graph-rest-beta). " +
					"See [apiApplication](https://learn.microsoft.com/en-us/graph/api/resources/apiapplication?view=graph-rest-beta).",
				Optional: true,
				Attributes: map[string]schema.Attribute{
					"accept_mapped_claims": schema.BoolAttribute{
						MarkdownDescription: "When `true`, allows an application to use claims mapping without specifying a custom signing key.",
						Optional:            true,
					},
					"known_client_applications": schema.SetAttribute{
						MarkdownDescription: "Used for bundling consent if you have a solution that contains two parts: a client app and a custom web API app. " +
							"If you set the appID of the client app to this value, the user only consents once to the client app. " +
							"Microsoft Entra ID knows that consenting to the client means implicitly consenting to the web API and automatically provisions service principals for both APIs at the same time.",
						ElementType: schema.StringAttribute{}.GetType(),
						Optional:    true,
					},
					"pre_authorized_applications": schema.ListNestedAttribute{
						MarkdownDescription: "Lists the client applications that are pre-authorized with the specified delegated permissions to access this application's APIs.",
						Optional:            true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"app_id": schema.StringAttribute{
									MarkdownDescription: "The unique identifier for the client application.",
									Optional:            true,
								},
								"delegated_permission_ids": schema.SetAttribute{
									MarkdownDescription: "The unique identifier for the delegated permissions the client application is authorized to use.",
									ElementType:         schema.StringAttribute{}.GetType(),
									Optional:            true,
								},
							},
						},
					},
					"requested_access_token_version": schema.Int64Attribute{
						MarkdownDescription: "Specifies the access token version expected by this resource. Possible values are 1, 2, or null.",
						Optional:            true,
					},
					"oauth2_permission_scopes": schema.ListNestedAttribute{
						MarkdownDescription: "The definition of the delegated permissions exposed by the web API represented by this application registration.",
						Optional:            true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"admin_consent_description": schema.StringAttribute{
									MarkdownDescription: "A description of the delegated permissions, intended to be read by an administrator granting the permission on behalf of all users.",
									Optional:            true,
								},
								"admin_consent_display_name": schema.StringAttribute{
									MarkdownDescription: "The permission's title, intended to be read by an administrator granting the permission on behalf of all users.",
									Optional:            true,
								},
								"id": schema.StringAttribute{
									MarkdownDescription: "Unique delegated permission identifier inside the collection of delegated permissions defined for a resource application.",
									Optional:            true,
									Computed:            true,
								},
								"is_enabled": schema.BoolAttribute{
									MarkdownDescription: "When you create or update a permission, this property must be set to true (which is the default). To delete a permission, this property must first be set to false.",
									Optional:            true,
								},
								"origin": schema.StringAttribute{
									MarkdownDescription: "For internal use only. Don't write or rely on this property. May be removed in future versions.",
									Computed:            true,
								},
								"type": schema.StringAttribute{
									MarkdownDescription: "The possible values are: User and Admin. Specifies whether this delegated permission should be considered safe for non-admin users to consent to on behalf of themselves, or whether an administrator consent should be required for the permission to be granted.",
									Optional:            true,
								},
								"user_consent_description": schema.StringAttribute{
									MarkdownDescription: "A description of the delegated permissions, intended to be read by a user granting the permission on their own behalf.",
									Optional:            true,
								},
								"user_consent_display_name": schema.StringAttribute{
									MarkdownDescription: "A title for the permission, intended to be read by a user granting the permission on their own behalf.",
									Optional:            true,
								},
								"value": schema.StringAttribute{
									MarkdownDescription: "Specifies the value to include in the scp (scope) claim in access tokens.",
									Optional:            true,
								},
							},
						},
					},
				},
			},
			"optional_claims": schema.SingleNestedAttribute{
				MarkdownDescription: "Application developers can configure optional claims in their Microsoft Entra agent identity blueprints to specify the claims that are sent to their application by the Microsoft security token service. " +
					"For more information, see [How to: Provide optional claims to your app](https://learn.microsoft.com/en-us/entra/identity-platform/optional-claims). " +
					"Inherited from [application](https://learn.microsoft.com/en-us/graph/api/resources/application?view=graph-rest-beta). " +
					"See [optionalClaims](https://learn.microsoft.com/en-us/graph/api/resources/optionalclaims?view=graph-rest-beta).",
				Optional: true,
				Attributes: map[string]schema.Attribute{
					"access_token": schema.SetAttribute{
						MarkdownDescription: "The optional claims returned in the JWT access token. " +
							"See [optionalClaim](https://learn.microsoft.com/en-us/graph/api/resources/optionalclaim?view=graph-rest-beta).",
						ElementType: schema.StringAttribute{}.GetType(),
						Optional:    true,
					},
					"id_token": schema.SetAttribute{
						MarkdownDescription: "The optional claims returned in the JWT ID token. " +
							"See [optionalClaim](https://learn.microsoft.com/en-us/graph/api/resources/optionalclaim?view=graph-rest-beta).",
						ElementType: schema.StringAttribute{}.GetType(),
						Optional:    true,
					},
					"saml2_token": schema.SetAttribute{
						MarkdownDescription: "The optional claims returned in the SAML token. " +
							"See [optionalClaim](https://learn.microsoft.com/en-us/graph/api/resources/optionalclaim?view=graph-rest-beta).",
						ElementType: schema.StringAttribute{}.GetType(),
						Optional:    true,
					},
				},
			},
			"verified_publisher": schema.SingleNestedAttribute{
				MarkdownDescription: "Specifies the verified publisher of the agent identity blueprint. " +
					"For more information, see [Publisher verification](https://learn.microsoft.com/en-us/entra/identity-platform/publisher-verification-overview). " +
					"Inherited from [application](https://learn.microsoft.com/en-us/graph/api/resources/application?view=graph-rest-beta). " +
					"See [verifiedPublisher](https://learn.microsoft.com/en-us/graph/api/resources/verifiedpublisher?view=graph-rest-beta).",
				Optional: true,
				Attributes: map[string]schema.Attribute{
					"added_date_time": schema.StringAttribute{
						MarkdownDescription: "The timestamp when the verified publisher was first added or most recently updated. " +
							"The Timestamp type represents date and time information using ISO 8601 format and is always in UTC time.",
						Computed: true,
					},
					"display_name": schema.StringAttribute{
						MarkdownDescription: "The verified publisher name from the app publisher's Microsoft Partner Network (MPN) account.",
						Computed:            true,
					},
					"verified_publisher_id": schema.StringAttribute{
						MarkdownDescription: "The ID of the verified publisher from the app publisher's Partner Center account.",
						Computed:            true,
					},
				},
			},
			"web": schema.SingleNestedAttribute{
				MarkdownDescription: "Specifies settings for a web application. " +
					"Inherited from [application](https://learn.microsoft.com/en-us/graph/api/resources/application?view=graph-rest-beta). " +
					"See [webApplication](https://learn.microsoft.com/en-us/graph/api/resources/webapplication?view=graph-rest-beta).",
				Optional: true,
				Attributes: map[string]schema.Attribute{
					"home_page_url": schema.StringAttribute{
						MarkdownDescription: "Home page or landing page of the application.",
						Optional:            true,
					},
					"implicit_grant_settings": schema.SingleNestedAttribute{
						MarkdownDescription: "Specifies whether this web application can request tokens using the OAuth 2.0 implicit flow. " +
							"See [implicitGrantSettings](https://learn.microsoft.com/en-us/graph/api/resources/implicitgrantsettings?view=graph-rest-beta).",
						Optional: true,
						Attributes: map[string]schema.Attribute{
							"enable_access_token_issuance": schema.BoolAttribute{
								MarkdownDescription: "Specifies whether this web application can request an access token using the OAuth 2.0 implicit flow.",
								Optional:            true,
							},
							"enable_id_token_issuance": schema.BoolAttribute{
								MarkdownDescription: "Specifies whether this web application can request an ID token using the OAuth 2.0 implicit flow.",
								Optional:            true,
							},
						},
					},
					"logout_url": schema.StringAttribute{
						MarkdownDescription: "Specifies the URL that is used by Microsoft's authorization service to sign out a user using [front-channel](https://openid.net/specs/openid-connect-frontchannel-1_0.html), " +
							"[back-channel](https://openid.net/specs/openid-connect-backchannel-1_0.html) or SAML logout protocols.",
						Optional: true,
					},
					"redirect_uris": schema.SetAttribute{
						MarkdownDescription: "Specifies the URLs where user tokens are sent for sign-in, or the redirect URIs where OAuth 2.0 authorization codes and access tokens are sent.",
						ElementType:         schema.StringAttribute{}.GetType(),
						Optional:            true,
					},
					"redirect_uri_settings": schema.ListNestedAttribute{
						MarkdownDescription: "Specifies the index of the URLs where user tokens are sent for sign-in. This is only valid for applications using SAML.",
						Optional:            true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"index": schema.Int64Attribute{
									MarkdownDescription: "Identifies the specific URI to use from the list of redirect_uris.",
									Optional:            true,
								},
								"uri": schema.StringAttribute{
									MarkdownDescription: "Specifies the URI that tokens are sent to.",
									Optional:            true,
								},
							},
						},
					},
				},
			},
			"certification": schema.SingleNestedAttribute{
				MarkdownDescription: "Specifies the certification status of the agent identity blueprint. " +
					"Supports `$filter` (`eq`, `ne`, `not`). Read-only. " +
					"Inherited from [application](https://learn.microsoft.com/en-us/graph/api/resources/application?view=graph-rest-beta). " +
					"See [certification](https://learn.microsoft.com/en-us/graph/api/resources/certification?view=graph-rest-beta).",
				Optional: true,
				Attributes: map[string]schema.Attribute{
					"certification_details_url": schema.StringAttribute{
						MarkdownDescription: "URL that shows certification details for the application.",
						Computed:            true,
					},
					"certification_expiration_date_time": schema.StringAttribute{
						MarkdownDescription: "The timestamp when the current certification for the application expires. " +
							"The DateTimeOffset type represents date and time information using ISO 8601 format and is always in UTC time.",
						Computed: true,
					},
					"is_certified_by_microsoft": schema.BoolAttribute{
						MarkdownDescription: "Indicates whether the application is certified by Microsoft.",
						Computed:            true,
					},
					"is_publisher_attested": schema.BoolAttribute{
						MarkdownDescription: "Indicates whether the application has been self-attested by the application developer or the publisher.",
						Computed:            true,
					},
					"last_certification_date_time": schema.StringAttribute{
						MarkdownDescription: "The timestamp when the certification for the application was most recently added or updated. " +
							"The DateTimeOffset type represents date and time information using ISO 8601 format and is always in UTC time.",
						Computed: true,
					},
				},
			},
			"timeouts": commonschema.Timeouts(ctx),
		},
	}
}
