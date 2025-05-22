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

# Helper Function to get settings data with pagination support
function Get-AllSettingsCatalogMetadata {
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

# Script Setup
Import-Module Microsoft.Graph.Authentication

$secureClientSecret = ConvertTo-SecureString -String $ClientSecret -AsPlainText -Force
$clientSecretCredential = New-Object -TypeName System.Management.Automation.PSCredential -ArgumentList $ClientId, $secureClientSecret

Write-Host "Connecting to Microsoft Graph..."
Connect-MgGraph -ClientSecretCredential $clientSecretCredential -TenantId $TenantId

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

# Vars
$allSettingsData = Get-AllSettingsCatalogMetadata

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

    # Create setting object with ordered base fields
    $settingInfo = [ordered]@{
        # Base fields first
        oDataType = $_.'@odata.type'
        id = $_.id
        name = $_.name  
        displayName = $_.displayName
        description = $_.description
        keywords = $_.keywords
        helpText = $_.helpText
        visibility = $_.visibility
        accessTypes = $_.accessTypes
        
        # Complex objects
        settingDefinition = [ordered]@{
            rootDefinitionId = $_.rootDefinitionId
            categoryId = $_.categoryId
            settingUsage = $_.settingUsage
            baseUri = $_.baseUri
            offsetUri = $_.offsetUri
            version = $_.version
        }
        applicability = $_.applicability
    }

    # Add optional fields in a specific order if they exist
    $optionalFieldsOrder = @(
        @{Name='options'; Field='options'},
        @{Name='defaultValue'; Field='defaultValue'},
        @{Name='valueDefinition'; Field='valueDefinition'},
        @{Name='occurrence'; Field='occurrence'},
        @{Name='dependentOn'; Field='dependentOn'},
        @{Name='dependedOnBy'; Field='dependedOnBy'},
        @{Name='infoUrls'; Field='infoUrls'},
        @{Name='uxBehaviour'; Field='uxBehaviour'}
    )

    foreach ($field in $optionalFieldsOrder) {
        if ($null -ne $_."$($field.Field)") {
            if ($field.Field -eq 'options') {
                $settingInfo[$field.Name] = $_.options | ForEach-Object {
                    [ordered]@{
                        displayName = $_.displayName
                        description = $_.description
                        value = $_.value
                        optionId = $_.optionId
                    }
                }
            }
            else {
                $settingInfo[$field.Name] = $_."$($field.Field)"
            }
        }
    }

    # Add defaultOptionId to settingDefinition if it exists
    if ($null -ne $_.defaultOptionId) {
        $settingInfo.settingDefinition['defaultOptionId'] = $_.defaultOptionId
    }

    # Remove any null properties while preserving order
    $cleanedSettingInfo = [ordered]@{}
    $settingInfo.GetEnumerator() | Where-Object { $null -ne $_.Value } | ForEach-Object {
        $cleanedSettingInfo[$_.Key] = $_.Value
    }

    $settingsCatalogMetadata.platformDefinitions[$sortedPlatforms] += $cleanedSettingInfo
}

# Before export, sort the platform groups and definitions
$sortedPlatformGroups = [ordered]@{}
$settingsCatalogMetadata.metadata.platformGroups.GetEnumerator() | 
    Sort-Object -Property Name | 
    ForEach-Object { $sortedPlatformGroups[$_.Name] = $_.Value }

$settingsCatalogMetadata.metadata.platformGroups = $sortedPlatformGroups

$sortedPlatformDefinitions = [ordered]@{}
$settingsCatalogMetadata.platformDefinitions.GetEnumerator() | 
    Sort-Object -Property Name | 
    ForEach-Object { $sortedPlatformDefinitions[$_.Name] = $_.Value }

$settingsCatalogMetadata.platformDefinitions = $sortedPlatformDefinitions

# Export each platform definition to a separate file
foreach ($platform in $sortedPlatformDefinitions.GetEnumerator()) {
    $platformName = $platform.Name
    if ($platformName -like "*,*") {
        $platformName = $platformName.Replace(",", "_").Replace(" ", "")
        $fileName = "$($platformName)_shared_settings_catalog_metadata.json"
    } else {
        $fileName = "$($platformName)_settings_catalog_metadata.json"
    }
    
    $platformMetadata = @{
        metadata = @{
            version = $settingsCatalogMetadata.metadata.version
            exportTimestamp = $settingsCatalogMetadata.metadata.exportTimestamp
            totalSettingsCount = $platform.Value.Count
            platformGroups = @{
                $platform.Name = $settingsCatalogMetadata.metadata.platformGroups[$platform.Name]
            }
        }
        platformDefinitions = @{
            $platform.Name = $platform.Value
        }
    }

    $outputFile = Join-Path -Path $exportFolder -ChildPath $fileName
    $platformMetadata | ConvertTo-Json -Depth 10 | Out-File -FilePath $outputFile
    Write-Host "Exported $($platform.Name) settings to: $fileName"
}

Write-Host "`nExport completed:"
Write-Host "- Total settings exported: $($settingsCatalogMetadata.metadata.totalSettingsCount)"
Write-Host "`nPlatform groups exported:"
foreach ($group in $settingsCatalogMetadata.metadata.platformGroups.GetEnumerator() | Sort-Object Name) {
    Write-Host "  - $($group.Key): $($group.Value.count) settings"
    Write-Host "    Platforms: $($group.Value.platforms -join ', ')"
}
Write-Host "`n- Export location: $exportFolder"