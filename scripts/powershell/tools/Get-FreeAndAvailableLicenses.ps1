<#
.SYNOPSIS
    Gets all available and free license SKUs in the tenant.

.DESCRIPTION
    Retrieves all subscribed SKUs in the Microsoft 365 tenant and identifies those with available 
    capacity or free licenses. Displays detailed information about each license including SKU details,
    available capacity, and service plans.

.PARAMETER TenantId
    The Entra ID tenant ID (Directory ID).

.PARAMETER ClientId
    The application (client) ID for authentication.

.PARAMETER ClientSecret
    The client secret for authentication.

.PARAMETER ShowAllLicenses
    Show all licenses regardless of availability. Default: $false (only show licenses with available capacity)

.PARAMETER ExportToJson
    Export the results to a JSON file in the output directory.

.PARAMETER RequiredPermissions
    Required Microsoft Graph application permissions to validate.

.EXAMPLE
    # Get all available licenses
    .\Get-FreeAndAvailableLicenses.ps1 `
        -TenantId "your-tenant-id" `
        -ClientId "your-client-id" `
        -ClientSecret "your-secret"

.EXAMPLE
    # Get all licenses (including those without capacity) and export to JSON
    .\Get-FreeAndAvailableLicenses.ps1 `
        -TenantId "your-tenant-id" `
        -ClientId "your-client-id" `
        -ClientSecret "your-secret" `
        -ShowAllLicenses $true `
        -ExportToJson $true

.NOTES
    Author: Deployment Theory
    Requires: Microsoft.Graph.Authentication module
    Purpose: Discover available licenses for testing and configuration
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
    HelpMessage="Show all licenses regardless of availability")]
    [switch]$ShowAllLicenses,
    
    [Parameter(Mandatory=$false,
    HelpMessage="Export results to JSON file")]
    [switch]$ExportToJson,
    
    [Parameter(Mandatory=$false,
    HelpMessage="Required Microsoft Graph application permissions to validate")]
    [string[]]$RequiredPermissions = @("Organization.Read.All")
)

Import-Module Microsoft.Graph.Authentication

# Function to get all subscribed SKUs
function Get-SubscribedSKUs {
    try {
        Write-Host "🔍 Retrieving subscribed SKUs..." -ForegroundColor Cyan
        
        $allSkus = @()
        $uri = "https://graph.microsoft.com/v1.0/subscribedSkus"
        
        do {
            $response = Invoke-MgGraphRequest -Method GET -Uri $uri
            
            if ($response.value) {
                $allSkus += $response.value
            }
            
            $uri = $response.'@odata.nextLink'
        } while ($uri)
        
        Write-Host "✅ Retrieved $($allSkus.Count) license SKU(s)" -ForegroundColor Green
        return $allSkus
    }
    catch {
        Write-Host "❌ Failed to retrieve subscribed SKUs: $_" -ForegroundColor Red
        throw
    }
}

# Function to calculate license availability
function Get-LicenseAvailability {
    param (
        [Parameter(Mandatory=$true)]
        $Sku
    )
    
    $enabled = if ($Sku.prepaidUnits.enabled) { $Sku.prepaidUnits.enabled } else { 0 }
    $consumed = if ($Sku.consumedUnits) { $Sku.consumedUnits } else { 0 }
    $available = $enabled - $consumed
    
    return @{
        Enabled = $enabled
        Consumed = $consumed
        Available = $available
    }
}

# Function to determine if a license is considered "free"
function Test-IsFreeOrAvailable {
    param (
        [Parameter(Mandatory=$true)]
        $Sku,
        
        [Parameter(Mandatory=$true)]
        [hashtable]$Availability
    )
    
    # Check if it's a free SKU by name
    $isFree = $Sku.skuPartNumber -match "FREE" -or 
              $Sku.skuPartNumber -match "VIRAL"
    
    # Check if it has available capacity
    $hasCapacity = $Availability.Available -gt 0
    
    # Include unlimited licenses (enabled = 0 means unlimited for some SKUs)
    $isUnlimited = $Availability.Enabled -eq 0 -and $isFree
    
    return ($isFree -or $hasCapacity -or $isUnlimited)
}

