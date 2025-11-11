# Unit Test: Both sid and object_id provided (should fail)
data "microsoft365_utility_entra_id_sid_converter" "test" {
  sid       = "S-1-12-1-1943430372-1249052806-2496021943-3034400218"
  object_id = "73d664e4-0886-4a73-b745-c694da45ddb4"
}

