<#
.SYNOPSIS
    Exports all Microsoft 365 SKUs and Service Plans from tenant

.DESCRIPTION
    Queries Microsoft Graph API to retrieve all subscribed SKUs and their service plans,
    then exports them to JSON and CSV formats for use in generating Go constants.

.PARAMETER TenantId
    Entra ID Tenant ID (Directory ID)

.PARAMETER ClientId
    Application (Client) ID of the Entra ID app registration

.PARAMETER ClientSecret
    Client Secret of the Entra ID app registration

.PARAMETER OutputPath
    Path where the export files will be saved (default: current directory)

.EXAMPLE
    .\generate_constants.ps1 -TenantId "your-tenant-id" -ClientId "your-client-id" -ClientSecret "your-client-secret"

.EXAMPLE
    .\generate_constants.ps1 -TenantId "your-tenant-id" -ClientId "your-client-id" -ClientSecret "your-client-secret" -OutputPath "./exports"

.NOTES
    Requires: Microsoft.Graph.Authentication module
    Required Permission: Organization.Read.All
#>

[CmdletBinding()]
param (
    [Parameter(Mandatory=$true,
    HelpMessage="Specify the Entra ID tenant ID (Directory ID)")]
    [ValidateNotNullOrEmpty()]
    [string]$TenantId,

    [Parameter(Mandatory=$true,
    HelpMessage="Specify the application (client) ID of the Entra ID app registration")]
    [ValidateNotNullOrEmpty()]
    [string]$ClientId,
    
    [Parameter(Mandatory=$true,
    HelpMessage="Specify the client secret of the Entra ID app registration")]
    [ValidateNotNullOrEmpty()]
    [string]$ClientSecret,

    [Parameter(Mandatory=$false,
    HelpMessage="Path where the export files will be saved")]
    [string]$OutputPath = "."
)

#region Helper Functions

function Get-AllSubscribedSkus {
    <#
    .SYNOPSIS
        Retrieves all subscribed SKUs from the tenant
    #>
    try {
        Write-Host "📊 Fetching subscribed SKUs from tenant..." -ForegroundColor Cyan
        
        $uri = "https://graph.microsoft.com/beta/subscribedSkus"
        $response = Invoke-MgGraphRequest -Method GET -Uri $uri
        
        if ($response.value) {
            Write-Host "✅ Retrieved $($response.value.Count) SKUs" -ForegroundColor Green
            return $response.value
        }
        else {
            Write-Host "⚠️  No SKUs found in tenant" -ForegroundColor Yellow
            return @()
        }
    }
    catch {
        Write-Host "❌ Error fetching SKUs: $_" -ForegroundColor Red
        throw
    }
}

function Export-SkusToJson {
    <#
    .SYNOPSIS
        Exports SKU data to JSON format
    #>
    param (
        [Parameter(Mandatory=$true)]
        [array]$Skus,
        
        [Parameter(Mandatory=$true)]
        [string]$FilePath
    )
    
    try {
        $exportData = @{
            exportDate = (Get-Date -Format "yyyy-MM-dd HH:mm:ss")
            tenantId = $TenantId
            totalSkus = $Skus.Count
            skus = $Skus | Select-Object skuId, skuPartNumber, capabilityStatus, consumedUnits, @{
                Name = 'prepaidUnitsEnabled'
                Expression = { $_.prepaidUnits.enabled }
            }, servicePlans
        }
        
        $exportData | ConvertTo-Json -Depth 10 | Out-File -FilePath $FilePath -Encoding UTF8
        Write-Host "✅ Exported SKUs to: $FilePath" -ForegroundColor Green
    }
    catch {
        Write-Host "❌ Error exporting SKUs to JSON: $_" -ForegroundColor Red
        throw
    }
}


