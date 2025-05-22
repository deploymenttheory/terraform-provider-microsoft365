param (
    [Parameter(Mandatory=$true)]
    [string]$TenantId,
    
    [Parameter(Mandatory=$true)]
    [string]$ClientId,
    
    [Parameter(Mandatory=$true)]
    [string]$ClientSecret,
    
    [Parameter(Mandatory=$false)]
    [string]$OutputFile = ".\IntuneEndpointData.json",
    
    [Parameter(Mandatory=$false)]
    [string]$AppId
)

# Collection of Intune Graph API GET functions

# Function to get all mobile apps
function Get-IntuneMobileApps {
    try {
        $uri = "https://graph.microsoft.com/beta/deviceAppManagement/mobileApps"
        $response = Invoke-MgGraphRequest -Method GET -Uri $uri
        
        return $response
    }
    catch {
        Write-Error "Error getting mobile apps: $_"
        throw
    }
}

# Function to get specific mobile app details
function Get-IntuneMobileApp {
    param (
        [Parameter(Mandatory=$true)]
        [string]$AppId
    )
    
    try {
        $uri = "https://graph.microsoft.com/beta/deviceAppManagement/mobileApps/$AppId"
        $response = Invoke-MgGraphRequest -Method GET -Uri $uri
        
        return $response
    }
    catch {
        Write-Error "Error getting mobile app details: $_"
        throw
    }
}

# Function to get content versions for an app
function Get-AppContentVersions {
    param (
        [Parameter(Mandatory=$true)]
        [string]$AppId,
        
        [Parameter(Mandatory=$true)]
        [string]$AppType
    )
    
    try {
        $uri = "https://graph.microsoft.com/beta/deviceAppManagement/mobileApps/$AppId/microsoft.graph.$AppType/contentVersions"
        $response = Invoke-MgGraphRequest -Method GET -Uri $uri
        
        return $response
    }
    catch {
        Write-Error "Error getting content versions: $_"
        throw
    }
}

# Function to get files for a content version
function Get-AppContentFiles {
    param (
        [Parameter(Mandatory=$true)]
        [string]$AppId,
        
        [Parameter(Mandatory=$true)]
        [string]$AppType,
        
        [Parameter(Mandatory=$true)]
        [string]$ContentVersionId
    )
    
    try {
        $uri = "https://graph.microsoft.com/beta/deviceAppManagement/mobileApps/$AppId/microsoft.graph.$AppType/contentVersions/$ContentVersionId/files"
        $response = Invoke-MgGraphRequest -Method GET -Uri $uri
        
        return $response
    }
    catch {
        Write-Error "Error getting content files: $_"
        throw
    }
}

# Function to get content file status
function Get-AppContentFileStatus {
    param (
        [Parameter(Mandatory=$true)]
        [string]$AppId,
        
        [Parameter(Mandatory=$true)]
        [string]$AppType,
        
        [Parameter(Mandatory=$true)]
        [string]$ContentVersionId,
        
        [Parameter(Mandatory=$true)]
        [string]$ContentFileId
    )
    
    try {
        $uri = "https://graph.microsoft.com/beta/deviceAppManagement/mobileApps/$AppId/microsoft.graph.$AppType/contentVersions/$ContentVersionId/files/$ContentFileId"
        $response = Invoke-MgGraphRequest -Method GET -Uri $uri
        return $response
    }
    catch {
        Write-Host "  ❌ Error retrieving content file status for file ${ContentFileId}: $_" -ForegroundColor Red
        $allApiData.mobileapp["fileStatus_${ContentVersionId}_${ContentFileId}"] = @{
            uri = $uri
            error = $_.Exception.Message
        }
        return $null
    }
}


# Function to get app assignments
function Get-AppAssignments {
    param (
        [Parameter(Mandatory=$true)]
        [string]$AppId
    )
    
    try {
        $uri = "https://graph.microsoft.com/beta/deviceAppManagement/mobileApps/$AppId/assignments"
        $response = Invoke-MgGraphRequest -Method GET -Uri $uri
        
        return $response
    }
    catch {
        Write-Error "Error getting app assignments: $_"
        throw
    }
}

# Function to get app categories
function Get-AppCategories {
    param (
        [Parameter(Mandatory=$true)]
        [string]$AppId
    )
    
    try {
        $uri = "https://graph.microsoft.com/beta/deviceAppManagement/mobileApps/$AppId/categories"
        $response = Invoke-MgGraphRequest -Method GET -Uri $uri
        
        return $response
    }
    catch {
        Write-Error "Error getting app categories: $_"
        throw
    }
}

# Function to get all mobile app categories
function Get-IntuneMobileAppCategories {
    try {
        $uri = "https://graph.microsoft.com/beta/deviceAppManagement/mobileAppCategories"
        $response = Invoke-MgGraphRequest -Method GET -Uri $uri
        
        return $response
    }
    catch {
        Write-Error "Error getting mobile app categories: $_"
        throw
    }
}

