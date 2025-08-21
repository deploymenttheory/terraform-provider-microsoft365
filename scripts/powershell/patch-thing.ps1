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
    HelpMessage="Display name for the enrollment notification configuration")]
    [ValidateNotNullOrEmpty()]
    [string]$DisplayName = "Default Enrollment Notification",
    
    [Parameter(Mandatory=$false,
    HelpMessage="Description for the enrollment notification configuration")]
    [ValidateNotNullOrEmpty()]
    [string]$Description = "Default enrollment notification configuration",
    
    [Parameter(Mandatory=$false,
    HelpMessage="Platform type for the enrollment notification")]
    [ValidateSet('ios', 'android', 'windows', 'mac', 'androidForWork')]
    [string]$PlatformType = "androidForWork",
    
    [Parameter(Mandatory=$false,
    HelpMessage="Default locale for the notification")]
    [ValidateNotNullOrEmpty()]
    [string]$DefaultLocale = "en-US",
    
    [Parameter(Mandatory=$false,
    HelpMessage="Branding options for the notification")]
    [ValidateSet('none', 'includeCompanyLogo', 'includeCompanyName', 'includeCompanyLogoAndName', 'includeContactInformation')]
    [string]$BrandingOptions = "none",
    
    [Parameter(Mandatory=$false,
    HelpMessage="Notification template IDs")]
    [ValidateNotNullOrEmpty()]
    [string[]]$NotificationTemplates = @("push_00000000-0000-0000-0000-000000000000"),
    
    [Parameter(Mandatory=$false,
    HelpMessage="Role scope tag IDs")]
    [ValidateNotNullOrEmpty()]
    [string[]]$RoleScopeTagIds = @("0"),

    [Parameter(Mandatory=$false,
    HelpMessage="Subject for localized notification message")]
    [ValidateNotNullOrEmpty()]
    [string]$NotificationSubject = "Enrollment Notification",

    [Parameter(Mandatory=$false,
    HelpMessage="Message template for localized notification")]
    [ValidateNotNullOrEmpty()]
    [string]$NotificationMessage = "Your device enrollment is in progress."
)

# Example usage:
# ./Create-IntuneEnrollmentNotifications.ps1 -TenantId "your-tenant-id" -ClientId "your-client-id" -ClientSecret "your-client-secret"
# ./Create-IntuneEnrollmentNotifications.ps1 -TenantId "your-tenant-id" -ClientId "your-client-id" -ClientSecret "your-client-secret" -DisplayName "Custom Notification" -PlatformType "ios" -NotificationSubject "Custom Subject" -NotificationMessage "Custom message for enrollment"

Import-Module Microsoft.Graph.Authentication

function New-EnrollmentNotificationConfiguration {
    param (
        [Parameter(Mandatory=$true)]
        [string]$DisplayName,
        [Parameter(Mandatory=$true)]
        [string]$Description,
        [Parameter(Mandatory=$true)]
        [string]$PlatformType,
        [Parameter(Mandatory=$true)]
        [string]$DefaultLocale,
        [Parameter(Mandatory=$true)]
        [string]$BrandingOptions,
        [Parameter(Mandatory=$true)]
        [string[]]$NotificationTemplates,
        [Parameter(Mandatory=$true)]
        [string[]]$RoleScopeTagIds
    )
    
    try {
        Write-Host "📝 Creating enrollment notification configuration..." -ForegroundColor Cyan
        Write-Host "   Display Name: $DisplayName" -ForegroundColor Gray
        Write-Host "   Platform Type: $PlatformType" -ForegroundColor Gray
        
        $createUri = "https://graph.microsoft.com/beta/deviceManagement/deviceEnrollmentConfigurations"
        Write-Host "   Endpoint: $createUri" -ForegroundColor Gray
        
        $body = @{
            "@odata.type" = "#microsoft.graph.deviceEnrollmentNotificationConfiguration"
            "displayName" = $DisplayName
            "description" = $Description
            "platformType" = $PlatformType
            "defaultLocale" = $DefaultLocale
            "brandingOptions" = $BrandingOptions
            "notificationTemplates" = $NotificationTemplates
            "roleScopeTagIds" = $RoleScopeTagIds
        }
        
        Write-Host "   Request Body:" -ForegroundColor Gray
        Write-Host "   $($body | ConvertTo-Json)" -ForegroundColor Gray
        
        $response = Invoke-MgGraphRequest -Method POST -Uri $createUri -Body ($body | ConvertTo-Json)
        Write-Host "   ✅ Created enrollment notification configuration with ID: $($response.id)" -ForegroundColor Green
        
        return $response
    }
    catch {
        Write-Host "❌ Error creating enrollment notification configuration: $_" -ForegroundColor Red
        Write-Host ""
        if ($_.Exception.Response) {
            $statusCode = $_.Exception.Response.StatusCode
            $statusDescription = $_.Exception.Response.StatusDescription
            Write-Host "   Status Code: $statusCode" -ForegroundColor Red
            Write-Host "   Status Description: $statusDescription" -ForegroundColor Red
            try {
                $responseContent = $_.Exception.Response.GetResponseStream()
                if ($responseContent) {
                    $reader = [System.IO.StreamReader]::new($responseContent)
                    $errorDetails = $reader.ReadToEnd()
                    $reader.Close()
                    Write-Host "   Error Details: $errorDetails" -ForegroundColor Red
                }
            } catch {
                Write-Host "   Could not read error details" -ForegroundColor Red
            }
        }
        return $null
    }
}

