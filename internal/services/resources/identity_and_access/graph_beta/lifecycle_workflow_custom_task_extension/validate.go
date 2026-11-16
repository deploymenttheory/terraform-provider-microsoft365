package graphBetaIdentityAndAccessLifecycleWorkflowCustomTaskExtension

import (
	"context"
	"errors"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	errEndpointConfigurationRequired       = errors.New("endpoint_configuration is required")
	errAuthenticationConfigurationRequired = errors.New("authentication_configuration is required")
	errResourceIdRequired                  = errors.New("resource_id is required when type is azureAdTokenAuthentication")
	errResourceIdNotAllowed                = errors.New("resource_id must not be set when type is azureAdPopTokenAuthentication")
)

// ValidateConfig checks the authentication configuration at plan time so an invalid combination
// fails before any Graph call is made.
func (r *LifecycleWorkflowCustomTaskExtensionResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var authentication types.Object
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("authentication_configuration"), &authentication)...)
	if resp.Diagnostics.HasError() || authentication.IsNull() || authentication.IsUnknown() {
		return
	}

	var data *AuthenticationConfiguration
	if diags := req.Config.GetAttribute(ctx, path.Root("authentication_configuration"), &data); diags.HasError() || data == nil {
		return
	}

	if err := validateAuthenticationConfiguration(data); err != nil {
		resp.Diagnostics.AddAttributeError(path.Root("authentication_configuration"), "Invalid authentication configuration", err.Error())
	}
}

// validateRequest validates the model before it is sent to Graph.
func validateRequest(data *LifecycleWorkflowCustomTaskExtensionResourceModel) error {
	if data.EndpointConfiguration == nil {
		return errEndpointConfigurationRequired
	}
	if data.AuthenticationConfiguration == nil {
		return errAuthenticationConfigurationRequired
	}
	return validateAuthenticationConfiguration(data.AuthenticationConfiguration)
}

// validateAuthenticationConfiguration enforces that resource_id is set only for azureAdTokenAuthentication.
// Values that are still unknown are skipped; they are checked again when the request is built.
func validateAuthenticationConfiguration(data *AuthenticationConfiguration) error {
	if data.Type.IsUnknown() || data.ResourceId.IsUnknown() {
		return nil
	}

	hasResourceId := !data.ResourceId.IsNull() && data.ResourceId.ValueString() != ""

	switch data.Type.ValueString() {
	case authTypeAzureAd:
		if !hasResourceId {
			return errResourceIdRequired
		}
	case authTypeProofOfPossession:
		if hasResourceId {
			return errResourceIdNotAllowed
		}
	}

	return nil
}
