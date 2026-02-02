package graphBetaTenantAppManagementPolicy

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// constructResource constructs a TenantAppManagementPolicy from the Terraform model
func constructResource(ctx context.Context, data *TenantAppManagementPolicyResourceModel) (graphmodels.TenantAppManagementPolicyable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))

	requestBody := graphmodels.NewTenantAppManagementPolicy()

	convert.FrameworkToGraphString(data.DisplayName, requestBody.SetDisplayName)
	convert.FrameworkToGraphString(data.Description, requestBody.SetDescription)
	convert.FrameworkToGraphBool(data.IsEnabled, requestBody.SetIsEnabled)

	if data.ApplicationRestrictions != nil {
		appRestrictions := graphmodels.NewAppManagementApplicationConfiguration()

		if len(data.ApplicationRestrictions.PasswordCredentials) > 0 {
			passwordCreds := make([]graphmodels.PasswordCredentialConfigurationable, 0, len(data.ApplicationRestrictions.PasswordCredentials))
			for _, pc := range data.ApplicationRestrictions.PasswordCredentials {
				cred, err := constructPasswordCredentialConfiguration(ctx, &pc)
				if err != nil {
					return nil, err
				}
				passwordCreds = append(passwordCreds, cred)
			}
			appRestrictions.SetPasswordCredentials(passwordCreds)
		}

		if len(data.ApplicationRestrictions.KeyCredentials) > 0 {
			keyCreds := make([]graphmodels.KeyCredentialConfigurationable, 0, len(data.ApplicationRestrictions.KeyCredentials))
			for _, kc := range data.ApplicationRestrictions.KeyCredentials {
				cred, err := constructKeyCredentialConfiguration(ctx, &kc)
				if err != nil {
					return nil, err
				}
				keyCreds = append(keyCreds, cred)
			}
			appRestrictions.SetKeyCredentials(keyCreds)
		}

		if data.ApplicationRestrictions.IdentifierUris != nil {
			identifierUris := graphmodels.NewIdentifierUriConfiguration()
			if data.ApplicationRestrictions.IdentifierUris.NonDefaultUriAddition != nil {
				restriction := graphmodels.NewIdentifierUriRestriction()

				if err := convert.FrameworkToGraphTime(
					data.ApplicationRestrictions.IdentifierUris.NonDefaultUriAddition.RestrictForAppsCreatedAfterDateTime,
					restriction.SetRestrictForAppsCreatedAfterDateTime,
				); err != nil {
					return nil, fmt.Errorf("error parsing restrict_for_apps_created_after_date_time: %v", err)
				}

				convert.FrameworkToGraphBool(
					data.ApplicationRestrictions.IdentifierUris.NonDefaultUriAddition.ExcludeAppsReceivingV2Tokens,
					restriction.SetExcludeAppsReceivingV2Tokens,
				)

				convert.FrameworkToGraphBool(
					data.ApplicationRestrictions.IdentifierUris.NonDefaultUriAddition.ExcludeSaml,
					restriction.SetExcludeSaml,
				)

				if data.ApplicationRestrictions.IdentifierUris.NonDefaultUriAddition.ExcludeActors != nil {
					exemptions, err := constructActorExemptions(ctx, data.ApplicationRestrictions.IdentifierUris.NonDefaultUriAddition.ExcludeActors)
					if err != nil {
						return nil, err
					}
					restriction.SetExcludeActors(exemptions)
				}

				identifierUris.SetNonDefaultUriAddition(restriction)
			}
			appRestrictions.SetIdentifierUris(identifierUris)
		}

		requestBody.SetApplicationRestrictions(appRestrictions)
	}

	if data.ServicePrincipalRestrictions != nil {
		spRestrictions := graphmodels.NewAppManagementServicePrincipalConfiguration()

		if len(data.ServicePrincipalRestrictions.PasswordCredentials) > 0 {
			passwordCreds := make([]graphmodels.PasswordCredentialConfigurationable, 0, len(data.ServicePrincipalRestrictions.PasswordCredentials))
			for _, pc := range data.ServicePrincipalRestrictions.PasswordCredentials {
				cred, err := constructPasswordCredentialConfiguration(ctx, &pc)
				if err != nil {
					return nil, err
				}
				passwordCreds = append(passwordCreds, cred)
			}
			spRestrictions.SetPasswordCredentials(passwordCreds)
		}

		if len(data.ServicePrincipalRestrictions.KeyCredentials) > 0 {
			keyCreds := make([]graphmodels.KeyCredentialConfigurationable, 0, len(data.ServicePrincipalRestrictions.KeyCredentials))
			for _, kc := range data.ServicePrincipalRestrictions.KeyCredentials {
				cred, err := constructKeyCredentialConfiguration(ctx, &kc)
				if err != nil {
					return nil, err
				}
				keyCreds = append(keyCreds, cred)
			}
			spRestrictions.SetKeyCredentials(keyCreds)
		}

		requestBody.SetServicePrincipalRestrictions(spRestrictions)
	}

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s resource", ResourceName))

	return requestBody, nil
}

