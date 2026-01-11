package graphBetaApplicationsIpApplicationSegment_test

import (
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	ipSegmentMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/applications/graph_beta/ip_application_segment/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

func setupMockEnvironment() (*mocks.Mocks, *ipSegmentMocks.IpApplicationSegmentMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	ipSegmentMock := &ipSegmentMocks.IpApplicationSegmentMock{}
	ipSegmentMock.RegisterMocks()
	return mockClient, ipSegmentMock
}

func testCheckExists(resourceName string) resource.TestCheckFunc {
	return resource.TestCheckResourceAttrSet(resourceName, "id")
}

func TestIpApplicationSegmentResource_Minimal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, ipSegmentMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer ipSegmentMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigIpSegmentMinimal(),
				Check: resource.ComposeTestCheckFunc(
					// Basic attributes
					resource.TestCheckResourceAttr("microsoft365_graph_beta_applications_ip_application_segment.ip_segment_minimal", "application_id", "12345678-1234-1234-1234-123456789012"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_applications_ip_application_segment.ip_segment_minimal", "destination_host", "192.168.1.100"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_applications_ip_application_segment.ip_segment_minimal", "destination_type", "ipAddress"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_applications_ip_application_segment.ip_segment_minimal", "protocol", "tcp"),
					resource.TestMatchResourceAttr("microsoft365_graph_beta_applications_ip_application_segment.ip_segment_minimal", "id", regexp.MustCompile(`^[0-9a-fA-F-]+$`)),

					// Ports
					resource.TestCheckResourceAttr("microsoft365_graph_beta_applications_ip_application_segment.ip_segment_minimal", "ports.#", "1"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_applications_ip_application_segment.ip_segment_minimal", "ports.*", "80-80"),
				),
			},
		},
	})
}

func TestIpApplicationSegmentResource_Maximal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, ipSegmentMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer ipSegmentMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigIpSegmentMaximal(),
				Check: resource.ComposeTestCheckFunc(
					// Basic attributes
					testCheckExists("microsoft365_graph_beta_applications_ip_application_segment.ip_segment_maximal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_applications_ip_application_segment.ip_segment_maximal", "application_id", "12345678-1234-1234-1234-123456789012"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_applications_ip_application_segment.ip_segment_maximal", "destination_host", "*.example.com"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_applications_ip_application_segment.ip_segment_maximal", "destination_type", "dnsSuffix"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_applications_ip_application_segment.ip_segment_maximal", "protocol", "tcp"),

					// Ports
					resource.TestCheckResourceAttr("microsoft365_graph_beta_applications_ip_application_segment.ip_segment_maximal", "ports.#", "4"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_applications_ip_application_segment.ip_segment_maximal", "ports.*", "80-80"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_applications_ip_application_segment.ip_segment_maximal", "ports.*", "443-443"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_applications_ip_application_segment.ip_segment_maximal", "ports.*", "8080-8080"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_applications_ip_application_segment.ip_segment_maximal", "ports.*", "8443-8443"),
				),
			},
		},
	})
}

func TestIpApplicationSegmentResource_IpRange(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, ipSegmentMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer ipSegmentMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigIpSegmentRange(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_applications_ip_application_segment.ip_segment_range"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_applications_ip_application_segment.ip_segment_range", "destination_host", "192.168.1.0/24"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_applications_ip_application_segment.ip_segment_range", "destination_type", "ipRangeCidr"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_applications_ip_application_segment.ip_segment_range", "protocol", "tcp"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_applications_ip_application_segment.ip_segment_range", "ports.#", "1"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_applications_ip_application_segment.ip_segment_range", "ports.*", "443-443"),
				),
			},
		},
	})
}

func TestIpApplicationSegmentResource_FQDN(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, ipSegmentMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer ipSegmentMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigIpSegmentFQDN(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_applications_ip_application_segment.ip_segment_fqdn"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_applications_ip_application_segment.ip_segment_fqdn", "destination_host", "app.example.com"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_applications_ip_application_segment.ip_segment_fqdn", "destination_type", "fqdn"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_applications_ip_application_segment.ip_segment_fqdn", "protocol", "tcp"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_applications_ip_application_segment.ip_segment_fqdn", "ports.#", "2"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_applications_ip_application_segment.ip_segment_fqdn", "ports.*", "443-443"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_applications_ip_application_segment.ip_segment_fqdn", "ports.*", "8443-8443"),
				),
			},
		},
	})
}

func TestIpApplicationSegmentResource_Update(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, ipSegmentMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer ipSegmentMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// Create minimal config
				Config: testConfigIpSegmentMinimal(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("microsoft365_graph_beta_applications_ip_application_segment.ip_segment_minimal", "ports.#", "1"),
					testCheckExists("microsoft365_graph_beta_applications_ip_application_segment.ip_segment_minimal"),
				),
			},
			{
				// Update to maximal config
				Config: testConfigIpSegmentMaximal(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("microsoft365_graph_beta_applications_ip_application_segment.ip_segment_maximal", "ports.#", "4"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_applications_ip_application_segment.ip_segment_maximal", "destination_host", "*.example.com"),
				),
			},
		},
	})
}

// Configuration helper functions
func testConfigIpSegmentMinimal() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/resource_minimal.tf")
	if err != nil {
		panic("failed to load ip application segment minimal config: " + err.Error())
	}
	return unitTestConfig
}

func testConfigIpSegmentMaximal() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/resource_maximal.tf")
	if err != nil {
		panic("failed to load ip application segment maximal config: " + err.Error())
	}
	return unitTestConfig
}

func testConfigIpSegmentRange() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/resource_ip_range.tf")
	if err != nil {
		panic("failed to load ip application segment range config: " + err.Error())
	}
	return unitTestConfig
}

func testConfigIpSegmentFQDN() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/resource_fqdn.tf")
	if err != nil {
		panic("failed to load ip application segment fqdn config: " + err.Error())
	}
	return unitTestConfig
}
