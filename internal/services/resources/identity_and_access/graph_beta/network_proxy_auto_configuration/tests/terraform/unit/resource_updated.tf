resource "microsoft365_graph_beta_identity_and_access_network_proxy_auto_configuration" "test" {
  name       = "test-pac"
  is_enabled = false
  content    = "// updated\nfunction FindProxyForURL(url, host) { return \"DIRECT\"; }"
}
