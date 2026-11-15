package graphBetaIdentityAndAccessAuthenticationMethodsPolicy

import (
	"context"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/identityschema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
)

const (
	ResourceName  = "microsoft365_graph_beta_identity_and_access_authentication_methods_policy"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180

	// singletonID is the fixed identifier for the tenant-wide authentication methods policy.
	singletonID = "authenticationMethodsPolicy"

	// allUsersTargetID is the special target ID that represents every user in the tenant.
	allUsersTargetID = "all_users"
)

var (
	// Basic resource interface (CRUD operations)
	_ resource.Resource = &AuthenticationMethodsPolicyResource{}

	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &AuthenticationMethodsPolicyResource{}

	// Enables import functionality
	_ resource.ResourceWithImportState = &AuthenticationMethodsPolicyResource{}

	// Enables identity schema for list resource support
	_ resource.ResourceWithIdentity = &AuthenticationMethodsPolicyResource{}
)

// https://learn.microsoft.com/en-us/graph/api/authenticationmethodspolicy-get?view=graph-rest-beta#permissions
// https://learn.microsoft.com/en-us/graph/api/authenticationmethodspolicy-update?view=graph-rest-beta#permissions
// https://learn.microsoft.com/en-us/graph/api/authenticationmethodconfiguration-update?view=graph-rest-beta#permissions
func NewAuthenticationMethodsPolicyResource() resource.Resource {
	return &AuthenticationMethodsPolicyResource{
		ReadPermissions: []string{
			"Policy.Read.AuthenticationMethod",
		},
		WritePermissions: []string{
			"Policy.ReadWrite.AuthenticationMethod",
		},
	}
}

type AuthenticationMethodsPolicyResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
}

// Metadata returns the resource type name.
func (r *AuthenticationMethodsPolicyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = ResourceName
}

// Configure sets the client for the resource.
func (r *AuthenticationMethodsPolicyResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, ResourceName)
}

