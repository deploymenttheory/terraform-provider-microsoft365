---
page_title: "Batch Import with List Resources and Terraform Generate - terraform-provider-microsoft365"
subcategory: "Guides"
description: |-
  Step-by-step guide for discovering and importing existing Microsoft 365 resources using list resources and Terraform's configuration generation feature.
---

# Batch Import with List Resources and Terraform Generate

This guide demonstrates how to discover existing Microsoft 365 resources using list resources and automatically generate Terraform configuration using the `terraform plan -generate-config-out` flag.

## Prerequisites

- Terraform 1.14.0 or later (required for list resource support)
- Microsoft 365 provider configured with appropriate credentials
- Read access to the resources you want to discover and import

## Overview

The workflow consists of three phases:

1. **Discover** - Use list resources to query and filter existing infrastructure
2. **Select** - Choose which resources to import based on list results
3. **Generate** - Use Terraform to automatically create HCL configuration from selected resources

## Step 1: Discover Resources

Create a `.tfquery.hcl` file to query existing resources. Use filters to narrow results to your target resources.

### Example: Discover Settings Catalog Policies

Create `discover.tfquery.hcl`:

```hcl
# List all Windows 10 baseline policies
list "microsoft365_graph_beta_device_management_settings_catalog_configuration_policy" "baseline_policies" {
  provider = microsoft365
  config {
    platform_filter         = ["windows10"]
    template_family_filter  = "baseline"
  }
}
```

Execute the query:

```bash
terraform query
```

### Review Results

Terraform returns a structured list of matching resources with their IDs and names:

```
list.microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.baseline_policies   id=892f843c-4660-4cac-9f94-a94265c72c8f   Windows Security Baseline
```

## Step 2: Create Import Blocks

Based on the list results, create import blocks for the resources you want to manage with Terraform.

### Define Import Targets

Create `import.tf`:

```hcl
import {
  to = microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.windows_baseline
  id = "892f843c-4660-4cac-9f94-a94265c72c8f"
}

import {
  to = microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.edge_baseline
  id = "a1b2c3d4-5678-90ab-cdef-1234567890ab"
}

import {
  to = microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.defender_baseline
  id = "b2c3d4e5-6789-01bc-def0-234567890abc"
}
```

### Import Block Structure

Each import block requires:
- `to` - The Terraform resource address (resource type + name)
- `id` - The unique identifier from the list resource output

## Step 3: Generate Configuration

Use Terraform's configuration generation feature to automatically create HCL from the imported resources.

### Execute Plan with Generation

Run the plan command with the `-generate-config-out` flag:

```bash
terraform plan -generate-config-out=generated.tf
```

### What Happens

Terraform performs the following actions:

1. Reads each import block from `import.tf`
2. Queries the Microsoft Graph API for the current state of each resource
3. Converts the API response to HCL syntax
4. Writes the configuration to `generated.tf`

### Review Generated Configuration

Examine `generated.tf` to verify the generated configuration:

```hcl
resource "microsoft365_graph_beta_device_management_settings_catalog_configuration_policy" "windows_baseline" {
  name        = "Windows Security Baseline"
  description = "Microsoft security baseline for Windows 10"
  platforms   = ["windows10"]
  technologies = ["mdm"]
  
  settings = [
    {
      # ... generated settings configuration
    }
  ]
}

# Additional imported resources...
```

## Step 4: Apply the Import

Execute the import by applying the configuration:

```bash
terraform apply
```

Terraform displays the import plan. Review the changes and confirm:

```
Plan: 3 to import, 0 to add, 0 to change, 0 to destroy.

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes
```

After confirmation, Terraform imports the resources into state without modifying them.

## Step 5: Validate and Refine

### Verify State

Confirm all resources are in Terraform state:

```bash
terraform state list
```

Expected output:

```
microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.windows_baseline
microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.edge_baseline
microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.defender_baseline
```

### Remove Import Blocks

After successful import, remove the import blocks from `import.tf`:

```bash
# Remove the import.tf file
rm import.tf

# Or comment out import blocks if keeping for reference
```

### Verify No Configuration Drift

Re run `terraform plan` to confirm the imported resources match the generated configuration with no differences:

```bash
terraform plan
```

Expected output:

```
No changes. Your infrastructure matches the configuration.

Terraform has compared your real infrastructure against your configuration
and found no differences, so no changes are needed.
```

