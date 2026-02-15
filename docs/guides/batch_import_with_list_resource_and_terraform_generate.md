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

## Supported List Resources

The following list resources support this workflow:

| Resource Type | Use Case |
|---------------|----------|
| `microsoft365_graph_beta_device_management_settings_catalog_configuration_policy` | Settings Catalog policies |
| `microsoft365_graph_beta_device_management_windows_platform_script` | Windows PowerShell scripts |
| `microsoft365_graph_beta_identity_and_access_conditional_access_policy` | Conditional Access policies |
| `microsoft365_graph_beta_users_user` | Users |

## Overview

The workflow consists of three phases:

1. **Discover** - Query existing resources using list resources with filters
2. **Generate** - Create HCL configuration from discovered resources
3. **Import** - Bring resources under Terraform management

## Example: Importing Settings Catalog Policies

This example demonstrates importing Windows 10 baseline policies that exist in Intune but aren't managed by Terraform.

### Step 1: Discover Resources

Create `discover.tfquery.hcl` to query for baseline policies:

```hcl
# Query for Windows 10 baseline policies
list "microsoft365_graph_beta_device_management_settings_catalog_configuration_policy" "baseline_policies" {
  provider = microsoft365
  config {
    platform_filter        = ["windows10"]
    template_family_filter = "baseline"
  }
}
```

Execute the query:

```bash
terraform query
```

Output shows discovered policies:

```
list.microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.baseline_policies
  id=892f843c-4660-4cac-9f94-a94265c72c8f   Windows Security Baseline
list.microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.baseline_policies
  id=a1b2c3d4-5678-90ab-cdef-1234567890ab   Edge Baseline  
list.microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.baseline_policies
  id=b2c3d4e5-6789-01bc-def0-234567890abc   Defender Baseline
```

### Step 2: Generate Configuration

Generate Terraform configuration for all discovered policies:

```bash
terraform plan -generate-config-out=baselines.tf \
  -target 'list.microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.baseline_policies'
```

This creates `baselines.tf` containing:
- `resource` blocks with complete configuration
- `import` blocks with resource identities

### Step 3: Review Generated Configuration

The generated file contains verbose configuration. Review and refine:

```hcl
# __generated__ by Terraform
resource "microsoft365_graph_beta_device_management_settings_catalog_configuration_policy" "windows_baseline" {
  name         = "Windows Security Baseline"
  description  = "Microsoft security baseline for Windows 10"
  platforms    = ["windows10"]
  technologies = ["mdm"]
  
  settings = [
    # ... generated settings
  ]
  
  assignments = [
    {
      target = {
        group_id = "1c4f3adf-ebe8-422c-97b1-f174632d7538"
      }
    }
  ]
}

import {
  to = microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.windows_baseline
  id = "892f843c-4660-4cac-9f94-a94265c72c8f"
}

# Additional policies...
```

**Common refinements:**
- Remove null/default values
- Replace literal IDs with resource references
- Extract common values to locals
- Add comments for documentation

### Step 4: Import Resources

Apply the configuration to import resources into state:

```bash
terraform apply
```

Terraform displays the import plan:

```
Plan: 3 to import, 0 to add, 0 to change, 0 to destroy.
```

After confirmation, resources are imported without modification.

### Step 5: Verify and Clean Up

Verify resources are in state:

```bash
terraform state list
```

Remove the import blocks from `baselines.tf` after successful import.

Verify no drift:

```bash
terraform plan
```

Expected output: `No changes. Your infrastructure matches the configuration.`

## Refining Generated Configuration

Generated configuration is literal and verbose. Transform it for production use.

### Express Dependencies

Replace literal IDs with resource references:

```hcl
# Before: Literal group ID
resource "microsoft365_graph_beta_device_management_settings_catalog_configuration_policy" "policy" {
  assignments = [
    {
      target = {
        group_id = "1c4f3adf-ebe8-422c-97b1-f174632d7538"  # Hard-coded ID
      }
    }
  ]
}

# After: Resource reference
resource "microsoft365_graph_beta_groups_group" "security_team" {
  display_name     = "Security Team"
  security_enabled = true
}

resource "microsoft365_graph_beta_device_management_settings_catalog_configuration_policy" "policy" {
  assignments = [
    {
      target = {
        group_id = microsoft365_graph_beta_groups_group.security_team.id
      }
    }
  ]
}
```

### Use for_each for Similar Resources

Convert repetitive resources to use `for_each`:

```hcl
locals {
  security_policies = {
    bitlocker = "BitLocker Policy"
    defender  = "Defender Policy"
    firewall  = "Firewall Policy"
  }
}

resource "microsoft365_graph_beta_device_management_settings_catalog_configuration_policy" "security" {
  for_each = local.security_policies
  
  name         = each.value
  platforms    = ["windows10"]
  technologies = ["mdm"]
}
```

### Centralize Common Values

Extract repeated values:

```hcl
locals {
  default_platforms   = ["windows10"]
  default_technologies = ["mdm"]
}

resource "microsoft365_graph_beta_device_management_settings_catalog_configuration_policy" "policy" {
  name         = "Policy Name"
  platforms    = local.default_platforms
  technologies = local.default_technologies
}
```

## Advanced Filtering

### Multiple Filters

Combine filters to target specific resources:

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

## Best Practices

### Start Small

Test with small batches before large-scale imports:

```hcl
# Good: Start with 1-2 policies
config {
  name_filter = "Kerberos"
}

# Then expand gradually
config {
  platform_filter = ["windows10"]
}
```

### Name Resources Consistently

Use descriptive names in generated configuration:

```hcl
# Good
resource "..." "prod_win10_bitlocker" { }

# Avoid
resource "..." "policy1" { }
```

### Version Control

Commit generated configuration after review:

```bash
git add baselines.tf
git commit -m "Import production security baselines"
```

### Organize Files

Keep discovery queries separate from resource management:

```
.
├── discover.tfquery.hcl    # List resource queries
├── baselines.tf            # Generated and refined configuration
├── main.tf                 # Primary configuration
└── variables.tf            # Variables
```

## Troubleshooting

### Import Block Errors

**Error**: `Missing Resource Identity After Read`

**Solution**: Verify the resource type supports import and update to the latest provider version.

### Generated Configuration Issues

**Error**: Invalid syntax in generated configuration

**Solution**:
1. Run `terraform validate` to identify errors
2. Adjust complex nested blocks manually
3. Test with `terraform plan`

### List Resource Returns No Results

**Solution**: Simplify filters progressively:

```hcl
# Start broad
config {}

# Add filters incrementally
config {
  platform_filter = ["windows10"]
}
```

## Additional Resources

- [Terraform Import Documentation](https://developer.hashicorp.com/terraform/language/import)
- [Terraform List Resources](https://developer.hashicorp.com/terraform/plugin/framework/list-resources)
- [Microsoft Graph API Documentation](https://learn.microsoft.com/en-us/graph/overview)

