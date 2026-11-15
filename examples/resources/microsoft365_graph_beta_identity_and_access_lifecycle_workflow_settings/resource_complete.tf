# Run scheduled lifecycle workflows every 12 hours and send email notifications
# from a verified tenant domain using the company branding.
resource "microsoft365_graph_beta_identity_and_access_lifecycle_workflow_settings" "example" {
  workflow_schedule_interval_in_hours = 12

  email_settings = {
    # Must be a verified domain of the tenant, written exactly as in the tenant domain list.
    # Use "microsoft.com" to send from the default Microsoft sender.
    sender_domain        = "contoso.com"
    use_company_branding = true
  }

  timeouts = {
    create = "5m"
    read   = "5m"
    update = "5m"
    delete = "5m"
  }
}
