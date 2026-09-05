package graphBetaNetworkPromptPolicy

import "fmt"

// Verify readback because successful portal-backed PATCH responses can ignore fields.
func verifyObserved(plan, observed NetworkPromptPolicyResourceModel) error {
	if !plan.Name.IsUnknown() && !plan.Name.Equal(observed.Name) {
		return fmt.Errorf("%w: the API did not retain the requested Name value; inspect the observed state before retrying", errInvalidResponse)
	}
	if !plan.Description.IsUnknown() && !plan.Description.Equal(observed.Description) {
		return fmt.Errorf("%w: the API did not retain the requested Description value; inspect the observed state before retrying", errInvalidResponse)
	}
	if !plan.DefaultAction.IsUnknown() && !plan.DefaultAction.Equal(observed.DefaultAction) {
		return fmt.Errorf("%w: the API did not retain the requested DefaultAction value; inspect the observed state before retrying", errInvalidResponse)
	}
	return nil
}
