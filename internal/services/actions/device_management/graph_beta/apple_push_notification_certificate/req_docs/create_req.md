Request URL
https://graph.microsoft.com/beta/deviceManagement/applePushNotificationCertificate
Request Method
PATCH

{
    "appleIdentifier":"dafydd.watkins@bankofscotland.appleaccount.com",
    "certificate":"-----BEGIN CERTIFICATE-----\nMIIFdjCCWU1HjDih\nqUgw3TNZ7z+h/Q==\n-----END CERTIFICATE-----"}


resp

{
    "@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/applePushNotificationCertificate/$entity",
    "id": "54fac284-7866-43e5-860a-9c8e10fa3d7d",
    "appleIdentifier": "dafydd.watkins@bankofscotland.appleaccount.com",
    "topicIdentifier": "com.apple.mgmt.External.2f42b322-d7c9-4734-8e2e-eeadadd36f20",
    "lastModifiedDateTime": "2025-10-06T13:11:29Z",
    "expirationDateTime": "2026-10-06T13:11:28Z",
    "certificateUploadStatus": null,
    "certificateUploadFailureReason": null,
    "certificateSerialNumber": "0B7002EFEA6BA6D0",
    "certificate": null
}

Request URL
https://graph.microsoft.com/beta/deviceManagement/applePushNotificationCertificate
Request Method
GET

{
    "@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/applePushNotificationCertificate/$entity",
    "id": "5047ff77-0e37-495a-a672-ea7a04a319e0",
    "appleIdentifier": "dafydd.watkins@bankofscotland.appleaccount.com",
    "topicIdentifier": "com.apple.mgmt.External.2f42b322-d7c9-4734-8e2e-eeadadd36f20",
    "lastModifiedDateTime": "2025-10-06T13:22:48Z",
    "expirationDateTime": "2026-10-06T13:11:28Z",
    "certificateUploadStatus": null,
    "certificateUploadFailureReason": null,
    "certificateSerialNumber": "0B7002EFEA6BA6D0",
    "certificate": null
}

Request URL
https://graph.microsoft.com/beta/deviceManagement/dataSharingConsents/appleMDMPushCertificate
Request Method
GET

{
    "@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/dataSharingConsents/$entity",
    "id": "appleMDMPushCertificate",
    "serviceDisplayName": "Apple MDM Push Certificate",
    "termsUrl": "https://go.microsoft.com/fwlink/?linkid=866314",
    "granted": true,
    "grantDateTime": "2021-08-20T13:48:33.669234Z",
    "grantedByUpn": "admin_d.watkins@deploymenttheory.com",
    "grantedByUserId": "40b18b93-ae0d-45f5-98fe-6579fc792faa"
}