| Action Name | Unit Test Harness | Managed Device API | Comanaged Device API | Send Messages | Validate Device Type |
|---|---|---|---|---|---|
| activate_device_esim | ✅ | POST /managedDevices/{id}/activateDeviceEsim | POST /deviceManagement/comanagedDevices/{id}/activateDeviceEsim | ❌ | ✅ iPhone, iPad only |
| bypass_activation_lock | ✅ | POST /managedDevices/{id}/bypassActivationLock | ❌ | ❌ | ✅ iPhone, iPad, Mac (enum validated) |
| clean_windows_device | ✅ | POST /managedDevices/{id}/cleanWindowsDevice | POST /deviceManagement/comanagedDevices/{id}/cleanWindowsDevice | ❌ | ✅ Windows only |
| create_device_log_collection_request | ❌ | POST /managedDevices/{id}/createDeviceLogCollectionRequest | POST /deviceManagement/comanagedDevices/{id}/createDeviceLogCollectionRequest | ❌ | ✅ Windows only (ADR-001 aligned) |
| delete_user_from_shared_apple_device | ❌ | POST /managedDevices/{id}/deleteUserFromSharedAppleDevice | POST /deviceManagement/comanagedDevices/{id}/deleteUserFromSharedAppleDevice | ❌ | ✅ iPad only (ADR-001 aligned) |
| deprovision | ❌ | POST /managedDevices/{id}/deprovision | ❌ | ❌ | ❓ |
| disable | ❌ | POST /managedDevices/{id}/disable | ❌ | ❌ | ❓ |
| disable_lost_mode | ❌ | POST /managedDevices/{id}/disableLostMode | ❌ | ❌ | ❓ |
| enable_lost_mode | ❌ | POST /managedDevices/{id}/enableLostMode | ❌ | ❌ | ❓ |
| get_file_vault_key | ❌ | GET /managedDevices/{id}/getFileVaultKey | ❌ | ❌ | ❓ |
| initiate_device_attestation | ❌ | POST /managedDevices/{id}/initiateDeviceAttestation | ❌ | ❌ | ❓ |
| initiate_mobile_device_management_key_recovery | ❌ | POST /managedDevices/{id}/initiateMobileDeviceManagementKeyRecovery | ❌ | ❌ | ❓ |
| initiate_on_demand_proactive_remediation | ❌ | POST /managedDevices/{id}/initiateOnDemandProactiveRemediation | ❌ | ❌ | ❓ |
| locate_device | ❌ | POST /managedDevices/{id}/locateDevice | ❌ | ❌ | ❓ |
| logout_shared_apple_device_active_user | ❌ | POST /managedDevices/{id}/logoutSharedAppleDeviceActiveUser | ❌ | ❌ | ❓ |
| move_devices_to_ou | ❌ | POST /managedDevices/{id}/moveDevicesToOU | ❌ | ❌ | ❓ |
| pause_configuration_refresh | ❌ | POST /managedDevices/{id}/pauseConfigurationRefresh | ❌ | ❌ | ❓ |
| play_lost_mode_sound | ❌ | POST /managedDevices/{id}/playLostModeSound | ❌ | ❌ | ❓ |
| reboot_now | ❌ | POST /managedDevices/{id}/rebootNow | ❌ | ❌ | ❓ |
| recover_passcode | ❌ | POST /managedDevices/{id}/recoverPasscode | ❌ | ❌ | ❓ |
| reenable | ❌ | POST /managedDevices/{id}/reenable | ❌ | ❌ | ❓ |
| remote_lock | ❌ | POST /managedDevices/{id}/remoteLock | ❌ | ❌ | ❓ |
| remove_device_firmware_configuration_interface_management | ❌ | POST /managedDevices/{id}/removeDeviceFirmwareConfigurationInterfaceManagement | ❌ | ❌ | ❓ |
| reset_passcode | ❌ | POST /managedDevices/{id}/resetPasscode | ❌ | ❌ | ❓ |
| retire | ❌ | POST /managedDevices/{id}/retire | ❌ | ❌ | ❓ |
| revoke_apple_vpp_licenses | ❌ | POST /managedDevices/{id}/revokeAppleVppLicenses | ❌ | ❌ | ❓ |
| rotate_bitlocker_keys | ❌ | POST /managedDevices/{id}/rotateBitlockerKeys | ❌ | ❌ | ❓ |
| rotate_file_vault_key | ❌ | POST /managedDevices/{id}/rotateFileVaultKey | ❌ | ❌ | ❓ |
| rotate_local_admin_password | ❌ | POST /managedDevices/{id}/rotateLocalAdminPassword | ❌ | ❌ | ❓ |
| send_custom_notification_to_company_portal | ❌ | POST /managedDevices/{id}/sendCustomNotificationToCompanyPortal | ❌ | ✅ | ❓ |
| set_device_name | ❌ | POST /managedDevices/{id}/setDeviceName | ❌ | ❌ | ❓ |
| shutdown | ❌ | POST /managedDevices/{id}/shutdown | ❌ | ❌ | ❓ |
| sync_device | ❌ | POST /managedDevices/{id}/syncDevice | ❌ | ❌ | ❓ |
| trigger_configuration_manager_action | ❌ | POST /managedDevices/{id}/triggerConfigurationManagerAction | ❌ | ❌ | ❓ |
| update_windows_device_account | ❌ | POST /managedDevices/{id}/updateWindowsDeviceAccount | ❌ | ❌ | ❓ |
| windows_defender_scan | ❌ | POST /managedDevices/{id}/windowsDefenderScan | ❌ | ❌ | ❓ |
| windows_defender_update_signatures | ❌ | POST /managedDevices/{id}/windowsDefenderUpdateSignatures | ❌ | ❌ | ❓ |
| wipe | ❌ | POST /managedDevices/{id}/wipe | ❌ | ❌ | ❓ |
