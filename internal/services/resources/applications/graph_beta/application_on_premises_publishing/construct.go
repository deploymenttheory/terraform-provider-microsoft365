package graphBetaApplicationsOnPremisesPublishing

import "github.com/hashicorp/terraform-plugin-framework/types"

func constructOnPremisesPublishingPatchPayload(data OnPremisesPublishingResourceModel) map[string]any {
	onPremisesPublishing := map[string]any{}

	addStringIfKnown(onPremisesPublishing, "alternateUrl", data.AlternateUrl)
	addStringIfKnown(onPremisesPublishing, "applicationServerTimeout", data.ApplicationServerTimeout)
	addStringIfKnown(onPremisesPublishing, "applicationType", data.ApplicationType)
	addStringIfKnown(onPremisesPublishing, "externalAuthenticationType", data.ExternalAuthenticationType)
	addStringIfKnown(onPremisesPublishing, "internalUrl", data.InternalUrl)
	addStringIfKnown(onPremisesPublishing, "externalUrl", data.ExternalUrl)
	addStringIfKnown(onPremisesPublishing, "trafficRoutingMethod", data.TrafficRoutingMethod)
	addStringIfKnown(onPremisesPublishing, "wafProvider", data.WafProvider)

	addBoolIfKnown(onPremisesPublishing, "isAccessibleViaZTNAClient", data.IsAccessibleViaZTNAClient)
	addBoolIfKnown(onPremisesPublishing, "isBackendCertificateValidationEnabled", data.IsBackendCertificateValidationEnabled)
	addBoolIfKnown(onPremisesPublishing, "isContinuousAccessEvaluationEnabled", data.IsContinuousAccessEvaluationEnabled)
	addBoolIfKnown(onPremisesPublishing, "isDnsResolutionEnabled", data.IsDnsResolutionEnabled)
	addBoolIfKnown(onPremisesPublishing, "isHttpOnlyCookieEnabled", data.IsHttpOnlyCookieEnabled)
	addBoolIfKnown(onPremisesPublishing, "isOnPremPublishingEnabled", data.IsOnPremPublishingEnabled)
	addBoolIfKnown(onPremisesPublishing, "isPersistentCookieEnabled", data.IsPersistentCookieEnabled)
	addBoolIfKnown(onPremisesPublishing, "isSecureCookieEnabled", data.IsSecureCookieEnabled)
	addBoolIfKnown(onPremisesPublishing, "isStateSessionEnabled", data.IsStateSessionEnabled)
	addBoolIfKnown(onPremisesPublishing, "isTranslateHostHeaderEnabled", data.IsTranslateHostHeaderEnabled)
	addBoolIfKnown(onPremisesPublishing, "isTranslateLinksInBodyEnabled", data.IsTranslateLinksInBodyEnabled)
	addBoolIfKnown(onPremisesPublishing, "useAlternateUrlForTranslationAndRedirect", data.UseAlternateUrlForTranslationAndRedirect)

	// The Graph beta applications PATCH endpoint rejects the top-level @odata.type
	// emitted by the Kiota Application model for this nested property update.
	return map[string]any{
		"onPremisesPublishing": onPremisesPublishing,
	}
}

func addStringIfKnown(values map[string]any, key string, value types.String) {
	if !value.IsNull() && !value.IsUnknown() {
		values[key] = value.ValueString()
	}
}

func addBoolIfKnown(values map[string]any, key string, value types.Bool) {
	if !value.IsNull() && !value.IsUnknown() {
		values[key] = value.ValueBool()
	}
}