function Get-EnrollmentNotificationConfigurationById {
    param (
        [Parameter(Mandatory=$true)]
        [string]$ConfigurationId
    )
    
    try {
        Write-Host "🔍 Getting enrollment notification configuration by ID..." -ForegroundColor Cyan
        Write-Host "   Configuration ID: $ConfigurationId" -ForegroundColor Gray
        
        $configUri = "https://graph.microsoft.com/beta/deviceManagement/deviceEnrollmentConfigurations/$ConfigurationId"
        Write-Host "   Endpoint: $configUri" -ForegroundColor Gray
        
        $configuration = Invoke-MgGraphRequest -Method GET -Uri $configUri
        Write-Host "   ✅ Configuration retrieved successfully" -ForegroundColor Green
        Write-Host ""
        
        return $configuration
    }
    catch {
        Write-Host "❌ Error retrieving enrollment notification configuration: $_" -ForegroundColor Red
        Write-Host ""
        if ($_.Exception.Response) {
            $statusCode = $_.Exception.Response.StatusCode
            $statusDescription = $_.Exception.Response.StatusDescription
            Write-Host "   Status Code: $statusCode" -ForegroundColor Red
            Write-Host "   Status Description: $statusDescription" -ForegroundColor Red
            if ($statusCode -eq 404) {
                Write-Host "   Configuration not found with ID: $ConfigurationId" -ForegroundColor Red
            }
        }
        throw
    }
}

function Update-EnrollmentNotificationConfiguration {
    param (
        [Parameter(Mandatory=$true)]
        [string]$ConfigurationId,
        [Parameter(Mandatory=$true)]
        [string]$DisplayName,
        [Parameter(Mandatory=$true)]
        [string]$Description,
        [Parameter(Mandatory=$true)]
        [string]$PlatformType,
        [Parameter(Mandatory=$true)]
        [string]$DefaultLocale,
        [Parameter(Mandatory=$true)]
        [string]$BrandingOptions,
        [Parameter(Mandatory=$true)]
        [string[]]$NotificationTemplates,
        [Parameter(Mandatory=$true)]
        [string[]]$RoleScopeTagIds
    )
    
    try {
        Write-Host "🔄 Updating enrollment notification configuration..." -ForegroundColor Cyan
        Write-Host "   Configuration ID: $ConfigurationId" -ForegroundColor Gray
        Write-Host "   Display Name: $DisplayName" -ForegroundColor Gray
        
        $updateUri = "https://graph.microsoft.com/beta/deviceManagement/deviceEnrollmentConfigurations/$ConfigurationId"
        Write-Host "   Endpoint: $updateUri" -ForegroundColor Gray
        
        $body = @{
            "@odata.type" = "#microsoft.graph.deviceEnrollmentNotificationConfiguration"
            "displayName" = $DisplayName
            "description" = $Description
            "platformType" = $PlatformType
            "defaultLocale" = $DefaultLocale
            "brandingOptions" = $BrandingOptions
            "notificationTemplates" = $NotificationTemplates
            "roleScopeTagIds" = $RoleScopeTagIds
        }
        
        Write-Host "   Request Body:" -ForegroundColor Gray
        Write-Host "   $($body | ConvertTo-Json)" -ForegroundColor Gray
        
        $response = Invoke-MgGraphRequest -Method PATCH -Uri $updateUri -Body ($body | ConvertTo-Json)
        Write-Host "   ✅ Updated enrollment notification configuration successfully" -ForegroundColor Green
        
        return $response
    }
    catch {
        Write-Host "❌ Error updating enrollment notification configuration: $_" -ForegroundColor Red
        Write-Host ""
        if ($_.Exception.Response) {
            $statusCode = $_.Exception.Response.StatusCode
            $statusDescription = $_.Exception.Response.StatusDescription
            Write-Host "   Status Code: $statusCode" -ForegroundColor Red
            Write-Host "   Status Description: $statusDescription" -ForegroundColor Red
            try {
                $responseContent = $_.Exception.Response.GetResponseStream()
                if ($responseContent) {
                    $reader = [System.IO.StreamReader]::new($responseContent)
                    $errorDetails = $reader.ReadToEnd()
                    $reader.Close()
                    Write-Host "   Error Details: $errorDetails" -ForegroundColor Red
                }
            } catch {
                Write-Host "   Could not read error details" -ForegroundColor Red
            }
        }
        return $null
    }
}

