resource "microsoft365_graph_beta_identity_and_access_authentication_strength" "auth_strength_with_configs" {
  # Display name must be 30 characters or less
  display_name = "unit-test-auth-w-configs"
  description  = "Unit test authentication strength with FIDO2 and X509 combination configurations"

  allowed_combinations = [
    "fido2",
    "x509CertificateMultiFactor"
  ]

  combination_configurations = [
    {
      odata_type              = "#microsoft.graph.fido2CombinationConfiguration"
      applies_to_combinations = ["fido2"]
      allowed_aaguids = [
        "90a3ccdf-635c-4729-a248-9b709135078f",  # YubiKey 5 Series
        "de1e552d-db1d-4423-a619-566b625cdc84"   # Feitian ePass FIDO
      ]
    },
    {
      odata_type              = "#microsoft.graph.x509CertificateCombinationConfiguration"
      applies_to_combinations = ["x509CertificateMultiFactor"]
      allowed_issuer_skis = [
        "1A2B3C4D5E6F7A8B9C0D1E2F3A4B5C6D7E8F9A0B",  # Corporate Root CA
        "9F8E7D6C5B4A3F2E1D0C9B8A7F6E5D4C3B2A1F0E"   # Backup Root CA
      ]
      allowed_issuers = [
        "CUSTOMIDENTIFIER:1A2B3C4D5E6F7A8B9C0D1E2F3A4B5C6D7E8F9A0B",
        "CUSTOMIDENTIFIER:9F8E7D6C5B4A3F2E1D0C9B8A7F6E5D4C3B2A1F0E"
      ]
      allowed_policy_oids = [
        "1.3.6.1.4.1.311.21.8.1.1",  # Microsoft Smart Card Logon
        "1.3.6.1.5.5.7.3.2"          # Client Authentication
      ]
    }
  ]
}

