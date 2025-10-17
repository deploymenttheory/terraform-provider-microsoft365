# Example 1: Simple String Attribute (Most Common Use Case)
resource "microsoft365_graph_beta_identity_and_access_custom_security_attribute_definition" "project_name" {
  attribute_set               = "Engineering"
  name                        = "ProjectName"
  description                 = "Name of the project the user is assigned to"
  type                        = "String"
  status                      = "Available"
  is_collection               = false
  is_searchable               = true
  use_pre_defined_values_only = false

  # Optional: Define custom timeouts
  timeouts = {
    create = "10m"
    read   = "5m"
    update = "10m"
    delete = "10m"
  }
}

# Example 2: Integer Attribute
resource "microsoft365_graph_beta_identity_and_access_custom_security_attribute_definition" "cost_center" {
  attribute_set               = microsoft365_graph_beta_identity_and_access_attribute_set.example.id
  name                        = "CostCenter"
  description                 = "Cost center number for budget tracking"
  type                        = "Integer"
  status                      = "Available"
  is_collection               = false
  is_searchable               = true
  use_pre_defined_values_only = false
}

# Example 3: Boolean Attribute
resource "microsoft365_graph_beta_identity_and_access_custom_security_attribute_definition" "security_clearance" {
  attribute_set               = "Security"
  name                        = "HasClearance"
  description                 = "Indicates if user has security clearance"
  type                        = "Boolean"
  status                      = "Available"
  is_collection               = false # Must be false for Boolean type
  is_searchable               = true
  use_pre_defined_values_only = false # Must be false for Boolean type
}

# Example 4: Collection Attribute (Multiple Values)
resource "microsoft365_graph_beta_identity_and_access_custom_security_attribute_definition" "skills" {
  attribute_set               = "HumanResources"
  name                        = "Skills"
  description                 = "Skills and competencies of the employee"
  type                        = "String"
  status                      = "Available"
  is_collection               = true # Allows multiple values
  is_searchable               = true
  use_pre_defined_values_only = false
}

# Example 5: Attribute with Predefined Values Only
resource "microsoft365_graph_beta_identity_and_access_custom_security_attribute_definition" "department" {
  attribute_set               = "Organization"
  name                        = "Department"
  description                 = "Department assignment with predefined values"
  type                        = "String"
  status                      = "Available"
  is_collection               = false
  is_searchable               = true
  use_pre_defined_values_only = true # Only predefined values allowed

  # Note: After creating this definition, you would need to add allowed values
  # using the microsoft365_graph_beta_identity_and_access_custom_security_attribute_allowed_value resource
}

# Example 6: Multi-Collection Attribute for Team Assignments
resource "microsoft365_graph_beta_identity_and_access_custom_security_attribute_definition" "team_assignments" {
  attribute_set               = "Engineering"
  name                        = "Teams"
  description                 = "Engineering teams the user belongs to"
  type                        = "String"
  status                      = "Available"
  is_collection               = true
  is_searchable               = true
  use_pre_defined_values_only = true # Restrict to predefined team names

  timeouts = {
    create = "15m"
    read   = "5m"
    update = "15m"
    delete = "15m"
  }
}

# Example 7: Deprecated Attribute (Inactive State)
resource "microsoft365_graph_beta_identity_and_access_custom_security_attribute_definition" "legacy_code" {
  attribute_set               = "Legacy"
  name                        = "OldSystemCode"
  description                 = "Legacy system identifier (deprecated)"
  type                        = "String"
  status                      = "Deprecated" # Marks as inactive
  is_collection               = false
  is_searchable               = false
  use_pre_defined_values_only = false
}

# Example 8: Employee Classification with Limited Search
resource "microsoft365_graph_beta_identity_and_access_custom_security_attribute_definition" "employee_classification" {
  attribute_set               = "HumanResources"
  name                        = "Classification"
  description                 = "Employee classification level"
  type                        = "String"
  status                      = "Available"
  is_collection               = false
  is_searchable               = false # Not searchable for privacy
  use_pre_defined_values_only = true
}

# Example 9: Attribute for Compliance Tracking
resource "microsoft365_graph_beta_identity_and_access_custom_security_attribute_definition" "compliance_training_date" {
  attribute_set               = "Compliance"
  name                        = "LastTrainingDate"
  description                 = "Date of last compliance training completion (format: YYYY-MM-DD)"
  type                        = "String"
  status                      = "Available"
  is_collection               = false
  is_searchable               = true
  use_pre_defined_values_only = false

  timeouts = {
    create = "10m"
    read   = "5m"
    update = "10m"
    delete = "10m"
  }
}

# Example 10: Office Location Tracking
resource "microsoft365_graph_beta_identity_and_access_custom_security_attribute_definition" "office_locations" {
  attribute_set               = "Facilities"
  name                        = "Locations"
  description                 = "Office locations where the employee works"
  type                        = "String"
  status                      = "Available"
  is_collection               = true # Employee can work from multiple locations
  is_searchable               = true
  use_pre_defined_values_only = true

  # After creation, add specific office locations as allowed values
}

