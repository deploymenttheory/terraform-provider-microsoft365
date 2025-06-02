# Example 1: Assign Office 365 E3 license to a user
resource "microsoft365_graph_beta_user_license_assignment" "user_e3_license" {
  user_id = "john.doe@example.com"  # Can be user ID or UPN

  add_licenses = [
    {
      sku_id = "6fd2c87f-b296-42f0-b197-1e91e994b900"  # Office 365 E3
      disabled_plans = [
        "efb87545-963c-4e0d-99df-69c6916d9eb0"  # Disable specific service plan
      ]
    }
  ]

  remove_licenses = []

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}

# Example 2: Assign multiple licenses to a user
resource "microsoft365_graph_beta_user_license_assignment" "user_multiple_licenses" {
  user_id = "jane.smith@example.com"

  add_licenses = [
    {
      sku_id = "6fd2c87f-b296-42f0-b197-1e91e994b900"  # Office 365 E3
      disabled_plans = []
    },
    {
      sku_id = "b05e124f-c7cc-45a0-a6aa-8cf78c946968"  # Enterprise Mobility + Security E5
      disabled_plans = [
        "8a256a2b-b617-496d-b51b-e76466e88db0"  # Disable Microsoft Defender for Cloud Apps
      ]
    }
  ]

  remove_licenses = []

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}

# Example 3: Remove specific licenses from a user
resource "microsoft365_graph_beta_user_license_assignment" "user_license_removal" {
  user_id = "bob.johnson@example.com"

  add_licenses = []

  remove_licenses = [
    "6fd2c87f-b296-42f0-b197-1e91e994b900",  # Remove Office 365 E3
    "b05e124f-c7cc-45a0-a6aa-8cf78c946968"   # Remove Enterprise Mobility + Security E5
  ]

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}

# Example 4: Replace existing licenses (remove old, add new)
resource "microsoft365_graph_beta_user_license_assignment" "user_license_replacement" {
  user_id = "alice.wilson@example.com"

  add_licenses = [
    {
      sku_id = "c7df2760-2c81-4ef7-b578-5b5392b571df"  # Office 365 E5
      disabled_plans = []
    }
  ]

  remove_licenses = [
    "6fd2c87f-b296-42f0-b197-1e91e994b900"  # Remove Office 365 E3
  ]

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}

# Example 5: Using a data source to get user ID dynamically
data "microsoft365_graph_beta_user" "target_user" {
  user_principal_name = "dynamic.user@example.com"
}

resource "microsoft365_graph_beta_user_license_assignment" "dynamic_user_license" {
  user_id = data.microsoft365_graph_beta_user.target_user.id

  add_licenses = [
    {
      sku_id = "6fd2c87f-b296-42f0-b197-1e91e994b900"  # Office 365 E3
      disabled_plans = []
    }
  ]

  remove_licenses = []

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
} 