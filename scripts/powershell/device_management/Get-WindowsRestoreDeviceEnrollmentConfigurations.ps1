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
    HelpMessage="Specific Windows Restore Device Enrollment Configuration ID (if not provided, will list all configurations)")]
    [string]$ConfigurationId,
    
    [Parameter(Mandatory=$true,
    HelpMessage="Export results to JSON file")]
    [bool]$ExportToJson
)

# Import required modules
Import-Module Microsoft.Graph.Authentication

# Function to get Windows Restore device enrollment configurations
function Get-WindowsRestoreDeviceEnrollmentConfigurations {
    param (
        [Parameter(Mandatory=$false)]
        [string]$SpecificConfigurationId
    )
    
    try {
        if ($SpecificConfigurationId) {
            # GET specific configuration - can use $expand for a single resource
            $uri = "https://graph.microsoft.com/beta/deviceManagement/deviceEnrollmentConfigurations/$SpecificConfigurationId"
            Write-Host "🔍 Getting specific Windows Restore device enrollment configuration..." -ForegroundColor Cyan
            Write-Host "   Configuration ID: $SpecificConfigurationId" -ForegroundColor Gray
        } else {
            # GET all Windows Restore configurations using filter
            $uri = "https://graph.microsoft.com/beta/deviceManagement/deviceEnrollmentConfigurations?`$filter=deviceEnrollmentConfigurationType eq 'windowsRestore'"
            Write-Host "🔍 Getting all Windows Restore device enrollment configurations..." -ForegroundColor Cyan
        }
        
        Write-Host "   Endpoint: $uri" -ForegroundColor Gray
        Write-Host ""
        
        $response = Invoke-MgGraphRequest -Method GET -Uri $uri
        
        return $response
    }
    catch {
        Write-Host "❌ Error getting Windows Restore device enrollment configurations: $_" -ForegroundColor Red
        Write-Host ""
        
        # Enhanced error handling
        if ($_.Exception.Response) {
            $statusCode = $_.Exception.Response.StatusCode
            $statusDescription = $_.Exception.Response.StatusDescription
            Write-Host "   Status Code: $statusCode" -ForegroundColor Red
            Write-Host "   Status Description: $statusDescription" -ForegroundColor Red
            
            # Try to get the response content
            try {
                $responseContent = $_.Exception.Response.GetResponseStream()
                if ($responseContent) {
                    $reader = [System.IO.StreamReader]::new($responseContent)
                    $errorDetails = $reader.ReadToEnd()
                    $reader.Close()
                    Write-Host "   Error Details: $errorDetails" -ForegroundColor Red
                }
            }
            catch {
                Write-Host "   Could not read error details" -ForegroundColor Red
            }
        }
        
        throw
    }
}

# Function to export configurations to JSON
function Export-ConfigurationsToJson {
    param (
        [Parameter(Mandatory=$true)]
        $Configurations,
        
        [Parameter(Mandatory=$false)]
        [string]$SpecificConfigurationId
    )
    
    try {
        # Create output directory if it doesn't exist
        $outputDir = Join-Path -Path (Get-Location) -ChildPath "output"
        if (-not (Test-Path -Path $outputDir)) {
            New-Item -Path $outputDir -ItemType Directory | Out-Null
        }
        
        # Generate filename with timestamp
        $timestamp = Get-Date -Format "yyyyMMdd_HHmmss"
        if ($SpecificConfigurationId) {
            $filename = "WindowsRestoreDeviceEnrollmentConfiguration_$SpecificConfigurationId`_$timestamp.json"
        } else {
            $filename = "WindowsRestoreDeviceEnrollmentConfigurations_$timestamp.json"
        }
        
        $filePath = Join-Path -Path $outputDir -ChildPath $filename
        
        # Convert to JSON with proper formatting
        $jsonContent = $Configurations | ConvertTo-Json -Depth 10
        
        # Write to file
        $jsonContent | Out-File -FilePath $filePath -Encoding UTF8
        
        Write-Host "✅ Configuration(s) exported to: $filePath" -ForegroundColor Green
        Write-Host ""
        
        return $filePath
    }
    catch {
        Write-Host "❌ Error exporting configurations to JSON: $_" -ForegroundColor Red
        Write-Host ""
        throw
    }
}

