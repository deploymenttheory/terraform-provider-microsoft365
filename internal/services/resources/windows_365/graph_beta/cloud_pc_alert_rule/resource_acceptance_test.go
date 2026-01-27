package graphBetaCloudPcAlertRule_test

import (
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceCloudPcAlertRule_01_Complete(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create with minimal configuration
			{
				Config: testAccCloudPcAlertRuleConfig_minimal(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_windows_365_cloud_pc_alert_rule.test", "id"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_alert_rule.test", "alert_rule_template", "cloudPcProvisionScenario"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_alert_rule.test", "display_name", "Test Acceptance Cloud PC Alert Rule"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_alert_rule.test", "severity", "warning"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_alert_rule.test", "enabled", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_alert_rule.test", "notification_channels.#", "1"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_alert_rule.test", "threshold.aggregation", "count"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_alert_rule.test", "conditions.#", "1"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "microsoft365_graph_beta_windows_365_cloud_pc_alert_rule.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update to maximal configuration
			{
				Config: testAccCloudPcAlertRuleConfig_maximal(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_windows_365_cloud_pc_alert_rule.test", "id"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_alert_rule.test", "alert_rule_template", "cloudPcProvisionScenario"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_alert_rule.test", "display_name", "Test Acceptance Cloud PC Alert Rule - Updated"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_alert_rule.test", "description", "Updated description for acceptance testing"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_alert_rule.test", "severity", "critical"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_alert_rule.test", "enabled", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_alert_rule.test", "notification_channels.#", "2"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_alert_rule.test", "threshold.target", "5"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_alert_rule.test", "conditions.#", "2"),
				),
			},
			// Update back to minimal configuration
			{
				Config: testAccCloudPcAlertRuleConfig_minimal(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_windows_365_cloud_pc_alert_rule.test", "id"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_alert_rule.test", "display_name", "Test Acceptance Cloud PC Alert Rule"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_alert_rule.test", "severity", "warning"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_alert_rule.test", "notification_channels.#", "1"),
				),
			},
		},
	})
}

func TestAccResourceCloudPcAlertRule_02_RequiredFields(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccCloudPcAlertRuleConfig_missingAlertRuleTemplate(),
				ExpectError: regexp.MustCompile("Missing required argument"),
			},
			{
				Config:      testAccCloudPcAlertRuleConfig_missingDisplayName(),
				ExpectError: regexp.MustCompile("Missing required argument"),
			},
			{
				Config:      testAccCloudPcAlertRuleConfig_missingSeverity(),
				ExpectError: regexp.MustCompile("Missing required argument"),
			},
			{
				Config:      testAccCloudPcAlertRuleConfig_missingNotificationChannels(),
				ExpectError: regexp.MustCompile("Missing required argument"),
			},
			{
				Config:      testAccCloudPcAlertRuleConfig_missingThreshold(),
				ExpectError: regexp.MustCompile("Missing required argument"),
			},
		},
	})
}

func TestAccResourceCloudPcAlertRule_03_InvalidValues(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccCloudPcAlertRuleConfig_invalidSeverity(),
				ExpectError: regexp.MustCompile("Attribute severity value must be one of"),
			},
			{
				Config:      testAccCloudPcAlertRuleConfig_invalidAlertRuleTemplate(),
				ExpectError: regexp.MustCompile("Attribute alert_rule_template value must be one of"),
			},
		},
	})
}

func testAccCloudPcAlertRuleConfig_minimal() string {
	return `
resource "microsoft365_graph_beta_windows_365_cloud_pc_alert_rule" "test" {
  alert_rule_template = "cloudPcProvisionScenario"
  display_name        = "Test Acceptance Cloud PC Alert Rule"
  severity            = "warning"
  enabled             = true
  is_system_rule      = false

  notification_channels = [
    {
      notification_channel_type = "portal"
      notification_receivers = [
        {
          contact_information = "admin@test.com"
          locale             = "en-US"
        }
      ]
    }
  ]

  threshold = {
    aggregation = "count"
    operator    = "greaterOrEqual"
    target      = 1
  }

  conditions = [
    {
      relationship_type   = "and"
      condition_category  = "provisionFailures"
      aggregation        = "count"
      operator           = "greaterOrEqual"
      threshold_value    = "1"
    }
  ]
}
`
}

func testAccCloudPcAlertRuleConfig_maximal() string {
	return `
resource "microsoft365_graph_beta_windows_365_cloud_pc_alert_rule" "test" {
  alert_rule_template = "cloudPcProvisionScenario"
  display_name        = "Test Acceptance Cloud PC Alert Rule - Updated"
  description         = "Updated description for acceptance testing"
  severity            = "critical"
  enabled             = true
  is_system_rule      = false

  notification_channels = [
    {
      notification_channel_type = "portal"
      notification_receivers = [
        {
          contact_information = "admin@test.com"
          locale             = "en-US"
        },
        {
          contact_information = "manager@test.com"
          locale             = "en-US"
        }
      ]
    },
    {
      notification_channel_type = "email"
      notification_receivers = [
        {
          contact_information = "alerts@test.com"
          locale             = "en-US"
        }
      ]
    }
  ]

  threshold = {
    aggregation = "count"
    operator    = "greaterOrEqual"
    target      = 5
  }

  conditions = [
    {
      relationship_type   = "and"
      condition_category  = "provisionFailures"
      aggregation        = "count"
      operator           = "greaterOrEqual"
      threshold_value    = "3"
    },
    {
      relationship_type   = "or"
      condition_category  = "cloudPcConnectionErrors"
      aggregation        = "percentage"
      operator           = "less"
      threshold_value    = "95"
    }
  ]
}
`
}

