resource "microsoft365_graph_beta_identity_and_access_authentication_strength_policy" "auth_strength_example" {
  display_name = "example auth strength policy"
  description  = "test maximal authentication strength policy with all combinations and configurations"

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

  combination_configurations = [
    {
      applies_to_combinations = "fido2"
      allowed_aaguids = [
        "12345678-0000-0000-0000-123456780000",
        "90a3ccdf-635c-4729-a248-9b709135078f",
        "de1e552d-db1d-4423-a619-566b625cdc84"
      ]
    },
    {
      applies_to_combinations = "x509CertificateMultiFactor"
      allowed_issuer_skis     = ["1A2B3C4D5E6F7A8B9C0D1E2F3A4B5C6D7E8F9A0A", "1A2B3C4D5E6F7A8B9C0D1E2F3A4B5C6D7E8F9A0B"]
      allowed_policy_oids     = ["1.3.6.1.4.1.311.21.8.1.5", "1.2.3.4.5.8"]
    },
    {
      applies_to_combinations = "x509CertificateSingleFactor"
      allowed_issuer_skis     = ["1A2B3C4D5E6F7A8B9C0D1E2F3A4B5C6D7E8F9A0C", "1A2B3C4D5E6F7A8B9C0D1E2F3A4B5C6D7E8F9A0D"]
      allowed_policy_oids     = ["1.3.6.1.4.1.311.21.8.1.4", "1.2.3.4.5.8"]
    }
  ]
}