# Function to display configuration summary
function Show-ConfigurationSummary {
    param (
        [Parameter(Mandatory=$true)]
        $Configurations
    )
    
    try {
        Write-Host "📊 Configuration Summary" -ForegroundColor Yellow
        Write-Host "========================" -ForegroundColor Yellow
        
        if ($Configurations.value) {
            $configCount = $Configurations.value.Count
            Write-Host "   Total Configurations Found: $configCount" -ForegroundColor White
            
            foreach ($config in $Configurations.value) {
                Write-Host ""
                Write-Host "   Configuration ID: $($config.id)" -ForegroundColor Cyan
                Write-Host "   Display Name: $($config.displayName)" -ForegroundColor White
                Write-Host "   Description: $($config.description)" -ForegroundColor Gray
                Write-Host "   State: $($config.state)" -ForegroundColor $(if ($config.state -eq "enabled") { "Green" } else { "Red" })
                Write-Host "   Priority: $($config.priority)" -ForegroundColor White
                Write-Host "   Created: $($config.createdDateTime)" -ForegroundColor Gray
                Write-Host "   Last Modified: $($config.lastModifiedDateTime)" -ForegroundColor Gray
                Write-Host "   Version: $($config.version)" -ForegroundColor White
                
                if ($config.roleScopeTagIds) {
                    Write-Host "   Role Scope Tag IDs: $($config.roleScopeTagIds -join ', ')" -ForegroundColor Gray
                }
                
                if ($config.assignments) {
                    Write-Host "   Assignments: $($config.assignments.Count) assignment(s)" -ForegroundColor Gray
                }
            }
        } else {
            Write-Host "   No configurations found" -ForegroundColor Red
        }
        
        Write-Host ""
    }
    catch {
        Write-Host "❌ Error displaying configuration summary: $_" -ForegroundColor Red
        Write-Host ""
    }
}

# Main execution
try {
    Write-Host "🚀 Windows Restore Device Enrollment Configuration Retrieval Script" -ForegroundColor Magenta
    Write-Host "=================================================================" -ForegroundColor Magenta
    Write-Host ""
    
    # Connect to Microsoft Graph
    Write-Host "🔐 Connecting to Microsoft Graph..." -ForegroundColor Cyan
    Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
    
    $secureClientSecret = ConvertTo-SecureString -String $ClientSecret -AsPlainText -Force
    $clientSecretCredential = New-Object -TypeName System.Management.Automation.PSCredential -ArgumentList $ClientId, $secureClientSecret
    
    Connect-MgGraph -ClientSecretCredential $clientSecretCredential -TenantId $TenantId
    
    Write-Host "✅ Successfully connected to Microsoft Graph" -ForegroundColor Green
    Write-Host ""
    
    # Get configurations
    $configurations = Get-WindowsRestoreDeviceEnrollmentConfigurations -SpecificConfigurationId $ConfigurationId
    
    # Display summary
    Show-ConfigurationSummary -Configurations $configurations
    
    # Export to JSON if requested
    if ($ExportToJson) {
        $exportedFile = Export-ConfigurationsToJson -Configurations $configurations -SpecificConfigurationId $ConfigurationId
    }
    
    Write-Host "✅ Script completed successfully!" -ForegroundColor Green
    Write-Host ""
    
    # Return the configurations for further processing if needed
    return $configurations
}
catch {
    Write-Host "❌ Script failed with error: $_" -ForegroundColor Red
    Write-Host ""
    exit 1
}
finally {
    # Disconnect from Microsoft Graph
    Write-Host "🔌 Disconnecting from Microsoft Graph..." -ForegroundColor Cyan
    try {
        Disconnect-MgGraph 2>$null
        Write-Host "✅ Disconnected from Microsoft Graph" -ForegroundColor Green
    }
    catch {
        # Ignore disconnect errors
    }
}
