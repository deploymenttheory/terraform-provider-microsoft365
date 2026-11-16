variable "service_principal_id" {
  type = string
}

variable "job_id" {
  type = string
}

variable "synchronization_rules_json" {
  type = string
}

variable "directories_json" {
  type = string
}

resource "microsoft365_graph_beta_applications_service_principal_synchronization_schema_json" "test" {
  service_principal_id       = var.service_principal_id
  job_id                     = var.job_id
  synchronization_rules_json = var.synchronization_rules_json
  directories_json           = var.directories_json
}
