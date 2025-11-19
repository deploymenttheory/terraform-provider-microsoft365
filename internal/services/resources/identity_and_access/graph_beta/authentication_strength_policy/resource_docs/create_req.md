// Code snippets are only available for the latest major version. Current major version is $v0.*

// Dependencies
import (
	  "context"
	  msgraphsdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
	  graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
	  //other-imports
)

requestBody := graphmodels.NewAuthenticationStrengthPolicy()
displayName := "test-1"
requestBody.SetDisplayName(&displayName) 
description := "test-1"
requestBody.SetDescription(&description) 
allowedCombinations := []graphmodels.AuthenticationMethodModesable {
	authenticationMethodModes := graphmodels.DEVICEBASEDPUSH_AUTHENTICATIONMETHODMODES 
	requestBody.SetAuthenticationMethodModes(&authenticationMethodModes) 
	authenticationMethodModes := graphmodels.FEDERATEDMULTIFACTOR_AUTHENTICATIONMETHODMODES 
	requestBody.SetAuthenticationMethodModes(&authenticationMethodModes) 
	authenticationMethodModes := graphmodels.FEDERATEDSINGLEFACTOR_AUTHENTICATIONMETHODMODES 
	requestBody.SetAuthenticationMethodModes(&authenticationMethodModes) 
	authenticationMethodModes := graphmodels.FIDO2_AUTHENTICATIONMETHODMODES 
	requestBody.SetAuthenticationMethodModes(&authenticationMethodModes) 
	authenticationMethodModes := graphmodels.HARDWAREOATH,FEDERATEDSINGLEFACTOR_AUTHENTICATIONMETHODMODES 
	requestBody.SetAuthenticationMethodModes(&authenticationMethodModes) 
	authenticationMethodModes := graphmodels.MICROSOFTAUTHENTICATORPUSH,FEDERATEDSINGLEFACTOR_AUTHENTICATIONMETHODMODES 
	requestBody.SetAuthenticationMethodModes(&authenticationMethodModes) 
	authenticationMethodModes := graphmodels.PASSWORD_AUTHENTICATIONMETHODMODES 
	requestBody.SetAuthenticationMethodModes(&authenticationMethodModes) 
	authenticationMethodModes := graphmodels.PASSWORD,HARDWAREOATH_AUTHENTICATIONMETHODMODES 
	requestBody.SetAuthenticationMethodModes(&authenticationMethodModes) 
	authenticationMethodModes := graphmodels.PASSWORD,MICROSOFTAUTHENTICATORPUSH_AUTHENTICATIONMETHODMODES 
	requestBody.SetAuthenticationMethodModes(&authenticationMethodModes) 
	authenticationMethodModes := graphmodels.PASSWORD,SMS_AUTHENTICATIONMETHODMODES 
	requestBody.SetAuthenticationMethodModes(&authenticationMethodModes) 
	authenticationMethodModes := graphmodels.PASSWORD,SOFTWAREOATH_AUTHENTICATIONMETHODMODES 
	requestBody.SetAuthenticationMethodModes(&authenticationMethodModes) 
	authenticationMethodModes := graphmodels.PASSWORD,VOICE_AUTHENTICATIONMETHODMODES 
	requestBody.SetAuthenticationMethodModes(&authenticationMethodModes) 
	authenticationMethodModes := graphmodels.QRCODEPIN_AUTHENTICATIONMETHODMODES 
	requestBody.SetAuthenticationMethodModes(&authenticationMethodModes) 
	authenticationMethodModes := graphmodels.SMS_AUTHENTICATIONMETHODMODES 
	requestBody.SetAuthenticationMethodModes(&authenticationMethodModes) 
	authenticationMethodModes := graphmodels.SMS,FEDERATEDSINGLEFACTOR_AUTHENTICATIONMETHODMODES 
	requestBody.SetAuthenticationMethodModes(&authenticationMethodModes) 
	authenticationMethodModes := graphmodels.SOFTWAREOATH,FEDERATEDSINGLEFACTOR_AUTHENTICATIONMETHODMODES 
	requestBody.SetAuthenticationMethodModes(&authenticationMethodModes) 
	authenticationMethodModes := graphmodels.TEMPORARYACCESSPASSMULTIUSE_AUTHENTICATIONMETHODMODES 
	requestBody.SetAuthenticationMethodModes(&authenticationMethodModes) 
	authenticationMethodModes := graphmodels.TEMPORARYACCESSPASSONETIME_AUTHENTICATIONMETHODMODES 
	requestBody.SetAuthenticationMethodModes(&authenticationMethodModes) 
	authenticationMethodModes := graphmodels.VOICE,FEDERATEDSINGLEFACTOR_AUTHENTICATIONMETHODMODES 
	requestBody.SetAuthenticationMethodModes(&authenticationMethodModes) 
	authenticationMethodModes := graphmodels.WINDOWSHELLOFORBUSINESS_AUTHENTICATIONMETHODMODES 
	requestBody.SetAuthenticationMethodModes(&authenticationMethodModes) 
	authenticationMethodModes := graphmodels.X509CERTIFICATEMULTIFACTOR_AUTHENTICATIONMETHODMODES 
	requestBody.SetAuthenticationMethodModes(&authenticationMethodModes) 
	authenticationMethodModes := graphmodels.X509CERTIFICATESINGLEFACTOR_AUTHENTICATIONMETHODMODES 
	requestBody.SetAuthenticationMethodModes(&authenticationMethodModes)
}
requestBody.SetAllowedCombinations(allowedCombinations)


