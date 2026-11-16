// REF: https://learn.microsoft.com/en-us/graph/api/resources/identitygovernance-customtaskextension?view=graph-rest-beta
package graphBetaIdentityAndAccessLifecycleWorkflowCustomTaskExtension

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// LifecycleWorkflowCustomTaskExtensionResourceModel represents the schema for the lifecycle workflow custom task extension resource
type LifecycleWorkflowCustomTaskExtensionResourceModel struct {
	ID                          types.String                 `tfsdk:"id"`
	DisplayName                 types.String                 `tfsdk:"display_name"`
	Description                 types.String                 `tfsdk:"description"`
	EndpointConfiguration       *EndpointConfiguration       `tfsdk:"endpoint_configuration"`
	AuthenticationConfiguration *AuthenticationConfiguration `tfsdk:"authentication_configuration"`
	ClientConfiguration         *ClientConfiguration         `tfsdk:"client_configuration"`
	CallbackConfiguration       *CallbackConfiguration       `tfsdk:"callback_configuration"`
	ReplyMode                   types.String                 `tfsdk:"reply_mode"`
	CreatedDateTime             types.String                 `tfsdk:"created_date_time"`
	LastModifiedDateTime        types.String                 `tfsdk:"last_modified_date_time"`
	Timeouts                    timeouts.Value               `tfsdk:"timeouts"`
}

// EndpointConfiguration represents a logicAppTriggerEndpointConfiguration, the only endpoint type
// lifecycle workflows support.
type EndpointConfiguration struct {
	SubscriptionId       types.String `tfsdk:"subscription_id"`
	ResourceGroupName    types.String `tfsdk:"resource_group_name"`
	LogicAppWorkflowName types.String `tfsdk:"logic_app_workflow_name"`
	Url                  types.String `tfsdk:"url"`
}

// AuthenticationConfiguration represents the polymorphic customExtensionAuthenticationConfiguration.
// type selects azureAdPopTokenAuthentication or azureAdTokenAuthentication; resource_id is only
// used by azureAdTokenAuthentication.
type AuthenticationConfiguration struct {
	Type       types.String `tfsdk:"type"`
	ResourceId types.String `tfsdk:"resource_id"`
}

// ClientConfiguration represents a customExtensionClientConfiguration
type ClientConfiguration struct {
	MaximumRetries        types.Int64 `tfsdk:"maximum_retries"`
	TimeoutInMilliseconds types.Int64 `tfsdk:"timeout_in_milliseconds"`
}

// CallbackConfiguration represents an identityGovernance.customTaskExtensionCallbackConfiguration
type CallbackConfiguration struct {
	TimeoutDuration  types.String `tfsdk:"timeout_duration"`
	AuthorizedAppIds types.Set    `tfsdk:"authorized_app_ids"`
}
