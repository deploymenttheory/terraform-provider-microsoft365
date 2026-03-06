package graphBetaApplicationsOnPremisesPublishing

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// OnPremisesPublishingResourceModel describes the resource data model.
type OnPremisesPublishingResourceModel struct {
	ApplicationID                            types.String   `tfsdk:"application_id"`
	AlternateUrl                             types.String   `tfsdk:"alternate_url"`
	ApplicationServerTimeout                 types.String   `tfsdk:"application_server_timeout"`
	ApplicationType                          types.String   `tfsdk:"application_type"`
	ExternalAuthenticationType               types.String   `tfsdk:"external_authentication_type"`
	InternalUrl                              types.String   `tfsdk:"internal_url"`
	ExternalUrl                              types.String   `tfsdk:"external_url"`
	IsAccessibleViaZTNAClient                types.Bool     `tfsdk:"is_accessible_via_ztna_client"`
	IsBackendCertificateValidationEnabled    types.Bool     `tfsdk:"is_backend_certificate_validation_enabled"`
	IsContinuousAccessEvaluationEnabled      types.Bool     `tfsdk:"is_continuous_access_evaluation_enabled"`
	IsDnsResolutionEnabled                   types.Bool     `tfsdk:"is_dns_resolution_enabled"`
	IsHttpOnlyCookieEnabled                  types.Bool     `tfsdk:"is_http_only_cookie_enabled"`
	IsOnPremPublishingEnabled                types.Bool     `tfsdk:"is_on_prem_publishing_enabled"`
	IsPersistentCookieEnabled                types.Bool     `tfsdk:"is_persistent_cookie_enabled"`
	IsSecureCookieEnabled                    types.Bool     `tfsdk:"is_secure_cookie_enabled"`
	IsStateSessionEnabled                    types.Bool     `tfsdk:"is_state_session_enabled"`
	IsTranslateHostHeaderEnabled             types.Bool     `tfsdk:"is_translate_host_header_enabled"`
	IsTranslateLinksInBodyEnabled            types.Bool     `tfsdk:"is_translate_links_in_body_enabled"`
	UseAlternateUrlForTranslationAndRedirect types.Bool     `tfsdk:"use_alternate_url_for_translation_and_redirect"`
	WafProvider                              types.String   `tfsdk:"waf_provider"`
	Timeouts                                 timeouts.Value `tfsdk:"timeouts"`
}
