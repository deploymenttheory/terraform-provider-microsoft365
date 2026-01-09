resource "microsoft365_graph_beta_device_management_terms_and_conditions" "maximal" {
  display_name         = "unit-test-terms-and-conditions-maximal"
  description          = "Comprehensive terms and conditions for testing with all features"
  title                = "Complete Company Terms and Conditions"
  body_text            = "These are the comprehensive terms and conditions that all users must read and accept before accessing company resources. This includes detailed policies about data usage, privacy, security requirements, and acceptable use of company systems."
  acceptance_statement = "I have read and agree to abide by all terms and conditions outlined above"
  role_scope_tag_ids   = ["0", "1"]

  timeouts = {
    create = "180s"
    read   = "180s"
    update = "180s"
    delete = "180s"
  }
}