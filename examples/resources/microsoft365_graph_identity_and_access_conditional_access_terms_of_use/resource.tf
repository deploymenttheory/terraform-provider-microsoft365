resource "microsoft365_graph_identity_and_access_conditional_access_terms_of_use" "example" {
  display_name                          = "Example Conditional Access Terms of Use"
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
        file_name    = "test.pdf"
        display_name = "Terms of Use - English (US)"
        language     = "en-US"
        is_default   = true
        file_data = {
          data = filebase64("${path.module}/test.pdf")
        }
      },
      {
        file_name        = "test.pdf"
        display_name     = "Terms of Use - English"
        language         = "en"
        is_default       = false
        is_major_version = false
        file_data = {
          data = filebase64("${path.module}/test.pdf")
        }
      },
      {
        file_name        = "test.pdf"
        display_name     = "Terms of Use - English (UK)"
        language         = "en-GB"
        is_default       = false
        is_major_version = false
        file_data = {
          data = filebase64("${path.module}/test.pdf")
        }
      },
      {
        file_name        = "test.pdf"
        display_name     = "Terms of Use - Afrikaans"
        language         = "af"
        is_default       = false
        is_major_version = false
        file_data = {
          data = filebase64("${path.module}/test.pdf")
        }
      },
      {
        file_name        = "test.pdf"
        display_name     = "Terms of Use - Amharic"
        language         = "am"
        is_default       = false
        is_major_version = false
        file_data = {
          data = filebase64("${path.module}/test.pdf")
        }
      },
      {
        file_name        = "test.pdf"
        display_name     = "Terms of Use - Arabic (Saudi Arabia)"
        language         = "ar-SA"
        is_default       = false
        is_major_version = false
        file_data = {
          data = filebase64("${path.module}/test.pdf")
        }
      },
      {
        file_name        = "test.pdf"
        display_name     = "Terms of Use - Armenian"
        language         = "hy"
        is_default       = false
        is_major_version = false
        file_data = {
          data = filebase64("${path.module}/test.pdf")
        }
      },
      {
        file_name        = "test.pdf"
        display_name     = "Terms of Use - Assamese"
        language         = "as"
        is_default       = false
        is_major_version = false
        file_data = {
          data = filebase64("${path.module}/test.pdf")
        }
      },
      {
        file_name        = "test.pdf"
        display_name     = "Terms of Use - Azerbaijani"
        language         = "az"
        is_default       = false
        is_major_version = false
        file_data = {
          data = filebase64("${path.module}/test.pdf")
        }
      },
      {
        file_name        = "test.pdf"
        display_name     = "Terms of Use - Belarusian"
        language         = "be"
        is_default       = false
        is_major_version = false
        file_data = {
          data = filebase64("${path.module}/test.pdf")
        }
      },
      {
        file_name        = "test.pdf"
        display_name     = "Terms of Use - Bangla"
        language         = "bn"
        is_default       = false
        is_major_version = false
        file_data = {
          data = filebase64("${path.module}/test.pdf")
        }
      },
      {
        file_name        = "test.pdf"
        display_name     = "Terms of Use - Bangla (India)"
        language         = "bn-IN"
        is_default       = false
        is_major_version = false
        file_data = {
          data = filebase64("${path.module}/test.pdf")
        }
      },
      {
        file_name        = "test.pdf"
        display_name     = "Terms of Use - Basque"
        language         = "eu"
        is_default       = false
        is_major_version = false
        file_data = {
          data = filebase64("${path.module}/test.pdf")
        }
      },
      {
        file_name        = "test.pdf"
        display_name     = "Terms of Use - Central Kurdish"
        language         = "ku-Arab"
        is_default       = false
        is_major_version = false
        file_data = {
          data = filebase64("${path.module}/test.pdf")
        }
      },
      {
        file_name        = "test.pdf"
        display_name     = "Terms of Use - Japanese"
        language         = "ja"
        is_default       = false
        is_major_version = false
        file_data = {
          data = filebase64("${path.module}/test.pdf")
        }
      },
      {
        file_name        = "test.pdf"
        display_name     = "Terms of Use - Zulu"
        language         = "zu"
        is_default       = false
        is_major_version = false
        file_data = {
          data = filebase64("${path.module}/test.pdf")
        }
      },
      {
        file_name        = "test.pdf"
        display_name     = "Terms of Use - Telugu"
        language         = "te"
        is_default       = false
        is_major_version = false
        file_data = {
          data = filebase64("${path.module}/test.pdf")
        }
      },
      {
        file_name        = "test.pdf"
        display_name     = "Terms of Use - Thai"
        language         = "th"
        is_default       = false
        is_major_version = false
        file_data = {
          data = filebase64("${path.module}/test.pdf")
        }
      },
      {
        file_name        = "test.pdf"
        display_name     = "Terms of Use - Tigrinya"
        language         = "ti"
        is_default       = false
        is_major_version = false
        file_data = {
          data = filebase64("${path.module}/test.pdf")
        }
      },
      {
        file_name        = "test.pdf"
        display_name     = "Terms of Use - Turkish"
        language         = "tr"
        is_default       = false
        is_major_version = false
        file_data = {
          data = filebase64("${path.module}/test.pdf")
        }
      },
      {
        file_name        = "test.pdf"
        display_name     = "Terms of Use - Turkmen"
        language         = "tk"
        is_default       = false
        is_major_version = false
        file_data = {
          data = filebase64("${path.module}/test.pdf")
        }
      },
      {
        file_name        = "test.pdf"
        display_name     = "Terms of Use - Ukrainian"
        language         = "uk"
        is_default       = false
        is_major_version = false
        file_data = {
          data = filebase64("${path.module}/test.pdf")
        }
      },
      {
        file_name        = "test.pdf"
        display_name     = "Terms of Use - Urdu"
        language         = "ur"
        is_default       = false
        is_major_version = false
        file_data = {
          data = filebase64("${path.module}/test.pdf")
        }
      },
      {
        file_name        = "test.pdf"
        display_name     = "Terms of Use - Uyghur"
        language         = "ug"
        is_default       = false
        is_major_version = false
        file_data = {
          data = filebase64("${path.module}/test.pdf")
        }
      },
      {
        file_name        = "test.pdf"
        display_name     = "Terms of Use - Uzbek"
        language         = "uz"
        is_default       = false
        is_major_version = false
        file_data = {
          data = filebase64("${path.module}/test.pdf")
        }
      },
      {
        file_name        = "test.pdf"
        display_name     = "Terms of Use - Valencian (Spain)"
        language         = "ca-ES-valencia"
        is_default       = false
        is_major_version = false
        file_data = {
          data = filebase64("${path.module}/test.pdf")
        }
      },
      {
        file_name        = "test.pdf"
        display_name     = "Terms of Use - Vietnamese"
        language         = "vi"
        is_default       = false
        is_major_version = false
        file_data = {
          data = filebase64("${path.module}/test.pdf")
        }
      },
      {
        file_name        = "test.pdf"
        display_name     = "Terms of Use - Welsh"
        language         = "cy"
        is_default       = false
        is_major_version = false
        file_data = {
          data = filebase64("${path.module}/test.pdf")
        }
      },
      {
        file_name        = "test.pdf"
        display_name     = "Terms of Use - Wolof"
        language         = "wo"
        is_default       = false
        is_major_version = false
        file_data = {
          data = filebase64("${path.module}/test.pdf")
        }
      },
      {
        file_name        = "test.pdf"
        display_name     = "Terms of Use - Yoruba"
        language         = "yo"
        is_default       = false
        is_major_version = false
        file_data = {
          data = filebase64("${path.module}/test.pdf")
        }
      }
    ]
  }
}
