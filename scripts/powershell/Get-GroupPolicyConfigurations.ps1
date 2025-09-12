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
    HelpMessage="Specific Group Policy Configuration ID (if not provided, will list all configurations)")]
    [string]$PolicyId,
    
    [Parameter(Mandatory=$false,
    HelpMessage="Get definition values for a specific configuration")]
    [string]$GetDefinitionValuesForId,
    
    [Parameter(Mandatory=$false,
    HelpMessage="Export results to JSON file")]
    [switch]$ExportToJson
)

# List all configurations and export to JSON
#./Get-GroupPolicyConfigurations.ps1 -TenantId "your-tenant-id" -ClientId "your-client-id" -ClientSecret "your-secret" -ExportToJson

# Get specific configuration and its definition values with export
#./Get-GroupPolicyConfigurations.ps1 -TenantId "your-tenant-id" -ClientId "your-client-id" -ClientSecret "your-secret" -PolicyId "specific-id" -GetDefinitionValuesForId "specific-id" -ExportToJson

# Get all configurations without export
#./Get-GroupPolicyConfigurations.ps1 -TenantId "your-tenant-id" -ClientId "your-client-id" -ClientSecret "your-secret"

# NEW API CALLS ADDED:
# - GET /deviceManagement/groupPolicyConfigurations('{id}')/definitionValues('{defValueId}')/presentationValues
# - GET /deviceManagement/groupPolicyConfigurations('{id}')/definitionValues('{defValueId}')/presentationValues('{presValueId}')?$expand=presentation
# NOTE: The /definitionValue endpoint doesn't exist - using $expand=presentation instead

Import-Module Microsoft.Graph.Authentication

function Get-GroupPolicyConfigurations {
    param (
        [Parameter(Mandatory=$false)]
        [string]$SpecificConfigurationId
    )
    try {
        if ($SpecificConfigurationId) {
            $uri = "https://graph.microsoft.com/beta/deviceManagement/groupPolicyConfigurations/$SpecificConfigurationId"
            Write-Host "🔍 Getting specific Group Policy Configuration..." -ForegroundColor Cyan
            Write-Host "   Configuration ID: $SpecificConfigurationId" -ForegroundColor Gray
        } else {
            $uri = "https://graph.microsoft.com/beta/deviceManagement/groupPolicyConfigurations"
            Write-Host "🔍 Getting all Group Policy Configurations..." -ForegroundColor Cyan
        }
        Write-Host "   Endpoint: $uri" -ForegroundColor Gray
        Write-Host ""
        $response = Invoke-MgGraphRequest -Method GET -Uri $uri
        return @{
            Response = $response
            Endpoint = $uri
        }
    }
    catch {
        Write-Host "❌ Error getting Group Policy Configurations: $_" -ForegroundColor Red
        Write-Host ""
        return $null
    }
}

function Get-GroupPolicyDefinitionValues {
    param (
        [Parameter(Mandatory=$true)]
        [string]$ConfigurationId
    )
    try {
        $uri = "https://graph.microsoft.com/beta/deviceManagement/groupPolicyConfigurations('$ConfigurationId')/definitionValues?`$expand=definition(`$select=id,classType,displayName,policyType,hasRelatedDefinitions,version,minUserCspVersion,minDeviceCspVersion)"
        Write-Host "🔍 Getting definition values for configuration: $ConfigurationId" -ForegroundColor Cyan
        Write-Host "   Endpoint: $uri" -ForegroundColor Gray
        Write-Host ""
        $response = Invoke-MgGraphRequest -Method GET -Uri $uri
        return @{
            Response = $response
            Endpoint = $uri
        }
    }
    catch {
        Write-Host "❌ Error getting definition values: $_" -ForegroundColor Red
        Write-Host ""
        return $null
    }
}

function Get-GroupPolicyConfigurationWithAssignments {
    param (
        [Parameter(Mandatory=$true)]
        [string]$ConfigurationId
    )
    try {
        $uri = "https://graph.microsoft.com/beta/deviceManagement/groupPolicyConfigurations/$ConfigurationId" + "?`$expand=assignments"
        Write-Host "🔍 Getting configuration with assignments: $ConfigurationId" -ForegroundColor Cyan
        Write-Host "   Endpoint: $uri" -ForegroundColor Gray
        Write-Host ""
        $response = Invoke-MgGraphRequest -Method GET -Uri $uri
        return @{
            Response = $response
            Endpoint = $uri
        }
    }
    catch {
        Write-Host "❌ Error getting configuration with assignments: $_" -ForegroundColor Red
        Write-Host ""
        return $null
    }
}

