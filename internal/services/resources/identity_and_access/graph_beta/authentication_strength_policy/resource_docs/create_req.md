Request URL
https://graph.microsoft.com/v1.0/$batch
Request Method
POST

{"requests":[{"id":"0178675b-e250-46a9-babd-b92881217c0b","method":"GET","url":"/organization/2fd6bb84-ad40-4ec5-9369-a215b25c9952/certificateBasedAuthConfiguration","headers":{"x-ms-command-name":"AuthenticationStrengths - getCertIssuersForAuthStrengths","x-ms-client-request-id":"a6cd687c-233d-4078-8654-232d3b0e5c0f","client-request-id":"a6cd687c-233d-4078-8654-232d3b0e5c0f","x-ms-client-session-id":"4bdab402b3f54080b4890147387920a3"}}]}


Request URL
https://graph.microsoft.com/beta/$batch
Request Method
POST

{"requests":[{"id":"f628f58b-dd22-4128-bee0-3acb96537864","method":"POST","url":"/identity/conditionalAccess/authenticationStrength/policies/","body":{"displayName":"gui-test","description":"gui-test","allowedCombinations":["windowsHelloForBusiness"],"combinationConfigurations":[]},"headers":{"x-ms-command-name":"AuthenticationStrengths - AddCustomAuthStrength","x-ms-client-request-id":"5e21e5d7-c4d8-4d3e-aea5-1b5f2c2c09a8","client-request-id":"5e21e5d7-c4d8-4d3e-aea5-1b5f2c2c09a8","x-ms-client-session-id":"4bdab402b3f54080b4890147387920a3","Content-Type":"application/json"}}]}

Request URL
https://graph.microsoft.com/beta/identity/conditionalAccess/authenticationStrength/policies/?
Request Method
GET

{
    "@odata.context": "https://graph.microsoft.com/beta/$metadata#identity/conditionalAccess/authenticationStrength/policies",
    "value": [
        {
            "id": "8484189b-f18b-4aea-8401-737bf3069677",
            "createdDateTime": "2025-08-26T10:05:38.43381Z",
            "modifiedDateTime": "2025-08-26T10:05:38.4348109Z",
            "displayName": "gui-test",
            "description": "gui-test",
            "policyType": "custom",
            "requirementsSatisfied": "mfa",
            "allowedCombinations": [
                "windowsHelloForBusiness"
            ],
            "combinationConfigurations@odata.context": "https://graph.microsoft.com/beta/$metadata#identity/conditionalAccess/authenticationStrength/policies('8484189b-f18b-4aea-8401-737bf3069677')/combinationConfigurations",
            "combinationConfigurations": []
        },
        {
            "id": "00000000-0000-0000-0000-000000000002",
            "createdDateTime": "2021-12-01T00:00:00Z",
            "modifiedDateTime": "2021-12-01T00:00:00Z",
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
                "password,sms",
                "password,voice",
                "federatedMultiFactor",
                "microsoftAuthenticatorPush,federatedSingleFactor",
                "softwareOath,federatedSingleFactor",
                "hardwareOath,federatedSingleFactor",
                "sms,federatedSingleFactor",
                "voice,federatedSingleFactor"
            ],
            "combinationConfigurations@odata.context": "https://graph.microsoft.com/beta/$metadata#identity/conditionalAccess/authenticationStrength/policies('00000000-0000-0000-0000-000000000002')/combinationConfigurations",
            "combinationConfigurations": []
        },
        {
            "id": "00000000-0000-0000-0000-000000000003",
            "createdDateTime": "2021-12-01T00:00:00Z",
            "modifiedDateTime": "2021-12-01T00:00:00Z",
            "displayName": "Passwordless MFA",
            "description": "Passwordless methods that satisfy strong authentication, such as Passwordless sign-in with the Microsoft Authenticator",
            "policyType": "builtIn",
            "requirementsSatisfied": "mfa",
            "allowedCombinations": [
                "windowsHelloForBusiness",
                "fido2",
                "x509CertificateMultiFactor",
                "deviceBasedPush"
            ],
            "combinationConfigurations@odata.context": "https://graph.microsoft.com/beta/$metadata#identity/conditionalAccess/authenticationStrength/policies('00000000-0000-0000-0000-000000000003')/combinationConfigurations",
            "combinationConfigurations": []
        },
        {
            "id": "00000000-0000-0000-0000-000000000004",
            "createdDateTime": "2021-12-01T00:00:00Z",
            "modifiedDateTime": "2021-12-01T00:00:00Z",
            "displayName": "Phishing-resistant MFA",
            "description": "Phishing-resistant, Passwordless methods for the strongest authentication, such as a FIDO2 security key",
            "policyType": "builtIn",
            "requirementsSatisfied": "mfa",
            "allowedCombinations": [
                "windowsHelloForBusiness",
                "fido2",
                "x509CertificateMultiFactor"
            ],
            "combinationConfigurations@odata.context": "https://graph.microsoft.com/beta/$metadata#identity/conditionalAccess/authenticationStrength/policies('00000000-0000-0000-0000-000000000004')/combinationConfigurations",
            "combinationConfigurations": []
        }
    ]
}