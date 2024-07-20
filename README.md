# Community Terraform Provider for Microsoft 365

The Microsoft 365 Terraform Provider allows managing environments and other resources within [Intune](https://intune.microsoft.com/) , [Office365](https://www.office.com/
), [MicrosoftTeams](https://teams.microsoft.com/
) and [MicrosoftSecurity](https://security.microsoft.com/
).

> [!WARNING]
> This code is experimental and provided solely for evaluation purposes. It is **NOT** intended for production use and may contain bugs, incomplete features, or other issues. Use at your own risk, as it may undergo significant changes without notice until it reaches general availability, and no guarantees or support are provided. By using this code, you acknowledge and agree to these conditions. Consult the documentation or contact the maintainer if you have questions or concerns.

# To-Do List

## Exchange

- [ ] AddressBookPolicy
- [ ] AddressList
- [ ] AdminAuditLogConfig
- [ ] Application access policies
- [ ] Availability address spaces
- [ ] Availability config
- [ ] CAS mailbox plan
- [ ] Client access rules
- [ ] DomainKeys identified mail signing config
- [ ] Dynamic Distribution Groups
- [ ] Email address policies
- [ ] GlobalAddressList
- [ ] Inbound connectors
- [ ] Inbound IntraOrganizationConnectors
- [ ] Journal Rules
- [ ] Mail flow
- [ ] Accepted Domains
- [ ] Remote Domains
- [ ] Mailboxes
- [ ] MailboxPlans
- [ ] MalwareFilterPolicies
- [ ] ManagementRole
- [ ] Mobile device access
- [ ] Device Access Rules
- [ ] MobileDeviceMailboxPolicies
- [ ] Modern authentication
- [ ] OfflineAddressBook
- [ ] On-Premises Organizations
- [ ] Organization Relationship
- [ ] OrganizationConfig
- [ ] Outbound connectors
- [ ] Outlook Web App policies
- [ ] Partner Applications
- [ ] PolicyTipConfig
- [ ] Role Assignment Policies
- [ ] Transport Rules

## Security & Compliance

- [ ] Audit configuration policy
- [ ] Case hold policies
- [ ] Case hold rules
- [ ] Compliance cases
- [ ] ComplianceTags
- [ ] Content search actions
- [ ] Content searches
- [ ] Device conditional access policies
- [ ] Device configuration policies
- [ ] DLP compliance policies
- [ ] DLP senstitive information types
- [ ] File plan property authorities
- [ ] File plan property categories
- [ ] File plan property citations
- [ ] File plan property departments
- [ ] File plan property reference ids
- [ ] File plan property reference sub categories
- [ ] Hosted connection filter policies
- [ ] Hosted content filter policies
- [ ] Hosted content filter rules
- [ ] Hosted outbound spam filter policies
- [ ] Hosted outbound spam filter rules
- [ ] Information governance
- [ ] Compliance Retention Event Types
- [ ] Retention
- [ ] Label Policy
- [ ] Labels
- [ ] Threat management
- [ ] Policy
- [ ] ATP Anti-Phishing
- [ ] ATP Safe Attachments
- [ ] ATP Safe Links
- [ ] Global Settings
- [ ] Quarantine Policies

## Teams
- [ ] Apps
- [ ] Permission policies
- [ ] Meetings
- [ ] Meeting policies
- [ ] Meeting settings
- [ ] Messaging policies
- [ ] Org-wide settings
- [ ] Teams settings
- [ ] Voice
- [ ] Calling policies

## Intune

### Apps

- [ ] App configuration policies
- [ ] App protection policies (Platform = Android)
- [ ] App protection policies (Platform = iOS/iPadOS)
- [ ] App protection policies (Platform = Windows 10)
- [ ] Diagnostic settings
- [ ] Endpoint security
- [ ] Mobile Threat Defense
- [ ] Policy Sets

### Devices

- [ ] Compliance policies
- [ ] Compliance policy settings
- [ ] Locations
- [ ] Notifications
- [ ] Configuration profiles
- [ ] Configuration profiles (Profile Type = Administrative Templates)
- [ ] Configuration profiles (Settings Catalog)
- [ ] Device clean-up rules
- [ ] Enrollment restrictions
- [ ] Scripts
- [ ] Windows Autopilot deployment profiles
- [ ] Quality updates for Windows 10 and later
- [ ] Feature updates for Windows 10 and later

### Reports

- [ ] Endpoint analytics
- [ ] Proactive Remediations

### Tenant administration

- [ ] Filters

# Template

This repository serves as a **Default Template Repository** according official [GitHub Contributing Guidelines][ProjectSetup] for healthy contributions. It brings you clean default Templates for several areas:

