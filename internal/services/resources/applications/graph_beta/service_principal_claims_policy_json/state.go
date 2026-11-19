package graphBetaServicePrincipalClaimsPolicyJSON

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func MapRemoteStateToTerraform(data *ServicePrincipalClaimsPolicyJSONResourceModel, content []byte) error {
	object, err := decodePolicyJSON(content)
	if err != nil {
		return err
	}
	data.PolicyID = types.StringNull()
	if id, ok := object["id"].(string); ok {
		data.PolicyID = types.StringValue(id)
	}
	canonical, err := normalisePolicyJSON(content)
	if err != nil {
		return err
	}
	data.ID = data.ServicePrincipalID
	data.PolicyJSON = NewPolicyJSONValue(string(canonical))
	return nil
}

func policyConsistencyPredicate(expected *ServicePrincipalClaimsPolicyJSONResourceModel) func(context.Context, tfsdk.State) bool {
	matches := 0
	return func(ctx context.Context, state tfsdk.State) bool {
		var actual ServicePrincipalClaimsPolicyJSONResourceModel
		if state.Get(ctx, &actual).HasError() {
			matches = 0
			return false
		}
		equal, _ := expected.PolicyJSON.StringSemanticEquals(ctx, actual.PolicyJSON)
		if !actual.ID.Equal(expected.ID) || !equal {
			matches = 0
			return false
		}
		matches++
		return matches >= 5
	}
}
