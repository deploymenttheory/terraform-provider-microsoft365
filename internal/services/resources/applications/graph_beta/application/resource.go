package graphBetaApplication

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	validate "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/validate/attribute"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	ResourceName  = "microsoft365_graph_beta_applications_application"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	// Basic resource interface (CRUD operations)
	_ resource.Resource = &ApplicationResource{}

	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &ApplicationResource{}

	// Enables import functionality
	_ resource.ResourceWithImportState = &ApplicationResource{}
)

func NewApplicationResource() resource.Resource {
	return &ApplicationResource{
		ReadPermissions: []string{
			"Application.Read.All",
			"Directory.Read.All",
		},
		WritePermissions: []string{
			"Application.ReadWrite.All",
			"Directory.ReadWrite.All",
		},
		ResourcePath: "/applications",
	}
}

type ApplicationResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

// Metadata returns the resource type name.
func (r *ApplicationResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = ResourceName
}

// Configure sets the client for the resource.
func (r *ApplicationResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, ResourceName)
}

// ImportState handles importing the resource with an extended ID format.
//
// Supported formats:
//   - Simple:   "resource_id" (prevent_duplicate_names and hard_delete default to false)
//   - Extended: "resource_id:prevent_duplicate_names=true:hard_delete=true"
//
// Example:
//
//	terraform import microsoft365_graph_beta_applications_application.example "12345678-1234-1234-1234-123456789012:prevent_duplicate_names=true:hard_delete=true"
func (r *ApplicationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	idParts := strings.Split(req.ID, ":")
	resourceID := idParts[0]
	preventDuplicateNames := false // Default
	hardDelete := false            // Default

	if len(idParts) > 1 {
		for _, part := range idParts[1:] {
			if strings.HasPrefix(part, "prevent_duplicate_names=") {
				value := strings.TrimPrefix(part, "prevent_duplicate_names=")
				switch strings.ToLower(value) {
				case "true":
					preventDuplicateNames = true
				case "false":
					preventDuplicateNames = false
				default:
					resp.Diagnostics.AddError(
						"Invalid Import ID",
						fmt.Sprintf("Invalid prevent_duplicate_names value '%s'. Must be 'true' or 'false'.", value),
					)
					return
				}
			}
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

	tflog.Info(ctx, fmt.Sprintf("Importing %s with ID: %s, prevent_duplicate_names: %t, hard_delete: %t", ResourceName, resourceID, preventDuplicateNames, hardDelete))

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), resourceID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("prevent_duplicate_names"), preventDuplicateNames)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("hard_delete"), hardDelete)...)
}

