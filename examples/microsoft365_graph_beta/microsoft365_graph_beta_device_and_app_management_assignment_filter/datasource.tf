# Basic usage - looking up a single filter by display name
data "microsoft365_graph_beta_device_and_app_management_assignment_filter" "by_name" {
  display_name = "Filter | Android Enterprise Device Status Is Rooted"
}

# Look up by ID
data "microsoft365_graph_beta_device_and_app_management_assignment_filter" "windows_vdi" {
  id = "2983b1c2-8ec2-45d3-84ed-deca619d2c04"
}

# Example: Create new filter based on existing one (using name lookup)
resource "microsoft365_graph_beta_device_and_app_management_assignment_filter" "clone_android" {
  display_name = "Clone - Android Rooted Device Filter"
  description  = "Cloned from: ${data.microsoft365_graph_beta_device_and_app_management_assignment_filter.by_name.description}"

  # Copy configuration from existing filter
  platform                          = data.microsoft365_graph_beta_device_and_app_management_assignment_filter.by_name.platform
  rule                              = data.microsoft365_graph_beta_device_and_app_management_assignment_filter.by_name.rule
  assignment_filter_management_type = data.microsoft365_graph_beta_device_and_app_management_assignment_filter.by_name.assignment_filter_management_type
  role_scope_tags                   = data.microsoft365_graph_beta_device_and_app_management_assignment_filter.by_name.role_scope_tags

  timeouts = {
    create = "10s"
    read   = "10s"
    update = "10s"
    delete = "10s"
  }
}

# Output showing all available attributes
output "filter_details" {
  value = {
    # Basic details
    id           = data.microsoft365_graph_beta_device_and_app_management_assignment_filter.by_name.id
    display_name = data.microsoft365_graph_beta_device_and_app_management_assignment_filter.by_name.display_name
    description  = data.microsoft365_graph_beta_device_and_app_management_assignment_filter.by_name.description

    # Filter configuration
    platform                          = data.microsoft365_graph_beta_device_and_app_management_assignment_filter.by_name.platform
    rule                              = data.microsoft365_graph_beta_device_and_app_management_assignment_filter.by_name.rule
    assignment_filter_management_type = data.microsoft365_graph_beta_device_and_app_management_assignment_filter.by_name.assignment_filter_management_type

    # Additional metadata
    created_date_time       = data.microsoft365_graph_beta_device_and_app_management_assignment_filter.by_name.created_date_time
    last_modified_date_time = data.microsoft365_graph_beta_device_and_app_management_assignment_filter.by_name.last_modified_date_time
    role_scope_tags         = data.microsoft365_graph_beta_device_and_app_management_assignment_filter.by_name.role_scope_tags
  }
}


# Example: Create new filter based on Windows VDI filter (using ID lookup)
resource "microsoft365_graph_beta_device_and_app_management_assignment_filter" "clone_windows_vdi" {
  display_name = "Clone - Windows VDI Device Filter"
  description  = "Cloned from: ${data.microsoft365_graph_beta_device_and_app_management_assignment_filter.windows_vdi.description}"

  # Copy configuration from existing filter
  platform                          = data.microsoft365_graph_beta_device_and_app_management_assignment_filter.windows_vdi.platform
  rule                              = data.microsoft365_graph_beta_device_and_app_management_assignment_filter.windows_vdi.rule
  assignment_filter_management_type = data.microsoft365_graph_beta_device_and_app_management_assignment_filter.windows_vdi.assignment_filter_management_type
  role_scope_tags                   = data.microsoft365_graph_beta_device_and_app_management_assignment_filter.windows_vdi.role_scope_tags

  timeouts = {
    create = "10s"
    read   = "10s"
    update = "10s"
    delete = "10s"
  }
}

# Output showing Windows VDI filter attributes
output "vdi_filter_details" {
  value = {
    # Basic details
    id           = data.microsoft365_graph_beta_device_and_app_management_assignment_filter.windows_vdi.id
    display_name = data.microsoft365_graph_beta_device_and_app_management_assignment_filter.windows_vdi.display_name
    description  = data.microsoft365_graph_beta_device_and_app_management_assignment_filter.windows_vdi.description

    # Filter configuration
    platform                          = data.microsoft365_graph_beta_device_and_app_management_assignment_filter.windows_vdi.platform
    rule                              = data.microsoft365_graph_beta_device_and_app_management_assignment_filter.windows_vdi.rule
    assignment_filter_management_type = data.microsoft365_graph_beta_device_and_app_management_assignment_filter.windows_vdi.assignment_filter_management_type

    # Additional metadata
    created_date_time       = data.microsoft365_graph_beta_device_and_app_management_assignment_filter.windows_vdi.created_date_time
    last_modified_date_time = data.microsoft365_graph_beta_device_and_app_management_assignment_filter.windows_vdi.last_modified_date_time
    role_scope_tags         = data.microsoft365_graph_beta_device_and_app_management_assignment_filter.windows_vdi.role_scope_tags
  }
}


