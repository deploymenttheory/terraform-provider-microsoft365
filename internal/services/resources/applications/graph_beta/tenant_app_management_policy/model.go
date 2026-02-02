// REF: https://learn.microsoft.com/en-us/graph/api/resources/tenantappmanagementpolicy?view=graph-rest-beta
// REF: https://learn.microsoft.com/en-us/graph/api/tenantappmanagementpolicy-get?view=graph-rest-beta
// REF: https://learn.microsoft.com/en-us/graph/api/tenantappmanagementpolicy-update?view=graph-rest-beta
package graphBetaTenantAppManagementPolicy

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type TenantAppManagementPolicyResourceModel struct {
	ID                           types.String                                     `tfsdk:"id"`
	DisplayName                  types.String                                     `tfsdk:"display_name"`
	Description                  types.String                                     `tfsdk:"description"`
	IsEnabled                    types.Bool                                       `tfsdk:"is_enabled"`
	ApplicationRestrictions      *AppManagementApplicationConfigurationModel      `tfsdk:"application_restrictions"`
	ServicePrincipalRestrictions *AppManagementServicePrincipalConfigurationModel `tfsdk:"service_principal_restrictions"`
	RestoreToDefaultUponDelete   types.Bool                                       `tfsdk:"restore_to_default_upon_delete"`
	Timeouts                     timeouts.Value                                   `tfsdk:"timeouts"`
}

type AppManagementApplicationConfigurationModel struct {
	PasswordCredentials []PasswordCredentialConfigurationModel `tfsdk:"password_credentials"`
	KeyCredentials      []KeyCredentialConfigurationModel      `tfsdk:"key_credentials"`
	IdentifierUris      *IdentifierUriConfigurationModel       `tfsdk:"identifier_uris"`
}

type AppManagementServicePrincipalConfigurationModel struct {
	PasswordCredentials []PasswordCredentialConfigurationModel `tfsdk:"password_credentials"`
	KeyCredentials      []KeyCredentialConfigurationModel      `tfsdk:"key_credentials"`
}

type PasswordCredentialConfigurationModel struct {
	RestrictionType                     types.String                             `tfsdk:"restriction_type"`
	State                               types.String                             `tfsdk:"state"`
	RestrictForAppsCreatedAfterDateTime types.String                             `tfsdk:"restrict_for_apps_created_after_date_time"`
	MaxLifetime                         types.String                             `tfsdk:"max_lifetime"`
	ExcludeActors                       *AppManagementPolicyActorExemptionsModel `tfsdk:"exclude_actors"`
}

type KeyCredentialConfigurationModel struct {
	RestrictionType                             types.String                             `tfsdk:"restriction_type"`
	State                                       types.String                             `tfsdk:"state"`
	RestrictForAppsCreatedAfterDateTime         types.String                             `tfsdk:"restrict_for_apps_created_after_date_time"`
	MaxLifetime                                 types.String                             `tfsdk:"max_lifetime"`
	CertificateBasedApplicationConfigurationIds []types.String                           `tfsdk:"certificate_based_application_configuration_ids"`
	ExcludeActors                               *AppManagementPolicyActorExemptionsModel `tfsdk:"exclude_actors"`
}

type IdentifierUriConfigurationModel struct {
	NonDefaultUriAddition *IdentifierUriRestrictionModel `tfsdk:"non_default_uri_addition"`
}

type IdentifierUriRestrictionModel struct {
	RestrictForAppsCreatedAfterDateTime types.String                             `tfsdk:"restrict_for_apps_created_after_date_time"`
	ExcludeAppsReceivingV2Tokens        types.Bool                               `tfsdk:"exclude_apps_receiving_v2_tokens"`
	ExcludeSaml                         types.Bool                               `tfsdk:"exclude_saml"`
	ExcludeActors                       *AppManagementPolicyActorExemptionsModel `tfsdk:"exclude_actors"`
}

type AppManagementPolicyActorExemptionsModel struct {
	CustomSecurityAttributes []CustomSecurityAttributeExemptionModel `tfsdk:"custom_security_attributes"`
}

type CustomSecurityAttributeExemptionModel struct {
	ID       types.String `tfsdk:"id"`
	Operator types.String `tfsdk:"operator"`
	Value    types.String `tfsdk:"value"`
}
