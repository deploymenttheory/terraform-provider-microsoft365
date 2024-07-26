package provider

import (
	"context"
	"fmt"
	"net/url"
	"regexp"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

/* tenant_id and client_id schema validator */
type guidValidator struct {
	attributeName string
}

func (v guidValidator) Description(ctx context.Context) string {
	return fmt.Sprintf("%s must be a valid GUID if provided", v.attributeName)
}

func (v guidValidator) MarkdownDescription(ctx context.Context) string {
	return fmt.Sprintf("%s must be a valid GUID if provided", v.attributeName)
}

func (v guidValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	// If the value is null, unknown, or empty string, return without validation
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	value := req.ConfigValue.ValueString()
	if value == "" {
		return
	}

	// Proceed with GUID validation only if the string is non-empty
	match, err := regexp.MatchString(helpers.GuidRegex, value)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Invalid %s format", v.attributeName),
			fmt.Sprintf("Error matching GUID format for %s: %s", v.attributeName, err),
		)
		return
	}

	if !match {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Invalid %s format", v.attributeName),
			fmt.Sprintf("The value %q for %s is not a valid GUID. It must match the format: xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx", value, v.attributeName),
		)
	}
}

func validateGUID(attributeName string) validator.String {
	return guidValidator{attributeName: attributeName}
}

/* use_proxy field schema validator */

type useProxyValidator struct{}

func (v useProxyValidator) Description(ctx context.Context) string {
	return "Ensures that proxy_url is set if use_proxy is true."
}

func (v useProxyValidator) MarkdownDescription(ctx context.Context) string {
	return "Ensures that proxy_url is set if use_proxy is true."
}

func validateUseProxy() validator.Bool {
	return useProxyValidator{}
}

func (v useProxyValidator) ValidateBool(ctx context.Context, request validator.BoolRequest, response *validator.BoolResponse) {
	var useProxy types.Bool
	if diags := request.Config.GetAttribute(ctx, path.Root("use_proxy"), &useProxy); diags.HasError() {
		response.Diagnostics.Append(diags...)
		return
	}

	var proxyURL types.String
	if diags := request.Config.GetAttribute(ctx, path.Root("proxy_url"), &proxyURL); diags.HasError() {
		response.Diagnostics.Append(diags...)
		return
	}

	if useProxy.ValueBool() && (proxyURL.IsNull() || proxyURL.IsUnknown() || proxyURL.ValueString() == "") {
		response.Diagnostics.AddError(
			"Invalid Configuration",
			"The 'proxy_url' attribute must be set when 'use_proxy' is true.",
		)
	}
}

/* redirect_url, proxy_url, token_endpoint fields schema validator */

type urlValidator struct{}

func (v urlValidator) Description(ctx context.Context) string {
	return "Validates that the value is a valid URL."
}

func (v urlValidator) MarkdownDescription(ctx context.Context) string {
	return "Validates that the value is a valid URL."
}

func validateURL() validator.String {
	return urlValidator{}
}

// ValidateString validates the "redirect_url", "proxy_url", or any URL attribute.
func (v urlValidator) ValidateString(ctx context.Context, request validator.StringRequest, response *validator.StringResponse) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}

	re := regexp.MustCompile(helpers.UrlValidStringRegex)

	if !re.MatchString(request.ConfigValue.ValueString()) {
		response.Diagnostics.AddError(
			"Invalid URL",
			"The value must be a valid URL.",
		)
		return
	}

	u, err := url.ParseRequestURI(request.ConfigValue.ValueString())
	if err != nil || u.Scheme == "" || u.Host == "" {
		response.Diagnostics.AddError(
			"Invalid URL",
			"The value must be a valid URL.",
		)
	}
}