function Get-GroupPolicyDefinitionPresentations {
    param (
        [Parameter(Mandatory=$true)]
        [string]$DefinitionId
    )
    try {
        $uri = "https://graph.microsoft.com/beta/deviceManagement/groupPolicyDefinitions('$DefinitionId')/presentations"
        Write-Host "🔍 Getting presentations for definition: $DefinitionId" -ForegroundColor Cyan
        Write-Host "   Endpoint: $uri" -ForegroundColor Gray
        Write-Host ""
        $response = Invoke-MgGraphRequest -Method GET -Uri $uri
        return @{
            Response = $response
            Endpoint = $uri
        }
    }
    catch {
        Write-Host "❌ Error getting definition presentations: $_" -ForegroundColor Red
        Write-Host ""
        return $null
    }
}

function Get-GroupPolicyPresentationValues {
    param (
        [Parameter(Mandatory=$true)]
        [string]$ConfigurationId,
        [Parameter(Mandatory=$true)]
        [string]$DefinitionValueId
    )
    try {
        $uri = "https://graph.microsoft.com/beta/deviceManagement/groupPolicyConfigurations('$ConfigurationId')/definitionValues('$DefinitionValueId')/presentationValues"
        Write-Host "🔍 Getting presentation values for definition value: $DefinitionValueId" -ForegroundColor Cyan
        Write-Host "   Configuration ID: $ConfigurationId" -ForegroundColor Gray
        Write-Host "   Endpoint: $uri" -ForegroundColor Gray
        Write-Host ""
        $response = Invoke-MgGraphRequest -Method GET -Uri $uri
        return @{
            Response = $response
            Endpoint = $uri
        }
    }
    catch {
        Write-Host "❌ Error getting presentation values: $_" -ForegroundColor Red
        Write-Host ""
        return $null
    }
}

function Get-GroupPolicyPresentationValueDetails {
    param (
        [Parameter(Mandatory=$true)]
        [string]$ConfigurationId,
        [Parameter(Mandatory=$true)]
        [string]$DefinitionValueId,
        [Parameter(Mandatory=$true)]
        [string]$PresentationValueId
    )
    try {
        # The /definitionValue endpoint doesn't exist - instead get the presentation value with expand
        $uri = "https://graph.microsoft.com/beta/deviceManagement/groupPolicyConfigurations('$ConfigurationId')/definitionValues('$DefinitionValueId')/presentationValues('$PresentationValueId')?`$expand=presentation"
        Write-Host "🔍 Getting detailed presentation value: $PresentationValueId" -ForegroundColor Cyan
        Write-Host "   Configuration ID: $ConfigurationId" -ForegroundColor Gray
        Write-Host "   Definition Value ID: $DefinitionValueId" -ForegroundColor Gray
        Write-Host "   Endpoint: $uri" -ForegroundColor Gray
        Write-Host ""
        $response = Invoke-MgGraphRequest -Method GET -Uri $uri
        return @{
            Response = $response
            Endpoint = $uri
        }
    }
    catch {
        Write-Host "❌ Error getting presentation value details: $_" -ForegroundColor Red
        Write-Host ""
        return $null
    }
}

function Get-GroupPolicyDefinitionDetails {
    param (
        [Parameter(Mandatory=$true)]
        [string]$DefinitionId
    )
    try {
        $uri = "https://graph.microsoft.com/beta/deviceManagement/groupPolicyDefinitions('$DefinitionId')"
        Write-Host "🔍 Getting group policy definition details: $DefinitionId" -ForegroundColor Cyan
        Write-Host "   Endpoint: $uri" -ForegroundColor Gray
        Write-Host ""
        $response = Invoke-MgGraphRequest -Method GET -Uri $uri
        return @{
            Response = $response
            Endpoint = $uri
        }
    }
    catch {
        Write-Host "❌ Error getting group policy definition details: $_" -ForegroundColor Red
        Write-Host ""
        return $null
    }
}

