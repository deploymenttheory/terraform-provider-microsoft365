package graphBetaNetworkPromptPolicyRule

import "fmt"

// Verify readback because successful portal-backed PATCH responses can ignore fields.
func verifyObserved(plan, observed NetworkPromptPolicyRuleResourceModel) error {
	if !plan.Name.IsUnknown() && !plan.Name.Equal(observed.Name) {
		return fmt.Errorf(
			"%w: the API did not retain the requested Name value; inspect the observed state before retrying",
			errInvalidResponse,
		)
	}
	if !plan.Description.IsUnknown() && !plan.Description.Equal(observed.Description) {
		return fmt.Errorf(
			"%w: the API did not retain the requested Description value; inspect the observed state before retrying",
			errInvalidResponse,
		)
	}
	if !plan.Action.IsUnknown() && !plan.Action.Equal(observed.Action) {
		return fmt.Errorf(
			"%w: the API did not retain the requested Action value; inspect the observed state before retrying",
			errInvalidResponse,
		)
	}
	if !plan.Priority.IsUnknown() && !plan.Priority.Equal(observed.Priority) {
		return fmt.Errorf(
			"%w: the API did not retain the requested Priority value; inspect the observed state before retrying",
			errInvalidResponse,
		)
	}
	if !plan.Enabled.IsUnknown() && !plan.Enabled.Equal(observed.Enabled) {
		return fmt.Errorf(
			"%w: the API did not retain the requested Enabled value; inspect the observed state before retrying",
			errInvalidResponse,
		)
	}
	if !plan.PromptLogging.IsUnknown() && !plan.PromptLogging.Equal(observed.PromptLogging) {
		return fmt.Errorf(
			"%w: the API did not retain the requested PromptLogging value; inspect the observed state before retrying",
			errInvalidResponse,
		)
	}
	if !plan.ScanResult.IsUnknown() && !plan.ScanResult.Equal(observed.ScanResult) {
		return fmt.Errorf(
			"%w: the API did not retain the requested ScanResult value; inspect the observed state before retrying",
			errInvalidResponse,
		)
	}
	if !plan.ConversationSchemes.IsUnknown() &&
		!plan.ConversationSchemes.Equal(observed.ConversationSchemes) {
		return fmt.Errorf(
			"%w: the API did not retain the requested ConversationSchemes value; inspect the observed state before retrying",
			errInvalidResponse,
		)
	}
	return nil
}
