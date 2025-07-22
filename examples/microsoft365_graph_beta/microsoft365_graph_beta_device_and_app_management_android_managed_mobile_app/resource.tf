########################################################################################
# Android Managed Mobile App Examples
########################################################################################

# Resource for adding Microsoft Outlook to an Android managed app protection policy
resource "microsoft365_graph_beta_device_and_app_management_android_managed_mobile_app" "outlook" {
  managed_app_protection_id = "00000000-0000-0000-0000-000000000001" # Replace with your Android managed app protection policy ID
  
  mobile_app_identifier = {
    package_id = "com.microsoft.office.outlook"
  }

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}

# Resource for adding Microsoft Teams to an Android managed app protection policy
resource "microsoft365_graph_beta_device_and_app_management_android_managed_mobile_app" "teams" {
  managed_app_protection_id = "00000000-0000-0000-0000-000000000001" # Replace with your Android managed app protection policy ID
  
  mobile_app_identifier = {
    package_id = "com.microsoft.teams"
  }

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}

# Resource for adding Microsoft OneDrive to an Android managed app protection policy
resource "microsoft365_graph_beta_device_and_app_management_android_managed_mobile_app" "onedrive" {
  managed_app_protection_id = "00000000-0000-0000-0000-000000000001" # Replace with your Android managed app protection policy ID
  
  mobile_app_identifier = {
    package_id = "com.microsoft.skydrive"
  }

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}

# Resource for adding Microsoft Word to an Android managed app protection policy
resource "microsoft365_graph_beta_device_and_app_management_android_managed_mobile_app" "word" {
  managed_app_protection_id = "00000000-0000-0000-0000-000000000001" # Replace with your Android managed app protection policy ID
  
  mobile_app_identifier = {
    package_id = "com.microsoft.office.word"
  }

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}

# Resource for adding Microsoft Excel to an Android managed app protection policy
resource "microsoft365_graph_beta_device_and_app_management_android_managed_mobile_app" "excel" {
  managed_app_protection_id = "00000000-0000-0000-0000-000000000001" # Replace with your Android managed app protection policy ID
  
  mobile_app_identifier = {
    package_id = "com.microsoft.office.excel"
  }

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
} 