function Add-ApiCallMetadata {
    param (
        [Parameter(Mandatory=$true)]
        $Response,
        [Parameter(Mandatory=$true)]
        [string]$ApiEndpoint,
        [Parameter(Mandatory=$true)]
        [string]$HttpMethod,
        [Parameter(Mandatory=$false)]
        [string]$Description = ""
    )
    try {
        $timestamp = Get-Date -Format "yyyy-MM-ddTHH:mm:ss.fffZ"
        
        # Create wrapper object with metadata and original response separate
        $wrappedResponse = [ordered]@{
            "_api_call_info" = [ordered]@{
                "timestamp" = $timestamp
                "http_method" = $HttpMethod.ToUpper()
                "endpoint" = $ApiEndpoint
                "description" = $Description
                "graph_api_version" = "beta"
            }
            "response" = $Response
        }
        
        return $wrappedResponse
    }
    catch {
        Write-Host "⚠️ Warning: Could not add API metadata, returning original response: $_" -ForegroundColor Yellow
        return $Response
    }
}

function Save-ApiResponseToJson {
    param (
        [Parameter(Mandatory=$true)]
        $Response,
        [Parameter(Mandatory=$true)]
        [string]$FileName,
        [Parameter(Mandatory=$true)]
        [string]$ApiEndpoint,
        [Parameter(Mandatory=$false)]
        [string]$HttpMethod = "GET",
        [Parameter(Mandatory=$false)]
        [string]$Description = ""
    )
    try {
        $timestamp = Get-Date -Format "yyyy-MM-dd-HH-mm-ss"
        $fullFileName = "$($FileName)_$timestamp.json"
        
        # Add API call metadata to the response
        $responseWithMetadata = Add-ApiCallMetadata -Response $Response -ApiEndpoint $ApiEndpoint -HttpMethod $HttpMethod -Description $Description
        
        $responseWithMetadata | ConvertTo-Json -Depth 15 | Out-File -FilePath $fullFileName -Encoding UTF8
        
        if ($Description) {
            Write-Host "📄 $Description saved to: $fullFileName" -ForegroundColor Green
        } else {
            Write-Host "📄 Response saved to: $fullFileName" -ForegroundColor Green
        }
        Write-Host "   API Endpoint: $ApiEndpoint" -ForegroundColor Gray
        Write-Host ""
        
        return $fullFileName
    }
    catch {
        Write-Host "❌ Error saving to JSON: $_" -ForegroundColor Red
        Write-Host ""
        return $null
    }
}


function Export-ConfigurationsToJson {
    param (
        [Parameter(Mandatory=$true)]
        $Configurations,
        [Parameter(Mandatory=$false)]
        [string]$SpecificConfigurationId
    )
    try {
        $timestamp = Get-Date -Format "yyyy-MM-dd-HH-mm-ss"
        if ($SpecificConfigurationId) {
            $filename = "GroupPolicyConfiguration_$($SpecificConfigurationId)_$timestamp.json"
        } else {
            $filename = "GroupPolicyConfigurations_All_$timestamp.json"
        }
        $Configurations | ConvertTo-Json -Depth 10 | Out-File -FilePath $filename -Encoding UTF8
        Write-Host "📄 Exported to: $filename" -ForegroundColor Green
        Write-Host ""
    }
    catch {
        Write-Host "❌ Error exporting to JSON: $_" -ForegroundColor Red
        Write-Host ""
    }
}

function Show-ConfigurationDetails {
    param (
        [Parameter(Mandatory=$true)]
        $Configuration
    )
    Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
    Write-Host "📋 Group Policy Configuration Details" -ForegroundColor Magenta
    Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
    
    # Basic Information
    Write-Host "🔷 Basic Information:" -ForegroundColor Green
    $basicFields = @('id', 'displayName', 'description', 'createdDateTime', 'lastModifiedDateTime')
    foreach ($field in $basicFields) {
        if ($Configuration.PSObject.Properties[$field] -and $Configuration.$field) {
            Write-Host ("   · {0}: {1}" -f $field, $Configuration.$field) -ForegroundColor Yellow
        }
    }
    
    # Role Scope Tags
    if ($Configuration.roleScopeTagIds -and $Configuration.roleScopeTagIds.Count -gt 0) {
        Write-Host "🔷 Role Scope Tag IDs:" -ForegroundColor Green
        foreach ($tagId in $Configuration.roleScopeTagIds) {
            Write-Host "   · $tagId" -ForegroundColor Yellow
        }
    }
    
    Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
    Write-Host ""
}

