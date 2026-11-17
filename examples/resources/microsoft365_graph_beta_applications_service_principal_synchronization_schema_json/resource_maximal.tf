# Starts from the exported default schema and changes it in Terraform:
#   - disables group provisioning,
#   - adds a custom SCIM extension attribute to the target user object,
#   - maps the Microsoft Entra ID department attribute to it.
locals {
  default_rules       = jsondecode(file("${path.module}/synchronization_rules.json"))
  default_directories = jsondecode(file("${path.module}/directories.json"))

  custom_attribute = {
    anchor            = false
    caseExact         = false
    defaultValue      = null
    flowNullValues    = false
    multivalued       = false
    mutability        = "ReadWrite"
    name              = "urn:ietf:params:scim:schemas:extension:example:2.0:User:costCenter"
    required          = false
    type              = "String"
    apiExpressions    = []
    metadata          = []
    referencedObjects = []
  }

  custom_attribute_mapping = {
    defaultValue            = null
    exportMissingReferences = false
    flowBehavior            = "FlowWhenChanged"
    flowType                = "Always"
    matchingPriority        = 0
    targetAttributeName     = local.custom_attribute.name
    source = {
      expression = "[department]"
      name       = "department"
      type       = "Attribute"
      parameters = []
    }
  }

  synchronization_rules = [
    for rule in local.default_rules : merge(rule, {
      objectMappings = [
        for mapping in rule.objectMappings : merge(mapping, {
          enabled           = mapping.sourceObjectName == "Group" ? false : mapping.enabled
          attributeMappings = concat(mapping.attributeMappings, [for extra in [local.custom_attribute_mapping] : extra if mapping.sourceObjectName == "User"])
        })
      ]
    })
  ]

  directories = [
    for directory in local.default_directories : merge(directory, {
      objects = [
        for object in directory.objects : merge(object, {
          attributes = concat(object.attributes, [for extra in [local.custom_attribute] : extra if directory.name == "customappsso" && endswith(object.name, ":User")])
        })
      ]
    })
  ]
}

resource "microsoft365_graph_beta_applications_service_principal_synchronization_schema_json" "example" {
  service_principal_id = "00000000-0000-0000-0000-000000000000"
  job_id               = "scim.00000000000000000000000000000000.00000000-0000-0000-0000-000000000000"

  synchronization_rules_json = jsonencode(local.synchronization_rules)
  directories_json           = jsonencode(local.directories)

  timeouts = {
    create = "10m"
    read   = "5m"
    update = "10m"
    delete = "5m"
  }
}
