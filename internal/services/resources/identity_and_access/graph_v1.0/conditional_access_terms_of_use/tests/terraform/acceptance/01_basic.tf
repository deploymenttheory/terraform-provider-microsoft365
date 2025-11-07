resource "microsoft365_graph_identity_and_access_conditional_access_terms_of_use" "test" {
  display_name                          = "acc_test"
  is_viewing_before_acceptance_required = true
  is_per_device_acceptance_required     = false
  user_reaccept_required_frequency      = "P10D"

  terms_expiration = {
    start_date_time = "2025-11-06"
    frequency       = "P180D"
  }

  file = {
    localizations = [
      {
        file_name    = "test-terms.pdf"
        display_name = "Test Terms of Use"
        language     = "en-US"
        is_default   = true
        file_data = {
          data = "JVBERi0xLjQKVGVzdCBUZXJtcyBvZiBVc2UgRG9jdW1lbnQ="
        }
      },
      {
        file_name    = "english.pdf"
        display_name = "English Terms of Use"
        language     = "en"
        is_default   = false
        is_major_version = true
        file_data = {
          data = "JVBERi0xLjQKVGVzdCBUZXJtcyBvZiBVc2UgRG9jdW1lbnQ="
        }
      }
    ]
  }
}
