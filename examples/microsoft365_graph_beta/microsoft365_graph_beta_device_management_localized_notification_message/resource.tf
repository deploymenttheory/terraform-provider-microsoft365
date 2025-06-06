# Example device enrollment notification configuration
resource "microsoft365_graph_beta_device_management_device_enrollment_notification_configuration" "example" {
  display_name   = "Example Notification Configuration"
  description    = "Example configuration for notifications"
  template_types = ["email", "push"]
  priority       = 1

  timeouts = {
    create = "180s"
    read   = "180s"
    update = "180s"
    delete = "180s"
  }
}

# Basic English localized notification message for email
resource "microsoft365_graph_beta_device_management_localized_notification_message" "english_email_message" {
  device_enrollment_notification_configuration_id = microsoft365_graph_beta_device_management_device_enrollment_notification_configuration.example.id
  template_type                                   = "email"
  locale                                         = "en-us"
  subject                                        = "Device Enrollment Complete"
  message_template                               = "Your device has been successfully enrolled in our mobile device management system. Please contact IT support if you have any questions."
  is_default                                     = true

  timeouts = {
    create = "180s"
    read   = "180s"
    update = "180s"
    delete = "180s"
  }
}

# Spanish localized notification message for email
resource "microsoft365_graph_beta_device_management_localized_notification_message" "spanish_email_message" {
  device_enrollment_notification_configuration_id = microsoft365_graph_beta_device_management_device_enrollment_notification_configuration.example.id
  template_type                                   = "email"
  locale                                         = "es-es"
  subject                                        = "Inscripción de dispositivo completada"
  message_template                               = "Su dispositivo ha sido inscrito exitosamente en nuestro sistema de gestión de dispositivos móviles. Por favor contacte a soporte técnico si tiene alguna pregunta."
  is_default                                     = false

  timeouts = {
    create = "180s"
    read   = "180s"
    update = "180s"
    delete = "180s"
  }
}

# Push notification message (shorter content for mobile notifications)
resource "microsoft365_graph_beta_device_management_localized_notification_message" "push_message" {
  device_enrollment_notification_configuration_id = microsoft365_graph_beta_device_management_device_enrollment_notification_configuration.example.id
  template_type                                   = "push"
  locale                                         = "en-us"
  subject                                        = "Enrollment Success"
  message_template                               = "Device enrolled successfully!"
  is_default                                     = true

  timeouts = {
    create = "180s"
    read   = "180s"
    update = "180s"
    delete = "180s"
  }
} 