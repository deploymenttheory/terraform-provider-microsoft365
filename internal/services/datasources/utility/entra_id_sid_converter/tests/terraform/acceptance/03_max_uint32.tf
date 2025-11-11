# Acceptance Test: Maximum uint32 values
data "microsoft365_utility_entra_id_sid_converter" "test" {
  sid = "S-1-12-1-4294967295-4294967295-4294967295-4294967295"
}

