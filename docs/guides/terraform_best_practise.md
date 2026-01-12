---
page_title: "Terraform Best Practices"
subcategory: "Guides"
description: |-
  Best practices for operating the Microsoft 365 Terraform provider at scale, including performance, concurrency, and reliability considerations.
---

# Microsoft 365 Provider execution Best Practices

This guide covers operational best practices when using the Microsoft 365 Terraform provider, focusing on parallelism, API throttling, and handling large workloads efficiently.

!> Following these recommendations will help you avoid common operational issues when managing Microsoft 365 resources at scale, particularly for large deployments or when working with file-based resources.

## Understanding Parallelism in Terraform

By default, Terraform runs up to 10 operations concurrently. While this parallelism improves performance in many scenarios, it can cause issues with the Microsoft 365 provider due to:

1. Microsoft Graph API throttling limits
2. Header map collision issues in the Kiota SDK HTTP client
3. Resource contention when uploading and processing large files

## Recommended Terraform Settings

### Setting Appropriate Parallelism

For Microsoft 365 provider operations, we strongly recommend reducing Terraform's default parallelism to 1 to prevent API throttling and client-side errors:

```bash
# Run terraform with reduced parallelism
terraform apply -parallelism=1
```

### When to Limit Parallelism

**Always use `-parallelism=1` for both `terraform plan` and `terraform apply` operations.**