If the plan shows differences, review and adjust the generated configuration to match the actual resource state.

### Refine Generated Configuration

Review and adjust `generated.tf` as needed:

- Remove default or null values
- Reorder fields for readability, imported fields are always imported alphabetically by default.
- Add variables for dynamic values
- Reorganize into modules if importing many resources
- Add comments for documentation

## Optimizing Generated Configuration

Generated configuration is functional but not optimized. Terraform generates verbose, literal configurations 
that require manual refinement for production use. Values are exactly what the API returns.

### Generated Configuration Limitations

**Output Example**: Note:Terraform does not generate meta-arguments like `for_each` or `count`. Each resource is explicitly defined:

```hcl
# __generated__ by Terraform from "9af36d64-2294-4c46-8f78-f30bf4d17061"
resource "microsoft365_graph_beta_device_management_settings_catalog_configuration_policy" "kerberos_policy" {
  assignments = [
    {
      filter_id   = "00000000-0000-0000-0000-000000000000"  # <----- Repeated literal value, use local
      filter_type = "none"                                  # <----- Repeated literal value, use local
      group_id    = "1c4f3adf-ebe8-422c-97b1-f174632d7538" # <----- Literal ID, no dependency on group resource
      type        = "groupAssignmentTarget"                 # <----- Repeated literal value, use local
    },
  ]
  configuration_policy = {
    settings = [
      {
        id = "0"
        setting_instance = {
          choice_setting_collection_value     = null        # <----- Null values from API, can be safely removed
          choice_setting_value                = null        # <----- Null values from API, can be safely removed
          group_setting_collection_value      = null        # <----- Null values from API, can be safely removed
          odata_type                          = "#microsoft.graph.deviceManagementConfigurationSimpleSettingCollectionInstance"
          setting_definition_id               = "device_vendor_msft_policy_config_kerberos_upnnamehints"
          setting_instance_template_reference = null        # <----- Null values from API, can be safely removed
          simple_setting_collection_value = [
            {
              odata_type                       = "#microsoft.graph.deviceManagementConfigurationStringSettingValue"
              setting_value_template_reference = null       # <----- Null values from API, can be safely removed
              value                            = "deploymenttheory.com"
            },
          ]
          simple_setting_value = null                       # <----- Null values from API, can be safely removed
        }
      },
    ]
  }
  name               = "[Base] Prod | Windows - Settings Catalog | Kerberos ver1.0"
  platforms          = "windows10"                          # <----- Repeated value across policies, use local
  role_scope_tag_ids = ["5"]                                # <----- Literal scope tag ID, no dependency on scope tag resource
  technologies       = ["mdm"]                              # <----- Repeated value across policies, use local
}

# __generated__ by Terraform from "4a126db0-997f-4156-98ed-07e4449730b8"
resource "microsoft365_graph_beta_device_management_settings_catalog_configuration_policy" "licensing_policy" {
  assignments = [
    {
      filter_id   = "00000000-0000-0000-0000-000000000000"  # <----- Same values as above, no shared definition
      filter_type = "none"                                  # <----- Same values as above, no shared definition
      group_id    = "1c4f3adf-ebe8-422c-97b1-f174632d7538"  # <----- Same group ID as above, should reference same resource
      type        = "groupAssignmentTarget"                 # <----- Same values as above, no shared definition
    },
  ]
  # ... additional configuration
  platforms          = "windows10"                          # <----- Duplicate values across all policies
  role_scope_tag_ids = ["0"]                                # <----- Duplicate values across all policies
  technologies       = ["mdm"]                              # <----- Duplicate values across all policies
}
```

**Literal Values**: Generated configuration uses actual values from the API rather than references to other resources:

```hcl
# __generated__ by Terraform
resource "microsoft365_graph_beta_device_management_settings_catalog_configuration_policy" "kerberos_policy" {
  assignments = [
    {
      group_id    = "1c4f3adf-ebe8-422c-97b1-f174632d7538" # <----- Literal group ID, breaks if group changes
      type        = "groupAssignmentTarget"
    },
  ]
  role_scope_tag_ids = ["0"]                                # <----- Literal scope tag ID, no relationship to actual tag resource
  platforms          = "windows10"                          # <----- Could reference variable for consistency
  technologies       = ["mdm"]
}
```

**Missing Dependencies**: Resource relationships are not expressed, so Terraform cannot determine the correct order of operations:

