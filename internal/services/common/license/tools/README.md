# Microsoft 365 License Export Tools

This directory contains tools to export and analyze Microsoft 365 licenses from your tenant.

## Tools

### 1. PowerShell Script (`generate_constants.ps1`)

Queries Microsoft Graph API to retrieve all subscribed SKUs and service plans, then exports them to multiple formats.

#### Prerequisites

- PowerShell 7.0 or later
- Microsoft.Graph.Authentication module
- Service Principal with `Organization.Read.All` permission

#### Installation

```powershell
# Install required module
Install-Module Microsoft.Graph.Authentication -Scope CurrentUser
```

#### Usage

```powershell
.\generate_constants.ps1 `
    -TenantId "your-tenant-id" `
    -ClientId "your-client-id" `
    -ClientSecret "your-client-secret"
```

With custom output path:

```powershell
.\generate_constants.ps1 `
    -TenantId "your-tenant-id" `
    -ClientId "your-client-id" `
    -ClientSecret "your-client-secret" `
    -OutputPath "./exports"
```

#### Outputs

The script generates 2 files:

1. **`licenses_YYYYMMDD_HHMMSS.json`** - Complete SKU data with all details
2. **`generated_constants.go`** - Go constants ready to use in the provider, grouped by SKU

---

## Output Formats Explained

### JSON Export

Contains complete raw data from Microsoft Graph API:

```json
{
  "exportDate": "2025-11-18 14:30:00",
  "tenantId": "...",
  "totalSkus": 11,
  "skus": [
    {
      "skuId": "...",
      "skuPartNumber": "Microsoft_Entra_Suite",
      "capabilityStatus": "Enabled",
      "consumedUnits": 10,
      "prepaidUnitsEnabled": 25,
      "servicePlans": [...]
    }
  ]
}
```

### Go Constants

Service plans are intelligently organized:
- **Shared plans** appear in a single section with comments listing all SKUs that include them
- **SKU-specific plans** are grouped under their parent SKU

This prevents duplicate constant definitions while maintaining logical grouping.

```go
// Auto-generated from Microsoft Graph API
// Generated: 2025-11-18 14:30:00

package constants

// ============================================================================
// SKU Part Numbers (11 total)
// ============================================================================
const (
    SKUMicrosoftEntraSuite = "Microsoft_Entra_Suite"
    SKUMicrosoftIntuneSuite = "Microsoft_Intune_Suite"
    // ...
)

// ============================================================================
// Shared Service Plans (15 plans)
// These service plans appear in multiple SKUs
// ============================================================================
const (
    ServicePlanEXCHANGESFOUNDATION = "EXCHANGE_S_FOUNDATION" // Shared: CPC_B_2C_8RAM_128GB, CPC_E_2C_8GB_128GB, FLOW_FREE, POWER_BI_STANDARD, RMSBASIC, WINDOWS_STORE
    ServicePlanM365LIGHTHOUSECUSTOMERPLAN1 = "M365_LIGHTHOUSE_CUSTOMER_PLAN1" // Shared: CPC_B_2C_8RAM_128GB, SPE_E5
    // ...
)

// ============================================================================
// Service Plans from: Microsoft_Entra_Suite (5 plans)
// ============================================================================
const (
    ServicePlanEntraIDentityGovernance = "Entra_Identity_Governance"
    ServicePlanEntraPremiumInternetAccess = "Entra_Premium_Internet_Access"
    // ...
)

// ============================================================================
// Service Plans from: Microsoft_Intune_Suite (8 plans)
// ============================================================================
const (
    ServicePlan3PARTYAPPPATCH = "3_PARTY_APP_PATCH"
    ServicePlanCLOUDPKI = "CLOUD_PKI"
    // ...
)
```

---

## Architecture

The license constants are organized into two locations:

### 1. **Centralized Constants** (`internal/constants/sku_service_plan.go`)
- Contains **all** SKU and Service Plan constants generated from the Microsoft Graph API
- Single source of truth for license identifiers
- Used throughout the provider for license checks and validations

### 2. **License Feature Mappings** (`internal/services/common/license/required.go`)
- Imports centralized constants from `internal/constants`
- Defines `RequiredLicenses` map: which licenses enable which features
- Contains business logic for license requirements
- References [m365maps.com](https://m365maps.com/) for accurate license-to-feature mappings

**Example usage:**
```go
import (
    "github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
    "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/license"
)

