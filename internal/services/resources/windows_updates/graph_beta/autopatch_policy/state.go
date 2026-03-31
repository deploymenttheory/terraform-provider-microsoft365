package graphBetaWindowsUpdatesAutopatchPolicy

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/sentinels"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodelswindowsupdates "github.com/microsoftgraph/msgraph-beta-sdk-go/models/windowsupdates"
)

// MapRemoteResourceStateToTerraform maps a remote policy resource to the Terraform state model.
func MapRemoteResourceStateToTerraform(ctx context.Context, data *WindowsUpdatesAutopatchPolicyResourceModel, remoteResource graphmodelswindowsupdates.Policyable) {
	if remoteResource == nil {
		tflog.Debug(ctx, "Remote resource is nil")
		return
	}

	tflog.Debug(ctx, "Starting to map remote state to Terraform state", map[string]any{
		"resourceId": convert.GraphToFrameworkString(remoteResource.GetId()).ValueString(),
	})

	data.ID = convert.GraphToFrameworkString(remoteResource.GetId())
	data.DisplayName = convert.GraphToFrameworkString(remoteResource.GetDisplayName())
	data.Description = convert.GraphToFrameworkString(remoteResource.GetDescription())
	data.CreatedDateTime = convert.GraphToFrameworkTime(remoteResource.GetCreatedDateTime())
	data.LastModifiedDateTime = convert.GraphToFrameworkTime(remoteResource.GetLastModifiedDateTime())

	approvalRulesSet, err := mapApprovalRulesToState(remoteResource.GetApprovalRules())
	if err != nil {
		tflog.Error(ctx, "Failed to map approval rules to state", map[string]any{
			"error": err.Error(),
		})
		data.ApprovalRules = types.SetNull(approvalRuleAttrType())
	} else {
		data.ApprovalRules = approvalRulesSet
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished mapping resource %s with id %s", ResourceName, data.ID.ValueString()))
}

// approvalRuleAttrType returns the object type for ApprovalRuleModel.
func approvalRuleAttrType() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"deferral_in_days": types.Int32Type,
			"classification":   types.StringType,
			"cadence":          types.StringType,
		},
	}
}

// mapApprovalRulesToState maps the SDK approval rules collection to a Terraform types.Set.
func mapApprovalRulesToState(rules []graphmodelswindowsupdates.ApprovalRuleable) (types.Set, error) {
	ruleType := approvalRuleAttrType()

	if len(rules) == 0 {
		return types.SetNull(ruleType), nil
	}

	ruleValues := make([]attr.Value, 0, len(rules))

	for _, rule := range rules {
		if rule == nil {
			continue
		}

		deferralInDays := convert.GraphToFrameworkInt32(rule.GetDeferralInDays())

		var classification types.String
		if qualRule, ok := rule.(graphmodelswindowsupdates.QualityUpdateApprovalRuleable); ok {
			classification = convert.GraphToFrameworkEnum(qualRule.GetClassification())
		} else {
			classification = types.StringNull()
		}

		var cadence types.String
		if qualRule, ok := rule.(graphmodelswindowsupdates.QualityUpdateApprovalRuleable); ok {
			cadence = convert.GraphToFrameworkEnum(qualRule.GetCadence())
		} else {
			cadence = types.StringNull()
		}

		ruleAttrs := map[string]attr.Value{
			"deferral_in_days": deferralInDays,
			"classification":   classification,
			"cadence":          cadence,
		}

		ruleValue, diags := types.ObjectValue(ruleType.(types.ObjectType).AttrTypes, ruleAttrs)
		if diags.HasError() {
			return types.SetNull(ruleType), sentinels.ErrCreateApprovalRuleObject
		}
		ruleValues = append(ruleValues, ruleValue)
	}

	set, diags := types.SetValue(ruleType, ruleValues)
	if diags.HasError() {
		return types.SetNull(ruleType), sentinels.ErrCreateApprovalRulesSet
	}

	return set, nil
}
