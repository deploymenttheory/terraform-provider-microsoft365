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
    
    [Parameter(Mandatory = $false,
        HelpMessage = "Search term to find packages by display name")]
    [string]$SearchTerm,
    
    [Parameter(Mandatory = $false,
        HelpMessage = "Maximum number of results to return")]
    [int]$MaxResults = 50,
    
    [Parameter(Mandatory = $false,
        HelpMessage = "Export results to JSON file")]
    [bool]$ExportToJson = $true,
    
    [Parameter(Mandatory = $false,
        HelpMessage = "Use `$apply parameter to group results")]
    [bool]$UseGroupBy = $false,
    
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

# Function to get enterprise app catalog applications
function Get-EnterpriseAppCatalogPackages {
    param (
        [Parameter(Mandatory = $false)]
        [string]$SearchTerm,
        
        [Parameter(Mandatory = $false)]
        [int]$MaxResults = 50,
        
        [Parameter(Mandatory = $false)]
        [bool]$UseGroupBy = $false
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
        
        # Handle groupBy if requested
        if ($UseGroupBy) {
            $queryParams += "`$apply=groupby((productId,productDisplayName,publisherDisplayName))"
        }
        
        # Build the final URL
        $url = $baseUrl
        if ($queryParams.Count -gt 0) {
            $url += "?" + ($queryParams -join "&")
        }
        
        Write-Host "🔍 Fetching enterprise app catalog applications..." -ForegroundColor Cyan
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
        Write-Host "❌ Error getting enterprise app catalog applications: $_" -ForegroundColor Red
        throw
    }
}

# Function to export packages to JSON
function Export-PackagesToJson {
    param (
        [Parameter(Mandatory = $true)]
        $Packages,
        
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
                "AllPackages" 
            }
            $fileName = "EnterpriseAppCatalogPackages_${searchTermClean}_${timestamp}.json"
            $filePath = Join-Path -Path $outputDir -ChildPath $fileName
        } else {
            $filePath = $OutputFile
        }
        
        # Export to JSON
        $Packages | ConvertTo-Json -Depth 20 | Out-File -FilePath $filePath -Encoding UTF8
        
        Write-Host "💾 Exported packages to: $filePath" -ForegroundColor Green
        return $filePath
    }
    catch {
        Write-Host "❌ Error exporting packages to JSON: $_" -ForegroundColor Red
        return $null
    }
}

# Function to display package details
function Show-PackageDetails {
    param (
        [Parameter(Mandatory = $true)]
        $Package
    )
    
    Write-Host "📦 Enterprise App Catalog Package Details:" -ForegroundColor Cyan
    Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
    
    # Display all properties dynamically
    $properties = $Package | Get-Member -MemberType NoteProperty | Select-Object -ExpandProperty Name
    
    foreach ($prop in $properties) {
        $value = $Package.$prop
        
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
    Write-Host "📱 Enterprise App Catalog Package Explorer" -ForegroundColor Cyan
    Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
    Write-Host ""
    
    # Connect to Microsoft Graph
    Connect-MicrosoftGraph -TenantId $TenantId -ClientId $ClientId -ClientSecret $ClientSecret
    
    # Get enterprise app catalog applications
    if ($SearchTerm) {
        Write-Host "🔍 Searching for packages with display name containing: $SearchTerm" -ForegroundColor Cyan
    }
    else {
        Write-Host "🔍 Retrieving all enterprise app catalog applications (limited to $MaxResults)" -ForegroundColor Cyan
    }
    
    $packages = Get-EnterpriseAppCatalogPackages -SearchTerm $SearchTerm -MaxResults $MaxResults -UseGroupBy $UseGroupBy
    
    if ($packages.Count -eq 0) {
        Write-Host "⚠️ No packages found matching the criteria" -ForegroundColor Yellow
        exit 0
    }
    
    # Export to JSON if requested
    if ($ExportToJson) {
        $jsonPath = Export-PackagesToJson -Packages $packages -SearchTerm $SearchTerm -OutputFile $OutputFile
    }
    
    # Display package details
    Write-Host "📊 Found $($packages.Count) package(s)" -ForegroundColor Green
    Write-Host ""
    
    $displayLimit = [Math]::Min($packages.Count, 10)
    Write-Host "📋 Displaying details for first $displayLimit packages:" -ForegroundColor Cyan
    Write-Host ""
    
    for ($i = 0; $i -lt $displayLimit; $i++) {
        Write-Host "Package $($i + 1) of ${displayLimit}:" -ForegroundColor Magenta
        Show-PackageDetails -Package $packages[$i]
    }
    
    if ($packages.Count -gt $displayLimit) {
        Write-Host "ℹ️ $($packages.Count - $displayLimit) more packages were retrieved but not displayed." -ForegroundColor Yellow
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