func testAccCloudPcAlertRuleConfig_missingAlertRuleTemplate() string {
	return `
resource "microsoft365_graph_beta_windows_365_cloud_pc_alert_rule" "test" {
  display_name   = "Test Alert Rule"
  severity       = "warning"
  enabled        = true
  is_system_rule = false

  notification_channels = [
    {
      notification_channel_type = "portal"
      notification_receivers = [
        {
          contact_information = "admin@test.com"
          locale             = "en-US"
        }
      ]
    }
  ]

  threshold = {
    aggregation = "count"
    operator    = "greaterOrEqual"
    target      = 1
  }

  conditions = [
    {
      relationship_type   = "and"
      condition_category  = "provisionFailures"
      aggregation        = "count"
      operator           = "greaterOrEqual"
      threshold_value    = "1"
    }
  ]
}
`
}

func testAccCloudPcAlertRuleConfig_missingDisplayName() string {
	return `
resource "microsoft365_graph_beta_windows_365_cloud_pc_alert_rule" "test" {
  alert_rule_template = "cloudPcProvisionScenario"
  severity            = "warning"
  enabled             = true
  is_system_rule      = false

  notification_channels = [
    {
      notification_channel_type = "portal"
      notification_receivers = [
        {
          contact_information = "admin@test.com"
          locale             = "en-US"
        }
      ]
    }
  ]

  threshold = {
    aggregation = "count"
    operator    = "greaterOrEqual"
    target      = 1
  }

  conditions = [
    {
      relationship_type   = "and"
      condition_category  = "provisionFailures"
      aggregation        = "count"
      operator           = "greaterOrEqual"
      threshold_value    = "1"
    }
  ]
}
`
}

func testAccCloudPcAlertRuleConfig_missingSeverity() string {
	return `
resource "microsoft365_graph_beta_windows_365_cloud_pc_alert_rule" "test" {
  alert_rule_template = "cloudPcProvisionScenario"
  display_name        = "Test Alert Rule"
  enabled             = true

  notification_channels = [
    {
      notification_channel_type = "portal"
      notification_receivers = [
        {
          contact_information = "admin@test.com"
          locale             = "en-US"
        }
      ]
    }
  ]

  threshold = {
    aggregation = "count"
    operator    = "greaterOrEqual"
    target      = 1
  }

  conditions = [
    {
      relationship_type   = "and"
      condition_category  = "provisionFailures"
      aggregation        = "count"
      operator           = "greaterOrEqual"
      threshold_value    = "1"
    }
  ]
}
`
}

func testAccCloudPcAlertRuleConfig_missingNotificationChannels() string {
	return `
resource "microsoft365_graph_beta_windows_365_cloud_pc_alert_rule" "test" {
  alert_rule_template = "cloudPcProvisionScenario"
  display_name        = "Test Alert Rule"
  severity            = "warning"
  enabled             = true
  is_system_rule      = false

  threshold = {
    aggregation = "count"
    operator    = "greaterOrEqual"
    target      = 1
  }

  conditions = [
    {
      relationship_type   = "and"
      condition_category  = "provisionFailures"
      aggregation        = "count"
      operator           = "greaterOrEqual"
      threshold_value    = "1"
    }
  ]
}
`
}

func testAccCloudPcAlertRuleConfig_missingThreshold() string {
	return `
resource "microsoft365_graph_beta_windows_365_cloud_pc_alert_rule" "test" {
  alert_rule_template = "cloudPcProvisionScenario"
  display_name        = "Test Alert Rule"
  severity            = "warning"
  enabled             = true
  is_system_rule      = false

  notification_channels = [
    {
      notification_channel_type = "portal"
      notification_receivers = [
        {
          contact_information = "admin@test.com"
          locale             = "en-US"
        }
      ]
    }
  ]

  conditions = [
    {
      relationship_type   = "and"
      condition_category  = "provisionFailures"
      aggregation        = "count"
      operator           = "greaterOrEqual"
      threshold_value    = "1"
    }
  ]
}
`
}

func testAccCloudPcAlertRuleConfig_invalidSeverity() string {
	return `
resource "microsoft365_graph_beta_windows_365_cloud_pc_alert_rule" "test" {
  alert_rule_template = "cloudPcProvisionScenario"
  display_name        = "Test Alert Rule"
  severity            = "invalid"
  enabled             = true
  is_system_rule      = false

  notification_channels = [
    {
      notification_channel_type = "portal"
      notification_receivers = [
        {
          contact_information = "admin@test.com"
          locale             = "en-US"
        }
      ]
    }
  ]

  threshold = {
    aggregation = "count"
    operator    = "greaterOrEqual"
    target      = 1
  }

  conditions = [
    {
      relationship_type   = "and"
      condition_category  = "provisionFailures"
      aggregation        = "count"
      operator           = "greaterOrEqual"
      threshold_value    = "1"
    }
  ]
}
`
}

func testAccCloudPcAlertRuleConfig_invalidAlertRuleTemplate() string {
	return `
resource "microsoft365_graph_beta_windows_365_cloud_pc_alert_rule" "test" {
  alert_rule_template = "invalidTemplate"
  display_name        = "Test Alert Rule"
  severity            = "warning"
  enabled             = true
  is_system_rule      = false

  notification_channels = [
    {
      notification_channel_type = "portal"
      notification_receivers = [
        {
          contact_information = "admin@test.com"
          locale             = "en-US"
        }
      ]
    }
  ]

  threshold = {
    aggregation = "count"
    operator    = "greaterOrEqual"
    target      = 1
  }

  conditions = [
    {
      relationship_type   = "and"
      condition_category  = "provisionFailures"
      aggregation        = "count"
      operator           = "greaterOrEqual"
      threshold_value    = "1"
    }
  ]
}
`
}