```hcl
# __generated__ by Terraform - These resources exist but no relationship is defined
resource "microsoft365_graph_beta_groups_group" "production_users" {
  display_name     = "Production Users"
  id               = "1c4f3adf-ebe8-422c-97b1-f174632d7538"
  # ... other fields
}

resource "microsoft365_graph_beta_device_management_role_scope_tag" "production" {
  display_name = "Production"
  id           = "5"
  # ... other fields
}

resource "microsoft365_graph_beta_device_management_settings_catalog_configuration_policy" "kerberos_policy" {
  assignments = [
    {
      group_id = "1c4f3adf-ebe8-422c-97b1-f174632d7538" # <----- Literal ID, no dependency on group resource above
    },
  ]
  role_scope_tag_ids = ["0"]                            # <----- Literal ID, no dependency on scope tag resource above
}

# Terraform's dependency graph cannot determine:
# - That the policy depends on the group existing
# - That the policy depends on the scope tag existing
# - The correct order to create/destroy these resources
```

### Refactoring for Production

Transform generated configuration for maintainability and correct dependency management.

#### Use for_each for Similar Resources

Convert repetitive resources to use `for_each`:

```hcl
# Refactored: Use for_each with a map
locals {
  security_policies = {
    bitlocker = {
      name        = "BitLocker Policy"
      description = "Production BitLocker configuration"
    }
    defender = {
      name        = "Defender Policy"
      description = "Production Defender configuration"
    }
    firewall = {
      name        = "Firewall Policy"
      description = "Production Firewall configuration"
    }
  }
}

resource "microsoft365_graph_beta_device_management_settings_catalog_configuration_policy" "security" {
  for_each = local.security_policies
  
  name        = each.value.name
  description = each.value.description
  platforms   = ["windows10"]
  technologies = ["mdm"]
}
```

#### Express Resource Dependencies

Replace literal values with resource references to establish dependencies:

```hcl
# Refactored: Use references to establish dependencies
resource "microsoft365_graph_beta_groups_group" "security_team" {
  display_name     = "Security Team"
  mail_enabled     = false
  security_enabled = true
}

resource "microsoft365_graph_beta_device_management_settings_catalog_configuration_policy" "security_policy" {
  name = "Security Policy"
  
  assignments = [
    {
      target = {
        group_id = microsoft365_graph_beta_groups_group.security_team.id
      }
    }
  ]
}
```

Now Terraform understands that the policy depends on the group and will create them in the correct order.

#### Create Resource References for Scope Tags

Replace hardcoded scope tag IDs with references:

```hcl
# Refactored: Express dependency on scope tag
resource "microsoft365_graph_beta_device_management_role_scope_tag" "production" {
  display_name = "Production"
}

resource "microsoft365_graph_beta_device_management_settings_catalog_configuration_policy" "policy" {
  name            = "Production Policy"
  role_scope_tags = [microsoft365_graph_beta_device_management_role_scope_tag.production.id]
}
```

#### Centralize Common Values

Extract repeated values into locals or variables:

```hcl
# Refactored: Use locals for common values
locals {
  default_platforms   = ["windows10"]
  default_technologies = ["mdm"]
  production_prefix   = "[Prod]"
}

resource "microsoft365_graph_beta_device_management_settings_catalog_configuration_policy" "bitlocker" {
  name         = "${local.production_prefix} BitLocker"
  platforms    = local.default_platforms
  technologies = local.default_technologies
}

resource "microsoft365_graph_beta_device_management_settings_catalog_configuration_policy" "defender" {
  name         = "${local.production_prefix} Defender"
  platforms    = local.default_platforms
  technologies = local.default_technologies
}
```

### Benefits of Refactoring

**Correct Dependency Graph**: Terraform understands the relationships between resources and manages them in the correct order.

**Maintainability**: Changes to shared values (like group IDs or scope tags) only need to be updated in one place.

**Reduced Errors**: Dependencies prevent resources from being destroyed or modified in the wrong order.

**Version Control**: Smaller, more focused diffs make changes easier to review.

## Advanced Filtering Techniques

### Multiple Filters

Combine filters to precisely target resources:

```hcl
list "microsoft365_graph_beta_device_management_settings_catalog_configuration_policy" "assigned_defender" {
  provider = microsoft365
  config {
    name_filter        = "Defender"
    platform_filter    = ["windows10"]
    is_assigned_filter = true
  }
}
```

