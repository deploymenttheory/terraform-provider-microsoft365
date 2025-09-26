# Known Issues & Bugs

This document tracks known bugs and limitations in the provider, particularly those related to Microsoft Graph API issues and other external dependencies that are outside of the provider's direct control.

## How to Report Issues

If you encounter a new issue:

1. **Search existing issues** in this document and [GitHub Issues](https://github.com/deploymenttheory/terraform-provider-microsoft365/issues)
2. **Verify reproducibility** in your environment
3. **Report new issues** using the [Issue Template](#issue-template) below
4. **Include all required information** for proper tracking and resolution

## Issue Tracking Template

Use this template when adding new known issues to this document:

```markdown
### [BUG-XXX] Brief Issue Title

**Bug ID:** BUG-XXX  
**Resource:** `resource_name_here`  
**Provider Version:** vX.X.X (first affected version)  
**Date Reported:** YYYY-MM-DD  
**Last Verified:** YYYY-MM-DD  
**Impact Level:** Critical/High/Medium/Low  

#### Affected API Endpoints
- `/path/to/graph/endpoint`
- `/another/affected/endpoint`

#### Expected Behavior
Description of what should happen according to documentation or expected functionality.

#### Observed Behavior
Description of what actually happens, including error messages, unexpected responses, etc.

#### Reproduction Steps
1. Step one
2. Step two
3. Step three

#### Environment Details
- **Affected Clouds:** Public/GCC/GCC High/DoD/China/All
- **Terraform Version:** vX.X.X
- **Provider Version:** vX.X.X
- **Operating System:** Windows/macOS/Linux

#### Status
- **Current Status:** Open/Investigating/Microsoft Acknowledged/Fixed in Graph/Workaround Available/Closed
- **Microsoft Status:** Not Reported/Reported/Acknowledged/In Progress/Fixed/Won't Fix
- **Upstream Links:**
  - Microsoft Graph Known Issues: [link if applicable]
  - GitHub Issues: [link to related issues]
  - Microsoft Feedback: [link if reported to Microsoft]

#### Workarounds
Description of any available workarounds, including:
- Alternative approaches
- Manual steps
- Configuration changes
- Version constraints

#### Additional Notes
Any other relevant information, context, or dependencies.

---
```

## Current Known Issues

> **Note:** This section will be populated as issues are discovered and verified. Each issue should follow the template format above.

### Example Issue (Template Demonstration)

### [BUG-001] Windows Update Ring Assignments Not Reflecting Filter Changes

**Bug ID:** BUG-001  
**Resource:** `microsoft365_graph_beta_device_management_windows_update_ring`  
**Provider Version:** v0.1.0+  
**Date Reported:** 2025-01-15  
**Last Verified:** 2025-01-15  
**Impact Level:** Medium  

#### Affected API Endpoints
- `/deviceManagement/deviceConfigurations/{id}/assign`
- `/deviceManagement/deviceConfigurations/{id}?$expand=assignments`

#### Expected Behavior
When updating assignment filter configurations for Windows Update Rings, the changes should be immediately reflected in subsequent read operations and take effect for device targeting.

#### Observed Behavior
Assignment filter changes are accepted by the API but may not be immediately visible in read operations due to Microsoft Graph eventual consistency. Devices may continue to receive the old assignment targeting for up to 24 hours.

#### Reproduction Steps
1. Create a Windows Update Ring with group assignments
2. Update the assignment to include filter criteria  
3. Immediately read the resource state
4. Observer that filter changes may not be reflected

#### Environment Details
- **Affected Clouds:** All
- **Terraform Version:** v1.5.0+
- **Provider Version:** v0.1.0+
- **Operating System:** All

#### Status
- **Current Status:** Microsoft Acknowledged
- **Microsoft Status:** Known Limitation
- **Upstream Links:**
  - Microsoft Graph Known Issues: [Microsoft Graph known issues - Microsoft Graph | Microsoft Learn](https://docs.microsoft.com/en-us/graph/known-issues)

#### Workarounds
1. **Wait Period:** Allow 15-30 minutes after assignment changes before expecting full consistency
2. **Retry Logic:** Implement retry logic in deployment pipelines when checking assignment state
3. **Manual Verification:** Verify assignment changes in the Microsoft Endpoint Manager console

#### Additional Notes
This is a known limitation of Microsoft Graph's eventual consistency model, not a provider bug. The provider correctly sends the API requests, but Microsoft's backend systems require time to propagate changes across all services.

---

## Issue Status Definitions

| Status | Description |
|--------|-------------|
| **Open** | Issue identified but not yet investigated |
| **Investigating** | Provider team is researching the issue |
| **Microsoft Acknowledged** | Microsoft has acknowledged the issue |
| **Fixed in Graph** | Microsoft has fixed the underlying API issue |
| **Workaround Available** | Temporary solution is available |
| **Closed** | Issue has been resolved |

## Impact Level Definitions

| Level | Description |
|-------|-------------|
| **Critical** | Complete functionality failure, no workaround available |
| **High** | Major functionality impaired, limited workarounds |
| **Medium** | Functionality affected but workarounds available |
| **Low** | Minor inconvenience, easy workarounds available |

## Contributing to Issue Tracking

When contributing to this document:

1. **Use the template** provided above for consistency
2. **Assign sequential bug IDs** (BUG-001, BUG-002, etc.)
3. **Include all required fields** - don't leave sections empty
4. **Verify issues** before documenting to avoid duplicates
5. **Update status** regularly as issues evolve
6. **Link to related GitHub issues** when available
7. **Keep workarounds current** and test their effectiveness

## Resources

- [Microsoft Graph Known Issues](https://docs.microsoft.com/en-us/graph/known-issues)
- [Microsoft Graph Feedback Portal](https://aka.ms/graphfeedback)
- [Provider GitHub Issues](https://github.com/deploymenttheory/terraform-provider-microsoft365/issues)
- [Microsoft 365 Service Health](https://admin.microsoft.com/servicehealth)