# Function to display license details
function Show-LicenseDetails {
    param (
        [Parameter(Mandatory=$true)]
        $Licenses,
        
        [Parameter(Mandatory=$false)]
        [bool]$ShowAll = $false
    )
    
    Write-Host "`n╔════════════════════════════════════════════════════════════════════╗" -ForegroundColor Cyan
    Write-Host "║  📦 Microsoft 365 License Inventory                                ║" -ForegroundColor Cyan
    Write-Host "╚════════════════════════════════════════════════════════════════════╝" -ForegroundColor Cyan
    
    $freeLicenses = @()
    $availableLicenses = @()
    $totalLicenses = 0
    
    foreach ($sku in $Licenses | Sort-Object skuPartNumber) {
        $availability = Get-LicenseAvailability -Sku $sku
        $isAvailable = Test-IsFreeOrAvailable -Sku $sku -Availability $availability
        
        # Skip if not showing all and license is not available
        if (-not $ShowAll -and -not $isAvailable) {
            continue
        }
        
        $totalLicenses++
        
        # Determine license type
        $isFree = $sku.skuPartNumber -match "FREE" -or $sku.skuPartNumber -match "VIRAL"
        if ($isFree) {
            $freeLicenses += $sku
        }
        if ($availability.Available -gt 0) {
            $availableLicenses += $sku
        }
        
        Write-Host "`n──────────────────────────────────────────────────────────────────────" -ForegroundColor Gray
        Write-Host "📋 SKU Part Number: " -NoNewline -ForegroundColor White
        Write-Host "$($sku.skuPartNumber)" -ForegroundColor Yellow
        
        Write-Host "   SKU ID: " -NoNewline -ForegroundColor Gray
        Write-Host "$($sku.skuId)" -ForegroundColor Cyan
        
        # Display availability with color coding
        Write-Host "   Licenses:" -ForegroundColor Gray
        Write-Host "      Enabled:  " -NoNewline -ForegroundColor Gray
        Write-Host "$($availability.Enabled)" -ForegroundColor White
        Write-Host "      Consumed: " -NoNewline -ForegroundColor Gray
        Write-Host "$($availability.Consumed)" -ForegroundColor White
        Write-Host "      Available: " -NoNewline -ForegroundColor Gray
        
        $availColor = if ($availability.Available -gt 0) { "Green" } 
                     elseif ($availability.Available -eq 0 -and $availability.Enabled -gt 0) { "Yellow" }
                     else { "Red" }
        Write-Host "$($availability.Available)" -ForegroundColor $availColor
        
        # Display license type indicators
        $indicators = @()
        if ($isFree) { $indicators += "🆓 FREE" }
        if ($availability.Available -gt 0) { $indicators += "✅ AVAILABLE" }
        if ($availability.Enabled -eq 0 -and $isFree) { $indicators += "♾️  UNLIMITED" }
        
        if ($indicators.Count -gt 0) {
            Write-Host "   Type: " -NoNewline -ForegroundColor Gray
            Write-Host ($indicators -join " | ") -ForegroundColor Green
        }
        
        # Display service plans count
        if ($sku.servicePlans -and $sku.servicePlans.Count -gt 0) {
            Write-Host "   Service Plans: " -NoNewline -ForegroundColor Gray
            Write-Host "$($sku.servicePlans.Count) plan(s)" -ForegroundColor White
            
            # Show first few service plans as examples
            $displayCount = [Math]::Min(3, $sku.servicePlans.Count)
            for ($i = 0; $i -lt $displayCount; $i++) {
                $plan = $sku.servicePlans[$i]
                Write-Host "      • " -NoNewline -ForegroundColor DarkGray
                Write-Host "$($plan.servicePlanName)" -NoNewline -ForegroundColor White
                Write-Host " ($($plan.servicePlanId))" -ForegroundColor DarkGray
            }
            
            if ($sku.servicePlans.Count -gt $displayCount) {
                Write-Host "      ... and $($sku.servicePlans.Count - $displayCount) more" -ForegroundColor DarkGray
            }
        }
    }
    
    # Display summary
    Write-Host "`n╔════════════════════════════════════════════════════════════════════╗" -ForegroundColor Green
    Write-Host "║  📊 Summary                                                         ║" -ForegroundColor Green
    Write-Host "╚════════════════════════════════════════════════════════════════════╝" -ForegroundColor Green
    
    Write-Host "`n   Total Licenses Found: " -NoNewline -ForegroundColor White
    Write-Host "$totalLicenses" -ForegroundColor Cyan
    
    Write-Host "   Free Licenses: " -NoNewline -ForegroundColor White
    Write-Host "$($freeLicenses.Count)" -ForegroundColor Green
    
    Write-Host "   Licenses with Available Capacity: " -NoNewline -ForegroundColor White
    Write-Host "$($availableLicenses.Count)" -ForegroundColor Green
    
    # Recommend best licenses for testing
    Write-Host "`n╔════════════════════════════════════════════════════════════════════╗" -ForegroundColor Yellow
    Write-Host "║  💡 Recommended for Testing                                        ║" -ForegroundColor Yellow
    Write-Host "╚════════════════════════════════════════════════════════════════════╝" -ForegroundColor Yellow
    
    # Find best candidates for testing (free with service plans and capacity)
    $recommended = $Licenses | Where-Object {
        $avail = Get-LicenseAvailability -Sku $_
        $isFree = $_.skuPartNumber -match "FREE" -or $_.skuPartNumber -match "VIRAL"
        ($isFree -or $avail.Available -gt 100) -and $_.servicePlans.Count -gt 0
    } | Sort-Object { Get-LicenseAvailability -Sku $_ | Select-Object -ExpandProperty Available } -Descending | Select-Object -First 3
    
    if ($recommended) {
        foreach ($license in $recommended) {
            $avail = Get-LicenseAvailability -Sku $license
            Write-Host "`n   📦 $($license.skuPartNumber)" -ForegroundColor Cyan
            Write-Host "      SKU ID: $($license.skuId)" -ForegroundColor White
            Write-Host "      Available: $($avail.Available) license(s)" -ForegroundColor Green
            Write-Host "      Service Plans: $($license.servicePlans.Count)" -ForegroundColor Gray
            
            Write-Host "`n      Terraform Configuration:" -ForegroundColor Yellow
            Write-Host "      sku_id = `"$($license.skuId)`"  # $($license.skuPartNumber)" -ForegroundColor White
            
            if ($license.servicePlans.Count -gt 0) {
                Write-Host "`n      Example disabled_plans:" -ForegroundColor Yellow
                Write-Host "      disabled_plans = [" -ForegroundColor White
                $license.servicePlans | Select-Object -First 2 | ForEach-Object {
                    Write-Host "        `"$($_.servicePlanId)`",  # $($_.servicePlanName)" -ForegroundColor White
                }
                Write-Host "      ]" -ForegroundColor White
            }
        }
    } else {
        Write-Host "`n   ⚠️  No recommended licenses found for testing" -ForegroundColor Yellow
    }
    
    return @{
        Total = $totalLicenses
        Free = $freeLicenses.Count
        Available = $availableLicenses.Count
    }
}

# Function to export results to JSON
function Export-ResultsToJson {
    param (
        [Parameter(Mandatory=$true)]
        $Licenses,
        
        [Parameter(Mandatory=$true)]
        [hashtable]$Summary
    )
    
    try {
        # Create output directory if it doesn't exist
        $outputDir = Join-Path -Path (Get-Location) -ChildPath "output"
        if (-not (Test-Path -Path $outputDir)) {
            New-Item -Path $outputDir -ItemType Directory | Out-Null
            Write-Host "`n📁 Created output directory: $outputDir" -ForegroundColor Gray
        }
        
        # Generate timestamp for filename
        $timestamp = Get-Date -Format "yyyyMMdd_HHmmss"
        $fileName = "FreeAndAvailableLicenses_${timestamp}.json"
        $filePath = Join-Path -Path $outputDir -ChildPath $fileName
        
        # Build export object
        $exportData = @{
            Timestamp = Get-Date -Format "o"
            TenantId = $TenantId
            Summary = $Summary
            Licenses = @()
        }
        
        foreach ($sku in $Licenses) {
            $availability = Get-LicenseAvailability -Sku $sku
            $isFree = $sku.skuPartNumber -match "FREE" -or $sku.skuPartNumber -match "VIRAL"
            
            $exportData.Licenses += @{
                SkuPartNumber = $sku.skuPartNumber
                SkuId = $sku.skuId
                IsFree = $isFree
                Enabled = $availability.Enabled
                Consumed = $availability.Consumed
                Available = $availability.Available
                ServicePlans = $sku.servicePlans | ForEach-Object {
                    @{
                        ServicePlanId = $_.servicePlanId
                        ServicePlanName = $_.servicePlanName
                        AppliesTo = $_.appliesTo
                        ProvisioningStatus = $_.provisioningStatus
                    }
                }
            }
        }
        
        $exportData | ConvertTo-Json -Depth 10 | Out-File -FilePath $filePath -Encoding UTF8
        
        Write-Host "💾 Exported results to: $filePath" -ForegroundColor Green
        return $filePath
    }
    catch {
        Write-Host "❌ Error exporting results to JSON: $_" -ForegroundColor Red
        return $null
    }
}

# Main Script Execution
try {
    Write-Host ""
    Write-Host "╔════════════════════════════════════════════════════════════════════╗" -ForegroundColor Magenta
    Write-Host "║  🎫 Microsoft 365 License Discovery Tool                          ║" -ForegroundColor Magenta
    Write-Host "╚════════════════════════════════════════════════════════════════════╝" -ForegroundColor Magenta
    Write-Host ""
    
    Write-Host "🔐 Connecting to Microsoft Graph..." -ForegroundColor Cyan
    
    $secureSecret = ConvertTo-SecureString -String $ClientSecret -AsPlainText -Force
    $credential = New-Object -TypeName System.Management.Automation.PSCredential -ArgumentList $ClientId, $secureSecret
    
    Connect-MgGraph -ClientSecretCredential $credential -TenantId $TenantId -NoWelcome
    
    Write-Host "✅ Connected to Microsoft Graph" -ForegroundColor Green
    Write-Host ""
    
    # Get all subscribed SKUs
    $allLicenses = Get-SubscribedSKUs
    
    if (-not $allLicenses -or $allLicenses.Count -eq 0) {
        Write-Host "⚠️  No licenses found in tenant" -ForegroundColor Yellow
        return
    }
    
    # Filter licenses if not showing all
    if (-not $ShowAllLicenses) {
        $filteredLicenses = $allLicenses | Where-Object {
            $avail = Get-LicenseAvailability -Sku $_
            Test-IsFreeOrAvailable -Sku $_ -Availability $avail
        }
    } else {
        $filteredLicenses = $allLicenses
    }
    
    # Display license details
    $summary = Show-LicenseDetails -Licenses $filteredLicenses -ShowAll $ShowAllLicenses
    
    # Export if requested
    if ($ExportToJson) {
        Export-ResultsToJson -Licenses $filteredLicenses -Summary $summary
    }
    
    Write-Host "`n🎉 License discovery completed successfully!" -ForegroundColor Green
    Write-Host ""
}
catch {
    Write-Host ""
    Write-Host "❌ Error: $_" -ForegroundColor Red
    Write-Host ""
    exit 1
}
finally {
    Write-Host "🔌 Disconnecting from Microsoft Graph..." -ForegroundColor Cyan
    Disconnect-MgGraph 2>$null
    Write-Host "✅ Disconnected" -ForegroundColor Green
    Write-Host ""
}
