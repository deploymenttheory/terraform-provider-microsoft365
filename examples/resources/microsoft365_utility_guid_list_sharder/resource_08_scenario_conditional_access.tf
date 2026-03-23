# ==============================================================================
# Scenario: Sharder → Security Groups → Conditional Access Policy (Phased)
#
# This example demonstrates how to wire sharder outputs directly into a phased
# Conditional Access rollout. The pattern is:
#
#   sharder  → shard GUIDs (users)
#   groups   → one security group per ring (stable IDs, group-based targeting)
#   members  → populate each group from the corresponding shard
#   CA policy per ring → target the ring's group; start in report-only mode
#
# Why one CA policy per ring rather than one policy with all groups?
#   Each ring can be independently promoted from report-only → enabled without
#   touching the other rings. Ring 0 can be fully enforced while Ring 1 is still
#   in observation mode. This also makes rollback surgical — disable Ring 1's
#   policy without any impact on Ring 0 or Ring 2.
#
# Use case: Phased Require MFA rollout — 10% pilot, then 30%, then 60%
# ==============================================================================

resource "microsoft365_utility_guid_list_sharder" "mfa_ca_rings" {
  resource_type           = "users"
  odata_filter            = "accountEnabled eq true and userType eq 'Member'"
  shard_percentages       = [10, 30, 60]
  strategy                = "percentage"
  seed                    = "ca-mfa-2026" # Unique seed — different population than other rollouts
  recalculate_on_next_run = false
}

# ─────────────────────────────────────────────────────────────────────────────
# Security Groups (one per ring)
# ─────────────────────────────────────────────────────────────────────────────

resource "microsoft365_graph_beta_groups_group" "ca_mfa_ring_0" {
  display_name     = "CA - Require MFA - Ring 0 (10% Pilot)"
  mail_nickname    = "ca-mfa-ring-0"
  description      = "Pilot cohort for Require MFA CA policy — 10% of members"
  security_enabled = true
  mail_enabled     = false
  hard_delete      = true
}

resource "microsoft365_graph_beta_groups_group" "ca_mfa_ring_1" {
  display_name     = "CA - Require MFA - Ring 1 (30% Broader)"
  mail_nickname    = "ca-mfa-ring-1"
  description      = "Second wave for Require MFA CA policy — 30% of members"
  security_enabled = true
  mail_enabled     = false
  hard_delete      = true
}

resource "microsoft365_graph_beta_groups_group" "ca_mfa_ring_2" {
  display_name     = "CA - Require MFA - Ring 2 (60% Full)"
  mail_nickname    = "ca-mfa-ring-2"
  description      = "Final wave for Require MFA CA policy — remaining 60% of members"
  security_enabled = true
  mail_enabled     = false
  hard_delete      = true
}

# Break-glass account group — always excluded from every CA policy to ensure
# emergency access is never blocked by a misconfigured policy.
# This group must be created separately and pre-populated with break-glass accounts.
data "microsoft365_graph_beta_groups_group" "breakglass" {
  display_name = "Break Glass - Emergency Access"
}

# ─────────────────────────────────────────────────────────────────────────────
# Group Member Assignments — populate each ring group from sharder output
# ─────────────────────────────────────────────────────────────────────────────

resource "microsoft365_graph_beta_groups_group_member_assignment" "ca_mfa_ring_0_members" {
  for_each = microsoft365_utility_guid_list_sharder.mfa_ca_rings.shards["shard_0"]

  group_id           = microsoft365_graph_beta_groups_group.ca_mfa_ring_0.id
  member_id          = each.value
  member_object_type = "User"
}

resource "microsoft365_graph_beta_groups_group_member_assignment" "ca_mfa_ring_1_members" {
  for_each = microsoft365_utility_guid_list_sharder.mfa_ca_rings.shards["shard_1"]

  group_id           = microsoft365_graph_beta_groups_group.ca_mfa_ring_1.id
  member_id          = each.value
  member_object_type = "User"
}

resource "microsoft365_graph_beta_groups_group_member_assignment" "ca_mfa_ring_2_members" {
  for_each = microsoft365_utility_guid_list_sharder.mfa_ca_rings.shards["shard_2"]

  group_id           = microsoft365_graph_beta_groups_group.ca_mfa_ring_2.id
  member_id          = each.value
  member_object_type = "User"
}

# ─────────────────────────────────────────────────────────────────────────────
# Conditional Access Policies — one per ring, independently promotable
#
# Lifecycle:
#   1. Deploy all three policies in "enabledForReportingButNotEnforced"
#   2. Review Ring 0 sign-in logs for 1–2 weeks; no unexpected blocks → proceed
#   3. Change Ring 0 state to "enabled"; leave Ring 1 in report-only
#   4. Repeat for Ring 1 (1–2 weeks observation), then Ring 2
#   5. Once Ring 2 is enabled, all members are covered; rings can be consolidated
# ─────────────────────────────────────────────────────────────────────────────

