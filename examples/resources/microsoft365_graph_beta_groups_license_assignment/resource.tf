# Example 1: Minimal - Assign a single Office 365 E3 license to a group
resource "microsoft365_graph_beta_groups_license_assignment" "e3_license" {
  group_id = "1132b215-826f-42a9-8cfe-1643d19d17fd"
  sku_id   = "6fd2c87f-b296-42f0-b197-1e91e994b900" # Office 365 E3
}

# Example 2: Single license with disabled service plans
resource "microsoft365_graph_beta_groups_license_assignment" "e3_custom" {
  group_id = "1132b215-826f-42a9-8cfe-1643d19d17fd"
  sku_id   = "6fd2c87f-b296-42f0-b197-1e91e994b900" # Office 365 E3

  # Disable specific service plans within the license
  disabled_plans = [
    "efb87545-963c-4e0d-99df-69c6916d9eb0", # Azure Information Protection Premium P1
    "9f431833-0334-42de-a7dc-70aa40db46db"  # Microsoft Stream
  ]
}

# Example 3: Assign multiple licenses to the same group
# Each license requires a separate resource instance
resource "microsoft365_graph_beta_groups_license_assignment" "group_e3" {
  group_id = "2243c326-937g-53f0-c9df-2e68f106b901"
  sku_id   = "6fd2c87f-b296-42f0-b197-1e91e994b900" # Office 365 E3
}

resource "microsoft365_graph_beta_groups_license_assignment" "group_ems_e5" {
  group_id = "2243c326-937g-53f0-c9df-2e68f106b901"
  sku_id   = "b05e124f-c7cc-45a0-a6aa-8cf78c946968" # Enterprise Mobility + Security E5

  disabled_plans = [
    "113feb6c-3fe4-4440-bddc-54d774bf0318", # Exchange Foundation
    "14ab5db5-e6c4-4b20-b4bc-13e36fd2227f"  # Intune for Education
  ]
}

resource "microsoft365_graph_beta_groups_license_assignment" "group_power_bi" {
  group_id = "2243c326-937g-53f0-c9df-2e68f106b901"
  sku_id   = "f30db892-07e9-47e9-837c-80727f46fd3d" # Power BI Free
}

# Example 4: Using a data source to get group ID dynamically
data "microsoft365_graph_beta_groups_group" "sales_team" {
  display_name = "Sales Team"
}

resource "microsoft365_graph_beta_groups_license_assignment" "sales_e3" {
  group_id = data.microsoft365_graph_beta_groups_group.sales_team.id
  sku_id   = "6fd2c87f-b296-42f0-b197-1e91e994b900" # Office 365 E3
}

# Example 5: Assign license to a newly created group
resource "microsoft365_graph_beta_groups_group" "engineering" {
  display_name     = "Engineering Team"
  mail_nickname    = "engineering"
  mail_enabled     = false
  security_enabled = true
}

resource "microsoft365_graph_beta_groups_license_assignment" "engineering_e5" {
  group_id = microsoft365_graph_beta_groups_group.engineering.id
  sku_id   = "c7df2760-2c81-4ef7-b578-5b5392b571df" # Office 365 E5

  disabled_plans = [
    "57ff2da0-773e-42df-b2af-ffb7a2317929" # Teams
  ]
}

# Example 6: With custom timeouts
resource "microsoft365_graph_beta_groups_license_assignment" "custom_timeout" {
  group_id = "4465e548-159i-75h2-e1fg-4g80h328d123"
  sku_id   = "c7df2760-2c81-4ef7-b578-5b5392b571df" # Office 365 E5

  timeouts {
    create = "300s"
    read   = "180s"
    update = "300s"
    delete = "300s"
  }
}

# Example 7: Department-based license assignments
# Marketing department gets Office 365 E3
resource "microsoft365_graph_beta_groups_license_assignment" "marketing_e3" {
  group_id = "marketing-group-uuid-here"
  sku_id   = "6fd2c87f-b296-42f0-b197-1e91e994b900" # Office 365 E3
}

# Engineering department gets Office 365 E5 with full features
resource "microsoft365_graph_beta_groups_license_assignment" "engineering_e5_full" {
  group_id = "engineering-group-uuid-here"
  sku_id   = "c7df2760-2c81-4ef7-b578-5b5392b571df" # Office 365 E5

  # No disabled plans - all features enabled
}

# Sales department gets E3 with some features disabled
resource "microsoft365_graph_beta_groups_license_assignment" "sales_e3_limited" {
  group_id = "sales-group-uuid-here"
  sku_id   = "6fd2c87f-b296-42f0-b197-1e91e994b900" # Office 365 E3

  disabled_plans = [
    "efb87545-963c-4e0d-99df-69c6916d9eb0", # Azure Information Protection
    "9f431833-0334-42de-a7dc-70aa40db46db", # Microsoft Stream
    "b737dad2-2f6c-4c65-90e3-ca563267e8b9"  # Yammer Enterprise
  ]
}