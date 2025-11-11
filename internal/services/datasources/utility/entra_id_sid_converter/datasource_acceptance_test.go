package entra_id_sid_converter_test

import (
	"fmt"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccEntraIdSidConverterDataSource_SidToObjectId(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccConfigSidToObjectId(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.microsoft365_utility_entra_id_sid_converter.test", "sid", "S-1-12-1-1943430372-1249052806-2496021943-3034400218"),
					resource.TestCheckResourceAttr("data.microsoft365_utility_entra_id_sid_converter.test", "object_id", "73d664e4-0886-4a73-b745-c694da45ddb4"),
					resource.TestCheckResourceAttrSet("data.microsoft365_utility_entra_id_sid_converter.test", "id"),
				),
			},
		},
	})
}

func TestAccEntraIdSidConverterDataSource_ObjectIdToSid(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccConfigObjectIdToSid(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.microsoft365_utility_entra_id_sid_converter.test", "object_id", "73d664e4-0886-4a73-b745-c694da45ddb4"),
					resource.TestCheckResourceAttr("data.microsoft365_utility_entra_id_sid_converter.test", "sid", "S-1-12-1-1943430372-1249052806-2496021943-3034400218"),
					resource.TestCheckResourceAttrSet("data.microsoft365_utility_entra_id_sid_converter.test", "id"),
				),
			},
		},
	})
}

func TestAccEntraIdSidConverterDataSource_MaxUint32Values(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccConfigMaxUint32(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.microsoft365_utility_entra_id_sid_converter.test", "sid", "S-1-12-1-4294967295-4294967295-4294967295-4294967295"),
					resource.TestCheckResourceAttrSet("data.microsoft365_utility_entra_id_sid_converter.test", "object_id"),
					resource.TestCheckResourceAttrSet("data.microsoft365_utility_entra_id_sid_converter.test", "id"),
				),
			},
		},
	})
}

func TestAccEntraIdSidConverterDataSource_BidirectionalConversion(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccConfigBidirectional(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.microsoft365_utility_entra_id_sid_converter.sid_to_oid", "sid", "S-1-12-1-1943430372-1249052806-2496021943-3034400218"),
					resource.TestCheckResourceAttr("data.microsoft365_utility_entra_id_sid_converter.sid_to_oid", "object_id", "73d664e4-0886-4a73-b745-c694da45ddb4"),
					resource.TestCheckResourceAttr("data.microsoft365_utility_entra_id_sid_converter.oid_to_sid", "object_id", "73d664e4-0886-4a73-b745-c694da45ddb4"),
					resource.TestCheckResourceAttr("data.microsoft365_utility_entra_id_sid_converter.oid_to_sid", "sid", "S-1-12-1-1943430372-1249052806-2496021943-3034400218"),
				),
			},
		},
	})
}

// Acceptance test configuration functions
func testAccConfigSidToObjectId() string {
	accTestConfig, err := helpers.ParseHCLFile("tests/terraform/acceptance/01_sid_to_objectid.tf")
	if err != nil {
		panic(fmt.Sprintf("failed to load acceptance test config: %s", err.Error()))
	}
	return acceptance.ConfiguredM365ProviderBlock(accTestConfig)
}

func testAccConfigObjectIdToSid() string {
	accTestConfig, err := helpers.ParseHCLFile("tests/terraform/acceptance/02_objectid_to_sid.tf")
	if err != nil {
		panic(fmt.Sprintf("failed to load acceptance test config: %s", err.Error()))
	}
	return acceptance.ConfiguredM365ProviderBlock(accTestConfig)
}

func testAccConfigMaxUint32() string {
	accTestConfig, err := helpers.ParseHCLFile("tests/terraform/acceptance/03_max_uint32.tf")
	if err != nil {
		panic(fmt.Sprintf("failed to load acceptance test config: %s", err.Error()))
	}
	return acceptance.ConfiguredM365ProviderBlock(accTestConfig)
}

func testAccConfigBidirectional() string {
	accTestConfig, err := helpers.ParseHCLFile("tests/terraform/acceptance/04_bidirectional.tf")
	if err != nil {
		panic(fmt.Sprintf("failed to load acceptance test config: %s", err.Error()))
	}
	return acceptance.ConfiguredM365ProviderBlock(accTestConfig)
}
