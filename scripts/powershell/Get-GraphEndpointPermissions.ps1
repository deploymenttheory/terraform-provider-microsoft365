<#
.SYNOPSIS
    Discovers required Microsoft Graph API permissions for a specific endpoint and optionally compares against an enterprise application.

.DESCRIPTION
    Uses Find-MgGraphCommand to identify the exact APPLICATION permissions needed for any Graph API endpoint.
    Can optionally compare these requirements against an existing enterprise application to identify gaps.
    
    The script provides two modes of operation:
    1. Discovery Mode: Shows what permissions are required for a given endpoint
    2. Evaluation Mode: Compares an app's current permissions against the requirements

.PARAMETER Uri
    The Graph API URI to check (e.g., "/groups", "/groups/{id}", "/users/{id}/memberOf")
    Use {id} as a placeholder for resource identifiers.

.PARAMETER Method
    The HTTP method (GET, POST, PATCH, DELETE). If not specified, shows all methods.

.PARAMETER ApiVersion
    The API version to check (v1.0 or beta). Default: beta

.PARAMETER AppIdToBeEvaluated
    Optional: The Application (Client) ID of the enterprise application to evaluate against the endpoint requirements.
    When provided, the script will compare the app's current permissions against what's required.

.PARAMETER TenantId
    Required when AppIdToBeEvaluated is provided. The Tenant ID for authentication.

.PARAMETER ClientSecret
    Required when AppIdToBeEvaluated is provided. The Client Secret for authentication.

.EXAMPLE
    # Scenario 1: Check permissions required for creating groups
    .\Get-GraphEndpointPermissions.ps1 -Uri "/groups" -Method POST
    
    Output: Shows Directory.ReadWrite.All and Group.ReadWrite.All are required

.EXAMPLE
    # Scenario 2: Check permissions for updating groups
    .\Get-GraphEndpointPermissions.ps1 -Uri "/groups/{id}" -Method PATCH
    
    Output: Shows permissions needed to modify group properties

.EXAMPLE
    # Scenario 3: Check all operations on a groups endpoint
    .\Get-GraphEndpointPermissions.ps1 -Uri "/groups/{id}"
    
    Output: Shows permissions for GET, PATCH, DELETE operations

.EXAMPLE
    # Scenario 4: Check permissions for user updates
    .\Get-GraphEndpointPermissions.ps1 -Uri "/users/{id}" -Method PATCH
    
    Output: Shows User.ReadWrite.All and related permissions

.EXAMPLE
    # Scenario 5: Check permissions for reading group members
    .\Get-GraphEndpointPermissions.ps1 -Uri "/groups/{id}/members" -Method GET
    
    Output: Shows GroupMember.Read.All or Group.Read.All

.EXAMPLE
    # Scenario 6: Evaluate app against group creation requirements
    .\Get-GraphEndpointPermissions.ps1 -Uri "/groups" -Method POST `
        -AppIdToBeEvaluated "00000000-0000-0000-0000-000000000001" `
        -TenantId "00000000-0000-0000-0000-000000000002" `
        -ClientSecret "xxx~xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
    
    Output: Shows which permissions are granted vs missing, plus directory role status

.EXAMPLE
    # Scenario 7: Evaluate app against user update requirements
    .\Get-GraphEndpointPermissions.ps1 -Uri "/users/{id}" -Method PATCH `
        -AppIdToBeEvaluated "00000000-0000-0000-0000-000000000003" `
        -TenantId "00000000-0000-0000-0000-000000000002" `
        -ClientSecret "xxx~xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
    
    Output: Compares required vs granted permissions for user management

.EXAMPLE
    # Scenario 8: Check v1.0 API version instead of beta
    .\Get-GraphEndpointPermissions.ps1 -Uri "/groups" -Method POST -ApiVersion "v1.0"
    
    Output: Shows permissions for the v1.0 endpoint

.EXAMPLE
    # Scenario 9: Check permissions for app role assignments
    .\Get-GraphEndpointPermissions.ps1 -Uri "/groups/{id}/appRoleAssignments" -Method POST
    
    Output: Shows AppRoleAssignment.ReadWrite.All permission

.EXAMPLE
    # Scenario 10: Check permissions for directory roles
    .\Get-GraphEndpointPermissions.ps1 -Uri "/directoryRoles/{id}/members" -Method POST
    
    Output: Shows RoleManagement.ReadWrite.Directory permission

