Request URL
https://graph.microsoft.com/beta/deviceAppManagement/windowsManagementApp
Request Method
GET

response

{
    "@odata.context": "https://graph.microsoft.com/beta/$metadata#microsoft.graph.windowsManagementApp",
    "id": "54fac284-7866-43e5-860a-9c8e10fa3d7d",
    "availableVersion": "1.93.102.0",
    "managedInstaller": "disabled",
    "managedInstallerConfiguredDateTime": null
}

then to enable it

Request URL
https://graph.microsoft.com/beta/deviceAppManagement/windowsManagementApp/setAsManagedInstaller
Request Method
POST

empty body.

then get again

Request URL
https://graph.microsoft.com/beta/deviceAppManagement/windowsManagementApp
Request Method
GET

now it's enabled

{
    "@odata.context": "https://graph.microsoft.com/beta/$metadata#microsoft.graph.windowsManagementApp",
    "id": "54fac284-7866-43e5-860a-9c8e10fa3d7d",
    "availableVersion": "1.93.102.0",
    "managedInstaller": "enabled",
    "managedInstallerConfiguredDateTime": "8/23/2025 7:51:54 AM +00:00"
}
post again to disable it. 

Request URL
https://graph.microsoft.com/beta/deviceAppManagement/windowsManagementApp/setAsManagedInstaller
Request Method
POST

what it all means - https://learn.microsoft.com/en-us/windows/security/application-security/application-control/app-control-for-business/design/configure-authorized-apps-deployed-with-a-managed-installer