### Custom OData Queries

Use OData filters for complex queries:

```hcl
list "microsoft365_graph_beta_device_management_settings_catalog_configuration_policy" "prod_policies" {
  provider = microsoft365
  config {
    odata_filter = "startsWith(name, '[Prod]') and platforms eq 'windows10'"
  }
}
```

## Batch Import Workflow Example

### Complete End-to-End Process

**Step 1**: Query production policies

```bash
# Create discover.tfquery.hcl
cat > discover.tfquery.hcl << 'EOF'
list "microsoft365_graph_beta_device_management_settings_catalog_configuration_policy" "prod" {
  provider = microsoft365
  config {
    odata_filter = "startsWith(name, '[Prod]')"
  }
}
EOF

# Execute query
terraform query > policy_list.txt
```

**Step 2**: Extract IDs and create imports

```bash
# Parse list output and create import blocks
# (Manual selection based on business requirements)
cat > import.tf << 'EOF'
import {
  to = microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.prod_bitlocker
  id = "id-from-query-output"
}
import {
  to = microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.prod_defender
  id = "id-from-query-output"
}
EOF
```

**Step 3**: Generate configuration

```bash
terraform plan -generate-config-out=generated_policies.tf
```

**Step 4**: Review and import

```bash
# Review generated configuration
cat generated_policies.tf

# Apply import
terraform apply

# Verify
terraform state list
```

**Step 5**: Organize configuration

```bash
# Move generated config to organized structure
mkdir -p policies/security
mv generated_policies.tf policies/security/baselines.tf

# Update references if needed
terraform init
terraform plan
```

## Supported List Resources

The following list resources support this workflow:

| Resource Type | Use Case |
|---------------|----------|
| `microsoft365_graph_beta_device_management_settings_catalog_configuration_policy` | Settings Catalog policies |

## Best Practices

### Test with Small Batches First

Start with small, targeted imports to validate your workflow before attempting large-scale imports:

```hcl
# Good: Start with a single policy or small subset
config {
  name_filter = "Kerberos"  # Returns 1-2 policies for testing
}

# After validating the process, expand to larger batches
config {
  name_filter     = "Security"
  platform_filter = ["windows10"]  # Returns 10-20 policies
}

# Finally, proceed with full imports
config {
  platform_filter = ["windows10"]  # Returns all policies for platform
}
```

### Name Resources Consistently

Use descriptive, consistent naming in import blocks:

```hcl
# Good: Clear, descriptive names
import {
  to = microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.prod_win10_bitlocker
  id = "..."
}

# Avoid: Generic names
import {
  to = microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.policy1
  id = "..."
}
```

### Version Control Generated Files

Commit generated configuration to version control after review and refinement:

```bash
# Review changes
terraform plan

# Commit generated configuration
git add generated.tf import.tf
git commit -m "Import production security baselines"
```

### Separate Discovery from Management

Keep discovery queries separate from resource management:

```
.
├── discover.tfquery.hcl    # List resource queries (not applied)
├── import.tf               # Import blocks (temporary)
├── generated.tf            # Generated configuration (permanent)
├── main.tf                 # Primary configuration
└── variables.tf            # Variable definitions
```

## Troubleshooting

### Import Block Errors

**Error**: `Missing Resource Identity After Read`

**Cause**: The resource doesn't implement the required identity schema.

**Solution**: Verify the resource type supports import and update to the latest provider version.

### Generated Configuration Issues

**Error**: Generated configuration contains invalid syntax

**Cause**: Complex nested structures may require manual adjustment.

**Solution**: Review and refine the generated configuration:

1. Run `terraform validate` to identify syntax errors
2. Adjust complex nested blocks manually
3. Test with `terraform plan` before committing

### List Resource Returns No Results

**Cause**: Filters are too restrictive or resource doesn't exist.

**Solution**: Simplify filters progressively:

```hcl
# Start broad
config {}

# Add filters incrementally
config {
  platform_filter = ["windows10"]
}

# Refine further
config {
  platform_filter = ["windows10"]
  name_filter     = "Security"
}
```

## Additional Resources

- [Terraform Import Documentation](https://developer.hashicorp.com/terraform/language/import)
- [Terraform List Resources](https://developer.hashicorp.com/terraform/plugin/framework/list-resources)
- [Microsoft Graph API Documentation](https://learn.microsoft.com/en-us/graph/overview)

