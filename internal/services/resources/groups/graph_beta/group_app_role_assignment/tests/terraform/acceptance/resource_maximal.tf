resource "random_string" "maximal_suffix" {
  length  = 8
  special = false
  upper   = false
}

resource "microsoft365_graph_beta_groups_group" "maximal" {
  display_name     = "App Role Assignment Test Maximal ${random_string.maximal_suffix.result}"
  mail_nickname    = "approletestmax${random_string.maximal_suffix.result}"
  mail_enabled     = false
  security_enabled = true
}

# Get the MileIQ Admin Center service principal (has app roles that support groups)
data "microsoft365_graph_beta_applications_service_principal" "mileiq" {
  filter_type  = "display_name"
  filter_value = "MileIQ Admin Center"
}

# Assign the MileIQ admin role to the group with custom timeouts
resource "microsoft365_graph_beta_groups_group_app_role_assignment" "maximal" {
  target_group_id    = microsoft365_graph_beta_groups_group.maximal.id
  resource_object_id = data.microsoft365_graph_beta_applications_service_principal.mileiq.items[0].id
  app_role_id        = "ea358ccf-c4a8-48ac-8b94-2558ae2f7a5c" # mdladmincenterrole.admin

  timeouts = {
    create = "5m"
    read   = "5m"
    update = "5m"
    delete = "5m"
  }
}

