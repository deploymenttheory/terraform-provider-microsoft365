# MMD API Docs

https://skiptotheendpoint.co.uk/under-the-hood-pt-1-autopatch/

# Permissions

## List Devices


https://mmdls.microsoft.com/api/v1.0/tenant
GET
{
  "directoryId": "2fd6bb84-ad40-4ec5-9369-a215b25c9952",
  "domain": "deploymenttheory.com",
  "state": "Enrolled",
  "appPackagingEnrolled": false,
  "appLockerVersion": "0.0",
  "sfBDisabled": true,
  "testVersion": 0.0,
  "testTarget": 0.0,
  "firstVersion": 0.0,
  "firstTarget": 0.0,
  "fastVersion": 0.0,
  "fastTarget": 0.0,
  "broadVersion": 0.0,
  "broadTarget": 0.0,
  "enableWindowsDeviceLocation": false,
  "enableReassignGroupTag": true,
  "enableRenameDevices": true,
  "entitlementSegment": "S_50",
  "readOnly": false,
  "plans": [
    {
      "plan": "Starter",
      "state": "Enrolled",
      "onboardedTimestampUtc": "2024-01-19T08:33:19.5281115",
      "partiallyEnrolledTimestampUtc": "2024-01-19T08:33:37.8781168",
      "enrolledTimestampUtc": "2024-01-19T08:36:22.9335105",
      "enrollmentType": "ClassicEnrollment",
      "failedEnrollmentAttempts": 0
    }
  ],
  "@self": {
    "rel": "self",
    "href": "https://customerapi.eu03.mmdprod.trafficmanager.net/api/v1.0/tenant"
  },
  "@links": [
    {
      "rel": "related",
      "title": "View devices for this tenant.",
      "href": "https://customerapi.eu03.mmdprod.trafficmanager.net/api/v1.0/tenant/devices"
    }
  ]
}


https://mmdls.microsoft.com/device/v1/windows365/autopatchGroups
GET

[
  {
    "id": "4aa9b805-9494-4eed-a04b-ed51ec9e631e",
    "name": "Windows Autopatch",
    "tenantId": "2fd6bb84-ad40-4ec5-9369-a215b25c9952",
    "description": "Windows Autopatch, the default Autopatch Group",
    "type": "Default"
  }
]

https://mmdls.microsoft.com/device/v1/windows365/autopatchGroups