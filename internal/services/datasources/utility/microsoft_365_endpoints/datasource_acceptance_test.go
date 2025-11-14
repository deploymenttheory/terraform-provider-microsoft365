package utilityMicrosoft365Endpoints_test

import (
	"os"
	"path/filepath"
	"regexp"
	"testing"

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
func TestAccMicrosoft365EndpointsDataSource_Worldwide(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccConfigWorldwide(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.microsoft365_utility_microsoft_365_endpoints.test", "instance", "worldwide"),
					resource.TestCheckResourceAttrSet("data.microsoft365_utility_microsoft_365_endpoints.test", "id"),
					resource.TestCheckResourceAttrSet("data.microsoft365_utility_microsoft_365_endpoints.test", "endpoints.#"),
					// Verify we get a reasonable number of endpoints (should be >= 50)
					resource.TestMatchResourceAttr("data.microsoft365_utility_microsoft_365_endpoints.test", "endpoints.#", regexp.MustCompile(`^([5-9][0-9]|[1-9][0-9]{2,})$`)),
				),
			},
		},
	})
}

func TestAccMicrosoft365EndpointsDataSource_FilterExchangeOptimize(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccConfigFilterExchangeOptimize(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.microsoft365_utility_microsoft_365_endpoints.test", "instance", "worldwide"),
					resource.TestCheckResourceAttr("data.microsoft365_utility_microsoft_365_endpoints.test", "service_areas.#", "1"),
					resource.TestCheckTypeSetElemAttr("data.microsoft365_utility_microsoft_365_endpoints.test", "service_areas.*", "Exchange"),
					resource.TestCheckResourceAttr("data.microsoft365_utility_microsoft_365_endpoints.test", "categories.#", "1"),
					resource.TestCheckTypeSetElemAttr("data.microsoft365_utility_microsoft_365_endpoints.test", "categories.*", "Optimize"),
					resource.TestCheckResourceAttrSet("data.microsoft365_utility_microsoft_365_endpoints.test", "endpoints.#"),
					// Should get at least 1 Exchange Optimize endpoint
					resource.TestMatchResourceAttr("data.microsoft365_utility_microsoft_365_endpoints.test", "endpoints.#", regexp.MustCompile(`^[1-9][0-9]*$`)),
				),
			},
		},
	})
}

func TestAccMicrosoft365EndpointsDataSource_TeamsMedia(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccConfigTeamsMedia(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.microsoft365_utility_microsoft_365_endpoints.test", "instance", "worldwide"),
					resource.TestCheckResourceAttr("data.microsoft365_utility_microsoft_365_endpoints.test", "service_areas.#", "1"),
					resource.TestCheckTypeSetElemAttr("data.microsoft365_utility_microsoft_365_endpoints.test", "service_areas.*", "Skype"),
					resource.TestCheckResourceAttr("data.microsoft365_utility_microsoft_365_endpoints.test", "categories.#", "1"),
					resource.TestCheckTypeSetElemAttr("data.microsoft365_utility_microsoft_365_endpoints.test", "categories.*", "Optimize"),
					resource.TestCheckResourceAttrSet("data.microsoft365_utility_microsoft_365_endpoints.test", "endpoints.#"),
					// Verify UDP ports are present for Teams media
					resource.TestMatchResourceAttr("data.microsoft365_utility_microsoft_365_endpoints.test", "endpoints.0.udp_ports", regexp.MustCompile(`3478|3479|3480|3481`)),
				),
			},
		},
	})
}

func TestAccMicrosoft365EndpointsDataSource_RequiredOnly(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccConfigRequiredOnly(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.microsoft365_utility_microsoft_365_endpoints.test", "instance", "worldwide"),
					resource.TestCheckResourceAttr("data.microsoft365_utility_microsoft_365_endpoints.test", "required_only", "true"),
					resource.TestCheckResourceAttrSet("data.microsoft365_utility_microsoft_365_endpoints.test", "endpoints.#"),
					// Required endpoints should be fewer than all endpoints
					resource.TestMatchResourceAttr("data.microsoft365_utility_microsoft_365_endpoints.test", "endpoints.#", regexp.MustCompile(`^[1-9][0-9]+$`)),
				),
			},
		},
	})
}

func TestAccMicrosoft365EndpointsDataSource_ExpressRoute(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccConfigExpressRoute(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.microsoft365_utility_microsoft_365_endpoints.test", "instance", "worldwide"),
					resource.TestCheckResourceAttr("data.microsoft365_utility_microsoft_365_endpoints.test", "express_route", "true"),
					resource.TestCheckResourceAttrSet("data.microsoft365_utility_microsoft_365_endpoints.test", "endpoints.#"),
					// ExpressRoute endpoints should be subset of all endpoints
					resource.TestMatchResourceAttr("data.microsoft365_utility_microsoft_365_endpoints.test", "endpoints.#", regexp.MustCompile(`^[1-9][0-9]*$`)),
				),
			},
		},
	})
}
