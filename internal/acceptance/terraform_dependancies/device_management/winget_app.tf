# WinGet App dependencies for acceptance testing
# These apps serve as dependencies for Windows Autopilot Device Preparation policies

resource "random_string" "winget_app_suffix" {
  length  = 8
  special = false
  upper   = false
}

# Application category resource for reference
resource "microsoft365_graph_beta_device_and_app_management_application_category" "test_productivity" {
  display_name = "acc-test-productivity-${random_string.winget_app_suffix.result}"
}

# Test WinGet App 1 - Firefox
resource "microsoft365_graph_beta_device_and_app_management_win_get_app" "test_firefox" {
  package_identifier              = "9NZVDKPMR9RD" # Firefox from Microsoft Store
  automatically_generate_metadata = true

  # Optional app information
  is_featured             = true
  privacy_information_url = "https://www.mozilla.org/en-US/privacy/firefox/"
  information_url         = "https://support.mozilla.org/en-US/"
  owner                   = "Acceptance Test Suite"
  developer               = "Mozilla Foundation"
  notes                   = "Test browser for autopilot device preparation acceptance tests"

  # Required install experience settings
  install_experience = {
    run_as_account = "user"
  }

  # Optional role scope tag IDs
  role_scope_tag_ids = ["0"]

  categories = [
    microsoft365_graph_beta_device_and_app_management_application_category.test_productivity.id,
    "Business",
    "Productivity",
  ]

  # Optional timeouts
  timeouts = {
    create = "5m"
    update = "5m"
    read   = "2m"
    delete = "5m"
  }
}

# Test WinGet App 2 - Alternative app for testing
resource "microsoft365_graph_beta_device_and_app_management_win_get_app" "test_notepadplusplus" {
  package_identifier              = "9MSMLRH6LZF3" # Notepad++ from Microsoft Store
  automatically_generate_metadata = true

  # Optional app information
  is_featured = false
  owner       = "Acceptance Test Suite"
  developer   = "Notepad++ Team"
  notes       = "Test text editor for autopilot device preparation acceptance tests"

  # Required install experience settings
  install_experience = {
    run_as_account = "system"
  }

  # Optional role scope tag IDs
  role_scope_tag_ids = ["0"]

  categories = [
    microsoft365_graph_beta_device_and_app_management_application_category.test_productivity.id,
    "Business",
  ]

  # Optional timeouts
  timeouts = {
    create = "5m"
    update = "5m"
    read   = "2m"
    delete = "5m"
  }
}

# Outputs for easy reference in tests
output "test_winget_app_ids" {
  description = "WinGet app IDs for use in Windows Autopilot Device Preparation policy tests"
  value = {
    firefox         = microsoft365_graph_beta_device_and_app_management_win_get_app.test_firefox.id
    notepadplusplus = microsoft365_graph_beta_device_and_app_management_win_get_app.test_notepadplusplus.id
  }
}