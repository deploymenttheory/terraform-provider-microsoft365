# Secure account recovery with identity verification (Preview)
# Template discovery is independent of policy creation and its tenant prerequisites.
data "microsoft365_graph_beta_identity_and_access_conditional_access_template" "secure_account_recovery" {
  template_id = "6acdf4c3-6815-485c-a57d-2c349d517ba0"
}

output "secure_account_recovery_details" {
  value = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.secure_account_recovery.details
}
