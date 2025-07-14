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
    HelpMessage="Specify the ID of the enterprise app catalog package to retrieve configuration for")]
    [string]$PackageId,
    
    [Parameter(Mandatory = $false,
        HelpMessage = "Search term to find app configurations by display name")]
    [string]$SearchTerm,
    
    [Parameter(Mandatory = $false,
        HelpMessage = "Maximum number of results to return")]
    [int]$MaxResults = 50,
    
    [Parameter(Mandatory = $false,
        HelpMessage = "Export results to JSON file")]
    [bool]$ExportToJson = $true,
    
    [Parameter(Mandatory = $false,
        HelpMessage = "Path to output JSON file")]
    [string]$OutputFile = ""
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
        
        # Check if the Microsoft.Graph module is installed
        if (-not (Get-Module -ListAvailable -Name Microsoft.Graph)) {
            Write-Host "📦 Installing Microsoft.Graph PowerShell module..." -ForegroundColor Cyan
            Install-Module -Name Microsoft.Graph -Scope CurrentUser -Force
        }
        
        # Import the Microsoft.Graph module
        Import-Module Microsoft.Graph.Authentication
        
        # Create secure credential
        $secureClientSecret = ConvertTo-SecureString -String $ClientSecret -AsPlainText -Force
        $clientSecretCredential = New-Object -TypeName System.Management.Automation.PSCredential -ArgumentList $ClientId, $secureClientSecret
        
        # Connect to Microsoft Graph
        Connect-MgGraph -ClientSecretCredential $clientSecretCredential -TenantId $TenantId
        
        Write-Host "✅ Connected to Microsoft Graph" -ForegroundColor Green
    }
    catch {
        Write-Host "❌ Error connecting to Microsoft Graph: $_" -ForegroundColor Red
        throw
    }
}

# Function to make Graph API requests with retry logic
function Invoke-GraphApiRequest {
    param (
        [Parameter(Mandatory = $true)]
        [string]$Uri,
        
        [Parameter(Mandatory = $false)]
        [string]$Method = "GET",
        
        [Parameter(Mandatory = $false)]
        [hashtable]$Headers = @{},
        
        [Parameter(Mandatory = $false)]
        [object]$Body = $null
    )
    
    $retryCount = 0
    $maxRetries = 3
    $baseWaitTime = 2
    
    while ($retryCount -le $maxRetries) {
        try {
            $authToken = Get-MgContext
            if (-not $authToken) {
                throw "Not authenticated to Microsoft Graph"
            }
            
            $Headers["ConsistencyLevel"] = "eventual"
            
            $params = @{
                Uri     = $Uri
                Method  = $Method
                Headers = $Headers
            }
            
            if ($Body -and $Method -ne "GET") {
                $params["Body"] = $Body | ConvertTo-Json -Depth 20
                $params["ContentType"] = "application/json"
            }
            
            $response = Invoke-MgGraphRequest @params -ErrorAction Stop
            return $response
        }
        catch {
            $retryCount++
            $statusCode = $_.Exception.Response.StatusCode.value__
            
            # Don't retry if it's an authentication or authorization issue
            if ($statusCode -eq 401 -or $statusCode -eq 403) {
                Write-Host "❌ Authentication error ($statusCode): $_" -ForegroundColor Red
                throw
            }
            
            if ($retryCount -gt $maxRetries) {
                Write-Host "❌ Max retry attempts reached for request to $Uri. Error: $_" -ForegroundColor Red
                throw
            }
            
            $waitTime = $baseWaitTime * [Math]::Pow(2, $retryCount - 1)
            Write-Host "⚠️ Request attempt $retryCount failed (Status: $statusCode). Retrying in $waitTime seconds..." -ForegroundColor Yellow
            Start-Sleep -Seconds $waitTime
        }
    }
}

# Function to get app configuration by package ID
function Get-AppConfigurationByPackageId {
    param (
        [Parameter(Mandatory = $true)]
        [string]$PackageId
    )
    
    try {
        $url = "https://graph.microsoft.com/beta/deviceAppManagement/mobileApps/convertFromMobileAppCatalogPackage(mobileAppCatalogPackageId='$PackageId')"
        
        Write-Host "🔍 Fetching app configuration for package ID: $PackageId" -ForegroundColor Cyan
        Write-Host "   URL: $url" -ForegroundColor Gray
        
        $response = Invoke-GraphApiRequest -Uri $url
        
        if ($response) {
            Write-Host "✅ Retrieved app configuration" -ForegroundColor Green
            return $response
        }
        else {
            Write-Host "⚠️ No app configuration found for package ID: $PackageId" -ForegroundColor Yellow
            return $null
        }
    }
    catch {
        Write-Host "❌ Error getting app configuration: $_" -ForegroundColor Red
        throw
    }
}

