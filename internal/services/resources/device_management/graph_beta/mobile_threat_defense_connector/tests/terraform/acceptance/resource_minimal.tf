# Pradeo partner connector. The acceptance tests only use partners that have no
# connector in the test tenant, so creating and deleting it leaves existing
# Mobile Threat Defense connections untouched.
resource "microsoft365_graph_beta_device_management_mobile_threat_defense_connector" "test" {
  id = "8d0ed095-8191-4bd3-8a41-953b22d51ff7"
}