function New-LocalizedNotificationMessage {
    param (
        [Parameter(Mandatory=$true)]
        [string]$NotificationTemplateId,
        [Parameter(Mandatory=$true)]
        [string]$Locale,
        [Parameter(Mandatory=$true)]
        [string]$Subject,
        [Parameter(Mandatory=$true)]
        [string]$MessageTemplate,
        [Parameter(Mandatory=$false)]
        [bool]$IsDefault = $true
    )
    
    try {
        Write-Host "📝 Creating localized notification message..." -ForegroundColor Cyan
        Write-Host "   Template ID: $NotificationTemplateId" -ForegroundColor Gray
        Write-Host "   Locale: $Locale" -ForegroundColor Gray
        Write-Host "   Subject: $Subject" -ForegroundColor Gray
        
        $createUri = "https://graph.microsoft.com/beta/deviceManagement/notificationMessageTemplates/$NotificationTemplateId/localizedNotificationMessages"
        Write-Host "   Endpoint: $createUri" -ForegroundColor Gray
        
        $body = @{
            "locale" = $Locale
            "isDefault" = $IsDefault
            "subject" = $Subject
            "messageTemplate" = $MessageTemplate
        }
        
        Write-Host "   Request Body:" -ForegroundColor Gray
        Write-Host "   $($body | ConvertTo-Json)" -ForegroundColor Gray
        
        $response = Invoke-MgGraphRequest -Method POST -Uri $createUri -Body ($body | ConvertTo-Json)
        Write-Host "   ✅ Created localized message with ID: $($response.id)" -ForegroundColor Green
        
        return $response
    }
    catch {
        Write-Host "❌ Error creating localized notification message: $_" -ForegroundColor Red
        Write-Host "   Template ID: $NotificationTemplateId" -ForegroundColor Red
        if ($_.Exception.Response) {
            $statusCode = $_.Exception.Response.StatusCode
            $statusDescription = $_.Exception.Response.StatusDescription
            Write-Host "   Status Code: $statusCode" -ForegroundColor Red
            Write-Host "   Status Description: $statusDescription" -ForegroundColor Red
        }
        return $null
    }
}

function Get-NotificationTemplateId {
    param (
        [Parameter(Mandatory=$true)]
        [string]$NotificationTemplate
    )
    
    try {
        # Check if the template starts with "push_" and remove it
        if ($NotificationTemplate -match "^push_(.+)$") {
            $templateId = $Matches[1]
            Write-Host "🔍 Extracted notification template ID: $templateId" -ForegroundColor Cyan
            return $templateId
        } else {
            Write-Host "⚠️ Notification template does not have expected 'push_' prefix: $NotificationTemplate" -ForegroundColor Yellow
            return $NotificationTemplate
        }
    }
    catch {
        Write-Host "❌ Error extracting notification template ID: $_" -ForegroundColor Red
        return $null
    }
}

