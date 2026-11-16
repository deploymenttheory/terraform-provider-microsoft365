package graphBetaIdentityAndAccessLifecycleWorkflowCustomTaskExtension

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
	igmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models/identitygovernance"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
)

// constructResource converts the Terraform resource model to a Kiota SDK custom task extension.
// The same body is used for POST and PATCH; Graph accepts every writable property on PATCH.
//
// isUpdate controls how an unset callback_configuration is sent. On create it is omitted, which
// gives reply mode none. On update it is sent as null, because omitting it from a PATCH keeps the
// existing callback configuration; Graph then switches reply mode to none.
func constructResource(ctx context.Context, data *LifecycleWorkflowCustomTaskExtensionResourceModel, isUpdate bool) (igmodels.CustomTaskExtensionable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))

	if err := validateRequest(data); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	requestBody := igmodels.NewCustomTaskExtension()

	convert.FrameworkToGraphString(data.DisplayName, requestBody.SetDisplayName)
	convert.FrameworkToGraphString(data.Description, requestBody.SetDescription)

	requestBody.SetEndpointConfiguration(constructEndpointConfiguration(data.EndpointConfiguration))
	requestBody.SetAuthenticationConfiguration(constructAuthenticationConfiguration(data.AuthenticationConfiguration))
	requestBody.SetClientConfiguration(constructClientConfiguration(data.ClientConfiguration))

	if data.CallbackConfiguration != nil {
		callbackConfiguration, err := constructCallbackConfiguration(data.CallbackConfiguration)
		if err != nil {
			return nil, fmt.Errorf("failed to construct callback configuration: %w", err)
		}
		requestBody.SetCallbackConfiguration(callbackConfiguration)
	} else if isUpdate {
		requestBody.SetAdditionalData(map[string]any{"callbackConfiguration": nil})
	}

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s resource", ResourceName))

	return requestBody, nil
}

// constructEndpointConfiguration builds a logicAppTriggerEndpointConfiguration. An unknown or
// null url is omitted so Graph keeps or derives its own value.
func constructEndpointConfiguration(data *EndpointConfiguration) graphmodels.CustomExtensionEndpointConfigurationable {
	endpoint := graphmodels.NewLogicAppTriggerEndpointConfiguration()

	convert.FrameworkToGraphString(data.SubscriptionId, endpoint.SetSubscriptionId)
	convert.FrameworkToGraphString(data.ResourceGroupName, endpoint.SetResourceGroupName)
	convert.FrameworkToGraphString(data.LogicAppWorkflowName, endpoint.SetLogicAppWorkflowName)
	convert.FrameworkToGraphString(data.Url, endpoint.SetUrl)

	return endpoint
}

// constructAuthenticationConfiguration builds azureAdPopTokenAuthentication or azureAdTokenAuthentication.
func constructAuthenticationConfiguration(data *AuthenticationConfiguration) graphmodels.CustomExtensionAuthenticationConfigurationable {
	if data.Type.ValueString() == authTypeAzureAd {
		authentication := graphmodels.NewAzureAdTokenAuthentication()
		convert.FrameworkToGraphString(data.ResourceId, authentication.SetResourceId)
		return authentication
	}

	return graphmodels.NewAzureAdPopTokenAuthentication()
}

// constructClientConfiguration builds the customExtensionClientConfiguration. Graph rejects a
// custom task extension without one ("Invalid model state"), so a missing block falls back to the
// schema defaults.
func constructClientConfiguration(data *ClientConfiguration) graphmodels.CustomExtensionClientConfigurationable {
	client := graphmodels.NewCustomExtensionClientConfiguration()

	maximumRetries := types.Int64Value(defaultMaximumRetries)
	timeoutInMilliseconds := types.Int64Value(defaultTimeoutInMilliseconds)
	if data != nil {
		if !data.MaximumRetries.IsNull() && !data.MaximumRetries.IsUnknown() {
			maximumRetries = data.MaximumRetries
		}
		if !data.TimeoutInMilliseconds.IsNull() && !data.TimeoutInMilliseconds.IsUnknown() {
			timeoutInMilliseconds = data.TimeoutInMilliseconds
		}
	}

	retries := int32(maximumRetries.ValueInt64())        //nolint:gosec // maximum_retries is validated to 0..1 by the schema
	timeout := int32(timeoutInMilliseconds.ValueInt64()) //nolint:gosec // timeout_in_milliseconds is validated to 200..2000 by the schema
	client.SetMaximumRetries(&retries)
	client.SetTimeoutInMilliseconds(&timeout)

	return client
}

// constructCallbackConfiguration builds an identityGovernance.customTaskExtensionCallbackConfiguration.
// authorizedApps is always sent, as an empty array when no app IDs are configured, so removing
// every app ID clears the list.
func constructCallbackConfiguration(data *CallbackConfiguration) (graphmodels.CustomExtensionCallbackConfigurationable, error) {
	callback := igmodels.NewCustomTaskExtensionCallbackConfiguration()

	if err := convert.FrameworkToGraphISODuration(data.TimeoutDuration, callback.SetTimeoutDuration); err != nil {
		return nil, fmt.Errorf("failed to set timeout duration: %w", err)
	}

	apps := make([]graphmodels.Applicationable, 0, len(data.AuthorizedAppIds.Elements()))
	for _, element := range data.AuthorizedAppIds.Elements() {
		value, ok := element.(types.String)
		if !ok || value.IsNull() || value.IsUnknown() {
			continue
		}
		id := value.ValueString()
		app := graphmodels.NewApplication()
		app.SetId(&id)
		apps = append(apps, app)
	}
	callback.SetAuthorizedApps(apps)

	return callback, nil
}