// Schema returns the schema for the resource.
func (r *ApplicationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages an application in Microsoft Entra ID using the `/applications` endpoint. Any application " +
			"that outsources authentication to Microsoft Entra ID must be registered in the Microsoft identity platform. Application " +
			"registration involves telling Microsoft Entra ID about your application, including the URL where it's located, the URL to " +
			"send replies after authentication, the URI to identify your application, and more.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Unique identifier for the application object. This property is referred to as **Object ID** in the Microsoft Entra admin center. Key. Not nullable. Read-only. Supports `$filter` (`eq`, `ne`, `not`, `in`).",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"app_id": schema.StringAttribute{
				MarkdownDescription: "The unique identifier for the application that is assigned by Microsoft Entra ID. Not nullable. Read-only. Alternate key. Supports `$filter` (`eq`).",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"display_name": schema.StringAttribute{
				MarkdownDescription: "The display name for the application. Maximum length is 256 characters. Supports `$filter` (`eq`, `ne`, `not`, `ge`, `le`, `in`, `startsWith`, and `eq` on null values), `$search`, and `$orderby`.",
				Required:            true,
				Validators: []validator.String{
					validate.StringLengthAtMost(256),
				},
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Free text field to provide a description of the application object to end users. The maximum allowed size is 1,024 characters. Returned by default. Supports `$filter` (`eq`, `ne`, `not`, `ge`, `le`, `startsWith`) and `$search`.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					validate.StringLengthAtMost(1024),
				},
			},
			"sign_in_audience": schema.StringAttribute{
				MarkdownDescription: "Specifies the Microsoft accounts that are supported for the current application. The possible values are: `AzureADMyOrg` (default), `AzureADMultipleOrgs`, `AzureADandPersonalMicrosoftAccount`, and `PersonalMicrosoftAccount`. The value of this object also limits the number of permissions an app can request. For more information, see Limits on requested permissions per app. The value for this property has implications on other app object properties. As a result, if you change this property, you may need to change other properties first. For more information, see Validation differences for signInAudience. Supports `$filter` (`eq`, `ne`, `not`).",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("AzureADMyOrg"),
				Validators: []validator.String{
					stringvalidator.OneOf(
						"AzureADMyOrg",
						"AzureADMultipleOrgs",
						"AzureADandPersonalMicrosoftAccount",
						"PersonalMicrosoftAccount",
					),
				},
			},
		"identifier_uris": schema.SetAttribute{
			MarkdownDescription: "Also known as App ID URI, this value is set when an application is used as a resource app. The identifierUris acts as the prefix for the scopes you reference in your API's code, and it must be globally unique across Microsoft Entra ID. For more information on valid identifierUris patterns and best practices, see Microsoft Entra application registration security best practices. Not nullable. Supports `$filter` (`eq`, `ne`, `ge`, `le`, `startsWith`).",
			Optional:            true,
			Computed:            true,
			ElementType:         types.StringType,
			PlanModifiers: []planmodifier.Set{
				setplanmodifier.UseStateForUnknown(),
			},
			Validators: []validator.Set{
				setvalidator.SizeAtLeast(1),
			},
		},
		"group_membership_claims": schema.SetAttribute{
			MarkdownDescription: "Configures the groups claim issued in a user or OAuth 2.0 access token that the application expects. To set this attribute, use one of the following string values: `None`, `SecurityGroup` (for security groups and Microsoft Entra roles), `All` (this gets all security groups, distribution groups, and Microsoft Entra directory roles that the signed-in user is a member of).",
			Optional:            true,
			Computed:            true,
			ElementType:         types.StringType,
			PlanModifiers: []planmodifier.Set{
				setplanmodifier.UseStateForUnknown(),
			},
			Validators: []validator.Set{
				setvalidator.SizeAtLeast(1),
				setvalidator.ValueStringsAre(
					stringvalidator.OneOf(
						"None",
						"SecurityGroup",
						"DirectoryRole",
						"ApplicationGroup",
						"All",
					),
				),
			},
		},
			"notes": schema.StringAttribute{
				MarkdownDescription: "Notes relevant for the management of the application.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"is_device_only_auth_supported": schema.BoolAttribute{
				MarkdownDescription: "Specifies whether this application supports device authentication without a user. The default is false.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"is_fallback_public_client": schema.BoolAttribute{
				MarkdownDescription: "Specifies the fallback application type as public client, such as an installed application running on a mobile device. The default value is false, which means the fallback application type is confidential client such as a web app. There are certain scenarios where Microsoft Entra ID can't determine the client application type. For example, the ROPC flow where the application is configured without specifying a redirect URI. In those cases Microsoft Entra ID interprets the application type based on the value of this property.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			// "oauth2_require_post_response": schema.BoolAttribute{
			// 	MarkdownDescription: "Specifies whether, as part of OAuth 2.0 token requests, Microsoft Entra ID allows POST requests, as opposed to GET requests. The default is false, which specifies that only GET requests are allowed.",
			// 	Optional:            true,
			// 	Computed:            true,
			// 	Default:             booldefault.StaticBool(false),
			// }, // Field doesn't exist in SDK v0.158.0.
			"service_management_reference": schema.StringAttribute{
				MarkdownDescription: "References application or service contact information from a Service or Asset Management database. Nullable.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		"tags": schema.SetAttribute{
			MarkdownDescription: "Custom strings that can be used to categorize and identify the application. Not nullable. Strings added here also appear in the tags property of any associated service principals. Supports `$filter` (`eq`, `not`, `ge`, `le`, `startsWith`) and `$search`.",
			Optional:            true,
			Computed:            true,
			ElementType:         types.StringType,
			PlanModifiers: []planmodifier.Set{
				setplanmodifier.UseStateForUnknown(),
			},
			Validators: []validator.Set{
				setvalidator.SizeAtLeast(1),
			},
		},
			"disabled_by_microsoft_status": schema.StringAttribute{
				MarkdownDescription: "Specifies whether Microsoft has disabled the registered application. The possible values are: null (default value), `NotDisabled`, and `DisabledDueToViolationOfServicesAgreement` (reasons may include suspicious, abusive, or malicious activity, or a violation of the Microsoft Services Agreement). Supports `$filter` (`eq`, `ne`, `not`). Read-only.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"publisher_domain": schema.StringAttribute{
				MarkdownDescription: "The verified publisher domain for the application. Read-only. Supports `$filter` (`eq`, `ne`, `ge`, `le`, `startsWith`).",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"created_date_time": schema.StringAttribute{
				MarkdownDescription: "The date and time the application was registered. The DateTimeOffset type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z. Read-only. Supports `$filter` (`eq`, `ne`, `not`, `ge`, `le`, `in`, and `eq` on null values) and `$orderby`.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"deleted_date_time": schema.StringAttribute{
				MarkdownDescription: "The date and time the application was deleted. The DateTimeOffset type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z. Read-only.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"api": schema.SingleNestedAttribute{
				MarkdownDescription: "Specifies settings for an application that implements a web API.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
				Attributes: map[string]schema.Attribute{
					"accept_mapped_claims": schema.BoolAttribute{
						MarkdownDescription: "Allows an application to use claims mapping without specifying a custom signing key.",
						Optional:            true,
						Computed:            true,
					},
				"known_client_applications": schema.SetAttribute{
					MarkdownDescription: "Used for bundling consent if you have a solution that contains two parts: a client app and a custom web API app. If you set the appID of the client app to this value, the user only consents once to the client app. Microsoft Entra ID knows that consenting to the client means implicitly consenting to the web API and automatically provisions service principals for both APIs at the same time. Both the client and the web API app must be registered in the same tenant.",
					Optional:            true,
					Computed:            true,
					ElementType:         types.StringType,
					Validators: []validator.Set{
						setvalidator.SizeAtLeast(1),
						setvalidator.ValueStringsAre(
							stringvalidator.RegexMatches(
								regexp.MustCompile(constants.GuidRegex),
								"must be a valid GUID in the format 00000000-0000-0000-0000-000000000000",
							),
						),
					},
				},
				"oauth2_permission_scopes": schema.SetNestedAttribute{
					MarkdownDescription: "The definition of the delegated permissions exposed by the web API represented by this application registration. These delegated permissions may be requested by a client application, and may be granted by users or administrators during consent. Delegated permissions are sometimes referred to as OAuth 2.0 scopes.",
					Optional:            true,
					Computed:            true,
					Validators: []validator.Set{
						setvalidator.SizeAtLeast(1),
					},
					NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"id": schema.StringAttribute{
									MarkdownDescription: "Unique scope permission identifier inside the oauth2PermissionScopes collection. Required.",
									Required:            true,
									Validators: []validator.String{
										stringvalidator.RegexMatches(
											regexp.MustCompile(constants.GuidRegex),
											"must be a valid GUID in the format 00000000-0000-0000-0000-000000000000",
										),
									},
								},
								"admin_consent_description": schema.StringAttribute{
									MarkdownDescription: "A description of the delegated permissions, intended to be read by an administrator granting the permission on behalf of all users. This text appears in tenant-wide admin consent experiences.",
									Optional:            true,
									Computed:            true,
								},
								"admin_consent_display_name": schema.StringAttribute{
									MarkdownDescription: "The permission's title, intended to be read by an administrator granting the permission on behalf of all users.",
									Optional:            true,
									Computed:            true,
								},
								"is_enabled": schema.BoolAttribute{
									MarkdownDescription: "When you create or update a permission, this property must be set to true (which is the default). To delete a permission, this property must first be set to false. At that point, in a subsequent call, the permission may be removed.",
									Optional:            true,
									Computed:            true,
									//Default:             booldefault.StaticBool(true),
								},
								"type": schema.StringAttribute{
									MarkdownDescription: "The possible values are: `User` and `Admin`. Specifies whether this delegated permission should be considered safe for non-admin users to consent to on behalf of themselves, or whether an administrator consent should be required for the permissions. While Microsoft Graph defines the default consent requirement for each permission, the tenant administrator may override the behavior in their organization (by allowing, restricting, or limiting user consent to this delegated permission). For more information, see Configure how users consent to applications.",
									Optional:            true,
									Computed:            true,
									Default:             stringdefault.StaticString("User"),
									Validators: []validator.String{
										stringvalidator.OneOf("User", "Admin"),
									},
								},
								"user_consent_description": schema.StringAttribute{
									MarkdownDescription: "A description of the delegated permissions, intended to be read by a user granting the permission on their own behalf. This text appears in consent experiences where the user is consenting only on behalf of themselves.",
									Optional:            true,
									Computed:            true,
								},
								"user_consent_display_name": schema.StringAttribute{
									MarkdownDescription: "A title for the permission, intended to be read by a user granting the permission on their own behalf. This text appears in consent experiences where the user is consenting only on behalf of themselves.",
									Optional:            true,
									Computed:            true,
								},
								"value": schema.StringAttribute{
									MarkdownDescription: "Specifies the value to include in the scp (scope) claim in access tokens. Must not exceed 120 characters in length. Allowed characters are : ! # $ % & ' ( ) * + , - . / : ;  =  ? @ [ ] ^ + _  {  } ~, and characters in the ranges 0-9, A-Z and a-z. Any other character, including the space character, aren't allowed. May not begin with `.`.",
									Optional:            true,
									Computed:            true,
									Validators: []validator.String{
										validate.StringLengthAtMost(120),
									},
								},
							},
						},
					},
				"pre_authorized_applications": schema.SetNestedAttribute{
					MarkdownDescription: "Lists the client applications that are preauthorized with the specified delegated permissions to access this application's APIs. Users aren't required to consent to any preauthorized application (for the permissions specified). However, any other permissions not listed in preAuthorizedApplications (requested through incremental consent for example) will require user consent.",
					Optional:            true,
					Computed:            true,
					Validators: []validator.Set{
						setvalidator.SizeAtLeast(1),
					},
					NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"app_id": schema.StringAttribute{
									MarkdownDescription: "The unique identifier for the client application.",
									Required:            true,
									Validators: []validator.String{
										stringvalidator.RegexMatches(
											regexp.MustCompile(constants.GuidRegex),
											"must be a valid GUID in the format 00000000-0000-0000-0000-000000000000",
										),
									},
								},
								"delegated_permission_ids": schema.SetAttribute{
									MarkdownDescription: "The unique identifier for the scopes the client application is granted.",
									Required:            true,
									ElementType:         types.StringType,
									Validators: []validator.Set{
										setvalidator.ValueStringsAre(
											stringvalidator.RegexMatches(
												regexp.MustCompile(constants.GuidRegex),
												"must be a valid GUID in the format 00000000-0000-0000-0000-000000000000",
											),
										),
									},
								},
							},
						},
					},
					"requested_access_token_version": schema.Int32Attribute{
						MarkdownDescription: "Specifies the access token version expected by this resource. This changes the version and format of the JWT produced independent of the endpoint or client used to request the access token. The endpoint used, v1.0 or v2.0, is chosen by the client and only impacts the version of id_tokens. Resources need to explicitly configure requestedAccessTokenVersion to indicate the supported access token format. Possible values for requestedAccessTokenVersion are 1, 2, or null. If the value is null, this defaults to 1, which corresponds to the v1.0 endpoint. If signInAudience on the application is configured as AzureADandPersonalMicrosoftAccount or PersonalMicrosoftAccount, the value for this property must be 2.",
						Optional:            true,
						Computed:            true,
						Default:             int32default.StaticInt32(1),
						Validators: []validator.Int32{
							int32validator.Between(1, 2),
						},
					},
				},
			},
		"app_roles": schema.SetNestedAttribute{
			MarkdownDescription: "The collection of roles defined for the application. With app role assignments, these roles can be assigned to users, groups, or service principals associated with other applications. Not nullable.",
			Optional:            true,
			Computed:            true,
			PlanModifiers: []planmodifier.Set{
				setplanmodifier.UseStateForUnknown(),
			},
			Validators: []validator.Set{
				setvalidator.SizeAtLeast(1),
			},
			NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							MarkdownDescription: "Unique role identifier inside the appRoles collection. When creating a new app role, a new GUID identifier must be provided. Required.",
							Required:            true,
							Validators: []validator.String{
								stringvalidator.RegexMatches(
									regexp.MustCompile(constants.GuidRegex),
									"must be a valid GUID in the format 00000000-0000-0000-0000-000000000000",
								),
							},
						},
						"allowed_member_types": schema.SetAttribute{
							MarkdownDescription: "Specifies whether this app role can be assigned to users and groups (by setting to `['User']`), to other application's (by setting to `['Application']`, or both (by setting to `['User', 'Application']`). App roles supporting assignment to other applications' service principals are also known as application permissions. The 'Application' value is only supported for app roles defined on application entities. Required.",
							Required:            true,
							ElementType:         types.StringType,
							Validators: []validator.Set{
								setvalidator.ValueStringsAre(
									stringvalidator.OneOf("User", "Application"),
								),
							},
						},
						"description": schema.StringAttribute{
							MarkdownDescription: "The description for the app role. This is displayed when the app role is being assigned and, if the app role functions as an application permission, during consent experiences. Required.",
							Required:            true,
						},
						"display_name": schema.StringAttribute{
							MarkdownDescription: "Display name for the permission that appears in the app role assignment and consent experiences. Required.",
							Required:            true,
						},
						"is_enabled": schema.BoolAttribute{
							MarkdownDescription: "Defines whether the application's app role is enabled or disabled. Required.",
							Required:            true,
							PlanModifiers: []planmodifier.Bool{
								boolplanmodifier.UseStateForUnknown(),
							},
						},
						"origin": schema.StringAttribute{
							MarkdownDescription: "Specifies if the app role is defined on the application object or on the servicePrincipal entity. Must not be included in any POST or PATCH requests. Read-only.",
							Computed:            true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.UseStateForUnknown(),
							},
						},
						"value": schema.StringAttribute{
							MarkdownDescription: "Specifies the value to include in the roles claim in ID tokens and access tokens authenticating an assigned user or service principal. Must not exceed 120 characters in length. Allowed characters are : ! # $ % & ' ( ) * + , - . / : ;  =  ? @ [ ] ^ + _  {  } ~, and characters in the ranges 0-9, A-Z and a-z. Any other character, including the space character, aren't allowed. May not begin with `.`. Nullable.",
							Optional:            true,
							Computed:            true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.UseStateForUnknown(),
							},
							Validators: []validator.String{
								validate.StringLengthAtMost(120),
							},
						},
					},
				},
			},
			"info": schema.SingleNestedAttribute{
				MarkdownDescription: "Basic profile information of the application, such as it's marketing, support, terms of service, and privacy statement URLs. The terms of service and privacy statement are surfaced to users through the user consent experience. For more information, see How to: Add Terms of service and privacy statement for registered Microsoft Entra apps. Supports `$filter` (`eq`, `ne`, `not`, `ge`, `le`, and `eq` on null values).",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
				Attributes: map[string]schema.Attribute{
					"logo_url": schema.StringAttribute{
						MarkdownDescription: "CDN URL to the application's logo. Read-only.",
						Computed:            true,
					},
					"marketing_url": schema.StringAttribute{
						MarkdownDescription: "Link to the application's marketing page. For example, https://www.contoso.com/app/marketing.",
						Optional:            true,
						Computed:            true,
						Validators: []validator.String{
							stringvalidator.RegexMatches(
								regexp.MustCompile(constants.HttpOrHttpsUrlRegex),
								"must be a valid HTTP or HTTPS URL",
							),
						},
					},
					"privacy_statement_url": schema.StringAttribute{
						MarkdownDescription: "Link to the application's privacy statement. For example, https://www.contoso.com/app/privacy.",
						Optional:            true,
						Computed:            true,
						Validators: []validator.String{
							stringvalidator.RegexMatches(
								regexp.MustCompile(constants.HttpOrHttpsUrlRegex),
								"must be a valid HTTP or HTTPS URL",
							),
						},
					},
					"support_url": schema.StringAttribute{
						MarkdownDescription: "Link to the application's support page. For example, https://www.contoso.com/app/support.",
						Optional:            true,
						Computed:            true,
						Validators: []validator.String{
							stringvalidator.RegexMatches(
								regexp.MustCompile(constants.HttpOrHttpsUrlRegex),
								"must be a valid HTTP or HTTPS URL",
							),
						},
					},
					"terms_of_service_url": schema.StringAttribute{
						MarkdownDescription: "Link to the application's terms of service statement. For example, https://www.contoso.com/app/termsofservice.",
						Optional:            true,
						Computed:            true,
						Validators: []validator.String{
							stringvalidator.RegexMatches(
								regexp.MustCompile(constants.HttpOrHttpsUrlRegex),
								"must be a valid HTTP or HTTPS URL",
							),
						},
					},
				},
			},
			"optional_claims": schema.SingleNestedAttribute{
				MarkdownDescription: "Application developers can configure optional claims in their Microsoft Entra applications to specify the claims that are sent to their application by the Microsoft security token service. For more information, see How to: Provide optional claims to your app.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
				Attributes: map[string]schema.Attribute{
				"access_token": schema.SetNestedAttribute{
					MarkdownDescription: "The optional claims returned in the JWT access token.",
					Optional:            true,
					Computed:            true,
					Validators: []validator.Set{
						setvalidator.SizeAtLeast(1),
					},
					NestedObject: schema.NestedAttributeObject{
						Attributes: optionalClaimAttributes(),
					},
				},
				"id_token": schema.SetNestedAttribute{
					MarkdownDescription: "The optional claims returned in the JWT ID token.",
					Optional:            true,
					Computed:            true,
					Validators: []validator.Set{
						setvalidator.SizeAtLeast(1),
					},
					NestedObject: schema.NestedAttributeObject{
						Attributes: optionalClaimAttributes(),
					},
				},
				"saml2_token": schema.SetNestedAttribute{
					MarkdownDescription: "The optional claims returned in the SAML token.",
					Optional:            true,
					Computed:            true,
					Validators: []validator.Set{
						setvalidator.SizeAtLeast(1),
					},
					NestedObject: schema.NestedAttributeObject{
						Attributes: optionalClaimAttributes(),
					},
				},
				},
			},
			"parental_control_settings": schema.SingleNestedAttribute{
				MarkdownDescription: "Specifies parental control settings for an application.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
				Attributes: map[string]schema.Attribute{
				"countries_blocked_for_minors": schema.SetAttribute{
					MarkdownDescription: "Specifies the two-letter ISO country codes. Access to the application will be blocked for minors from the countries specified in this list.",
					Optional:            true,
					Computed:            true,
					ElementType:         types.StringType,
					Validators: []validator.Set{
						setvalidator.SizeAtLeast(1),
						setvalidator.ValueStringsAre(
							stringvalidator.RegexMatches(
								regexp.MustCompile(`^[A-Z]{2}$`),
								"must be a two-letter ISO country code",
							),
						),
					},
				},
					"legal_age_group_rule": schema.StringAttribute{
						MarkdownDescription: "Specifies the legal age group rule that applies to users of the app. Can be set to one of the following values: `Allow`, `RequireConsentForPrivacyServices`, `RequireConsentForMinors`, `RequireConsentForKids`, `BlockMinors`.",
						Optional:            true,
						Computed:            true,
						Validators: []validator.String{
							stringvalidator.OneOf(
								"Allow",
								"RequireConsentForPrivacyServices",
								"RequireConsentForMinors",
								"RequireConsentForKids",
								"BlockMinors",
							),
						},
					},
				},
			},
			"public_client": schema.SingleNestedAttribute{
				MarkdownDescription: "Specifies settings for installed clients such as desktop or mobile devices.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
				Attributes: map[string]schema.Attribute{
					"redirect_uris": schema.SetAttribute{
						MarkdownDescription: "Specifies the URLs where user tokens are sent for sign-in, or the redirect URIs where OAuth 2.0 authorization codes and access tokens are sent.",
						Optional:            true,
						Computed:            true,
						ElementType:         types.StringType,
						Validators: []validator.Set{
							setvalidator.SizeAtLeast(1),
							setvalidator.SizeAtMost(256),
						},
					},
				},
			},
		"required_resource_access": schema.SetNestedAttribute{
			MarkdownDescription: "Specifies the resources that the application needs to access. This property also specifies the set of delegated permissions and application roles that it needs for each of those resources. This configuration of access to the required resources drives the consent experience. No more than 50 resource services (APIs) can be configured. Beginning mid-October 2021, the total number of required permissions must not exceed 400. For more information, see Limits on requested permissions per app. Not nullable. Supports `$filter` (`eq`, `not`, `ge`, `le`).",
			Optional:            true,
			Computed:            true,
			PlanModifiers: []planmodifier.Set{
				setplanmodifier.UseStateForUnknown(),
			},
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"resource_app_id": schema.StringAttribute{
						MarkdownDescription: "The unique identifier for the resource that the application requires access to. This should be equal to the appId declared on the target resource application. Required.",
						Required:            true,
						Validators: []validator.String{
							stringvalidator.RegexMatches(
								regexp.MustCompile(constants.GuidRegex),
								"must be a valid GUID in the format 00000000-0000-0000-0000-000000000000",
							),
						},
					},
					"resource_access": schema.SetNestedAttribute{
						MarkdownDescription: "The list of OAuth2.0 permission scopes and app roles that the application requires from the specified resource. Required.",
						Required:            true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"id": schema.StringAttribute{
									MarkdownDescription: "The unique identifier of an app role or delegated permission exposed by the resource application. For delegated permissions, this should match the id property of one of the delegated permissions in the oauth2PermissionScopes collection of the resource application's service principal. For app roles (application permissions), this should match the id property of an app role in the appRoles collection of the resource application's service principal. Required.",
									Required:            true,
									Validators: []validator.String{
										stringvalidator.RegexMatches(
											regexp.MustCompile(constants.GuidRegex),
											"must be a valid GUID in the format 00000000-0000-0000-0000-000000000000",
										),
									},
								},
								"type": schema.StringAttribute{
									MarkdownDescription: "Specifies whether the id property references a delegated permission or an app role (application permission). The possible values are: `Scope` (for delegated permissions) or `Role` (for app roles). Required.",
									Required:            true,
									Validators: []validator.String{
										stringvalidator.OneOf("Scope", "Role"),
									},
								},
							},
						},
					},
				},
			},
			Validators: []validator.Set{
				setvalidator.SizeAtLeast(1),
				setvalidator.SizeAtMost(50),
			},
		},
			"spa": schema.SingleNestedAttribute{
				MarkdownDescription: "Specifies settings for a single-page application, including sign out URLs and redirect URIs for authorization codes and access tokens.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
				Attributes: map[string]schema.Attribute{
					"redirect_uris": schema.SetAttribute{
						MarkdownDescription: "Specifies the URLs where user tokens are sent for sign-in, or the redirect URIs where OAuth 2.0 authorization codes and access tokens are sent.",
						Optional:            true,
						Computed:            true,
						ElementType:         types.StringType,
						Validators: []validator.Set{
							setvalidator.SizeAtLeast(1),
							setvalidator.SizeAtMost(256),
						},
					},
				},
			},
			"web": schema.SingleNestedAttribute{
				MarkdownDescription: "Specifies settings for a web application.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
				Attributes: map[string]schema.Attribute{
					"home_page_url": schema.StringAttribute{
						MarkdownDescription: "Home page or landing page of the application.",
						Optional:            true,
						Computed:            true,
						Validators: []validator.String{
							stringvalidator.RegexMatches(
								regexp.MustCompile(constants.HttpOrHttpsUrlRegex),
								"must be a valid HTTP or HTTPS URL",
							),
						},
					},
					"logout_url": schema.StringAttribute{
						MarkdownDescription: "Specifies the URL that is used by Microsoft's authorization service to log out a user using front-channel, back-channel or SAML logout protocols.",
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
						Validators: []validator.String{
							stringvalidator.RegexMatches(
								regexp.MustCompile(constants.HttpOrHttpsUrlRegex),
								"must be a valid HTTP or HTTPS URL",
							),
						},
					},
					"redirect_uris": schema.SetAttribute{
						MarkdownDescription: "Specifies the URLs where user tokens are sent for sign-in, or the redirect URIs where OAuth 2.0 authorization codes and access tokens are sent.",
						Optional:            true,
						Computed:            true,
						ElementType:         types.StringType,
						Validators: []validator.Set{
							setvalidator.SizeAtLeast(1),
							setvalidator.SizeAtMost(256),
						},
					},
					"implicit_grant_settings": schema.SingleNestedAttribute{
						MarkdownDescription: "Specifies whether this web application can request tokens using the OAuth 2.0 implicit flow.",
						Optional:            true,
						Computed:            true,
						Attributes: map[string]schema.Attribute{
							"enable_access_token_issuance": schema.BoolAttribute{
								MarkdownDescription: "Specifies whether this web application can request an access token using the OAuth 2.0 implicit flow.",
								Optional:            true,
								Computed:            true,
							},
							"enable_id_token_issuance": schema.BoolAttribute{
								MarkdownDescription: "Specifies whether this web application can request an ID token using the OAuth 2.0 implicit flow.",
								Optional:            true,
								Computed:            true,
							},
						},
					},
					"redirect_uri_settings": schema.SetNestedAttribute{
						MarkdownDescription: "Specifies the index of the URLs where user tokens are sent for sign-in. This is only valid for applications using SAML. Note: If not specified, the API may auto-generate settings based on redirect_uris. To manage this field, you must provide at least one entry; empty arrays are not supported as the API auto-generates values.",
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.Set{
							setplanmodifier.UseStateForUnknown(),
						},
						Validators: []validator.Set{
							setvalidator.SizeAtLeast(1),
						},
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"uri": schema.StringAttribute{
									MarkdownDescription: "Specifies the URI that tokens are sent to.",
									Optional:            true,
									Computed:            true,
									Validators: []validator.String{
										stringvalidator.RegexMatches(
											regexp.MustCompile(constants.HttpOrHttpsUrlRegex),
											"must be a valid HTTP or HTTPS URL",
										),
									},
								},
								"index": schema.Int32Attribute{
									MarkdownDescription: "Identifies the specific URI within the redirectURIs collection in SAML SSO flows. Defaults to null. The index is unique across all the redirectUris for the application.",
									Optional:            true,
									Computed:            true,
								},
							},
						},
					},
				},
			},
			"sign_in_audience_restrictions": schema.SingleNestedAttribute{
				MarkdownDescription: "Specifies restrictions on the supported account types specified in signInAudience. The value type determines the restrictions that can be applied: unrestrictedAudience (There are no additional restrictions on the supported account types allowed by signInAudience) or allowedTenantsAudience (The application can only be used in the specified Entra tenants. Only supported when signInAudience is AzureADMultipleOrgs). Default is a value of type unrestrictedAudience. Returned only on `$select`.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
				Attributes: map[string]schema.Attribute{
					"odata_type": schema.StringAttribute{
						MarkdownDescription: "The OData type. Must be `#microsoft.graph.allowedTenantsAudience` or `#microsoft.graph.unrestrictedAudience`.",
						Optional:            true,
						Computed:            true,
						Validators: []validator.String{
							stringvalidator.OneOf(
								"#microsoft.graph.allowedTenantsAudience",
								"#microsoft.graph.unrestrictedAudience",
							),
						},
					},
					"is_home_tenant_allowed": schema.BoolAttribute{
						MarkdownDescription: "Indicates whether the home tenant is allowed. Only applicable when odata_type is `#microsoft.graph.allowedTenantsAudience`.",
						Optional:            true,
						Computed:            true,
					},
				"allowed_tenant_ids": schema.SetAttribute{
					MarkdownDescription: "The list of allowed tenant IDs. Only applicable when odata_type is `#microsoft.graph.allowedTenantsAudience`.",
					Optional:            true,
					Computed:            true,
					ElementType:         types.StringType,
					Validators: []validator.Set{
						setvalidator.SizeAtLeast(1),
						setvalidator.ValueStringsAre(
							stringvalidator.RegexMatches(
								regexp.MustCompile(constants.GuidRegex),
								"must be a valid GUID in the format 00000000-0000-0000-0000-000000000000",
							),
						),
					},
				},
				},
			},
		"owner_user_ids": schema.SetAttribute{
			MarkdownDescription: "The user IDs of the owners for the application. At least one owner is typically required when creating an application. Owners are a set of non-admin users or service principals allowed to modify this object.",
			Optional:            true,
			Computed:            true,
			ElementType:         types.StringType,
			PlanModifiers: []planmodifier.Set{
				setplanmodifier.UseStateForUnknown(),
			},
			Validators: []validator.Set{
				setvalidator.SizeAtLeast(1),
				setvalidator.SizeAtMost(100),
				setvalidator.ValueStringsAre(
					stringvalidator.RegexMatches(
						regexp.MustCompile(constants.GuidRegex),
						"must be a valid GUID in the format 00000000-0000-0000-0000-000000000000",
					),
				),
			},
		},
			"prevent_duplicate_names": schema.BoolAttribute{
				MarkdownDescription: "If set to `true`, the provider will check for existing applications with the same display " +
					"name and return an error if one is found. This helps prevent accidentally creating duplicate applications. Note: " +
					"This field defaults to `false` on import since the API does not return this value.",
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(false),
			},
			"hard_delete": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(false),
				MarkdownDescription: "When `true`, the application will be permanently deleted (hard delete) during destroy. " +
					"When `false` (default), the application will only be soft deleted and moved to the deleted items container " +
					"where it can be restored within 30 days. " +
					"Note: This field defaults to `false` on import since the API does not return this value.",
			},
			"timeouts": commonschema.ResourceTimeouts(ctx),
		},
	}
}

// optionalClaimAttributes returns the common attributes for optional claims
func optionalClaimAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"name": schema.StringAttribute{
			MarkdownDescription: "The name of the optional claim. Required.",
			Required:            true,
		},
		"source": schema.StringAttribute{
			MarkdownDescription: "The source (directory object) of the claim. There are predefined claims and user-defined claims from extension properties. If the source value is null, the claim is a predefined optional claim. If the source value is user, the value in the name property is the extension property from the user object.",
			Optional:            true,
			Computed:            true,
		},
		"essential": schema.BoolAttribute{
			MarkdownDescription: "If the value is true, the claim specified by the client is necessary to ensure a smooth authorization experience for the specific task requested by the end user. The default value is false.",
			Optional:            true,
			Computed:            true,
			Default:             booldefault.StaticBool(false),
		},
		"additional_properties": schema.SetAttribute{
			MarkdownDescription: "Additional properties of the claim. If a property exists in this collection, it modifies the behavior of the optional claim specified in the name property.",
			Optional:            true,
			Computed:            true,
			ElementType:         types.StringType,
			Validators: []validator.Set{
				setvalidator.SizeAtLeast(1),
			},
		},
	}
}