// Check if a feature is available
requiredLicenses := license.GetRequiredLicensesForFeature("NetworkFilteringPolicy")
// Returns: [constants.SKUMicrosoftEntraSuite, constants.ServicePlanEntraPremiumInternetAccess]
```

---

## Integration Workflow

1. **Run the export tool** to discover all licenses in your tenant
2. **Review the JSON file** to understand which licenses and service plans are available
3. **Review the generated Go constants** to see all SKUs and their service plans organized by parent SKU
4. **Copy the generated constants** from `exports/generated_constants.go` to `../../../constants/sku_service_plan.go`
5. **Update `RequiredLicenses` map** in `../required.go` to reference the centralized constants with feature-to-license mappings
   - Use [m365maps.com](https://m365maps.com/) to verify which features are included in each license tier
   - Cross-reference Microsoft Learn documentation for specific feature requirements
6. **Re-run periodically** as your tenant licenses change or Microsoft updates their licensing structure

---

## Comprehensive License Mappings

The `required.go` file contains an extensive mapping of Microsoft 365 features to their required licenses, organized by product category:

### Categories Covered

1. **Microsoft Entra** (Identity & Access Management)
   - Conditional Access, PIM, Identity Governance, Verified ID, Global Secure Access, MFA

2. **Microsoft Intune** (Endpoint Management)
   - Basic device management, Endpoint Privilege Management, Cloud PKI, Advanced Analytics, Remote Help, MAM Tunnel, ServiceNow integration, Third-party app patching, Windows Autopatch

3. **Windows 365** (Cloud PC)
   - Business, Enterprise, and Frontline editions, Windows 10 ESU

4. **Microsoft Defender** (Threat Protection)
   - Defender for Endpoint, Office 365, Cloud Apps, Identity, XDR, Threat Intelligence, IoT

5. **Microsoft Purview** (Compliance & Data Governance)
   - Data Lifecycle Management, Information Protection, DLP, Insider Risk, Communication Compliance, eDiscovery, Advanced Audit, Content Explorer, Information Barriers, Customer Key, Customer Lockbox, Premium Encryption, ML Classification, Safe Documents

6. **Microsoft 365 Apps & Productivity**
   - Office apps, Exchange, SharePoint, OneDrive, Teams (including Phone System, Advanced Meetings, Mesh), Loop, Search, Bookings, Forms, Sway, Clipchamp, Whiteboard, To Do, Stream, Yammer/Viva Engage, Project, Universal Print, Excel Premium

7. **Power Platform**
   - Power Apps, Power Automate, Power Virtual Agents, Dataverse

8. **Viva & Analytics**
   - Viva Insights, Learning, Nucleus, Engage, People Skills, Places

9. **Power BI**
   - Standard and Premium tiers

10. **Additional Services**
    - Bing Chat Enterprise, Microsoft 365 Lighthouse (MSP tooling)

### Reference Sources

All license mappings are validated against:
- **[m365maps.com](https://m365maps.com/)** - Comprehensive visual licensing diagrams by Aaron Dinnage
- **[Microsoft Learn](https://learn.microsoft.com/en-us/entra/fundamentals/licensing)** - Official licensing documentation

### Usage in Provider

Resources can check for required licenses before making API calls:

```go
if !license.HasRequiredLicense(ctx, r.client, "ConditionalAccessPolicy") {
    resp.Diagnostics.AddError(
        "Missing Required License",
        fmt.Sprintf(
            "This resource requires a tenant license that was not found.\n\n%s",
            license.FormatRequiredLicensesMessage("ConditionalAccessPolicy"),
        ),
    )
    return
}
```

This provides users with clear, actionable error messages instead of cryptic 403 Forbidden errors.

---

## Important Notes

### Limitations

- Only discovers licenses **present in your tenant**
- Does not show licenses you don't subscribe to
- Microsoft doesn't provide a "catalog" API of all possible licenses

For a complete license catalog, you would need:
- Access to multiple tenants with different licenses
- Microsoft's internal documentation
- Community-sourced license names

---

## Troubleshooting

### "Organization.Read.All permission is required"

Grant the required permission to your service principal:

```powershell
# Using Microsoft Graph PowerShell
Connect-MgGraph -Scopes "Application.ReadWrite.All"

$sp = Get-MgServicePrincipal -Filter "appId eq 'your-client-id'"
$graphApp = Get-MgServicePrincipal -Filter "appId eq '00000003-0000-0000-c000-000000000000'"
$orgReadAll = $graphApp.AppRoles | Where-Object { $_.Value -eq "Organization.Read.All" }

New-MgServicePrincipalAppRoleAssignment `
    -ServicePrincipalId $sp.Id `
    -PrincipalId $sp.Id `
    -ResourceId $graphApp.Id `
    -AppRoleId $orgReadAll.Id
```

### "No SKUs found"

Verify:
1. Service principal has correct permissions
2. Tenant has active subscriptions
3. Credentials are correct

### "Authentication failed"

Check:
1. Client Secret hasn't expired
2. Tenant ID is correct
3. Client ID is correct
4. Service Principal is enabled

---

## Example Session

```powershell
PS> .\generate_constants.ps1 -TenantId "..." -ClientId "..." -ClientSecret "..."

╔════════════════════════════════════════════════════════════════╗
║         Microsoft 365 License Export Tool                     ║
╚════════════════════════════════════════════════════════════════╝

🔐 Authenticating to Microsoft Graph...
✅ Authentication successful

📊 Fetching subscribed SKUs from tenant...
✅ Retrieved 11 SKUs

📦 Export Summary:
   • Total SKUs: 11
   • Total Service Plans: 116

📄 Exporting data...

✅ Exported SKUs to: licenses_20251118_143000.json
✅ Generated Go constants to: generated_constants.go
   • 11 SKUs
   • 116 Service Plans (15 shared, 8 SKU-specific)

╔════════════════════════════════════════════════════════════════╗
║                    Export Complete!                            ║
╚════════════════════════════════════════════════════════════════╝

📁 Exported files:
   • JSON (full data):   licenses_20251118_143000.json
   • Go constants:       generated_constants.go

🔓 Disconnected from Microsoft Graph
```

---

## References

- [Microsoft Graph API - subscribedSkus](https://learn.microsoft.com/en-us/graph/api/subscribedsku-list)
- [Microsoft 365 licensing overview](https://learn.microsoft.com/en-us/entra/fundamentals/licensing)
- [Product names and service plan identifiers](https://learn.microsoft.com/en-us/entra/identity/users/licensing-service-plan-reference)