// ImportState imports the resource state using the fixed singleton identifier.
func (r *AuthenticationMethodsPolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	if req.ID != singletonID {
		resp.Diagnostics.AddError(
			"Invalid import ID",
			"Import this singleton using the literal ID authenticationMethodsPolicy.",
		)
		return
	}
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// IdentitySchema defines the identity schema for this resource, used by list operations to uniquely identify instances
func (r *AuthenticationMethodsPolicyResource) IdentitySchema(ctx context.Context, req resource.IdentitySchemaRequest, resp *resource.IdentitySchemaResponse) {
	resp.IdentitySchema = identityschema.Schema{
		Attributes: map[string]identityschema.Attribute{
			"id": identityschema.StringAttribute{
				RequiredForImport: true,
			},
		},
	}
}

// Schema defines the schema for the resource.
func (r *AuthenticationMethodsPolicyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages the tenant-wide Microsoft Entra authentication methods policy using the Microsoft Graph beta `/policies/authenticationMethodsPolicy` endpoint " +
			"and its `authenticationMethodConfigurations`. The policy controls which authentication methods users can register and use, the registration campaign, " +
			"system-preferred multifactor authentication and suspicious activity reporting.\n\n" +
			"This is a **singleton resource**: one policy exists per tenant. `create` adopts the existing policy and applies the configuration; `update` applies changes. " +
			"Every setting and every method block is optional. A setting or method block that is not configured keeps its current tenant value and is never written. " +
			"Each changed method configuration is updated with its own PATCH request, so method configurations that are not in the configuration are left untouched. " +
			"On `destroy`, Terraform only removes the resource from state and leaves the tenant policy as it is. Manage only one instance per tenant across all Terraform states.\n\n" +
			"This version manages the `email`, `temporary_access_pass`, `sms`, `software_oath` and `hardware_oath` method configurations. " +
			"The `fido2`, `microsoftAuthenticator`, `voice`, `x509Certificate`, `qrCodePin`, `password`, `verifiableCredentials`, `federatedIdentityCredential` and " +
			"external method configurations are not managed by this resource and are never changed by it.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The fixed singleton identifier `authenticationMethodsPolicy`.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"display_name": schema.StringAttribute{
				MarkdownDescription: "The display name of the policy. Read-only.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The description of the policy. Read-only.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"policy_migration_state": schema.StringAttribute{
				MarkdownDescription: "The state of migration of the authentication methods policy from the legacy multifactor authentication and self-service password reset (SSPR) policies. " +
					"Possible values are: `preMigration`, `migrationInProgress`, `migrationComplete`. When not configured, the current tenant value is kept.",
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf("preMigration", "migrationInProgress", "migrationComplete"),
				},
			},
			"registration_enforcement": schema.SingleNestedAttribute{
				MarkdownDescription: "Enforces registration at sign-in time. When not configured, the current tenant value is kept.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
				Attributes: map[string]schema.Attribute{
					"authentication_methods_registration_campaign": schema.SingleNestedAttribute{
						MarkdownDescription: "Run campaigns to remind users to set up targeted authentication methods.",
						Required:            true,
						Attributes: map[string]schema.Attribute{
							"snooze_duration_in_days": schema.Int32Attribute{
								MarkdownDescription: "The number of days a user can postpone registration. Accepts 0 to 14. When not configured, the current tenant value is kept.",
								Optional:            true,
								Computed:            true,
								PlanModifiers: []planmodifier.Int32{
									int32planmodifier.UseStateForUnknown(),
								},
								Validators: []validator.Int32{
									int32validator.Between(0, 14),
								},
							},
							"enforce_registration_after_allowed_snoozes": schema.BoolAttribute{
								MarkdownDescription: "Whether users must register after they use up their allowed snoozes. When not configured, the current tenant value is kept.",
								Optional:            true,
								Computed:            true,
								PlanModifiers: []planmodifier.Bool{
									boolplanmodifier.UseStateForUnknown(),
								},
							},
							"state":           advancedConfigStateAttribute("The state of the registration campaign."),
							"include_targets": registrationCampaignIncludeTargetsAttribute(),
							"exclude_targets": targetSetAttribute("Users and groups excluded from the registration campaign."),
						},
					},
				},
			},
			"report_suspicious_activity_settings": schema.SingleNestedAttribute{
				MarkdownDescription: "Settings that let users report suspicious multifactor authentication prompts. When not configured, the current tenant value is kept.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
				Attributes: map[string]schema.Attribute{
					"state": advancedConfigStateAttribute("The state of suspicious activity reporting."),
					"voice_reporting_code": schema.Int32Attribute{
						MarkdownDescription: "The keypad code a user enters during a voice call to report suspicious activity. Accepts 0 to 9; `0` uses the default. " +
							"When not configured, the current tenant value is kept.",
						Optional: true,
						Computed: true,
						PlanModifiers: []planmodifier.Int32{
							int32planmodifier.UseStateForUnknown(),
						},
						Validators: []validator.Int32{
							int32validator.Between(0, 9),
						},
					},
					"include_target": schema.SingleNestedAttribute{
						MarkdownDescription: "The user or group the setting applies to. When not configured, the current tenant value is kept.",
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.Object{
							objectplanmodifier.UseStateForUnknown(),
						},
						Attributes: targetAttributes(),
					},
				},
			},
			"system_credential_preferences": schema.SingleNestedAttribute{
				MarkdownDescription: "System-preferred multifactor authentication settings, which prompt users for the most secure method they registered. " +
					"When not configured, the current tenant value is kept.",
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
				Attributes: map[string]schema.Attribute{
					"state":           advancedConfigStateAttribute("The state of system-preferred multifactor authentication."),
					"include_targets": targetSetAttribute("Users and groups included in system-preferred multifactor authentication."),
					"exclude_targets": targetSetAttribute("Users and groups excluded from system-preferred multifactor authentication."),
				},
			},
			"email": methodAttribute("Email one-time passcode authentication method configuration (`Email`).", map[string]schema.Attribute{
				"allow_external_id_to_use_email_otp": schema.StringAttribute{
					MarkdownDescription: "Whether external users can use email one-time passcodes for sign-in. Possible values are: `default`, `enabled`, `disabled`. " +
						"When not configured, the current tenant value is kept.",
					Optional: true,
					Computed: true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.UseStateForUnknown(),
					},
					Validators: []validator.String{
						stringvalidator.OneOf("default", "enabled", "disabled"),
					},
				},
				"include_targets": methodTargetSetAttribute("Users and groups that can use email one-time passcodes."),
			}),
			"temporary_access_pass": methodAttribute("Temporary Access Pass authentication method configuration (`TemporaryAccessPass`).", map[string]schema.Attribute{
				"default_lifetime_in_minutes": lifetimeAttribute("The default lifetime in minutes of a new pass. Accepts 10 to 43200."),
				"default_length": schema.Int32Attribute{
					MarkdownDescription: "The default length in characters of a new pass. Accepts 8 to 48. When not configured, the current tenant value is kept.",
					Optional:            true,
					Computed:            true,
					PlanModifiers: []planmodifier.Int32{
						int32planmodifier.UseStateForUnknown(),
					},
					Validators: []validator.Int32{
						int32validator.Between(8, 48),
					},
				},
				"minimum_lifetime_in_minutes": lifetimeAttribute("The minimum lifetime in minutes of any pass. Accepts 10 to 43200."),
				"maximum_lifetime_in_minutes": lifetimeAttribute("The maximum lifetime in minutes of any pass. Accepts 10 to 43200."),
				"is_usable_once": schema.BoolAttribute{
					MarkdownDescription: "Whether passes are limited to a one-time use. When not configured, the current tenant value is kept.",
					Optional:            true,
					Computed:            true,
					PlanModifiers: []planmodifier.Bool{
						boolplanmodifier.UseStateForUnknown(),
					},
				},
				"include_targets": methodTargetSetAttribute("Users and groups that can use a Temporary Access Pass."),
			}),
			"sms": methodAttribute("Text message authentication method configuration (`Sms`).", map[string]schema.Attribute{
				"include_targets": smsMethodTargetSetAttribute(),
			}),
			"software_oath": methodAttribute("Third-party software OATH token authentication method configuration (`SoftwareOath`).", map[string]schema.Attribute{
				"include_targets": methodTargetSetAttribute("Users and groups that can use software OATH tokens."),
			}),
			"hardware_oath": methodAttribute("Hardware OATH token authentication method configuration (`HardwareOath`).", map[string]schema.Attribute{
				"include_targets": methodTargetSetAttribute("Users and groups that can use hardware OATH tokens."),
			}),
			"timeouts": commonschema.ResourceTimeouts(ctx),
		},
	}
}