.EXAMPLE
    # Scenario 11: Check Intune device management permissions
    .\Get-GraphEndpointPermissions.ps1 -Uri "/deviceManagement/managedDevices/{id}" -Method PATCH
    
    Output: Shows DeviceManagementManagedDevices.ReadWrite.All

.EXAMPLE
    # Scenario 12: Check service principal permissions
    .\Get-GraphEndpointPermissions.ps1 -Uri "/servicePrincipals/{id}" -Method PATCH
    
    Output: Shows Application.ReadWrite.All permission

.NOTES
    Author: Deployment Theory
    Reference: https://practical365.com/microsoft-graph-api-permission/
    Version: 1.0
    
    Common URI Patterns:
    - /groups                           - Collection operations
    - /groups/{id}                      - Single resource operations
    - /groups/{id}/members              - Nested collection
    - /users/{id}/memberOf              - Relationship navigation
    - /deviceManagement/configurationPolicies - Service-specific endpoints
    
    Permission Types:
    - Application: For service principals (daemon apps, background services)
    - Delegated: For user context (interactive apps, user sign-in)
    
    Special Requirements:
    - Some operations require directory roles in addition to API permissions
    - Role-assignable groups require Privileged Role Administrator role
    - The script will identify these special requirements automatically
    #>

[CmdletBinding()]
param (
    [Parameter(Mandatory = $true,
        HelpMessage = "The Graph API URI to check (e.g., '/groups', '/groups/{id}')")]
    [ValidateNotNullOrEmpty()]
    [string]$Uri,

    [Parameter(Mandatory = $false,
        HelpMessage = "The HTTP method (GET, POST, PATCH, DELETE)")]
    [ValidateSet("GET", "POST", "PATCH", "PUT", "DELETE", "")]
    [string]$Method = "",

    [Parameter(Mandatory = $false,
        HelpMessage = "The API version to check (v1.0 or beta)")]
    [ValidateSet("v1.0", "beta")]
    [string]$ApiVersion = "beta",

    [Parameter(Mandatory = $false,
        HelpMessage = "The Application (Client) ID to evaluate against the endpoint requirements")]
    [string]$AppIdToBeEvaluated,

    [Parameter(Mandatory = $false,
        HelpMessage = "The Tenant ID (required when AppIdToBeEvaluated is provided)")]
    [string]$TenantId,

    [Parameter(Mandatory = $false,
        HelpMessage = "The Client Secret (required when AppIdToBeEvaluated is provided)")]
    [string]$ClientSecret
)

#region Helper Functions

function Write-LogHeader {
    param (
        [Parameter(Mandatory = $true)]
        [string]$Message
    )
    
    Write-Host "`n========================================" -ForegroundColor Cyan
    Write-Host $Message -ForegroundColor Cyan
    Write-Host "========================================" -ForegroundColor Cyan
}

function Write-LogSection {
    param (
        [Parameter(Mandatory = $true)]
        [string]$Message
    )
    
    Write-Host "`n$Message" -ForegroundColor Yellow
    Write-Host ("-" * $Message.Length) -ForegroundColor Yellow
}

function Write-LogInfo {
    param (
        [Parameter(Mandatory = $true)]
        [string]$Message,
        
        [Parameter(Mandatory = $false)]
        [int]$Indent = 0
    )
    
    $indentString = "  " * $Indent
    Write-Host "$indentString$Message" -ForegroundColor White
}

function Write-LogSuccess {
    param (
        [Parameter(Mandatory = $true)]
        [string]$Message,
        
        [Parameter(Mandatory = $false)]
        [int]$Indent = 0
    )
    
    $indentString = "  " * $Indent
    Write-Host "$indentString✓ $Message" -ForegroundColor Green
}

function Write-LogWarning {
    param (
        [Parameter(Mandatory = $true)]
        [string]$Message,
        
        [Parameter(Mandatory = $false)]
        [int]$Indent = 0
    )
    
    $indentString = "  " * $Indent
    Write-Host "$indentString⚠ $Message" -ForegroundColor Yellow
}

function Write-LogError {
    param (
        [Parameter(Mandatory = $true)]
        [string]$Message,
        
        [Parameter(Mandatory = $false)]
        [int]$Indent = 0
    )
    
    $indentString = "  " * $Indent
    Write-Host "$indentString✗ $Message" -ForegroundColor Red
}

function Test-ModuleInstalled {
    param (
        [Parameter(Mandatory = $true)]
        [string]$ModuleName
    )
    
    return (Get-Module -ListAvailable -Name $ModuleName) -ne $null
}

