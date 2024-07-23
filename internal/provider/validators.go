package provider

import (
	"context"
	"fmt"
	"net/url"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (p *M365Provider) ValidateConfig(ctx context.Context, req provider.ValidateConfigRequest, resp *provider.ValidateConfigResponse) {
	var data M365ProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Check if both client_certificate and client_certificate_file_path are provided.
	if !data.ClientCertificateBase64.IsNull() && !data.ClientCertificateBase64.IsUnknown() &&
		!data.ClientCertificateFilePath.IsNull() && !data.ClientCertificateFilePath.IsUnknown() {
		resp.Diagnostics.AddError(
			"Conflicting Configuration",
			"Only one of 'client_certificate_base64' or 'client_certificate_file_path' can be provided. Please choose one.",
		)
	}
}

/* auth_method schema validator */

type authMethodValidator struct{}

func (v authMethodValidator) Description(ctx context.Context) string {
	return "Validates that the necessary fields are set based on the authentication method"
}

func (v authMethodValidator) MarkdownDescription(ctx context.Context) string {
	return "Validates that the necessary fields are set based on the authentication method"
}

func validateAuthMethod() validator.String {
	return authMethodValidator{}
}

func (v authMethodValidator) ValidateString(ctx context.Context, request validator.StringRequest, response *validator.StringResponse) {
	var authMethod types.String
	if diags := request.Config.GetAttribute(ctx, path.Root("auth_method"), &authMethod); diags.HasError() {
		response.Diagnostics.Append(diags...)
		return
	}

	validAuthMethods := map[string]bool{
		"device_code":         true,
		"client_secret":       true,
		"client_certificate":  true,
		"interactive_browser": true,
		"username_password":   true,
	}

	if _, valid := validAuthMethods[authMethod.ValueString()]; !valid {
		response.Diagnostics.AddError(
			"Invalid Authentication Method",
			"The 'auth_method' must be one of 'device_code', 'client_secret', 'client_certificate', 'interactive_browser', 'username_password'.",
		)
		return
	}

	isSet := func(attrName string) bool {
		var attr types.String
		if diags := request.Config.GetAttribute(ctx, path.Root(attrName), &attr); diags.HasError() {
			return false
		}
		return !attr.IsNull() && !attr.IsUnknown()
	}

	switch authMethod.ValueString() {
	case "client_secret":
		if !isSet("client_secret") {
			response.Diagnostics.AddError(
				"Invalid Configuration",
				"The 'client_secret' attribute must be set when 'auth_method' is 'client_secret'.",
			)
		}
	case "client_certificate":
		if !isSet("client_certificate") && !isSet("client_certificate_file_path") {
			response.Diagnostics.AddError(
				"Invalid Configuration",
				"Either 'client_certificate' or 'client_certificate_file_path' must be set when 'auth_method' is 'client_certificate'.",
			)
		}
		if isSet("client_certificate") && isSet("client_certificate_file_path") {
			response.Diagnostics.AddError(
				"Invalid Configuration",
				"Only one of 'client_certificate' or 'client_certificate_file_path' can be set when 'auth_method' is 'client_certificate'.",
			)
		}
	case "interactive_browser":
		if !isSet("redirect_url") {
			response.Diagnostics.AddError(
				"Invalid Configuration",
				"The 'redirect_url' attribute must be set when 'auth_method' is 'interactive_browser'.",
			)
		}
	case "username_password":
		if !isSet("username") || !isSet("password") {
			response.Diagnostics.AddError(
				"Invalid Configuration",
				"The 'username' and 'password' attributes must be set when 'auth_method' is 'username_password'.",
			)
		}
	}
}

/* tenant_id and client_id schema validator */

type guidValidator struct{}

func (v guidValidator) Description(ctx context.Context) string {
	return "Validates that the value is a valid GUID."
}

func (v guidValidator) MarkdownDescription(ctx context.Context) string {
	return "Validates that the value is a valid GUID."
}

func validateGUID() validator.String {
	return guidValidator{}
}

func (v guidValidator) ValidateString(ctx context.Context, request validator.StringRequest, response *validator.StringResponse) {
	guidRegex := `^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`
	re := regexp.MustCompile(guidRegex)

	if !re.MatchString(request.ConfigValue.ValueString()) {
		response.Diagnostics.AddError(
			"Invalid GUID",
			"The value must be a valid GUID.",
		)
	}
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

func (v urlValidator) ValidateString(ctx context.Context, request validator.StringRequest, response *validator.StringResponse) {
	u, err := url.ParseRequestURI(request.ConfigValue.ValueString())
	if err != nil || u.Scheme == "" || u.Host == "" {
		response.Diagnostics.AddError(
			"Invalid URL",
			"The value must be a valid URL.",
		)
	}
}

/* cloud field schema validator */

// cloudValidator is a custom validator for the "cloud" attribute.
type cloudValidator struct{}

// Description returns a plain text description of the validator's behavior.
func (v cloudValidator) Description(ctx context.Context) string {
	return "Ensures the cloud attribute is one of 'public', 'gcc', 'gcchigh', 'china', 'dod', 'ex', or 'rx'."
}

// MarkdownDescription returns a markdown description of the validator's behavior.
func (v cloudValidator) MarkdownDescription(ctx context.Context) string {
	return "Ensures the cloud attribute is one of `public`, `gcc`, `gcchigh`, `china`, `dod`, `ex`, or `rx`."
}

// ValidateString validates the "cloud" attribute.
func (v cloudValidator) ValidateString(ctx context.Context, request validator.StringRequest, response *validator.StringResponse) {
	validClouds := []string{"public", "gcc", "gcchigh", "china", "dod", "ex", "rx"}
	for _, validCloud := range validClouds {
		if request.ConfigValue.ValueString() == validCloud {
			return
		}
	}
	response.Diagnostics.AddError(
		"Invalid Cloud Value",
		fmt.Sprintf("The 'cloud' attribute must be one of 'public', 'gcc', 'gcchigh', 'china', 'dod', 'ex', or 'rx'. Got: %s", request.ConfigValue.ValueString()),
	)
}

// validateCloud returns an instance of the cloudValidator.
func validateCloud() validator.String {
	return cloudValidator{}
}
