package graphBetaWindowsUpdatesAutopatchPolicyApproval

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodelswindowsupdates "github.com/microsoftgraph/msgraph-beta-sdk-go/models/windowsupdates"
)

func constructResource(ctx context.Context, data *WindowsUpdatesAutopatchPolicyApprovalResourceModel) (graphmodelswindowsupdates.PolicyApprovalable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))

	requestBody := graphmodelswindowsupdates.NewPolicyApproval()

	catalogEntryId := data.CatalogEntryId.ValueString()
	requestBody.SetCatalogEntryId(&catalogEntryId)

	status, err := graphmodelswindowsupdates.ParseApprovalStatus(data.Status.ValueString())
	if err != nil || status == nil {
		return nil, fmt.Errorf("invalid approval status %q: %w", data.Status.ValueString(), err)
	}
	approvalStatus := status.(*graphmodelswindowsupdates.ApprovalStatus)
	requestBody.SetStatus(approvalStatus)

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing deployment settings for %s resource", ResourceName))
	return requestBody, nil
}

// constructUpdateResource builds a PATCH body containing only the mutable field (status).
// catalogEntryId is not patchable after creation.
func constructUpdateResource(ctx context.Context, data *WindowsUpdatesAutopatchPolicyApprovalResourceModel) (graphmodelswindowsupdates.PolicyApprovalable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing update request for %s resource", ResourceName))

	requestBody := graphmodelswindowsupdates.NewPolicyApproval()

	status, err := graphmodelswindowsupdates.ParseApprovalStatus(data.Status.ValueString())
	if err != nil || status == nil {
		return nil, fmt.Errorf("invalid approval status %q: %w", data.Status.ValueString(), err)
	}
	approvalStatus := status.(*graphmodelswindowsupdates.ApprovalStatus)
	requestBody.SetStatus(approvalStatus)

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing update request for %s resource", ResourceName))
	return requestBody, nil
}
