resource "microsoft365_graph_beta_applications_service_principal_synchronization_schema_json" "test" {
  service_principal_id = "11111111-1111-1111-1111-111111111111"
  job_id               = "scim.00000000000000000000000000000001.22222222-2222-2222-2222-222222222222"

  synchronization_rules_json = jsonencode([
    {
      "editable"            = true
      "id"                  = "03f7d90d-bf71-41b1-bda6-aaf0ddbee5d8"
      "name"                = "USERGROUP_OUTBOUND_USERGROUP"
      "priority"            = 1
      "sourceDirectoryName" = "Microsoft Entra ID"
      "targetDirectoryName" = "customappsso"
      "containerFilter"     = null
      "groupFilter"         = null
      "metadata"            = []
      "objectMappings"      = [
        {
          "enabled"           = false
          "flowTypes"         = "Add, Update, Delete"
          "name"              = "Provision Microsoft Entra ID Groups"
          "sourceObjectName"  = "Group"
          "targetObjectName"  = "urn:ietf:params:scim:schemas:core:2.0:Group"
          "scope"             = null
          "attributeMappings" = [
            {
              "defaultValue"            = null
              "exportMissingReferences" = false
              "flowBehavior"            = "FlowWhenChanged"
              "flowType"                = "Always"
              "matchingPriority"        = 1
              "targetAttributeName"     = "displayName"
              "source"                  = {
                "expression" = "[displayName]"
                "name"       = "displayName"
                "type"       = "Attribute"
                "parameters" = []
              }
            }
          ]
          "metadata"          = [
            {
              "key"   = "IsCustomerDefined"
              "value" = "false"
            },
            {
              "key"   = "Unsynchronized"
              "value" = "false"
            }
          ]
        },
        {
          "enabled"           = true
          "flowTypes"         = "Add, Update, Delete"
          "name"              = "Provision Microsoft Entra ID Users"
          "sourceObjectName"  = "User"
          "targetObjectName"  = "urn:ietf:params:scim:schemas:extension:enterprise:2.0:User"
          "scope"             = null
          "attributeMappings" = [
            {
              "defaultValue"            = null
              "exportMissingReferences" = false
              "flowBehavior"            = "FlowWhenChanged"
              "flowType"                = "Always"
              "matchingPriority"        = 1
              "targetAttributeName"     = "userName"
              "source"                  = {
                "expression" = "[userPrincipalName]"
                "name"       = "userPrincipalName"
                "type"       = "Attribute"
                "parameters" = []
              }
            },
            {
              "defaultValue"            = null
              "exportMissingReferences" = false
              "flowBehavior"            = "FlowWhenChanged"
              "flowType"                = "Always"
              "matchingPriority"        = 0
              "targetAttributeName"     = "active"
              "source"                  = {
                "expression" = "Not([IsSoftDeleted])"
                "name"       = "Not"
                "type"       = "Function"
                "parameters" = [
                  {
                    "key"   = "source"
                    "value" = {
                      "expression" = "[IsSoftDeleted]"
                      "name"       = "IsSoftDeleted"
                      "type"       = "Attribute"
                      "parameters" = []
                    }
                  }
                ]
              }
            }
          ]
          "metadata"          = [
            {
              "key"   = "IsCustomerDefined"
              "value" = "false"
            },
            {
              "key"   = "Unsynchronized"
              "value" = "false"
            }
          ]
        }
      ]
    }
  ])

  directories_json = jsonencode([
    {
      "id"                = "66e4a8cc-1b7b-435e-95f8-f06cea133828"
      "discoveryDateTime" = null
      "discoverabilities" = "None"
      "name"              = "Microsoft Entra ID"
      "readOnly"          = true
      "version"           = "69F1E34D-A39F-4A48-B627-3AE0AD8E0DF6"
      "objects"           = [
        {
          "name"          = "User"
          "supportedApis" = []
          "metadata"      = []
          "attributes"    = [
            {
              "anchor"            = true
              "caseExact"         = false
              "defaultValue"      = null
              "flowNullValues"    = false
              "multivalued"       = false
              "mutability"        = "ReadWrite"
              "name"              = "objectId"
              "required"          = false
              "type"              = "String"
              "apiExpressions"    = []
              "metadata"          = []
              "referencedObjects" = []
            }
          ]
        }
      ]
    },
    {
      "id"                = "09760868-cafb-47ac-9031-0a3262300427"
      "discoveryDateTime" = null
      "discoverabilities" = "None"
      "name"              = "customappsso"
      "readOnly"          = false
      "version"           = "6.0"
      "objects"           = [
        {
          "name"          = "urn:ietf:params:scim:schemas:extension:enterprise:2.0:User"
          "supportedApis" = []
          "metadata"      = []
          "attributes"    = [
            {
              "anchor"            = true
              "caseExact"         = false
              "defaultValue"      = null
              "flowNullValues"    = false
              "multivalued"       = false
              "mutability"        = "ReadWrite"
              "name"              = "userName"
              "required"          = true
              "type"              = "String"
              "apiExpressions"    = []
              "metadata"          = []
              "referencedObjects" = []
            }
          ]
        }
      ]
    }
  ])
}
