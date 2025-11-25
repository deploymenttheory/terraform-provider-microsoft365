// REF: https://learn.microsoft.com/en-us/graph/api/resources/grouplifecyclepolicy?view=graph-rest-beta
package graphBetaGroupLifecycleExpirationPolicy

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type GroupLifecycleExpirationPolicyResourceModel struct {
	ID                          types.String   `tfsdk:"id"`
	AlternateNotificationEmails types.String   `tfsdk:"alternate_notification_emails"`
	GroupLifetimeInDays         types.Int32    `tfsdk:"group_lifetime_in_days"`
	ManagedGroupTypes           types.String   `tfsdk:"managed_group_types"`
	OverwriteExistingPolicy     types.Bool     `tfsdk:"overwrite_existing_policy"`
	Timeouts                    timeouts.Value `tfsdk:"timeouts"`
}