# Function to get app assignments (all)
function Get-IntuneAppAssignments {
    try {
        $uri = "https://graph.microsoft.com/beta/deviceAppManagement/mobileApps/assignments"
        $response = Invoke-MgGraphRequest -Method GET -Uri $uri
        
        return $response
    }
    catch {
        Write-Error "Error getting app assignments: $_"
        throw
    }
}

# Function to get specific app assignment details
function Get-AppAssignmentById {
    param (
        [Parameter(Mandatory=$true)]
        [string]$AppId,
        
        [Parameter(Mandatory=$true)]
        [string]$AssignmentId
    )
    
    try {
        $uri = "https://graph.microsoft.com/beta/deviceAppManagement/mobileApps/$AppId/assignments/$AssignmentId"
        $response = Invoke-MgGraphRequest -Method GET -Uri $uri
        
        return $response
    }
    catch {
        Write-Error "Error getting assignment details: $_"
        throw
    }
}


# Function to get managed device overview
function Get-IntuneManagedDeviceOverview {
    try {
        $uri = "https://graph.microsoft.com/beta/deviceManagement/managedDeviceOverview"
        $response = Invoke-MgGraphRequest -Method GET -Uri $uri
        
        return $response
    }
    catch {
        Write-Error "Error getting managed device overview: $_"
        throw
    }
}

# Helper function to generate reference ID from path
function Get-PathBasedReferenceId {
    param (
        [string]$Path
    )
    
    $segments = $Path -split '/' | Where-Object { $_ -ne '' }
    $refParts = @()
    
    foreach ($segment in $segments) {
        $segmentStr = [string]$segment
        if ($segmentStr -match '{.*}') {
            $paramName = $segmentStr -replace '{|}'
            $paramName = $paramName -replace '-', '_'
            $refParts += "BY_$([string]($paramName.ToUpper()))"
        } else {
            $refParts += [string]($segmentStr.ToUpper())
        }
    }
    
    return $refParts -join '_'
}

# Import required module
Import-Module Microsoft.Graph.Authentication

# Connect to Microsoft Graph
Write-Host "Connecting to Microsoft Graph..." -ForegroundColor Cyan
$secureClientSecret = ConvertTo-SecureString -String $ClientSecret -AsPlainText -Force
$clientSecretCredential = New-Object -TypeName System.Management.Automation.PSCredential -ArgumentList $ClientId, $secureClientSecret
Connect-MgGraph -ClientSecretCredential $clientSecretCredential -TenantId $TenantId

# Initialize collection for all API data
$allApiData = @{
    timestamp = Get-Date -Format "yyyy-MM-dd HH:mm:ss"
    tenant_id = $TenantId
    mobileapp = @{}
}

# Helper function to add API response to the collection
function Add-ApiResponse {
    param (
        [string]$EndpointName,
        [string]$Uri,
        [PSObject]$Response
    )
    
    $refId = Get-PathBasedReferenceId -Path ($Uri -replace "https://graph.microsoft.com/beta", "")
    
    $allApiData.mobileapp[$EndpointName] = @{
        uri = $Uri
        reference_id = $refId
        data = $Response
    }
}

# 1. Get all mobile apps or specific app data based on AppId
Write-Host "`n📱 Retrieving mobile apps..." -ForegroundColor Cyan

