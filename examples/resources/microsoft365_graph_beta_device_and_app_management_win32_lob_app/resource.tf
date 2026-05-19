resource "microsoft365_graph_beta_device_and_app_management_win32_lob_app" "seven_zip" {
  display_name    = "7-Zip 24.09 x64"
  description     = "7-Zip is a file archiver with a high compression ratio. Supports 7z, ZIP, CAB, RAR, ARJ, LZH, CHM, CPIO, CramFS, DEB, DMG, FAT, HFS, ISO, LZH, LZMA, MBR, MSI, NSIS, NTFS, RAR, RPM, SquashFS, UDF, VHD, WIM, XAR and Z formats."
  publisher       = "Igor Pavlov"
  product_code    = "{23170F69-40C1-2702-2409-000001000000}"
  product_version = "24.09.0.0"
  file_name       = "7z2409-x64.msi.intunewin"

  categories = [
    "Business",
    "Productivity",
  ]

  information_url = "https://www.7-zip.org/"
  developer       = "Igor Pavlov"
  owner           = "IT"
  notes           = "Standard MSI deployment. Silent install handled by Intune LOB app engine."

  app_installer = {
    installer_file_path_source = "/path/to/7z2409-x64.msi.intunewin"
  }

  app_icon = {
    icon_url_source = "https://upload.wikimedia.org/wikipedia/commons/thumb/2/21/7-zip-logo.svg/480px-7-zip-logo.svg.png"
  }

  command_line = "/qn REBOOT=ReallySuppress"
}
