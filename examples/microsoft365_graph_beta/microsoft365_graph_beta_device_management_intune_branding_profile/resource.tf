resource "microsoft365_graph_beta_device_management_intune_branding_profile" "example" {
  profile_name                               = "Corporate Branding"
  profile_description                        = "Corporate branding profile"
  display_name                               = "Company Portal Branding"
  contact_it_name                            = "IT Support Team"
  contact_it_phone_number                    = "+1 (555) 123-4567"
  contact_it_email_address                   = "support@example.com"
  contact_it_notes                           = "Available Monday to Friday, 9am to 5pm"
  online_support_site_url                    = "https://support.example.com"
  online_support_site_name                   = "Company IT Support Portal"
  privacy_url                                = "https://www.example.com/privacy"
  custom_privacy_message                     = "Your privacy is important to us. Please read our privacy policy."
  custom_can_see_privacy_message             = "You can view our privacy policy at any time."
  custom_cant_see_privacy_message            = "Please contact IT support for privacy information."
  show_logo                                  = true
  show_display_name_next_to_logo             = true
  role_scope_tag_ids = ["0"]
  
  # Optional: Branding colors
  theme_color = {
    r = 0    # Red
    g = 120  # Green
    b = 212  # Blue
  }
  
  # Optional: Branding images
  theme_color_logo = {
    image_url_source = "https://mailmeteor.com/logos/assets/PNG/Microsoft_Logo_256px.png"
  }
  
  light_background_logo = {
    image_url_source = "https://mailmeteor.com/logos/assets/PNG/Microsoft_Logo_256px.png"
  }
  
  landing_page_customized_image = {
    image_url_source = "https://mailmeteor.com/logos/assets/PNG/Microsoft_Logo_256px.png"
  }
  
  # Company portal settings
  is_remove_device_disabled                  = false
  is_factory_reset_disabled                  = false
  show_azure_ad_enterprise_apps              = true
  show_office_web_apps                       = true
  send_device_ownership_change_push_notification = true
  enrollment_availability                    = "availableWithPrompts"
  disable_client_telemetry                   = false
  is_default_profile                         = false

  # Company portal blocked actions
  company_portal_blocked_actions = [
    {
      platform  = "windows10AndLater"
      owner_type = "company"
      action    = "remove"
    },
    {
      platform  = "iOS"
      owner_type = "company"
      action    = "reset"
    }
  ]
  
  # Assignments
  assignments = [
    # Optional: Assignment targeting a specific group
    {
      type        = "groupAssignmentTarget"
      group_id    = "51a96cdd-4b9b-4849-b416-8c94a6d88797"
    },
    {
      type        = "groupAssignmentTarget"
      group_id    = "b15228f4-9d49-41ed-9b4f-0e7c721fd9c2"
    },
    # Optional: Exclusion group assignments
    {
      type     = "exclusionGroupAssignmentTarget"
      group_id = "b8c661c2-fa9a-4351-af86-adc1729c343f"
    },
    {
      type     = "exclusionGroupAssignmentTarget"
      group_id = "f6ebd6ff-501e-4b3d-a00b-a2e102c3fa0f"
    },

  ]
  
  # Timeouts
  timeouts ={
    create = "30m"
    read   = "10m"
    update = "30m"
    delete = "10m"
  }
} 