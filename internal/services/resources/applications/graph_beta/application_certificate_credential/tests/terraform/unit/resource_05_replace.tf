# Unit test: Application Certificate Credential - Replace Existing Certificates
# Tests that replace_existing_certificates=true can be set and the resource is created

resource "microsoft365_graph_beta_applications_application_certificate_credential" "test_replace" {
  application_id = "55555555-5555-5555-5555-555555555555"
  display_name   = "unit-test-certificate-replace"

  # Sample PEM certificate
  key      = <<EOT
-----BEGIN CERTIFICATE-----
MIIDvzCCAqegAwIBAgIUIMw1dI1Z8ZIX+G1MiRL4BbR4lQYwDQYJKoZIhvcNAQEL
BQAwbzELMAkGA1UEBhMCVVMxDTALBgNVBAgMBFRlc3QxDTALBgNVBAcMBFRlc3Qx
GjAYBgNVBAoMEVRlc3QgT3JnYW5pemF0aW9uMSYwJAYDVQQDDB1BZ2VudCBJZGVu
dGl0eSBCbHVlcHJpbnQgVGVzdDAeFw0yNTEyMDUxMTU5MDVaFw0yNjEyMDUxMTU5
MDVaMG8xCzAJBgNVBAYTAlVTMQ0wCwYDVQQIDARUZXN0MQ0wCwYDVQQHDARUZXN0
MRowGAYDVQQKDBFUZXN0IE9yZ2FuaXphdGlvbjEmMCQGA1UEAwwdQWdlbnQgSWRl
bnRpdHkgQmx1ZXByaW50IFRlc3QwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEK
AoIBAQC2tHd5hTDWFYIo9BWC6H0aXoS5OvOHlTxEEPxDbL2xdTGXBCl3x0FAieW9
tlD+jzpZzX8XD7a/4VynB0vCb+5Bi4txa3xU+ObtCC7Z7Q+SgczzCY2jEMQ8DaKu
FD+pyzscXkuXHO8cGNokjd3ULChVUi7MucF60DJaQCGeXWai/GvJ5BC1Ywn3lXjj
yDuVJqtf+P0x4/IWMqnH5uOj68pGKtnx+k8Ome8qgwRrhFKtSVM5TyZzkfSGZSgt
4EKQEl5/2IuvmbCMEw4m+5o6RzrlFN7F1KuwUEpbx5X8X/KueRfLGEjeuGksz8ZZ
2rA2wi2PU032b6RjHyi5DS29AlihAgMBAAGjUzBRMB0GA1UdDgQWBBSTU3e9u/ju
oMwqMjgQt/iO9M+CmzAfBgNVHSMEGDAWgBSTU3e9u/juoMwqMjgQt/iO9M+CmzAP
BgNVHRMBAf8EBTADAQH/MA0GCSqGSIb3DQEBCwUAA4IBAQBBSX+hqIOZmjU8SXY0
iPwAYFZeonDwn8kUHUBem0IdBlNc4dewIbdiUmLE6+765bmjHBfZNavNF9vdF2Hb
rDquUorIQvJQheNK6dTWTS+ehuIHFICUdznwUB3HmBCIPU69vtA8TSOq3ZrTZL4Z
h4D51mPm/ePE7IZ/MJDAR8ZB0/5FRyu+RdXDG7i1TuDawlz2uYk10Iv3G0fxrfTm
VuTGU88nFUDtXw8A0lHJOMdF5xx0gyou8IXM/EVfHHESNcyLJRkYD1wdHy2+uq1m
+HlBmHypstf57Bnd8hghN/VW3NcTuyAzISmq11hx4r8SxmQ/Fuq3QvfqNQj4o/6T
/+gw
-----END CERTIFICATE-----
EOT
  encoding = "pem"
  type     = "AsymmetricX509Cert"
  usage    = "Verify"

  # Test that replace_existing_certificates=true can be set
  # When true, this will remove ALL existing certificates and leave only this one
  replace_existing_certificates = true
}