Microsoft Graph imposes [varying throttling limits](https://learn.microsoft.com/en-us/graph/throttling-limits) across different services and operation types. During a typical Terraform workflow:

- **Planning phase**: Terraform performs multiple GET requests to read current resource state, consuming API quota
- **Apply phase**: Terraform executes CREATE (POST), UPDATE (PATCH/PUT), and DELETE operations, each with their own throttling limits
- **Refresh operations**: Concurrent read operations can trigger service-specific throttling before write operations even begin

Since all operation types (GET, POST, PATCH, PUT, DELETE) count toward service-specific throttling limits, and these limits vary significantly across Microsoft Graph services (e.g., 150 requests/minute for security alerts vs. 10,000 requests/10 minutes for Outlook), running operations concurrently can quickly exhaust API quotas and trigger 429 (Too Many Requests) errors.

By limiting parallelism to 1, you ensure operations execute sequentially, staying well within all throttling limits and avoiding both client-side errors and server-side API throttling.

## Environment Variable Configuration

You can set environment variables to enforce parallelism limits for all Terraform operations:

```bash
# Enforce sequential operations for both plan and apply
export TF_CLI_ARGS_apply="-parallelism=1"
export TF_CLI_ARGS_plan="-parallelism=1"

# Optional: Configure retry behavior for transient errors
export M365_MAX_RETRY_ATTEMPTS=15
export M365_RETRY_MAX_DELAY_MS=60000
export M365_RETRY_MIN_DELAY_MS=1000
```

Setting these environment variables ensures that all Terraform commands respect Microsoft Graph API throttling limits without requiring explicit flags on each command.

## Continuous Integration Considerations

For CI/CD pipelines, add explicit parallelism settings:

```yaml
# Example GitHub Actions workflow
jobs:
  terraform:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4.2.2
      - uses: hashicorp/setup-terraform@v3.1.2
      
      - name: Terraform Init
        run: terraform init
      
      - name: Terraform Plan
        run: terraform plan -parallelism=1
        env:
          # Authentication variables...
      
      - name: Terraform Apply
        run: terraform apply -auto-approve -parallelism=1
        env:
          # Authentication variables...
```

## State Management Best Practices

### Use Remote State Storage

**Always use remote state backends for production environments.** HashiCorp strongly recommends storing Terraform state remotely to enable collaboration, state locking, and secure access control.

```hcl
terraform {
  backend "azurerm" {
    resource_group_name  = "terraform-state-rg"
    storage_account_name = "tfstateaccount"
    container_name       = "tfstate"
    key                  = "m365.terraform.tfstate"
  }
}
```

Remote state provides:
- **State locking**: Prevents concurrent modifications that could corrupt state
- **Encryption at rest**: Protects sensitive data in state files
- **Team collaboration**: Enables multiple team members to work on the same infrastructure
- **Audit trail**: Tracks who made changes and when

### Enable Refresh-Only Mode for Drift Detection

Use `terraform plan -refresh-only` to detect configuration drift without proposing changes. This is particularly useful for monitoring Microsoft 365 resources that may be modified outside Terraform:

```bash
# Check for drift without planning changes
terraform plan -refresh-only -parallelism=1

# Update state to match real-world resources
terraform apply -refresh-only -auto-approve -parallelism=1
```

## Credential Security

### Never Hardcode Credentials

**Never store credentials in Terraform configuration files or version control.** Use one of these secure methods:

1. **Environment variables** (recommended for local development):
```bash
export M365_CLIENT_ID="your-client-id"
export M365_CLIENT_SECRET="your-client-secret"
export M365_TENANT_ID="your-tenant-id"
```

2. **Terraform Cloud/Enterprise Variables** (recommended for CI/CD)
3. **Azure Key Vault** or other secrets management solutions
4. **Managed identities** when running in Azure environments

### Apply Least Privilege Principle

Grant only the minimum required Microsoft Graph API permissions to your service principal. Review and audit permissions regularly:

- Start with read-only permissions during initial development
- Add write permissions only for resources you need to manage
- Use application permissions (not delegated) for unattended operations
- Document required permissions in your project README

## Monitoring and Logging

### Enable Debug Logging for Troubleshooting

When investigating issues with Microsoft Graph API calls, enable debug logging to capture detailed request/response information:

```bash
# Enable Terraform debug logging
export TF_LOG=DEBUG
export TF_LOG_PATH="./terraform.log"

# Enable provider-specific logging
export M365_DEBUG=true

terraform apply -parallelism=1
```

**Important**: Debug logs contain sensitive information including access tokens. Never commit log files to version control or share them publicly.

### Monitor API Rate Limit Headers

The provider automatically handles retry logic, but you can monitor rate limit consumption by reviewing the `Retry-After` headers in debug logs. This helps identify resources that consume significant API quota.

## Advanced Configurations for Large Environments

For very large Microsoft 365 tenants (1000+ resources), consider these additional strategies:

1. **Use multiple workspaces**: Split resources across multiple Terraform workspaces by service area (e.g., Intune policies, Conditional Access, users)
2. **Implement resource tagging**: Use consistent naming and tagging conventions for easier management
3. **Use multiple client IDs**: Distribute operations across multiple registered applications with separate throttling quotas
4. **Schedule deployments during off-peak hours**: Reduce contention with other automated systems
5. **Implement phased rollouts**: Deploy changes in stages to minimize impact and facilitate rollback if needed

## Summary

**Always use `-parallelism=1` when working with the Microsoft 365 Terraform provider.** Microsoft Graph's [variable throttling limits](https://learn.microsoft.com/en-us/graph/throttling-limits) across services mean that concurrent operations—whether during `terraform plan` (GET operations) or `terraform apply` (POST/PATCH/DELETE operations)—can quickly exhaust API quotas and trigger 429 errors.

Setting `-parallelism=1` helps avoid:
- Server-side API throttling across all Microsoft Graph services
- Client-side header collision issues in the Kiota SDK
- Resource contention during file uploads and processing

While this approach increases total deployment time for large configurations, it substantially improves reliability, reduces troubleshooting efforts, and ensures consistent operation regardless of which Microsoft Graph services your configuration uses.

## References

- [Microsoft Graph Throttling Limits](https://learn.microsoft.com/en-us/graph/throttling-limits) - Official documentation on service-specific API rate limits
- [Microsoft Graph Best Practices](https://learn.microsoft.com/en-us/graph/best-practices-concept) - Microsoft's recommended patterns for Graph API usage
- [Terraform State Management](https://developer.hashicorp.com/terraform/language/state) - HashiCorp documentation on state files and remote backends
- [Terraform Security Best Practices](https://developer.hashicorp.com/terraform/cloud-docs/recommended-practices/part2) - HashiCorp's guide to secure Terraform usage
- [Terraform Provider Design Principles](https://developer.hashicorp.com/terraform/plugin/hashicorp-provider-design-principles) - HashiCorp's principles for provider development
- [Terraform Recommended Practices](https://developer.hashicorp.com/terraform/cloud-docs/recommended-practices) - HashiCorp's comprehensive guide to Terraform workflows