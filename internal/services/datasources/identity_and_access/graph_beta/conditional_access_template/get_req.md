Request URL
https://graph.microsoft.com/beta/conditionalAccess/templates
Request Method
GET

{
    "@odata.context": "https://graph.microsoft.com/beta/$metadata#conditionalAccess/templates",
    "@odata.count": 21,
    "value": [
        {
            "name": "Require multifactor authentication for admins",
            "description": "Require multifactor authentication for privileged administrative accounts to reduce risk of compromise. This policy will target the same roles as security defaults.",
            "id": "c7503427-338e-4c5e-902d-abe252abfb43",
            "scenarios": "secureFoundation,zeroTrust,protectAdmins",
            "details": {
                "sessionControls": null,
                "conditions": {
                    "userRiskLevels": [],
                    "signInRiskLevels": [],
                    "clientAppTypes": [
                        "all"
                    ],
                    "servicePrincipalRiskLevels": [],
                    "agentIdRiskLevels": null,
                    "insiderRiskLevels": null,
                    "clients": null,
                    "platforms": null,
                    "locations": null,
                    "times": null,
                    "deviceStates": null,
                    "devices": null,
                    "clientApplications": null,
                    "authenticationFlows": null,
                    "applications": {
                        "includeApplications": [
                            "All"
                        ],
                        "excludeApplications": [],
                        "includeUserActions": [],
                        "includeAuthenticationContextClassReferences": [],
                        "applicationFilter": null,
                        "networkAccess": null,
                        "globalSecureAccess": null
                    },
                    "users": {
                        "includeUsers": [],
                        "excludeUsers": [
                            "Current administrator will be excluded"
                        ],
                        "includeGroups": [],
                        "excludeGroups": [],
                        "includeRoles": [
                            "62e90394-69f5-4237-9190-012177145e10",
                            "194ae4cb-b126-40b2-bd5b-6091b380977d",
                            "f28a1f50-f6e7-4571-818b-6a12f2af6b6c",
                            "29232cdf-9323-42fd-ade2-1d097af3e4de",
                            "b1be1c3e-b65d-4f19-8427-f6fa0d97feb9",
                            "729827e3-9c14-49f7-bb1b-9608f156bbb8",
                            "b0f54661-2d74-4c50-afa3-1ec803f12efe",
                            "fe930be7-5e62-47db-91af-98c3a49a38b1",
                            "c4e39bd9-1100-46d3-8c65-fb160da0071f",
                            "9b895d92-2cd3-44c7-9d02-a6ac2d5ea5c3",
                            "158c047a-c907-4556-b7ef-446551a6b5f7",
                            "966707d0-3269-4727-9be2-8c3a10f19b9d",
                            "7be44c8a-adaf-4e2a-84d6-ab2649e08a13",
                            "e8611ab8-c189-46e8-94e1-60213ab1f814"
                        ],
                        "excludeRoles": [],
                        "includeGuestsOrExternalUsers": null,
                        "excludeGuestsOrExternalUsers": null
                    }
                },
                "grantControls": {
                    "operator": "OR",
                    "builtInControls": [
                        "mfa"
                    ],
                    "customAuthenticationFactors": [],
                    "termsOfUse": [],
                    "authenticationStrength@odata.context": "https://graph.microsoft.com/beta/$metadata#conditionalAccess/templates('c7503427-338e-4c5e-902d-abe252abfb43')/details/grantControls/authenticationStrength/$entity",
                    "authenticationStrength": null
                }
            }
        },
        {
            "name": "Securing security info registration",
            "description": "Secure when and how users register for Azure AD multifactor authentication and self-service password reset.",
            "id": "b8bda7f8-6584-4446-bce9-d871480e53fa",
            "scenarios": "secureFoundation,zeroTrust,remoteWork",
            "details": {
                "sessionControls": null,
                "conditions": {
                    "userRiskLevels": [],
                    "signInRiskLevels": [],
                    "clientAppTypes": [
                        "all"
                    ],
                    "servicePrincipalRiskLevels": [],
                    "agentIdRiskLevels": null,
                    "insiderRiskLevels": null,
                    "clients": null,
                    "platforms": null,
                    "times": null,
                    "deviceStates": null,
                    "devices": null,
                    "clientApplications": null,
                    "authenticationFlows": null,
                    "applications": {
                        "includeApplications": [],
                        "excludeApplications": [],
                        "includeUserActions": [
                            "urn:user:registersecurityinfo"
                        ],
                        "includeAuthenticationContextClassReferences": [],
                        "applicationFilter": null,
                        "networkAccess": null,
                        "globalSecureAccess": null
                    },
                    "users": {
                        "includeUsers": [
                            "All"
                        ],
                        "excludeUsers": [
                            "GuestsOrExternalUsers",
                            "Current administrator will be excluded"
                        ],
                        "includeGroups": [],
                        "excludeGroups": [],
                        "includeRoles": [],
                        "excludeRoles": [
                            "62e90394-69f5-4237-9190-012177145e10"
                        ],
                        "includeGuestsOrExternalUsers": null,
                        "excludeGuestsOrExternalUsers": null
                    },
                    "locations": {
                        "includeLocations": [
                            "All"
                        ],
                        "excludeLocations": [
                            "AllTrusted"
                        ]
                    }
                },
                "grantControls": {
                    "operator": "OR",
                    "builtInControls": [
                        "mfa"
                    ],
                    "customAuthenticationFactors": [],
                    "termsOfUse": [],
                    "authenticationStrength@odata.context": "https://graph.microsoft.com/beta/$metadata#conditionalAccess/templates('b8bda7f8-6584-4446-bce9-d871480e53fa')/details/grantControls/authenticationStrength/$entity",
                    "authenticationStrength": null
                }
            }
        },
        {
            "name": "Block legacy authentication",
            "description": "Block legacy authentication endpoints that can be used to bypass multifactor authentication.",
            "id": "0b2282f9-2862-4178-88b5-d79340b36cb8",
            "scenarios": "secureFoundation,zeroTrust,remoteWork,protectAdmins",
            "details": {
                "sessionControls": null,
                "conditions": {
                    "userRiskLevels": [],
                    "signInRiskLevels": [],
                    "clientAppTypes": [
                        "exchangeActiveSync",
                        "other"
                    ],
                    "servicePrincipalRiskLevels": [],
                    "agentIdRiskLevels": null,
                    "insiderRiskLevels": null,
                    "clients": null,
                    "platforms": null,
                    "locations": null,
                    "times": null,
                    "deviceStates": null,
                    "devices": null,
                    "clientApplications": null,
                    "authenticationFlows": null,
                    "applications": {
                        "includeApplications": [
                            "All"
                        ],
                        "excludeApplications": [],
                        "includeUserActions": [],
                        "includeAuthenticationContextClassReferences": [],
                        "applicationFilter": null,
                        "networkAccess": null,
                        "globalSecureAccess": null
                    },
                    "users": {
                        "includeUsers": [
                            "All"
                        ],
                        "excludeUsers": [
                            "Current administrator will be excluded"
                        ],
                        "includeGroups": [],
                        "excludeGroups": [],
                        "includeRoles": [],
                        "excludeRoles": [],
                        "includeGuestsOrExternalUsers": null,
                        "excludeGuestsOrExternalUsers": null
                    }
                },
                "grantControls": {
                    "operator": "OR",
                    "builtInControls": [
                        "block"
                    ],
                    "customAuthenticationFactors": [],
                    "termsOfUse": [],
                    "authenticationStrength@odata.context": "https://graph.microsoft.com/beta/$metadata#conditionalAccess/templates('0b2282f9-2862-4178-88b5-d79340b36cb8')/details/grantControls/authenticationStrength/$entity",
                    "authenticationStrength": null
                }
            }
        },
        {
            "name": "Require multifactor authentication for all users",
            "description": "Require multifactor authentication for all user accounts to reduce risk of compromise. Directory Synchronization Accounts are excluded for on-premise directory synchronization tasks.",
            "id": "a3d0a415-b068-4326-9251-f9cdf9feeb64",
            "scenarios": "secureFoundation,zeroTrust,remoteWork",
            "details": {
                "sessionControls": null,
                "conditions": {
                    "userRiskLevels": [],
                    "signInRiskLevels": [],
                    "clientAppTypes": [
                        "all"
                    ],
                    "servicePrincipalRiskLevels": [],
                    "agentIdRiskLevels": null,
                    "insiderRiskLevels": null,
                    "clients": null,
                    "platforms": null,
                    "locations": null,
                    "times": null,
                    "deviceStates": null,
                    "devices": null,
                    "clientApplications": null,
                    "authenticationFlows": null,
                    "applications": {
                        "includeApplications": [
                            "All"
                        ],
                        "excludeApplications": [],
                        "includeUserActions": [],
                        "includeAuthenticationContextClassReferences": [],
                        "applicationFilter": null,
                        "networkAccess": null,
                        "globalSecureAccess": null
                    },
                    "users": {
                        "includeUsers": [
                            "All"
                        ],
                        "excludeUsers": [
                            "Current administrator will be excluded"
                        ],
                        "includeGroups": [],
                        "excludeGroups": [],
                        "includeRoles": [],
                        "excludeRoles": [
                            "d29b2b05-8046-44ba-8758-1e26182fcf32"
                        ],
                        "includeGuestsOrExternalUsers": null,
                        "excludeGuestsOrExternalUsers": null
                    }
                },
                "grantControls": {
                    "operator": "OR",
                    "builtInControls": [
                        "mfa"
                    ],
                    "customAuthenticationFactors": [],
                    "termsOfUse": [],
                    "authenticationStrength@odata.context": "https://graph.microsoft.com/beta/$metadata#conditionalAccess/templates('a3d0a415-b068-4326-9251-f9cdf9feeb64')/details/grantControls/authenticationStrength/$entity",
                    "authenticationStrength": null
                }
            }
        },
        {
            "name": "Require multifactor authentication for guest access",
            "description": "Require guest users perform multifactor authentication when accessing your company resources.",
            "id": "a4072ac0-722b-4991-981b-7f9755daef14",
            "scenarios": "zeroTrust,remoteWork",
            "details": {
                "sessionControls": null,
                "conditions": {
                    "userRiskLevels": [],
                    "signInRiskLevels": [],
                    "clientAppTypes": [
                        "all"
                    ],
                    "servicePrincipalRiskLevels": [],
                    "agentIdRiskLevels": null,
                    "insiderRiskLevels": null,
                    "clients": null,
                    "platforms": null,
                    "locations": null,
                    "times": null,
                    "deviceStates": null,
                    "devices": null,
                    "clientApplications": null,
                    "authenticationFlows": null,
                    "applications": {
                        "includeApplications": [
                            "All"
                        ],
                        "excludeApplications": [],
                        "includeUserActions": [],
                        "includeAuthenticationContextClassReferences": [],
                        "applicationFilter": null,
                        "networkAccess": null,
                        "globalSecureAccess": null
                    },
                    "users": {
                        "includeUsers": [
                            "GuestsOrExternalUsers"
                        ],
                        "excludeUsers": [
                            "Current administrator will be excluded"
                        ],
                        "includeGroups": [],
                        "excludeGroups": [],
                        "includeRoles": [],
                        "excludeRoles": [],
                        "includeGuestsOrExternalUsers": null,
                        "excludeGuestsOrExternalUsers": null
                    }
                },
                "grantControls": {
                    "operator": "OR",
                    "builtInControls": [
                        "mfa"
                    ],
                    "customAuthenticationFactors": [],
                    "termsOfUse": [],
                    "authenticationStrength@odata.context": "https://graph.microsoft.com/beta/$metadata#conditionalAccess/templates('a4072ac0-722b-4991-981b-7f9755daef14')/details/grantControls/authenticationStrength/$entity",
                    "authenticationStrength": null
                }
            }
        },
        {
            "name": "Require multifactor authentication for Azure management",
            "description": "Require multifactor authentication to protect privileged access to Azure management.",
            "id": "d8c51a9a-e6b1-454d-86af-554e7872e2c1",
            "scenarios": "secureFoundation,zeroTrust,protectAdmins",
            "details": {
                "sessionControls": null,
                "conditions": {
                    "userRiskLevels": [],
                    "signInRiskLevels": [],
                    "clientAppTypes": [
                        "all"
                    ],
                    "servicePrincipalRiskLevels": [],
                    "agentIdRiskLevels": null,
                    "insiderRiskLevels": null,
                    "clients": null,
                    "platforms": null,
                    "locations": null,
                    "times": null,
                    "deviceStates": null,
                    "devices": null,
                    "clientApplications": null,
                    "authenticationFlows": null,
                    "applications": {
                        "includeApplications": [
                            "797f4846-ba00-4fd7-ba43-dac1f8f63013"
                        ],
                        "excludeApplications": [],
                        "includeUserActions": [],
                        "includeAuthenticationContextClassReferences": [],
                        "applicationFilter": null,
                        "networkAccess": null,
                        "globalSecureAccess": null
                    },
                    "users": {
                        "includeUsers": [
                            "All"
                        ],
                        "excludeUsers": [
                            "Current administrator will be excluded"
                        ],
                        "includeGroups": [],
                        "excludeGroups": [],
                        "includeRoles": [],
                        "excludeRoles": [],
                        "includeGuestsOrExternalUsers": null,
                        "excludeGuestsOrExternalUsers": null
                    }
                },
                "grantControls": {
                    "operator": "OR",
                    "builtInControls": [
                        "mfa"
                    ],
                    "customAuthenticationFactors": [],
                    "termsOfUse": [],
                    "authenticationStrength@odata.context": "https://graph.microsoft.com/beta/$metadata#conditionalAccess/templates('d8c51a9a-e6b1-454d-86af-554e7872e2c1')/details/grantControls/authenticationStrength/$entity",
                    "authenticationStrength": null
                }
            }
        },
        {
            "name": "Require multifactor authentication for risky sign-ins",
            "description": "Require multifactor authentication if the sign-in risk is detected to be medium or high. (Requires a Microsoft Entra ID P2 license)",
            "id": "6b619f55-792e-45dc-9711-d83ec9d7ae90",
            "scenarios": "zeroTrust,remoteWork",
            "details": {
                "conditions": {
                    "userRiskLevels": [],
                    "signInRiskLevels": [
                        "high",
                        "medium"
                    ],
                    "clientAppTypes": [
                        "all"
                    ],
                    "servicePrincipalRiskLevels": [],
                    "agentIdRiskLevels": null,
                    "insiderRiskLevels": null,
                    "clients": null,
                    "platforms": null,
                    "locations": null,
                    "times": null,
                    "deviceStates": null,
                    "devices": null,
                    "clientApplications": null,
                    "authenticationFlows": null,
                    "applications": {
                        "includeApplications": [
                            "All"
                        ],
                        "excludeApplications": [],
                        "includeUserActions": [],
                        "includeAuthenticationContextClassReferences": [],
                        "applicationFilter": null,
                        "networkAccess": null,
                        "globalSecureAccess": null
                    },
                    "users": {
                        "includeUsers": [
                            "All"
                        ],
                        "excludeUsers": [
                            "Current administrator will be excluded"
                        ],
                        "includeGroups": [],
                        "excludeGroups": [],
                        "includeRoles": [],
                        "excludeRoles": [],
                        "includeGuestsOrExternalUsers": null,
                        "excludeGuestsOrExternalUsers": null
                    }
                },
                "grantControls": {
                    "operator": "OR",
                    "builtInControls": [
                        "mfa"
                    ],
                    "customAuthenticationFactors": [],
                    "termsOfUse": [],
                    "authenticationStrength@odata.context": "https://graph.microsoft.com/beta/$metadata#conditionalAccess/templates('6b619f55-792e-45dc-9711-d83ec9d7ae90')/details/grantControls/authenticationStrength/$entity",
                    "authenticationStrength": null
                },
                "sessionControls": {
                    "disableResilienceDefaults": null,
                    "applicationEnforcedRestrictions": null,
                    "cloudAppSecurity": null,
                    "persistentBrowser": null,
                    "continuousAccessEvaluation": null,
                    "secureSignInSession": null,
                    "networkAccessSecurity": null,
                    "globalSecureAccessFilteringProfile": null,
                    "signInFrequency": {
                        "value": null,
                        "type": null,
                        "authenticationType": "primaryAndSecondaryAuthentication",
                        "frequencyInterval": "everyTime",
                        "isEnabled": true
                    }
                }
            }
        },
        {
            "name": "Require password change for high-risk users",
            "description": "Require the user to change their password if the user risk is detected to be high. (Requires a Microsoft Entra ID P2 license)",
            "id": "634b6de7-c38d-4357-a2c7-3842706eedd7",
            "scenarios": "zeroTrust,remoteWork",
            "details": {
                "conditions": {
                    "userRiskLevels": [
                        "high"
                    ],
                    "signInRiskLevels": [],
                    "clientAppTypes": [
                        "all"
                    ],
                    "servicePrincipalRiskLevels": [],
                    "agentIdRiskLevels": null,
                    "insiderRiskLevels": null,
                    "clients": null,
                    "platforms": null,
                    "locations": null,
                    "times": null,
                    "deviceStates": null,
                    "devices": null,
                    "clientApplications": null,
                    "authenticationFlows": null,
                    "applications": {
                        "includeApplications": [
                            "All"
                        ],
                        "excludeApplications": [],
                        "includeUserActions": [],
                        "includeAuthenticationContextClassReferences": [],
                        "applicationFilter": null,
                        "networkAccess": null,
                        "globalSecureAccess": null
                    },
                    "users": {
                        "includeUsers": [
                            "All"
                        ],
                        "excludeUsers": [
                            "Current administrator will be excluded"
                        ],
                        "includeGroups": [],
                        "excludeGroups": [],
                        "includeRoles": [],
                        "excludeRoles": [],
                        "includeGuestsOrExternalUsers": null,
                        "excludeGuestsOrExternalUsers": null
                    }
                },
                "grantControls": {
                    "operator": "AND",
                    "builtInControls": [
                        "mfa",
                        "passwordChange"
                    ],
                    "customAuthenticationFactors": [],
                    "termsOfUse": [],
                    "authenticationStrength@odata.context": "https://graph.microsoft.com/beta/$metadata#conditionalAccess/templates('634b6de7-c38d-4357-a2c7-3842706eedd7')/details/grantControls/authenticationStrength/$entity",
                    "authenticationStrength": null
                },
                "sessionControls": {
                    "disableResilienceDefaults": null,
                    "applicationEnforcedRestrictions": null,
                    "cloudAppSecurity": null,
                    "persistentBrowser": null,
                    "continuousAccessEvaluation": null,
                    "secureSignInSession": null,
                    "networkAccessSecurity": null,
                    "globalSecureAccessFilteringProfile": null,
                    "signInFrequency": {
                        "value": null,
                        "type": null,
                        "authenticationType": "primaryAndSecondaryAuthentication",
                        "frequencyInterval": "everyTime",
                        "isEnabled": true
                    }
                }
            }
        },
        {
            "name": "Require compliant or hybrid Azure AD joined device for admins",
            "description": "Require privileged administrators to only access resources when using a compliant or hybrid Azure AD joined device.",
            "id": "c26a510a-3b8b-4023-8c44-d4f4c854e9f9",
            "scenarios": "remoteWork,protectAdmins",
            "details": {
                "sessionControls": null,
                "conditions": {
                    "userRiskLevels": [],
                    "signInRiskLevels": [],
                    "clientAppTypes": [
                        "all"
                    ],
                    "servicePrincipalRiskLevels": [],
                    "agentIdRiskLevels": null,
                    "insiderRiskLevels": null,
                    "clients": null,
                    "platforms": null,
                    "locations": null,
                    "times": null,
                    "deviceStates": null,
                    "devices": null,
                    "clientApplications": null,
                    "authenticationFlows": null,
                    "applications": {
                        "includeApplications": [
                            "All"
                        ],
                        "excludeApplications": [],
                        "includeUserActions": [],
                        "includeAuthenticationContextClassReferences": [],
                        "applicationFilter": null,
                        "networkAccess": null,
                        "globalSecureAccess": null
                    },
                    "users": {
                        "includeUsers": [
                            "None"
                        ],
                        "excludeUsers": [
                            "Current administrator will be excluded"
                        ],
                        "includeGroups": [],
                        "excludeGroups": [],
                        "includeRoles": [
                            "62e90394-69f5-4237-9190-012177145e10",
                            "194ae4cb-b126-40b2-bd5b-6091b380977d",
                            "f28a1f50-f6e7-4571-818b-6a12f2af6b6c",
                            "29232cdf-9323-42fd-ade2-1d097af3e4de",
                            "b1be1c3e-b65d-4f19-8427-f6fa0d97feb9",
                            "729827e3-9c14-49f7-bb1b-9608f156bbb8",
                            "b0f54661-2d74-4c50-afa3-1ec803f12efe",
                            "fe930be7-5e62-47db-91af-98c3a49a38b1",
                            "c4e39bd9-1100-46d3-8c65-fb160da0071f",
                            "9b895d92-2cd3-44c7-9d02-a6ac2d5ea5c3",
                            "158c047a-c907-4556-b7ef-446551a6b5f7",
                            "966707d0-3269-4727-9be2-8c3a10f19b9d",
                            "7be44c8a-adaf-4e2a-84d6-ab2649e08a13",
                            "e8611ab8-c189-46e8-94e1-60213ab1f814"
                        ],
                        "excludeRoles": [],
                        "includeGuestsOrExternalUsers": null,
                        "excludeGuestsOrExternalUsers": null
                    }
                },
                "grantControls": {
                    "operator": "OR",
                    "builtInControls": [
                        "compliantDevice",
                        "domainJoinedDevice"
                    ],
                    "customAuthenticationFactors": [],
                    "termsOfUse": [],
                    "authenticationStrength@odata.context": "https://graph.microsoft.com/beta/$metadata#conditionalAccess/templates('c26a510a-3b8b-4023-8c44-d4f4c854e9f9')/details/grantControls/authenticationStrength/$entity",
                    "authenticationStrength": null
                }
            }
        },
        {
            "name": "Block access for unknown or unsupported device platform",
            "description": "Users will be blocked from accessing company resources when the device type is unknown or unsupported.",
            "id": "4e39a309-931e-4cb1-a371-e2beea168002",
            "scenarios": "zeroTrust,remoteWork",
            "details": {
                "sessionControls": null,
                "conditions": {
                    "userRiskLevels": [],
                    "signInRiskLevels": [],
                    "clientAppTypes": [
                        "all"
                    ],
                    "servicePrincipalRiskLevels": [],
                    "agentIdRiskLevels": null,
                    "insiderRiskLevels": null,
                    "clients": null,
                    "locations": null,
                    "times": null,
                    "deviceStates": null,
                    "devices": null,
                    "clientApplications": null,
                    "authenticationFlows": null,
                    "applications": {
                        "includeApplications": [
                            "All"
                        ],
                        "excludeApplications": [],
                        "includeUserActions": [],
                        "includeAuthenticationContextClassReferences": [],
                        "applicationFilter": null,
                        "networkAccess": null,
                        "globalSecureAccess": null
                    },
                    "users": {
                        "includeUsers": [
                            "All"
                        ],
                        "excludeUsers": [
                            "Current administrator will be excluded"
                        ],
                        "includeGroups": [],
                        "excludeGroups": [],
                        "includeRoles": [],
                        "excludeRoles": [],
                        "includeGuestsOrExternalUsers": null,
                        "excludeGuestsOrExternalUsers": null
                    },
                    "platforms": {
                        "includePlatforms": [
                            "all"
                        ],
                        "excludePlatforms": [
                            "android",
                            "iOS",
                            "windows",
                            "macOS",
                            "linux",
                            "windowsPhone"
                        ]
                    }
                },
                "grantControls": {
                    "operator": "OR",
                    "builtInControls": [
                        "block"
                    ],
                    "customAuthenticationFactors": [],
                    "termsOfUse": [],
                    "authenticationStrength@odata.context": "https://graph.microsoft.com/beta/$metadata#conditionalAccess/templates('4e39a309-931e-4cb1-a371-e2beea168002')/details/grantControls/authenticationStrength/$entity",
                    "authenticationStrength": null
                }
            }
        },
        {
            "name": "No persistent browser session",
            "description": "Protect user access on unmanaged devices by preventing browser sessions from remaining signed in after the browser is closed and setting a sign-in frequency to 1 hour.",
            "id": "62e51ccc-c9c3-4554-ac70-066172c81007",
            "scenarios": "zeroTrust,remoteWork",
            "details": {
                "grantControls": null,
                "conditions": {
                    "userRiskLevels": [],
                    "signInRiskLevels": [],
                    "clientAppTypes": [
                        "all"
                    ],
                    "servicePrincipalRiskLevels": [],
                    "agentIdRiskLevels": null,
                    "insiderRiskLevels": null,
                    "clients": null,
                    "platforms": null,
                    "locations": null,
                    "times": null,
                    "deviceStates": null,
                    "clientApplications": null,
                    "authenticationFlows": null,
                    "applications": {
                        "includeApplications": [
                            "All"
                        ],
                        "excludeApplications": [],
                        "includeUserActions": [],
                        "includeAuthenticationContextClassReferences": [],
                        "applicationFilter": null,
                        "networkAccess": null,
                        "globalSecureAccess": null
                    },
                    "users": {
                        "includeUsers": [
                            "All"
                        ],
                        "excludeUsers": [
                            "Current administrator will be excluded"
                        ],
                        "includeGroups": [],
                        "excludeGroups": [],
                        "includeRoles": [],
                        "excludeRoles": [],
                        "includeGuestsOrExternalUsers": null,
                        "excludeGuestsOrExternalUsers": null
                    },
                    "devices": {
                        "includeDeviceStates": [],
                        "excludeDeviceStates": [],
                        "includeDevices": [],
                        "excludeDevices": [],
                        "deviceFilter": {
                            "mode": "include",
                            "rule": "device.trustType -ne \"ServerAD\" -or device.isCompliant -ne True"
                        }
                    }
                },
                "sessionControls": {
                    "disableResilienceDefaults": null,
                    "applicationEnforcedRestrictions": null,
                    "cloudAppSecurity": null,
                    "continuousAccessEvaluation": null,
                    "secureSignInSession": null,
                    "networkAccessSecurity": null,
                    "globalSecureAccessFilteringProfile": null,
                    "signInFrequency": {
                        "value": 1,
                        "type": "hours",
                        "authenticationType": "primaryAndSecondaryAuthentication",
                        "frequencyInterval": "timeBased",
                        "isEnabled": true
                    },
                    "persistentBrowser": {
                        "mode": "never",
                        "isEnabled": true
                    }
                }
            }
        },
        {
            "name": "Require compliant or hybrid Azure AD joined device or multifactor authentication for all users",
            "description": "Protect access to company resources by requiring users to use a managed device or perform multifactor authentication. Directory Synchronization Accounts are excluded for on-premise directory synchronization tasks.",
            "id": "927c884e-7888-4e81-abc4-bd56ded28985",
            "scenarios": "secureFoundation,zeroTrust",
            "details": {
                "sessionControls": null,
                "conditions": {
                    "userRiskLevels": [],
                    "signInRiskLevels": [],
                    "clientAppTypes": [
                        "all"
                    ],
                    "servicePrincipalRiskLevels": [],
                    "agentIdRiskLevels": null,
                    "insiderRiskLevels": null,
                    "clients": null,
                    "platforms": null,
                    "locations": null,
                    "times": null,
                    "deviceStates": null,
                    "devices": null,
                    "clientApplications": null,
                    "authenticationFlows": null,
                    "applications": {
                        "includeApplications": [
                            "All"
                        ],
                        "excludeApplications": [],
                        "includeUserActions": [],
                        "includeAuthenticationContextClassReferences": [],
                        "applicationFilter": null,
                        "networkAccess": null,
                        "globalSecureAccess": null
                    },
                    "users": {
                        "includeUsers": [
                            "All"
                        ],
                        "excludeUsers": [
                            "Current administrator will be excluded"
                        ],
                        "includeGroups": [],
                        "excludeGroups": [],
                        "includeRoles": [],
                        "excludeRoles": [
                            "d29b2b05-8046-44ba-8758-1e26182fcf32"
                        ],
                        "includeGuestsOrExternalUsers": null,
                        "excludeGuestsOrExternalUsers": null
                    }
                },
                "grantControls": {
                    "operator": "OR",
                    "builtInControls": [
                        "mfa",
                        "compliantDevice",
                        "domainJoinedDevice"
                    ],
                    "customAuthenticationFactors": [],
                    "termsOfUse": [],
                    "authenticationStrength@odata.context": "https://graph.microsoft.com/beta/$metadata#conditionalAccess/templates('927c884e-7888-4e81-abc4-bd56ded28985')/details/grantControls/authenticationStrength/$entity",
                    "authenticationStrength": null
                }
            }
        },
        {
            "name": "Use application enforced restrictions for O365 apps",
            "description": "Block or limit access to O365 apps, including SharePoint Online, OneDrive, and Exchange Online content. This policy requires SharePoint admin center configuration.",
            "id": "81fd2072-4876-42b6-8157-c6000693046b",
            "scenarios": "remoteWork",
            "details": {
                "grantControls": null,
                "conditions": {
                    "userRiskLevels": [],
                    "signInRiskLevels": [],
                    "clientAppTypes": [
                        "all"
                    ],
                    "servicePrincipalRiskLevels": [],
                    "agentIdRiskLevels": null,
                    "insiderRiskLevels": null,
                    "clients": null,
                    "platforms": null,
                    "locations": null,
                    "times": null,
                    "deviceStates": null,
                    "devices": null,
                    "clientApplications": null,
                    "authenticationFlows": null,
                    "applications": {
                        "includeApplications": [
                            "Office365"
                        ],
                        "excludeApplications": [],
                        "includeUserActions": [],
                        "includeAuthenticationContextClassReferences": [],
                        "applicationFilter": null,
                        "networkAccess": null,
                        "globalSecureAccess": null
                    },
                    "users": {
                        "includeUsers": [
                            "All"
                        ],
                        "excludeUsers": [
                            "Current administrator will be excluded"
                        ],
                        "includeGroups": [],
                        "excludeGroups": [],
                        "includeRoles": [],
                        "excludeRoles": [],
                        "includeGuestsOrExternalUsers": null,
                        "excludeGuestsOrExternalUsers": null
                    }
                },
                "sessionControls": {
                    "disableResilienceDefaults": null,
                    "cloudAppSecurity": null,
                    "signInFrequency": null,
                    "persistentBrowser": null,
                    "continuousAccessEvaluation": null,
                    "secureSignInSession": null,
                    "networkAccessSecurity": null,
                    "globalSecureAccessFilteringProfile": null,
                    "applicationEnforcedRestrictions": {
                        "isEnabled": true
                    }
                }
            }
        },
        {
            "name": "Require phishing-resistant multifactor authentication for admins",
            "description": "Require phishing-resistant multifactor authentication for privileged administrative accounts to reduce risk of compromise and phishing attacks. This policy requires admins to have at least one phishing resistant authentication method registered.",
            "id": "76c03f19-ea37-4656-a772-a183b4ddb81d",
            "scenarios": "protectAdmins,emergingThreats",
            "details": {
                "sessionControls": null,
                "conditions": {
                    "userRiskLevels": [],
                    "signInRiskLevels": [],
                    "clientAppTypes": [
                        "all"
                    ],
                    "servicePrincipalRiskLevels": [],
                    "agentIdRiskLevels": null,
                    "insiderRiskLevels": null,
                    "clients": null,
                    "platforms": null,
                    "locations": null,
                    "times": null,
                    "deviceStates": null,
                    "devices": null,
                    "clientApplications": null,
                    "authenticationFlows": null,
                    "applications": {
                        "includeApplications": [
                            "All"
                        ],
                        "excludeApplications": [],
                        "includeUserActions": [],
                        "includeAuthenticationContextClassReferences": [],
                        "applicationFilter": null,
                        "networkAccess": null,
                        "globalSecureAccess": null
                    },
                    "users": {
                        "includeUsers": [],
                        "excludeUsers": [
                            "Current administrator will be excluded"
                        ],
                        "includeGroups": [],
                        "excludeGroups": [],
                        "includeRoles": [
                            "62e90394-69f5-4237-9190-012177145e10",
                            "194ae4cb-b126-40b2-bd5b-6091b380977d",
                            "f28a1f50-f6e7-4571-818b-6a12f2af6b6c",
                            "29232cdf-9323-42fd-ade2-1d097af3e4de",
                            "b1be1c3e-b65d-4f19-8427-f6fa0d97feb9",
                            "729827e3-9c14-49f7-bb1b-9608f156bbb8",
                            "b0f54661-2d74-4c50-afa3-1ec803f12efe",
                            "fe930be7-5e62-47db-91af-98c3a49a38b1",
                            "c4e39bd9-1100-46d3-8c65-fb160da0071f",
                            "9b895d92-2cd3-44c7-9d02-a6ac2d5ea5c3",
                            "158c047a-c907-4556-b7ef-446551a6b5f7",
                            "966707d0-3269-4727-9be2-8c3a10f19b9d",
                            "7be44c8a-adaf-4e2a-84d6-ab2649e08a13",
                            "e8611ab8-c189-46e8-94e1-60213ab1f814",
                            "17315797-102d-40b4-93e0-432062caca18",
                            "e6d1a23a-da11-4be4-9570-befc86d067a7",
                            "3a2c62db-5318-420d-8d74-23affee5d9d5",
                            "44367163-eba1-44c3-98af-f5787879f96a",
                            "11648597-926c-4cf3-9c36-bcebb0ba8dcc"
                        ],
                        "excludeRoles": [],
                        "includeGuestsOrExternalUsers": null,
                        "excludeGuestsOrExternalUsers": null
                    }
                },
                "grantControls": {
                    "operator": "AND",
                    "builtInControls": [],
                    "customAuthenticationFactors": [],
                    "termsOfUse": [],
                    "authenticationStrength@odata.context": "https://graph.microsoft.com/beta/$metadata#conditionalAccess/templates('76c03f19-ea37-4656-a772-a183b4ddb81d')/details/grantControls/authenticationStrength/$entity",
                    "authenticationStrength": {
                        "id": "00000000-0000-0000-0000-000000000004",
                        "createdDateTime": "2021-12-01T08:00:00Z",
                        "modifiedDateTime": "2021-12-01T08:00:00Z",
                        "displayName": "Phishing-resistant MFA",
                        "description": "Phishing-resistant, Passwordless methods for the strongest authentication, such as a FIDO2 security key",
                        "policyType": "builtIn",
                        "requirementsSatisfied": "mfa",
                        "allowedCombinations": [
                            "windowsHelloForBusiness",
                            "fido2",
                            "x509CertificateMultiFactor"
                        ],
                        "combinationConfigurations@odata.context": "https://graph.microsoft.com/beta/$metadata#conditionalAccess/templates('76c03f19-ea37-4656-a772-a183b4ddb81d')/details/grantControls/authenticationStrength/combinationConfigurations",
                        "combinationConfigurations": []
                    }
                }
            }
        },
        {
            "name": "Require multifactor authentication for Microsoft admin portals",
            "description": "Use this template to protect sign-ins to admin portals if you are unable to use the \"Require MFA for admins\" template.",
            "id": "6364131e-bc4a-47c4-a20b-33492d1fff6c",
            "scenarios": "zeroTrust,protectAdmins",
            "details": {
                "sessionControls": null,
                "conditions": {
                    "userRiskLevels": [],
                    "signInRiskLevels": [],
                    "clientAppTypes": [
                        "all"
                    ],
                    "servicePrincipalRiskLevels": [],
                    "agentIdRiskLevels": null,
                    "insiderRiskLevels": null,
                    "clients": null,
                    "platforms": null,
                    "locations": null,
                    "times": null,
                    "deviceStates": null,
                    "devices": null,
                    "clientApplications": null,
                    "authenticationFlows": null,
                    "applications": {
                        "includeApplications": [
                            "MicrosoftAdminPortals"
                        ],
                        "excludeApplications": [],
                        "includeUserActions": [],
                        "includeAuthenticationContextClassReferences": [],
                        "applicationFilter": null,
                        "networkAccess": null,
                        "globalSecureAccess": null
                    },
                    "users": {
                        "includeUsers": [],
                        "excludeUsers": [
                            "Current administrator will be excluded"
                        ],
                        "includeGroups": [],
                        "excludeGroups": [],
                        "includeRoles": [
                            "62e90394-69f5-4237-9190-012177145e10",
                            "194ae4cb-b126-40b2-bd5b-6091b380977d",
                            "f28a1f50-f6e7-4571-818b-6a12f2af6b6c",
                            "29232cdf-9323-42fd-ade2-1d097af3e4de",
                            "b1be1c3e-b65d-4f19-8427-f6fa0d97feb9",
                            "729827e3-9c14-49f7-bb1b-9608f156bbb8",
                            "b0f54661-2d74-4c50-afa3-1ec803f12efe",
                            "fe930be7-5e62-47db-91af-98c3a49a38b1",
                            "c4e39bd9-1100-46d3-8c65-fb160da0071f",
                            "9b895d92-2cd3-44c7-9d02-a6ac2d5ea5c3",
                            "158c047a-c907-4556-b7ef-446551a6b5f7",
                            "966707d0-3269-4727-9be2-8c3a10f19b9d",
                            "7be44c8a-adaf-4e2a-84d6-ab2649e08a13",
                            "e8611ab8-c189-46e8-94e1-60213ab1f814"
                        ],
                        "excludeRoles": [],
                        "includeGuestsOrExternalUsers": null,
                        "excludeGuestsOrExternalUsers": null
                    }
                },
                "grantControls": {
                    "operator": "OR",
                    "builtInControls": [],
                    "customAuthenticationFactors": [],
                    "termsOfUse": [],
                    "authenticationStrength@odata.context": "https://graph.microsoft.com/beta/$metadata#conditionalAccess/templates('6364131e-bc4a-47c4-a20b-33492d1fff6c')/details/grantControls/authenticationStrength/$entity",
                    "authenticationStrength": {
                        "id": "00000000-0000-0000-0000-000000000002",
                        "createdDateTime": "2021-12-01T08:00:00Z",
                        "modifiedDateTime": "2021-12-01T08:00:00Z",
                        "displayName": "Multifactor authentication",
                        "description": "Combinations of methods that satisfy strong authentication, such as a password + SMS",
                        "policyType": "builtIn",
                        "requirementsSatisfied": "mfa",
                        "allowedCombinations": [
                            "windowsHelloForBusiness",
                            "fido2",
                            "x509CertificateMultiFactor",
                            "deviceBasedPush",
                            "temporaryAccessPassOneTime",
                            "temporaryAccessPassMultiUse",
                            "password,microsoftAuthenticatorPush",
                            "password,softwareOath",
                            "password,hardwareOath",
                            "password,x509CertificateSingleFactor",
                            "password,x509CertificateMultiFactor",
                            "password,sms",
                            "password,voice",
                            "federatedMultiFactor",
                            "microsoftAuthenticatorPush,federatedSingleFactor",
                            "softwareOath,federatedSingleFactor",
                            "hardwareOath,federatedSingleFactor",
                            "sms,federatedSingleFactor",
                            "voice,federatedSingleFactor"
                        ],
                        "combinationConfigurations@odata.context": "https://graph.microsoft.com/beta/$metadata#conditionalAccess/templates('6364131e-bc4a-47c4-a20b-33492d1fff6c')/details/grantControls/authenticationStrength/combinationConfigurations",
                        "combinationConfigurations": []
                    }
                }
            }
        },
        {
            "name": "Block access to Office365 apps for users with insider risk",
            "description": "Configure insider risk as a condition to identify potential risky behavior (Requires a Microsoft Entra ID P2 license).",
            "id": "16aaa400-bfdf-4756-a420-ad2245d4cde8",
            "scenarios": "zeroTrust",
            "details": {
                "sessionControls": null,
                "conditions": {
                    "userRiskLevels": [],
                    "signInRiskLevels": [],
                    "clientAppTypes": [
                        "all"
                    ],
                    "servicePrincipalRiskLevels": [],
                    "agentIdRiskLevels": null,
                    "insiderRiskLevels": "elevated",
                    "clients": null,
                    "platforms": null,
                    "locations": null,
                    "times": null,
                    "deviceStates": null,
                    "devices": null,
                    "clientApplications": null,
                    "authenticationFlows": null,
                    "applications": {
                        "includeApplications": [
                            "Office365"
                        ],
                        "excludeApplications": [],
                        "includeUserActions": [],
                        "includeAuthenticationContextClassReferences": [],
                        "applicationFilter": null,
                        "networkAccess": null,
                        "globalSecureAccess": null
                    },
                    "users": {
                        "includeUsers": [
                            "All"
                        ],
                        "excludeUsers": [
                            "Current administrator will be excluded"
                        ],
                        "includeGroups": [],
                        "excludeGroups": [],
                        "includeRoles": [],
                        "excludeRoles": [],
                        "includeGuestsOrExternalUsers": null,
                        "excludeGuestsOrExternalUsers": {
                            "guestOrExternalUserTypes": "b2bDirectConnectUser,otherExternalUser,serviceProvider",
                            "externalTenants": null
                        }
                    }
                },
                "grantControls": {
                    "operator": "OR",
                    "builtInControls": [
                        "block"
                    ],
                    "customAuthenticationFactors": [],
                    "termsOfUse": [],
                    "authenticationStrength@odata.context": "https://graph.microsoft.com/beta/$metadata#conditionalAccess/templates('16aaa400-bfdf-4756-a420-ad2245d4cde8')/details/grantControls/authenticationStrength/$entity",
                    "authenticationStrength": null
                }
            }
        },
        {
            "name": "Require MDM-enrolled and compliant device to access cloud apps for all users (Preview)",
            "description": "Require devices to be enrolled in mobile device management (MDM) and be compliant for all users and devices accessing company resources. This improves data security by reducing risks of breaches, malware, and unauthorized access. Directory Synchronization Accounts are excluded for on-premise directory synchronization tasks.",
            "id": "a297dd1a-21fe-4016-99a0-ba43ba64378c",
            "scenarios": "secureFoundation,zeroTrust",
            "details": {
                "sessionControls": null,
                "conditions": {
                    "userRiskLevels": [],
                    "signInRiskLevels": [],
                    "clientAppTypes": [
                        "all"
                    ],
                    "servicePrincipalRiskLevels": [],
                    "agentIdRiskLevels": null,
                    "insiderRiskLevels": null,
                    "clients": null,
                    "platforms": null,
                    "locations": null,
                    "times": null,
                    "deviceStates": null,
                    "devices": null,
                    "clientApplications": null,
                    "authenticationFlows": null,
                    "applications": {
                        "includeApplications": [
                            "All"
                        ],
                        "excludeApplications": [],
                        "includeUserActions": [],
                        "includeAuthenticationContextClassReferences": [],
                        "applicationFilter": null,
                        "networkAccess": null,
                        "globalSecureAccess": null
                    },
                    "users": {
                        "includeUsers": [
                            "All"
                        ],
                        "excludeUsers": [
                            "Current administrator will be excluded"
                        ],
                        "includeGroups": [],
                        "excludeGroups": [],
                        "includeRoles": [],
                        "excludeRoles": [
                            "d29b2b05-8046-44ba-8758-1e26182fcf32"
                        ],
                        "includeGuestsOrExternalUsers": null,
                        "excludeGuestsOrExternalUsers": null
                    }
                },
                "grantControls": {
                    "operator": "OR",
                    "builtInControls": [
                        "compliantDevice"
                    ],
                    "customAuthenticationFactors": [],
                    "termsOfUse": [],
                    "authenticationStrength@odata.context": "https://graph.microsoft.com/beta/$metadata#conditionalAccess/templates('a297dd1a-21fe-4016-99a0-ba43ba64378c')/details/grantControls/authenticationStrength/$entity",
                    "authenticationStrength": null
                }
            }
        },
        {
            "name": "Secure account recovery with identity verification (Preview)",
            "description": "Secure self-service account recovery by requiring users to verify their real-world identity with Microsoft Entra Verified ID and Face Check.",
            "id": "6acdf4c3-6815-485c-a57d-2c349d517ba0",
            "scenarios": "secureFoundation,zeroTrust,remoteWork,emergingThreats",
            "details": {
                "sessionControls": null,
                "conditions": {
                    "userRiskLevels": [],
                    "signInRiskLevels": [],
                    "clientAppTypes": [
                        "all"
                    ],
                    "servicePrincipalRiskLevels": [],
                    "agentIdRiskLevels": null,
                    "insiderRiskLevels": null,
                    "clients": null,
                    "platforms": null,
                    "locations": null,
                    "times": null,
                    "deviceStates": null,
                    "devices": null,
                    "clientApplications": null,
                    "authenticationFlows": null,
                    "applications": {
                        "includeApplications": [],
                        "excludeApplications": [],
                        "includeUserActions": [
                            "urn:user:accountrecovery"
                        ],
                        "includeAuthenticationContextClassReferences": [],
                        "applicationFilter": null,
                        "networkAccess": null,
                        "globalSecureAccess": null
                    },
                    "users": {
                        "includeUsers": [
                            "All"
                        ],
                        "excludeUsers": [
                            "Current administrator will be excluded"
                        ],
                        "includeGroups": [],
                        "excludeGroups": [],
                        "includeRoles": [],
                        "excludeRoles": [],
                        "includeGuestsOrExternalUsers": null,
                        "excludeGuestsOrExternalUsers": {
                            "guestOrExternalUserTypes": "b2bCollaborationGuest,b2bCollaborationMember,b2bDirectConnectUser,otherExternalUser,serviceProvider",
                            "externalTenants": {
                                "@odata.type": "#microsoft.graph.conditionalAccessAllExternalTenants",
                                "membershipKind": "all"
                            }
                        }
                    }
                },
                "grantControls": {
                    "operator": "AND",
                    "builtInControls": [
                        "verifiedID"
                    ],
                    "customAuthenticationFactors": [],
                    "termsOfUse": [],
                    "authenticationStrength@odata.context": "https://graph.microsoft.com/beta/$metadata#conditionalAccess/templates('6acdf4c3-6815-485c-a57d-2c349d517ba0')/details/grantControls/authenticationStrength/$entity",
                    "authenticationStrength": null
                }
            }
        },
        {
            "name": "Block high risk agent identities from accessing resources",
            "description": "This policy blocks agent identities with a high risk level from accessing resources in your tenant.",
            "id": "6bdcbbb7-ebd1-4f8f-a02f-652db0e3665d",
            "scenarios": "aiAgents",
            "details": {
                "sessionControls": null,
                "conditions": {
                    "userRiskLevels": [],
                    "signInRiskLevels": [],
                    "clientAppTypes": [
                        "all"
                    ],
                    "servicePrincipalRiskLevels": [],
                    "agentIdRiskLevels": "high",
                    "insiderRiskLevels": null,
                    "clients": null,
                    "platforms": null,
                    "locations": null,
                    "times": null,
                    "deviceStates": null,
                    "devices": null,
                    "authenticationFlows": null,
                    "applications": {
                        "includeApplications": [
                            "All"
                        ],
                        "excludeApplications": [],
                        "includeUserActions": [],
                        "includeAuthenticationContextClassReferences": [],
                        "applicationFilter": null,
                        "networkAccess": null,
                        "globalSecureAccess": null
                    },
                    "users": {
                        "includeUsers": [
                            "None"
                        ],
                        "excludeUsers": [],
                        "includeGroups": [],
                        "excludeGroups": [],
                        "includeRoles": [],
                        "excludeRoles": [],
                        "includeGuestsOrExternalUsers": null,
                        "excludeGuestsOrExternalUsers": null
                    },
                    "clientApplications": {
                        "includeServicePrincipals": [],
                        "includeAgentIdServicePrincipals": [
                            "All"
                        ],
                        "excludeServicePrincipals": [],
                        "excludeAgentIdServicePrincipals": [],
                        "servicePrincipalFilter": null,
                        "agentIdServicePrincipalFilter": null
                    }
                },
                "grantControls": {
                    "operator": "OR",
                    "builtInControls": [
                        "block"
                    ],
                    "customAuthenticationFactors": [],
                    "termsOfUse": [],
                    "authenticationStrength@odata.context": "https://graph.microsoft.com/beta/$metadata#conditionalAccess/templates('6bdcbbb7-ebd1-4f8f-a02f-652db0e3665d')/details/grantControls/authenticationStrength/$entity",
                    "authenticationStrength": null
                }
            }
        },
        {
            "name": "Block all agent identities from accessing resources",
            "description": "This policy blocks all agent identities from accessing resources in your tenant. You can update this policy to exclude specific agents.",
            "id": "c5bf9137-a43c-48e1-9d13-8258b31b855d",
            "scenarios": "aiAgents",
            "details": {
                "sessionControls": null,
                "conditions": {
                    "userRiskLevels": [],
                    "signInRiskLevels": [],
                    "clientAppTypes": [
                        "all"
                    ],
                    "servicePrincipalRiskLevels": [],
                    "agentIdRiskLevels": null,
                    "insiderRiskLevels": null,
                    "clients": null,
                    "platforms": null,
                    "locations": null,
                    "times": null,
                    "deviceStates": null,
                    "devices": null,
                    "authenticationFlows": null,
                    "applications": {
                        "includeApplications": [
                            "All"
                        ],
                        "excludeApplications": [],
                        "includeUserActions": [],
                        "includeAuthenticationContextClassReferences": [],
                        "applicationFilter": null,
                        "networkAccess": null,
                        "globalSecureAccess": null
                    },
                    "users": {
                        "includeUsers": [
                            "None"
                        ],
                        "excludeUsers": [],
                        "includeGroups": [],
                        "excludeGroups": [],
                        "includeRoles": [],
                        "excludeRoles": [],
                        "includeGuestsOrExternalUsers": null,
                        "excludeGuestsOrExternalUsers": null
                    },
                    "clientApplications": {
                        "includeServicePrincipals": [],
                        "includeAgentIdServicePrincipals": [
                            "All"
                        ],
                        "excludeServicePrincipals": [],
                        "excludeAgentIdServicePrincipals": [],
                        "servicePrincipalFilter": null,
                        "agentIdServicePrincipalFilter": null
                    }
                },
                "grantControls": {
                    "operator": "OR",
                    "builtInControls": [
                        "block"
                    ],
                    "customAuthenticationFactors": [],
                    "termsOfUse": [],
                    "authenticationStrength@odata.context": "https://graph.microsoft.com/beta/$metadata#conditionalAccess/templates('c5bf9137-a43c-48e1-9d13-8258b31b855d')/details/grantControls/authenticationStrength/$entity",
                    "authenticationStrength": null
                }
            }
        },
        {
            "name": "Block all agent users from accessing resources",
            "description": "This policy blocks all agent users from accessing resources in your tenant. Note that this policy does not support any exclusions during public preview.",
            "id": "376d9923-7f15-40c2-830b-3312120f1eed",
            "scenarios": "aiAgents",
            "details": {
                "sessionControls": null,
                "conditions": {
                    "userRiskLevels": [],
                    "signInRiskLevels": [],
                    "clientAppTypes": [
                        "all"
                    ],
                    "servicePrincipalRiskLevels": [],
                    "agentIdRiskLevels": null,
                    "insiderRiskLevels": null,
                    "clients": null,
                    "platforms": null,
                    "locations": null,
                    "times": null,
                    "deviceStates": null,
                    "devices": null,
                    "clientApplications": null,
                    "authenticationFlows": null,
                    "applications": {
                        "includeApplications": [
                            "All"
                        ],
                        "excludeApplications": [],
                        "includeUserActions": [],
                        "includeAuthenticationContextClassReferences": [],
                        "applicationFilter": null,
                        "networkAccess": null,
                        "globalSecureAccess": null
                    },
                    "users": {
                        "includeUsers": [
                            "AllAgentIdUsers"
                        ],
                        "excludeUsers": [],
                        "includeGroups": [],
                        "excludeGroups": [],
                        "includeRoles": [],
                        "excludeRoles": [],
                        "includeGuestsOrExternalUsers": null,
                        "excludeGuestsOrExternalUsers": null
                    }
                },
                "grantControls": {
                    "operator": "OR",
                    "builtInControls": [
                        "block"
                    ],
                    "customAuthenticationFactors": [],
                    "termsOfUse": [],
                    "authenticationStrength@odata.context": "https://graph.microsoft.com/beta/$metadata#conditionalAccess/templates('376d9923-7f15-40c2-830b-3312120f1eed')/details/grantControls/authenticationStrength/$entity",
                    "authenticationStrength": null
                }
            }
        }
    ]
}