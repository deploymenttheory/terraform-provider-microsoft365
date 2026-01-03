resource "microsoft365_graph_beta_device_management_terms_and_conditions" "description" {
  display_name         = "acc-test-terms-and-conditions-description"
  description          = "This is a test terms and conditions with description"
  title                = "Terms with Description"
  body_text            = "These are terms and conditions with a description field."
  acceptance_statement = "I accept these terms with description"

  timeouts = {
    create = "180s"
    read   = "180s"
    update = "180s"
    delete = "180s"
  }
}