package graphBetaIdentityAndAccessLifecycleWorkflowCustomTaskExtension

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoft/kiota-abstractions-go/serialization"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
	igmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models/identitygovernance"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
)

// MapRemoteResourceStateToTerraform maps the remote custom task extension from the Kiota SDK to Terraform state
func MapRemoteResourceStateToTerraform(
	ctx context.Context,
	data LifecycleWorkflowCustomTaskExtensionResourceModel,
	remoteResource igmodels.CustomTaskExtensionable,
) LifecycleWorkflowCustomTaskExtensionResourceModel {
	if remoteResource == nil {
		tflog.Debug(ctx, "Remote resource is nil")
		return data
	}

	tflog.Debug(ctx, "Starting to map remote state to Terraform state", map[string]any{
		"resourceName": remoteResource.GetDisplayName(),
		"resourceId":   remoteResource.GetId(),
	})

	data.ID = convert.GraphToFrameworkString(remoteResource.GetId())
	data.DisplayName = convert.GraphToFrameworkString(remoteResource.GetDisplayName())
	data.Description = convert.GraphToFrameworkStringWithDefault(remoteResource.GetDescription(), "")
	data.ReplyMode = convert.GraphToFrameworkEnum(remoteResource.GetReplyMode())
	data.CreatedDateTime = convert.GraphToFrameworkTime(remoteResource.GetCreatedDateTime())
	data.LastModifiedDateTime = convert.GraphToFrameworkTime(remoteResource.GetLastModifiedDateTime())

	data.EndpointConfiguration = stateEndpointConfiguration(ctx, remoteResource.GetEndpointConfiguration())
	data.AuthenticationConfiguration = stateAuthenticationConfiguration(ctx, remoteResource.GetAuthenticationConfiguration())
	data.ClientConfiguration = stateClientConfiguration(remoteResource.GetClientConfiguration())
	data.CallbackConfiguration = stateCallbackConfiguration(remoteResource.GetCallbackConfiguration(), data.CallbackConfiguration)

	tflog.Debug(ctx, fmt.Sprintf("Finished mapping resource %s with id %s", ResourceName, data.ID.ValueString()))
	return data
}

// stateEndpointConfiguration maps a logicAppTriggerEndpointConfiguration. Other endpoint types are
// not used by lifecycle workflows.
func stateEndpointConfiguration(ctx context.Context, remote graphmodels.CustomExtensionEndpointConfigurationable) *EndpointConfiguration {
	if remote == nil {
		return nil
	}

	endpoint, ok := remote.(graphmodels.LogicAppTriggerEndpointConfigurationable)
	if !ok {
		tflog.Warn(ctx, fmt.Sprintf("Unsupported endpoint configuration type %v", remote.GetOdataType()))
		return nil
	}

	return &EndpointConfiguration{
		SubscriptionId:       convert.GraphToFrameworkString(endpoint.GetSubscriptionId()),
		ResourceGroupName:    convert.GraphToFrameworkString(endpoint.GetResourceGroupName()),
		LogicAppWorkflowName: convert.GraphToFrameworkString(endpoint.GetLogicAppWorkflowName()),
		Url:                  convert.GraphToFrameworkString(endpoint.GetUrl()),
	}
}

// stateAuthenticationConfiguration maps the polymorphic authentication configuration to type and resource_id.
func stateAuthenticationConfiguration(ctx context.Context, remote graphmodels.CustomExtensionAuthenticationConfigurationable) *AuthenticationConfiguration {
	switch authentication := remote.(type) {
	case nil:
		return nil
	case graphmodels.AzureAdTokenAuthenticationable:
		return &AuthenticationConfiguration{
			Type:       types.StringValue(authTypeAzureAd),
			ResourceId: convert.GraphToFrameworkString(authentication.GetResourceId()),
		}
	case graphmodels.AzureAdPopTokenAuthenticationable:
		return &AuthenticationConfiguration{
			Type:       types.StringValue(authTypeProofOfPossession),
			ResourceId: types.StringNull(),
		}
	}

	tflog.Warn(ctx, fmt.Sprintf("Unsupported authentication configuration type %v", remote.GetOdataType()))
	return nil
}

// stateClientConfiguration maps the client configuration.
func stateClientConfiguration(remote graphmodels.CustomExtensionClientConfigurationable) *ClientConfiguration {
	if remote == nil {
		return nil
	}

	return &ClientConfiguration{
		MaximumRetries:        convert.GraphToFrameworkInt32AsInt64(remote.GetMaximumRetries()),
		TimeoutInMilliseconds: convert.GraphToFrameworkInt32AsInt64(remote.GetTimeoutInMilliseconds()),
	}
}

// stateCallbackConfiguration maps the callback configuration.
//
// The SDK deserializes callbackConfiguration with the base customExtensionCallbackConfiguration
// discriminator, which does not know the identityGovernance subtype, so authorizedApps arrives in
// AdditionalData as []any of map[string]any with *string IDs.
//
// timeout_duration keeps the prior value when it is the same duration Graph returned in a shorter
// form (PT60M and PT1H). Graph returns authorizedApps as [] when none are set, so an empty list
// stays null unless the prior state held an explicit empty set.
func stateCallbackConfiguration(remote graphmodels.CustomExtensionCallbackConfigurationable, prior *CallbackConfiguration) *CallbackConfiguration {
	if remote == nil {
		return nil
	}

	callback := &CallbackConfiguration{
		TimeoutDuration: convert.GraphToFrameworkISODuration(remote.GetTimeoutDuration()),
	}

	if prior != nil && sameDuration(prior.TimeoutDuration, remote.GetTimeoutDuration()) {
		callback.TimeoutDuration = prior.TimeoutDuration
	}

	ids := authorizedAppIds(remote)
	switch {
	case len(ids) > 0:
		callback.AuthorizedAppIds = types.SetValueMust(types.StringType, ids)
	case prior != nil && !prior.AuthorizedAppIds.IsNull() && !prior.AuthorizedAppIds.IsUnknown():
		callback.AuthorizedAppIds = types.SetValueMust(types.StringType, []attr.Value{})
	default:
		callback.AuthorizedAppIds = types.SetNull(types.StringType)
	}

	return callback
}

// authorizedAppIds reads authorizedApps from the typed model when the SDK resolved the subtype,
// otherwise from AdditionalData.
func authorizedAppIds(remote graphmodels.CustomExtensionCallbackConfigurationable) []attr.Value {
	ids := make([]attr.Value, 0)

	if typed, ok := remote.(igmodels.CustomTaskExtensionCallbackConfigurationable); ok {
		for _, app := range typed.GetAuthorizedApps() {
			if app != nil && app.GetId() != nil {
				ids = append(ids, types.StringValue(*app.GetId()))
			}
		}
		return ids
	}

	apps, _ := remote.GetAdditionalData()["authorizedApps"].([]any)
	for _, item := range apps {
		app, _ := item.(map[string]any)
		if id, ok := app["id"].(*string); ok && id != nil {
			ids = append(ids, types.StringValue(*id))
		}
	}

	return ids
}

// sameDuration reports whether the prior timeout_duration and the remote value are the same duration.
func sameDuration(prior types.String, remote *serialization.ISODuration) bool {
	if prior.IsNull() || prior.IsUnknown() || remote == nil {
		return false
	}

	parsed, err := serialization.ParseISODuration(prior.ValueString())
	if err != nil {
		return false
	}

	return parsed.String() == remote.String()
}
