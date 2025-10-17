# Example 1: Basic Allowed Value
resource "microsoft365_graph_beta_identity_and_access_custom_security_attribute_allowed_value" "engineering_project_alpha" {
  custom_security_attribute_definition_id = microsoft365_graph_beta_identity_and_access_custom_security_attribute_definition.department.id
  id                                      = "Alpha"
  is_active                               = true

  # Optional: Define custom timeouts
  timeouts = {
    create = "10m"
    read   = "5m"
    update = "10m"
    delete = "10m"
  }
}

# Example 2: Multiple Allowed Values for Department
resource "microsoft365_graph_beta_identity_and_access_custom_security_attribute_allowed_value" "dept_engineering" {
  custom_security_attribute_definition_id = microsoft365_graph_beta_identity_and_access_custom_security_attribute_definition.department.id
  id                                      = "Engineering"
  is_active                               = true
}

resource "microsoft365_graph_beta_identity_and_access_custom_security_attribute_allowed_value" "dept_sales" {
  custom_security_attribute_definition_id = microsoft365_graph_beta_identity_and_access_custom_security_attribute_definition.department.id
  id                                      = "Sales"
  is_active                               = true
}

resource "microsoft365_graph_beta_identity_and_access_custom_security_attribute_allowed_value" "dept_marketing" {
  custom_security_attribute_definition_id = microsoft365_graph_beta_identity_and_access_custom_security_attribute_definition.department.id
  id                                      = "Marketing"
  is_active                               = true
}

resource "microsoft365_graph_beta_identity_and_access_custom_security_attribute_allowed_value" "dept_hr" {
  custom_security_attribute_definition_id = microsoft365_graph_beta_identity_and_access_custom_security_attribute_definition.department.id
  id                                      = "Human Resources"
  is_active                               = true
}

# Example 3: Deprecated/Inactive Allowed Value
resource "microsoft365_graph_beta_identity_and_access_custom_security_attribute_allowed_value" "legacy_dept" {
  custom_security_attribute_definition_id = microsoft365_graph_beta_identity_and_access_custom_security_attribute_definition.department.id
  id                                      = "Legacy Department"
  is_active                               = false # Marked as inactive
}

# Example 4: Office Locations with Allowed Values
resource "microsoft365_graph_beta_identity_and_access_custom_security_attribute_allowed_value" "location_seattle" {
  custom_security_attribute_definition_id = microsoft365_graph_beta_identity_and_access_custom_security_attribute_definition.office_locations.id
  id                                      = "Seattle"
  is_active                               = true
}

resource "microsoft365_graph_beta_identity_and_access_custom_security_attribute_allowed_value" "location_new_york" {
  custom_security_attribute_definition_id = microsoft365_graph_beta_identity_and_access_custom_security_attribute_definition.office_locations.id
  id                                      = "New York"
  is_active                               = true
}

resource "microsoft365_graph_beta_identity_and_access_custom_security_attribute_allowed_value" "location_london" {
  custom_security_attribute_definition_id = microsoft365_graph_beta_identity_and_access_custom_security_attribute_definition.office_locations.id
  id                                      = "London"
  is_active                               = true
}

resource "microsoft365_graph_beta_identity_and_access_custom_security_attribute_allowed_value" "location_tokyo" {
  custom_security_attribute_definition_id = microsoft365_graph_beta_identity_and_access_custom_security_attribute_definition.office_locations.id
  id                                      = "Tokyo"
  is_active                               = true
}

# Example 5: Project Names with Allowed Values
resource "microsoft365_graph_beta_identity_and_access_custom_security_attribute_allowed_value" "project_apollo" {
  custom_security_attribute_definition_id = "Engineering_ProjectName"
  id                                      = "Apollo"
  is_active                               = true
}

resource "microsoft365_graph_beta_identity_and_access_custom_security_attribute_allowed_value" "project_orion" {
  custom_security_attribute_definition_id = "Engineering_ProjectName"
  id                                      = "Orion"
  is_active                               = true
}

# Example 6: Values with Special Characters
resource "microsoft365_graph_beta_identity_and_access_custom_security_attribute_allowed_value" "classification_level_1" {
  custom_security_attribute_definition_id = "Security_Classification"
  id                                      = "Level-1"
  is_active                               = true

  timeouts = {
    create = "10m"
    read   = "5m"
    update = "10m"
    delete = "10m"
  }
}

resource "microsoft365_graph_beta_identity_and_access_custom_security_attribute_allowed_value" "classification_level_2" {
  custom_security_attribute_definition_id = "Security_Classification"
  id                                      = "Level-2"
  is_active                               = true
}

# Note: You can define up to 100 allowed values per custom security attribute definition
# The id is the identifier for the predefined value (e.g., "Alpine", "Engineering")
# Allowed values cannot be deleted, only deactivated by setting is_active to false
