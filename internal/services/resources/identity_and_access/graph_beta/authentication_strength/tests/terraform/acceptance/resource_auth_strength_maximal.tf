resource "microsoft365_graph_beta_identity_and_access_authentication_strength" "auth_strength_maximal" {
  display_name = "acc-test-authentication-strength-maximal"
  description  = "Acceptance test maximal authentication strength policy with all combinations"

  allowed_combinations = [
    "deviceBasedPush",
    "federatedMultiFactor",
    "federatedSingleFactor",
    "fido2",
    "hardwareOath,federatedSingleFactor",
    "microsoftAuthenticatorPush,federatedSingleFactor",
    "password",
    "password,hardwareOath",
    "password,microsoftAuthenticatorPush",
    "password,sms",
    "password,softwareOath",
    "password,voice",
    "qrCodePin",
    "sms",
    "sms,federatedSingleFactor",
    "softwareOath,federatedSingleFactor",
    "temporaryAccessPassMultiUse",
    "temporaryAccessPassOneTime",
    "voice,federatedSingleFactor",
    "windowsHelloForBusiness",
    "x509CertificateMultiFactor",
    "x509CertificateSingleFactor"
  ]
}