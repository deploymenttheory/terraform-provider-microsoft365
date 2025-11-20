variable "test_user_id" {
  description = "Test user ID for acceptance tests"
  type        = string
}

variable "test_license_sku_id_1" {
  description = "First test license SKU ID for acceptance tests"
  type        = string
}

variable "test_license_sku_id_2" {
  description = "Second test license SKU ID for acceptance tests"
  type        = string
}

variable "test_service_plan_id" {
  description = "Test service plan ID for acceptance tests (optional)"
  type        = string
  default     = ""
}

resource "microsoft365_graph_beta_users_user_license_assignment" "maximal" {
  user_id = var.test_user_id
  add_licenses = [{
    sku_id = var.test_license_sku_id_1
    disabled_plans = var.test_service_plan_id != "" ? [var.test_service_plan_id] : []
  }, {
    sku_id = var.test_license_sku_id_2
  }]
}