# Use Case 1: Filter Migration - Export multiple filters as JSON for documentation/migration
output "all_filters_export" {
  value = {
    android_filter = {
      name = data.microsoft365_graph_beta_device_and_app_management_assignment_filter.by_name.display_name
      config = {
        platform = data.microsoft365_graph_beta_device_and_app_management_assignment_filter.by_name.platform
        rule     = data.microsoft365_graph_beta_device_and_app_management_assignment_filter.by_name.rule
        type     = data.microsoft365_graph_beta_device_and_app_management_assignment_filter.by_name.assignment_filter_management_type
      }
    }
    vdi_filter = {
      name = data.microsoft365_graph_beta_device_and_app_management_assignment_filter.windows_vdi.display_name
      config = {
        platform = data.microsoft365_graph_beta_device_and_app_management_assignment_filter.windows_vdi.platform
        rule     = data.microsoft365_graph_beta_device_and_app_management_assignment_filter.windows_vdi.rule
        type     = data.microsoft365_graph_beta_device_and_app_management_assignment_filter.windows_vdi.assignment_filter_management_type
      }
    }
  }
}

# Use Case 2: Create multiple environment-specific clones with prefix
resource "microsoft365_graph_beta_device_and_app_management_assignment_filter" "prod_clone" {
  display_name = "PROD - ${data.microsoft365_graph_beta_device_and_app_management_assignment_filter.windows_vdi.display_name}"
  description  = "Production clone of: ${data.microsoft365_graph_beta_device_and_app_management_assignment_filter.windows_vdi.description}"
  
  platform                          = data.microsoft365_graph_beta_device_and_app_management_assignment_filter.windows_vdi.platform
  rule                              = data.microsoft365_graph_beta_device_and_app_management_assignment_filter.windows_vdi.rule
  assignment_filter_management_type = data.microsoft365_graph_beta_device_and_app_management_assignment_filter.windows_vdi.assignment_filter_management_type
  role_scope_tags                   = data.microsoft365_graph_beta_device_and_app_management_assignment_filter.windows_vdi.role_scope_tags

  timeouts = {
    create = "10s"
    read   = "10s"
    update = "10s"
    delete = "10s"
  }
}

resource "microsoft365_graph_beta_device_and_app_management_assignment_filter" "dev_clone" {
  display_name = "DEV - ${data.microsoft365_graph_beta_device_and_app_management_assignment_filter.windows_vdi.display_name}"
  description  = "Development clone of: ${data.microsoft365_graph_beta_device_and_app_management_assignment_filter.windows_vdi.description}"
  
  platform                          = data.microsoft365_graph_beta_device_and_app_management_assignment_filter.windows_vdi.platform
  rule                              = data.microsoft365_graph_beta_device_and_app_management_assignment_filter.windows_vdi.rule
  assignment_filter_management_type = data.microsoft365_graph_beta_device_and_app_management_assignment_filter.windows_vdi.assignment_filter_management_type
  role_scope_tags                   = data.microsoft365_graph_beta_device_and_app_management_assignment_filter.windows_vdi.role_scope_tags

  timeouts = {
    create = "10s"
    read   = "10s"
    update = "10s"
    delete = "10s"
  }
}

# Use Case 3: Create a modified clone with an enhanced rule
resource "microsoft365_graph_beta_device_and_app_management_assignment_filter" "enhanced_vdi_filter" {
  display_name = "Enhanced - ${data.microsoft365_graph_beta_device_and_app_management_assignment_filter.windows_vdi.display_name}"
  description  = "Enhanced version with additional conditions"
  
  platform                          = data.microsoft365_graph_beta_device_and_app_management_assignment_filter.windows_vdi.platform
  # Original rule with additional conditions
  rule                              = "${data.microsoft365_graph_beta_device_and_app_management_assignment_filter.windows_vdi.rule} and (device.manufacturer -eq \"Microsoft\")"
  assignment_filter_management_type = data.microsoft365_graph_beta_device_and_app_management_assignment_filter.windows_vdi.assignment_filter_management_type
  role_scope_tags                   = data.microsoft365_graph_beta_device_and_app_management_assignment_filter.windows_vdi.role_scope_tags

  timeouts = {
    create = "10s"
    read   = "10s"
    update = "10s"
    delete = "10s"
  }
}

# Use Case 4: Output comparing multiple filters
output "filter_comparison" {
  value = {
    original_vs_enhanced = {
      original_rule  = data.microsoft365_graph_beta_device_and_app_management_assignment_filter.windows_vdi.rule
      enhanced_rule = "${data.microsoft365_graph_beta_device_and_app_management_assignment_filter.windows_vdi.rule} and (device.manufacturer -eq \"Microsoft\")"
      differences = {
        platform_same = (
          data.microsoft365_graph_beta_device_and_app_management_assignment_filter.windows_vdi.platform == 
          data.microsoft365_graph_beta_device_and_app_management_assignment_filter.by_name.platform
        )
        management_type_same = (
          data.microsoft365_graph_beta_device_and_app_management_assignment_filter.windows_vdi.assignment_filter_management_type == 
          data.microsoft365_graph_beta_device_and_app_management_assignment_filter.by_name.assignment_filter_management_type
        )
      }
    }
  }
}