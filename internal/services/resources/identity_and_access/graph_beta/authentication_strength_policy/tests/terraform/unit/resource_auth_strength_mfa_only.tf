resource "microsoft365_graph_beta_identity_and_access_authentication_strength" "auth_strength_mfa_only" {
  # Display name must be 30 characters or less
  display_name = "unit-test-auth-strength-mfa"
  description  = "Unit test MFA-only authentication strength policy"

  allowed_combinations = [
    "fido2",
    "windowsHelloForBusiness",
    "microsoftAuthenticatorPush,federatedSingleFactor",
    "x509CertificateMultiFactor"
  ]
}