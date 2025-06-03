# Example 1: Assign Office 365 E3 license to a group
resource "microsoft365_graph_beta_group_license_assignment" "group_e3_license" {
  group_id = "1132b215-826f-42a9-8cfe-1643d19d17fd"  # Group UUID

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
    create = "180s"
    read   = "180s"
    update = "180s"
    delete = "180s"
  }
}

# Example 2: Assign multiple licenses to a group
resource "microsoft365_graph_beta_group_license_assignment" "group_multiple_licenses" {
  group_id = "2243c326-937g-53f0-c9df-2e68f106b901"

  add_licenses = [
    {
      sku_id = "6fd2c87f-b296-42f0-b197-1e91e994b900"  # Office 365 E3
      disabled_plans = []
    },
    {
      sku_id = "b05e124f-c7cc-45a0-a6aa-8cf78c946968"  # Enterprise Mobility + Security E5
      disabled_plans = [
        "113feb6c-3fe4-4440-bddc-54d774bf0318",  # Disable Exchange Foundation
        "14ab5db5-e6c4-4b20-b4bc-13e36fd2227f"   # Disable another service plan
      ]
    }
  ]

  remove_licenses = []

  timeouts = {
    create = "180s"
    read   = "180s"
    update = "180s"
    delete = "180s"
  }
}

# Example 3: Remove specific licenses from a group
resource "microsoft365_graph_beta_group_license_assignment" "group_license_removal" {
  group_id = "3354d437-048h-64g1-d0ef-3f79g217c012"

  add_licenses = []

  remove_licenses = [
    "6fd2c87f-b296-42f0-b197-1e91e994b900",  # Remove Office 365 E3
    "b05e124f-c7cc-45a0-a6aa-8cf78c946968"   # Remove Enterprise Mobility + Security E5
  ]

  timeouts = {
    create = "180s"
    read   = "180s"
    update = "180s"
    delete = "180s"
  }
}

# Example 4: Replace existing licenses (remove old, add new)
resource "microsoft365_graph_beta_group_license_assignment" "group_license_replacement" {
  group_id = "4465e548-159i-75h2-e1fg-4g80h328d123"

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
    create = "180s"
    read   = "180s"
    update = "180s"
    delete = "180s"
  }
}

# Example 5: Using a data source to get group ID dynamically
data "microsoft365_graph_beta_group" "target_group" {
  display_name = "Sales Team"
}

resource "microsoft365_graph_beta_group_license_assignment" "dynamic_group_license" {
  group_id = data.microsoft365_graph_beta_group.target_group.id

  add_licenses = [
    {
      sku_id = "6fd2c87f-b296-42f0-b197-1e91e994b900"  # Office 365 E3
      disabled_plans = []
    }
  ]

  remove_licenses = []

  timeouts = {
    create = "180s"
    read   = "180s"
    update = "180s"
    delete = "180s"
  }
} 