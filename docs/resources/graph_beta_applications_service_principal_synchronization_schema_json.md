---
page_title: "microsoft365_graph_beta_applications_service_principal_synchronization_schema_json Resource - terraform-provider-microsoft365"
subcategory: "Applications"
description: |-
  Manages the synchronization (provisioning) schema of a service principal's synchronization job using the /servicePrincipals/{id}/synchronization/jobs/{jobId}/schema endpoint. The schema defines which objects are provisioned, the attribute mappings and scoping filters, and the directory definitions of the source and target systems. The schema is sent in full with PUT on every create and update, so the configuration must contain the complete synchronizationRules and directories collections. Destroying the resource calls DELETE on the schema, which resets it to the default schema of the job's synchronization template.
---

# microsoft365_graph_beta_applications_service_principal_synchronization_schema_json (Resource)

Manages the synchronization (provisioning) schema of a service principal's synchronization job using the `/servicePrincipals/{id}/synchronization/jobs/{jobId}/schema` endpoint. The schema defines which objects are provisioned, the attribute mappings and scoping filters, and the directory definitions of the source and target systems. The schema is sent in full with `PUT` on every create and update, so the configuration must contain the complete `synchronizationRules` and `directories` collections. Destroying the resource calls `DELETE` on the schema, which resets it to the default schema of the job's synchronization template.

## Microsoft Documentation

- [synchronizationSchema resource type](https://learn.microsoft.com/en-us/graph/api/resources/synchronization-synchronizationschema?view=graph-rest-beta)
- [Get synchronizationSchema](https://learn.microsoft.com/en-us/graph/api/synchronization-synchronizationschema-get?view=graph-rest-beta)
- [Update synchronizationSchema](https://learn.microsoft.com/en-us/graph/api/synchronization-synchronizationschema-update?view=graph-rest-beta)
- [Delete synchronizationSchema](https://learn.microsoft.com/en-us/graph/api/synchronization-synchronizationschema-delete?view=graph-rest-beta)
- [Configure synchronization with directory extension attributes](https://learn.microsoft.com/en-us/graph/synchronization-configure-with-directory-extension-attributes)

## Microsoft Graph API Permissions

The following client `application` permissions are needed in order to use this resource:

**Required:**
- `Synchronization.Read.All`
- `Application.ReadWrite.OwnedBy` (the calling application must own the service principal) or `Synchronization.ReadWrite.All`

**Optional:**
- `None` `[N/A]`

## Version History

| Version | Status | Notes |
|---------|--------|-------|
| v1.7.0 | Experimental | Initial release |

## Example Usage

### Minimal

```terraform
# Manages the provisioning schema of an existing synchronization job.
#
# Export the job's current schema with
#   GET https://graph.microsoft.com/beta/servicePrincipals/{id}/synchronization/jobs/{jobId}/schema
# and save its "synchronizationRules" and "directories" arrays as JSON files next to this configuration.
resource "microsoft365_graph_beta_applications_service_principal_synchronization_schema_json" "example" {
  service_principal_id = "00000000-0000-0000-0000-000000000000"
  job_id               = "scim.00000000000000000000000000000000.00000000-0000-0000-0000-000000000000"

  synchronization_rules_json = file("${path.module}/synchronization_rules.json")
  directories_json           = file("${path.module}/directories.json")
}
```

### Maximal

```terraform
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
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `directories_json` (String) The complete `directories` collection of the schema as a JSON array, for example from `jsonencode()` or `file()`. Each directory definition lists the objects and attributes that can take part in the synchronization rules. The value is compared semantically: whitespace and property order are ignored, and so are the `version` properties that Microsoft Graph updates on every schema change.
- `job_id` (String) The ID of the synchronization job whose schema is managed, for example `scim.{tenantId}.{appId}`. Changing it replaces the resource.
- `service_principal_id` (String) The object ID of the service principal that owns the synchronization job. Changing it replaces the resource.
- `synchronization_rules_json` (String) The complete `synchronizationRules` collection of the schema as a JSON array, for example from `jsonencode()` or `file()`. Each rule defines the source and target directories, object mappings, attribute mappings and scoping filters. The value is compared semantically: whitespace and property order are ignored, and so are the `Unsynchronized` object mapping metadata entries that Microsoft Graph derives from `enabled`.

### Optional

- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))

### Read-Only

- `id` (String) The ID of the resource in the format `{service_principal_id}/{job_id}`.
- `version` (String) The version of the schema, updated by Microsoft Graph with every schema change.

<a id="nestedatt--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).
- `delete` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Setting a timeout for a Delete operation is only applicable if changes are saved into state before the destroy operation occurs.
- `read` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Read operations occur during any refresh or planning operation when refresh is enabled.
- `update` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).

## Important Notes

### Existing Synchronization Job

- The resource manages the schema of an existing synchronization (provisioning) job. Create the job first, for example in the Microsoft Entra admin center or with `POST /servicePrincipals/{id}/synchronization/jobs`.
- The schema always exists while the job exists, so creating the resource replaces the job's current schema.

### Full Replacement

- Every create and update sends the complete schema with `PUT`. Configure the complete `synchronizationRules` and `directories` collections, not only the parts you change.
- The easiest starting point is the job's current schema from `GET /servicePrincipals/{id}/synchronization/jobs/{jobId}/schema`, edited with `jsondecode()` and `jsonencode()` as in the maximal example.

### Semantic Comparison

- The JSON values are compared semantically. Whitespace and property order are ignored.
- Microsoft Graph changes every directory `version` on each schema change and sets the `Unsynchronized` object mapping metadata from the mapping's `enabled` flag. These values are ignored, so they never cause drift.
- Array order is significant, because the order of expression parameters matters.
- Changing only the formatting of the JSON plans a one-time update that sends the same schema again.

### Import

- After an import, state holds the schema as returned by Microsoft Graph, encoded like `jsonencode()`. Unless the configuration produces exactly the same string, the first plan shows an in-place update that sends the configured schema once.

### Deletion

- Destroying the resource calls `DELETE` on the schema, which resets it to the default schema of the job's synchronization template. The job itself is not deleted.

## Import

Import is supported using the following syntax:

```shell
#!/bin/bash
# Import using the service principal object ID and the synchronization job ID, separated by a forward slash.
terraform import microsoft365_graph_beta_applications_service_principal_synchronization_schema_json.example 00000000-0000-0000-0000-000000000000/scim.00000000000000000000000000000000.00000000-0000-0000-0000-000000000000
```
