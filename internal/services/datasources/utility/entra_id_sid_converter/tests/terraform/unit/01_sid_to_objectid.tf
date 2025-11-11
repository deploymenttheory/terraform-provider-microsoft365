# Unit Test: Convert SID to Object ID
data "microsoft365_utility_entra_id_sid_converter" "test" {
  sid = "S-1-12-1-1943430372-1249052806-2496021943-3034400218"
}