function Install-RequiredModule {
    param (
        [Parameter(Mandatory = $true)]
        [string]$ModuleName
    )
    
    if (-not (Test-ModuleInstalled -ModuleName $ModuleName)) {
        Write-LogInfo "Installing $ModuleName..." -Indent 1
        try {
            Install-Module $ModuleName -Scope CurrentUser -Force -AllowClobber -ErrorAction Stop
            Write-LogSuccess "$ModuleName installed successfully" -Indent 1
        }
        catch {
            Write-LogError "Failed to install $ModuleName`: $($_.Exception.Message)" -Indent 1
            throw
        }
    }
    else {
        Write-LogSuccess "$ModuleName is already installed" -Indent 1
    }
}

function Get-PermissionName {
    param (
        [Parameter(Mandatory = $true)]
        $PermissionObject
    )
    
    # Extract permission name - handle both string and object types
    if ($PermissionObject -is [string]) {
        return $PermissionObject
    }
    elseif ($PermissionObject.Name) {
        return $PermissionObject.Name
    }
    elseif ($PermissionObject.ToString() -ne $PermissionObject.GetType().FullName) {
        return $PermissionObject.ToString()
    }
    
    return $null
}

function Test-ApplicationPermission {
    param (
        [Parameter(Mandatory = $true)]
        [string]$PermissionName
    )
    
    # Application permissions typically don't have "AccessAsUser" and end with ".All"
    return ($PermissionName -notlike "*AccessAsUser*" -and $PermissionName -like "*.All")
}

#endregion

#region Main Script

