# Unit Test: Invalid GUID format
data "microsoft365_utility_entra_id_sid_converter" "test" {
  object_id = "not-a-valid-guid"
}