// constructPasswordCredentialConfiguration constructs a password credential configuration
func constructPasswordCredentialConfiguration(ctx context.Context, data *PasswordCredentialConfigurationModel) (graphmodels.PasswordCredentialConfigurationable, error) {
	config := graphmodels.NewPasswordCredentialConfiguration()

	if !data.RestrictionType.IsNull() && !data.RestrictionType.IsUnknown() {
		restrictionType := data.RestrictionType.ValueString()
		restrictionTypeEnum, err := graphmodels.ParseAppCredentialRestrictionType(restrictionType)
		if err != nil {
			return nil, fmt.Errorf("invalid restriction type: %s", restrictionType)
		}
		if typedEnum, ok := restrictionTypeEnum.(*graphmodels.AppCredentialRestrictionType); ok {
			config.SetRestrictionType(typedEnum)
		}
	}

	if err := convert.FrameworkToGraphTime(data.RestrictForAppsCreatedAfterDateTime, config.SetRestrictForAppsCreatedAfterDateTime); err != nil {
		return nil, fmt.Errorf("error parsing restrict_for_apps_created_after_date_time: %v", err)
	}

	if err := convert.FrameworkToGraphISODuration(data.MaxLifetime, config.SetMaxLifetime); err != nil {
		return nil, fmt.Errorf("error parsing max_lifetime: %v", err)
	}

	if data.ExcludeActors != nil {
		exemptions, err := constructActorExemptions(ctx, data.ExcludeActors)
		if err != nil {
			return nil, err
		}
		config.SetExcludeActors(exemptions)
	}

	return config, nil
}

// constructKeyCredentialConfiguration constructs a key credential configuration
func constructKeyCredentialConfiguration(ctx context.Context, data *KeyCredentialConfigurationModel) (graphmodels.KeyCredentialConfigurationable, error) {
	config := graphmodels.NewKeyCredentialConfiguration()

	if !data.RestrictionType.IsNull() && !data.RestrictionType.IsUnknown() {
		restrictionType := data.RestrictionType.ValueString()
		restrictionTypeEnum, err := graphmodels.ParseAppKeyCredentialRestrictionType(restrictionType)
		if err != nil {
			return nil, fmt.Errorf("invalid restriction type: %s", restrictionType)
		}
		if typedEnum, ok := restrictionTypeEnum.(*graphmodels.AppKeyCredentialRestrictionType); ok {
			config.SetRestrictionType(typedEnum)
		}
	}

	if err := convert.FrameworkToGraphTime(data.RestrictForAppsCreatedAfterDateTime, config.SetRestrictForAppsCreatedAfterDateTime); err != nil {
		return nil, fmt.Errorf("error parsing restrict_for_apps_created_after_date_time: %v", err)
	}

	if err := convert.FrameworkToGraphISODuration(data.MaxLifetime, config.SetMaxLifetime); err != nil {
		return nil, fmt.Errorf("error parsing max_lifetime: %v", err)
	}

	if len(data.CertificateBasedApplicationConfigurationIds) > 0 {
		certIds := make([]string, 0, len(data.CertificateBasedApplicationConfigurationIds))
		for _, id := range data.CertificateBasedApplicationConfigurationIds {
			certIds = append(certIds, id.ValueString())
		}
		config.SetCertificateBasedApplicationConfigurationIds(certIds)
	}

	if data.ExcludeActors != nil {
		exemptions, err := constructActorExemptions(ctx, data.ExcludeActors)
		if err != nil {
			return nil, err
		}
		config.SetExcludeActors(exemptions)
	}

	return config, nil
}

// constructActorExemptions constructs actor exemptions
func constructActorExemptions(ctx context.Context, data *AppManagementPolicyActorExemptionsModel) (graphmodels.AppManagementPolicyActorExemptionsable, error) {
	exemptions := graphmodels.NewAppManagementPolicyActorExemptions()

	if len(data.CustomSecurityAttributes) > 0 {
		csaExemptions := make([]graphmodels.CustomSecurityAttributeExemptionable, 0, len(data.CustomSecurityAttributes))
		for _, csa := range data.CustomSecurityAttributes {
			exemption := graphmodels.NewCustomSecurityAttributeStringValueExemption()

			convert.FrameworkToGraphString(csa.ID, exemption.SetId)

			if !csa.Operator.IsNull() && !csa.Operator.IsUnknown() {
				operator := csa.Operator.ValueString()
				operatorEnum, err := graphmodels.ParseCustomSecurityAttributeComparisonOperator(operator)
				if err != nil {
					return nil, fmt.Errorf("invalid operator: %s", operator)
				}
				if typedEnum, ok := operatorEnum.(*graphmodels.CustomSecurityAttributeComparisonOperator); ok {
					exemption.SetOperator(typedEnum)
				}
			}

			convert.FrameworkToGraphString(csa.Value, exemption.SetValue)

			csaExemptions = append(csaExemptions, exemption)
		}
		exemptions.SetCustomSecurityAttributes(csaExemptions)
	}

	return exemptions, nil
}

// constructDefaultPolicy constructs a default policy to restore settings
func constructDefaultPolicy() (graphmodels.TenantAppManagementPolicyable, error) {
	policy := graphmodels.NewTenantAppManagementPolicy()

	isEnabled := false
	policy.SetIsEnabled(&isEnabled)

	appRestrictions := graphmodels.NewAppManagementApplicationConfiguration()
	appRestrictions.SetPasswordCredentials([]graphmodels.PasswordCredentialConfigurationable{})
	appRestrictions.SetKeyCredentials([]graphmodels.KeyCredentialConfigurationable{})
	policy.SetApplicationRestrictions(appRestrictions)

	spRestrictions := graphmodels.NewAppManagementServicePrincipalConfiguration()
	spRestrictions.SetPasswordCredentials([]graphmodels.PasswordCredentialConfigurationable{})
	spRestrictions.SetKeyCredentials([]graphmodels.KeyCredentialConfigurationable{})
	policy.SetServicePrincipalRestrictions(spRestrictions)

	return policy, nil
}
