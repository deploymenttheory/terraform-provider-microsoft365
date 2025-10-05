resource "microsoft365_graph_identity_and_access_conditional_access_terms_of_use" "minimal" {
  display_name                          = "Minimal Terms of Use"
  is_viewing_before_acceptance_required = false
  is_per_device_acceptance_required     = false

  file = {
    localizations = [
      {
        file_name        = "minimal-terms.pdf"
        display_name     = "Minimal Terms"
        language         = "en-US"
        is_default       = true
        is_major_version = false
        file_data = {
          data = "%PDF-1.4\nMinimal PDF content"
        }
      }
    ]
  }
}
