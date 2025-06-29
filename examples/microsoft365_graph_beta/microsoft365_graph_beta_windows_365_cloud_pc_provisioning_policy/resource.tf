resource "microsoft365_graph_beta_windows_365_cloud_pc_provisioning_policy" "example" {
  display_name              = "Example Cloud PC Provisioning Policy"
  description               = "Provisioning policy for Windows 365 Cloud PCs."
  cloud_pc_naming_template  = "CPC-%USERNAME:5%-%RAND:5%"
  image_id                  = "microsoftwindowsdesktop_windows-ent-cpc_win11-24H2-ent-cpc-m365"

  # Defaults: provisioning_type = "dedicated", managed_by = "windows365", image_type = "gallery", enable_single_sign_on = true

  domain_join_configurations = [
    {
      domain_join_type          = "azureADJoin"
      on_premises_connection_id = "00000000-0000-0000-0000-000000000000" # Replace with your connection ID
      region_name               = "automatic"
      region_group              = "unitedKingdom" # Valid values: default, australia, canada, usCentral, usEast, usWest, france, germany, europeUnion, unitedKingdom, japan, asia, india, southAmerica, euap, usGovernment, usGovernmentDOD, unknownFutureValue, norway, switzerland, southKorea, middleEast, mexico, australasia, europe
    }
  ]

  windows_setting {
    locale = "en-US"
  }

  microsoft_managed_desktop {
    managed_type = "starterManaged"
    # profile    = "" # Optional
  }

  autopatch {
    autopatch_group_id = "11111111-1111-1111-1111-111111111111" # Replace with your autopatch group ID
  }

  autopilot_configuration {
    device_preparation_profile_id    = "22222222-2222-2222-2222-222222222222" # Replace with your profile ID
    application_timeout_in_minutes   = 60
    on_failure_device_access_denied  = false
  }

  scope_ids = ["0", "8"]
}

resource "microsoft365_graph_beta_windows_365_cloud_pc_provisioning_policy" "entra_id" {
  display_name              = "test"
  description               = "test"
  image_id                  = "microsoftwindowsdesktop_windows-ent-cpc_win11-24H2-ent-cpc-m365"

  domain_join_configurations = [
    {
      domain_join_type = "azureADJoin"
      region_name      = "japaneast"
      region_group     = "japan" # Valid values: default, australia, canada, usCentral, usEast, usWest, france, germany, europeUnion, unitedKingdom, japan, asia, india, southAmerica, euap, usGovernment, usGovernmentDOD, unknownFutureValue, norway, switzerland, southKorea, middleEast, mexico, australasia, europe
    }
  ]

  windows_setting {
    locale = "ja-JP"
  }

  microsoft_managed_desktop {
    managed_type = "starterManaged"
    # profile    = null # Omitted as null
  }

  autopatch {
    autopatch_group_id = "4aa9b805-9494-4eed-a04b-ed51ec9e631e"
  }

  scope_ids = ["0"]
  # cloud_pc_naming_template = null # Omitted as null
} 