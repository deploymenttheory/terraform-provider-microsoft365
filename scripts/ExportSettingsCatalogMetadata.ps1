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
    [string]$ClientSecret
)

# Script Setup
Import-Module Microsoft.Graph.Authentication

$secureClientSecret = ConvertTo-SecureString -String $ClientSecret -AsPlainText -Force
$clientSecretCredential = New-Object -TypeName System.Management.Automation.PSCredential -ArgumentList $ClientId, $secureClientSecret

Write-Host "Connecting to Microsoft Graph..."
Connect-MgGraph -ClientSecretCredential $clientSecretCredential -TenantId $TenantId

# Create export directory
$scriptPath = $PSScriptRoot
if (-not $scriptPath) {
    $scriptPath = Split-Path -Parent -Path $MyInvocation.MyCommand.Definition
}
if (-not $scriptPath) {
    $scriptPath = Get-Location
}

$exportFolder = Join-Path -Path $scriptPath -ChildPath "SettingsCatalogMetadataExport"
if (-not (Test-Path -Path $exportFolder)) {
    New-Item -ItemType Directory -Path $exportFolder | Out-Null
}

# Function to get settings data with pagination support
function Get-AllSettingsData {
    $baseUri = "https://graph.microsoft.com/beta/deviceManagement/configurationSettings"
    $allSettings = @()
    $nextLink = $baseUri

    Write-Host "Retrieving settings catalog metadata..."
    
    while ($nextLink) {
        try {
            $response = Invoke-MgGraphRequest -Method GET -Uri $nextLink
            $allSettings += $response.value
            $nextLink = $response.'@odata.nextLink'
            
            Write-Host "Retrieved $($allSettings.Count) settings so far..."
        }
        catch {
            Write-Error "Error retrieving settings: $_"
            break
        }
    }

    return $allSettings
}

# Get all settings
$allSettingsData = Get-AllSettingsData

# Create organized settings structure
$settingsCatalogMetadata = @{
    metadata = @{
        version = Get-Date -Format "yyyy-MM-dd"
        totalSettingsCount = $allSettingsData.Count
        platformGroups = @{}
        exportTimestamp = (Get-Date).ToUniversalTime().ToString("o")
    }
    platformDefinitions = @{}
}

Write-Host "Processing settings data..."

# First pass: Count settings per platform for metadata
$allSettingsData | ForEach-Object {
    $platforms = if ($_.applicability.platform) {
        $_.applicability.platform -split ','
    } else {
        @("Unspecified")
    }
    
    $sortedPlatforms = ($platforms | Sort-Object) -join ','
    
    if (-not $settingsCatalogMetadata.metadata.platformGroups.ContainsKey($sortedPlatforms)) {
        $settingsCatalogMetadata.metadata.platformGroups[$sortedPlatforms] = @{
            count = 0
            platforms = $platforms
        }
    }
    $settingsCatalogMetadata.metadata.platformGroups[$sortedPlatforms].count++
}

# Second pass: Organize settings by platform
$allSettingsData | ForEach-Object {
    $platforms = if ($_.applicability.platform) {
        $_.applicability.platform -split ','
    } else {
        @("Unspecified")
    }
    
    $sortedPlatforms = ($platforms | Sort-Object) -join ','
    
    if (-not $settingsCatalogMetadata.platformDefinitions.ContainsKey($sortedPlatforms)) {
        $settingsCatalogMetadata.platformDefinitions[$sortedPlatforms] = @()
    }

    # Create setting object with specific fields
    $settingInfo = @{
        oDataType = $_.'@odata.type'
        id = $_.id
        name = $_.name
        displayName = $_.displayName
        description = $_.description
        settingDefinition = @{
            rootDefinitionId = $_.rootDefinitionId
            categoryId = $_.categoryId
            settingUsage = $_.settingUsage
            baseUri = $_.baseUri
            offsetUri = $_.offsetUri
            version = $_.version
        }
        applicability = $_.applicability
        keywords = $_.keywords
        accessTypes = $_.accessTypes
    }

    # Add optional fields if they exist
    if ($null -ne $_.defaultOptionId) {
        $settingInfo.settingDefinition['defaultOptionId'] = $_.defaultOptionId
    }
    if ($null -ne $_.infoUrls) {
        $settingInfo['infoUrls'] = $_.infoUrls
    }
    if ($null -ne $_.uxBehaviour) {
        $settingInfo['uxBehaviour'] = $_.uxBehaviour
    }
    if ($null -ne $_.occurrence) {
        $settingInfo['occurrence'] = $_.occurrence
    }
    if ($null -ne $_.options) {
        $settingInfo['options'] = $_.options | ForEach-Object {
            @{
                displayName = $_.displayName
                description = $_.description
                value = $_.value
                optionId = $_.optionId
            }
        }
    }

    # Remove any null properties
    $settingInfo = $settingInfo.GetEnumerator() | Where-Object { $null -ne $_.Value } | 
                  ForEach-Object { $h = @{} } { $h[$_.Key] = $_.Value } { $h }

    # Add to platform group
    $settingsCatalogMetadata.platformDefinitions[$sortedPlatforms] += $settingInfo
}

# Export the complete metadata to a single file
$outputFile = Join-Path -Path $exportFolder -ChildPath "SettingsCatalogMetadata.json"
$settingsCatalogMetadata | ConvertTo-Json -Depth 10 | Out-File -FilePath $outputFile

Write-Host "`nExport completed:"
Write-Host "- Total settings exported: $($settingsCatalogMetadata.metadata.totalSettingsCount)"
Write-Host "`nPlatform groups found:"
foreach ($group in $settingsCatalogMetadata.metadata.platformGroups.GetEnumerator() | Sort-Object Name) {
    Write-Host "  - $($group.Key): $($group.Value.count) settings"
    Write-Host "    Platforms: $($group.Value.platforms -join ', ')"
}
Write-Host "`n- Export location: $outputFile"

Disconnect-MgGraph
Write-Host "Disconnected from Microsoft Graph"