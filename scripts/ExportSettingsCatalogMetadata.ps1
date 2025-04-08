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
    HelpMessage="Specify the full path for the combined output JSON file.")]
    [string]$OutputFile # Optional: Defaults to a path relative to the script if not provided
)

# --- Helper Function to get settings data with pagination support ---
function Get-AllSettingsCatalogMetadata {
    [CmdletBinding()]
    param()

    $baseUri = "https://graph.microsoft.com/beta/deviceManagement/configurationSettings"
    $allSettings = [System.Collections.Generic.List[object]]::new() # Use a generic list for better performance
    $nextLink = $baseUri
    $totalCountEstimate = 0 # For progress indication

    # Optional: Get an estimated count first (might not be perfectly accurate)
    try {
        $countUri = $baseUri + '/$count'
        # Note: $count might require specific permissions or might not always be supported
        # $countResponse = Invoke-MgGraphRequest -Method GET -Uri $countUri -OutputType PSObject -ErrorAction SilentlyContinue
        # if ($countResponse -is [int]) {
        #     $totalCountEstimate = $countResponse
        #     Write-Host "Estimated total settings count: $totalCountEstimate"
        # } else {
             Write-Host "Retrieving settings catalog metadata (pagination enabled)..."
        # }
    } catch {
        Write-Warning "Could not retrieve estimated count. Proceeding with pagination."
        Write-Host "Retrieving settings catalog metadata (pagination enabled)..."
    }


    while ($null -ne $nextLink) {
        try {
            # Use -OutputType PSObject for easier property access
            $response = Invoke-MgGraphRequest -Method GET -Uri $nextLink -OutputType PSObject

            if ($null -ne $response.value) {
                $batchCount = $response.value.Count
                $allSettings.AddRange($response.value)
                # Provide progress based on count if available, otherwise just batches
                # if ($totalCountEstimate -gt 0) {
                #     Write-Host ("Retrieved {0} / ~{1} settings..." -f $allSettings.Count, $totalCountEstimate)
                # } else {
                     Write-Host ("Retrieved {0} settings (batch size: {1})..." -f $allSettings.Count, $batchCount)
                # }
            } else {
                 Write-Warning "Received a response with no 'value' array from URI: $nextLink"
            }

            # Check for nextLink using PSObject property access
            $nextLink = $response.'@odata.nextLink'

            # Safety break / throttling (optional)
            # Start-Sleep -Milliseconds 100

        }
        catch {
            Write-Error "Error retrieving settings batch from URI '$nextLink': $_"
            # Decide whether to break or try to continue
            $nextLink = $null # Stop pagination on error
        }
    }

    Write-Host "Finished retrieving settings. Total found: $($allSettings.Count)"
    return $allSettings
}

# --- Script Setup ---
# Ensure the necessary module is available
if (-not (Get-Module -ListAvailable -Name Microsoft.Graph.Authentication)) {
    Write-Error "Microsoft.Graph.Authentication module is required. Please install it using 'Install-Module Microsoft.Graph.Authentication -Scope CurrentUser'."
    return
}
Import-Module Microsoft.Graph.Authentication

# Securely handle the client secret
$secureClientSecret = ConvertTo-SecureString -String $ClientSecret -AsPlainText -Force
$clientSecretCredential = New-Object -TypeName System.Management.Automation.PSCredential -ArgumentList $ClientId, $secureClientSecret

# Connect to Microsoft Graph
Write-Host "Connecting to Microsoft Graph (Tenant: $TenantId, ClientId: $ClientId)..."
try {
    Connect-MgGraph -ClientSecretCredential $clientSecretCredential -TenantId $TenantId -ErrorAction Stop
    Write-Host "Successfully connected to Microsoft Graph."
} catch {
    Write-Error "Failed to connect to Microsoft Graph: $_"
    return
}

# Determine script path and output folder/file
$scriptPath = $PSScriptRoot
if (-not $scriptPath) {
    $scriptPath = Split-Path -Parent -Path $MyInvocation.MyCommand.Definition
}
if (-not $scriptPath) {
    $scriptPath = Get-Location
}

$exportFolder = Join-Path -Path $scriptPath -ChildPath "SettingsCatalogMetadataExport"
if (-not (Test-Path -Path $exportFolder)) {
    Write-Verbose "Creating export directory: $exportFolder"
    New-Item -ItemType Directory -Path $exportFolder | Out-Null
}

# Set default output file path if not provided
if (-not $PSBoundParameters.ContainsKey('OutputFile')) {
    $OutputFile = Join-Path -Path $exportFolder -ChildPath "all_platforms_settings_catalog_metadata.json"
} elseif (-not (Split-Path -Path $OutputFile -IsAbsolute)) {
    # Make relative output path relative to script location
    $OutputFile = Join-Path -Path $scriptPath -ChildPath $OutputFile
}

