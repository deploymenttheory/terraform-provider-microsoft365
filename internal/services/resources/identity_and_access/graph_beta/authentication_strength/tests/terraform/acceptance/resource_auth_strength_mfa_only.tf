resource "microsoft365_graph_beta_identity_and_access_authentication_strength" "auth_strength_mfa_only" {
  display_name = "acc-test-authentication-strength-mfa-only"
  description  = "Acceptance test MFA-only authentication strength policy"

  allowed_combinations = [
    "fido2",
    "windowsHelloForBusiness",
    "microsoftAuthenticatorPush,federatedSingleFactor",
    "x509CertificateMultiFactor"
  ]
}