# Check if AppId is provided
if (-not [string]::IsNullOrEmpty($AppId)) {
    Write-Host "`n🔍 Retrieving specific app data for ID: $AppId" -ForegroundColor Cyan
    
    # 1a. Get app details
    try {
        $appDetails = Get-IntuneMobileApp -AppId $AppId
        Add-ApiResponse -EndpointName "appDetails" -Uri "https://graph.microsoft.com/beta/deviceAppManagement/mobileApps/$AppId" -Response $appDetails
        Write-Host "  ✅ Successfully retrieved app details" -ForegroundColor Green
        
        # Get app type (if available)
        $appType = "Unknown"
        if ($appDetails -and $appDetails.'@odata.type') {
            $appType = $appDetails.'@odata.type' -replace "#microsoft.graph.", ""
        }
        
        # 1b. Get content versions
        try {
            $contentVersions = Get-AppContentVersions -AppId $AppId -AppType $appType
            Add-ApiResponse -EndpointName "contentVersions" -Uri "https://graph.microsoft.com/beta/deviceAppManagement/mobileApps/$AppId/microsoft.graph.$appType/contentVersions" -Response $contentVersions
            Write-Host "  ✅ Successfully retrieved content versions" -ForegroundColor Green
            
            # For each content version, get files
            if ($contentVersions -and $contentVersions.value) {
                foreach ($version in $contentVersions.value) {
                    $versionId = $version.id
                    
                    # 1c. Get files for this content version
                    try {
                        $contentFiles = Get-AppContentFiles -AppId $AppId -AppType $appType -ContentVersionId $versionId
                        Add-ApiResponse -EndpointName "contentFiles_$versionId" -Uri "https://graph.microsoft.com/beta/deviceAppManagement/mobileApps/$AppId/microsoft.graph.$appType/contentVersions/$versionId/files" -Response $contentFiles
                        Write-Host "  ✅ Successfully retrieved content files for version $versionId" -ForegroundColor Green
                        
                        # For each file, get status
                        if ($contentFiles -and $contentFiles.value) {
                            foreach ($file in $contentFiles.value) {
                                $fileId = $file.id
                                
                                # 1d. Get file status
                                try {
                                    $fileStatus = Get-AppContentFileStatus -AppId $AppId -AppType $appType -ContentVersionId $versionId -ContentFileId $fileId
                                    Add-ApiResponse -EndpointName "fileStatus_${versionId}_$fileId" -Uri "https://graph.microsoft.com/beta/deviceAppManagement/mobileApps/$AppId/microsoft.graph.$appType/contentVersions/$versionId/files/$fileId" -Response $fileStatus
                                    Write-Host "  ✅ Successfully retrieved file status for file $fileId" -ForegroundColor Green
                                }
                                catch {
                                    Write-Host "  ❌ Error retrieving file status: $_" -ForegroundColor Red
                                    $allApiData.mobileapp["fileStatus_${versionId}_$fileId"] = @{
                                        uri = "https://graph.microsoft.com/beta/deviceAppManagement/mobileApps/$AppId/microsoft.graph.$appType/contentVersions/$versionId/files/$fileId"
                                        error = $_.Exception.Message
                                    }
                                }
                            }
                        }
                    }
                    catch {
                        Write-Host "  ❌ Error retrieving content files: $_" -ForegroundColor Red
                        $allApiData.mobileapp["contentFiles_$versionId"] = @{
                            uri = "https://graph.microsoft.com/beta/deviceAppManagement/mobileApps/$AppId/microsoft.graph.$appType/contentVersions/$versionId/files"
                            error = $_.Exception.Message
                        }
                    }
                }
            }
        }
        catch {
            Write-Host "  ❌ Error retrieving content versions: $_" -ForegroundColor Red
            $allApiData.mobileapp["contentVersions"] = @{
                uri = "https://graph.microsoft.com/beta/deviceAppManagement/mobileApps/$AppId/microsoft.graph.$appType/contentVersions"
                error = $_.Exception.Message
            }
        }
        
        # 1e. Get assignments
        try {
            $assignments = Get-AppAssignments -AppId $AppId
            Add-ApiResponse -EndpointName "assignments" -Uri "https://graph.microsoft.com/beta/deviceAppManagement/mobileApps/$AppId/assignments" -Response $assignments
            Write-Host "  ✅ Successfully retrieved app assignments" -ForegroundColor Green
        }
        catch {
            Write-Host "  ❌ Error retrieving app assignments: $_" -ForegroundColor Red
            $allApiData.mobileapp["assignments"] = @{
                uri = "https://graph.microsoft.com/beta/deviceAppManagement/mobileApps/$AppId/assignments"
                error = $_.Exception.Message
            }
        }
        
        # 1f. Get categories
        try {
            $categories = Get-AppCategories -AppId $AppId
            Add-ApiResponse -EndpointName "categories" -Uri "https://graph.microsoft.com/beta/deviceAppManagement/mobileApps/$AppId/categories" -Response $categories
            Write-Host "  ✅ Successfully retrieved app categories" -ForegroundColor Green
        }
        catch {
            Write-Host "  ❌ Error retrieving app categories: $_" -ForegroundColor Red
            $allApiData.mobileapp["categories"] = @{
                uri = "https://graph.microsoft.com/beta/deviceAppManagement/mobileApps/$AppId/categories"
                error = $_.Exception.Message
            }
        }
    }
    catch {
        Write-Host "  ❌ Error retrieving app details: $_" -ForegroundColor Red
        $allApiData.mobileapp["appDetails"] = @{
            uri = "https://graph.microsoft.com/beta/deviceAppManagement/mobileApps/$AppId"
            error = $_.Exception.Message
        }
    }
}
else {
    # 2. If no AppId, get all mobile apps
    try {
        $mobileApps = Get-IntuneMobileApps
        Add-ApiResponse -EndpointName "mobileApps" -Uri "https://graph.microsoft.com/beta/deviceAppManagement/mobileApps" -Response $mobileApps
        Write-Host "  ✅ Successfully retrieved mobile apps" -ForegroundColor Green
    }
    catch {
        Write-Host "  ❌ Error retrieving mobile apps: $_" -ForegroundColor Red
        $allApiData.mobileapp["mobileApps"] = @{
            uri = "https://graph.microsoft.com/beta/deviceAppManagement/mobileApps"
            error = $_.Exception.Message
        }
    }
}

# Export to JSON file
Write-Host "`n💾 Exporting data to JSON file..." -ForegroundColor Cyan
$jsonData = $allApiData | ConvertTo-Json -Depth 20
$jsonData | Out-File -FilePath $OutputFile -Force

Write-Host "✅ Successfully exported data to: $OutputFile" -ForegroundColor Green
Write-Host "   File Size: $([math]::Round((Get-Item $OutputFile).Length / 1KB, 2)) KB" -ForegroundColor Cyan

# Disconnect from Microsoft Graph
Write-Host "`nDisconnecting from Microsoft Graph..." -ForegroundColor Cyan
Disconnect-MgGraph > $null 2>&1
Write-Host "✅ Disconnected from Microsoft Graph" -ForegroundColor Green
