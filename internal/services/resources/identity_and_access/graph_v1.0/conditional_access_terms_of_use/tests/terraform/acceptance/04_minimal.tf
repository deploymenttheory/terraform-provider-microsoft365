resource "microsoft365_graph_identity_and_access_conditional_access_terms_of_use" "test" {
  display_name = "Minimal Terms of Use"

  file = {
    localizations = [
      {
        file_name    = "minimal-terms.pdf"
        display_name = "Minimal Terms"
        language     = "en-US"
        is_default   = true
        file_data = {
          data = "%PDF-1.4\nMinimal PDF content"
        }
      }
    ]
  }
}