try {
    # Validate parameters
    if ($AppIdToBeEvaluated) {
        if (-not $TenantId -or -not $ClientSecret) {
            Write-LogError "When AppIdToBeEvaluated is provided, TenantId and ClientSecret are required"
            exit 1
        }
    }

    # Check and install required modules
    Write-LogHeader "Checking Required Modules"
    Install-RequiredModule -ModuleName "Microsoft.Graph.Authentication"

    # Discover commands for the URI
    Write-LogHeader "Discovering Commands for URI"
    Write-LogInfo "URI: $Uri" -Indent 1
    Write-LogInfo "API Version: $ApiVersion" -Indent 1
    if ($Method) {
        Write-LogInfo "Method: $Method" -Indent 1
    }

    $commands = Find-MgGraphCommand -Uri $Uri | Where-Object { $_.APIVersion -eq $ApiVersion }

    if ($Method) {
        $commands = $commands | Where-Object { $_.Method -eq $Method }
    }

    if (-not $commands) {
        Write-LogError "No commands found for URI: $Uri"
        Write-LogWarning "Try different variations:" -Indent 1
        Write-LogInfo "- Use {id} for resource IDs: /groups/{id}" -Indent 2
        Write-LogInfo "- Check both v1.0 and beta versions" -Indent 2
        exit 1
    }

    Write-LogSuccess "Found $($commands.Count) command(s)" -Indent 1

    # Collect permissions from all commands
    $allPermissions = @{}
    $applicationPerms = @()
    $delegatedPerms = @()

    Write-LogSection "Analyzing Commands and Permissions"

    foreach ($cmd in $commands) {
        Write-Host ""
        Write-LogInfo "Command: $($cmd.Command)" -Indent 1
        Write-LogInfo "Method:  $($cmd.Method)" -Indent 1
        Write-LogInfo "URI:     $($cmd.URI)" -Indent 1
        
        if ($cmd.Permissions -and $cmd.Permissions.Count -gt 0) {
            Write-LogInfo "Permissions:" -Indent 1
            
            foreach ($permObj in $cmd.Permissions) {
                $perm = Get-PermissionName -PermissionObject $permObj
                
                if (-not $perm) {
                    continue
                }
                
                $isAppPermission = Test-ApplicationPermission -PermissionName $perm
                
                if ($isAppPermission) {
                    if ($perm -notin $applicationPerms) {
                        $applicationPerms += $perm
                    }
                    Write-LogInfo "[App] $perm" -Indent 2
                }
                else {
                    if ($perm -notin $delegatedPerms) {
                        $delegatedPerms += $perm
                    }
                    Write-LogInfo "[Del] $perm" -Indent 2
                }
                
                if (-not $allPermissions.ContainsKey($perm)) {
                    $allPermissions[$perm] = @()
                }
                $allPermissions[$perm] += "$($cmd.Method) $($cmd.URI)"
            }
        }
        else {
            Write-LogWarning "No permissions documented" -Indent 2
        }
    }

    # Display permission summary
    Write-LogHeader "Permission Requirements Summary"

    if ($applicationPerms.Count -gt 0) {
        Write-LogSection "Application Permissions (for service principals)"
        foreach ($perm in $applicationPerms | Sort-Object) {
            Write-LogSuccess $perm -Indent 1
            $usage = $allPermissions[$perm] | Select-Object -Unique
            foreach ($use in $usage) {
                Write-LogInfo "→ $use" -Indent 2
            }
        }
    }

    if ($delegatedPerms.Count -gt 0) {
        Write-LogSection "Delegated Permissions (for user context)"
        foreach ($perm in $delegatedPerms | Sort-Object) {
            Write-LogInfo $perm -Indent 1
            $usage = $allPermissions[$perm] | Select-Object -Unique
            foreach ($use in $usage) {
                Write-LogInfo "→ $use" -Indent 2
            }
        }
    }

    # Get detailed permission information
    Write-LogHeader "Detailed Permission Information"

    $uniquePerms = ($applicationPerms + $delegatedPerms) | Select-Object -Unique

    foreach ($permName in $uniquePerms | Sort-Object) {
        $searchTerm = ($permName -split '\.')[0]
        
        $permDetails = Find-MgGraphPermission -SearchString $searchTerm | 
            Where-Object { $_.Name -eq $permName }
        
        if ($permDetails) {
            Write-Host ""
            Write-LogInfo "$($permDetails.Name)" -Indent 1
            Write-LogInfo "Type: $($permDetails.PermissionType)" -Indent 2
            Write-LogInfo "Consent: $($permDetails.Consent)" -Indent 2
            Write-LogInfo "Description: $($permDetails.Description)" -Indent 2
        }
    }

    # Evaluate enterprise application if requested
    if ($AppIdToBeEvaluated) {
        Write-LogHeader "Evaluating Enterprise Application"
        
        try {
            Write-LogInfo "Connecting to Microsoft Graph..." -Indent 1
            $secureSecret = ConvertTo-SecureString $ClientSecret -AsPlainText -Force
            $credential = New-Object System.Management.Automation.PSCredential($AppIdToBeEvaluated, $secureSecret)
            
            Connect-MgGraph -TenantId $TenantId -ClientSecretCredential $credential -NoWelcome -ErrorAction Stop
            Write-LogSuccess "Connected successfully" -Indent 1
            
            # Get the service principal
            Write-LogInfo "Retrieving service principal..." -Indent 1
            $sp = Get-MgServicePrincipal -Filter "appId eq '$AppIdToBeEvaluated'" -ErrorAction Stop
            
            if (-not $sp) {
                Write-LogError "Service principal not found with AppId: $AppIdToBeEvaluated" -Indent 1
                throw "Service principal not found"
            }
            
            Write-LogSuccess "Found: $($sp.DisplayName)" -Indent 1
            Write-LogInfo "Object ID: $($sp.Id)" -Indent 2
            
            # Get app role assignments
            Write-LogInfo "Retrieving granted permissions..." -Indent 1
            $appRoleAssignments = Get-MgServicePrincipalAppRoleAssignment -ServicePrincipalId $sp.Id
            
            # Get Microsoft Graph service principal
            $mgSp = Get-MgServicePrincipal -Filter "appId eq '00000003-0000-0000-c000-000000000000'"
            
            # Build list of granted permissions
            $grantedPerms = @()
            foreach ($assignment in $appRoleAssignments) {
                $appRole = $mgSp.AppRoles | Where-Object { $_.Id -eq $assignment.AppRoleId }
                if ($appRole) {
                    $grantedPerms += $appRole.Value
                }
            }
            
            # Compare required vs granted
            Write-LogSection "Permission Comparison"
            
            $missingPerms = @()
            $coveredPerms = @()
            
            foreach ($requiredPerm in $applicationPerms | Sort-Object) {
                if ($requiredPerm -in $grantedPerms) {
                    $coveredPerms += $requiredPerm
                    Write-LogSuccess "$requiredPerm (GRANTED)" -Indent 1
                }
                else {
                    $missingPerms += $requiredPerm
                    Write-LogError "$requiredPerm (MISSING)" -Indent 1
                }
            }
            
            # Check directory role memberships
            Write-LogSection "Directory Role Assignments"
            $memberOf = Get-MgServicePrincipalMemberOf -ServicePrincipalId $sp.Id
            $directoryRoles = $memberOf | Where-Object { $_.AdditionalProperties.'@odata.type' -eq '#microsoft.graph.directoryRole' }
            
            if ($directoryRoles.Count -eq 0) {
                Write-LogError "No directory roles assigned" -Indent 1
            }
            else {
                foreach ($roleRef in $directoryRoles) {
                    $role = Get-MgDirectoryRole -DirectoryRoleId $roleRef.Id
                    Write-LogSuccess $role.DisplayName -Indent 1
                }
            }
            
            # Final recommendations
            Write-LogHeader "Evaluation Summary"
            
            Write-LogSection "Coverage Status"
            Write-LogInfo "Total Required Permissions: $($applicationPerms.Count)" -Indent 1
            Write-LogSuccess "Covered: $($coveredPerms.Count)" -Indent 1
            Write-LogError "Missing: $($missingPerms.Count)" -Indent 1
            
            if ($missingPerms.Count -gt 0) {
                Write-LogSection "Missing Permissions"
                Write-LogInfo "The following permissions need to be granted:" -Indent 1
                foreach ($perm in $missingPerms) {
                    Write-LogError $perm -Indent 2
                }
                
                Write-Host ""
                Write-LogInfo "To grant these permissions:" -Indent 1
                Write-LogInfo "1. Go to Azure Portal > Entra ID > Enterprise Applications" -Indent 2
                Write-LogInfo "2. Find: $($sp.DisplayName)" -Indent 2
                Write-LogInfo "3. Go to: API Permissions" -Indent 2
                Write-LogInfo "4. Add the missing permissions listed above" -Indent 2
                Write-LogInfo "5. Grant admin consent" -Indent 2
            }
            else {
                Write-LogSuccess "All required API permissions are granted!" -Indent 1
            }
            
            # Check for special requirements
            if ($Uri -like "*/groups*" -and ($Method -eq "POST" -or $Method -eq "")) {
                Write-LogSection "Special Requirements for Role-Assignable Groups"
                if ($directoryRoles.Count -eq 0) {
                    Write-LogWarning "To create groups with isAssignableToRole=true:" -Indent 1
                    Write-LogInfo "Service principal needs one of these directory roles:" -Indent 2
                    Write-LogInfo "- Privileged Role Administrator" -Indent 3
                    Write-LogInfo "- Global Administrator" -Indent 3
                }
                else {
                    $hasRequiredRole = $false
                    foreach ($roleRef in $directoryRoles) {
                        $role = Get-MgDirectoryRole -DirectoryRoleId $roleRef.Id
                        if ($role.DisplayName -in @("Privileged Role Administrator", "Global Administrator")) {
                            $hasRequiredRole = $true
                            break
                        }
                    }
                    
                    if ($hasRequiredRole) {
                        Write-LogSuccess "Service principal has required directory role for role-assignable groups" -Indent 1
                    }
                    else {
                        Write-LogWarning "Service principal has directory roles but not the required ones for role-assignable groups" -Indent 1
                    }
                }
            }
            
            Disconnect-MgGraph | Out-Null
        }
        catch {
            Write-LogError "Failed to evaluate enterprise application: $($_.Exception.Message)"
            throw
        }
    }
    else {
        # No app evaluation - just show recommendations
        Write-LogHeader "Recommendations"
        
        if ($applicationPerms.Count -gt 0) {
            Write-LogSection "For Service Principal (Application) Authentication"
            Write-LogInfo "Grant these API permissions in Azure Portal:" -Indent 1
            foreach ($perm in $applicationPerms | Sort-Object) {
                Write-LogInfo "- $perm" -Indent 2
            }
            Write-Host ""
            Write-LogInfo "Admin consent is required for all permissions" -Indent 1
            
            if ($Uri -like "*/groups*" -and ($Method -eq "POST" -or $Method -eq "")) {
                Write-Host ""
                Write-LogInfo "For groups with isAssignableToRole=true:" -Indent 1
                Write-LogInfo "Service principal needs 'Privileged Role Administrator' OR" -Indent 2
                Write-LogInfo "'Global Administrator' directory role assigned" -Indent 2
            }
        }
    }

    # Documentation references
    Write-LogHeader "Documentation References"
    Write-LogInfo "Microsoft Graph Permissions Reference:" -Indent 1
    Write-LogInfo "https://learn.microsoft.com/en-us/graph/permissions-reference" -Indent 2
    Write-Host ""
    Write-LogInfo "How to Figure Out Graph Permissions:" -Indent 1
    Write-LogInfo "https://practical365.com/microsoft-graph-api-permission/" -Indent 2
}
catch {
    Write-LogError "Script execution failed: $($_.Exception.Message)"
    exit 1
}

#endregion