# Function to get enterprise app catalog packages
function Get-EnterpriseAppCatalogPackages {
    param (
        [Parameter(Mandatory = $false)]
        [string]$SearchTerm,
        
        [Parameter(Mandatory = $false)]
        [int]$MaxResults = 50
    )
    
    try {
        $baseUrl = "https://graph.microsoft.com/beta/deviceAppManagement/mobileAppCatalogPackages"
        $queryParams = @()
        
        # Add $top parameter
        $queryParams += "`$top=$MaxResults"
        
        # Handle search term if provided
        if ($SearchTerm) {
            $queryParams += "`$search=`"$SearchTerm`""
        }
        
        # Build the final URL
        $url = $baseUrl
        if ($queryParams.Count -gt 0) {
            $url += "?" + ($queryParams -join "&")
        }
        
        Write-Host "🔍 Fetching enterprise app catalog packages..." -ForegroundColor Cyan
        Write-Host "   URL: $url" -ForegroundColor Gray
        
        $response = Invoke-GraphApiRequest -Uri $url
        
        if ($response -and $response.value) {
            Write-Host "✅ Retrieved $($response.value.Count) package(s)" -ForegroundColor Green
            return $response.value
        }
        else {
            Write-Host "⚠️ No packages found" -ForegroundColor Yellow
            return @()
        }
    }
    catch {
        Write-Host "❌ Error getting enterprise app catalog packages: $_" -ForegroundColor Red
        throw
    }
}

# Function to get app configurations for multiple packages
function Get-AppConfigurations {
    param (
        [Parameter(Mandatory = $true)]
        [array]$Packages
    )
    
    try {
        $appConfigurations = @()
        $totalPackages = $Packages.Count
        
        Write-Host "🔄 Retrieving app configurations for $totalPackages packages..." -ForegroundColor Cyan
        
        for ($i = 0; $i -lt $totalPackages; $i++) {
            $package = $Packages[$i]
            $progress = [math]::Round(($i + 1) / $totalPackages * 100)
            
            Write-Host "   Processing package $($i + 1) of $totalPackages ($progress%): $($package.productDisplayName)" -ForegroundColor Yellow
            
            $appConfig = Get-AppConfigurationByPackageId -PackageId $package.id
            if ($appConfig) {
                # Add the original package ID for reference
                $appConfig | Add-Member -MemberType NoteProperty -Name "mobileAppCatalogPackageId" -Value $package.id -Force
                $appConfigurations += $appConfig
            }
        }
        
        Write-Host "✅ Retrieved $($appConfigurations.Count) app configurations" -ForegroundColor Green
        return $appConfigurations
    }
    catch {
        Write-Host "❌ Error getting app configurations: $_" -ForegroundColor Red
        throw
    }
}

# Function to export app configurations to JSON
function Export-AppConfigurationsToJson {
    param (
        [Parameter(Mandatory = $true)]
        $AppConfigurations,
        
        [Parameter(Mandatory = $false)]
        [string]$SearchTerm,
        
        [Parameter(Mandatory = $false)]
        [string]$OutputFile
    )
    
    try {
        # Create output directory if it doesn't exist
        $outputDir = Join-Path -Path (Get-Location) -ChildPath "output"
        if (-not (Test-Path -Path $outputDir)) {
            New-Item -Path $outputDir -ItemType Directory | Out-Null
            Write-Host "📁 Created output directory: $outputDir" -ForegroundColor Gray
        }
        
        # Generate timestamp for filename
        $timestamp = Get-Date -Format "yyyyMMdd_HHmmss"
        
        # Create filename based on search term
        if ([string]::IsNullOrEmpty($OutputFile)) {
            $searchTermClean = if ($SearchTerm) { 
                ($SearchTerm -replace '[\\\/\:\*\?\"\<\>\|]', '_') 
            } else { 
                "AllAppConfigurations" 
            }
            $fileName = "EnterpriseAppCatalogConfigurations_${searchTermClean}_${timestamp}.json"
            $filePath = Join-Path -Path $outputDir -ChildPath $fileName
        } else {
            $filePath = $OutputFile
        }
        
        # Export to JSON
        $AppConfigurations | ConvertTo-Json -Depth 20 | Out-File -FilePath $filePath -Encoding UTF8
        
        Write-Host "💾 Exported app configurations to: $filePath" -ForegroundColor Green
        return $filePath
    }
    catch {
        Write-Host "❌ Error exporting app configurations to JSON: $_" -ForegroundColor Red
        return $null
    }
}

# Function to display app configuration details
function Show-AppConfigurationDetails {
    param (
        [Parameter(Mandatory = $true)]
        $AppConfiguration
    )
    
    Write-Host "📱 Enterprise App Catalog App Configuration Details:" -ForegroundColor Cyan
    Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
    
    # Display all properties dynamically
    $properties = $AppConfiguration | Get-Member -MemberType NoteProperty | Select-Object -ExpandProperty Name
    
    foreach ($prop in $properties) {
        $value = $AppConfiguration.$prop
        
        # Handle different property types
        if ($null -eq $value) {
            Write-Host "   • ${prop}: null" -ForegroundColor Yellow
        }
        elseif ($value -is [System.Collections.ICollection]) {
            if ($value.Count -gt 0) {
                Write-Host "   • ${prop}: [$($value.Count) items]" -ForegroundColor Green
                foreach ($item in $value) {
                    if ($item -is [PSCustomObject]) {
                        Write-Host "     - Object:" -ForegroundColor Yellow
                        $itemProps = $item | Get-Member -MemberType NoteProperty | Select-Object -ExpandProperty Name
                        foreach ($itemProp in $itemProps) {
                            Write-Host "       · ${itemProp}: $($item.$itemProp)" -ForegroundColor White
                        }
                    }
                    else {
                        Write-Host "     - $item" -ForegroundColor White
                    }
                }
            }
            else {
                Write-Host "   • ${prop}: [empty collection]" -ForegroundColor Yellow
            }
        }
        elseif ($value -is [PSCustomObject]) {
            Write-Host "   • ${prop}:" -ForegroundColor Green
            $nestedProps = $value | Get-Member -MemberType NoteProperty | Select-Object -ExpandProperty Name
            foreach ($nestedProp in $nestedProps) {
                Write-Host "     - ${nestedProp}: $($value.$nestedProp)" -ForegroundColor White
            }
        }
        else {
            # Handle simple properties
            Write-Host "   • ${prop}: $value" -ForegroundColor Green
        }
    }
    
    Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
    Write-Host ""
}

# Main script execution
try {
    Write-Host "📱 Enterprise App Catalog App Configuration Explorer" -ForegroundColor Cyan
    Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
    Write-Host ""
    
    # Connect to Microsoft Graph
    Connect-MicrosoftGraph -TenantId $TenantId -ClientId $ClientId -ClientSecret $ClientSecret
    
    $appConfigurations = @()
    
    # If a specific package ID is provided, get that configuration
    if ($PackageId) {
        Write-Host "🔍 Retrieving app configuration for package ID: $PackageId" -ForegroundColor Cyan
        $appConfig = Get-AppConfigurationByPackageId -PackageId $PackageId
        if ($appConfig) {
            $appConfigurations += $appConfig
        }
    }
    # Otherwise, get packages and their configurations
    else {
        if ($SearchTerm) {
            Write-Host "🔍 Searching for packages with display name containing: $SearchTerm" -ForegroundColor Cyan
        }
        else {
            Write-Host "🔍 Retrieving all enterprise app catalog packages (limited to $MaxResults)" -ForegroundColor Cyan
        }
        
        $packages = Get-EnterpriseAppCatalogPackages -SearchTerm $SearchTerm -MaxResults $MaxResults
        
        if ($packages.Count -eq 0) {
            Write-Host "⚠️ No packages found matching the criteria" -ForegroundColor Yellow
            exit 0
        }
        
        $appConfigurations = Get-AppConfigurations -Packages $packages
    }
    
    if ($appConfigurations.Count -eq 0) {
        Write-Host "⚠️ No app configurations found matching the criteria" -ForegroundColor Yellow
        exit 0
    }
    
    # Export to JSON if requested
    if ($ExportToJson) {
        $jsonPath = Export-AppConfigurationsToJson -AppConfigurations $appConfigurations -SearchTerm $SearchTerm -OutputFile $OutputFile
    }
    
    # Display app configuration details
    Write-Host "📊 Found $($appConfigurations.Count) app configuration(s)" -ForegroundColor Green
    Write-Host ""
    
    $displayLimit = [Math]::Min($appConfigurations.Count, 5)
    Write-Host "📋 Displaying details for first $displayLimit app configurations:" -ForegroundColor Cyan
    Write-Host ""
    
    for ($i = 0; $i -lt $displayLimit; $i++) {
        Write-Host "App Configuration $($i + 1) of ${displayLimit}:" -ForegroundColor Magenta
        Show-AppConfigurationDetails -AppConfiguration $appConfigurations[$i]
    }
    
    if ($appConfigurations.Count -gt $displayLimit) {
        Write-Host "ℹ️ $($appConfigurations.Count - $displayLimit) more app configurations were retrieved but not displayed." -ForegroundColor Yellow
        Write-Host "   Check the exported JSON file for complete data." -ForegroundColor Yellow
    }
    
    Write-Host "🎉 Operation completed successfully!" -ForegroundColor Green
}
catch {
    Write-Host "❌ Script execution failed: $_" -ForegroundColor Red
    exit 1
}
finally {
    # Disconnect from Microsoft Graph
    Write-Host "`nDisconnecting from Microsoft Graph..." -ForegroundColor Cyan
    Disconnect-MgGraph > $null 2>&1
    Write-Host "✅ Disconnected from Microsoft Graph" -ForegroundColor Green
} 