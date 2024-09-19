[CmdletBinding()]
param (
    [Parameter(Mandatory=$true)]
    [string]$ClientId,
    
    [Parameter(Mandatory=$true)]
    [string]$ClientSecret,
    
    [Parameter(Mandatory=$true)]
    [string]$TenantId
)

# Import the Microsoft.Graph.Authentication module
Import-Module Microsoft.Graph.Authentication

# Connect to Microsoft Graph using client credentials
$secureClientSecret = ConvertTo-SecureString -String $ClientSecret -AsPlainText -Force
$clientSecretCredential = New-Object -TypeName System.Management.Automation.PSCredential -ArgumentList $ClientId, $secureClientSecret

Connect-MgGraph -ClientSecretCredential $clientSecretCredential -TenantId $TenantId

# Determine the script's location
$scriptPath = $PSScriptRoot
if (-not $scriptPath) {
    $scriptPath = Split-Path -Parent -Path $MyInvocation.MyCommand.Definition
}
if (-not $scriptPath) {
    $scriptPath = Get-Location
}

# Create the Export folder if it doesn't exist
$exportFolder = Join-Path -Path $scriptPath -ChildPath "Export"
if (-not (Test-Path -Path $exportFolder)) {
    New-Item -ItemType Directory -Path $exportFolder | Out-Null
}

# Define the output file path
$outputFile = Join-Path -Path $exportFolder -ChildPath "ExportedGraphPermissions.json"

# Get all permissions
$allPermissions = Find-MgGraphPermission

# Get all commands for both v1.0 and beta
$allCommands = @()
$allCommands += Find-MgGraphCommand -Command "*" -ApiVersion "v1.0"
$allCommands += Find-MgGraphCommand -Command "*" -ApiVersion "beta"

# Initialize the result hashtable
$permissionDetails = [ordered]@{}

# Process each permission
foreach ($permission in $allPermissions) {
    $permissionName = $permission.Name
    $permissionDetails[$permissionName] = [ordered]@{
        Description = $permission.Description
        Id = $permission.Id
        Consent = $permission.Consent
        Type = $permission.PermissionType
        Scope = ($permissionName -split '\.')[0]
        IsAdmin = $permission.IsAdmin
        AssociatedURIs = @()
        PowerShellCmdlets = @()
    }

    # Filter commands associated with this permission
    $associatedCommands = $allCommands | Where-Object { $_.Permissions.Name -contains $permissionName }

    foreach ($command in $associatedCommands) {
        $uri = $command.URI
        $method = $command.Method
        $apiVersions = @($command.ApiVersion)

        # Check if this URI already exists in AssociatedURIs
        $existingUri = $permissionDetails[$permissionName].AssociatedURIs | Where-Object { $_.URI -eq $uri -and $_.Method -eq $method }

        if ($existingUri) {
            # Update existing URI entry
            $existingUri.ApiVersions += $apiVersions | Where-Object { $_ -notin $existingUri.ApiVersions }
        } else {
            # Add new URI entry
            $permissionDetails[$permissionName].AssociatedURIs += @{
                URI = $uri
                Method = $method
                IsLeastPrivileged = $command.Permissions.Count -eq 1
                ApiVersions = $apiVersions
            }
        }

        # Add PowerShell cmdlet if not already present
        if ($command.Command -notin $permissionDetails[$permissionName].PowerShellCmdlets) {
            $permissionDetails[$permissionName].PowerShellCmdlets += $command.Command
        }
    }
}

# Sort the permissions alphabetically
$sortedPermissionDetails = [ordered]@{}
$permissionDetails.GetEnumerator() | Sort-Object Name | ForEach-Object {
    $sortedPermissionDetails[$_.Name] = $_.Value
}

# Convert to JSON and export
$sortedPermissionDetails | ConvertTo-Json -Depth 6 | Out-File -FilePath $outputFile

Write-Host "Detailed Graph permissions data has been exported to $outputFile"

# Disconnect from Microsoft Graph
Disconnect-MgGraph