function Show-DefinitionValuesDetails {
    param (
        [Parameter(Mandatory=$true)]
        $DefinitionValues
    )
    Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
    Write-Host "⚙️ Definition Values Details" -ForegroundColor Magenta
    Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
    
    if ($DefinitionValues.value -and $DefinitionValues.value.Count -gt 0) {
        Write-Host "📊 Found $($DefinitionValues.value.Count) definition value(s)" -ForegroundColor Green
        Write-Host ""
        
        for ($i = 0; $i -lt $DefinitionValues.value.Count; $i++) {
            $defValue = $DefinitionValues.value[$i]
            Write-Host "Definition Value $($i + 1):" -ForegroundColor Magenta
            
            $defFields = @('id', 'enabled', 'configurationType', 'createdDateTime', 'lastModifiedDateTime')
            foreach ($field in $defFields) {
                if ($defValue.PSObject.Properties[$field] -and $null -ne $defValue.$field) {
                    Write-Host ("   · {0}: {1}" -f $field, $defValue.$field) -ForegroundColor Yellow
                }
            }
            
            if ($defValue.definition) {
                Write-Host "   • Definition Details:" -ForegroundColor Green
                $def = $defValue.definition
                $defDetailFields = @('id', 'displayName', 'classType', 'policyType', 'version')
                foreach ($field in $defDetailFields) {
                    if ($def.PSObject.Properties[$field] -and $def.$field) {
                        Write-Host ("     · {0}: {1}" -f $field, $def.$field) -ForegroundColor Cyan
                    }
                }
            }
            Write-Host ""
        }
    } else {
        Write-Host "📊 No definition values found" -ForegroundColor Yellow
    }
    
    Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
    Write-Host ""
}