# Ring 0 — Pilot (10%): report-only until validated
resource "microsoft365_graph_beta_identity_and_access_conditional_access_policy" "require_mfa_ring_0" {
  display_name = "CA - Require MFA - Ring 0 (10% Pilot)"
  state        = "enabledForReportingButNotEnforced" # Promote to "enabled" after validation

  conditions = {
    client_app_types = ["browser", "mobileAppsAndDesktopClients"]

    users = {
      include_users  = []
      exclude_users  = []
      include_groups = [microsoft365_graph_beta_groups_group.ca_mfa_ring_0.id]
      exclude_groups = [data.microsoft365_graph_beta_groups_group.breakglass.id]
      include_roles  = []
      exclude_roles  = []
    }

    applications = {
      include_applications                            = ["All"]
      exclude_applications                            = []
      include_user_actions                            = []
      include_authentication_context_class_references = []
    }

    sign_in_risk_levels = []
  }

  grant_controls = {
    operator                      = "OR"
    built_in_controls             = ["mfa"]
    custom_authentication_factors = []
  }
}

# Ring 1 — Broader (30%): disabled until Ring 0 is fully enforced
resource "microsoft365_graph_beta_identity_and_access_conditional_access_policy" "require_mfa_ring_1" {
  display_name = "CA - Require MFA - Ring 1 (30% Broader)"
  state        = "disabled" # Enable after Ring 0 reaches "enabled" and is stable

  conditions = {
    client_app_types = ["browser", "mobileAppsAndDesktopClients"]

    users = {
      include_users  = []
      exclude_users  = []
      include_groups = [microsoft365_graph_beta_groups_group.ca_mfa_ring_1.id]
      exclude_groups = [data.microsoft365_graph_beta_groups_group.breakglass.id]
      include_roles  = []
      exclude_roles  = []
    }

    applications = {
      include_applications                            = ["All"]
      exclude_applications                            = []
      include_user_actions                            = []
      include_authentication_context_class_references = []
    }

    sign_in_risk_levels = []
  }

  grant_controls = {
    operator                      = "OR"
    built_in_controls             = ["mfa"]
    custom_authentication_factors = []
  }
}

# Ring 2 — Full (60%): disabled until Ring 1 is fully enforced
resource "microsoft365_graph_beta_identity_and_access_conditional_access_policy" "require_mfa_ring_2" {
  display_name = "CA - Require MFA - Ring 2 (60% Full)"
  state        = "disabled" # Enable after Ring 1 reaches "enabled" and is stable

  conditions = {
    client_app_types = ["browser", "mobileAppsAndDesktopClients"]

    users = {
      include_users  = []
      exclude_users  = []
      include_groups = [microsoft365_graph_beta_groups_group.ca_mfa_ring_2.id]
      exclude_groups = [data.microsoft365_graph_beta_groups_group.breakglass.id]
      include_roles  = []
      exclude_roles  = []
    }

    applications = {
      include_applications                            = ["All"]
      exclude_applications                            = []
      include_user_actions                            = []
      include_authentication_context_class_references = []
    }

    sign_in_risk_levels = []
  }

  grant_controls = {
    operator                      = "OR"
    built_in_controls             = ["mfa"]
    custom_authentication_factors = []
  }
}

# ─────────────────────────────────────────────────────────────────────────────
# Outputs
# ─────────────────────────────────────────────────────────────────────────────

output "ca_mfa_rollout_summary" {
  description = "CA policy IDs and ring population counts — share with security team before promoting states"
  value = {
    ring_0 = {
      policy_id = microsoft365_graph_beta_identity_and_access_conditional_access_policy.require_mfa_ring_0.id
      group_id  = microsoft365_graph_beta_groups_group.ca_mfa_ring_0.id
      headcount = length(microsoft365_utility_guid_list_sharder.mfa_ca_rings.shards["shard_0"])
      state     = "enabledForReportingButNotEnforced"
    }
    ring_1 = {
      policy_id = microsoft365_graph_beta_identity_and_access_conditional_access_policy.require_mfa_ring_1.id
      group_id  = microsoft365_graph_beta_groups_group.ca_mfa_ring_1.id
      headcount = length(microsoft365_utility_guid_list_sharder.mfa_ca_rings.shards["shard_1"])
      state     = "disabled"
    }
    ring_2 = {
      policy_id = microsoft365_graph_beta_identity_and_access_conditional_access_policy.require_mfa_ring_2.id
      group_id  = microsoft365_graph_beta_groups_group.ca_mfa_ring_2.id
      headcount = length(microsoft365_utility_guid_list_sharder.mfa_ca_rings.shards["shard_2"])
      state     = "disabled"
    }
  }
}
