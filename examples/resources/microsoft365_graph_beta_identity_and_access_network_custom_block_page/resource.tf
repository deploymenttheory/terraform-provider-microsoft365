# One custom block page exists per tenant. Initial apply updates that singleton.
# Destroy leaves these settings in Entra and removes Terraform management only.
resource "microsoft365_graph_beta_identity_and_access_network_custom_block_page" "example" {
  state = "enabled"
  configuration = {
    body = "Access to this website is blocked by your organization. [Contact support](https://support.example.com)."
  }
}
