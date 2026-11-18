resource "microsoft365_graph_beta_applications_service_principal_claims_policy_json" "test" {
  service_principal_id = "%s"
  policy_json = jsonencode({
    includeBasicClaimSet         = true
    includeApplicationIdInIssuer = false
    claims = [
      {
        "@odata.type" = "#microsoft.graph.samlNameIdClaim"
        nameIdFormat  = "persistent"
        configurations = [{
          condition = null
          attribute = {
            "@odata.type"        = "#microsoft.graph.sourcedAttribute"
            id                   = "userprincipalname"
            source               = "user"
            isExtensionAttribute = false
          }
          transformations = []
        }]
      },
      {
        "@odata.type"           = "#microsoft.graph.customClaim"
        name                    = "codex4894transformed"
        namespace               = "urn:codex:4894"
        tokenFormat             = ["saml"]
        samlAttributeNameFormat = "basic"
        configurations = [{
          condition = {
            "@odata.type" = "#microsoft.graph.customClaimCondition"
            userType      = "members"
            memberOf      = []
          }
          attribute = {
            "@odata.type" = "#microsoft.graph.valueBasedAttribute"
            value         = "fallback"
          }
          transformations = [
            {
              "@odata.type" = "#microsoft.graph.toLowercaseTransformation"
              input = {
                treatAsMultiValue = false
                attribute = {
                  "@odata.type"        = "#microsoft.graph.sourcedAttribute"
                  id                   = "userprincipalname"
                  source               = "user"
                  isExtensionAttribute = false
                }
              }
            },
            {
              "@odata.type" = "#microsoft.graph.trimTransformation"
              type          = "leadingAndTrailing"
              value         = " "
            }
          ]
        }]
      }
    ]
  })
}
