---
page_title: "Terraform Best Practices"
subcategory: "Operations"
description: |-
  Best practices for operating the Microsoft 365 Terraform provider at scale, including performance, concurrency, and reliability considerations.
---

# Microsoft 365 Provider Operations Best Practices

This guide covers operational best practices when using the Microsoft 365 Terraform provider, focusing on parallelism, API throttling, and handling large workloads efficiently.

> [!NOTE]
> Following these recommendations will help you avoid common operational issues when managing Microsoft 365 resources at scale, particularly for large deployments or when working with file-based resources.

## Understanding Parallelism in Terraform

By default, Terraform runs up to 10 operations concurrently. While this parallelism improves performance in many scenarios, it can cause issues with the Microsoft 365 provider due to:

1. Microsoft Graph API throttling limits
2. Header map collision issues in the Graph SDK HTTP client
3. Resource contention when uploading and processing large files

## Recommended Terraform Settings

### Setting Appropriate Parallelism

For Microsoft 365 provider operations, we strongly recommend reducing Terraform's default parallelism to 1 prevent API throttling and client-side errors:

```bash
# Run terraform with reduced parallelism
terraform apply -parallelism=1
```

### When to Limit Parallelism

Limit parallelism specifically when:

- It's encouraged to appply this during all terraform apply's.

## Environment Variable Configuration

You can set environment variables to adjust provider behavior:

```bash
# Use a single operation at a time
export TF_CLI_ARGS_apply="-parallelism=1"
export TF_CLI_ARGS_plan="-parallelism=1"

# More granular configuration
export M365_MAX_RETRY_ATTEMPTS=15
export M365_RETRY_MAX_DELAY_MS=60000
export M365_RETRY_MIN_DELAY_MS=1000
```

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
      
      - name: Terraform Apply
        run: terraform apply -auto-approve -parallelism=1
        env:
          # Authentication variables...
```

## Advanced Configurations for Large Environments

For very large Microsoft 365 tenants, consider these additional strategies:

1. **Use multiple client IDs**: Distribute operations across multiple registered applications
2. **Implement rate limiting in your code**: Add deliberate delays between resource operations
3. **Monitor API usage**: Track Graph API usage metrics to optimize your approach
4. **Batch similar operations**: Group similar resource types together

## Summary

Microsoft 365 Graph API operations, especially those involving file uploads and multiple polling cycles, benefit significantly from limiting Terraform parallelism. Setting `-parallelism=1` helps avoid both client-side header collision issues and server-side API throttling, resulting in more reliable deployments.

While this approach may increase total deployment time for large configurations, it substantially improves reliability and reduces troubleshooting efforts.