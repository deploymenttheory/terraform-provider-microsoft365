resource "microsoft365_graph_beta_device_management_linux_device_compliance_policy" "test_004" {
  name                       = "unit-test-linux-compliance-004"
  device_encryption_required = true
  description                = "Maximal Linux compliance policy"
  custom_compliance_required = false
  distribution_allowed_distros = [
    {
      type            = "ubuntu"
      minimum_version = "22.04"
      maximum_version = "24.04"
    },
    {
      type            = "rhel"
      minimum_version = "8.0"
      maximum_version = "9.9"
    }
  ]
  password_policy_minimum_digits    = 2
  password_policy_minimum_length    = 12
  password_policy_minimum_lowercase = 2
  password_policy_minimum_symbols   = 1
  password_policy_minimum_uppercase = 2
  role_scope_tag_ids                = ["0", "1"]
  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}