- [Azure DevOps Pull Requests](.azuredevops/PULL_REQUEST_TEMPLATE.md) ([`.azuredevops\PULL_REQUEST_TEMPLATE.md`](`.azuredevops\PULL_REQUEST_TEMPLATE.md`))
- [Azure Pipelines](.pipelines/pipeline.yml) ([`.pipelines/pipeline.yml`](`.pipelines/pipeline.yml`))
- [GitHub Workflows](.github/workflows/)
  - [Super Linter](.github/workflows/linter.yml) ([`.github/workflows/linter.yml`](`.github/workflows/linter.yml`))
  - [Sample Workflows](.github/workflows/workflow.yml) ([`.github/workflows/workflow.yml`](`.github/workflows/workflow.yml`))
- [GitHub Pull Requests](.github/PULL_REQUEST_TEMPLATE.md) ([`.github/PULL_REQUEST_TEMPLATE.md`](`.github/PULL_REQUEST_TEMPLATE.md`))
- [GitHub Issues](.github/ISSUE_TEMPLATE/)
  - [Feature Requests](.github/ISSUE_TEMPLATE/FEATURE_REQUEST.md) ([`.github/ISSUE_TEMPLATE/FEATURE_REQUEST.md`](`.github/ISSUE_TEMPLATE/FEATURE_REQUEST.md`))
  - [Bug Reports](.github/ISSUE_TEMPLATE/BUG_REPORT.md) ([`.github/ISSUE_TEMPLATE/BUG_REPORT.md`](`.github/ISSUE_TEMPLATE/BUG_REPORT.md`))
- [Codeowners](.github/CODEOWNERS) ([`.github/CODEOWNERS`](`.github/CODEOWNERS`)) _adjust usernames once cloned_
- [Wiki and Documentation](docs/) ([`docs/`](`docs/`))
- [gitignore](.gitignore) ([`.gitignore`](.gitignore))
- [gitattributes](.gitattributes) ([`.gitattributes`](.gitattributes))
- [Changelog](CHANGELOG.md) ([`CHANGELOG.md`](`CHANGELOG.md`))
- [Code of Conduct](CODE_OF_CONDUCT.md) ([`CODE_OF_CONDUCT.md`](`CODE_OF_CONDUCT.md`))
- [Contribution](CONTRIBUTING.md) ([`CONTRIBUTING.md`](`CONTRIBUTING.md`))
- [License](LICENSE) ([`LICENSE`](`LICENSE`)) _adjust projectname once cloned_
- [Readme](README.md) ([`README.md`](`README.md`))
- [Security](SECURITY.md) ([`SECURITY.md`](`SECURITY.md`))


## Status

[![Super Linter](<https://github.com/segraef/Template/actions/workflows/linter.yml/badge.svg>)](<https://github.com/segraef/Template/actions/workflows/linter.yml>)

[![Sample Workflow](<https://github.com/segraef/Template/actions/workflows/workflow.yml/badge.svg>)](<https://github.com/segraef/Template/actions/workflows/workflow.yml>)

## Creating a repository from a template

You can [generate](https://github.com/segraef/Template/generate) a new repository with the same directory structure and files as an existing repository. More details can be found [here][CreateFromTemplate].

## Reporting Issues and Feedback

### Issues and Bugs

If you find any bugs, please file an issue in the [GitHub Issues][GitHubIssues] page. Please fill out the provided template with the appropriate information.

If you are taking the time to mention a problem, even a seemingly minor one, it is greatly appreciated, and a totally valid contribution to this project. **Thank you!**

## Feedback

If there is a feature you would like to see in here, please file an issue or feature request in the [GitHub Issues][GitHubIssues] page to provide direct feedback.

## Contribution

If you would like to become an active contributor to this repository or project, please follow the instructions provided in [`CONTRIBUTING.md`][Contributing].

## Learn More

* [GitHub Documentation][GitHubDocs]
* [Azure DevOps Documentation][AzureDevOpsDocs]
* [Microsoft Azure Documentation][MicrosoftAzureDocs]

<!-- References -->

<!-- Local -->
[ProjectSetup]: <https://docs.github.com/en/communities/setting-up-your-project-for-healthy-contributions>
[CreateFromTemplate]: <https://docs.github.com/en/github/creating-cloning-and-archiving-repositories/creating-a-repository-on-github/creating-a-repository-from-a-template>
[GitHubDocs]: <https://docs.github.com/>
[AzureDevOpsDocs]: <https://docs.microsoft.com/en-us/azure/devops/?view=azure-devops>
[GitHubIssues]: <https://github.com/segraef/Template/issues>
[Contributing]: CONTRIBUTING.md

<!-- External -->
[Az]: <https://img.shields.io/powershellgallery/v/Az.svg?style=flat-square&label=Az>
[AzGallery]: <https://www.powershellgallery.com/packages/Az/>
[PowerShellCore]: <https://github.com/PowerShell/PowerShell/releases/latest>

<!-- Docs -->
[MicrosoftAzureDocs]: <https://docs.microsoft.com/en-us/azure/>
[PowerShellDocs]: <https://docs.microsoft.com/en-us/powershell/>