# --- Fetch Data ---
$allSettingsData = Get-AllSettingsCatalogMetadata
if ($null -eq $allSettingsData -or $allSettingsData.Count -eq 0) {
    Write-Error "Failed to retrieve any settings catalog metadata from Microsoft Graph."
    return
}

# --- Process and Structure Data ---
Write-Host "Processing $($allSettingsData.Count) retrieved settings..."

# Create the root structure
$settingsCatalogMetadata = @{
    metadata = [ordered]@{
        version = Get-Date -Format "yyyy-MM-dd"
        totalSettingsCount = $allSettingsData.Count # Correct total count
        platformGroups = [ordered]@{} # Will be populated
        exportTimestamp = (Get-Date).ToUniversalTime().ToString("o")
    }
    platformDefinitions = [ordered]@{} # Will be populated
}

# Group settings by unique platform combinations
$groupedSettings = $allSettingsData | Group-Object {
    $platforms = if ($_.applicability.platform) { $_.applicability.platform -split ',' } else { @("Unspecified") }
    ($platforms | Sort-Object) -join ',' # Use sorted comma-separated string as the group key
}

# Populate platformDefinitions and platformGroups
Write-Host "Organizing settings by platform applicability..."
foreach ($group in $groupedSettings) {
    $platformKey = $group.Name # e.g., "iOS", "macOS", "iOS,macOS"
    $platformSettingsList = $group.Group

    # Store the actual settings under this key in platformDefinitions
    $settingsCatalogMetadata.platformDefinitions[$platformKey] = $platformSettingsList | ForEach-Object {
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
            @{Name='options'; Field='options'}, # Handle options separately below
            @{Name='defaultValue'; Field='defaultValue'},
            @{Name='valueDefinition'; Field='valueDefinition'},
            @{Name='occurrence'; Field='occurrence'},
            @{Name='dependentOn'; Field='dependentOn'},
            @{Name='dependedOnBy'; Field='dependedOnBy'},
            @{Name='infoUrls'; Field='infoUrls'},
            @{Name='uxBehaviour'; Field='uxBehaviour'} # Added uxBehaviour based on original script comments
        )

        foreach ($field in $optionalFieldsOrder) {
            if ($_.PSObject.Properties.Name -contains $field.Field) { # More robust check if property exists
                 $propertyValue = $_.$($field.Field)
                 if ($null -ne $propertyValue) {
                    if ($field.Field -eq 'options' -and $propertyValue -is [array]) {
                        # Ensure options preserve order if needed (though likely not critical for schema gen)
                        $settingInfo[$field.Name] = $propertyValue | ForEach-Object {
                            [ordered]@{
                                displayName = $_.displayName
                                description = $_.description
                                value = $_.value
                                optionId = $_.optionId
                            }
                        }
                    }
                    else {
                        $settingInfo[$field.Name] = $propertyValue
                    }
                }
            }
        }

        # Add defaultOptionId to settingDefinition if it exists on the main object
        if ($_.PSObject.Properties.Name -contains 'defaultOptionId' -and $null -ne $_.defaultOptionId) {
            $settingInfo.settingDefinition['defaultOptionId'] = $_.defaultOptionId
        }

        # Return the ordered hashtable for this setting
        $settingInfo
    }

    # Add entry to platformGroups metadata
     $originalPlatformNames = if ($platformKey -eq "Unspecified") { @("Unspecified") } else { $platformKey -split ',' }
     $settingsCatalogMetadata.metadata.platformGroups[$platformKey] = @{
         count = $platformSettingsList.Count
         platforms = $originalPlatformNames # Store original platform names as an array
     }
     Write-Host "  - Processed group '$platformKey': $($platformSettingsList.Count) settings."
}

# Sort the platform groups and definitions keys for consistent output
Write-Host "Sorting final structure..."
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

# --- Export Combined Data ---
Write-Host "Exporting combined data to: $OutputFile"
try {
    # Use Depth 10 or higher to ensure nested objects are fully converted
    $settingsCatalogMetadata | ConvertTo-Json -Depth 10 | Out-File -FilePath $OutputFile -Encoding UTF8 -ErrorAction Stop
    Write-Host "Successfully exported all $($settingsCatalogMetadata.metadata.totalSettingsCount) settings to '$OutputFile'."
} catch {
    Write-Error "Failed to write output file '$OutputFile': $_"
    return
}

Write-Host "`nExport Summary:"
Write-Host "- Total settings exported: $($settingsCatalogMetadata.metadata.totalSettingsCount)"
Write-Host "`nPlatform groups included:"
foreach ($group in $settingsCatalogMetadata.metadata.platformGroups.GetEnumerator()) {
    Write-Host ("  - [{0}] Count: {1}, Platforms: {2}" -f $group.Key, $group.Value.count, ($group.Value.platforms -join ', '))
}
Write-Host "`n- Export location: $OutputFile"

# Disconnect (Optional - good practice in scripts)
# Disconnect-MgGraph