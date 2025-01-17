# Microsoft 365 Terraform Provider Roadmap

This roadmap outlines the development trajectory for the Microsoft 365 Terraform Provider, which aims to enable infrastructure-as-code management of Microsoft 365 services and configurations through HashiCorp Terraform.

## Core Vision

Our goal is to provide a robust, enterprise-ready Terraform provider that allows organizations to manage their Microsoft 365 environment programmatically, bringing the benefits of infrastructure as code to Microsoft 365 administration. This includes user management, security configurations, compliance policies, device management and service-specific settings across the Microsoft 365 suite.

## Current Focus

The initial development phase focuses on implementing core Microsoft 365 management capabilities through Microsoft Graph API integration, with particular emphasis on:

Identity and Access Management (Azure AD)
Security and Compliance Policies
Intune Device Management
Microsoft 365 Apps Configuration
Teams Management

## Development Principles

- API-First: Built on Microsoft Graph API and Graph Beta for future-proof integration
- Enterprise-Ready: Focusing on features most requested by enterprise environments
- Security-Focused: Implementing secure defaults and best practices
- Idempotent Operations: Ensuring consistent and predictable resource management
- Comprehensive Testing: Maintaining high test coverage and validation

## Completion Status

## Microsoft Exchange Online

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
- [ ] Outlook on the web policies
- [ ] Partner Applications
- [ ] PolicyTipConfig
- [ ] Role Assignment Policies
- [ ] Transport Rules

## Microsoft Purview

- [ ] Audit configuration policy
- [ ] Case hold policies
- [ ] Case hold rules
- [ ] Compliance cases
- [ ] ComplianceTags
- [ ] Content search actions
- [ ] Content searches
- [ ] DLP compliance policies
- [ ] DLP sensitive information types
- [ ] File plan property authorities
- [ ] File plan property categories
- [ ] File plan property citations
- [ ] File plan property departments
- [ ] File plan property reference ids
- [ ] File plan property reference sub categories
- [ ] Information governance
- [ ] Compliance Retention Event Types
- [ ] Retention
- [ ] Label Policy
- [ ] Labels

## Microsoft Defender for Office 365

- [ ] Threat management
- [ ] Policy
- [ ] Anti-Phishing
- [ ] Safe Attachments
- [ ] Safe Links
- [ ] Global Settings
- [ ] Quarantine Policies
- [ ] Hosted connection filter policies
- [ ] Hosted content filter policies
- [ ] Hosted content filter rules
- [ ] Hosted outbound spam filter policies
- [ ] Hosted outbound spam filter rules

## Microsoft Entra ID (formerly Azure Active Directory)

- [ ] Device conditional access policies
- [ ] Device configuration policies

## Microsoft Teams

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

## Microsoft Intune

### Apps

- [ ] Applications
- [ ] App configuration policies
- [ ] App protection policies (Platform = Android)
- [ ] App protection policies (Platform = iOS/iPadOS)
- [ ] App protection policies (Platform = Windows 10/11)
- [ ] Diagnostic settings
- [ ] Endpoint security
- [ ] Mobile Threat Defense
- [ ] Policy Sets

### Device Management

- [ ] Compliance policies
- [ ] Locations
- [ ] Notifications
- [x] Configuration profiles (Settings Catalog = Administrative Templates)
- [x] Configuration profiles (Settings Catalog)
- [ ] Configuration profiles (Other legacy types)
- [ ] Enrollment restrictions
- [x] Platform Scripts (Windows, macOS, Linux)
- [ ] Proactive Remediations
- [ ] Windows Autopilot deployment profiles
- [ ] Quality updates for Windows 10 and later
- [ ] Feature updates for Windows 10 and later

### Reports

- [ ] Endpoint analytics

### Intune Suite

- [ ] Advanced Analytics
- [ ] Advanced Endpoint Analytics
- [ ] Remote Help
- [x] Endpoint Privilege Management (Reuseable policies, elevation policies, elevation rules)
- [ ] Microsoft Tunnel for Mobile Application Management
- [ ] Automated App Patching

## Windows 365

- [ ] Cloud PC Provisioning Policies
- [ ] User Settings Policies
- [ ] Network Settings
- [ ] Image Management
- [ ] Windows 365 Frontline
- [ ] Windows 365 Boot
- [ ] Windows 365 Switch
- [ ] Cloud PC Health Monitoring
- [ ] Cloud PC Restore
- [ ] Cloud PC Resize
- [ ] Cloud PC Reprovision
- [ ] Usage Analytics
- [ ] Point-in-Time Restore
- [ ] Universal Print Integration
- [ ] Regional Settings
- [ ] Custom Language Packs

### Intune Tenant administration

- [x] Assignment Filters
- [x] Role Scope Tags
- [x] Role Definitions
- [ ] Device clean-up rules