resource "microsoft365_graph_beta_applications_service_principal_claims_policy_json" "example" {
  service_principal_id = "11111111-1111-4111-8111-111111111111"
  policy_json = jsonencode({
    includeBasicClaimSet         = true
    includeApplicationIdInIssuer = false
    claims = [{
      "@odata.type"           = "#microsoft.graph.customClaim"
      name                    = "codex4894"
      namespace               = "urn:codex:4894"
      tokenFormat             = ["saml"]
      samlAttributeNameFormat = "basic"
      configurations = [{
        condition = null
        attribute = {
          "@odata.type" = "#microsoft.graph.valueBasedAttribute"
          value         = "fixture-value-1"
        }
        transformations = []
      }]
    }]
  })
}
