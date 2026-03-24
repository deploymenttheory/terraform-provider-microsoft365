package graphBetaChangeNotificationsSubscription

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// subscriptionResourcePathValidator ensures the value matches Microsoft Graph expectations for the
// subscription **resource** property: no URL scheme, no Graph host, no API version path prefix.
// See https://learn.microsoft.com/en-us/graph/api/resources/subscription?view=graph-rest-beta
type subscriptionResourcePathValidator struct{}

func subscriptionResourcePath() validator.String {
	return subscriptionResourcePathValidator{}
}

func (subscriptionResourcePathValidator) Description(_ context.Context) string {
	return "must be a Graph subscription resource path without https://, a Graph host name, or a beta/v1.0 prefix (see subscription resource property)"
}

func (subscriptionResourcePathValidator) MarkdownDescription(_ context.Context) string {
	return "Must satisfy the **resource** property rules in the [subscription resource type](https://learn.microsoft.com/en-us/graph/api/resources/subscription?view=graph-rest-beta): do not include the base URL (`https://graph.microsoft.com/beta/` or `https://graph.microsoft.com/v1.0/`), scheme, or hostname."
}

func (subscriptionResourcePathValidator) ValidateString(_ context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}
	s := req.ConfigValue.ValueString()
	if s == "" {
		return
	}
	if strings.TrimSpace(s) != s {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid resource",
			"The value must not have leading or trailing whitespace.",
		)
		return
	}
	if len(s) > 8192 {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid resource",
			"The value exceeds the maximum allowed length.",
		)
		return
	}
	if strings.Contains(s, "://") {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid resource",
			"Do not include a URL scheme or full URL. Supply only the path stored in the subscription **resource** property.",
		)
		return
	}
	lower := strings.ToLower(s)
	for _, host := range graphHostSubstrings {
		if strings.Contains(lower, host) {
			resp.Diagnostics.AddAttributeError(
				req.Path,
				"Invalid resource",
				"Do not include the Microsoft Graph host name in the path.",
			)
			return
		}
	}
	trimmed := strings.TrimLeft(s, "/")
	lt := strings.ToLower(trimmed)
	if strings.HasPrefix(lt, "beta/") || strings.HasPrefix(lt, "v1.0/") {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid resource",
			"Do not prefix the path with `beta/` or `v1.0/`; the API version is determined by the Graph client.",
		)
		return
	}
}

var graphHostSubstrings = []string{
	"graph.microsoft.com",
	"graph.microsoft.us",
	"dod-graph.microsoft.us",
	"graph.microsoft.de",
	"microsoftgraph.chinacloudapi.cn",
	"canary.graph.microsoft.com",
}
