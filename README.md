# Community Terraform Provider for Microsoft 365

[![Release](https://img.shields.io/github/v/release/deploymenttheory/terraform-provider-microsoft365)](https://github.com/deploymenttheory/terraform-provider-microsoft365/releases)
[![Installs](https://img.shields.io/badge/dynamic/json?logo=terraform&label=installs&query=$.data.attributes.downloads&url=https%3A%2F%2Fregistry.terraform.io%2Fv2%2Fproviders%2F5565)](https://registry.terraform.io/providers/deploymenttheory/microsoft365)
[![Registry](https://img.shields.io/badge/registry-doc%40latest-lightgrey?logo=terraform)](https://registry.terraform.io/providers/deploymenttheory/microsoft365/latest/docs)
[![Lint Status](https://github.com/deploymenttheory/terraform-provider-microsoft365/workflows/go%20%7C%20Linter/badge.svg)](https://github.com/deploymenttheory/terraform-provider-microsoft365/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/deploymenttheory/terraform-provider-microsoft365)](https://goreportcard.com/report/github.com/deploymenttheory/terraform-provider-microsoft365)
[![codecov](https://codecov.io/gh/deploymenttheory/terraform-provider-microsoft365/graph/badge.svg?token=STSFXO1NL6)](https://codecov.io/gh/deploymenttheory/terraform-provider-microsoft365)
[![Go Version](https://img.shields.io/github/go-mod/go-version/deploymenttheory/terraform-provider-microsoft365)](https://go.dev/)
[![License](https://img.shields.io/github/license/deploymenttheory/terraform-provider-microsoft365)](LICENSE)
![Status: Tech Preview](https://img.shields.io/badge/status-experimental-EAAA32)

The community Microsoft 365 Terraform Provider allows managing environments and other resources within [Microsoft Intune](https://www.microsoft.com/en-us/security/business/endpoint-management/microsoft-intune), [Microsoft 365](https://www.microsoft.com/en-us/microsoft-365), [Microsoft Teams](https://www.microsoft.com/en-us/microsoft-teams/group-chat-software), and [Microsoft Defender](https://www.microsoft.com/en-us/security/business/microsoft-defender).

> [!WARNING]
> This code is experimental and provided solely for evaluation purposes. It is **NOT** intended for production use and may contain bugs, incomplete features, or other issues. Use at your own risk, as it may undergo significant changes without notice until it reaches general availability, and no guarantees or support is provided. By using this code, you acknowledge and agree to these conditions. Consult the documentation or contact the maintainer if you have questions or concerns.


> [!TIP]
> This is a community-driven project and is not officially supported by Microsoft.
> If you need help, want to ask questions, or connect with other users and contributors, join our community
> [Discord](https://discord.gg/Uq8zG6g7WE)

## Overview

The Community Terraform Provider for Microsoft 365 empowers workplace teams and administrators to manage their Microsoft 365 environments using Infrastructure as Code (IaC) principles. This provider bridges the gap between Terraform's powerful resource management capabilities and the extensive features of Microsoft 365, allowing for automated, version-controlled, and repeatable deployments across various Microsoft cloud services.

## Use Cases

- **Infrastructure as Code for Microsoft 365**  
  Manage Microsoft 365 configuration (users, groups, policies, device management, and more) as code, enabling version control, peer review, and repeatable deployments—just as you would for cloud infrastructure in Azure or GCP.

- **Automated, Auditable Change Management**  
  Use Terraform's plan and apply in gitOps workflows to preview, approve, and track changes to your Microsoft 365 environment, ensuring all modifications are intentional, reviewed, and logged.

- **Environment Replication and Drift Detection**
  Reproduce Microsoft 365 tenant configurations across multiple environments (development, staging, production) or tenants, and detect configuration drift over time using Terraform’s state management.

- **Disaster Recovery and Rapid Rebuilds**  
  Store your Microsoft 365 configuration in code, allowing for rapid recovery or migration of tenant settings, policies, and assignments in the event of accidental changes or tenant loss.

- **Collaboration and Delegation**
  Empower teams to collaborate on Microsoft 365 configuration using pull requests, code reviews, and CI/CD pipelines, reducing bottlenecks and enabling safe delegation of administrative tasks.

- **Bulk and Consistent Policy Enforcement**
  Apply security, compliance, and device management policies at scale, ensuring consistency and reducing manual configuration errors across large organizations or multiple tenants.

- **Self-Service via Terraform Modules**  
  Build reusable Terraform modules for common Microsoft 365 workloads, enabling service-owning teams to provide self-service provisioning to other engineering teams while maintaining standards and reducing manual effort.

- **Integration with Policy-as-Code (OPA/Conftest)**  
  Integrate with Open Policy Agent (OPA) or Conftest to enforce organizational standards, compliance, and guardrails on Microsoft 365 resources before deployment, ensuring only approved configurations are applied in production.

- **Guardrailed Deployments**  
  Implement automated checks and guardrails in CI/CD pipelines to prevent misconfiguration and enforce best practices, reducing risk and improving governance for Microsoft 365 administration.

## Getting Started

Please refer to the [Getting Started](https://registry.terraform.io/providers/deploymenttheory/microsoft365/latest/docs) guide in the terraform registry for more information on how to get started.

## Provider Key Features

- **Wide Resource Coverage**: Supports management of resources across Microsoft Intune, Microsoft 365, Microsoft Teams, Microsoft Defender, and related services. This includes device and app management, user and group management, and administrative resources.
- **Multi-Cloud Compatibility**: Operates with Microsoft public cloud, US Government (GCC, GCC High, DoD), China, and other national cloud environments.
- **Multiple Authentication Methods**: Provides support for client secret, client certificate, device code, interactive browser, managed identity, workload identity, OIDC (including GitHub Actions and Azure DevOps), and Azure Developer CLI authentication.
- **Proxy and Network Configuration**: Allows configuration of HTTP proxy settings, custom user agents, request timeouts, and retry/redirect/compression policies.
- **Microsoft Graph API Support**: Integrates with both v1.0 and beta Microsoft Graph API endpoint sets, to support both generally available and preview features.
- **Microsoft Graph SDK Adoption**: The provider leverages the Microsoft Graph API through the Kiota-generated graphSDKs, allowing for a strongly typed development experience.

## Provider Comparison

For information and a comparison between this provider in relation to the msft official `terraform-provider-msgraph provider`, see the [Provider Comparison](./docs/development/provider_comparison.md) documentation.

## Community Contributions

As a community-driven project, contributions, feedback, and issue reports are welcome and encouraged. Together, we can enhance and expand the capabilities of this provider to meet the evolving needs of Microsoft 365 administrators and DevOps professionals.

### Development Guide

The style guidelines and the design decisions for this provider can be found here

[Development Guide](./docs/development/guide.md)

### Known Issues

For information about known bugs and limitations, including Microsoft Graph API issues and workarounds, see:

[Known Issues & Bugs](./docs/development/known_bugs.md)

## Community Terraform Provider for Microsoft 365 Provider Roadmap

Please see the roadmap below on the intended provider resource suppport

[Provider Roadmap](./docs/providerroadmap/roadmap.md)

## Disclaimer

> [!IMPORTANT]  
> While every effort is made to maintain accuracy and reliability, users should thoroughly test configurations in non-production environments before deploying to production. Always refer to official Microsoft documentation for the most up-to-date information on Microsoft 365 services and features.

## Data Collection

The software may collect information about you and your use of the software and send it to Microsoft. Microsoft may use this information to provide services and improve their products and services. You may turn off the telemetry as described in the repository. There are also some features in the software that may enable you and Microsoft to collect data from users of your applications. If you use these features, you must comply with applicable law, including providing appropriate notices to users of your applications together with a copy of Microsoft's privacy statement. Microsoft's privacy statement is located at https://go.microsoft.com/fwlink/?LinkID=824704. You can learn more about data collection and use in the help documentation and their privacy statement. Your use of the software operates as your consent to these practices.

