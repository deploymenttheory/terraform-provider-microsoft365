{requests: [{id: "a47e5cac-b696-4d66-b790-7a5741187620", method: "GET",…},…]}
requests
: 
[{id: "a47e5cac-b696-4d66-b790-7a5741187620", method: "GET",…},…]
0
: 
{id: "a47e5cac-b696-4d66-b790-7a5741187620", method: "GET",…}
headers
: 
{x-ms-command-name: "TenantManagement - GetDefaultXTAP",…}
id
: 
"a47e5cac-b696-4d66-b790-7a5741187620"
method
: 
"GET"
url
: 
"/policies/crossTenantAccessPolicy/default"
1
: 
{id: "240cf715-589e-485d-a2b0-8503710341a7", method: "POST",…}
body
: 
{,…}
resourceActionAuthorizationChecks
: 
[{directoryScopeId: "/",…}, {directoryScopeId: "/",…}, {directoryScopeId: "/",…},…]
0
: 
{directoryScopeId: "/",…}
directoryScopeId
: 
"/"
resourceAction
: 
"microsoft.directory/crossTenantAccessPolicy/default/b2bCollaboration/update"
1
: 
{directoryScopeId: "/",…}
directoryScopeId
: 
"/"
resourceAction
: 
"microsoft.directory/crossTenantAccessPolicy/default/b2bDirectConnect/update"
2
: 
{directoryScopeId: "/",…}
directoryScopeId
: 
"/"
resourceAction
: 
"microsoft.directory/crossTenantAccessPolicy/partners/b2bCollaboration/update"
3
: 
{directoryScopeId: "/",…}
directoryScopeId
: 
"/"
resourceAction
: 
"microsoft.directory/crossTenantAccessPolicy/partners/b2bDirectConnect/update"
headers
: 
{x-ms-command-name: "RBACv2 - estimateAccess",…}
Content-Type
: 
"application/json"
client-request-id
: 
"551f48f6-f0c8-44bf-9d17-46d321d576f9"
x-ms-client-request-id
: 
"551f48f6-f0c8-44bf-9d17-46d321d576f9"
x-ms-client-session-id
: 
"2fed463aef6c449ba1abf3f42e20e0f4"
x-ms-command-name
: 
"RBACv2 - estimateAccess"
id
: 
"240cf715-589e-485d-a2b0-8503710341a7"
method
: 
"POST"
url
: 
"/roleManagement/directory/estimateAccess"

"responses": [
        {
            "id": "a47e5cac-b696-4d66-b790-7a5741187620",
            "status": 200,
            "headers": {
                "Cache-Control": "no-cache",
                "x-ms-resource-unit": "4",
                "OData-Version": "4.0",
                "Content-Type": "application/json;odata.metadata=minimal;odata.streaming=true;IEEE754Compatible=false;charset=utf-8"
            },
            "body": {
                "@odata.context": "https://graph.microsoft.com/beta/$metadata#policies/crossTenantAccessPolicy/default/$entity",
                "id": "5c027107-6ba7-42e7-9583-7df394bffd0a",
                "isServiceDefault": false,
                "inboundTrust": {
                    "isMfaAccepted": false,
                    "isCompliantDeviceAccepted": false,
                    "isHybridAzureADJoinedDeviceAccepted": false,
                    "isCompliantNetworkAccepted": false
                },
                "b2bCollaborationOutbound": {
                    "usersAndGroups": {
                        "accessType": "allowed",
                        "targets": [
                            {
                                "target": "AllUsers",
                                "targetType": "user"
                            }
                        ]
                    },
                    "applications": {
                        "accessType": "allowed",
                        "targets": [
                            {
                                "target": "AllApplications",
                                "targetType": "application"
                            }
                        ]
                    }
                },
                "b2bCollaborationInbound": {
                    "usersAndGroups": {
                        "accessType": "allowed",
                        "targets": [
                            {
                                "target": "AllUsers",
                                "targetType": "user"
                            }
                        ]
                    },
                    "applications": {
                        "accessType": "allowed",
                        "targets": [
                            {
                                "target": "AllApplications",
                                "targetType": "application"
                            }
                        ]
                    }
                },
                "b2bDirectConnectOutbound": {
                    "usersAndGroups": {
                        "accessType": "blocked",
                        "targets": [
                            {
                                "target": "AllUsers",
                                "targetType": "user"
                            }
                        ]
                    },
                    "applications": {
                        "accessType": "blocked",
                        "targets": [
                            {
                                "target": "AllApplications",
                                "targetType": "application"
                            }
                        ]
                    }
                },
                "b2bDirectConnectInbound": {
                    "usersAndGroups": {
                        "accessType": "blocked",
                        "targets": [
                            {
                                "target": "AllUsers",
                                "targetType": "user"
                            }
                        ]
                    },
                    "applications": {
                        "accessType": "blocked",
                        "targets": [
                            {
                                "target": "AllApplications",
                                "targetType": "application"
                            }
                        ]
                    }
                },
                "crossCloudMeetingConfiguration": {
                    "inboundAllowed": false,
                    "outboundAllowed": false
                },
                "automaticUserConsentSettings": {
                    "inboundAllowed": false,
                    "outboundAllowed": false
                },
                "protectedContentSharing": {
                    "inboundAllowed": true,
                    "outboundAllowed": true
                },
                "tenantRestrictions": {
                    "devices": null,
                    "usersAndGroups": {
                        "accessType": "blocked",
                        "targets": [
                            {
                                "target": "AllUsers",
                                "targetType": "user"
                            }
                        ]
                    },
                    "applications": {
                        "accessType": "blocked",
                        "targets": [
                            {
                                "target": "AllApplications",
                                "targetType": "application"
                            }
                        ]
                    }
                },
                "invitationRedemptionIdentityProviderConfiguration": {
                    "primaryIdentityProviderPrecedenceOrder": [
                        "externalFederation",
                        "azureActiveDirectory",
                        "socialIdentityProviders"
                    ],
                    "fallbackIdentityProvider": "defaultConfiguredIdp"
                },
                "m365CollaborationInbound": {
                    "users": {
                        "accessType": "blocked",
                        "targets": [
                            {
                                "target": "AllUsers",
                                "targetType": "user"
                            }
                        ]
                    }
                },
                "m365CollaborationOutbound": {
                    "usersAndGroups": {
                        "accessType": "allowed",
                        "targets": [
                            {
                                "target": "AllUsers",
                                "targetType": "user"
                            }
                        ]
                    }
                },
                "appServiceConnectInbound": {
                    "applications": {
                        "accessType": "blocked",
                        "targets": [
                            {
                                "target": "AllApplications",
                                "targetType": "application"
                            }
                        ]
                    }
                }
            }
        },