function Export-GoConstants {
    <#
    .SYNOPSIS
        Generates Go constants from SKUs and service plans, grouped by domain
    #>
    param (
        [Parameter(Mandatory=$true)]
        [array]$Skus,
        
        [Parameter(Mandatory=$true)]
        [string]$FilePath
    )
    
    try {
        # Collect unique SKU part numbers (sanitize invisible Unicode)
        $skuPartNumbers = @{}
        foreach ($sku in $Skus) {
            if ($sku.skuPartNumber) {
                $sanitized = $sku.skuPartNumber -replace '[\u200B-\u200D\uFEFF]', ''
                $skuPartNumbers[$sanitized] = $true
            }
        }
        
        # Group service plans by their parent SKU and identify shared plans
        $groupedData = Group-ServicePlansBySku -Skus $Skus
        $sharedPlans = $groupedData.SharedPlans
        $skuSpecificPlans = $groupedData.SkuSpecificPlans
        
        # Sort SKUs
        $sortedSkus = $skuPartNumbers.Keys | Sort-Object
        
        # Generate Go constants
        $output = @"
// Auto-generated from Microsoft Graph API
// Generated: $(Get-Date -Format "yyyy-MM-dd HH:mm:ss")

package constants

// ============================================================================
// SKU Part Numbers ($($sortedSkus.Count) total)
// ============================================================================
const (
"@

        foreach ($skuName in $sortedSkus) {
            $constName = Convert-ToGoConstantName -Name $skuName
            $output += "`n`tSKU$constName = `"$skuName`""
        }

        $output += "`n)"
        $output += "`n"
        
        # Output shared service plans first
        if ($sharedPlans.Count -gt 0) {
            $output += "`n// ============================================================================"
            $output += "`n// Shared Service Plans ($($sharedPlans.Count) plans)"
            $output += "`n// These service plans appear in multiple SKUs"
            $output += "`n// ============================================================================"
            $output += "`nconst ("
            
            foreach ($planName in $sharedPlans.Keys | Sort-Object) {
                $constName = Convert-ToGoConstantName -Name $planName
                $skuList = $sharedPlans[$planName] -join ", "
                $output += "`n`tServicePlan$constName = `"$planName`" // Shared: $skuList"
            }
            
            $output += "`n)"
            $output += "`n"
        }
        
        # Output SKU-specific service plans
        $totalPlans = $sharedPlans.Count
        foreach ($skuName in $skuSpecificPlans.Keys | Sort-Object) {
            $plans = $skuSpecificPlans[$skuName] | Sort-Object
            $totalPlans += $plans.Count
            
            $output += "`n// ============================================================================"
            $output += "`n// Service Plans from: $skuName ($($plans.Count) plans)"
            $output += "`n// ============================================================================"
            $output += "`nconst ("
            
            foreach ($planName in $plans) {
                $constName = Convert-ToGoConstantName -Name $planName
                $output += "`n`tServicePlan$constName = `"$planName`""
            }
            
            $output += "`n)"
            $output += "`n"
        }

        $output | Out-File -FilePath $FilePath -Encoding UTF8
        Write-Host "✅ Generated Go constants to: $FilePath" -ForegroundColor Green
        Write-Host "   • $($sortedSkus.Count) SKUs" -ForegroundColor Gray
        Write-Host "   • $totalPlans Service Plans ($($sharedPlans.Count) shared, $($skuSpecificPlans.Count) SKU-specific)" -ForegroundColor Gray
    }
    catch {
        Write-Host "❌ Error generating Go constants: $_" -ForegroundColor Red
        throw
    }
}

function Group-ServicePlansBySku {
    <#
    .SYNOPSIS
        Groups service plans by their parent SKU and identifies shared plans
    #>
    param (
        [Parameter(Mandatory=$true)]
        [array]$Skus
    )
    
    # Track which SKUs each service plan appears in
    $planToSkus = @{}
    $skuToPlans = @{}
    
    foreach ($sku in $Skus) {
        if (-not $sku.skuPartNumber) {
            continue
        }
        
        # Sanitize SKU name
        $skuName = $sku.skuPartNumber -replace '[\u200B-\u200D\uFEFF]', ''
        
        # Initialize array for this SKU if needed
        if (-not $skuToPlans.ContainsKey($skuName)) {
            $skuToPlans[$skuName] = @()
        }
        
        # Track all service plans from this SKU
        foreach ($plan in $sku.servicePlans) {
            if ($plan.servicePlanName) {
                $sanitizedPlan = $plan.servicePlanName -replace '[\u200B-\u200D\uFEFF]', ''
                
                # Track which SKUs this plan belongs to
                if (-not $planToSkus.ContainsKey($sanitizedPlan)) {
                    $planToSkus[$sanitizedPlan] = @()
                }
                if ($planToSkus[$sanitizedPlan] -notcontains $skuName) {
                    $planToSkus[$sanitizedPlan] += $skuName
                }
                
                # Add plan to SKU list
                if ($skuToPlans[$skuName] -notcontains $sanitizedPlan) {
                    $skuToPlans[$skuName] += $sanitizedPlan
                }
            }
        }
    }
    
    # Identify shared plans (appear in 2+ SKUs)
    $sharedPlans = @{}
    foreach ($plan in $planToSkus.Keys) {
        if ($planToSkus[$plan].Count -gt 1) {
            $sharedPlans[$plan] = $planToSkus[$plan] | Sort-Object
        }
    }
    
    # Create SKU-specific plans (only unique to that SKU)
    $skuSpecificPlans = @{}
    foreach ($skuName in $skuToPlans.Keys) {
        $uniquePlans = @()
        foreach ($plan in $skuToPlans[$skuName]) {
            if (-not $sharedPlans.ContainsKey($plan)) {
                $uniquePlans += $plan
            }
        }
        if ($uniquePlans.Count -gt 0) {
            $skuSpecificPlans[$skuName] = $uniquePlans
        }
    }
    
    return @{
        SharedPlans = $sharedPlans
        SkuSpecificPlans = $skuSpecificPlans
    }
}

function Convert-ToGoConstantName {
    <#
    .SYNOPSIS
        Converts a string to a readable PascalCase Go constant name
    #>
    param (
        [Parameter(Mandatory=$true)]
        [string]$Name
    )
    
    # Remove zero-width spaces and other invisible Unicode characters
    $Name = $Name -replace '[\u200B-\u200D\uFEFF]', ''
    
    # Split on underscores, hyphens, and spaces
    $parts = $Name -split '[_\-\s]' | Where-Object { $_ -ne '' }
    
    $result = ""
    
    foreach ($part in $parts) {
        if ($part -match '^\d+$') {
            # If part is all digits, keep as-is
            $result += $part
        }
        elseif ($part -match '^\d') {
            # If starts with digit, spell out the number or prefix
            $result += $part
        }
        else {
            # Capitalize first letter, lowercase rest (unless already all caps)
            if ($part -ceq $part.ToUpper() -and $part.Length -gt 1) {
                # If ALL CAPS and longer than 1 char, keep it (e.g., "AAD", "MTP")
                $result += $part
            }
            else {
                # Convert to PascalCase
                $result += $part.Substring(0,1).ToUpper() + $part.Substring(1).ToLower()
            }
        }
    }
    
    # Handle special cases for common acronyms/abbreviations
    $result = $result -replace 'Aad', 'AAD'
    $result = $result -replace 'Atp', 'ATP'
    $result = $result -replace 'Ems', 'EMS'
    $result = $result -replace 'Mfa', 'MFA'
    $result = $result -replace 'Mtp', 'MTP'
    $result = $result -replace 'Rms', 'RMS'
    $result = $result -replace 'Dlp', 'DLP'
    $result = $result -replace 'Pam', 'PAM'
    $result = $result -replace 'Pim', 'PIM'
    $result = $result -replace 'Api', 'API'
    $result = $result -replace 'Sso', 'SSO'
    $result = $result -replace 'Id([A-Z]|$)', 'ID$1'
    $result = $result -replace 'O365', 'O365'
    $result = $result -replace 'M365', 'M365'
    $result = $result -replace 'E5', 'E5'
    $result = $result -replace 'E3', 'E3'
    $result = $result -replace 'P1', 'P1'
    $result = $result -replace 'P2', 'P2'
    $result = $result -replace 'P3', 'P3'
    
    return $result
}

#endregion

#region Main Script

Write-Host ""
Write-Host "╔════════════════════════════════════════════════════════════════╗" -ForegroundColor Cyan
Write-Host "║         Microsoft 365 License Export Tool                     ║" -ForegroundColor Cyan
Write-Host "╚════════════════════════════════════════════════════════════════╝" -ForegroundColor Cyan
Write-Host ""

try {
    # Authenticate to Microsoft Graph
    Write-Host "🔐 Authenticating to Microsoft Graph..." -ForegroundColor Cyan
    $secureClientSecret = ConvertTo-SecureString -String $ClientSecret -AsPlainText -Force
    $clientSecretCredential = New-Object -TypeName System.Management.Automation.PSCredential -ArgumentList $ClientId, $secureClientSecret
    
    Connect-MgGraph -ClientSecretCredential $clientSecretCredential -TenantId $TenantId -NoWelcome
    Write-Host "✅ Authentication successful" -ForegroundColor Green
    Write-Host ""
    
    # Create output directory if it doesn't exist
    if (-not (Test-Path $OutputPath)) {
        New-Item -ItemType Directory -Path $OutputPath -Force | Out-Null
    }
    
    # Fetch all SKUs
    $skus = Get-AllSubscribedSkus
    
    if ($skus.Count -eq 0) {
        Write-Host "⚠️  No SKUs found. Exiting." -ForegroundColor Yellow
        exit 0
    }
    
    Write-Host ""
    Write-Host "📦 Export Summary:" -ForegroundColor Cyan
    Write-Host "   • Total SKUs: $($skus.Count)" -ForegroundColor White
    
    $totalServicePlans = 0
    foreach ($sku in $skus) {
        $totalServicePlans += $sku.servicePlans.Count
    }
    Write-Host "   • Total Service Plans: $totalServicePlans" -ForegroundColor White
    Write-Host ""
    
    # Generate timestamp for filenames
    $timestamp = Get-Date -Format "yyyyMMdd_HHmmss"
    
    # Export to different formats
    Write-Host "📄 Exporting data..." -ForegroundColor Cyan
    Write-Host ""
    
    $jsonPath = Join-Path $OutputPath "licenses_$timestamp.json"
    $goConstantsPath = Join-Path $OutputPath "generated_constants.go"
    
    Export-SkusToJson -Skus $skus -FilePath $jsonPath
    Export-GoConstants -Skus $skus -FilePath $goConstantsPath
    
    Write-Host ""
    Write-Host "╔════════════════════════════════════════════════════════════════╗" -ForegroundColor Green
    Write-Host "║                    Export Complete!                            ║" -ForegroundColor Green
    Write-Host "╚════════════════════════════════════════════════════════════════╝" -ForegroundColor Green
    Write-Host ""
    Write-Host "📁 Exported files:" -ForegroundColor Cyan
    Write-Host "   • JSON (full data):   $jsonPath" -ForegroundColor White
    Write-Host "   • Go constants:       $goConstantsPath" -ForegroundColor White
    Write-Host ""
}
catch {
    Write-Host ""
    Write-Host "❌ Script failed with error:" -ForegroundColor Red
    Write-Host $_.Exception.Message -ForegroundColor Red
    Write-Host ""
    Write-Host "Stack Trace:" -ForegroundColor Yellow
    Write-Host $_.ScriptStackTrace -ForegroundColor Yellow
    exit 1
}
finally {
    # Disconnect from Microsoft Graph
    try {
        Disconnect-MgGraph | Out-Null
        Write-Host "🔓 Disconnected from Microsoft Graph" -ForegroundColor Gray
        Write-Host ""
    }
    catch {
        # Silently ignore disconnect errors
    }
}

#endregion
