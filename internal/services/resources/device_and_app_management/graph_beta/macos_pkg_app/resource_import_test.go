package graphBetaMacOSPKGApp_test

import (
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/jarcoal/httpmock"
)

func TestUnitResourceMacOSPkgApp_07_GeneratedImportURLs(t *testing.T) {
	for _, test := range []struct {
		name    string
		privacy any
		want    knownvalue.Check
	}{
		{name: "empty privacy", privacy: "", want: knownvalue.Null()},
		{name: "null privacy", want: knownvalue.Null()},
		{name: "valid privacy", privacy: "https://example.com/privacy", want: knownvalue.StringExact("https://example.com/privacy")},
	} {
		t.Run(test.name, func(t *testing.T) {
			mocks.SetupUnitTestEnvironment(t)
			httpmock.Activate()
			defer httpmock.DeactivateAndReset()
			minimumOS := map[string]bool{}
			for _, name := range []string{"v10_7", "v10_8", "v10_9", "v10_10", "v10_11", "v10_12", "v10_13", "v10_14", "v10_15", "v11_0", "v12_0", "v13_0", "v14_0", "v15_0"} {
				minimumOS[name] = name == "v12_0"
			}
			body := map[string]any{
				"@odata.type": "#microsoft.graph.macOSPkgApp",
				"id":          "11111111-1111-1111-1111-111111111111",
				"displayName": "unit-test-macos-pkg-url-import", "description": "URL import regression", "publisher": "MacAdmins",
				"isFeatured": false, "privacyInformationUrl": test.privacy, "informationUrl": nil,
				"primaryBundleId": "com.github.macadmins.Nudge", "primaryBundleVersion": "2.1.4.81863",
				"ignoreVersionDetection": false, "minimumSupportedOperatingSystem": minimumOS,
				"roleScopeTagIds": []string{"0"},
			}
			responder, err := httpmock.NewJsonResponder(http.StatusOK, body)
			if err != nil {
				t.Fatal(err)
			}
			httpmock.RegisterRegexpResponder("GET", regexp.MustCompile("^https://graph.microsoft.com/beta/deviceAppManagement/mobileApps/11111111-1111-1111-1111-111111111111"), responder)
			httpmock.RegisterResponder("DELETE", "https://graph.microsoft.com/beta/deviceAppManagement/mobileApps/11111111-1111-1111-1111-111111111111", httpmock.NewStringResponder(http.StatusNoContent, ""))
			config, err := helpers.ParseHCLFile("tests/terraform/unit/import_urls.tf")
			if err != nil {
				t.Fatal(err)
			}
			resource.UnitTest(t, resource.TestCase{
				TerraformVersionChecks:   []tfversion.TerraformVersionCheck{tfversion.SkipBelow(tfversion.Version1_5_0)},
				ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
				Steps: []resource.TestStep{{
					ResourceName:  resourceType + ".test",
					ImportStateId: "11111111-1111-1111-1111-111111111111",
					ImportState:   true, GenerateConfig: true,
					ImportStateKind:        resource.ImportBlockWithID,
					ImportStateConfigExact: true,
					Config:                 config,
					ImportPlanChecks: resource.ImportPlanChecks{
						PreApply: []plancheck.PlanCheck{
							plancheck.ExpectResourceAction(resourceType+".test", plancheck.ResourceActionNoop),
							plancheck.ExpectKnownValue(resourceType+".test", tfjsonpath.New("privacy_information_url"), test.want),
							plancheck.ExpectKnownValue(resourceType+".test", tfjsonpath.New("information_url"), knownvalue.Null()),
						},
					},
				}},
			})
		})
	}
}

func TestUnitResourceMacOSPkgApp_08_InvalidPrivacyURLs(t *testing.T) {
	for _, value := range []string{"", "relative/path", "not-a-url"} {
		t.Run(value, func(t *testing.T) {
			mocks.SetupUnitTestEnvironment(t)
			config, err := helpers.ParseHCLFile("tests/terraform/unit/resource_minimal.tf")
			if err != nil {
				t.Fatal(err)
			}
			config = strings.Replace(config, "# FEATURED", "privacy_information_url = "+strconv.Quote(value), 1)
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
				Steps: []resource.TestStep{{
					Config: config, PlanOnly: true,
					ExpectError: regexp.MustCompile("must be a valid URL starting with http://"),
				}},
			})
		})
	}
}