// methodAttribute returns the schema for an authentication method configuration block. Every block has a
// required state and an exclude_targets set; specific holds the method's include_targets and other settings.
func methodAttribute(description string, specific map[string]schema.Attribute) schema.SingleNestedAttribute {
	attributes := map[string]schema.Attribute{
		"state": schema.StringAttribute{
			MarkdownDescription: "Whether the method is enabled in the tenant. Possible values are: `enabled`, `disabled`.",
			Required:            true,
			Validators: []validator.String{
				stringvalidator.OneOf("enabled", "disabled"),
			},
		},
		"exclude_targets": targetSetAttribute("Users and groups excluded from the method."),
	}
	for name, attribute := range specific {
		attributes[name] = attribute
	}
	return schema.SingleNestedAttribute{
		MarkdownDescription: description + " When the block is not configured, the method configuration is read from the tenant and is never written.",
		Optional:            true,
		Computed:            true,
		PlanModifiers: []planmodifier.Object{
			objectplanmodifier.UseStateForUnknown(),
		},
		Attributes: attributes,
	}
}

// advancedConfigStateAttribute returns the schema for an advancedConfigState value.
func advancedConfigStateAttribute(description string) schema.StringAttribute {
	return schema.StringAttribute{
		MarkdownDescription: description + " Possible values are: `default`, `enabled`, `disabled`. When not configured, the current tenant value is kept.",
		Optional:            true,
		Computed:            true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
		Validators: []validator.String{
			stringvalidator.OneOf("default", "enabled", "disabled"),
		},
	}
}

// lifetimeAttribute returns the schema for a Temporary Access Pass lifetime in minutes.
func lifetimeAttribute(description string) schema.Int32Attribute {
	return schema.Int32Attribute{
		MarkdownDescription: description + " When not configured, the current tenant value is kept.",
		Optional:            true,
		Computed:            true,
		PlanModifiers: []planmodifier.Int32{
			int32planmodifier.UseStateForUnknown(),
		},
		Validators: []validator.Int32{
			int32validator.Between(10, 43200),
		},
	}
}

