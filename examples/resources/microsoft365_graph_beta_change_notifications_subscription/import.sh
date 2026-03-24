#!/usr/bin/env bash
# Replace {subscription_id} with the Microsoft Graph subscription id (GUID).
terraform import microsoft365_graph_beta_change_notifications_subscription.example '{subscription_id}'
