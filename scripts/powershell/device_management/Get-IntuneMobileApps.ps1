# Get-IntuneMobileApps.ps1
# Script to get mobile apps from Intune via Microsoft Graph API and save to JSON file

[CmdletBinding()]
param (
    [Parameter(Mandatory=$true,
    HelpMessage="Specify the Entra ID tenant ID (Directory ID) where the application is registered")]
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
    HelpMessage="Optional ID of a specific mobile app to retrieve")]
    [string]$MobileAppId,
    
    [Parameter(Mandatory=$false,
    HelpMessage="Optional filter query for mobile apps")]
    [string]$Filter,
    
    [Parameter(Mandatory=$false,
    HelpMessage="Optional select query to specify which properties to retrieve")]
    [string]$Select,
    
    [Parameter(Mandatory=$false,
    HelpMessage="Optional expand query to include related entities")]
    [string]$Expand,
    
    [Parameter(Mandatory=$false,
    HelpMessage="Maximum number of apps to return (for pagination)")]
    [int]$Top = 100,
    
    [Parameter(Mandatory=$true,
    HelpMessage="Path to output JSON file")]
    [string]$OutputFile
)

# Function to authenticate and get access token
function Connect-MicrosoftGraph {
    param (
        [Parameter(Mandatory=$true)]
        [string]$TenantId,
        
        [Parameter(Mandatory=$true)]
        [string]$ClientId,
        
        [Parameter(Mandatory=$true)]
        [string]$ClientSecret
    )
    
    try {
        Write-Host "Connecting to Microsoft Graph..." -ForegroundColor Cyan
        
        # Create secure credential
        $secureClientSecret = ConvertTo-SecureString -String $ClientSecret -AsPlainText -Force
        $clientSecretCredential = New-Object -TypeName System.Management.Automation.PSCredential -ArgumentList $ClientId, $secureClientSecret
        
        # Connect to Microsoft Graph
        Import-Module Microsoft.Graph.Authentication
        Connect-MgGraph -ClientSecretCredential $clientSecretCredential -TenantId $TenantId
        
        Write-Host "✅ Connected to Microsoft Graph" -ForegroundColor Green
    }
    catch {
        Write-Host "❌ Error connecting to Microsoft Graph: $_" -ForegroundColor Red
        throw
    }
}

# Function to get mobile apps from Intune and save as JSON file
function Get-MobileAppsToFile {
    param (
        [Parameter(Mandatory=$false)]
        [string]$MobileAppId,
        
        [Parameter(Mandatory=$false)]
        [string]$Filter,
        
        [Parameter(Mandatory=$false)]
        [string]$Select,
        
        [Parameter(Mandatory=$false)]
        [string]$Expand,
        
        [Parameter(Mandatory=$false)]
        [int]$Top = 100,
        
        [Parameter(Mandatory=$true)]
        [string]$OutputFile
    )
    
    try {
        $baseUrl = "https://graph.microsoft.com/beta/deviceAppManagement/mobileApps"
        $queryParams = @()
        
        # Build the URL based on parameters
        if ($MobileAppId) {
            $url = "$baseUrl/$MobileAppId"
        } else {
            $url = $baseUrl
            
            if ($Filter) {
                $queryParams += "`$filter=$([System.Web.HttpUtility]::UrlEncode($Filter))"
            }
            
            if ($Select) {
                $queryParams += "`$select=$([System.Web.HttpUtility]::UrlEncode($Select))"
            }
            
            if ($Expand) {
                $queryParams += "`$expand=$([System.Web.HttpUtility]::UrlEncode($Expand))"
            }
            
            if ($Top -gt 0) {
                $queryParams += "`$top=$Top"
            }
            
            if ($queryParams.Count -gt 0) {
                $url += "?" + ($queryParams -join "&")
            }
        }
        
        Write-Host "Retrieving mobile apps from Intune..." -ForegroundColor Cyan
        Write-Host "URL: $url" -ForegroundColor Gray
        
        # Get response and save directly to file
        $response = Invoke-MgGraphRequest -Method GET -Uri $url
        $prettyJson = ConvertTo-Json -InputObject $response -Depth 10 -Compress:$false
        $prettyJson | Out-File -FilePath $OutputFile -Encoding utf8
        
        # Get app count for summary
        if ($MobileAppId) {
            $appCount = 1
        } else {
            $appCount = if ($response.value) { $response.value.Count } else { 0 }
        }
        
        return $appCount
    }
    catch {
        Write-Host "❌ Error retrieving mobile apps: $_" -ForegroundColor Red
        throw
    }
}

# Main script execution
try {
    # Connect to Microsoft Graph
    Connect-MicrosoftGraph -TenantId $TenantId -ClientId $ClientId -ClientSecret $ClientSecret
    
    # Get mobile apps and save to file
    Write-Host "`n📱 Retrieving Intune mobile apps..." -ForegroundColor Cyan
    Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
    
    $appCount = Get-MobileAppsToFile -MobileAppId $MobileAppId -Filter $Filter -Select $Select -Expand $Expand -Top $Top -OutputFile $OutputFile
    
    # Summary message
    if ($MobileAppId) {
        Write-Host "`n✨ Successfully saved mobile app details to: $OutputFile" -ForegroundColor Green
    } else {
        Write-Host "`n✨ Successfully saved $appCount mobile apps to: $OutputFile" -ForegroundColor Green
    }
}
catch {
    Write-Host "`n❌ Script execution failed: $_" -ForegroundColor Red
    exit 1
}
finally {
    # Disconnect from Microsoft Graph
    Write-Host "`nDisconnecting from Microsoft Graph..." -ForegroundColor Cyan
    Disconnect-MgGraph > $null 2>&1
    Write-Host "✅ Disconnected from Microsoft Graph" -ForegroundColor Green
}