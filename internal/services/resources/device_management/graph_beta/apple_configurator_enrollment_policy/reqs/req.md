Request URL
https://graph.microsoft.com/beta/deviceManagement
Request Method
GET

resp

{
    "@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement",
    "id": "deviceManagement",
    "settings": null,
    "maximumDepTokens": 100,
    "intuneAccountId": "54fac284-7866-43e5-860a-9c8e10fa3d7d",
    "lastReportAggregationDateTime": "0001-01-01T00:00:00Z",
    "deviceComplianceReportSummarizationDateTime": "0001-01-01T00:00:00Z",
    "legacyPcManangementEnabled": false,
    "unlicensedAdminstratorsEnabled": true
}

Request URL
https://graph.microsoft.com/beta/deviceManagement/depOnboardingSettings/54fac284-7866-43e5-860a-9c8e10fa3d7d/enrollmentProfiles
Request Method
POST

{
  "@odata.type":"#microsoft.graph.enrollmentProfile",
  "displayName":"test",
  "description":"test",
  "requiresUserAuthentication":true,
  "enableAuthenticationViaCompanyPortal":true,
  "requireCompanyPortalOnSetupAssistantEnrolledDevices":false
}

Request URL
https://graph.microsoft.com/beta/deviceManagement/depOnboardingSettings/54fac284-7866-43e5-860a-9c8e10fa3d7d/enrollmentProfiles/54fac284-7866-43e5-860a-9c8e10fa3d7d_b4cda2e6-3d9f-4eb6-aac3-bab920178f4a
Request Method
GET

{
    "@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/depOnboardingSettings('54fac284-7866-43e5-860a-9c8e10fa3d7d')/enrollmentProfiles/$entity",
    "id": "54fac284-7866-43e5-860a-9c8e10fa3d7d_b4cda2e6-3d9f-4eb6-aac3-bab920178f4a",
    "displayName": "test",
    "description": "test",
    "requiresUserAuthentication": true,
    "configurationEndpointUrl": "https://appleconfigurator2.manage.microsoft.com/EnrollmentServer/MDMServiceConfig?id=54fac284-7866-43e5-860a-9c8e10fa3d7d&AADTenantId=2fd6bb84-ad40-4ec5-9369-a215b25c9952",
    "enableAuthenticationViaCompanyPortal": true,
    "requireCompanyPortalOnSetupAssistantEnrolledDevices": false
}

Request URL
https://graph.microsoft.com/beta/deviceManagement/depOnboardingSettings/54fac284-7866-43e5-860a-9c8e10fa3d7d/enrollmentProfiles/54fac284-7866-43e5-860a-9c8e10fa3d7d_b4cda2e6-3d9f-4eb6-aac3-bab920178f4a
Request Method
PATCH

{
  "@odata.type":"#microsoft.graph.enrollmentProfile",
  "id":"54fac284-7866-43e5-860a-9c8e10fa3d7d_b4cda2e6-3d9f-4eb6-aac3-bab920178f4a",
  "displayName":"test",
  "description":"test",
  "requiresUserAuthentication":true,
  "enableAuthenticationViaCompanyPortal":false,
  "requireCompanyPortalOnSetupAssistantEnrolledDevices":true
}

Request URL
https://graph.microsoft.com/beta/deviceManagement/depOnboardingSettings/54fac284-7866-43e5-860a-9c8e10fa3d7d/importedAppleDeviceIdentities?$count=true&$filter=discoverySource%20eq%20%27AdminImport%27%20and%20requestedEnrollmentProfileId%20eq%20%27b4cda2e6-3d9f-4eb6-aac3-bab920178f4a%27
Request Method
GET

resp

{
    "@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/depOnboardingSettings('54fac284-7866-43e5-860a-9c8e10fa3d7d')/importedAppleDeviceIdentities",
    "@odata.count": 0,
    "value": []
}


Request URL
https://graph.microsoft.com/beta/deviceManagement/depOnboardingSettings/54fac284-7866-43e5-860a-9c8e10fa3d7d/enrollmentProfiles/54fac284-7866-43e5-860a-9c8e10fa3d7d_b4cda2e6-3d9f-4eb6-aac3-bab920178f4a
Request Method
DELETE