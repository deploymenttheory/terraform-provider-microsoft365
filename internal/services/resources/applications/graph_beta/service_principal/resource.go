package graphBetaServicePrincipal

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/validate/attribute"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/identityschema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	ResourceName  = "microsoft365_graph_beta_applications_service_principal"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	// Basic resource interface (CRUD operations)
	_ resource.Resource = &ServicePrincipalResource{}

	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &ServicePrincipalResource{}

	_ resource.ResourceWithImportState = &ServicePrincipalResource{}

	// Enables identity schema for list resource support
	_ resource.ResourceWithIdentity = &ServicePrincipalResource{}
)

func NewServicePrincipalResource() resource.Resource {
	return &ServicePrincipalResource{
		ReadPermissions: []string{
			"Application.Read.All",
			"Directory.Read.All",
		},
		WritePermissions: []string{
			"Application.ReadWrite.All",
			"Directory.ReadWrite.All",
		},
		ResourcePath: "/servicePrincipals",
	}
}

type ServicePrincipalResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

func (r *ServicePrincipalResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = ResourceName
}

func (r *ServicePrincipalResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, ResourceName)
}

// ImportState handles importing the resource with an extended ID format.
//
// Supported formats:
//   - Simple:   "resource_id" (hard_delete defaults to false)
//   - Extended: "resource_id:hard_delete=true"
//
// Example:
//
//	terraform import microsoft365_graph_beta_applications_service_principal.example "12345678-1234-1234-1234-123456789012:hard_delete=true"
func (r *ServicePrincipalResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	idParts := strings.Split(req.ID, ":")
	resourceID := idParts[0]
	hardDelete := false // Default

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

// IdentitySchema defines the identity schema for this resource, used by list operations to uniquely identify instances
func (r *ServicePrincipalResource) IdentitySchema(ctx context.Context, req resource.IdentitySchemaRequest, resp *resource.IdentitySchemaResponse) {
	resp.IdentitySchema = identityschema.Schema{
		Attributes: map[string]identityschema.Attribute{
			"id": identityschema.StringAttribute{
				RequiredForImport: true,
			},
		},
	}
}

func (r *ServicePrincipalResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a Service Principal in Microsoft Entra ID. " +
			"Service principals are the local representation of an application object in a specific tenant. " +
			"They define what the app can do in the specific tenant, who can access the app, and what resources the app can access.\n\n" +
			"This resource models only the properties the service principal itself owns. Properties that Microsoft Entra reflects from the " +
			"backing application — such as its display name, home page, logout URL, reply URLs, app roles and permission scopes — are " +
			"deliberately not exposed here and will not be added: they change whenever the application changes. Read them from the " +
			"`microsoft365_graph_beta_applications_application` resource, or from the service principal data source for an application " +
			"this configuration does not manage.\n\n" +
			"For more information, see the [Microsoft Graph API documentation](https://learn.microsoft.com/en-us/graph/api/serviceprincipal-post-serviceprincipals?view=graph-rest-beta).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique identifier (object ID) for the service principal. Read-only.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"app_id": schema.StringAttribute{
				MarkdownDescription: "The application (client) ID of the application for which to create the service principal. Required.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"alternative_names": schema.SetAttribute{
				MarkdownDescription: "Used to retrieve service principals by subscription, identify resource group and full resource IDs for managed identities. " +
					"Set to `[]` to clear previously configured values. Supports `$filter` (`eq`, `not`, `ge`, `le`, `startsWith`).",
				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"app_owner_organization_id": schema.StringAttribute{
				MarkdownDescription: "Contains the tenant ID where the application is registered. Applicable only to service principals backed by applications. Read-only. " +
					"Equivalent to `application_tenant_id` in the azuread provider.",
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"application_template_id": schema.StringAttribute{
				MarkdownDescription: "Unique identifier of the application template that the associated application was created from. Read-only. `null` if the app wasn't created from an application template.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"created_by_app_id": schema.StringAttribute{
				MarkdownDescription: "The appId of the application that created this service principal. Set internally by Microsoft Entra ID. Read-only.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"is_disabled": schema.BoolAttribute{
				MarkdownDescription: "Specifies whether the service principal is deactivated so the app can't obtain new access tokens or access protected resources. Read-only; the API rejects writes to this property.",
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"key_credentials": schema.SetNestedAttribute{
				MarkdownDescription: "The collection of key credentials associated with the service principal. Read-only on this resource; certificates are added through dedicated credential resources or the addTokenSigningCertificate API. " +
					"Private key material is never returned by the API, and this resource does not expose the public `key` field.",
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"custom_key_identifier": schema.StringAttribute{
							MarkdownDescription: "A base64-encoded custom key identifier.",
							Computed:            true,
						},
						"display_name": schema.StringAttribute{
							MarkdownDescription: "The friendly name for the key.",
							Computed:            true,
						},
						"end_date_time": schema.StringAttribute{
							MarkdownDescription: "The date and time at which the credential expires.",
							Computed:            true,
						},
						"key_id": schema.StringAttribute{
							MarkdownDescription: "The unique identifier for the key.",
							Computed:            true,
						},
						"start_date_time": schema.StringAttribute{
							MarkdownDescription: "The date and time at which the credential becomes valid.",
							Computed:            true,
						},
						"type": schema.StringAttribute{
							MarkdownDescription: "The type of key credential, for example `Symmetric` or `AsymmetricX509Cert`.",
							Computed:            true,
						},
						"usage": schema.StringAttribute{
							MarkdownDescription: "A string that describes the purpose for which the key can be used, for example `Verify` or `Sign`.",
							Computed:            true,
						},
					},
				},
			},
			"password_credentials": schema.SetNestedAttribute{
				MarkdownDescription: "The collection of password credentials associated with the service principal. Read-only. The secret itself is never returned by the API.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"custom_key_identifier": schema.StringAttribute{
							MarkdownDescription: "A base64-encoded custom key identifier.",
							Computed:            true,
						},
						"display_name": schema.StringAttribute{
							MarkdownDescription: "The friendly name for the password.",
							Computed:            true,
						},
						"end_date_time": schema.StringAttribute{
							MarkdownDescription: "The date and time at which the password expires.",
							Computed:            true,
						},
						"hint": schema.StringAttribute{
							MarkdownDescription: "Contains the first three characters of the password. Read-only.",
							Computed:            true,
						},
						"key_id": schema.StringAttribute{
							MarkdownDescription: "The unique identifier for the password.",
							Computed:            true,
						},
						"start_date_time": schema.StringAttribute{
							MarkdownDescription: "The date and time at which the password becomes valid.",
							Computed:            true,
						},
					},
				},
			},
			"preferred_token_signing_key_end_date_time": schema.StringAttribute{
				MarkdownDescription: "Specifies the expiration date of the key credential used for token signing, marked by `preferred_token_signing_key_thumbprint`. Read-only.",
				Computed:            true,
			},
			"preferred_token_signing_key_thumbprint": schema.StringAttribute{
				MarkdownDescription: "The thumbprint of the certificate used to sign SAML responses for applications with `preferred_single_sign_on_mode` set to `saml`. " +
					"Read-only on this resource; it is set when a token signing certificate is activated on the service principal.",
				Computed: true,
			},
			"token_encryption_key_id": schema.StringAttribute{
				MarkdownDescription: "Specifies the keyId of a public key from the key credentials collection. When configured, Microsoft Entra ID issues tokens for this application encrypted using the key specified by this property. " +
					"The referenced key credential must already exist on the service principal. " +
					"When this attribute is absent from the configuration, updates actively clear the property on the service principal.",
				Optional: true,
				Validators: []validator.String{
					attribute.RegexMatches(regexp.MustCompile(constants.GuidRegex), "must be a valid GUID"),
				},
			},
			"account_enabled": schema.BoolAttribute{
				MarkdownDescription: "True if the service principal account is enabled; otherwise, false. If set to false, then no users are able to sign in to this app, even if they're assigned to it. Supports `$filter` (`eq`, `ne`, `not`, `in`).",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"app_role_assignment_required": schema.BoolAttribute{
				MarkdownDescription: "Specifies whether users or other service principals need to be granted an app role assignment for this service principal before users can sign in or apps can get tokens. The default value is false. Not nullable. Supports `$filter` (`eq`, `ne`, `NOT`).",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Free text field to provide an internal end-user facing description of the service principal. End-user portals such MyApps displays the application description in this field. The maximum allowed size is 1,024 characters. Supports `$filter` (`eq`, `ne`, `not`, `ge`, `le`, `startsWith`) and `$search`.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(0, 1024),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"login_url": schema.StringAttribute{
				MarkdownDescription: "Specifies the URL where the service provider redirects the user to Microsoft Entra ID to authenticate. Microsoft Entra ID uses the URL to launch the application from Microsoft 365 or the Microsoft Entra My Apps. When blank, Microsoft Entra ID performs IdP-initiated sign-on for applications configured with SAML-based single sign-on.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					attribute.RegexMatches(regexp.MustCompile(constants.HttpOrHttpsUrlRegex), "must be a valid HTTP or HTTPS URL"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"notes": schema.StringAttribute{
				MarkdownDescription: "Free text field to capture information about the service principal, typically used for operational purposes. Maximum allowed size is 1,024 characters.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(0, 1024),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"notification_email_addresses": schema.SetAttribute{
				MarkdownDescription: "Specifies the list of email addresses where Microsoft Entra ID sends a notification when the active certificate is near the expiration date. This is only for the certificates used to sign the SAML token issued for Microsoft Entra Gallery applications.",
				ElementType:         types.StringType,
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.Set{
					setvalidator.SizeAtLeast(1),
				},
			},
			"preferred_single_sign_on_mode": schema.StringAttribute{
				MarkdownDescription: "Specifies the single sign-on mode configured for this application. Microsoft Entra ID uses the preferred single sign-on mode to launch the application from Microsoft 365 or the Microsoft Entra My Apps. The supported values are `password`, `saml`, `notSupported`, and `oidc`.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("password", "saml", "notSupported", "oidc"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"saml_single_sign_on_settings": schema.SingleNestedAttribute{
				MarkdownDescription: "The collection for settings related to SAML single sign-on. " +
					"When this block is absent from the configuration, updates actively clear the property on the service principal.",
				Optional: true,
				Attributes: map[string]schema.Attribute{
					"relay_state": schema.StringAttribute{
						MarkdownDescription: "The relative URI the service provider would redirect to after completion of the single sign-on flow.",
						Required:            true,
					},
				},
			},
			"service_principal_type": schema.StringAttribute{
				MarkdownDescription: "Identifies if the service principal represents an application or a managed identity. Read-only.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"tags": schema.SetAttribute{
				MarkdownDescription: "Custom strings that can be used to categorize and identify the service principal. " +
					"Note: Microsoft may automatically add system-managed tags in addition to the tags you specify.",
				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.Set{
					setvalidator.SizeAtLeast(1),
				},
			},
			"hard_delete": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(false),
				MarkdownDescription: "When `true`, the service principal will be permanently deleted (hard delete) during destroy. " +
					"When `false` (default), the service principal will only be soft deleted and moved to the deleted items container " +
					"where it can be restored within 30 days. " +
					"Note: This field defaults to `false` on import since the API does not return this value.",
			},
			"timeouts": commonschema.ResourceTimeouts(ctx),
		},
	}
}
