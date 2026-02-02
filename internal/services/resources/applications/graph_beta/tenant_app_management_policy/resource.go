package graphBetaTenantAppManagementPolicy

import (
	"context"
	"regexp"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	planmodifiers "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/plan_modifiers"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	ResourceName  = "microsoft365_graph_beta_applications_tenant_app_management_policy"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	// Basic resource interface (CRUD operations)
	_ resource.Resource = &TenantAppManagementPolicyResource{}

	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &TenantAppManagementPolicyResource{}

	// Enables import functionality
	_ resource.ResourceWithImportState = &TenantAppManagementPolicyResource{}
)

func NewTenantAppManagementPolicyResource() resource.Resource {
	return &TenantAppManagementPolicyResource{
		ReadPermissions: []string{
			"Policy.Read.All",
			"Policy.Read.ApplicationConfiguration",
		},
		WritePermissions: []string{
			"Policy.ReadWrite.ApplicationConfiguration",
		},
		ResourcePath: "/policies/defaultAppManagementPolicy",
	}
}

type TenantAppManagementPolicyResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

// Metadata returns the resource type name.
func (r *TenantAppManagementPolicyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = ResourceName
}

// Configure sets the client for the resource.
func (r *TenantAppManagementPolicyResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, ResourceName)
}

// ImportState imports the resource state.
func (r *TenantAppManagementPolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Schema defines the schema for the resource.
func (r *TenantAppManagementPolicyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages the tenant-wide default app management policy for applications and service principals using the `/policies/defaultAppManagementPolicy` endpoint. " +
			"This policy enforces app management restrictions such as password/key credential lifetimes, identifier URI restrictions, and certificate authority requirements. " +
			"The policy applies to all applications and service principals unless overridden by a specific appManagementPolicy.\n\n" +
			"For more information, see the [Microsoft Graph API documentation](https://learn.microsoft.com/en-us/graph/api/resources/tenantappmanagementpolicy?view=graph-rest-beta).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
				MarkdownDescription: "The unique identifier for the tenant app management policy. This is always '00000000-0000-0000-0000-000000000000' for the default policy.",
			},
			"display_name": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
				MarkdownDescription: "The display name of the policy. Inherited from policyBase.",
			},
			"description": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
				MarkdownDescription: "The description of the policy. Inherited from policyBase.",
			},
			"is_enabled": schema.BoolAttribute{
				Required:            true,
				MarkdownDescription: "Denotes whether the policy is enabled. Default value is false. When false, restrictions are not evaluated or enforced.",
			},
			"restore_to_default_upon_delete": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(false),
				MarkdownDescription: "Terraform-specific attribute. When true, deleting this resource will restore the policy to Microsoft's default settings (isEnabled=false, empty restrictions). " +
					"When false (default), the resource is only removed from Terraform state without changing the policy in Microsoft Graph.",
			},
			"application_restrictions": schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{
					"password_credentials": schema.ListNestedAttribute{
						Optional: true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: getPasswordCredentialConfigurationAttributes(),
						},
						MarkdownDescription: "Collection of password credential restrictions for applications.",
					},
					"key_credentials": schema.ListNestedAttribute{
						Optional: true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: getKeyCredentialConfigurationAttributes(),
						},
						MarkdownDescription: "Collection of key credential restrictions for applications.",
					},
					"identifier_uris": schema.SingleNestedAttribute{
						Optional: true,
						Attributes: map[string]schema.Attribute{
							"non_default_uri_addition": schema.SingleNestedAttribute{
								Optional: true,
								Attributes: map[string]schema.Attribute{
									"restrict_for_apps_created_after_date_time": schema.StringAttribute{
										Required: true,
										Validators: []validator.String{
											stringvalidator.RegexMatches(
												regexp.MustCompile(constants.ISO8601DateTimeRegex),
												"must be a valid ISO 8601 datetime in the format YYYY-MM-DDTHH:MM:SSZ or with timezone offset",
											),
										},
										MarkdownDescription: "Specifies the date from which the policy restriction applies to newly created applications. For existing applications, the enforcement date can be retroactively applied. Format: RFC3339/DateTimeOffset (e.g., '2024-01-01T10:37:00Z').",
									},
									"exclude_apps_receiving_v2_tokens": schema.BoolAttribute{
										Optional:            true,
										MarkdownDescription: "When true, excludes applications receiving v2 tokens from this restriction.",
									},
									"exclude_saml": schema.BoolAttribute{
										Optional:            true,
										MarkdownDescription: "When true, excludes SAML applications from this restriction.",
									},
									"exclude_actors": schema.SingleNestedAttribute{
										Optional:            true,
										Attributes:          getActorExemptionsAttributes(),
										MarkdownDescription: "Exemptions based on custom security attributes.",
									},
								},
								MarkdownDescription: "Restrictions on adding non-default identifier URIs.",
							},
						},
						MarkdownDescription: "Identifier URI configuration restrictions.",
					},
				},
				MarkdownDescription: "Restrictions that apply to all application objects in the tenant unless overridden by a specific appManagementPolicy.",
			},
			"service_principal_restrictions": schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{
					"password_credentials": schema.ListNestedAttribute{
						Optional: true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: getPasswordCredentialConfigurationAttributes(),
						},
						MarkdownDescription: "Collection of password credential restrictions for service principals.",
					},
					"key_credentials": schema.ListNestedAttribute{
						Optional: true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: getKeyCredentialConfigurationAttributes(),
						},
						MarkdownDescription: "Collection of key credential restrictions for service principals.",
					},
				},
				MarkdownDescription: "Restrictions that apply to all service principal objects in the tenant unless overridden by a specific appManagementPolicy.",
			},
			"timeouts": commonschema.ResourceTimeouts(ctx),
		},
	}
}

func getPasswordCredentialConfigurationAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"restriction_type": schema.StringAttribute{
			Required: true,
			Validators: []validator.String{
				stringvalidator.OneOf(
					"passwordAddition",
					"passwordLifetime",
					"symmetricKeyAddition",
					"symmetricKeyLifetime",
					"customPasswordAddition",
				),
			},
			MarkdownDescription: "A unique identifier key for passwordCredentialConfiguration. This value also represents the type of restriction being applied. The possible values are: passwordAddition, passwordLifetime, symmetricKeyAddition, symmetricKeyLifetime, customPasswordAddition, and unknownFutureValue. Each value of restrictionType can be used only once per policy.",
		},
		"state": schema.StringAttribute{
			Computed:            true,
			MarkdownDescription: "Indicates whether the restriction is evaluated. The possible values are: enabled, disabled, unknownFutureValue. If enabled, the restriction is evaluated. If disabled, the restriction isn't evaluated or enforced. Read-only.",
		},
		"restrict_for_apps_created_after_date_time": schema.StringAttribute{
			Required: true,
			Validators: []validator.String{
				stringvalidator.RegexMatches(
					regexp.MustCompile(constants.ISO8601DateTimeRegex),
					"must be a valid ISO 8601 datetime in the format YYYY-MM-DDTHH:MM:SSZ or with timezone offset",
				),
			},
			MarkdownDescription: "Specifies the date from which the policy restriction applies to newly created applications. For existing applications, the enforcement date can be retroactively applied. Format: RFC3339/DateTimeOffset (e.g., '2021-01-01T10:37:00Z').",
		},
		"max_lifetime": schema.StringAttribute{
			Optional: true,
			Validators: []validator.String{
				stringvalidator.RegexMatches(
					regexp.MustCompile(constants.ISO8601DurationRegex),
					"must be a valid ISO 8601 duration (e.g., P90D, P4DT12H30M5S)",
				),
			},
			MarkdownDescription: "String value that indicates the maximum lifetime for password expiration, defined as an ISO 8601 duration. For example, P4DT12H30M5S represents four days, 12 hours, 30 minutes, and five seconds. This property is required when restriction_type is set to passwordLifetime or symmetricKeyLifetime.",
		},
		"exclude_actors": schema.SingleNestedAttribute{
			Optional:            true,
			Attributes:          getActorExemptionsAttributes(),
			MarkdownDescription: "Collection of custom security attribute exemptions. If an actor user or service principal has the custom security attribute defined in this section, they're exempted from the restriction. This means that calls the user or service principal makes to create or update apps are exempt from this policy enforcement.",
		},
	}
}

func getKeyCredentialConfigurationAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"restriction_type": schema.StringAttribute{
			Required: true,
			Validators: []validator.String{
				stringvalidator.OneOf("asymmetricKeyLifetime", "trustedCertificateAuthority"),
			},
			MarkdownDescription: "A unique identifier key for keyCredentialConfiguration. This value also represents the type of restriction being applied. Possible values are asymmetricKeyLifetime, trustedCertificateAuthority, and unknownFutureValue. Each value of restrictionType can be used only once per policy.",
		},
		"state": schema.StringAttribute{
			Computed:            true,
			MarkdownDescription: "Indicates whether the restriction is evaluated. The possible values are: enabled, disabled, unknownFutureValue. If enabled, the restriction is evaluated. If disabled, the restriction isn't evaluated or enforced. Read-only.",
		},
		"restrict_for_apps_created_after_date_time": schema.StringAttribute{
			Required: true,
			Validators: []validator.String{
				stringvalidator.RegexMatches(
					regexp.MustCompile(constants.ISO8601DateTimeRegex),
					"must be a valid ISO 8601 datetime in the format YYYY-MM-DDTHH:MM:SSZ or with timezone offset",
				),
			},
			MarkdownDescription: "Specifies the date from which the policy restriction applies to newly created applications. For existing applications, the enforcement date can be retroactively applied. Format: RFC3339/DateTimeOffset (e.g., '2019-10-19T10:37:00Z').",
		},
		"max_lifetime": schema.StringAttribute{
			Optional: true,
			Validators: []validator.String{
				stringvalidator.RegexMatches(
					regexp.MustCompile(constants.ISO8601DurationRegex),
					"must be a valid ISO 8601 duration (e.g., P90D, P4DT12H30M5S)",
				),
			},
			MarkdownDescription: "String value that indicates the maximum lifetime for key expiration, defined as an ISO 8601 duration. For example, P4DT12H30M5S represents four days, 12 hours, 30 minutes, and five seconds. This property is required when restriction_type is set to asymmetricKeyLifetime.",
		},
		"certificate_based_application_configuration_ids": schema.ListAttribute{
			ElementType:         types.StringType,
			Optional:            true,
			MarkdownDescription: "Collection of GUIDs that represent certificateBasedApplicationConfiguration that is allowed as root and intermediate certificate authorities.",
		},
		"exclude_actors": schema.SingleNestedAttribute{
			Optional:            true,
			Attributes:          getActorExemptionsAttributes(),
			MarkdownDescription: "Exemptions based on custom security attributes.",
		},
	}
}

func getActorExemptionsAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"custom_security_attributes": schema.ListNestedAttribute{
			Optional: true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Required:            true,
						MarkdownDescription: "The custom security attribute ID (e.g., 'PolicyExemptions_AppManagementExemption').",
					},
					"operator": schema.StringAttribute{
						Required: true,
						Validators: []validator.String{
							stringvalidator.OneOf("equals"),
						},
						MarkdownDescription: "The comparison operator. Currently only 'equals' is supported.",
					},
					"value": schema.StringAttribute{
						Required:            true,
						MarkdownDescription: "The value to compare against the custom security attribute.",
					},
				},
			},
			MarkdownDescription: "Collection of custom security attribute exemptions. If an actor user or service principal has the custom security attribute defined in this section, they're exempted from the restriction.",
		},
	}
}