authenticationCombinationConfiguration := graphmodels.NewFido2CombinationConfiguration()
appliesToCombinations := []graphmodels.AuthenticationMethodModesable {
	authenticationMethodModes := graphmodels.FIDO2_AUTHENTICATIONMETHODMODES 
	authenticationCombinationConfiguration.SetAuthenticationMethodModes(&authenticationMethodModes)
}
authenticationCombinationConfiguration.SetAppliesToCombinations(appliesToCombinations)
allowedAAGUIDs := []string {
	"12345678-0000-0000-0000-123456780000",
	"12345678-0000-0000-0000-123456780001",
	"90a3ccdf-635c-4729-a248-9b709135078f",
	"de1e552d-db1d-4423-a619-566b625cdc84",
}
authenticationCombinationConfiguration.SetAllowedAAGUIDs(allowedAAGUIDs)
authenticationCombinationConfiguration1 := graphmodels.NewX509CertificateCombinationConfiguration()
appliesToCombinations := []graphmodels.AuthenticationMethodModesable {
	authenticationMethodModes := graphmodels.X509CERTIFICATEMULTIFACTOR_AUTHENTICATIONMETHODMODES 
	authenticationCombinationConfiguration1.SetAuthenticationMethodModes(&authenticationMethodModes)
}
authenticationCombinationConfiguration1.SetAppliesToCombinations(appliesToCombinations)
allowedIssuerSkis := []string {
	"9A4248C6AC8C2931AB2A86537818E92E7B6C97B6",
	"9A4248C6AC8C2931AB2A86537818E92E7B6C97B7",
}
authenticationCombinationConfiguration1.SetAllowedIssuerSkis(allowedIssuerSkis)
allowedPolicyOIDs := []string {
	"1.3.6.1.4.1.311.21.8.1.4",
	"1.3.6.1.4.1.311.21.8.1.7",
}
authenticationCombinationConfiguration1.SetAllowedPolicyOIDs(allowedPolicyOIDs)
additionalData := map[string]interface{}{
	allowedIssuers := []string {
		"CUSTOMIDENTIFIER:9A4248C6AC8C2931AB2A86537818E92E7B6C97B6",
		"CUSTOMIDENTIFIER:9A4248C6AC8C2931AB2A86537818E92E7B6C97B7",
	}
}
authenticationCombinationConfiguration1.SetAdditionalData(additionalData)
authenticationCombinationConfiguration2 := graphmodels.NewX509CertificateCombinationConfiguration()
appliesToCombinations := []graphmodels.AuthenticationMethodModesable {
	authenticationMethodModes := graphmodels.X509CERTIFICATESINGLEFACTOR_AUTHENTICATIONMETHODMODES 
	authenticationCombinationConfiguration2.SetAuthenticationMethodModes(&authenticationMethodModes)
}
authenticationCombinationConfiguration2.SetAppliesToCombinations(appliesToCombinations)
allowedIssuerSkis := []string {
	"1A2B3C4D5E6F7A8B9C0D1E2F3A4B5C6D7E8F9A0C",
	"1A2B3C4D5E6F7A8B9C0D1E2F3A4B5C6D7E8F9A0E",
}
authenticationCombinationConfiguration2.SetAllowedIssuerSkis(allowedIssuerSkis)
allowedPolicyOIDs := []string {
	"1.3.6.1.4.1.311.21.8.1.4",
	"1.3.6.1.4.1.311.21.8.1.7",
}
authenticationCombinationConfiguration2.SetAllowedPolicyOIDs(allowedPolicyOIDs)
additionalData := map[string]interface{}{
	allowedIssuers := []string {
		"CUSTOMIDENTIFIER:1A2B3C4D5E6F7A8B9C0D1E2F3A4B5C6D7E8F9A0C",
		"CUSTOMIDENTIFIER:1A2B3C4D5E6F7A8B9C0D1E2F3A4B5C6D7E8F9A0E",
	}
}
authenticationCombinationConfiguration2.SetAdditionalData(additionalData)

combinationConfigurations := []graphmodels.AuthenticationCombinationConfigurationable {
	authenticationCombinationConfiguration,
	authenticationCombinationConfiguration1,
	authenticationCombinationConfiguration2,
}
requestBody.SetCombinationConfigurations(combinationConfigurations)

// To initialize your graphClient, see https://learn.microsoft.com/en-us/graph/sdks/create-client?from=snippets&tabs=go
policies, err := graphClient.Identity().ConditionalAccess().AuthenticationStrength().Policies().Post(context.Background(), requestBody, nil)