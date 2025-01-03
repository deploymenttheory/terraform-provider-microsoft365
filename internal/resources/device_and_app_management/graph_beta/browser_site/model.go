// REF: https://learn.microsoft.com/en-us/graph/api/resources/browsersite?view=graph-rest-beta
package graphBetaBrowserSite

import (
	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/shared_models/graph_beta"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type BrowserSiteResourceModel struct {
	ID                          types.String                          `tfsdk:"id"`
	BrowserSiteListAssignmentID types.String                          `tfsdk:"browser_site_list_assignment_id"`
	AllowRedirect               types.Bool                            `tfsdk:"allow_redirect"`
	Comment                     types.String                          `tfsdk:"comment"`
	CompatibilityMode           types.String                          `tfsdk:"compatibility_mode"`
	CreatedDateTime             types.String                          `tfsdk:"created_date_time"`
	DeletedDateTime             types.String                          `tfsdk:"deleted_date_time"`
	History                     []BrowserSiteHistoryResourceModel     `tfsdk:"history"`
	LastModifiedBy              sharedmodels.IdentitySetResourceModel `tfsdk:"last_modified_by"`
	LastModifiedDateTime        types.String                          `tfsdk:"last_modified_date_time"`
	MergeType                   types.String                          `tfsdk:"merge_type"`
	Status                      types.String                          `tfsdk:"status"`
	TargetEnvironment           types.String                          `tfsdk:"target_environment"`
	WebUrl                      types.String                          `tfsdk:"web_url"`
	Timeouts                    timeouts.Value                        `tfsdk:"timeouts"`
}

type BrowserSiteHistoryResourceModel struct {
	AllowRedirect     types.Bool                            `tfsdk:"allow_redirect"`
	Comment           types.String                          `tfsdk:"comment"`
	CompatibilityMode types.String                          `tfsdk:"compatibility_mode"`
	LastModifiedBy    sharedmodels.IdentitySetResourceModel `tfsdk:"last_modified_by"`
	MergeType         types.String                          `tfsdk:"merge_type"`
	PublishedDateTime types.String                          `tfsdk:"published_date_time"`
	TargetEnvironment types.String                          `tfsdk:"target_environment"`
}
