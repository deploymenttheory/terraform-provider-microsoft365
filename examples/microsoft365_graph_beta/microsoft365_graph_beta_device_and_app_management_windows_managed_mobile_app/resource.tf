########################################################################################
# Windows Managed Mobile App Examples
########################################################################################

# Resource for adding Microsoft Outlook to a Windows managed app protection policy
resource "microsoft365_graph_beta_device_and_app_management_windows_managed_mobile_app" "outlook" {
  managed_app_protection_id = "00000000-0000-0000-0000-000000000001" # Replace with your Windows managed app protection policy ID

  mobile_app_identifier = {
    windows_app_id = "Microsoft.Office.Outlook_8wekyb3d8bbwe"
  }

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}

# Resource for adding Microsoft Teams to a Windows managed app protection policy
resource "microsoft365_graph_beta_device_and_app_management_windows_managed_mobile_app" "teams" {
  managed_app_protection_id = "00000000-0000-0000-0000-000000000001" # Replace with your Windows managed app protection policy ID

  mobile_app_identifier = {
    windows_app_id = "Microsoft.Teams_8wekyb3d8bbwe"
  }

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}

# Resource for adding Microsoft OneDrive to a Windows managed app protection policy
resource "microsoft365_graph_beta_device_and_app_management_windows_managed_mobile_app" "onedrive" {
  managed_app_protection_id = "00000000-0000-0000-0000-000000000001" # Replace with your Windows managed app protection policy ID

  mobile_app_identifier = {
    windows_app_id = "Microsoft.OneDrive_8wekyb3d8bbwe"
  }

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}

# Resource for adding Microsoft Word to a Windows managed app protection policy
resource "microsoft365_graph_beta_device_and_app_management_windows_managed_mobile_app" "word" {
  managed_app_protection_id = "00000000-0000-0000-0000-000000000001" # Replace with your Windows managed app protection policy ID

  mobile_app_identifier = {
    windows_app_id = "Microsoft.Office.Word_8wekyb3d8bbwe"
  }

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}

# Resource for adding Microsoft Excel to a Windows managed app protection policy
resource "microsoft365_graph_beta_device_and_app_management_windows_managed_mobile_app" "excel" {
  managed_app_protection_id = "00000000-0000-0000-0000-000000000001" # Replace with your Windows managed app protection policy ID

  mobile_app_identifier = {
    windows_app_id = "Microsoft.Office.Excel_8wekyb3d8bbwe"
  }

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
} 