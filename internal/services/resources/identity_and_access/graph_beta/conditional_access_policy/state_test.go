package graphBetaConditionalAccessPolicy

import (
	"context"
	"testing"

	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// A grant without controls can't be configured, so it must map to no block like a null one.
func TestMapRemoteResourceStateMapsGrantWithoutControlsToNoBlock(t *testing.T) {
	operator := "OR"
	empty := graphmodels.NewConditionalAccessGrantControls()
	empty.SetOperator(&operator)
	empty.SetBuiltInControls([]graphmodels.ConditionalAccessGrantControl{})
	empty.SetCustomAuthenticationFactors([]string{})
	empty.SetTermsOfUse([]string{})

	mfa := graphmodels.MFA_CONDITIONALACCESSGRANTCONTROL
	withControl := graphmodels.NewConditionalAccessGrantControls()
	withControl.SetOperator(&operator)
	withControl.SetBuiltInControls([]graphmodels.ConditionalAccessGrantControl{mfa})

	for name, tc := range map[string]struct {
		grantControls graphmodels.ConditionalAccessGrantControlsable
		wantBlock     bool
	}{
		"null":             {grantControls: nil},
		"without controls": {grantControls: empty},
		"with a control":   {grantControls: withControl, wantBlock: true},
	} {
		t.Run(name, func(t *testing.T) {
			policy := graphmodels.NewConditionalAccessPolicy()
			policy.SetGrantControls(tc.grantControls)

			state := MapRemoteResourceStateToTerraform(context.Background(), ConditionalAccessPolicyResourceModel{}, policy)
			if (state.GrantControls != nil) != tc.wantBlock {
				t.Fatalf("expected grant_controls block %v, got %+v", tc.wantBlock, state.GrantControls)
			}
		})
	}
}
