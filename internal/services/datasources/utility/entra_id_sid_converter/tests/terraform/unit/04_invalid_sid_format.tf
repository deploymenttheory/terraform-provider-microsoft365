# Unit Test: Invalid SID format (on-prem AD SID instead of Entra ID SID)
data "microsoft365_utility_entra_id_sid_converter" "test" {
  sid = "S-1-5-21-1234567890-1234567891-1234567892-1234567893"
}

