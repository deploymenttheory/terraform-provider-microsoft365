# Unit Test: Invalid RID exceeds uint32 maximum
data "microsoft365_utility_entra_id_sid_converter" "test" {
  sid = "S-1-12-1-1234567890-9876543210-1111111111-2222222222"
}

