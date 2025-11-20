variable "test_user_id" {
  description = "Test user ID for acceptance tests"
  type        = string
}

variable "test_license_sku_id" {
  description = "Test license SKU ID for acceptance tests"
  type        = string
}

resource "microsoft365_graph_beta_users_user_license_assignment" "minimal" {
  user_id = var.test_user_id
  add_licenses = [{
    sku_id = var.test_license_sku_id
  }]
}

