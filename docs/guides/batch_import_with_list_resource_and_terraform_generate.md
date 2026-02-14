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

## Supported List Resources

The following list resources support this workflow:

| Resource Type | Use Case |
|---------------|----------|
| `microsoft365_graph_beta_device_management_settings_catalog_configuration_policy` | Settings Catalog policies |
| `microsoft365_graph_beta_identity_and_access_conditional_access_policy` | Conditional Access policies |
| `microsoft365_graph_beta_users_user` | Users |

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

## Advanced Filtering Techniques

### Multiple Filters

Each list resource supports different filters. Combine filters to precisely target resources:

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

