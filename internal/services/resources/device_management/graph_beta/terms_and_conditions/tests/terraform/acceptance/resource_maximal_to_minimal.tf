resource "microsoft365_graph_beta_device_management_terms_and_conditions" "transition" {
  display_name         = "acc-test-terms-and-conditions-transition"
  title                = "Company Terms"
  body_text            = "These are the basic terms and conditions."
  acceptance_statement = "I accept these terms"

  timeouts = {
    create = "180s"
    read   = "180s"
    update = "180s"
    delete = "180s"
  }
}

