# Acceptance test: Application Certificate Credential with HEX encoding
# Full dependency chain: random_string -> application -> time_sleep -> certificate_credential

resource "random_string" "test_id_hex" {
  length  = 8
  special = false
  upper   = false
}

# Minimal application for certificate testing
resource "microsoft365_graph_beta_applications_application" "test_app_hex" {
  display_name = "acc-test-app-cert-hex-${random_string.test_id_hex.result}"
  hard_delete  = true
}

# Wait for eventual consistency after application creation
resource "time_sleep" "wait_for_app_hex" {
  depends_on = [microsoft365_graph_beta_applications_application.test_app_hex]

  create_duration = "30s"
}

resource "microsoft365_graph_beta_applications_application_certificate_credential" "test_hex" {
  application_id = microsoft365_graph_beta_applications_application.test_app_hex.id
  display_name   = "acc-test-certificate-hex-${random_string.test_id_hex.result}"

  # Pre-generated hex-encoded self-signed certificate, valid until 2029-11-19. Keep the lifetime under the lab tenant's 3-year (P1095D) asymmetric key policy. Regenerate before then with:
  #   openssl req -x509 -newkey rsa:2048 -nodes -keyout key.pem -out cert.pem -days 1090 -subj "/C=US/O=Terraform Provider Test/CN=acc-test-certificate-hex"
  #   openssl x509 -in cert.pem -outform DER | xxd -p | tr -d '[:space:]'
  key      = "308203853082026da003020102021463a3727390e71445c945a07fe138b537ae9949ff300d06092a864886f70d01010b05003052310b30090603550406130255533120301e060355040a0c175465727261666f726d2050726f766964657220546573743121301f06035504030c186163632d746573742d63657274696669636174652d686578301e170d3236313132353037303632355a170d3239313131393037303632355a3052310b30090603550406130255533120301e060355040a0c175465727261666f726d2050726f766964657220546573743121301f06035504030c186163632d746573742d63657274696669636174652d68657830820122300d06092a864886f70d01010105000382010f003082010a0282010100c4eb01e78ae03a376a4645c202b1a693f2dff2796f5fd62b50aefa5cadcf168793798d05ee31aeebc13c8b9fb352c740670fe8cb07756527a36d84953a78ba750c459fa0168c7efd7e688eee54d094d911339428efe33a00f9f7550f075b4935c14c14b10816dce1d241fc7a6c25e117f390ca2eeddcaf91696ec5cc0751b8adef0a4d05ca06bb2e8b43e1854bb140fc26561cc052673c90a0dd979d4fcddb046a37d7486bc7e0b9f1d72742d10fda2b1b8e7e4482e2f772caa407c72c55d00edf4043af3e018c4af4f7a8e6bcb095c771651d29f929487084072d8e722031470e3c440ec3be3c6749a9d1dc7c246b5c948714c7f3640fc9e7db978cf199bca90203010001a3533051301d0603551d0e041604145a58ec6506065e0254d718583a84f992ecf74439301f0603551d230418301680145a58ec6506065e0254d718583a84f992ecf74439300f0603551d130101ff040530030101ff300d06092a864886f70d01010b0500038201010050451b845c717fb7f10fd669dd84ffdf251ac085b217f5bfd4b55d0a6faebf07378b63180dc6a7c8982d759aaa9e736676c1bac5bb291acaefc3e065f28b58d97bd7d24484afc2eb531731bc9baa45c243edc4656fcba3e5dad4812773cd0a565d981ef213c26ff5708ca0707ec66eb3f60c77c059e66bfc1218852a75c49a3c2c4d6c298b6fe83e0acc1b212cb6a554dd3d39b29a7f5f665a5aa95349b69164534107e9c6d3fe3f8af0a523fb95182c3c92eb0f9b7ab88293c32a4225e9dccf52eb8884822d99931d92d1eedb2c2c625786d3ed7a99fce2d20955e857853b746e78ff248b2d22caefa8dab80b47e9e5bbcfa7043a4a42c4647d6bc44107be37"
  encoding = "hex"
  type     = "AsymmetricX509Cert"
  usage    = "Verify"

  depends_on = [time_sleep.wait_for_app_hex]
}
