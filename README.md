# Community Terraform Provider for Microsoft 365

The community Microsoft 365 Terraform Provider allows managing environments and other resources within [Microsoft Intune](https://www.microsoft.com/en-us/security/business/endpoint-management/microsoft-intune), [Microsoft 365](https://www.microsoft.com/en-us/microsoft-365), [Microsoft Teams](https://www.microsoft.com/en-us/microsoft-teams/group-chat-software), and [Microsoft Defender](https://www.microsoft.com/en-us/security/business/microsoft-defender).

> [!WARNING]
> This code is experimental and provided solely for evaluation purposes. It is **NOT** intended for production use and may contain bugs, incomplete features, or other issues. Use at your own risk, as it may undergo significant changes without notice until it reaches general availability, and no guarantees or support is provided. By using this code, you acknowledge and agree to these conditions. Consult the documentation or contact the maintainer if you have questions or concerns.

## Overview

The Community Terraform Provider for Microsoft 365 empowers DevOps teams and administrators to manage their Microsoft 365 environments using Infrastructure as Code (IaC) principles. This provider bridges the gap between Terraform's powerful resource management capabilities and the extensive features of Microsoft 365, allowing for automated, version-controlled, and repeatable deployments across various Microsoft cloud services.

## Key Features

- **Comprehensive Resource Management**: Manage resources across Microsoft Intune, Microsoft 365, Microsoft Teams, and Microsoft Defender.
- **Multi-Cloud Support**: Compatible with various Microsoft cloud environments, including public, government, and national clouds.
- **Flexible Authentication**: Supports multiple authentication methods, including client credentials, certificate-based, and interactive browser flows.
- **Beta API Access**: Includes a beta client for accessing cutting-edge features and APIs still in development.
- **Enhanced Security Options**: Offers proxy support and various security configurations to align with organizational policies.

## Use Cases

- Automate the creation and management of user accounts, groups, and permissions.
- Deploy and configure Microsoft Teams environments at scale.
- Manage security policies and compliance settings across your Microsoft 365 tenant.
- Provision and configure Intune policies for device management.

## Community Contributions

As a community-driven project, contributions, feedback, and issue reports are welcome and encouraged. Together, we can enhance and expand the capabilities of this provider to meet the evolving needs of Microsoft 365 administrators and DevOps professionals.

## Community Terraform Provider for Microsoft 365 Provider Roadmap

Please see the roadmap below on the intended provider resource suppport

[Provider Roadmap](./docs/providerroadmap/roadmap.md)

## Development Guide

[Development Guide](./docs/developmentguide/guide.md)

## Disclaimer

While every effort is made to maintain accuracy and reliability, users should thoroughly test configurations in non-production environments before deploying to production. Always refer to official Microsoft documentation for the most up-to-date information on Microsoft 365 services and features.

## Data Collection

The software may collect information about you and your use of the software and send it to Microsoft. Microsoft may use this information to provide services and improve our products and services. You may turn off the telemetry as described in the repository. There are also some features in the software that may enable you and Microsoft to collect data from users of your applications. If you use these features, you must comply with applicable law, including providing appropriate notices to users of your applications together with a copy of Microsoft’s privacy statement. Our privacy statement is located at https://go.microsoft.com/fwlink/?LinkID=824704. You can learn more about data collection and use in the help documentation and our privacy statement. Your use of the software operates as your consent to these practices.

