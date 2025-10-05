resource "microsoft365_graph_identity_and_access_conditional_access_terms_of_use" "maximal" {
  display_name                          = "Maximal Terms of Use Agreement"
  is_viewing_before_acceptance_required = true
  is_per_device_acceptance_required     = true
  user_reaccept_required_frequency      = "P90D"

  terms_expiration = {
    start_date_time = "2025-12-31"
    frequency       = "P365D"
  }

  file = {
    localizations = [
      {
        file_name        = "terms-en.pdf"
        display_name     = "Terms of Use - English"
        language         = "en-US"
        is_default       = true
        is_major_version = true
        file_data = {
          data = "%PDF-1.4\nComprehensive PDF content for maximal testing"
        }
      },
      {
        file_name        = "terms-fr.pdf"
        display_name     = "Terms of Use - French"
        language         = "fr-FR"
        is_default       = false
        is_major_version = false
        file_data = {
          data = "%PDF-1.4\nContenu PDF français pour les tests"
        }
      }
    ]
  }
}
