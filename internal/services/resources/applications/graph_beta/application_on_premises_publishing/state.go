package graphBetaApplicationsOnPremisesPublishing

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteStateToTerraform maps the remote state from Microsoft Graph API to the Terraform state
func MapRemoteStateToTerraform(ctx context.Context, data OnPremisesPublishingResourceModel, application graphmodels.Applicationable) OnPremisesPublishingResourceModel {
	tflog.Debug(ctx, fmt.Sprintf("Mapping %s remote state to Terraform state", ResourceName))

	if application == nil {
		tflog.Debug(ctx, "Application is nil")
		return data
	}

	onPremPub := application.GetOnPremisesPublishing()
	if onPremPub == nil {
		tflog.Debug(ctx, "OnPremisesPublishing is nil")
		return data
	}

	// Map string fields
	data.AlternateUrl = convert.GraphToFrameworkString(onPremPub.GetAlternateUrl())
	data.ApplicationServerTimeout = convert.GraphToFrameworkString(onPremPub.GetApplicationServerTimeout())
	data.ApplicationType = convert.GraphToFrameworkString(onPremPub.GetApplicationType())
	data.InternalUrl = convert.GraphToFrameworkString(onPremPub.GetInternalUrl())
	data.ExternalUrl = convert.GraphToFrameworkString(onPremPub.GetExternalUrl())
	data.WafProvider = convert.GraphToFrameworkString(onPremPub.GetWafProvider())

	// Map enum fields
	if externalAuthType := onPremPub.GetExternalAuthenticationType(); externalAuthType != nil {
		data.ExternalAuthenticationType = convert.GraphToFrameworkEnum(externalAuthType)
	}

	// Map boolean fields
	data.IsAccessibleViaZTNAClient = convert.GraphToFrameworkBool(onPremPub.GetIsAccessibleViaZTNAClient())
	data.IsBackendCertificateValidationEnabled = convert.GraphToFrameworkBool(onPremPub.GetIsBackendCertificateValidationEnabled())
	data.IsContinuousAccessEvaluationEnabled = convert.GraphToFrameworkBool(onPremPub.GetIsContinuousAccessEvaluationEnabled())
	data.IsDnsResolutionEnabled = convert.GraphToFrameworkBool(onPremPub.GetIsDnsResolutionEnabled())
	data.IsHttpOnlyCookieEnabled = convert.GraphToFrameworkBool(onPremPub.GetIsHttpOnlyCookieEnabled())
	data.IsOnPremPublishingEnabled = convert.GraphToFrameworkBool(onPremPub.GetIsOnPremPublishingEnabled())
	data.IsPersistentCookieEnabled = convert.GraphToFrameworkBool(onPremPub.GetIsPersistentCookieEnabled())
	data.IsSecureCookieEnabled = convert.GraphToFrameworkBool(onPremPub.GetIsSecureCookieEnabled())
	data.IsStateSessionEnabled = convert.GraphToFrameworkBool(onPremPub.GetIsStateSessionEnabled())
	data.IsTranslateHostHeaderEnabled = convert.GraphToFrameworkBool(onPremPub.GetIsTranslateHostHeaderEnabled())
	data.IsTranslateLinksInBodyEnabled = convert.GraphToFrameworkBool(onPremPub.GetIsTranslateLinksInBodyEnabled())
	data.UseAlternateUrlForTranslationAndRedirect = convert.GraphToFrameworkBool(onPremPub.GetUseAlternateUrlForTranslationAndRedirect())

	tflog.Debug(ctx, fmt.Sprintf("Finished mapping %s remote state to Terraform state", ResourceName))

	return data
}
