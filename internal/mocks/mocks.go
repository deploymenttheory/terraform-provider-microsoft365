package mocks

// import (
// 	"context"
// 	"net/http"
// 	"path/filepath"
// 	"regexp"
// 	"runtime"
// 	"strings"

// 	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
// 	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/provider"
// 	"github.com/hashicorp/terraform-plugin-framework/providerserver"
// 	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
// 	"github.com/jarcoal/httpmock"
// )

// func TestName() string {
// 	pc, _, _, _ := runtime.Caller(1)
// 	nameFull := runtime.FuncForPC(pc).Name()
// 	nameEnd := filepath.Ext(nameFull)
// 	name := strings.TrimPrefix(nameEnd, ".")
// 	return name
// }

// func TestsEntraLicesingGroupName() string {
// 	return "pptestusers"
// }

// var TestUnitTestProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
// 	"microsoft365": providerserver.NewProtocol6WithError(provider.NewMicrosoft365Provider(helpers.UnitTestContext(context.Background(), ""), true)()),
// }

// var TestAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
// 	"microsoft365": providerserver.NewProtocol6WithError(provider.NewMicrosoft365Provider(context.Background(), false)()),
// }

// func ActivateEnvironmentHttpMocks() {
// 	httpmock.RegisterResponder("GET", `=~^https://([\d-]+)\.crm4\.dynamics\.com/api/data/v9\.2/transactioncurrencies\z`,
// 		func(req *http.Request) (*http.Response, error) {
// 			return httpmock.NewStringResponse(http.StatusOK, `{
// 				"value": [
// 					{
// 						"isocurrencycode": "PLN"
// 					}]}`), nil
// 		})

// 	httpmock.RegisterResponder("GET", `=~^https://([\d-]+)\.crm4\.dynamics\.com/api/data/v9\.2/organizations\z`,
// 		func(req *http.Request) (*http.Response, error) {
// 			return httpmock.NewStringResponse(http.StatusOK, `{
// 				"value": [
// 					{
// 						"_basecurrencyid_value": "xyz"
// 					}]}`), nil
// 		})

// 	httpmock.RegisterRegexpResponder("GET", regexp.MustCompile(`^https://api\.bap\.microsoft\.com/providers/Microsoft\.BusinessAppPlatform/locations/(europe|unitedstates)/environmentLanguages\?api-version=2023-06-01$`),
// 		func(req *http.Request) (*http.Response, error) {
// 			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("../../services/languages/tests/datasource/Validate_Read/get_languages.json").String()), nil
// 		})

// 	httpmock.RegisterRegexpResponder("GET", regexp.MustCompile(`^https://api\.bap\.microsoft\.com/providers/Microsoft\.BusinessAppPlatform/locations/(europe|unitedstates)/environmentCurrencies\?api-version=2023-06-01$`),
// 		func(req *http.Request) (*http.Response, error) {
// 			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("../../services/currencies/tests/datasource/Validate_Read/get_currencies.json").String()), nil
// 		})

// 	httpmock.RegisterResponder("GET", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/locations?api-version=2023-06-01",
// 		func(req *http.Request) (*http.Response, error) {
// 			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("../../services/locations/tests/datasource/Validate_Read/get_locations.json").String()), nil
// 		})

// 	httpmock.RegisterResponder("POST", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/validateEnvironmentDetails?api-version=2021-04-01",
// 		func(req *http.Request) (*http.Response, error) {
// 			return httpmock.NewStringResponse(http.StatusOK, ""), nil
// 		})

// 	httpmock.RegisterResponder("GET", `=~^https://([a-z\d-]+)\.(crm[\d]*)\.dynamics\.com/api/data/v9\.2/WhoAmI\z`,
// 		func(req *http.Request) (*http.Response, error) {
// 			return httpmock.NewStringResponse(http.StatusOK, `{
// 				"@odata.context": "https://org000001.crm.dynamics.com/api/data/v9.2/$metadata#Microsoft.Dynamics.CRM.WhoAmIResponse",
// 				"BusinessUnitId": "00000000-0000-0000-0000-000000000002",
// 				"UserId": "00000000-0000-0000-0000-000000000001",
// 				"OrganizationId": "00000000-0000-0000-0000-000000000003"
// 			}`), nil
// 		})

// 	httpmock.RegisterResponder("GET", `https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/tenant?api-version=2021-04-01`,
// 		func(_ *http.Request) (*http.Response, error) {
// 			return httpmock.NewStringResponse(http.StatusOK, `{
// 				"tenantId": "00000000-0000-0000-0000-000000000001",
// 				"state": "Enabled",
// 				"location": "unitedstates",
// 				"aadCountryGeo": "unitedstates",
// 				"dataStorageGeo": "unitedstates",
// 				"defaultEnvironmentGeo": "unitedstates",
// 				"aadDataBoundary": "none",
// 				"fedRAMPHighCertificationRequired": false
// 			}`), nil
// 		})
// }
