resource "microsoft365_graph_beta_device_management_linux_device_compliance_script" "example" {
  display_name = "Linux Security Compliance Check"
  description  = "Checks for critical security configurations on Linux devices"

  detection_script_content = <<-EOT
    #!/bin/bash
    # Detection script for Linux security compliance
    
    # Check if firewall is active
    if ! systemctl is-active --quiet ufw && ! systemctl is-active --quiet firewalld; then
      echo "No firewall service is active - Non-compliant"
      exit 1
    fi
    
    # Check for automatic updates
    if command -v unattended-upgrade >/dev/null 2>&1; then
      if [ ! -f /etc/apt/apt.conf.d/20auto-upgrades ]; then
        echo "Automatic updates not configured - Non-compliant"
        exit 1
      fi
    elif command -v dnf >/dev/null 2>&1; then
      if ! systemctl is-enabled --quiet dnf-automatic.timer; then
        echo "Automatic updates not enabled - Non-compliant"
        exit 1
      fi
    fi
    
    # Check SSH configuration
    if [ -f /etc/ssh/sshd_config ]; then
      if grep -q "^PermitRootLogin yes" /etc/ssh/sshd_config; then
        echo "Root SSH login is enabled - Non-compliant"
        exit 1
      fi
    fi
    
    # Check for required security packages
    if command -v apt >/dev/null 2>&1; then
      if ! dpkg -l | grep -q fail2ban; then
        echo "fail2ban not installed - Non-compliant"
        exit 1
      fi
    elif command -v rpm >/dev/null 2>&1; then
      if ! rpm -q fail2ban >/dev/null 2>&1; then
        echo "fail2ban not installed - Non-compliant"
        exit 1
      fi
    fi
    
    echo "All security checks passed - Compliant"
    exit 0
  EOT

  timeouts = {
    create = "30m"
    read   = "5m"
    update = "30m"
    delete = "30m"
  }
}