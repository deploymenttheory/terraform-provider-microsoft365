package graphBetaTenantAppManagementPolicy

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteStateToTerraform maps the remote state from Microsoft Graph API to the Terraform state
func MapRemoteStateToTerraform(ctx context.Context, data *TenantAppManagementPolicyResourceModel, policy graphmodels.TenantAppManagementPolicyable) {
	tflog.Debug(ctx, fmt.Sprintf("Mapping %s remote state to Terraform state", ResourceName))

	if policy == nil {
		tflog.Debug(ctx, "Policy is nil")
		return
	}

	data.ID = convert.GraphToFrameworkString(policy.GetId())
	data.DisplayName = convert.GraphToFrameworkString(policy.GetDisplayName())
	data.Description = convert.GraphToFrameworkString(policy.GetDescription())

	if isEnabled := policy.GetIsEnabled(); isEnabled != nil {
		data.IsEnabled = types.BoolValue(*isEnabled)
	}

	// Map application restrictions
	if appRestrictions := policy.GetApplicationRestrictions(); appRestrictions != nil {
		data.ApplicationRestrictions = &AppManagementApplicationConfigurationModel{}

		// Map password credentials
		if passwordCreds := appRestrictions.GetPasswordCredentials(); passwordCreds != nil {
			data.ApplicationRestrictions.PasswordCredentials = make([]PasswordCredentialConfigurationModel, 0, len(passwordCreds))
			for _, pc := range passwordCreds {
				data.ApplicationRestrictions.PasswordCredentials = append(
					data.ApplicationRestrictions.PasswordCredentials,
					mapPasswordCredentialConfiguration(pc),
				)
			}
		}

		// Map key credentials
		if keyCreds := appRestrictions.GetKeyCredentials(); keyCreds != nil {
			data.ApplicationRestrictions.KeyCredentials = make([]KeyCredentialConfigurationModel, 0, len(keyCreds))
			for _, kc := range keyCreds {
				data.ApplicationRestrictions.KeyCredentials = append(
					data.ApplicationRestrictions.KeyCredentials,
					mapKeyCredentialConfiguration(kc),
				)
			}
		}

		// Map identifier URIs
		if identifierUris := appRestrictions.GetIdentifierUris(); identifierUris != nil {
			data.ApplicationRestrictions.IdentifierUris = &IdentifierUriConfigurationModel{}

			if nonDefaultUriAddition := identifierUris.GetNonDefaultUriAddition(); nonDefaultUriAddition != nil {
				restriction := &IdentifierUriRestrictionModel{}

				restriction.RestrictForAppsCreatedAfterDateTime = convert.GraphToFrameworkTime(nonDefaultUriAddition.GetRestrictForAppsCreatedAfterDateTime())

				if exclude := nonDefaultUriAddition.GetExcludeAppsReceivingV2Tokens(); exclude != nil {
					restriction.ExcludeAppsReceivingV2Tokens = types.BoolValue(*exclude)
				}

				if exclude := nonDefaultUriAddition.GetExcludeSaml(); exclude != nil {
					restriction.ExcludeSaml = types.BoolValue(*exclude)
				}

				if exemptions := nonDefaultUriAddition.GetExcludeActors(); exemptions != nil {
					restriction.ExcludeActors = mapActorExemptions(exemptions)
				}

				data.ApplicationRestrictions.IdentifierUris.NonDefaultUriAddition = restriction
			}
		}
	}

	// Map service principal restrictions
	if spRestrictions := policy.GetServicePrincipalRestrictions(); spRestrictions != nil {
		data.ServicePrincipalRestrictions = &AppManagementServicePrincipalConfigurationModel{}

		// Map password credentials
		if passwordCreds := spRestrictions.GetPasswordCredentials(); passwordCreds != nil {
			data.ServicePrincipalRestrictions.PasswordCredentials = make([]PasswordCredentialConfigurationModel, 0, len(passwordCreds))
			for _, pc := range passwordCreds {
				data.ServicePrincipalRestrictions.PasswordCredentials = append(
					data.ServicePrincipalRestrictions.PasswordCredentials,
					mapPasswordCredentialConfiguration(pc),
				)
			}
		}

		// Map key credentials
		if keyCreds := spRestrictions.GetKeyCredentials(); keyCreds != nil {
			data.ServicePrincipalRestrictions.KeyCredentials = make([]KeyCredentialConfigurationModel, 0, len(keyCreds))
			for _, kc := range keyCreds {
				data.ServicePrincipalRestrictions.KeyCredentials = append(
					data.ServicePrincipalRestrictions.KeyCredentials,
					mapKeyCredentialConfiguration(kc),
				)
			}
		}
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished mapping %s remote state to Terraform state", ResourceName))
}

func mapPasswordCredentialConfiguration(config graphmodels.PasswordCredentialConfigurationable) PasswordCredentialConfigurationModel {
	model := PasswordCredentialConfigurationModel{}

	if restrictionType := config.GetRestrictionType(); restrictionType != nil {
		model.RestrictionType = types.StringValue(restrictionType.String())
	}

	if state := config.GetState(); state != nil {
		model.State = types.StringValue(state.String())
	}

	model.RestrictForAppsCreatedAfterDateTime = convert.GraphToFrameworkTime(config.GetRestrictForAppsCreatedAfterDateTime())
	model.MaxLifetime = convert.GraphToFrameworkISODuration(config.GetMaxLifetime())

	if exemptions := config.GetExcludeActors(); exemptions != nil {
		model.ExcludeActors = mapActorExemptions(exemptions)
	}

	return model
}

func mapKeyCredentialConfiguration(config graphmodels.KeyCredentialConfigurationable) KeyCredentialConfigurationModel {
	model := KeyCredentialConfigurationModel{}

	if restrictionType := config.GetRestrictionType(); restrictionType != nil {
		model.RestrictionType = types.StringValue(restrictionType.String())
	}

	if state := config.GetState(); state != nil {
		model.State = types.StringValue(state.String())
	}

	model.RestrictForAppsCreatedAfterDateTime = convert.GraphToFrameworkTime(config.GetRestrictForAppsCreatedAfterDateTime())
	model.MaxLifetime = convert.GraphToFrameworkISODuration(config.GetMaxLifetime())

	if certIds := config.GetCertificateBasedApplicationConfigurationIds(); certIds != nil {
		model.CertificateBasedApplicationConfigurationIds = make([]types.String, 0, len(certIds))
		for _, id := range certIds {
			model.CertificateBasedApplicationConfigurationIds = append(
				model.CertificateBasedApplicationConfigurationIds,
				types.StringValue(id),
			)
		}
	}

	if exemptions := config.GetExcludeActors(); exemptions != nil {
		model.ExcludeActors = mapActorExemptions(exemptions)
	}

	return model
}

func mapActorExemptions(exemptions graphmodels.AppManagementPolicyActorExemptionsable) *AppManagementPolicyActorExemptionsModel {
	model := &AppManagementPolicyActorExemptionsModel{}

	if csaExemptions := exemptions.GetCustomSecurityAttributes(); csaExemptions != nil {
		model.CustomSecurityAttributes = make([]CustomSecurityAttributeExemptionModel, 0, len(csaExemptions))
		for _, csa := range csaExemptions {
			exemption := CustomSecurityAttributeExemptionModel{}

			exemption.ID = convert.GraphToFrameworkString(csa.GetId())

			if operator := csa.GetOperator(); operator != nil {
				exemption.Operator = types.StringValue(operator.String())
			}

			if csaStringValue, ok := csa.(graphmodels.CustomSecurityAttributeStringValueExemptionable); ok {
				exemption.Value = convert.GraphToFrameworkString(csaStringValue.GetValue())
			}

			model.CustomSecurityAttributes = append(model.CustomSecurityAttributes, exemption)
		}
	}

	return model
}
