package graphBetaOperationApprovalPolicy

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// constructResource constructs and returns an OperationApprovalPolicy
func constructResource(ctx context.Context, data OperationApprovalPolicyResourceModel) (graphmodels.OperationApprovalPolicyable, error) {
	tflog.Debug(ctx, "Starting operation approval policy construction")

	policy := graphmodels.NewOperationApprovalPolicy()

	constructors.SetStringProperty(data.DisplayName, policy.SetDisplayName)
	constructors.SetStringProperty(data.Description, policy.SetDescription)

	if err := constructors.SetEnumProperty(data.PolicyType,
		func(s string) (any, error) { return graphmodels.ParseOperationApprovalPolicyType(s) },
		policy.SetPolicyType); err != nil {
		return nil, fmt.Errorf("error setting policy type: %v", err)
	}

	if err := constructors.SetEnumProperty(data.PolicyPlatform,
		func(s string) (any, error) { return graphmodels.ParseOperationApprovalPolicyPlatform(s) },
		policy.SetPolicyPlatform); err != nil {
		return nil, fmt.Errorf("error setting policy platform: %v", err)
	}

	policySet, err := constructPolicySet(ctx, &data.PolicySet)
	if err != nil {
		return nil, fmt.Errorf("error constructing policy set: %v", err)
	}
	policy.SetPolicySet(policySet)

	if err := constructors.SetStringSet(ctx, data.ApproverGroupIds, policy.SetApproverGroupIds); err != nil {
		return nil, fmt.Errorf("error setting approver group IDs: %v", err)
	}

	if err := constructors.DebugLogGraphObject(ctx, "Constructed operation approval policy", policy); err != nil {
		tflog.Error(ctx, "Failed to log operation approval policy", map[string]interface{}{
			"error": err.Error(),
		})
	}

	return policy, nil
}

// constructPolicySet constructs the operation approval policy set
func constructPolicySet(ctx context.Context, data *OperationApprovalPolicySetResourceModel) (graphmodels.OperationApprovalPolicySetable, error) {
	if data == nil {
		return nil, fmt.Errorf("policy set data is required")
	}

	policySet := graphmodels.NewOperationApprovalPolicySet()

	if err := constructors.SetEnumProperty(data.PolicyType,
		func(s string) (any, error) { return graphmodels.ParseOperationApprovalPolicyType(s) },
		policySet.SetPolicyType); err != nil {
		return nil, fmt.Errorf("error setting policy set policy type: %v", err)
	}

	if err := constructors.SetEnumProperty(data.PolicyPlatform,
		func(s string) (any, error) { return graphmodels.ParseOperationApprovalPolicyPlatform(s) },
		policySet.SetPolicyPlatform); err != nil {
		return nil, fmt.Errorf("error setting policy set policy platform: %v", err)
	}

	tflog.Debug(ctx, "Finished constructing policy set")
	return policySet, nil
}