function Show-ConfigurationDetails {
    param (
        [Parameter(Mandatory=$true)]
        $Configuration
    )
    
    Write-Host "📋 Enrollment Notification Configuration Details:" -ForegroundColor Cyan
    Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
    
    # Basic configuration information
    foreach ($field in @('id', 'displayName', 'description', 'deviceEnrollmentConfigurationType', 'platformType', 'defaultLocale', 'brandingOptions', 'roleScopeTagIds', 'createdDateTime', 'lastModifiedDateTime', 'version')) {
        if ($Configuration.PSObject.Properties[$field]) {
            if ($field -eq 'roleScopeTagIds' -and $Configuration.$field) {
                Write-Host "   • roleScopeTagIds: $($Configuration.$field -join ', ')" -ForegroundColor Green
            } else {
                Write-Host ("   • {0}: {1}" -f $field, $Configuration.$field) -ForegroundColor Green
            }
        }
    }
    
    # Notification templates information
    if ($Configuration.notificationTemplates -and $Configuration.notificationTemplates.Count -gt 0) {
        Write-Host "   • notificationTemplates: $($Configuration.notificationTemplates -join ', ')" -ForegroundColor Green
    }
    
    Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
    Write-Host ""
}

function Show-LocalizedMessageDetails {
    param (
        [Parameter(Mandatory=$true)]
        $LocalizedMessage
    )
    
    Write-Host "📋 Localized Notification Message Details:" -ForegroundColor Cyan
    Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
    
    # Basic message information
    foreach ($field in @('id', 'locale', 'subject', 'messageTemplate', 'isDefault', 'lastModifiedDateTime')) {
        if ($LocalizedMessage.PSObject.Properties[$field]) {
            Write-Host ("   • {0}: {1}" -f $field, $LocalizedMessage.$field) -ForegroundColor Green
        }
    }
    
    Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
    Write-Host ""
}

try {
    Write-Host "🔐 Connecting to Microsoft Graph..." -ForegroundColor Cyan
    Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
    
    $secureClientSecret = ConvertTo-SecureString -String $ClientSecret -AsPlainText -Force
    $clientSecretCredential = New-Object -TypeName System.Management.Automation.PSCredential -ArgumentList $ClientId, $secureClientSecret
    
    Connect-MgGraph -ClientSecretCredential $clientSecretCredential -TenantId $TenantId
    Write-Host "✅ Connected to Microsoft Graph successfully" -ForegroundColor Green
    Write-Host ""
    
    # Step 1: Create enrollment notification configuration
    Write-Host "🔄 Step 1: Creating enrollment notification configuration..." -ForegroundColor Cyan
    $newConfiguration = New-EnrollmentNotificationConfiguration -DisplayName $DisplayName -Description $Description -PlatformType $PlatformType -DefaultLocale $DefaultLocale -BrandingOptions $BrandingOptions -NotificationTemplates $NotificationTemplates -RoleScopeTagIds $RoleScopeTagIds
    
    if ($newConfiguration) {
        Write-Host "🎉 Enrollment notification configuration created successfully!" -ForegroundColor Green
        
        # Get and show the created configuration details
        $configDetails = Get-EnrollmentNotificationConfigurationById -ConfigurationId $newConfiguration.id
        Show-ConfigurationDetails -Configuration $configDetails
        
        # Step 2: Create localized notification message for each notification template
        Write-Host "🔄 Step 2: Creating localized notification messages..." -ForegroundColor Cyan
        
        foreach ($template in $newConfiguration.notificationTemplates) {
            $templateId = Get-NotificationTemplateId -NotificationTemplate $template
            
            if ($templateId) {
                Write-Host "   Processing template: $template (ID: $templateId)" -ForegroundColor Gray
                
                $localizedMessage = New-LocalizedNotificationMessage -NotificationTemplateId $templateId -Locale $DefaultLocale -Subject $NotificationSubject -MessageTemplate $NotificationMessage -IsDefault $true
                
                if ($localizedMessage) {
                    Write-Host "   ✅ Created localized message for template: $template" -ForegroundColor Green
                    Show-LocalizedMessageDetails -LocalizedMessage $localizedMessage
                } else {
                    Write-Host "   ❌ Failed to create localized message for template: $template" -ForegroundColor Red
                }
            } else {
                Write-Host "   ❌ Failed to extract template ID from: $template" -ForegroundColor Red
            }
        }
        
        Write-Host "🎉 Script execution completed successfully!" -ForegroundColor Green
    } else {
        Write-Host "❌ Failed to create enrollment notification configuration" -ForegroundColor Red
    }
}
catch {
    Write-Host "❌ Script execution failed: $_" -ForegroundColor Red
    exit 1
}
finally {
    Write-Host "🔌 Disconnecting from Microsoft Graph..." -ForegroundColor Cyan
    try {
        Disconnect-MgGraph 2>$null
        Write-Host "✅ Disconnected from Microsoft Graph" -ForegroundColor Green
    } catch {}
}