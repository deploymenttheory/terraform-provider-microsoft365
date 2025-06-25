resource "microsoft365_graph_beta_groups_group" "maximal" {
  display_name                   = "Maximal Group"
  description                    = "This is a maximal group configuration for testing"
  mail_nickname                  = "maximal.group"
  mail_enabled                   = true
  security_enabled               = true
  group_types                    = ["Unified", "DynamicMembership"]
  visibility                     = "Private"
  is_assignable_to_role          = false
  membership_rule                = "user.department -eq \"Engineering\""
  membership_rule_processing_state = "On"
  preferred_data_location        = "NAM"
  preferred_language             = "en-US"
  theme                          = "Blue"
  classification                 = "High"
  
} 