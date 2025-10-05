Request URL
https://graph.microsoft.com/beta/deviceManagement/deviceEnrollmentConfigurations/54fac284-7866-43e5-860a-9c8e10fa3d7d_WindowsRestore
Request Method
PATCH

{"@odata.type":"#microsoft.graph.windowsRestoreDeviceEnrollmentConfiguration","state":"enabled"}

Request URL
https://graph.microsoft.com/beta/deviceManagement/deviceEnrollmentConfigurations?$filter=deviceEnrollmentConfigurationType%20eq%20%27WindowsRestore%27
Request Method
GET

{
    "@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/deviceEnrollmentConfigurations",
    "value": [
        {
            "@odata.type": "#microsoft.graph.windowsRestoreDeviceEnrollmentConfiguration",
            "id": "54fac284-7866-43e5-860a-9c8e10fa3d7d_WindowsRestore",
            "displayName": "All users and all devices",
            "description": "This is the default Windows Restore configuration applied with the lowest priority to all users and all devices regardless of group membership.",
            "priority": 0,
            "createdDateTime": "0001-01-01T00:00:00Z",
            "lastModifiedDateTime": "2025-10-03T13:01:20Z",
            "version": 0,
            "roleScopeTagIds": [],
            "deviceEnrollmentConfigurationType": "windowsRestore",
            "state": "enabled"
        }
    ]
}

Request URL
https://graph.microsoft.com/beta/deviceManagement/deviceEnrollmentConfigurations/54fac284-7866-43e5-860a-9c8e10fa3d7d_WindowsRestore
Request Method
PATCH

{"@odata.type":"#microsoft.graph.windowsRestoreDeviceEnrollmentConfiguration","state":"disabled"}

Request URL
https://graph.microsoft.com/beta/deviceManagement/deviceEnrollmentConfigurations?$filter=deviceEnrollmentConfigurationType%20eq%20%27WindowsRestore%27
Request Method
GET

{
    "@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/deviceEnrollmentConfigurations",
    "value": [
        {
            "@odata.type": "#microsoft.graph.windowsRestoreDeviceEnrollmentConfiguration",
            "id": "54fac284-7866-43e5-860a-9c8e10fa3d7d_WindowsRestore",
            "displayName": "All users and all devices",
            "description": "This is the default Windows Restore configuration applied with the lowest priority to all users and all devices regardless of group membership.",
            "priority": 0,
            "createdDateTime": "0001-01-01T00:00:00Z",
            "lastModifiedDateTime": "2025-10-03T13:02:14Z",
            "version": 0,
            "roleScopeTagIds": [],
            "deviceEnrollmentConfigurationType": "windowsRestore",
            "state": "disabled"
        }
    ]
}

{"@odata.type":"#microsoft.graph.windowsRestoreDeviceEnrollmentConfiguration","state":"notConfigured"}