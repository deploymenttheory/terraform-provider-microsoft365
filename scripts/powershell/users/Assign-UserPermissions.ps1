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
    
    [Parameter(Mandatory=$true,
    HelpMessage="User ID (UPN or Object ID)")]
    [ValidateNotNullOrEmpty()]
    [string]$UserId,
    
    [Parameter(Mandatory=$true,
    HelpMessage="Assignment type: AppRole or DirectoryRole")]
    [ValidateSet("AppRole", "DirectoryRole")]
    [string]$AssignmentType,
    
    [Parameter(Mandatory=$false,
    HelpMessage="Service Principal ID (required for AppRole assignments)")]
    [string]$ResourceId,
    
    [Parameter(Mandatory=$false,
    HelpMessage="App Role ID to assign (required for AppRole assignments)")]
    [string]$AppRoleId,
    
    [Parameter(Mandatory=$false,
    HelpMessage="Directory Role ID (required for DirectoryRole assignments)")]
    [string]$DirectoryRoleId
)

# Import required modules
Import-Module Microsoft.Graph.Authentication

# Function to assign app role to user
function New-UserAppRoleAssignment {
    param (
        [Parameter(Mandatory=$true)]
        [string]$UserObjectId,
        
        [Parameter(Mandatory=$true)]
        [string]$ServicePrincipalId,
        
        [Parameter(Mandatory=$true)]
        [string]$RoleId
    )
    
    try {
        $uri = "https://graph.microsoft.com/v1.0/users/$UserObjectId/appRoleAssignments"
        
        Write-Host "🔨 Creating app role assignment..." -ForegroundColor Cyan
        Write-Host "   User ID: $UserObjectId" -ForegroundColor Gray
        Write-Host "   Resource ID (Service Principal): $ServicePrincipalId" -ForegroundColor Gray
        Write-Host "   App Role ID: $RoleId" -ForegroundColor Gray
        Write-Host "   Endpoint: $uri" -ForegroundColor Gray
        Write-Host ""
        
        $body = @{
            principalId = $UserObjectId
            resourceId = $ServicePrincipalId
            appRoleId = $RoleId
        }
        
        $response = Invoke-MgGraphRequest -Method POST -Uri $uri -Body $body
        
        Write-Host "✅ App role assignment created successfully!" -ForegroundColor Green
        Write-Host ""
        
        return $response
    }
    catch {
        Write-Host "❌ Error creating app role assignment: $_" -ForegroundColor Red
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

# Function to assign user to directory role
function Add-UserToDirectoryRole {
    param (
        [Parameter(Mandatory=$true)]
        [string]$UserObjectId,
        
        [Parameter(Mandatory=$true)]
        [string]$RoleId
    )
    
    try {
        $uri = "https://graph.microsoft.com/v1.0/directoryRoles/$RoleId/members/`$ref"
        
        Write-Host "🔨 Adding user to directory role..." -ForegroundColor Cyan
        Write-Host "   User ID: $UserObjectId" -ForegroundColor Gray
        Write-Host "   Directory Role ID: $RoleId" -ForegroundColor Gray
        Write-Host "   Endpoint: $uri" -ForegroundColor Gray
        Write-Host ""
        
        $body = @{
            "@odata.id" = "https://graph.microsoft.com/v1.0/users/$UserObjectId"
        }
        
        $response = Invoke-MgGraphRequest -Method POST -Uri $uri -Body $body
        
        Write-Host "✅ User added to directory role successfully!" -ForegroundColor Green
        Write-Host ""
        
        return $response
    }
    catch {
        Write-Host "❌ Error adding user to directory role: $_" -ForegroundColor Red
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

# Function to display app role assignment details
function Show-AppRoleAssignmentDetails {
    param (
        [Parameter(Mandatory=$true)]
        $Assignment
    )
    
    Write-Host "📋 App Role Assignment Details:" -ForegroundColor Cyan
    Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
    
    if ($Assignment.id) {
        Write-Host "   • Assignment ID: $($Assignment.id)" -ForegroundColor Green
    }
    
    if ($Assignment.principalId) {
        Write-Host "   • Principal ID (User): $($Assignment.principalId)" -ForegroundColor Green
    }
    
    if ($Assignment.principalDisplayName) {
        Write-Host "   • Principal Display Name: $($Assignment.principalDisplayName)" -ForegroundColor Green
    }
    
    if ($Assignment.resourceId) {
        Write-Host "   • Resource ID (Service Principal): $($Assignment.resourceId)" -ForegroundColor Green
    }
    
    if ($Assignment.resourceDisplayName) {
        Write-Host "   • Resource Display Name: $($Assignment.resourceDisplayName)" -ForegroundColor Green
    }
    
    if ($Assignment.appRoleId) {
        Write-Host "   • App Role ID: $($Assignment.appRoleId)" -ForegroundColor Green
    }
    
    if ($Assignment.principalType) {
        Write-Host "   • Principal Type: $($Assignment.principalType)" -ForegroundColor Green
    }
    
    if ($Assignment.createdDateTime) {
        Write-Host "   • Created: $($Assignment.createdDateTime)" -ForegroundColor Green
    }
    
    Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
    Write-Host ""
}

# Validate parameters based on assignment type
if ($AssignmentType -eq "AppRole") {
    if (-not $ResourceId -or -not $AppRoleId) {
        Write-Host "❌ Error: ResourceId and AppRoleId are required for AppRole assignments" -ForegroundColor Red
        exit 1
    }
}

if ($AssignmentType -eq "DirectoryRole") {
    if (-not $DirectoryRoleId) {
        Write-Host "❌ Error: DirectoryRoleId is required for DirectoryRole assignments" -ForegroundColor Red
        exit 1
    }
}

# Script Setup
try {
    Write-Host "🔐 Connecting to Microsoft Graph..." -ForegroundColor Cyan
    Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
    
    $secureClientSecret = ConvertTo-SecureString -String $ClientSecret -AsPlainText -Force
    $clientSecretCredential = New-Object -TypeName System.Management.Automation.PSCredential -ArgumentList $ClientId, $secureClientSecret
    
    Connect-MgGraph -ClientSecretCredential $clientSecretCredential -TenantId $TenantId
    
    Write-Host "✅ Connected to Microsoft Graph successfully" -ForegroundColor Green
    Write-Host ""
    
    # Process based on assignment type
    switch ($AssignmentType) {
        "AppRole" {
            # Create app role assignment
            $assignment = New-UserAppRoleAssignment -UserObjectId $UserId -ServicePrincipalId $ResourceId -RoleId $AppRoleId
            
            if ($assignment) {
                Show-AppRoleAssignmentDetails -Assignment $assignment
            }
        }
        
        "DirectoryRole" {
            # Add user to directory role
            Add-UserToDirectoryRole -UserObjectId $UserId -RoleId $DirectoryRoleId
            
            Write-Host "📋 Directory Role Assignment Completed:" -ForegroundColor Cyan
            Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
            Write-Host "   • User ID: $UserId" -ForegroundColor Green
            Write-Host "   • Directory Role ID: $DirectoryRoleId" -ForegroundColor Green
            Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
            Write-Host ""
        }
    }
    
    Write-Host "🎉 Operation completed successfully!" -ForegroundColor Green
}
catch {
    Write-Host "❌ Script execution failed: $_" -ForegroundColor Red
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
