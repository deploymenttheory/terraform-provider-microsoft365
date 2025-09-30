# Basic resource operation for device management
resource "microsoft365_graph_beta_device_management_rbac_resource_operation" "device_read" {
  resource_name = "Device"
  action_name   = "Read"
  description   = "Allows reading device information"
}

# Resource operation for application management
resource "microsoft365_graph_beta_device_management_rbac_resource_operation" "app_deploy" {
  resource_name = "MobileApplication"
  action_name   = "Deploy"
  description   = "Allows deploying mobile applications to devices"
}

# Resource operation for compliance policy management
resource "microsoft365_graph_beta_device_management_rbac_resource_operation" "compliance_create" {
  resource_name = "CompliancePolicy"
  action_name   = "Create"
  description   = "Allows creating new compliance policies"
}

# Resource operation for script execution
resource "microsoft365_graph_beta_device_management_rbac_resource_operation" "script_execute" {
  resource_name = "DeviceManagementScript"
  action_name   = "Execute"
  description   = "Allows executing PowerShell scripts on managed devices"
}

# Resource operation for role assignment
resource "microsoft365_graph_beta_device_management_rbac_resource_operation" "role_assign" {
  resource_name = "RoleAssignment"
  action_name   = "Assign"
  description   = "Allows assigning roles to users and groups"
}

# Resource operation with custom timeouts
resource "microsoft365_graph_beta_device_management_rbac_resource_operation" "policy_update" {
  resource_name = "DeviceConfigurationPolicy"
  action_name   = "Update"
  description   = "Allows updating device configuration policies"

  timeouts = {
    create = "10m"
    read   = "5m"
    update = "10m"
    delete = "5m"
  }
}

# Example using variables for reusability
variable "resource_operations" {
  description = "Map of resource operations to create"
  type = map(object({
    resource_name = string
    action_name   = string
    description   = string
  }))
  default = {
    "device_wipe" = {
      resource_name = "Device"
      action_name   = "Wipe"
      description   = "Allows wiping corporate data from devices"
    }
    "app_uninstall" = {
      resource_name = "MobileApplication"
      action_name   = "Uninstall"
      description   = "Allows uninstalling applications from devices"
    }
  }
}

resource "microsoft365_graph_beta_device_management_rbac_resource_operation" "bulk_operations" {
  for_each = var.resource_operations

  resource_name = each.value.resource_name
  action_name   = each.value.action_name
  description   = each.value.description
}