// targetIDAttribute returns the schema for a target ID, which is a user or group object ID or all_users.
func targetIDAttribute() schema.StringAttribute {
	return schema.StringAttribute{
		MarkdownDescription: "The object ID of the user or group, or `all_users` for every user in the tenant.",
		Required:            true,
		Validators: []validator.String{
			stringvalidator.Any(
				stringvalidator.RegexMatches(regexp.MustCompile(constants.GuidRegex), "must be a valid GUID in the format '00000000-0000-0000-0000-000000000000'"),
				stringvalidator.OneOf(allUsersTargetID),
			),
		},
	}
}

// targetTypeAttribute returns the schema for an authenticationMethodTargetType value.
func targetTypeAttribute() schema.StringAttribute {
	return schema.StringAttribute{
		MarkdownDescription: "The type of the target. Possible values are: `user`, `group`.",
		Required:            true,
		Validators: []validator.String{
			stringvalidator.OneOf("user", "group"),
		},
	}
}

// targetAttributes returns the attributes of an includeTarget or excludeTarget.
func targetAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id":          targetIDAttribute(),
		"target_type": targetTypeAttribute(),
	}
}

// isRegistrationRequiredAttribute returns the schema for authenticationMethodTarget.isRegistrationRequired.
func isRegistrationRequiredAttribute() schema.BoolAttribute {
	return schema.BoolAttribute{
		MarkdownDescription: "Whether users in the target must register the method. Defaults to `false`.",
		Optional:            true,
		Computed:            true,
		Default:             booldefault.StaticBool(false),
	}
}

// targetSetAttribute returns the schema for a set of includeTarget or excludeTarget values.
func targetSetAttribute(description string) schema.SetNestedAttribute {
	return schema.SetNestedAttribute{
		MarkdownDescription: description + " When not configured, the current tenant value is kept; set `[]` to clear it.",
		Optional:            true,
		Computed:            true,
		PlanModifiers: []planmodifier.Set{
			setplanmodifier.UseStateForUnknown(),
		},
		NestedObject: schema.NestedAttributeObject{
			Attributes: targetAttributes(),
		},
	}
}

// methodTargetSetAttribute returns the schema for a set of authenticationMethodTarget values.
func methodTargetSetAttribute(description string) schema.SetNestedAttribute {
	return schema.SetNestedAttribute{
		MarkdownDescription: description + " When not configured, the current tenant value is kept; set `[]` to clear it.",
		Optional:            true,
		Computed:            true,
		PlanModifiers: []planmodifier.Set{
			setplanmodifier.UseStateForUnknown(),
		},
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"id":                       targetIDAttribute(),
				"target_type":              targetTypeAttribute(),
				"is_registration_required": isRegistrationRequiredAttribute(),
			},
		},
	}
}

// smsMethodTargetSetAttribute returns the schema for a set of smsAuthenticationMethodTarget values.
func smsMethodTargetSetAttribute() schema.SetNestedAttribute {
	return schema.SetNestedAttribute{
		MarkdownDescription: "Users and groups that can use text messages. When not configured, the current tenant value is kept; set `[]` to clear it.",
		Optional:            true,
		Computed:            true,
		PlanModifiers: []planmodifier.Set{
			setplanmodifier.UseStateForUnknown(),
		},
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"id":                       targetIDAttribute(),
				"target_type":              targetTypeAttribute(),
				"is_registration_required": isRegistrationRequiredAttribute(),
				"is_usable_for_sign_in": schema.BoolAttribute{
					// Required rather than defaulted: a second default in the same set element is resolved
					// against the element value after the first default is applied, which drops the configured value.
					MarkdownDescription: "Whether users in the target can use text messages as a first factor (SMS sign-in).",
					Required:            true,
				},
			},
		},
	}
}

// registrationCampaignIncludeTargetsAttribute returns the schema for the registration campaign include targets.
func registrationCampaignIncludeTargetsAttribute() schema.SetNestedAttribute {
	return schema.SetNestedAttribute{
		MarkdownDescription: "Users and groups the registration campaign targets. When not configured, the current tenant value is kept; set `[]` to clear it.",
		Optional:            true,
		Computed:            true,
		PlanModifiers: []planmodifier.Set{
			setplanmodifier.UseStateForUnknown(),
		},
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"id":          targetIDAttribute(),
				"target_type": targetTypeAttribute(),
				"targeted_authentication_method": schema.StringAttribute{
					MarkdownDescription: "The authentication method users are prompted to register, for example `microsoftAuthenticator`.",
					Required:            true,
				},
			},
		},
	}
}
