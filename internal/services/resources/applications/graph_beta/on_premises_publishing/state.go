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
	if appType := onPremPub.GetApplicationType(); appType != nil {
		data.ApplicationType = convert.GraphToFrameworkString(appType)
	}

	if externalAuthType := onPremPub.GetExternalAuthenticationType(); externalAuthType != nil {
		data.ExternalAuthenticationType = convert.GraphToFrameworkEnum(externalAuthType)
	}

	data.InternalUrl = convert.GraphToFrameworkString(onPremPub.GetInternalUrl())
	data.ExternalUrl = convert.GraphToFrameworkString(onPremPub.GetExternalUrl())

	// Map boolean fields
	data.IsAccessibleViaZTNAClient = convert.GraphToFrameworkBool(onPremPub.GetIsAccessibleViaZTNAClient())
	data.IsHttpOnlyCookieEnabled = convert.GraphToFrameworkBool(onPremPub.GetIsHttpOnlyCookieEnabled())
	data.IsOnPremPublishingEnabled = convert.GraphToFrameworkBool(onPremPub.GetIsOnPremPublishingEnabled())
	data.IsPersistentCookieEnabled = convert.GraphToFrameworkBool(onPremPub.GetIsPersistentCookieEnabled())
	data.IsSecureCookieEnabled = convert.GraphToFrameworkBool(onPremPub.GetIsSecureCookieEnabled())
	data.IsStateSessionEnabled = convert.GraphToFrameworkBool(onPremPub.GetIsStateSessionEnabled())
	data.IsTranslateHostHeaderEnabled = convert.GraphToFrameworkBool(onPremPub.GetIsTranslateHostHeaderEnabled())
	data.IsTranslateLinksInBodyEnabled = convert.GraphToFrameworkBool(onPremPub.GetIsTranslateLinksInBodyEnabled())

	tflog.Debug(ctx, fmt.Sprintf("Finished mapping %s remote state to Terraform state", ResourceName))

	return data
}
