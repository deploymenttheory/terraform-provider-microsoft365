resource "microsoft365_graph_identity_and_access_conditional_access_terms_of_use" "test" {
  display_name                          = "Updated Acceptance Test Terms of Use"
  is_viewing_before_acceptance_required = false
  is_per_device_acceptance_required     = true
  user_reaccept_required_frequency      = "P180D"

  terms_expiration = {
    start_date_time = "2026-06-30T23:59:59.000Z"
    frequency       = "P180D"
  }

  file = {
    localizations = [
      {
        file_name    = "updated-terms.pdf"
        display_name = "Updated Test Terms of Use"
        language     = "en-US"
        is_default   = true
        file_data = {
          data = "%PDF-1.4\nUpdated PDF content for acceptance testing"
        }
      }
    ]
  }
}
