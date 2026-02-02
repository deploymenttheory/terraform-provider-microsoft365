resource "microsoft365_graph_beta_applications_application" "minimal" {
  display_name = "my-minimal-app"
  description  = "A minimal application with only required fields"
}
