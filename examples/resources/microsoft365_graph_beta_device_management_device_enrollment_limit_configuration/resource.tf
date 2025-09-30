# Basic device enrollment limit configuration
resource "microsoft365_graph_beta_device_management_device_enrollment_limit_configuration" "basic" {
  display_name = "Basic Device Limit"
  description  = "Limits users to 5 devices maximum"
  limit        = 5
  priority     = 1

  timeouts = {
    create = "180s"
    read   = "180s"
    update = "180s"
    delete = "180s"
  }
}

# Advanced device enrollment limit configuration with role scope tags
resource "microsoft365_graph_beta_device_management_device_enrollment_limit_configuration" "advanced" {
  display_name = "Department Device Limit"
  description  = "Allows IT department users to enroll up to 10 devices"
  limit        = 10
  priority     = 2

  role_scope_tag_ids = [
    "0", # Default scope tag
    "1"  # Custom scope tag
  ]

  timeouts = {
    create = "180s"
    read   = "180s"
    update = "180s"
    delete = "180s"
  }
}

# Minimal device enrollment limit configuration
resource "microsoft365_graph_beta_device_management_device_enrollment_limit_configuration" "minimal" {
  display_name = "Minimal Device Limit"
  limit        = 3
}

# High limit configuration for power users
resource "microsoft365_graph_beta_device_management_device_enrollment_limit_configuration" "power_users" {
  display_name = "Power User Device Limit"
  description  = "Allows power users to enroll up to 25 devices"
  limit        = 25
  priority     = 10

  timeouts = {
    create = "300s"
    read   = "300s"
    update = "300s"
    delete = "300s"
  }
} 