# Main execution
try {
    Write-Host "🔐 Connecting to Microsoft Graph..." -ForegroundColor Cyan
    Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
    $secureClientSecret = ConvertTo-SecureString -String $ClientSecret -AsPlainText -Force
    $clientSecretCredential = New-Object -TypeName System.Management.Automation.PSCredential -ArgumentList $ClientId, $secureClientSecret
    Connect-MgGraph -ClientSecretCredential $clientSecretCredential -TenantId $TenantId
    Write-Host "✅ Connected to Microsoft Graph successfully" -ForegroundColor Green
    Write-Host ""

    # Get configurations
    $configurationsResult = Get-GroupPolicyConfigurations -SpecificConfigurationId $PolicyId
    $configurations = $configurationsResult.Response
    
    if ($ExportToJson) {
        Write-Host "🚀 Starting comprehensive data model export..." -ForegroundColor Magenta
        Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Magenta
        
        # 1. Export main configurations
        Save-ApiResponseToJson -Response $configurations -FileName "GroupPolicyConfigurations_All" -ApiEndpoint $configurationsResult.Endpoint -Description "All Group Policy Configurations"
        
        # 2. If specific configuration, get detailed data
        if ($PolicyId) {
            # Get configuration with assignments
            $configWithAssignmentsResult = Get-GroupPolicyConfigurationWithAssignments -ConfigurationId $PolicyId
            if ($configWithAssignmentsResult) {
                Save-ApiResponseToJson -Response $configWithAssignmentsResult.Response -FileName "GroupPolicyConfiguration_WithAssignments_$PolicyId" -ApiEndpoint $configWithAssignmentsResult.Endpoint -Description "Configuration with assignments"
                Show-ConfigurationDetails -Configuration $configWithAssignmentsResult.Response
            }
            
            # Get definition values
            $definitionValuesResult = Get-GroupPolicyDefinitionValues -ConfigurationId $PolicyId
            if ($definitionValuesResult) {
                $definitionValues = $definitionValuesResult.Response
                Save-ApiResponseToJson -Response $definitionValues -FileName "GroupPolicyDefinitionValues_$PolicyId" -ApiEndpoint $definitionValuesResult.Endpoint -Description "Definition values for configuration"
                Show-DefinitionValuesDetails -DefinitionValues $definitionValues
                
                # Get presentations and presentation values for each definition
                if ($definitionValues.value) {
                    foreach ($defValue in $definitionValues.value) {
                        if ($defValue.definition -and $defValue.definition.id) {
                            # Get group policy definition details (including categoryPath)
                            $definitionDetailsResult = Get-GroupPolicyDefinitionDetails -DefinitionId $defValue.definition.id
                            if ($definitionDetailsResult) {
                                Save-ApiResponseToJson -Response $definitionDetailsResult.Response -FileName "GroupPolicyDefinitionDetails_$($defValue.definition.id)" -ApiEndpoint $definitionDetailsResult.Endpoint -Description "Definition details for $($defValue.definition.displayName) including categoryPath"
                            }
                            
                            # Get available presentations
                            $presentationsResult = Get-GroupPolicyDefinitionPresentations -DefinitionId $defValue.definition.id
                            if ($presentationsResult) {
                                Save-ApiResponseToJson -Response $presentationsResult.Response -FileName "GroupPolicyDefinitionPresentations_$($defValue.definition.id)" -ApiEndpoint $presentationsResult.Endpoint -Description "Presentations for definition $($defValue.definition.displayName)"
                            }
                            
                            # Get actual presentation values (current configuration)
                            if ($defValue.id) {
                                $presentationValuesResult = Get-GroupPolicyPresentationValues -ConfigurationId $PolicyId -DefinitionValueId $defValue.id
                                if ($presentationValuesResult -and $presentationValuesResult.Response.value) {
                                    Save-ApiResponseToJson -Response $presentationValuesResult.Response -FileName "GroupPolicyPresentationValues_$($PolicyId)_$($defValue.id)" -ApiEndpoint $presentationValuesResult.Endpoint -Description "Current presentation values for $($defValue.definition.displayName)"
                                    
                                    # Get detailed presentation value information
                                    foreach ($presValue in $presentationValuesResult.Response.value) {
                                        if ($presValue.id) {
                                            $presValueDetailsResult = Get-GroupPolicyPresentationValueDetails -ConfigurationId $PolicyId -DefinitionValueId $defValue.id -PresentationValueId $presValue.id
                                            if ($presValueDetailsResult) {
                                                Save-ApiResponseToJson -Response $presValueDetailsResult.Response -FileName "GroupPolicyPresentationValueDetails_$($PolicyId)_$($defValue.id)_$($presValue.id)" -ApiEndpoint $presValueDetailsResult.Endpoint -Description "Detailed info for presentation value $($presValue.id)"
                                            }
                                        }
                                    }
                                }
                            }
                        }
                    }
                }
            }
        } else {
            # For all configurations, get detailed data for each
            if ($configurations.value -and $configurations.value.Count -gt 0) {
                Write-Host "📊 Found $($configurations.value.Count) Group Policy Configuration(s)" -ForegroundColor Green
                Write-Host "🔄 Getting detailed data for each configuration..." -ForegroundColor Cyan
                Write-Host ""
                
                foreach ($config in $configurations.value) {
                    Write-Host "📋 Processing configuration: $($config.displayName) ($($config.id))" -ForegroundColor Yellow
                    
                    # Get configuration with assignments
                    $configWithAssignmentsResult = Get-GroupPolicyConfigurationWithAssignments -ConfigurationId $config.id
                    if ($configWithAssignmentsResult) {
                        Save-ApiResponseToJson -Response $configWithAssignmentsResult.Response -FileName "GroupPolicyConfiguration_WithAssignments_$($config.id)" -ApiEndpoint $configWithAssignmentsResult.Endpoint -Description "Configuration with assignments for $($config.displayName)"
                    }
                    
                    # Get definition values
                    $definitionValuesResult = Get-GroupPolicyDefinitionValues -ConfigurationId $config.id
                    if ($definitionValuesResult) {
                        $definitionValues = $definitionValuesResult.Response
                        Save-ApiResponseToJson -Response $definitionValues -FileName "GroupPolicyDefinitionValues_$($config.id)" -ApiEndpoint $definitionValuesResult.Endpoint -Description "Definition values for $($config.displayName)"
                        
                        # Get presentations and presentation values for each definition
                        if ($definitionValues.value) {
                            foreach ($defValue in $definitionValues.value) {
                                if ($defValue.definition -and $defValue.definition.id) {
                                    # Get group policy definition details (including categoryPath)
                                    $definitionDetailsResult = Get-GroupPolicyDefinitionDetails -DefinitionId $defValue.definition.id
                                    if ($definitionDetailsResult) {
                                        Save-ApiResponseToJson -Response $definitionDetailsResult.Response -FileName "GroupPolicyDefinitionDetails_$($defValue.definition.id)" -ApiEndpoint $definitionDetailsResult.Endpoint -Description "Definition details for $($defValue.definition.displayName) including categoryPath"
                                    }
                                    
                                    # Get available presentations
                                    $presentationsResult = Get-GroupPolicyDefinitionPresentations -DefinitionId $defValue.definition.id
                                    if ($presentationsResult) {
                                        Save-ApiResponseToJson -Response $presentationsResult.Response -FileName "GroupPolicyDefinitionPresentations_$($defValue.definition.id)" -ApiEndpoint $presentationsResult.Endpoint -Description "Presentations for definition $($defValue.definition.displayName)"
                                    }
                                    
                                    # Get actual presentation values (current configuration)
                                    if ($defValue.id) {
                                        $presentationValuesResult = Get-GroupPolicyPresentationValues -ConfigurationId $config.id -DefinitionValueId $defValue.id
                                        if ($presentationValuesResult -and $presentationValuesResult.Response.value) {
                                            Save-ApiResponseToJson -Response $presentationValuesResult.Response -FileName "GroupPolicyPresentationValues_$($config.id)_$($defValue.id)" -ApiEndpoint $presentationValuesResult.Endpoint -Description "Current presentation values for $($defValue.definition.displayName) in $($config.displayName)"
                                            
                                            # Get detailed presentation value information
                                            foreach ($presValue in $presentationValuesResult.Response.value) {
                                                if ($presValue.id) {
                                                    $presValueDetailsResult = Get-GroupPolicyPresentationValueDetails -ConfigurationId $config.id -DefinitionValueId $defValue.id -PresentationValueId $presValue.id
                                                    if ($presValueDetailsResult) {
                                                        Save-ApiResponseToJson -Response $presValueDetailsResult.Response -FileName "GroupPolicyPresentationValueDetails_$($config.id)_$($defValue.id)_$($presValue.id)" -ApiEndpoint $presValueDetailsResult.Endpoint -Description "Detailed info for presentation value $($presValue.id) in $($config.displayName)"
                                                    }
                                                }
                                            }
                                        }
                                    }
                                }
                            }
                        }
                    }
                    
                    Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Gray
                }
            } else {
                Write-Host "📊 No Group Policy Configurations found" -ForegroundColor Yellow
            }
        }
        
        # Get definition values if specifically requested (and it's a different ID than already processed)
        if ($GetDefinitionValuesForId -and $GetDefinitionValuesForId -ne "" -and $GetDefinitionValuesForId -ne $PolicyId) {
            Write-Host "📋 Processing additional definition values for ID: $GetDefinitionValuesForId" -ForegroundColor Yellow
            $definitionValuesResult = Get-GroupPolicyDefinitionValues -ConfigurationId $GetDefinitionValuesForId
            if ($definitionValuesResult) {
                Save-ApiResponseToJson -Response $definitionValuesResult.Response -FileName "GroupPolicyDefinitionValues_$GetDefinitionValuesForId" -ApiEndpoint $definitionValuesResult.Endpoint -Description "Definition values for requested configuration"
                Show-DefinitionValuesDetails -DefinitionValues $definitionValuesResult.Response
            }
        }
        
        Write-Host "✅ Comprehensive data model export completed successfully!" -ForegroundColor Green
        Write-Host "🎯 All API responses saved to separate JSON files for analysis" -ForegroundColor Green
        return
    }

    # Display configurations (only when not exporting)
    if ($PolicyId) {
        if ($configurations) {
            Show-ConfigurationDetails -Configuration $configurations
        }
    } else {
        if ($configurations.value -and $configurations.value.Count -gt 0) {
            Write-Host "📊 Found $($configurations.value.Count) Group Policy Configuration(s)" -ForegroundColor Green
            Write-Host ""
            for ($i = 0; $i -lt $configurations.value.Count; $i++) {
                Write-Host "Configuration $($i + 1):" -ForegroundColor Magenta
                Show-ConfigurationDetails -Configuration $configurations.value[$i]
            }
        } elseif ($configurations -and -not $configurations.value) {
            Write-Host "📊 Found 1 Group Policy Configuration" -ForegroundColor Green
            Write-Host ""
            Show-ConfigurationDetails -Configuration $configurations
        } else {
            Write-Host "📊 No Group Policy Configurations found" -ForegroundColor Yellow
        }
    }

    # Get definition values if requested (only when not exporting)
    if ($GetDefinitionValuesForId) {
        $definitionValuesResult = Get-GroupPolicyDefinitionValues -ConfigurationId $GetDefinitionValuesForId
        if ($definitionValuesResult) {
            Show-DefinitionValuesDetails -DefinitionValues $definitionValuesResult.Response
        }
    }

    Write-Host "✅ Group Policy Configuration data retrieval completed successfully!" -ForegroundColor Green
    Write-Host ""
}
catch {
    Write-Host "❌ An error occurred: $_" -ForegroundColor Red
    Write-Host ""
}
finally {
    try {
        Disconnect-MgGraph | Out-Null
        Write-Host "🔓 Disconnected from Microsoft Graph" -ForegroundColor Gray
    }
    catch {
        # Ignore disconnection errors
    }
}
