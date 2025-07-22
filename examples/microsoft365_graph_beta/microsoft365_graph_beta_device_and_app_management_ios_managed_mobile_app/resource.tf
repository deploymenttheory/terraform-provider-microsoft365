########################################################################################
# iOS Managed Mobile App Examples
########################################################################################

# Resource for adding Microsoft Outlook to an iOS managed app protection policy
resource "microsoft365_graph_beta_device_and_app_management_ios_managed_mobile_app" "outlook" {
  managed_app_protection_id = "00000000-0000-0000-0000-000000000001" # Replace with your iOS managed app protection policy ID
  
  mobile_app_identifier = {
    bundle_id = "com.microsoft.Office.Outlook"
  }

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}

# Resource for adding Microsoft Teams to an iOS managed app protection policy
resource "microsoft365_graph_beta_device_and_app_management_ios_managed_mobile_app" "teams" {
  managed_app_protection_id = "00000000-0000-0000-0000-000000000001" # Replace with your iOS managed app protection policy ID
  
  mobile_app_identifier = {
    bundle_id = "com.microsoft.teamsmobile"
  }

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}

# Resource for adding Microsoft OneDrive to an iOS managed app protection policy
resource "microsoft365_graph_beta_device_and_app_management_ios_managed_mobile_app" "onedrive" {
  managed_app_protection_id = "00000000-0000-0000-0000-000000000001" # Replace with your iOS managed app protection policy ID
  
  mobile_app_identifier = {
    bundle_id = "com.microsoft.OneDrive"
  }

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}

# Resource for adding Microsoft Word to an iOS managed app protection policy
resource "microsoft365_graph_beta_device_and_app_management_ios_managed_mobile_app" "word" {
  managed_app_protection_id = "00000000-0000-0000-0000-000000000001" # Replace with your iOS managed app protection policy ID
  
  mobile_app_identifier = {
    bundle_id = "com.microsoft.Office.Word"
  }

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}

# Resource for adding Microsoft Excel to an iOS managed app protection policy
resource "microsoft365_graph_beta_device_and_app_management_ios_managed_mobile_app" "excel" {
  managed_app_protection_id = "00000000-0000-0000-0000-000000000001" # Replace with your iOS managed app protection policy ID
  
  mobile_app_identifier = {
    bundle_id = "com.microsoft.Office.Excel"
  }

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
} 