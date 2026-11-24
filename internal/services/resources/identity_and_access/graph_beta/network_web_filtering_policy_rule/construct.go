package graphBetaNetworkWebFilteringPolicyRule

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const (
	// These OData discriminator values are observed from the Entra portal Web
	// content filtering blade. They are not exposed as webFiltering* models in
	// the Microsoft Graph beta SDK.
	webFilteringRuleODataType                   = "#microsoft.graph.networkaccess.webFilteringRule"
	webFilteringActionAllowODataType            = "#microsoft.graph.networkaccess.webFilteringActionAllow"
	webFilteringActionBlockODataType            = "#microsoft.graph.networkaccess.webFilteringActionBlock"
	webFilteringURLDestinationODataType         = "#microsoft.graph.networkaccess.webFilteringUrlDestination"
	webFilteringWebCategoryDestinationODataType = "#microsoft.graph.networkaccess.webFilteringWebCategoryDestination"
	headerModificationAddODataType              = "#microsoft.graph.networkaccess.headerModificationAdd"
)

var (
	errReadAttribute             = errors.New("failed to read")
	errNoDestinations            = errors.New("at least one destination must be specified using urls_or_fqdns or web_categories")
	errInvalidCustomHeaders      = errors.New("invalid custom_headers")
	errCustomHeadersRequireAllow = errors.New(`custom_headers can only be set when action is "allow"`)
)

// constructResource builds the portal-observed webFilteringRule payload.
//
// Microsoft Learn documents the generic filteringRule hierarchy, but not
// /networkaccess/webFilteringPolicies/{id}/policyRules or the webFiltering*
// OData discriminator values used by the Entra Global Secure Access Web content
// filtering blade:
// https://learn.microsoft.com/graph/api/resources/networkaccess-filteringrule
//
// The shape below intentionally follows captured portal traffic:
//
//   - URL/FQDN and category destinations are sibling entries in
//     matchingConditions.destinations.targets.
//
//   - URL/FQDN destinations are sent to Graph as a values array. The portal
//     offers a comma-delimited text box for usability, but Terraform keeps the
//     API shape as set(string).
//
//   - HTTP methods and session types are serialized back to comma-separated
//     strings because that is what the endpoint returns and accepts.
//
//   - Custom headers are sent as action.headerSettings.modifications, one
//     headerModificationAdd entry per header, in configuration order. Allow
//     rules without headers send headerSettings: null; block rules omit it.
func constructResource(ctx context.Context, data *NetworkWebFilteringPolicyRuleResourceModel) (webFilteringRuleRequest, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))

	setValues := func(name string, value types.Set) ([]string, error) {
		if value.IsNull() || value.IsUnknown() {
			return nil, nil
		}
		var result []string
		if diags := value.ElementsAs(ctx, &result, false); diags.HasError() {
			return nil, fmt.Errorf("%w %s: %s", errReadAttribute, name, diags.Errors()[0].Detail())
		}
		return result, nil
	}
	commaSeparated := func(values []string) *string {
		if len(values) == 0 {
			return nil
		}
		value := strings.Join(values, ", ")
		return &value
	}

	urlsOrFqdns, err := setValues("urls_or_fqdns", data.UrlsOrFqdns)
	if err != nil {
		return webFilteringRuleRequest{}, err
	}
	webCategories, err := setValues("web_categories", data.WebCategories)
	if err != nil {
		return webFilteringRuleRequest{}, err
	}
	httpMethods, err := setValues("http_methods", data.HTTPMethods)
	if err != nil {
		return webFilteringRuleRequest{}, err
	}
	sessionTypes, err := setValues("session_types", data.SessionTypes)
	if err != nil {
		return webFilteringRuleRequest{}, err
	}
	if len(urlsOrFqdns) == 0 && len(webCategories) == 0 {
		return webFilteringRuleRequest{}, errNoDestinations
	}

	// The Entra portal and Graph endpoint reject rules without a destination.
	// Portal UI text asks for at least one URL/FQDN or web category. Terraform
	// follows the Graph payload shape and models URL/FQDN values as set(string)
	// even though the portal combines them into one comma-delimited text box.
	targets := make([]webFilteringDestinationTarget, 0, 2)
	if len(urlsOrFqdns) > 0 {
		targets = append(targets, webFilteringDestinationTarget{ODataType: webFilteringURLDestinationODataType, Values: urlsOrFqdns})
	}
	if len(webCategories) > 0 {
		targets = append(targets, webFilteringDestinationTarget{ODataType: webFilteringWebCategoryDestinationODataType, Values: webCategories})
	}

	var customHeaders []CustomHeaderModel
	if !data.CustomHeaders.IsNull() && !data.CustomHeaders.IsUnknown() {
		if diags := data.CustomHeaders.ElementsAs(ctx, &customHeaders, false); diags.HasError() {
			return webFilteringRuleRequest{}, fmt.Errorf("failed to read custom_headers: %w: %s", errInvalidCustomHeaders, diags.Errors()[0].Detail())
		}
	}

	// Block actions reject any headerSettings, even null. Allow actions keep
	// existing headers when it is omitted on PATCH, so send null to clear them.
	var action any = webFilteringBlockAction{ODataType: webFilteringActionBlockODataType}
	if data.Action.ValueString() != "block" {
		allow := webFilteringRuleAction{ODataType: webFilteringActionAllowODataType}
		if len(customHeaders) > 0 {
			allow.HeaderSettings = &webFilteringHeaderSettings{Modifications: make([]webFilteringHeaderModification, 0, len(customHeaders))}
			for _, header := range customHeaders {
				modification := webFilteringHeaderModification{
					ODataType:   headerModificationAddODataType,
					HeaderName:  header.Name.ValueString(),
					HeaderValue: header.Value.ValueString(),
				}
				allow.HeaderSettings.Modifications = append(allow.HeaderSettings.Modifications, modification)
			}
		}
		action = allow
	} else if len(customHeaders) > 0 {
		// Fail before calling Graph; the config validator skips an action unknown at plan time.
		return webFilteringRuleRequest{}, fmt.Errorf("%w, got %q", errCustomHeadersRequireAllow, data.Action.ValueString())
	}

	requestBody := webFilteringRuleRequest{
		ODataType:   webFilteringRuleODataType,
		Name:        data.Name.ValueStringPointer(),
		Description: data.Description.ValueStringPointer(),
		Action:      action,
		Priority:    data.Priority.ValueInt64Pointer(),
		Settings:    webFilteringRuleSettings{Status: data.Status.ValueStringPointer()},
		MatchingConditions: webFilteringMatchingConditions{
			Destinations: &webFilteringDestinations{Targets: targets, HTTPRequestMethod: commaSeparated(httpMethods)},
		},
	}
	if len(sessionTypes) > 0 {
		requestBody.MatchingConditions.Sources = &webFilteringSources{SessionType: commaSeparated(sessionTypes)}
	}

	return requestBody, nil
}
