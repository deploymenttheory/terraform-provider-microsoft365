package graphBetaWindowsUpdatesAutopatchPolicy

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodelswindowsupdates "github.com/microsoftgraph/msgraph-beta-sdk-go/models/windowsupdates"
)

// constructResource builds the qualityUpdatePolicy request body from the Terraform model.
func constructResource(ctx context.Context, data *WindowsUpdatesAutopatchPolicyResourceModel) (graphmodelswindowsupdates.Policyable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))

	requestBody := graphmodelswindowsupdates.NewQualityUpdatePolicy()

	convert.FrameworkToGraphString(data.DisplayName, requestBody.SetDisplayName)
	convert.FrameworkToGraphString(data.Description, requestBody.SetDescription)

	if !data.ApprovalRules.IsNull() && !data.ApprovalRules.IsUnknown() && len(data.ApprovalRules.Elements()) > 0 {
		var ruleModels []ApprovalRuleModel
		diags := data.ApprovalRules.ElementsAs(ctx, &ruleModels, false)
		if diags.HasError() {
			return nil, fmt.Errorf("failed to extract approval_rules: %s", diags.Errors())
		}

		approvalRules := make([]graphmodelswindowsupdates.ApprovalRuleable, 0, len(ruleModels))
		for _, rm := range ruleModels {
			rule := graphmodelswindowsupdates.NewQualityUpdateApprovalRule()

			deferralInDays := rm.DeferralInDays.ValueInt32()
			rule.SetDeferralInDays(&deferralInDays)

			if err := convert.FrameworkToGraphEnum(rm.Classification, graphmodelswindowsupdates.ParseQualityUpdateClassification, rule.SetClassification); err != nil {
				return nil, fmt.Errorf("failed to set classification: %v", err)
			}

			if err := convert.FrameworkToGraphEnum(rm.Cadence, graphmodelswindowsupdates.ParseQualityUpdateCadence, rule.SetCadence); err != nil {
				return nil, fmt.Errorf("failed to set cadence: %v", err)
			}

			approvalRules = append(approvalRules, rule)
		}

		requestBody.SetApprovalRules(approvalRules)
	}

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s resource", ResourceName))

	return requestBody, nil
}
