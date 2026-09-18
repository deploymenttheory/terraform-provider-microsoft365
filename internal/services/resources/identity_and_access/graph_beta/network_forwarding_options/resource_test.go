package graphBetaNetworkForwardingOptions_test

import (
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	target "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/identity_and_access/graph_beta/network_forwarding_options"
	optionMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/identity_and_access/graph_beta/network_forwarding_options/mocks"
)

const resourceType = target.ResourceName

func setup(t *testing.T) *optionMocks.ForwardingOptionsMock {
	t.Helper()
	mocks.SetupUnitTestEnvironment(t)
	httpmock.Activate()
	mocks.NewMocks().AuthMocks.RegisterMocks()
	m := &optionMocks.ForwardingOptionsMock{}
	m.RegisterMocks()
	t.Cleanup(func() { httpmock.DeactivateAndReset(); m.CleanupMockState() })
	return m
}

func config(t *testing.T, name string) string {
	t.Helper()
	s, e := helpers.ParseHCLFile("tests/terraform/unit/" + name)
	require.NoError(t, e)
	return s
}

func TestUnitResourceNetworkForwardingOptions_01_Lifecycle(t *testing.T) {
	m := setup(t)
	c := config(t, "resource.tf")
	u := config(t, "resource_updated.tf")
	resource.UnitTest(
		t,
		resource.TestCase{
			ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: c,
					Check: check.That(resourceType + ".test").
						Key("skip_dns_lookup_state").
						HasValue("enabled"),
				},
				{
					ResourceName:      resourceType + ".test",
					ImportState:       true,
					ImportStateId:     "forwardingOptions",
					ImportStateVerify: true,
				},
				{Config: c, PlanOnly: true},
				{
					Config: u,
					Check: check.That(resourceType + ".test").
						Key("skip_dns_lookup_state").
						HasValue("disabled"),
				},
				{Config: u, PlanOnly: true},
			},
		},
	)
	m.Lock()
	defer m.Unlock()
	require.Equal(t, "disabled", m.State["skipDnsLookupState"])
	require.Equal(t, "preserve-me", m.State["unmanagedFutureSetting"])
	require.Len(t, m.Patches, 2)
}

func TestUnitResourceNetworkForwardingOptions_02_Validation(t *testing.T) {
	for _, v := range []string{`""`, `"unknownFutureValue"`, `null`} {
		t.Run(v, func(t *testing.T) {
			setup(t)
			c := strings.Replace(config(t, "resource.tf"), `"enabled"`, v, 1)
			resource.UnitTest(
				t,
				resource.TestCase{
					ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
					Steps: []resource.TestStep{
						{
							Config: c,
							ExpectError: regexp.MustCompile(
								`value must be one of|Missing Configuration for Required Attribute`,
							),
						},
					},
				},
			)
		})
	}
	t.Run("omitted", func(t *testing.T) {
		setup(t)
		resource.UnitTest(
			t,
			resource.TestCase{
				ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{
						Config:      `resource "` + resourceType + `" "test" {}`,
						ExpectError: regexp.MustCompile(`Missing required argument`),
					},
				},
			},
		)
	})
}
