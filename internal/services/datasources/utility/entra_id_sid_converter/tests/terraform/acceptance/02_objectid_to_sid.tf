# Acceptance Test: Convert Object ID to SID
data "microsoft365_utility_entra_id_sid_converter" "test" {
  object_id = "73d664e4-0886-4a73-b745-c694da45ddb4"
}

