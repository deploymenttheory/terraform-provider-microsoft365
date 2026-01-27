package utilityMicrosoft365EndpointReference_test

import (
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Helper functions to load each test configuration from acceptance directory
func testAccConfigWorldwide() string {
	content, err := os.ReadFile(filepath.Join("tests", "terraform", "acceptance", "datasource_worldwide.tf"))
	if err != nil {
		return ""
	}
	return string(content)
}

func testAccConfigFilterExchangeOptimize() string {
	content, err := os.ReadFile(filepath.Join("tests", "terraform", "acceptance", "datasource_filter_exchange_optimize.tf"))
	if err != nil {
		return ""
	}
	return string(content)
}

func testAccConfigTeamsMedia() string {
	content, err := os.ReadFile(filepath.Join("tests", "terraform", "acceptance", "datasource_teams_media.tf"))
	if err != nil {
		return ""
	}
	return string(content)
}

func testAccConfigRequiredOnly() string {
	content, err := os.ReadFile(filepath.Join("tests", "terraform", "acceptance", "datasource_required_only.tf"))
	if err != nil {
		return ""
	}
	return string(content)
}

func testAccConfigExpressRoute() string {
	content, err := os.ReadFile(filepath.Join("tests", "terraform", "acceptance", "datasource_expressroute.tf"))
	if err != nil {
		return ""
	}
	return string(content)
}

// Acceptance test cases - these call the real Microsoft endpoints API
func TestAccDatasourceMicrosoft365EndpointReference_01_Worldwide(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccConfigWorldwide(),
				Check: resource.ComposeTestCheckFunc(
					check.That("data."+dataSourceType+".test").Key("instance").HasValue("worldwide"),
					check.That("data."+dataSourceType+".test").Key("id").IsSet(),
					check.That("data."+dataSourceType+".test").Key("endpoints.#").IsSet(),
					// Verify we get a reasonable number of endpoints (should be >= 50)
					check.That("data."+dataSourceType+".test").Key("endpoints.#").MatchesRegex(regexp.MustCompile(`^([5-9][0-9]|[1-9][0-9]{2,})$`)),
				),
			},
		},
	})
}

func TestAccDatasourceMicrosoft365EndpointReference_02_FilterExchangeOptimize(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccConfigFilterExchangeOptimize(),
				Check: resource.ComposeTestCheckFunc(
					check.That("data."+dataSourceType+".test").Key("instance").HasValue("worldwide"),
					check.That("data."+dataSourceType+".test").Key("service_areas.#").HasValue("1"),
					check.That("data."+dataSourceType+".test").Key("service_areas.*").ContainsTypeSetElement("Exchange"),
					check.That("data."+dataSourceType+".test").Key("categories.#").HasValue("1"),
					check.That("data."+dataSourceType+".test").Key("categories.*").ContainsTypeSetElement("Optimize"),
					check.That("data."+dataSourceType+".test").Key("endpoints.#").IsSet(),
					// Should get at least 1 Exchange Optimize endpoint
					check.That("data."+dataSourceType+".test").Key("endpoints.#").MatchesRegex(regexp.MustCompile(`^[1-9][0-9]*$`)),
				),
			},
		},
	})
}

func TestAccDatasourceMicrosoft365EndpointReference_03_TeamsMedia(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccConfigTeamsMedia(),
				Check: resource.ComposeTestCheckFunc(
					check.That("data."+dataSourceType+".test").Key("instance").HasValue("worldwide"),
					check.That("data."+dataSourceType+".test").Key("service_areas.#").HasValue("1"),
					check.That("data."+dataSourceType+".test").Key("service_areas.*").ContainsTypeSetElement("Skype"),
					check.That("data."+dataSourceType+".test").Key("categories.#").HasValue("1"),
					check.That("data."+dataSourceType+".test").Key("categories.*").ContainsTypeSetElement("Optimize"),
					check.That("data."+dataSourceType+".test").Key("endpoints.#").IsSet(),
					// Verify UDP ports are present for Teams media
					check.That("data."+dataSourceType+".test").Key("endpoints.0.udp_ports").MatchesRegex(regexp.MustCompile(`3478|3479|3480|3481`)),
				),
			},
		},
	})
}

func TestAccDatasourceMicrosoft365EndpointReference_04_RequiredOnly(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccConfigRequiredOnly(),
				Check: resource.ComposeTestCheckFunc(
					check.That("data."+dataSourceType+".test").Key("instance").HasValue("worldwide"),
					check.That("data."+dataSourceType+".test").Key("required_only").HasValue("true"),
					check.That("data."+dataSourceType+".test").Key("endpoints.#").IsSet(),
					// Required endpoints should be fewer than all endpoints
					check.That("data."+dataSourceType+".test").Key("endpoints.#").MatchesRegex(regexp.MustCompile(`^[1-9][0-9]+$`)),
				),
			},
		},
	})
}

func TestAccDatasourceMicrosoft365EndpointReference_05_ExpressRoute(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccConfigExpressRoute(),
				Check: resource.ComposeTestCheckFunc(
					check.That("data."+dataSourceType+".test").Key("instance").HasValue("worldwide"),
					check.That("data."+dataSourceType+".test").Key("express_route").HasValue("true"),
					check.That("data."+dataSourceType+".test").Key("endpoints.#").IsSet(),
					// ExpressRoute endpoints should be subset of all endpoints
					check.That("data."+dataSourceType+".test").Key("endpoints.#").MatchesRegex(regexp.MustCompile(`^[1-9][0-9]*$`)),
				),
			},
		},
	})
}
