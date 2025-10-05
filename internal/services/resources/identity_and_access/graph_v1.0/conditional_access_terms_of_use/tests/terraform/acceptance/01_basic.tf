resource "microsoft365_graph_identity_and_access_conditional_access_terms_of_use" "test" {
  display_name                          = "Acceptance Test Terms of Use"
  is_viewing_before_acceptance_required = true
  is_per_device_acceptance_required     = false
  user_reaccept_required_frequency      = "P90D"

  terms_expiration = {
    start_date_time = "2025-12-31"
    frequency       = "P365D"
  }

  file = {
    localizations = [
      {
        file_name    = "test-terms.pdf"
        display_name = "Test Terms of Use"
        language     = "en-US"
        is_default   = true
        file_data = {
          data = "%PDF-1.4\nTest PDF content for acceptance testing"
        }
      }
    ]
  }
}
