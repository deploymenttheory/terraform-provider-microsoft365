# Acceptance Test: Bidirectional conversion (round trip)
data "microsoft365_utility_entra_id_sid_converter" "sid_to_oid" {
  sid = "S-1-12-1-1943430372-1249052806-2496021943-3034400218"
}

data "microsoft365_utility_entra_id_sid_converter" "oid_to_sid" {
  object_id = "73d664e4-0886-4a73-b745-c694da45ddb4"
}

